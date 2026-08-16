package service

import (
	"log/slog"
	"time"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/util"

	"gorm.io/gorm"
)

// StockRecordService 库存记录业务逻辑。
type StockRecordService struct {
	db            *gorm.DB
	repo          *repository.StockRecordRepository
	productRepo   *repository.ProductRepository
	warehouseRepo *repository.WarehouseRepository
	logger        *slog.Logger
}

// NewStockRecordService 构造库存记录服务。
func NewStockRecordService(db *gorm.DB, repo *repository.StockRecordRepository, productRepo *repository.ProductRepository,
	warehouseRepo *repository.WarehouseRepository, logger *slog.Logger) *StockRecordService {
	return &StockRecordService{db: db, repo: repo, productRepo: productRepo, warehouseRepo: warehouseRepo, logger: logger}
}

// Query 组合条件查询库存。
func (s *StockRecordService) Query(productID, warehouseID, shelfID uint64, page, pageSize int) ([]model.StockRecord, int64, error) {
	list, total, err := s.repo.Query(productID, warehouseID, shelfID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	s.logger.Info(constants.LogStockQuery, "total", total)
	return list, total, nil
}

// Inbound 入库：在目标库位增加库存记录。
func (s *StockRecordService) Inbound(productID, warehouseID, shelfID uint64, quantity int, batchNo string) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] inbound: quantity must > 0")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.InboundTx(tx, productID, warehouseID, shelfID, quantity, batchNo)
	})
}

// InboundTx 在给定事务内执行入库。
func (s *StockRecordService) InboundTx(tx *gorm.DB, productID, warehouseID, shelfID uint64, quantity int, batchNo string) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] inbound: quantity must > 0")
	}
	now := time.Now()
	rec := &model.StockRecord{
		ProductID: productID, WarehouseID: warehouseID, ShelfID: shelfID,
		BatchNo: batchNo, Quantity: quantity,
		LastOpType: constants.StockOpInbound, LastOpAt: now, InboundDate: &now,
	}
	if err := s.repo.CreateTx(tx, rec); err != nil {
		return util.Wrap(err, "StockRecord[product_id=%d] inbound create failed", productID)
	}
	return nil
}

// Outbound 出库：从指定库位扣减库存，不足时报错。
func (s *StockRecordService) Outbound(productID, warehouseID, shelfID uint64, quantity int) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] outbound: quantity must > 0")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.OutboundTx(tx, productID, warehouseID, shelfID, quantity)
	})
}

// OutboundTx 在给定事务内执行出库。
func (s *StockRecordService) OutboundTx(tx *gorm.DB, productID, warehouseID, shelfID uint64, quantity int) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] outbound: quantity must > 0")
	}
	records, err := s.repo.FindByProductWarehouseForUpdate(tx, productID, warehouseID)
	if err != nil {
		return err
	}
	var available int
	for _, r := range records {
		available += r.Quantity
	}
	if available < quantity {
		s.logger.Warn(constants.LogStockOutboundFailed, "product_id", productID, "available", available, "need", quantity)
		return util.NewAppError(constants.CodeInsufficientStock, constants.MsgInsufficientStock)
	}
	remain := quantity
	for i := range records {
		if remain <= 0 {
			break
		}
		take := records[i].Quantity
		if take > remain {
			take = remain
		}
		records[i].Quantity -= take
		records[i].LastOpType = constants.StockOpOutbound
		records[i].LastOpAt = time.Now()
		remain -= take
		if err := s.repo.UpdateTx(tx, &records[i]); err != nil {
			return util.Wrap(err, "StockRecord[id=%d] outbound update failed", records[i].ID)
		}
	}
	return nil
}

// Adjust 盘点校准：设置实际数量（差异自动修正）。
func (s *StockRecordService) Adjust(productID, warehouseID, shelfID uint64, actualQuantity int) error {
	if actualQuantity < 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[actual="+itoa(uint64(actualQuantity))+"] adjust: actual must >= 0")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.AdjustTx(tx, productID, warehouseID, shelfID, actualQuantity)
	})
}

// AdjustTx 在给定事务内执行盘点校准。
func (s *StockRecordService) AdjustTx(tx *gorm.DB, productID, warehouseID, shelfID uint64, actualQuantity int) error {
	if actualQuantity < 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[actual="+itoa(uint64(actualQuantity))+"] adjust: actual must >= 0")
	}
	records, err := s.repo.FindByProductWarehouseForUpdate(tx, productID, warehouseID)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		now := time.Now()
		rec := &model.StockRecord{ProductID: productID, WarehouseID: warehouseID, ShelfID: shelfID,
			Quantity: actualQuantity, LastOpType: constants.StockOpAdjust, LastOpAt: now, InboundDate: &now}
		if err := s.repo.CreateTx(tx, rec); err != nil {
			return util.Wrap(err, "StockRecord[product_id=%d] adjust create failed", productID)
		}
		return nil
	}
	total := 0
	for _, r := range records {
		total += r.Quantity
	}
	diff := actualQuantity - total
	records[0].Quantity += diff
	records[0].LastOpType = constants.StockOpAdjust
	records[0].LastOpAt = time.Now()
	if err := s.repo.UpdateTx(tx, &records[0]); err != nil {
		return util.Wrap(err, "StockRecord[id=%d] adjust update failed", records[0].ID)
	}
	return nil
}

// Transfer 调拨：源仓库出库 + 目标仓库入库。
func (s *StockRecordService) Transfer(productID, fromWarehouse, toWarehouse, shelfID uint64, quantity int) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] transfer: quantity must > 0")
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.TransferTx(tx, productID, fromWarehouse, toWarehouse, shelfID, quantity)
	})
}

// TransferTx 在给定事务内执行调拨。
func (s *StockRecordService) TransferTx(tx *gorm.DB, productID, fromWarehouse, toWarehouse, shelfID uint64, quantity int) error {
	if quantity <= 0 {
		return util.NewAppError(constants.CodeValidationFailed, "StockRecord[quantity="+itoa(uint64(quantity))+"] transfer: quantity must > 0")
	}
	if err := s.OutboundTx(tx, productID, fromWarehouse, 0, quantity); err != nil {
		return err
	}
	if err := s.InboundTx(tx, productID, toWarehouse, shelfID, quantity, "TRANSFER"); err != nil {
		return err
	}
	return nil
}
