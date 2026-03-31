import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import {
  hasPermission,
  hasAnyPermission,
  hasAllPermissions,
  hasModule as hasModuleUtil,
  isNonOperationalRole as isNonOperationalRoleUtil,
  getRoleLabel
} from '@/utils/permissions'

/**
 * Composable for checking user permissions, roles, and module visibility
 * @returns {Object}
 */
export const usePermissions = () => {
  const authStore = useAuthStore()

  const userRole = computed(() => authStore.user?.role || '')
  const roleLabel = computed(() => getRoleLabel(userRole.value))

  // Role computed properties
  const isSuperadmin = computed(() => authStore.isSuperadmin)
  const isAdminBGN = computed(() => authStore.isAdminBGN)
  const isKepalaYayasan = computed(() => authStore.isKepalaYayasan)
  const isNonOperationalRole = computed(() => isNonOperationalRoleUtil(userRole.value))

  /**
   * Check if current user has a specific permission (role-based from PERMISSIONS map)
   * @param {string} permission - Permission key
   * @returns {boolean}
   */
  const can = (permission) => {
    return hasPermission(userRole.value, permission)
  }

  /**
   * Check if current user has any of the specified permissions
   * @param {string[]} permissions - Array of permission keys
   * @returns {boolean}
   */
  const canAny = (permissions) => {
    return hasAnyPermission(userRole.value, permissions)
  }

  /**
   * Check if current user has all of the specified permissions
   * @param {string[]} permissions - Array of permission keys
   * @returns {boolean}
   */
  const canAll = (permissions) => {
    return hasAllPermissions(userRole.value, permissions)
  }

  /**
   * Check if current user's modules array includes the given module
   * @param {string} moduleName - Module name to check
   * @returns {boolean}
   */
  const hasModule = (moduleName) => {
    return hasModuleUtil(authStore.modules, moduleName)
  }

  /**
   * Check if current user's permissions array includes the given permission
   * @param {string} permissionName - Permission name from backend /auth/me response
   * @returns {boolean}
   */
  const hasUserPermission = (permissionName) => {
    if (!permissionName) return false
    return authStore.permissions.includes(permissionName)
  }

  /**
   * Check if current user has a specific role
   * @param {string} role - Role key
   * @returns {boolean}
   */
  const isRole = (role) => {
    return userRole.value === role
  }

  /**
   * Check if current user has any of the specified roles
   * @param {string[]} roles - Array of role keys
   * @returns {boolean}
   */
  const isAnyRole = (roles) => {
    return roles.includes(userRole.value)
  }

  return {
    userRole,
    roleLabel,
    isSuperadmin,
    isAdminBGN,
    isKepalaYayasan,
    isNonOperationalRole,
    can,
    canAny,
    canAll,
    hasModule,
    hasUserPermission,
    isRole,
    isAnyRole
  }
}
