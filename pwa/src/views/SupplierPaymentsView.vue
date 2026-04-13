<template>
  <div class="supplier-payments-page">
    <!-- Header -->
    <van-nav-bar title="Riwayat Pembayaran" fixed />

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <!-- Loading -->
        <van-loading v-if="loading" type="spinner" vertical class="page-loading">
          Memuat pembayaran...
        </van-loading>

        <template v-else>
          <!-- Total Summary -->
          <div class="total-card">
            <span class="total-label">Total Pembayaran Diterima</span>
            <span class="total-value">{{ formatRupiah(totalPayments) }}</span>
          </div>

          <!-- Empty -->
          <van-empty
            v-if="payments.length === 0"
            image="search"
            description="Belum ada riwayat pembayaran"
          />

          <!-- Payment List -->
          <van-cell-group v-else inset>
            <van-cell
              v-for="payment in payments"
              :key="payment.id"
              :title="payment.invoice_number || payment.invoice?.invoice_number || '-'"
              :label="formatDate(payment.payment_date)"
            >
              <template #value>
                <div class="payment-value">
                  <span class="payment-amount">{{ formatRupiah(payment.amount) }}</span>
                  <span class="payment-method">{{ payment.payment_method || '-' }}</span>
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </template>
      </van-pull-refresh>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'

const loading = ref(false)
const refreshing = ref(false)
const payments = ref([])

const totalPayments = computed(() =>
  payments.value.reduce((sum, p) => sum + (Number(p.amount) || 0), 0)
)

function formatRupiah(value) {
  if (!value && value !== 0) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

async function loadPayments() {
  loading.value = true
  try {
    const res = await supplierPortalService.getPayments()
    payments.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal memuat pembayaran')
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  try {
    const res = await supplierPortalService.getPayments()
    payments.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast('Gagal memperbarui data')
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  loadPayments()
})
</script>

<style scoped>
.supplier-payments-page {
  min-height: 100vh;
  background: #F7F8FA;
  padding-top: 46px;
  padding-bottom: 88px;
}

.page-content {
  padding: 16px 0;
}

.page-loading {
  padding: 60px 0;
}

.total-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  margin: 0 16px 16px;
  padding: 20px;
  background: linear-gradient(135deg, #07c160, #06ad56);
  border-radius: 14px;
  color: #fff;
}

.total-label {
  font-size: 13px;
  opacity: 0.9;
}

.total-value {
  font-size: 24px;
  font-weight: 700;
}

.payment-value {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.payment-amount {
  font-size: 14px;
  font-weight: 600;
  color: #303030;
}

.payment-method {
  font-size: 11px;
  color: #969799;
}
</style>
