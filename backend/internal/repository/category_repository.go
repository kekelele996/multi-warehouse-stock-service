package repository

import (
	"errors"
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
)

// CategoryRepository 商品分类仓储。
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 构造商品分类仓储。
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create 创建分类。
func (r *CategoryRepository) Create(c *model.Category) error {
	if err := r.db.Create(c).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询分类。
func (r *CategoryRepository) FindByID(id uint64) (*model.Category, error) {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	return &c, nil
}

// List 查询全部分类（树形）。
func (r *CategoryRepository) List() ([]model.Category, error) {
	var list []model.Category
	if err := r.db.Order("sort_order ASC, id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return list, nil
}

// Update 更新分类。
func (r *CategoryRepository) Update(c *model.Category) error {
	if err := r.db.Save(c).Error; err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}
