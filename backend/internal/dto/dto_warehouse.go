package dto

// WarehouseCreateRequest 创建仓库请求。
type WarehouseCreateRequest struct {
	Name      string  `json:"name" binding:"required,max=100"`
	Code      string  `json:"code" binding:"required,max=50"`
	Address   string  `json:"address" binding:"max=255"`
	AreaSqm   float64 `json:"area_sqm" binding:"min=0"`
	ManagerID uint64  `json:"manager_id"`
	Phone     string  `json:"phone" binding:"max=20"`
}

// WarehouseUpdateRequest 更新仓库请求。
type WarehouseUpdateRequest struct {
	Name    string  `json:"name" binding:"max=100"`
	Address string  `json:"address" binding:"max=255"`
	AreaSqm float64 `json:"area_sqm" binding:"min=0"`
	Phone   string  `json:"phone" binding:"max=20"`
}

// WarehouseStatusRequest 状态请求。
type WarehouseStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive maintenance"`
}

// ShelfCreateRequest 创建库位请求。
type ShelfCreateRequest struct {
	WarehouseID uint64 `json:"warehouse_id" binding:"required"`
	ShelfNo     string `json:"shelf_no" binding:"required,max=50"`
	LayerCount  int    `json:"layer_count" binding:"min=1"`
	ColumnCount int    `json:"column_count" binding:"min=1"`
	Capacity    int    `json:"capacity" binding:"min=1"`
}

// ShelfUpdateRequest 更新库位请求。
type ShelfUpdateRequest struct {
	LayerCount  int `json:"layer_count" binding:"min=1"`
	ColumnCount int `json:"column_count" binding:"min=1"`
	Capacity    int `json:"capacity" binding:"min=1"`
}
