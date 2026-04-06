<template>
  <div class="bottom-navigation-wrapper">
    <div class="bottom-navigation">
      <div
        v-for="item in allItems"
        :key="item.name"
        :class="['nav-item', { 'nav-item--active': isActive(item) }]"
        @click="onNavClick(item)"
      >
        <van-icon
          :name="isActive(item) ? item.activeIcon : item.icon"
          :color="isActive(item) ? activeColor : inactiveColor"
          size="22"
        />
        <span class="nav-item__label">{{ item.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeColor = '#C94A3A'
const inactiveColor = '#6B6B6B'

const NAV_CONFIG = {
  dashboard: { name: 'dashboard', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/dashboard' },
  dashboardYayasan: { name: 'dashboardYayasan', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/dashboard-yayasan' },
  monitoring: { name: 'monitoring', label: 'Monitoring', icon: 'search', activeIcon: 'search', route: '/monitoring' },
  menu: { name: 'menu', label: 'Menu', icon: 'notes-o', activeIcon: 'notes-o', route: '/menu-planning' },
  tasks: { name: 'tasks', label: 'Tugas', icon: 'todo-list-o', activeIcon: 'todo-list-o', route: '/tasks' },
  riskAssessment: { name: 'riskAssessment', label: 'Audit', icon: 'shield-o', activeIcon: 'shield-o', route: '/risk-assessment' },
  attendance: { name: 'attendance', label: 'Absensi', icon: 'clock-o', activeIcon: 'clock', route: '/attendance' },
  profile: { name: 'profile', label: 'Profil', icon: 'user-o', activeIcon: 'user-o', route: '/profile' },
  schoolMonitoring: { name: 'schoolMonitoring', label: 'Monitoring', icon: 'eye-o', activeIcon: 'eye', route: '/school-monitoring' }
}

const ROLE_NAV_MAP = {
  driver: [
    { left: ['tasks'], center: 'attendance', right: ['profile'] }
  ],
  asisten_lapangan: [
    { left: ['tasks'], center: 'attendance', right: ['profile'] }
  ],
  kepala_sppg: [
    { left: ['dashboard', 'monitoring'], center: 'attendance', right: ['menu', 'profile'] }
  ],
  kepala_yayasan: [
    { left: ['dashboardYayasan', 'riskAssessment'], center: 'attendance', right: ['profile'] }
  ],
  ahli_gizi: [
    { left: ['menu'], center: 'attendance', right: ['profile'] }
  ],
  sekolah: [
    { left: ['schoolMonitoring'], center: 'attendance', right: ['profile'] }
  ]
}

const DEFAULT_NAV = { left: [], center: 'attendance', right: ['profile'] }

const navConfig = computed(() => {
  const role = authStore.user?.role?.toLowerCase() || ''
  return ROLE_NAV_MAP[role]?.[0] || DEFAULT_NAV
})

const leftItems = computed(() => {
  return navConfig.value.left.map(key => NAV_CONFIG[key])
})

const centerItem = computed(() => {
  return NAV_CONFIG[navConfig.value.center]
})

const rightItems = computed(() => {
  return navConfig.value.right.map(key => NAV_CONFIG[key])
})

const allItems = computed(() => {
  return [...leftItems.value, centerItem.value, ...rightItems.value]
})

const activeItem = ref(null)

function findActiveItem() {
  const currentPath = route.path
  const item = allItems.value.find(item => currentPath.startsWith(item.route))
  return item || allItems.value[0]
}

function isActive(item) {
  return activeItem.value?.name === item?.name
}

function onNavClick(item) {
  if (item) {
    activeItem.value = item
    router.push(item.route)
  }
}

// Sync active item with current route
watch(
  () => route.path,
  () => {
    activeItem.value = findActiveItem()
  },
  { immediate: true }
)

onMounted(() => {
  activeItem.value = findActiveItem()
})
</script>

<style scoped>
.bottom-navigation-wrapper {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
}

.bottom-navigation {
  height: 56px;
  background: #FFFFFF;
  border-top: 1px solid #D8D8DB;
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  padding: 0;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 8px 12px;
}

.nav-item__label {
  font-size: 11px;
  font-weight: 500;
  color: var(--h-text-secondary);
  transition: color 0.2s ease;
}

.nav-item--active .nav-item__label {
  color: #C94A3A;
  font-weight: 600;
}

.nav-item--active :deep(.van-icon) {
  transform: scale(1.1);
}
</style>
