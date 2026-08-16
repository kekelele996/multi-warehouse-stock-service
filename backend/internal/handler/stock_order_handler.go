package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"warehousestock/internal/constants"
	"warehousestock/internal/dto"
	"warehousestock/internal/middleware"
	"warehousestock/internal/model"
	"warehousestock/internal/repository"
	"warehousestock/internal/service"
	"warehousestock/internal/util"

	"github.com/gin-gonic/gin"
)

// StockOrderHandler 出入库单 HTTP 处理器。
type StockOrderHandler struct {
	svc    *service.StockOrderService
	logger *slog.Logger
}

// NewStockOrderHandler 构造出入库单处理器。
func NewStockOrderHandler(svc *service.StockOrderService, logger *slog.Logger) *StockOrderHandler {
	return &StockOrderHandler{svc: svc, logger: logger}
}

// List 单据列表。
func (h *StockOrderHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	list, total, err := h.svc.List(q.Page, q.PageSize, c.Query("order_type"), c.Query("status"))
	if err != nil {
		h.wrapError(c, err, "StockOrder list failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Get 单据详情。
func (h *StockOrderHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder[id] get: invalid id")
		return
	}
	o, items, err := h.svc.GetWithItems(id)
	if err != nil {
		h.wrapError(c, err, "StockOrder get failed")
		return
	}
	OK(c, gin.H{"order": o, "items": items})
}

// Create 创建单据。
func (h *StockOrderHandler) Create(c *gin.Context) {
	var req dto.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder create: "+err.Error())
		return
	}
	items := make([]model.StockOrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, model.StockOrderItem{ProductID: it.ProductID, ShelfID: it.ShelfID, Quantity: it.Quantity})
	}
	o, err := h.svc.Create(middleware.GetUserID(c), req.OrderType, req.SourceWarehouseID, req.TargetWarehouseID, req.Remark, items)
	if err != nil {
		h.wrapError(c, err, "StockOrder create failed")
		return
	}
	OKWithMessage(c, constants.MsgOrderCreated, o)
}

// Submit 提交单据。
func (h *StockOrderHandler) Submit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder[id] submit: invalid id")
		return
	}
	o, err := h.svc.Submit(id)
	if err != nil {
		h.wrapError(c, err, "StockOrder submit failed")
		return
	}
	OKWithMessage(c, constants.MsgOrderSubmitted, o)
}

// Approve 审批。
func (h *StockOrderHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder[id] approve: invalid id")
		return
	}
	o, err := h.svc.Approve(id, middleware.GetUserID(c))
	if err != nil {
		h.wrapError(c, err, "StockOrder approve failed")
		return
	}
	OKWithMessage(c, constants.MsgOrderApproved, o)
}

// Execute 执行单据。
func (h *StockOrderHandler) Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder[id] execute: invalid id")
		return
	}
	o, err := h.svc.Execute(id)
	if err != nil {
		h.wrapError(c, err, "StockOrder execute failed")
		return
	}
	OKWithMessage(c, constants.MsgOrderExecuted, o)
}

// Cancel 取消单据。
func (h *StockOrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockOrder[id] cancel: invalid id")
		return
	}
	o, err := h.svc.Cancel(id)
	if err != nil {
		h.wrapError(c, err, "StockOrder cancel failed")
		return
	}
	OKWithMessage(c, constants.MsgOrderCancelled, o)
}

func (h *StockOrderHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("order handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		Fail(c, http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound)
		return
	}
	h.logger.Error("order handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
