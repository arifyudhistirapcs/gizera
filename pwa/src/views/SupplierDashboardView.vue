<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="supplier-dashboard-page">
      <!-- Header -->
      <div class="dash-header">
        <div class="dash-header__content">
          <h1 class="dash-header__title">Dashboard Supplier</h1>
          <p class="dash-header__subtitle">Ringkasan pesanan & pembayaran</p>
        </div>
      </div>

      <!-- Loading -->
      <van-loading v-if="loading" type="spinner" vertical class="page-loading">
        Memuat data...
      </van-loading>

      <!-- Error -->
      <div v-else-if="error" class="error-state">
        <van-icon name="warning-o" size="48" color="#ee0a24" />
        <p class="error-state__message">{{ error }}</p>
        <van-button type="primary" size="normal" @click="loadDashboard">
          Coba Lagi
        </van-button>
      </div>

      <!-- Content -->
      <template v-else>
        <!-- Summary Grid -->
        <van-grid :column-num="2" :gutter="12" class="summary-grid">
          <van-grid-item>
            <div class="stat-card">
              <van-icon name="orders-o" size="24" color="#1989fa" />
              <span class="stat-value">{{ dashboard.total_po_active ?? 0 }}</span>
              <span class="stat-label">PO Aktif</span>
            </div>
          </van-grid-item>
          <van-grid-item>
            <div class="stat-card">
              <van-icon name="passed" size="24" color="#07c160" />
              <span class="stat-value">{{ dashboard.total_po_completed ?? 0 }}</span>
              <span class="stat-label">PO Selesai</span>
            </div>
          </van-grid-item>
          <van-grid-item>
            <div class="stat-card">
              <van-icon name="bill-o" size="24" color="#ff976a" />
              <span class="stat-value">{{ dashboard.total_invoice_pending ?? 0 }}</span>
              <span class="stat-label">Invoice Pending</span>
            </div>
          </van-grid-item>
          <van-grid-item>
            <div class="stat-card">
              <van-icon name="gold-coin-o" size="24" color="#07c160" />
              <span class="stat-value">{{ formatRupiah(dashboard.total_payment_received ?? 0) }}</span>
              <span class="stat-label">Pembayaran Diterima</span>
            </div>
          </van-grid-item>
        </van-grid>

        <!-- Quick Actions -->
        <div class="section-card">
          <h3 class="section-title">Aksi Cepat</h3>
          <van-cell-group inset>
            <van-cell
              title="Lihat PO"
              icon="orders-o"
              is-link
              @click="$router.push('/supplier-po')"
            />
            <van-cell
              title="Buat Invoice"
              icon="bill-o"
              is-link
              @click="$router.push('/supplier-invoices')"
            />
            <van-cell
              title="Katalog Produk"
              icon="shop-o"
              is-link
              @click="$router.push('/supplier-products')"
            />
          </van-cell-group>
        </div>
      </template>
    </div>
  </van-pull-refresh>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { showToast } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'

const loading = ref(false)
const refreshing = ref(false)
const error = ref(null)
const dashboard = ref({})

function formatRupiah(value) {
  if (!value && value !== 0) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

async function loadDashboard() {
  loading.value = true
  error.value = null
  try {
    const res = await supplierPortalService.getDashboard()
    dashboard.value = res.data?.data ?? res.data ?? {}
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal memuat dashboard'
    showToast(error.value)
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  try {
    const res = await supplierPortalService.getDashboard()
    dashboard.value = res.data?.data ?? res.data ?? {}
  } catch (e) {
    showToast('Gagal memperbarui data')
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.supplier-dashboard-page {
  min-height: 100vh;
  background: #F7F8FA;
  padding-bottom: 88px;
}

.dash-header {
  background: #F82C17;
  padding: 20px 20px 28px;
  border-radius: 0 0 24px 24px;
}

.dash-header__title {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.dash-header__subtitle {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  margin: 4px 0 0;
}

.page-loading {
  padding: 60px 0;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 20px;
  text-align: center;
}

.error-state__message {
  font-size: 14px;
  color: #969799;
  margin: 16px 0;
}

.summary-grid {
  margin: 16px;
}

.summary-grid :deep(.van-grid-item__content) {
  padding: 0;
  background: transparent;
}

.stat-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 16px 12px;
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  width: 100%;
}

.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #303030;
}

.stat-label {
  font-size: 12px;
  color: #969799;
}

.section-card {
  margin: 0 16px 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: #303030;
  margin: 0 0 12px 4px;
}
</style>
