<template>
  <div class="rab-detail">
    <a-page-header
      title="Detail RAB"
      :sub-title="rab?.rab_number || ''"
      @back="() => $router.push('/rab')"
    >
      <template #extra>
        <!-- Action buttons based on status + role -->
        <template v-if="rab">
          <a-button
            v-if="rab.status === 'draft' && can('RAB_APPROVE_SPPG')"
            type="primary"
            :loading="actionLoading"
            @click="handleApproveSPPG"
          >
            Approve SPPG
          </a-button>
          <a-button
            v-if="rab.status === 'approved_sppg' && can('RAB_APPROVE_YAYASAN')"
            type="primary"
            :loading="actionLoading"
            @click="handleApproveYayasan"
          >
            Approve Yayasan
          </a-button>
          <a-button
            v-if="rab.status === 'approved_sppg' && can('RAB_APPROVE_YAYASAN')"
            danger
            @click="showRejectModal"
          >
            Tolak
          </a-button>
        </template>
      </template>
    </a-page-header>

    <a-spin :spinning="loading">
      <template v-if="rab">
        <!-- Revision warning -->
        <a-alert
          v-if="rab.status === 'revision_requested' && rab.revision_notes"
          type="warning"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #message>Catatan Revisi dari Yayasan</template>
          <template #description>
            {{ rab.revision_notes }}
            <br /><br />
            <strong>Silakan revisi menu plan di halaman Perencanaan Menu, lalu approve ulang untuk generate RAB baru.</strong>
          </template>
        </a-alert>

        <!-- Header info -->
        <a-card style="margin-bottom: 16px">
          <a-descriptions bordered :column="{ xs: 1, sm: 2, md: 3 }">
            <a-descriptions-item label="Nomor RAB">
              <strong>{{ rab.rab_number }}</strong>
            </a-descriptions-item>
            <a-descriptions-item label="Menu Plan">
              {{ rab.menu_plan ? `Minggu ${formatDate(rab.menu_plan.week_start)} - ${formatDate(rab.menu_plan.week_end)}` : '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="SPPG">
              {{ rab.sppg?.nama || rab.sppg?.name || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Status">
              <a-tag :color="getStatusColor(rab.status)">
                {{ getStatusLabel(rab.status) }}
              </a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="Total">
              <strong>{{ formatRupiah(rab.total_amount) }}</strong>
            </a-descriptions-item>
            <a-descriptions-item label="Dibuat Oleh">
              {{ rab.creator?.full_name || rab.creator?.name || '-' }}
            </a-descriptions-item>
            <a-descriptions-item label="Tanggal Dibuat">
              {{ formatDate(rab.created_at) }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- Items table -->
        <a-card title="Item RAB" style="margin-bottom: 16px">
          <a-table
            :columns="itemColumns"
            :data-source="rab.items || []"
            :pagination="false"
            row-key="id"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'ingredient'">
                {{ record.ingredient?.name || '-' }}
              </template>
              <template v-else-if="column.key === 'unit_price'">
                {{ formatRupiah(record.unit_price) }}
              </template>
              <template v-else-if="column.key === 'subtotal'">
                {{ formatRupiah(record.subtotal) }}
              </template>
              <template v-else-if="column.key === 'recommended_supplier'">
                {{ record.recommended_supplier?.name || '-' }}
              </template>
              <template v-else-if="column.key === 'item_status'">
                <a-tag :color="getItemStatusColor(record.status)">
                  {{ getItemStatusLabel(record.status) }}
                </a-tag>
              </template>
              <template v-else-if="column.key === 'po_id'">
                {{ record.po_id || '-' }}
              </template>
              <template v-else-if="column.key === 'grn_id'">
                {{ record.grn_id || '-' }}
              </template>
            </template>
          </a-table>
        </a-card>

        <!-- Tabs: PO Tracking & Perbandingan -->
        <a-card>
          <a-tabs v-model:activeKey="activeTab" @change="handleTabChange">
            <a-tab-pane key="po-tracking" tab="PO Tracking">
              <a-spin :spinning="trackingLoading">
                <a-table
                  :columns="poTrackingColumns"
                  :data-source="poTracking"
                  :pagination="false"
                  row-key="id"
                  size="small"
                >
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'status'">
                      <a-tag :color="getPOStatusColor(record.status)">
                        {{ record.status }}
                      </a-tag>
                    </template>
                    <template v-else-if="column.key === 'total_amount'">
                      {{ formatRupiah(record.total_amount) }}
                    </template>
                    <template v-else-if="column.key === 'grn_status'">
                      <a-tag v-if="record.grn_status" color="green">{{ record.grn_status }}</a-tag>
                      <span v-else>-</span>
                    </template>
                  </template>
                </a-table>
                <a-empty v-if="!trackingLoading && poTracking.length === 0" description="Belum ada PO terkait" />
              </a-spin>
            </a-tab-pane>
            <a-tab-pane key="comparison" tab="Perbandingan">
              <a-spin :spinning="comparisonLoading">
                <a-table
                  :columns="comparisonColumns"
                  :data-source="comparison"
                  :pagination="false"
                  row-key="id"
                  size="small"
                >
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'planned_amount'">
                      {{ formatRupiah(record.planned_amount) }}
                    </template>
                    <template v-else-if="column.key === 'actual_amount'">
                      {{ formatRupiah(record.actual_amount) }}
                    </template>
                    <template v-else-if="column.key === 'difference'">
                      <span :style="{ color: (record.actual_amount - record.planned_amount) > 0 ? '#f5222d' : '#52c41a' }">
                        {{ formatRupiah(record.actual_amount - record.planned_amount) }}
                      </span>
                    </template>
                  </template>
                </a-table>
                <a-empty v-if="!comparisonLoading && comparison.length === 0" description="Belum ada data perbandingan" />
              </a-spin>
            </a-tab-pane>
          </a-tabs>
        </a-card>
      </template>
    </a-spin>

    <!-- Reject Modal -->
    <a-modal
      v-model:open="rejectModalVisible"
      title="Tolak RAB"
      :confirm-loading="actionLoading"
      @ok="handleReject"
      @cancel="rejectModalVisible = false"
      ok-text="Tolak"
      :ok-button-props="{ danger: true }"
    >
      <a-form layout="vertical">
        <a-form-item label="Alasan Penolakan" required>
          <a-textarea
            v-model:value="revisionNotes"
            :rows="4"
            placeholder="Masukkan alasan penolakan atau catatan revisi..."
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import rabService from '@/services/rabService'
import { usePermissions } from '@/composables/usePermissions'

const route = useRoute()
const router = useRouter()
const { can } = usePermissions()

const loading = ref(false)
const actionLoading = ref(false)
const trackingLoading = ref(false)
const comparisonLoading = ref(false)
const rab = ref(null)
const poTracking = ref([])
const comparison = ref([])
const activeTab = ref('po-tracking')
const rejectModalVisible = ref(false)
const revisionNotes = ref('')

const itemColumns = [
  { title: 'Bahan', key: 'ingredient' },
  { title: 'Qty', dataIndex: 'quantity', key: 'quantity', width: 80 },
  { title: 'Satuan', dataIndex: 'unit', key: 'unit', width: 80 },
  { title: 'Harga Satuan', key: 'unit_price', width: 140 },
  { title: 'Subtotal', key: 'subtotal', width: 140 },
  { title: 'Supplier Rekomendasi', key: 'recommended_supplier' },
  { title: 'Status Item', key: 'item_status', width: 130 },
  { title: 'PO', key: 'po_id', width: 60 },
  { title: 'GRN', key: 'grn_id', width: 60 }
]

const poTrackingColumns = [
  { title: 'Nomor PO', dataIndex: 'po_number', key: 'po_number' },
  { title: 'Supplier', dataIndex: ['supplier', 'name'], key: 'supplier' },
  { title: 'Total', key: 'total_amount' },
  { title: 'Status PO', key: 'status' },
  { title: 'Status GRN', key: 'grn_status' }
]

const comparisonColumns = [
  { title: 'Bahan', dataIndex: ['ingredient', 'name'], key: 'ingredient' },
  { title: 'Qty Rencana', dataIndex: 'planned_quantity', key: 'planned_quantity' },
  { title: 'Qty Aktual', dataIndex: 'actual_quantity', key: 'actual_quantity' },
  { title: 'Biaya Rencana', key: 'planned_amount' },
  { title: 'Biaya Aktual', key: 'actual_amount' },
  { title: 'Selisih', key: 'difference' }
]

const fetchRABDetail = async () => {
  loading.value = true
  try {
    const response = await rabService.getRABDetail(route.params.id)
    rab.value = response.data.rab || response.data.data || response.data
  } catch (error) {
    message.error('Gagal memuat detail RAB')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchPOTracking = async () => {
  trackingLoading.value = true
  try {
    const response = await rabService.getRABPOTracking(route.params.id)
    poTracking.value = response.data.purchase_orders || response.data.data || []
  } catch (error) {
    console.error('Gagal memuat PO tracking:', error)
  } finally {
    trackingLoading.value = false
  }
}

const fetchComparison = async () => {
  comparisonLoading.value = true
  try {
    const response = await rabService.getRABComparison(route.params.id)
    comparison.value = response.data.items || response.data.data || []
  } catch (error) {
    console.error('Gagal memuat perbandingan:', error)
  } finally {
    comparisonLoading.value = false
  }
}

const handleTabChange = (key) => {
  if (key === 'po-tracking' && poTracking.value.length === 0) {
    fetchPOTracking()
  } else if (key === 'comparison' && comparison.value.length === 0) {
    fetchComparison()
  }
}

const handleApproveSPPG = async () => {
  actionLoading.value = true
  try {
    await rabService.approveSPPG(route.params.id)
    message.success('RAB berhasil di-approve SPPG')
    fetchRABDetail()
  } catch (error) {
    message.error(error.response?.data?.message || 'Gagal approve RAB')
  } finally {
    actionLoading.value = false
  }
}

const handleApproveYayasan = async () => {
  actionLoading.value = true
  try {
    await rabService.approveYayasan(route.params.id)
    message.success('RAB berhasil di-approve Yayasan')
    fetchRABDetail()
  } catch (error) {
    message.error(error.response?.data?.message || 'Gagal approve RAB')
  } finally {
    actionLoading.value = false
  }
}

const showRejectModal = () => {
  revisionNotes.value = ''
  rejectModalVisible.value = true
}

const handleReject = async () => {
  if (!revisionNotes.value.trim()) {
    message.error('Alasan penolakan wajib diisi')
    return
  }
  actionLoading.value = true
  try {
    await rabService.rejectRAB(route.params.id, { notes: revisionNotes.value })
    message.success('RAB berhasil ditolak')
    rejectModalVisible.value = false
    fetchRABDetail()
  } catch (error) {
    message.error(error.response?.data?.message || 'Gagal menolak RAB')
  } finally {
    actionLoading.value = false
  }
}

const handleResubmit = async () => {
  actionLoading.value = true
  try {
    await rabService.resubmitRAB(route.params.id)
    message.success('RAB berhasil dikirim ulang')
    fetchRABDetail()
  } catch (error) {
    message.error(error.response?.data?.message || 'Gagal mengirim ulang RAB')
  } finally {
    actionLoading.value = false
  }
}

const handleEdit = () => {
  // Edit inline — for now just show message; could open modal
  message.info('Fitur edit RAB akan segera tersedia')
}

const getStatusColor = (status) => {
  const colors = { draft: 'default', approved_sppg: 'blue', approved_yayasan: 'green', revision_requested: 'orange', completed: 'purple' }
  return colors[status] || 'default'
}

const getStatusLabel = (status) => {
  const labels = { draft: 'Draft', approved_sppg: 'Approved SPPG', approved_yayasan: 'Approved Yayasan', revision_requested: 'Revisi', completed: 'Selesai' }
  return labels[status] || status
}

const getItemStatusColor = (status) => {
  const colors = { pending: 'default', po_created: 'blue', grn_received: 'green' }
  return colors[status] || 'default'
}

const getItemStatusLabel = (status) => {
  const labels = { pending: 'Pending', po_created: 'PO Dibuat', grn_received: 'GRN Diterima' }
  return labels[status] || status
}

const getPOStatusColor = (status) => {
  const colors = { pending: 'orange', approved: 'blue', received: 'green', cancelled: 'red' }
  return colors[status] || 'default'
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(() => {
  fetchRABDetail()
  fetchPOTracking()
})
</script>

<style scoped>
.rab-detail {
  padding: 24px;
}
</style>
