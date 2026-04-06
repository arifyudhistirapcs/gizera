<template>
  <div class="kds-cooking">
    <!-- Header Controls -->
    <div class="cooking-header">
      <div class="header-left">
        <h2 class="page-title">
          <FireOutlined class="title-icon" />
          KDS Dapur Memasak
        </h2>
        <div class="status-summary">
          <span class="summary-chip chip-pending">
            <ClockCircleOutlined /> {{ pendingCount }} Menunggu
          </span>
          <span class="summary-chip chip-cooking">
            <FireOutlined /> {{ cookingCount }} Dimasak
          </span>
          <span class="summary-chip chip-ready">
            <CheckCircleOutlined /> {{ readyCount }} Selesai
          </span>
        </div>
      </div>
      <div class="header-controls">
        <KDSDatePicker
          v-model="selectedDate"
          :loading="loading"
          @change="handleDateChange"
        />
        <a-tag :color="isConnected ? 'green' : 'red'" class="connection-tag">
          <template #icon>
            <wifi-outlined v-if="isConnected" />
            <disconnect-outlined v-else />
          </template>
          {{ isConnected ? 'Terhubung' : 'Terputus' }}
        </a-tag>
        <a-button @click="refreshData" :loading="loading" type="default">
          <template #icon><reload-outlined /></template>
          Refresh
        </a-button>
      </div>
    </div>

    <!-- Alert -->
    <a-alert
      v-if="error"
      type="error"
      :message="error"
      closable
      show-icon
      @close="error = null"
      class="cooking-alert"
    >
      <template #action>
        <a-button size="small" type="primary" @click="retryLoad">
          Coba Lagi
        </a-button>
      </template>
    </a-alert>

    <!-- Bento Grid -->
    <a-spin :spinning="loading" tip="Memuat data...">
      <a-empty v-if="!loading && recipes.length === 0" :description="emptyMessage" />

      <div v-else class="bento-grid">
        <!-- COOKING cards first (prominent, larger) -->
        <div
          v-for="recipe in cookingRecipes"
          :key="'c-' + recipe.recipe_id"
          class="bento-card bento-card--cooking"
        >
          <div class="card-top">
            <div class="card-header">
              <div class="card-title-row">
                <span class="card-name">{{ recipe.name }}</span>
                <span class="status-badge badge-cooking">
                  <span class="badge-dot"></span>
                  Sedang Dimasak
                </span>
              </div>
              <div v-if="recipe.start_time" class="card-timer">
                <ClockCircleOutlined />
                Mulai {{ formatTime(recipe.start_time) }}
              </div>
            </div>
            <div class="card-portions">
              <span class="portions-number">{{ recipe.portions_required }}</span>
              <span class="portions-label">porsi</span>
            </div>
          </div>

          <div class="card-body">
            <div class="ingredients-section">
              <div class="section-label">Bahan</div>
              <div class="ingredients-compact">
                <div v-for="item in recipe.items" :key="item.name" class="ingredient-row">
                  <span class="ing-name">{{ item.name }}</span>
                  <span class="ing-qty">{{ item.quantity }} {{ item.unit }}</span>
                </div>
              </div>
            </div>

            <div v-if="recipe.instructions" class="instructions-section">
              <div
                class="section-label section-label--toggle"
                @click="toggleInstructions(recipe.recipe_id)"
                role="button"
                tabindex="0"
                :aria-expanded="expandedInstructions[recipe.recipe_id] ? 'true' : 'false'"
                @keydown.enter.space.prevent="toggleInstructions(recipe.recipe_id)"
              >
                Instruksi
                <DownOutlined class="toggle-icon" :class="{ 'toggle-icon--open': expandedInstructions[recipe.recipe_id] }" />
              </div>
              <div v-if="expandedInstructions[recipe.recipe_id]" class="instructions-content">
                {{ recipe.instructions }}
              </div>
            </div>
          </div>

          <a-button type="primary" block size="large" @click="finishCooking(recipe)" :loading="updatingRecipeId === recipe.recipe_id" class="action-btn action-btn--finish">
            <template #icon><CheckCircleOutlined /></template>
            Selesai Masak
          </a-button>
        </div>

        <!-- PENDING cards -->
        <div v-for="recipe in pendingRecipes" :key="'p-' + recipe.recipe_id" class="bento-card bento-card--pending">
          <div class="card-top">
            <div class="card-header">
              <div class="card-title-row">
                <span class="card-name">{{ recipe.name }}</span>
                <span class="status-badge badge-pending">
                  <span class="badge-dot"></span>
                  Belum Dimulai
                </span>
              </div>
            </div>
            <div class="card-portions">
              <span class="portions-number">{{ recipe.portions_required }}</span>
              <span class="portions-label">porsi</span>
            </div>
          </div>

          <div class="card-body">
            <div class="ingredients-section">
              <div class="section-label">Bahan</div>
              <div class="ingredients-compact">
                <div v-for="item in recipe.items" :key="item.name" class="ingredient-row">
                  <span class="ing-name">{{ item.name }}</span>
                  <span class="ing-qty">{{ item.quantity }} {{ item.unit }}</span>
                </div>
              </div>
            </div>

            <div v-if="recipe.instructions" class="instructions-section">
              <div class="section-label section-label--toggle" @click="toggleInstructions(recipe.recipe_id)" role="button" tabindex="0" :aria-expanded="expandedInstructions[recipe.recipe_id] ? 'true' : 'false'" @keydown.enter.space.prevent="toggleInstructions(recipe.recipe_id)">
                Instruksi
                <DownOutlined class="toggle-icon" :class="{ 'toggle-icon--open': expandedInstructions[recipe.recipe_id] }" />
              </div>
              <div v-if="expandedInstructions[recipe.recipe_id]" class="instructions-content">
                {{ recipe.instructions }}
              </div>
            </div>
          </div>

          <a-button type="primary" block size="large" @click="startCooking(recipe)" :loading="updatingRecipeId === recipe.recipe_id" class="action-btn action-btn--start">
            <template #icon><PlayCircleOutlined /></template>
            Mulai Masak
          </a-button>
        </div>

        <!-- READY cards (compact, muted) -->
        <div v-for="recipe in readyRecipes" :key="'r-' + recipe.recipe_id" class="bento-card bento-card--ready">
          <div class="card-top card-top--ready">
            <div class="card-header">
              <div class="card-title-row">
                <span class="card-name">{{ recipe.name }}</span>
                <span class="status-badge badge-ready">
                  <CheckOutlined />
                  Selesai
                </span>
              </div>
            </div>
            <div class="card-portions card-portions--ready">
              <span class="portions-number">{{ recipe.portions_required }}</span>
              <span class="portions-label">porsi</span>
            </div>
          </div>
          <div v-if="recipe.start_time || recipe.end_time" class="card-times">
            <span v-if="recipe.start_time" class="time-chip">Mulai {{ formatTime(recipe.start_time) }}</span>
            <span v-if="recipe.end_time" class="time-chip">Selesai {{ formatTime(recipe.end_time) }}</span>
            <span v-if="recipe.duration_minutes" class="time-chip time-chip--duration">{{ recipe.duration_minutes }} menit</span>
          </div>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  WifiOutlined, DisconnectOutlined, ReloadOutlined, PlayCircleOutlined,
  CheckCircleOutlined, CheckOutlined, ClockCircleOutlined, FireOutlined, DownOutlined
} from '@ant-design/icons-vue'
import KDSDatePicker from '@/components/KDSDatePicker.vue'
import { getCookingToday, updateCookingStatus } from '@/services/kdsService'
import { database, firebasePaths } from '@/services/firebase'
import { ref as dbRef, onValue, off } from 'firebase/database'

const recipes = ref([])
const loading = ref(false)
const updatingRecipeId = ref(null)
const isConnected = ref(true)
const selectedDate = ref(new Date())
const error = ref(null)
const expandedInstructions = reactive({})
let firebaseListener = null

const pendingRecipes = computed(() => recipes.value.filter(r => r.status === 'pending'))
const cookingRecipes = computed(() => recipes.value.filter(r => r.status === 'cooking'))
const readyRecipes = computed(() => recipes.value.filter(r => r.status === 'ready'))
const pendingCount = computed(() => pendingRecipes.value.length)
const cookingCount = computed(() => cookingRecipes.value.length)
const readyCount = computed(() => readyRecipes.value.length)

const emptyMessage = computed(() => {
  const today = new Date()
  const isToday = selectedDate.value.toDateString() === today.toDateString()
  return isToday ? 'Tidak ada menu untuk hari ini' : 'Tidak ada menu untuk tanggal ini'
})

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

const toggleInstructions = (recipeId) => {
  expandedInstructions[recipeId] = !expandedInstructions[recipeId]
}

const loadData = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await getCookingToday(selectedDate.value)
    if (response.success) {
      recipes.value = response.data || []
    } else {
      error.value = response.message || 'Gagal memuat data'
    }
  } catch (err) {
    console.error('Error loading cooking data:', err)
    error.value = err.response?.data?.message || 'Gagal memuat data menu. Silakan coba lagi.'
  } finally {
    loading.value = false
  }
}

const retryLoad = () => loadData()
const refreshData = () => loadData()

const startCooking = async (recipe) => {
  updatingRecipeId.value = recipe.recipe_id
  try {
    const response = await updateCookingStatus(recipe.recipe_id, 'cooking')
    if (response.success) {
      message.success('Status berhasil diperbarui: Mulai Masak')
      await loadData()
    } else {
      message.error(response.message || 'Gagal memperbarui status')
    }
  } catch (err) {
    console.error('Error updating status:', err)
    message.error(err.response?.data?.message || 'Gagal memperbarui status')
  } finally {
    updatingRecipeId.value = null
  }
}

const finishCooking = async (recipe) => {
  updatingRecipeId.value = recipe.recipe_id
  try {
    const response = await updateCookingStatus(recipe.recipe_id, 'ready')
    if (response.success) {
      message.success('Status berhasil diperbarui: Selesai')
      await loadData()
    } else {
      message.error(response.message || 'Gagal memperbarui status')
    }
  } catch (err) {
    console.error('Error updating status:', err)
    message.error(err.response?.data?.message || 'Gagal memperbarui status')
  } finally {
    updatingRecipeId.value = null
  }
}

const setupFirebaseListener = () => {
  cleanupFirebaseListener()
  const dateStr = selectedDate.value.toISOString().split('T')[0]
  const cookingRef = dbRef(database, firebasePaths.kdsCooking(dateStr))
  firebaseListener = onValue(cookingRef, (snapshot) => {
    isConnected.value = true
    const data = snapshot.val()
    if (data) {
      const firebaseRecipes = Object.values(data)
      recipes.value = recipes.value.map(recipe => {
        const fr = firebaseRecipes.find(f => f.recipe_id === recipe.recipe_id)
        if (fr) {
          return { ...recipe, status: fr.status, start_time: fr.start_time, end_time: fr.end_time, duration_minutes: fr.duration_minutes }
        }
        return recipe
      })
    }
  }, (err) => {
    console.error('Firebase listener error:', err)
    isConnected.value = false
  })
}

const cleanupFirebaseListener = () => {
  if (firebaseListener) {
    const dateStr = selectedDate.value.toISOString().split('T')[0]
    const cookingRef = dbRef(database, firebasePaths.kdsCooking(dateStr))
    off(cookingRef)
    firebaseListener = null
  }
}

const handleDateChange = (date) => {
  selectedDate.value = date
  loadData()
  setupFirebaseListener()
}

watch(selectedDate, () => { setupFirebaseListener() })
onMounted(() => { loadData(); setupFirebaseListener() })
onUnmounted(() => { cleanupFirebaseListener() })
</script>

<style scoped>
.kds-cooking { padding: 0; }

.cooking-header { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; margin-bottom: 20px; flex-wrap: wrap; }
.header-left { display: flex; flex-direction: column; gap: 8px; }
.page-title { margin: 0; font-size: 20px; font-weight: 700; color: #303030; display: flex; align-items: center; gap: 8px; }
.title-icon { color: #C94A3A; }
.status-summary { display: flex; gap: 8px; flex-wrap: wrap; }
.summary-chip { display: inline-flex; align-items: center; gap: 4px; padding: 2px 10px; border-radius: 999px; font-size: 12px; font-weight: 600; }
.chip-pending { background: #FFF7E6; color: #D48806; }
.chip-cooking { background: #FDEAE7; color: #C94A3A; }
.chip-ready { background: #D1FAE5; color: #1E8A6E; }
.header-controls { display: flex; align-items: center; gap: 12px; }
.connection-tag { font-size: 14px; padding: 4px 12px; border-radius: 6px; }
.cooking-alert { margin-bottom: 20px; }

.bento-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 20px; align-items: start; }

.bento-card { background: #fff; border-radius: 16px; border: 1px solid #E8E8E8; overflow: hidden; display: flex; flex-direction: column; transition: transform 0.2s, box-shadow 0.2s; }
.bento-card:hover { box-shadow: 0 4px 16px rgba(0,0,0,0.08); }

.card-top { display: flex; justify-content: space-between; align-items: flex-start; padding: 20px; gap: 16px; }
.card-header { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 8px; }
.card-title-row { display: flex; align-items: flex-start; gap: 8px; flex-wrap: wrap; }
.card-name { font-size: 18px; font-weight: 700; color: #303030; line-height: 1.25; }

.status-badge { display: inline-flex; align-items: center; gap: 4px; padding: 2px 10px; border-radius: 999px; font-size: 12px; font-weight: 600; white-space: nowrap; flex-shrink: 0; }
.badge-dot { width: 6px; height: 6px; border-radius: 50%; background: currentColor; }
.badge-pending { background: #FFF7E6; color: #D48806; }
.badge-cooking { background: #FDEAE7; color: #C94A3A; }
.badge-ready { background: #D1FAE5; color: #1E8A6E; }

.card-timer { display: flex; align-items: center; gap: 4px; font-size: 14px; color: #C94A3A; font-weight: 500; }

.card-portions { display: flex; flex-direction: column; align-items: center; justify-content: center; min-width: 80px; height: 80px; background: #303030; border-radius: 12px; color: #fff; flex-shrink: 0; }
.card-portions--ready { background: #1E8A6E; opacity: 0.85; }
.portions-number { font-size: 30px; font-weight: 700; line-height: 1; }
.portions-label { font-size: 12px; font-weight: 500; opacity: 0.8; text-transform: uppercase; letter-spacing: 0.5px; }

.card-body { padding: 0 20px 16px; display: flex; flex-direction: column; gap: 12px; flex: 1; }
.section-label { font-size: 12px; font-weight: 700; color: #6B6B6B; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 4px; }
.section-label--toggle { cursor: pointer; display: flex; align-items: center; gap: 4px; user-select: none; }
.section-label--toggle:hover { color: #303030; }
.toggle-icon { font-size: 10px; transition: transform 0.2s; }
.toggle-icon--open { transform: rotate(180deg); }

.ingredients-compact { display: flex; flex-direction: column; gap: 2px; }
.ingredient-row { display: flex; justify-content: space-between; align-items: center; padding: 4px 8px; border-radius: 6px; background: #F7F8FA; font-size: 14px; }
.ing-name { color: #303030; font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-right: 8px; }
.ing-qty { color: #C94A3A; font-weight: 600; white-space: nowrap; flex-shrink: 0; }

.instructions-content { white-space: pre-wrap; line-height: 1.5; color: #6B6B6B; padding: 8px 12px; background: #F7F8FA; border-radius: 6px; font-size: 14px; max-height: 120px; overflow-y: auto; }

.card-times { display: flex; gap: 8px; padding: 0 20px 16px; flex-wrap: wrap; }
.time-chip { display: inline-flex; align-items: center; gap: 4px; padding: 2px 10px; border-radius: 999px; font-size: 12px; font-weight: 500; background: #F7F8FA; color: #6B6B6B; }
.time-chip--duration { background: #D1FAE5; color: #1E8A6E; font-weight: 600; }

.action-btn { margin: 0 20px 20px; height: 44px; font-weight: 700; border-radius: 12px; font-size: 14px; }
.action-btn--start { background: #303030 !important; border-color: #303030 !important; }
.action-btn--start:hover { background: #1a1a1a !important; border-color: #1a1a1a !important; }
.action-btn--finish { background: #C94A3A !important; border-color: #C94A3A !important; }
.action-btn--finish:hover { background: #A33D30 !important; border-color: #A33D30 !important; }

.bento-card--cooking { border-left: 5px solid #C94A3A; box-shadow: 0 2px 12px rgba(201,74,58,0.10); }
.bento-card--cooking .card-top { background: #FDEAE7; }
.bento-card--pending { border-left: 5px solid #D48806; }
.bento-card--ready { border-left: 5px solid #1E8A6E; opacity: 0.75; }
.bento-card--ready:hover { opacity: 1; }
.bento-card--ready .card-top--ready { background: #D1FAE5; }

@media (min-width: 1536px) { .bento-grid { grid-template-columns: repeat(4, 1fr); } }
@media (min-width: 1280px) and (max-width: 1535px) { .bento-grid { grid-template-columns: repeat(3, 1fr); } }
@media (min-width: 768px) and (max-width: 1279px) { .bento-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 767px) {
  .cooking-header { flex-direction: column; gap: 12px; }
  .header-controls { flex-wrap: wrap; gap: 8px; width: 100%; }
  .bento-grid { grid-template-columns: 1fr; gap: 16px; }
  .card-name { font-size: 16px; }
  .portions-number { font-size: 24px; }
  .card-portions { min-width: 64px; height: 64px; }
}
</style>
