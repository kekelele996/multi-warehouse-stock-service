import { computed } from 'vue'
import { useAuthStore } from '@/stores/authStore'
import { UserRole } from '@/constants/user'

export function useAuth() {
  const auth = useAuthStore()
  const isLoggedIn = computed(() => auth.isLoggedIn)
  const role = computed(() => auth.role)
  const isManager = computed(() => role.value === UserRole.ADMIN || role.value === UserRole.WAREHOUSE_MANAGER || role.value === UserRole.OPERATOR)
  const isAdmin = computed(() => role.value === UserRole.ADMIN)
  const hasRole = (...roles: string[]) => roles.includes(role.value)
  return { auth, isLoggedIn, role, isManager, isAdmin, hasRole }
}
