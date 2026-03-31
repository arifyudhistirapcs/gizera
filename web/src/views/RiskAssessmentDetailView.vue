<template>
  <div class="risk-assessment-detail">
    <a-page-header
      title="Detail Risk Assessment"
      sub-title="Hasil Audit Kepatuhan SOP"
      @back="$router.push('/risk-assessment')"
    />
    <a-spin :spinning="loading" tip="Memuat data...">
      <template v-if="form">
        <a-card class="detail-card">
          <a-descriptions title="Informasi Audit" bordered :column="{ xs: 1, sm: 2, md: 3 }">
            <a-descriptions-item label="SPPG">
              {{ form.sppg?.nama || `SPPG #${form.sppg_id}` }}
            </a-descriptions-item>
            <a-descriptions-item label="Tanggal Audit">
              {{ formatDate(form.created_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="Status">
              <a-tag :color="form.status === 'submitted' ? 'blue' : 'default'">{{ form.status }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="Skor Risiko">
              <span :style="{ color: scoreColor(form.overall_risk_score), fontWeight: 'bold', fontSize: '18px' }">
                {{ form.overall_risk_score != null ? form.overall_risk_score.toFixed(1) : '-' }}
              </span>
            </a-descriptions-item>
            <a-descriptions-item label="Risk Level">
              <a-tag :color="riskLevelColor(form.risk_level)" size="large">{{ form.risk_level || '-' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="Tanggal Submit">
              {{ form.submitted_at ? formatDate(form.submitted_at) : '-' }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card title="Skor per Kategori SOP" class="detail-card" v-if="form.category_scores?.length">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :md="8" v-for="cs in form.category_scores" :key="cs.id">
              <a-card size="small" class="category-score-card">
                <div class="category-name">{{ cs.category_nama }}</div>
                <a-progress :percent="(cs.average_score / 5) * 100" :stroke-color="scoreColor(cs.average_score)" :format="() => cs.average_score.toFixed(1)" />
                <a-tag :color="riskLevelColor(cs.risk_level)" size="small" style="margin-top: 8px">{{ cs.risk_level }}</a-tag>
                <div class="category-item-count">{{ cs.item_count }} item</div>
              </a-card>
            </a-col>
          </a-row>
        </a-card>

        <a-card title="Detail Item Checklist" class="detail-card">
          <a-collapse v-model:activeKey="activeCategories">
            <a-collapse-panel v-for="(items, categoryName) in groupedItems" :key="categoryName" :header="categoryName">
              <template #extra>
                <a-tag :color="getCategoryTagColor(categoryName)">{{ getCategoryAvgScore(categoryName) }}</a-tag>
              </template>
              <a-table :columns="itemColumns" :data-source="items" :pagination="false" row-key="id" size="small">
                <template #bodyCell="{ column, record }">
                  <template v-if="column.key === 'compliance_score'">
                    <a-tag v-if="record.compliance_score != null" :color="itemScoreColor(record.compliance_score)">{{ record.compliance_score }} / 5</a-tag>
                    <span v-else style="color: #999">Belum dinilai</span>
                  </template>
                  <template v-else-if="column.key === 'catatan'">{{ record.catatan || '-' }}</template>
                  <template v-else-if="column.key === 'evidence'">
                    <a-image v-if="record.evidence_url" :src="record.evidence_url" :width="60" :height="60" style="object-fit: cover; border-radius: 4px" />
                    <span v-else style="color: #999">-</span>
                  </template>
                </template>
              </a-table>
            </a-collapse-panel>
          </a-collapse>
        </a-card>

        <a-card title="Data Operasional SPPG Saat Audit" class="detail-card" v-if="form.snapshot">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :md="8" v-for="m in snapshotMetrics" :key="m.key">
              <div class="metric-card">
                <div class="metric-label">{{ m.label }}</div>
                <div class="metric-value">
                  <span :style="{ color: metricHex(m.key, m.value) }">{{ m.display }}</span>
                </div>
                <a-progress v-if="m.showProgress" :percent="Math.round(Math.min(m.value || 0, 100))" :stroke-color="metricHex(m.key, m.value)" size="small" />
                <a-tag :color="metricColor(m.key, m.value)">{{ metricLabel(m.key, m.value) }}</a-tag>
                <div v-if="m.sub" class="metric-sub">{{ m.sub }}</div>
              </div>
            </a-col>
          </a-row>
          <div class="snapshot-period" v-if="form.snapshot.captured_at">
            Periode snapshot: {{ formatDate(form.snapshot.snapshot_period_start) }} &mdash; {{ formatDate(form.snapshot.snapshot_period_end) }}
          </div>
        </a-card>
      </template>
      <a-empty v-else-if="!loading" description="Data tidak ditemukan" />
    </a-spin>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'
import riskAssessmentService from '@/services/riskAssessmentService'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const form = ref(null)
const activeCategories = ref([])

const itemColumns = [
  { title: 'Item', dataIndex: 'item_nama', key: 'item_nama', ellipsis: true },
  { title: 'Skor', key: 'compliance_score', width: 100, align: 'center' },
  { title: 'Catatan', key: 'catatan', width: 250 },
  { title: 'Evidence', key: 'evidence', width: 80, align: 'center' }
]

const groupedItems = computed(() => {
  if (!form.value?.items?.length) return {}
  const groups = {}
  for (const item of form.value.items) {
    const cat = item.category_nama || 'Lainnya'
    if (!groups[cat]) groups[cat] = []
    groups[cat].push(item)
  }
  return groups
})

const snapshotMetrics = computed(() => {
  const s = form.value?.snapshot
  if (!s) return []
  return [
    { key: 'rating', label: 'Rating Review', value: s.average_overall_rating, display: `${round1(s.average_overall_rating)} / 5.0`, sub: `${s.total_reviews || 0} ulasan`, showProgress: false },
    { key: 'budget', label: 'Penyerapan Anggaran', value: clampPct(s.budget_absorption_rate), display: `${round1(s.budget_absorption_rate)}%`, showProgress: true },
    { key: 'ontime', label: 'On-Time Delivery', value: clampPct(s.on_time_delivery_rate), display: `${round1(s.on_time_delivery_rate)}%`, sub: `${s.completed_deliveries || 0}/${s.total_deliveries || 0} pengiriman`, showProgress: true },
    { key: 'completion', label: 'Completion Rate', value: clampPct(s.delivery_completion_rate), display: `${round1(s.delivery_completion_rate)}%`, showProgress: true },
    { key: 'stock', label: 'Stok Kritis', value: s.critical_stock_items, display: `${s.critical_stock_items || 0} item`, sub: `dari ${s.total_inventory_items || 0} total item`, showProgress: false },
    { key: 'attendance', label: 'Kehadiran', value: clampPct(s.attendance_rate), display: `${round1(s.attendance_rate)}%`, sub: `${s.total_active_employees || 0} karyawan aktif`, showProgress: true }
  ]
})

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY HH:mm')
}

const round1 = (v) => Math.round((v || 0) * 10) / 10
const clampPct = (v) => Math.min(Math.max(v || 0, 0), 100)

const riskLevelColor = (level) => {
  if (!level) return 'default'
  switch (level.toLowerCase()) {
    case 'rendah': return 'green'
    case 'sedang': return 'orange'
    case 'tinggi': return 'red'
    default: return 'default'
  }
}

const scoreColor = (score) => {
  if (score == null) return '#999'
  if (score >= 4.0) return '#52c41a'
  if (score >= 2.5) return '#faad14'
  return '#f5222d'
}

const itemScoreColor = (score) => {
  if (score >= 4) return 'green'
  if (score >= 3) return 'orange'
  return 'red'
}

const metricColor = (key, value) => {
  const v = value ?? 0
  switch (key) {
    case 'rating': return v >= 4.0 ? 'green' : v >= 3.0 ? 'orange' : 'red'
    case 'budget': return v >= 80 ? 'green' : v >= 50 ? 'orange' : 'red'
    case 'ontime': return v >= 90 ? 'green' : v >= 70 ? 'orange' : 'red'
    case 'completion': return v >= 90 ? 'green' : v >= 70 ? 'orange' : 'red'
    case 'stock': return v === 0 ? 'green' : v <= 3 ? 'orange' : 'red'
    case 'attendance': return v >= 90 ? 'green' : v >= 75 ? 'orange' : 'red'
    default: return 'default'
  }
}

const metricHex = (key, value) => {
  const c = metricColor(key, value)
  if (c === 'green') return '#52c41a'
  if (c === 'orange') return '#faad14'
  if (c === 'red') return '#f5222d'
  return '#999'
}

const metricLabel = (key, value) => {
  const c = metricColor(key, value)
  if (c === 'green') return 'Baik'
  if (c === 'orange') return 'Perlu Perhatian'
  if (c === 'red') return 'Kritis'
  return '-'
}

const getCategoryTagColor = (categoryName) => {
  const cs = form.value?.category_scores?.find(c => c.category_nama === categoryName)
  return cs ? riskLevelColor(cs.risk_level) : 'default'
}

const getCategoryAvgScore = (categoryName) => {
  const cs = form.value?.category_scores?.find(c => c.category_nama === categoryName)
  return cs ? cs.average_score.toFixed(1) : '-'
}

const fetchForm = async () => {
  loading.value = true
  try {
    const id = route.params.id
    const response = await riskAssessmentService.getForm(id)
    const data = response.data || response
    form.value = data.data || data
    if (form.value?.items) {
      const cats = [...new Set(form.value.items.map(i => i.category_nama || 'Lainnya'))]
      activeCategories.value = cats
    }
  } catch (error) {
    message.error('Gagal memuat detail risk assessment')
    console.error(error)
    if (error.response?.status === 404) {
      router.push('/risk-assessment')
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchForm()
})
</script>

<style scoped>
.risk-assessment-detail {
  padding: 24px;
}

.detail-card {
  margin-bottom: 16px;
}

.category-score-card {
  text-align: center;
}

.category-name {
  font-weight: 600;
  margin-bottom: 8px;
  font-size: 13px;
}

.category-item-count {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.metric-card {
  background: #fafafa;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
}

.metric-label {
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
  font-weight: 500;
}

.metric-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
}

.metric-sub {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.snapshot-period {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #f0f0f0;
  font-size: 12px;
  color: #999;
  text-align: center;
}
</style>
