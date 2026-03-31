<template>
  <div class="risk-assessment-list">
    <a-page-header
      title="Risk Assessment"
      sub-title="Audit Kepatuhan SOP SPPG"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Summary Statistics -->
        <a-row :gutter="16" v-if="stats.length > 0">
          <a-col :span="6" v-for="stat in stats" :key="stat.sppg_id">
            <a-card size="small" class="stat-card">
              <a-statistic
                :title="stat.sppg_name || `SPPG #${stat.sppg_id}`"
                :value="stat.average_score || 0"
                :precision="1"
                :value-style="{ color: scoreColor(stat.average_score) }"
              >
                <template #suffix>
                  <a-tag :color="riskLevelColor(stat.risk_level_trend)" size="small">
                    {{ stat.risk_level_trend || '-' }}
                  </a-tag>
                </template>
              </a-statistic>
              <div class="stat-sub">{{ stat.total_audits || 0 }} audit</div>
            </a-card>
          </a-col>
        </a-row>

        <!-- Filter Section -->
        <a-card size="small" title="Filter">
          <a-form layout="inline" :model="filters">
            <a-form-item label="SPPG">
              <a-select
                v-model:value="filters.sppgId"
                placeholder="Semua SPPG"
                style="width: 200px"
                allow-clear
                show-search
                option-filter-prop="label"
                :options="sppgOptions"
              />
            </a-form-item>

            <a-form-item label="Periode">
              <a-range-picker
                v-model:value="dateRange"
                format="DD/MM/YYYY"
                style="width: 280px"
              />
            </a-form-item>

            <a-form-item label="Risk Level">
              <a-select
                v-model:value="filters.riskLevel"
                placeholder="Semua Level"
                style="width: 150px"
                allow-clear
                :options="riskLevelOptions"
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

        <!-- Forms Table -->
        <a-table
          :columns="columns"
          :data-source="forms"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'sppg'">
              {{ record.sppg?.nama || `SPPG #${record.sppg_id}` }}
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDate(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'overall_risk_score'">
              {{ record.overall_risk_score != null ? record.overall_risk_score.toFixed(1) : '-' }}
            </template>
            <template v-else-if="column.key === 'risk_level'">
              <a-tag :color="riskLevelColor(record.risk_level)">
                {{ record.risk_level || '-' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'status'">
              <a-tag :color="record.status === 'submitted' ? 'blue' : 'default'">
                {{ record.status }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'actions'">
              <router-link :to="`/risk-assessment/${record.id}`">
                <a-button type="link" size="small">
                  <EyeOutlined /> Detail
                </a-button>
              </router-link>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, ClearOutlined, EyeOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import riskAssessmentService from '@/services/riskAssessmentService'

const loading = ref(false)
const forms = ref([])
const stats = ref([])
const sppgList = ref([])
const dateRange = ref([])

const filters = reactive({
  sppgId: undefined,
  riskLevel: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total, range) => `${range[0]}-${range[1]} dari ${total} data`
})

const columns = [
  { title: 'SPPG', key: 'sppg', width: 200 },
  { title: 'Tanggal Audit', key: 'created_at', width: 140 },
  { title: 'Skor Risiko', key: 'overall_risk_score', width: 120, align: 'center' },
  { title: 'Risk Level', key: 'risk_level', width: 120, align: 'center' },
  { title: 'Status', key: 'status', width: 100, align: 'center' },
  { title: 'Aksi', key: 'actions', width: 100, align: 'center' }
]

const riskLevelOptions = [
  { value: 'rendah', label: 'Rendah' },
  { value: 'sedang', label: 'Sedang' },
  { value: 'tinggi', label: 'Tinggi' }
]

const sppgOptions = ref([])

const riskLevelColor = (level) => {
  if (!level) return 'default'
  switch (level.toLowerCase()) {
    case 'rendah': return 'green'
    case 'sedang': return 'orange'
    case 'tinggi': return 'red'
    default: return 'default'
  }
}

const scoreColor = (score) => {
  if (score == null) return '#999'
  if (score >= 4.0) return '#52c41a'
  if (score >= 2.5) return '#faad14'
  return '#f5222d'
}

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY HH:mm')
}

const fetchSPPGs = async () => {
  try {
    const response = await riskAssessmentService.getSPPGList()
    const list = response.data?.sppgs || response.data || []
    sppgList.value = list
    sppgOptions.value = list.map(s => ({ value: s.id, label: s.nama || s.name }))
  } catch (error) {
    console.error('Failed to load SPPGs:', error)
  }
}

const fetchForms = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize
    }

    if (filters.sppgId) params.sppg_id = filters.sppgId
    if (filters.riskLevel) params.risk_level = filters.riskLevel
    if (dateRange.value?.length === 2) {
      params.date_from = dateRange.value[0].format('YYYY-MM-DD')
      params.date_to = dateRange.value[1].format('YYYY-MM-DD')
    }

    const response = await riskAssessmentService.getForms(params)
    const data = response.data || response
    forms.value = data.forms || data.data || []
    pagination.total = data.total || 0
  } catch (error) {
    message.error('Gagal memuat data risk assessment')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const ids = sppgList.value.map(s => s.id)
    if (ids.length === 0) return
    const response = await riskAssessmentService.getStats(ids)
    const data = response.data || response
    stats.value = data.stats || data.data || []
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchForms()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchForms()
}

const resetFilters = () => {
  filters.sppgId = undefined
  filters.riskLevel = undefined
  dateRange.value = []
  pagination.current = 1
  fetchForms()
  fetchStats()
}

onMounted(async () => {
  await fetchSPPGs()
  await Promise.all([fetchForms(), fetchStats()])
})
</script>

<style scoped>
.risk-assessment-list {
  padding: 24px;
}

.stat-card {
  text-align: center;
}

.stat-sub {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>
