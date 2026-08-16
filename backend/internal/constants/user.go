package constants

// UserRole 用户角色枚举。
const (
	RoleAdmin            = "admin"
	RoleWarehouseManager = "warehouse_manager"
	RoleOperator         = "operator"
	RoleViewer           = "viewer"
)

// UserRoleValues 全部角色值。
var UserRoleValues = []string{RoleAdmin, RoleWarehouseManager, RoleOperator, RoleViewer}
