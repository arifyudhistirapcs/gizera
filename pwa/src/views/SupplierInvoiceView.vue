<template>
  <div class="supplier-invoice-page">
    <!-- Header -->
    <van-nav-bar title="Invoice" fixed>
      <template #right>
        <van-icon name="plus" size="20" @click="showCreateForm = true" />
      </template>
    </van-nav-bar>

    <div class="page-content">
      <!-- Loading -->
      <van-loading v-if="loading" type="spinner" vertical class="page-loading">
        Memuat invoice...
      </van-loading>

      <!-- Empty -->
      <van-empty
        v-else-if="invoices.length === 0"
        image="search"
        description="Belum ada invoice"
      >
        <van-button type="primary" size="small" @click="showCreateForm = true">
          Buat Invoice
        </van-button>
      </van-empty>

      <!-- Invoice List -->
      <van-cell-group v-else inset>
        <van-cell
          v-for="inv in invoices"
          :key="inv.id"
          :title="inv.invoice_number || '-'"
          clickable
          @click="toggleDetail(inv)"
        >
          <template #label>
            <div class="inv-label">
              <span>PO: {{ inv.purchase_order?.po_number || inv.po_number || '-' }}</span>
              <span>Jatuh tempo: {{ formatDate(inv.due_date) }}</span>
            </div>
          </template>
          <template #value>
            <div class="inv-value">
              <span class="inv-amount">{{ formatRupiah(inv.amount) }}</span>
              <van-tag :type="inv.status === 'paid' ? 'success' : 'warning'" size="medium">
                {{ inv.status === 'paid' ? 'Dibayar' : 'Pending' }}
              </van-tag>
            </div>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Payment Info (expanded) -->
      <van-cell-group
        v-if="selectedInvoice?.status === 'paid' && selectedInvoice?.payment"
        inset
        title="Info Pembayaran"
        class="payment-info"
      >
        <van-cell title="Tanggal Bayar" :value="formatDate(selectedInvoice.payment.payment_date)" />
        <van-cell title="Metode" :value="selectedInvoice.payment.payment_method || '-'" />
        <van-cell title="Jumlah" :value="formatRupiah(selectedInvoice.payment.amount)" />
      </van-cell-group>
    </div>

    <!-- Create Invoice Popup -->
    <van-popup v-model:show="showCreateForm" position="bottom" round style="max-height: 80%">
      <div class="create-form-wrapper">
        <h3 class="form-title">Buat Invoice</h3>
        <van-form @submit="onCreateInvoice">
          <!-- PO Picker -->
          <van-field
            v-model="poLabel"
            is-link
            readonly
            label="Purchase Order"
            placeholder="Pilih PO"
            :rules="[{ required: true, message: 'Pilih PO' }]"
            @click="showPOPicker = true"
          />

          <van-popup v-model:show="showPOPicker" position="bottom" round>
            <van-picker
              :columns="poColumns"
              @confirm="onPOConfirm"
              @cancel="showPOPicker = false"
            />
          </van-popup>

          <!-- Amount -->
          <van-field
            v-model="createForm.amount"
            type="number"
            label="Jumlah (Rp)"
            placeholder="Masukkan jumlah"
            :rules="[{ required: true, message: 'Masukkan jumlah' }]"
          />

          <!-- Due Date -->
          <van-field
            v-model="dueDateLabel"
            is-link
            readonly
            label="Jatuh Tempo"
            placeholder="Pilih tanggal"
            :rules="[{ required: true, message: 'Pilih tanggal' }]"
            @click="showDatePicker = true"
          />

          <van-popup v-model:show="showDatePicker" position="bottom" round>
            <van-date-picker
              v-model="dueDate"
              :min-date="new Date()"
              @confirm="onDateConfirm"
              @cancel="showDatePicker = false"
            />
          </van-popup>

          <div class="form-actions">
            <van-button
              type="primary"
              block
              native-type="submit"
              :loading="submitting"
              :disabled="submitting"
            >
              Buat Invoice
            </van-button>
          </div>
        </van-form>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'

const loading = ref(false)
const submitting = ref(false)
const invoices = ref([])
const selectedInvoice = ref(null)
const showCreateForm = ref(false)
const showPOPicker = ref(false)
const showDatePicker = ref(false)
const receivedPOs = ref([])
const poLabel = ref('')
const dueDateLabel = ref('')
const dueDate = ref([])

const createForm = ref({
  po_id: null,
  amount: '',
  due_date: ''
})

const poColumns = computed(() =>
  receivedPOs.value.map(po => ({
    text: `${po.po_number} - ${formatRupiah(po.total_amount)}`,
    value: po.id
  }))
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
    month: 'short',
    year: 'numeric'
  })
}

function toggleDetail(inv) {
  selectedInvoice.value = selectedInvoice.value?.id === inv.id ? null : inv
}

function onPOConfirm({ selectedOptions }) {
  const selected = selectedOptions[0]
  if (selected) {
    createForm.value.po_id = selected.value
    poLabel.value = selected.text
    // Auto-fill amount from PO
    const po = receivedPOs.value.find(p => p.id === selected.value)
    if (po) {
      createForm.value.amount = String(po.total_amount ?? '')
    }
  }
  showPOPicker.value = false
}

function onDateConfirm({ selectedValues }) {
  const [year, month, day] = selectedValues
  createForm.value.due_date = `${year}-${month}-${day}`
  dueDateLabel.value = `${day}/${month}/${year}`
  showDatePicker.value = false
}

async function loadInvoices() {
  loading.value = true
  try {
    const res = await supplierPortalService.getInvoices()
    invoices.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal memuat invoice')
  } finally {
    loading.value = false
  }
}

async function loadReceivedPOs() {
  try {
    const res = await supplierPortalService.getPurchaseOrders({ status: 'received' })
    receivedPOs.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    // Silent fail — PO list for picker
  }
}

async function onCreateInvoice() {
  if (!createForm.value.po_id) {
    showToast('Pilih PO terlebih dahulu')
    return
  }

  submitting.value = true
  try {
    await supplierPortalService.createInvoice({
      po_id: Number(createForm.value.po_id),
      amount: Number(createForm.value.amount),
      due_date: createForm.value.due_date
    })
    showToast('Invoice berhasil dibuat')
    showCreateForm.value = false
    createForm.value = { po_id: null, amount: '', due_date: '' }
    poLabel.value = ''
    dueDateLabel.value = ''
    await loadInvoices()
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal membuat invoice')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadInvoices()
  loadReceivedPOs()
})
</script>

<style scoped>
.supplier-invoice-page {
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

.inv-label {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
  color: #969799;
}

.inv-value {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.inv-amount {
  font-size: 14px;
  font-weight: 600;
  color: #303030;
}

.payment-info {
  margin-top: 16px;
}

.create-form-wrapper {
  padding: 20px 16px;
}

.form-title {
  font-size: 18px;
  font-weight: 700;
  color: #303030;
  margin: 0 0 16px;
}

.form-actions {
  padding: 16px 0;
}
</style>
