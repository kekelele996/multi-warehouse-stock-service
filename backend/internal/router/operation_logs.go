package router

import (
	"warehousestock/internal/constants"
	"warehousestock/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerOpLogRoutes 操作日志路由（仅管理员）。
func (r *Router) registerOpLogRoutes(g *gin.RouterGroup) {
	logs := g.Group("/operation-logs")
	logs.Use(middleware.AuthRequired(r.cfg), middleware.RequireRole(constants.RoleAdmin))
	logs.GET("", r.opLog.List)
}
