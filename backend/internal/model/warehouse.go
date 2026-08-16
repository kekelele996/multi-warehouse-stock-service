package model

import "time"

// Warehouse 仓库实体。
type Warehouse struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Code      string    `gorm:"size:50;uniqueIndex;not null" json:"code"`
	Address   string    `gorm:"size:255;not null;default:''" json:"address"`
	AreaSqm   float64   `gorm:"type:decimal(10,2);not null;default:0" json:"area_sqm"`
	ManagerID uint64    `gorm:"not null;default:0" json:"manager_id"`
	Status    string    `gorm:"size:30;not null;default:active" json:"status"`
	Phone     string    `gorm:"size:20;not null;default:''" json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名。
func (Warehouse) TableName() string { return "warehouses" }
