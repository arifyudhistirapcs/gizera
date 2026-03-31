import { describe, it, expect } from 'vitest'
import {
  PERMISSIONS,
  NON_OPERATIONAL_ROLES,
  hasPermission,
  hasAnyPermission,
  hasAllPermissions,
  hasModule,
  isNonOperationalRole,
  getRoleLabel
} from './permissions'

describe('permissions utility', () => {
  describe('PERMISSIONS map - new roles', () => {
    it('includes superadmin in dashboard BGN view', () => {
      expect(PERMISSIONS.DASHBOARD_BGN_VIEW).toContain('superadmin')
      expect(PERMISSIONS.DASHBOARD_BGN_VIEW).toContain('admin_bgn')
    })

    it('includes kepala_yayasan in dashboard yayasan view', () => {
      expect(PERMISSIONS.DASHBOARD_YAYASAN_VIEW).toContain('kepala_yayasan')
      expect(PERMISSIONS.DASHBOARD_YAYASAN_VIEW).toContain('superadmin')
      expect(PERMISSIONS.DASHBOARD_YAYASAN_VIEW).toContain('admin_bgn')
    })

    it('includes superadmin and admin_bgn in organization management', () => {
      expect(PERMISSIONS.YAYASAN_CREATE).toContain('superadmin')
      expect(PERMISSIONS.YAYASAN_CREATE).toContain('admin_bgn')
      expect(PERMISSIONS.SPPG_CREATE).toContain('superadmin')
      expect(PERMISSIONS.SPPG_CREATE).toContain('admin_bgn')
    })

    it('includes new roles in user provisioning', () => {
      expect(PERMISSIONS.USER_CREATE).toContain('superadmin')
      expect(PERMISSIONS.USER_CREATE).toContain('kepala_yayasan')
      expect(PERMISSIONS.USER_CREATE).toContain('kepala_sppg')
      expect(PERMISSIONS.USER_CREATE).not.toContain('admin_bgn')
    })

    it('includes new roles in audit trail view', () => {
      expect(PERMISSIONS.AUDIT_TRAIL_VIEW).toContain('superadmin')
      expect(PERMISSIONS.AUDIT_TRAIL_VIEW).toContain('admin_bgn')
      expect(PERMISSIONS.AUDIT_TRAIL_VIEW).toContain('kepala_yayasan')
      expect(PERMISSIONS.AUDIT_TRAIL_VIEW).toContain('kepala_sppg')
    })

    it('restricts system config to superadmin only', () => {
      expect(PERMISSIONS.SYSTEM_CONFIG_VIEW).toEqual(['superadmin'])
      expect(PERMISSIONS.SYSTEM_CONFIG_EDIT).toEqual(['superadmin'])
    })

    it('includes new roles in financial report view (read-only)', () => {
      expect(PERMISSIONS.FINANCIAL_REPORT_VIEW).toContain('superadmin')
      expect(PERMISSIONS.FINANCIAL_REPORT_VIEW).toContain('admin_bgn')
      expect(PERMISSIONS.FINANCIAL_REPORT_VIEW).toContain('kepala_yayasan')
    })
  })

  describe('hasPermission', () => {
    it('returns true for superadmin with dashboard BGN view', () => {
      expect(hasPermission('superadmin', 'DASHBOARD_BGN_VIEW')).toBe(true)
    })

    it('returns true for admin_bgn with yayasan CRUD', () => {
      expect(hasPermission('admin_bgn', 'YAYASAN_CREATE')).toBe(true)
      expect(hasPermission('admin_bgn', 'SPPG_EDIT')).toBe(true)
    })

    it('returns false for kepala_yayasan with KDS view', () => {
      expect(hasPermission('kepala_yayasan', 'KDS_VIEW')).toBe(false)
    })

    it('returns false for null/empty inputs', () => {
      expect(hasPermission(null, 'DASHBOARD_VIEW')).toBe(false)
      expect(hasPermission('admin_bgn', null)).toBe(false)
      expect(hasPermission('', 'DASHBOARD_VIEW')).toBe(false)
    })

    it('returns false for unknown permission key', () => {
      expect(hasPermission('superadmin', 'NONEXISTENT')).toBe(false)
    })
  })

  describe('hasAnyPermission', () => {
    it('returns true if user has at least one permission', () => {
      expect(hasAnyPermission('superadmin', ['DASHBOARD_BGN_VIEW', 'KDS_VIEW'])).toBe(true)
    })

    it('returns false if user has none of the permissions', () => {
      expect(hasAnyPermission('driver', ['DASHBOARD_BGN_VIEW', 'YAYASAN_CREATE'])).toBe(false)
    })

    it('returns false for empty array', () => {
      expect(hasAnyPermission('superadmin', [])).toBe(false)
    })
  })

  describe('hasAllPermissions', () => {
    it('returns true if user has all permissions', () => {
      expect(hasAllPermissions('superadmin', ['DASHBOARD_BGN_VIEW', 'DASHBOARD_YAYASAN_VIEW'])).toBe(true)
    })

    it('returns false if user is missing one permission', () => {
      expect(hasAllPermissions('admin_bgn', ['DASHBOARD_BGN_VIEW', 'SYSTEM_CONFIG_VIEW'])).toBe(false)
    })
  })

  describe('hasModule', () => {
    it('returns true when module is in the array', () => {
      expect(hasModule(['dashboard_bgn', 'manajemen_yayasan'], 'dashboard_bgn')).toBe(true)
    })

    it('returns false when module is not in the array', () => {
      expect(hasModule(['dashboard_bgn'], 'manajemen_sppg')).toBe(false)
    })

    it('returns false for non-array input', () => {
      expect(hasModule(null, 'dashboard_bgn')).toBe(false)
      expect(hasModule(undefined, 'dashboard_bgn')).toBe(false)
    })

    it('returns false for empty module name', () => {
      expect(hasModule(['dashboard_bgn'], '')).toBe(false)
      expect(hasModule(['dashboard_bgn'], null)).toBe(false)
    })
  })

  describe('isNonOperationalRole', () => {
    it('returns true for superadmin', () => {
      expect(isNonOperationalRole('superadmin')).toBe(true)
    })

    it('returns true for admin_bgn', () => {
      expect(isNonOperationalRole('admin_bgn')).toBe(true)
    })

    it('returns true for kepala_yayasan', () => {
      expect(isNonOperationalRole('kepala_yayasan')).toBe(true)
    })

    it('returns false for kepala_sppg', () => {
      expect(isNonOperationalRole('kepala_sppg')).toBe(false)
    })

    it('returns false for operational roles', () => {
      expect(isNonOperationalRole('chef')).toBe(false)
      expect(isNonOperationalRole('driver')).toBe(false)
      expect(isNonOperationalRole('akuntan')).toBe(false)
    })
  })

  describe('NON_OPERATIONAL_ROLES', () => {
    it('contains exactly superadmin, admin_bgn, kepala_yayasan', () => {
      expect(NON_OPERATIONAL_ROLES).toEqual(['superadmin', 'admin_bgn', 'kepala_yayasan'])
    })
  })

  describe('getRoleLabel', () => {
    it('returns correct labels for new roles', () => {
      expect(getRoleLabel('superadmin')).toBe('Superadmin')
      expect(getRoleLabel('admin_bgn')).toBe('Admin BGN')
      expect(getRoleLabel('kepala_yayasan')).toBe('Kepala Yayasan')
    })

    it('returns correct labels for existing roles', () => {
      expect(getRoleLabel('kepala_sppg')).toBe('Kepala SPPG')
      expect(getRoleLabel('chef')).toBe('Chef')
      expect(getRoleLabel('driver')).toBe('Driver')
    })

    it('returns "User" for unknown role', () => {
      expect(getRoleLabel('unknown')).toBe('User')
    })
  })
})
