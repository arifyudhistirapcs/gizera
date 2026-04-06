<template>
  <aside
    class="h-sidebar"
    :class="{ 'collapsed': isCollapsed, 'mobile': isMobile }"
  >
    <!-- Logo Area -->
    <div class="sidebar-logo">
      <div class="logo-content">
        <img
          v-if="!isCollapsed"
          :src="isDark ? '/gizera-dark.png' : '/gizera-light.png'"
          alt="Gizera"
          class="logo-img"
        />
        <img
          v-else
          :src="isDark ? '/gizera-dark.png' : '/gizera-light.png'"
          alt="Gizera"
          class="logo-img-collapsed"
        />
      </div>
      <button
        v-if="!isMobile"
        class="collapse-btn"
        @click="toggleCollapse"
        :aria-label="isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
      >
        <MenuUnfoldOutlined v-if="isCollapsed" />
        <MenuFoldOutlined v-else />
      </button>
    </div>

    <!-- Menu Items -->
    <nav class="sidebar-nav">
      <template v-for="item in filteredMenuItems" :key="item.key">
        <!-- Menu item dengan submenu -->
        <div v-if="item.children && item.children.length > 0" class="menu-group">
          <div
            class="menu-item"
            :class="{ 'active': isGroupActive(item), 'expanded': expandedGroups.includes(item.key) }"
            @click="toggleGroup(item.key)"
          >
            <span class="menu-emoji">{{ item.emoji }}</span>
            <span v-if="!isCollapsed" class="menu-label">{{ item.label }}</span>
            <DownOutlined
              v-if="!isCollapsed"
              class="menu-arrow"
              :class="{ 'rotated': expandedGroups.includes(item.key) }"
            />
          </div>
          
          <!-- Submenu -->
          <Transition name="submenu">
            <div
              v-if="!isCollapsed && expandedGroups.includes(item.key)"
              class="submenu"
            >
              <router-link
                v-for="child in item.children"
                :key="child.key"
                :to="child.route"
                class="submenu-item"
                :class="{ 'active': isActive(child.route) }"
              >
                <span class="submenu-emoji">{{ child.emoji }}</span>
                <span class="submenu-label">{{ child.label }}</span>
              </router-link>
            </div>
          </Transition>
        </div>

        <!-- Menu item tanpa submenu -->
        <router-link
          v-else
          :to="item.route"
          class="menu-item"
          :class="{ 'active': isActive(item.route) }"
        >
          <span class="menu-emoji">{{ item.emoji }}</span>
          <span v-if="!isCollapsed" class="menu-label">{{ item.label }}</span>
        </router-link>
      </template>
    </nav>

    <!-- Sidebar Promo Card (hidden) -->
    <!--
    <div v-if="!isCollapsed" class="sidebar-promo">
      <div class="sidebar-promo__icon">📊</div>
      <div class="sidebar-promo__text">Laporan Harian</div>
      <div class="sidebar-promo__sub">Lihat ringkasan hari ini</div>
      <button class="sidebar-promo__btn" @click="router.push('/monitoring-activity')">Lihat Laporan</button>
    </div>
    -->

    <!-- User Info -->
    <div v-if="!isCollapsed" class="user-info">
      <UserOutlined class="user-info__avatar" />
      <div class="user-info__details">
        <div class="user-info__name">{{ authStore.user?.full_name || 'User' }}</div>
        <div class="user-info__role">{{ userRoleLabel }}</div>
      </div>
    </div>
    <div v-else class="user-info user-info--collapsed">
      <UserOutlined class="user-info__avatar" />
    </div>

    <!-- Logout Button -->
    <button
      class="logout-button"
      @click="handleLogout"
      aria-label="Keluar"
    >
      <LogoutOutlined class="logout-icon" />
      <span v-if="!isCollapsed" class="logout-label">Keluar</span>
    </button>
  </aside>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePermissions } from '@/composables/usePermissions'
import { getRoleLabel } from '@/utils/permissions'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useDarkMode } from '@/composables/useDarkMode'
import { useAuthStore } from '@/stores/auth'
import { message } from 'ant-design-vue'
import {
  HomeOutlined,
  DashboardOutlined,
  MonitorOutlined,
  AppstoreOutlined,
  ShoppingOutlined,
  CarOutlined,
  TeamOutlined,
  DollarOutlined,
  SettingOutlined,
  CalendarOutlined,
  FileTextOutlined,
  InboxOutlined,
  ShopOutlined,
  ShoppingCartOutlined,
  DatabaseOutlined,
  EnvironmentOutlined,
  UserOutlined,
  ClockCircleOutlined,
  BankOutlined,
  LineChartOutlined,
  AuditOutlined,
  ControlOutlined,
  DownOutlined,
  LogoutOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  FireOutlined,
  ClearOutlined,
  StarOutlined,
  GlobalOutlined,
  FundOutlined,
  SafetyCertificateOutlined
} from '@ant-design/icons-vue'

/**
 * HSidebar Component
 * 
 * Main navigation sidebar untuk ERP SPPG dengan:
 * - Width 280px (expanded) / 80px (collapsed)
 * - Background white (light) / #111C44 (dark)
 * - Logo area 64px height
 * - Menu items 44px height
 * - Icon 20px, font 14px medium
 * - Active state (purple bg + white text)
 * - Hover state (#F4F7FE bg)
 * - Collapsible toggle button
 * - Nested menu support
 * - Role-based menu filtering
 * - Smooth width transition (300ms)
 * - Mobile: hidden by default (used inside MobileDrawer)
 */

const props = defineProps({
  /**
   * Collapsed state (controlled)
   */
  collapsed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:collapsed'])

const route = useRoute()
const router = useRouter()
const { can, isAnyRole } = usePermissions()
const { isMobile } = useBreakpoint()
const { isDark } = useDarkMode()
const authStore = useAuthStore()

// User role label for sidebar display
const userRoleLabel = computed(() => getRoleLabel(authStore.user?.role || ''))

// Local collapsed state
const isCollapsed = ref(props.collapsed)

// Expanded groups tracking
const expandedGroups = ref([])

// Roles that should NOT see operational modules (KDS, Menu Planning, PO, Cooking, Packing, Attendance)
const NON_OPERATIONAL_ROLES = ['superadmin', 'admin_bgn', 'kepala_yayasan']

/**
 * Multi-tenant management menu items — shown based on modules from auth store
 */
const multiTenantMenuItems = [
  {
    key: 'dashboard-bgn',
    label: 'Dashboard BGN',
    emoji: '🌐',
    icon: GlobalOutlined,
    route: '/dashboard-bgn',
    module: 'dashboard_bgn'
  },
  {
    key: 'dashboard-yayasan',
    label: 'Dashboard Yayasan',
    emoji: '📊',
    icon: FundOutlined,
    route: '/dashboard-yayasan',
    module: 'dashboard_yayasan'
  },
  {
    key: 'manajemen-yayasan',
    label: 'Manajemen Yayasan',
    emoji: '🏛️',
    icon: BankOutlined,
    route: '/yayasan',
    module: 'manajemen_yayasan'
  },
  {
    key: 'manajemen-sppg',
    label: 'Manajemen SPPG',
    emoji: '🏪',
    icon: ShopOutlined,
    route: '/sppg',
    module: 'manajemen_sppg'
  },
  {
    key: 'manajemen-user',
    label: 'Manajemen User',
    emoji: '👥',
    icon: TeamOutlined,
    route: '/users',
    module: 'manajemen_user'
  }
]

/**
 * Operational menu items — hidden from superadmin, admin_bgn, kepala_yayasan
 */
const operationalMenuItems = [
  {
    key: 'dashboard',
    label: 'Dashboard',
    emoji: '🏠',
    icon: HomeOutlined,
    route: '/dashboard',
    roles: ['kepala_sppg', 'akuntan', 'ahli_gizi', 'pengadaan'],
    operational: true
  },
  {
    key: 'monitoring',
    label: 'Monitoring Aktivitas',
    emoji: '📊',
    icon: MonitorOutlined,
    route: '/monitoring-activity',
    roles: ['kepala_sppg', 'akuntan', 'ahli_gizi', 'pengadaan', 'chef', 'packing', 'driver', 'asisten_lapangan'],
    operational: true
  },
  {
    key: 'reviews',
    label: 'Ulasan & Rating',
    emoji: '⭐',
    icon: StarOutlined,
    route: '/reviews',
    roles: ['kepala_sppg', 'akuntan'],
    operational: true
  },
  {
    key: 'display',
    label: 'Display / KDS',
    emoji: '🍳',
    icon: DashboardOutlined,
    roles: ['kepala_sppg', 'ahli_gizi', 'chef', 'packing', 'kebersihan'],
    operational: true,
    children: [
      {
        key: 'kds-cooking',
        label: 'Dapur',
        emoji: '🔥',
        icon: FireOutlined,
        route: '/kds/cooking',
        roles: ['kepala_sppg', 'ahli_gizi', 'chef']
      },
      {
        key: 'kds-packing',
        label: 'Pengemasan',
        emoji: '📦',
        icon: InboxOutlined,
        route: '/kds/packing',
        roles: ['kepala_sppg', 'ahli_gizi', 'chef', 'packing']
      },
      {
        key: 'kds-cleaning',
        label: 'Kebersihan',
        emoji: '🧹',
        icon: ClearOutlined,
        route: '/kds/cleaning',
        roles: ['kepala_sppg', 'kebersihan']
      }
    ]
  },
  {
    key: 'menu',
    label: 'Menu & Komponen',
    emoji: '📋',
    icon: FileTextOutlined,
    roles: ['kepala_sppg', 'ahli_gizi', 'chef'],
    operational: true,
    children: [
      {
        key: 'menu-planning',
        label: 'Perencanaan Menu',
        emoji: '📅',
        icon: CalendarOutlined,
        route: '/menu-planning',
        roles: ['kepala_sppg', 'ahli_gizi']
      },
      {
        key: 'recipes',
        label: 'Manajemen Menu',
        emoji: '📝',
        icon: FileTextOutlined,
        route: '/recipes',
        roles: ['kepala_sppg', 'ahli_gizi']
      },
      {
        key: 'semi-finished',
        label: 'Manajemen Komponen',
        emoji: '🥘',
        icon: InboxOutlined,
        route: '/semi-finished',
        roles: ['kepala_sppg', 'ahli_gizi', 'chef']
      }
    ]
  },
  {
    key: 'supply-chain',
    label: 'Supply Chain',
    emoji: '🛒',
    icon: ShoppingOutlined,
    roles: ['kepala_sppg', 'pengadaan'],
    operational: true,
    children: [
      {
        key: 'suppliers',
        label: 'Supplier',
        emoji: '🏪',
        icon: ShopOutlined,
        route: '/suppliers',
        roles: ['kepala_sppg', 'pengadaan']
      },
      {
        key: 'purchase-orders',
        label: 'Purchase Order',
        emoji: '🛍️',
        icon: ShoppingCartOutlined,
        route: '/purchase-orders',
        roles: ['kepala_sppg', 'pengadaan']
      },
      {
        key: 'goods-receipts',
        label: 'Penerimaan Barang',
        emoji: '📥',
        icon: InboxOutlined,
        route: '/goods-receipts',
        roles: ['kepala_sppg', 'pengadaan']
      },
      {
        key: 'inventory',
        label: 'Manajemen Bahan Baku',
        emoji: '🗄️',
        icon: DatabaseOutlined,
        route: '/inventory',
        roles: ['kepala_sppg', 'pengadaan', 'akuntan']
      }
    ]
  },
  {
    key: 'logistics',
    label: 'Logistik',
    emoji: '🚚',
    icon: CarOutlined,
    roles: ['kepala_sppg', 'driver', 'asisten'],
    operational: true,
    children: [
      {
        key: 'schools',
        label: 'Data Sekolah',
        emoji: '🏫',
        icon: EnvironmentOutlined,
        route: '/schools',
        roles: ['kepala_sppg', 'driver', 'asisten']
      },
      {
        key: 'delivery-tasks',
        label: 'Tugas Pengiriman & Pengambilan',
        emoji: '📍',
        icon: CarOutlined,
        route: '/delivery-tasks',
        roles: ['kepala_sppg', 'driver', 'asisten']
      }
    ]
  },
  {
    key: 'sdm',
    label: 'SDM',
    emoji: '👥',
    icon: TeamOutlined,
    roles: ['kepala_sppg', 'akuntan'],
    operational: true,
    children: [
      {
        key: 'employees',
        label: 'Data Karyawan',
        emoji: '👤',
        icon: UserOutlined,
        route: '/employees',
        roles: ['kepala_sppg', 'akuntan']
      },
      {
        key: 'attendance-report',
        label: 'Laporan Absensi',
        emoji: '⏰',
        icon: ClockCircleOutlined,
        route: '/attendance-report',
        roles: ['kepala_sppg', 'akuntan']
      },
      {
        key: 'attendance-config',
        label: 'Konfigurasi Absensi',
        emoji: '⚙️',
        icon: ControlOutlined,
        route: '/attendance-config',
        roles: ['kepala_sppg', 'akuntan']
      }
    ]
  },
  {
    key: 'finance',
    label: 'Keuangan',
    emoji: '💰',
    icon: DollarOutlined,
    roles: ['kepala_sppg', 'akuntan'],
    operational: true,
    children: [
      {
        key: 'assets',
        label: 'Aset Dapur',
        emoji: '🏦',
        icon: BankOutlined,
        route: '/assets',
        roles: ['kepala_sppg', 'akuntan']
      },
      {
        key: 'cash-flow',
        label: 'Arus Kas',
        emoji: '📈',
        icon: LineChartOutlined,
        route: '/cash-flow',
        roles: ['kepala_sppg', 'akuntan']
      },
      {
        key: 'financial-reports',
        label: 'Laporan Keuangan',
        emoji: '📄',
        icon: FileTextOutlined,
        route: '/financial-reports',
        roles: ['kepala_sppg', 'akuntan']
      }
    ]
  },
  {
    key: 'risk-assessment',
    label: 'Risk Assessment',
    emoji: '🛡️',
    icon: SafetyCertificateOutlined,
    route: '/risk-assessment',
    roles: ['kepala_yayasan', 'superadmin']
  },
  {
    key: 'system',
    label: 'Sistem',
    emoji: '⚙️',
    icon: SettingOutlined,
    roles: ['kepala_sppg', 'superadmin'],
    children: [
      {
        key: 'audit-trail',
        label: 'Audit Trail',
        emoji: '📋',
        icon: AuditOutlined,
        route: '/audit-trail',
        roles: ['kepala_sppg', 'superadmin']
      },
      {
        key: 'system-config',
        label: 'Konfigurasi',
        emoji: '🔧',
        icon: ControlOutlined,
        route: '/system-config',
        roles: ['kepala_sppg', 'superadmin']
      }
    ]
  }
]

/**
 * Filter menu items berdasarkan role user dan modules dari auth store
 * - Multi-tenant items: shown based on `modules` array from auth store
 * - Operational items: hidden from superadmin, admin_bgn, kepala_yayasan (Req 17.4)
 * - SPPG-level roles: see operational items as before (Req 17.5)
 */
const filteredMenuItems = computed(() => {
  const userModules = authStore.modules || []
  const userRole = authStore.role || ''
  const isNonOperational = NON_OPERATIONAL_ROLES.includes(userRole)

  // 1. Filter multi-tenant menu items based on modules
  const tenantItems = multiTenantMenuItems.filter(item => {
    return userModules.includes(item.module)
  })

  // 2. Filter operational menu items
  const opItems = operationalMenuItems.filter(item => {
    // Hide operational modules from non-operational roles (Req 17.4)
    if (item.operational && isNonOperational) {
      return false
    }

    // Check if user has required role
    if (item.roles && !isAnyRole(item.roles)) {
      return false
    }

    // Filter children jika ada
    if (item.children) {
      item.children = item.children.filter(child => {
        return !child.roles || isAnyRole(child.roles)
      })

      // Hide parent jika tidak ada children yang visible
      return item.children.length > 0
    }

    return true
  })

  // Multi-tenant items first, then operational items
  return [...tenantItems, ...opItems]
})

/**
 * Check if route is active
 */
const isActive = (routePath) => {
  if (!routePath) return false
  return route.path === routePath || route.path.startsWith(routePath + '/')
}

/**
 * Check if any child in group is active
 */
const isGroupActive = (item) => {
  if (!item.children) return false
  return item.children.some(child => isActive(child.route))
}

/**
 * Toggle group expansion
 */
const toggleGroup = (key) => {
  if (isCollapsed.value) {
    // Jika collapsed, expand sidebar dulu
    toggleCollapse()
  }
  
  const index = expandedGroups.value.indexOf(key)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(key)
  }
}

/**
 * Toggle sidebar collapse
 */
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  emit('update:collapsed', isCollapsed.value)
  if (isCollapsed.value) {
    expandedGroups.value = []
  }
}

/**
 * Handle logout
 */
const handleLogout = async () => {
  try {
    await authStore.logout()
    message.success('Berhasil keluar')
    router.push('/login')
  } catch (error) {
    message.error('Gagal keluar')
  }
}

/**
 * Auto-expand active group
 */
watch(() => route.path, () => {
  // Find and expand group containing active route (only operational items have children)
  operationalMenuItems.forEach(item => {
    if (item.children && isGroupActive(item)) {
      if (!expandedGroups.value.includes(item.key)) {
        expandedGroups.value.push(item.key)
      }
    }
  })
}, { immediate: true })

/**
 * Watch props.collapsed changes
 */
watch(() => props.collapsed, (newValue) => {
  isCollapsed.value = newValue
})
</script>

<style scoped>
.h-sidebar {
  position: fixed;
  top: 12px;
  left: 12px;
  bottom: 12px;
  width: 220px;
  background-color: #FFFFFF;
  border-radius: 20px;
  border: none;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  display: flex;
  flex-direction: column;
  transition: width 300ms ease;
  z-index: 100;
  overflow: hidden;
}

/* Dark mode */
.dark .h-sidebar {
  background-color: #252525;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.3);
}

/* Collapsed state */
.h-sidebar.collapsed {
  width: 72px;
}

/* Mobile: hidden by default (will be used inside MobileDrawer) */
.h-sidebar.mobile {
  position: static;
  width: 220px;
  border-radius: 0;
  top: 0;
  left: 0;
  bottom: 0;
  box-shadow: none;
}

/* Logo Area - compact */
.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-bottom: none;
}

.dark .sidebar-logo {
  /* no border in floating sidebar */
}

.collapsed .sidebar-logo {
  justify-content: center;
  padding: 10px 8px;
  gap: 0;
}

.collapse-btn {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--h-text-secondary, #74788C);
  border-radius: 6px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.collapse-btn:hover {
  background-color: #F7F8FA;
  color: #303030;
}

.dark .collapse-btn {
  color: #D8D8DB;
}

.dark .collapse-btn:hover {
  background-color: #303030;
  color: #F7F8FA;
}

.collapse-btn:active {
  transform: scale(0.9);
}

.logo-content {
  display: flex;
  align-items: center;
  min-width: 0;
  flex: 1;
  height: 100%;
}

.logo-img {
  height: 56px;
  max-height: 56px;
  width: auto;
  object-fit: contain;
}

.logo-img-collapsed {
  height: 40px;
  width: auto;
  object-fit: contain;
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 16px 0;
  
  /* Smooth scrolling */
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.15) transparent;
}

.sidebar-nav::-webkit-scrollbar {
  width: 4px;
}

.sidebar-nav::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-nav::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.15);
  border-radius: 2px;
}

/* Menu Group */
.menu-group {
  margin-bottom: 4px;
}

/* Menu Item - 42px height */
.menu-item {
  height: 42px;
  display: flex;
  align-items: center;
  padding: 0 14px;
  margin: 2px 8px;
  border-radius: 10px;
  cursor: pointer;
  text-decoration: none;
  color: #303030;
  transition: all 0.2s ease;
  position: relative;
  font-size: 13px;
}

.menu-item:hover {
  background-color: #F7F8FA;
  color: #303030;
}

.dark .menu-item:hover {
  background-color: #303030;
  color: #F7F8FA;
}

/* Active state - Culinary Coral bg + white text */
.menu-item.active {
  background-color: #C94A3A;
  color: #FFFFFF;
}

.dark .menu-item.active {
  background-color: #C94A3A;
  color: #FFFFFF;
}

/* Menu Emoji - 16px */
.menu-emoji {
  font-size: 16px;
  min-width: 22px;
  text-align: center;
  margin-right: 8px;
  line-height: 1;
  transition: margin 300ms ease;
}

.collapsed .menu-emoji {
  margin-right: 0;
}

/* Menu Label - 13px medium */
.menu-label {
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

/* Menu Arrow */
.menu-arrow {
  font-size: 12px;
  margin-left: auto;
  transition: transform 0.3s ease;
}

.menu-arrow.rotated {
  transform: rotate(180deg);
}

/* Submenu */
.submenu {
  padding-left: 0;
  overflow: hidden;
}

.submenu-item {
  height: 38px;
  display: flex;
  align-items: center;
  padding: 0 12px 0 28px;
  margin: 2px 8px;
  border-radius: 10px;
  cursor: pointer;
  text-decoration: none;
  color: #303030;
  transition: all 0.2s ease;
}

.submenu-item:hover {
  background-color: #F7F8FA;
  color: #303030;
}

.dark .submenu-item:hover {
  background-color: #303030;
  color: #F7F8FA;
}

.submenu-item.active {
  background-color: #C94A3A;
  color: #FFFFFF;
}

.dark .submenu-item.active {
  background-color: #C94A3A;
  color: #FFFFFF;
}

.submenu-emoji {
  font-size: 13px;
  min-width: 18px;
  text-align: center;
  margin-right: 8px;
  line-height: 1;
}

.submenu-label {
  font-size: 13px;
  font-weight: 400;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Submenu transition */
.submenu-enter-active,
.submenu-leave-active {
  transition: all 0.3s ease;
}

.submenu-enter-from,
.submenu-leave-to {
  opacity: 0;
  max-height: 0;
}

.submenu-enter-to,
.submenu-leave-from {
  opacity: 1;
  max-height: 500px;
}

/* User Info */
.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  margin: 0 8px;
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}

.dark .user-info {
  border-top-color: rgba(255, 255, 255, 0.08);
}

.user-info--collapsed {
  justify-content: center;
  padding: 12px 8px;
}

.user-info__avatar {
  font-size: 20px;
  color: #303030;
  flex-shrink: 0;
}

.dark .user-info__avatar {
  color: #F7F8FA;
}

.user-info__details {
  min-width: 0;
  flex: 1;
}

.user-info__name {
  font-size: 13px;
  font-weight: 600;
  color: var(--h-text-primary, #322837);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .user-info__name {
  color: #F7F8FA;
}

.user-info__role {
  font-size: 11px;
  color: var(--h-text-secondary, #74788C);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.dark .user-info__role {
  color: #D8D8DB;
}

/* Logout Button */
.logout-button {
  height: 42px;
  margin: 8px;
  border: none;
  border-radius: 10px;
  background-color: transparent;
  color: var(--h-text-secondary, #74788C);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s ease;
}

.logout-button:hover {
  background-color: rgba(238, 93, 80, 0.1);
  color: #EE5D50;
}

.dark .logout-button {
  color: #D8D8DB;
}

.dark .logout-button:hover {
  background-color: rgba(238, 93, 80, 0.15);
  color: #EE5D50;
}

.logout-button:active {
  transform: scale(0.95);
}

.logout-icon {
  font-size: 20px;
}

.logout-label {
  white-space: nowrap;
}

/* Sidebar Promo Card */
.sidebar-promo {
  margin: 8px;
  padding: 14px;
  background: #303030;
  border-radius: 12px;
  color: #fff;
  text-align: center;
}

.sidebar-promo__icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.sidebar-promo__text {
  font-size: 13px;
  font-weight: 600;
}

.sidebar-promo__sub {
  font-size: 11px;
  opacity: 0.7;
  margin-top: 4px;
}

.sidebar-promo__btn {
  display: inline-block;
  margin-top: 10px;
  padding: 6px 16px;
  background: #C94A3A;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  min-height: auto;
  min-width: auto;
  transition: background 0.2s ease;
}

.sidebar-promo__btn:hover {
  background: #A33D30;
}
</style>
