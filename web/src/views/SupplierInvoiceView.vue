<template>
  <div class="supplier-invoice">
    <a-page-header
      title="Invoice Supplier"
      sub-title="Kelola invoice dan lihat status pembayaran"
    >
      <template #extra>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Buat Invoice
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-table
        :columns="columns"
        :data-source="invoices"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'paid' ? 'green' : 'orange'">
              {{ record.status === 'paid' ? 'Sudah Dibayar' : 'Menunggu Pembayaran' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'amount'">
            {{ formatRupiah(record.amount) }}
          </template>
          <template v-else-if="column.key === 'due_date'">
            {{ formatDate(record.due_date) }}
          </template>
          <template v-else-if="column.key === 'po_number'">
            {{ record.purchase_order?.po_number || '-' }}
          </template>
          <template v-else-if="column.key === 'payment_info'">
            <template v-if="record.payment">
              <span>{{ formatDate(record.payment.payment_date) }}</span>
              <a-tag color="green" style="margin-left: 8px">{{ record.payment.payment_method }}</a-tag>
            </template>
            <span v-else>-</span>
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-button type="link" size="small" @click="viewInvoice(record)">
              Detail
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Invoice Detail Modal -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Detail Invoice"
      width="700px"
    >
      <template #footer>
        <a-space>
          <a-button @click="exportInvoiceToPDF(selectedInvoice)">
            <template #icon><FilePdfOutlined /></template>
            Export PDF
          </a-button>
          <a-button @click="detailModalVisible = false">Tutup</a-button>
        </a-space>
      </template>

      <a-descriptions v-if="selectedInvoice" bordered :column="2">
        <a-descriptions-item label="Nomor Invoice" :span="2">
          <strong>{{ selectedInvoice.invoice_number }}</strong>
        </a-descriptions-item>
        <a-descriptions-item label="Nomor PO">
          {{ selectedInvoice.purchase_order?.po_number || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="Status">
          <a-tag :color="selectedInvoice.status === 'paid' ? 'green' : 'orange'">
            {{ selectedInvoice.status === 'paid' ? 'Sudah Dibayar' : 'Menunggu Pembayaran' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="Jumlah" :span="2">
          <strong>{{ formatRupiah(selectedInvoice.amount) }}</strong>
        </a-descriptions-item>
        <a-descriptions-item label="Jatuh Tempo" :span="2">
          {{ formatDate(selectedInvoice.due_date) }}
        </a-descriptions-item>
        <template v-if="selectedInvoice.payment">
          <a-descriptions-item label="Tanggal Bayar">
            {{ formatDate(selectedInvoice.payment.payment_date) }}
          </a-descriptions-item>
          <a-descriptions-item label="Metode Pembayaran">
            {{ selectedInvoice.payment.payment_method || '-' }}
          </a-descriptions-item>
        </template>
      </a-descriptions>
    </a-modal>

    <!-- Create Invoice Modal -->
    <a-modal
      v-model:open="createModalVisible"
      title="Buat Invoice Baru"
      :confirm-loading="submitting"
      @ok="handleCreate"
      @cancel="createModalVisible = false"
      ok-text="Buat Invoice"
    >
      <a-form ref="formRef" :model="formData" :rules="rules" layout="vertical">
        <a-form-item label="Purchase Order" name="po_id">
          <a-select
            v-model:value="formData.po_id"
            placeholder="Pilih PO (yang sudah GRN)"
            show-search
            :filter-option="filterPO"
            @change="handlePOSelect"
          >
            <a-select-option v-for="po in eligiblePOs" :key="po.id" :value="po.id">
              {{ po.po_number }} — {{ formatRupiah(po.total_amount) }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Jumlah (Rp)" name="amount">
          <a-input-number v-model:value="formData.amount" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="Jatuh Tempo" name="due_date">
          <a-date-picker v-model:value="formData.due_date" style="width: 100%" format="DD/MM/YYYY" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, FilePdfOutlined } from '@ant-design/icons-vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import invoiceService from '@/services/invoiceService'
import purchaseOrderService from '@/services/purchaseOrderService'
import dayjs from 'dayjs'

const loading = ref(false)
const submitting = ref(false)
const createModalVisible = ref(false)
const invoices = ref([])
const eligiblePOs = ref([])
const formRef = ref()
const detailModalVisible = ref(false)
const selectedInvoice = ref(null)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  po_id: undefined,
  amount: 0,
  due_date: null
})

const rules = {
  po_id: [{ required: true, message: 'PO wajib dipilih' }],
  amount: [{ required: true, message: 'Jumlah wajib diisi' }],
  due_date: [{ required: true, message: 'Jatuh tempo wajib diisi' }]
}

const columns = [
  { title: 'Nomor Invoice', dataIndex: 'invoice_number', key: 'invoice_number' },
  { title: 'Nomor PO', key: 'po_number' },
  { title: 'Jumlah', key: 'amount' },
  { title: 'Status', key: 'status', width: 170 },
  { title: 'Jatuh Tempo', key: 'due_date' },
  { title: 'Info Pembayaran', key: 'payment_info' },
  { title: 'Aksi', key: 'actions', width: 120 }
]

const fetchInvoices = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, per_page: pagination.pageSize }
    const response = await invoiceService.getInvoices(params)
    invoices.value = response.data.invoices || response.data.data || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat data invoice')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchEligiblePOs = async () => {
  try {
    const response = await purchaseOrderService.getPurchaseOrders({ status: 'received' })
    eligiblePOs.value = response.data.purchase_orders || []
  } catch (error) {
    console.error('Gagal memuat PO:', error)
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchInvoices()
}

const showCreateModal = () => {
  formData.po_id = undefined
  formData.amount = 0
  formData.due_date = null
  formRef.value?.resetFields()
  createModalVisible.value = true
  fetchEligiblePOs()
}

const handlePOSelect = (poId) => {
  const po = eligiblePOs.value.find(p => p.id === poId)
  if (po) {
    formData.amount = po.total_amount
  }
}

const handleCreate = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    await invoiceService.createInvoice({
      po_id: formData.po_id,
      amount: formData.amount,
      due_date: formData.due_date.format('YYYY-MM-DD')
    })

    message.success('Invoice berhasil dibuat')
    createModalVisible.value = false
    fetchInvoices()
  } catch (error) {
    if (error.errorFields) return
    const errMsg = error.response?.data?.message || ''
    if (errMsg.includes('GRN_NOT_COMPLETED') || errMsg.includes('grn')) {
      message.error('GRN belum selesai untuk PO ini')
    } else {
      message.error(errMsg || 'Gagal membuat invoice')
    }
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const filterPO = (input, option) => {
  return option.children?.[0]?.toLowerCase?.()?.includes(input.toLowerCase()) ?? false
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const viewInvoice = (invoice) => {
  selectedInvoice.value = invoice
  detailModalVisible.value = true
}

const exportInvoiceToPDF = (invoice) => {
  if (!invoice) return
  const doc = new jsPDF()

  // Title
  doc.setFontSize(18)
  doc.setTextColor(0, 0, 0)
  doc.text('INVOICE', 105, 20, { align: 'center' })

  // Invoice Number subtitle
  doc.setFontSize(12)
  doc.text(invoice.invoice_number || '', 105, 28, { align: 'center' })

  // Invoice Info
  doc.setFontSize(10)
  let y = 42
  doc.text(`Nomor Invoice: ${invoice.invoice_number || '-'}`, 14, y); y += 7
  doc.text(`Supplier: ${invoice.supplier?.name || '-'}`, 14, y); y += 7
  doc.text(`Nomor PO: ${invoice.purchase_order?.po_number || '-'}`, 14, y); y += 7
  doc.text(`Tanggal: ${formatDate(invoice.created_at)}`, 14, y); y += 7
  doc.text(`Jatuh Tempo: ${formatDate(invoice.due_date)}`, 14, y); y += 7
  const statusText = invoice.status === 'paid' ? 'Sudah Dibayar' : 'Menunggu Pembayaran'
  doc.text(`Status: ${statusText}`, 14, y); y += 14

  // Amount table
  autoTable(doc, {
    startY: y,
    head: [['Keterangan', 'Jumlah']],
    body: [
      ['Total Tagihan', formatRupiah(invoice.amount)]
    ],
    theme: 'grid',
    headStyles: { fillColor: [248, 44, 23] },
    columnStyles: { 1: { halign: 'right' } }
  })

  y = (doc.lastAutoTable?.finalY || 200) + 10

  // Payment info if paid
  if (invoice.payment) {
    doc.setFontSize(11)
    doc.setTextColor(0, 0, 0)
    doc.text('Informasi Pembayaran', 14, y); y += 8
    doc.setFontSize(10)
    doc.text(`Tanggal Bayar: ${formatDate(invoice.payment.payment_date)}`, 14, y); y += 7
    doc.text(`Metode: ${invoice.payment.payment_method || '-'}`, 14, y); y += 7
    if (invoice.payment.proof_url) {
      doc.text(`Bukti: ${invoice.payment.proof_url}`, 14, y); y += 7
    }
    y += 7
  }

  // Footer
  const finalY = y + 10
  doc.setFontSize(8)
  doc.setTextColor(150)
  doc.text('Dokumen ini digenerate otomatis oleh sistem POSe', 105, finalY, { align: 'center' })

  doc.save(`Invoice-${invoice.invoice_number || 'draft'}.pdf`)
}

onMounted(() => {
  fetchInvoices()
})
</script>

<style scoped>
.supplier-invoice {
  padding: 24px;
}
</style>
