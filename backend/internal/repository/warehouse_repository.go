package repository

import (
	"errors"
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
)

// WarehouseRepository 仓库仓储。
type WarehouseRepository struct {
	db *gorm.DB
}

// NewWarehouseRepository 构造仓库仓储。
func NewWarehouseRepository(db *gorm.DB) *WarehouseRepository {
	return &WarehouseRepository{db: db}
}

// Create 创建仓库。
func (r *WarehouseRepository) Create(w *model.Warehouse) error {
	if err := r.db.Create(w).Error; err != nil {
		return fmt.Errorf("create warehouse: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询仓库。
func (r *WarehouseRepository) FindByID(id uint64) (*model.Warehouse, error) {
	var w model.Warehouse
	if err := r.db.First(&w, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find warehouse by id: %w", err)
	}
	return &w, nil
}

// List 分页查询仓库。
func (r *WarehouseRepository) List(page, pageSize int) ([]model.Warehouse, int64, error) {
	var list []model.Warehouse
	var total int64
	q := r.db.Model(&model.Warehouse{})
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count warehouses: %w", err)
	}
	if err := q.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list warehouses: %w", err)
	}
	return list, total, nil
}

// Update 更新仓库。
func (r *WarehouseRepository) Update(w *model.Warehouse) error {
	if err := r.db.Save(w).Error; err != nil {
		return fmt.Errorf("update warehouse: %w", err)
	}
	return nil
}
