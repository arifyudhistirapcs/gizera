<template>
  <div class="bottom-navigation-wrapper">
    <div class="bottom-navigation">
      <!-- Left Items -->
      <div
        v-for="(item, index) in leftItems"
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

      <!-- Center Attendance Button (Floating) -->
      <div class="nav-center-spacer"></div>

      <!-- Right Items -->
      <div
        v-for="(item, index) in rightItems"
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

    <!-- Floating Center Button -->
    <div
      :class="['floating-center-btn', { 'floating-center-btn--active': isActive(centerItem) }]"
      @click="onNavClick(centerItem)"
    >
      <van-icon
        :name="isActive(centerItem) ? centerItem.activeIcon : centerItem.icon"
        color="#FFFFFF"
        size="28"
      />
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

const activeColor = '#5A4372'
const inactiveColor = '#ACA9B0'

const NAV_CONFIG = {
  dashboard: { name: 'dashboard', label: 'Home', icon: 'wap-home-o', activeIcon: 'wap-home', route: '/dashboard' },
  monitoring: { name: 'monitoring', label: 'Monitoring', icon: 'search', activeIcon: 'search', route: '/monitoring' },
  menu: { name: 'menu', label: 'Menu', icon: 'notes-o', activeIcon: 'notes-o', route: '/menu-planning' },
  tasks: { name: 'tasks', label: 'Tugas', icon: 'todo-list-o', activeIcon: 'todo-list-o', route: '/tasks' },
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
  pointer-events: none;
}

.bottom-navigation {
  height: 64px;
  background: #FFFFFF;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.08);
  border-radius: 24px 24px 0 0;
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 0 16px;
  pointer-events: auto;
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 8px 0;
}

.nav-item__label {
  font-size: 11px;
  font-weight: 500;
  color: var(--h-text-secondary);
  transition: color 0.2s ease;
}

.nav-item--active .nav-item__label {
  color: #5A4372;
  font-weight: 600;
}

.nav-item--active :deep(.van-icon) {
  transform: scale(1.1);
}

.nav-center-spacer {
  width: 64px;
  flex-shrink: 0;
}

/* Floating Center Button */
.floating-center-btn {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #6B4E9B 0%, #5A4372 100%);
  box-shadow: 0 8px 24px rgba(90, 67, 114, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  pointer-events: auto;
}

.floating-center-btn:active {
  transform: translateX(-50%) scale(0.95);
}

.floating-center-btn--active {
  background: linear-gradient(135deg, #7B5EAB 0%, #6A5382 100%);
  box-shadow: 0 12px 32px rgba(90, 67, 114, 0.5);
}

/* Add pulse animation for center button */
@keyframes pulse {
  0% {
    box-shadow: 0 8px 24px rgba(90, 67, 114, 0.4);
  }
  50% {
    box-shadow: 0 8px 24px rgba(90, 67, 114, 0.6), 0 0 0 8px rgba(90, 67, 114, 0.1);
  }
  100% {
    box-shadow: 0 8px 24px rgba(90, 67, 114, 0.4);
  }
}

.floating-center-btn:hover {
  animation: pulse 2s infinite;
}
</style>
