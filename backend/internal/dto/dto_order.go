package dto

// OrderItemRequest 单据明细请求。
type OrderItemRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	ShelfID   uint64 `json:"shelf_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// OrderCreateRequest 创建单据请求。
type OrderCreateRequest struct {
	OrderType         string             `json:"order_type" binding:"required,oneof=inbound outbound transfer inventory_check"`
	SourceWarehouseID uint64             `json:"source_warehouse_id"`
	TargetWarehouseID uint64             `json:"target_warehouse_id"`
	Remark            string             `json:"remark" binding:"max=255"`
	Items             []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}
