<template>
  <header class="h-header" :class="{ 'mobile': isMobile }">
    <!-- Mobile: Hamburger Menu Button -->
    <button
      v-if="isMobile"
      class="hamburger-button"
      @click="$emit('toggle-drawer')"
      aria-label="Toggle menu"
      type="button"
    >
      <MenuOutlined />
    </button>

    <!-- Breadcrumb & Title Section -->
    <div class="header-left">
      <!-- Breadcrumb (hidden on mobile) -->
      <div v-if="!isMobile && breadcrumb" class="breadcrumb">
        <span class="breadcrumb-parent">{{ breadcrumb.parent }}</span>
        <span class="breadcrumb-separator">/</span>
        <span class="breadcrumb-current">{{ breadcrumb.current }}</span>
      </div>
      
      <!-- Page Title -->
      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>

    <!-- Right Section -->
    <div class="header-right">
      <!-- Dark mode toggle disabled -->
    </div>
  </header>
</template>

<script setup>
import { computed } from 'vue'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { MenuOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  pageTitle: {
    type: String,
    default: 'Dashboard'
  },
  breadcrumb: {
    type: Object,
    default: () => ({ parent: 'Pages', current: 'Dashboard' })
  }
})

defineEmits(['toggle-drawer'])

const { isMobile } = useBreakpoint()
</script>

<style scoped>
.h-header {
  display: flex;
  align-items: center;
  height: 56px;
  background-color: transparent;
  padding: 0 4px;
  border-bottom: none;
  position: relative;
  transition: all 0.2s ease;
}

.h-header.mobile {
  height: 52px;
  padding: 0 16px;
  background-color: #fff;
  border-bottom: 1px solid #D8D8DB;
}

.dark .h-header {
  background-color: transparent;
  border-bottom-color: transparent;
}

.dark .h-header.mobile {
  background-color: var(--h-bg-secondary-dark, #252525);
  border-bottom-color: #404040;
}

.hamburger-button {
  width: 44px;
  height: 44px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  color: var(--h-text-primary, #322837);
  border-radius: 8px;
  transition: all 0.2s ease;
  margin-right: 12px;
}

.hamburger-button:hover {
  background-color: var(--h-bg-light, #F0F0F0);
}

.hamburger-button:active {
  transform: scale(0.95);
}

.dark .hamburger-button {
  color: var(--h-text-primary-dark, #F7F8FA);
}

.dark .hamburger-button:hover {
  background-color: var(--h-bg-card-dark, #303030);
}

.header-left {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6B6B6B;
}

.breadcrumb-parent {
  color: #6B6B6B;
}

.breadcrumb-separator {
  color: #6B6B6B;
}

.breadcrumb-current {
  color: #303030;
  font-weight: 500;
}

.dark .breadcrumb-current {
  color: var(--h-text-primary-dark, #F7F8FA);
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #303030;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: -0.3px;
}

.mobile .page-title {
  font-size: 18px;
}

.dark .page-title {
  color: var(--h-text-primary-dark, #F7F8FA);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

@media (max-width: 768px) {
  .header-right {
    gap: 8px;
  }
}
</style>