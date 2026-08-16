package repository

import (
	"errors"
	"fmt"

	"warehousestock/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StockRecordRepository 库存记录仓储。
type StockRecordRepository struct {
	db *gorm.DB
}

// NewStockRecordRepository 构造库存记录仓储。
func NewStockRecordRepository(db *gorm.DB) *StockRecordRepository {
	return &StockRecordRepository{db: db}
}

// Create 创建库存记录。
func (r *StockRecordRepository) Create(rec *model.StockRecord) error {
	return r.CreateTx(r.db, rec)
}

// CreateTx 在事务内创建库存记录。
func (r *StockRecordRepository) CreateTx(tx *gorm.DB, rec *model.StockRecord) error {
	if err := tx.Create(rec).Error; err != nil {
		return fmt.Errorf("create stock record: %w", err)
	}
	return nil
}

// FindByID 按 ID 查询库存记录。
func (r *StockRecordRepository) FindByID(id uint64) (*model.StockRecord, error) {
	var rec model.StockRecord
	if err := r.db.First(&rec, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find stock record by id: %w", err)
	}
	return &rec, nil
}

// FindByProductWarehouse 查询商品在仓库的库存（按批次合并）。
func (r *StockRecordRepository) FindByProductWarehouse(productID, warehouseID uint64) ([]model.StockRecord, error) {
	return r.findByProductWarehouse(r.db, productID, warehouseID, false)
}

// FindByProductWarehouseForUpdate 在事务内锁定商品在仓库的库存记录。
func (r *StockRecordRepository) FindByProductWarehouseForUpdate(tx *gorm.DB, productID, warehouseID uint64) ([]model.StockRecord, error) {
	return r.findByProductWarehouse(tx, productID, warehouseID, true)
}

func (r *StockRecordRepository) findByProductWarehouse(db *gorm.DB, productID, warehouseID uint64, forUpdate bool) ([]model.StockRecord, error) {
	var list []model.StockRecord
	q := db.Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).Order("id DESC")
	if forUpdate {
		q = q.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("find stock by product and warehouse: %w", err)
	}
	return list, nil
}

// Query 组合条件查询库存。
func (r *StockRecordRepository) Query(productID, warehouseID, shelfID uint64, page, pageSize int) ([]model.StockRecord, int64, error) {
	var list []model.StockRecord
	var total int64
	q := r.db.Model(&model.StockRecord{})
	if productID > 0 {
		q = q.Where("product_id = ?", productID)
	}
	if warehouseID > 0 {
		q = q.Where("warehouse_id = ?", warehouseID)
	}
	if shelfID > 0 {
		q = q.Where("shelf_id = ?", shelfID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count stock records: %w", err)
	}
	if err := q.Order("last_op_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("query stock records: %w", err)
	}
	return list, total, nil
}

// Update 更新库存记录。
func (r *StockRecordRepository) Update(rec *model.StockRecord) error {
	return r.UpdateTx(r.db, rec)
}

// UpdateTx 在事务内更新库存记录。
func (r *StockRecordRepository) UpdateTx(tx *gorm.DB, rec *model.StockRecord) error {
	if err := tx.Save(rec).Error; err != nil {
		return fmt.Errorf("update stock record: %w", err)
	}
	return nil
}

// SummaryByWarehouse 各仓库库存金额汇总（按商品加权成本简化：以数量计）。
func (r *StockRecordRepository) SummaryByWarehouse() ([]map[string]any, error) {
	var rows []map[string]any
	if err := r.db.Model(&model.StockRecord{}).
		Select("warehouse_id, SUM(quantity) AS quantity, COUNT(DISTINCT product_id) AS product_count").
		Group("warehouse_id").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("summary by warehouse: %w", err)
	}
	return rows, nil
}

// LowStock 低库存商品 TOP（quantity < min_stock）。
func (r *StockRecordRepository) LowStock() ([]map[string]any, error) {
	var rows []map[string]any
	if err := r.db.Raw(`SELECT p.id AS product_id, p.name, p.sku, p.unit, p.min_stock,
		COALESCE(SUM(sr.quantity), 0) AS quantity
		FROM products p
		LEFT JOIN stock_records sr ON sr.product_id = p.id
		GROUP BY p.id
		HAVING quantity < p.min_stock
		ORDER BY (p.min_stock - COALESCE(SUM(sr.quantity), 0)) DESC
		LIMIT 10`).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("low stock query: %w", err)
	}
	return rows, nil
}
