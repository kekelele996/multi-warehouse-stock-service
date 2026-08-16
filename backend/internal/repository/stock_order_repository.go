package repository

import (
	"errors"
	"fmt"
	"time"

	"warehousestock/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StockOrderRepository 出入库单仓储。
type StockOrderRepository struct {
	db *gorm.DB
}

// NewStockOrderRepository 构造出入库单仓储。
func NewStockOrderRepository(db *gorm.DB) *StockOrderRepository {
	return &StockOrderRepository{db: db}
}

// Create 创建单据。
func (r *StockOrderRepository) Create(o *model.StockOrder) error {
	return r.CreateTx(r.db, o)
}

// CreateTx 在事务内创建单据。
func (r *StockOrderRepository) CreateTx(tx *gorm.DB, o *model.StockOrder) error {
	if err := tx.Create(o).Error; err != nil {
		return fmt.Errorf("create stock order: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询单据。
func (r *StockOrderRepository) FindByID(id uint64) (*model.StockOrder, error) {
	return r.findByID(r.db, id, false)
}

// FindByIDForUpdate 在事务内锁定单据行。
func (r *StockOrderRepository) FindByIDForUpdate(tx *gorm.DB, id uint64) (*model.StockOrder, error) {
	return r.findByID(tx, id, true)
}

func (r *StockOrderRepository) findByID(db *gorm.DB, id uint64, forUpdate bool) (*model.StockOrder, error) {
	var o model.StockOrder
	q := db
	if forUpdate {
		q = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.First(&o, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("find stock order by id: %v", ErrNotFound)
		}
		return nil, fmt.Errorf("find stock order by id: %w", err)
	}
	return &o, nil
}

// List 分页查询单据，支持类型/状态筛选。
func (r *StockOrderRepository) List(page, pageSize int, orderType, status string) ([]model.StockOrder, int64, error) {
	var list []model.StockOrder
	var total int64
	q := r.db.Model(&model.StockOrder{})
	if orderType != "" {
		q = q.Where("order_type = ?", orderType)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count stock orders: %w", err)
	}
	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("list stock orders: %w", err)
	}
	return list, total, nil
}

// Update 更新单据。
func (r *StockOrderRepository) Update(o *model.StockOrder) error {
	return r.UpdateTx(r.db, o)
}

// UpdateTx 在事务内更新单据。
func (r *StockOrderRepository) UpdateTx(tx *gorm.DB, o *model.StockOrder) error {
	if err := tx.Save(o).Error; err != nil {
		return fmt.Errorf("update stock order: %w", err)
	}
	return nil
}

// CountToday 今日单据数。
func (r *StockOrderRepository) CountToday() (int64, error) {
	start := time.Now()
	begin := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	end := begin.AddDate(0, 0, 1)
	var n int64
	if err := r.db.Model(&model.StockOrder{}).Where("created_at >= ? AND created_at < ?", begin, end).Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count today orders: %w", err)
	}
	return n, nil
}
