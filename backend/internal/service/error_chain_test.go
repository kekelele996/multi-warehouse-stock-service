package service

import (
	"errors"
	"log/slog"
	"testing"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestErrorChainSentinelsAndLogin(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.StockOrder{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewStockOrderRepository(db)
	userSvc := NewUserService(userRepo, slog.Default())

	if _, err := userSvc.Register("13900000001", "secret123", "小明", constants.RoleViewer); err != nil {
		t.Fatalf("Register new phone: %v", err)
	}

	_, _, err = userSvc.Login("secret", 1, "13900009999", "secret123")
	appErr, ok := err.(*util.AppError)
	if !ok || appErr.Code != constants.CodeInvalidCredentials {
		t.Fatalf("Login unknown phone = %v, want CodeInvalidCredentials", err)
	}

	if _, err := orderRepo.FindByID(99999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("FindByID missing = %v, want ErrNotFound", err)
	}
}
