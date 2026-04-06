<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="dashboard-page">
      <!-- NavBar -->
      <van-nav-bar title="Dashboard" />

      <!-- Loading State -->
      <template v-if="dashboardStore.loading">
        <div class="dashboard-loading">
          <LottiePlayer src="/lottie/loading-cooking.json" width="100px" height="100px" />
          <p class="dashboard-loading__text">Memuat data...</p>
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="dashboardStore.error" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ dashboardStore.error }}</p>
        <van-button type="primary" size="normal" @click="dashboardStore.retry()">
          Coba Lagi
        </van-button>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- Metrics Grid (7 cards in 2 rows) -->
        <div class="metrics-grid">
          <!-- Row 1: 4 cards -->
          <MetricCard
            icon="apps-o"
            iconColor="#FDEAE7"
            label="Porsi Disiapkan"
            :value="dashboardStore.summary.porsiDisiapkan"
            :trend="dashboardStore.summary.porsiDisiapkanTrend"
            trendUp
          />
          <MetricCard
            icon="logistics"
            iconColor="#D1FAE5"
            label="Delivery Rate"
            :value="`${dashboardStore.summary.deliveryRate}%`"
            :trend="dashboardStore.summary.deliveryRateTrend"
            trendUp
          />
          <MetricCard
            icon="bag-o"
            iconColor="#FEF3C7"
            label="Ketersediaan Stok"
            :value="`${dashboardStore.summary.ketersediaanStok}%`"
            :trend="dashboardStore.summary.stokKritisTrend"
            trendDown
          />
          <MetricCard
            icon="success"
            iconColor="#D1FAE5"
            label="On-Time Delivery"
            :value="`${dashboardStore.summary.onTimeDelivery}%`"
            :trend="dashboardStore.summary.onTimeDeliveryTrend"
            trendUp
          />
          
          <!-- Row 2: 3 cards -->
          <MetricCard
            icon="star-o"
            iconColor="#FDEAE7"
            label="Rating Keseluruhan"
            :value="`${dashboardStore.summary.ratingKeseluruhan}/5`"
            :trend="dashboardStore.summary.ratingKeseluruhanTrend"
            trendUp
          />
          <MetricCard
            icon="fire-o"
            iconColor="#D1FAE5"
            label="Rating Menu"
            :value="`${dashboardStore.summary.ratingMenu}/5`"
            :trend="dashboardStore.summary.ratingMenuTrend"
            trendUp
          />
          <MetricCard
            icon="logistics"
            iconColor="#D1FAE5"
            label="Rating Layanan"
            :value="`${dashboardStore.summary.ratingLayanan}/5`"
            :trend="dashboardStore.summary.ratingLayananTrend"
            trendUp
          />
        </div>

        <!-- Arus Kas Section -->
        <div class="detail-card h-card arus-kas-card">
          <h3 class="section-title">Arus Kas</h3>
          <div class="arus-kas-grid">
            <div class="arus-kas-item income">
              <div class="arus-kas-icon">
                <van-icon name="arrow-up" size="20" />
              </div>
              <div class="arus-kas-info">
                <span class="arus-kas-label">Total Pemasukan</span>
                <span class="arus-kas-value">{{ formatCurrency(dashboardStore.arusKas.totalPemasukan) }}</span>
                <span class="arus-kas-period">{{ dashboardStore.arusKas.period }}</span>
              </div>
            </div>
            <div class="arus-kas-item expense">
              <div class="arus-kas-icon">
                <van-icon name="arrow-down" size="20" />
              </div>
              <div class="arus-kas-info">
                <span class="arus-kas-label">Total Pengeluaran</span>
                <span class="arus-kas-value">{{ formatCurrency(dashboardStore.arusKas.totalPengeluaran) }}</span>
                <span class="arus-kas-period">{{ dashboardStore.arusKas.period }}</span>
              </div>
            </div>
            <div class="arus-kas-item net">
              <div class="arus-kas-icon">
                <van-icon name="balance-o" size="20" />
              </div>
              <div class="arus-kas-info">
                <span class="arus-kas-label">Arus Kas Bersih</span>
                <span class="arus-kas-value" :class="{ negative: dashboardStore.arusKas.netCashFlow < 0 }">
                  {{ formatCurrency(dashboardStore.arusKas.netCashFlow) }}
                </span>
                <span class="arus-kas-period">{{ dashboardStore.arusKas.period }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Top 5 Supplier Section -->
        <div class="detail-card h-card supplier-card">
          <h3 class="section-title">Top 5 Supplier</h3>
          <div v-if="dashboardStore.topSuppliers.length === 0" class="empty-state">
            <MobileEmptyState lottie="/lottie/empty-box.json" description="Belum ada data supplier" />
          </div>
          <div v-else class="supplier-list">
            <div
              v-for="(supplier, index) in dashboardStore.topSuppliers"
              :key="supplier.id"
              class="supplier-item"
            >
              <div class="supplier-rank" :class="`rank-${index + 1}`">
                {{ index + 1 }}
              </div>
              <div class="supplier-info">
                <span class="supplier-name">{{ supplier.name }}</span>
                <span class="supplier-meta">{{ supplier.total_orders }} Order • {{ formatCurrency(supplier.total_amount) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Detail Tables -->
        <div class="detail-section">
          <!-- Detail Produksi -->
          <div class="detail-card h-card">
            <h3 class="section-title">Detail Produksi</h3>
            <div v-if="dashboardStore.detailProduksi.length === 0" class="empty-state">
              <MobileEmptyState description="Belum ada data" />
            </div>
            <div v-else class="detail-table">
              <div class="detail-table-header">
                <span class="col-school">SEKOLAH</span>
                <span class="col-portion">PORSI</span>
                <span class="col-status">STATUS</span>
              </div>
              <div
                v-for="item in dashboardStore.detailProduksi"
                :key="item.id"
                class="detail-table-row"
              >
                <span class="col-school">{{ item.sekolah || item.school }}</span>
                <span class="col-portion">{{ item.porsi || item.portions }}</span>
                <span class="col-status">
                  <van-tag :type="getStatusType(item.status)" size="small">
                    {{ item.status }}
                  </van-tag>
                </span>
              </div>
            </div>
          </div>

          <!-- Detail Pengiriman & Pengambilan -->
          <div class="detail-card h-card">
            <h3 class="section-title">Detail Pengiriman & Pengambilan</h3>
            <div v-if="dashboardStore.detailPengiriman.length === 0" class="empty-state">
              <MobileEmptyState description="Belum ada data" />
            </div>
            <div v-else class="detail-table">
              <div class="detail-table-header">
                <span class="col-school">SEKOLAH</span>
                <span class="col-portion">PORSI</span>
                <span class="col-status">STATUS</span>
              </div>
              <div
                v-for="item in dashboardStore.detailPengiriman"
                :key="item.id"
                class="detail-table-row"
              >
                <span class="col-school">{{ item.sekolah || item.school }}</span>
                <span class="col-portion">{{ item.porsi || item.portions }}</span>
                <span class="col-status">
                  <van-tag :type="getStatusType(item.status)" size="small">
                    {{ item.status }}
                  </van-tag>
                </span>
              </div>
            </div>
          </div>

          <!-- Detail Pencucian -->
          <div class="detail-card h-card">
            <h3 class="section-title">Detail Pencucian</h3>
            <div v-if="dashboardStore.detailPencucian.length === 0" class="empty-state">
              <MobileEmptyState description="Belum ada data" />
            </div>
            <div v-else class="detail-table">
              <div class="detail-table-header">
                <span class="col-school">SEKOLAH</span>
                <span class="col-portion">PORSI</span>
                <span class="col-status">STATUS</span>
              </div>
              <div
                v-for="item in dashboardStore.detailPencucian"
                :key="item.id"
                class="detail-table-row"
              >
                <span class="col-school">{{ item.sekolah || item.school }}</span>
                <span class="col-portion">{{ item.porsi || item.portions }}</span>
                <span class="col-status">
                  <van-tag :type="getStatusType(item.status)" size="small">
                    {{ item.status }}
                  </van-tag>
                </span>
              </div>
            </div>
          </div>

          <!-- Stok Kritis -->
          <div class="detail-card h-card stok-kritis-card">
            <h3 class="section-title">Stok Kritis ({{ dashboardStore.stokKritis.length }} Item)</h3>
            <div v-if="dashboardStore.stokKritis.length === 0" class="empty-state">
              <MobileEmptyState description="Tidak ada stok kritis" />
            </div>
            <div v-else>
              <div class="stok-kritis-grid">
                <div
                  v-for="item in dashboardStore.stokKritis.slice(0, 6)"
                  :key="item.id"
                  class="stok-kritis-item"
                >
                  <div class="stok-kritis-info">
                    <span class="stok-kritis-name">{{ item.nama || item.name }}</span>
                    <span class="stok-kritis-stock">{{ item.stok || item.stock || 0 }} {{ item.satuan || item.unit }}</span>
                    <span class="stok-kritis-min">Min: {{ item.min || item.minimum || 0 }} {{ item.satuan || item.unit }}</span>
                  </div>
                  <van-tag type="danger" size="small">KRITIS</van-tag>
                </div>
              </div>
              <p v-if="dashboardStore.stokKritis.length > 6" class="stok-kritis-note">
                Menampilkan 6 item dari {{ dashboardStore.stokKritis.length }} item kritis
              </p>
            </div>
          </div>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import MetricCard from '@/components/mobile/MetricCard.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'
import MobileEmptyState from '@/components/common/MobileEmptyState.vue'
import LottiePlayer from '@/components/common/LottiePlayer.vue'

const dashboardStore = useDashboardStore()
const refreshing = ref(false)

function getStatusType(status) {
  const statusLower = status?.toLowerCase() || ''
  if (statusLower.includes('selesai') || statusLower.includes('completed')) return 'success'
  if (statusLower.includes('proses') || statusLower.includes('progress')) return 'primary'
  if (statusLower.includes('pending') || statusLower.includes('menunggu')) return 'warning'
  return 'default'
}

function formatCurrency(value) {
  if (!value) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

async function onRefresh() {
  await dashboardStore.fetchDashboardData()
  refreshing.value = false
}

onMounted(() => {
  dashboardStore.fetchDashboardData()
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

/* Dashboard Loading */
.dashboard-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  gap: 8px;
  padding: 40px 0;
}

.dashboard-loading__text {
  font-size: 13px;
  color: var(--h-text-secondary);
}

/* Metrics Grid (7 cards in 2 rows) */
.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--h-spacing-md);
  margin-bottom: var(--h-spacing-lg);
  margin-top: var(--h-spacing-md);
}

/* Section title */
.section-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 var(--h-spacing-md) 0;
}

.section-block {
  margin-bottom: var(--h-spacing-lg);
}

/* Detail Section */
.detail-section {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-lg);
}

.detail-card {
  margin-bottom: 0;
}

/* Detail Table */
.detail-table {
  display: flex;
  flex-direction: column;
}

.detail-table-header {
  display: grid;
  grid-template-columns: 2fr 1fr 1.2fr;
  gap: var(--h-spacing-sm);
  padding: var(--h-spacing-sm) 0;
  border-bottom: 2px solid #D8D8DB;
  font-size: 11px;
  font-weight: 600;
  color: var(--h-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.2fr;
  gap: var(--h-spacing-sm);
  padding: var(--h-spacing-md) 0;
  border-bottom: 1px solid var(--h-border-light);
  align-items: center;
}

.detail-table-row:last-child {
  border-bottom: none;
}

.col-school {
  font-size: 13px;
  color: var(--h-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.col-portion {
  font-size: 13px;
  color: var(--h-text-primary);
  text-align: center;
}

.col-status {
  display: flex;
  justify-content: flex-end;
}

/* Stok Kritis Card */
.stok-kritis-card {
  margin-bottom: var(--h-spacing-lg);
}

.stok-kritis-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--h-spacing-md);
}

.stok-kritis-header .section-title {
  margin: 0;
}

.stok-kritis-grid {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-md);
}

.stok-kritis-item {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: var(--h-spacing-md);
  background: #FFFFFF;
  border: 1px solid var(--h-border-color);
  border-radius: var(--h-radius-md);
  border-left: 3px solid var(--h-error);
}

.stok-kritis-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.stok-kritis-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--h-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stok-kritis-stock {
  font-size: 14px;
  font-weight: 600;
  color: var(--h-error);
}

.stok-kritis-min {
  font-size: 11px;
  color: var(--h-text-secondary);
}

.stok-kritis-note {
  font-size: 12px;
  color: var(--h-text-secondary);
  text-align: center;
  margin-top: var(--h-spacing-md);
  padding-top: var(--h-spacing-md);
  border-top: 1px solid var(--h-border-light);
}

/* Error State */
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

/* Empty State */
.empty-state {
  padding: var(--h-spacing-xl) 0;
  text-align: center;
}

.empty-state__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: 0;
}

/* Arus Kas Card */
.arus-kas-card {
  margin-bottom: var(--h-spacing-lg);
}

.arus-kas-grid {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-md);
}

.arus-kas-item {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-md);
  padding: var(--h-spacing-md);
  background: #FFFFFF;
  border: 1px solid var(--h-border-color);
  border-radius: var(--h-radius-md);
  border-left: 4px solid;
}

.arus-kas-item.income {
  border-left-color: #52c41a;
}

.arus-kas-item.expense {
  border-left-color: #ff4d4f;
}

.arus-kas-item.net {
  border-left-color: #1890ff;
}

.arus-kas-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.arus-kas-item.income .arus-kas-icon {
  background: rgba(82, 196, 26, 0.1);
  color: #52c41a;
}

.arus-kas-item.expense .arus-kas-icon {
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
}

.arus-kas-item.net .arus-kas-icon {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

.arus-kas-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.arus-kas-label {
  font-size: 12px;
  color: var(--h-text-secondary);
}

.arus-kas-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
}

.arus-kas-value.negative {
  color: var(--h-error);
}

.arus-kas-period {
  font-size: 11px;
  color: var(--h-text-light);
}

/* Supplier Card */
.supplier-card {
  margin-bottom: var(--h-spacing-lg);
}

.supplier-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--h-spacing-md);
}

.supplier-header .section-title {
  margin: 0;
}

.supplier-list {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-md);
}

.supplier-item {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-md);
  padding: var(--h-spacing-md);
  background: #FFFFFF;
  border: 1px solid var(--h-border-color);
  border-radius: var(--h-radius-md);
}

.supplier-rank {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  color: #303030;
  flex-shrink: 0;
}

.supplier-rank.rank-1 {
  background: #FDEAE7;
  color: #C94A3A;
}

.supplier-rank.rank-2 {
  background: #FEF3C7;
  color: #D97706;
}

.supplier-rank.rank-3 {
  background: #D1FAE5;
  color: #1E8A6E;
}

.supplier-rank.rank-4 {
  background: #F0F0F0;
}

.supplier-rank.rank-5 {
  background: #F0F0F0;
}

.supplier-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.supplier-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--h-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.supplier-meta {
  font-size: 12px;
  color: var(--h-text-secondary);
}
</style>
