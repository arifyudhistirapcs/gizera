<template>
  <div class="review-list">
    <a-page-header
      title="Ulasan & Rating"
      sub-title="Ulasan pengiriman makanan dari sekolah"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Summary Statistics -->
        <a-row :gutter="16">
          <a-col :span="6">
            <a-card size="small" class="stat-card">
              <a-statistic
                title="Total Ulasan"
                :value="summary.total_reviews || 0"
                :value-style="{ color: '#5A4372' }"
              >
                <template #prefix><MessageOutlined /></template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small" class="stat-card">
              <a-statistic
                title="Rating Keseluruhan"
                :value="summary.average_overall_rating || 0"
                :precision="1"
                :value-style="{ color: '#faad14' }"
              >
                <template #prefix><StarFilled /></template>
                <template #suffix>/ 5</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small" class="stat-card">
              <a-statistic
                title="Rating Menu"
                :value="summary.average_menu_rating || 0"
                :precision="1"
                :value-style="{ color: '#52c41a' }"
              >
                <template #prefix><CoffeeOutlined /></template>
                <template #suffix>/ 5</template>
              </a-statistic>
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small" class="stat-card">
              <a-statistic
                title="Rating Layanan"
                :value="summary.average_service_rating || 0"
                :precision="1"
                :value-style="{ color: '#1890ff' }"
              >
                <template #prefix><CarOutlined /></template>
                <template #suffix>/ 5</template>
              </a-statistic>
            </a-card>
          </a-col>
        </a-row>

        <!-- Filter Section -->
        <a-card size="small" title="Filter">
          <a-form layout="inline" :model="filters">
            <a-form-item label="Sekolah">
              <a-select
                v-model:value="filters.schoolId"
                placeholder="Semua Sekolah"
                style="width: 200px"
                allow-clear
                show-search
                option-filter-prop="children"
              >
                <a-select-option 
                  v-for="school in schools" 
                  :key="school.id" 
                  :value="school.id"
                >
                  {{ school.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
            
            <a-form-item label="Periode">
              <a-range-picker
                v-model:value="dateRange"
                format="DD/MM/YYYY"
                style="width: 280px"
              />
            </a-form-item>

            <a-form-item>
              <a-button type="primary" @click="handleSearch" :loading="loading">
                <template #icon><SearchOutlined /></template>
                Cari
              </a-button>
            </a-form-item>

            <a-form-item>
              <a-button @click="resetFilters">
                <template #icon><ClearOutlined /></template>
                Reset
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>

        <!-- Rating Breakdown -->
        <a-card size="small" title="Breakdown Rating" v-if="summary.total_reviews > 0">
          <a-row :gutter="[16, 16]">
            <a-col :span="12">
              <h4 style="color: #5A4372; margin-bottom: 12px;">Rating Menu</h4>
              <div class="rating-breakdown">
                <div class="rating-item">
                  <span>Rasa Makanan</span>
                  <a-rate :value="summary.avg_food_taste || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_food_taste || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Kebersihan Penyajian</span>
                  <a-rate :value="summary.avg_food_cleanliness || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_food_cleanliness || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Kesesuaian Menu</span>
                  <a-rate :value="summary.avg_menu_accuracy || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_menu_accuracy || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Porsi Makanan</span>
                  <a-rate :value="summary.avg_portion_size || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_portion_size || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Variasi Menu</span>
                  <a-rate :value="summary.avg_menu_variety || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_menu_variety || 0).toFixed(1) }}</span>
                </div>
              </div>
            </a-col>
            <a-col :span="12">
              <h4 style="color: #5A4372; margin-bottom: 12px;">Rating Layanan</h4>
              <div class="rating-breakdown">
                <div class="rating-item">
                  <span>Ketepatan Waktu</span>
                  <a-rate :value="summary.avg_delivery_time || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_delivery_time || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Sikap Driver</span>
                  <a-rate :value="summary.avg_driver_attitude || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_driver_attitude || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Kondisi Makanan</span>
                  <a-rate :value="summary.avg_food_condition || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_food_condition || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Kerapihan Driver</span>
                  <a-rate :value="summary.avg_driver_tidiness || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_driver_tidiness || 0).toFixed(1) }}</span>
                </div>
                <div class="rating-item">
                  <span>Konsistensi Layanan</span>
                  <a-rate :value="summary.avg_service_consistency || 0" disabled allow-half />
                  <span class="rating-value">{{ (summary.avg_service_consistency || 0).toFixed(1) }}</span>
                </div>
              </div>
            </a-col>
          </a-row>
        </a-card>

        <!-- Reviews Table -->
        <a-table
          :columns="columns"
          :data-source="reviews"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'school'">
              {{ record.school?.name || '-' }}
            </template>
            <template v-else-if="column.key === 'overall_rating'">
              <a-rate :value="record.overall_rating" disabled allow-half />
              <span style="margin-left: 8px;">{{ record.overall_rating?.toFixed(1) }}</span>
            </template>
            <template v-else-if="column.key === 'menu_rating'">
              <a-tag color="green">{{ record.average_menu_rating?.toFixed(1) }}</a-tag>
            </template>
            <template v-else-if="column.key === 'service_rating'">
              <a-tag color="blue">{{ record.average_service_rating?.toFixed(1) }}</a-tag>
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDate(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-button type="link" size="small" @click="showDetail(record)">
                <EyeOutlined /> Detail
              </a-button>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Detail Modal -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Detail Ulasan"
      width="700px"
      :footer="null"
    >
      <div v-if="selectedReview" class="review-detail">
        <a-descriptions bordered :column="2" size="small">
          <a-descriptions-item label="Sekolah" :span="2">
            {{ selectedReview.school?.name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Pengisi">
            {{ selectedReview.reviewer_name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Jabatan">
            {{ selectedReview.reviewer_role || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Tanggal" :span="2">
            {{ formatDate(selectedReview.created_at) }}
          </a-descriptions-item>
        </a-descriptions>

        <a-divider>Rating</a-divider>

        <a-row :gutter="16">
          <a-col :span="12">
            <h4>Menu</h4>
            <div class="detail-rating">
              <div><span>Rasa:</span> <a-rate :value="selectedReview.rating_food_taste" disabled /></div>
              <div><span>Kebersihan:</span> <a-rate :value="selectedReview.rating_food_cleanliness" disabled /></div>
              <div><span>Kesesuaian:</span> <a-rate :value="selectedReview.rating_menu_accuracy" disabled /></div>
              <div><span>Porsi:</span> <a-rate :value="selectedReview.rating_portion_size" disabled /></div>
              <div><span>Variasi:</span> <a-rate :value="selectedReview.rating_menu_variety" disabled /></div>
            </div>
          </a-col>
          <a-col :span="12">
            <h4>Layanan</h4>
            <div class="detail-rating">
              <div><span>Waktu:</span> <a-rate :value="selectedReview.rating_delivery_time" disabled /></div>
              <div><span>Sikap:</span> <a-rate :value="selectedReview.rating_driver_attitude" disabled /></div>
              <div><span>Kondisi:</span> <a-rate :value="selectedReview.rating_food_condition" disabled /></div>
              <div><span>Kerapihan:</span> <a-rate :value="selectedReview.rating_driver_tidiness" disabled /></div>
              <div><span>Konsistensi:</span> <a-rate :value="selectedReview.rating_service_consistency" disabled /></div>
            </div>
          </a-col>
        </a-row>

        <a-divider v-if="selectedReview.comments">Catatan</a-divider>
        <p v-if="selectedReview.comments">{{ selectedReview.comments }}</p>

        <div v-if="selectedReview.photo_url" style="margin-top: 16px;">
          <h4>Foto</h4>
          <a-image :src="getPhotoUrl(selectedReview.photo_url)" :width="200" />
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { 
  SearchOutlined, 
  ClearOutlined, 
  EyeOutlined,
  StarFilled,
  MessageOutlined,
  CoffeeOutlined,
  CarOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import reviewService from '@/services/reviewService'
import schoolService from '@/services/schoolService'

const loading = ref(false)
const reviews = ref([])
const schools = ref([])
const summary = ref({})
const dateRange = ref([])
const detailModalVisible = ref(false)
const selectedReview = ref(null)

const filters = reactive({
  schoolId: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total, range) => `${range[0]}-${range[1]} dari ${total} data`
})

const columns = [
  { title: 'Sekolah', key: 'school', width: 200 },
  { title: 'Pengisi', dataIndex: 'reviewer_name', key: 'reviewer_name', width: 150 },
  { title: 'Rating Keseluruhan', key: 'overall_rating', width: 200 },
  { title: 'Menu', key: 'menu_rating', width: 80, align: 'center' },
  { title: 'Layanan', key: 'service_rating', width: 80, align: 'center' },
  { title: 'Tanggal', key: 'created_at', width: 120 },
  { title: 'Aksi', key: 'actions', width: 100, align: 'center' }
]

const fetchSchools = async () => {
  try {
    const response = await schoolService.getSchools()
    schools.value = response.data.schools || []
  } catch (error) {
    console.error('Failed to load schools:', error)
  }
}

const fetchReviews = async () => {
  loading.value = true
  try {
    const params = {
      limit: pagination.pageSize,
      offset: (pagination.current - 1) * pagination.pageSize
    }

    if (filters.schoolId) {
      params.school_id = filters.schoolId
    }

    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0].format('YYYY-MM-DD')
      params.end_date = dateRange.value[1].format('YYYY-MM-DD')
    }

    const response = await reviewService.getReviews(params)
    reviews.value = response.reviews || []
    pagination.total = response.total || 0
  } catch (error) {
    message.error('Gagal memuat data ulasan')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchSummary = async () => {
  try {
    const params = {}
    
    if (filters.schoolId) {
      params.school_id = filters.schoolId
    }

    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0].format('YYYY-MM-DD')
      params.end_date = dateRange.value[1].format('YYYY-MM-DD')
    }

    const response = await reviewService.getSummary(params)
    summary.value = response.summary || {}
  } catch (error) {
    console.error('Failed to load summary:', error)
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchReviews()
  fetchSummary()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchReviews()
}

const resetFilters = () => {
  filters.schoolId = undefined
  dateRange.value = []
  pagination.current = 1
  fetchReviews()
  fetchSummary()
}

const showDetail = (record) => {
  selectedReview.value = record
  detailModalVisible.value = true
}

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY HH:mm')
}

const getPhotoUrl = (url) => {
  if (!url) return ''
  // If already a full URL, return as is
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  // Add backend base URL for relative paths
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  // Remove /api/v1 suffix if present
  const cleanBaseUrl = baseUrl.replace(/\/api\/v1$/, '')
  return `${cleanBaseUrl}${url}`
}

onMounted(async () => {
  await fetchSchools()
  await fetchReviews()
  await fetchSummary()
})
</script>

<style scoped>
.review-list {
  padding: 24px;
}

.stat-card {
  text-align: center;
}

.rating-breakdown {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rating-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rating-item > span:first-child {
  width: 140px;
  font-size: 13px;
}

.rating-value {
  font-weight: 600;
  color: #5A4372;
}

.detail-rating {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-rating > div {
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-rating > div > span {
  width: 80px;
  font-size: 13px;
}

:deep(.ant-rate) {
  font-size: 14px;
}

:deep(.ant-rate-star) {
  margin-right: 4px;
}
</style>
