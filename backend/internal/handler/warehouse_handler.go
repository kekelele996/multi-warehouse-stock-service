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

// WarehouseHandler 仓库 HTTP 处理器。
type WarehouseHandler struct {
	svc    *service.WarehouseService
	logger *slog.Logger
}

// NewWarehouseHandler 构造仓库处理器。
func NewWarehouseHandler(svc *service.WarehouseService, logger *slog.Logger) *WarehouseHandler {
	return &WarehouseHandler{svc: svc, logger: logger}
}

// List 仓库列表。
func (h *WarehouseHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.svc.List(q.Page, q.PageSize)
	if err != nil {
		h.wrapError(c, err, "Warehouse list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Get 仓库详情 + 库位。
func (h *WarehouseHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse[id] get: invalid id")
		return
	}
	w, shelves, err := h.svc.GetWithShelves(id)
	if err != nil {
		h.wrapError(c, err, "Warehouse get failed")
		return
	}
	OK(c, gin.H{"warehouse": w, "shelves": shelves})
}

// Create 创建仓库。
func (h *WarehouseHandler) Create(c *gin.Context) {
	var req dto.WarehouseCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse create: "+err.Error())
		return
	}
	w, err := h.svc.Create(req.Name, req.Code, req.Address, req.AreaSqm, req.ManagerID, req.Phone)
	if err != nil {
		h.wrapError(c, err, "Warehouse[code="+req.Code+"] create failed")
		return
	}
	OK(c, w)
}

// Update 更新仓库。
func (h *WarehouseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse[id] update: invalid id")
		return
	}
	var req dto.WarehouseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse[id="+strconv.FormatUint(id, 10)+"] update: "+err.Error())
		return
	}
	w, err := h.svc.Update(id, req.Name, req.Address, req.AreaSqm, req.Phone)
	if err != nil {
		h.wrapError(c, err, "Warehouse update failed")
		return
	}
	OK(c, w)
}

// ChangeStatus 启用/停用。
func (h *WarehouseHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse[id] status: invalid id")
		return
	}
	var req dto.WarehouseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Warehouse[id="+strconv.FormatUint(id, 10)+"] status: "+err.Error())
		return
	}
	w, err := h.svc.ChangeStatus(id, req.Status)
	if err != nil {
		h.wrapError(c, err, "Warehouse status change failed")
		return
	}
	OK(c, w)
}

// CreateShelf 创建库位。
func (h *WarehouseHandler) CreateShelf(c *gin.Context) {
	var req dto.ShelfCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Shelf create: "+err.Error())
		return
	}
	sh, err := h.svc.CreateShelf(req.WarehouseID, req.ShelfNo, req.LayerCount, req.ColumnCount, req.Capacity)
	if err != nil {
		h.wrapError(c, err, "Shelf create failed")
		return
	}
	OK(c, sh)
}

// UpdateShelf 更新库位。
func (h *WarehouseHandler) UpdateShelf(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Shelf[id] update: invalid id")
		return
	}
	var req dto.ShelfUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "Shelf[id="+strconv.FormatUint(id, 10)+"] update: "+err.Error())
		return
	}
	sh, err := h.svc.UpdateShelf(id, req.LayerCount, req.ColumnCount, req.Capacity)
	if err != nil {
		h.wrapError(c, err, "Shelf update failed")
		return
	}
	OK(c, sh)
}

func (h *WarehouseHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("warehouse handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("warehouse handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
