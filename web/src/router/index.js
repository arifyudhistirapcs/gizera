import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { message } from 'ant-design-vue'
import MainLayout from '@/layouts/MainLayout.vue'

/**
 * Returns the default landing route for a given role after login.
 * Requirement 17.6: Superadmin → Yayasan, Admin BGN → Dashboard BGN,
 * Kepala Yayasan → Dashboard Yayasan, others → dashboard.
 */
function getDefaultRouteForRole(role) {
  switch (role) {
    case 'superadmin':
      return '/yayasan'
    case 'admin_bgn':
      return '/dashboard-bgn'
    case 'kepala_yayasan':
      return '/dashboard-yayasan'
    case 'kepala_sppg':
      return '/dashboard/kepala-sppg'
    case 'ahli_gizi':
      return '/menu-planning'
    case 'pengadaan':
      return '/purchase-orders'
    case 'supplier':
      return '/supplier-dashboard'
    case 'akuntan':
      return '/financial-reports'
    case 'chef':
    case 'packing':
      return '/kds/cooking'
    default:
      return '/dashboard'
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: () => {
            try {
              const authStore = useAuthStore()
              if (authStore.user?.role) {
                return getDefaultRouteForRole(authStore.user.role)
              }
            } catch (e) { /* store not ready */ }
            return '/dashboard'
          }
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: {
            requiresAuth: true
          }
        },
        {
          path: 'dashboard/kepala-sppg',
          name: 'dashboard-kepala-sppg',
          component: () => import('@/views/DashboardKepalaSSPGView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg'],
            title: 'Dashboard'
          }
        },
        // --- New multi-tenancy routes ---
        {
          path: 'yayasan',
          name: 'yayasan',
          component: () => import('@/views/YayasanListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'admin_bgn'],
            requiredModule: 'manajemen_yayasan',
            title: 'Manajemen Yayasan'
          }
        },
        {
          path: 'sppg',
          name: 'sppg',
          component: () => import('@/views/SPPGListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'admin_bgn'],
            requiredModule: 'manajemen_sppg',
            title: 'Manajemen SPPG'
          }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('@/views/UserManagementView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'kepala_yayasan', 'kepala_sppg'],
            requiredModule: 'manajemen_user',
            title: 'Manajemen User'
          }
        },
        {
          path: 'dashboard-bgn',
          name: 'dashboard-bgn',
          component: () => import('@/views/DashboardBGNView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'admin_bgn'],
            requiredModule: 'dashboard_bgn',
            title: 'Dashboard BGN'
          }
        },
        {
          path: 'dashboard-yayasan',
          name: 'dashboard-yayasan',
          component: () => import('@/views/DashboardKepalaYayasanView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'admin_bgn', 'kepala_yayasan'],
            requiredModule: 'dashboard_yayasan',
            title: 'Dashboard Yayasan'
          }
        },
        // --- End new multi-tenancy routes ---
        {
          path: 'dashboard/kepala-yayasan',
          name: 'dashboard-kepala-yayasan',
          component: () => import('@/views/DashboardKepalaYayasanView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_yayasan'],
            title: 'Dashboard Kepala Yayasan'
          }
        },
        {
          path: 'recipes',
          name: 'recipes',
          component: () => import('@/views/RecipeListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi'],
            title: 'Manajemen Resep'
          }
        },
        {
          path: 'semi-finished',
          name: 'semi-finished',
          component: () => import('@/views/SemiFinishedGoodsView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi', 'chef'],
            title: 'Barang Setengah Jadi'
          }
        },
        {
          path: 'menu-planning',
          name: 'menu-planning',
          component: () => import('@/views/MenuPlanningView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi'],
            title: 'Perencanaan Menu'
          }
        },
        {
          path: 'kds/cooking',
          name: 'kds-cooking',
          component: () => import('@/views/KDSCookingView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi', 'chef'],
            title: 'KDS - Dapur Memasak'
          }
        },
        {
          path: 'kds/packing',
          name: 'kds-packing',
          component: () => import('@/views/KDSPackingView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi', 'chef', 'packing'],
            title: 'KDS - Packing'
          }
        },
        {
          path: 'kds/cleaning',
          name: 'kds-cleaning',
          component: () => import('@/views/KDSCleaningView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kebersihan'],
            title: 'KDS - Pencucian Ompreng'
          }
        },
        {
          path: 'suppliers',
          name: 'suppliers',
          component: () => import('@/views/SupplierListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan', 'kepala_yayasan'],
            title: 'Manajemen Supplier'
          }
        },
        {
          path: 'purchase-orders',
          name: 'purchase-orders',
          component: () => import('@/views/PurchaseOrderListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan', 'kepala_yayasan'],
            title: 'Purchase Order'
          }
        },
        {
          path: 'goods-receipts',
          name: 'goods-receipts',
          component: () => import('@/views/GoodsReceiptView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan'],
            title: 'Penerimaan Barang'
          }
        },
        {
          path: 'inventory',
          name: 'inventory',
          component: () => import('@/views/InventoryView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan', 'akuntan'],
            title: 'Manajemen Bahan Baku'
          }
        },
        {
          path: 'inventory/stok-opname/create',
          name: 'stok-opname-create',
          component: () => import('@/components/StokOpnameForm.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan'],
            title: 'Buat Stok Opname'
          }
        },
        {
          path: 'inventory/stok-opname/:id/edit',
          name: 'stok-opname-edit',
          component: () => import('@/components/StokOpnameForm.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan'],
            title: 'Edit Stok Opname'
          }
        },
        {
          path: 'inventory/stok-opname/:id',
          name: 'stok-opname-detail',
          component: () => import('@/components/StokOpnameDetail.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'pengadaan', 'akuntan'],
            title: 'Detail Stok Opname'
          }
        },
        {
          path: 'schools',
          name: 'schools',
          component: () => import('@/views/SchoolListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'driver', 'asisten'],
            title: 'Manajemen Sekolah'
          }
        },
        {
          path: 'schools/create',
          name: 'school-create',
          component: () => import('@/views/SchoolFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg'],
            title: 'Tambah Sekolah'
          }
        },
        {
          path: 'schools/:id/edit',
          name: 'school-edit',
          component: () => import('@/views/SchoolFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg'],
            title: 'Edit Sekolah'
          }
        },
        {
          path: 'delivery-tasks',
          name: 'delivery-tasks',
          component: () => import('@/views/DeliveryTaskListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'driver', 'asisten'],
            title: 'Manajemen Tugas Pengiriman'
          }
        },
        {
          path: 'delivery-tasks/create',
          name: 'delivery-task-create',
          component: () => import('@/views/DeliveryTaskFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg'],
            title: 'Buat Tugas Pengiriman'
          }
        },
        {
          path: 'delivery-tasks/:id/edit',
          name: 'delivery-task-edit',
          component: () => import('@/views/DeliveryTaskFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg'],
            title: 'Edit Tugas Pengiriman'
          }
        },
        {
          path: 'ompreng-tracking',
          name: 'ompreng-tracking',
          component: () => import('@/views/OmprengTrackingView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'driver', 'asisten'],
            title: 'Pelacakan Ompreng'
          }
        },
        {
          path: 'monitoring-activity',
          name: 'monitoring-activity',
          component: () => import('@/views/logistics/MonitoringDashboardView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kepala_yayasan', 'akuntan', 'ahli_gizi', 'pengadaan', 'chef', 'packing', 'driver', 'asisten_lapangan'],
            title: 'Monitoring Aktivitas'
          }
        },
        {
          path: 'monitoring-activity/deliveries/:id',
          name: 'delivery-detail',
          component: () => import('@/views/logistics/DeliveryDetailView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kepala_yayasan', 'akuntan', 'ahli_gizi', 'pengadaan', 'chef', 'packing', 'driver', 'asisten_lapangan'],
            title: 'Detail Aktivitas'
          }
        },
        {
          path: 'employees',
          name: 'employees',
          component: () => import('@/views/EmployeeListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Manajemen Karyawan'
          }
        },
        {
          path: 'employees/create',
          name: 'employee-create',
          component: () => import('@/views/EmployeeFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Tambah Karyawan'
          }
        },
        {
          path: 'employees/:id/edit',
          name: 'employee-edit',
          component: () => import('@/views/EmployeeFormView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Edit Karyawan'
          }
        },
        {
          path: 'attendance-report',
          name: 'attendance-report',
          component: () => import('@/views/AttendanceReportView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Laporan Absensi'
          }
        },
        {
          path: 'attendance-config',
          name: 'attendance-config',
          component: () => import('@/views/AttendanceConfigView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Konfigurasi Absensi'
          }
        },
        {
          path: 'assets',
          name: 'assets',
          component: () => import('@/views/AssetListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Manajemen Aset Dapur'
          }
        },
        {
          path: 'cash-flow',
          name: 'cash-flow',
          component: () => import('@/views/CashFlowListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Manajemen Arus Kas'
          }
        },
        {
          path: 'financial-reports',
          name: 'financial-reports',
          component: () => import('@/views/FinancialReportView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'akuntan'],
            title: 'Laporan Keuangan'
          }
        },
        {
          path: 'audit-trail',
          name: 'audit-trail',
          component: () => import('@/views/AuditTrailView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'admin_bgn', 'kepala_yayasan', 'kepala_sppg'],
            title: 'Audit Trail'
          }
        },
        {
          path: 'system-config',
          name: 'system-config',
          component: () => import('@/views/SystemConfigView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['superadmin', 'kepala_sppg'],
            title: 'Konfigurasi Sistem'
          }
        },
        {
          path: 'activity-tracker',
          name: 'activity-tracker',
          component: () => import('@/views/ActivityTrackerListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kepala_yayasan', 'akuntan'],
            title: 'Aktivitas Pelacakan'
          }
        },
        {
          path: 'activity-tracker/:id',
          name: 'activity-tracker-detail',
          component: () => import('@/views/ActivityTrackerDetailView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kepala_yayasan', 'akuntan'],
            title: 'Detail Aktivitas Pelacakan'
          }
        },
        {
          path: 'reviews',
          name: 'reviews',
          component: () => import('@/views/ReviewListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'kepala_yayasan', 'akuntan'],
            title: 'Ulasan & Rating'
          }
        },
        // --- RAB routes ---
        {
          path: 'rab',
          name: 'rab',
          component: () => import('@/views/RABListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi', 'kepala_yayasan'],
            title: 'Daftar RAB'
          }
        },
        {
          path: 'rab/:id',
          name: 'rab-detail',
          component: () => import('@/views/RABDetailView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_sppg', 'ahli_gizi', 'kepala_yayasan'],
            title: 'Detail RAB'
          }
        },
        // --- Invoice routes ---
        {
          path: 'invoices',
          name: 'invoices',
          component: () => import('@/views/InvoiceListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_yayasan', 'supplier'],
            title: 'Invoice'
          }
        },
        // --- Supplier Product Catalog (kepala_yayasan) ---
        {
          path: 'supplier-products',
          name: 'supplier-products',
          component: () => import('@/views/SupplierProductListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_yayasan'],
            title: 'Katalog Supplier'
          }
        },
        // --- Supplier Portal routes ---
        {
          path: 'supplier-dashboard',
          name: 'supplier-dashboard',
          component: () => import('@/views/SupplierDashboardView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['supplier'],
            title: 'Dashboard Supplier'
          }
        },
        {
          path: 'supplier-products/manage',
          name: 'supplier-products-manage',
          component: () => import('@/views/SupplierProductManageView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['supplier'],
            title: 'Kelola Produk'
          }
        },
        {
          path: 'supplier-po',
          name: 'supplier-po',
          component: () => import('@/views/SupplierPOListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['supplier'],
            title: 'Purchase Order'
          }
        },
        {
          path: 'supplier-invoices',
          name: 'supplier-invoices',
          component: () => import('@/views/SupplierInvoiceView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['supplier'],
            title: 'Invoice Supplier'
          }
        },
        // --- Risk Assessment routes ---
        {
          path: 'risk-assessment',
          name: 'risk-assessment',
          component: () => import('@/views/RiskAssessmentListView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_yayasan', 'superadmin'],
            title: 'Risk Assessment'
          }
        },
        {
          path: 'risk-assessment/:id',
          name: 'risk-assessment-detail',
          component: () => import('@/views/RiskAssessmentDetailView.vue'),
          meta: {
            requiresAuth: true,
            roles: ['kepala_yayasan', 'superadmin'],
            title: 'Detail Risk Assessment'
          }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: () => {
        try {
          const authStore = useAuthStore()
          if (authStore.isAuthenticated && authStore.user?.role) {
            return getDefaultRouteForRole(authStore.user.role)
          }
        } catch (e) { /* store not ready */ }
        return '/login'
      }
    }
  ]
})

// Check if user has required role
const hasRequiredRole = (userRole, requiredRoles) => {
  if (!requiredRoles || requiredRoles.length === 0) {
    return true
  }
  return requiredRoles.includes(userRole)
}

// Check if user has access to the required module
const hasRequiredModule = (userModules, requiredModule) => {
  if (!requiredModule) {
    return true
  }
  if (!Array.isArray(userModules) || userModules.length === 0) {
    return false
  }
  return userModules.includes(requiredModule)
}

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }

    if (!authStore.user) {
      try {
        await authStore.getCurrentUser()
      } catch (error) {
        console.error('Failed to fetch user data:', error)
        authStore.clearAuth()
        next('/login')
        return
      }
    }

    const userRole = authStore.user?.role

    // Check role-based access
    if (to.meta.roles) {
      if (!hasRequiredRole(userRole, to.meta.roles)) {
        message.error('Anda tidak memiliki akses ke halaman ini')
        if (from.name) {
          next(false)
        } else {
          next(getDefaultRouteForRole(userRole))
        }
        return
      }
    }

    // Check module-based access
    if (to.meta.requiredModule) {
      if (!hasRequiredModule(authStore.modules, to.meta.requiredModule)) {
        message.error('Anda tidak memiliki akses ke modul ini')
        if (from.name) {
          next(false)
        } else {
          next(getDefaultRouteForRole(userRole))
        }
        return
      }
    }

    next()
  } else {
    // For non-auth routes (e.g. /login), redirect authenticated users to their default page
    if (to.path === '/login' && authStore.isAuthenticated) {
      const userRole = authStore.user?.role
      // Ignore ?redirect query param — always use role-based default
      next(getDefaultRouteForRole(userRole))
    } else {
      next()
    }
  }
})

router.afterEach((to, from) => {
  const baseTitle = 'Dapur Sehat'
  const pageTitle = to.meta.title || to.name
  document.title = pageTitle ? `${pageTitle} - ${baseTitle}` : baseTitle
})

export default router
