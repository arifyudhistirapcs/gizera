<template>
  <div class="risk-assessment-select">
    <van-nav-bar
      title="Risk Assessment"
      left-arrow
      @click-left="goBack"
    />

    <van-loading v-if="loading" type="spinner" vertical class="loading-state">
      Memuat data...
    </van-loading>

    <div v-else class="content">
      <!-- Draft Forms -->
      <div v-if="draftForms.length" class="section">
        <p class="section-title">Draft — Lanjutkan Pengisian</p>
        <van-swipe-cell v-for="form in draftForms" :key="form.id">
          <van-cell-group inset class="card" style="margin-bottom: 0">
            <van-cell
              :title="form.sppg?.nama || `SPPG #${form.sppg_id}`"
              :label="formatDate(form.created_at)"
              is-link
              @click="router.push(`/risk-assessment/${form.id}`)"
            >
              <template #value>
                <van-tag type="warning" size="medium">Draft</van-tag>
              </template>
            </van-cell>
          </van-cell-group>
          <template #right>
            <van-button square type="danger" text="Hapus" class="swipe-delete-btn" @click="deleteDraft(form)" />
          </template>
        </van-swipe-cell>
      </div>

      <!-- New Audit -->
      <div class="section">
        <p class="section-title">Audit Baru</p>
        <p class="section-desc">Pilih SPPG yang akan diaudit</p>
        <van-cell-group inset class="card">
          <van-cell
            v-for="sppg in sppgList"
            :key="sppg.id"
            :title="sppg.nama || `SPPG #${sppg.id}`"
            :label="sppg.kode || ''"
            is-link
            @click="selectSPPG(sppg)"
          >
            <template #right-icon>
              <van-loading v-if="creatingSppgId === sppg.id" size="20" />
              <van-icon v-else name="add-o" color="#5A4372" size="20" />
            </template>
          </van-cell>
        </van-cell-group>
        <div v-if="!sppgList.length" class="empty-inline">
          <van-icon name="info-o" size="16" color="#999" />
          <span>Tidak ada SPPG ditemukan</span>
        </div>
      </div>

      <!-- Submitted History -->
      <div v-if="submittedForms.length" class="section">
        <p class="section-title">Riwayat Audit</p>
        <van-cell-group inset class="card">
          <van-cell
            v-for="form in submittedForms"
            :key="form.id"
            :title="form.sppg?.nama || `SPPG #${form.sppg_id}`"
            :label="formatDate(form.submitted_at || form.created_at)"
            is-link
            @click="router.push(`/risk-assessment/${form.id}`)"
          >
            <template #value>
              <div class="history-value">
                <span class="score" :style="{ color: scoreColor(form.overall_risk_score) }">
                  {{ form.overall_risk_score != null ? form.overall_risk_score.toFixed(1) : '-' }}
                </span>
                <van-tag :type="riskTagType(form.risk_level)" size="small">
                  {{ form.risk_level || '-' }}
                </van-tag>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
      </div>

      <div v-if="!draftForms.length && !submittedForms.length && !sppgList.length" class="empty-state">
        <van-icon name="search" size="48" color="#999" />
        <p class="empty-message">Belum ada data</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showFailToast, showConfirmDialog } from 'vant'
import riskAssessmentService from '@/services/riskAssessmentService'

const router = useRouter()
const loading = ref(false)
const sppgList = ref([])
const forms = ref([])
const creatingSppgId = ref(null)

const draftForms = computed(() => forms.value.filter(f => f.status === 'draft'))
const submittedForms = computed(() => forms.value.filter(f => f.status === 'submitted'))

function goBack() {
  router.back()
}

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function scoreColor(score) {
  if (score == null) return '#999'
  if (score >= 4.0) return '#05CD99'
  if (score >= 2.5) return '#FFB547'
  return '#EE5D50'
}

function riskTagType(level) {
  if (level === 'rendah') return 'success'
  if (level === 'sedang') return 'warning'
  if (level === 'tinggi') return 'danger'
  return 'default'
}

async function fetchData() {
  loading.value = true
  try {
    const [sppgRes, formsRes] = await Promise.all([
      riskAssessmentService.getSPPGList(),
      riskAssessmentService.getForms({ page_size: 50 })
    ])
    sppgList.value = sppgRes.data?.data || []
    forms.value = formsRes.data?.data || []
  } catch (err) {
    console.error('Error loading data:', err)
  } finally {
    loading.value = false
  }
}

async function deleteDraft(form) {
  try {
    await showConfirmDialog({ title: 'Hapus Draft', message: `Hapus draft audit untuk ${form.sppg?.nama || 'SPPG'}?` })
  } catch { return }
  try {
    await riskAssessmentService.deleteForm(form.id)
    showToast('Draft berhasil dihapus')
    forms.value = forms.value.filter(f => f.id !== form.id)
  } catch (err) {
    showFailToast(err.response?.data?.message || 'Gagal menghapus draft')
  }
}

async function selectSPPG(sppg) {
  if (creatingSppgId.value) return
  creatingSppgId.value = sppg.id
  try {
    const response = await riskAssessmentService.createForm({ sppg_id: sppg.id })
    const data = response.data
    if (data?.success && data.data?.id) {
      showToast('Formulir berhasil dibuat')
      router.push(`/risk-assessment/${data.data.id}`)
    } else {
      throw new Error(data?.message || 'Gagal membuat formulir')
    }
  } catch (err) {
    showFailToast(err.response?.data?.message || 'Gagal membuat formulir')
  } finally {
    creatingSppgId.value = null
  }
}

onMounted(() => fetchData())
</script>

<style scoped>
.risk-assessment-select {
  min-height: 100vh;
  background-color: #F8FDEA;
}

.loading-state {
  padding-top: 120px;
  text-align: center;
}

.content {
  padding: 16px;
  padding-bottom: 80px;
}

.section {
  margin-bottom: 20px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: #322837;
  margin: 0 0 4px;
}

.section-desc {
  font-size: 13px;
  color: #666;
  margin: 0 0 10px;
}

.card {
  border-radius: 16px !important;
  overflow: hidden;
}

.history-value {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.score {
  font-size: 18px;
  font-weight: 700;
}

.swipe-delete-btn {
  height: 100%;
}

.empty-inline {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 4px;
  font-size: 13px;
  color: #999;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 24px;
}

.empty-message {
  font-size: 14px;
  color: #999;
  margin: 12px 0;
}
</style>
