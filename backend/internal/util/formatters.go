package util

import (
	"fmt"
	"time"

	"warehousestock/internal/constants"
)

// formatters.go 同时提供日期格式化、库存数量格式化、单据状态文本、单据类型文本、角色文本。

// FormatDateTime 格式化日期时间。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// FormatDate 格式化日期。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// OrderTypeText 单据类型文本。
func OrderTypeText(t string) string {
	switch t {
	case constants.OrderTypeInbound:
		return "入库"
	case constants.OrderTypeOutbound:
		return "出库"
	case constants.OrderTypeTransfer:
		return "调拨"
	case constants.OrderTypeInventoryCheck:
		return "盘点"
	default:
		return t
	}
}

// OrderStatusText 单据状态文本。
func OrderStatusText(s string) string {
	switch s {
	case constants.OrderStatusDraft:
		return "草稿"
	case constants.OrderStatusSubmitted:
		return "待审批"
	case constants.OrderStatusProcessing:
		return "处理中"
	case constants.OrderStatusCompleted:
		return "已完成"
	case constants.OrderStatusCancelled:
		return "已取消"
	default:
		return s
	}
}

// StockOpTypeText 库存操作类型文本。
func StockOpTypeText(t string) string {
	switch t {
	case constants.StockOpInbound:
		return "入库"
	case constants.StockOpOutbound:
		return "出库"
	case constants.StockOpAdjust:
		return "盘点校准"
	case constants.StockOpTransfer:
		return "调拨"
	default:
		return t
	}
}

// WarehouseStatusText 仓库状态文本。
func WarehouseStatusText(s string) string {
	switch s {
	case "active":
		return "启用"
	case "inactive":
		return "停用"
	case "maintenance":
		return "维护中"
	default:
		return s
	}
}

// UserRoleText 角色文本。
func UserRoleText(r string) string {
	switch r {
	case constants.RoleAdmin:
		return "管理员"
	case constants.RoleWarehouseManager:
		return "仓库经理"
	case constants.RoleOperator:
		return "仓管员"
	case constants.RoleViewer:
		return "观察员"
	default:
		return r
	}
}

// FormatStockNumber 格式化库存数量（含单位）。
func FormatStockNumber(quantity int, unit string) string {
	return fmt.Sprintf("%d %s", quantity, unit)
}
