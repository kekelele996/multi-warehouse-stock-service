package service

import (
	"log/slog"

	"warehousestock/internal/constants"
	"warehousestock/internal/repository"
)

// DashboardService 仪表盘统计业务逻辑。
type DashboardService struct {
	recordSvc    *StockRecordService
	orderRepo    *repository.StockOrderRepository
	warehouseSvc *WarehouseService
	logger       *slog.Logger
}

// NewDashboardService 构造仪表盘服务。
func NewDashboardService(recordSvc *StockRecordService, orderRepo *repository.StockOrderRepository,
	warehouseSvc *WarehouseService, logger *slog.Logger) *DashboardService {
	return &DashboardService{recordSvc: recordSvc, orderRepo: orderRepo, warehouseSvc: warehouseSvc, logger: logger}
}

// Stats 汇总仪表盘数据。
func (s *DashboardService) Stats() (map[string]any, error) {
	byWarehouse, err := s.recordSvc.repo.SummaryByWarehouse()
	if err != nil {
		return nil, err
	}
	lowStock, err := s.recordSvc.repo.LowStock()
	if err != nil {
		return nil, err
	}
	today, err := s.orderRepo.CountToday()
	if err != nil {
		return nil, err
	}
	warehouses, _, err := s.warehouseSvc.List(1, 100)
	if err != nil {
		return nil, err
	}
	result := map[string]any{
		"warehouse_summary": byWarehouse,
		"low_stock":         lowStock,
		"today_orders":      today,
		"warehouse_count":   len(warehouses),
	}
	s.logger.Info(constants.LogDashboardStats, "data", result)
	return result, nil
}
