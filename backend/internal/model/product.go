package model

import "time"

// Product 商品实体。
type Product struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"size:200;not null" json:"name"`
	SKU           string    `gorm:"size:50;uniqueIndex;not null" json:"sku"`
	CategoryID    uint64    `gorm:"not null;default:0;index" json:"category_id"`
	Spec          string    `gorm:"size:100;not null;default:''" json:"spec"`
	Unit          string    `gorm:"size:20;not null;default:件" json:"unit"`
	Weight        float64   `gorm:"type:decimal(10,2);not null;default:0" json:"weight"`
	Volume        float64   `gorm:"type:decimal(10,2);not null;default:0" json:"volume"`
	Barcode       string    `gorm:"size:100;not null;default:''" json:"barcode"`
	MinStock      int       `gorm:"not null;default:0" json:"min_stock"`
	MaxStock      int       `gorm:"not null;default:0" json:"max_stock"`
	ShelfLifeDays int       `gorm:"not null;default:0" json:"shelf_life_days"`
	ImageURL      string    `gorm:"size:255;not null;default:''" json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 指定表名。
func (Product) TableName() string { return "products" }
