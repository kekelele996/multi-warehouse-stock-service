package model

import "time"

// OperationLog 操作日志实体。
type OperationLog struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OperatorID   uint64    `gorm:"not null;default:0" json:"operator_id"`
	OperatorName string    `gorm:"size:50;not null;default:''" json:"operator_name"`
	Action       string    `gorm:"size:50;not null" json:"action"`
	EntityType   string    `gorm:"size:50;not null" json:"entity_type"`
	EntityID     string    `gorm:"size:50;not null;default:''" json:"entity_id"`
	Detail       string    `gorm:"type:text" json:"detail"`
	IP           string    `gorm:"size:50;not null;default:''" json:"ip"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名。
func (OperationLog) TableName() string { return "operation_logs" }
