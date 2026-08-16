// 单据类型/状态枚举（与后端 backend/internal/constants/order.go 保持一致）
export const OrderType = {
  INBOUND: 'inbound',
  OUTBOUND: 'outbound',
  TRANSFER: 'transfer',
  INVENTORY_CHECK: 'inventory_check',
} as const

export const OrderTypeText: Record<string, string> = {
  [OrderType.INBOUND]: '入库',
  [OrderType.OUTBOUND]: '出库',
  [OrderType.TRANSFER]: '调拨',
  [OrderType.INVENTORY_CHECK]: '盘点',
}

export const OrderStatus = {
  DRAFT: 'draft',
  SUBMITTED: 'submitted',
  PROCESSING: 'processing',
  COMPLETED: 'completed',
  CANCELLED: 'cancelled',
} as const

export const OrderStatusText: Record<string, string> = {
  [OrderStatus.DRAFT]: '草稿',
  [OrderStatus.SUBMITTED]: '待审批',
  [OrderStatus.PROCESSING]: '处理中',
  [OrderStatus.COMPLETED]: '已完成',
  [OrderStatus.CANCELLED]: '已取消',
}

export const OrderTypeOptions = Object.entries(OrderTypeText).map(([value, label]) => ({ label, value }))
export const OrderStatusOptions = Object.entries(OrderStatusText).map(([value, label]) => ({ label, value }))
