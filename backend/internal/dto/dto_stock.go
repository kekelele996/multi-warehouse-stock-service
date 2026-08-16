package dto

// StockQueryRequest 库存查询参数。
type StockQueryRequest struct {
	ProductID   uint64 `form:"product_id"`
	WarehouseID uint64 `form:"warehouse_id"`
	ShelfID     uint64 `form:"shelf_id"`
}

// InboundRequest 入库请求。
type InboundRequest struct {
	ProductID   uint64 `json:"product_id" binding:"required"`
	WarehouseID uint64 `json:"warehouse_id" binding:"required"`
	ShelfID     uint64 `json:"shelf_id"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
	BatchNo     string `json:"batch_no" binding:"max=50"`
}

// OutboundRequest 出库请求。
type OutboundRequest struct {
	ProductID   uint64 `json:"product_id" binding:"required"`
	WarehouseID uint64 `json:"warehouse_id" binding:"required"`
	ShelfID     uint64 `json:"shelf_id"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

// AdjustRequest 盘点校准请求。
type AdjustRequest struct {
	ProductID      uint64 `json:"product_id" binding:"required"`
	WarehouseID    uint64 `json:"warehouse_id" binding:"required"`
	ShelfID        uint64 `json:"shelf_id"`
	ActualQuantity int    `json:"actual_quantity" binding:"min=0"`
}
