# Dokumen Persyaratan: RAB, Procurement Yayasan & Supplier Portal

## Pendahuluan

Fitur ini mengubah secara fundamental alur pengadaan (procurement) dalam sistem ERP SPPG. Perubahan utama meliputi: (1) auto-generasi RAB (Rencana Anggaran Belanja) dari menu plan yang dikonfirmasi, (2) alur persetujuan multi-level (Kepala SPPG → Kepala Yayasan), (3) pemindahan kepemilikan supplier dan pembuatan PO dari SPPG ke Yayasan, (4) portal supplier baru dengan katalog produk, notifikasi, dan invoice, (5) modifikasi GRN agar validasi terhadap PO dan update checklist RAB, serta (6) alur pembayaran invoice oleh yayasan.

## Glosarium

- **Sistem**: Aplikasi ERP SPPG secara keseluruhan (backend + frontend + PWA mobile)
- **RAB_Generator**: Komponen sistem yang bertanggung jawab membuat RAB secara otomatis dari menu plan
- **RAB**: Rencana Anggaran Belanja — dokumen perencanaan biaya belanja bahan baku berdasarkan menu plan
- **RAB_Item**: Baris item dalam RAB yang merepresentasikan satu kebutuhan bahan baku beserta harga dan rekomendasi supplier
- **Menu_Plan**: Rencana menu mingguan yang dibuat oleh ahli_gizi atau kepala_sppg
- **Supplier_Catalog**: Daftar produk yang dijual oleh supplier beserta harga, stok, dan ketersediaan
- **Supplier_Product**: Satu entri produk dalam Supplier_Catalog (ingredient + harga + ketersediaan)
- **Approval_Engine**: Komponen sistem yang mengelola alur persetujuan multi-level untuk RAB
- **PO_Manager**: Komponen sistem yang mengelola pembuatan dan pengelolaan Purchase Order oleh yayasan
- **Supplier_Portal**: Antarmuka web dan PWA mobile untuk role supplier
- **GRN_Processor**: Komponen sistem yang memproses penerimaan barang dan validasi terhadap PO
- **Invoice_Manager**: Komponen sistem yang mengelola pembuatan invoice oleh supplier dan pembayaran oleh yayasan
- **Notification_Service**: Komponen sistem yang mengirim notifikasi push (FCM) dan in-app notification
- **Supplier_Yayasan**: Tabel junction yang menghubungkan supplier dengan yayasan (relasi many-to-many)
- **kepala_sppg**: Role pengguna yang memimpin operasional SPPG
- **kepala_yayasan**: Role pengguna yang memimpin yayasan, sekarang memiliki izin tulis pada supplier, PO, persetujuan RAB, dan invoice
- **supplier**: Role pengguna baru yang dapat login ke web dan PWA mobile untuk mengelola katalog dan menerima PO
- **ahli_gizi**: Role pengguna yang membuat dan mengkonfirmasi menu plan
- **pengadaan**: Role pengguna yang sebelumnya mengelola procurement di level SPPG
- **YayasanTenantScope**: Pola tenant isolation baru untuk entitas yang dimiliki yayasan (supplier, PO, RAB approval, invoice)

## Persyaratan

### Persyaratan 1: Auto-Generasi RAB dari Menu Plan

**User Story:** Sebagai kepala_sppg, saya ingin RAB dibuat secara otomatis ketika menu plan dikonfirmasi, sehingga saya tidak perlu menghitung kebutuhan bahan baku dan biaya secara manual.

#### Acceptance Criteria

1. WHEN menu plan berstatus "approved" oleh ahli_gizi atau kepala_sppg, THE RAB_Generator SHALL membuat satu RAB baru yang terhubung ke menu plan tersebut
2. THE RAB_Generator SHALL menghitung kebutuhan bahan baku dari setiap recipe item dalam menu plan dan mengagregasi total kuantitas per ingredient
3. WHEN menghitung harga per ingredient, THE RAB_Generator SHALL mengambil harga dari Supplier_Catalog dan merekomendasikan supplier dengan harga termurah yang memiliki quality_rating tertinggi
4. THE RAB SHALL menyimpan rab_number (format: RAB-YYYYMMDD-XXXX), menu_plan_id, sppg_id, status ("draft"), total_amount, dan created_by
5. THE RAB_Item SHALL menyimpan rab_id, ingredient_id, quantity, unit, unit_price, subtotal, recommended_supplier_id, dan status ("pending")
6. IF tidak ada Supplier_Product yang tersedia untuk suatu ingredient, THEN THE RAB_Generator SHALL menandai RAB_Item tersebut dengan recommended_supplier_id NULL dan unit_price 0
7. WHEN RAB berhasil dibuat, THE Notification_Service SHALL mengirim notifikasi in-app ke kepala_sppg dari SPPG terkait

### Persyaratan 2: Alur Persetujuan Multi-Level RAB

**User Story:** Sebagai kepala_yayasan, saya ingin mereview dan menyetujui RAB beserta menu plan sebelum pengadaan dilakukan, sehingga anggaran yayasan terkontrol.

#### Acceptance Criteria

1. WHEN kepala_sppg menyetujui RAB, THE Approval_Engine SHALL mengubah status RAB menjadi "approved_sppg" dan mengirim RAB ke kepala_yayasan untuk review
2. WHEN kepala_yayasan menyetujui RAB, THE Approval_Engine SHALL mengubah status RAB menjadi "approved_yayasan" dan RAB siap untuk pembuatan PO
3. WHEN kepala_yayasan menolak RAB, THE Approval_Engine SHALL mengubah status RAB menjadi "revision_requested" dan menyimpan alasan penolakan (revision_notes)
4. WHEN RAB berstatus "revision_requested", THE Sistem SHALL mengharuskan SPPG merevisi menu plan DAN RAB, kemudian mengirim ulang untuk persetujuan
5. THE Approval_Engine SHALL mencatat setiap perubahan status RAB dalam audit trail (user_id, timestamp, from_status, to_status, notes)
6. WHILE RAB berstatus "draft" atau "revision_requested", THE Sistem SHALL mengizinkan kepala_sppg untuk mengedit RAB dan menu plan terkait
7. WHILE RAB berstatus "approved_sppg", THE Sistem SHALL melarang perubahan pada RAB dan menu plan terkait oleh SPPG
8. WHEN status RAB berubah, THE Notification_Service SHALL mengirim notifikasi in-app dan push notification (FCM) ke pihak terkait (kepala_sppg atau kepala_yayasan)

### Persyaratan 3: Migrasi Kepemilikan Supplier ke Yayasan

**User Story:** Sebagai kepala_yayasan, saya ingin mengelola database supplier di level yayasan, sehingga semua SPPG di bawah yayasan menggunakan supplier yang sama dan pengadaan terpusat.

#### Acceptance Criteria

1. THE Sistem SHALL mengubah relasi Supplier dari SPPG-level (sppg_id) menjadi Yayasan-level menggunakan tabel junction Supplier_Yayasan (supplier_id, yayasan_id)
2. THE Sistem SHALL mendukung satu supplier terdaftar di beberapa yayasan (relasi many-to-many melalui Supplier_Yayasan)
3. WHEN kepala_yayasan membuat atau mengedit supplier, THE Sistem SHALL menyimpan data supplier dan membuat relasi di Supplier_Yayasan untuk yayasan tersebut
4. THE Sistem SHALL menerapkan YayasanTenantScope pada query supplier sehingga kepala_yayasan hanya melihat supplier yang terhubung dengan yayasan miliknya
5. THE Sistem SHALL menyediakan migration script untuk memindahkan data supplier existing dari relasi SPPG ke relasi Yayasan
6. WHEN migration dijalankan, THE Sistem SHALL memetakan setiap supplier ke yayasan berdasarkan sppg.yayasan_id dari SPPG yang memiliki supplier tersebut
7. THE Sistem SHALL mengubah TenantMiddleware agar kepala_yayasan memiliki izin tulis (POST, PUT, DELETE) pada endpoint supplier

### Persyaratan 4: Supplier Catalog (Produk Supplier)

**User Story:** Sebagai supplier, saya ingin mengelola katalog produk yang saya jual beserta harga dan ketersediaan, sehingga yayasan dapat melihat dan memilih produk saya.

#### Acceptance Criteria

1. THE Supplier_Portal SHALL menyediakan CRUD untuk Supplier_Product (ingredient_id, unit_price, min_order_qty, is_available, stock_quantity)
2. WHEN supplier membuat atau mengubah Supplier_Product, THE Sistem SHALL memvalidasi bahwa ingredient_id merujuk ke Ingredient yang valid
3. THE Supplier_Product SHALL memiliki constraint unique pada kombinasi (supplier_id, ingredient_id) untuk mencegah duplikasi
4. WHEN supplier mengubah is_available menjadi false, THE Sistem SHALL menandai produk sebagai tidak tersedia tanpa menghapus data
5. THE Sistem SHALL menerapkan isolasi data sehingga supplier hanya dapat melihat dan mengelola Supplier_Product miliknya sendiri (berdasarkan supplier_id dari JWT claims)
6. THE Sistem SHALL menyediakan endpoint publik (untuk role yayasan) untuk melihat katalog supplier yang terhubung dengan yayasan tersebut

### Persyaratan 5: Pembuatan Purchase Order oleh Yayasan

**User Story:** Sebagai kepala_yayasan, saya ingin membuat Purchase Order ke supplier berdasarkan RAB yang sudah disetujui, sehingga pengadaan bahan baku dapat dilakukan secara terpusat.

#### Acceptance Criteria

1. WHEN RAB berstatus "approved_yayasan", THE PO_Manager SHALL mengizinkan kepala_yayasan membuat satu atau lebih PO dari RAB tersebut
2. THE PurchaseOrder SHALL menyimpan yayasan_id, rab_id, target_sppg_id, supplier_id, dan po_number
3. THE Sistem SHALL memvalidasi bahwa satu PO hanya ditujukan untuk satu SPPG (target_sppg_id) karena GRN terjadi di level SPPG
4. THE Sistem SHALL memvalidasi bahwa supplier_id dalam PO terhubung dengan yayasan pembuat PO melalui Supplier_Yayasan
5. WHEN PO dibuat, THE RAB_Item SHALL diupdate dengan po_id untuk setiap item yang tercakup dalam PO tersebut
6. THE Sistem SHALL mengizinkan satu RAB menghasilkan beberapa PO ke supplier yang berbeda
7. WHEN PO dibuat, THE Notification_Service SHALL mengirim push notification (FCM) dan in-app notification ke supplier terkait
8. WHEN PO dibuat, THE Notification_Service SHALL mengirim notifikasi in-app ke kepala_sppg dari SPPG target, sehingga SPPG mengetahui supplier mana yang akan mengirim barang dan untuk PO mana
9. THE PO_Manager SHALL menerapkan YayasanTenantScope sehingga kepala_yayasan hanya melihat PO dari yayasan miliknya

### Persyaratan 6: Supplier Portal — Autentikasi dan Dashboard

**User Story:** Sebagai supplier, saya ingin login ke sistem melalui web dan PWA mobile untuk melihat ringkasan pesanan dan pembayaran saya.

#### Acceptance Criteria

1. THE Sistem SHALL mendukung role baru "supplier" dalam JWT claims dengan field supplier_id (bukan sppg_id atau yayasan_id)
2. WHEN supplier login, THE Sistem SHALL menghasilkan JWT token dengan claims: user_id, role="supplier", supplier_id
3. THE Supplier_Portal SHALL menampilkan dashboard dengan ringkasan: total PO aktif, total PO selesai, total invoice pending, total pembayaran diterima
4. THE Sistem SHALL menerapkan entity-scoped isolation untuk supplier (filter berdasarkan supplier_id, bukan sppg_id atau yayasan_id)
5. THE Supplier_Portal SHALL dapat diakses melalui web browser dan PWA mobile dengan tampilan responsif

### Persyaratan 7: Supplier Portal — Notifikasi

**User Story:** Sebagai supplier, saya ingin menerima notifikasi ketika ada PO baru dari yayasan, sehingga saya dapat segera memproses pesanan.

#### Acceptance Criteria

1. WHEN PO baru dibuat untuk supplier, THE Notification_Service SHALL mengirim push notification melalui FCM ke device supplier
2. THE Notification_Service SHALL menyimpan notifikasi dalam database (in-app notification) dengan type, title, message, dan link
3. THE Supplier_Portal SHALL menampilkan daftar notifikasi dengan status baca/belum baca
4. WHEN supplier membuka notifikasi, THE Sistem SHALL menandai notifikasi sebagai sudah dibaca
5. THE Supplier_Portal SHALL menampilkan badge jumlah notifikasi yang belum dibaca

### Persyaratan 8: Penerimaan Barang (GRN) yang Dimodifikasi

**User Story:** Sebagai kepala_sppg atau asisten_lapangan, saya ingin menerima dan memeriksa barang berdasarkan PO, sehingga kualitas dan kuantitas barang tervalidasi.

#### Acceptance Criteria

1. WHEN barang diterima di SPPG, THE GRN_Processor SHALL memvalidasi GRN terhadap PO (bukan langsung terhadap RAB)
2. THE Sistem SHALL memvalidasi bahwa satu PO hanya menghasilkan satu GRN (relasi 1:1, tidak ada partial delivery)
3. IF PO sudah memiliki GRN, THEN THE GRN_Processor SHALL menolak pembuatan GRN baru untuk PO tersebut
4. WHEN GRN selesai dibuat, THE Sistem SHALL mengupdate RAB_Item dengan grn_id untuk setiap item yang diterima melalui PO terkait
5. WHEN GRN selesai dibuat, THE Sistem SHALL menambah stok inventory sesuai kuantitas yang diterima
6. THE GRN_Processor SHALL mengizinkan SPPG memberikan quality_rating (skala 1-5) untuk pengiriman tersebut
7. WHEN quality_rating diberikan, THE Sistem SHALL mengupdate rata-rata quality_rating pada Supplier terkait
8. THE Sistem SHALL menampilkan checklist pada RAB yang menunjukkan status setiap item: "PO ID A diterima via GRN ID B"

### Persyaratan 9: Invoice dan Pembayaran

**User Story:** Sebagai supplier, saya ingin membuat invoice setelah barang diterima, dan sebagai kepala_yayasan, saya ingin memproses pembayaran, sehingga alur keuangan tercatat dengan baik.

#### Acceptance Criteria

1. WHEN GRN selesai untuk suatu PO, THE Supplier_Portal SHALL mengizinkan supplier membuat invoice untuk PO tersebut
2. THE Invoice SHALL menyimpan invoice_number (format: INV-YYYYMMDD-XXXX), po_id, supplier_id, yayasan_id, amount, status ("pending"), dan due_date
3. THE Sistem SHALL memvalidasi bahwa amount pada invoice sesuai dengan total_amount pada PO terkait
4. WHEN invoice dibuat, THE Notification_Service SHALL mengirim notifikasi ke kepala_yayasan
5. WHEN kepala_yayasan memproses pembayaran, THE Invoice_Manager SHALL menyimpan data pembayaran: payment_date, amount, proof_url (bukti transfer), dan payment_method ("bank_transfer")
6. WHEN bukti pembayaran diupload, THE Sistem SHALL menyimpan file ke filesystem lokal (/uploads/payment-proofs/)
7. WHEN pembayaran berhasil dicatat, THE Invoice_Manager SHALL mengubah status invoice menjadi "paid" dan membuat CashFlowEntry dengan category "pengadaan" dan type "expense" untuk yayasan
8. WHEN pembayaran dicatat, THE Notification_Service SHALL mengirim notifikasi ke supplier bahwa pembayaran telah dilakukan
9. IF invoice sudah berstatus "paid", THEN THE Sistem SHALL menolak pembayaran duplikat

### Persyaratan 10: Perubahan Tenant Middleware dan Otorisasi

**User Story:** Sebagai administrator sistem, saya ingin tenant middleware mendukung pola baru untuk entitas milik yayasan dan role supplier, sehingga isolasi data tetap terjaga.

#### Acceptance Criteria

1. THE Sistem SHALL menambahkan YayasanTenantScope yang memfilter data berdasarkan yayasan_id untuk entitas: Supplier (via Supplier_Yayasan), PurchaseOrder, RAB, Invoice, Payment
2. THE TenantMiddleware SHALL mengizinkan kepala_yayasan melakukan operasi tulis (POST, PUT, DELETE) pada endpoint: /api/v1/rab/*/approve, /api/v1/rab/*/reject, /api/v1/suppliers, /api/v1/purchase-orders, /api/v1/invoices/*/pay
3. THE Sistem SHALL menambahkan SupplierTenantScope yang memfilter data berdasarkan supplier_id dari JWT claims untuk role "supplier"
4. THE Sistem SHALL memvalidasi bahwa role "supplier" hanya dapat mengakses endpoint: supplier catalog CRUD, PO list (read-only), invoice CRUD, notifikasi, dan dashboard supplier
5. THE Sistem SHALL memperbarui PermissionChecker dengan permission baru: "rab_management", "supplier_portal", "yayasan_procurement", "invoice_management"
6. WHEN user dengan role "supplier" mengakses endpoint di luar izinnya, THE Sistem SHALL mengembalikan HTTP 403 Forbidden

### Persyaratan 11: Migrasi Data Existing

**User Story:** Sebagai administrator sistem, saya ingin data supplier yang sudah ada dimigrasikan dari level SPPG ke level Yayasan, sehingga transisi ke sistem baru berjalan lancar.

#### Acceptance Criteria

1. THE Sistem SHALL menyediakan migration script yang membuat tabel baru: rabs, rab_items, supplier_products, supplier_yayasans, invoices, payments
2. THE Sistem SHALL menyediakan migration script yang menambahkan kolom baru pada purchase_orders: yayasan_id, rab_id, target_sppg_id
3. WHEN migration dijalankan, THE Sistem SHALL memindahkan setiap supplier ke Supplier_Yayasan berdasarkan mapping sppg.yayasan_id
4. THE Sistem SHALL mempertahankan kolom sppg_id pada Supplier sebagai nullable untuk backward compatibility selama masa transisi
5. IF supplier memiliki sppg_id yang merujuk ke SPPG tanpa yayasan_id, THEN THE Sistem SHALL mencatat warning log dan melewati supplier tersebut dalam migrasi

### Persyaratan 12: RAB Tracking dan Reporting

**User Story:** Sebagai kepala_yayasan, saya ingin melihat status lengkap RAB termasuk PO dan GRN yang terkait, sehingga saya dapat memantau progres pengadaan.

#### Acceptance Criteria

1. THE Sistem SHALL menyediakan endpoint untuk melihat detail RAB beserta status setiap RAB_Item (pending, po_created, grn_received)
2. THE RAB_Item SHALL menghitung status berdasarkan: "pending" jika po_id NULL, "po_created" jika po_id terisi tapi grn_id NULL, "grn_received" jika grn_id terisi
3. THE Sistem SHALL menyediakan ringkasan RAB: total items, items with PO, items received, total budget, total spent
4. WHEN semua RAB_Item berstatus "grn_received", THE Sistem SHALL mengubah status RAB menjadi "completed"
5. THE Sistem SHALL menyediakan endpoint perbandingan RAB vs Aktual: planned_amount vs actual_amount per item

### Persyaratan 13: Visibilitas RAB untuk SPPG (PO & Supplier Tracking)

**User Story:** Sebagai kepala_sppg atau asisten_lapangan, saya ingin melihat di RAB supplier mana yang akan datang mengirim barang dan untuk PO mana, sehingga saya siap menerima dan melakukan GRN.

#### Acceptance Criteria

1. THE Sistem SHALL menampilkan pada halaman detail RAB di sisi SPPG: daftar RAB_Item beserta po_id, nama supplier, dan status GRN
2. WHEN yayasan membuat PO dari RAB, THE Sistem SHALL mengupdate tampilan RAB di SPPG dengan informasi PO (po_number, supplier_name, expected_delivery)
3. THE Sistem SHALL menyediakan endpoint GET /api/v1/rab/:id/po-tracking yang mengembalikan daftar PO terkait RAB beserta status GRN masing-masing
4. WHEN SPPG membuka halaman RAB, THE Sistem SHALL menampilkan ringkasan: total PO yang sudah dibuat, PO yang sudah diterima (GRN), dan PO yang masih pending delivery
