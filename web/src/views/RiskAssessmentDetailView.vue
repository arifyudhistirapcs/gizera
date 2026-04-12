<template>
  <div class="risk-assessment-detail">
    <a-page-header
      title="Detail Risk Assessment"
      sub-title="Hasil Audit Kepatuhan SOP"
      @back="$router.push('/risk-assessment')"
    >
      <template #extra>
        <a-space>
          <a-button
            type="default"
            :icon="h(FileExcelOutlined)"
            :loading="exportingExcel"
            :disabled="!form"
            @click="exportExcel"
            class="btn-export-excel"
          >
            Export Excel
          </a-button>
          <a-button
            type="default"
            :icon="h(FilePdfOutlined)"
            :loading="exportingPdf"
            :disabled="!form"
            @click="exportPdf"
            class="btn-export-pdf"
          >
            Export PDF
          </a-button>
        </a-space>
      </template>
    </a-page-header>
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
                    <a-image v-if="record.evidence_url" :src="resolveUrl(record.evidence_url)" :width="60" :height="60" style="object-fit: cover; border-radius: 4px" />
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
      <HEmptyState v-else-if="!loading" description="Data tidak ditemukan" />
    </a-spin>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { FileExcelOutlined, FilePdfOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import riskAssessmentService from '@/services/riskAssessmentService'
import HEmptyState from '@/components/common/HEmptyState.vue'

const route = useRoute()
const router = useRouter()

// Resolve backend URL for uploaded files
// In dev: use relative path (proxied by Vite), in prod: use backend base
const resolveUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  // Relative path like /uploads/... will be proxied by Vite dev server
  return url
}

const loading = ref(false)
const form = ref(null)
const activeCategories = ref([])
const exportingExcel = ref(false)
const exportingPdf = ref(false)

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

// --- Export helpers ---

const getExportFileName = (ext) => {
  const sppgName = (form.value?.sppg?.nama || `SPPG-${form.value?.sppg_id || 'unknown'}`)
    .replace(/[^a-zA-Z0-9_-]/g, '-')
  const date = dayjs(form.value?.created_at).format('YYYY-MM-DD')
  return `Risk-Assessment-${sppgName}-${date}.${ext}`
}

// --- Excel Export ---

const exportExcel = async () => {
  exportingExcel.value = true
  try {
    const XLSX = await import('xlsx')
    const wb = XLSX.utils.book_new()

    // Sheet 1: Informasi Audit
    const auditInfo = [
      ['SPPG', form.value.sppg?.nama || `SPPG #${form.value.sppg_id}`],
      ['Tanggal Audit', formatDate(form.value.created_at)],
      ['Status', form.value.status || '-'],
      ['Skor Risiko', form.value.overall_risk_score != null ? form.value.overall_risk_score.toFixed(1) : '-'],
      ['Risk Level', form.value.risk_level || '-'],
      ['Tanggal Submit', form.value.submitted_at ? formatDate(form.value.submitted_at) : '-']
    ]
    const wsAudit = XLSX.utils.aoa_to_sheet([['Field', 'Value'], ...auditInfo])
    wsAudit['!cols'] = [{ wch: 20 }, { wch: 40 }]
    XLSX.utils.book_append_sheet(wb, wsAudit, 'Informasi Audit')

    // Sheet 2: Skor per Kategori
    if (form.value.category_scores?.length) {
      const catHeaders = ['Kategori', 'Rata-rata Skor', 'Risk Level', 'Jumlah Item']
      const catRows = form.value.category_scores.map(cs => [
        cs.category_nama,
        cs.average_score?.toFixed(1) ?? '-',
        cs.risk_level || '-',
        cs.item_count ?? 0
      ])
      const wsCat = XLSX.utils.aoa_to_sheet([catHeaders, ...catRows])
      wsCat['!cols'] = [{ wch: 30 }, { wch: 15 }, { wch: 15 }, { wch: 15 }]
      XLSX.utils.book_append_sheet(wb, wsCat, 'Skor Kategori')
    }

    // Sheet 3: Detail Item Checklist
    if (form.value.items?.length) {
      const itemHeaders = ['Item', 'Kategori', 'Skor Kepatuhan', 'Catatan', 'Evidence URL']
      const itemRows = form.value.items.map(item => [
        item.item_nama || '-',
        item.category_nama || '-',
        item.compliance_score != null ? `${item.compliance_score} / 5` : 'Belum dinilai',
        item.catatan || '-',
        item.evidence_url || '-'
      ])
      const wsItems = XLSX.utils.aoa_to_sheet([itemHeaders, ...itemRows])
      wsItems['!cols'] = [{ wch: 35 }, { wch: 25 }, { wch: 15 }, { wch: 35 }, { wch: 50 }]

      // Add hyperlinks for evidence URLs
      form.value.items.forEach((item, idx) => {
        if (item.evidence_url) {
          const cellRef = XLSX.utils.encode_cell({ r: idx + 1, c: 4 })
          const fullUrl = resolveUrl(item.evidence_url)
          if (!wsItems[cellRef]) wsItems[cellRef] = { t: 's', v: fullUrl }
          wsItems[cellRef].l = { Target: fullUrl, Tooltip: 'Lihat Evidence' }
        }
      })

      XLSX.utils.book_append_sheet(wb, wsItems, 'Detail Checklist')
    }

    // Sheet 4: Data Operasional (Snapshot)
    if (form.value.snapshot) {
      const metrics = snapshotMetrics.value
      const snapHeaders = ['Metrik', 'Nilai', 'Status']
      const snapRows = metrics.map(m => [
        m.label,
        m.display,
        metricLabel(m.key, m.value)
      ])
      const wsSnap = XLSX.utils.aoa_to_sheet([snapHeaders, ...snapRows])
      wsSnap['!cols'] = [{ wch: 25 }, { wch: 20 }, { wch: 20 }]
      XLSX.utils.book_append_sheet(wb, wsSnap, 'Data Operasional')
    }

    XLSX.writeFile(wb, getExportFileName('xlsx'))
    message.success('Export Excel berhasil')
  } catch (error) {
    console.error('Excel export error:', error)
    message.error('Gagal export Excel')
  } finally {
    exportingExcel.value = false
  }
}

// --- PDF Export ---

const loadImageAsBase64 = async (url) => {
  try {
    const response = await fetch(url)
    if (!response.ok) return null
    const blob = await response.blob()
    return new Promise((resolve) => {
      const reader = new FileReader()
      reader.onloadend = () => resolve(reader.result)
      reader.onerror = () => resolve(null)
      reader.readAsDataURL(blob)
    })
  } catch {
    return null
  }
}

const exportPdf = async () => {
  exportingPdf.value = true
  try {
    const jsPDFModule = await import('jspdf')
    const jsPDF = jsPDFModule.default || jsPDFModule.jsPDF
    const { applyPlugin } = await import('jspdf-autotable')
    applyPlugin(jsPDF)

    const doc = new jsPDF('p', 'mm', 'a4')
    const pageWidth = doc.internal.pageSize.getWidth()
    const margin = 14
    let y = 15

    const PRIMARY_RED = [201, 74, 58]
    const GREEN = [30, 138, 110]
    const DARK = [51, 51, 51]
    const GRAY = [128, 128, 128]

    const addSectionTitle = (title) => {
      if (y > 260) { doc.addPage(); y = 15 }
      doc.setFontSize(13)
      doc.setTextColor(...PRIMARY_RED)
      doc.setFont(undefined, 'bold')
      doc.text(title, margin, y)
      y += 2
      doc.setDrawColor(...PRIMARY_RED)
      doc.setLineWidth(0.5)
      doc.line(margin, y, pageWidth - margin, y)
      y += 7
    }

    // --- Header ---
    doc.setFontSize(18)
    doc.setTextColor(...PRIMARY_RED)
    doc.setFont(undefined, 'bold')
    doc.text('Laporan Risk Assessment', pageWidth / 2, y, { align: 'center' })
    y += 7
    doc.setFontSize(10)
    doc.setTextColor(...GRAY)
    doc.setFont(undefined, 'normal')
    doc.text('Audit Kepatuhan SOP - Dapur Sehat', pageWidth / 2, y, { align: 'center' })
    y += 10

    // --- Informasi Audit ---
    addSectionTitle('Informasi Audit')
    const sppgName = form.value.sppg?.nama || `SPPG #${form.value.sppg_id}`

    doc.autoTable({
      startY: y,
      margin: { left: margin, right: margin },
      theme: 'plain',
      styles: { fontSize: 10, cellPadding: 3, textColor: DARK },
      columnStyles: { 0: { fontStyle: 'bold', cellWidth: 45, textColor: GRAY } },
      body: [
        ['SPPG', sppgName],
        ['Tanggal Audit', formatDate(form.value.created_at)],
        ['Status', form.value.status || '-'],
        ['Skor Risiko', form.value.overall_risk_score != null ? form.value.overall_risk_score.toFixed(1) : '-'],
        ['Risk Level', form.value.risk_level || '-'],
        ['Tanggal Submit', form.value.submitted_at ? formatDate(form.value.submitted_at) : '-']
      ]
    })
    y = doc.lastAutoTable.finalY + 10

    // --- Skor per Kategori ---
    if (form.value.category_scores?.length) {
      addSectionTitle('Skor per Kategori SOP')
      doc.autoTable({
        startY: y,
        margin: { left: margin, right: margin },
        headStyles: { fillColor: PRIMARY_RED, textColor: [255, 255, 255], fontStyle: 'bold', fontSize: 10 },
        bodyStyles: { fontSize: 9, textColor: DARK },
        alternateRowStyles: { fillColor: [252, 245, 243] },
        head: [['Kategori', 'Rata-rata Skor', 'Risk Level', 'Jumlah Item']],
        body: form.value.category_scores.map(cs => [
          cs.category_nama,
          cs.average_score?.toFixed(1) ?? '-',
          cs.risk_level || '-',
          cs.item_count ?? 0
        ])
      })
      y = doc.lastAutoTable.finalY + 10
    }

    // --- Detail Item Checklist ---
    if (form.value.items?.length) {
      addSectionTitle('Detail Item Checklist')

      // Pre-load evidence images
      const evidenceImages = {}
      const itemsWithEvidence = form.value.items.filter(item => item.evidence_url)
      console.log(`Loading ${itemsWithEvidence.length} evidence images...`)
      
      const imagePromises = itemsWithEvidence.map(async (item) => {
          const fullUrl = resolveUrl(item.evidence_url)
          const base64 = await loadImageAsBase64(fullUrl)
          if (base64) {
            evidenceImages[item.evidence_url] = base64
            console.log(`✓ Loaded: ${item.item_nama}`)
          } else {
            console.warn(`✗ Failed: ${item.item_nama} (${fullUrl})`)
          }
        })
      await Promise.all(imagePromises)
      console.log(`Loaded ${Object.keys(evidenceImages).length}/${itemsWithEvidence.length} images`)

      doc.autoTable({
        startY: y,
        margin: { left: margin, right: margin },
        headStyles: { fillColor: PRIMARY_RED, textColor: [255, 255, 255], fontStyle: 'bold', fontSize: 9 },
        bodyStyles: { fontSize: 8, textColor: DARK, minCellHeight: 12 },
        alternateRowStyles: { fillColor: [252, 245, 243] },
        columnStyles: {
          0: { cellWidth: 45 },
          1: { cellWidth: 30 },
          2: { cellWidth: 18, halign: 'center' },
          3: { cellWidth: 50 },
          4: { cellWidth: 30, halign: 'center' }
        },
        head: [['Item', 'Kategori', 'Skor', 'Catatan', 'Evidence']],
        body: form.value.items.map(item => [
          item.item_nama || '-',
          item.category_nama || '-',
          item.compliance_score != null ? `${item.compliance_score}/5` : '-',
          item.catatan || '-',
          item.evidence_url ? '' : '-'
        ]),
        didDrawCell: (data) => {
          if (data.section === 'body' && data.column.index === 4) {
            const item = form.value.items[data.row.index]
            if (item?.evidence_url && evidenceImages[item.evidence_url]) {
              const imgData = evidenceImages[item.evidence_url]
              const imgFormat = imgData.startsWith('data:image/png') ? 'PNG' : 'JPEG'
              const imgSize = Math.min(data.cell.height - 2, 20)
              const xPos = data.cell.x + (data.cell.width - imgSize) / 2
              const yPos = data.cell.y + (data.cell.height - imgSize) / 2
              try {
                doc.addImage(imgData, imgFormat, xPos, yPos, imgSize, imgSize)
              } catch (e) {
                console.warn('Failed to add image to PDF:', item.item_nama, e)
              }
            }
          }
        },
        didParseCell: (data) => {
          if (data.section === 'body' && data.column.index === 4) {
            const item = form.value.items[data.row.index]
            if (item?.evidence_url && evidenceImages[item.evidence_url]) {
              data.cell.styles.minCellHeight = 22
            }
          }
        }
      })
      y = doc.lastAutoTable.finalY + 10
    }

    // --- Data Operasional ---
    if (form.value.snapshot) {
      addSectionTitle('Data Operasional SPPG')
      const metrics = snapshotMetrics.value
      doc.autoTable({
        startY: y,
        margin: { left: margin, right: margin },
        headStyles: { fillColor: GREEN, textColor: [255, 255, 255], fontStyle: 'bold', fontSize: 10 },
        bodyStyles: { fontSize: 9, textColor: DARK },
        alternateRowStyles: { fillColor: [240, 249, 245] },
        head: [['Metrik', 'Nilai', 'Status', 'Keterangan']],
        body: metrics.map(m => [
          m.label,
          m.display,
          metricLabel(m.key, m.value),
          m.sub || '-'
        ])
      })
      y = doc.lastAutoTable.finalY + 10

      if (form.value.snapshot.captured_at) {
        if (y > 275) { doc.addPage(); y = 15 }
        doc.setFontSize(8)
        doc.setTextColor(...GRAY)
        doc.text(
          `Periode snapshot: ${formatDate(form.value.snapshot.snapshot_period_start)} — ${formatDate(form.value.snapshot.snapshot_period_end)}`,
          pageWidth / 2, y, { align: 'center' }
        )
        y += 8
      }
    }

    // --- Footer on each page ---
    const totalPages = doc.internal.getNumberOfPages()
    for (let i = 1; i <= totalPages; i++) {
      doc.setPage(i)
      doc.setFontSize(8)
      doc.setTextColor(...GRAY)
      doc.text(
        `Dicetak: ${dayjs().format('DD/MM/YYYY HH:mm')}  |  Halaman ${i} / ${totalPages}`,
        pageWidth / 2,
        doc.internal.pageSize.getHeight() - 8,
        { align: 'center' }
      )
    }

    doc.save(getExportFileName('pdf'))
    message.success('Export PDF berhasil')
  } catch (error) {
    console.error('PDF export error:', error)
    message.error('Gagal export PDF')
  } finally {
    exportingPdf.value = false
  }
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

.btn-export-excel {
  color: #764AF1;
  border-color: #764AF1;
}
.btn-export-excel:hover {
  color: #fff;
  background-color: #764AF1;
  border-color: #764AF1;
}

.btn-export-pdf {
  color: #F82C17;
  border-color: #F82C17;
}
.btn-export-pdf:hover {
  color: #fff;
  background-color: #F82C17;
  border-color: #F82C17;
}
</style>
