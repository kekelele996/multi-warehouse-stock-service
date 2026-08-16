package service

import (
	"log/slog"
	"testing"
	"time"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newOutboundFixture(t *testing.T) (*gorm.DB, *repository.ProductRepository, *repository.WarehouseRepository, *repository.StockRecordRepository, *StockRecordService, *StockOrderService) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	models := []any{&model.User{}, &model.Warehouse{}, &model.Shelf{}, &model.Category{},
		&model.Product{}, &model.StockOrder{}, &model.StockOrderItem{}, &model.StockRecord{}}
	if err := db.AutoMigrate(models...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	log := slog.Default()
	productRepo := repository.NewProductRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	recordRepo := repository.NewStockRecordRepository(db)
	orderRepo := repository.NewStockOrderRepository(db)
	itemRepo := repository.NewStockOrderItemRepository(db)
	recordSvc := NewStockRecordService(db, recordRepo, productRepo, warehouseRepo, log)
	orderSvc := NewStockOrderService(db, orderRepo, itemRepo, recordSvc, productRepo, warehouseRepo, log)
	return db, productRepo, warehouseRepo, recordRepo, recordSvc, orderSvc
}

func TestOutboundOrderFIFOAndExactStock(t *testing.T) {
	_, productRepo, warehouseRepo, recordRepo, _, orderSvc := newOutboundFixture(t)

	p := &model.Product{Name: "电阻", SKU: "SKU-OUT-1", Unit: "件", MinStock: 0, MaxStock: 0}
	if err := productRepo.Create(p); err != nil {
		t.Fatalf("create product: %v", err)
	}
	w := &model.Warehouse{Name: "华东仓", Code: "WH-OUT", Status: "active"}
	if err := warehouseRepo.Create(w); err != nil {
		t.Fatalf("create warehouse: %v", err)
	}
	now := time.Now()
	b1 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "B1", Quantity: 100, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	b2 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "B2", Quantity: 50, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	if err := recordRepo.Create(b1); err != nil {
		t.Fatalf("create batch1: %v", err)
	}
	if err := recordRepo.Create(b2); err != nil {
		t.Fatalf("create batch2: %v", err)
	}

	items := []model.StockOrderItem{{ProductID: p.ID, Quantity: 150}}
	order, err := orderSvc.Create(1, constants.OrderTypeOutbound, w.ID, 0, "", items)
	if err != nil {
		t.Fatalf("Create order: %v", err)
	}
	if _, err := orderSvc.Submit(order.ID); err != nil {
		t.Fatalf("Submit order: %v", err)
	}
	executed, err := orderSvc.Execute(order.ID)
	if err != nil {
		t.Fatalf("Execute order: %v", err)
	}
	if executed.Status != constants.OrderStatusCompleted {
		t.Fatalf("order status = %s, want completed", executed.Status)
	}

	records, err := recordRepo.FindByProductWarehouse(p.ID, w.ID)
	if err != nil {
		t.Fatalf("find records: %v", err)
	}
	total := 0
	var oldest *model.StockRecord
	for i := range records {
		total += records[i].Quantity
		if oldest == nil || records[i].ID < oldest.ID {
			oldest = &records[i]
		}
	}
	if total != 0 {
		t.Fatalf("total stock after outbound = %d, want 0", total)
	}
	if oldest == nil || oldest.Quantity != 0 {
		t.Fatalf("oldest batch quantity = %d, want 0 (FIFO)", oldest.Quantity)
	}
}

// TestOutboundOrderFIFOPartialDeduct 验证出库按 FIFO 从最老批次开始扣减，
// 当只需消耗部分老批次时，新批次应保持原样、老批次只扣减掉对应数量。
func TestOutboundOrderFIFOPartialDeduct(t *testing.T) {
	_, productRepo, warehouseRepo, recordRepo, _, orderSvc := newOutboundFixture(t)

	p := &model.Product{Name: "电容", SKU: "SKU-OUT-2", Unit: "件", MinStock: 0, MaxStock: 0}
	if err := productRepo.Create(p); err != nil {
		t.Fatalf("create product: %v", err)
	}
	w := &model.Warehouse{Name: "华南仓", Code: "WH-OUT2", Status: "active"}
	if err := warehouseRepo.Create(w); err != nil {
		t.Fatalf("create warehouse: %v", err)
	}
	now := time.Now()
	b1 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "OLD", Quantity: 100, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	b2 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "NEW", Quantity: 100, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	if err := recordRepo.Create(b1); err != nil {
		t.Fatalf("create batch1: %v", err)
	}
	if err := recordRepo.Create(b2); err != nil {
		t.Fatalf("create batch2: %v", err)
	}

	items := []model.StockOrderItem{{ProductID: p.ID, Quantity: 70}}
	order, err := orderSvc.Create(1, constants.OrderTypeOutbound, w.ID, 0, "", items)
	if err != nil {
		t.Fatalf("Create order: %v", err)
	}
	if _, err := orderSvc.Submit(order.ID); err != nil {
		t.Fatalf("Submit order: %v", err)
	}
	if _, err := orderSvc.Execute(order.ID); err != nil {
		t.Fatalf("Execute order: %v", err)
	}

	records, err := recordRepo.FindByProductWarehouse(p.ID, w.ID)
	if err != nil {
		t.Fatalf("find records: %v", err)
	}
	byID := map[uint64]int{}
	for i := range records {
		byID[records[i].ID] = records[i].Quantity
	}
	// FIFO: 70 全部来自老批次 b1（100 -> 30），新批次 b2 保持 100 不动。
	if byID[b1.ID] != 30 {
		t.Fatalf("oldest batch quantity = %d, want 30 (FIFO partial deduct)", byID[b1.ID])
	}
	if byID[b2.ID] != 100 {
		t.Fatalf("newest batch quantity = %d, want 100 (untouched)", byID[b2.ID])
	}
}

// TestOutboundOrderInsufficientStock 验证库存不足时返回库存不足错误，
// 并且不会改动任何批次数量（事务回滚）。
func TestOutboundOrderInsufficientStock(t *testing.T) {
	_, productRepo, warehouseRepo, recordRepo, _, orderSvc := newOutboundFixture(t)

	p := &model.Product{Name: "电感", SKU: "SKU-OUT-3", Unit: "件", MinStock: 0, MaxStock: 0}
	if err := productRepo.Create(p); err != nil {
		t.Fatalf("create product: %v", err)
	}
	w := &model.Warehouse{Name: "华北仓", Code: "WH-OUT3", Status: "active"}
	if err := warehouseRepo.Create(w); err != nil {
		t.Fatalf("create warehouse: %v", err)
	}
	now := time.Now()
	b1 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "OLD", Quantity: 100, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	b2 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "NEW", Quantity: 50, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	if err := recordRepo.Create(b1); err != nil {
		t.Fatalf("create batch1: %v", err)
	}
	if err := recordRepo.Create(b2); err != nil {
		t.Fatalf("create batch2: %v", err)
	}

	// 总库存 150，需 200，应报库存不足。
	items := []model.StockOrderItem{{ProductID: p.ID, Quantity: 200}}
	order, err := orderSvc.Create(1, constants.OrderTypeOutbound, w.ID, 0, "", items)
	if err != nil {
		t.Fatalf("Create order: %v", err)
	}
	if _, err := orderSvc.Submit(order.ID); err != nil {
		t.Fatalf("Submit order: %v", err)
	}
	if _, err := orderSvc.Execute(order.ID); err == nil {
		t.Fatalf("Execute should fail with insufficient stock")
	}

	records, err := recordRepo.FindByProductWarehouse(p.ID, w.ID)
	if err != nil {
		t.Fatalf("find records: %v", err)
	}
	byID := map[uint64]int{}
	for i := range records {
		byID[records[i].ID] = records[i].Quantity
	}
	if byID[b1.ID] != 100 || byID[b2.ID] != 50 {
		t.Fatalf("stock changed after failed outbound: b1=%d b2=%d, want 100/50", byID[b1.ID], byID[b2.ID])
	}
}
