# Dokumen Persyaratan: Frontend RAB, Procurement & Supplier Portal

## Pendahuluan

Dokumen ini mendefinisikan persyaratan frontend untuk fitur RAB (Rencana Anggaran Belanja), procurement yayasan, dan supplier portal. Backend API sudah sepenuhnya diimplementasikan. Persyaratan ini mencakup dua platform: (1) Web Dashboard (`web/`) menggunakan Vue 3 + Ant Design Vue + Pinia + Vue Router, dan (2) PWA Mobile (`pwa/`) menggunakan Vue 3 + Vant UI + Pinia + Vue Router. Fitur ini menambahkan halaman-halaman baru untuk manajemen RAB, katalog supplier, invoice & pembayaran, serta portal supplier lengkap dengan dashboard dan notifikasi.

## Glosarium

- **Web_Dashboard**: Aplikasi web dashboard di direktori `web/` menggunakan Vue 3, Ant Design Vue, Pinia, dan Vue Router
- **PWA_Mobile**: Aplikasi PWA mobile di direktori `pwa/` menggunakan Vue 3, Vant UI, Pinia, dan Vue Router
- **RAB_Service**: Modul service layer di frontend (`rabService.js`) yang membungkus panggilan API untuk endpoint RAB
- **Supplier_Product_Service**: Modul service layer di frontend (`supplierProductService.js`) yang membungkus panggilan API untuk endpoint supplier product catalog
- **Invoice_Service**: Modul service layer di frontend (`invoiceService.js`) yang membungkus panggilan API untuk endpoint invoice dan pembayaran
- **Supplier_Portal_Service**: Modul service layer di PWA mobile (`supplierPortalService.js`) yang membungkus panggilan API untuk endpoint supplier dashboard dan payments
- **RAB_List_View**: Halaman daftar RAB di web dashboard
- **RAB_Detail_View**: Halaman detail RAB di web dashboard yang menampilkan items, status approval, dan PO tracking
- **Permission_System**: Sistem otorisasi frontend menggunakan `permissions.js` dan `usePermissions.js` composable
- **Router_Guard**: Middleware navigasi Vue Router yang memvalidasi autentikasi dan role sebelum mengizinkan akses halaman
- **Sidebar_Navigation**: Komponen navigasi samping di MainLayout.vue yang menampilkan menu berdasarkan role pengguna
- **Bottom_Tab_Navigation**: Komponen navigasi bawah di MobileLayout.vue untuk PWA mobile
- **kepala_sppg**: Role pengguna yang memimpin operasional SPPG — melihat RAB SPPG-nya, approve RAB, edit RAB draft
- **kepala_yayasan**: Role pengguna yang memimpin yayasan — approve/reject RAB, kelola supplier, buat PO, bayar invoice
- **supplier**: Role pengguna baru — akses portal supplier di web dan PWA mobile untuk kelola katalog, lihat PO, buat invoice
- **ahli_gizi**: Role pengguna yang membuat menu plan — melihat RAB yang dihasilkan dari menu plan-nya
- **asisten_lapangan**: Role pengguna yang membantu penerimaan barang — melihat RAB checklist dan GRN

## Persyaratan

### Persyaratan 1: API Service Layer untuk RAB

**User Story:** Sebagai developer frontend, saya ingin service layer yang membungkus semua endpoint RAB API, sehingga komponen Vue dapat memanggil API secara konsisten mengikuti pola existing (`api.js` + service modules).

#### Acceptance Criteria

1. THE RAB_Service SHALL menyediakan method `getRABList(params)` yang memanggil `GET /api/v1/rab` dengan parameter filter opsional (status, sppg_id, page, per_page)
2. THE RAB_Service SHALL menyediakan method `getRABDetail(id)` yang memanggil `GET /api/v1/rab/:id` dan mengembalikan data RAB beserta items dan status approval
3. THE RAB_Service SHALL menyediakan method `updateRAB(id, data)` yang memanggil `PUT /api/v1/rab/:id` untuk mengedit RAB
4. THE RAB_Service SHALL menyediakan method `approveSPPG(id)` yang memanggil `POST /api/v1/rab/:id/approve-sppg`
5. THE RAB_Service SHALL menyediakan method `approveYayasan(id)` yang memanggil `POST /api/v1/rab/:id/approve-yayasan`
6. THE RAB_Service SHALL menyediakan method `rejectRAB(id, data)` yang memanggil `POST /api/v1/rab/:id/reject` dengan body berisi `revision_notes`
7. THE RAB_Service SHALL menyediakan method `resubmitRAB(id)` yang memanggil `POST /api/v1/rab/:id/resubmit`
8. THE RAB_Service SHALL menyediakan method `getRABComparison(id)` yang memanggil `GET /api/v1/rab/:id/comparison`
9. THE RAB_Service SHALL menyediakan method `getRABPOTracking(id)` yang memanggil `GET /api/v1/rab/:id/po-tracking`
10. THE RAB_Service SHALL menggunakan instance axios dari `api.js` yang sudah dikonfigurasi dengan base URL dan JWT interceptor

### Persyaratan 2: API Service Layer untuk Supplier Product Catalog

**User Story:** Sebagai developer frontend, saya ingin service layer untuk endpoint supplier product catalog, sehingga halaman katalog produk supplier dapat berinteraksi dengan API.

#### Acceptance Criteria

1. THE Supplier_Product_Service SHALL menyediakan method `getProducts(params)` yang memanggil `GET /api/v1/supplier-products` dengan parameter filter opsional
2. THE Supplier_Product_Service SHALL menyediakan method `createProduct(data)` yang memanggil `POST /api/v1/supplier-products`
3. THE Supplier_Product_Service SHALL menyediakan method `updateProduct(id, data)` yang memanggil `PUT /api/v1/supplier-products/:id`
4. THE Supplier_Product_Service SHALL menyediakan method `deleteProduct(id)` yang memanggil `DELETE /api/v1/supplier-products/:id`
5. THE Supplier_Product_Service SHALL menggunakan instance axios dari `api.js` yang sudah dikonfigurasi dengan base URL dan JWT interceptor

### Persyaratan 3: API Service Layer untuk Invoice & Pembayaran

**User Story:** Sebagai developer frontend, saya ingin service layer untuk endpoint invoice dan pembayaran, sehingga halaman invoice dapat berinteraksi dengan API.

#### Acceptance Criteria

1. THE Invoice_Service SHALL menyediakan method `getInvoices(params)` yang memanggil `GET /api/v1/invoices` dengan parameter filter opsional (status, page, per_page)
2. THE Invoice_Service SHALL menyediakan method `createInvoice(data)` yang memanggil `POST /api/v1/invoices`
3. THE Invoice_Service SHALL menyediakan method `getInvoiceDetail(id)` yang memanggil `GET /api/v1/invoices/:id`
4. THE Invoice_Service SHALL menyediakan method `payInvoice(id, data)` yang memanggil `POST /api/v1/invoices/:id/pay`
5. THE Invoice_Service SHALL menyediakan method `uploadPaymentProof(id, formData)` yang memanggil `POST /api/v1/invoices/:id/upload-proof` dengan header `Content-Type: multipart/form-data`
6. THE Invoice_Service SHALL menggunakan instance axios dari `api.js` yang sudah dikonfigurasi dengan base URL dan JWT interceptor

### Persyaratan 4: API Service Layer untuk Supplier Portal (PWA Mobile)

**User Story:** Sebagai developer frontend, saya ingin service layer di PWA mobile untuk endpoint supplier dashboard dan payments, sehingga halaman supplier portal di mobile dapat berinteraksi dengan API.

#### Acceptance Criteria

1. THE Supplier_Portal_Service SHALL menyediakan method `getDashboard()` yang memanggil `GET /api/v1/supplier/dashboard`
2. THE Supplier_Portal_Service SHALL menyediakan method `getPayments(params)` yang memanggil `GET /api/v1/supplier/payments`
3. THE Supplier_Portal_Service SHALL menyediakan method `getProducts(params)` yang memanggil `GET /api/v1/supplier-products`
4. THE Supplier_Portal_Service SHALL menyediakan method `createProduct(data)` yang memanggil `POST /api/v1/supplier-products`
5. THE Supplier_Portal_Service SHALL menyediakan method `updateProduct(id, data)` yang memanggil `PUT /api/v1/supplier-products/:id`
6. THE Supplier_Portal_Service SHALL menyediakan method `deleteProduct(id)` yang memanggil `DELETE /api/v1/supplier-products/:id`
7. THE Supplier_Portal_Service SHALL menyediakan method `getInvoices(params)` yang memanggil `GET /api/v1/invoices`
8. THE Supplier_Portal_Service SHALL menyediakan method `createInvoice(data)` yang memanggil `POST /api/v1/invoices`
9. THE Supplier_Portal_Service SHALL menyediakan method `getPurchaseOrders(params)` yang memanggil `GET /api/v1/purchase-orders`
10. THE Supplier_Portal_Service SHALL menggunakan instance axios dari `api.js` PWA yang sudah dikonfigurasi dengan base URL dan JWT interceptor


### Persyaratan 5: Halaman Daftar RAB (Web Dashboard)

**User Story:** Sebagai kepala_sppg atau kepala_yayasan, saya ingin melihat daftar RAB dengan filter dan status, sehingga saya dapat memantau semua RAB yang relevan dengan peran saya.

#### Acceptance Criteria

1. THE RAB_List_View SHALL menampilkan tabel RAB menggunakan komponen `a-table` dari Ant Design Vue dengan kolom: rab_number, menu_plan (nama/periode), SPPG, status, total_amount, tanggal dibuat, dan aksi
2. THE RAB_List_View SHALL menyediakan filter berdasarkan status RAB (draft, approved_sppg, approved_yayasan, revision_requested, completed) menggunakan komponen `a-select`
3. WHEN pengguna dengan role kepala_sppg mengakses halaman, THE RAB_List_View SHALL hanya menampilkan RAB dari SPPG miliknya (backend sudah memfilter berdasarkan tenant scope)
4. WHEN pengguna dengan role kepala_yayasan mengakses halaman, THE RAB_List_View SHALL menampilkan RAB dari semua SPPG di bawah yayasannya
5. THE RAB_List_View SHALL menampilkan status RAB menggunakan `a-tag` dengan warna berbeda per status: draft (default), approved_sppg (blue), approved_yayasan (green), revision_requested (orange), completed (purple)
6. WHEN pengguna mengklik baris RAB, THE RAB_List_View SHALL menavigasi ke halaman RAB_Detail_View
7. THE RAB_List_View SHALL mendukung pagination server-side menggunakan parameter `page` dan `per_page`
8. THE RAB_List_View SHALL menampilkan total_amount dalam format mata uang Rupiah (Rp) dengan pemisah ribuan

### Persyaratan 6: Halaman Detail RAB (Web Dashboard)

**User Story:** Sebagai kepala_sppg atau kepala_yayasan, saya ingin melihat detail RAB lengkap dengan items, status approval, dan tracking PO, sehingga saya dapat mengambil keputusan yang tepat.

#### Acceptance Criteria

1. THE RAB_Detail_View SHALL menampilkan informasi header RAB: rab_number, menu plan terkait, SPPG, status, total_amount, created_by, dan tanggal
2. THE RAB_Detail_View SHALL menampilkan tabel RAB items menggunakan `a-table` dengan kolom: ingredient, quantity, unit, unit_price, subtotal, recommended_supplier, status item (pending/po_created/grn_received), po_id, dan grn_id
3. WHEN RAB berstatus "draft" dan pengguna adalah kepala_sppg, THE RAB_Detail_View SHALL menampilkan tombol "Approve SPPG" yang memanggil `POST /api/v1/rab/:id/approve-sppg`
4. WHEN RAB berstatus "approved_sppg" dan pengguna adalah kepala_yayasan, THE RAB_Detail_View SHALL menampilkan tombol "Approve Yayasan" dan tombol "Tolak" yang masing-masing memanggil endpoint approve-yayasan dan reject
5. WHEN kepala_yayasan mengklik tombol "Tolak", THE RAB_Detail_View SHALL menampilkan modal `a-modal` dengan textarea untuk mengisi alasan penolakan (revision_notes) sebelum mengirim request reject
6. WHEN RAB berstatus "draft" atau "revision_requested", THE RAB_Detail_View SHALL menampilkan tombol "Edit RAB" yang membuka form edit inline atau modal untuk mengubah item RAB (quantity, unit_price)
7. WHEN RAB berstatus "revision_requested", THE RAB_Detail_View SHALL menampilkan revision_notes dari kepala_yayasan dalam `a-alert` berwarna warning
8. WHEN RAB berstatus "revision_requested" dan pengguna adalah kepala_sppg, THE RAB_Detail_View SHALL menampilkan tombol "Kirim Ulang" yang memanggil `POST /api/v1/rab/:id/resubmit`
9. THE RAB_Detail_View SHALL menampilkan tab "PO Tracking" yang memanggil `GET /api/v1/rab/:id/po-tracking` dan menampilkan daftar PO terkait beserta status GRN
10. THE RAB_Detail_View SHALL menampilkan tab "Perbandingan" yang memanggil `GET /api/v1/rab/:id/comparison` dan menampilkan tabel RAB vs Aktual (planned_amount vs actual_amount per item)
11. IF operasi approve/reject/resubmit gagal, THEN THE RAB_Detail_View SHALL menampilkan pesan error dari API menggunakan `message.error()` dari Ant Design Vue
12. WHEN operasi approve/reject/resubmit berhasil, THE RAB_Detail_View SHALL menampilkan pesan sukses menggunakan `message.success()` dan me-refresh data RAB

### Persyaratan 7: Halaman Daftar Invoice (Web Dashboard)

**User Story:** Sebagai kepala_yayasan, saya ingin melihat daftar invoice dari supplier dan memproses pembayaran, sehingga alur keuangan pengadaan tercatat dengan baik.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman InvoiceListView yang menampilkan tabel invoice menggunakan `a-table` dengan kolom: invoice_number, supplier, PO number, amount, status, due_date, dan aksi
2. THE InvoiceListView SHALL menyediakan filter berdasarkan status invoice (pending, paid) menggunakan `a-select`
3. WHEN pengguna dengan role kepala_yayasan mengakses halaman, THE InvoiceListView SHALL menampilkan invoice dari yayasannya (backend memfilter berdasarkan yayasan_id)
4. WHEN kepala_yayasan mengklik tombol "Bayar" pada invoice berstatus "pending", THE InvoiceListView SHALL menampilkan modal pembayaran dengan form: payment_date (date picker), payment_method (default: bank_transfer), dan upload bukti transfer (file upload)
5. WHEN bukti transfer diupload, THE Web_Dashboard SHALL memanggil `POST /api/v1/invoices/:id/upload-proof` dengan `Content-Type: multipart/form-data`
6. WHEN pembayaran berhasil dicatat, THE InvoiceListView SHALL memanggil `POST /api/v1/invoices/:id/pay` dan menampilkan pesan sukses
7. THE InvoiceListView SHALL menampilkan amount dalam format mata uang Rupiah (Rp) dengan pemisah ribuan
8. IF pembayaran gagal (misalnya invoice sudah dibayar), THEN THE InvoiceListView SHALL menampilkan pesan error dari API

### Persyaratan 8: Modifikasi Halaman Supplier (Web Dashboard)

**User Story:** Sebagai kepala_yayasan, saya ingin mengelola supplier di level yayasan, sehingga semua SPPG di bawah yayasan saya menggunakan database supplier yang sama.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL memodifikasi SupplierListView agar dapat diakses oleh role kepala_yayasan selain kepala_sppg dan pengadaan
2. WHEN pengguna dengan role kepala_yayasan mengakses halaman supplier, THE SupplierListView SHALL menampilkan supplier yang terhubung dengan yayasannya (backend memfilter via YayasanTenantScope)
3. THE SupplierListView SHALL mengizinkan kepala_yayasan melakukan CRUD supplier (tambah, edit, hapus) menggunakan form modal yang sudah ada
4. THE SupplierListView SHALL menampilkan kolom tambahan "Katalog Produk" yang menunjukkan jumlah produk yang ditawarkan supplier

### Persyaratan 9: Halaman Katalog Produk Supplier (Web Dashboard)

**User Story:** Sebagai kepala_yayasan, saya ingin melihat katalog produk dari supplier yang terhubung dengan yayasan saya, sehingga saya dapat memilih supplier terbaik untuk pengadaan.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman SupplierProductListView yang menampilkan tabel produk supplier menggunakan `a-table` dengan kolom: supplier_name, ingredient_name, unit_price, min_order_qty, stock_quantity, is_available
2. WHEN pengguna dengan role kepala_yayasan mengakses halaman, THE SupplierProductListView SHALL menampilkan produk dari semua supplier yang terhubung dengan yayasannya
3. THE SupplierProductListView SHALL menyediakan filter berdasarkan supplier dan ingredient menggunakan `a-select`
4. THE SupplierProductListView SHALL menampilkan unit_price dalam format mata uang Rupiah (Rp) dengan pemisah ribuan
5. THE SupplierProductListView SHALL menampilkan status ketersediaan menggunakan `a-tag` (hijau untuk tersedia, merah untuk tidak tersedia)

### Persyaratan 10: Modifikasi Halaman Purchase Order (Web Dashboard)

**User Story:** Sebagai kepala_yayasan, saya ingin membuat Purchase Order dari RAB yang sudah disetujui, sehingga pengadaan bahan baku dapat dilakukan secara terpusat.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL memodifikasi PurchaseOrderListView agar dapat diakses oleh role kepala_yayasan
2. WHEN kepala_yayasan membuat PO baru, THE PurchaseOrderListView SHALL menampilkan form yang mewajibkan field: rab_id (pilih dari RAB berstatus approved_yayasan), supplier_id (pilih dari supplier yayasan), target_sppg_id (pilih SPPG tujuan pengiriman), dan items
3. THE PurchaseOrderListView SHALL menampilkan kolom tambahan: rab_number, target_sppg (nama SPPG tujuan), dan yayasan pada tabel PO
4. WHEN kepala_yayasan memilih RAB pada form PO, THE Web_Dashboard SHALL menampilkan daftar RAB items yang belum memiliki PO (status "pending") untuk dipilih sebagai items PO

### Persyaratan 11: Modifikasi Halaman Goods Receipt (Web Dashboard)

**User Story:** Sebagai kepala_sppg atau asisten_lapangan, saya ingin form GRN memvalidasi relasi 1:1 dengan PO, sehingga tidak ada penerimaan ganda untuk satu PO.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL memodifikasi GoodsReceiptView agar menampilkan warning jika PO yang dipilih sudah memiliki GRN
2. WHEN pengguna memilih PO pada form GRN, THE GoodsReceiptView SHALL mengecek apakah PO sudah memiliki GRN dan menampilkan pesan "PO ini sudah memiliki GRN" jika ya
3. THE GoodsReceiptView SHALL menambahkan field quality_rating (skala 1-5) menggunakan komponen `a-rate` pada form GRN
4. IF backend mengembalikan error `PO_ALREADY_HAS_GRN`, THEN THE GoodsReceiptView SHALL menampilkan pesan error yang jelas kepada pengguna


### Persyaratan 12: Halaman Supplier Dashboard (Web Dashboard)

**User Story:** Sebagai supplier, saya ingin melihat ringkasan pesanan dan pembayaran saya di web dashboard, sehingga saya dapat memantau bisnis saya dengan yayasan.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman SupplierDashboardView yang menampilkan ringkasan: total PO aktif, total PO selesai, total invoice pending, dan total pembayaran diterima menggunakan komponen `a-statistic` atau `a-card`
2. THE SupplierDashboardView SHALL memanggil `GET /api/v1/supplier/dashboard` untuk mendapatkan data ringkasan
3. THE SupplierDashboardView SHALL menampilkan daftar PO terbaru (5 terakhir) dalam tabel ringkas
4. THE SupplierDashboardView SHALL menampilkan daftar invoice terbaru (5 terakhir) dalam tabel ringkas
5. THE SupplierDashboardView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 13: Halaman Katalog Produk Supplier — CRUD oleh Supplier (Web Dashboard)

**User Story:** Sebagai supplier, saya ingin mengelola katalog produk saya (tambah, edit, hapus) di web dashboard, sehingga yayasan dapat melihat produk terbaru yang saya tawarkan.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman SupplierProductManageView yang menampilkan tabel produk milik supplier menggunakan `a-table` dengan kolom: ingredient_name, unit_price, min_order_qty, stock_quantity, is_available, dan aksi (edit, hapus)
2. THE SupplierProductManageView SHALL menyediakan tombol "Tambah Produk" yang membuka modal `a-modal` dengan form: ingredient_id (select dari daftar ingredient), unit_price, min_order_qty, stock_quantity, is_available (switch)
3. WHEN supplier mengklik tombol "Edit" pada produk, THE SupplierProductManageView SHALL membuka modal edit dengan data produk yang sudah terisi
4. WHEN supplier mengklik tombol "Hapus" pada produk, THE SupplierProductManageView SHALL menampilkan konfirmasi `a-popconfirm` sebelum menghapus
5. IF backend mengembalikan error `DUPLICATE_SUPPLIER_PRODUCT`, THEN THE SupplierProductManageView SHALL menampilkan pesan "Produk untuk ingredient ini sudah ada"
6. THE SupplierProductManageView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 14: Halaman PO List untuk Supplier (Web Dashboard)

**User Story:** Sebagai supplier, saya ingin melihat daftar Purchase Order yang ditujukan kepada saya, sehingga saya dapat memproses pesanan.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman SupplierPOListView yang menampilkan tabel PO menggunakan `a-table` dengan kolom: po_number, yayasan, target_sppg, tanggal, total_amount, status, dan aksi
2. THE SupplierPOListView SHALL bersifat read-only (supplier tidak dapat mengedit PO)
3. WHEN supplier mengklik baris PO, THE SupplierPOListView SHALL menampilkan detail PO dalam modal atau halaman detail dengan daftar items
4. THE SupplierPOListView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 15: Halaman Invoice untuk Supplier (Web Dashboard)

**User Story:** Sebagai supplier, saya ingin membuat invoice untuk PO yang sudah diterima (GRN selesai), sehingga saya dapat menagih pembayaran ke yayasan.

#### Acceptance Criteria

1. THE Web_Dashboard SHALL menyediakan halaman SupplierInvoiceView yang menampilkan tabel invoice milik supplier menggunakan `a-table` dengan kolom: invoice_number, PO number, amount, status, due_date, dan payment info
2. THE SupplierInvoiceView SHALL menyediakan tombol "Buat Invoice" yang membuka modal dengan form: po_id (select dari PO yang sudah GRN), amount (auto-fill dari PO total_amount), due_date
3. WHEN invoice berhasil dibuat, THE SupplierInvoiceView SHALL menampilkan pesan sukses dan me-refresh daftar invoice
4. IF backend mengembalikan error `GRN_NOT_COMPLETED`, THEN THE SupplierInvoiceView SHALL menampilkan pesan "GRN belum selesai untuk PO ini"
5. THE SupplierInvoiceView SHALL menampilkan status pembayaran: "Menunggu Pembayaran" (pending) atau "Sudah Dibayar" (paid) dengan tanggal dan bukti transfer
6. THE SupplierInvoiceView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 16: Update Permission System (Web Dashboard)

**User Story:** Sebagai developer frontend, saya ingin permission system diperbarui dengan permission baru untuk RAB, supplier portal, dan invoice, sehingga akses halaman terkontrol berdasarkan role.

#### Acceptance Criteria

1. THE Permission_System SHALL menambahkan permission baru di `permissions.js`: RAB_VIEW (kepala_sppg, ahli_gizi, kepala_yayasan), RAB_APPROVE_SPPG (kepala_sppg), RAB_APPROVE_YAYASAN (kepala_yayasan), RAB_EDIT (kepala_sppg)
2. THE Permission_System SHALL menambahkan permission: INVOICE_VIEW (kepala_yayasan, supplier), INVOICE_CREATE (supplier), INVOICE_PAY (kepala_yayasan)
3. THE Permission_System SHALL menambahkan permission: SUPPLIER_PRODUCT_VIEW (kepala_yayasan, supplier), SUPPLIER_PRODUCT_MANAGE (supplier)
4. THE Permission_System SHALL menambahkan permission: SUPPLIER_PORTAL_VIEW (supplier), SUPPLIER_DASHBOARD_VIEW (supplier)
5. THE Permission_System SHALL memodifikasi permission SUPPLIER_VIEW dan SUPPLIER_MANAGE agar menyertakan role kepala_yayasan
6. THE Permission_System SHALL memodifikasi permission PO_VIEW dan PO_CREATE agar menyertakan role kepala_yayasan
7. THE Permission_System SHALL menambahkan role "supplier" ke fungsi `getRoleLabel()` dengan label "Supplier"

### Persyaratan 17: Update Router dan Navigasi (Web Dashboard)

**User Story:** Sebagai pengguna, saya ingin dapat mengakses halaman-halaman baru melalui navigasi sidebar dan URL, sehingga saya dapat menemukan fitur RAB, invoice, dan supplier portal dengan mudah.

#### Acceptance Criteria

1. THE Router_Guard SHALL menambahkan route baru: `/rab` (RABListView, roles: kepala_sppg, ahli_gizi, kepala_yayasan), `/rab/:id` (RABDetailView, roles: kepala_sppg, ahli_gizi, kepala_yayasan)
2. THE Router_Guard SHALL menambahkan route baru: `/invoices` (InvoiceListView, roles: kepala_yayasan, supplier)
3. THE Router_Guard SHALL menambahkan route baru: `/supplier-products` (SupplierProductListView, roles: kepala_yayasan), `/supplier-products/manage` (SupplierProductManageView, roles: supplier)
4. THE Router_Guard SHALL menambahkan route baru: `/supplier-dashboard` (SupplierDashboardView, roles: supplier), `/supplier-po` (SupplierPOListView, roles: supplier), `/supplier-invoices` (SupplierInvoiceView, roles: supplier)
5. THE Router_Guard SHALL menambahkan default route untuk role "supplier" ke `/supplier-dashboard` pada fungsi `getDefaultRouteForRole()`
6. THE Sidebar_Navigation SHALL menambahkan menu group "RAB & Pengadaan" dengan sub-menu: Daftar RAB (untuk kepala_sppg, ahli_gizi, kepala_yayasan), Invoice (untuk kepala_yayasan), Katalog Supplier (untuk kepala_yayasan)
7. THE Sidebar_Navigation SHALL menambahkan menu group "Supplier Portal" dengan sub-menu: Dashboard, Katalog Produk, Purchase Order, Invoice, yang hanya terlihat oleh role "supplier"
8. THE Sidebar_Navigation SHALL memodifikasi menu "Supplier" existing agar juga terlihat oleh role kepala_yayasan

### Persyaratan 18: Supplier Dashboard PWA Mobile

**User Story:** Sebagai supplier, saya ingin mengakses dashboard ringkasan di aplikasi mobile, sehingga saya dapat memantau pesanan dan pembayaran kapan saja.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierDashboardView yang menampilkan ringkasan menggunakan komponen Vant UI (`van-grid`, `van-cell-group`): total PO aktif, total PO selesai, total invoice pending, total pembayaran diterima
2. THE SupplierDashboardView SHALL memanggil `GET /api/v1/supplier/dashboard` untuk mendapatkan data ringkasan
3. THE SupplierDashboardView SHALL menampilkan quick action buttons: "Lihat PO", "Buat Invoice", "Katalog Produk"
4. THE SupplierDashboardView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 19: Halaman Katalog Produk Supplier PWA Mobile

**User Story:** Sebagai supplier, saya ingin mengelola katalog produk saya di aplikasi mobile, sehingga saya dapat memperbarui harga dan ketersediaan produk dari mana saja.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierProductsView yang menampilkan daftar produk menggunakan `van-cell-group` dan `van-cell` dengan informasi: ingredient_name, unit_price, stock_quantity, is_available
2. THE SupplierProductsView SHALL menyediakan tombol floating action "Tambah Produk" yang menavigasi ke form tambah produk
3. WHEN supplier mengklik produk, THE SupplierProductsView SHALL menavigasi ke form edit produk dengan data yang sudah terisi
4. THE SupplierProductsView SHALL menyediakan swipe action untuk menghapus produk menggunakan `van-swipe-cell`
5. THE SupplierProductsView SHALL menampilkan toggle ketersediaan menggunakan `van-switch` yang langsung memanggil API update
6. THE SupplierProductsView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 20: Halaman PO List Supplier PWA Mobile

**User Story:** Sebagai supplier, saya ingin melihat daftar Purchase Order yang masuk di aplikasi mobile, sehingga saya dapat memproses pesanan dengan cepat.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierPOListView yang menampilkan daftar PO menggunakan `van-cell-group` dengan informasi: po_number, yayasan, tanggal, total_amount, status
2. WHEN supplier mengklik PO, THE SupplierPOListView SHALL menampilkan detail PO dalam halaman baru dengan daftar items
3. THE SupplierPOListView SHALL bersifat read-only
4. THE SupplierPOListView SHALL mendukung pull-to-refresh menggunakan `van-pull-refresh`
5. THE SupplierPOListView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 21: Halaman Invoice Supplier PWA Mobile

**User Story:** Sebagai supplier, saya ingin membuat invoice dan melihat status pembayaran di aplikasi mobile, sehingga saya dapat menagih pembayaran dengan mudah.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierInvoiceView yang menampilkan daftar invoice menggunakan `van-cell-group` dengan informasi: invoice_number, PO number, amount, status, due_date
2. THE SupplierInvoiceView SHALL menyediakan tombol "Buat Invoice" yang membuka form menggunakan `van-form` dengan field: po_id (picker dari PO yang sudah GRN), amount, due_date
3. THE SupplierInvoiceView SHALL menampilkan status pembayaran dengan badge warna: pending (orange), paid (green)
4. WHEN invoice berstatus "paid", THE SupplierInvoiceView SHALL menampilkan informasi pembayaran: tanggal bayar dan link bukti transfer
5. THE SupplierInvoiceView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 22: Halaman Notifikasi Supplier PWA Mobile

**User Story:** Sebagai supplier, saya ingin melihat notifikasi PO baru dan pembayaran di aplikasi mobile, sehingga saya dapat merespons dengan cepat.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierNotificationsView yang menampilkan daftar notifikasi menggunakan `van-cell-group` dengan informasi: title, message, waktu, status baca/belum baca
2. WHEN supplier mengklik notifikasi, THE SupplierNotificationsView SHALL menandai notifikasi sebagai sudah dibaca dan menavigasi ke halaman terkait (PO detail atau invoice)
3. THE SupplierNotificationsView SHALL menampilkan badge jumlah notifikasi belum dibaca pada ikon notifikasi di Bottom_Tab_Navigation
4. THE PWA_Mobile SHALL menerima push notification melalui Firebase Cloud Messaging (FCM) untuk event: PO baru dan pembayaran diterima
5. THE SupplierNotificationsView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 23: Halaman Riwayat Pembayaran Supplier PWA Mobile

**User Story:** Sebagai supplier, saya ingin melihat riwayat pembayaran yang sudah diterima di aplikasi mobile, sehingga saya dapat melacak arus kas dari yayasan.

#### Acceptance Criteria

1. THE PWA_Mobile SHALL menyediakan halaman SupplierPaymentsView yang menampilkan daftar pembayaran menggunakan `van-cell-group` dengan informasi: invoice_number, amount, payment_date, payment_method
2. THE SupplierPaymentsView SHALL memanggil `GET /api/v1/supplier/payments` untuk mendapatkan data riwayat pembayaran
3. THE SupplierPaymentsView SHALL menampilkan total pembayaran diterima di bagian atas halaman
4. THE SupplierPaymentsView SHALL mendukung pull-to-refresh menggunakan `van-pull-refresh`
5. THE SupplierPaymentsView SHALL hanya dapat diakses oleh pengguna dengan role "supplier"

### Persyaratan 24: Update Router dan Navigasi PWA Mobile

**User Story:** Sebagai supplier, saya ingin navigasi mobile yang intuitif dengan bottom tab khusus supplier, sehingga saya dapat berpindah antar fitur dengan mudah.

#### Acceptance Criteria

1. THE PWA_Mobile Router SHALL menambahkan route baru: `/supplier-dashboard` (SupplierDashboardView), `/supplier-products` (SupplierProductsView), `/supplier-po` (SupplierPOListView), `/supplier-invoices` (SupplierInvoiceView), `/supplier-notifications` (SupplierNotificationsView), `/supplier-payments` (SupplierPaymentsView)
2. THE PWA_Mobile Router SHALL menambahkan role guard "supplier" pada semua route supplier
3. THE PWA_Mobile Router SHALL menambahkan default route untuk role "supplier" ke `/supplier-dashboard` pada fungsi `getDefaultRoute()`
4. THE Bottom_Tab_Navigation SHALL menampilkan tab khusus untuk role "supplier": Dashboard, PO, Invoice, Notifikasi, Profil
5. THE Bottom_Tab_Navigation SHALL menampilkan badge notifikasi belum dibaca pada tab Notifikasi

### Persyaratan 25: Update Auth Store untuk Role Supplier

**User Story:** Sebagai developer frontend, saya ingin auth store mendukung role supplier dengan field supplier_id, sehingga komponen dapat mengecek role dan scope data supplier.

#### Acceptance Criteria

1. THE Web_Dashboard auth store SHALL menambahkan computed property `supplierId` yang mengambil `supplier_id` dari user data
2. THE Web_Dashboard auth store SHALL menambahkan computed property `isSupplier` yang mengecek `role === 'supplier'`
3. THE PWA_Mobile auth store SHALL menambahkan computed property `supplierId` yang mengambil `supplier_id` dari user data
4. THE PWA_Mobile auth store SHALL menambahkan computed property `isSupplier` yang mengecek `userRole === 'supplier'`
5. THE Permission_System SHALL menambahkan "supplier" ke daftar `NON_OPERATIONAL_ROLES` karena supplier bukan role operasional SPPG
