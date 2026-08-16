package repository

import (
	"errors"
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
)

// ShelfRepository 库位仓储。
type ShelfRepository struct {
	db *gorm.DB
}

// NewShelfRepository 构造库位仓储。
func NewShelfRepository(db *gorm.DB) *ShelfRepository {
	return &ShelfRepository{db: db}
}

// Create 创建库位。
func (r *ShelfRepository) Create(s *model.Shelf) error {
	if err := r.db.Create(s).Error; err != nil {
		return fmt.Errorf("create shelf: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询库位。
func (r *ShelfRepository) FindByID(id uint64) (*model.Shelf, error) {
	var s model.Shelf
	if err := r.db.First(&s, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find shelf by id: %w", err)
	}
	return &s, nil
}

// ListByWarehouse 查询某仓库的库位。
func (r *ShelfRepository) ListByWarehouse(warehouseID uint64) ([]model.Shelf, error) {
	var list []model.Shelf
	if err := r.db.Where("warehouse_id = ?", warehouseID).Order("id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list shelves by warehouse: %w", err)
	}
	return list, nil
}

// Update 更新库位。
func (r *ShelfRepository) Update(s *model.Shelf) error {
	if err := r.db.Save(s).Error; err != nil {
		return fmt.Errorf("update shelf: %w", err)
	}
	return nil
}
