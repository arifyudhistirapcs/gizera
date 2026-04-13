<template>
  <div class="supplier-po-detail-page">
    <!-- Header -->
    <van-nav-bar title="Detail PO" left-arrow fixed @click-left="$router.back()" />

    <div class="page-content">
      <!-- Loading -->
      <van-loading v-if="loading" type="spinner" vertical class="page-loading">
        Memuat detail PO...
      </van-loading>

      <!-- Error -->
      <van-empty v-else-if="error" image="error" :description="error">
        <van-button type="primary" size="small" @click="loadDetail">Coba Lagi</van-button>
      </van-empty>

      <template v-else-if="po">
        <!-- PO Header -->
        <van-cell-group inset title="Informasi PO">
          <van-cell title="No. PO" :value="po.po_number || '-'" />
          <van-cell title="Yayasan" :value="po.yayasan?.name || po.yayasan_name || '-'" />
          <van-cell title="SPPG Tujuan" :value="po.target_sppg?.name || po.target_sppg_name || '-'" />
          <van-cell title="Tanggal" :value="formatDate(po.created_at || po.order_date)" />
          <van-cell title="Total" :value="formatRupiah(po.total_amount)" />
          <van-cell title="Status">
            <template #value>
              <van-tag :type="getStatusType(po.status)" size="medium">
                {{ getStatusText(po.status) }}
              </van-tag>
            </template>
          </van-cell>
        </van-cell-group>

        <!-- PO Items -->
        <van-cell-group inset title="Daftar Item" class="items-group">
          <van-cell
            v-for="(item, idx) in (po.items || [])"
            :key="item.id || idx"
            :title="item.ingredient?.name || item.ingredient_name || `Item ${idx + 1}`"
            :label="`${item.quantity ?? 0} ${item.unit || ''}`"
            :value="formatRupiah(item.subtotal || (item.quantity * item.unit_price))"
          />
          <van-empty
            v-if="!po.items || po.items.length === 0"
            image="search"
            description="Tidak ada item"
          />
        </van-cell-group>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { showToast } from 'vant'
import api from '@/services/api'

const route = useRoute()
const loading = ref(false)
const error = ref(null)
const po = ref(null)

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

async function loadDetail() {
  loading.value = true
  error.value = null
  try {
    const res = await api.get(`/purchase-orders/${route.params.id}`)
    po.value = res.data?.data ?? res.data ?? null
  } catch (e) {
    error.value = e.response?.data?.message || 'Gagal memuat detail PO'
    showToast(error.value)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.supplier-po-detail-page {
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

.items-group {
  margin-top: 16px;
}
</style>
