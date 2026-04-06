<template>
  <div class="dashboard-kepala-yayasan">
    <div v-if="loading && !dashboard" class="dashboard-loading">
      <a-skeleton active :paragraph="{ rows: 4 }" />
    </div>

    <template v-else>
      <div class="bento">
        <!-- Hero -->
        <div class="bento__hero">
          <div class="hero-content">
            <div class="hero-left">
              <h2 class="hero-greeting">🏛️ Dashboard Yayasan</h2>
              <p class="hero-subtitle">{{ dashboard?.yayasan_nama || 'Monitoring Yayasan' }}</p>
            </div>
            <div class="hero-metrics">
              <div class="hero-metric">
                <span class="hero-metric__label">SPPG</span>
                <span class="hero-metric__value">{{ dashboard?.total_sppg || 0 }}</span>
              </div>
              <div class="hero-metric">
                <span class="hero-metric__label">Porsi</span>
                <span class="hero-metric__value">{{ dashboard?.aggregated_production?.total_portions || 0 }}</span>
              </div>
              <div class="hero-metric">
                <span class="hero-metric__label">Rating</span>
                <span class="hero-metric__value">{{ (dashboard?.aggregated_review?.average_overall || 0).toFixed(1) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Filters -->
        <div class="bento__filters">
          <div class="filter-row">
            <a-select v-if="needsYayasanSelector" v-model:value="filters.yayasan_id" placeholder="Pilih Yayasan" show-search :filter-option="filterOption" :options="yayasanList.map(y => ({ value: y.id, label: y.kode + ' - ' + y.nama }))" @change="handleFilterChange" style="width:200px;" size="small" />
            <a-select v-model:value="filters.sppg_id" placeholder="Semua SPPG" allow-clear show-search :filter-option="filterOption" :options="sppgOptions" @change="handleFilterChange" style="width:180px;" size="small" />
            <a-range-picker v-model:value="dateRange" format="DD/MM" @change="handleFilterChange" size="small" style="width:200px;" />
            <a-button @click="fetchDashboard" :loading="loading" size="small"><template #icon><ReloadOutlined /></template></a-button>
            <a-button @click="handleExport" :loading="exporting" size="small"><template #icon><DownloadOutlined /></template></a-button>
          </div>
        </div>

        <!-- Metrics: 2 cols -->
        <div class="bento__metrics-wide">
          <h3 class="bento__title">📈 Performa Agregat</h3>
          <div class="metrics-grid">
            <div class="metric-compact">
              <span class="metric-compact__icon" style="background:#FDEAE7;">🍽️</span>
              <div>
                <div class="metric-compact__value">{{ dashboard?.aggregated_production?.total_portions || 0 }}</div>
                <div class="metric-compact__label">Total Porsi</div>
              </div>
            </div>
            <div class="metric-compact">
              <span class="metric-compact__icon" style="background:#D1FAE5;">🚚</span>
              <div>
                <div class="metric-compact__value">{{ dashboard?.aggregated_delivery?.total_deliveries || 0 }}</div>
                <div class="metric-compact__label">Pengiriman</div>
              </div>
            </div>
            <div class="metric-compact">
              <span class="metric-compact__icon" style="background:#FEF3C7;">💰</span>
              <div>
                <div class="metric-compact__value" style="font-size:16px;">{{ formatPercent(dashboard?.aggregated_financial?.absorption_rate) }}</div>
                <div class="metric-compact__label">Penyerapan</div>
              </div>
            </div>
            <div class="metric-compact">
              <span class="metric-compact__icon" style="background:#DBEAFE;">⭐</span>
              <div>
                <div class="metric-compact__value">{{ (dashboard?.aggregated_review?.average_overall || 0).toFixed(1) }}</div>
                <div class="metric-compact__label">Rating</div>
              </div>
            </div>
          </div>
        </div>

        <!-- SPPG Table: full width -->
        <div class="bento__table-full">
          <h3 class="bento__title">🏪 Performa SPPG</h3>
          <a-table :columns="sppgColumns" :data-source="dashboard?.sppg_summaries || []" :loading="loading" :pagination="{ pageSize: 10, size: 'small' }" row-key="sppg_id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'sppg_kode'"><a-tag color="blue" style="font-size:11px;">{{ record.sppg_kode }}</a-tag></template>
              <template v-if="column.key === 'delivery_rate'">{{ formatPercent(record.delivery_rate) }}</template>
              <template v-if="column.key === 'budget_absorption'">{{ formatPercent(record.budget_absorption) }}</template>
              <template v-if="column.key === 'average_review_rating'">{{ record.average_review_rating?.toFixed(1) || '-' }}/5</template>
              <template v-if="column.key === 'action'"><a-button size="small" type="link" @click="drillDownSPPG(record)">Detail →</a-button></template>
            </template>
          </a-table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/stores/auth'
import aggregatedDashboardService from '@/services/aggregatedDashboardService'
import organizationService from '@/services/organizationService'

const authStore = useAuthStore()
const loading = ref(false)
const exporting = ref(false)
const dashboard = ref(null)
const sppgList = ref([])
const yayasanList = ref([])
const dateRange = ref([])
const filters = ref({ sppg_id: null, yayasan_id: null })

// Superadmin/admin_bgn need to select a Yayasan first
const needsYayasanSelector = computed(() => {
  const role = authStore.user?.role
  return role === 'superadmin' || role === 'admin_bgn'
})

const sppgColumns = [
  { title: 'Kode', key: 'sppg_kode', width: 130 },
  { title: 'Nama SPPG', dataIndex: 'sppg_nama', key: 'sppg_nama' },
  { title: 'Total Porsi', dataIndex: 'total_portions', key: 'total_portions', width: 120, align: 'right' },
  { title: 'Tingkat Pengiriman', key: 'delivery_rate', width: 150, align: 'center' },
  { title: 'Penyerapan Anggaran', key: 'budget_absorption', width: 160, align: 'center' },
  { title: 'Rating Review', key: 'average_review_rating', width: 120, align: 'center' },
  { title: '', key: 'action', width: 100 }
]

const sppgOptions = computed(() =>
  sppgList.value.map(s => ({ value: s.id, label: `${s.kode} - ${s.nama}` }))
)

const filterOption = (input, option) =>
  option.label.toLowerCase().includes(input.toLowerCase())

const buildParams = () => {
  const params = {}
  if (filters.value.sppg_id) params.sppg_id = filters.value.sppg_id
  // For superadmin/admin_bgn, pass yayasan_id as query param
  if (needsYayasanSelector.value && filters.value.yayasan_id) {
    params.yayasan_id = filters.value.yayasan_id
  }
  if (dateRange.value && dateRange.value.length === 2) {
    params.start_date = dateRange.value[0].format('YYYY-MM-DD')
    params.end_date = dateRange.value[1].format('YYYY-MM-DD')
  }
  return params
}

const fetchDashboard = async () => {
  // For superadmin/admin_bgn, require yayasan selection first
  if (needsYayasanSelector.value && !filters.value.yayasan_id) {
    dashboard.value = null
    return
  }
  loading.value = true
  try {
    const params = buildParams()
    const res = await aggregatedDashboardService.getKepalaYayasanDashboard(params)
    dashboard.value = res.data?.dashboard || res.data?.data || res.data || null
  } catch (e) {
    message.error('Gagal memuat dashboard Kepala Yayasan')
  } finally {
    loading.value = false
  }
}

const fetchLookups = async () => {
  try {
    const promises = [organizationService.getSPPGList()]
    if (needsYayasanSelector.value) {
      promises.push(organizationService.getYayasanList())
    }
    const results = await Promise.all(promises)
    sppgList.value = results[0].data?.data || results[0].data?.sppgs || []
    if (results[1]) {
      yayasanList.value = results[1].data?.data || results[1].data?.yayasans || []
      // Auto-select first yayasan if only one
      if (yayasanList.value.length === 1 && !filters.value.yayasan_id) {
        filters.value.yayasan_id = yayasanList.value[0].id
      }
    }
  } catch (e) { /* silent */ }
}

const handleFilterChange = () => {
  fetchDashboard()
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params = buildParams()
    const res = await aggregatedDashboardService.exportKepalaYayasanDashboard(params)
    const data = res.data?.data || res.data || {}

    const XLSX = await import('xlsx')

    // Sheet 1: Ringkasan
    const summaryRows = [
      ['Yayasan', data.yayasan_nama || '-'],
      ['Total SPPG', data.total_sppg || 0],
      ['Total Porsi', data.aggregated_production?.total_portions || 0],
      ['Tingkat Penyelesaian', `${(data.aggregated_production?.completion_rate || 0).toFixed(1)}%`],
      ['Total Pengiriman', data.aggregated_delivery?.total_deliveries || 0],
      ['Pengiriman Selesai', data.aggregated_delivery?.completed_deliveries || 0],
      ['On-Time Rate', `${(data.aggregated_delivery?.on_time_rate || 0).toFixed(1)}%`],
      ['Total Anggaran', data.aggregated_financial?.total_budget || 0],
      ['Total Pengeluaran', data.aggregated_financial?.total_spent || 0],
      ['Penyerapan Anggaran', `${(data.aggregated_financial?.absorption_rate || 0).toFixed(1)}%`],
      ['Total Ulasan', data.aggregated_review?.total_reviews || 0],
      ['Rata-rata Rating', (data.aggregated_review?.average_overall || 0).toFixed(1)]
    ]
    const wsSummary = XLSX.utils.aoa_to_sheet([['Metrik', 'Nilai'], ...summaryRows])

    // Sheet 2: Daftar SPPG
    const sppgRows = (data.sppg_summaries || []).map(s => ({
      'Kode': s.sppg_kode,
      'Nama SPPG': s.sppg_nama,
      'Total Porsi': s.total_portions,
      'Delivery Rate': `${(s.delivery_rate || 0).toFixed(1)}%`,
      'Penyerapan Anggaran': `${(s.budget_absorption || 0).toFixed(1)}%`,
      'Rating Review': s.average_review_rating?.toFixed(1)
    }))
    const wsSPPG = XLSX.utils.json_to_sheet(sppgRows)

    const wb = XLSX.utils.book_new()
    XLSX.utils.book_append_sheet(wb, wsSummary, 'Ringkasan')
    XLSX.utils.book_append_sheet(wb, wsSPPG, 'SPPG')

    XLSX.writeFile(wb, `dashboard-yayasan-${new Date().toISOString().slice(0, 10)}.xlsx`)
    message.success('Dashboard berhasil diexport ke Excel')
  } catch (e) {
    message.error('Gagal mengexport dashboard')
  } finally {
    exporting.value = false
  }
}

const drillDownSPPG = (record) => {
  filters.value.sppg_id = record.sppg_id
  fetchDashboard()
}

const formatPercent = (val) => {
  if (val == null) return '-'
  return `${Number(val).toFixed(1)}%`
}

onMounted(() => {
  fetchLookups()
  fetchDashboard()
})
</script>

<style scoped>
.dashboard-kepala-yayasan { min-height: 100%; }
.dashboard-loading { padding: 40px; display: flex; flex-direction: column; align-items: center; gap: 12px; }

.bento { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.bento > * { background: #fff; border-radius: 14px; padding: 20px; border: 1px solid #F0F0F0; min-width: 0; overflow: hidden; }

.bento__hero { grid-column: 1 / -1; border: none; background: linear-gradient(135deg, #C94A3A 0%, #D4553E 50%, #1E8A6E 100%); color: #fff; padding: 24px 28px; }
.bento__filters { grid-column: 1 / -1; padding: 14px 20px; }
.bento__metrics-wide { grid-column: 1 / -1; }
.bento__table-full { grid-column: 1 / -1; }

@media (max-width: 768px) { .bento { grid-template-columns: 1fr; } .bento > * { grid-column: 1 / -1 !important; } }

.hero-content { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
.hero-greeting { font-size: 26px; font-weight: 700; margin: 0 0 4px; color: #fff; }
.hero-subtitle { font-size: 14px; opacity: 0.85; margin: 0; color: rgba(255,255,255,0.85); }
.hero-metrics { display: flex; gap: 12px; }
.hero-metric { background: rgba(255,255,255,0.18); backdrop-filter: blur(4px); border-radius: 10px; padding: 12px 18px; min-width: 90px; text-align: center; }
.hero-metric__label { display: block; font-size: 11px; opacity: 0.75; margin-bottom: 2px; }
.hero-metric__value { display: block; font-size: 26px; font-weight: 700; }
@media (max-width: 768px) { .hero-content { flex-direction: column; align-items: flex-start; } .hero-metrics { flex-wrap: wrap; } }

.bento__title { font-size: 16px; font-weight: 700; margin: 0 0 12px; color: #303030; }
.filter-row { display: flex; gap: 10px; align-items: center; flex-wrap: wrap; }

.metrics-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
@media (max-width: 768px) { .metrics-grid { grid-template-columns: repeat(2, 1fr); } }
.metric-compact { display: flex; align-items: center; gap: 12px; padding: 14px; background: #fff; border: 1px solid #F0F0F0; border-radius: 12px; }
.metric-compact__icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; font-size: 18px; flex-shrink: 0; }
.metric-compact__value { font-size: 22px; font-weight: 700; color: #303030; }
.metric-compact__label { font-size: 12px; color: #6B6B6B; }

.dashboard-kepala-yayasan :deep(.ant-table) { font-size: 13px; }
.dashboard-kepala-yayasan :deep(.ant-table-thead > tr > th) { background: #F7F8FA !important; color: #6B6B6B !important; font-weight: 500; font-size: 11px; padding: 8px 12px; }
.dashboard-kepala-yayasan :deep(.ant-table-tbody > tr > td) { padding: 8px 12px; border-bottom: 1px solid #F0F0F0; }
.dashboard-kepala-yayasan :deep(.ant-table-tbody > tr:hover > td) { background: #F7F8FA !important; }
</style>
