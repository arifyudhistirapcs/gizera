# Proposal Redesign V2 — Perubahan Struktural & Visual

## Masalah dengan V1
Redesign V1 hanya mengganti design token (warna, font, shadow, border-radius) tanpa mengubah layout dan struktur. Hasilnya terasa "sama saja" — hanya beda warna.

## Apa yang Akan Berubah di V2

### 1. LAYOUT DASHBOARD — Dari "List Panjang" ke "Grid Modular"

**Sebelum (V1):** Dashboard adalah scroll panjang vertikal — peta di atas, lalu stat cards, lalu tabel, lalu stok kritis. Semua ditumpuk ke bawah.

**Sesudah (V2):**
- Dashboard menggunakan **grid 2-3 kolom** yang memanfaatkan lebar layar
- **Bento-box layout**: kartu-kartu dengan ukuran berbeda (1x1, 2x1, 1x2) membentuk grid yang menarik
- Peta diperkecil jadi widget 2x1, bukan full-width
- Stat cards dikelompokkan dalam **overview card besar** dengan background hijau muda
- Chart dan tabel berdampingan, bukan ditumpuk

### 2. SIDEBAR — Dari "Daftar Menu Biasa" ke "Compact Icon Sidebar"

**Sebelum (V1):** Sidebar 280px dengan teks menu panjang, terlihat seperti template admin generik.

**Sesudah (V2):**
- **Default collapsed** (64px) — hanya ikon yang terlihat
- Hover/click untuk expand ke 240px dengan animasi smooth
- **Grouped sections** dengan divider dan label kategori kecil
- Active indicator: garis vertikal 3px di kiri (bukan background penuh)
- Logo area lebih compact (48px height)
- User avatar di bawah sidebar dengan tooltip nama

### 3. HEADER — Dari "Bar Biasa" ke "Contextual Header"

**Sebelum (V1):** Header 64px dengan judul halaman dan breadcrumb.

**Sesudah (V2):**
- Header lebih tipis (52px)
- **Search bar prominent** di tengah header (seperti Notion/Linear)
- Breadcrumb dihapus — diganti dengan page title yang lebih besar
- Quick actions di kanan: notifikasi bell, user avatar dropdown
- **Greeting banner** di dashboard: "Selamat Pagi, Arif" dengan tanggal dan cuaca

### 4. STAT CARDS — Dari "Kotak Biasa" ke "Metric Tiles Modern"

**Sebelum (V1):** Stat card dengan ikon kotak abu + angka. Semua ukuran sama.

**Sesudah (V2):**
- **Hero metric** (kartu besar 2x1) untuk KPI utama — dengan mini sparkline chart di dalamnya
- **Secondary metrics** (kartu 1x1) lebih compact
- Setiap metric punya **trend indicator** dengan mini chart (bukan cuma panah atas/bawah)
- Progress ring/donut kecil di dalam kartu untuk persentase
- **Color-coded left border** (4px) untuk kategori: hijau=produksi, biru=delivery, amber=stok

### 5. TABEL — Dari "Ant Design Default" ke "Clean Data Grid"

**Sebelum (V1):** Tabel Ant Design standar dengan header abu.

**Sesudah (V2):**
- **Borderless design** — hanya garis horizontal tipis antar baris
- Header tanpa background — hanya teks uppercase kecil dengan spacing besar
- **Inline status pills** yang lebih modern (rounded, dengan dot indicator)
- Row hover: subtle left-border highlight (bukan background penuh)
- **Sticky header** saat scroll
- Pagination di bawah dengan style minimal (1 2 3 ... 10)

### 6. CHART CARDS — Dari "Kotak dengan Chart" ke "Integrated Visualization"

**Sebelum (V1):** Chart card dengan judul di atas dan chart di bawah.

**Sesudah (V2):**
- Chart **tanpa border/card** — langsung embedded di grid
- **Gradient area fill** di bawah line chart (bukan garis saja)
- Donut chart dengan **center label** (angka besar di tengah)
- Bar chart dengan **rounded corners** dan spacing lebih lebar
- **Interactive tooltip** yang lebih informatif

### 7. LOGIN PAGE — Dari "Split Screen Biasa" ke "Immersive Login"

**Sebelum (V1):** Split screen — form kiri, branding gelap kanan.

**Sesudah (V2):**
- **Full-screen background** dengan pattern/texture sage subtle
- Form card **floating di tengah** dengan glassmorphism ringan (blur 8px, bukan berat)
- Logo besar di atas form
- **Animated illustration** (Lottie) di bawah form — bukan di samping
- Social proof: "Digunakan oleh 50+ SPPG di Indonesia" di bawah

### 8. OVERVIEW WIDGET — Dari "Kotak Hijau Biasa" ke "Dashboard Hero Section"

**Sebelum (V1):** Kotak hijau muda dengan 4 angka di dalamnya.

**Sesudah (V2):**
- **Full-width hero section** di atas dashboard
- Background gradient subtle (sage ke putih)
- **4 metric cards** dengan ikon besar dan angka prominent
- **Date range selector** terintegrasi di hero section
- **Quick action buttons**: "Lihat Produksi", "Cek Pengiriman", "Stok Kritis"

### 9. EMPTY STATES — Dari "Ikon + Teks" ke "Illustrated Empty States"

**Sebelum (V1):** Ikon sederhana dengan teks "Belum ada data".

**Sesudah (V2):**
- **Ilustrasi SVG besar** yang kontekstual (chef untuk produksi, truk untuk delivery)
- Teks deskriptif yang helpful: "Belum ada data produksi hari ini. Mulai dengan membuat menu planning."
- **CTA button** yang mengarahkan ke aksi: "Buat Menu Planning →"

### 10. MICRO-INTERACTIONS & POLISH

- **Page transitions**: fade + slide subtle saat navigasi antar halaman
- **Skeleton loading**: shimmer effect yang lebih smooth (bukan kotak abu)
- **Number animations**: angka di stat cards count-up saat pertama kali muncul
- **Hover effects**: kartu sedikit terangkat (2px) dengan shadow subtle
- **Toast notifications**: muncul dari atas dengan slide-down animation
- **Scroll animations**: elemen muncul dengan fade-in saat di-scroll ke view

---

## Prioritas Implementasi

| Fase | Perubahan | Impact |
|------|-----------|--------|
| 1 | Layout grid dashboard + Hero section | Tinggi — langsung terasa berbeda |
| 2 | Sidebar compact + Header baru | Tinggi — navigasi terasa modern |
| 3 | Stat cards modern + Chart improvements | Sedang — data visualization lebih baik |
| 4 | Login page immersive | Sedang — first impression |
| 5 | Tabel clean + Empty states | Sedang — detail polish |
| 6 | Micro-interactions | Rendah — finishing touch |

---

## Catatan Penting
- Semua perubahan tetap PURE UI/UX — tidak ada perubahan backend, API, atau business logic
- Menggunakan palet Datum yang sudah ada (#303030, #E8EDE5, #CCE2C8, #D8D8DB)
- Font Urbanist tetap digunakan
- Kompatibel dengan dark mode
