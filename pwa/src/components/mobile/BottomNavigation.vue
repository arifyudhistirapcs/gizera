<template>
  <div class="bottom-navigation-wrapper">
    <div class="bottom-navigation">
      <div
        v-for="item in allItems"
        :key="item.name"
        :class="['nav-item', { 'nav-item--active': isActive(item) }]"
        @click="onNavClick(item)"
      >
        <div class="nav-icon-wrap" :class="{ 'nav-icon-wrap--active': isActive(item) }">
          <van-badge v-if="item.name === 'supplierNotif' && unreadCount > 0" :content="unreadCount" :max="99">
            <van-icon
              :name="isActive(item) ? item.activeIcon : item.icon"
              :color="isActive(item) ? '#fff' : '#8C8C8C'"
              size="20"
            />
          </van-badge>
          <van-icon
            v-else
            :name="isActive(item) ? item.activeIcon : item.icon"
            :color="isActive(item) ? '#fff' : '#8C8C8C'"
            size="20"
          />
        </div>
        <span class="nav-item__label">{{ item.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const unreadCount = ref(0)
let unreadInterval = null

const NAV_CONFIG = {
  dashboard: { name: 'dashboard', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/dashboard' },
  dashboardYayasan: { name: 'dashboardYayasan', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/dashboard-yayasan' },
  monitoring: { name: 'monitoring', label: 'Monitoring', icon: 'search', activeIcon: 'search', route: '/monitoring' },
  menu: { name: 'menu', label: 'Menu', icon: 'notes-o', activeIcon: 'notes-o', route: '/menu-planning' },
  tasks: { name: 'tasks', label: 'Tugas', icon: 'todo-list-o', activeIcon: 'todo-list-o', route: '/tasks' },
  riskAssessment: { name: 'riskAssessment', label: 'Audit', icon: 'shield-o', activeIcon: 'shield-o', route: '/risk-assessment' },
  attendance: { name: 'attendance', label: 'Absensi', icon: 'clock-o', activeIcon: 'clock', route: '/attendance' },
  profile: { name: 'profile', label: 'Profil', icon: 'user-o', activeIcon: 'user-o', route: '/profile' },
  schoolMonitoring: { name: 'schoolMonitoring', label: 'Monitoring', icon: 'eye-o', activeIcon: 'eye', route: '/school-monitoring' },
  supplierDashboard: { name: 'supplierDashboard', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/supplier-dashboard' },
  supplierPO: { name: 'supplierPO', label: 'PO', icon: 'orders-o', activeIcon: 'orders-o', route: '/supplier-po' },
  supplierInvoice: { name: 'supplierInvoice', label: 'Invoice', icon: 'bill-o', activeIcon: 'bill-o', route: '/supplier-invoices' },
  supplierNotif: { name: 'supplierNotif', label: 'Notifikasi', icon: 'bell', activeIcon: 'bell', route: '/supplier-notifications' }
}

const ROLE_NAV_MAP = {
  driver: [{ left: ['tasks'], center: 'attendance', right: ['profile'] }],
  asisten_lapangan: [{ left: ['tasks'], center: 'attendance', right: ['profile'] }],
  kepala_sppg: [{ left: ['dashboard', 'monitoring'], center: 'attendance', right: ['menu', 'profile'] }],
  kepala_yayasan: [{ left: ['dashboardYayasan', 'riskAssessment'], center: 'attendance', right: ['profile'] }],
  ahli_gizi: [{ left: ['menu'], center: 'attendance', right: ['profile'] }],
  sekolah: [{ left: ['schoolMonitoring'], center: 'attendance', right: ['profile'] }],
  supplier: [{ left: ['supplierDashboard', 'supplierPO'], center: 'supplierInvoice', right: ['supplierNotif', 'profile'] }]
}

const DEFAULT_NAV = { left: [], center: 'attendance', right: ['profile'] }

const navConfig = computed(() => {
  const role = authStore.user?.role?.toLowerCase() || ''
  return ROLE_NAV_MAP[role]?.[0] || DEFAULT_NAV
})

const leftItems = computed(() => navConfig.value.left.map(key => NAV_CONFIG[key]))
const centerItem = computed(() => NAV_CONFIG[navConfig.value.center])
const rightItems = computed(() => navConfig.value.right.map(key => NAV_CONFIG[key]))
const allItems = computed(() => [...leftItems.value, centerItem.value, ...rightItems.value])

const activeItem = ref(null)

function findActiveItem() {
  const currentPath = route.path
  return allItems.value.find(item => currentPath.startsWith(item.route)) || allItems.value[0]
}

function isActive(item) { return activeItem.value?.name === item?.name }

function onNavClick(item) {
  if (item) {
    activeItem.value = item
    router.push(item.route)
  }
}

watch(() => route.path, () => { activeItem.value = findActiveItem() }, { immediate: true })
onMounted(() => {
  activeItem.value = findActiveItem()
  // Fetch unread notification count for supplier
  if (authStore.user?.role?.toLowerCase() === 'supplier') {
    fetchUnreadCount()
    unreadInterval = setInterval(fetchUnreadCount, 60000) // poll every 60s
  }
})

onUnmounted(() => {
  if (unreadInterval) clearInterval(unreadInterval)
})

async function fetchUnreadCount() {
  try {
    const res = await api.get('/notifications/unread-count')
    unreadCount.value = res.data?.data?.count ?? res.data?.count ?? 0
  } catch (e) {
    // Silent fail
  }
}
</script>

<style scoped>
.bottom-navigation-wrapper {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 0 12px 8px;
  pointer-events: none;
}

.bottom-navigation {
  height: 64px;
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 -2px 20px rgba(0, 0, 0, 0.08), 0 0 0 1px rgba(0, 0, 0, 0.03);
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  padding: 0 8px;
  pointer-events: auto;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 6px 14px;
  -webkit-tap-highlight-color: transparent;
}

.nav-icon-wrap {
  width: 36px;
  height: 36px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.nav-icon-wrap--active {
  background: #F82C17;
  box-shadow: 0 4px 12px rgba(248, 44, 23, 0.3);
  transform: translateY(-2px);
}

.nav-item__label {
  font-size: 10px;
  font-weight: 500;
  color: #8C8C8C;
  transition: color 0.2s ease;
  line-height: 1;
}

.nav-item--active .nav-item__label {
  color: #F82C17;
  font-weight: 700;
}

.nav-item:active .nav-icon-wrap {
  transform: scale(0.9);
}

.nav-item--active:active .nav-icon-wrap {
  transform: translateY(-2px) scale(0.95);
}
</style>
