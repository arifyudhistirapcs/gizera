<template>
  <div class="risk-form-container">
    <!-- Navigation Bar -->
    <van-nav-bar
      title="Form Risk Assessment"
      left-arrow
      @click-left="goBack"
      fixed
    />

    <!-- Loading State -->
    <van-loading v-if="loading" type="spinner" vertical class="loading-state">
      Memuat formulir...
    </van-loading>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <van-icon name="warning-o" size="48" color="#ee0a24" />
      <p class="error-message">{{ error }}</p>
      <van-button type="primary" size="small" @click="loadForm">
        Coba Lagi
      </van-button>
    </div>

    <!-- Form Content -->
    <div v-else-if="form" class="form-content">
      <!-- Form Header Info -->
      <van-cell-group inset class="info-card">
        <van-cell title="SPPG" :value="form.sppg?.nama || form.sppg?.sppg_nama || '-'" />
        <van-cell title="Tanggal" :value="formatDate(form.created_at)" />
        <van-cell title="Status">
          <template #value>
            <van-tag :type="form.status === 'submitted' ? 'success' : 'warning'" size="medium">
              {{ form.status === 'submitted' ? 'Submitted' : 'Draft' }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Progress Indicator -->
      <van-cell-group inset class="info-card" v-if="form.status === 'draft'">
        <van-cell title="Progres Penilaian" :value="`${scoredCount} / ${totalCount} item`" />
        <div class="progress-bar-wrapper">
          <van-progress
            :percentage="progressPercentage"
            :color="progressPercentage === 100 ? '#05CD99' : '#303030'"
            stroke-width="8"
          />
        </div>
      </van-cell-group>

      <!-- Submitted Score Summary -->
      <van-cell-group inset class="info-card" v-if="form.status === 'submitted'">
        <van-cell title="Skor Risiko" class="section-title" />
        <van-cell title="Overall Risk Score" :value="form.overall_risk_score?.toFixed(2) || '-'" />
        <van-cell title="Risk Level">
          <template #value>
            <van-tag :type="riskLevelTagType(form.risk_level)" size="medium">
              {{ riskLevelLabel(form.risk_level) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- Auto-save indicator -->
      <div v-if="saving" class="autosave-indicator">
        <van-loading size="14" />
        <span>Menyimpan...</span>
      </div>
      <div v-else-if="lastSaved" class="autosave-indicator saved">
        <van-icon name="success" size="14" />
        <span>Tersimpan</span>
      </div>

      <!-- Sync status indicator -->
      <div v-if="!online" class="sync-status-bar offline">
        <van-icon name="warning-o" size="14" />
        <span>Offline — perubahan disimpan lokal</span>
      </div>
      <div v-else-if="pendingSyncCount > 0" class="sync-status-bar pending">
        <van-icon name="replay" size="14" />
        <span>{{ pendingSyncCount }} draft menunggu sinkronisasi</span>
      </div>

      <!-- Checklist Items grouped by Category -->
      <van-collapse v-model="activeCategories" class="category-collapse">
        <van-collapse-item
          v-for="category in groupedItems"
          :key="category.name"
          :name="category.name"
          :title="category.name"
        >
          <template #title>
            <div class="category-header">
              <span class="category-title">{{ category.name }}</span>
              <van-tag size="small" :type="categoryProgressType(category)">
                {{ categoryScoredCount(category) }}/{{ category.items.length }}
              </van-tag>
            </div>
          </template>

          <div
            v-for="item in category.items"
            :key="item.id"
            class="checklist-item"
          >
            <div class="item-header">
              <span class="item-name">{{ item.item_nama }}</span>
              <van-icon
                v-if="item.compliance_score"
                name="success"
                color="#05CD99"
                size="16"
              />
            </div>

            <!-- Score Selector (1-5 stepper) -->
            <div class="score-section">
              <span class="score-label">Skor Kepatuhan (1-5):</span>
              <van-stepper
                :model-value="item.compliance_score || 0"
                @update:model-value="(val) => onScoreChange(item, val)"
                :min="0"
                :max="5"
                :step="1"
                theme="round"
                :disabled="form.status === 'submitted'"
                button-size="28"
              />
              <span class="score-text" v-if="item.compliance_score">
                {{ scoreLabel(item.compliance_score) }}
              </span>
            </div>

            <!-- Notes Input -->
            <van-field
              :model-value="item.catatan || ''"
              @update:model-value="(val) => onNotesChange(item, val)"
              type="textarea"
              rows="2"
              placeholder="Catatan (opsional)"
              :disabled="form.status === 'submitted'"
              class="notes-field"
            />

            <!-- Evidence Camera Capture -->
            <div class="evidence-section">
              <template v-if="form.status === 'draft'">
                <van-button
                  size="small"
                  icon="photograph"
                  :loading="uploadingItemId === item.id"
                  :disabled="uploadingItemId === item.id"
                  @click="openCamera(item)"
                >
                  {{ item.evidence_url ? 'Foto Ulang' : 'Ambil Foto' }}
                </van-button>
              </template>
              <div v-if="item.evidence_url" class="evidence-preview">
                <van-image
                  :src="resolveUrl(item.evidence_url)"
                  width="80"
                  height="80"
                  fit="cover"
                  radius="8"
                  @click="previewImage(resolveUrl(item.evidence_url))"
                />
              </div>
            </div>
          </div>
        </van-collapse-item>
      </van-collapse>

      <!-- Operational Snapshot Section -->
      <van-cell-group inset class="info-card snapshot-section" v-if="form.snapshot">
        <van-cell title="Data Operasional SPPG Saat Audit" class="section-title" />

        <!-- Review Metrics -->
        <van-cell title="Rating Review">
          <template #value>
            <span :class="metricClass('review', form.snapshot.average_overall_rating)">
              {{ form.snapshot.average_overall_rating?.toFixed(1) || '0.0' }}
            </span>
          </template>
        </van-cell>
        <van-cell title="Total Ulasan" :value="form.snapshot.total_reviews || 0" />

        <!-- Financial Metrics -->
        <van-cell title="Penyerapan Anggaran">
          <template #value>
            <span :class="metricClass('budget', form.snapshot.budget_absorption_rate)">
              {{ (form.snapshot.budget_absorption_rate || 0).toFixed(1) }}%
            </span>
          </template>
        </van-cell>

        <!-- Delivery Metrics -->
        <van-cell title="On-Time Delivery">
          <template #value>
            <span :class="metricClass('delivery', form.snapshot.on_time_delivery_rate)">
              {{ (form.snapshot.on_time_delivery_rate || 0).toFixed(1) }}%
            </span>
          </template>
        </van-cell>
        <van-cell title="Pengiriman Selesai" :value="`${form.snapshot.completed_deliveries || 0} / ${form.snapshot.total_deliveries || 0}`" />

        <!-- Production Metrics -->
        <van-cell title="Porsi Diproduksi" :value="form.snapshot.total_portions_produced || 0" />
        <van-cell title="Completion Rate">
          <template #value>
            <span :class="metricClass('completion', form.snapshot.production_completion_rate)">
              {{ (form.snapshot.production_completion_rate || 0).toFixed(1) }}%
            </span>
          </template>
        </van-cell>

        <!-- Inventory Metrics -->
        <van-cell title="Stok Kritis">
          <template #value>
            <span :class="metricClass('stock', form.snapshot.critical_stock_items)">
              {{ form.snapshot.critical_stock_items || 0 }} item
            </span>
          </template>
        </van-cell>

        <!-- HR Metrics -->
        <van-cell title="Karyawan Aktif" :value="form.snapshot.total_active_employees || 0" />
        <van-cell title="Kehadiran">
          <template #value>
            <span :class="metricClass('attendance', form.snapshot.attendance_rate)">
              {{ (form.snapshot.attendance_rate || 0).toFixed(1) }}%
            </span>
          </template>
        </van-cell>

        <van-cell
          v-if="form.snapshot.captured_at"
          title="Data diambil"
          :value="formatDate(form.snapshot.captured_at)"
        />
      </van-cell-group>

      <!-- Submit Button -->
      <div class="action-buttons" v-if="form.status === 'draft'">
        <van-button
          type="primary"
          block
          size="large"
          @click="submitForm"
          :loading="submitting"
          :disabled="!canSubmit"
        >
          Submit Penilaian
        </van-button>
        <p class="submit-hint" v-if="!canSubmit">
          Semua item harus memiliki skor sebelum submit
        </p>
      </div>
    </div>

    <!-- Image Preview -->
    <van-image-preview
      v-model:show="showImagePreview"
      :images="previewImages"
    />

    <!-- Camera Overlay -->
    <van-overlay :show="cameraActive" z-index="2000" @click="closeCamera">
      <div class="camera-overlay" @click.stop>
        <div class="camera-header">
          <span>Ambil Foto Evidence</span>
          <van-icon name="cross" size="24" color="#fff" @click="closeCamera" />
        </div>
        <div class="camera-viewfinder">
          <video ref="videoRef" autoplay playsinline muted class="camera-video"></video>
          <canvas ref="captureCanvas" style="display:none"></canvas>
        </div>
        <div class="camera-footer">
          <div class="camera-shutter" @click="capturePhoto">
            <div class="shutter-inner"></div>
          </div>
        </div>
      </div>
    </van-overlay>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showSuccessToast, showFailToast, showConfirmDialog } from 'vant'
import riskAssessmentService from '@/services/riskAssessmentService'
import {
  isOnline,
  saveDraft as saveDraftOffline,
  saveDraftLocally,
  loadDraftLocally,
  cacheResponse,
  getCachedResponse,
  setupOfflineSync,
  getPendingSyncCount
} from '@/utils/riskAssessmentOffline'

const router = useRouter()
const route = useRoute()

// Helper to resolve backend URLs for uploaded files
const backendBaseUrl = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1').replace('/api/v1', '')
function resolveUrl(url) {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return backendBaseUrl + url
}

// State
const loading = ref(false)
const error = ref(null)
const form = ref(null)
const saving = ref(false)
const lastSaved = ref(false)
const submitting = ref(false)
const uploadingItemId = ref(null)
const activeCategories = ref([])
const showImagePreview = ref(false)
const previewImages = ref([])
const online = ref(navigator.onLine)
const pendingSyncCount = ref(0)

// Cleanup for offline sync listener
let cleanupOfflineSync = null

// Debounce timer
let saveTimer = null

// ==================== Computed ====================

const groupedItems = computed(() => {
  if (!form.value?.items) return []
  const groups = {}
  for (const item of form.value.items) {
    const cat = item.category_nama || 'Lainnya'
    if (!groups[cat]) {
      groups[cat] = { name: cat, items: [] }
    }
    groups[cat].items.push(item)
  }
  // Sort items within each group by id to maintain order
  const result = Object.values(groups)
  result.forEach(g => g.items.sort((a, b) => a.id - b.id))
  return result
})

const totalCount = computed(() => form.value?.items?.length || 0)

const scoredCount = computed(() => {
  if (!form.value?.items) return 0
  return form.value.items.filter(i => i.compliance_score && i.compliance_score > 0).length
})

const progressPercentage = computed(() => {
  if (totalCount.value === 0) return 0
  return Math.round((scoredCount.value / totalCount.value) * 100)
})

const canSubmit = computed(() => {
  return scoredCount.value === totalCount.value && totalCount.value > 0
})

// ==================== Methods ====================

function goBack() {
  router.back()
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function scoreLabel(score) {
  const labels = { 1: 'Tidak Patuh', 2: 'Kurang', 3: 'Cukup', 4: 'Baik', 5: 'Sangat Patuh' }
  return labels[score] || ''
}

function riskLevelTagType(level) {
  if (level === 'rendah') return 'success'
  if (level === 'sedang') return 'warning'
  if (level === 'tinggi') return 'danger'
  return 'default'
}

function riskLevelLabel(level) {
  if (level === 'rendah') return 'Risiko Rendah'
  if (level === 'sedang') return 'Risiko Sedang'
  if (level === 'tinggi') return 'Risiko Tinggi'
  return level || '-'
}

function categoryProgressType(category) {
  const scored = categoryScoredCount(category)
  if (scored === category.items.length) return 'success'
  if (scored > 0) return 'warning'
  return 'default'
}

function categoryScoredCount(category) {
  return category.items.filter(i => i.compliance_score && i.compliance_score > 0).length
}

// Metric color classes based on design thresholds
function metricClass(type, value) {
  switch (type) {
    case 'review':
      if (value >= 4.0) return 'metric-green'
      if (value >= 3.0) return 'metric-yellow'
      return 'metric-red'
    case 'budget':
      if (value >= 80) return 'metric-green'
      if (value >= 50) return 'metric-yellow'
      return 'metric-red'
    case 'delivery':
      if (value >= 90) return 'metric-green'
      if (value >= 70) return 'metric-yellow'
      return 'metric-red'
    case 'completion':
      if (value >= 90) return 'metric-green'
      if (value >= 70) return 'metric-yellow'
      return 'metric-red'
    case 'stock':
      if (value === 0) return 'metric-green'
      if (value <= 3) return 'metric-yellow'
      return 'metric-red'
    case 'attendance':
      if (value >= 90) return 'metric-green'
      if (value >= 75) return 'metric-yellow'
      return 'metric-red'
    default:
      return ''
  }
}

function previewImage(url) {
  previewImages.value = [url]
  showImagePreview.value = true
}

// ==================== Score & Notes Handlers ====================

function onScoreChange(item, val) {
  // val=0 means "not scored" (stepper min is 0)
  item.compliance_score = val > 0 ? val : null
  debouncedSave()
}

function onNotesChange(item, val) {
  item.catatan = val
  debouncedSave()
}

function debouncedSave() {
  lastSaved.value = false
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    saveDraft()
  }, 2000)
}

async function saveDraft() {
  if (!form.value || form.value.status !== 'draft') return
  saving.value = true
  try {
    const items = form.value.items.map(i => ({
      item_id: i.id,
      compliance_score: i.compliance_score,
      catatan: i.catatan || ''
    }))
    const { savedOnline } = await saveDraftOffline(form.value.id, items, form.value)
    lastSaved.value = true
    if (!savedOnline) {
      pendingSyncCount.value = await getPendingSyncCount()
    }
    // Hide "saved" indicator after 3 seconds
    setTimeout(() => { lastSaved.value = false }, 3000)
  } catch (err) {
    console.error('Auto-save failed:', err)
    showToast('Gagal menyimpan draft')
  } finally {
    saving.value = false
  }
}

// ==================== Camera Capture & Evidence ====================

const videoRef = ref(null)
const captureCanvas = ref(null)
const cameraActive = ref(false)
let cameraStream = null
let cameraTargetItem = null

function getCurrentPosition() {
  return new Promise((resolve) => {
    if (!navigator.geolocation) { resolve(null); return }
    navigator.geolocation.getCurrentPosition(
      (pos) => resolve({ lat: pos.coords.latitude, lng: pos.coords.longitude }),
      () => resolve(null),
      { enableHighAccuracy: true, timeout: 8000 }
    )
  })
}

async function openCamera(item) {
  cameraTargetItem = item
  cameraActive.value = true
  try {
    cameraStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment', width: { ideal: 1280 }, height: { ideal: 960 } },
      audio: false
    })
    // Wait for DOM to render video element
    await new Promise(r => setTimeout(r, 100))
    if (videoRef.value) {
      videoRef.value.srcObject = cameraStream
    }
  } catch (err) {
    console.error('Camera access failed:', err)
    showFailToast('Tidak dapat mengakses kamera')
    closeCamera()
  }
}

function closeCamera() {
  cameraActive.value = false
  if (cameraStream) {
    cameraStream.getTracks().forEach(t => t.stop())
    cameraStream = null
  }
  cameraTargetItem = null
}

async function capturePhoto() {
  if (!videoRef.value || !cameraTargetItem) return
  const item = cameraTargetItem
  const video = videoRef.value

  // Capture frame to canvas
  const canvas = captureCanvas.value
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight
  const ctx = canvas.getContext('2d')
  ctx.drawImage(video, 0, 0)

  // Close camera immediately
  closeCamera()

  uploadingItemId.value = item.id
  try {
    // Get GPS
    const coords = await getCurrentPosition()

    // Stamp watermark
    const w = canvas.width
    const h = canvas.height
    const barH = Math.round(h * 0.06)
    const fontSize = Math.round(barH * 0.38)
    const smallFont = Math.round(barH * 0.3)

    ctx.fillStyle = 'rgba(0, 0, 0, 0.55)'
    ctx.fillRect(0, h - barH, w, barH)

    const now = new Date()
    const pad = (n) => String(n).padStart(2, '0')
    const ts = `${pad(now.getDate())}/${pad(now.getMonth() + 1)}/${now.getFullYear()} ${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`

    ctx.fillStyle = '#fff'
    ctx.font = `bold ${fontSize}px sans-serif`
    ctx.textBaseline = 'middle'
    ctx.textAlign = 'left'
    ctx.fillText(ts, 12, h - barH * 0.62)

    if (coords) {
      ctx.font = `${smallFont}px sans-serif`
      ctx.fillText(`📍 ${coords.lat.toFixed(6)}, ${coords.lng.toFixed(6)}`, 12, h - barH * 0.25)
    }

    ctx.font = `${smallFont}px sans-serif`
    ctx.textAlign = 'right'
    ctx.fillText('Dapur Sehat - Audit', w - 12, h - barH * 0.5)

    // Convert to blob and upload
    const blob = await new Promise(r => canvas.toBlob(r, 'image/jpeg', 0.85))

    const formData = new FormData()
    formData.append('photo', blob, `evidence-${Date.now()}.jpg`)
    formData.append('item_id', item.id)

    const response = await riskAssessmentService.uploadEvidence(form.value.id, formData)
    const data = response.data
    if (data?.success && data.data?.evidence_url) {
      item.evidence_url = data.data.evidence_url
      showSuccessToast('Foto berhasil diambil')
    } else {
      throw new Error(data?.message || 'Upload gagal')
    }
  } catch (err) {
    console.error('Evidence capture failed:', err)
    showFailToast(err.response?.data?.message || 'Gagal mengambil foto')
  } finally {
    uploadingItemId.value = null
  }
}

// ==================== Submit ====================

async function submitForm() {
  if (!canSubmit.value) {
    showToast('Semua item harus memiliki skor')
    return
  }

  try {
    await showConfirmDialog({
      title: 'Konfirmasi Submit',
      message: 'Setelah disubmit, formulir tidak dapat diubah lagi. Lanjutkan?'
    })
  } catch {
    return // User cancelled
  }

  // Flush any pending auto-save first
  if (saveTimer) {
    clearTimeout(saveTimer)
    await saveDraft()
  }

  submitting.value = true
  try {
    const response = await riskAssessmentService.submitForm(form.value.id)
    const data = response.data
    if (data?.success) {
      form.value = data.data
      showSuccessToast('Formulir berhasil disubmit')
    } else {
      throw new Error(data?.message || 'Submit gagal')
    }
  } catch (err) {
    console.error('Submit failed:', err)
    showFailToast(err.response?.data?.message || 'Gagal submit formulir')
  } finally {
    submitting.value = false
  }
}

// ==================== Load Form ====================

async function loadForm() {
  const formId = route.params.id
  if (!formId) {
    error.value = 'ID formulir tidak valid'
    return
  }

  loading.value = true
  error.value = null
  try {
    if (isOnline()) {
      const response = await riskAssessmentService.getForm(formId)
      const data = response.data
      if (data?.success && data.data) {
        form.value = data.data
        // Cache for offline use
        await cacheResponse(`form_${formId}`, data.data)
        await saveDraftLocally(data.data)
        // Open all categories by default
        activeCategories.value = groupedItems.value.map(g => g.name)
      } else {
        throw new Error(data?.message || 'Gagal memuat formulir')
      }
    } else {
      // Offline — try loading from local draft first, then cache
      const localDraft = await loadDraftLocally(Number(formId))
      const cached = localDraft || await getCachedResponse(`form_${formId}`)
      if (cached) {
        form.value = cached
        activeCategories.value = groupedItems.value.map(g => g.name)
        showToast('Memuat data offline')
      } else {
        throw new Error('Tidak ada data offline tersedia')
      }
    }
  } catch (err) {
    console.error('Error loading form:', err)
    // If online load failed, try offline fallback
    if (isOnline()) {
      const localDraft = await loadDraftLocally(Number(formId))
      const cached = localDraft || await getCachedResponse(`form_${formId}`)
      if (cached) {
        form.value = cached
        activeCategories.value = groupedItems.value.map(g => g.name)
        error.value = null
        showToast('Memuat data dari cache lokal')
        return
      }
    }
    error.value = err.response?.data?.message || err.message || 'Gagal memuat formulir'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadForm()

  // Track online/offline status
  const setOnline = () => { online.value = true }
  const setOffline = () => { online.value = false }
  window.addEventListener('online', setOnline)
  window.addEventListener('offline', setOffline)

  // Setup background sync when coming back online
  cleanupOfflineSync = setupOfflineSync(async (result) => {
    pendingSyncCount.value = await getPendingSyncCount()
    if (result.synced > 0) {
      showSuccessToast(`${result.synced} draft berhasil disinkronkan`)
    }
  })

  // Store cleanup refs for onUnmounted
  window.__raOfflineCleanup = { setOnline, setOffline }
})

onUnmounted(() => {
  if (cleanupOfflineSync) cleanupOfflineSync()
  const c = window.__raOfflineCleanup
  if (c) {
    window.removeEventListener('online', c.setOnline)
    window.removeEventListener('offline', c.setOffline)
    delete window.__raOfflineCleanup
  }
})
</script>

<style scoped>
.risk-form-container {
  min-height: 100vh;
  background-color: #E8EDE5;
  padding-top: 46px;
  padding-bottom: 80px;
}

.loading-state {
  padding-top: 120px;
  text-align: center;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 24px;
  text-align: center;
}

.error-message {
  font-size: 14px;
  color: #666;
  margin: 16px 0;
}

.form-content {
  padding: 16px;
}

.info-card {
  margin-bottom: 16px;
  border-radius: 16px !important;
  overflow: hidden;
}

.section-title {
  font-weight: bold;
  background-color: #303030;
  color: white;
}

.section-title :deep(.van-cell__title) {
  color: white;
}

.progress-bar-wrapper {
  padding: 12px 16px 16px;
}

/* Auto-save indicator */
.autosave-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #999;
  padding: 0 4px 8px;
}

.autosave-indicator.saved {
  color: #05CD99;
}

/* Sync status bar */
.sync-status-bar {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 6px 12px;
  border-radius: 8px;
  margin-bottom: 8px;
}

.sync-status-bar.offline {
  background-color: #fff1f0;
  color: #ee0a24;
}

.sync-status-bar.pending {
  background-color: #fffbe8;
  color: #ed6a0c;
}

/* Category collapse */
.category-collapse {
  margin-bottom: 16px;
}

.category-collapse :deep(.van-collapse-item) {
  margin-bottom: 8px;
  border-radius: 12px;
  overflow: hidden;
}

.category-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding-right: 8px;
}

.category-title {
  font-size: 14px;
  font-weight: 600;
  color: #303030;
  flex: 1;
  margin-right: 8px;
}

/* Checklist item */
.checklist-item {
  padding: 12px 0;
  border-bottom: 1px solid #f0f0f0;
}

.checklist-item:last-child {
  border-bottom: none;
}

.item-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 8px;
}

.item-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
  flex: 1;
  margin-right: 8px;
}

/* Score section */
.score-section {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.score-label {
  font-size: 12px;
  color: #666;
}

.score-text {
  font-size: 12px;
  color: #303030;
  font-weight: 500;
}

.notes-field {
  margin-bottom: 8px;
  border-radius: 8px;
}

/* Evidence section */
.evidence-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.evidence-preview {
  flex-shrink: 0;
}

/* Snapshot section */
.snapshot-section {
  margin-top: 16px;
}

/* Metric color indicators */
.metric-green {
  color: #05CD99;
  font-weight: 600;
}

.metric-yellow {
  color: #FFB547;
  font-weight: 600;
}

.metric-red {
  color: #EE5D50;
  font-weight: 600;
}

/* Action buttons */
.action-buttons {
  margin-top: 24px;
  padding-bottom: 20px;
}

.submit-hint {
  font-size: 12px;
  color: #999;
  text-align: center;
  margin-top: 8px;
}

/* Vant overrides */
:deep(.van-stepper--round .van-stepper__minus),
:deep(.van-stepper--round .van-stepper__plus) {
  background-color: #303030;
}

:deep(.van-progress__pivot) {
  font-size: 11px;
}

/* Camera Overlay */
.camera-overlay {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #000;
}

.camera-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.camera-viewfinder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.camera-video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.camera-footer {
  padding: 24px;
  display: flex;
  justify-content: center;
  flex-shrink: 0;
}

.camera-shutter {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  border: 4px solid #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

.camera-shutter:active .shutter-inner {
  transform: scale(0.85);
}

.shutter-inner {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.15s;
}
</style>
