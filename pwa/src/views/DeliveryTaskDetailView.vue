<template>
  <div class="detail-view">
    <!-- Navigation Bar -->
    <van-nav-bar
      title="Detail Tugas Pengiriman"
      left-arrow
      fixed
      @click-left="goBack"
    >
      <template #right>
        <div class="detail-view__nav-actions">
          <SyncStatusIndicator />
          <van-icon
            name="replay"
            size="20"
            color="#ffffff"
            :class="{ 'detail-view__refresh--spinning': isRefreshing }"
            @click="refreshTask"
          />
        </div>
      </template>
    </van-nav-bar>

    <!-- Offline Indicator -->
    <van-notice-bar
      v-if="!isOnline"
      type="warning"
      text="Mode offline - Data mungkin tidak terbaru"
      left-icon="warning-o"
    />

    <!-- Loading State -->
    <div v-if="isLoading" class="detail-view__content">
      <SkeletonCard :rows="2" />
      <SkeletonCard :rows="3" />
      <SkeletonCard :rows="2" />
      <SkeletonCard :rows="2" />
    </div>

    <!-- Error State -->
    <van-empty
      v-else-if="!isLoading && !task"
      image="error"
      description="Tugas tidak ditemukan"
    >
      <van-button type="primary" @click="goBack">
        Kembali ke Daftar Tugas
      </van-button>
    </van-empty>

    <!-- Task Detail Content -->
    <div v-else-if="task" class="detail-view__content">
      <!-- Progress Steps -->
      <div class="detail-view__card">
        <van-steps :active="activeStep" active-color="#303030">
          <van-step>Assign</van-step>
          <van-step>Jalan</van-step>
          <van-step>Sampai</van-step>
          <van-step>Selesai</van-step>
        </van-steps>
      </div>

      <!-- Status Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Status</div>
        <div class="detail-view__status-row">
          <van-tag
            :type="getStatusType(task.status)"
            size="large"
            class="detail-view__status-tag"
          >
            {{ getStatusText(task.status) }}
          </van-tag>
          <van-tag
            :type="getTaskTypeTagType(task)"
            size="medium"
          >
            {{ getTaskTypeText(task) }}
          </van-tag>
        </div>
        <div class="detail-view__meta-row">
          <van-icon name="location" color="var(--h-primary)" />
          <span>Urutan Rute: {{ task.route_order }}</span>
        </div>
      </div>

      <!-- School Info Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Info Sekolah</div>
        <div class="detail-view__info-row">
          <van-icon name="shop-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Nama Sekolah</span>
            <span class="detail-view__info-value">{{ task.school?.name || 'Tidak tersedia' }}</span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="location-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Alamat</span>
            <span class="detail-view__info-value detail-view__info-value--link" @click="showFullAddress">
              {{ task.school?.address || 'Tidak tersedia' }}
            </span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="contact" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Kontak Person</span>
            <span class="detail-view__info-value">{{ task.school?.contact_person || 'Tidak tersedia' }}</span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="phone-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Nomor Telepon</span>
            <span class="detail-view__info-value detail-view__info-value--link" @click="callSchool">
              {{ task.school?.phone_number || 'Tidak tersedia' }}
            </span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="friends-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Jumlah Siswa</span>
            <span class="detail-view__info-value">{{ task.school?.student_count?.toString() || 'Tidak tersedia' }}</span>
          </div>
        </div>

        <!-- Google Maps Navigation Button -->
        <button
          class="detail-view__maps-btn"
          :disabled="!hasValidGPS"
          @click="openGPSNavigation"
        >
          <van-icon name="guide-o" size="20" />
          <span>Buka di Google Maps</span>
          <van-icon name="arrow" size="16" />
        </button>
      </div>

      <!-- Menu / Portion Info Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Info Menu & Porsi</div>
        <div class="detail-view__info-row">
          <van-icon name="shopping-cart-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Jumlah Porsi</span>
            <span class="detail-view__info-value">{{ task.portions?.toString() || '0' }} porsi</span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="calendar-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Tanggal Pengiriman</span>
            <span class="detail-view__info-value">{{ formatDate(task.task_date) }}</span>
          </div>
        </div>
        <template v-if="task.menu_items && task.menu_items.length > 0">
          <div class="detail-view__divider" />
          <div
            v-for="item in task.menu_items"
            :key="item.id"
            class="detail-view__info-row"
          >
            <van-icon name="shop-o" color="var(--h-success)" />
            <div class="detail-view__info-text">
              <span class="detail-view__info-label">{{ item.recipe?.name || 'Menu tidak diketahui' }}</span>
              <span class="detail-view__info-value">{{ item.portions }} porsi</span>
            </div>
          </div>
        </template>
      </div>

      <!-- Action Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Aksi</div>

        <van-button
          v-if="task.status === 'pending'"
          type="success"
          size="large"
          block
          icon="play-circle-o"
          :loading="isUpdatingStatus"
          class="detail-view__action-btn"
          @click="startDelivery"
        >
          Mulai Pengiriman
        </van-button>

        <van-button
          v-if="task.status === 'in_progress'"
          type="success"
          size="large"
          block
          icon="location-o"
          :loading="isUpdatingStatus"
          class="detail-view__action-btn"
          @click="arrivedAtSchool"
        >
          Sudah Sampai Sekolah
        </van-button>

        <van-button
          v-if="task.status === 'arrived'"
          type="warning"
          size="large"
          block
          icon="edit"
          :loading="isUpdatingStatus"
          class="detail-view__action-btn detail-view__action-btn--epod"
          @click="openePODForm"
        >
          Buat Bukti Pengiriman (e-POD)
        </van-button>

        <van-button
          v-if="task.status === 'completed'"
          type="default"
          size="large"
          block
          disabled
          icon="success"
          class="detail-view__action-btn"
        >
          Pengiriman Selesai
        </van-button>
      </div>
    </div>

    <!-- Full Address Dialog -->
    <van-dialog
      v-model:show="showAddressDialog"
      title="Alamat Lengkap"
      :message="task?.school?.address"
      show-cancel-button
      cancel-button-text="Tutup"
      confirm-button-text="Salin"
      @confirm="copyAddress"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDeliveryTasksStore } from '@/stores/deliveryTasks'
import SyncStatusIndicator from '@/components/SyncStatusIndicator.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'
import { showToast, showConfirmDialog, showSuccessToast } from 'vant'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const deliveryTasksStore = useDeliveryTasksStore()

// Reactive data
const isLoading = ref(false)
const isRefreshing = ref(false)
const isUpdatingStatus = ref(false)
const isOnline = ref(navigator.onLine)
const showAddressDialog = ref(false)
const task = ref(null)

// Computed: progress step based on status
const activeStep = computed(() => {
  if (!task.value) return 0
  const statusMap = {
    'pending': 0,
    'in_progress': 1,
    'arrived': 2,
    'received': 3,
    'completed': 3,
    'cancelled': 0
  }
  return statusMap[task.value.status] ?? 0
})

const hasValidGPS = computed(() => {
  return task.value?.school?.latitude &&
         task.value?.school?.longitude &&
         Math.abs(task.value.school.latitude) <= 90 &&
         Math.abs(task.value.school.longitude) <= 180
})

// Methods
const loadTask = async () => {
  const taskId = route.params.id
  if (!taskId) {
    showToast('ID tugas tidak valid')
    goBack()
    return
  }

  isLoading.value = true
  try {
    let foundTask = deliveryTasksStore.getTaskById(parseInt(taskId))

    if (!foundTask) {
      await deliveryTasksStore.fetchTodayTasks(authStore.user.id)
      foundTask = deliveryTasksStore.getTaskById(parseInt(taskId))
    }

    if (foundTask) {
      task.value = foundTask
    } else {
      showToast('Tugas tidak ditemukan')
      goBack()
    }
  } catch (error) {
    console.error('Error loading task:', error)
    showToast('Gagal memuat detail tugas')
  } finally {
    isLoading.value = false
  }
}

const refreshTask = async () => {
  if (!authStore.user?.id) return

  isRefreshing.value = true
  try {
    await deliveryTasksStore.fetchTodayTasks(authStore.user.id, true)
    const taskId = route.params.id
    const updatedTask = deliveryTasksStore.getTaskById(parseInt(taskId))

    if (updatedTask) {
      task.value = updatedTask
      showToast('Data berhasil diperbarui')
    }
  } catch (error) {
    console.error('Error refreshing task:', error)
    showToast('Gagal memperbarui data')
  } finally {
    isRefreshing.value = false
  }
}

const goBack = async () => {
  // Refresh tasks before going back to ensure list shows updated status
  if (authStore.user?.id) {
    try {
      await deliveryTasksStore.fetchTodayTasks(authStore.user.id, true)
    } catch (error) {
      console.error('Error refreshing tasks before navigation:', error)
    }
  }
  router.push('/tasks')
}

const formatDate = (dateString) => {
  if (!dateString) return 'Tidak tersedia'

  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('id-ID', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch (error) {
    return 'Format tanggal tidak valid'
  }
}

const getStatusType = (status) => {
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

const getStatusText = (status) => {
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
  const stage = task?.current_stage || 1
  if (stage >= 5 && stage <= 8) {
    return '📦 Pengambilan Ompreng'
  }
  return '🍱 Pengiriman Makanan'
}

const getTaskTypeTagType = (task) => {
  const stage = task?.current_stage || 1
  if (stage >= 5 && stage <= 8) {
    return 'warning'
  }
  return 'success'
}

const showFullAddress = () => {
  if (task.value?.school?.address) {
    showAddressDialog.value = true
  } else {
    showToast('Alamat tidak tersedia')
  }
}

const copyAddress = async () => {
  if (task.value?.school?.address) {
    try {
      await navigator.clipboard.writeText(task.value.school.address)
      showSuccessToast('Alamat berhasil disalin')
    } catch (error) {
      showToast('Gagal menyalin alamat')
    }
  }
}

const callSchool = () => {
  const phoneNumber = task.value?.school?.phone_number
  if (phoneNumber) {
    window.location.href = `tel:${phoneNumber}`
  } else {
    showToast('Nomor telepon tidak tersedia')
  }
}

const openGPSNavigation = () => {
  if (!hasValidGPS.value) {
    showToast('Koordinat GPS tidak tersedia atau tidak valid')
    return
  }

  const { latitude, longitude } = task.value.school
  const schoolName = task.value.school.name || 'Sekolah Tujuan'

  const mapsUrl = `https://www.google.com/maps/dir/?api=1&destination=${latitude},${longitude}&destination_place_id=${encodeURIComponent(schoolName)}`
  const androidMapsUrl = `google.navigation:q=${latitude},${longitude}`
  const iosMapsUrl = `maps://maps.google.com/maps?daddr=${latitude},${longitude}&amp;ll=`

  const userAgent = navigator.userAgent.toLowerCase()

  if (userAgent.includes('android')) {
    window.location.href = androidMapsUrl
    setTimeout(() => {
      window.open(mapsUrl, '_blank')
    }, 1000)
  } else if (userAgent.includes('iphone') || userAgent.includes('ipad')) {
    window.location.href = iosMapsUrl
    setTimeout(() => {
      window.open(mapsUrl, '_blank')
    }, 1000)
  } else {
    window.open(mapsUrl, '_blank')
  }

  showToast('Membuka navigasi GPS...')
}

const startDelivery = async () => {
  try {
    const confirmed = await showConfirmDialog({
      title: 'Mulai Pengiriman',
      message: `Apakah Anda yakin ingin memulai pengiriman ke ${task.value.school?.name}?`,
      confirmButtonText: 'Ya, Mulai',
      cancelButtonText: 'Batal'
    })

    if (confirmed) {
      isUpdatingStatus.value = true
      await deliveryTasksStore.updateTaskStatus(task.value.id, 'in_progress')
      task.value.status = 'in_progress'
      showSuccessToast('Status pengiriman diperbarui')
    }
  } catch (error) {
    console.error('Error starting delivery:', error)
    showToast('Gagal memperbarui status')
  } finally {
    isUpdatingStatus.value = false
  }
}

const arrivedAtSchool = async () => {
  try {
    const confirmed = await showConfirmDialog({
      title: 'Sudah Sampai Sekolah',
      message: `Apakah Anda sudah sampai di ${task.value.school?.name}?`,
      confirmButtonText: 'Ya, Sudah Sampai',
      cancelButtonText: 'Belum'
    })

    if (confirmed) {
      isUpdatingStatus.value = true
      await deliveryTasksStore.updateTaskStatus(task.value.id, 'arrived')
      task.value.status = 'arrived'
      showSuccessToast('Status diperbarui: Sudah Sampai')
    }
  } catch (error) {
    console.error('Error updating arrival status:', error)
    showToast('Gagal memperbarui status')
  } finally {
    isUpdatingStatus.value = false
  }
}

const completeDelivery = async () => {
  try {
    const confirmed = await showConfirmDialog({
      title: 'Selesaikan Pengiriman',
      message: `Apakah Anda yakin pengiriman ke ${task.value.school?.name} sudah selesai tanpa e-POD?`,
      confirmButtonText: 'Ya, Selesai',
      cancelButtonText: 'Belum'
    })

    if (confirmed) {
      isUpdatingStatus.value = true
      await deliveryTasksStore.updateTaskStatus(task.value.id, 'completed')
      task.value.status = 'completed'
      showSuccessToast('Pengiriman berhasil diselesaikan')
    }
  } catch (error) {
    console.error('Error completing delivery:', error)
    showToast('Gagal menyelesaikan pengiriman')
  } finally {
    isUpdatingStatus.value = false
  }
}

const openePODForm = () => {
  router.push(`/tasks/${task.value.id}/epod`)
}

// Network status handlers
const handleOnline = () => {
  isOnline.value = true
  showToast('Koneksi internet tersambung')
  deliveryTasksStore.syncAllOfflineData()
}

const handleOffline = () => {
  isOnline.value = false
  showToast('Mode offline - Data akan disinkronkan saat online')
}

// Lifecycle
onMounted(() => {
  loadTask()
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
})

onUnmounted(() => {
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
})
</script>

<style scoped>
.detail-view {
  min-height: 100vh;
  background: var(--h-bg-primary);
  padding-top: 46px;
}

.detail-view__content {
  padding: var(--h-spacing-lg);
  padding-bottom: 80px;
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-lg);
}

.detail-view__nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-view__refresh--spinning {
  animation: detail-spin 1s linear infinite;
}

@keyframes detail-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Card base */
.detail-view__card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
}

.detail-view__card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin-bottom: var(--h-spacing-md);
}

/* Status card */
.detail-view__status-row {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-sm);
  flex-wrap: wrap;
  margin-bottom: var(--h-spacing-sm);
}

.detail-view__status-tag {
  font-weight: 600;
}

.detail-view__meta-row {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-xs);
  font-size: 13px;
  color: var(--h-text-secondary);
}

/* Info rows */
.detail-view__info-row {
  display: flex;
  align-items: flex-start;
  gap: var(--h-spacing-md);
  padding: var(--h-spacing-sm) 0;
}

.detail-view__info-row:not(:last-child) {
  border-bottom: 1px solid var(--h-border-light);
}

.detail-view__info-row .van-icon {
  margin-top: 2px;
  flex-shrink: 0;
}

.detail-view__info-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
  min-width: 0;
}

.detail-view__info-label {
  font-size: 12px;
  color: var(--h-text-secondary);
}

.detail-view__info-value {
  font-size: 14px;
  font-weight: 500;
  color: var(--h-text-primary);
  word-break: break-word;
}

.detail-view__info-value--link {
  color: var(--h-primary);
  cursor: pointer;
}

.detail-view__divider {
  height: 1px;
  background: var(--h-border-color);
  margin: var(--h-spacing-sm) 0;
}

/* Google Maps button */
.detail-view__maps-btn {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-sm);
  width: 100%;
  padding: var(--h-spacing-md) var(--h-spacing-lg);
  margin-top: var(--h-spacing-md);
  background: var(--h-primary-lighter);
  border: 1px solid var(--h-primary);
  border-radius: var(--h-radius-md);
  color: var(--h-primary);
  font-family: var(--h-font-family);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--h-transition-base);
  min-height: 44px;
}

.detail-view__maps-btn:active {
  background: var(--h-primary);
  color: #ffffff;
}

.detail-view__maps-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.detail-view__maps-btn span {
  flex: 1;
  text-align: left;
}

/* Action buttons */
.detail-view__action-btn {
  border-radius: var(--h-radius-md) !important;
  font-weight: 600;
  height: 48px;
}

.detail-view__action-btn:not(:last-child) {
  margin-bottom: var(--h-spacing-sm);
}

.detail-view__action-btn--epod {
  background: linear-gradient(135deg, #ff976a, #ff6b35) !important;
  border: none !important;
  color: white !important;
}

/* Responsive */
@media (max-width: 375px) {
  .detail-view__content {
    padding: var(--h-spacing-md);
    padding-bottom: 80px;
  }

  .detail-view__card {
    padding: var(--h-spacing-md);
  }
}
</style>
