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

func newTransferFixture(t *testing.T) (*gorm.DB, *repository.ProductRepository, *repository.WarehouseRepository, *repository.StockRecordRepository, *StockRecordService, *StockOrderService) {
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

func TestTransferMovesStockBetweenWarehouses(t *testing.T) {
	_, productRepo, warehouseRepo, recordRepo, _, orderSvc := newTransferFixture(t)

	p := &model.Product{Name: "电容", SKU: "SKU-TR-1", Unit: "件", MinStock: 0, MaxStock: 0}
	if err := productRepo.Create(p); err != nil {
		t.Fatalf("create product: %v", err)
	}
	from := &model.Warehouse{Name: "华东仓", Code: "WH-FROM", Status: "active"}
	to := &model.Warehouse{Name: "华南仓", Code: "WH-TO", Status: "active"}
	if err := warehouseRepo.Create(from); err != nil {
		t.Fatalf("create from warehouse: %v", err)
	}
	if err := warehouseRepo.Create(to); err != nil {
		t.Fatalf("create to warehouse: %v", err)
	}
	now := time.Now()
	r1 := &model.StockRecord{ProductID: p.ID, WarehouseID: from.ID, BatchNo: "B1", Quantity: 60, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	r2 := &model.StockRecord{ProductID: p.ID, WarehouseID: from.ID, BatchNo: "B2", Quantity: 40, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now}
	if err := recordRepo.Create(r1); err != nil {
		t.Fatalf("create r1: %v", err)
	}
	if err := recordRepo.Create(r2); err != nil {
		t.Fatalf("create r2: %v", err)
	}

	items := []model.StockOrderItem{{ProductID: p.ID, Quantity: 40}}
	order, err := orderSvc.Create(1, constants.OrderTypeTransfer, from.ID, to.ID, "", items)
	if err != nil {
		t.Fatalf("Create transfer: %v", err)
	}
	if _, err := orderSvc.Submit(order.ID); err != nil {
		t.Fatalf("Submit transfer: %v", err)
	}
	if _, err := orderSvc.Execute(order.ID); err != nil {
		t.Fatalf("Execute transfer: %v", err)
	}

	fromRecords, err := recordRepo.FindByProductWarehouse(p.ID, from.ID)
	if err != nil {
		t.Fatalf("find from records: %v", err)
	}
	toRecords, err := recordRepo.FindByProductWarehouse(p.ID, to.ID)
	if err != nil {
		t.Fatalf("find to records: %v", err)
	}
	fromTotal, toTotal := 0, 0
	for _, r := range fromRecords {
		fromTotal += r.Quantity
	}
	for _, r := range toRecords {
		toTotal += r.Quantity
	}
	if fromTotal != 60 {
		t.Fatalf("source stock = %d, want 60", fromTotal)
	}
	if toTotal != 40 {
		t.Fatalf("target stock = %d, want 40", toTotal)
	}

	summary, err := recordRepo.SummaryByWarehouse()
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	got := map[uint64]float64{}
	for _, row := range summary {
		wid := num(row["warehouse_id"])
		qty := num(row["quantity"])
		got[uint64(wid)] = qty
	}
	if got[from.ID] != 60 || got[to.ID] != 40 {
		t.Fatalf("summary by warehouse = %v, want from=60 to=40", got)
	}
}

func num(v any) float64 {
	switch x := v.(type) {
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case uint64:
		return float64(x)
	case float64:
		return x
	case float32:
		return float64(x)
	default:
		return 0
	}
}
