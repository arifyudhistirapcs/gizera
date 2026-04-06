# Dokumen Persyaratan — Redesign UI/UX Datum

## Pendahuluan

Dokumen ini mendefinisikan persyaratan untuk redesign menyeluruh UI/UX aplikasi Gizera ERP SPPG, mencakup web dashboard (Vue 3 + Ant Design Vue) dan PWA mobile (Vue 3 + Vant 4). Redesign mengikuti bahasa desain proyek "Datum" dari Behance yang mengedepankan tampilan bersih, minimal, profesional, dan korporat — menggantikan tampilan saat ini yang dinilai seperti "template AI" dengan warna ungu berlebihan, shadow berat, dan gradien mencolok.

## Glosarium

- **Sistem_Desain**: Kumpulan design token (variabel CSS), komponen, dan aturan visual yang mengatur tampilan seluruh aplikasi
- **Design_Token**: Variabel CSS yang mendefinisikan warna, tipografi, spacing, border-radius, dan shadow
- **Web_Dashboard**: Aplikasi web ERP SPPG berbasis Vue 3 + Ant Design Vue yang diakses melalui browser desktop
- **PWA_Mobile**: Aplikasi Progressive Web App berbasis Vue 3 + Vant 4 yang diakses melalui perangkat mobile
- **Sidebar**: Panel navigasi utama di sisi kiri Web_Dashboard
- **Bottom_Navigation**: Komponen navigasi di bagian bawah layar PWA_Mobile
- **Stat_Card**: Komponen kartu yang menampilkan metrik/statistik ringkas di dashboard
- **Chart_Card**: Komponen kartu yang menampilkan grafik/chart di dashboard
- **Data_Table**: Komponen tabel untuk menampilkan data tabular
- **Palet_Datum**: Palet warna resmi Datum: #303030 (teks utama/gelap), #FFFFFF (putih), #CCE2C8 (hijau muda/aksen), #6B6B6B (abu medium/sekunder), #D8D8DB (abu terang/border)
- **Theme_CSS**: File `web/src/styles/theme.css` yang saat ini memuat tema merah Montserrat yang berkonflik
- **Variables_CSS**: File `web/src/styles/horizon/variables.css` yang memuat semua design token web
- **Horizon_Mobile_CSS**: File `pwa/src/styles/horizon-mobile.css` yang memuat semua design token mobile
- **Ilustrasi_Figure**: Gambar vektor/SVG bergaya flat illustration yang menggambarkan aktivitas operasional SPPG (memasak, pengiriman, pengemasan, dll)
- **Lottie_Animation**: Animasi ringan berbasis JSON dari LottieFiles yang diputar menggunakan library lottie-web atau @lottiefiles/vue-lottie, digunakan untuk empty state, loading, onboarding, dan elemen dekoratif

## Batasan dan Constraint

1. Redesign ini HANYA mencakup perubahan UI/UX (CSS, template Vue, komponen visual). TIDAK ADA perubahan pada:
   - Backend (Go/Gin/GORM) — tidak ada endpoint, service, model, atau migration yang diubah
   - Business logic — alur bisnis, validasi, kalkulasi, dan aturan tetap sama persis
   - API contract — request/response format, endpoint URL, dan parameter tidak berubah
   - Database schema — tidak ada perubahan tabel, kolom, atau relasi
   - Router/routing logic — path dan guard tetap sama, hanya tampilan yang berubah
   - State management (Pinia stores) — logic di store tidak diubah, hanya cara menampilkan data
   - Service layer (API calls) — tidak ada perubahan pada file service/API
2. Semua fungsionalitas yang ada HARUS tetap berfungsi identik setelah redesign
3. Perubahan HANYA pada file: CSS/style, Vue template (`<template>`), Vue `<style>`, dan asset statis (SVG, Lottie JSON)
4. Script logic (`<script>`) pada komponen Vue HANYA boleh diubah jika diperlukan untuk menambahkan Lottie player atau mengubah class binding — TIDAK untuk mengubah business logic

## Persyaratan

### Persyaratan 1: Penghapusan Konflik Tema Ganda

**User Story:** Sebagai developer, saya ingin menghilangkan konflik antara dua sistem tema yang dimuat bersamaan, sehingga seluruh aplikasi menggunakan satu sumber kebenaran desain yang konsisten.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL menggunakan satu file design token tunggal sebagai sumber kebenaran untuk seluruh Web_Dashboard
2. WHEN Web_Dashboard dimuat, THE Sistem_Desain SHALL memastikan tidak ada variabel CSS dari Theme_CSS (tema merah #f82c17 / Montserrat) yang aktif
3. THE Sistem_Desain SHALL menghapus atau menonaktifkan seluruh override Ant Design yang didefinisikan di Theme_CSS
4. IF Theme_CSS masih di-import oleh file lain, THEN THE Sistem_Desain SHALL menghapus import tersebut dan mengarahkan ke Variables_CSS

### Persyaratan 2: Palet Warna Datum

**User Story:** Sebagai pengguna, saya ingin melihat palet warna yang bersih, netral, dan profesional mengikuti gaya Datum, sehingga aplikasi terlihat korporat dan terpercaya.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL menggunakan warna background utama hijau muda sangat terang / sage (#E8EDE5) untuk area konten utama, mengikuti referensi Datum — menggantikan warna krem (#F8FDEA) saat ini. Bukan putih.
2. THE Sistem_Desain SHALL menggunakan warna teks utama gelap (#303030) dengan kontras rasio minimal 4.5:1 terhadap background sesuai WCAG AA
3. THE Sistem_Desain SHALL menggunakan warna teks sekunder abu medium (#6B6B6B) untuk label dan teks pendukung
4. THE Sistem_Desain SHALL menggunakan warna aksen hijau muda (#CCE2C8) untuk indikator positif, keberhasilan, dan elemen aksen
5. THE Sistem_Desain SHALL menggunakan warna abu terang (#D8D8DB) untuk border, divider, dan elemen pembatas
6. THE Sistem_Desain SHALL menghapus seluruh penggunaan warna ungu (#5A4372) sebagai warna primer dan menggantinya dengan warna gelap (#303030) untuk elemen interaktif utama
7. THE Sistem_Desain SHALL menghapus seluruh gradien ungu (linear-gradient #5A4372 ke #3D2B53) dari komponen UI
8. THE Sistem_Desain SHALL mendefinisikan palet semantik: success (#CCE2C8 / hijau muda Datum), warning (amber subtle), error (merah subtle), info (biru subtle) dengan saturasi rendah yang konsisten dengan estetika Datum

### Persyaratan 3: Sistem Tipografi Datum

**User Story:** Sebagai pengguna, saya ingin tipografi yang bersih dan modern, sehingga teks mudah dibaca dan tampilan terasa profesional.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL menggunakan font family Urbanist sebagai font utama dengan fallback ke system font stack (-apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif)
2. THE Sistem_Desain SHALL menghapus penggunaan font Montserrat dari Theme_CSS dan font DM Sans dari Variables_CSS
3. THE Sistem_Desain SHALL mendefinisikan skala ukuran font sesuai Datum: 32px, 24px, 18px, 14px, 12px — dengan mapping: xs (12px), sm (14px), base (14px), lg (18px), xl (24px), 2xl (32px)
4. THE Sistem_Desain SHALL membatasi penggunaan font-weight pada tiga level: regular (400), medium (500), dan semibold (600)
5. THE Sistem_Desain SHALL menggunakan letter-spacing 0% (normal) sesuai spesifikasi Datum
6. THE Sistem_Desain SHALL menggunakan line-height 1.5 sebagai default dan 1.25 untuk heading

### Persyaratan 4: Sistem Shadow dan Elevasi Datum

**User Story:** Sebagai pengguna, saya ingin tampilan kartu dan komponen yang flat dan bersih tanpa shadow berlebihan, sehingga UI terasa modern dan ringan.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL mengganti shadow kartu saat ini (`0px 18px 40px rgba(112, 144, 176, 0.12)`) dengan shadow minimal (`0 1px 3px rgba(0, 0, 0, 0.04)`) atau tanpa shadow
2. THE Sistem_Desain SHALL menggunakan border 1px solid dengan warna abu terang (#D8D8DB) sebagai pembatas utama antar komponen, menggantikan shadow
3. THE Sistem_Desain SHALL mendefinisikan maksimal tiga level shadow: none (default), sm (`0 1px 2px rgba(0, 0, 0, 0.05)`), dan md (`0 2px 8px rgba(0, 0, 0, 0.08)`)
4. THE Sistem_Desain SHALL menghapus shadow dekoratif dari navbar mobile (`0 4px 20px rgba(90, 67, 114, 0.25)`)

### Persyaratan 5: Sistem Border Radius Datum

**User Story:** Sebagai pengguna, saya ingin sudut komponen yang subtle dan konsisten, sehingga UI tidak terlihat "bubbly" atau seperti template.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL mengganti skala border-radius saat ini (8px-20px) dengan skala yang lebih kecil: sm (4px), md (6px), lg (8px), full (9999px)
2. THE Sistem_Desain SHALL menggunakan border-radius 8px sebagai radius maksimal untuk kartu dan kontainer
3. THE Sistem_Desain SHALL menggunakan border-radius 6px untuk input, button, dan elemen interaktif
4. THE Sistem_Desain SHALL menggunakan border-radius 4px untuk tag, badge, dan elemen kecil

### Persyaratan 6: Redesign Sidebar Web

**User Story:** Sebagai pengguna web, saya ingin sidebar navigasi yang bersih dan minimal mengikuti gaya Datum, sehingga navigasi terasa profesional dan tidak mencolok.

#### Kriteria Penerimaan

1. THE Sidebar SHALL menggunakan background putih (#FFFFFF) dengan border kanan 1px solid (#D8D8DB), tanpa shadow
2. THE Sidebar SHALL menampilkan menu item dengan teks warna abu gelap (#303030) dan ikon warna abu medium (#6B6B6B)
3. WHEN menu item aktif, THE Sidebar SHALL menandai item tersebut dengan background abu sangat terang (#F0F0F0) dan teks warna gelap (#303030), tanpa menggunakan warna ungu atau gradien
4. WHEN menu item di-hover, THE Sidebar SHALL menampilkan background abu sangat terang (#F7F8FA)
5. THE Sidebar SHALL menggunakan border-radius 6px untuk menu item, menggantikan 12px saat ini
6. THE Sidebar SHALL menampilkan logo area dengan tinggi proporsional dan tanpa border-bottom yang mencolok

### Persyaratan 7: Redesign Header Web

**User Story:** Sebagai pengguna web, saya ingin header yang bersih dan fungsional, sehingga informasi halaman mudah ditemukan tanpa elemen visual berlebihan.

#### Kriteria Penerimaan

1. THE HHeader SHALL menggunakan background putih (#FFFFFF) dengan border-bottom 1px solid (#D8D8DB), menggantikan box-shadow saat ini
2. THE HHeader SHALL menampilkan judul halaman dengan font-size 20px, font-weight 600, dan warna teks gelap (#303030)
3. THE HHeader SHALL menampilkan breadcrumb dengan font-size 13px dan warna teks abu (#6B6B6B)
4. THE HHeader SHALL memiliki tinggi 64px untuk desktop dan 56px untuk mobile, menggantikan 88px/68px saat ini

### Persyaratan 8: Redesign Stat Card

**User Story:** Sebagai pengguna, saya ingin kartu statistik yang bersih dan informatif mengikuti gaya Datum, sehingga data ringkas mudah dibaca.

#### Kriteria Penerimaan

1. THE Stat_Card SHALL menggunakan background putih dengan border 1px solid (#D8D8DB) dan tanpa shadow, atau shadow minimal
2. THE Stat_Card SHALL menampilkan ikon dalam kontainer dengan background abu terang (#F0F0F0) dan ikon berwarna gelap (#303030), menggantikan gradien ungu saat ini
3. THE Stat_Card SHALL menampilkan label dengan font-size 13px, warna abu (#6B6B6B), dan font-weight 500
4. THE Stat_Card SHALL menampilkan nilai dengan font-size 20px, warna gelap (#303030), dan font-weight 600
5. THE Stat_Card SHALL menggunakan border-radius 8px dan padding 20px
6. THE Stat_Card SHALL menampilkan indikator perubahan (naik/turun) dengan warna hijau (#CCE2C8) untuk positif dan merah (#EF4444) untuk negatif, menggunakan font-size 12px

### Persyaratan 9: Redesign Chart Card

**User Story:** Sebagai pengguna, saya ingin kartu grafik yang bersih dengan fokus pada data, sehingga visualisasi mudah dipahami.

#### Kriteria Penerimaan

1. THE Chart_Card SHALL menggunakan background putih dengan border 1px solid (#D8D8DB) dan tanpa shadow
2. THE Chart_Card SHALL menampilkan judul dengan font-size 16px, font-weight 600, dan warna gelap (#303030)
3. THE Chart_Card SHALL menggunakan border-radius 8px dan padding 20px
4. THE Chart_Card SHALL menyediakan area header-right untuk filter atau aksi tanpa styling berlebihan

### Persyaratan 10: Redesign Data Table

**User Story:** Sebagai pengguna, saya ingin tabel data yang bersih dan mudah di-scan, sehingga informasi tabular cepat ditemukan.

#### Kriteria Penerimaan

1. THE Data_Table SHALL menggunakan header tabel dengan background abu sangat terang (#F7F8FA), font-size 12px uppercase, font-weight 500, dan warna teks abu (#6B6B6B)
2. THE Data_Table SHALL menggunakan baris tabel dengan border-bottom 1px solid (#F0F0F0) dan tanpa background bergantian (striped)
3. THE Data_Table SHALL menampilkan teks sel dengan font-size 14px dan warna gelap (#303030)
4. WHEN baris tabel di-hover, THE Data_Table SHALL menampilkan background abu sangat terang (#F7F8FA)
5. THE Data_Table SHALL menghapus background merah (#ffeae8) dari header tabel yang berasal dari Theme_CSS

### Persyaratan 11: Redesign Halaman Login Web

**User Story:** Sebagai pengguna web, saya ingin halaman login yang bersih dan profesional, sehingga kesan pertama terhadap aplikasi positif.

#### Kriteria Penerimaan

1. THE LoginView SHALL menggunakan layout split-screen dengan sisi kiri form login berlatar putih dan sisi kanan branding berlatar gelap netral (#303030), menggantikan gradien ungu
2. THE LoginView SHALL menampilkan input field dengan border 1px solid (#D8D8DB), border-radius 6px, dan tinggi 44px
3. THE LoginView SHALL menampilkan tombol submit dengan background gelap (#303030), border-radius 6px, dan tanpa gradien
4. THE LoginView SHALL menghapus elemen dekoratif berlebihan (lingkaran transparan, efek hover translateY)
5. THE LoginView SHALL menampilkan branding sisi kanan dengan tipografi bersih tanpa efek glassmorphism atau backdrop-filter

### Persyaratan 12: Redesign Dashboard SPPG

**User Story:** Sebagai Kepala SPPG, saya ingin dashboard yang bersih dan informatif mengikuti gaya Datum, sehingga monitoring operasional efektif.

#### Kriteria Penerimaan

1. THE DashboardView SHALL menampilkan grid stat card dengan gap 16px dan menggunakan komponen Stat_Card yang sudah di-redesign
2. THE DashboardView SHALL menampilkan kartu konten (Status Produksi, Aktivitas Terbaru) menggunakan border dan tanpa shadow
3. THE DashboardView SHALL menggunakan warna netral untuk ikon stat card, menggantikan gradien warna-warni saat ini
4. THE DashboardView SHALL menampilkan welcome card dengan styling minimal: background putih, border, tanpa padding berlebihan
5. THE DashboardView SHALL menampilkan overview/highlight widget utama dengan background hijau muda (#CCE2C8) mengikuti style Overview card di Datum, berisi stat cards kecil dengan background putih semi-transparan di dalamnya

### Persyaratan 13: Redesign Dashboard BGN

**User Story:** Sebagai Admin BGN, saya ingin dashboard agregat yang konsisten dengan desain Datum, sehingga monitoring nasional terasa profesional.

#### Kriteria Penerimaan

1. THE DashboardBGNView SHALL menggunakan komponen Ant Design (a-card, a-table, a-statistic) yang di-override dengan styling Datum melalui design token
2. THE DashboardBGNView SHALL menampilkan filter card dengan background putih, border, dan tanpa shadow
3. THE DashboardBGNView SHALL menampilkan tabel performa dengan header abu terang dan baris bersih sesuai Persyaratan 10
4. THE DashboardBGNView SHALL menampilkan metric card dengan styling konsisten tanpa warna latar yang mencolok

### Persyaratan 14: Redesign PWA Mobile — Design Token

**User Story:** Sebagai pengguna mobile, saya ingin tampilan PWA yang bersih dan modern mengikuti gaya Datum, sehingga pengalaman mobile terasa profesional.

#### Kriteria Penerimaan

1. THE Horizon_Mobile_CSS SHALL menggunakan palet warna yang identik dengan Variables_CSS web (Palet_Datum)
2. THE Horizon_Mobile_CSS SHALL mengganti background utama dari krem (#F8FDEA) ke sage/hijau muda sangat terang (#E8EDE5) mengikuti Datum
3. THE Horizon_Mobile_CSS SHALL mengganti semua referensi warna ungu (#5A4372) dengan warna netral gelap (#303030)
4. THE Horizon_Mobile_CSS SHALL menggunakan border-radius yang lebih kecil: sm (4px), md (6px), lg (8px)
5. THE Horizon_Mobile_CSS SHALL mengganti shadow kartu dari `0px 18px 40px rgba(112, 144, 176, 0.12)` ke border 1px solid (#D8D8DB)
6. THE Horizon_Mobile_CSS SHALL menggunakan font Urbanist menggantikan DM Sans

### Persyaratan 15: Redesign PWA Mobile — Navbar

**User Story:** Sebagai pengguna mobile, saya ingin navbar atas yang bersih tanpa gradien, sehingga navigasi terasa ringan dan modern.

#### Kriteria Penerimaan

1. THE Horizon_Mobile_CSS SHALL mengganti navbar background dari gradien ungu (`linear-gradient(135deg, #5A4372, #7B5E99)`) ke background putih (#FFFFFF) dengan border-bottom 1px solid (#D8D8DB)
2. THE Horizon_Mobile_CSS SHALL menampilkan judul navbar dengan warna teks gelap (#303030), menggantikan warna putih
3. THE Horizon_Mobile_CSS SHALL menghapus border-radius bawah navbar (20px) sehingga navbar berbentuk persegi
4. THE Horizon_Mobile_CSS SHALL menampilkan ikon navigasi navbar dengan warna gelap (#303030), menggantikan warna putih

### Persyaratan 16: Redesign Bottom Navigation Mobile

**User Story:** Sebagai pengguna mobile, saya ingin navigasi bawah yang bersih dan fungsional tanpa elemen flashy, sehingga navigasi mudah dan profesional.

#### Kriteria Penerimaan

1. THE Bottom_Navigation SHALL menggunakan background putih (#FFFFFF) dengan border-top 1px solid (#D8D8DB), menggantikan shadow dan border-radius atas 24px
2. THE Bottom_Navigation SHALL menampilkan item navigasi dengan ikon dan label berwarna abu (#6B6B6B) untuk state default
3. WHEN item navigasi aktif, THE Bottom_Navigation SHALL menandai item dengan warna gelap (#303030) dan font-weight 600, tanpa warna ungu
4. THE Bottom_Navigation SHALL menghapus floating center button dengan gradien ungu dan animasi pulse, menggantinya dengan item navigasi inline yang sejajar dengan item lainnya
5. THE Bottom_Navigation SHALL menggunakan tinggi 56px dengan border-radius atas 0px (persegi)

### Persyaratan 17: Redesign Metric Card Mobile

**User Story:** Sebagai pengguna mobile, saya ingin kartu metrik yang bersih dan informatif, sehingga data ringkas mudah dibaca di layar kecil.

#### Kriteria Penerimaan

1. THE MetricCard SHALL menggunakan background putih dengan border 1px solid (#D8D8DB), menggantikan shadow
2. THE MetricCard SHALL menampilkan ikon dalam kontainer dengan background abu terang (#F0F0F0) dan border-radius 6px, menggantikan warna ungu
3. THE MetricCard SHALL menggunakan border-radius 8px dan padding 16px
4. THE MetricCard SHALL menampilkan label dengan font-size 12px warna abu dan nilai dengan font-size 18px warna gelap

### Persyaratan 18: Redesign Task Card Mobile

**User Story:** Sebagai pengguna mobile, saya ingin kartu tugas yang bersih dan mudah di-scan, sehingga informasi tugas cepat dipahami.

#### Kriteria Penerimaan

1. THE TaskCard SHALL menggunakan background putih dengan border 1px solid (#D8D8DB), menggantikan shadow
2. THE TaskCard SHALL menampilkan nomor urut route dalam lingkaran dengan background abu terang (#F0F0F0) dan teks gelap, menggantikan warna ungu
3. THE TaskCard SHALL menggunakan border-radius 8px
4. THE TaskCard SHALL menampilkan tag tipe tugas dengan background subtle (hijau terang untuk delivery, amber terang untuk pickup) dan border-radius 4px

### Persyaratan 19: Redesign Halaman Login Mobile

**User Story:** Sebagai pengguna mobile, saya ingin halaman login yang bersih dan profesional, sehingga pengalaman pertama di mobile positif.

#### Kriteria Penerimaan

1. THE PWA LoginView SHALL menggunakan background putih atau abu sangat terang untuk seluruh halaman, menggantikan gradien ungu
2. THE PWA LoginView SHALL menampilkan logo dengan ukuran proporsional tanpa efek drop-shadow berlebihan
3. THE PWA LoginView SHALL menampilkan form card dengan border 1px solid (#D8D8DB) dan border-radius 8px, tanpa shadow berat
4. THE PWA LoginView SHALL menampilkan tombol login dengan background gelap (#303030), border-radius 6px, dan tanpa shadow dekoratif
5. THE PWA LoginView SHALL menghapus elemen dekoratif (lingkaran transparan, animasi float)

### Persyaratan 20: Redesign Dashboard Mobile

**User Story:** Sebagai pengguna mobile, saya ingin dashboard yang bersih dan terorganisir, sehingga monitoring di perangkat mobile efektif.

#### Kriteria Penerimaan

1. THE PWA DashboardView SHALL menampilkan grid metrik menggunakan komponen MetricCard yang sudah di-redesign
2. THE PWA DashboardView SHALL menampilkan section card (Arus Kas, Top Supplier, Detail Produksi) dengan border dan tanpa shadow
3. THE PWA DashboardView SHALL menampilkan detail table dengan header abu terang dan baris bersih
4. THE PWA DashboardView SHALL menampilkan item stok kritis dengan border-left warna merah subtle dan background putih, tanpa background abu
5. THE PWA DashboardView SHALL menampilkan ranking supplier dengan nomor urut berlatar abu terang dan teks gelap, menggantikan gradien emas/perak/perunggu

### Persyaratan 21: Konsistensi Override Komponen Ant Design (Web)

**User Story:** Sebagai developer, saya ingin override komponen Ant Design yang konsisten dengan Palet_Datum, sehingga seluruh komponen bawaan terlihat seragam.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL meng-override warna primer Ant Design dari ungu (#5A4372) ke gelap netral (#303030)
2. THE Sistem_Desain SHALL meng-override background header tabel Ant Design dari merah (#ffeae8) ke abu terang (#F7F8FA)
3. THE Sistem_Desain SHALL meng-override border-radius komponen Ant Design (button, input, card, modal) ke skala Datum (4px-8px)
4. THE Sistem_Desain SHALL meng-override shadow komponen Ant Design (card, dropdown, modal) ke shadow minimal atau border
5. IF komponen Ant Design menggunakan warna biru default (#1890ff), THEN THE Sistem_Desain SHALL menggantinya dengan warna netral gelap (#303030) untuk elemen interaktif

### Persyaratan 22: Konsistensi Override Komponen Vant (Mobile)

**User Story:** Sebagai developer, saya ingin override komponen Vant yang konsisten dengan Palet_Datum, sehingga seluruh komponen mobile terlihat seragam.

#### Kriteria Penerimaan

1. THE Horizon_Mobile_CSS SHALL meng-override warna primer Vant dari ungu (#5A4372) ke gelap netral (#303030)
2. THE Horizon_Mobile_CSS SHALL meng-override border-radius komponen Vant (button, cell-group, dialog, popup) ke skala Datum (4px-8px)
3. THE Horizon_Mobile_CSS SHALL meng-override shadow komponen Vant (card, cell-group, tabbar) ke border atau shadow minimal
4. THE Horizon_Mobile_CSS SHALL meng-override background switch, checkbox, dan radio dari ungu ke gelap netral (#303030)

### Persyaratan 23: Konsistensi Spacing

**User Story:** Sebagai developer, saya ingin sistem spacing yang konsisten di seluruh aplikasi, sehingga layout terasa teratur dan rapi.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL menggunakan skala spacing berbasis 4px: 4, 8, 12, 16, 20, 24, 32, 40, 48
2. THE Sistem_Desain SHALL menggunakan gap 16px sebagai default antar kartu dan section
3. THE Sistem_Desain SHALL menggunakan padding 20px untuk konten kartu di desktop dan 16px di mobile
4. THE Sistem_Desain SHALL menggunakan padding area konten utama 24px di desktop dan 16px di mobile

### Persyaratan 24: Dukungan Dark Mode

**User Story:** Sebagai pengguna, saya ingin dark mode yang konsisten dengan estetika Datum, sehingga pengalaman di kondisi cahaya rendah tetap nyaman.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL mendefinisikan variabel dark mode dengan background gelap (#1A1A1A), kartu gelap (#252525), dan teks terang (#F7F8FA)
2. THE Sistem_Desain SHALL memastikan dark mode tidak menggunakan warna ungu atau gradien ungu
3. THE Sistem_Desain SHALL memastikan kontras teks di dark mode memenuhi WCAG AA (rasio minimal 4.5:1)
4. WHILE dark mode aktif, THE Sistem_Desain SHALL menggunakan border warna (#303030) menggantikan border terang (#D8D8DB)

### Persyaratan 25: Migrasi View yang Ada

**User Story:** Sebagai developer, saya ingin seluruh view yang ada (40+ web, 15+ mobile) menggunakan design token baru, sehingga tidak ada inkonsistensi visual.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL memastikan seluruh view web menggunakan variabel CSS dari Variables_CSS yang sudah di-update, tanpa hardcoded warna ungu atau merah
2. THE Sistem_Desain SHALL memastikan seluruh view mobile menggunakan variabel CSS dari Horizon_Mobile_CSS yang sudah di-update
3. IF sebuah view menggunakan class utility `.h-card`, THEN THE Sistem_Desain SHALL memastikan class tersebut menggunakan styling Datum (border, tanpa shadow berat)
4. IF sebuah view menggunakan class `.h-gradient-primary`, THEN THE Sistem_Desain SHALL menghapus atau mengganti class tersebut dengan background solid gelap

### Persyaratan 26: Ilustrasi Figure untuk Halaman Kunci

**User Story:** Sebagai pengguna, saya ingin melihat ilustrasi/figure yang relevan dengan konteks SPPG di halaman-halaman kunci, sehingga aplikasi terasa lebih hidup, profesional, dan memiliki identitas visual yang kuat — bukan sekadar template kosong.

#### Kriteria Penerimaan

1. THE Web_Dashboard LoginView SHALL menampilkan ilustrasi figure bergaya flat/modern di sisi branding (kanan) yang menggambarkan aktivitas operasional SPPG (contoh: petugas dapur, pengiriman makanan, atau manajemen gizi), menggantikan logo circle "ERP" saat ini
2. THE Web_Dashboard DashboardView SHALL menampilkan ilustrasi figure kecil di welcome card yang relevan dengan role pengguna (contoh: chef untuk ahli_gizi, truk untuk driver, grafik untuk akuntan)
3. THE PWA_Mobile LoginView SHALL menampilkan ilustrasi figure di area header yang menggambarkan aktivitas SPPG, menggantikan logo saja
4. THE Sistem_Desain SHALL menyimpan semua file ilustrasi dalam format SVG di folder `web/src/assets/illustrations/` dan `pwa/src/assets/illustrations/` untuk optimasi performa
5. THE Sistem_Desain SHALL memastikan ilustrasi menggunakan palet warna yang konsisten dengan Palet_Datum (netral, dengan aksen warna subtle)
6. THE Sistem_Desain SHALL menampilkan ilustrasi figure pada empty state di seluruh halaman yang memiliki kondisi "tidak ada data" (contoh: tabel kosong, daftar kosong), menggantikan komponen `<a-empty>` atau `<van-empty>` default
7. THE Sistem_Desain SHALL memastikan ilustrasi memiliki ukuran responsif yang proporsional di berbagai ukuran layar

### Persyaratan 27: Lottie Animations untuk Interaksi dan Feedback

**User Story:** Sebagai pengguna, saya ingin animasi ringan yang memberikan feedback visual pada momen-momen penting, sehingga pengalaman menggunakan aplikasi terasa lebih smooth dan engaging tanpa berlebihan.

#### Kriteria Penerimaan

1. THE Sistem_Desain SHALL mengintegrasikan library Lottie (lottie-web atau @lottiefiles/dotlottie-vue) di kedua project web dan PWA untuk memutar animasi JSON
2. THE Web_Dashboard SHALL menampilkan Lottie animation pada loading state halaman (menggantikan spinner default Ant Design) dengan animasi yang relevan dengan konteks SPPG (contoh: animasi memasak, animasi pengiriman)
3. THE PWA_Mobile SHALL menampilkan Lottie animation pada pull-to-refresh dan loading state, menggantikan spinner default Vant
4. THE Sistem_Desain SHALL menampilkan Lottie animation pada empty state halaman sebagai alternatif atau pelengkap ilustrasi statis (contoh: animasi kotak kosong, animasi pencarian)
5. THE Web_Dashboard LoginView SHALL menampilkan Lottie animation subtle di sisi branding sebagai elemen dekoratif yang menambah kesan profesional (contoh: animasi data flow, animasi dashboard)
6. THE PWA_Mobile SHALL menampilkan Lottie animation pada success state setelah aksi berhasil (contoh: animasi centang setelah check-in absensi, animasi berhasil setelah submit form)
7. THE Sistem_Desain SHALL menyimpan semua file Lottie JSON di folder `web/src/assets/lottie/` dan `pwa/src/assets/lottie/`
8. THE Sistem_Desain SHALL menggunakan IconScout (https://iconscout.com/lottie-animations) sebagai sumber utama untuk mencari dan mengunduh animasi Lottie yang relevan dengan konteks SPPG (food, delivery, cooking, dashboard, analytics, success, empty state, loading)
9. THE Sistem_Desain SHALL memastikan file Lottie JSON berukuran maksimal 100KB per file untuk menjaga performa
10. THE Sistem_Desain SHALL memastikan animasi Lottie menggunakan palet warna yang konsisten dengan Palet_Datum
10. THE Sistem_Desain SHALL menyediakan komponen wrapper Vue reusable (`LottiePlayer.vue`) untuk web dan mobile yang mendukung props: src, autoplay, loop, width, height
11. THE Sistem_Desain SHALL memastikan animasi Lottie menghormati preferensi `prefers-reduced-motion` — WHEN pengguna mengaktifkan reduced motion, THEN animasi SHALL ditampilkan sebagai frame statis pertama
