<template>
  <div class="supplier-po-list">
    <a-page-header
      title="Purchase Order"
      sub-title="Daftar PO yang ditujukan kepada Anda"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari nomor PO..."
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </a-input>
          </a-col>
        </a-row>

        <a-table
          :columns="columns"
          :data-source="purchaseOrders"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
          :custom-row="customRow"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'total_amount'">
              {{ formatRupiah(record.total_amount) }}
            </template>
            <template v-else-if="column.key === 'order_date'">
              {{ formatDate(record.order_date) }}
            </template>
            <template v-else-if="column.key === 'yayasan'">
              {{ record.yayasan?.nama || record.yayasan?.name || '-' }}
            </template>
            <template v-else-if="column.key === 'target_sppg'">
              {{ record.target_sppg?.nama || record.target_sppg?.name || '-' }}
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- PO Detail Modal -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Detail Purchase Order"
      width="900px"
      @cancel="closeDetailModal"
    >
      <template #footer>
        <a-space>
          <a-button @click="exportPOToPDF(selectedPO)">
            <template #icon><FilePdfOutlined /></template>
            Export PDF
          </a-button>

          <!-- Actions for pending status -->
          <template v-if="selectedPO?.status === 'pending' && !showRevisionForm">
            <a-button
              type="primary"
              style="background-color: #52c41a; border-color: #52c41a"
              :loading="actionLoading"
              @click="handleConfirmPO"
            >
              <template #icon><CheckOutlined /></template>
              Konfirmasi PO
            </a-button>
            <a-button
              type="primary"
              style="background-color: #fa8c16; border-color: #fa8c16"
              @click="openRevisionForm"
            >
              <template #icon><EditOutlined /></template>
              Ajukan Perubahan
            </a-button>
          </template>

          <!-- Submit revision button -->
          <template v-if="showRevisionForm">
            <a-button @click="cancelRevisionForm">Batal</a-button>
            <a-button
              type="primary"
              style="background-color: #fa8c16; border-color: #fa8c16"
              :loading="actionLoading"
              :disabled="!revisionNotes.trim()"
              @click="handleSubmitRevision"
            >
              Kirim Perubahan
            </a-button>
          </template>

          <!-- Actions for approved status: mark as shipping -->
          <a-button
            v-if="(selectedPO?.status === 'approved' || selectedPO?.status === 'confirmed') && !showRevisionForm"
            type="primary"
            style="background-color: #1890ff; border-color: #1890ff"
            :loading="actionLoading"
            @click="handleMarkAsShipping"
          >
            <template #icon><SendOutlined /></template>
            Kirim Barang
          </a-button>

          <a-button v-if="!showRevisionForm" @click="closeDetailModal">Tutup</a-button>
        </a-space>
      </template>

      <template v-if="selectedPO">
        <a-descriptions bordered :column="2">
          <a-descriptions-item label="Nomor PO" :span="2">
            <strong>{{ selectedPO.po_number }}</strong>
          </a-descriptions-item>
          <a-descriptions-item label="Yayasan">
            {{ selectedPO.yayasan?.nama || selectedPO.yayasan?.name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="SPPG Tujuan">
            {{ selectedPO.target_sppg?.nama || selectedPO.target_sppg?.name || '-' }}
          </a-descriptions-item>
          <a-descriptions-item label="Tanggal Order">
            {{ formatDate(selectedPO.order_date) }}
          </a-descriptions-item>
          <a-descriptions-item label="Status">
            <a-tag :color="getStatusColor(selectedPO.status)">
              {{ getStatusText(selectedPO.status) }}
            </a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="Total" :span="2">
            <strong>{{ formatRupiah(selectedPO.total_amount) }}</strong>
          </a-descriptions-item>
        </a-descriptions>

        <!-- Status: revision_by_supplier -->
        <a-alert
          v-if="selectedPO.status === 'revision_by_supplier'"
          type="info"
          show-icon
          style="margin-top: 16px"
          message="Menunggu review yayasan"
          :description="selectedPO.supplier_revision_notes ? `Catatan revisi Anda: ${selectedPO.supplier_revision_notes}` : 'Perubahan Anda sedang ditinjau oleh yayasan.'"
        />

        <!-- Status: approved (after supplier confirms) -->
        <a-alert
          v-if="selectedPO.status === 'approved'"
          type="success"
          show-icon
          style="margin-top: 16px"
          message="PO disetujui"
          description="PO telah disetujui. Silakan klik 'Kirim Barang' untuk memulai pengiriman ke SPPG tujuan."
        />

        <!-- Status: shipping -->
        <a-alert
          v-if="selectedPO.status === 'shipping'"
          type="info"
          show-icon
          style="margin-top: 16px"
          message="Barang sedang dikirim"
          description="Barang sedang dalam proses pengiriman ke SPPG tujuan. Menunggu penerimaan (GRN) dari pihak SPPG."
        />

        <!-- Revision Form (shown when supplier clicks "Ajukan Perubahan") -->
        <template v-if="showRevisionForm">
          <a-divider>Form Perubahan</a-divider>

          <a-alert
            type="warning"
            show-icon
            style="margin-bottom: 16px"
            message="Anda dapat mengubah jumlah, harga satuan, atau menghapus item. Perubahan akan dikirim ke yayasan untuk ditinjau."
          />

          <a-table
            :columns="revisionItemColumns"
            :data-source="revisionItems"
            :pagination="false"
            size="small"
            row-key="id"
          >
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'ingredient_name'">
                {{ record.ingredient?.name || '-' }}
              </template>
              <template v-else-if="column.key === 'unit'">
                {{ record.ingredient?.unit || '-' }}
              </template>
              <template v-else-if="column.key === 'quantity'">
                <a-input-number
                  v-model:value="record.quantity"
                  :min="0.01"
                  :step="0.1"
                  style="width: 100%"
                  @change="updateRevisionSubtotal(index)"
                />
              </template>
              <template v-else-if="column.key === 'unit_price'">
                <a-input-number
                  v-model:value="record.unit_price"
                  :min="0"
                  :step="1000"
                  style="width: 100%"
                  :formatter="value => `Rp ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')"
                  :parser="value => value.replace(/Rp\s?|(,*)/g, '')"
                  @change="updateRevisionSubtotal(index)"
                />
              </template>
              <template v-else-if="column.key === 'subtotal'">
                {{ formatRupiah((record.quantity || 0) * (record.unit_price || 0)) }}
              </template>
              <template v-else-if="column.key === 'actions'">
                <a-popconfirm
                  title="Hapus item ini?"
                  ok-text="Ya"
                  cancel-text="Batal"
                  @confirm="removeRevisionItem(index)"
                >
                  <a-button type="link" size="small" danger :disabled="revisionItems.length <= 1">
                    Hapus
                  </a-button>
                </a-popconfirm>
              </template>
            </template>
          </a-table>

          <div style="text-align: right; margin: 12px 0">
            <strong>Total: {{ formatRupiah(revisionTotal) }}</strong>
          </div>

          <a-form-item
            label="Catatan Revisi"
            :validate-status="!revisionNotes.trim() ? 'error' : ''"
            :help="!revisionNotes.trim() ? 'Catatan revisi wajib diisi' : ''"
            style="margin-top: 16px"
          >
            <a-textarea
              v-model:value="revisionNotes"
              placeholder="Jelaskan alasan perubahan..."
              :rows="3"
              :maxlength="500"
              show-count
            />
          </a-form-item>
        </template>

        <!-- Normal item table (hidden when revision form is open) -->
        <template v-if="!showRevisionForm">
          <a-divider>Item Pesanan</a-divider>

          <a-table
            :columns="itemColumns"
            :data-source="selectedPO.po_items || []"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'unit_price'">
                {{ formatRupiah(record.unit_price) }}
              </template>
              <template v-else-if="column.key === 'subtotal'">
                {{ formatRupiah(record.subtotal) }}
              </template>
            </template>
          </a-table>
        </template>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, FilePdfOutlined, CheckOutlined, EditOutlined, SendOutlined } from '@ant-design/icons-vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import purchaseOrderService from '@/services/purchaseOrderService'

const loading = ref(false)
const actionLoading = ref(false)
const purchaseOrders = ref([])
const searchText = ref('')
const detailModalVisible = ref(false)
const selectedPO = ref(null)

// Revision form state
const showRevisionForm = ref(false)
const revisionItems = ref([])
const revisionNotes = ref('')

const revisionTotal = computed(() => {
  return revisionItems.value.reduce((sum, item) => {
    return sum + (item.quantity || 0) * (item.unit_price || 0)
  }, 0)
})

const revisionItemColumns = [
  { title: 'Bahan', key: 'ingredient_name', width: 200 },
  { title: 'Satuan', key: 'unit', width: 80 },
  { title: 'Jumlah', key: 'quantity', width: 120 },
  { title: 'Harga Satuan', key: 'unit_price', width: 160 },
  { title: 'Subtotal', key: 'subtotal', width: 160 },
  { title: 'Aksi', key: 'actions', width: 80 }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'Nomor PO', dataIndex: 'po_number', key: 'po_number' },
  { title: 'Yayasan', key: 'yayasan' },
  { title: 'SPPG Tujuan', key: 'target_sppg' },
  { title: 'Tanggal', key: 'order_date' },
  { title: 'Total', key: 'total_amount' },
  { title: 'Status', key: 'status', width: 120 }
]

const itemColumns = [
  { title: 'Bahan', dataIndex: ['ingredient', 'name'], key: 'ingredient_name' },
  { title: 'Jumlah', dataIndex: 'quantity', key: 'quantity' },
  { title: 'Satuan', dataIndex: ['ingredient', 'unit'], key: 'unit' },
  { title: 'Harga Satuan', key: 'unit_price' },
  { title: 'Subtotal', key: 'subtotal' }
]

const fetchPurchaseOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value || undefined
    }
    const response = await purchaseOrderService.getPurchaseOrders(params)
    purchaseOrders.value = response.data.purchase_orders || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat data purchase order')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchPurchaseOrders()
}

const handleSearch = () => {
  pagination.current = 1
  fetchPurchaseOrders()
}

const customRow = (record) => ({
  onClick: () => viewPO(record),
  style: { cursor: 'pointer' }
})

const viewPO = async (po) => {
  try {
    const response = await purchaseOrderService.getPurchaseOrder(po.id)
    selectedPO.value = response.data.purchase_order
    detailModalVisible.value = true
  } catch (error) {
    message.error('Gagal memuat detail PO')
    console.error(error)
  }
}

const getStatusColor = (status) => {
  const colors = {
    pending: 'orange',
    revision_by_supplier: 'orange',
    approved: 'blue',
    shipping: 'geekblue',
    received: 'green',
    cancelled: 'red'
  }
  return colors[status] || 'default'
}

const getStatusText = (status) => {
  const texts = {
    pending: 'Pending',
    revision_by_supplier: 'Revisi oleh Supplier',
    approved: 'Disetujui',
    shipping: 'Sedang Dikirim',
    received: 'Diterima',
    cancelled: 'Dibatalkan'
  }
  return texts[status] || status
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const exportPOToPDF = (po) => {
  if (!po) return
  const doc = new jsPDF()

  // Title
  doc.setFontSize(18)
  doc.setTextColor(0, 0, 0)
  doc.text('PURCHASE ORDER', 105, 20, { align: 'center' })

  // PO Number subtitle
  doc.setFontSize(12)
  doc.text(po.po_number || '', 105, 28, { align: 'center' })

  // PO Info
  doc.setFontSize(10)
  doc.text(`Nomor PO: ${po.po_number || '-'}`, 14, 42)
  doc.text(`Supplier: ${po.supplier?.name || '-'}`, 14, 49)
  doc.text(`Yayasan: ${po.yayasan?.nama || po.yayasan?.name || '-'}`, 14, 56)
  doc.text(`SPPG Tujuan: ${po.target_sppg?.nama || po.target_sppg?.name || '-'}`, 14, 63)
  doc.text(`Tanggal Order: ${formatDate(po.order_date)}`, 14, 70)
  doc.text(`Tanggal Pengiriman: ${formatDate(po.expected_delivery)}`, 14, 77)
  doc.text(`Status: ${getStatusText(po.status)}`, 14, 84)

  // Items table
  const items = (po.po_items || []).map(item => [
    item.ingredient?.name || '-',
    item.quantity?.toString() || '0',
    item.ingredient?.unit || '-',
    formatRupiah(item.unit_price),
    formatRupiah(item.subtotal)
  ])

  const tableResult = autoTable(doc, {
    startY: 92,
    head: [['Bahan', 'Jumlah', 'Satuan', 'Harga Satuan', 'Subtotal']],
    body: items,
    theme: 'grid',
    headStyles: { fillColor: [248, 44, 23] },
    foot: [['', '', '', 'TOTAL', formatRupiah(po.total_amount)]],
    footStyles: { fillColor: [240, 240, 240], textColor: [0, 0, 0], fontStyle: 'bold' }
  })

  // Footer
  const finalY = (doc.lastAutoTable?.finalY || tableResult?.finalY || 200) + 20
  doc.setFontSize(8)
  doc.setTextColor(150)
  doc.text('Dokumen ini digenerate otomatis oleh sistem Dapur Sehat', 105, finalY, { align: 'center' })

  doc.save(`PO-${po.po_number || 'draft'}.pdf`)
}

const closeDetailModal = () => {
  detailModalVisible.value = false
  showRevisionForm.value = false
  revisionItems.value = []
  revisionNotes.value = ''
}

const handleConfirmPO = async () => {
  actionLoading.value = true
  try {
    await purchaseOrderService.confirmBySupplier(selectedPO.value.id)
    message.success('PO berhasil dikonfirmasi')
    closeDetailModal()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal mengkonfirmasi PO'
    message.error(errMsg)
    console.error(error)
  } finally {
    actionLoading.value = false
  }
}

const handleMarkAsShipping = async () => {
  actionLoading.value = true
  try {
    await purchaseOrderService.markAsShipping(selectedPO.value.id)
    message.success('Status PO berhasil diubah ke sedang dikirim')
    closeDetailModal()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal mengubah status pengiriman'
    message.error(errMsg)
    console.error(error)
  } finally {
    actionLoading.value = false
  }
}

const openRevisionForm = () => {
  // Clone PO items into editable revision items
  revisionItems.value = (selectedPO.value.po_items || []).map(item => ({
    id: item.id,
    ingredient_id: item.ingredient_id,
    ingredient: item.ingredient ? { ...item.ingredient } : null,
    quantity: item.quantity,
    unit_price: item.unit_price
  }))
  revisionNotes.value = ''
  showRevisionForm.value = true
}

const cancelRevisionForm = () => {
  showRevisionForm.value = false
  revisionItems.value = []
  revisionNotes.value = ''
}

const updateRevisionSubtotal = () => {
  // Subtotal is computed inline in the template, nothing extra needed
}

const removeRevisionItem = (index) => {
  revisionItems.value.splice(index, 1)
}

const handleSubmitRevision = async () => {
  if (!revisionNotes.value.trim()) {
    message.warning('Catatan revisi wajib diisi')
    return
  }
  if (revisionItems.value.length === 0) {
    message.warning('Minimal satu item harus ada')
    return
  }

  actionLoading.value = true
  try {
    const payload = {
      notes: revisionNotes.value.trim(),
      items: revisionItems.value.map(item => ({
        ingredient_id: item.ingredient_id,
        quantity: item.quantity,
        unit_price: item.unit_price
      }))
    }
    await purchaseOrderService.requestRevision(selectedPO.value.id, payload)
    message.success('Perubahan berhasil dikirim ke yayasan')
    closeDetailModal()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal mengirim perubahan'
    message.error(errMsg)
    console.error(error)
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  fetchPurchaseOrders()
})
</script>

<style scoped>
.supplier-po-list {
  padding: 24px;
}
</style>
