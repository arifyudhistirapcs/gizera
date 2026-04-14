<template>
  <div class="dashboard-sspg">
    <!-- Loading -->
    <div v-if="loading" class="dashboard-loading">
      <a-skeleton active :paragraph="{ rows: 4 }" />
    </div>

    <template v-else>
      <!-- BENTO GRID -->
      <div class="bento">

        <!-- Hero: spans full width -->
        <div class="bento__hero">
          <div class="hero-content">
            <div class="hero-left">
              <h2 class="hero-greeting">Selamat {{ getGreeting() }}, {{ userName }}! 👋</h2>
              <p class="hero-subtitle">{{ todayFormatted }}</p>
            </div>
            <img src="@/assets/illustrations/hero-cooking.svg" alt="" class="hero-illustration" />
            <div class="hero-metrics">
              <div class="hero-metric">
                <span class="hero-metric__label">Total Porsi</span>
                <span class="hero-metric__value">{{ kpis.portions_prepared }}</span>
              </div>
              <div class="hero-metric">
                <span class="hero-metric__label">Pengiriman</span>
                <span class="hero-metric__value">{{ kpis.delivery_rate }}%</span>
              </div>
              <div class="hero-metric">
                <span class="hero-metric__label">Tepat Waktu</span>
                <span class="hero-metric__value">{{ kpis.on_time_delivery_rate }}%</span>
              </div>
              <div class="hero-metric">
                <span class="hero-metric__label">Ketersediaan Stok</span>
                <span class="hero-metric__value">{{ kpis.stock_availability }}%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Peta: large, spans 2 cols -->
        <div class="bento__map">
          <div class="map-header">
            <h3 class="bento__title">🗺️ Peta Sebaran</h3>
            <div class="map-legend">
              <label class="map-legend__item" @click="toggleMapLayer('Sekolah')">
                <span class="map-legend__dot" :style="{ background: mapLayers.includes('Sekolah') ? '#3B82F6' : '#D8D8DB' }"></span>
                <span class="map-legend__label">Sekolah ({{ schoolsList.length }})</span>
              </label>
              <label class="map-legend__item" @click="toggleMapLayer('Supplier')">
                <span class="map-legend__dot" :style="{ background: mapLayers.includes('Supplier') ? '#fa8c16' : '#D8D8DB' }"></span>
                <span class="map-legend__label">Supplier ({{ suppliersList.length }})</span>
              </label>
            </div>
          </div>
          <OrganizationMap :markers="filteredSppgMapMarkers" :height="280" />
        </div>

        <!-- Rating: 1 col, stacked -->
        <div class="bento__rating">
          <h3 class="bento__title">⭐ Ulasan</h3>
          <div class="compact-stat">
            <span class="compact-stat__label">Keseluruhan</span>
            <span class="compact-stat__value">{{ reviewSummary.average_overall_rating?.toFixed(1) || '0.0' }}<small>/5</small></span>
          </div>
          <div class="compact-stat">
            <span class="compact-stat__label">Menu</span>
            <span class="compact-stat__value">{{ reviewSummary.average_menu_rating?.toFixed(1) || '0.0' }}<small>/5</small></span>
          </div>
          <div class="compact-stat">
            <span class="compact-stat__label">Layanan</span>
            <span class="compact-stat__value">{{ reviewSummary.average_service_rating?.toFixed(1) || '0.0' }}<small>/5</small></span>
          </div>
          <div class="compact-stat__footer">{{ reviewSummary.total_reviews || 0 }} ulasan total</div>
        </div>

        <!-- Chart Produksi: 1 col -->
        <div class="bento__chart1">
          <h3 class="bento__title">🍳 Produksi</h3>
          <div ref="productionChartRef" style="width:100%;height:220px;"></div>
        </div>

        <!-- Chart Pengiriman: 1 col -->
        <div class="bento__chart2">
          <h3 class="bento__title">🚚 Pengiriman</h3>
          <div ref="deliveryChartRef" style="width:100%;height:220px;"></div>
        </div>

        <!-- Chart Pencucian: 1 col -->
        <div class="bento__chart3">
          <h3 class="bento__title">🧹 Pencucian</h3>
          <div ref="cleaningChartRef" style="width:100%;height:220px;"></div>
        </div>

        <!-- Arus Kas: spans 2 cols -->
        <div class="bento__cashflow">
          <div class="bento__title-row">
            <h3 class="bento__title">💰 Arus Kas</h3>
            <a-range-picker
              v-model:value="cashFlowDateRange"
              format="DD/MM"
              @change="handleCashFlowDateChange"
              size="small"
              style="width: 200px;"
            />
          </div>
          <div class="cashflow-grid">
            <div class="cashflow-item cashflow-item--income">
              <span class="cashflow-item__label">Pemasukan</span>
              <span class="cashflow-item__value">{{ formatCurrency(cashFlowSummary.total_income || 0) }}</span>
            </div>
            <div class="cashflow-item cashflow-item--expense">
              <span class="cashflow-item__label">Pengeluaran</span>
              <span class="cashflow-item__value">{{ formatCurrency(cashFlowSummary.total_expense || 0) }}</span>
            </div>
            <div class="cashflow-item cashflow-item--net">
              <span class="cashflow-item__label">Bersih</span>
              <span class="cashflow-item__value" :class="{ 'text-error': (cashFlowSummary.net_cash_flow || 0) < 0 }">{{ formatCurrency(cashFlowSummary.net_cash_flow || 0) }}</span>
            </div>
          </div>
        </div>

        <!-- Supplier: 1 col -->
        <div class="bento__supplier">
          <div class="bento__title-row">
            <h3 class="bento__title">🏪 Top Supplier</h3>
            <a-button type="link" size="small" @click="goToSuppliers" style="color:#764AF1;padding:0;">Semua →</a-button>
          </div>
          <div v-if="topSuppliers.length === 0" style="color:#6B6B6B;font-size:13px;padding:12px 0;">Belum ada data</div>
          <div v-else class="supplier-compact-list">
            <div v-for="(s, i) in topSuppliers.slice(0, 3)" :key="s.id" class="supplier-compact">
              <span class="supplier-compact__rank" :class="`rank-${i+1}`">{{ i+1 }}</span>
              <span class="supplier-compact__name">{{ s.name }}</span>
              <span class="supplier-compact__amount">{{ formatCurrency(s.total_amount) }}</span>
            </div>
          </div>
        </div>

        <!-- Stok Kritis: spans full width, compact -->
        <div v-if="criticalStockItems.length > 0" class="bento__stock">
          <div class="bento__title-row">
            <h3 class="bento__title">⚠️ Stok Kritis ({{ criticalStockItems.length }})</h3>
            <a-button v-if="criticalStockItems.length > 4" type="link" size="small" @click="goToInventory" style="color:#764AF1;padding:0;">Lihat Semua →</a-button>
          </div>
          <div class="stock-compact-grid">
            <div v-for="item in criticalStockItems.slice(0, 4)" :key="item.ingredient_id" class="stock-compact">
              <div class="stock-compact__info">
                <span class="stock-compact__name">{{ item.ingredient_name }}</span>
                <span class="stock-compact__qty">{{ item.current_stock }} {{ item.unit }}</span>
              </div>
              <a-progress
                :percent="Math.min(100, Math.round((item.current_stock / item.min_threshold) * 100))"
                :stroke-color="item.current_stock <= item.min_threshold * 0.5 ? '#F82C17' : '#D97706'"
                :show-info="false"
                size="small"
                style="width:80px;"
              />
            </div>
          </div>
        </div>
        <div v-else class="bento__stock-ok">
          ✅ Semua stok dalam kondisi aman
        </div>

      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  AppstoreOutlined,
  CarOutlined,
  InboxOutlined,
  CheckCircleOutlined,
  WarningOutlined,
  RightOutlined,
  StarOutlined,
  CoffeeOutlined,
  ArrowUpOutlined,
  ArrowDownOutlined,
  DollarOutlined
} from '@ant-design/icons-vue'
import HStatCard from '@/components/horizon/HStatCard.vue'
import HChartCard from '@/components/horizon/HChartCard.vue'
import HDataTable from '@/components/horizon/HDataTable.vue'
import { useHorizonChart } from '@/composables/useHorizonChart'
import { getKepalaSSPGDashboard } from '@/services/dashboardService'
import reviewService from '@/services/reviewService'
import supplierService from '@/services/supplierService'
import cashFlowService from '@/services/cashFlowService'
import { database, firebasePaths } from '@/services/firebase'
import { ref as dbRef, onValue, off } from 'firebase/database'
import dayjs from 'dayjs'
import OrganizationMap from '@/components/OrganizationMap.vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()
const dashboard = ref(null)
const loading = ref(true)
const loadingCashFlow = ref(false)
const reviewSummary = ref({
  total_reviews: 0,
  average_overall_rating: 0,
  average_menu_rating: 0,
  average_service_rating: 0
})
const topSuppliers = ref([])
const cashFlowSummary = ref({
  total_income: 0,
  total_expense: 0,
  net_cash_flow: 0
})
const cashFlowDateRange = ref([])

const productionChartRef = ref(null)
const deliveryChartRef = ref(null)
const cleaningChartRef = ref(null)

let dashboardListener = null

// Map layer toggles
const mapLayers = ref(['Sekolah', 'Supplier'])

const toggleMapLayer = (layer) => {
  const idx = mapLayers.value.indexOf(layer)
  if (idx > -1) {
    mapLayers.value.splice(idx, 1)
  } else {
    mapLayers.value.push(layer)
  }
}

const filteredSppgMapMarkers = computed(() => {
  return sppgMapMarkers.value.filter(m => mapLayers.value.includes(m.type))
})

// Hero section helpers
const userName = computed(() => authStore.user?.full_name || 'User')

const todayFormatted = computed(() => {
  return new Date().toLocaleDateString('id-ID', {
    weekday: 'long', day: 'numeric', month: 'long', year: 'numeric'
  })
})

const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 12) return 'Pagi'
  if (hour < 15) return 'Siang'
  if (hour < 18) return 'Sore'
  return 'Malam'
}

// Computed helpers for safe access
const kpis = computed(() => dashboard.value?.today_kpis || {
  portions_prepared: 0, delivery_rate: 0, stock_availability: 0, on_time_delivery_rate: 0
})

const production = computed(() => dashboard.value?.production_status || {
  total_recipes: 0, recipes_pending: 0, recipes_cooking: 0, recipes_ready: 0,
  packing_pending: 0, packing_in_progress: 0, packing_ready: 0, completion_rate: 0
})

const delivery = computed(() => dashboard.value?.delivery_status || {
  total_deliveries: 0, status_breakdown: [], completion_rate: 0
})

const cleaning = computed(() => dashboard.value?.cleaning_status || {
  total_items: 0, items_pending: 0, items_in_progress: 0, items_completed: 0, completion_rate: 0
})

const criticalStockItems = computed(() => dashboard.value?.critical_stock || [])

// Display only first 6 critical stock items
const displayedCriticalStock = computed(() => criticalStockItems.value.slice(0, 6))

// Computed text for cash flow date range
const cashFlowDateRangeText = computed(() => {
  if (cashFlowDateRange.value && cashFlowDateRange.value.length === 2) {
    const start = cashFlowDateRange.value[0].format('DD/MM/YYYY')
    const end = cashFlowDateRange.value[1].format('DD/MM/YYYY')
    return `${start} - ${end}`
  }
  return 'Pilih rentang tanggal'
})

// Navigate to inventory with low stock filter
const goToInventory = () => {
  router.push('/inventory')
}

// Navigate to suppliers page
const goToSuppliers = () => {
  router.push('/suppliers')
}

// Format currency
const formatCurrency = (value) => {
  if (!value) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

// Get rank color for progress bar
const getRankColor = (index) => {
  const colors = ['#FFD700', '#C0C0C0', '#CD7F32', '#303030', '#6B6B6B']
  return colors[index] || '#6B6B6B'
}

// Load cash flow summary
const loadCashFlowSummary = async () => {
  if (!cashFlowDateRange.value || cashFlowDateRange.value.length !== 2) return
  
  loadingCashFlow.value = true
  try {
    const startDate = cashFlowDateRange.value[0].format('YYYY-MM-DD')
    const endDate = cashFlowDateRange.value[1].format('YYYY-MM-DD')
    
    const response = await cashFlowService.getCashFlowSummary(startDate, endDate)
    if (response.summary) {
      cashFlowSummary.value = response.summary
    }
  } catch (error) {
    console.warn('Could not load cash flow summary:', error)
  } finally {
    loadingCashFlow.value = false
  }
}

// Handle cash flow date range change
const handleCashFlowDateChange = () => {
  loadCashFlowSummary()
}

// Set default cash flow date range to current month
const setDefaultCashFlowDateRange = () => {
  const startOfMonth = dayjs().startOf('month')
  const endOfMonth = dayjs().endOf('month')
  cashFlowDateRange.value = [startOfMonth, endOfMonth]
}

// Production table data from real API
const productionColumns = [
  { title: 'Sekolah', dataIndex: 'school', key: 'school' },
  { title: 'Porsi', dataIndex: 'portions', key: 'portions' },
  { title: 'Status', dataIndex: 'status', key: 'status', type: 'status' }
]

const productionTableData = computed(() => {
  const details = dashboard.value?.production_details || []
  return details.map((detail, index) => ({
    key: String(index + 1),
    school: detail.school_name,
    portions: detail.portions,
    status: detail.status
  }))
})

// Delivery table data from real API
const deliveryColumns = [
  { title: 'Sekolah', dataIndex: 'school', key: 'school' },
  { title: 'Porsi', dataIndex: 'portions', key: 'portions' },
  { title: 'Status', dataIndex: 'status', key: 'status', type: 'status' }
]

const deliveryTableData = computed(() => {
  const details = dashboard.value?.delivery_details || []
  return details.map((detail, index) => ({
    key: String(index + 1),
    school: detail.school_name,
    portions: detail.portions,
    status: detail.status
  }))
})

// Cleaning table data from real API
const cleaningColumns = [
  { title: 'Sekolah', dataIndex: 'school', key: 'school' },
  { title: 'Porsi', dataIndex: 'portions', key: 'portions' },
  { title: 'Status', dataIndex: 'status', key: 'status', type: 'status' }
]

const cleaningTableData = computed(() => {
  const details = dashboard.value?.cleaning_details || []
  return details.map((detail, index) => ({
    key: String(index + 1),
    school: detail.school_name,
    portions: detail.portions, // Backend already returns portions field
    status: detail.status
  }))
})

// Chart composables
const { setOption: setProductionOption } = useHorizonChart(productionChartRef, {})
const { setOption: setDeliveryOption } = useHorizonChart(deliveryChartRef, {})
const { setOption: setCleaningOption } = useHorizonChart(cleaningChartRef, {})

const updateCharts = () => {
  const p = production.value
  const d = delivery.value

  // Production donut chart
  const pieData = [
    { value: p.recipes_pending, name: 'Menunggu', itemStyle: { color: '#A3AED0' } },
    { value: p.recipes_cooking, name: 'Sedang Dimasak', itemStyle: { color: '#FFB547' } },
    { value: p.recipes_ready, name: 'Selesai Dimasak', itemStyle: { color: '#764AF1' } },
    { value: p.packing_pending, name: 'Siap Packing', itemStyle: { color: '#6B6B6B' } },
    { value: p.packing_in_progress, name: 'Sedang Packing', itemStyle: { color: '#F82C17' } },
    { value: p.packing_ready, name: 'Selesai Packing', itemStyle: { color: '#5A38C4' } },
  ].filter(item => item.value > 0)

  // If no data, show a "no data" placeholder
  const finalPieData = pieData.length > 0 ? pieData : [{ value: 1, name: 'Belum ada data', itemStyle: { color: '#F0F0F0' } }]

  setProductionOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, data: finalPieData.map(i => i.name) },
    series: [{
      name: 'Produksi',
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{c}' },
      data: finalPieData
    }]
  })

  // Delivery bar chart - use status_breakdown from API
  const deliveryBreakdown = d.status_breakdown || []
  const deliveryLabels = deliveryBreakdown.map(item => item.status_label)
  const deliveryData = deliveryBreakdown.map((item, index) => {
    // Assign colors based on status — mix red and purple
    const colors = ['#A3AED0', '#FFB547', '#764AF1', '#F82C17', '#5A38C4', '#303030']
    return {
      value: item.count,
      itemStyle: { color: colors[index % colors.length], borderRadius: [0, 4, 4, 0] }
    }
  })

  // If no data, show placeholder
  const finalDeliveryLabels = deliveryLabels.length > 0 ? deliveryLabels : ['Belum ada data']
  const finalDeliveryData = deliveryData.length > 0 ? deliveryData : [{ value: 0, itemStyle: { color: '#F0F0F0', borderRadius: [0, 4, 4, 0] } }]

  setDeliveryOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '8%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'value',
      axisLabel: { formatter: '{value}' }
    },
    yAxis: {
      type: 'category',
      data: finalDeliveryLabels
    },
    series: [{
      name: 'Pengiriman',
      type: 'bar',
      data: finalDeliveryData,
      barWidth: '40%',
      label: { show: true, position: 'right', formatter: '{c}' }
    }]
  })

  // Cleaning pie chart (same style as production)
  const c = cleaning.value
  const cleaningPieData = [
    { value: c.items_pending, name: 'Menunggu', itemStyle: { color: '#A3AED0' } },
    { value: c.items_in_progress, name: 'Sedang Dicuci', itemStyle: { color: '#FFB547' } },
    { value: c.items_completed, name: 'Selesai', itemStyle: { color: '#764AF1' } },
  ].filter(item => item.value > 0)

  // If no data, show a "no data" placeholder
  const finalCleaningData = cleaningPieData.length > 0 ? cleaningPieData : [{ value: 1, name: 'Belum ada data', itemStyle: { color: '#F0F0F0' } }]

  setCleaningOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, data: finalCleaningData.map(i => i.name) },
    series: [{
      name: 'Pencucian',
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '45%'],
      avoidLabelOverlap: true,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}\n{c}' },
      data: finalCleaningData
    }]
  })
}

const loadDashboardData = async () => {
  loading.value = true
  try {
    // Load dashboard data
    const response = await getKepalaSSPGDashboard()
    if (response.success) {
      dashboard.value = response.dashboard
    }

    // Load review summary
    try {
      const reviewResponse = await reviewService.getSummary()
      if (reviewResponse.summary) {
        reviewSummary.value = reviewResponse.summary
      }
    } catch (reviewError) {
      console.warn('Could not load review summary:', reviewError)
    }

    // Load top 5 suppliers
    try {
      const supplierResponse = await supplierService.getSupplierStats()
      if (supplierResponse.data && supplierResponse.data.data) {
        topSuppliers.value = (supplierResponse.data.data.topSuppliers || []).slice(0, 5)
      }
    } catch (supplierError) {
      console.warn('Could not load top suppliers:', supplierError)
    }
  } catch (error) {
    console.error('Error loading dashboard:', error)
    message.error('Gagal memuat data dashboard')
  } finally {
    loading.value = false
    await nextTick()
    updateCharts()
  }
}

const setupFirebaseListeners = () => {
  try {
    const dashRef = dbRef(database, firebasePaths.dashboardKepalaSPPG())
    dashboardListener = onValue(dashRef, (snapshot) => {
      const data = snapshot.val()
      if (data) {
        dashboard.value = data
        updateCharts()
      }
    }, (error) => {
      console.error('Firebase listener error:', error)
    })
  } catch (e) {
    console.warn('Firebase not available:', e)
  }
}

const cleanupFirebaseListeners = () => {
  if (dashboardListener) {
    try { off(dbRef(database, firebasePaths.dashboardKepalaSPPG())) } catch (e) { /* ignore */ }
    dashboardListener = null
  }
}

const getStatusType = (status) => {
  const s = String(status).toLowerCase()
  if (s.includes('selesai') || s.includes('completed')) return 'success'
  if (s.includes('packing') || s.includes('memasak') || s.includes('%')) return 'warning'
  return 'default'
}

// Map data for schools and suppliers
const schoolsList = ref([])
const suppliersList = ref([])

// School delivery/review stats cache
const schoolStats = ref({})

const sppgMapMarkers = computed(() => {
  const markers = []
  schoolsList.value.forEach(s => {
    if (s.latitude !== 0 || s.longitude !== 0) {
      const stats = schoolStats.value[s.id]
      markers.push({
        id: `school-${s.id}`,
        name: s.name,
        kode: s.npsn || '',
        type: 'Sekolah',
        latitude: s.latitude,
        longitude: s.longitude,
        details: `${s.category || ''} • ${s.student_count || 0} siswa`,
        stats: stats ? {
          portionsSmall: stats.portions_small,
          portionsLarge: stats.portions_large,
          totalDelivered: stats.total_delivered,
          rating: stats.avg_rating
        } : {
          portionsSmall: s.student_count_grade_1_3 || 0,
          portionsLarge: (s.student_count_grade_4_6 || 0) + (s.category !== 'SD' ? (s.student_count || 0) : 0),
          totalDelivered: 0,
          rating: 0
        }
      })
    }
  })
  suppliersList.value.forEach(s => {
    if (s.latitude !== 0 || s.longitude !== 0) {
      markers.push({
        id: `supplier-${s.id}`,
        name: s.name,
        kode: s.product_category || '',
        type: 'Supplier',
        latitude: s.latitude,
        longitude: s.longitude,
        details: s.contact_person ? `CP: ${s.contact_person}` : '',
        stats: null
      })
    }
  })
  return markers
})

const loadMapData = async () => {
  try {
    const [schoolsRes, suppliersRes] = await Promise.all([
      api.get('/schools'),
      api.get('/suppliers?active_only=true')
    ])
    schoolsList.value = schoolsRes.data?.schools || schoolsRes.data?.data || []
    suppliersList.value = suppliersRes.data?.suppliers || suppliersRes.data?.data || []

    // Fetch per-school delivery + review stats (all-time, completed only)
    try {
      const statsRes = await api.get('/schools/map-stats')
      const statsArr = statsRes.data?.data || []
      const result = {}
      statsArr.forEach(s => {
        result[s.school_id] = {
          portions_small: s.portions_small || 0,
          portions_large: s.portions_large || 0,
          total_delivered: s.total_delivered || 0,
          avg_rating: s.avg_rating || 0
        }
      })
      schoolStats.value = result
    } catch (e) { /* stats optional */ }
  } catch (e) {
    // silent — map is optional
  }
}

onMounted(() => {
  setDefaultCashFlowDateRange()
  loadDashboardData()
  loadCashFlowSummary()
  setupFirebaseListeners()
  loadMapData()
})

onUnmounted(() => {
  cleanupFirebaseListeners()
})
</script>

<style scoped>
.dashboard-sspg { min-height: 100%; }
.dashboard-loading { padding: 40px; }

/* === BENTO GRID === */
.bento {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-auto-rows: auto;
  gap: 16px;
}

.bento > * {
  background: #fff;
  border-radius: 14px;
  padding: 20px;
  border: 1px solid #F0F0F0;
  min-width: 0;
  overflow: hidden;
}

/* Grid placement */
.bento__hero    { grid-column: 1 / -1; border: none; }
.bento__map     { grid-column: 1 / 3; }
.bento__rating  { grid-column: 3 / 4; }
.bento__chart1  { grid-column: 1 / 2; }
.bento__chart2  { grid-column: 2 / 3; }
.bento__chart3  { grid-column: 3 / 4; }
.bento__cashflow { grid-column: 1 / 3; }
.bento__supplier { grid-column: 3 / 4; }
.bento__stock   { grid-column: 1 / -1; }
.bento__stock-ok { grid-column: 1 / -1; text-align: center; color: #764AF1; font-weight: 600; font-size: 14px; }

@media (max-width: 1024px) {
  .bento { grid-template-columns: 1fr 1fr; }
  .bento__hero { grid-column: 1 / -1; }
  .bento__map { grid-column: 1 / -1; }
  .bento__rating { grid-column: 1 / -1; }
  .bento__chart3 { grid-column: 1 / -1; }
  .bento__cashflow { grid-column: 1 / -1; }
  .bento__supplier { grid-column: 1 / -1; }
  .bento__stock { grid-column: 1 / -1; }
}

@media (max-width: 768px) {
  .bento { grid-template-columns: 1fr; }
  .bento > * { grid-column: 1 / -1 !important; }
}

/* === HERO === */
.bento__hero {
  background: #F82C17;
  color: #fff;
  padding: 24px 28px;
}
.hero-content { display: flex; justify-content: space-between; align-items: center; gap: 20px; }
.hero-greeting { font-size: 26px; font-weight: 700; margin: 0 0 4px; letter-spacing: -0.5px; color: #FFFFFF; }
.hero-subtitle { font-size: 14px; opacity: 0.9; margin: 0; color: rgba(255,255,255,0.85); }
.hero-illustration { height: 90px; width: auto; opacity: 0.9; flex-shrink: 0; }
.hero-metrics { display: flex; gap: 12px; }
.hero-metric { background: rgba(255,255,255,0.18); backdrop-filter: blur(4px); border-radius: 10px; padding: 12px 18px; min-width: 100px; text-align: center; }
.hero-metric__label { display: block; font-size: 11px; opacity: 0.75; margin-bottom: 2px; }
.hero-metric__value { display: block; font-size: 26px; font-weight: 700; letter-spacing: -0.5px; }
@media (max-width: 768px) {
  .hero-content { flex-direction: column; align-items: flex-start; }
  .hero-metrics { flex-wrap: wrap; }
  .hero-metric { min-width: 80px; padding: 10px 14px; }
  .hero-metric__value { font-size: 20px; }
  .hero-illustration { display: none; }
}

/* === TITLES === */
.bento__title { font-size: 16px; font-weight: 700; margin: 0 0 12px; color: #764AF1; }
.bento__title-row { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.bento__title-row .bento__title { margin-bottom: 0; }

/* === MAP HEADER === */
.map-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px; }
.map-header .bento__title { margin-bottom: 0; }
.map-legend { display: flex; gap: 16px; }
.map-legend__item { display: flex; align-items: center; gap: 6px; cursor: pointer; user-select: none; font-size: 13px; color: #303030; }
.map-legend__item:hover { opacity: 0.8; }
.map-legend__dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; transition: background 0.2s; }
.map-legend__label { font-weight: 500; }

/* === COMPACT STATS (Rating) === */
.compact-stat { display: flex; justify-content: space-between; align-items: baseline; padding: 10px 0; border-bottom: 1px solid #F0F0F0; }
.compact-stat:last-of-type { border-bottom: none; }
.compact-stat__label { font-size: 14px; color: #6B6B6B; }
.compact-stat__value { font-size: 24px; font-weight: 700; color: #303030; }
.compact-stat__value small { font-size: 14px; font-weight: 400; color: #6B6B6B; }
.compact-stat__footer { font-size: 12px; color: #6B6B6B; margin-top: 8px; text-align: center; }

/* === CASHFLOW === */
.cashflow-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.cashflow-item { padding: 14px; border-radius: 10px; }
.cashflow-item--income { background: #D1FAE5; }
.cashflow-item--expense { background: #FDEAE7; }
.cashflow-item--net { background: #E8DAFE; }
.cashflow-item__label { display: block; font-size: 12px; color: #6B6B6B; margin-bottom: 4px; }
.cashflow-item__value { display: block; font-size: 18px; font-weight: 700; color: #303030; }
.text-error { color: #F82C17 !important; }

/* === SUPPLIER COMPACT === */
.supplier-compact-list { display: flex; flex-direction: column; gap: 8px; }
.supplier-compact { display: flex; align-items: center; gap: 10px; padding: 8px 0; border-bottom: 1px solid #F0F0F0; }
.supplier-compact:last-child { border-bottom: none; }
.supplier-compact__rank { width: 28px; height: 28px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 700; flex-shrink: 0; }
.supplier-compact__rank.rank-1 { background: #FDEAE7; color: #F82C17; }
.supplier-compact__rank.rank-2 { background: #FEF3C7; color: #D97706; }
.supplier-compact__rank.rank-3 { background: #D1FAE5; color: #764AF1; }
.supplier-compact__name { flex: 1; font-size: 13px; font-weight: 600; color: #303030; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.supplier-compact__amount { font-size: 12px; color: #6B6B6B; flex-shrink: 0; }

/* === STOCK COMPACT === */
.stock-compact-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; }
@media (max-width: 1024px) { .stock-compact-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 768px) { .stock-compact-grid { grid-template-columns: 1fr; } }
.stock-compact { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 10px 14px; background: #FFF; border: 1px solid #F0F0F0; border-radius: 10px; border-left: 3px solid #764AF1; }
.stock-compact__info { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.stock-compact__name { font-size: 13px; font-weight: 600; color: #303030; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.stock-compact__qty { font-size: 12px; color: #F82C17; font-weight: 600; }

/* === SKELETON === */
.stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
</style>
