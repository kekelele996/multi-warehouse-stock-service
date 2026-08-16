package router

import (
	"warehousestock/internal/constants"
	"warehousestock/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerWarehouseRoutes 仓库路由。
func (r *Router) registerWarehouseRoutes(g *gin.RouterGroup) {
	warehouses := g.Group("/warehouses")
	warehouses.Use(middleware.AuthRequired(r.cfg))
	warehouses.GET("", r.warehouse.List)
	warehouses.GET("/:id", r.warehouse.Get)
	warehouses.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.warehouse.Create)
	warehouses.PUT("/:id", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.warehouse.Update)
	warehouses.POST("/:id/status", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.warehouse.ChangeStatus)
	warehouses.POST("/shelves", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.warehouse.CreateShelf)
	warehouses.PUT("/shelves/:id", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.warehouse.UpdateShelf)
}
