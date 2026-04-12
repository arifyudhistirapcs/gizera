<template>
  <div class="review-form-container">
    <!-- Success Animation Overlay -->
    <div v-if="showSuccessAnimation" class="success-overlay">
      <LottiePlayer src="/lottie/success-check.json" width="120px" height="120px" :loop="false" />
    </div>

    <!-- Navigation Bar -->
    <van-nav-bar 
      title="Form Ulasan Pengiriman" 
      left-arrow 
      @click-left="goBack"
      fixed
    />

    <!-- Loading State -->
    <van-loading v-if="isLoading" type="spinner" vertical class="loading-state">
      Memuat...
    </van-loading>

    <!-- Form Content -->
    <div v-else class="form-content">
      <!-- School Info -->
      <van-cell-group inset class="info-card">
        <van-cell title="Sekolah" :value="schoolName" />
        <van-cell title="Tanggal" :value="deliveryDate" />
      </van-cell-group>

      <!-- Reviewer Info -->
      <van-cell-group inset class="info-card">
        <van-field
          v-model="formData.reviewer_name"
          label="Nama Pengisi"
          placeholder="Nama lengkap"
          :error-message="errors.reviewer_name"
        />
        <van-field
          v-model="formData.reviewer_role"
          label="Jabatan"
          placeholder="Contoh: Guru, Kepala Sekolah"
          :error-message="errors.reviewer_role"
        />
      </van-cell-group>

      <!-- Menu Ratings -->
      <van-cell-group inset class="info-card">
        <van-cell title="Penilaian Menu" class="section-title" />
        
        <div class="rating-item">
          <span class="rating-label">Rasa Makanan</span>
          <van-rate v-model="formData.rating_food_taste" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Kebersihan & Kerapian Penyajian</span>
          <van-rate v-model="formData.rating_food_cleanliness" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Kesesuaian Menu dengan Jadwal</span>
          <van-rate v-model="formData.rating_menu_accuracy" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Porsi Makanan</span>
          <van-rate v-model="formData.rating_portion_size" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Variasi Menu</span>
          <van-rate v-model="formData.rating_menu_variety" :count="5" allow-half />
        </div>
      </van-cell-group>

      <!-- Service Ratings -->
      <van-cell-group inset class="info-card">
        <van-cell title="Penilaian Layanan" class="section-title" />
        
        <div class="rating-item">
          <span class="rating-label">Ketepatan Waktu Pengantaran</span>
          <van-rate v-model="formData.rating_delivery_time" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Sikap & Keramahan Driver/Kader</span>
          <van-rate v-model="formData.rating_driver_attitude" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Kondisi Makanan Saat Diterima</span>
          <van-rate v-model="formData.rating_food_condition" :count="5" allow-half />
        </div>

        <div class="rating-item">
          <span class="rating-label">Kerapihan & Kebersihan Pengantar</span>
          <van-rate v-model="formData.rating_driver_tidiness" :count="5" allow-half />
        </div>
        
        <div class="rating-item">
          <span class="rating-label">Konsistensi Layanan</span>
          <van-rate v-model="formData.rating_service_consistency" :count="5" allow-half />
        </div>
      </van-cell-group>

      <!-- Comments & Photo -->
      <van-cell-group inset class="info-card">
        <van-field
          v-model="formData.comments"
          label="Catatan"
          type="textarea"
          rows="3"
          placeholder="Catatan tambahan (opsional)"
        />
        
        <van-cell title="Foto (Opsional)">
          <template #value>
            <van-uploader
              v-model="photoFiles"
              :max-count="1"
              :after-read="onPhotoRead"
              accept="image/*"
            />
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Submit Button -->
      <div class="action-buttons">
        <van-button 
          type="primary" 
          block 
          size="large"
          @click="submitReview"
          :loading="isSubmitting"
          :disabled="!isFormValid"
        >
          Kirim Ulasan
        </van-button>
      </div>
    </div>

    <!-- Bottom Navigation -->
    <van-tabbar v-model="active" route fixed>
      <van-tabbar-item to="/tasks" icon="orders-o">Tugas</van-tabbar-item>
      <van-tabbar-item to="/attendance" icon="clock-o">Absensi</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">Profil</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showSuccessToast } from 'vant'
import api from '@/services/api'
import LottiePlayer from '@/components/common/LottiePlayer.vue'

const router = useRouter()
const route = useRoute()

const active = ref(0)
const isLoading = ref(false)
const isSubmitting = ref(false)
const schoolName = ref('')
const showSuccessAnimation = ref(false)
const schoolId = ref(null)
const deliveryRecordId = ref(null)
const deliveryDate = ref('')
const photoFiles = ref([])
const photoBase64 = ref('')

const formData = ref({
  reviewer_name: '',
  reviewer_role: '',
  rating_food_taste: 0,
  rating_food_cleanliness: 0,
  rating_menu_accuracy: 0,
  rating_portion_size: 0,
  rating_menu_variety: 0,
  rating_delivery_time: 0,
  rating_driver_attitude: 0,
  rating_food_condition: 0,
  rating_driver_tidiness: 0,
  rating_service_consistency: 0,
  comments: ''
})

const errors = ref({
  reviewer_name: '',
  reviewer_role: ''
})

const goBack = () => {
  router.back()
}

const isFormValid = computed(() => {
  return (
    formData.value.reviewer_name.trim() !== '' &&
    formData.value.rating_food_taste > 0 &&
    formData.value.rating_food_cleanliness > 0 &&
    formData.value.rating_menu_accuracy > 0 &&
    formData.value.rating_portion_size > 0 &&
    formData.value.rating_menu_variety > 0 &&
    formData.value.rating_delivery_time > 0 &&
    formData.value.rating_driver_attitude > 0 &&
    formData.value.rating_food_condition > 0 &&
    formData.value.rating_driver_tidiness > 0 &&
    formData.value.rating_service_consistency > 0
  )
})

const onPhotoRead = (file) => {
  photoBase64.value = file.content
}

const loadData = () => {
  // Get data from route query params
  schoolName.value = route.query.school || '-'
  schoolId.value = parseInt(route.query.school_id) || null
  deliveryRecordId.value = parseInt(route.query.delivery_record_id) || null
  deliveryDate.value = route.query.date || new Date().toLocaleDateString('id-ID')
}

const submitReview = async () => {
  if (!isFormValid.value) {
    showToast('Mohon lengkapi semua rating')
    return
  }

  if (!deliveryRecordId.value || !schoolId.value) {
    showToast('Data pengiriman tidak valid')
    return
  }

  isSubmitting.value = true

  try {
    const payload = {
      delivery_record_id: deliveryRecordId.value,
      school_id: schoolId.value,
      reviewer_name: formData.value.reviewer_name,
      reviewer_role: formData.value.reviewer_role,
      rating_food_taste: Math.round(formData.value.rating_food_taste),
      rating_food_cleanliness: Math.round(formData.value.rating_food_cleanliness),
      rating_menu_accuracy: Math.round(formData.value.rating_menu_accuracy),
      rating_portion_size: Math.round(formData.value.rating_portion_size),
      rating_menu_variety: Math.round(formData.value.rating_menu_variety),
      rating_delivery_time: Math.round(formData.value.rating_delivery_time),
      rating_driver_attitude: Math.round(formData.value.rating_driver_attitude),
      rating_food_condition: Math.round(formData.value.rating_food_condition),
      rating_driver_tidiness: Math.round(formData.value.rating_driver_tidiness),
      rating_service_consistency: Math.round(formData.value.rating_service_consistency),
      comments: formData.value.comments,
      photo_url: photoBase64.value || ''
    }

    await api.post('/reviews', payload)
    
    showSuccessAnimation.value = true
    showSuccessToast('Ulasan berhasil dikirim')
    
    // Navigate back with success flag after animation
    setTimeout(() => {
      showSuccessAnimation.value = false
      router.replace({
        path: route.query.return_path || '/tasks',
        query: { review_submitted: 'true' }
      })
    }, 1500)
  } catch (err) {
    console.error('Error submitting review:', err)
    if (err.response?.data?.error_code === 'REVIEW_EXISTS') {
      showToast('Ulasan sudah pernah dikirim')
    } else {
      showToast('Gagal mengirim ulasan')
    }
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
/* Success Animation Overlay */
.success-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.85);
  z-index: 9999;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.review-form-container {
  min-height: 100vh;
  background-color: #F7F8FA;
  padding-top: 46px;
  padding-bottom: 70px;
}

.loading-state {
  padding-top: 100px;
}

.form-content {
  padding: 16px;
}

.info-card {
  margin-bottom: 16px;
  border-radius: 16px !important;
  overflow: hidden;
}

.section-title {
  font-weight: bold;
  background-color: #303030;
  color: white;
}

.section-title :deep(.van-cell__title) {
  color: white;
}

.rating-item {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.rating-item:last-child {
  border-bottom: none;
}

.rating-label {
  font-size: 14px;
  color: #333;
}

.action-buttons {
  margin-top: 24px;
  padding-bottom: 20px;
}

:deep(.van-rate__icon) {
  font-size: 24px;
}

:deep(.van-rate__icon--full) {
  color: #303030;
}
</style>
