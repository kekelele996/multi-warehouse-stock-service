package service

import (
	"log/slog"

	"warehousestock/internal/constants"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/util"
)

// WarehouseService 仓库业务逻辑。
type WarehouseService struct {
	repo      *repository.WarehouseRepository
	shelfRepo *repository.ShelfRepository
	userRepo  *repository.UserRepository
	logger    *slog.Logger
}

// NewWarehouseService 构造仓库服务。
func NewWarehouseService(repo *repository.WarehouseRepository, shelfRepo *repository.ShelfRepository,
	userRepo *repository.UserRepository, logger *slog.Logger) *WarehouseService {
	return &WarehouseService{repo: repo, shelfRepo: shelfRepo, userRepo: userRepo, logger: logger}
}

// Create 创建仓库。
func (s *WarehouseService) Create(name, code, address string, areaSqm float64, managerID uint64, phone string) (*model.Warehouse, error) {
	w := &model.Warehouse{Name: name, Code: code, Address: address, AreaSqm: areaSqm, ManagerID: managerID, Status: "active", Phone: phone}
	if err := s.repo.Create(w); err != nil {
		return nil, util.Wrap(err, "Warehouse[code=%s] create failed", code)
	}
	s.logger.Info(constants.LogWarehouseCreateSuccess, "warehouse_id", w.ID)
	return w, nil
}

// Update 更新仓库。
func (s *WarehouseService) Update(id uint64, name, address string, areaSqm float64, phone string) (*model.Warehouse, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Warehouse[id=%d] update find failed", id)
	}
	if name != "" {
		w.Name = name
	}
	if address != "" {
		w.Address = address
	}
	if areaSqm > 0 {
		w.AreaSqm = areaSqm
	}
	if phone != "" {
		w.Phone = phone
	}
	if err := s.repo.Update(w); err != nil {
		return nil, util.Wrap(err, "Warehouse[id=%d] update save failed", id)
	}
	s.logger.Info(constants.LogWarehouseUpdateSuccess, "warehouse_id", w.ID)
	return w, nil
}

// ChangeStatus 启用/停用仓库。
func (s *WarehouseService) ChangeStatus(id uint64, status string) (*model.Warehouse, error) {
	if status != "active" && status != "inactive" && status != "maintenance" {
		return nil, util.NewAppError(constants.CodeValidationFailed, "Warehouse[id="+u64(id)+"] status invalid: "+status)
	}
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Warehouse[id=%d] status find failed", id)
	}
	w.Status = status
	if err := s.repo.Update(w); err != nil {
		return nil, util.Wrap(err, "Warehouse[id=%d] status save failed", id)
	}
	s.logger.Info(constants.LogWarehouseStatusChange, "warehouse_id", w.ID, "status", status)
	return w, nil
}

// List 分页查询仓库。
func (s *WarehouseService) List(page, pageSize int) ([]model.Warehouse, int64, error) {
	return s.repo.List(page, pageSize)
}

// GetWithShelves 仓库详情 + 库位列表。
func (s *WarehouseService) GetWithShelves(id uint64) (*model.Warehouse, []model.Shelf, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, util.Wrap(err, "Warehouse[id=%d] get failed", id)
	}
	shelves, err := s.shelfRepo.ListByWarehouse(id)
	if err != nil {
		return nil, nil, err
	}
	return w, shelves, nil
}

// CreateShelf 创建库位。
func (s *WarehouseService) CreateShelf(warehouseID uint64, shelfNo string, layerCount, columnCount, capacity int) (*model.Shelf, error) {
	if _, err := s.repo.FindByID(warehouseID); err != nil {
		return nil, util.Wrap(err, "Shelf[warehouse_id=%d] create: warehouse not found", warehouseID)
	}
	sh := &model.Shelf{WarehouseID: warehouseID, ShelfNo: shelfNo, LayerCount: layerCount, ColumnCount: columnCount, Capacity: capacity}
	if err := s.shelfRepo.Create(sh); err != nil {
		return nil, util.Wrap(err, "Shelf[shelf_no=%s] create failed", shelfNo)
	}
	s.logger.Info(constants.LogShelfCreateSuccess, "shelf_id", sh.ID)
	return sh, nil
}

// UpdateShelf 更新库位。
func (s *WarehouseService) UpdateShelf(id uint64, layerCount, columnCount, capacity int) (*model.Shelf, error) {
	sh, err := s.shelfRepo.FindByID(id)
	if err != nil {
		return nil, util.Wrap(err, "Shelf[id=%d] update find failed", id)
	}
	if layerCount > 0 {
		sh.LayerCount = layerCount
	}
	if columnCount > 0 {
		sh.ColumnCount = columnCount
	}
	if capacity > 0 {
		sh.Capacity = capacity
	}
	if err := s.shelfRepo.Update(sh); err != nil {
		return nil, util.Wrap(err, "Shelf[id=%d] update save failed", id)
	}
	s.logger.Info(constants.LogShelfUpdateSuccess, "shelf_id", sh.ID)
	return sh, nil
}

func itoa(v uint64) string {
	return u64(v)
}

func u64(v uint64) string {
	if v == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for v > 0 {
		i--
		b[i] = byte('0' + v%10)
		v /= 10
	}
	return string(b[i:])
}
