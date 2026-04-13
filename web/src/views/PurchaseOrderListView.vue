<template>
  <div class="purchase-order-list">
    <a-page-header
      title="Purchase Order (PO)"
      sub-title="Kelola pemesanan barang ke supplier"
    >
      <template #extra>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Buat PO Baru
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Search and Filter -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari nomor PO atau supplier..."
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
              placeholder="Status"
              style="width: 100%"
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <a-select-option value="pending">Pending</a-select-option>
              <a-select-option value="revision_by_supplier">Revisi oleh Supplier</a-select-option>
              <a-select-option value="approved">Disetujui</a-select-option>
              <a-select-option value="shipping">Sedang Dikirim</a-select-option>
              <a-select-option value="received">Diterima</a-select-option>
              <a-select-option value="cancelled">Dibatalkan</a-select-option>
            </a-select>
          </a-col>
        </a-row>

        <!-- Table -->
        <a-table
          :columns="columns"
          :data-source="purchaseOrders"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusText(record.status) }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'rab_number'">
              {{ record.rab?.rab_number || '-' }}
            </template>
            <template v-else-if="column.key === 'target_sppg'">
              {{ record.target_sppg?.nama || record.target_sppg?.name || '-' }}
            </template>
            <template v-else-if="column.key === 'total_amount'">
              {{ formatCurrency(record.total_amount) }}
            </template>
            <template v-else-if="column.key === 'order_date'">
              {{ formatDate(record.order_date) }}
            </template>
            <template v-else-if="column.key === 'expected_delivery'">
              {{ formatDate(record.expected_delivery) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="viewPO(record)">
                  Detail
                </a-button>
                <a-button
                  v-if="record.status === 'pending' && canApprove"
                  type="link"
                  size="small"
                  @click="approvePO(record)"
                >
                  Setujui
                </a-button>
                <a-button
                  v-if="record.status === 'pending' && !isKepalaYayasan"
                  type="link"
                  size="small"
                  @click="editPO(record)"
                >
                  Edit
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- ============================================ -->
    <!-- BATCH CREATE MODAL (kepala_yayasan only)     -->
    <!-- ============================================ -->
    <a-modal
      v-if="isKepalaYayasan && !editingPO"
      v-model:open="modalVisible"
      :title="'Buat PO Batch dari RAB'"
      :confirm-loading="submitting"
      :ok-text="batchCurrentStep === 0 ? 'Lanjut' : 'Buat PO Batch'"
      :ok-button-props="{ disabled: batchCurrentStep === 0 ? !batchSelectedRABId : !batchExpectedDelivery || batchSupplierGroups.length === 0 }"
      @ok="handleBatchStepNext"
      @cancel="handleBatchCancel"
      width="960px"
    >
      <a-steps :current="batchCurrentStep" style="margin-bottom: 24px" size="small">
        <a-step title="Pilih RAB" />
        <a-step title="Review & Tanggal Kirim" />
      </a-steps>

      <!-- Step 0: Select RAB -->
      <div v-if="batchCurrentStep === 0">
        <a-form layout="vertical">
          <a-form-item label="RAB (Status: Approved Yayasan)" required>
            <a-select
              v-model:value="batchSelectedRABId"
              placeholder="Pilih RAB yang sudah disetujui yayasan"
              show-search
              allow-clear
              :filter-option="filterRAB"
              :loading="batchLoadingRABs"
              @change="handleBatchRABChange"
              style="width: 100%"
              size="large"
            >
              <a-select-option v-for="r in approvedRABs" :key="r.id" :value="r.id">
                {{ r.rab_number }} — {{ r.menu_plan ? formatDate(r.menu_plan.week_start) + ' s/d ' + formatDate(r.menu_plan.week_end) : '' }}
              </a-select-option>
            </a-select>
          </a-form-item>

          <a-form-item label="SPPG Tujuan">
            <a-select
              :value="batchSPPGId"
              disabled
              placeholder="Otomatis dari RAB"
              style="width: 100%"
              size="large"
            >
              <a-select-option v-if="batchSPPGId" :value="batchSPPGId">
                {{ batchSPPGName }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-form>

        <!-- Loading state -->
        <div v-if="batchLoadingDetail" style="text-align: center; padding: 40px 0">
          <a-spin size="large" />
          <p style="margin-top: 12px; color: rgba(0,0,0,0.45)">Memuat detail RAB...</p>
        </div>

        <!-- Preview grouped items after RAB selected -->
        <div v-if="batchSelectedRABId && !batchLoadingDetail && batchAllItems.length > 0">
          <a-divider>Preview Item RAB</a-divider>
          <a-row :gutter="16" style="margin-bottom: 16px">
            <a-col :span="8">
              <a-statistic
                title="Total Supplier"
                :value="batchSupplierGroups.length"
                suffix="supplier"
                :value-style="{ color: '#1890ff' }"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="Total Item (ada supplier)"
                :value="batchTotalItemsWithSupplier"
                suffix="item"
                :value-style="{ color: '#52c41a' }"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="Tanpa Supplier"
                :value="batchItemsWithoutSupplier.length"
                suffix="item"
                :value-style="{ color: batchItemsWithoutSupplier.length > 0 ? '#faad14' : '#52c41a' }"
              />
            </a-col>
          </a-row>
        </div>
      </div>

      <!-- Step 1: Review grouped items + delivery date -->
      <div v-if="batchCurrentStep === 1">
        <a-form layout="vertical">
          <a-form-item label="Tanggal Pengiriman Diharapkan" required>
            <a-date-picker
              v-model:value="batchExpectedDelivery"
              style="width: 100%"
              format="DD/MM/YYYY"
              placeholder="Pilih tanggal pengiriman"
              :disabled-date="disabledDate"
              size="large"
            />
          </a-form-item>
        </a-form>

        <a-divider>Ringkasan PO yang Akan Dibuat</a-divider>

        <!-- Summary stats -->
        <a-row :gutter="16" style="margin-bottom: 16px">
          <a-col :span="8">
            <a-statistic
              title="Jumlah PO"
              :value="batchSupplierGroups.length"
              :value-style="{ color: '#1890ff' }"
            />
          </a-col>
          <a-col :span="8">
            <a-statistic
              title="Total Item"
              :value="batchTotalItemsWithSupplier"
              :value-style="{ color: '#52c41a' }"
            />
          </a-col>
          <a-col :span="8">
            <a-statistic
              title="Total Nilai"
              :value="formatCurrency(batchGrandTotal)"
              :value-style="{ color: '#1890ff', fontSize: '20px' }"
            />
          </a-col>
        </a-row>

        <!-- Supplier groups -->
        <a-collapse v-model:activeKey="batchActiveKeys" style="margin-bottom: 16px">
          <a-collapse-panel
            v-for="group in batchSupplierGroups"
            :key="group.supplierName"
            :header="`📦 ${group.supplierName} (${group.items.length} item) — Total: ${formatCurrency(group.total)}`"
          >
            <a-table
              :columns="batchItemColumns"
              :data-source="group.items"
              :pagination="false"
              size="small"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'name'">
                  {{ record.ingredient?.name || record.ingredient_name || '-' }}
                </template>
                <template v-else-if="column.key === 'quantity'">
                  {{ formatNumber(record.quantity) }} {{ record.unit }}
                </template>
                <template v-else-if="column.key === 'unit_price'">
                  {{ formatCurrency(record.unit_price || record.estimated_price) }}
                </template>
                <template v-else-if="column.key === 'subtotal'">
                  {{ formatCurrency((record.quantity || 0) * (record.unit_price || record.estimated_price || 0)) }}
                </template>
              </template>
            </a-table>
          </a-collapse-panel>
        </a-collapse>

        <!-- Items without supplier -->
        <a-alert
          v-if="batchItemsWithoutSupplier.length > 0"
          type="warning"
          show-icon
          style="margin-bottom: 16px"
        >
          <template #message>
            ⚠️ Tanpa Supplier ({{ batchItemsWithoutSupplier.length }} item) — Tidak akan dibuat PO
          </template>
          <template #description>
            <ul style="margin: 8px 0 0 0; padding-left: 20px">
              <li v-for="item in batchItemsWithoutSupplier" :key="item.id">
                {{ item.ingredient?.name || item.ingredient_name || '-' }}: {{ formatNumber(item.quantity) }} {{ item.unit }}
              </li>
            </ul>
          </template>
        </a-alert>
      </div>
    </a-modal>

    <!-- ============================================ -->
    <!-- STANDARD CREATE/EDIT MODAL (non-yayasan)     -->
    <!-- ============================================ -->
    <a-modal
      v-if="!isKepalaYayasan || editingPO"
      v-model:open="modalVisible"
      :title="editingPO ? 'Edit Purchase Order' : 'Buat Purchase Order Baru'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="900px"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="RAB" name="rab_id">
              <a-select
                v-model:value="formData.rab_id"
                placeholder="Pilih RAB (approved yayasan)"
                show-search
                allow-clear
                :filter-option="filterRAB"
                @change="handleRABChange"
              >
                <a-select-option v-for="r in approvedRABs" :key="r.id" :value="r.id">
                  {{ r.rab_number }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="SPPG Tujuan" name="target_sppg_id">
              <a-select
                v-model:value="formData.target_sppg_id"
                placeholder="Pilih SPPG tujuan"
                show-search
                allow-clear
                :filter-option="filterSPPG"
              >
                <a-select-option v-for="s in sppgList" :key="s.id" :value="s.id">
                  {{ s.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Supplier" name="supplier_id">
              <a-select
                v-model:value="formData.supplier_id"
                placeholder="Pilih supplier"
                show-search
                :filter-option="filterSupplier"
                @change="handleSupplierChange"
              >
                <a-select-option
                  v-for="supplier in activeSuppliers"
                  :key="supplier.id"
                  :value="supplier.id"
                >
                  {{ supplier.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Tanggal Pengiriman Diharapkan" name="expected_delivery">
              <a-date-picker
                v-model:value="formData.expected_delivery"
                style="width: 100%"
                format="DD/MM/YYYY"
                :disabled-date="disabledDate"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-divider>Item Pesanan</a-divider>

        <a-form-item label="Item" name="items">
          <a-table
            :columns="itemColumns"
            :data-source="formData.items"
            :pagination="false"
            size="small"
          >
            <template #bodyCell="{ column, record, index }">
              <template v-if="column.key === 'ingredient_id'">
                <a-select
                  v-model:value="record.ingredient_id"
                  placeholder="Pilih bahan"
                  show-search
                  style="width: 100%"
                  :filter-option="filterIngredient"
                  @change="calculateSubtotal(index)"
                >
                  <a-select-option
                    v-for="ingredient in ingredients"
                    :key="ingredient.id"
                    :value="ingredient.id"
                  >
                    {{ ingredient.name }} ({{ ingredient.unit }})
                  </a-select-option>
                </a-select>
              </template>
              <template v-else-if="column.key === 'quantity'">
                <a-input-number
                  v-model:value="record.quantity"
                  :min="0.01"
                  :step="0.1"
                  style="width: 100%"
                  @change="calculateSubtotal(index)"
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
                  @change="calculateSubtotal(index)"
                />
              </template>
              <template v-else-if="column.key === 'subtotal'">
                {{ formatCurrency(record.subtotal || 0) }}
              </template>
              <template v-else-if="column.key === 'actions'">
                <a-button type="link" size="small" danger @click="removeItem(index)">
                  Hapus
                </a-button>
              </template>
            </template>
          </a-table>

          <a-button type="dashed" block @click="addItem" style="margin-top: 8px">
            <template #icon><PlusOutlined /></template>
            Tambah Item
          </a-button>
        </a-form-item>

        <a-row justify="end">
          <a-col>
            <div style="text-align: right">
              <div style="color: rgba(0,0,0,0.45); font-size: 14px">Total</div>
              <div style="font-size: 24px; font-weight: 500">{{ formatCurrency(totalAmount) }}</div>
            </div>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>

    <!-- Detail Modal -->
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

          <!-- Pending: kepala_sppg can approve -->
          <a-button v-if="selectedPO?.status === 'pending' && canApprove" type="primary" @click="approvePO(selectedPO)">
            Setujui PO
          </a-button>

          <!-- revision_by_supplier: kepala_yayasan actions -->
          <template v-if="selectedPO?.status === 'revision_by_supplier' && isKepalaYayasan && !showReviseForm">
            <a-button
              type="primary"
              :loading="revisionActionLoading"
              @click="handleAcceptRevision"
            >
              <template #icon><CheckOutlined /></template>
              Terima Perubahan
            </a-button>
            <a-button
              type="primary"
              style="background-color: #fa8c16; border-color: #fa8c16"
              @click="openReviseForm"
            >
              <template #icon><EditOutlined /></template>
              Revisi Ulang
            </a-button>
          </template>

          <!-- Revise form submit -->
          <template v-if="showReviseForm">
            <a-button @click="cancelReviseForm">Batal</a-button>
            <a-button
              type="primary"
              :loading="revisionActionLoading"
              :disabled="reviseItems.length === 0"
              @click="handleSubmitRevise"
            >
              Kirim Revisi ke Supplier
            </a-button>
          </template>

          <!-- confirmed: no longer needed, supplier confirm goes straight to approved -->

          <a-button v-if="!showReviseForm" @click="closeDetailModal">Tutup</a-button>
        </a-space>
      </template>

      <a-descriptions v-if="selectedPO" bordered :column="2">
        <a-descriptions-item label="Nomor PO" :span="2">
          <strong>{{ selectedPO.po_number }}</strong>
        </a-descriptions-item>
        <a-descriptions-item label="Supplier">
          {{ selectedPO.supplier?.name }}
        </a-descriptions-item>
        <a-descriptions-item label="Status">
          <a-tag :color="getStatusColor(selectedPO.status)">
            {{ getStatusText(selectedPO.status) }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="Tanggal Order">
          {{ formatDate(selectedPO.order_date) }}
        </a-descriptions-item>
        <a-descriptions-item label="Tanggal Pengiriman">
          {{ formatDate(selectedPO.expected_delivery) }}
        </a-descriptions-item>
        <a-descriptions-item label="Total" :span="2">
          <strong>{{ formatCurrency(selectedPO.total_amount) }}</strong>
        </a-descriptions-item>
        <a-descriptions-item v-if="selectedPO.approved_by" label="Disetujui Oleh" :span="2">
          {{ selectedPO.approver?.name }} pada {{ formatDate(selectedPO.approved_at) }}
        </a-descriptions-item>
      </a-descriptions>

      <!-- Alert for revision_by_supplier status -->
      <a-alert
        v-if="selectedPO?.status === 'revision_by_supplier' && !showReviseForm"
        type="warning"
        show-icon
        style="margin-top: 16px"
      >
        <template #message>Supplier mengajukan perubahan pada PO ini</template>
        <template #description>
          <div v-if="selectedPO.supplier_revision_notes">
            <strong>Catatan dari supplier:</strong>
            <p style="margin: 8px 0 0 0; white-space: pre-wrap">{{ selectedPO.supplier_revision_notes }}</p>
          </div>
        </template>
      </a-alert>

      <!-- Alert for approved status -->
      <a-alert
        v-if="selectedPO?.status === 'approved' && !showReviseForm"
        type="success"
        show-icon
        style="margin-top: 16px"
        message="PO telah disetujui"
        description="PO telah disetujui oleh supplier dan siap untuk pengiriman."
      />

      <!-- Alert for shipping status -->
      <a-alert
        v-if="selectedPO?.status === 'shipping' && !showReviseForm"
        type="info"
        show-icon
        style="margin-top: 16px"
        message="Barang sedang dikirim"
        description="Supplier sedang mengirim barang ke SPPG tujuan. Menunggu penerimaan (GRN) dari pihak SPPG."
      />

      <!-- Revise form (kepala_yayasan edits items to send back to supplier) -->
      <template v-if="showReviseForm">
        <a-divider>Revisi Item PO</a-divider>

        <a-alert
          type="info"
          show-icon
          style="margin-bottom: 16px"
          message="Ubah item PO lalu kirim kembali ke supplier. Status PO akan kembali ke Pending."
        />

        <a-table
          :columns="reviseItemColumns"
          :data-source="reviseItems"
          :pagination="false"
          size="small"
          row-key="ingredient_id"
        >
          <template #bodyCell="{ column, record, index }">
            <template v-if="column.key === 'ingredient_name'">
              {{ record.ingredient?.name || record.ingredient_name || '-' }}
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
              />
            </template>
            <template v-else-if="column.key === 'subtotal'">
              {{ formatCurrency((record.quantity || 0) * (record.unit_price || 0)) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-popconfirm
                title="Hapus item ini?"
                ok-text="Ya"
                cancel-text="Batal"
                @confirm="removeReviseItem(index)"
              >
                <a-button type="link" size="small" danger :disabled="reviseItems.length <= 1">
                  Hapus
                </a-button>
              </a-popconfirm>
            </template>
          </template>
        </a-table>

        <div style="text-align: right; margin: 12px 0">
          <strong>Total: {{ formatCurrency(reviseTotal) }}</strong>
        </div>
      </template>

      <!-- Normal item table (hidden when revise form is open) -->
      <template v-if="!showReviseForm">
        <a-divider>Item Pesanan</a-divider>

        <a-table
          :columns="detailItemColumns"
          :data-source="selectedPO?.po_items || []"
          :pagination="false"
          size="small"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'unit_price'">
              {{ formatCurrency(record.unit_price) }}
            </template>
            <template v-else-if="column.key === 'subtotal'">
              {{ formatCurrency(record.subtotal) }}
            </template>
          </template>
        </a-table>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, SearchOutlined, FilePdfOutlined, CheckOutlined, EditOutlined } from '@ant-design/icons-vue'
import jsPDF from 'jspdf'
import autoTable from 'jspdf-autotable'
import purchaseOrderService from '@/services/purchaseOrderService'
import supplierService from '@/services/supplierService'
import recipeService from '@/services/recipeService'
import rabService from '@/services/rabService'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'
import dayjs from 'dayjs'

const authStore = useAuthStore()
const canApprove = computed(() => authStore.user?.role === 'kepala_sppg')
const isKepalaYayasan = computed(() => authStore.user?.role === 'kepala_yayasan')

const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const detailModalVisible = ref(false)
const editingPO = ref(null)
const selectedPO = ref(null)

// Revision handling state (kepala_yayasan)
const revisionActionLoading = ref(false)
const showReviseForm = ref(false)
const reviseItems = ref([])

const reviseTotal = computed(() => {
  return reviseItems.value.reduce((sum, item) => {
    return sum + (item.quantity || 0) * (item.unit_price || 0)
  }, 0)
})

const reviseItemColumns = [
  { title: 'Bahan', key: 'ingredient_name', width: 200 },
  { title: 'Satuan', key: 'unit', width: 80 },
  { title: 'Jumlah', key: 'quantity', width: 120 },
  { title: 'Harga Satuan', key: 'unit_price', width: 160 },
  { title: 'Subtotal', key: 'subtotal', width: 160 },
  { title: 'Aksi', key: 'actions', width: 80 }
]
const purchaseOrders = ref([])
const activeSuppliers = ref([])
const ingredients = ref([])
const searchText = ref('')
const filterStatus = ref(undefined)
const formRef = ref()
const totalAmount = ref(0)
const approvedRABs = ref([])
const sppgList = ref([])
const rabItems = ref([])

// --- Batch creation state (kepala_yayasan) ---
const batchCurrentStep = ref(0)
const batchSelectedRABId = ref(undefined)
const batchLoadingRABs = ref(false)
const batchLoadingDetail = ref(false)
const batchExpectedDelivery = ref(null)
const batchAllItems = ref([])
const batchSPPGId = ref(undefined)
const batchSPPGName = ref('')
const batchActiveKeys = ref([])

const batchSupplierGroups = computed(() => {
  const groups = {}
  for (const item of batchAllItems.value) {
    const supplierName = item.recommended_supplier?.name
    if (!supplierName) continue
    if (!groups[supplierName]) {
      groups[supplierName] = { supplierName, items: [], total: 0 }
    }
    groups[supplierName].items.push(item)
    const price = item.unit_price || item.estimated_price || 0
    groups[supplierName].total += (item.quantity || 0) * price
  }
  return Object.values(groups)
})

const batchItemsWithoutSupplier = computed(() => {
  return batchAllItems.value.filter(item => !item.recommended_supplier?.name)
})

const batchTotalItemsWithSupplier = computed(() => {
  return batchSupplierGroups.value.reduce((sum, g) => sum + g.items.length, 0)
})

const batchGrandTotal = computed(() => {
  return batchSupplierGroups.value.reduce((sum, g) => sum + g.total, 0)
})

const batchItemColumns = [
  { title: 'Bahan', key: 'name' },
  { title: 'Jumlah', key: 'quantity', width: 160 },
  { title: 'Harga Satuan', key: 'unit_price', width: 160 },
  { title: 'Subtotal', key: 'subtotal', width: 180 }
]

const handleBatchRABChange = async (rabId) => {
  batchAllItems.value = []
  batchSPPGId.value = undefined
  batchSPPGName.value = ''
  if (!rabId) return

  batchLoadingDetail.value = true
  try {
    const response = await rabService.getRABDetail(rabId)
    const rab = response.data.rab || response.data.data || response.data

    // Auto-populate SPPG from RAB
    batchSPPGId.value = rab.sppg_id
    batchSPPGName.value = rab.sppg?.nama || rab.sppg?.name || `SPPG #${rab.sppg_id}`

    // Filter items that don't have po_id yet (status "pending")
    const items = (rab.items || []).filter(item => !item.po_id && item.status !== 'ordered')
    batchAllItems.value = items
  } catch (error) {
    message.error('Gagal memuat detail RAB')
    console.error(error)
  } finally {
    batchLoadingDetail.value = false
  }
}

const handleBatchStepNext = async () => {
  if (batchCurrentStep.value === 0) {
    if (!batchSelectedRABId.value) {
      message.warning('Pilih RAB terlebih dahulu')
      return
    }
    if (batchSupplierGroups.value.length === 0) {
      message.warning('Tidak ada item dengan supplier yang bisa dibuat PO')
      return
    }
    // Open all collapse panels by default
    batchActiveKeys.value = batchSupplierGroups.value.map(g => g.supplierName)
    batchCurrentStep.value = 1
    return
  }

  // Step 1: Submit batch
  if (!batchExpectedDelivery.value) {
    message.warning('Pilih tanggal pengiriman terlebih dahulu')
    return
  }

  submitting.value = true
  try {
    const payload = {
      rab_id: batchSelectedRABId.value,
      expected_delivery: batchExpectedDelivery.value.format('YYYY-MM-DD')
    }
    const response = await purchaseOrderService.createBatchFromRAB(payload)
    const count = response.data.created_count || response.data.total || batchSupplierGroups.value.length
    message.success(`Berhasil membuat ${count} Purchase Order dari RAB`)
    modalVisible.value = false
    resetBatchForm()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal membuat PO batch'
    message.error(errMsg)
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleBatchCancel = () => {
  modalVisible.value = false
  resetBatchForm()
}

const resetBatchForm = () => {
  batchCurrentStep.value = 0
  batchSelectedRABId.value = undefined
  batchExpectedDelivery.value = null
  batchAllItems.value = []
  batchSPPGId.value = undefined
  batchSPPGName.value = ''
  batchActiveKeys.value = []
}

// --- End batch creation state ---

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = ref({
  supplier_id: undefined,
  expected_delivery: null,
  rab_id: undefined,
  target_sppg_id: undefined,
  items: []
})

const updateTotal = () => {
  const total = formData.value.items.reduce((sum, item) => {
    return sum + (parseFloat(item.subtotal) || 0)
  }, 0)
  totalAmount.value = total
}

// Watch for items changes and update total
watch(() => formData.value.items, () => {
  updateTotal()
}, { deep: true })

const rules = {
  supplier_id: [{ required: true, message: 'Supplier wajib dipilih' }],
  expected_delivery: [{ required: true, message: 'Tanggal pengiriman wajib diisi' }],
  items: [
    {
      required: true,
      validator: (rule, value) => {
        if (!value || value.length === 0) {
          return Promise.reject('Minimal satu item harus ditambahkan')
        }
        return Promise.resolve()
      }
    }
  ]
}

const columns = [
  { title: 'Nomor PO', dataIndex: 'po_number', key: 'po_number' },
  { title: 'RAB', key: 'rab_number' },
  { title: 'Supplier', dataIndex: ['supplier', 'name'], key: 'supplier_name' },
  { title: 'SPPG Tujuan', key: 'target_sppg' },
  { title: 'Tanggal Order', key: 'order_date' },
  { title: 'Tanggal Pengiriman', key: 'expected_delivery' },
  { title: 'Total', key: 'total_amount' },
  { title: 'Status', key: 'status', width: 120 },
  { title: 'Aksi', key: 'actions', width: 200 }
]

const itemColumns = [
  { title: 'Bahan', key: 'ingredient_id', width: 250 },
  { title: 'Jumlah', key: 'quantity', width: 120 },
  { title: 'Harga Satuan', key: 'unit_price', width: 150 },
  { title: 'Subtotal', key: 'subtotal', width: 150 },
  { title: 'Aksi', key: 'actions', width: 80 }
]

const detailItemColumns = [
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
      search: searchText.value || undefined,
      status: filterStatus.value
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

const fetchSuppliers = async () => {
  try {
    const response = await supplierService.getSuppliers({ is_active: 'active' })
    activeSuppliers.value = response.data.suppliers || []
  } catch (error) {
    console.error('Gagal memuat data supplier:', error)
  }
}

const fetchIngredients = async () => {
  try {
    const response = await recipeService.getIngredients()
    ingredients.value = response.data.data || []
  } catch (error) {
    console.error('Gagal memuat data bahan:', error)
  }
}

const fetchApprovedRABs = async () => {
  batchLoadingRABs.value = true
  try {
    const response = await rabService.getRABList({ status: 'approved_yayasan' })
    approvedRABs.value = response.data.rabs || response.data.data || []
  } catch (error) {
    console.error('Gagal memuat data RAB:', error)
  } finally {
    batchLoadingRABs.value = false
  }
}

const fetchSPPGList = async () => {
  try {
    const response = await api.get('/sppg')
    sppgList.value = response.data.sppgs || response.data.data || []
  } catch (error) {
    console.error('Gagal memuat data SPPG:', error)
  }
}

const handleRABChange = async (rabId) => {
  if (!rabId) {
    rabItems.value = []
    return
  }
  try {
    const response = await rabService.getRABDetail(rabId)
    const rab = response.data.rab || response.data.data || response.data
    rabItems.value = (rab.items || []).filter(item => !item.po_id)
  } catch (error) {
    console.error('Gagal memuat item RAB:', error)
  }
}

const filterRAB = (input, option) => {
  const r = approvedRABs.value.find(rab => rab.id === option.value)
  return r?.rab_number?.toLowerCase().includes(input.toLowerCase()) ?? false
}

const filterSPPG = (input, option) => {
  const s = sppgList.value.find(sppg => sppg.id === option.value)
  return s?.name?.toLowerCase().includes(input.toLowerCase()) ?? false
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

const showCreateModal = () => {
  editingPO.value = null
  if (isKepalaYayasan.value) {
    resetBatchForm()
  } else {
    resetForm()
  }
  modalVisible.value = true
}

const editPO = (po) => {
  editingPO.value = po
  formData.value = {
    supplier_id: po.supplier_id,
    expected_delivery: dayjs(po.expected_delivery),
    rab_id: undefined,
    target_sppg_id: undefined,
    items: po.po_items.map(item => ({
      ingredient_id: item.ingredient_id,
      quantity: item.quantity,
      unit_price: item.unit_price,
      subtotal: item.subtotal
    }))
  }
  modalVisible.value = true
}

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

const closeDetailModal = () => {
  detailModalVisible.value = false
  showReviseForm.value = false
  reviseItems.value = []
}

const handleAcceptRevision = async () => {
  revisionActionLoading.value = true
  try {
    await purchaseOrderService.acceptRevision(selectedPO.value.id)
    message.success('Perubahan supplier diterima, PO dikonfirmasi')
    closeDetailModal()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal menerima perubahan'
    message.error(errMsg)
    console.error(error)
  } finally {
    revisionActionLoading.value = false
  }
}

const openReviseForm = () => {
  reviseItems.value = (selectedPO.value.po_items || []).map(item => ({
    id: item.id,
    ingredient_id: item.ingredient_id,
    ingredient: item.ingredient ? { ...item.ingredient } : null,
    ingredient_name: item.ingredient?.name || '-',
    quantity: item.quantity,
    unit_price: item.unit_price
  }))
  showReviseForm.value = true
}

const cancelReviseForm = () => {
  showReviseForm.value = false
  reviseItems.value = []
}

const removeReviseItem = (index) => {
  reviseItems.value.splice(index, 1)
}

const handleSubmitRevise = async () => {
  if (reviseItems.value.length === 0) {
    message.warning('Minimal satu item harus ada')
    return
  }

  revisionActionLoading.value = true
  try {
    const payload = {
      items: reviseItems.value.map(item => ({
        ingredient_id: item.ingredient_id,
        quantity: item.quantity,
        unit_price: item.unit_price
      }))
    }
    await purchaseOrderService.revisePO(selectedPO.value.id, payload)
    message.success('PO direvisi dan dikirim kembali ke supplier')
    closeDetailModal()
    fetchPurchaseOrders()
  } catch (error) {
    const errMsg = error.response?.data?.error || error.response?.data?.message || 'Gagal merevisi PO'
    message.error(errMsg)
    console.error(error)
  } finally {
    revisionActionLoading.value = false
  }
}

const approvePO = async (po) => {
  try {
    await purchaseOrderService.approvePurchaseOrder(po.id)
    message.success('Purchase Order berhasil disetujui')
    detailModalVisible.value = false
    fetchPurchaseOrders()
  } catch (error) {
    message.error('Gagal menyetujui PO')
    console.error(error)
  }
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
  doc.text(`Yayasan: ${po.yayasan?.nama || '-'}`, 14, 56)
  doc.text(`SPPG Tujuan: ${po.target_sppg?.nama || po.target_sppg?.name || '-'}`, 14, 63)
  doc.text(`Tanggal Order: ${formatDate(po.order_date)}`, 14, 70)
  doc.text(`Tanggal Pengiriman: ${formatDate(po.expected_delivery)}`, 14, 77)
  doc.text(`Status: ${getStatusText(po.status)}`, 14, 84)

  // Items table
  const items = (po.po_items || []).map(item => [
    item.ingredient?.name || '-',
    item.quantity?.toString() || '0',
    item.ingredient?.unit || '-',
    formatCurrency(item.unit_price),
    formatCurrency(item.subtotal)
  ])

  autoTable(doc, {
    startY: 92,
    head: [['Bahan', 'Jumlah', 'Satuan', 'Harga Satuan', 'Subtotal']],
    body: items,
    theme: 'grid',
    headStyles: { fillColor: [248, 44, 23] },
    foot: [['', '', '', 'TOTAL', formatCurrency(po.total_amount)]],
    footStyles: { fillColor: [240, 240, 240], textColor: [0, 0, 0], fontStyle: 'bold' }
  })

  // Footer
  const finalY = (doc.lastAutoTable?.finalY || 200) + 20
  doc.setFontSize(8)
  doc.setTextColor(150)
  doc.text('Dokumen ini digenerate otomatis oleh sistem Dapur Sehat', 105, finalY, { align: 'center' })

  doc.save(`PO-${po.po_number || 'draft'}.pdf`)
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()

    if (formData.value.items.length === 0) {
      message.error('Minimal satu item harus ditambahkan')
      return
    }

    submitting.value = true

    const payload = {
      supplier_id: formData.value.supplier_id,
      expected_delivery: formData.value.expected_delivery.format('YYYY-MM-DD'),
      rab_id: formData.value.rab_id || undefined,
      target_sppg_id: formData.value.target_sppg_id || undefined,
      items: formData.value.items.map(item => ({
        ingredient_id: item.ingredient_id,
        quantity: item.quantity,
        unit_price: item.unit_price
      }))
    }

    if (editingPO.value) {
      await purchaseOrderService.updatePurchaseOrder(editingPO.value.id, payload)
      message.success('Purchase Order berhasil diperbarui')
    } else {
      await purchaseOrderService.createPurchaseOrder(payload)
      message.success('Purchase Order berhasil dibuat')
    }

    modalVisible.value = false
    fetchPurchaseOrders()
  } catch (error) {
    if (error.errorFields) return
    message.error('Gagal menyimpan purchase order')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleCancel = () => {
  modalVisible.value = false
  resetForm()
}

const resetForm = () => {
  formData.value = {
    supplier_id: undefined,
    expected_delivery: null,
    rab_id: undefined,
    target_sppg_id: undefined,
    items: []
  }
  rabItems.value = []
  formRef.value?.resetFields()
}

const addItem = () => {
  formData.value.items.push({
    ingredient_id: undefined,
    quantity: 1,
    unit_price: 0,
    subtotal: 0
  })
  updateTotal()
}

const removeItem = (index) => {
  formData.value.items.splice(index, 1)
  updateTotal()
}

const calculateSubtotal = (index) => {
  const item = formData.value.items[index]
  const qty = parseFloat(item.quantity) || 0
  const price = parseCurrency(item.unit_price)
  item.subtotal = qty * price
  updateTotal()
}

const handleSupplierChange = () => {
  // Could fetch supplier-specific pricing here
}

const filterSupplier = (input, option) => {
  const supplier = activeSuppliers.value.find(s => s.id === option.value)
  return supplier?.name.toLowerCase().includes(input.toLowerCase())
}

const filterIngredient = (input, option) => {
  const ingredient = ingredients.value.find(i => i.id === option.value)
  if (!ingredient) return false
  return ingredient.name.toLowerCase().includes(input.toLowerCase())
}

const disabledDate = (current) => {
  return current && current < dayjs().startOf('day')
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

const formatCurrency = (value) => {
  if (value === undefined || value === null || isNaN(value)) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}

const formatNumber = (value) => {
  if (value === undefined || value === null || isNaN(value)) return '0'
  return new Intl.NumberFormat('id-ID').format(value)
}

const parseCurrency = (value) => {
  if (!value) return 0
  const cleaned = value.toString().replace(/Rp\s?|[.]/g, '').replace(',', '.')
  return parseFloat(cleaned) || 0
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchPurchaseOrders()
  fetchSuppliers()
  fetchIngredients()
  fetchApprovedRABs()
  fetchSPPGList()
})
</script>

<style scoped>
.purchase-order-list {
  padding: 24px;
}
</style>
