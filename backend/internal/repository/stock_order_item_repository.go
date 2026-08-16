package repository

import (
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
)

// StockOrderItemRepository 单据明细仓储。
type StockOrderItemRepository struct {
	db *gorm.DB
}

// NewStockOrderItemRepository 构造单据明细仓储。
func NewStockOrderItemRepository(db *gorm.DB) *StockOrderItemRepository {
	return &StockOrderItemRepository{db: db}
}

// CreateMany 批量创建明细。
func (r *StockOrderItemRepository) CreateMany(items []model.StockOrderItem) error {
	return r.CreateManyTx(r.db, items)
}

// CreateManyTx 在事务内批量创建明细。
func (r *StockOrderItemRepository) CreateManyTx(tx *gorm.DB, items []model.StockOrderItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := tx.Create(&items).Error; err != nil {
		return fmt.Errorf("create stock order items: %w", err)
	}
	return nil
}

// ListByOrder 查询某单据明细。
func (r *StockOrderItemRepository) ListByOrder(orderID uint64) ([]model.StockOrderItem, error) {
	return r.ListByOrderTx(r.db, orderID)
}

// ListByOrderTx 在事务内查询某单据明细。
func (r *StockOrderItemRepository) ListByOrderTx(tx *gorm.DB, orderID uint64) ([]model.StockOrderItem, error) {
	var list []model.StockOrderItem
	if err := tx.Where("order_id = ?", orderID).Order("id ASC").Find(&list).Error; err != nil {
		return nil, fmt.Errorf("list stock order items: %w", err)
	}
	return list, nil
}

// Update 更新明细。
func (r *StockOrderItemRepository) Update(item *model.StockOrderItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return fmt.Errorf("update stock order item: %w", err)
	}
	return nil
}
