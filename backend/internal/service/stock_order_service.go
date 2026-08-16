package service

import (
	"fmt"
	"log/slog"
	"time"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/util"

	"gorm.io/gorm"
)

// StockOrderService 出入库单业务逻辑。
type StockOrderService struct {
	db            *gorm.DB
	repo          *repository.StockOrderRepository
	itemRepo      *repository.StockOrderItemRepository
	stockSvc      *StockRecordService
	productRepo   *repository.ProductRepository
	warehouseRepo *repository.WarehouseRepository
	logger        *slog.Logger
}

// NewStockOrderService 构造出入库单服务。
func NewStockOrderService(db *gorm.DB, repo *repository.StockOrderRepository, itemRepo *repository.StockOrderItemRepository,
	stockSvc *StockRecordService, productRepo *repository.ProductRepository,
	warehouseRepo *repository.WarehouseRepository, logger *slog.Logger) *StockOrderService {
	return &StockOrderService{db: db, repo: repo, itemRepo: itemRepo, stockSvc: stockSvc,
		productRepo: productRepo, warehouseRepo: warehouseRepo, logger: logger}
}

// Create 创建单据（草稿）。
func (s *StockOrderService) Create(creatorID uint64, orderType string, sourceWarehouseID, targetWarehouseID uint64,
	remark string, items []model.StockOrderItem) (*model.StockOrder, error) {
	if !constants.IsValidOrderType(orderType) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "StockOrder[order_type="+orderType+"] create: invalid type")
	}
	for _, it := range items {
		if _, err := s.productRepo.FindByID(it.ProductID); err != nil {
			return nil, util.Wrap(err, "StockOrder[product_id=%d] create: product not found", it.ProductID)
		}
	}
	if orderType == constants.OrderTypeTransfer && (sourceWarehouseID == 0 || targetWarehouseID == 0) {
		return nil, util.NewAppError(constants.CodeValidationFailed, "StockOrder[transfer] create: need source and target warehouse")
	}
	if orderType == constants.OrderTypeInbound && targetWarehouseID == 0 {
		return nil, util.NewAppError(constants.CodeValidationFailed, "StockOrder[inbound] create: need target warehouse")
	}
	if orderType == constants.OrderTypeOutbound && sourceWarehouseID == 0 {
		return nil, util.NewAppError(constants.CodeValidationFailed, "StockOrder[outbound] create: need source warehouse")
	}
	o := &model.StockOrder{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		o.OrderNo = genOrderNo(orderType)
		o.OrderType = orderType
		o.SourceWarehouseID = sourceWarehouseID
		o.TargetWarehouseID = targetWarehouseID
		o.Status = constants.OrderStatusDraft
		o.CreatorID = creatorID
		o.Remark = remark
		if err := s.repo.CreateTx(tx, o); err != nil {
			s.logger.Error(constants.LogOrderCreateFailed, "error", err.Error())
			return util.Wrap(err, "StockOrder[type=%s] create failed", orderType)
		}
		for i := range items {
			items[i].OrderID = o.ID
		}
		if err := s.itemRepo.CreateManyTx(tx, items); err != nil {
			return util.Wrap(err, "StockOrder[id=%d] create items failed", o.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderCreateSuccess, "order_id", o.ID, "order_no", o.OrderNo)
	return o, nil
}

// Submit 提交单据（draft -> submitted）。
func (s *StockOrderService) Submit(id uint64) (*model.StockOrder, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] submit find failed", id)
	}
	if o.Status != constants.OrderStatusDraft {
		s.logger.Warn(constants.LogOrderSubmitFailed, "order_id", id, "status", o.Status)
		return nil, util.NewAppError(constants.CodeOrderStatusConflict, "StockOrder[id="+u64(id)+"] submit conflict: status="+o.Status)
	}
	o.Status = constants.OrderStatusSubmitted
	if err := s.repo.Update(o); err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] submit save failed", id)
	}
	s.logger.Info(constants.LogOrderSubmitSuccess, "order_id", o.ID)
	return o, nil
}

// Approve 审批（submitted -> processing）。
func (s *StockOrderService) Approve(id, approverID uint64) (*model.StockOrder, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] approve find failed", id)
	}
	if o.Status != constants.OrderStatusSubmitted {
		return nil, util.NewAppError(constants.CodeOrderStatusConflict, "StockOrder[id="+u64(id)+"] approve conflict: status="+o.Status)
	}
	o.Status = constants.OrderStatusProcessing
	o.ApproverID = approverID
	if err := s.repo.Update(o); err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] approve save failed", id)
	}
	s.logger.Info(constants.LogOrderApproveSuccess, "order_id", o.ID)
	return o, nil
}

// Execute 执行单据（processing -> completed），按类型应用库存变化。
func (s *StockOrderService) Execute(id uint64) (*model.StockOrder, error) {
	o := &model.StockOrder{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		cur, err := s.repo.FindByIDForUpdate(tx, id)
		if err != nil {
			return util.Wrap(err, "StockOrder[id=%d] execute find failed", id)
		}
		if cur.Status != constants.OrderStatusProcessing && cur.Status != constants.OrderStatusSubmitted {
			return util.NewAppError(constants.CodeOrderStatusConflict, "StockOrder[id="+u64(id)+"] execute conflict: status="+cur.Status)
		}
		items, err := s.itemRepo.ListByOrderTx(tx, id)
		if err != nil {
			return err
		}
		switch cur.OrderType {
		case constants.OrderTypeInbound:
			for _, it := range items {
				if err := s.stockSvc.InboundTx(tx, it.ProductID, cur.TargetWarehouseID, it.ShelfID, it.Quantity, cur.OrderNo); err != nil {
					s.logger.Error(constants.LogOrderExecuteFailed, "error", err.Error())
					return err
				}
			}
		case constants.OrderTypeOutbound:
			for _, it := range items {
				if err := s.stockSvc.OutboundTx(tx, it.ProductID, cur.SourceWarehouseID, it.ShelfID, it.Quantity); err != nil {
					s.logger.Error(constants.LogOrderExecuteFailed, "error", err.Error())
					return err
				}
			}
		case constants.OrderTypeTransfer:
			for _, it := range items {
				if err := s.stockSvc.TransferTx(tx, it.ProductID, cur.SourceWarehouseID, cur.TargetWarehouseID, it.ShelfID, it.Quantity); err != nil {
					return err
				}
			}
		case constants.OrderTypeInventoryCheck:
			for _, it := range items {
				if err := s.stockSvc.AdjustTx(tx, it.ProductID, cur.SourceWarehouseID, it.ShelfID, it.ActualQuantity); err != nil {
					return err
				}
			}
		}
		cur.Status = constants.OrderStatusCompleted
		if err := s.repo.UpdateTx(tx, cur); err != nil {
			return util.Wrap(err, "StockOrder[id=%d] execute save failed", id)
		}
		o = cur
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderExecuteSuccess, "order_id", o.ID, "order_type", o.OrderType)
	return o, nil
}

// Cancel 取消单据（draft/submitted -> cancelled）。
func (s *StockOrderService) Cancel(id uint64) (*model.StockOrder, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] cancel find failed", id)
	}
	if o.Status != constants.OrderStatusDraft && o.Status != constants.OrderStatusSubmitted {
		return nil, util.NewAppError(constants.CodeOrderStatusConflict, "StockOrder[id="+u64(id)+"] cancel conflict: status="+o.Status)
	}
	o.Status = constants.OrderStatusCancelled
	if err := s.repo.Update(o); err != nil {
		return nil, util.Wrap(err, "StockOrder[id=%d] cancel save failed", id)
	}
	s.logger.Info(constants.LogOrderCancelSuccess, "order_id", o.ID)
	return o, nil
}

// List 分页查询单据。
func (s *StockOrderService) List(page, pageSize int, orderType, status string) ([]model.StockOrder, int64, error) {
	return s.repo.List(page, pageSize, orderType, status)
}

// GetWithItems 单据详情 + 明细。
func (s *StockOrderService) GetWithItems(id uint64) (*model.StockOrder, []model.StockOrderItem, error) {
	o, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, util.Wrap(err, "StockOrder[id=%d] get failed", id)
	}
	items, err := s.itemRepo.ListByOrder(id)
	if err != nil {
		return nil, nil, err
	}
	return o, items, nil
}

func genOrderNo(orderType string) string {
	prefix := "IN"
	switch orderType {
	case constants.OrderTypeOutbound:
		prefix = "OUT"
	case constants.OrderTypeTransfer:
		prefix = "TR"
	case constants.OrderTypeInventoryCheck:
		prefix = "CK"
	}
	return fmt.Sprintf("%s%d%06d", prefix, time.Now().Year(), time.Now().UnixNano()%1000000)
}
