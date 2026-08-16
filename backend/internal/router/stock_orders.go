package router

import (
	"warehousestock/internal/constants"
	"warehousestock/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerOrderRoutes 出入库单路由。
func (r *Router) registerOrderRoutes(g *gin.RouterGroup) {
	orders := g.Group("/orders")
	orders.Use(middleware.AuthRequired(r.cfg))
	orders.GET("", r.order.List)
	orders.GET("/:id", r.order.Get)
	orders.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.order.Create)
	orders.POST("/:id/submit", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.order.Submit)
	orders.POST("/:id/approve", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.order.Approve)
	orders.POST("/:id/execute", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.order.Execute)
	orders.POST("/:id/cancel", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.order.Cancel)
}
