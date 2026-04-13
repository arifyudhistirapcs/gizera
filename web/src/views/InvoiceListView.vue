<template>
  <div class="invoice-list">
    <a-page-header
      title="Invoice"
      sub-title="Kelola invoice dan pembayaran"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Filters -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari nomor invoice..."
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </a-input>
          </a-col>
          <a-col :span="6">
            <a-select
              v-model:value="filterStatus"
              placeholder="Filter Status"
              style="width: 100%"
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <a-select-option value="pending">Pending</a-select-option>
              <a-select-option value="paid">Sudah Dibayar</a-select-option>
            </a-select>
          </a-col>
        </a-row>

        <!-- Table -->
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
            <template v-else-if="column.key === 'supplier'">
              {{ record.supplier?.name || '-' }}
            </template>
            <template v-else-if="column.key === 'po_number'">
              {{ record.purchase_order?.po_number || '-' }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="viewInvoice(record)">
                  Detail
                </a-button>
                <a-button
                  v-if="record.status === 'pending' && can('INVOICE_PAY')"
                  type="primary"
                  size="small"
                  @click="showPayModal(record)"
                >
                  Bayar
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Invoice Detail Modal -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Detail Invoice"
      width="700px"
    >
      <template #footer>
        <a-space>
          <a-button @click="exportInvoiceToPDF(selectedDetailInvoice)">
            <template #icon><FilePdfOutlined /></template>
            Export PDF
          </a-button>
          <a-button @click="detailModalVisible = false">Tutup</a-button>
        </a-space>
      </template>

      <a-descriptions v-if="selectedDetailInvoice" bordered :column="2">
        <a-descriptions-item label="Nomor Invoice" :span="2">
          <strong>{{ selectedDetailInvoice.invoice_number }}</strong>
        </a-descriptions-item>
        <a-descriptions-item label="Supplier">
          {{ selectedDetailInvoice.supplier?.name || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="Nomor PO">
          {{ selectedDetailInvoice.purchase_order?.po_number || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="Yayasan">
          {{ selectedDetailInvoice.yayasan?.nama || selectedDetailInvoice.yayasan?.name || '-' }}
        </a-descriptions-item>
        <a-descriptions-item label="Status">
          <a-tag :color="selectedDetailInvoice.status === 'paid' ? 'green' : 'orange'">
            {{ selectedDetailInvoice.status === 'paid' ? 'Sudah Dibayar' : 'Menunggu Pembayaran' }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="Jumlah" :span="2">
          <strong>{{ formatRupiah(selectedDetailInvoice.amount) }}</strong>
        </a-descriptions-item>
        <a-descriptions-item label="Jatuh Tempo" :span="2">
          {{ formatDate(selectedDetailInvoice.due_date) }}
        </a-descriptions-item>
        <template v-if="selectedDetailInvoice.payment">
          <a-descriptions-item label="Tanggal Bayar">
            {{ formatDate(selectedDetailInvoice.payment.payment_date) }}
          </a-descriptions-item>
          <a-descriptions-item label="Metode Pembayaran">
            {{ selectedDetailInvoice.payment.payment_method || '-' }}
          </a-descriptions-item>
          <a-descriptions-item v-if="selectedDetailInvoice.payment.proof_url" label="Bukti Pembayaran" :span="2">
            <a :href="selectedDetailInvoice.payment.proof_url" target="_blank">Lihat Bukti</a>
          </a-descriptions-item>
        </template>
      </a-descriptions>
    </a-modal>

    <!-- Payment Modal -->
    <a-modal
      v-model:open="payModalVisible"
      title="Pembayaran Invoice"
      :confirm-loading="submitting"
      @ok="handlePay"
      @cancel="payModalVisible = false"
      ok-text="Bayar"
    >
      <a-form layout="vertical">
        <a-form-item label="Invoice">
          <strong>{{ selectedInvoice?.invoice_number }}</strong> — {{ formatRupiah(selectedInvoice?.amount) }}
        </a-form-item>
        <a-form-item label="Tanggal Pembayaran" required>
          <a-date-picker
            v-model:value="paymentDate"
            style="width: 100%"
            format="DD/MM/YYYY"
          />
        </a-form-item>
        <a-form-item label="Metode Pembayaran">
          <a-select v-model:value="paymentMethod" style="width: 100%">
            <a-select-option value="bank_transfer">Transfer Bank</a-select-option>
            <a-select-option value="cash">Tunai</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="Bukti Pembayaran" required>
          <a-upload
            v-model:file-list="proofFileList"
            :before-upload="beforeUpload"
            :max-count="1"
            accept="image/*,.pdf"
          >
            <a-button>
              <template #icon><UploadOutlined /></template>
              Upload Bukti
            </a-button>
          </a-upload>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, UploadOutlined, FilePdfOutlined } from '@ant-design/icons-vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import invoiceService from '@/services/invoiceService'
import { usePermissions } from '@/composables/usePermissions'
import dayjs from 'dayjs'

const { can } = usePermissions()

const loading = ref(false)
const submitting = ref(false)
const invoices = ref([])
const searchText = ref('')
const filterStatus = ref(undefined)
const payModalVisible = ref(false)
const selectedInvoice = ref(null)
const detailModalVisible = ref(false)
const selectedDetailInvoice = ref(null)
const paymentDate = ref(dayjs())
const paymentMethod = ref('bank_transfer')
const proofFileList = ref([])
const proofFile = ref(null)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'Nomor Invoice', dataIndex: 'invoice_number', key: 'invoice_number' },
  { title: 'Supplier', key: 'supplier' },
  { title: 'Nomor PO', key: 'po_number' },
  { title: 'Jumlah', key: 'amount' },
  { title: 'Status', key: 'status', width: 160 },
  { title: 'Jatuh Tempo', key: 'due_date' },
  { title: 'Aksi', key: 'actions', width: 180 }
]

const fetchInvoices = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      per_page: pagination.pageSize,
      status: filterStatus.value || undefined,
      search: searchText.value || undefined
    }
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

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchInvoices()
}

const handleSearch = () => {
  pagination.current = 1
  fetchInvoices()
}

const showPayModal = (invoice) => {
  selectedInvoice.value = invoice
  paymentDate.value = dayjs()
  paymentMethod.value = 'bank_transfer'
  proofFileList.value = []
  proofFile.value = null
  payModalVisible.value = true
}

const beforeUpload = (file) => {
  const isValid = file.type.startsWith('image/') || file.type === 'application/pdf'
  if (!isValid) {
    message.error('Hanya file gambar atau PDF yang diperbolehkan!')
    return false
  }
  const isLt5M = file.size / 1024 / 1024 < 5
  if (!isLt5M) {
    message.error('Ukuran file harus kurang dari 5MB!')
    return false
  }
  proofFile.value = file
  return false
}

const handlePay = async () => {
  if (!paymentDate.value) {
    message.error('Tanggal pembayaran wajib diisi')
    return
  }
  if (!proofFile.value) {
    message.error('Bukti pembayaran wajib diupload')
    return
  }

  submitting.value = true
  try {
    // 1. Upload proof first
    const formData = new FormData()
    formData.append('proof', proofFile.value)
    const uploadRes = await invoiceService.uploadPaymentProof(selectedInvoice.value.id, formData)
    const proofUrl = uploadRes.data.proof_url || uploadRes.data.url || ''

    // 2. Pay invoice
    await invoiceService.payInvoice(selectedInvoice.value.id, {
      payment_date: paymentDate.value.format('YYYY-MM-DD'),
      payment_method: paymentMethod.value,
      proof_url: proofUrl
    })

    message.success('Pembayaran berhasil dicatat')
    payModalVisible.value = false
    fetchInvoices()
  } catch (error) {
    message.error(error.response?.data?.message || 'Gagal memproses pembayaran')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

const getInvoiceStatusText = (status) => {
  return status === 'paid' ? 'Sudah Dibayar' : 'Menunggu Pembayaran'
}

const viewInvoice = (invoice) => {
  selectedDetailInvoice.value = invoice
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
  doc.text(`Yayasan: ${invoice.yayasan?.nama || invoice.yayasan?.name || '-'}`, 14, y); y += 7
  doc.text(`Nomor PO: ${invoice.purchase_order?.po_number || '-'}`, 14, y); y += 7
  doc.text(`Tanggal: ${formatDate(invoice.created_at)}`, 14, y); y += 7
  doc.text(`Jatuh Tempo: ${formatDate(invoice.due_date)}`, 14, y); y += 7
  doc.text(`Status: ${getInvoiceStatusText(invoice.status)}`, 14, y); y += 14

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
  doc.text('Dokumen ini digenerate otomatis oleh sistem Dapur Sehat', 105, finalY, { align: 'center' })

  doc.save(`Invoice-${invoice.invoice_number || 'draft'}.pdf`)
}

onMounted(() => {
  fetchInvoices()
})
</script>

<style scoped>
.invoice-list {
  padding: 24px;
}
</style>
