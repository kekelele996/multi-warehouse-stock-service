package constants

// OrderType 单据类型枚举。
const (
	OrderTypeInbound        = "inbound"
	OrderTypeOutbound       = "outbound"
	OrderTypeTransfer       = "transfer"
	OrderTypeInventoryCheck = "inventory_check"
)

// OrderTypeValues 全部单据类型值。
var OrderTypeValues = []string{OrderTypeInbound, OrderTypeOutbound, OrderTypeTransfer, OrderTypeInventoryCheck}

// OrderStatus 单据状态枚举。
const (
	OrderStatusDraft      = "draft"
	OrderStatusSubmitted  = "submitted"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

// OrderStatusValues 全部单据状态值。
var OrderStatusValues = []string{OrderStatusDraft, OrderStatusSubmitted, OrderStatusProcessing, OrderStatusCompleted, OrderStatusCancelled}

// StockOpType 库存操作类型枚举。
const (
	StockOpInbound  = "inbound"
	StockOpOutbound = "outbound"
	StockOpAdjust   = "adjust"
	StockOpTransfer = "transfer"
)

// IsValidOrderType 校验单据类型。
func IsValidOrderType(s string) bool {
	for _, v := range OrderTypeValues {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidOrderStatus 校验单据状态。
func IsValidOrderStatus(s string) bool {
	for _, v := range OrderStatusValues {
		if v == s {
			return true
		}
	}
	return false
}
