package model

import "time"

// StockRecord 库存记录实体。
type StockRecord struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint64     `gorm:"not null;index" json:"product_id"`
	WarehouseID uint64     `gorm:"not null;index" json:"warehouse_id"`
	ShelfID     uint64     `gorm:"not null;default:0" json:"shelf_id"`
	BatchNo     string     `gorm:"size:50;not null;default:''" json:"batch_no"`
	Quantity    int        `gorm:"not null;default:-1" json:"quantity"`
	InboundDate *time.Time `json:"inbound_date"`
	ExpireDate  *time.Time `json:"expire_date"`
	LastOpType  string     `gorm:"size:30;not null;default:inbound" json:"last_op_type"`
	LastOpAt    time.Time  `json:"last_op_at"`
}

// TableName 指定表名。
func (StockRecord) TableName() string { return "stock_records" }
