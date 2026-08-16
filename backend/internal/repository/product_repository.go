package repository

import (
	"errors"
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
)

// ProductRepository 商品仓储。
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 构造商品仓储。
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create 创建商品。
func (r *ProductRepository) Create(p *model.Product) error {
	if err := r.db.Create(p).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询商品。
func (r *ProductRepository) FindByID(id uint64) (*model.Product, error) {
	var p model.Product
	if err := r.db.First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find product by id: %w", err)
	}
	return &p, nil
}

// List 分页查询商品，支持分类/关键词筛选。
func (r *ProductRepository) List(page, pageSize int, categoryID uint64, keyword string) ([]model.Product, int64, error) {
	var list []model.Product
	var total int64
	q := r.db.Model(&model.Product{})
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q = q.Where("name LIKE ? OR sku LIKE ? OR barcode LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return list, total, nil
}

// Update 更新商品。
func (r *ProductRepository) Update(p *model.Product) error {
	if err := r.db.Save(p).Error; err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}
