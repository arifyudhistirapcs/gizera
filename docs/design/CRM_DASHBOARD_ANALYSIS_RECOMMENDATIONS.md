# Analisa CRM Dashboard Behance → Rekomendasi Gizera ERP SPPG

> Referensi: [CRM Dashboard for SaaS Platform](https://www.behance.net/gallery/241565381/CRM-Dashboard-for-SaaS-Platform-UXUI)
> Tanggal analisa: Juli 2025
> Konteks: Gizera ERP SPPG — Web (Ant Design Vue) + PWA (Vant UI)

---

## 1. Layout & Visual Hierarchy

### Temuan dari CRM Dashboard Behance

CRM dashboard ini menggunakan pendekatan "calm, structured interface" dengan:
- **F-pattern layout**: KPI cards di atas, charts di tengah, detail tables di bawah
- **Progressive disclosure**: Data primer (KPI angka besar) langsung terlihat, data sekunder (breakdown, detail) tersembunyi di bawah atau di-expand
- **Generous whitespace**: Spacing antar section besar (~32px), memberi ruang napas visual
- **Clear section grouping**: Setiap blok data punya header + subtitle yang jelas

### Gap di Gizera Saat Ini

| Aspek | Gizera Sekarang | Masalah |
|-------|----------------|---------|
| Dashboard BGN | `a-card` + `a-statistic` standar Ant Design | Flat, tidak ada visual hierarchy yang kuat |
| Section spacing | `margin-bottom: 20px` seragam | Kurang breathing room, terasa padat |
| Section headers | Hanya `<h3>` atau card title | Tidak ada subtitle/context untuk setiap section |
| Progressive disclosure | Semua data ditampilkan sekaligus | Overwhelming untuk Admin BGN yang monitor nasional |

### Rekomendasi Implementasi

#### 1.1 Tambahkan Section Header Component dengan Subtitle

```vue
<!-- web/src/components/horizon/HSectionHeader.vue -->
<template>
  <div class="h-section-header">
    <div class="h-section-header__left">
      <h3 class="h-section-header__title">
        <slot name="icon" />
        {{ title }}
        <span v-if="badge" class="h-section-header__badge">{{ badge }}</span>
      </h3>
      <p v-if="subtitle" class="h-section-header__subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.action" class="h-section-header__action">
      <slot name="action" />
    </div>
  </div>
</template>

<style scoped>
.h-section-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.h-section-header__title {
  font-size: 18px;
  font-weight: 700;
  color: var(--h-text-primary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}
.h-section-header__subtitle {
  font-size: 13px;
  color: var(--h-text-secondary);
  margin: 4px 0 0 0;
}
.h-section-header__badge {
  font-size: 12px;
  font-weight: 600;
  background: var(--h-primary-lighter, rgba(90, 67, 114, 0.1));
  color: var(--h-primary);
  padding: 2px 10px;
  border-radius: 12px;
}
</style>
```

**Penggunaan di Dashboard:**
```vue
<HSectionHeader
  title="Performa Hari Ini"
  subtitle="Data real-time operasional dapur"
  badge="Live"
>
  <template #action>
    <a-button type="link">Lihat Detail →</a-button>
  </template>
</HSectionHeader>
```

#### 1.2 Perbaiki Section Spacing di Dashboard

```css
/* Tambahkan ke dashboard views */
.dashboard-sspg,
.dashboard-bgn {
  display: flex;
  flex-direction: column;
  gap: 32px; /* Naik dari 20-28px → 32px untuk breathing room */
}

/* Section internal spacing */
.kpi-section,
.activity-section,
.tables-row-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
```

#### 1.3 Progressive Disclosure untuk Dashboard BGN

Dashboard BGN saat ini menampilkan semua tabel sekaligus. Adopsi pattern tabs/collapse:

```vue
<!-- DashboardBGNView.vue — ganti 2 tabel terpisah jadi tabbed view -->
<a-card>
  <template #title>
    <HSectionHeader
      title="Ringkasan Performa"
      subtitle="Drill-down per yayasan atau SPPG"
    />
  </template>
  <a-tabs v-model:activeKey="performanceTab" type="card">
    <a-tab-pane key="yayasan" tab="Per Yayasan">
      <!-- Tabel yayasan existing -->
    </a-tab-pane>
    <a-tab-pane key="sppg" tab="Per SPPG">
      <!-- Tabel SPPG existing -->
    </a-tab-pane>
  </a-tabs>
</a-card>
```

---

## 2. Card Design & Data Presentation

### Temuan dari CRM Dashboard

- **Stat cards** dengan icon gradient di kiri, angka besar di kanan — mirip pattern HStatCard Gizera
- **Mini sparkline/trend** di dalam stat card untuk konteks temporal
- **Card elevation hierarchy**: Primary cards (shadow besar) vs secondary cards (shadow kecil/border only)
- **Color-coded left border** pada cards untuk kategori visual cepat
- **Compact data cards** untuk list items (supplier, customer) dengan progress bar inline

### Gap di Gizera

| Aspek | Gizera Sekarang | Rekomendasi |
|-------|----------------|-------------|
| HStatCard | Sudah bagus, tapi `change` hanya teks | Tambah mini sparkline atau trend arrow |
| Card hierarchy | Semua card pakai `h-card` sama | Bedakan primary vs secondary card |
| Critical stock cards | Custom inline styling | Standardisasi jadi reusable component |
| Supplier cards | Custom inline styling | Standardisasi jadi `HRankCard` component |

### Rekomendasi Implementasi

#### 2.1 Tambah Trend Indicator ke HStatCard

Extend HStatCard dengan prop `trend` untuk mini visual context:

```vue
<!-- Tambah prop di HStatCard.vue -->
props: {
  // ... existing props
  trend: {
    type: String, // 'up' | 'down' | 'stable'
    default: null
  },
  trendValue: {
    type: String, // e.g. "+12% dari kemarin"
    default: ''
  }
}
```

```html
<!-- Tambah di template setelah change indicator -->
<div v-if="trend" class="h-stat-card__trend" :class="`h-stat-card__trend--${trend}`">
  <RiseOutlined v-if="trend === 'up'" />
  <FallOutlined v-if="trend === 'down'" />
  <MinusOutlined v-if="trend === 'stable'" />
  <span>{{ trendValue }}</span>
</div>
```

#### 2.2 Card dengan Color-Coded Left Border

Pattern dari CRM: card dengan border-left berwarna untuk quick visual scanning.

```css
/* Tambahkan utility classes di horizon/utilities.css */
.h-card--accent-left {
  border-left: 4px solid var(--h-primary);
}
.h-card--accent-success {
  border-left: 4px solid var(--h-success);
}
.h-card--accent-warning {
  border-left: 4px solid var(--h-warning);
}
.h-card--accent-error {
  border-left: 4px solid var(--h-error);
}
```

**Penggunaan di Critical Stock:**
```vue
<div class="h-card h-card--accent-error">
  <!-- Stok kritis content -->
</div>
<div class="h-card h-card--accent-warning">
  <!-- Stok rendah content -->
</div>
```

#### 2.3 Standardisasi HAlertCard untuk Critical Items

```vue
<!-- web/src/components/horizon/HAlertCard.vue -->
<template>
  <div class="h-alert-card h-card" :class="`h-alert-card--${severity}`">
    <div class="h-alert-card__indicator"></div>
    <div class="h-alert-card__content">
      <div class="h-alert-card__header">
        <span class="h-alert-card__title">{{ title }}</span>
        <span class="h-alert-card__badge">{{ badge }}</span>
      </div>
      <div class="h-alert-card__body">
        <slot />
      </div>
      <div v-if="$slots.footer" class="h-alert-card__footer">
        <slot name="footer" />
      </div>
    </div>
  </div>
</template>
```

---

## 3. Color System & Status Indicators

### Temuan dari CRM Dashboard

- **Semantic color palette** yang konsisten: hijau (success/completed), kuning/oranye (pending/warning), merah (error/overdue), biru (info/in-progress)
- **Color-coded status badges** dengan dot indicator + background tint (10% opacity)
- **Gradient icon backgrounds** untuk stat cards — memberi depth tanpa overwhelming
- **Muted/tinted backgrounds** (bukan solid color) untuk status — lebih calm dan scannable

### Evaluasi Color System Gizera Saat Ini

Gizera sudah punya foundation yang baik:

| Token | Nilai | Status |
|-------|-------|--------|
| `--h-primary` | `#5A4372` (ungu) | ✅ Baik, distinctive |
| `--h-success` | `#05CD99` (hijau) | ✅ Baik |
| `--h-warning` | `#FFB547` (oranye) | ✅ Baik |
| `--h-error` | `#EE5D50` (merah) | ✅ Baik |
| `--h-info` | `#5A4372` (= primary) | ⚠️ Perlu warna terpisah |
| `--h-bg-primary` | `#F8FDEA` (krem) | ⚠️ Terlalu kuning-hijau, kurang netral |

### Rekomendasi

#### 3.1 Tambah Warna Info Terpisah dari Primary

```css
/* Di variables.css, ganti: */
--h-info: #5A4372; /* Sama dengan primary — ambigu */

/* Menjadi: */
--h-info: #4A90D9;           /* Biru yang distinct dari primary */
--h-info-light: #6BA3E3;
--h-info-dark: #3A7BC8;
```

Ini penting karena di dashboard, "info" (e.g. "Sedang Diproses") dan "primary" (brand) harus bisa dibedakan secara visual.

#### 3.2 Standardisasi Status Badge System

CRM dashboard menggunakan pattern dot + tinted background yang sudah ada di HDataTable. Extend ini jadi reusable standalone component:

```vue
<!-- web/src/components/horizon/HStatusBadge.vue -->
<template>
  <span class="h-status-badge" :class="`h-status-badge--${type}`">
    <span class="h-status-badge__dot"></span>
    <span class="h-status-badge__text">{{ label }}</span>
  </span>
</template>

<script setup>
defineProps({
  label: { type: String, required: true },
  type: {
    type: String,
    default: 'default',
    validator: v => ['success', 'warning', 'error', 'info', 'default', 'processing'].includes(v)
  }
})
</script>

<style scoped>
.h-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
}
.h-status-badge__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
/* Success */
.h-status-badge--success { background: rgba(5, 205, 153, 0.1); color: #04b888; }
.h-status-badge--success .h-status-badge__dot { background: #05CD99; }
/* Warning */
.h-status-badge--warning { background: rgba(255, 181, 71, 0.1); color: #d4940a; }
.h-status-badge--warning .h-status-badge__dot { background: #FFB547; }
/* Error */
.h-status-badge--error { background: rgba(238, 93, 80, 0.1); color: #d43b2e; }
.h-status-badge--error .h-status-badge__dot { background: #EE5D50; }
/* Info */
.h-status-badge--info { background: rgba(74, 144, 217, 0.1); color: #3A7BC8; }
.h-status-badge--info .h-status-badge__dot { background: #4A90D9; }
/* Processing (animated) */
.h-status-badge--processing { background: rgba(90, 67, 114, 0.1); color: #5A4372; }
.h-status-badge--processing .h-status-badge__dot {
  background: #5A4372;
  animation: pulse 1.5s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}
/* Default */
.h-status-badge--default { background: var(--h-bg-light); color: var(--h-text-secondary); }
.h-status-badge--default .h-status-badge__dot { background: var(--h-text-light); }
</style>
```

#### 3.3 Mapping Status Operasional ke Warna

Buat utility function untuk konsistensi di seluruh app:

```js
// web/src/utils/statusColors.js
export const STATUS_MAP = {
  // Produksi
  'Menunggu':        { type: 'default',    label: 'Menunggu' },
  'Sedang Dimasak':  { type: 'processing', label: 'Sedang Dimasak' },
  'Selesai Dimasak': { type: 'success',    label: 'Selesai Dimasak' },
  'Siap Packing':    { type: 'info',       label: 'Siap Packing' },
  'Sedang Packing':  { type: 'processing', label: 'Sedang Packing' },
  'Selesai Packing': { type: 'success',    label: 'Selesai Packing' },

  // Pengiriman
  'Belum Dikirim':   { type: 'default',    label: 'Belum Dikirim' },
  'Dalam Perjalanan':{ type: 'processing', label: 'Dalam Perjalanan' },
  'Terkirim':        { type: 'success',    label: 'Terkirim' },
  'Gagal':           { type: 'error',      label: 'Gagal' },

  // Pencucian
  'Belum Dicuci':    { type: 'default',    label: 'Belum Dicuci' },
  'Sedang Dicuci':   { type: 'processing', label: 'Sedang Dicuci' },
  'Selesai':         { type: 'success',    label: 'Selesai' },

  // PO
  'Draft':           { type: 'default',    label: 'Draft' },
  'Pending':         { type: 'warning',    label: 'Pending' },
  'Approved':        { type: 'success',    label: 'Approved' },
  'Rejected':        { type: 'error',      label: 'Rejected' },

  // Stok
  'Aman':            { type: 'success',    label: 'Aman' },
  'Rendah':          { type: 'warning',    label: 'Rendah' },
  'Kritis':          { type: 'error',      label: 'Kritis' },
}

export function getStatusConfig(status) {
  return STATUS_MAP[status] || { type: 'default', label: status }
}
```

---

## 4. Navigation & Sidebar Patterns

### Temuan dari CRM Dashboard

- **Clean sidebar** dengan icon + label, grouped by category
- **Active state** yang jelas: filled background (bukan hanya text color change)
- **Collapsible groups** dengan smooth animation
- **User profile** di bottom sidebar dengan avatar + role
- **Breadcrumb** di header untuk context awareness

### Evaluasi Gizera Sidebar (HSidebar.vue)

| Aspek | Status | Catatan |
|-------|--------|---------|
| Grouped navigation | ✅ | Sudah ada submenu groups |
| Active state | ✅ | Purple bg + white text — bagus |
| Collapsible | ✅ | Smooth transition 300ms |
| User info | ✅ | Name + role di bottom |
| Role-based filtering | ✅ | Sudah filter per role |
| Breadcrumb | ✅ | Ada di HHeader |
| Section dividers | ⚠️ | Tidak ada visual separator antar group |
| Badge/counter | ❌ | Tidak ada notification badge di menu items |

### Rekomendasi

#### 4.1 Tambah Section Dividers di Sidebar

CRM dashboard memisahkan menu groups dengan label kategori. Gizera bisa adopsi ini:

```vue
<!-- Di HSidebar.vue, tambahkan divider antara multi-tenant dan operational items -->
<template v-if="tenantItems.length > 0 && opItems.length > 0">
  <div class="sidebar-divider">
    <span class="sidebar-divider__label">Operasional</span>
  </div>
</template>
```

```css
.sidebar-divider {
  padding: 16px 20px 8px;
  margin: 8px 12px 0;
}
.sidebar-divider__label {
  font-size: 11px;
  font-weight: 700;
  color: var(--h-text-light);
  text-transform: uppercase;
  letter-spacing: 1px;
}
```

#### 4.2 Tambah Badge Counter untuk Menu Items Kritis

Pattern CRM: badge merah di menu item untuk items yang butuh perhatian.

```vue
<!-- Di menu item template, tambahkan badge -->
<span v-if="item.badge" class="menu-badge">{{ item.badge }}</span>
```

```css
.menu-badge {
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 10px;
  background: var(--h-error);
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: auto;
}
```

Contoh use case: badge di "Manajemen Bahan Baku" menunjukkan jumlah stok kritis.

---

## 5. Table & Data Visualization Patterns

### Temuan dari CRM Dashboard

- **Clean table headers**: uppercase, small font, secondary color — bukan bold dark headers
- **Row hover**: subtle background tint, bukan bold highlight
- **Inline status badges** di dalam tabel (dot + tinted bg)
- **Progress bars** inline di tabel untuk completion rates
- **Sortable columns** dengan visual indicator
- **Compact row height** untuk data-dense views
- **Action buttons** di kolom terakhir: icon-only, grouped

### Evaluasi HDataTable Gizera

| Aspek | Status | Catatan |
|-------|--------|---------|
| Clean headers | ✅ | Uppercase, secondary color, 12px |
| Row hover | ✅ | `#F8FDEA` background |
| Status badges | ✅ | Dot + tinted bg pattern |
| Progress bars | ✅ | Built-in support |
| Mobile card view | ✅ | Responsive card layout |
| Sortable columns | ❌ | Tidak ada sort indicator |
| Row actions | ⚠️ | Slot-based tapi tidak standardized |
| Empty state | ❌ | Tidak ada designed empty state |
| Pagination styling | ✅ | Purple active page |

### Rekomendasi

#### 5.1 Tambah Empty State Component

CRM dashboards selalu punya designed empty state. Gizera perlu ini:

```vue
<!-- web/src/components/horizon/HEmptyState.vue -->
<template>
  <div class="h-empty-state">
    <div class="h-empty-state__icon">
      <slot name="icon">
        <InboxOutlined />
      </slot>
    </div>
    <h4 class="h-empty-state__title">{{ title }}</h4>
    <p v-if="description" class="h-empty-state__description">{{ description }}</p>
    <div v-if="$slots.action" class="h-empty-state__action">
      <slot name="action" />
    </div>
  </div>
</template>

<style scoped>
.h-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  text-align: center;
}
.h-empty-state__icon {
  font-size: 48px;
  color: var(--h-text-light);
  margin-bottom: 16px;
}
.h-empty-state__title {
  font-size: 16px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 8px;
}
.h-empty-state__description {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: 0;
  max-width: 320px;
}
.h-empty-state__action {
  margin-top: 20px;
}
</style>
```

#### 5.2 Standardisasi Table Action Buttons

```vue
<!-- Tambah di HDataTable sebagai built-in action pattern -->
<template #actions="{ record }">
  <a-space :size="4">
    <a-tooltip title="Lihat Detail">
      <a-button type="text" size="small" @click="$emit('view', record)">
        <EyeOutlined />
      </a-button>
    </a-tooltip>
    <a-tooltip title="Edit">
      <a-button type="text" size="small" @click="$emit('edit', record)">
        <EditOutlined />
      </a-button>
    </a-tooltip>
  </a-space>
</template>
```

#### 5.3 Chart Improvements — Adopsi CRM Patterns

CRM dashboard menggunakan chart yang lebih clean. Rekomendasi untuk ECharts di Gizera:

```js
// Tambahkan di useHorizonChart.js sebagai default theme
const crmInspiredDefaults = {
  // Softer grid lines
  grid: {
    left: '3%',
    right: '4%',
    bottom: '8%',
    top: '12%',
    containLabel: true
  },
  // Cleaner axis
  xAxis: {
    axisLine: { show: false },
    axisTick: { show: false },
    splitLine: { lineStyle: { color: '#F4F7FE', type: 'dashed' } }
  },
  yAxis: {
    axisLine: { show: false },
    axisTick: { show: false },
    splitLine: { lineStyle: { color: '#F4F7FE', type: 'dashed' } }
  },
  // Rounded bar corners
  series: {
    bar: { itemStyle: { borderRadius: [4, 4, 0, 0] } }
  }
}
```

---

## 6. Specific Recommendations per Dashboard Level

### 6.1 Dashboard Admin BGN (Monitoring Nasional)

**Current state**: Flat `a-card` + `a-statistic`, 2 tabel terpisah, peta sebaran.

**Rekomendasi adopsi CRM patterns:**

```
┌─────────────────────────────────────────────────────────┐
│ [Filter Bar: Yayasan | SPPG | Date Range | Refresh]    │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │
│  │ 📊 Total │ │ 🚚 Total │ │ 💰 Total │ │ ⭐ Rata² │  │
│  │  Porsi   │ │Pengiriman│ │Pengeluaran│ │  Review  │  │
│  │  12,450  │ │   342    │ │ Rp 1.2M  │ │  4.2/5   │  │
│  │ ↑ 85%    │ │ ↑ 92%    │ │ 78% serap│ │ 156 rev  │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Peta Sebaran [Yayasan ✓] [SPPG ✓] [Sekolah ✓]  │   │
│  │ ┌─────────────────────────────────────────────┐ │   │
│  │ │              🗺️ Map                         │ │   │
│  │ └─────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Ringkasan Performa                               │   │
│  │ [Tab: Per Yayasan] [Tab: Per SPPG]               │   │
│  │ ┌───────────────────────────────────────────────┐│   │
│  │ │ Table with inline progress bars & badges      ││   │
│  │ └───────────────────────────────────────────────┘│   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

**Perubahan kunci:**
1. Ganti `a-statistic` → `HStatCard` (konsisten dengan dashboard SPPG)
2. Gabungkan 2 tabel jadi tabbed view (progressive disclosure)
3. Tambah inline progress bars di tabel untuk delivery rate & budget absorption
4. Pindahkan peta ke posisi setelah KPI cards (bukan sebelum)

#### Implementasi KPI Cards untuk BGN:

```vue
<!-- Ganti a-statistic cards di DashboardBGNView.vue -->
<div class="stats-row">
  <HStatCard
    :icon="AppstoreOutlined"
    icon-bg="linear-gradient(135deg, #5A4372 0%, #3D2B53 100%)"
    label="Total Porsi Diproduksi"
    :value="String(dashboard?.aggregated_production?.total_portions || 0)"
    :change="`Completion: ${formatPercent(dashboard?.aggregated_production?.completion_rate)}`"
    change-type="increase"
    :loading="loading"
  />
  <!-- ... repeat for other KPIs -->
</div>
```

### 6.2 Dashboard Kepala SPPG (Operasional Dapur)

**Current state**: Sudah menggunakan HStatCard, HChartCard, HDataTable — foundation bagus.

**Rekomendasi adopsi CRM patterns:**

1. **Reorder sections** berdasarkan priority (CRM pattern: most actionable first):
   ```
   1. KPI Performa (sudah ada) ← keep
   2. Critical Stock Alert (pindah ke atas!) ← actionable
   3. Aktivitas Charts (produksi, delivery, cleaning)
   4. Detail Tables (produksi, delivery, cleaning)
   5. Arus Kas (less urgent, bisa di-collapse)
   6. Ulasan & Top Supplier (informational)
   7. Peta Sebaran (pindah ke bawah — less frequent use)
   ```

2. **Critical stock sebagai alert banner**, bukan section di bawah:
   ```vue
   <!-- Pindah ke atas, setelah KPI cards -->
   <a-alert
     v-if="criticalStockItems.length > 0"
     type="warning"
     show-icon
     banner
     style="border-radius: 12px; margin-bottom: 8px;"
   >
     <template #message>
       <strong>{{ criticalStockItems.length }} bahan baku</strong> dalam kondisi kritis
       <a-button type="link" size="small" @click="goToInventory">
         Lihat Detail →
       </a-button>
     </template>
   </a-alert>
   ```

3. **Collapsible sections** untuk data yang jarang dilihat:
   ```vue
   <a-collapse ghost>
     <a-collapse-panel header="Ulasan & Rating" key="reviews">
       <!-- Rating stat cards -->
     </a-collapse-panel>
     <a-collapse-panel header="Top 5 Supplier" key="suppliers">
       <!-- Supplier list -->
     </a-collapse-panel>
   </a-collapse>
   ```

### 6.3 Dashboard Kepala Yayasan (Multi-SPPG Monitoring)

**Rekomendasi layout mengikuti CRM drill-down pattern:**

```
┌─────────────────────────────────────────────────────────┐
│ Dashboard Yayasan: [Nama Yayasan]                       │
├─────────────────────────────────────────────────────────┤
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │
│  │ Total    │ │ Avg      │ │ Total    │ │ Audit    │  │
│  │ SPPG: 5  │ │ Rating   │ │ Porsi    │ │ Score    │  │
│  │          │ │ 4.3/5    │ │ 45,200   │ │ 87%      │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘  │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Performa per SPPG                                │   │
│  │ ┌─────────────────────────────────────────────┐ │   │
│  │ │ SPPG Card 1: [Name] [Rating] [Progress]    │ │   │
│  │ │ SPPG Card 2: [Name] [Rating] [Progress]    │ │   │
│  │ │ SPPG Card 3: [Name] [Rating] [Progress]    │ │   │
│  │ └─────────────────────────────────────────────┘ │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ Riwayat Audit Terbaru                            │   │
│  │ [Table: SPPG | Tanggal | Skor | Status]         │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

---

## 7. Prioritas Implementasi

### Phase 1: Quick Wins (1-2 hari)

| # | Task | File | Impact |
|---|------|------|--------|
| 1 | Buat `HStatusBadge.vue` standalone component | `web/src/components/horizon/` | Konsistensi status di seluruh app |
| 2 | Buat `statusColors.js` utility | `web/src/utils/` | Single source of truth untuk status mapping |
| 3 | Tambah `--h-info` color terpisah dari primary | `web/src/styles/horizon/variables.css` | Disambiguasi info vs brand color |
| 4 | Perbaiki section spacing di dashboards (gap: 32px) | Dashboard views | Breathing room, less cluttered |
| 5 | Tambah color-coded left border utility classes | `web/src/styles/horizon/utilities.css` | Quick visual scanning |

### Phase 2: Component Improvements (3-5 hari)

| # | Task | File | Impact |
|---|------|------|--------|
| 6 | Buat `HSectionHeader.vue` dengan subtitle + badge | `web/src/components/horizon/` | Consistent section headers |
| 7 | Buat `HEmptyState.vue` | `web/src/components/horizon/` | Better empty data UX |
| 8 | Extend `HStatCard` dengan trend indicator | `web/src/components/horizon/HStatCard.vue` | Temporal context di KPI |
| 9 | Tambah sidebar section dividers + badge counter | `web/src/components/layout/HSidebar.vue` | Better navigation grouping |
| 10 | Standardisasi chart defaults (cleaner axes, rounded bars) | `web/src/composables/useHorizonChart.js` | Calmer chart aesthetic |

### Phase 3: Dashboard Restructuring (5-7 hari)

| # | Task | File | Impact |
|---|------|------|--------|
| 11 | Refactor Dashboard BGN: HStatCard + tabbed tables | `web/src/views/DashboardBGNView.vue` | Consistent with SPPG dashboard |
| 12 | Reorder Dashboard SPPG sections (critical stock ke atas) | `web/src/views/DashboardKepalaSSPGView.vue` | Actionable items first |
| 13 | Tambah progressive disclosure (collapsible sections) | Dashboard views | Reduce cognitive load |
| 14 | Buat `HAlertCard.vue` untuk critical items | `web/src/components/horizon/` | Standardized alert pattern |

---

## 8. Design Token Additions

Tambahkan tokens berikut ke `variables.css`:

```css
/* === TAMBAHAN TOKENS === */

/* Info color (terpisah dari primary) */
--h-info: #4A90D9;
--h-info-light: #6BA3E3;
--h-info-dark: #3A7BC8;

/* Processing color (untuk animated states) */
--h-processing: #5A4372;

/* Card variants */
--h-card-border-accent-width: 4px;

/* Section spacing */
--h-section-gap: 32px;
--h-section-gap-mobile: 24px;

/* Badge */
--h-badge-height: 20px;
--h-badge-font-size: 11px;
--h-badge-radius: 10px;

/* Alert card */
--h-alert-indicator-width: 4px;
```

---

## 9. Accessibility Checklist (CRM-Inspired)

CRM dashboard yang baik juga memperhatikan accessibility. Checklist untuk Gizera:

| Check | Status | Action |
|-------|--------|--------|
| Status badge color + text (tidak hanya warna) | ✅ | Sudah ada dot + text |
| Contrast ratio status text on tinted bg | ⚠️ | Verify `#d4940a` on `rgba(255,181,71,0.1)` — mungkin perlu darken |
| Focus indicators di sidebar menu items | ⚠️ | Tambah `outline` on `:focus-visible` |
| ARIA labels di chart cards | ❌ | Tambah `aria-label` dengan summary data |
| Screen reader text untuk trend arrows | ❌ | Tambah `sr-only` text: "naik 12%" |
| Keyboard navigation di tabbed tables | ✅ | Ant Design tabs sudah support |

---

## Ringkasan

Pattern utama dari CRM Dashboard Behance yang paling berdampak untuk Gizera:

1. **Progressive disclosure** — jangan tampilkan semua data sekaligus, gunakan tabs dan collapsible sections
2. **Consistent status badges** — standardisasi `HStatusBadge` + `statusColors.js` di seluruh app
3. **Visual hierarchy melalui spacing** — naikkan section gap ke 32px, tambah section headers dengan subtitle
4. **Actionable items first** — critical stock dan alerts harus di atas, bukan di bawah
5. **Card accent borders** — color-coded left border untuk quick visual scanning
6. **Calm aesthetic** — cleaner chart axes, softer grid lines, generous whitespace

Semua rekomendasi di atas kompatibel dengan Ant Design Vue dan design system Horizon yang sudah ada di Gizera.
