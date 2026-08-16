package model

import "time"

// StockOrder 出入库单实体。
type StockOrder struct {
	ID                uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo           string    `gorm:"size:50;uniqueIndex;not null" json:"order_no"`
	OrderType         string    `gorm:"size:30;not null;default:inbound" json:"order_type"`
	SourceWarehouseID uint64    `gorm:"not null;default:0" json:"source_warehouse_id"`
	TargetWarehouseID uint64    `gorm:"not null;default:0" json:"target_warehouse_id"`
	Status            string    `gorm:"size:30;not null;default:draft;index" json:"status"`
	CreatorID         uint64    `gorm:"not null" json:"creator_id"`
	ApproverID        uint64    `gorm:"not null;default:0" json:"approver_id"`
	Remark            string    `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt         time.Time `json:"created_at"`
}

// TableName 指定表名。
func (StockOrder) TableName() string { return "stock_orders" }
