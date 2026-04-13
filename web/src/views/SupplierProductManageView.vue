<template>
  <div class="supplier-product-manage">
    <a-page-header
      title="Kelola Katalog Produk"
      sub-title="Tambah, edit, dan hapus produk Anda"
    >
      <template #extra>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Tambah Produk
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-table
        :columns="columns"
        :data-source="products"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'ingredient_name'">
            {{ record.ingredient?.name || '-' }}
          </template>
          <template v-else-if="column.key === 'unit_price'">
            {{ formatRupiah(record.unit_price) }}
          </template>
          <template v-else-if="column.key === 'is_available'">
            <a-tag :color="record.is_available ? 'green' : 'red'">
              {{ record.is_available ? 'Tersedia' : 'Tidak Tersedia' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="showEditModal(record)">
                Edit
              </a-button>
              <a-popconfirm
                title="Yakin ingin menghapus produk ini?"
                ok-text="Ya"
                cancel-text="Tidak"
                @confirm="handleDelete(record.id)"
              >
                <a-button type="link" size="small" danger>
                  Hapus
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingProduct ? 'Edit Produk' : 'Tambah Produk'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="600px"
    >
      <a-form ref="formRef" :model="formData" :rules="rules" layout="vertical">
        <a-form-item label="Ingredient" name="ingredient_id">
          <a-select
            v-model:value="formData.ingredient_id"
            placeholder="Pilih ingredient"
            show-search
            :filter-option="filterIngredient"
            :disabled="!!editingProduct"
          >
            <a-select-option v-for="ing in ingredients" :key="ing.id" :value="ing.id">
              {{ ing.name }} ({{ ing.unit }})
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Harga Satuan (Rp)" name="unit_price">
              <a-input-number v-model:value="formData.unit_price" :min="0" :step="1000" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Min. Order Qty" name="min_order_qty">
              <a-input-number v-model:value="formData.min_order_qty" :min="0" :step="1" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Stok" name="stock_quantity">
              <a-input-number v-model:value="formData.stock_quantity" :min="0" :step="1" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Tersedia" name="is_available">
              <a-switch v-model:checked="formData.is_available" checked-children="Ya" un-checked-children="Tidak" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import supplierProductService from '@/services/supplierProductService'
import recipeService from '@/services/recipeService'

const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const editingProduct = ref(null)
const products = ref([])
const ingredients = ref([])
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  ingredient_id: undefined,
  unit_price: 0,
  min_order_qty: 1,
  stock_quantity: 0,
  is_available: true
})

const rules = {
  ingredient_id: [{ required: true, message: 'Ingredient wajib dipilih' }],
  unit_price: [{ required: true, message: 'Harga satuan wajib diisi' }],
  min_order_qty: [{ required: true, message: 'Min. order qty wajib diisi' }]
}

const columns = [
  { title: 'Ingredient', key: 'ingredient_name' },
  { title: 'Harga Satuan', key: 'unit_price' },
  { title: 'Min. Order', dataIndex: 'min_order_qty', key: 'min_order_qty' },
  { title: 'Stok', dataIndex: 'stock_quantity', key: 'stock_quantity' },
  { title: 'Ketersediaan', key: 'is_available', width: 140 },
  { title: 'Aksi', key: 'actions', width: 160 }
]

const fetchProducts = async () => {
  loading.value = true
  try {
    const params = { page: pagination.current, per_page: pagination.pageSize }
    const response = await supplierProductService.getProducts(params)
    products.value = response.data.products || response.data.data || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat data produk')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchIngredients = async () => {
  try {
    const response = await recipeService.getIngredients()
    ingredients.value = response.data.data || []
  } catch (error) {
    console.error('Gagal memuat data ingredient:', error)
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchProducts()
}

const showCreateModal = () => {
  editingProduct.value = null
  resetForm()
  modalVisible.value = true
}

const showEditModal = (product) => {
  editingProduct.value = product
  Object.assign(formData, {
    ingredient_id: product.ingredient_id,
    unit_price: product.unit_price,
    min_order_qty: product.min_order_qty,
    stock_quantity: product.stock_quantity,
    is_available: product.is_available
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    if (editingProduct.value) {
      await supplierProductService.updateProduct(editingProduct.value.id, formData)
      message.success('Produk berhasil diperbarui')
    } else {
      await supplierProductService.createProduct(formData)
      message.success('Produk berhasil ditambahkan')
    }

    modalVisible.value = false
    fetchProducts()
  } catch (error) {
    if (error.errorFields) return
    const errMsg = error.response?.data?.message || ''
    if (errMsg.includes('DUPLICATE_SUPPLIER_PRODUCT') || errMsg.includes('duplicate')) {
      message.error('Produk untuk ingredient ini sudah ada')
    } else {
      message.error(errMsg || 'Gagal menyimpan produk')
    }
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await supplierProductService.deleteProduct(id)
    message.success('Produk berhasil dihapus')
    fetchProducts()
  } catch (error) {
    message.error('Gagal menghapus produk')
    console.error(error)
  }
}

const handleCancel = () => {
  modalVisible.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    ingredient_id: undefined,
    unit_price: 0,
    min_order_qty: 1,
    stock_quantity: 0,
    is_available: true
  })
  formRef.value?.resetFields()
}

const filterIngredient = (input, option) => {
  const ing = ingredients.value.find(i => i.id === option.value)
  return ing?.name?.toLowerCase().includes(input.toLowerCase()) ?? false
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

onMounted(() => {
  fetchProducts()
  fetchIngredients()
})
</script>

<style scoped>
.supplier-product-manage {
  padding: 24px;
}
</style>
