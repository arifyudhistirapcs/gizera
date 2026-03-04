<template>
  <div class="detail-view">
    <!-- Navigation Bar -->
    <van-nav-bar
      title="Detail Tugas Pengambilan"
      left-arrow
      fixed
      @click-left="goBack"
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
      v-else-if="error"
      image="error"
      :description="error"
    >
      <van-button type="primary" @click="loadTaskDetail">Coba Lagi</van-button>
    </van-empty>

    <!-- Task Detail -->
    <div v-else-if="task" class="detail-view__content">
      <!-- Progress Steps -->
      <div class="detail-view__card">
        <van-steps :active="activeStep" active-color="#5A4372">
          <van-step>Menuju</van-step>
          <van-step>Tiba</van-step>
          <van-step>Kembali</van-step>
          <van-step>Selesai</van-step>
        </van-steps>
      </div>

      <!-- Status Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Status</div>
        <div class="detail-view__status-row">
          <van-tag
            :type="getStageType(deliveryRecord?.current_stage)"
            size="large"
            class="detail-view__status-tag"
          >
            {{ getStageName(deliveryRecord?.current_stage) }}
          </van-tag>
          <van-tag v-if="hasReview" type="success" size="medium">
            Ulasan Sudah Diisi
          </van-tag>
        </div>
      </div>

      <!-- School Info Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Info Sekolah</div>
        <div class="detail-view__info-row">
          <van-icon name="shop-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Sekolah</span>
            <span class="detail-view__info-value">{{ deliveryRecord?.school?.name || '-' }}</span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="location-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Alamat</span>
            <span class="detail-view__info-value">{{ deliveryRecord?.school?.address || '-' }}</span>
          </div>
        </div>
        <div class="detail-view__info-row">
          <van-icon name="aim" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Koordinat GPS</span>
            <span class="detail-view__info-value detail-view__info-value--mono">
              {{ deliveryRecord?.school?.latitude?.toFixed(6) }}, {{ deliveryRecord?.school?.longitude?.toFixed(6) }}
            </span>
          </div>
        </div>

        <!-- Google Maps Navigation Button -->
        <button
          class="detail-view__maps-btn"
          :disabled="!hasValidGPS"
          @click="openMaps"
        >
          <van-icon name="guide-o" size="20" />
          <span>Buka di Google Maps</span>
          <van-icon name="arrow" size="16" />
        </button>
      </div>

      <!-- Menu / Portion Info Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Info Ompreng</div>
        <div class="detail-view__info-row">
          <van-icon name="shopping-cart-o" color="var(--h-primary)" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Jumlah Ompreng</span>
            <span class="detail-view__info-value">{{ deliveryRecord?.ompreng_count || 0 }} wadah</span>
          </div>
        </div>
      </div>

      <!-- Ompreng Received Input (shown at stage 11) -->
      <div v-if="showOmprengInput" class="detail-view__card">
        <div class="detail-view__card-title">Input Ompreng Diterima</div>
        <van-field
          v-model="omprengReceived"
          type="digit"
          label="Jumlah Diterima"
          placeholder="Masukkan jumlah"
          :rules="[{ required: true, message: 'Wajib diisi' }]"
        >
          <template #button>
            <span class="detail-view__unit-text">wadah</span>
          </template>
        </van-field>

        <!-- Show difference if any -->
        <div v-if="omprengDifference !== 0" class="detail-view__info-row" style="padding-top: 8px;">
          <van-icon name="info-o" :color="omprengDifference < 0 ? 'var(--h-error)' : 'var(--h-warning)'" />
          <div class="detail-view__info-text">
            <span class="detail-view__info-label">Selisih</span>
            <van-tag :type="omprengDifference < 0 ? 'danger' : 'warning'">
              {{ omprengDifference > 0 ? '+' : '' }}{{ omprengDifference }} wadah
            </van-tag>
          </div>
        </div>

        <!-- Reason input if there's difference -->
        <van-field
          v-if="omprengDifference !== 0"
          v-model="omprengDifferenceReason"
          type="textarea"
          label="Alasan Selisih"
          placeholder="Jelaskan alasan selisih jumlah ompreng"
          rows="2"
          :rules="[{ required: true, message: 'Alasan wajib diisi jika ada selisih' }]"
          style="margin-top: 8px;"
        />
      </div>

      <!-- Action Card -->
      <div class="detail-view__card">
        <div class="detail-view__card-title">Aksi</div>

        <!-- Review Button (shown at stage 11 if no review yet) -->
        <van-button
          v-if="showReviewButton"
          type="primary"
          block
          size="large"
          class="detail-view__action-btn detail-view__action-btn--review"
          @click="goToReviewForm"
        >
          📝 Isi Form Ulasan
        </van-button>

        <!-- Stage Update Button -->
        <van-button
          v-if="canUpdateStage"
          type="primary"
          block
          size="large"
          :loading="isUpdating"
          :disabled="needsReviewFirst || !canProceedWithOmpreng"
          class="detail-view__action-btn"
          @click="updateToNextStage"
        >
          {{ getNextStageButtonText() }}
        </van-button>

        <!-- Info message if review needed -->
        <van-notice-bar
          v-if="needsReviewFirst"
          color="#5A4372"
          background="var(--h-primary-lighter)"
          left-icon="info-o"
          class="detail-view__notice"
        >
          Mohon isi form ulasan terlebih dahulu sebelum melanjutkan
        </van-notice-bar>

        <!-- Info message if ompreng input needed -->
        <van-notice-bar
          v-if="showOmprengInput && !canProceedWithOmpreng && !needsReviewFirst"
          color="#5A4372"
          background="var(--h-primary-lighter)"
          left-icon="info-o"
          class="detail-view__notice"
        >
          {{ omprengDifference !== 0 && !omprengDifferenceReason ? 'Mohon isi alasan selisih ompreng' : 'Mohon isi jumlah ompreng diterima' }}
        </van-notice-bar>

        <!-- No actions available -->
        <div
          v-if="!canUpdateStage && !showReviewButton && deliveryRecord?.current_stage >= 13"
          class="detail-view__completed-msg"
        >
          <van-icon name="success" color="var(--h-success)" size="24" />
          <span>Pengambilan selesai</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import api from '@/services/api'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

const router = useRouter()
const route = useRoute()

const isLoading = ref(false)
const isUpdating = ref(false)
const error = ref(null)
const task = ref(null)
const deliveryRecord = ref(null)
const hasReview = ref(false)

// Ompreng input
const omprengReceived = ref('')
const omprengDifferenceReason = ref('')

const goBack = async () => {
  // Force refresh tasks when going back to ensure latest data
  const { useDeliveryTasksStore } = await import('@/stores/deliveryTasks')
  const { useAuthStore } = await import('@/stores/auth')
  const deliveryTasksStore = useDeliveryTasksStore()
  const authStore = useAuthStore()
  
  if (authStore.user?.id) {
    // Clear existing tasks and fetch fresh data
    deliveryTasksStore.clearTasks()
    await deliveryTasksStore.fetchTodayTasks(authStore.user.id, true)
  }
  
  router.push('/tasks')
}

// Computed: progress step based on stage
const activeStep = computed(() => {
  if (!deliveryRecord.value) return 0
  const stage = deliveryRecord.value.current_stage
  if (stage >= 13) return 3
  if (stage >= 12) return 2
  if (stage >= 11) return 1
  return 0
})

// Computed: valid GPS
const hasValidGPS = computed(() => {
  const school = deliveryRecord.value?.school
  return school?.latitude && school?.longitude &&
         Math.abs(school.latitude) <= 90 &&
         Math.abs(school.longitude) <= 180
})

// Calculate ompreng difference
const omprengDifference = computed(() => {
  if (!omprengReceived.value || !deliveryRecord.value?.ompreng_count) return 0
  return parseInt(omprengReceived.value) - (deliveryRecord.value.ompreng_count || 0)
})

// Show ompreng input at stage 11
const showOmprengInput = computed(() => {
  if (!deliveryRecord.value) return false
  return deliveryRecord.value.current_stage === 11
})

// Can proceed with ompreng validation
const canProceedWithOmpreng = computed(() => {
  if (!showOmprengInput.value) return true
  if (!omprengReceived.value) return false
  if (omprengDifference.value !== 0 && !omprengDifferenceReason.value.trim()) return false
  return true
})

const checkReviewExists = async () => {
  if (!deliveryRecord.value?.id) return

  try {
    const response = await api.get('/reviews/check', {
      params: { delivery_record_id: deliveryRecord.value.id }
    })
    hasReview.value = response.data.exists || false
  } catch (err) {
    console.error('Error checking review:', err)
    hasReview.value = false
  }
}

const loadTaskDetail = async () => {
  isLoading.value = true
  error.value = null

  try {
    const pickupTaskId = route.params.id
    const recordId = route.query.record

    const response = await api.get(`/pickup-tasks/${pickupTaskId}`)

    if (response.data.pickup_task) {
      task.value = response.data.pickup_task

      // Find the specific delivery record
      if (recordId && task.value.delivery_records) {
        deliveryRecord.value = task.value.delivery_records.find(
          dr => dr.id === parseInt(recordId)
        )
      } else if (task.value.delivery_records?.length > 0) {
        deliveryRecord.value = task.value.delivery_records[0]
      }

      // Pre-fill ompreng received with expected count
      if (deliveryRecord.value?.ompreng_count) {
        omprengReceived.value = String(deliveryRecord.value.ompreng_count)
      }

      // Check if review exists
      await checkReviewExists()
    } else {
      error.value = 'Tugas tidak ditemukan'
    }
  } catch (err) {
    console.error('Error loading pickup task:', err)
    error.value = 'Gagal memuat detail tugas'
  } finally {
    isLoading.value = false
  }
}

const getStageType = (stage) => {
  if (stage >= 13) return 'success'
  if (stage >= 10) return 'primary'
  return 'warning'
}

const getStageName = (stage) => {
  const stages = {
    10: 'Menuju Lokasi',
    11: 'Tiba di Lokasi',
    12: 'Kembali ke SPPG',
    13: 'Tiba di SPPG'
  }
  return stages[stage] || `Stage ${stage}`
}

const canUpdateStage = computed(() => {
  if (!deliveryRecord.value) return false
  const stage = deliveryRecord.value.current_stage
  return stage >= 10 && stage < 13
})

// Show review button when at stage 11 and no review yet
const showReviewButton = computed(() => {
  if (!deliveryRecord.value) return false
  return deliveryRecord.value.current_stage === 11 && !hasReview.value
})

// Need review before proceeding from stage 11 to 12
const needsReviewFirst = computed(() => {
  if (!deliveryRecord.value) return false
  return deliveryRecord.value.current_stage === 11 && !hasReview.value
})

const getNextStageButtonText = () => {
  if (!deliveryRecord.value) return 'Update Status'
  const stage = deliveryRecord.value.current_stage
  const nextStages = {
    10: 'Sudah Tiba di Lokasi',
    11: 'Mulai Kembali ke SPPG',
    12: 'Sudah Tiba di SPPG'
  }
  return nextStages[stage] || 'Update Status'
}

const goToReviewForm = () => {
  const school = deliveryRecord.value?.school
  router.push({
    path: '/review',
    query: {
      school: school?.name || '-',
      school_id: school?.id || deliveryRecord.value?.school_id,
      delivery_record_id: deliveryRecord.value?.id,
      date: new Date().toLocaleDateString('id-ID'),
      return_path: route.fullPath
    }
  })
}

const updateToNextStage = async () => {
  if (!deliveryRecord.value) return

  const currentStage = deliveryRecord.value.current_stage

  // Block stage 11 -> 12 if no review
  if (currentStage === 11 && !hasReview.value) {
    showToast('Mohon isi form ulasan terlebih dahulu')
    return
  }

  // Block if ompreng validation not complete
  if (currentStage === 11 && !canProceedWithOmpreng.value) {
    if (!omprengReceived.value) {
      showToast('Mohon isi jumlah ompreng diterima')
    } else {
      showToast('Mohon isi alasan selisih ompreng')
    }
    return
  }

  const nextStage = currentStage + 1

  const stageStatusMap = {
    11: 'driver_tiba_di_lokasi_pengambilan',
    12: 'driver_kembali_ke_sppg',
    13: 'driver_tiba_di_sppg'
  }

  // Build confirmation message
  let confirmMessage = `Apakah Anda yakin ingin mengupdate status ke "${getStageName(nextStage)}"?`
  if (currentStage === 11 && omprengDifference.value !== 0) {
    confirmMessage += `\n\nJumlah ompreng diterima: ${omprengReceived.value} wadah (selisih ${omprengDifference.value})`
  }

  try {
    await showConfirmDialog({
      title: 'Konfirmasi',
      message: confirmMessage
    })

    isUpdating.value = true

    // Build request payload
    const payload = {
      stage: nextStage,
      status: stageStatusMap[nextStage]
    }

    // Add ompreng data if at stage 11
    if (currentStage === 11) {
      payload.ompreng_received = parseInt(omprengReceived.value)
      if (omprengDifference.value !== 0) {
        payload.ompreng_difference_reason = omprengDifferenceReason.value
      }
    }

    const response = await api.put(
      `/pickup-tasks/${task.value.id}/delivery-records/${deliveryRecord.value.id}/stage`,
      payload
    )

    if (response.data.delivery_record) {
      deliveryRecord.value = response.data.delivery_record
      showToast('Status berhasil diupdate')
    }
  } catch (err) {
    if (err !== 'cancel') {
      console.error('Error updating stage:', err)
      showToast('Gagal mengupdate status')
    }
  } finally {
    isUpdating.value = false
  }
}

const openMaps = () => {
  const school = deliveryRecord.value?.school
  if (!school?.latitude || !school?.longitude) {
    showToast('Koordinat GPS tidak tersedia')
    return
  }

  const url = `https://www.google.com/maps/dir/?api=1&destination=${school.latitude},${school.longitude}`
  window.open(url, '_blank')
}

// Watch for review_submitted query param (coming back from review form)
watch(() => route.query.review_submitted, async (submitted) => {
  if (submitted === 'true') {
    hasReview.value = true
    router.replace({ query: { ...route.query, review_submitted: undefined } })
  }
})

onMounted(() => {
  loadTaskDetail()
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

/* Card base */
.detail-view__card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
}

.detail-view__card-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--h-text-primary);
  margin-bottom: var(--h-spacing-md);
}

/* Status card */
.detail-view__status-row {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-sm);
  flex-wrap: wrap;
}

.detail-view__status-tag {
  font-weight: 600;
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

.detail-view__info-value--mono {
  font-family: monospace;
  font-size: 12px;
}

.detail-view__unit-text {
  color: var(--h-text-light);
  font-size: 14px;
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

.detail-view__action-btn--review {
  background: var(--h-primary) !important;
  border-color: var(--h-primary) !important;
}

.detail-view__notice {
  margin-top: var(--h-spacing-md);
  border-radius: var(--h-radius-sm);
}

.detail-view__completed-msg {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--h-spacing-sm);
  padding: var(--h-spacing-lg);
  font-size: 15px;
  font-weight: 600;
  color: var(--h-success);
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
