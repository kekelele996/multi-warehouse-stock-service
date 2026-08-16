package service

import (
	"log/slog"
	"time"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedService 启动时幂等写入预置数据。
type SeedService struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewSeedService 构造种子服务。
func NewSeedService(db *gorm.DB, logger *slog.Logger) *SeedService {
	return &SeedService{db: db, logger: logger}
}

// Seed 当 users 表为空时写入种子数据。
func (s *SeedService) Seed() error {
	var count int64
	if err := s.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	userHash, _ := bcrypt.GenerateFromPassword([]byte("User@123"), bcrypt.DefaultCost)
	users := []model.User{
		{Phone: "13800000001", PasswordHash: string(adminHash), Name: "系统管理员", Role: constants.RoleAdmin},
		{Phone: "13800000002", PasswordHash: string(userHash), Name: "仓库经理", Role: constants.RoleWarehouseManager},
		{Phone: "13800000003", PasswordHash: string(userHash), Name: "仓管员", Role: constants.RoleOperator},
		{Phone: "13800000004", PasswordHash: string(userHash), Name: "观察员", Role: constants.RoleViewer},
	}
	for i := range users {
		if err := s.db.Create(&users[i]).Error; err != nil {
			return err
		}
	}
	warehouses := []model.Warehouse{
		{Name: "华东一号仓", Code: "WH-001", Address: "上海市青浦区物流园 1 号", AreaSqm: 12000, ManagerID: 2, Status: "active", Phone: "021-88880001"},
		{Name: "华南二号仓", Code: "WH-002", Address: "广州市黄埔区仓储大道 8 号", AreaSqm: 8000, ManagerID: 2, Status: "active", Phone: "020-88880002"},
	}
	for i := range warehouses {
		if err := s.db.Create(&warehouses[i]).Error; err != nil {
			return err
		}
	}
	shelves := []model.Shelf{
		{WarehouseID: 1, ShelfNo: "A-01", LayerCount: 3, ColumnCount: 4, Capacity: 120},
		{WarehouseID: 1, ShelfNo: "A-02", LayerCount: 3, ColumnCount: 4, Capacity: 120},
		{WarehouseID: 2, ShelfNo: "B-01", LayerCount: 2, ColumnCount: 4, Capacity: 80},
	}
	for i := range shelves {
		if err := s.db.Create(&shelves[i]).Error; err != nil {
			return err
		}
	}
	categories := []model.Category{{Name: "电子元器件", SortOrder: 1}, {Name: "包装材料", SortOrder: 2}, {Name: "成品", SortOrder: 3}}
	for i := range categories {
		if err := s.db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}
	products := []model.Product{
		{Name: "电阻 10KΩ", SKU: "ELEC-R10K", CategoryID: 1, Spec: "0805", Unit: "件", MinStock: 100, MaxStock: 5000},
		{Name: "电容 100uF", SKU: "ELEC-C100", CategoryID: 1, Spec: "1206", Unit: "件", MinStock: 200, MaxStock: 8000},
		{Name: "气泡膜", SKU: "PKG-BUBBLE", CategoryID: 2, Spec: "50cm*100m", Unit: "卷", MinStock: 10, MaxStock: 200, ShelfLifeDays: 365},
	}
	for i := range products {
		if err := s.db.Create(&products[i]).Error; err != nil {
			return err
		}
	}
	now := time.Now()
	records := []model.StockRecord{
		{ProductID: 1, WarehouseID: 1, ShelfID: 1, BatchNo: "B20260801", Quantity: 2000, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now},
		{ProductID: 2, WarehouseID: 1, ShelfID: 2, BatchNo: "B20260802", Quantity: 3500, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now},
		{ProductID: 3, WarehouseID: 2, ShelfID: 3, BatchNo: "B20260803", Quantity: 60, InboundDate: &now, LastOpType: constants.StockOpInbound, LastOpAt: now},
	}
	for i := range records {
		if err := s.db.Create(&records[i]).Error; err != nil {
			return err
		}
	}
	orders := []model.StockOrder{
		{OrderNo: "IN202608160001", OrderType: constants.OrderTypeInbound, TargetWarehouseID: 1, Status: constants.OrderStatusCompleted, CreatorID: 3, ApproverID: 2, Remark: "电阻入库"},
		{OrderNo: "OUT202608160001", OrderType: constants.OrderTypeOutbound, SourceWarehouseID: 1, Status: constants.OrderStatusSubmitted, CreatorID: 3, Remark: "电容出库"},
	}
	for i := range orders {
		if err := s.db.Create(&orders[i]).Error; err != nil {
			return err
		}
	}
	orderItems := []model.StockOrderItem{
		{OrderID: 1, ProductID: 1, ShelfID: 1, Quantity: 500, ActualQuantity: 500},
		{OrderID: 2, ProductID: 2, ShelfID: 2, Quantity: 200},
	}
	for i := range orderItems {
		if err := s.db.Create(&orderItems[i]).Error; err != nil {
			return err
		}
	}
	s.logger.Info("seed data created")
	return nil
}
