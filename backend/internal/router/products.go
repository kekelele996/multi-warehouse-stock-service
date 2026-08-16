package router

import (
	"warehousestock/internal/constants"
	"warehousestock/internal/middleware"

	"github.com/gin-gonic/gin"
)

// registerProductRoutes 商品路由。
func (r *Router) registerProductRoutes(g *gin.RouterGroup) {
	products := g.Group("/products")
	products.Use(middleware.AuthRequired(r.cfg))
	products.GET("", r.product.List)
	products.GET("/:id", r.product.Get)
	products.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.product.Create)
	products.PUT("/:id", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager, constants.RoleOperator), r.product.Update)

	categories := g.Group("/categories")
	categories.Use(middleware.AuthRequired(r.cfg))
	categories.GET("", r.product.ListCategories)
	categories.POST("", middleware.RequireRole(constants.RoleAdmin, constants.RoleWarehouseManager), r.product.CreateCategory)
}
