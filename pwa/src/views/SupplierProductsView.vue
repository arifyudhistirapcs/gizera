<template>
  <div class="supplier-products-page">
    <!-- Header -->
    <van-nav-bar title="Katalog Produk" fixed />

    <div class="page-content">
      <!-- Loading -->
      <van-loading v-if="loading" type="spinner" vertical class="page-loading">
        Memuat produk...
      </van-loading>

      <!-- Empty -->
      <van-empty
        v-else-if="products.length === 0"
        image="search"
        description="Belum ada produk"
      />

      <!-- Product List -->
      <van-cell-group v-else inset>
        <van-swipe-cell v-for="product in products" :key="product.id">
          <van-cell
            :title="product.ingredient?.name || product.ingredient_name || '-'"
            :label="`Stok: ${product.stock_quantity ?? 0} • Min: ${product.min_order_qty ?? 0}`"
            :value="formatRupiah(product.unit_price)"
            is-link
            @click="editProduct(product)"
          >
            <template #right-icon>
              <van-switch
                :model-value="product.is_available"
                size="20px"
                @click.stop
                @update:model-value="(val) => toggleAvailability(product, val)"
              />
            </template>
          </van-cell>
          <template #right>
            <van-button
              square
              type="danger"
              text="Hapus"
              class="swipe-btn"
              @click="deleteProduct(product)"
            />
          </template>
        </van-swipe-cell>
      </van-cell-group>
    </div>

    <!-- FAB -->
    <div class="fab" @click="$router.push('/supplier-products/form')">
      <van-icon name="plus" size="24" color="#fff" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import supplierPortalService from '@/services/supplierPortalService'

const router = useRouter()
const loading = ref(false)
const products = ref([])

function formatRupiah(value) {
  if (!value && value !== 0) return 'Rp 0'
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value)
}

async function loadProducts() {
  loading.value = true
  try {
    const res = await supplierPortalService.getProducts()
    products.value = res.data?.data ?? res.data ?? []
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal memuat produk')
  } finally {
    loading.value = false
  }
}

function editProduct(product) {
  router.push(`/supplier-products/${product.id}/edit`)
}

async function toggleAvailability(product, val) {
  try {
    await supplierPortalService.updateProduct(product.id, { is_available: val })
    product.is_available = val
    showToast(val ? 'Produk diaktifkan' : 'Produk dinonaktifkan')
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal mengubah status')
  }
}

async function deleteProduct(product) {
  try {
    await showConfirmDialog({
      title: 'Hapus Produk',
      message: `Hapus "${product.ingredient?.name || product.ingredient_name || 'produk ini'}"?`
    })
    await supplierPortalService.deleteProduct(product.id)
    products.value = products.value.filter(p => p.id !== product.id)
    showToast('Produk dihapus')
  } catch (e) {
    if (e !== 'cancel' && e?.message !== 'cancel') {
      showToast(e.response?.data?.message || 'Gagal menghapus produk')
    }
  }
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.supplier-products-page {
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

.swipe-btn {
  height: 100%;
}

.fab {
  position: fixed;
  bottom: 100px;
  right: 20px;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  background: #F82C17;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(248, 44, 23, 0.4);
  z-index: 99;
  cursor: pointer;
  transition: transform 0.2s;
}

.fab:active {
  transform: scale(0.92);
}
</style>
