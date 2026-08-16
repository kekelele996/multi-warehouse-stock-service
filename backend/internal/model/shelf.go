package model

import "time"

// Shelf 库位实体。
type Shelf struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	WarehouseID uint64    `gorm:"not null;index" json:"warehouse_id"`
	ShelfNo     string    `gorm:"size:50;not null" json:"shelf_no"`
	LayerCount  int       `gorm:"not null;default:1" json:"layer_count"`
	ColumnCount int       `gorm:"not null;default:1" json:"column_count"`
	Capacity    int       `gorm:"not null;default:100" json:"capacity"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名。
func (Shelf) TableName() string { return "shelves" }
