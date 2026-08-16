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
