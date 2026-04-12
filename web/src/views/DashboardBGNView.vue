<template>
  <div class="dashboard-bgn">
    <div v-if="loading && !dashboard" class="dashboard-loading">
      <a-skeleton active :paragraph="{ rows: 4 }" />
    </div>

    <template v-else>
      <div class="bento">
        <!-- Hero -->
        <div class="bento__hero">
          <div class="hero-content">
            <div class="hero-left">
              <h2 class="hero-greeting">📊 Dashboard BGN</h2>
              <p class="hero-subtitle">Monitoring agregat nasional</p>
            </div>
            <div class="hero-metrics">
              <div class="hero-metric">
                <span class="hero-metric__label">Yayasan</span>
                <span class="hero-metric__value">{{ dashboard?.total_yayasan || 0 }}</span>
              </div>
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

        <!-- Filters: full width -->
        <div class="bento__filters">
          <div class="filter-row">
            <a-select v-model:value="filters.yayasan_id" placeholder="Semua Yayasan" allow-clear show-search :filter-option="filterOption" :options="yayasanOptions" @change="handleFilterChange" style="width:180px;" size="small" />
            <a-select v-model:value="filters.sppg_id" placeholder="Semua SPPG" allow-clear show-search :filter-option="filterOption" :options="sppgFilteredOptions" @change="handleFilterChange" style="width:180px;" size="small" />
            <a-range-picker v-model:value="dateRange" format="DD/MM" @change="handleFilterChange" size="small" style="width:200px;" />
            <a-button @click="fetchDashboard" :loading="loading" size="small"><template #icon><ReloadOutlined /></template></a-button>
            <a-button @click="handleExport" :loading="exporting" size="small"><template #icon><DownloadOutlined /></template></a-button>
          </div>
        </div>

        <!-- Drill-down alert -->
        <div v-if="filters.yayasan_id" class="bento__alert">
          <a-alert type="info" show-icon closable @close="resetFilters">
            <template #message>
              Yayasan: <strong>{{ selectedYayasanName }}</strong>
              <a-button type="link" size="small" @click="resetFilters">← Reset</a-button>
            </template>
          </a-alert>
        </div>

        <!-- Peta: 2 cols -->
        <div class="bento__map">
          <div class="map-header">
            <h3 class="bento__title">🗺️ Peta Sebaran</h3>
            <div class="map-legend">
              <label v-for="opt in mapLayerOptions" :key="opt.value" class="map-legend__item" @click="toggleLayer(opt.value)">
                <span class="map-legend__dot" :style="{ background: mapLayers.includes(opt.value) ? opt.color : '#D8D8DB' }"></span>
                <span class="map-legend__label">{{ opt.label }}</span>
              </label>
            </div>
          </div>
          <OrganizationMap :markers="filteredMapMarkers" :height="280" />
        </div>

        <!-- Metrics: 1 col -->
        <div class="bento__metrics">
          <h3 class="bento__title">📈 Performa</h3>
          <div class="compact-stat">
            <span class="compact-stat__label">Total Pengiriman</span>
            <span class="compact-stat__value">{{ dashboard?.aggregated_delivery?.total_deliveries || 0 }}</span>
          </div>
          <div class="compact-stat">
            <span class="compact-stat__label">Tepat Waktu</span>
            <span class="compact-stat__value">{{ formatPercent(dashboard?.aggregated_delivery?.on_time_rate) }}</span>
          </div>
          <div class="compact-stat">
            <span class="compact-stat__label">Pengeluaran</span>
            <span class="compact-stat__value" style="font-size:18px;">{{ formatCurrency(dashboard?.aggregated_financial?.total_spent || 0) }}</span>
          </div>
          <div class="compact-stat">
            <span class="compact-stat__label">Penyerapan</span>
            <span class="compact-stat__value">{{ formatPercent(dashboard?.aggregated_financial?.absorption_rate) }}</span>
          </div>
        </div>

        <!-- Yayasan Table: 2 cols -->
        <div class="bento__table1">
          <h3 class="bento__title">🏛️ Performa Yayasan</h3>
          <a-table :columns="yayasanColumns" :data-source="dashboard?.yayasan_summaries || []" :loading="loading" :pagination="{ pageSize: 5, size: 'small' }" row-key="yayasan_id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'yayasan_kode'"><a-tag color="blue" style="font-size:11px;">{{ record.yayasan_kode }}</a-tag></template>
              <template v-if="column.key === 'total_spent'">{{ formatCurrency(record.total_spent) }}</template>
              <template v-if="column.key === 'average_review_rating'">{{ record.average_review_rating?.toFixed(1) || '-' }}/5</template>
              <template v-if="column.key === 'action'"><a-button size="small" type="link" @click="drillDownYayasan(record)">Detail →</a-button></template>
            </template>
          </a-table>
        </div>

        <!-- SPPG Table: 1 col -->
        <div class="bento__table2">
          <h3 class="bento__title">🏪 Performa SPPG</h3>
          <a-table :columns="sppgColumns" :data-source="dashboard?.sppg_summaries || []" :loading="loading" :pagination="{ pageSize: 5, size: 'small' }" row-key="sppg_id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'sppg_kode'"><a-tag color="green" style="font-size:11px;">{{ record.sppg_kode }}</a-tag></template>
              <template v-if="column.key === 'delivery_rate'">{{ formatPercent(record.delivery_rate) }}</template>
              <template v-if="column.key === 'budget_absorption'">{{ formatPercent(record.budget_absorption) }}</template>
              <template v-if="column.key === 'average_review_rating'">{{ record.average_review_rating?.toFixed(1) || '-' }}/5</template>
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
import aggregatedDashboardService from '@/services/aggregatedDashboardService'
import organizationService from '@/services/organizationService'
import schoolService from '@/services/schoolService'
import supplierService from '@/services/supplierService'
import OrganizationMap from '@/components/OrganizationMap.vue'
import LottiePlayer from '@/components/common/LottiePlayer.vue'

const loading = ref(false)
const exporting = ref(false)
const dashboard = ref(null)
const yayasanList = ref([])
const sppgList = ref([])
const schoolList = ref([])
const supplierList = ref([])
const dateRange = ref([])

const mapLayers = ref(['Yayasan', 'SPPG', 'Sekolah', 'Supplier'])

const toggleLayer = (layer) => {
  const idx = mapLayers.value.indexOf(layer)
  if (idx > -1) mapLayers.value.splice(idx, 1)
  else mapLayers.value.push(layer)
}

const mapLayerOptions = [
  { label: 'Yayasan', value: 'Yayasan', color: '#303030' },
  { label: 'SPPG', value: 'SPPG', color: '#52c41a' },
  { label: 'Sekolah', value: 'Sekolah', color: '#3B82F6' },
  { label: 'Supplier', value: 'Supplier', color: '#fa8c16' }
]

const filters = ref({
  yayasan_id: null,
  sppg_id: null
})

const yayasanColumns = [
  { title: 'Kode', key: 'yayasan_kode', width: 120 },
  { title: 'Nama Yayasan', dataIndex: 'yayasan_nama', key: 'yayasan_nama' },
  { title: 'Jumlah SPPG', dataIndex: 'total_sppg', key: 'total_sppg', width: 120, align: 'center' },
  { title: 'Total Porsi', dataIndex: 'total_portions', key: 'total_portions', width: 120, align: 'right' },
  { title: 'Total Pengeluaran', key: 'total_spent', width: 160, align: 'right' },
  { title: 'Rating Review', key: 'average_review_rating', width: 120, align: 'center' },
  { title: '', key: 'action', width: 100 }
]

const sppgColumns = [
  { title: 'Kode', key: 'sppg_kode', width: 130 },
  { title: 'Nama SPPG', dataIndex: 'sppg_nama', key: 'sppg_nama' },
  { title: 'Total Porsi', dataIndex: 'total_portions', key: 'total_portions', width: 120, align: 'right' },
  { title: 'Delivery Rate', key: 'delivery_rate', width: 130, align: 'center' },
  { title: 'Penyerapan Anggaran', key: 'budget_absorption', width: 160, align: 'center' },
  { title: 'Rating Review', key: 'average_review_rating', width: 120, align: 'center' }
]

const yayasanOptions = computed(() =>
  yayasanList.value.map(y => ({ value: y.id, label: `${y.kode} - ${y.nama}` }))
)

const sppgFilteredOptions = computed(() => {
  let list = sppgList.value
  if (filters.value.yayasan_id) {
    list = list.filter(s => s.yayasan_id === filters.value.yayasan_id)
  }
  return list.map(s => ({ value: s.id, label: `${s.kode} - ${s.nama}` }))
})

const filterOption = (input, option) =>
  option.label.toLowerCase().includes(input.toLowerCase())

const buildParams = () => {
  const params = {}
  if (filters.value.yayasan_id) params.yayasan_id = filters.value.yayasan_id
  if (filters.value.sppg_id) params.sppg_id = filters.value.sppg_id
  if (dateRange.value && dateRange.value.length === 2) {
    params.start_date = dateRange.value[0].format('YYYY-MM-DD')
    params.end_date = dateRange.value[1].format('YYYY-MM-DD')
  }
  return params
}

const fetchDashboard = async () => {
  loading.value = true
  try {
    const params = buildParams()
    const res = await aggregatedDashboardService.getAdminBGNDashboard(params)
    dashboard.value = res.data?.dashboard || res.data?.data || res.data || null
  } catch (e) {
    message.error('Gagal memuat dashboard BGN')
  } finally {
    loading.value = false
  }
}

const fetchLookups = async () => {
  try {
    const [yRes, sRes, schRes, supRes] = await Promise.all([
      organizationService.getYayasanList(),
      organizationService.getSPPGList(),
      schoolService.getSchools().catch(() => ({ data: { data: [] } })),
      supplierService.getSuppliers().catch(() => ({ data: { data: [] } }))
    ])
    yayasanList.value = yRes.data?.data || yRes.data?.yayasans || []
    sppgList.value = sRes.data?.data || sRes.data?.sppgs || []
    schoolList.value = schRes.data?.data || schRes.data?.schools || []
    supplierList.value = supRes.data?.data || supRes.data?.suppliers || []
  } catch (e) {
    // silent
  }
}

const handleFilterChange = () => {
  fetchDashboard()
}

const handleExport = async () => {
  exporting.value = true
  try {
    const params = buildParams()
    const res = await aggregatedDashboardService.exportAdminBGNDashboard(params)
    const data = res.data?.data || res.data || {}

    // Build Excel workbook
    const XLSX = await import('xlsx')

    // Sheet 1: Ringkasan
    const summaryRows = [
      ['Total Yayasan', data.total_yayasan || 0],
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

    // Sheet 2: Daftar Yayasan
    const yayasanRows = (data.yayasan_summaries || []).map(y => ({
      'Kode': y.yayasan_kode,
      'Nama Yayasan': y.yayasan_nama,
      'Jumlah SPPG': y.total_sppg,
      'Total Porsi': y.total_portions,
      'Total Pengeluaran': y.total_spent,
      'Rating Review': y.average_review_rating?.toFixed(1)
    }))
    const wsYayasan = XLSX.utils.json_to_sheet(yayasanRows)

    // Sheet 3: Daftar SPPG
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
    XLSX.utils.book_append_sheet(wb, wsYayasan, 'Yayasan')
    XLSX.utils.book_append_sheet(wb, wsSPPG, 'SPPG')

    XLSX.writeFile(wb, `dashboard-bgn-${new Date().toISOString().slice(0, 10)}.xlsx`)
    message.success('Dashboard berhasil diexport ke Excel')
  } catch (e) {
    message.error('Gagal mengexport dashboard')
  } finally {
    exporting.value = false
  }
}

const drillDownYayasan = (record) => {
  filters.value.yayasan_id = record.yayasan_id
  filters.value.sppg_id = null
  fetchDashboard()
}

const resetFilters = () => {
  filters.value.yayasan_id = null
  filters.value.sppg_id = null
  fetchDashboard()
}

const mapMarkers = computed(() => {
  const markers = []
  const summaries = dashboard.value?.yayasan_summaries || []
  const selectedYayasan = filters.value.yayasan_id
  const selectedSPPG = filters.value.sppg_id

  // Yayasan markers — filter by selected yayasan
  yayasanList.value.forEach(y => {
    if (y.latitude && y.longitude) {
      if (selectedYayasan && y.id !== selectedYayasan) return
      const summary = summaries.find(s => s.yayasan_id === y.id)
      markers.push({
        id: `yayasan-${y.id}`,
        name: y.nama,
        kode: y.kode,
        type: 'Yayasan',
        latitude: y.latitude,
        longitude: y.longitude,
        details: y.penanggung_jawab ? `PJ: ${y.penanggung_jawab}` : '',
        stats: summary ? {
          sppgCount: summary.total_sppg,
          totalPortions: summary.total_portions,
          totalSpent: summary.total_spent,
          rating: summary.average_review_rating
        } : null
      })
    }
  })

  // SPPG markers — filter by selected yayasan and/or sppg
  sppgList.value.forEach(s => {
    if (s.latitude && s.longitude) {
      if (selectedYayasan && s.yayasan_id !== selectedYayasan) return
      if (selectedSPPG && s.id !== selectedSPPG) return
      const sppgSummary = (dashboard.value?.sppg_summaries || []).find(ss => ss.sppg_id === s.id)
      markers.push({
        id: `sppg-${s.id}`,
        name: s.nama,
        kode: s.kode,
        type: 'SPPG',
        latitude: s.latitude,
        longitude: s.longitude,
        details: s.yayasan?.nama ? `Yayasan: ${s.yayasan.nama}` : '',
        stats: sppgSummary ? {
          totalPortions: sppgSummary.total_portions,
          deliveryRate: sppgSummary.delivery_rate,
          budgetAbsorption: sppgSummary.budget_absorption,
          rating: sppgSummary.average_review_rating
        } : null
      })
    }
  })
  
  // School markers
  schoolList.value.forEach(sch => {
    if (sch.latitude && sch.longitude) {
      if (selectedSPPG && sch.sppg_id !== selectedSPPG) return
      markers.push({
        id: `school-${sch.id}`,
        name: sch.name,
        kode: sch.npsn || '',
        type: 'Sekolah',
        latitude: sch.latitude,
        longitude: sch.longitude,
        details: sch.address || '',
        stats: {
          portionsSmall: sch.student_count_grade_1_3 || 0,
          portionsLarge: sch.student_count_grade_4_6 || 0,
          totalDelivered: sch.student_count || 0,
          rating: 0
        }
      })
    }
  })

  // Supplier markers
  supplierList.value.forEach(sup => {
    if (sup.latitude && sup.longitude) {
      markers.push({
        id: `supplier-${sup.id}`,
        name: sup.name,
        kode: sup.product_category || '',
        type: 'Supplier',
        latitude: sup.latitude,
        longitude: sup.longitude,
        details: sup.address || '',
        stats: {
          onTimeDelivery: sup.on_time_delivery || 0,
          qualityRating: sup.quality_rating || 0,
          contactPerson: sup.contact_person || ''
        }
      })
    }
  })

  return markers
})

const filteredMapMarkers = computed(() => {
  return mapMarkers.value.filter(m => mapLayers.value.includes(m.type))
})

const selectedYayasanName = computed(() => {
  if (!filters.value.yayasan_id) return ''
  const y = yayasanList.value.find(y => y.id === filters.value.yayasan_id)
  return y ? `${y.kode} - ${y.nama}` : `ID ${filters.value.yayasan_id}`
})

const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(val)
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
.dashboard-bgn { min-height: 100%; }
.dashboard-loading { padding: 40px; display: flex; flex-direction: column; align-items: center; gap: 12px; }

.bento { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.bento > * { background: #fff; border-radius: 14px; padding: 20px; border: 1px solid #F0F0F0; min-width: 0; overflow: hidden; }

.bento__hero { grid-column: 1 / -1; border: none; background: #F82C17; color: #fff; padding: 24px 28px; }
.bento__filters { grid-column: 1 / -1; padding: 14px 20px; }
.bento__alert { grid-column: 1 / -1; padding: 0; background: transparent; border: none; }
.bento__map { grid-column: 1 / 3; }
.bento__metrics { grid-column: 3 / 4; }
.bento__table1 { grid-column: 1 / 3; }
.bento__table2 { grid-column: 3 / 4; }

@media (max-width: 1024px) {
  .bento { grid-template-columns: 1fr 1fr; }
  .bento__hero, .bento__filters, .bento__alert, .bento__map, .bento__metrics, .bento__table1, .bento__table2 { grid-column: 1 / -1; }
}
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

.map-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.map-header .bento__title { margin-bottom: 0; }
.map-legend { display: flex; gap: 12px; }
.map-legend__item { display: flex; align-items: center; gap: 5px; cursor: pointer; font-size: 12px; }
.map-legend__dot { width: 10px; height: 10px; border-radius: 50%; transition: background 0.2s; }
.map-legend__label { font-weight: 500; color: #303030; }

.compact-stat { display: flex; justify-content: space-between; align-items: baseline; padding: 10px 0; border-bottom: 1px solid #F0F0F0; }
.compact-stat:last-child { border-bottom: none; }
.compact-stat__label { font-size: 13px; color: #6B6B6B; }
.compact-stat__value { font-size: 22px; font-weight: 700; color: #303030; }

.dashboard-bgn :deep(.ant-table) { font-size: 13px; }
.dashboard-bgn :deep(.ant-table-thead > tr > th) { background: #F7F8FA !important; color: #6B6B6B !important; font-weight: 500; font-size: 11px; padding: 8px 12px; }
.dashboard-bgn :deep(.ant-table-tbody > tr > td) { padding: 8px 12px; border-bottom: 1px solid #F0F0F0; }
.dashboard-bgn :deep(.ant-table-tbody > tr:hover > td) { background: #F7F8FA !important; }
</style>
