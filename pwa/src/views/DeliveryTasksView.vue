<template>
  <div class="delivery-tasks-container">
    <!-- Navigation Bar -->
    <van-nav-bar title="Tugas Pengiriman dan Pengambilan" fixed>
      <template #right>
        <van-icon 
          name="refresh" 
          @click="refreshTasks" 
          :class="{ 'rotating': isRefreshing }"
        />
      </template>
    </van-nav-bar>

    <!-- Access Denied -->
    <van-empty 
      v-if="!hasAccess"
      image="error"
      description="Anda tidak memiliki akses ke halaman ini"
    >
      <van-button type="primary" @click="goToProfile">Kembali ke Profil</van-button>
    </van-empty>

    <!-- Offline Indicator -->
    <van-notice-bar 
      v-else-if="!isOnline" 
      type="warning" 
      text="Mode offline - Data mungkin tidak terbaru"
      left-icon="warning-o"
    />

    <div v-else class="delivery-tasks-content">
      <!-- Tab Selection -->
      <div class="tab-selection">
        <button
          :class="['tab-button', { 'tab-button--active': activeTab === 'delivery' }]"
          @click="activeTab = 'delivery'"
        >
          Pengiriman
        </button>
        <button
          :class="['tab-button', { 'tab-button--active': activeTab === 'pickup' }]"
          @click="activeTab = 'pickup'"
        >
          Pengambilan
        </button>
      </div>

      <!-- Date Filter -->
      <div class="date-filter">
        <van-field
          v-model="selectedDateFormatted"
          label="Tanggal"
          placeholder="Pilih tanggal"
          readonly
          right-icon="calendar-o"
          @click="showDatePicker = true"
        />
      </div>

      <!-- Date Picker Popup -->
      <van-popup v-model:show="showDatePicker" position="bottom">
        <van-date-picker
          v-model="selectedDate"
          :max-date="today"
          @confirm="onDateConfirm"
          @cancel="showDatePicker = false"
        />
      </van-popup>

      <!-- Loading State -->
      <van-loading v-if="isLoading" type="spinner" vertical>
        Memuat tugas...
      </van-loading>

      <!-- Empty State -->
      <van-empty 
        v-else-if="!isLoading && filteredTasks.length === 0"
        image="search"
        :description="`Tidak ada tugas ${activeTab === 'delivery' ? 'pengiriman' : 'pengambilan'} untuk tanggal yang dipilih`"
      />

      <!-- Tasks List -->
      <div v-else class="tasks-list">
      <van-card
        v-for="task in sortedFilteredTasks"
        :key="task.id"
        :title="task.school?.name || 'Sekolah tidak diketahui'"
        :desc="formatAddress(task.school?.address)"
        class="task-card"
        @click="showTaskDetail(task)"
      >
        <template #tags>
          <van-tag 
            :type="task.task_type === 'pickup' ? 'warning' : 'success'" 
            size="medium"
            class="task-type-tag"
          >
            {{ getTaskTypeText(task) }}
          </van-tag>
          <van-tag 
            :type="getStatusType(task.status, task)" 
            size="medium"
          >
            {{ getStatusText(task.status, task) }}
          </van-tag>
          <van-tag 
            type="primary" 
            size="medium" 
            class="route-tag"
          >
            Urutan: {{ task.route_order }}
          </van-tag>
        </template>

        <template #footer>
          <div class="task-info">
            <div class="info-row">
              <van-icon name="location-o" />
              <span class="info-text">
                {{ task.school?.latitude?.toFixed(6) }}, {{ task.school?.longitude?.toFixed(6) }}
              </span>
              <van-button 
                size="mini" 
                type="primary" 
                @click.stop="openMaps(task.school)"
                icon="guide-o"
              >
                Navigasi
              </van-button>
            </div>
            
            <div class="info-row">
              <van-icon name="friends-o" />
              <span class="info-text" v-if="task.task_type === 'pickup'">
                {{ task.ompreng_count || 0 }} wadah ompreng
              </span>
              <span class="info-text" v-else>{{ task.portions }} porsi</span>
            </div>

            <div class="info-row" v-if="task.menu_items && task.menu_items.length > 0">
              <van-icon name="shop-o" />
              <span class="info-text">
                {{ task.menu_items.map(item => item.recipe?.name).join(', ') }}
              </span>
            </div>
          </div>
        </template>
      </van-card>
    </div>
    </div>

    <!-- Bottom Navigation -->
    <van-tabbar v-model="active" route fixed>
      <van-tabbar-item v-if="hasAccess" to="/tasks" icon="orders-o">Tugas</van-tabbar-item>
      <van-tabbar-item to="/attendance" icon="clock-o">Absensi</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">Profil</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDeliveryTasksStore } from '@/stores/deliveryTasks'
import { showToast, showConfirmDialog } from 'vant'

const router = useRouter()

const authStore = useAuthStore()
const deliveryTasksStore = useDeliveryTasksStore()

// Reactive data
const active = ref(0)
const isLoading = ref(false)
const isRefreshing = ref(false)
const isOnline = ref(navigator.onLine)
const activeTab = ref('delivery')
const selectedDate = ref(new Date())
const showDatePicker = ref(false)
const today = new Date()

// Computed properties
const tasks = computed(() => deliveryTasksStore.tasks)

// Filter tasks by active tab
const filteredTasks = computed(() => {
  return tasks.value.filter(task => task.task_type === activeTab.value)
})

// Sort filtered tasks by route_order
const sortedFilteredTasks = computed(() => {
  return [...filteredTasks.value].sort((a, b) => a.route_order - b.route_order)
})

const sortedTasks = computed(() => {
  return [...tasks.value].sort((a, b) => a.route_order - b.route_order)
})

// Format selected date for display
const selectedDateFormatted = computed(() => {
  return selectedDate.value.toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

// Check if user has access (driver or asisten_lapangan only)
const hasAccess = computed(() => {
  const allowedRoles = ['driver', 'asisten_lapangan']
  const userRole = authStore.user?.role?.toLowerCase()
  return allowedRoles.includes(userRole)
})

// Methods
const goToProfile = () => {
  router.push('/profile')
}

const onDateConfirm = () => {
  showDatePicker.value = false
  // TODO: Fetch tasks for selected date
  // For now, we only show today's tasks from store
  showToast('Filter tanggal akan diterapkan')
}

const loadTasks = async () => {
  if (!authStore.user?.id || !hasAccess.value) {
    console.log('[DeliveryTasksView] Cannot load tasks - user:', authStore.user?.id, 'hasAccess:', hasAccess.value)
    return
  }
  
  console.log('[DeliveryTasksView] Loading tasks for user:', authStore.user?.id, authStore.user?.full_name, authStore.user?.role)
  isLoading.value = true
  try {
    await deliveryTasksStore.fetchTodayTasks(authStore.user.id)
    console.log('[DeliveryTasksView] Tasks loaded:', tasks.value.length)
  } catch (error) {
    console.error('Error loading tasks:', error)
    showToast('Gagal memuat tugas pengiriman')
  } finally {
    isLoading.value = false
  }
}

const refreshTasks = async () => {
  if (!authStore.user?.id) return
  
  isRefreshing.value = true
  try {
    await deliveryTasksStore.fetchTodayTasks(authStore.user.id, true) // force refresh
    showToast('Data berhasil diperbarui')
  } catch (error) {
    console.error('Error refreshing tasks:', error)
    showToast('Gagal memperbarui data')
  } finally {
    isRefreshing.value = false
  }
}

const showTaskDetail = (task) => {
  if (task.task_type === 'pickup') {
    // For pickup tasks, navigate to pickup task detail
    router.push(`/pickup-tasks/${task.pickup_task_id}?record=${task.delivery_record_id}`)
  } else {
    router.push(`/tasks/${task.id}`)
  }
}

const formatAddress = (address) => {
  if (!address) return 'Alamat tidak tersedia'
  return address.length > 50 ? address.substring(0, 50) + '...' : address
}

const getStatusType = (status, task) => {
  // For pickup tasks, use stage-based status
  if (task?.task_type === 'pickup') {
    const stage = task.current_stage
    if (stage >= 13) return 'success'
    if (stage >= 10) return 'primary'
    return 'warning'
  }
  
  const statusTypes = {
    'pending': 'warning',
    'in_progress': 'primary',
    'arrived': 'success',
    'received': 'success',
    'completed': 'success',
    'cancelled': 'danger'
  }
  return statusTypes[status] || 'default'
}

const getStatusText = (status, task) => {
  // For pickup tasks, show stage-based status
  if (task?.task_type === 'pickup') {
    const stage = task.current_stage
    const stageTexts = {
      10: 'Menuju Lokasi',
      11: 'Tiba di Lokasi',
      12: 'Kembali ke SPPG',
      13: 'Tiba di SPPG'
    }
    return stageTexts[stage] || `Stage ${stage}`
  }
  
  const statusTexts = {
    'pending': 'Menunggu',
    'in_progress': 'Dalam Perjalanan',
    'arrived': 'Sudah Sampai',
    'received': 'Sudah Diterima',
    'completed': 'Sudah Diterima',
    'cancelled': 'Dibatalkan'
  }
  return statusTexts[status] || status
}

const getTaskTypeText = (task) => {
  // Tentukan jenis tugas berdasarkan task_type atau current_stage
  if (task.task_type === 'pickup') {
    return '📦 Pengambilan Ompreng'
  }
  // Stage 1-4, 9: Pengiriman Makanan (siap dikirim → sudah diterima)
  // Stage 5-8: Pengambilan Ompreng (driver menuju lokasi pengambilan → driver tiba di SPPG)
  const stage = task.current_stage || 1
  if (stage >= 5 && stage <= 8) {
    return '📦 Pengambilan Ompreng'
  }
  return '🍱 Pengiriman Makanan'
}

const openMaps = (school) => {
  if (!school?.latitude || !school?.longitude) {
    showToast('Koordinat GPS tidak tersedia')
    return
  }
  
  const url = `https://www.google.com/maps/dir/?api=1&destination=${school.latitude},${school.longitude}`
  window.open(url, '_blank')
}

// Network status handlers
const handleOnline = () => {
  isOnline.value = true
  showToast('Koneksi internet tersambung')
  // Sync offline data when back online
  deliveryTasksStore.syncAllOfflineData()
}

const handleOffline = () => {
  isOnline.value = false
  showToast('Mode offline - Data akan disinkronkan saat online')
}

// Lifecycle
onMounted(() => {
  loadTasks()
  
  // Listen for network status changes
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
})
</script>

<style scoped>
.delivery-tasks-container {
  min-height: 100vh;
  background-color: #E8EDE5;
  padding-top: 46px;
  padding-bottom: 50px;
}

.delivery-tasks-content {
  padding: 16px;
}

/* Tab Selection */
.tab-selection {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.tab-button {
  flex: 1;
  padding: 10px 20px;
  border-radius: 20px;
  border: none;
  font-family: var(--van-base-font, -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, Segoe UI, Arial, Roboto, 'PingFang SC', 'miui', 'Hiragino Sans GB', 'Microsoft Yahei', sans-serif);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #FFFFFF;
  color: #6B6B6B;
  box-shadow: 0px 2px 8px rgba(0, 0, 0, 0.08);
}

.tab-button--active {
  background: #303030;
  color: #FFFFFF;
  box-shadow: none;
}

/* Date Filter */
.date-filter {
  margin-bottom: 16px;
  background: #FFFFFF;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0px 2px 8px rgba(0, 0, 0, 0.08);
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-card {
  margin-bottom: 16px;
  border-radius: 16px !important;
  box-shadow: 0px 18px 40px rgba(112, 144, 176, 0.12) !important;
  overflow: hidden;
  transition: all 0.3s ease;
}

.task-card:active {
  transform: scale(0.98);
}

.task-type-tag {
  margin-right: 6px;
}

.task-type-tag[type="success"] {
  background-color: #e8f5e9 !important;
  color: #2e7d32 !important;
}

.task-type-tag[type="warning"] {
  background-color: #fff3e0 !important;
  color: #e65100 !important;
}

.route-tag {
  margin-left: 8px;
}

.task-info {
  padding: 12px 0;
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
  color: #6B6B6B;
}

.info-row .van-icon {
  margin-right: 8px;
  color: #303030;
  font-size: 18px;
}

.info-text {
  flex: 1;
  margin-right: 8px;
  color: #303030;
  font-weight: 500;
}

.rotating {
  animation: rotate 1s linear infinite;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* Responsive adjustments */
@media (max-width: 375px) {
  .info-row {
    font-size: 12px;
  }
  
  .task-card {
    margin-bottom: 12px;
  }
  
  .tasks-list {
    padding: 12px;
  }
}
</style>