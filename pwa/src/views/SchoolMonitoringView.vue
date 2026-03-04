<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="school-monitoring-page">
      <!-- NavBar -->
      <van-nav-bar title="Monitoring" />

      <!-- School Header -->
      <div class="school-header">
        <h2 class="school-header__name">{{ schoolName }}</h2>
        <p class="school-header__date">{{ formattedDate }}</p>
      </div>

      <!-- Loading State -->
      <template v-if="schoolStore.loading">
        <div class="section-block">
          <SkeletonCard :rows="3" />
        </div>
        <div class="section-block">
          <SkeletonCard :rows="3" />
        </div>
        <div class="section-block">
          <SkeletonCard :rows="3" />
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="schoolStore.error" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ schoolStore.error }}</p>
        <van-button type="primary" size="normal" @click="schoolStore.retry()">
          Coba Lagi
        </van-button>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- Today's Menu Card -->
        <div v-if="schoolStore.todayMenu" class="h-card">
          <h3 class="section-title">Menu Hari Ini</h3>
          <p class="menu-name">{{ schoolStore.todayMenu.menuName }}</p>
          <p class="menu-detail">
            <span class="menu-detail__label">Komponen:</span>
            {{ formatComponents(schoolStore.todayMenu.components) }}
          </p>
          <p class="menu-detail">
            <span class="menu-detail__label">Porsi:</span>
            {{ schoolStore.todayMenu.portions }}
          </p>
        </div>
        <div v-else class="h-card empty-card">
          <h3 class="section-title">Menu Hari Ini</h3>
          <p class="empty-card__text">Belum ada data menu untuk hari ini</p>
        </div>

        <!-- Delivery Status Card -->
        <div v-if="schoolStore.deliveryStatus" class="h-card">
          <h3 class="section-title">Status Pengiriman</h3>
          <div class="status-row">
            <van-tag :type="statusTagType(schoolStore.deliveryStatus.status)" round size="medium">
              {{ statusLabel(schoolStore.deliveryStatus.status) }}
            </van-tag>
          </div>
          <p class="status-detail">
            <van-icon name="manager-o" class="status-detail__icon" />
            Driver: {{ schoolStore.deliveryStatus.driverName }}
          </p>
          <p v-if="schoolStore.deliveryStatus.estimatedTime" class="status-detail">
            <van-icon name="clock-o" class="status-detail__icon" />
            Estimasi: {{ schoolStore.deliveryStatus.estimatedTime }}
          </p>
        </div>
        <div v-else class="h-card empty-card">
          <h3 class="section-title">Status Pengiriman</h3>
          <p class="empty-card__text">Belum ada data pengiriman untuk hari ini</p>
        </div>

        <!-- Pickup Status Card (conditional) -->
        <div v-if="schoolStore.pickupStatus" class="h-card">
          <h3 class="section-title">Status Pengambilan</h3>
          <div class="status-row">
            <van-tag :type="statusTagType(schoolStore.pickupStatus.status)" round size="medium">
              {{ statusLabel(schoolStore.pickupStatus.status) }}
            </van-tag>
          </div>
          <p class="status-detail">
            <van-icon name="manager-o" class="status-detail__icon" />
            Driver: {{ schoolStore.pickupStatus.driverName }}
          </p>
          <p v-if="schoolStore.pickupStatus.estimatedTime" class="status-detail">
            <van-icon name="clock-o" class="status-detail__icon" />
            Estimasi: {{ schoolStore.pickupStatus.estimatedTime }}
          </p>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useSchoolMonitoringStore } from '@/stores/schoolMonitoring'
import { useAuthStore } from '@/stores/auth'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

const schoolStore = useSchoolMonitoringStore()
const authStore = useAuthStore()
const refreshing = ref(false)

const schoolName = computed(() => authStore.user?.name || 'Sekolah')

const formattedDate = computed(() => {
  const now = new Date()
  const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
  return `${days[now.getDay()]}, ${now.getDate()} ${months[now.getMonth()]} ${now.getFullYear()}`
})

function formatComponents(components) {
  if (Array.isArray(components)) {
    return components.join(', ')
  }
  return components || '-'
}

function statusLabel(status) {
  const labels = {
    pending: 'Menunggu',
    in_progress: 'Dalam Perjalanan',
    arrived: 'Sampai',
    completed: 'Selesai'
  }
  return labels[status] || status
}

function statusTagType(status) {
  const types = {
    pending: 'warning',
    in_progress: 'primary',
    arrived: 'primary',
    completed: 'success'
  }
  return types[status] || 'default'
}

async function onRefresh() {
  await schoolStore.fetchSchoolData()
  refreshing.value = false
}

onMounted(() => {
  schoolStore.fetchSchoolData()
})
</script>

<style scoped>
.school-monitoring-page {
  padding: var(--h-spacing-lg);
  padding-bottom: 80px;
  min-height: 100vh;
}

/* School Header */
.school-header {
  margin-bottom: var(--h-spacing-lg);
}

.school-header__name {
  font-size: 20px;
  font-weight: 700;
  color: var(--h-text-primary);
  margin: 0 0 var(--h-spacing-xs) 0;
}

.school-header__date {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: 0;
}

/* Section title */
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 var(--h-spacing-md) 0;
}

.section-block {
  margin-bottom: var(--h-spacing-lg);
}

/* Menu Card */
.menu-name {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-primary);
  margin: 0 0 var(--h-spacing-sm) 0;
}

.menu-detail {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: 0 0 var(--h-spacing-xs) 0;
  line-height: 1.5;
}

.menu-detail__label {
  font-weight: 500;
  color: var(--h-text-primary);
}

/* Status Cards */
.status-row {
  margin-bottom: var(--h-spacing-md);
}

.status-detail {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: 0 0 var(--h-spacing-xs) 0;
  display: flex;
  align-items: center;
  gap: var(--h-spacing-sm);
  line-height: 1.5;
}

.status-detail__icon {
  color: var(--h-text-light);
  font-size: 16px;
  flex-shrink: 0;
}

/* Empty Card */
.empty-card__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: 0;
}

/* Error State */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px var(--h-spacing-xl);
  text-align: center;
}

.error-state__message {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: var(--h-spacing-lg) 0;
  line-height: 1.5;
}
</style>
