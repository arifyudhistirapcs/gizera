<template>
  <div class="horizon-layout" :class="{ 'dark': isDark }">
    <!-- Desktop: Fixed Sidebar -->
    <HSidebar
      v-if="isDesktop"
      v-model:collapsed="sidebarCollapsed"
    />

    <!-- Mobile: Drawer Sidebar -->
    <MobileDrawer v-model="drawerOpen">
      <HSidebar :collapsed="false" />
    </MobileDrawer>

    <!-- Main Content Area -->
    <div
      class="horizon-main"
      :class="{ 'sidebar-collapsed': sidebarCollapsed && isDesktop }"
    >
      <!-- Header -->
      <HHeader
        :page-title="pageTitle"
        :breadcrumb="breadcrumb"
        :notification-count="notificationCount"
        @toggle-drawer="toggleDrawer"
        @search="handleSearch"
        @notification-click="handleNotificationClick"
      />

      <!-- Content Area -->
      <main class="horizon-content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useDarkMode } from '@/composables/useDarkMode'
import { useBreakpoint } from '@/composables/useBreakpoint'
import HSidebar from '@/components/layout/HSidebar.vue'
import HHeader from '@/components/layout/HHeader.vue'
import MobileDrawer from '@/components/layout/MobileDrawer.vue'

/**
 * HorizonLayout Component
 * 
 * Main layout wrapper yang menggabungkan HSidebar, HHeader, dan content area:
 * - Desktop: sidebar always visible, content area adjusts for sidebar width
 * - Mobile: sidebar hidden, opens via MobileDrawer, hamburger button in header
 * - Content area: bg #F4F7FE (light) / #0B1437 (dark), padding 24px, gap 20px
 * - Mobile: full width content, reduced padding 16px
 * - Integrates useDarkMode + useBreakpoint composables
 * - Smooth transitions
 */

const props = defineProps({
  /**
   * Page title untuk header
   */
  pageTitle: {
    type: String,
    default: 'Dashboard'
  },
  
  /**
   * Breadcrumb object { parent: 'POSe', current: 'Dashboard' }
   */
  breadcrumb: {
    type: Object,
    default: () => ({ parent: 'POSe', current: 'Dashboard' })
  },
  
  /**
   * Notification count untuk badge
   */
  notificationCount: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['search', 'notification-click'])

// Composables
const { isDark } = useDarkMode()
const { isDesktop, isMobile } = useBreakpoint()

// Sidebar state
const sidebarCollapsed = ref(false)
const drawerOpen = ref(false)

/**
 * Toggle mobile drawer
 */
const toggleDrawer = () => {
  drawerOpen.value = !drawerOpen.value
}

/**
 * Handle search dari header
 */
const handleSearch = (query) => {
  emit('search', query)
}

/**
 * Handle notification click dari header
 */
const handleNotificationClick = () => {
  emit('notification-click')
}
</script>

<style scoped>
.horizon-layout {
  display: flex;
  min-height: 100vh;
  background-color: var(--h-bg-primary, #F7F8FA);
  padding: 12px;
  transition: background-color 0.2s ease;
}

/* Dark mode background */
.horizon-layout.dark {
  background-color: var(--h-bg-primary);
}

/* Main Content Area */
.horizon-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  margin-left: 240px; /* sidebar width (220px) + gap */
  transition: margin-left 300ms ease;
}

.horizon-main.sidebar-collapsed {
  margin-left: 92px; /* collapsed (72px) + gap */
}

@media (max-width: 1024px) {
  .horizon-layout {
    padding: 0;
  }

  .horizon-main {
    margin-left: 0;
  }
}

/* Content Area */
.horizon-content {
  flex: 1;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 0;
  overflow-x: hidden;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

@media (max-width: 768px) {
  .horizon-content {
    padding: 16px;
  }
}

/* Smooth transitions for theme switching */
* {
  transition: background-color 0.2s ease, color 0.2s ease, border-color 0.2s ease;
}
</style>
