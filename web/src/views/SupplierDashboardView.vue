<template>
  <div class="supplier-dashboard">
    <a-page-header
      title="Dashboard Supplier"
      sub-title="Ringkasan pesanan dan pembayaran"
    />

    <!-- Stat Cards -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic title="PO Aktif" :value="dashboard.total_po_active || 0">
            <template #prefix>🛍️</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic title="PO Selesai" :value="dashboard.total_po_completed || 0">
            <template #prefix>✅</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic title="Invoice Pending" :value="dashboard.total_invoice_pending || 0">
            <template #prefix>🧾</template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card :loading="loading">
          <a-statistic title="Pembayaran Diterima" :value="dashboard.total_payment_received || 0" :precision="0">
            <template #prefix>💰</template>
            <template #formatter="{ value }">
              {{ formatRupiah(value) }}
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <!-- Recent POs -->
      <a-col :span="12">
        <a-card title="PO Terbaru" :loading="loading">
          <a-table
            :columns="poColumns"
            :data-source="recentPOs"
            :pagination="false"
            row-key="id"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="getPOStatusColor(record.status)">{{ record.status }}</a-tag>
              </template>
              <template v-else-if="column.key === 'total_amount'">
                {{ formatRupiah(record.total_amount) }}
              </template>
            </template>
          </a-table>
          <a-empty v-if="recentPOs.length === 0" description="Belum ada PO" />
        </a-card>
      </a-col>

      <!-- Recent Invoices -->
      <a-col :span="12">
        <a-card title="Invoice Terbaru" :loading="loading">
          <a-table
            :columns="invoiceColumns"
            :data-source="recentInvoices"
            :pagination="false"
            row-key="id"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.status === 'paid' ? 'green' : 'orange'">
                  {{ record.status === 'paid' ? 'Dibayar' : 'Pending' }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'amount'">
                {{ formatRupiah(record.amount) }}
              </template>
            </template>
          </a-table>
          <a-empty v-if="recentInvoices.length === 0" description="Belum ada invoice" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import api from '@/services/api'

const loading = ref(false)
const dashboard = reactive({
  total_po_active: 0,
  total_po_completed: 0,
  total_invoice_pending: 0,
  total_payment_received: 0
})
const recentPOs = ref([])
const recentInvoices = ref([])

const poColumns = [
  { title: 'Nomor PO', dataIndex: 'po_number', key: 'po_number' },
  { title: 'Total', key: 'total_amount' },
  { title: 'Status', key: 'status', width: 100 }
]

const invoiceColumns = [
  { title: 'Nomor Invoice', dataIndex: 'invoice_number', key: 'invoice_number' },
  { title: 'Jumlah', key: 'amount' },
  { title: 'Status', key: 'status', width: 100 }
]

const fetchDashboard = async () => {
  loading.value = true
  try {
    const response = await api.get('/supplier/dashboard')
    const data = response.data.data || response.data
    Object.assign(dashboard, data)
    recentPOs.value = data.recent_pos || data.recent_purchase_orders || []
    recentInvoices.value = data.recent_invoices || []
  } catch (error) {
    message.error('Gagal memuat dashboard supplier')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const getPOStatusColor = (status) => {
  const colors = { pending: 'orange', approved: 'blue', received: 'green', cancelled: 'red' }
  return colors[status] || 'default'
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

onMounted(() => {
  fetchDashboard()
})
</script>

<style scoped>
.supplier-dashboard {
  padding: 24px;
}
</style>
