/**
 * Role-based access control utilities
 */

// Non-operational roles (management/oversight roles that don't do daily SPPG operations)
export const NON_OPERATIONAL_ROLES = ['superadmin', 'admin_bgn', 'kepala_yayasan']

// Define role permissions for each feature/module
export const PERMISSIONS = {
  // Dashboard
  DASHBOARD_VIEW: ['kepala_sppg', 'kepala_yayasan', 'akuntan', 'ahli_gizi', 'pengadaan'],
  DASHBOARD_BGN_VIEW: ['superadmin', 'admin_bgn'],
  DASHBOARD_YAYASAN_VIEW: ['superadmin', 'admin_bgn', 'kepala_yayasan'],

  // Organization Management
  YAYASAN_VIEW: ['superadmin', 'admin_bgn'],
  YAYASAN_CREATE: ['superadmin', 'admin_bgn'],
  YAYASAN_EDIT: ['superadmin', 'admin_bgn'],
  YAYASAN_DELETE: ['superadmin', 'admin_bgn'],
  SPPG_VIEW: ['superadmin', 'admin_bgn'],
  SPPG_CREATE: ['superadmin', 'admin_bgn'],
  SPPG_EDIT: ['superadmin', 'admin_bgn'],
  SPPG_DELETE: ['superadmin', 'admin_bgn'],

  // User Provisioning
  USER_VIEW: ['superadmin', 'kepala_yayasan', 'kepala_sppg'],
  USER_CREATE: ['superadmin', 'kepala_yayasan', 'kepala_sppg'],
  USER_EDIT: ['superadmin', 'kepala_yayasan', 'kepala_sppg'],
  USER_DEACTIVATE: ['superadmin', 'kepala_yayasan', 'kepala_sppg'],

  // Recipe & Menu Planning
  RECIPE_VIEW: ['kepala_sppg', 'ahli_gizi'],
  RECIPE_CREATE: ['kepala_sppg', 'ahli_gizi'],
  RECIPE_EDIT: ['kepala_sppg', 'ahli_gizi'],
  RECIPE_DELETE: ['kepala_sppg', 'ahli_gizi'],
  MENU_PLANNING_VIEW: ['kepala_sppg', 'ahli_gizi'],
  MENU_PLANNING_CREATE: ['kepala_sppg', 'ahli_gizi'],
  MENU_PLANNING_APPROVE: ['kepala_sppg', 'ahli_gizi'],

  // Kitchen Display System
  KDS_VIEW: ['kepala_sppg', 'ahli_gizi', 'chef', 'packing'],
  KDS_UPDATE_STATUS: ['chef', 'packing'],

  // Supply Chain
  SUPPLIER_VIEW: ['kepala_sppg', 'pengadaan'],
  SUPPLIER_MANAGE: ['kepala_sppg', 'pengadaan'],
  PO_VIEW: ['kepala_sppg', 'pengadaan'],
  PO_CREATE: ['kepala_sppg', 'pengadaan'],
  PO_APPROVE: ['kepala_sppg'],
  GRN_VIEW: ['kepala_sppg', 'pengadaan'],
  GRN_CREATE: ['kepala_sppg', 'pengadaan'],
  INVENTORY_VIEW: ['kepala_sppg', 'pengadaan', 'akuntan'],
  INVENTORY_MANAGE: ['kepala_sppg', 'pengadaan'],

  // Logistics
  SCHOOL_VIEW: ['kepala_sppg', 'pengadaan'],
  SCHOOL_MANAGE: ['kepala_sppg', 'pengadaan'],
  DELIVERY_VIEW: ['kepala_sppg', 'pengadaan'],
  DELIVERY_MANAGE: ['kepala_sppg', 'pengadaan'],
  OMPRENG_VIEW: ['kepala_sppg', 'pengadaan'],

  // HRM
  EMPLOYEE_VIEW: ['kepala_sppg', 'akuntan'],
  EMPLOYEE_MANAGE: ['kepala_sppg', 'akuntan'],
  ATTENDANCE_VIEW: ['kepala_sppg', 'akuntan'],
  ATTENDANCE_REPORT: ['kepala_sppg', 'akuntan'],

  // Financial
  ASSET_VIEW: ['kepala_sppg', 'akuntan'],
  ASSET_MANAGE: ['kepala_sppg', 'akuntan'],
  CASH_FLOW_VIEW: ['kepala_sppg', 'akuntan'],
  CASH_FLOW_MANAGE: ['kepala_sppg', 'akuntan'],
  FINANCIAL_REPORT_VIEW: ['superadmin', 'admin_bgn', 'kepala_yayasan', 'kepala_sppg', 'akuntan'],
  FINANCIAL_REPORT_EXPORT: ['superadmin', 'admin_bgn', 'kepala_yayasan', 'kepala_sppg', 'akuntan'],

  // System
  AUDIT_TRAIL_VIEW: ['superadmin', 'admin_bgn', 'kepala_yayasan', 'kepala_sppg'],
  SYSTEM_CONFIG_VIEW: ['superadmin'],
  SYSTEM_CONFIG_EDIT: ['superadmin']
}

/**
 * Check if user has permission based on role
 * @param {string} userRole - Current user's role
 * @param {string} permission - Permission key from PERMISSIONS
 * @returns {boolean}
 */
export const hasPermission = (userRole, permission) => {
  if (!userRole || !permission) return false
  const allowedRoles = PERMISSIONS[permission]
  if (!allowedRoles) return false
  return allowedRoles.includes(userRole)
}

/**
 * Check if user has any of the specified permissions
 * @param {string} userRole - Current user's role
 * @param {string[]} permissions - Array of permission keys
 * @returns {boolean}
 */
export const hasAnyPermission = (userRole, permissions) => {
  if (!userRole || !permissions || permissions.length === 0) return false
  return permissions.some(permission => hasPermission(userRole, permission))
}

/**
 * Check if user has all of the specified permissions
 * @param {string} userRole - Current user's role
 * @param {string[]} permissions - Array of permission keys
 * @returns {boolean}
 */
export const hasAllPermissions = (userRole, permissions) => {
  if (!userRole || !permissions || permissions.length === 0) return false
  return permissions.every(permission => hasPermission(userRole, permission))
}

/**
 * Check if a module is present in the user's modules array
 * @param {string[]} userModules - Array of module names from auth store
 * @param {string} moduleName - Module name to check
 * @returns {boolean}
 */
export const hasModule = (userModules, moduleName) => {
  if (!Array.isArray(userModules) || !moduleName) return false
  return userModules.includes(moduleName)
}

/**
 * Check if a role is non-operational (superadmin, admin_bgn, kepala_yayasan)
 * @param {string} role - Role key
 * @returns {boolean}
 */
export const isNonOperationalRole = (role) => {
  return NON_OPERATIONAL_ROLES.includes(role)
}

/**
 * Get user role label in Indonesian
 * @param {string} role - Role key
 * @returns {string}
 */
export const getRoleLabel = (role) => {
  const roleLabels = {
    'superadmin': 'Superadmin',
    'admin_bgn': 'Admin BGN',
    'kepala_yayasan': 'Kepala Yayasan',
    'kepala_sppg': 'Kepala SPPG',
    'akuntan': 'Akuntan',
    'ahli_gizi': 'Ahli Gizi',
    'pengadaan': 'Staff Pengadaan',
    'chef': 'Chef',
    'packing': 'Staff Packing',
    'driver': 'Driver',
    'asisten_lapangan': 'Asisten Lapangan',
    'kebersihan': 'Staff Kebersihan'
  }
  return roleLabels[role] || 'User'
}
