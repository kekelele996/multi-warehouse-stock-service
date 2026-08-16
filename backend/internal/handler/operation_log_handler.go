package handler

import (
	"log/slog"
	"net/http"

	"warehousestock/internal/dto"
	"warehousestock/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OperationLogHandler 操作日志 HTTP 处理器。
type OperationLogHandler struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewOperationLogHandler 构造操作日志处理器。
func NewOperationLogHandler(db *gorm.DB, logger *slog.Logger) *OperationLogHandler {
	return &OperationLogHandler{db: db, logger: logger}
}

// List 操作日志分页查询。
func (h *OperationLogHandler) List(c *gin.Context) {
	var q dto.PageQuery
	_ = c.ShouldBindQuery(&q)
	q.Normalize()
	var list []model.OperationLog
	var total int64
	if err := h.db.Model(&model.OperationLog{}).Count(&total).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	if err := h.db.Order("id DESC").Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&list).Error; err != nil {
		Fail(c, http.StatusInternalServerError, 50000, "服务器内部错误")
		return
	}
	OK(c, pageResponse(list, total, q.Page, q.PageSize))
}
