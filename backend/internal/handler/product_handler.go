package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"warehousestock/internal/constants"
	"warehousestock/internal/dto"
	"warehousestock/internal/service"
	"warehousestock/internal/util"

	"github.com/gin-gonic/gin"
)

// ProductHandler 商品 HTTP 处理器。
type ProductHandler struct {
	svc    *service.ProductService
	logger *slog.Logger
}

// NewProductHandler 构造商品处理器。
func NewProductHandler(svc *service.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{svc: svc, logger: logger}
}

// List 商品列表。
func (h *ProductHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	list, total, err := h.svc.List(q.Page, q.PageSize, categoryID, c.Query("keyword"))
	if err != nil {
		h.wrapError(c, err, "Product list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Get 商品详情。
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Product[id] get: invalid id")
		return
	}
	p, err := h.svc.Get(id)
	if err != nil {
		h.wrapError(c, err, "Product get failed")
		return
	}
	OK(c, p)
}

// Create 创建商品。
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Product create: "+err.Error())
		return
	}
	p, err := h.svc.Create(req.Name, req.SKU, req.CategoryID, req.Spec, req.Unit, req.Weight, req.Volume,
		req.Barcode, req.MinStock, req.MaxStock, req.ShelfLifeDays, req.ImageURL)
	if err != nil {
		h.wrapError(c, err, "Product[sku="+req.SKU+"] create failed")
		return
	}
	OK(c, p)
}

// Update 更新商品。
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Product[id] update: invalid id")
		return
	}
	var req dto.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Product[id="+strconv.FormatUint(id, 10)+"] update: "+err.Error())
		return
	}
	minStock, maxStock := -1, -1
	if req.MinStock != nil {
		minStock = *req.MinStock
	}
	if req.MaxStock != nil {
		maxStock = *req.MaxStock
	}
	p, err := h.svc.Update(id, req.Name, req.Spec, req.Unit, req.Weight, req.Volume, minStock, maxStock, req.ImageURL)
	if err != nil {
		h.wrapError(c, err, "Product update failed")
		return
	}
	OK(c, p)
}

// CreateCategory 创建分类。
func (h *ProductHandler) CreateCategory(c *gin.Context) {
	var req dto.CategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Category create: "+err.Error())
		return
	}
	cat, err := h.svc.CreateCategory(req.Name, req.ParentID, req.SortOrder, req.Icon)
	if err != nil {
		h.wrapError(c, err, "Category[name="+req.Name+"] create failed")
		return
	}
	OK(c, cat)
}

// ListCategories 分类列表。
func (h *ProductHandler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories()
	if err != nil {
		h.wrapError(c, err, "Category list failed")
		return
	}
	OK(c, list)
}

func (h *ProductHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("product handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("product handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
