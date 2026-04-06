<template>
  <div class="delivery-detail-view">
    <a-page-header
      title="Detail Aktivitas"
      @back="goBack"
    >
      <template #extra>
        <a-button @click="refreshData" :loading="loading">
          <template #icon><reload-outlined /></template>
          Refresh
        </a-button>
      </template>
    </a-page-header>

    <div class="content-wrapper">
      <a-spin :spinning="loading">
        <!-- School Information Section -->
        <a-card title="Informasi Sekolah" style="margin-bottom: 16px">
          <a-descriptions bordered :column="{ xs: 1, sm: 2, md: 2 }">
            <a-descriptions-item label="Nama Sekolah">
              {{ deliveryDetail?.school?.name || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Jumlah Porsi">
              {{ deliveryDetail?.portions || 0 }} porsi
            </a-descriptions-item>
            <a-descriptions-item label="Alamat" :span="2">
              {{ deliveryDetail?.school?.address || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Kontak">
              {{ deliveryDetail?.school?.contact_person || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Nomor Telepon">
              {{ deliveryDetail?.school?.phone || '-' }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- Driver Information Section -->
        <a-card title="Informasi Driver" style="margin-bottom: 16px">
          <a-descriptions bordered :column="{ xs: 1, sm: 2, md: 2 }">
            <a-descriptions-item label="Nama Driver">
              {{ deliveryDetail?.driver?.full_name || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Jenis Kendaraan">
              {{ deliveryDetail?.driver?.vehicle_type || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Nomor Telepon">
              {{ deliveryDetail?.driver?.phone || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Email">
              {{ deliveryDetail?.driver?.email || '-' }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- Timeline Section -->
        <a-card title="Timeline Aktivitas" style="margin-bottom: 16px">
          <DeliveryTimeline
            v-if="deliveryDetail"
            :current-status="deliveryDetail.current_status"
            :activity-log="activityLog"
            @view-epod="showEPODModal"
          />
        </a-card>

        <!-- Activity Log Section -->
        <a-card title="Log Aktivitas">
          <ActivityLogTable
            v-if="activityLog.length > 0"
            :activity-log="activityLog"
          />
          <HEmptyState v-else description="Belum ada aktivitas" />
        </a-card>
      </a-spin>
    </div>

    <!-- e-POD Modal -->
    <a-modal
      v-model:open="epodModalVisible"
      title="Detail e-POD (Bukti Pengiriman)"
      :footer="null"
      width="700px"
    >
      <a-spin :spinning="epodLoading">
        <div v-if="epodData">
          <a-descriptions bordered :column="1" style="margin-bottom: 16px">
            <a-descriptions-item label="Nama Penerima">
              {{ epodData.recipient_name || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Waktu Selesai">
              {{ formatDateTime(epodData.completed_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="Lokasi GPS">
              <a :href="getGoogleMapsUrl(epodData.latitude, epodData.longitude)" target="_blank">
                {{ epodData.latitude?.toFixed(6) }}, {{ epodData.longitude?.toFixed(6) }}
                <environment-outlined style="margin-left: 4px" />
              </a>
            </a-descriptions-item>
          </a-descriptions>

          <!-- Photo Section -->
          <div style="margin-bottom: 16px">
            <h4 style="margin-bottom: 8px">📷 Foto Bukti Pengiriman</h4>
            <div v-if="epodData.photo_url" class="epod-image-container">
              <a-image
                :src="getImageUrl(epodData.photo_url)"
                :preview="{ src: getImageUrl(epodData.photo_url) }"
                style="max-width: 100%; max-height: 300px; object-fit: contain"
              />
            </div>
            <HEmptyState v-else description="Tidak ada foto" />
          </div>

          <!-- Signature Section -->
          <div>
            <h4 style="margin-bottom: 8px">✍️ Tanda Tangan Penerima</h4>
            <div v-if="epodData.signature_url" class="epod-image-container">
              <a-image
                :src="getImageUrl(epodData.signature_url)"
                :preview="{ src: getImageUrl(epodData.signature_url) }"
                style="max-width: 100%; max-height: 200px; object-fit: contain; background: #f5f5f5; border-radius: 8px"
              />
            </div>
            <HEmptyState v-else description="Tidak ada tanda tangan" />
          </div>
        </div>
        <HEmptyState v-else description="Data e-POD tidak ditemukan" />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { ReloadOutlined, EnvironmentOutlined } from '@ant-design/icons-vue'
import { getDeliveryDetail, getActivityLog, getEPODByDeliveryTask } from '@/services/monitoringService'
import DeliveryTimeline from '@/components/DeliveryTimeline.vue'
import ActivityLogTable from '@/components/ActivityLogTable.vue'
import HEmptyState from '@/components/common/HEmptyState.vue'

const router = useRouter()
const route = useRoute()

// API base URL for images
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// State
const loading = ref(false)
const deliveryDetail = ref(null)
const activityLog = ref([])
const epodModalVisible = ref(false)
const epodLoading = ref(false)
const epodData = ref(null)

// Methods
const goBack = () => {
  router.push('/monitoring-activity')
}

const refreshData = () => {
  fetchData()
}

const fetchData = async () => {
  await Promise.all([
    fetchDeliveryDetail(),
    fetchActivityLog()
  ])
}

const fetchDeliveryDetail = async () => {
  loading.value = true
  try {
    const id = route.params.id
    const response = await getDeliveryDetail(id)
    
    if (response.success) {
      deliveryDetail.value = response.data
    } else {
      message.error(response.message || 'Gagal memuat detail aktivitas')
    }
  } catch (error) {
    console.error('Error fetching delivery detail:', error)
    if (error.response?.status === 404) {
      message.error('Data pengiriman tidak ditemukan')
      goBack()
    } else {
      message.error(error.response?.data?.message || 'Gagal memuat detail aktivitas')
    }
  } finally {
    loading.value = false
  }
}

const fetchActivityLog = async () => {
  try {
    const id = route.params.id
    const response = await getActivityLog(id)
    
    if (response.success) {
      activityLog.value = response.data || []
    } else {
      message.error(response.message || 'Gagal memuat log aktivitas')
    }
  } catch (error) {
    console.error('Error fetching activity log:', error)
    message.error(error.response?.data?.message || 'Gagal memuat log aktivitas')
  }
}

const showEPODModal = async () => {
  epodModalVisible.value = true
  epodLoading.value = true
  
  try {
    const id = route.params.id
    // Use 'record' type since we're in monitoring view which uses delivery_record_id
    const response = await getEPODByDeliveryTask(id, 'record')
    
    if (response.success) {
      epodData.value = response.epod
    } else {
      message.error(response.message || 'Gagal memuat data e-POD')
      epodData.value = null
    }
  } catch (error) {
    console.error('Error fetching e-POD:', error)
    if (error.response?.status === 404) {
      message.warning('Data e-POD belum tersedia')
    } else {
      message.error(error.response?.data?.message || 'Gagal memuat data e-POD')
    }
    epodData.value = null
  } finally {
    epodLoading.value = false
  }
}

const formatDateTime = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getGoogleMapsUrl = (lat, lng) => {
  return `https://www.google.com/maps?q=${lat},${lng}`
}

const getImageUrl = (url) => {
  if (!url) return ''
  // If URL is already absolute, return as is
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  // Otherwise, prepend API base URL
  return `${API_BASE_URL}${url}`
}

// Lifecycle
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.delivery-detail-view {
  /* HorizonLayout handles padding/bg */
}

.content-wrapper {
  margin-top: 16px;
}

:deep(.ant-descriptions-item-label) {
  font-weight: 600;
  background-color: #fafafa;
}

.epod-image-container {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 12px;
  background-color: #fafafa;
  text-align: center;
}

.epod-image-container :deep(.ant-image) {
  display: block;
  margin: 0 auto;
}
</style>
