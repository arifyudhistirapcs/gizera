<template>
  <div class="supplier-po-page">
    <!-- Header -->
    <van-nav-bar title="Purchase Order" fixed />

    <div class="page-content">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <!-- Loading -->
        <van-loading v-if="loading" type="spinner" vertical class="page-loading">
          Memuat PO...
        </van-loading>

        <!-- Empty -->
        <van-empty
          v-else-if="orders.length === 0"
          image="search"
          description="Belum ada Purchase Order"
        />

        <!-- PO List -->
        <van-cell-group v-else inset>
          <van-cell
            v-for="po in orders"
            :key="po.id"
            :title="po.po_number || '-'"
            :label="formatDate(po.created_at || po.order_date)"
            is-link
            @click="$router.push(`/supplier-po/${po.id}`)"
          >
            <template #value>
              <div class="po-cell-value">
                <span class="po-amount">{{ formatRupiah(po.total_amount) }}</span>
                <van-tag :type="getStatusType(po.status)" size="medium">
                  {{ getStatusText(po.status) }}
                </van-tag>
              </div>
            </template>
            <template #label>
              <div class="po-label">
                <span>{{ po.yayasan?.name || po.yayasan_name || '' }}</span>
                <span>{{ formatDate(po.created_at || po.order_date) }}</span>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
      </van-pull-refresh>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { showToast } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'

const loading = ref(false)
const refreshing = ref(false)
const orders = ref([])

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
    month: 'short',
    year: 'numeric'
  })
}

function getStatusType(status) {
  const map = {
    draft: 'default',
    pending: 'warning',
    approved: 'primary',
    sent: 'primary',
    received: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return map[status] || 'default'
}

function getStatusText(status) {
  const map = {
    draft: 'Draft',
    pending: 'Menunggu',
    approved: 'Disetujui',
    sent: 'Dikirim',
    received: 'Diterima',
    completed: 'Selesai',
    cancelled: 'Dibatalkan'
  }
  return map[status] || status || '-'
}

async function loadOrders() {
  loading.value = true
  try {
    const res = await supplierPortalService.getPurchaseOrders()
    orders.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal memuat PO')
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  try {
    const res = await supplierPortalService.getPurchaseOrders()
    orders.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast('Gagal memperbarui data')
  } finally {
    refreshing.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.supplier-po-page {
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

.po-cell-value {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.po-amount {
  font-size: 14px;
  font-weight: 600;
  color: #303030;
}

.po-label {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
  color: #969799;
}
</style>
