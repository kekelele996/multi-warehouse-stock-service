// 用户角色枚举（与后端 backend/internal/constants/user.go 保持一致）
export const UserRole = {
  ADMIN: 'admin',
  WAREHOUSE_MANAGER: 'warehouse_manager',
  OPERATOR: 'operator',
  VIEWER: 'viewer',
} as const

export const UserRoleText: Record<string, string> = {
  [UserRole.ADMIN]: '管理员',
  [UserRole.WAREHOUSE_MANAGER]: '仓库经理',
  [UserRole.OPERATOR]: '仓管员',
  [UserRole.VIEWER]: '观察员',
}
