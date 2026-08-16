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

func newAdjustFixture(t *testing.T) (*gorm.DB, *repository.ProductRepository, *repository.WarehouseRepository, *repository.StockRecordRepository, *StockRecordService, *StockOrderService) {
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

func TestAdjustSetsStockToActual(t *testing.T) {
	_, productRepo, warehouseRepo, recordRepo, recordSvc, orderSvc := newAdjustFixture(t)

	p := &model.Product{Name: "气泡膜", SKU: "SKU-ADJ-1", Unit: "卷", MinStock: 0, MaxStock: 0}
	if err := productRepo.Create(p); err != nil {
		t.Fatalf("create product: %v", err)
	}
	w := &model.Warehouse{Name: "华东仓", Code: "WH-ADJ", Status: "active"}
	if err := warehouseRepo.Create(w); err != nil {
		t.Fatalf("create warehouse: %v", err)
	}
	now := time.Now()
	r1 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "B1", Quantity: 60, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	r2 := &model.StockRecord{ProductID: p.ID, WarehouseID: w.ID, BatchNo: "B2", Quantity: 40, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	if err := recordRepo.Create(r1); err != nil {
		t.Fatalf("create r1: %v", err)
	}
	if err := recordRepo.Create(r2); err != nil {
		t.Fatalf("create r2: %v", err)
	}

	if err := recordSvc.Adjust(p.ID, w.ID, 0, 120); err != nil {
		t.Fatalf("Adjust: %v", err)
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
	if total != 120 {
		t.Fatalf("total after adjust = %d, want 120", total)
	}
	if oldest == nil || oldest.Quantity != 80 {
		t.Fatalf("oldest batch after adjust = %d, want 80", oldest.Quantity)
	}

	items := []model.StockOrderItem{{ProductID: p.ID, Quantity: 130}}
	order, err := orderSvc.Create(1, constants.OrderTypeInventoryCheck, w.ID, 0, "", items)
	if err != nil {
		t.Fatalf("Create inventory_check: %v", err)
	}
	if _, err := orderSvc.Submit(order.ID); err != nil {
		t.Fatalf("Submit inventory_check: %v", err)
	}
	if _, err := orderSvc.Execute(order.ID); err != nil {
		t.Fatalf("Execute inventory_check: %v", err)
	}
	records, err = recordRepo.FindByProductWarehouse(p.ID, w.ID)
	if err != nil {
		t.Fatalf("find records after order: %v", err)
	}
	total = 0
	for _, r := range records {
		total += r.Quantity
	}
	if total != 130 {
		t.Fatalf("total after inventory_check order = %d, want 130", total)
	}
}
