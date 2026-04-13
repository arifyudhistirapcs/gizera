import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Allowed roles for tasks page
const TASK_ALLOWED_ROLES = ['driver', 'asisten_lapangan']

/**
 * Returns the default route path for a given user role.
 * Used for post-login redirect and unauthorized access fallback.
 */
function getDefaultRoute(role) {
  const roleStr = role?.toLowerCase()
  if (roleStr === 'admin_bgn' || roleStr === 'superadmin') return '/dashboard-bgn'
  if (roleStr === 'kepala_yayasan') return '/dashboard-yayasan'
  if (roleStr === 'kepala_sppg') return '/dashboard'
  if (roleStr === 'ahli_gizi') return '/menu-planning'
  if (roleStr === 'sekolah') return '/school-monitoring'
  if (roleStr === 'supplier') return '/supplier-dashboard'
  if (TASK_ALLOWED_ROLES.includes(roleStr)) return '/tasks'
  return '/attendance'
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
      component: () => import('@/layouts/MobileLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: () => {
            const authStore = useAuthStore()
            return getDefaultRoute(authStore.user?.role)
          }
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/DashboardView.vue'),
          meta: { roles: ['kepala_sppg'] }
        },
        // Dashboard Yayasan routes
        {
          path: 'dashboard-yayasan',
          name: 'dashboard-yayasan',
          component: () => import('@/views/DashboardYayasanView.vue'),
          meta: { roles: ['kepala_yayasan'] }
        },
        {
          path: 'dashboard-yayasan/:sppg_id',
          name: 'dashboard-yayasan-sppg',
          component: () => import('@/views/DashboardYayasanView.vue'),
          meta: { roles: ['kepala_yayasan'] }
        },
        // Dashboard BGN routes
        {
          path: 'dashboard-bgn',
          name: 'dashboard-bgn',
          component: () => import('@/views/DashboardBGNView.vue'),
          meta: { roles: ['admin_bgn', 'superadmin'] }
        },
        {
          path: 'dashboard-bgn/:yayasan_id',
          name: 'dashboard-bgn-yayasan',
          component: () => import('@/views/DashboardBGNView.vue'),
          meta: { roles: ['admin_bgn', 'superadmin'] }
        },
        {
          path: 'dashboard-bgn/:yayasan_id/:sppg_id',
          name: 'dashboard-bgn-sppg',
          component: () => import('@/views/DashboardBGNView.vue'),
          meta: { roles: ['admin_bgn', 'superadmin'] }
        },
        {
          path: 'monitoring',
          name: 'monitoring',
          component: () => import('@/views/MonitoringView.vue'),
          meta: { roles: ['kepala_sppg'] }
        },
        {
          path: 'monitoring/:id',
          name: 'monitoring-detail',
          component: () => import('@/views/MonitoringDetailView.vue'),
          meta: { roles: ['kepala_sppg'] }
        },
        {
          path: 'menu-planning',
          name: 'menu-planning',
          component: () => import('@/views/MenuPlanningView.vue'),
          meta: { roles: ['kepala_sppg', 'ahli_gizi'] }
        },
        {
          path: 'school-monitoring',
          name: 'school-monitoring',
          component: () => import('@/views/SchoolMonitoringView.vue'),
          meta: { roles: ['sekolah'] }
        },
        {
          path: 'tasks',
          name: 'tasks',
          component: () => import('@/views/TasksView.vue'),
          meta: { roles: TASK_ALLOWED_ROLES }
        },
        {
          path: 'tasks/:id',
          name: 'task-detail',
          component: () => import('@/views/DeliveryTaskDetailView.vue'),
          meta: { roles: TASK_ALLOWED_ROLES }
        },
        {
          path: 'pickup-tasks/:id',
          name: 'pickup-task-detail',
          component: () => import('@/views/PickupTaskDetailView.vue'),
          meta: { roles: TASK_ALLOWED_ROLES }
        },
        {
          path: 'tasks/:taskId/epod',
          name: 'epod-form',
          component: () => import('@/views/ePODFormView.vue'),
          meta: { roles: TASK_ALLOWED_ROLES }
        },
        {
          path: 'review',
          name: 'review-form',
          component: () => import('@/views/ReviewFormView.vue'),
          meta: { roles: TASK_ALLOWED_ROLES }
        },
        {
          path: 'attendance',
          name: 'attendance',
          component: () => import('@/views/AttendanceView.vue')
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/ProfileView.vue')
        },
        // Risk Assessment routes (kepala_yayasan)
        {
          path: 'risk-assessment',
          name: 'risk-assessment',
          component: () => import('@/views/RiskAssessmentSelectView.vue'),
          meta: { roles: ['kepala_yayasan'] }
        },
        {
          path: 'risk-assessment/:id',
          name: 'risk-assessment-form',
          component: () => import('@/views/RiskAssessmentFormView.vue'),
          meta: { roles: ['kepala_yayasan'] }
        },
        // Supplier Portal routes
        {
          path: 'supplier-dashboard',
          name: 'supplier-dashboard',
          component: () => import('@/views/SupplierDashboardView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-products',
          name: 'supplier-products',
          component: () => import('@/views/SupplierProductsView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-products/form',
          name: 'supplier-product-form',
          component: () => import('@/views/SupplierProductFormView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-products/:id/edit',
          name: 'supplier-product-edit',
          component: () => import('@/views/SupplierProductFormView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-po',
          name: 'supplier-po',
          component: () => import('@/views/SupplierPOListView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-po/:id',
          name: 'supplier-po-detail',
          component: () => import('@/views/SupplierPODetailView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-invoices',
          name: 'supplier-invoices',
          component: () => import('@/views/SupplierInvoiceView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-notifications',
          name: 'supplier-notifications',
          component: () => import('@/views/SupplierNotificationsView.vue'),
          meta: { roles: ['supplier'] }
        },
        {
          path: 'supplier-payments',
          name: 'supplier-payments',
          component: () => import('@/views/SupplierPaymentsView.vue'),
          meta: { roles: ['supplier'] }
        }
      ]
    }
  ]
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // Check if any matched route requires auth
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }

  if (to.path === '/login' && authStore.isAuthenticated) {
    next(getDefaultRoute(authStore.user?.role))
    return
  }

  // Check role-based access: look through matched routes for meta.roles
  const requiredRoles = to.matched.reduce((roles, record) => {
    return record.meta.roles || roles
  }, null)

  if (requiredRoles) {
    const userRole = authStore.user?.role?.toLowerCase()
    if (!requiredRoles.includes(userRole)) {
      next(getDefaultRoute(authStore.user?.role))
      return
    }
  }

  next()
})

export { getDefaultRoute }
export default router
