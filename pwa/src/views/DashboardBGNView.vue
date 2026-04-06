<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="dashboard-page">
      <!-- NavBar -->
      <van-nav-bar
        :title="navTitle"
        :left-arrow="!!route.params.yayasan_id"
        @click-left="goBack"
      />

      <!-- Breadcrumb indicator -->
      <div class="breadcrumb">
        <van-tag type="primary" size="medium">
          <van-icon name="flag-o" /> BGN
        </van-tag>
        <van-tag v-if="selectedYayasan" type="primary" size="medium" style="margin-left: 6px;">
          <van-icon name="shop-o" /> {{ selectedYayasan.yayasan_nama || selectedYayasan.yayasan_kode }}
        </van-tag>
        <van-tag v-if="selectedSppg" type="success" size="medium" style="margin-left: 6px;">
          <van-icon name="location-o" /> {{ selectedSppg.sppg_nama || selectedSppg.sppg_kode }}
        </van-tag>
      </div>

      <!-- Loading State -->
      <template v-if="loading">
        <div class="metrics-grid">
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
          <SkeletonCard :rows="2" />
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ error }}</p>
        <van-button type="primary" size="normal" @click="fetchData">
          Coba Lagi
        </van-button>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- Aggregated Metrics -->
        <div class="metrics-grid">
          <MetricCard
            icon="apps-o"
            iconColor="#F0F0F0"
            label="Total Porsi"
            :value="dashboard.aggregated_production?.total_portions || 0"
          />
          <MetricCard
            icon="logistics"
            iconColor="#4CAF50"
            label="Delivery Rate"
            :value="`${(dashboard.aggregated_delivery?.completion_rate || 0).toFixed(1)}%`"
          />
          <MetricCard
            icon="balance-o"
            iconColor="#FF9800"
            label="Penyerapan Anggaran"
            :value="`${(dashboard.aggregated_financial?.absorption_rate || 0).toFixed(1)}%`"
          />
          <MetricCard
            icon="star-o"
            iconColor="#FFA726"
            label="Rata-rata Rating"
            :value="`${(dashboard.aggregated_review?.average_overall || 0).toFixed(1)}/5`"
          />
        </div>

        <!-- Overview counts -->
        <div class="detail-card h-card">
          <h3 class="section-title">Ringkasan Nasional</h3>
          <van-cell-group inset>
            <van-cell title="Total Yayasan" :value="dashboard.total_yayasan || 0" />
            <van-cell title="Total SPPG" :value="dashboard.total_sppg || 0" />
            <van-cell title="Total Pengiriman" :value="dashboard.aggregated_delivery?.total_deliveries || 0" />
            <van-cell title="Pengiriman Selesai" :value="dashboard.aggregated_delivery?.completed_deliveries || 0" />
            <van-cell title="On-Time Rate" :value="`${(dashboard.aggregated_delivery?.on_time_rate || 0).toFixed(1)}%`" />
          </van-cell-group>
        </div>

        <!-- Financial Summary -->
        <div class="detail-card h-card">
          <h3 class="section-title">Ringkasan Keuangan</h3>
          <van-cell-group inset>
            <van-cell title="Total Anggaran" :value="formatCurrency(dashboard.aggregated_financial?.total_budget)" />
            <van-cell title="Total Pengeluaran" :value="formatCurrency(dashboard.aggregated_financial?.total_spent)" />
          </van-cell-group>
        </div>

        <!-- Yayasan List (top-level drill-down) -->
        <div v-if="!route.params.yayasan_id" class="detail-card h-card">
          <h3 class="section-title">Daftar Yayasan ({{ dashboard.total_yayasan || 0 }})</h3>
          <div v-if="!yayasanList.length" class="empty-state">
            <MobileEmptyState description="Belum ada data Yayasan" />
          </div>
          <van-cell-group v-else inset>
            <van-cell
              v-for="yayasan in yayasanList"
              :key="yayasan.yayasan_id"
              :title="yayasan.yayasan_nama"
              :label="`${yayasan.yayasan_kode} • ${yayasan.total_sppg} SPPG • ${yayasan.total_portions} porsi`"
              is-link
              @click="drillDownYayasan(yayasan)"
            >
              <template #value>
                <div class="item-metrics">
                  <van-tag type="warning" size="small">{{ (yayasan.average_review_rating || 0).toFixed(1) }}★</van-tag>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </div>

        <!-- SPPG List (when drilled into a Yayasan but not a specific SPPG) -->
        <div v-if="route.params.yayasan_id && !route.params.sppg_id && sppgList.length" class="detail-card h-card">
          <h3 class="section-title">Daftar SPPG</h3>
          <van-cell-group inset>
            <van-cell
              v-for="sppg in sppgList"
              :key="sppg.sppg_id"
              :title="sppg.sppg_nama"
              :label="`${sppg.sppg_kode} • ${sppg.total_portions} porsi`"
              is-link
              @click="drillDownSppg(sppg)"
            >
              <template #value>
                <div class="item-metrics">
                  <van-tag type="primary" size="small">{{ (sppg.delivery_rate || 0).toFixed(0) }}%</van-tag>
                  <van-tag type="warning" size="small">{{ (sppg.average_review_rating || 0).toFixed(1) }}★</van-tag>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </div>

        <!-- Review Summary -->
        <div class="detail-card h-card">
          <h3 class="section-title">Monitoring Ulasan</h3>
          <van-cell-group inset>
            <van-cell title="Total Ulasan" :value="dashboard.aggregated_review?.total_reviews || 0" />
            <van-cell title="Rating Menu" :value="`${(dashboard.aggregated_review?.average_menu_rating || 0).toFixed(1)}/5`" />
            <van-cell title="Rating Layanan" :value="`${(dashboard.aggregated_review?.average_service_rating || 0).toFixed(1)}/5`" />
          </van-cell-group>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import dashboardAggregatedService from '@/services/dashboardAggregatedService'
import { cacheDashboardData, getCachedDashboardData } from '@/services/db'
import MetricCard from '@/components/mobile/MetricCard.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'
import MobileEmptyState from '@/components/common/MobileEmptyState.vue'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const error = ref(null)
const refreshing = ref(false)
const dashboard = ref({})

const navTitle = computed(() => {
  if (route.params.sppg_id && selectedSppg.value) {
    return selectedSppg.value.sppg_nama || 'Detail SPPG'
  }
  if (route.params.yayasan_id && selectedYayasan.value) {
    return selectedYayasan.value.yayasan_nama || 'Detail Yayasan'
  }
  return 'Dashboard BGN'
})

const yayasanList = computed(() => dashboard.value.yayasan_summaries || [])

// When drilled into a yayasan, the API returns sppg_summaries
const sppgList = computed(() => dashboard.value.sppg_summaries || [])

const selectedYayasan = computed(() => {
  if (!route.params.yayasan_id) return null
  return yayasanList.value.find(y => String(y.yayasan_id) === String(route.params.yayasan_id))
    || { yayasan_kode: `Yayasan #${route.params.yayasan_id}` }
})

const selectedSppg = computed(() => {
  if (!route.params.sppg_id) return null
  return sppgList.value.find(s => String(s.sppg_id) === String(route.params.sppg_id))
    || { sppg_kode: `SPPG #${route.params.sppg_id}` }
})

function formatCurrency(value) {
  if (!value) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

function drillDownYayasan(yayasan) {
  router.push(`/dashboard-bgn/${yayasan.yayasan_id}`)
}

function drillDownSppg(sppg) {
  router.push(`/dashboard-bgn/${route.params.yayasan_id}/${sppg.sppg_id}`)
}

function goBack() {
  if (route.params.sppg_id) {
    router.push(`/dashboard-bgn/${route.params.yayasan_id}`)
  } else {
    router.push('/dashboard-bgn')
  }
}

async function fetchData() {
  loading.value = true
  error.value = null
  const params = {}
  if (route.params.yayasan_id) params.yayasan_id = route.params.yayasan_id
  if (route.params.sppg_id) params.sppg_id = route.params.sppg_id

  const cacheKey = `bgn_dashboard_${route.params.yayasan_id || 'all'}_${route.params.sppg_id || 'all'}`

  try {
    const response = await dashboardAggregatedService.getAdminBGNDashboard(params)
    if (response.data?.success) {
      dashboard.value = response.data.data || {}
      // Cache for offline use
      await cacheDashboardData(cacheKey, dashboard.value)
    } else {
      throw new Error(response.data?.message || 'Gagal memuat data')
    }
  } catch (err) {
    console.error('Error fetching BGN dashboard:', err)
    // Try loading from cache
    const cached = await getCachedDashboardData(cacheKey)
    if (cached) {
      dashboard.value = cached
      error.value = null
    } else {
      error.value = err.message || 'Gagal memuat data dashboard'
    }
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  await fetchData()
  refreshing.value = false
}

watch(
  () => [route.params.yayasan_id, route.params.sppg_id],
  () => { fetchData() }
)

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard-page {
  padding: 0;
  padding-bottom: 80px;
  min-height: 100vh;
  background: var(--h-bg-page);
}

.dashboard-page > :not(.van-nav-bar) {
  padding-left: var(--h-spacing-lg);
  padding-right: var(--h-spacing-lg);
}

.breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 4px;
  padding: var(--h-spacing-sm) var(--h-spacing-lg);
  background: var(--h-bg-light);
}

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--h-spacing-md);
  margin-bottom: var(--h-spacing-lg);
  margin-top: var(--h-spacing-md);
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 var(--h-spacing-md) 0;
}

.detail-card {
  margin-bottom: var(--h-spacing-lg);
}

.item-metrics {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px var(--h-spacing-xl);
  text-align: center;
}

.error-state__message {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: var(--h-spacing-lg) 0;
  line-height: 1.5;
}

.empty-state {
  padding: var(--h-spacing-xl) 0;
  text-align: center;
}

.empty-state__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: 0;
}
</style>
