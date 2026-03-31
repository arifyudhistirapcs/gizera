<template>
  <div class="dashboard-bgn">
    <a-page-header title="Dashboard BGN" sub-title="Monitoring agregat nasional" />

    <!-- Filters -->
    <a-card style="margin-bottom: 20px">
      <a-row :gutter="16" align="middle">
        <a-col :xs="24" :sm="8" :md="6">
          <a-form-item label="Yayasan" style="margin-bottom: 0">
            <a-select
              v-model:value="filters.yayasan_id"
              placeholder="Semua Yayasan"
              allow-clear
              show-search
              :filter-option="filterOption"
              :options="yayasanOptions"
              @change="handleFilterChange"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :xs="24" :sm="8" :md="6">
          <a-form-item label="SPPG" style="margin-bottom: 0">
            <a-select
              v-model:value="filters.sppg_id"
              placeholder="Semua SPPG"
              allow-clear
              show-search
              :filter-option="filterOption"
              :options="sppgFilteredOptions"
              @change="handleFilterChange"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :xs="24" :sm="8" :md="6">
          <a-form-item label="Rentang Tanggal" style="margin-bottom: 0">
            <a-range-picker
              v-model:value="dateRange"
              format="DD/MM/YYYY"
              @change="handleFilterChange"
              style="width: 100%"
            />
          </a-form-item>
        </a-col>
        <a-col :xs="24" :sm="8" :md="6">
          <a-space>
            <a-button @click="fetchDashboard" :loading="loading">
              <template #icon><ReloadOutlined /></template>
              Refresh
            </a-button>
            <a-button @click="handleExport" :loading="exporting">
              <template #icon><DownloadOutlined /></template>
              Export
            </a-button>
          </a-space>
        </a-col>
      </a-row>
    </a-card>

    <!-- Peta Sebaran (paling atas setelah filter) -->
    <a-card title="Peta Sebaran Yayasan & SPPG" style="margin-bottom: 20px">
      <OrganizationMap :markers="mapMarkers" :height="450" />
    </a-card>

    <!-- Drill-down indicator -->
    <a-alert
      v-if="filters.yayasan_id"
      type="info"
      show-icon
      closable
      style="margin-bottom: 20px"
      @close="resetFilters"
    >
      <template #message>
        Menampilkan data untuk Yayasan: <strong>{{ selectedYayasanName }}</strong>
        <a-button type="link" size="small" @click="resetFilters">← Kembali ke semua Yayasan</a-button>
      </template>
    </a-alert>

    <!-- Metric Cards -->
    <a-row :gutter="[16, 16]" style="margin-bottom: 20px">
      <a-col :xs="12" :sm="6">
        <a-card :loading="loading">
          <a-statistic
            title="Total Porsi Diproduksi"
            :value="dashboard?.aggregated_production?.total_portions || 0"
            :value-style="{ color: '#3f8600' }"
          />
          <div class="stat-sub">
            Tingkat penyelesaian: {{ formatPercent(dashboard?.aggregated_production?.completion_rate) }}
          </div>
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card :loading="loading">
          <a-statistic
            title="Total Pengiriman"
            :value="dashboard?.aggregated_delivery?.total_deliveries || 0"
          />
          <div class="stat-sub">
            Tepat waktu: {{ formatPercent(dashboard?.aggregated_delivery?.on_time_rate) }}
          </div>
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card :loading="loading">
          <a-statistic
            title="Total Pengeluaran"
            :value="dashboard?.aggregated_financial?.total_spent || 0"
            prefix="Rp"
            :precision="0"
          />
          <div class="stat-sub">
            Penyerapan: {{ formatPercent(dashboard?.aggregated_financial?.absorption_rate) }}
          </div>
        </a-card>
      </a-col>
      <a-col :xs="12" :sm="6">
        <a-card :loading="loading">
          <a-statistic
            title="Rata-rata Review"
            :value="dashboard?.aggregated_review?.average_overall || 0"
            :precision="1"
            suffix="/ 5"
          />
          <div class="stat-sub">
            {{ dashboard?.aggregated_review?.total_reviews || 0 }} ulasan
          </div>
        </a-card>
      </a-col>
    </a-row>

    <!-- Summary Info -->
    <a-row :gutter="16" style="margin-bottom: 20px">
      <a-col :span="12">
        <a-card :loading="loading">
          <a-statistic title="Total Yayasan" :value="dashboard?.total_yayasan || 0" />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card :loading="loading">
          <a-statistic title="Total SPPG" :value="dashboard?.total_sppg || 0" />
        </a-card>
      </a-col>
    </a-row>

    <!-- Yayasan Performance Table -->
    <a-card title="Daftar Yayasan — Ringkasan Performa" style="margin-bottom: 20px">
      <a-table
        :columns="yayasanColumns"
        :data-source="dashboard?.yayasan_summaries || []"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="yayasan_id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'yayasan_kode'">
            <a-tag color="blue">{{ record.yayasan_kode }}</a-tag>
          </template>
          <template v-if="column.key === 'total_spent'">
            {{ formatCurrency(record.total_spent) }}
          </template>
          <template v-if="column.key === 'average_review_rating'">
            {{ record.average_review_rating?.toFixed(1) || '-' }} / 5
          </template>
          <template v-if="column.key === 'action'">
            <a-button size="small" type="link" @click="drillDownYayasan(record)">
              Detail →
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- SPPG Performance Table -->
    <a-card title="Daftar SPPG — Ringkasan Performa" style="margin-bottom: 20px">
      <a-table
        :columns="sppgColumns"
        :data-source="dashboard?.sppg_summaries || []"
        :loading="loading"
        :pagination="{ pageSize: 10 }"
        row-key="sppg_id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'sppg_kode'">
            <a-tag color="green">{{ record.sppg_kode }}</a-tag>
          </template>
          <template v-if="column.key === 'delivery_rate'">
            {{ formatPercent(record.delivery_rate) }}
          </template>
          <template v-if="column.key === 'budget_absorption'">
            {{ formatPercent(record.budget_absorption) }}
          </template>
          <template v-if="column.key === 'average_review_rating'">
            {{ record.average_review_rating?.toFixed(1) || '-' }} / 5
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import aggregatedDashboardService from '@/services/aggregatedDashboardService'
import organizationService from '@/services/organizationService'
import OrganizationMap from '@/components/OrganizationMap.vue'

const loading = ref(false)
const exporting = ref(false)
const dashboard = ref(null)
const yayasanList = ref([])
const sppgList = ref([])
const dateRange = ref([])

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
    const [yRes, sRes] = await Promise.all([
      organizationService.getYayasanList(),
      organizationService.getSPPGList()
    ])
    yayasanList.value = yRes.data?.data || yRes.data?.yayasans || []
    sppgList.value = sRes.data?.data || sRes.data?.sppgs || []
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
  return markers
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
.stat-sub {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}
</style>
