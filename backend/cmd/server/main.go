package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"warehousestock/internal/config"
	"warehousestock/internal/handler"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/router"
	"warehousestock/internal/service"
	"warehousestock/internal/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	logger := util.NewLogger(slog.LevelInfo)

	db, err := gorm.Open(mysql.Open(cfg.DBDSN()), &gorm.Config{})
	if err != nil {
		logger.Error("connect database failed", "error", err.Error())
		os.Exit(1)
	}
	if err := db.AutoMigrate(
		&model.User{}, &model.Warehouse{}, &model.Shelf{}, &model.Category{}, &model.Product{},
		&model.StockRecord{}, &model.StockOrder{}, &model.StockOrderItem{}, &model.OperationLog{},
	); err != nil {
		logger.Error("auto migrate failed", "error", err.Error())
		os.Exit(1)
	}
	if err := service.NewSeedService(db, logger).Seed(); err != nil {
		logger.Error("seed failed", "error", err.Error())
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	shelfRepo := repository.NewShelfRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	recordRepo := repository.NewStockRecordRepository(db)
	orderRepo := repository.NewStockOrderRepository(db)
	orderItemRepo := repository.NewStockOrderItemRepository(db)

	userSvc := service.NewUserService(userRepo, logger)
	warehouseSvc := service.NewWarehouseService(warehouseRepo, shelfRepo, userRepo, logger)
	productSvc := service.NewProductService(productRepo, categoryRepo, logger)
	stockSvc := service.NewStockRecordService(db, recordRepo, productRepo, warehouseRepo, logger)
	orderSvc := service.NewStockOrderService(db, orderRepo, orderItemRepo, stockSvc, productRepo, warehouseRepo, logger)
	dashboardSvc := service.NewDashboardService(stockSvc, orderRepo, warehouseSvc, logger)

	userHandler := handler.NewUserHandler(userSvc, logger)
	warehouseHandler := handler.NewWarehouseHandler(warehouseSvc, logger)
	productHandler := handler.NewProductHandler(productSvc, logger)
	stockHandler := handler.NewStockRecordHandler(stockSvc, logger)
	orderHandler := handler.NewStockOrderHandler(orderSvc, logger)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc, logger)
	uploadHandler := handler.NewUploadHandler(cfg, logger)
	opLogHandler := handler.NewOperationLogHandler(db, logger)

	r := router.New(cfg, db, logger, userHandler, warehouseHandler, productHandler,
		stockHandler, orderHandler, dashboardHandler, uploadHandler, opLogHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r.Setup(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server run failed", "error", err.Error())
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("server shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err.Error())
	}
	logger.Info("server stopped")
}
