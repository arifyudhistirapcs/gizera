<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="monitoring-page">
      <!-- NavBar -->
      <van-nav-bar title="Monitoring" />

      <!-- Date Selector -->
      <div class="date-section">
        <DateSelector
          v-model="monitoringStore.selectedDate"
          @update:modelValue="onDateChange"
        />
      </div>

      <!-- Filter Chips -->
      <div class="filter-chips">
        <van-tag
          v-for="filter in filters"
          :key="filter.key"
          :class="['filter-chip', { 'filter-chip--active': monitoringStore.filterType === filter.key }]"
          :plain="monitoringStore.filterType !== filter.key"
          round
          size="large"
          @click="onFilterChange(filter.key)"
        >
          {{ filter.label }}
        </van-tag>
      </div>

      <!-- Loading State (initial load) -->
      <template v-if="monitoringStore.loading && !monitoringStore.activities.length">
        <div class="activity-list">
          <SkeletonCard v-for="n in 4" :key="n" :rows="2" :avatar="true" />
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="monitoringStore.error && !monitoringStore.activities.length" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ monitoringStore.error }}</p>
        <van-button type="primary" size="normal" @click="monitoringStore.retry()">
          Coba Lagi
        </van-button>
      </div>

      <!-- Activity List with Infinite Scroll -->
      <template v-else>
        <van-list
          v-model:loading="listLoading"
          :finished="!monitoringStore.hasMore"
          finished-text="Tidak ada aktivitas lagi"
          @load="onLoadMore"
          class="activity-list"
        >
          <ActivityLogItem
            v-for="activity in monitoringStore.activities"
            :key="activity.id"
            :employeeName="activity.employeeName"
            :activityType="activity.activityType"
            :timestamp="activity.timestamp"
            :status="activity.status"
          />
        </van-list>

        <!-- Empty State -->
        <div v-if="!monitoringStore.loading && !monitoringStore.activities.length" class="empty-state">
          <van-icon name="search" size="48" color="var(--h-text-light)" />
          <p class="empty-state__text">Tidak ada aktivitas ditemukan</p>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMonitoringStore } from '@/stores/monitoring'
import DateSelector from '@/components/mobile/DateSelector.vue'
import ActivityLogItem from '@/components/mobile/ActivityLogItem.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

const monitoringStore = useMonitoringStore()
const refreshing = ref(false)
const listLoading = ref(false)

const filters = [
  { key: 'all', label: 'Semua' },
  { key: 'attendance', label: 'Absensi' },
  { key: 'delivery', label: 'Pengiriman' },
  { key: 'pickup', label: 'Pengambilan' }
]

function onDateChange(date) {
  monitoringStore.setDate(date)
}

function onFilterChange(type) {
  monitoringStore.setFilter(type)
}

async function onLoadMore() {
  await monitoringStore.loadMore()
  listLoading.value = false
}

async function onRefresh() {
  await monitoringStore.fetchActivities()
  refreshing.value = false
}

onMounted(() => {
  monitoringStore.fetchActivities()
})
</script>

<style scoped>
.monitoring-page {
  padding: var(--h-spacing-lg);
  padding-bottom: 80px;
  min-height: 100vh;
}

/* Date Section */
.date-section {
  margin-bottom: var(--h-spacing-md);
}

/* Filter Chips */
.filter-chips {
  display: flex;
  gap: var(--h-spacing-sm);
  margin-bottom: var(--h-spacing-lg);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  padding-bottom: 2px;
}

.filter-chips::-webkit-scrollbar {
  display: none;
}

.filter-chip {
  cursor: pointer;
  flex-shrink: 0;
  min-height: 44px;
  padding: 0 16px;
  font-size: 13px;
  font-weight: 500;
  transition: all var(--h-transition-base);
  border-color: var(--h-primary) !important;
  color: var(--h-primary) !important;
  background: transparent !important;
}

.filter-chip--active {
  background: var(--h-primary) !important;
  color: #ffffff !important;
  border-color: var(--h-primary) !important;
}

/* Activity List */
.activity-list {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-md);
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

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px var(--h-spacing-xl);
  text-align: center;
}

.empty-state__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: var(--h-spacing-md) 0 0 0;
}
</style>
