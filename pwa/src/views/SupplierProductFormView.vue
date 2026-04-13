<template>
  <div class="supplier-product-form-page">
    <!-- Header -->
    <van-nav-bar
      :title="isEdit ? 'Edit Produk' : 'Tambah Produk'"
      left-arrow
      fixed
      @click-left="$router.back()"
    />

    <div class="page-content">
      <!-- Loading existing product -->
      <van-loading v-if="loadingProduct" type="spinner" vertical class="page-loading">
        Memuat data produk...
      </van-loading>

      <van-form v-else @submit="onSubmit" class="product-form">
        <!-- Ingredient Picker -->
        <van-field
          v-model="ingredientLabel"
          is-link
          readonly
          label="Bahan"
          placeholder="Pilih bahan"
          :rules="[{ required: true, message: 'Pilih bahan' }]"
          @click="showIngredientPicker = true"
        />

        <van-popup v-model:show="showIngredientPicker" position="bottom" round>
          <van-picker
            :columns="ingredientColumns"
            @confirm="onIngredientConfirm"
            @cancel="showIngredientPicker = false"
          />
        </van-popup>

        <!-- Unit Price -->
        <van-field
          v-model="form.unit_price"
          type="number"
          label="Harga Satuan (Rp)"
          placeholder="Masukkan harga"
          :rules="[{ required: true, message: 'Masukkan harga' }]"
        />

        <!-- Min Order Qty -->
        <van-field
          v-model="form.min_order_qty"
          type="number"
          label="Min. Order"
          placeholder="Jumlah minimum order"
          :rules="[{ required: true, message: 'Masukkan min order' }]"
        />

        <!-- Stock Quantity -->
        <van-field
          v-model="form.stock_quantity"
          type="number"
          label="Stok"
          placeholder="Jumlah stok tersedia"
          :rules="[{ required: true, message: 'Masukkan stok' }]"
        />

        <!-- Availability -->
        <van-cell title="Tersedia" center>
          <template #right-icon>
            <van-switch v-model="form.is_available" size="24px" />
          </template>
        </van-cell>

        <!-- Submit -->
        <div class="form-actions">
          <van-button
            type="primary"
            block
            native-type="submit"
            :loading="submitting"
            :disabled="submitting"
          >
            {{ isEdit ? 'Simpan Perubahan' : 'Tambah Produk' }}
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showToast } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'
import api from '@/services/api'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const loadingProduct = ref(false)
const submitting = ref(false)
const showIngredientPicker = ref(false)
const ingredientLabel = ref('')
const ingredients = ref([])

const form = ref({
  ingredient_id: null,
  unit_price: '',
  min_order_qty: '',
  stock_quantity: '',
  is_available: true
})

const ingredientColumns = computed(() =>
  ingredients.value.map(i => ({ text: i.name, value: i.id }))
)

function onIngredientConfirm({ selectedOptions }) {
  const selected = selectedOptions[0]
  if (selected) {
    form.value.ingredient_id = selected.value
    ingredientLabel.value = selected.text
  }
  showIngredientPicker.value = false
}

async function loadIngredients() {
  try {
    const res = await api.get('/ingredients', { params: { per_page: 500 } })
    ingredients.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    // Fallback: ingredients may not be available
    ingredients.value = []
  }
}

async function loadProduct() {
  if (!isEdit.value) return
  loadingProduct.value = true
  try {
    const res = await supplierPortalService.getProducts({ id: route.params.id })
    // Try to get single product from list or direct
    const data = res.data?.data
    const product = Array.isArray(data)
      ? data.find(p => String(p.id) === String(route.params.id))
      : data
    if (product) {
      form.value.ingredient_id = product.ingredient_id
      form.value.unit_price = String(product.unit_price ?? '')
      form.value.min_order_qty = String(product.min_order_qty ?? '')
      form.value.stock_quantity = String(product.stock_quantity ?? '')
      form.value.is_available = product.is_available ?? true
      ingredientLabel.value = product.ingredient?.name || product.ingredient_name || ''
    }
  } catch (e) {
    showToast('Gagal memuat data produk')
  } finally {
    loadingProduct.value = false
  }
}

async function onSubmit() {
  if (!form.value.ingredient_id) {
    showToast('Pilih bahan terlebih dahulu')
    return
  }

  submitting.value = true
  const payload = {
    ingredient_id: Number(form.value.ingredient_id),
    unit_price: Number(form.value.unit_price),
    min_order_qty: Number(form.value.min_order_qty),
    stock_quantity: Number(form.value.stock_quantity),
    is_available: form.value.is_available
  }

  try {
    if (isEdit.value) {
      await supplierPortalService.updateProduct(route.params.id, payload)
      showToast('Produk berhasil diperbarui')
    } else {
      await supplierPortalService.createProduct(payload)
      showToast('Produk berhasil ditambahkan')
    }
    router.back()
  } catch (e) {
    const msg = e.response?.data?.message || 'Gagal menyimpan produk'
    showToast(msg)
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadIngredients()
  loadProduct()
})
</script>

<style scoped>
.supplier-product-form-page {
  min-height: 100vh;
  background: #F7F8FA;
  padding-top: 46px;
  padding-bottom: 88px;
}

.page-content {
  padding: 16px;
}

.page-loading {
  padding: 60px 0;
}

.product-form {
  background: #fff;
  border-radius: 14px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.form-actions {
  padding: 16px;
}
</style>
