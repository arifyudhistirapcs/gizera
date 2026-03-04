<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="dashboard-page">
      <!-- NavBar -->
      <van-nav-bar title="Dashboard" />

      <!-- Loading State -->
      <template v-if="dashboardStore.loading">
        <div class="summary-grid">
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
        </div>
        <div class="section-block">
          <SkeletonCard :rows="4" />
        </div>
        <div class="section-block">
          <SkeletonCard :rows="3" />
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="dashboardStore.error" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ dashboardStore.error }}</p>
        <van-button type="primary" size="normal" @click="dashboardStore.retry()">
          Coba Lagi
        </van-button>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- 2x2 Summary Grid -->
        <div class="summary-grid">
          <SummaryCard
            icon="friends-o"
            iconColor="#4CAF50"
            label="Hadir"
            :value="dashboardStore.summary.totalHadir"
          />
          <SummaryCard
            icon="logistics"
            iconColor="#5A4372"
            label="Pengiriman"
            :value="dashboardStore.summary.totalPengiriman"
          />
          <SummaryCard
            icon="success"
            iconColor="#2196F3"
            label="Selesai"
            :value="dashboardStore.summary.totalSelesai"
          />
          <SummaryCard
            icon="shop-o"
            iconColor="#FF9800"
            label="Sekolah"
            :value="dashboardStore.summary.totalSekolah"
          />
        </div>

        <!-- Attendance Chart (7-day CSS bars) -->
        <div class="chart-card h-card">
          <h3 class="section-title">Kehadiran 7 Hari Terakhir</h3>
          <div class="chart-bars">
            <div
              v-for="(day, index) in dashboardStore.attendanceChart"
              :key="index"
              class="chart-bar-col"
            >
              <span class="chart-bar-value">{{ day.count }}</span>
              <div class="chart-bar-track">
                <div
                  class="chart-bar-fill"
                  :style="{ height: barHeight(day.count) }"
                />
              </div>
              <span class="chart-bar-label">{{ getDayLabel(day.date) }}</span>
            </div>
          </div>
        </div>

        <!-- Recent Tasks -->
        <div class="recent-tasks-card h-card">
          <h3 class="section-title">Tugas Terbaru</h3>
          <div v-if="dashboardStore.recentTasks.length === 0" class="empty-state">
            <p class="empty-state__text">Belum ada tugas terbaru</p>
          </div>
          <div
            v-for="task in dashboardStore.recentTasks"
            :key="task.id"
            class="task-item"
          >
            <div class="task-item__info">
              <span class="task-item__school">{{ task.schoolName }}</span>
              <span class="task-item__type">{{ taskTypeLabel(task.taskType) }}</span>
            </div>
            <van-tag
              :type="statusTagType(task.status)"
              round
              size="medium"
            >
              {{ statusLabel(task.status) }}
            </van-tag>
          </div>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import SummaryCard from '@/components/mobile/SummaryCard.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

const dashboardStore = useDashboardStore()
const refreshing = ref(false)

const maxCount = computed(() => {
  if (!dashboardStore.attendanceChart.length) return 1
  const max = Math.max(...dashboardStore.attendanceChart.map(d => d.count))
  return max || 1
})

function barHeight(count) {
  const percentage = (count / maxCount.value) * 100
  return `${Math.max(percentage, 4)}%`
}

const dayNames = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']

function getDayLabel(dateStr) {
  const date = new Date(dateStr)
  return dayNames[date.getDay()] || ''
}

function taskTypeLabel(type) {
  return type === 'delivery' ? 'Pengiriman' : 'Pengambilan'
}

function statusLabel(status) {
  const labels = {
    pending: 'Menunggu',
    in_progress: 'Dalam Perjalanan',
    completed: 'Selesai'
  }
  return labels[status] || status
}

function statusTagType(status) {
  const types = {
    pending: 'warning',
    in_progress: 'primary',
    completed: 'success'
  }
  return types[status] || 'default'
}

async function onRefresh() {
  await dashboardStore.fetchDashboardData()
  refreshing.value = false
}

onMounted(() => {
  dashboardStore.fetchDashboardData()
})
</script>

<style scoped>
.dashboard-page {
  padding: var(--h-spacing-lg);
  padding-bottom: 80px;
  min-height: 100vh;
}

/* 2x2 Summary Grid */
.summary-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--h-spacing-md);
  margin-bottom: var(--h-spacing-lg);
}

/* Section title */
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 var(--h-spacing-lg) 0;
}

.section-block {
  margin-bottom: var(--h-spacing-lg);
}

/* Chart Card */
.chart-card {
  margin-bottom: var(--h-spacing-lg);
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--h-spacing-sm);
  height: 140px;
}

.chart-bar-col {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  height: 100%;
}

.chart-bar-value {
  font-size: 11px;
  font-weight: 600;
  color: var(--h-text-secondary);
  margin-bottom: var(--h-spacing-xs);
}

.chart-bar-track {
  flex: 1;
  width: 100%;
  max-width: 32px;
  background: var(--h-bg-light);
  border-radius: var(--h-radius-sm);
  display: flex;
  align-items: flex-end;
  overflow: hidden;
}

.chart-bar-fill {
  width: 100%;
  background: linear-gradient(180deg, var(--h-primary) 0%, var(--h-accent) 100%);
  border-radius: var(--h-radius-sm);
  transition: height var(--h-transition-base);
  min-height: 4px;
}

.chart-bar-label {
  font-size: 11px;
  color: var(--h-text-secondary);
  margin-top: var(--h-spacing-xs);
}

/* Recent Tasks Card */
.recent-tasks-card {
  margin-bottom: var(--h-spacing-lg);
}

.task-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--h-spacing-md) 0;
  border-bottom: 1px solid var(--h-border-light);
}

.task-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.task-item__info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
  margin-right: var(--h-spacing-sm);
}

.task-item__school {
  font-size: 14px;
  font-weight: 500;
  color: var(--h-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.task-item__type {
  font-size: 12px;
  color: var(--h-text-secondary);
  margin-top: 2px;
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
  padding: var(--h-spacing-xl) 0;
  text-align: center;
}

.empty-state__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: 0;
}
</style>
