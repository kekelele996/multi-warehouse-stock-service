package service

import (
	"errors"
	"log/slog"
	"testing"

	"warehousestock/internal/model"
	"warehousestock/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func expectErrNotFoundNoPanic(t *testing.T, name string, fn func() error) error {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("%s panicked: %v", name, r)
			}
		}()
		err = fn()
	}()
	return err
}

func TestMissingOrderAndWarehouseReturnNotFound(t *testing.T) {
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
	orderRepo := repository.NewStockOrderRepository(db)
	itemRepo := repository.NewStockOrderItemRepository(db)
	recordRepo := repository.NewStockRecordRepository(db)
	productRepo := repository.NewProductRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	shelfRepo := repository.NewShelfRepository(db)
	userRepo := repository.NewUserRepository(db)

	recordSvc := NewStockRecordService(db, recordRepo, productRepo, warehouseRepo, log)
	orderSvc := NewStockOrderService(db, orderRepo, itemRepo, recordSvc, productRepo, warehouseRepo, log)
	warehouseSvc := NewWarehouseService(warehouseRepo, shelfRepo, userRepo, log)

	err = expectErrNotFoundNoPanic(t, "StockOrderService.Submit", func() error {
		_, err := orderSvc.Submit(99999)
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Submit missing = %v, want ErrNotFound", err)
	}

	err = expectErrNotFoundNoPanic(t, "WarehouseService.ChangeStatus", func() error {
		_, err := warehouseSvc.ChangeStatus(99999, "active")
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("ChangeStatus missing = %v, want ErrNotFound", err)
	}
}
