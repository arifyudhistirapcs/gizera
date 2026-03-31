<template>
  <div class="dashboard-kepala-yayasan">
    <a-page-header title="Dashboard Kepala Yayasan" :sub-title="dashboard?.yayasan_nama || ''" />

    <!-- Filters -->
    <a-card style="margin-bottom: 20px">
      <a-row :gutter="16" align="middle">
        <a-col v-if="needsYayasanSelector" :xs="24" :sm="8" :md="6">
          <a-form-item label="Yayasan" style="margin-bottom: 0">
            <a-select
              v-model:value="filters.yayasan_id"
              placeholder="Pilih Yayasan"
              show-search
              :filter-option="filterOption"
              :options="yayasanList.map(y => ({ value: y.id, label: y.kode + ' - ' + y.nama }))"
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
              :options="sppgOptions"
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

    <!-- SPPG Summary -->
    <a-card :loading="loading" style="margin-bottom: 20px">
      <a-statistic title="Total SPPG di bawah Yayasan" :value="dashboard?.total_sppg || 0" />
    </a-card>

    <!-- SPPG Performance Table -->
    <a-card title="Daftar SPPG — Ringkasan Performa">
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
            <a-tag color="blue">{{ record.sppg_kode }}</a-tag>
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
          <template v-if="column.key === 'action'">
            <a-button size="small" type="link" @click="drillDownSPPG(record)">
              Detail →
            </a-button>
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
.stat-sub {
  font-size: 12px;
  color: #888;
  margin-top: 4px;
}
</style>
