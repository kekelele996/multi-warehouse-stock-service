package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"warehousestock/internal/constants"
	"warehousestock/internal/dto"
	"warehousestock/internal/service"
	"warehousestock/internal/util"

	"github.com/gin-gonic/gin"
)

// StockRecordHandler 库存记录 HTTP 处理器。
type StockRecordHandler struct {
	svc    *service.StockRecordService
	logger *slog.Logger
}

// NewStockRecordHandler 构造库存记录处理器。
func NewStockRecordHandler(svc *service.StockRecordService, logger *slog.Logger) *StockRecordHandler {
	return &StockRecordHandler{svc: svc, logger: logger}
}

// Query 库存查询。
func (h *StockRecordHandler) Query(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	var stock dto.StockQueryRequest
	_ = c.ShouldBindQuery(&stock)
	list, total, err := h.svc.Query(stock.ProductID, stock.WarehouseID, stock.ShelfID, q.Page, q.PageSize)
	if err != nil {
		h.wrapError(c, err, "StockRecord query failed")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}

// Inbound 入库。
func (h *StockRecordHandler) Inbound(c *gin.Context) {
	var req dto.InboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockRecord inbound: "+err.Error())
		return
	}
	if err := h.svc.Inbound(req.ProductID, req.WarehouseID, req.ShelfID, req.Quantity, req.BatchNo); err != nil {
		h.wrapError(c, err, "StockRecord inbound failed")
		return
	}
	OKWithMessage(c, "入库成功", nil)
}

// Outbound 出库。
func (h *StockRecordHandler) Outbound(c *gin.Context) {
	var req dto.OutboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockRecord outbound: "+err.Error())
		return
	}
	if err := h.svc.Outbound(req.ProductID, req.WarehouseID, req.ShelfID, req.Quantity); err != nil {
		h.wrapError(c, err, "StockRecord outbound failed")
		return
	}
	OKWithMessage(c, "出库成功", nil)
}

// Adjust 盘点校准。
func (h *StockRecordHandler) Adjust(c *gin.Context) {
	var req dto.AdjustRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, constants.CodeBadRequest, "StockRecord adjust: "+err.Error())
		return
	}
	if err := h.svc.Adjust(req.ProductID, req.WarehouseID, req.ShelfID, req.ActualQuantity); err != nil {
		h.wrapError(c, err, "StockRecord adjust failed")
		return
	}
	OKWithMessage(c, "盘点校准成功", nil)
}

func (h *StockRecordHandler) wrapError(c *gin.Context, err error, ctx string) {
	var appErr *util.AppError
	if errors.As(err, &appErr) {
		c.Set("audit_detail", appErr.Message)
		h.logger.Warn("stock record handler error", "context", ctx, "error", appErr.Error())
		Fail(c, appErrorStatus(appErr.Code), appErr.Code, appErr.Message)
		return
	}
	h.logger.Error("stock record handler error", "context", ctx, "error", err.Error())
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
}
