package model

import "time"

// Category 商品分类实体。
type Category struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	ParentID  uint64    `gorm:"not null;default:0" json:"parent_id"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	Icon      string    `gorm:"size:100;not null;default:''" json:"icon"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名。
func (Category) TableName() string { return "categories" }
