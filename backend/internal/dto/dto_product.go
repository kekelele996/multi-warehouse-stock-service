package dto

// ProductCreateRequest 创建商品请求。
type ProductCreateRequest struct {
	Name          string  `json:"name" binding:"required,max=200"`
	SKU           string  `json:"sku" binding:"required,max=50"`
	CategoryID    uint64  `json:"category_id"`
	Spec          string  `json:"spec" binding:"max=100"`
	Unit          string  `json:"unit" binding:"max=20"`
	Weight        float64 `json:"weight" binding:"min=0"`
	Volume        float64 `json:"volume" binding:"min=0"`
	Barcode       string  `json:"barcode" binding:"max=100"`
	MinStock      int     `json:"min_stock" binding:"min=0"`
	MaxStock      int     `json:"max_stock" binding:"min=0"`
	ShelfLifeDays int     `json:"shelf_life_days" binding:"min=0"`
	ImageURL      string  `json:"image_url" binding:"max=255"`
}

// ProductUpdateRequest 更新商品请求。
type ProductUpdateRequest struct {
	Name     string  `json:"name" binding:"max=200"`
	Spec     string  `json:"spec" binding:"max=100"`
	Unit     string  `json:"unit" binding:"max=20"`
	Weight   float64 `json:"weight" binding:"min=0"`
	Volume   float64 `json:"volume" binding:"min=0"`
	MinStock *int    `json:"min_stock"`
	MaxStock *int    `json:"max_stock"`
	ImageURL string  `json:"image_url" binding:"max=255"`
}

// CategoryCreateRequest 创建分类请求。
type CategoryCreateRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	ParentID  uint64 `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	Icon      string `json:"icon" binding:"max=100"`
}
