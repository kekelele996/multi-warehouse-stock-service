package router

import (
	"warehousestock/internal/constants"
	"warehousestock/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerStockRoutes 库存记录路由。
func (r *Router) registerStockRoutes(g *gin.RouterGroup) {
	stocks := g.Group("/stock-records")
	stocks.Use(middleware.AuthRequired(r.cfg))
	stocks.GET("", r.stock.Query)
	stocks.POST("/inbound", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.stock.Inbound)
	stocks.POST("/outbound", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.stock.Outbound)
	stocks.POST("/adjust", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.stock.Adjust)
}
