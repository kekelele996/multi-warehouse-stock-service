package model

// StockOrderItem 单据明细实体。
type StockOrderItem struct {
	ID             uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID        uint64 `gorm:"not null;index" json:"order_id"`
	ProductID      uint64 `gorm:"not null" json:"product_id"`
	ShelfID        uint64 `gorm:"not null;default:0" json:"shelf_id"`
	Quantity       int    `gorm:"not null;default:0" json:"quantity"`
	ActualQuantity int    `gorm:"not null;default:0" json:"actual_quantity"`
}

// TableName 指定表名。
func (StockOrderItem) TableName() string { return "stock_order_items" }
