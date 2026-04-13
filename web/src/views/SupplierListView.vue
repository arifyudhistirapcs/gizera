<template>
  <div class="supplier-list">
    <a-page-header
      title="Manajemen Supplier"
      sub-title="Kelola data supplier dan lihat performa"
    >
      <template #extra>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Tambah Supplier
        </a-button>
      </template>
    </a-page-header>

    <!-- Statistics Cards -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="Total Supplier"
            :value="stats.totalSuppliers"
            :loading="loadingStats"
          >
            <template #prefix>
              <ShopOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="Total Pengeluaran"
            :value="stats.totalSpending"
            :loading="loadingStats"
            :precision="0"
          >
            <template #prefix>
              <DollarOutlined />
            </template>
            <template #formatter="{ value }">
              {{ formatCurrency(value) }}
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="Supplier Aktif"
            :value="stats.activeSuppliers"
            :loading="loadingStats"
          >
            <template #prefix>
              <CheckCircleOutlined style="color: #52c41a" />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="6">
        <a-card>
          <a-statistic
            title="Rata-rata Rating"
            :value="stats.averageRating"
            :loading="loadingStats"
            :precision="1"
            suffix="/ 5"
          >
            <template #prefix>
              <StarOutlined style="color: #faad14" />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <!-- Top 3 Suppliers Card -->
    <a-card title="Top 3 Supplier" style="margin-bottom: 16px" :loading="loadingStats">
      <a-list
        :data-source="stats.topSuppliers"
        :loading="loadingStats"
      >
        <template #renderItem="{ item, index }">
          <a-list-item>
            <a-list-item-meta>
              <template #avatar>
                <a-avatar :style="{ backgroundColor: getTrophyColor(index) }">
                  {{ index + 1 }}
                </a-avatar>
              </template>
              <template #title>
                {{ item.name }}
              </template>
              <template #description>
                <a-space direction="vertical" size="small" style="width: 100%">
                  <div>
                    <a-tag color="blue">{{ item.total_orders }} Order</a-tag>
                    <a-tag color="green">{{ formatCurrency(item.total_amount) }}</a-tag>
                  </div>
                  <a-progress
                    :percent="Math.round((item.total_amount / stats.totalSpending) * 100)"
                    :show-info="true"
                    size="small"
                  />
                </a-space>
              </template>
            </a-list-item-meta>
          </a-list-item>
        </template>
      </a-list>
    </a-card>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Search and Filter -->
        <a-row :gutter="16">
          <a-col :span="12">
            <a-input
              v-model:value="searchText"
              placeholder="Cari nama supplier..."
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
              <a-select-option value="active">Aktif</a-select-option>
              <a-select-option value="inactive">Tidak Aktif</a-select-option>
            </a-select>
          </a-col>
        </a-row>

        <!-- Table -->
        <a-table
          :columns="columns"
          :data-source="suppliers"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'is_active'">
              <a-tag :color="record.is_active ? 'green' : 'red'">
                {{ record.is_active ? 'Aktif' : 'Tidak Aktif' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'on_time_delivery'">
              <a-progress
                :percent="Math.round(record.on_time_delivery || 0)"
                :status="record.on_time_delivery >= 80 ? 'success' : 'normal'"
                size="small"
              />
            </template>
            <template v-else-if="column.key === 'quality_rating'">
              <a-rate :value="record.quality_rating || 0" disabled allow-half />
            </template>
            <template v-else-if="column.key === 'product_count'">
              <a-tag color="blue">{{ record.product_count || 0 }} produk</a-tag>
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="viewSupplier(record)">
                  Detail
                </a-button>
                <a-button type="link" size="small" @click="editSupplier(record)">
                  Edit
                </a-button>
                <a-popconfirm
                  title="Yakin ingin menghapus supplier ini?"
                  ok-text="Ya"
                  cancel-text="Tidak"
                  @confirm="deleteSupplier(record.id)"
                >
                  <a-button type="link" size="small" danger>
                    Hapus
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingSupplier ? 'Edit Supplier' : 'Tambah Supplier'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="650px"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="Nama Supplier" name="name">
          <a-input v-model:value="formData.name" placeholder="Masukkan nama supplier" />
        </a-form-item>

        <a-form-item label="Kategori Produk" name="product_category">
          <a-input v-model:value="formData.product_category" placeholder="Contoh: Sayuran, Daging, Bumbu" />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Nama Kontak" name="contact_person">
              <a-input v-model:value="formData.contact_person" placeholder="Nama kontak person" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Nomor Telepon" name="phone_number">
              <a-input v-model:value="formData.phone_number" placeholder="08xxxxxxxxxx" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="Email" name="email">
          <a-input v-model:value="formData.email" type="email" placeholder="email@supplier.com" />
        </a-form-item>

        <a-form-item label="Alamat" name="address">
          <a-textarea v-model:value="formData.address" :rows="3" placeholder="Alamat lengkap supplier" />
        </a-form-item>

        <a-form-item label="Lokasi">
          <LocationPicker
            v-model:latitude="formData.latitude"
            v-model:longitude="formData.longitude"
          />
        </a-form-item>

        <a-form-item label="Status" name="is_active">
          <a-switch v-model:checked="formData.is_active" checked-children="Aktif" un-checked-children="Tidak Aktif" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Detail Modal -->
    <a-modal
      v-model:open="detailModalVisible"
      title="Detail Supplier"
      :footer="null"
      width="800px"
    >
      <a-descriptions v-if="selectedSupplier" bordered :column="2">
        <a-descriptions-item label="Nama Supplier" :span="2">
          {{ selectedSupplier.name }}
        </a-descriptions-item>
        <a-descriptions-item label="Kategori Produk" :span="2">
          {{ selectedSupplier.product_category }}
        </a-descriptions-item>
        <a-descriptions-item label="Kontak Person">
          {{ selectedSupplier.contact_person }}
        </a-descriptions-item>
        <a-descriptions-item label="Telepon">
          {{ selectedSupplier.phone_number }}
        </a-descriptions-item>
        <a-descriptions-item label="Email" :span="2">
          {{ selectedSupplier.email }}
        </a-descriptions-item>
        <a-descriptions-item label="Alamat" :span="2">
          {{ selectedSupplier.address }}
        </a-descriptions-item>
        <a-descriptions-item label="Status">
          <a-tag :color="selectedSupplier.is_active ? 'green' : 'red'">
            {{ selectedSupplier.is_active ? 'Aktif' : 'Tidak Aktif' }}
          </a-tag>
        </a-descriptions-item>
      </a-descriptions>

      <a-divider>Metrik Performa</a-divider>

      <a-row :gutter="16">
        <a-col :span="12">
          <a-statistic title="Pengiriman Tepat Waktu" :value="selectedSupplier.on_time_delivery || 0" suffix="%" />
          <a-typography-text type="secondary" style="font-size: 12px">
            Dihitung otomatis saat barang diterima
          </a-typography-text>
        </a-col>
        <a-col :span="12">
          <a-statistic title="Rating Kualitas">
            <template #formatter>
              <a-rate :value="selectedSupplier.quality_rating || 0" disabled allow-half />
            </template>
          </a-statistic>
          <a-typography-text type="secondary" style="font-size: 12px">
            Rata-rata dari rating saat penerimaan barang (GRN)
          </a-typography-text>
        </a-col>
      </a-row>

      <a-divider>Riwayat Transaksi</a-divider>

      <a-table
        :columns="transactionColumns"
        :data-source="transactionHistory"
        :loading="loadingTransactions"
        :pagination="{ pageSize: 5 }"
        size="small"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'amount'">
            {{ formatCurrency(record.amount) }}
          </template>
          <template v-else-if="column.key === 'order_date'">
            {{ formatDate(record.order_date) }}
          </template>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  SearchOutlined,
  ShopOutlined,
  DollarOutlined,
  CheckCircleOutlined,
  StarOutlined
} from '@ant-design/icons-vue'
import supplierService from '@/services/supplierService'
import LocationPicker from '@/components/LocationPicker.vue'

const loading = ref(false)
const loadingStats = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const detailModalVisible = ref(false)
const editingSupplier = ref(null)
const selectedSupplier = ref(null)
const suppliers = ref([])
const transactionHistory = ref([])
const loadingTransactions = ref(false)
const searchText = ref('')
const filterStatus = ref(undefined)
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const stats = reactive({
  totalSuppliers: 0,
  totalSpending: 0,
  activeSuppliers: 0,
  averageRating: 0,
  topSuppliers: []
})

const formData = reactive({
  name: '',
  product_category: '',
  contact_person: '',
  phone_number: '',
  email: '',
  address: '',
  latitude: 0,
  longitude: 0,
  is_active: true
})

const rules = {
  name: [{ required: true, message: 'Nama supplier wajib diisi' }],
  product_category: [{ required: true, message: 'Kategori produk wajib diisi' }],
  contact_person: [{ required: true, message: 'Nama kontak wajib diisi' }],
  phone_number: [{ required: true, message: 'Nomor telepon wajib diisi' }],
  email: [
    { required: true, message: 'Email wajib diisi' },
    { type: 'email', message: 'Format email tidak valid' }
  ]
}

const columns = [
  {
    title: 'Nama Supplier',
    dataIndex: 'name',
    key: 'name',
    sorter: true
  },
  {
    title: 'Kategori Produk',
    dataIndex: 'product_category',
    key: 'product_category'
  },
  {
    title: 'Kontak',
    dataIndex: 'contact_person',
    key: 'contact_person'
  },
  {
    title: 'Telepon',
    dataIndex: 'phone_number',
    key: 'phone_number'
  },
  {
    title: 'Pengiriman Tepat Waktu',
    key: 'on_time_delivery',
    width: 180
  },
  {
    title: 'Rating Kualitas',
    key: 'quality_rating',
    width: 150
  },
  {
    title: 'Katalog Produk',
    key: 'product_count',
    width: 130
  },
  {
    title: 'Status',
    key: 'is_active',
    width: 100
  },
  {
    title: 'Aksi',
    key: 'actions',
    width: 200
  }
]

const transactionColumns = [
  {
    title: 'Nomor PO',
    dataIndex: 'po_number',
    key: 'po_number'
  },
  {
    title: 'Tanggal',
    key: 'order_date'
  },
  {
    title: 'Jumlah',
    key: 'amount'
  },
  {
    title: 'Status',
    dataIndex: 'status',
    key: 'status'
  }
]

const fetchSuppliers = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      page_size: pagination.pageSize,
      search: searchText.value || undefined,
      is_active: filterStatus.value
    }
    const response = await supplierService.getSuppliers(params)
    suppliers.value = response.data.suppliers || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat data supplier')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchSupplierStats = async () => {
  loadingStats.value = true
  try {
    const response = await supplierService.getSupplierStats()
    console.log('Supplier stats response:', response)
    console.log('Response data:', response.data)
    // Backend returns { success: true, data: { totalSuppliers, ... } }
    // So we need response.data.data
    if (response.data && response.data.data) {
      Object.assign(stats, response.data.data)
    }
    console.log('Stats after assign:', stats)
  } catch (error) {
    console.error('Gagal memuat statistik supplier:', error)
    console.error('Error details:', error.response || error)
  } finally {
    loadingStats.value = false
  }
}

const getTrophyColor = (index) => {
  const colors = ['#FFD700', '#C0C0C0', '#CD7F32'] // Gold, Silver, Bronze
  return colors[index] || '#1890ff'
}

const handleTableChange = (pag, filters, sorter) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchSuppliers()
}

const handleSearch = () => {
  pagination.current = 1
  fetchSuppliers()
}

const showCreateModal = () => {
  editingSupplier.value = null
  resetForm()
  modalVisible.value = true
}

const editSupplier = (supplier) => {
  editingSupplier.value = supplier
  Object.assign(formData, {
    name: supplier.name,
    product_category: supplier.product_category,
    contact_person: supplier.contact_person,
    phone_number: supplier.phone_number,
    email: supplier.email,
    address: supplier.address,
    latitude: supplier.latitude || 0,
    longitude: supplier.longitude || 0,
    is_active: supplier.is_active
  })
  modalVisible.value = true
}

const viewSupplier = async (supplier) => {
  selectedSupplier.value = supplier
  detailModalVisible.value = true
  
  // Fetch transaction history and performance metrics
  loadingTransactions.value = true
  try {
    const response = await supplierService.getSupplierPerformance(supplier.id)
    console.log('Performance response:', response.data)
    // Backend returns { success: true, data: { ...performance } }
    if (response.data && response.data.data) {
      transactionHistory.value = response.data.data.transactions || []
      // Update supplier with latest performance metrics
      selectedSupplier.value = {
        ...supplier,
        on_time_delivery: response.data.data.on_time_rate || 0,
        quality_rating: response.data.data.quality_rating || 0
      }
    }
  } catch (error) {
    console.error('Gagal memuat riwayat transaksi:', error)
  } finally {
    loadingTransactions.value = false
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    if (editingSupplier.value) {
      await supplierService.updateSupplier(editingSupplier.value.id, formData)
      message.success('Supplier berhasil diperbarui')
    } else {
      await supplierService.createSupplier(formData)
      message.success('Supplier berhasil ditambahkan')
    }

    modalVisible.value = false
    fetchSuppliers()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    message.error('Gagal menyimpan data supplier')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const deleteSupplier = async (id) => {
  try {
    await supplierService.deleteSupplier(id)
    message.success('Supplier berhasil dihapus')
    fetchSuppliers()
  } catch (error) {
    message.error('Gagal menghapus supplier')
    console.error(error)
  }
}

const handleCancel = () => {
  modalVisible.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    name: '',
    product_category: '',
    contact_person: '',
    phone_number: '',
    email: '',
    address: '',
    latitude: 0,
    longitude: 0,
    is_active: true
  })
  formRef.value?.resetFields()
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR'
  }).format(value)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

onMounted(() => {
  fetchSuppliers()
  fetchSupplierStats()
})
</script>

<style scoped>
.supplier-list {
  padding: 24px;
}
</style>
