<template>
  <div class="supplier-product-list">
    <a-page-header
      title="Katalog Produk Supplier"
      sub-title="Lihat produk dari supplier yang terhubung"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Filters -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-select
              v-model:value="filterSupplier"
              placeholder="Filter Supplier"
              style="width: 100%"
              @change="handleSearch"
              allow-clear
              size="large"
              show-search
              :filter-option="filterOption"
            >
              <a-select-option v-for="s in suppliers" :key="s.id" :value="s.id">
                {{ s.name }}
              </a-select-option>
            </a-select>
          </a-col>
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari ingredient..."
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

        <!-- Table -->
        <a-table
          :columns="columns"
          :data-source="products"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'supplier_name'">
              {{ record.supplier?.name || '-' }}
            </template>
            <template v-else-if="column.key === 'ingredient_name'">
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
          </template>
        </a-table>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import supplierProductService from '@/services/supplierProductService'
import supplierService from '@/services/supplierService'

const loading = ref(false)
const products = ref([])
const suppliers = ref([])
const searchText = ref('')
const filterSupplier = ref(undefined)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'Supplier', key: 'supplier_name' },
  { title: 'Ingredient', key: 'ingredient_name' },
  { title: 'Harga Satuan', key: 'unit_price' },
  { title: 'Min. Order', dataIndex: 'min_order_qty', key: 'min_order_qty' },
  { title: 'Stok', dataIndex: 'stock_quantity', key: 'stock_quantity' },
  { title: 'Ketersediaan', key: 'is_available', width: 140 }
]

const fetchProducts = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      per_page: pagination.pageSize,
      supplier_id: filterSupplier.value || undefined,
      search: searchText.value || undefined
    }
    const response = await supplierProductService.getProducts(params)
    products.value = response.data.products || response.data.data || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat katalog produk')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchSuppliers = async () => {
  try {
    const response = await supplierService.getSuppliers({ is_active: 'active' })
    suppliers.value = response.data.suppliers || []
  } catch (error) {
    console.error('Gagal memuat data supplier:', error)
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchProducts()
}

const handleSearch = () => {
  pagination.current = 1
  fetchProducts()
}

const filterOption = (input, option) => {
  return option.children?.[0]?.toLowerCase?.()?.includes(input.toLowerCase()) ?? false
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

onMounted(() => {
  fetchProducts()
  fetchSuppliers()
})
</script>

<style scoped>
.supplier-product-list {
  padding: 24px;
}
</style>
