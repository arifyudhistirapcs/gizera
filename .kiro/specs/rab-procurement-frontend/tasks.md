# Implementation Plan: Frontend RAB, Procurement & Supplier Portal

## Overview

Implementasi frontend untuk fitur RAB, procurement yayasan, dan supplier portal di dua platform: Web Dashboard (Vue 3 + Ant Design Vue) dan PWA Mobile (Vue 3 + Vant UI). Backend API sudah tersedia — fokus pada service layer, views, routes, permissions, dan navigasi. Urutan implementasi berdasarkan dependency: foundation → web dashboard → PWA mobile → integrasi.

## Tasks

- [x] 1. Foundation: Service Layer, Permissions, dan Auth Store
  - [x] 1.1 Buat `web/src/services/rabService.js` — service layer RAB
    - Implementasi semua method: `getRABList`, `getRABDetail`, `updateRAB`, `approveSPPG`, `approveYayasan`, `rejectRAB`, `resubmitRAB`, `getRABComparison`, `getRABPOTracking`
    - Gunakan axios instance dari `api.js` yang sudah ada
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9, 1.10_

  - [x] 1.2 Buat `web/src/services/supplierProductService.js` — service layer supplier product catalog
    - Implementasi method: `getProducts`, `createProduct`, `updateProduct`, `deleteProduct`
    - Gunakan axios instance dari `api.js`
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [x] 1.3 Buat `web/src/services/invoiceService.js` — service layer invoice & pembayaran
    - Implementasi method: `getInvoices`, `createInvoice`, `getInvoiceDetail`, `payInvoice`, `uploadPaymentProof`
    - `uploadPaymentProof` harus menggunakan header `Content-Type: multipart/form-data`
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

  - [x] 1.4 Buat `pwa/src/services/supplierPortalService.js` — service layer supplier portal PWA
    - Implementasi method: `getDashboard`, `getPayments`, `getProducts`, `createProduct`, `updateProduct`, `deleteProduct`, `getInvoices`, `createInvoice`, `getPurchaseOrders`
    - Gunakan axios instance dari `pwa/src/services/api.js`
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8, 4.9, 4.10_

  - [x] 1.5 Update `web/src/utils/permissions.js` — tambah permission baru dan modifikasi existing
    - Tambah permission: RAB_VIEW, RAB_APPROVE_SPPG, RAB_APPROVE_YAYASAN, RAB_EDIT, INVOICE_VIEW, INVOICE_CREATE, INVOICE_PAY, SUPPLIER_PRODUCT_VIEW, SUPPLIER_PRODUCT_MANAGE, SUPPLIER_PORTAL_VIEW, SUPPLIER_DASHBOARD_VIEW
    - Modifikasi: SUPPLIER_VIEW, SUPPLIER_MANAGE, PO_VIEW, PO_CREATE — tambah `kepala_yayasan`
    - Tambah `supplier` ke `NON_OPERATIONAL_ROLES`
    - Tambah `supplier: 'Supplier'` ke `getRoleLabel()`
    - _Requirements: 16.1, 16.2, 16.3, 16.4, 16.5, 16.6, 16.7, 25.5_

  - [x] 1.6 Update `web/src/stores/auth.js` — tambah support role supplier
    - Tambah `supplierId` ref yang di-sync dari `userData.supplier_id` di `_syncTenantFields()`
    - Tambah `isSupplier` computed property (`role === 'supplier'`)
    - Export kedua property baru
    - _Requirements: 25.1, 25.2_

  - [x] 1.7 Update `pwa/src/stores/auth.js` — tambah support role supplier
    - Tambah `supplierId` computed property (`user.value?.supplier_id ?? null`)
    - Tambah `isSupplier` computed property (`userRole.value === 'supplier'`)
    - Export kedua property baru
    - _Requirements: 25.3, 25.4_

  - [ ]* 1.8 Tulis property tests untuk utility functions (fast-check)
    - **Property 1: Format mata uang Rupiah konsisten** — generate bilangan non-negatif, verifikasi output `formatRupiah()` dimulai dengan "Rp" dan merepresentasikan nilai yang sama
    - **Property 4: Permission map mengembalikan hasil benar** — generate pasangan (role, permission) acak, verifikasi konsistensi dengan PERMISSIONS map
    - **Property 5: Default route mapping untuk setiap role** — generate role acak, verifikasi `getDefaultRouteForRole(role)` mengembalikan path valid
    - **Validates: Requirements 5.8, 7.7, 9.4, 16.1-16.6, 17.5, 24.3**

- [x] 2. Checkpoint — Pastikan foundation layer berfungsi
  - Pastikan semua service files, permissions, dan auth store updates tidak ada error. Tanyakan user jika ada pertanyaan.

- [x] 3. Web Dashboard: Halaman RAB
  - [x] 3.1 Buat `web/src/views/RABListView.vue` — halaman daftar RAB
    - Tabel `a-table` dengan kolom: rab_number, menu_plan, SPPG, status, total_amount, tanggal, aksi
    - Filter status menggunakan `a-select` (draft, approved_sppg, approved_yayasan, revision_requested, completed)
    - Status badge menggunakan `a-tag` dengan warna per status (draft=default, approved_sppg=blue, approved_yayasan=green, revision_requested=orange, completed=purple)
    - Format total_amount dalam Rupiah (Rp) dengan pemisah ribuan
    - Pagination server-side (page, per_page)
    - Klik baris navigasi ke `/rab/:id`
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 5.8_

  - [x] 3.2 Buat `web/src/views/RABDetailView.vue` — halaman detail RAB dengan approval flow
    - Header: rab_number, menu_plan, SPPG, status, total_amount, created_by, tanggal
    - Tabel RAB items: ingredient, quantity, unit, unit_price, subtotal, recommended_supplier, status item, po_id, grn_id
    - Tombol aksi berdasarkan status+role: Approve SPPG (draft+kepala_sppg), Approve Yayasan/Tolak (approved_sppg+kepala_yayasan), Edit (draft/revision_requested+kepala_sppg), Kirim Ulang (revision_requested+kepala_sppg)
    - Modal rejection notes (`a-modal` + textarea) untuk tombol Tolak
    - `a-alert` warning untuk revision_notes saat status revision_requested
    - Edit inline/modal untuk quantity dan unit_price
    - Tab "PO Tracking" — panggil `getRABPOTracking(id)`, tampilkan tabel PO + status GRN
    - Tab "Perbandingan" — panggil `getRABComparison(id)`, tampilkan tabel RAB vs Aktual
    - Error handling: `message.error()` saat gagal, `message.success()` + refresh saat berhasil
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 6.9, 6.10, 6.11, 6.12_

  - [ ]* 3.3 Tulis property test untuk visibilitas tombol aksi RAB (fast-check)
    - **Property 2: Visibilitas tombol aksi RAB berdasarkan status dan role** — generate kombinasi (status, role) acak, verifikasi `getVisibleActions(status, role)` sesuai spesifikasi
    - **Property 3: Mapping status ke warna tag konsisten** — generate status acak, verifikasi `getStatusColor(status)` bukan undefined/null
    - **Validates: Requirements 5.5, 6.3, 6.4, 6.6, 6.8**

- [x] 4. Web Dashboard: Halaman Invoice & Katalog Supplier
  - [x] 4.1 Buat `web/src/views/InvoiceListView.vue` — halaman daftar invoice + pembayaran
    - Tabel `a-table`: invoice_number, supplier, PO number, amount, status, due_date, aksi
    - Filter status (pending, paid) menggunakan `a-select`
    - Modal pembayaran: payment_date (date picker), payment_method (default: bank_transfer), upload bukti transfer (file upload)
    - Upload bukti → `uploadPaymentProof()`, lalu `payInvoice()` → `message.success()` + refresh
    - Format amount dalam Rupiah, error handling untuk invoice sudah dibayar
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8_

  - [x] 4.2 Buat `web/src/views/SupplierProductListView.vue` — katalog produk supplier (kepala_yayasan)
    - Tabel `a-table`: supplier_name, ingredient_name, unit_price, min_order_qty, stock_quantity, is_available
    - Filter berdasarkan supplier dan ingredient menggunakan `a-select`
    - Format unit_price dalam Rupiah, status ketersediaan menggunakan `a-tag` (hijau/merah)
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5_

- [x] 5. Web Dashboard: Supplier Portal Views
  - [x] 5.1 Buat `web/src/views/SupplierDashboardView.vue` — dashboard supplier
    - Ringkasan menggunakan `a-statistic`/`a-card`: total PO aktif, PO selesai, invoice pending, pembayaran diterima
    - Panggil `GET /api/v1/supplier/dashboard` via `rabService` atau service baru
    - Tabel ringkas 5 PO terbaru dan 5 invoice terbaru
    - Hanya role "supplier"
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5_

  - [x] 5.2 Buat `web/src/views/SupplierProductManageView.vue` — CRUD produk supplier
    - Tabel `a-table`: ingredient_name, unit_price, min_order_qty, stock_quantity, is_available, aksi (edit, hapus)
    - Tombol "Tambah Produk" → modal `a-modal` dengan form: ingredient_id (select), unit_price, min_order_qty, stock_quantity, is_available (switch)
    - Edit via modal, hapus via `a-popconfirm`
    - Handle error `DUPLICATE_SUPPLIER_PRODUCT`
    - Hanya role "supplier"
    - _Requirements: 13.1, 13.2, 13.3, 13.4, 13.5, 13.6_

  - [x] 5.3 Buat `web/src/views/SupplierPOListView.vue` — daftar PO untuk supplier (read-only)
    - Tabel `a-table`: po_number, yayasan, target_sppg, tanggal, total_amount, status, aksi
    - Klik baris → modal/halaman detail PO dengan daftar items
    - Read-only, hanya role "supplier"
    - _Requirements: 14.1, 14.2, 14.3, 14.4_

  - [x] 5.4 Buat `web/src/views/SupplierInvoiceView.vue` — invoice supplier (buat + lihat)
    - Tabel `a-table`: invoice_number, PO number, amount, status, due_date, payment info
    - Tombol "Buat Invoice" → modal: po_id (select dari PO yang sudah GRN), amount (auto-fill), due_date
    - Handle error `GRN_NOT_COMPLETED`, tampilkan status pembayaran (pending/paid)
    - Hanya role "supplier"
    - _Requirements: 15.1, 15.2, 15.3, 15.4, 15.5, 15.6_

- [x] 6. Web Dashboard: Router, Sidebar, dan Modifikasi Views Existing
  - [x] 6.1 Update `web/src/router/index.js` — tambah route baru dan default route supplier
    - Tambah routes: `/rab`, `/rab/:id`, `/invoices`, `/supplier-products`, `/supplier-dashboard`, `/supplier-products/manage`, `/supplier-po`, `/supplier-invoices`
    - Setiap route dengan meta.roles sesuai desain
    - Tambah `case 'supplier': return '/supplier-dashboard'` di `getDefaultRouteForRole()`
    - _Requirements: 17.1, 17.2, 17.3, 17.4, 17.5_

  - [x] 6.2 Update sidebar navigation (`HSidebar.vue` atau `MainLayout.vue`) — tambah menu groups baru
    - Tambah menu group "RAB & Pengadaan" dengan sub-menu: Daftar RAB, Invoice, Katalog Supplier (dengan roles masing-masing)
    - Tambah menu group "Supplier Portal" dengan sub-menu: Dashboard, Katalog Produk, Purchase Order, Invoice (hanya role supplier)
    - Modifikasi menu "Supplier" existing — tambah `kepala_yayasan` ke roles
    - Kedua menu group TIDAK ditandai `operational: true` karena diakses oleh non-operational roles
    - _Requirements: 17.6, 17.7, 17.8_

  - [x] 6.3 Modifikasi `web/src/views/SupplierListView.vue` — tambah akses kepala_yayasan
    - Tambah role `kepala_yayasan` ke akses halaman
    - Tambah kolom "Katalog Produk" yang menunjukkan jumlah produk supplier
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

  - [x] 6.4 Modifikasi `web/src/views/PurchaseOrderListView.vue` — tambah akses kepala_yayasan + field RAB
    - Tambah role `kepala_yayasan` ke akses halaman
    - Form PO baru: tambah field rab_id (select RAB approved_yayasan), target_sppg_id
    - Tampilkan kolom tambahan: rab_number, target_sppg, yayasan
    - Filter RAB items pending (po_id === null) untuk dipilih sebagai items PO
    - _Requirements: 10.1, 10.2, 10.3, 10.4_

  - [x] 6.5 Modifikasi `web/src/views/GoodsReceiptView.vue` — validasi PO-GRN 1:1 + quality rating
    - Warning jika PO yang dipilih sudah memiliki GRN
    - Handle error `PO_ALREADY_HAS_GRN`
    - Tambah field quality_rating (skala 1-5) menggunakan `a-rate`
    - _Requirements: 11.1, 11.2, 11.3, 11.4_

  - [ ]* 6.6 Tulis property test untuk filter RAB items pending (fast-check)
    - **Property 6: Filter RAB items pending hanya mengembalikan items tanpa PO** — generate daftar RAB items acak, verifikasi filter hanya mengembalikan items dengan `po_id === null`
    - **Validates: Requirements 10.4**

- [x] 7. Checkpoint — Pastikan semua halaman web dashboard berfungsi
  - Pastikan semua views, routes, sidebar, dan modifikasi existing tidak ada error. Tanyakan user jika ada pertanyaan.

- [x] 8. PWA Mobile: Supplier Portal Views
  - [x] 8.1 Buat `pwa/src/views/SupplierDashboardView.vue` — dashboard supplier mobile
    - Ringkasan menggunakan `van-grid` dan `van-cell-group`: total PO aktif, PO selesai, invoice pending, pembayaran diterima
    - Quick action buttons: "Lihat PO", "Buat Invoice", "Katalog Produk"
    - Panggil `getDashboard()` dari `supplierPortalService`
    - Hanya role "supplier"
    - _Requirements: 18.1, 18.2, 18.3, 18.4_

  - [x] 8.2 Buat `pwa/src/views/SupplierProductsView.vue` — katalog produk supplier mobile
    - Daftar produk menggunakan `van-cell-group` + `van-cell`: ingredient_name, unit_price, stock_quantity, is_available
    - FAB "Tambah Produk" → navigasi ke form
    - Klik produk → navigasi ke form edit
    - Swipe hapus menggunakan `van-swipe-cell`
    - Toggle ketersediaan menggunakan `van-switch` → langsung panggil API update
    - Hanya role "supplier"
    - _Requirements: 19.1, 19.2, 19.3, 19.4, 19.5, 19.6_

  - [x] 8.3 Buat `pwa/src/views/SupplierProductFormView.vue` — form tambah/edit produk supplier
    - `van-form` dengan field: ingredient picker, unit_price, min_order_qty, stock_quantity, is_available
    - Mode tambah dan edit (berdasarkan route param `:id`)
    - _Requirements: 19.2, 19.3_

  - [x] 8.4 Buat `pwa/src/views/SupplierPOListView.vue` — daftar PO supplier mobile
    - Daftar PO menggunakan `van-cell-group`: po_number, yayasan, tanggal, total_amount, status
    - Klik PO → navigasi ke detail
    - Pull-to-refresh menggunakan `van-pull-refresh`
    - Read-only, hanya role "supplier"
    - _Requirements: 20.1, 20.2, 20.3, 20.4, 20.5_

  - [x] 8.5 Buat `pwa/src/views/SupplierPODetailView.vue` — detail PO supplier mobile
    - Tampilkan detail PO + daftar items menggunakan `van-cell-group`
    - _Requirements: 20.2_

  - [x] 8.6 Buat `pwa/src/views/SupplierInvoiceView.vue` — invoice supplier mobile
    - Daftar invoice menggunakan `van-cell-group`: invoice_number, PO number, amount, status, due_date
    - Tombol "Buat Invoice" → `van-form`: po_id (picker PO yang sudah GRN), amount, due_date
    - Status badge: pending (orange), paid (green) menggunakan `van-tag`
    - Tampilkan info pembayaran saat status "paid"
    - Hanya role "supplier"
    - _Requirements: 21.1, 21.2, 21.3, 21.4, 21.5_

  - [x] 8.7 Buat `pwa/src/views/SupplierNotificationsView.vue` — notifikasi supplier mobile
    - Daftar notifikasi menggunakan `van-cell-group`: title, message, waktu, status baca/belum baca
    - Klik notifikasi → tandai sudah dibaca + navigasi ke halaman terkait
    - Badge jumlah notifikasi belum dibaca pada Bottom Tab
    - Integrasi Firebase Cloud Messaging (FCM) untuk push notification: PO baru, pembayaran diterima
    - Hanya role "supplier"
    - _Requirements: 22.1, 22.2, 22.3, 22.4, 22.5_

  - [x] 8.8 Buat `pwa/src/views/SupplierPaymentsView.vue` — riwayat pembayaran supplier mobile
    - Daftar pembayaran menggunakan `van-cell-group`: invoice_number, amount, payment_date, payment_method
    - Total pembayaran diterima di bagian atas
    - Pull-to-refresh menggunakan `van-pull-refresh`
    - Panggil `getPayments()` dari `supplierPortalService`
    - Hanya role "supplier"
    - _Requirements: 23.1, 23.2, 23.3, 23.4, 23.5_

- [x] 9. PWA Mobile: Router dan Bottom Tab Navigation
  - [x] 9.1 Update `pwa/src/router/index.js` — tambah route supplier dan default route
    - Tambah routes: `/supplier-dashboard`, `/supplier-products`, `/supplier-products/form`, `/supplier-products/:id/edit`, `/supplier-po`, `/supplier-po/:id`, `/supplier-invoices`, `/supplier-notifications`, `/supplier-payments`
    - Semua route dengan `meta: { roles: ['supplier'] }`
    - Tambah `if (roleStr === 'supplier') return '/supplier-dashboard'` di `getDefaultRoute()`
    - _Requirements: 24.1, 24.2, 24.3_

  - [x] 9.2 Update bottom tab navigation (`BottomNavigation.vue` atau `MobileLayout.vue`) — tambah tab supplier
    - Tambah konfigurasi navigasi supplier: Dashboard (Home), PO, Invoice, Notifikasi, Profil
    - Badge notifikasi belum dibaca pada tab Notifikasi menggunakan `van-badge`
    - _Requirements: 24.4, 24.5_

  - [ ]* 9.3 Tulis property test untuk konfigurasi navigasi bottom tab (fast-check)
    - **Property 7: Konfigurasi navigasi bottom tab per role** — generate role acak, verifikasi `ROLE_NAV_MAP[role]` memiliki minimal 2 item navigasi dengan properti lengkap
    - **Validates: Requirements 24.4**

- [x] 10. Final Checkpoint — Pastikan semua implementasi berfungsi
  - Pastikan semua tests pass, tidak ada error di kedua platform (web + PWA). Tanyakan user jika ada pertanyaan.

## Notes

- Tasks bertanda `*` bersifat opsional dan dapat dilewati untuk MVP lebih cepat
- Setiap task mereferensikan persyaratan spesifik untuk traceability
- Checkpoints memastikan validasi inkremental di setiap fase
- Property tests memvalidasi correctness properties universal dari design document
- Urutan implementasi: foundation → web views → PWA views → integrasi
