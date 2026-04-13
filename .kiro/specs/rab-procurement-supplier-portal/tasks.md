# Implementation Plan: RAB, Procurement Yayasan & Supplier Portal

## Overview

Implementasi fitur RAB (Rencana Anggaran Belanja), procurement terpusat di level yayasan, dan portal supplier. Urutan task mengikuti dependency: fondasi data model → tenant/auth → service layer → handler/router → integrasi. Bahasa implementasi: Go (Gin + GORM + PostgreSQL).

## Tasks

- [x] 1. Definisi model baru dan modifikasi model existing
  - [x] 1.1 Tambahkan model RAB, RABItem, SupplierProduct, SupplierYayasan, Invoice, Payment di `backend/internal/models/supply_chain.go`
    - Definisikan struct RAB dengan field: ID, RABNumber, MenuPlanID, SPPGID, YayasanID, Status, TotalAmount, RevisionNotes, ApprovedBySPPG, ApprovedAtSPPG, ApprovedByYayasan, ApprovedAtYayasan, CreatedBy, CreatedAt, UpdatedAt, dan relasi (MenuPlan, Creator, Items)
    - Definisikan struct RABItem dengan field: ID, RABID, IngredientID, Quantity, Unit, UnitPrice, Subtotal, RecommendedSupplierID, POID, GRNID, Status, CreatedAt, UpdatedAt, dan relasi
    - Definisikan struct SupplierProduct dengan unique constraint pada (SupplierID, IngredientID)
    - Definisikan struct SupplierYayasan sebagai junction table dengan unique constraint pada (SupplierID, YayasanID)
    - Definisikan struct Invoice dan Payment sesuai design document
    - _Requirements: 1.4, 1.5, 3.1, 3.2, 4.1, 4.3, 9.2, 9.5_

  - [x] 1.2 Modifikasi model PurchaseOrder — tambahkan field YayasanID, RABID, TargetSPPGID di `backend/internal/models/supply_chain.go`
    - Tambahkan field nullable: YayasanID *uint, RABID *uint, TargetSPPGID *uint dengan gorm index tag
    - _Requirements: 5.2, 5.3, 11.2_

  - [x] 1.3 Modifikasi model User — tambahkan SupplierID dan perluas Role enum di `backend/internal/models/user.go`
    - Tambahkan field SupplierID *uint dengan gorm index tag
    - Perluas validate tag Role: tambahkan "supplier" ke oneof
    - _Requirements: 6.1, 6.2_

  - [x] 1.4 Modifikasi model Supplier — tambahkan relasi many2many Yayasans di `backend/internal/models/supply_chain.go`
    - Tambahkan field: `Yayasans []Yayasan gorm:"many2many:supplier_yayasans"`
    - _Requirements: 3.1, 3.2_

  - [x] 1.5 Modifikasi model CashFlowEntry — tambahkan YayasanID dan perluas Category di `backend/internal/models/financial.go`
    - Tambahkan field YayasanID *uint dengan gorm index tag
    - Perluas validate tag Category: tambahkan "pengadaan" ke oneof
    - _Requirements: 9.7_

  - [x] 1.6 Registrasi model baru di AllModels() di `backend/internal/models/models.go`
    - Tambahkan &RAB{}, &RABItem{}, &SupplierProduct{}, &SupplierYayasan{}, &Invoice{}, &Payment{} ke slice AllModels
    - _Requirements: 11.1_

  - [ ]* 1.7 Tulis unit test untuk validasi struct tag dan relasi model baru
    - Test bahwa unique constraint SupplierProduct (supplier_id, ingredient_id) terdefinisi
    - Test bahwa unique constraint SupplierYayasan (supplier_id, yayasan_id) terdefinisi
    - _Requirements: 4.3, 3.2_

- [x] 2. Database migration dan migrasi data supplier
  - [x] 2.1 Tambahkan fungsi migration untuk tabel baru dan kolom baru di `backend/internal/database/migrate.go`
    - Buat fungsi MigrateRABProcurement(db *gorm.DB) yang menambahkan kolom baru pada purchase_orders (yayasan_id, rab_id, target_sppg_id), users (supplier_id), cash_flow_entries (yayasan_id)
    - Buat index pada tabel baru: idx_rabs_status, idx_rabs_sppg_id, idx_rabs_yayasan_id, idx_rab_items_status, idx_supplier_products_available, idx_invoices_status
    - Panggil MigrateRABProcurement dari fungsi Migrate utama
    - _Requirements: 11.1, 11.2_

  - [x] 2.2 Implementasi migrasi data supplier existing dari SPPG-level ke Yayasan-level di `backend/internal/database/migrate.go`
    - Buat fungsi MigrateSupplierToYayasan(db *gorm.DB) yang: query semua supplier dengan sppg_id, join ke sppgs untuk mendapatkan yayasan_id, buat record SupplierYayasan
    - Jika supplier memiliki sppg_id yang merujuk ke SPPG tanpa yayasan_id, log warning dan skip
    - Jalankan dalam transaction untuk atomicity
    - _Requirements: 3.5, 3.6, 11.3, 11.4, 11.5_

  - [ ]* 2.3 Tulis unit test untuk migrasi data supplier
    - Test mapping supplier ke yayasan berdasarkan sppg.yayasan_id
    - Test skip supplier dengan SPPG tanpa yayasan_id
    - _Requirements: 11.3, 11.5_

- [x] 3. Checkpoint — Pastikan migration berjalan dan semua test lulus
  - Ensure all tests pass, ask the user if questions arise.

- [x] 4. Perubahan Tenant Middleware dan Otorisasi
  - [x] 4.1 Modifikasi JWTClaims — tambahkan SupplierID di `backend/internal/services/auth_service.go`
    - Tambahkan field SupplierID *uint pada struct JWTClaims
    - Modifikasi GenerateToken agar menerima dan menyimpan supplierID
    - Modifikasi Login agar mengambil SupplierID dari User dan memasukkannya ke token
    - _Requirements: 6.1, 6.2_

  - [x] 4.2 Modifikasi JWTAuth middleware — set supplier_id di context di `backend/internal/middleware/auth.go`
    - Setelah ValidateToken, jika claims.SupplierID != nil, set c.Set("supplier_id", *claims.SupplierID)
    - _Requirements: 6.4, 10.3_

  - [x] 4.3 Modifikasi TenantMiddleware — tambahkan case "supplier" dan izin tulis kepala_yayasan di `backend/internal/middleware/tenant.go`
    - Tambahkan case "supplier" di switch: extract supplier_id dari context, set tenant_db
    - Modifikasi case "kepala_yayasan": izinkan write pada whitelist endpoints (/api/v1/suppliers, /api/v1/purchase-orders, /api/v1/rab, /api/v1/invoices)
    - Buat fungsi isYayasanWriteAllowed(path string) untuk cek whitelist
    - _Requirements: 10.1, 10.2, 10.3, 3.7_

  - [x] 4.4 Implementasi YayasanTenantScope dan SupplierTenantScope di `backend/internal/middleware/tenant.go`
    - Buat fungsi YayasanTenantScope(c *gin.Context) yang memfilter berdasarkan yayasan_id
    - Buat fungsi SupplierTenantScope(c *gin.Context) yang memfilter berdasarkan supplier_id dari JWT claims
    - _Requirements: 3.4, 4.5, 6.4, 10.1, 10.3_

  - [x] 4.5 Update PermissionChecker — tambahkan permission baru di `backend/internal/middleware/auth.go`
    - Tambahkan permission: "rab_management" (kepala_sppg, kepala_yayasan, ahli_gizi), "supplier_portal" (supplier), "yayasan_procurement" (kepala_yayasan), "invoice_management" (kepala_yayasan, supplier)
    - _Requirements: 10.5_

  - [ ]* 4.6 Tulis unit test untuk TenantMiddleware perubahan
    - Test bahwa supplier role mendapat SupplierTenantScope
    - Test bahwa kepala_yayasan bisa write pada whitelist endpoints
    - Test bahwa kepala_yayasan tetap read-only pada endpoint lain
    - Test bahwa supplier mendapat 403 pada endpoint di luar izinnya
    - _Requirements: 10.2, 10.4, 10.6_

  - [ ]* 4.7 Tulis property test untuk SupplierTenantScope dan YayasanTenantScope
    - **Property 9: YayasanTenantScope memfilter data dengan benar**
    - **Property 13: Isolasi data supplier (SupplierTenantScope)**
    - **Validates: Requirements 3.4, 4.5, 6.4, 10.1, 10.3**

- [x] 5. Checkpoint — Pastikan middleware dan auth berfungsi, semua test lulus
  - Ensure all tests pass, ask the user if questions arise.

- [x] 6. Implementasi SupplierProductService (Katalog Produk Supplier)
  - [x] 6.1 Buat `backend/internal/services/supplier_product_service.go`
    - Implementasi struct SupplierProductService dengan dependency db *gorm.DB
    - Implementasi CreateProduct: validasi ingredient_id valid, cek unique constraint (supplier_id, ingredient_id), create record
    - Implementasi UpdateProduct: validasi ownership (supplier_id), update fields
    - Implementasi GetProductsBySupplier: filter by supplier_id
    - Implementasi GetCatalogByYayasan: join SupplierYayasan untuk filter supplier yang terhubung dengan yayasan, filter is_available=true
    - Implementasi ToggleAvailability: update is_available tanpa hapus data
    - Implementasi DeleteProduct: soft delete atau hard delete
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

  - [ ]* 6.2 Tulis property test untuk SupplierProduct
    - **Property 11: Validasi ingredient pada SupplierProduct**
    - **Property 12: Unique constraint pada SupplierProduct (supplier_id, ingredient_id)**
    - **Validates: Requirements 4.2, 4.3**

  - [ ]* 6.3 Tulis unit test untuk SupplierProductService
    - Test CRUD operations
    - Test GetCatalogByYayasan hanya mengembalikan produk dari supplier yang terhubung
    - Test ToggleAvailability tidak menghapus data
    - _Requirements: 4.1, 4.4, 4.6_

- [x] 7. Implementasi RABGeneratorService (Auto-Generasi RAB)
  - [x] 7.1 Buat `backend/internal/services/rab_generator_service.go`
    - Implementasi struct RABGeneratorService dengan dependency db *gorm.DB, notif *NotificationService
    - Implementasi AggregateIngredients: traverse MenuPlan → MenuItems → Recipe → RecipeItems → SemiFinishedGoods → SemiFinishedRecipeIngredients, hitung total kuantitas per ingredient (dikalikan jumlah porsi)
    - Implementasi RecommendSupplier: filter SupplierProduct by ingredient_id + is_available=true + supplier terhubung yayasan, sort by unit_price ASC lalu quality_rating DESC, return supplier_id dan unit_price
    - Implementasi GenerateRABFromMenuPlan: panggil AggregateIngredients, untuk setiap ingredient panggil RecommendSupplier, buat RAB + RABItems dalam transaction, generate rab_number (format RAB-YYYYMMDD-XXXX), kirim notifikasi ke kepala_sppg
    - Handle edge case: ingredient tanpa supplier → set recommended_supplier_id NULL, unit_price 0
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7_

  - [ ]* 7.2 Tulis property test untuk agregasi ingredient
    - **Property 1: Agregasi ingredient dari menu plan menghasilkan total yang benar**
    - **Validates: Requirements 1.2**

  - [ ]* 7.3 Tulis property test untuk rekomendasi supplier
    - **Property 2: Rekomendasi supplier memilih harga termurah dengan quality terbaik**
    - **Validates: Requirements 1.3**

  - [ ]* 7.4 Tulis property test untuk integritas data RAB
    - **Property 3: Integritas data RAB dan RAB_Item**
    - **Validates: Requirements 1.4, 1.5**

- [x] 8. Implementasi ApprovalEngineService (Alur Persetujuan RAB)
  - [x] 8.1 Buat `backend/internal/services/approval_engine_service.go`
    - Implementasi struct ApprovalEngineService dengan dependency db *gorm.DB, notif *NotificationService
    - Implementasi ApproveByKepalaSPPG: validasi status="draft", ubah ke "approved_sppg", set approved_by_sppg dan approved_at_sppg, catat audit trail, kirim notifikasi ke kepala_yayasan
    - Implementasi ApproveByKepalaYayasan: validasi status="approved_sppg", ubah ke "approved_yayasan", set approved_by_yayasan dan approved_at_yayasan, catat audit trail, kirim notifikasi ke kepala_sppg
    - Implementasi RejectByKepalaYayasan: validasi status="approved_sppg", ubah ke "revision_requested", simpan revision_notes, catat audit trail, kirim notifikasi ke kepala_sppg
    - Implementasi ResubmitRAB: validasi status="revision_requested", ubah ke "draft", catat audit trail
    - Implementasi createAuditTrail: buat record AuditTrail dengan user_id, timestamp, from_status, to_status, notes
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8_

  - [ ]* 8.2 Tulis property test untuk state machine RAB
    - **Property 4: State machine RAB mengikuti transisi yang valid**
    - **Validates: Requirements 2.1, 2.2, 2.3**

  - [ ]* 8.3 Tulis property test untuk audit trail
    - **Property 5: Audit trail tercatat pada setiap perubahan status RAB**
    - **Validates: Requirements 2.5**

  - [ ]* 8.4 Tulis property test untuk izin edit RAB
    - **Property 6: Izin edit RAB berdasarkan status**
    - **Validates: Requirements 2.6, 2.7**

- [x] 9. Implementasi RABService (CRUD dan Tracking RAB)
  - [x] 9.1 Buat `backend/internal/services/rab_service.go`
    - Implementasi struct RABService dengan dependency db *gorm.DB
    - Implementasi GetRABByID: preload Items, MenuPlan, Creator, RecommendedSupplier
    - Implementasi GetRABList: filter by sppg_id atau yayasan_id (tergantung role), support pagination
    - Implementasi UpdateRAB: validasi status "draft" atau "revision_requested", update items dan total_amount
    - Implementasi GetRABComparison: hitung planned_amount vs actual_amount per item (dari GRN)
    - Implementasi GetRABSummary: hitung total_items, items_with_po, items_received, total_budget, total_spent
    - Implementasi GetPOTracking: list PO terkait RAB beserta status GRN masing-masing
    - Implementasi CheckAndCompleteRAB: cek apakah semua RABItem grn_received → auto-complete RAB
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 13.1, 13.2, 13.3, 13.4_

  - [ ]* 9.2 Tulis property test untuk derivasi status RABItem
    - **Property 20: Derivasi status RAB_Item berdasarkan po_id dan grn_id**
    - **Validates: Requirements 12.2, 8.8**

  - [ ]* 9.3 Tulis property test untuk ringkasan RAB
    - **Property 25: Ringkasan RAB menghitung dengan benar**
    - **Validates: Requirements 12.3**

  - [ ]* 9.4 Tulis property test untuk auto-complete RAB
    - **Property 26: RAB auto-complete saat semua item diterima**
    - **Validates: Requirements 12.4**

- [x] 10. Checkpoint — Pastikan semua service RAB dan approval berfungsi, semua test lulus
  - Ensure all tests pass, ask the user if questions arise.

- [x] 11. Modifikasi PurchaseOrderService dan GoodsReceiptService
  - [x] 11.1 Modifikasi PurchaseOrderService — tambahkan validasi RAB dan yayasan di `backend/internal/services/purchase_order_service.go`
    - Modifikasi CreatePurchaseOrder: tambahkan validasi RAB berstatus "approved_yayasan", validasi supplier terhubung yayasan via SupplierYayasan, set yayasan_id/rab_id/target_sppg_id, update RABItem.po_id setelah PO dibuat, kirim notifikasi ke supplier DAN kepala_sppg SPPG target
    - Tambahkan fungsi GetPurchaseOrdersByYayasan untuk YayasanTenantScope
    - Tambahkan fungsi GetPurchaseOrdersBySupplier yang sudah ada tapi pastikan support supplier role (read-only)
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 5.8, 5.9_

  - [x] 11.2 Modifikasi GoodsReceiptService — tambahkan validasi 1:1 PO-GRN dan update RAB di `backend/internal/services/goods_receipt_service.go`
    - Modifikasi CreateGoodsReceipt: tambahkan validasi 1:1 PO-GRN (tolak jika PO sudah punya GRN), update RABItem.grn_id setelah GRN dibuat, panggil CheckAndCompleteRAB, update supplier quality_rating average
    - Tambahkan field quality_rating pada request GRN (sudah ada di model)
    - Implementasi updateSupplierQualityRating: hitung rata-rata quality_rating dari semua GRN supplier
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7, 8.8_

  - [ ]* 11.3 Tulis property test untuk validasi PO
    - **Property 14: PO hanya bisa dibuat dari RAB yang approved_yayasan**
    - **Property 15: Validasi supplier-yayasan pada pembuatan PO**
    - **Property 16: RAB_Item diupdate dengan po_id saat PO dibuat**
    - **Validates: Requirements 5.1, 5.4, 5.5**

  - [ ]* 11.4 Tulis property test untuk GRN
    - **Property 17: Relasi 1:1 antara PO dan GRN**
    - **Property 18: RAB_Item diupdate dengan grn_id saat GRN dibuat**
    - **Property 19: Quality rating dalam range valid dan rata-rata supplier terupdate**
    - **Validates: Requirements 8.2, 8.3, 8.4, 8.6, 8.7**

- [x] 12. Implementasi InvoiceService (Invoice dan Pembayaran)
  - [x] 12.1 Buat `backend/internal/services/invoice_service.go`
    - Implementasi struct InvoiceService dengan dependency db *gorm.DB, cashFlow *CashFlowService, notif *NotificationService
    - Implementasi CreateInvoice: validasi PO sudah punya GRN (status "received"), validasi amount = PO total_amount, generate invoice_number (format INV-YYYYMMDD-XXXX), kirim notifikasi ke kepala_yayasan
    - Implementasi ProcessPayment: validasi invoice status="pending", simpan Payment record, upload bukti ke /uploads/payment-proofs/, update invoice status → "paid", buat CashFlowEntry (category="pengadaan", type="expense", yayasan_id), kirim notifikasi ke supplier
    - Implementasi GetInvoicesBySupplier: filter by supplier_id
    - Implementasi GetInvoicesByYayasan: filter by yayasan_id
    - Implementasi GetSupplierPaymentHistory: list pembayaran yang diterima supplier
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 9.7, 9.8, 9.9_

  - [ ]* 12.2 Tulis property test untuk invoice
    - **Property 21: Invoice hanya bisa dibuat setelah GRN selesai**
    - **Property 22: Amount invoice harus sesuai dengan total PO**
    - **Validates: Requirements 9.1, 9.3**

  - [ ]* 12.3 Tulis property test untuk pembayaran
    - **Property 23: Pembayaran mengubah status invoice dan membuat CashFlowEntry**
    - **Property 24: Pembayaran duplikat ditolak**
    - **Validates: Requirements 9.7, 9.9**

- [x] 13. Checkpoint — Pastikan semua service berfungsi, semua test lulus
  - Ensure all tests pass, ask the user if questions arise.

- [x] 14. Implementasi Handler Layer — RABHandler
  - [x] 14.1 Buat `backend/internal/handlers/rab_handler.go`
    - Implementasi struct RABHandler dengan dependency rabService, rabGenerator, approvalEngine
    - Implementasi GetRABList: GET /api/v1/rab — list RAB dengan tenant scope (SPPG atau Yayasan)
    - Implementasi GetRABDetail: GET /api/v1/rab/:id — detail RAB + items + status
    - Implementasi UpdateRAB: PUT /api/v1/rab/:id — edit RAB (hanya draft/revision_requested)
    - Implementasi ApproveSPPG: POST /api/v1/rab/:id/approve-sppg — kepala_sppg approve
    - Implementasi ApproveYayasan: POST /api/v1/rab/:id/approve-yayasan — kepala_yayasan approve
    - Implementasi RejectRAB: POST /api/v1/rab/:id/reject — kepala_yayasan reject dengan notes
    - Implementasi ResubmitRAB: POST /api/v1/rab/:id/resubmit — SPPG resubmit
    - Implementasi GetRABComparison: GET /api/v1/rab/:id/comparison — RAB vs Aktual
    - Implementasi GetPOTracking: GET /api/v1/rab/:id/po-tracking — PO tracking untuk SPPG
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.6, 12.1, 12.5, 13.3_

  - [ ]* 14.2 Tulis unit test untuk RABHandler
    - Test endpoint authorization (kepala_sppg vs kepala_yayasan)
    - Test error responses untuk invalid status transitions
    - _Requirements: 2.1, 2.2, 2.3_

- [x] 15. Implementasi Handler Layer — SupplierPortalHandler
  - [x] 15.1 Buat `backend/internal/handlers/supplier_portal_handler.go`
    - Implementasi struct SupplierPortalHandler dengan dependency supplierProductService, invoiceService, purchaseOrderService
    - Implementasi GetSupplierDashboard: GET /api/v1/supplier/dashboard — ringkasan: total PO aktif, PO selesai, invoice pending, pembayaran diterima
    - Implementasi GetSupplierPayments: GET /api/v1/supplier/payments — riwayat pembayaran
    - Implementasi CRUD SupplierProduct: GET/POST/PUT/DELETE /api/v1/supplier-products — dengan SupplierTenantScope
    - Implementasi GetCatalog: GET /api/v1/supplier-products (untuk role yayasan) — katalog supplier terhubung yayasan
    - _Requirements: 4.1, 4.5, 4.6, 6.3, 6.4, 6.5_

  - [ ]* 15.2 Tulis unit test untuk SupplierPortalHandler
    - Test isolasi data supplier (hanya lihat data miliknya)
    - Test dashboard calculation
    - _Requirements: 4.5, 6.3, 6.4_

- [x] 16. Implementasi Handler Layer — InvoiceHandler
  - [x] 16.1 Buat `backend/internal/handlers/invoice_handler.go`
    - Implementasi struct InvoiceHandler dengan dependency invoiceService
    - Implementasi GetInvoices: GET /api/v1/invoices — list invoice (scoped by role)
    - Implementasi CreateInvoice: POST /api/v1/invoices — supplier buat invoice
    - Implementasi GetInvoiceDetail: GET /api/v1/invoices/:id — detail invoice
    - Implementasi PayInvoice: POST /api/v1/invoices/:id/pay — kepala_yayasan bayar
    - Implementasi UploadPaymentProof: POST /api/v1/invoices/:id/upload-proof — upload bukti transfer ke /uploads/payment-proofs/
    - _Requirements: 9.1, 9.2, 9.4, 9.5, 9.6_

  - [ ]* 16.2 Tulis unit test untuk InvoiceHandler
    - Test validasi: invoice hanya bisa dibuat setelah GRN
    - Test validasi: pembayaran duplikat ditolak
    - _Requirements: 9.1, 9.9_

- [x] 17. Integrasi RAB Generator dengan Menu Plan Approval
  - [x] 17.1 Modifikasi MenuPlanningHandler.ApproveMenuPlan di `backend/internal/handlers/menu_planning_handler.go`
    - Setelah menu plan di-approve, panggil RABGeneratorService.GenerateRABFromMenuPlan
    - Inject RABGeneratorService sebagai dependency baru pada MenuPlanningHandler
    - Handle error gracefully: jika RAB generation gagal, log error tapi jangan gagalkan approval menu plan
    - _Requirements: 1.1, 1.7_

  - [ ]* 17.2 Tulis unit test untuk integrasi menu plan → RAB
    - Test bahwa approve menu plan memicu pembuatan RAB
    - Test bahwa RAB memiliki items yang benar dari menu plan
    - _Requirements: 1.1, 1.2_

- [x] 18. Modifikasi SupplyChainHandler — Supplier CRUD untuk Yayasan
  - [x] 18.1 Modifikasi `backend/internal/handlers/supply_chain_handler.go`
    - Modifikasi CreateSupplier: jika role=kepala_yayasan, buat SupplierYayasan record yang menghubungkan supplier dengan yayasan miliknya
    - Modifikasi GetAllSuppliers: jika role=kepala_yayasan, gunakan YayasanTenantScope (filter via SupplierYayasan)
    - Modifikasi UpdateSupplier: izinkan kepala_yayasan mengedit supplier yang terhubung yayasannya
    - _Requirements: 3.3, 3.4, 3.7_

  - [ ]* 18.2 Tulis property test untuk supplier-yayasan
    - **Property 7: Supplier many-to-many dengan yayasan**
    - **Property 8: Pembuatan supplier oleh kepala_yayasan membuat junction record**
    - **Validates: Requirements 3.2, 3.3**

  - [ ]* 18.3 Tulis property test untuk migrasi supplier
    - **Property 10: Migrasi supplier memetakan ke yayasan yang benar**
    - **Validates: Requirements 3.6, 11.3**

- [x] 19. Registrasi Route Baru di Router
  - [x] 19.1 Tambahkan semua route baru di `backend/internal/router/router.go`
    - Inisialisasi service baru: RABGeneratorService, ApprovalEngineService, RABService, SupplierProductService, InvoiceService
    - Inisialisasi handler baru: RABHandler, SupplierPortalHandler, InvoiceHandler
    - Registrasi route RAB: GET /rab, GET /rab/:id, PUT /rab/:id, POST /rab/:id/approve-sppg, POST /rab/:id/approve-yayasan, POST /rab/:id/reject, POST /rab/:id/resubmit, GET /rab/:id/comparison, GET /rab/:id/po-tracking
    - Registrasi route Supplier Products: GET/POST/PUT/DELETE /supplier-products
    - Registrasi route Invoice: GET/POST /invoices, GET /invoices/:id, POST /invoices/:id/pay, POST /invoices/:id/upload-proof
    - Registrasi route Supplier Dashboard: GET /supplier/dashboard, GET /supplier/payments
    - Terapkan RequireRole/RequirePermission middleware pada setiap route group
    - Terapkan SupplierTenantScope pada route supplier portal
    - _Requirements: 10.2, 10.4, 10.5_

- [x] 20. Implementasi Notifikasi untuk Alur Baru
  - [x] 20.1 Tambahkan notification types dan helper functions di `backend/internal/services/notification_service.go`
    - Tambahkan constants: NotificationTypeRABCreated, NotificationTypeRABApproved, NotificationTypeRABRejected, NotificationTypeNewPO, NotificationTypeInvoiceCreated, NotificationTypePaymentReceived
    - Implementasi SendRABNotification: kirim notifikasi saat RAB dibuat, disetujui, atau ditolak
    - Implementasi SendNewPONotification: kirim notifikasi ke supplier saat PO baru dibuat
    - Implementasi SendPOToSPPGNotification: kirim notifikasi ke kepala_sppg saat PO dibuat untuk SPPG-nya (agar tahu supplier mana yang akan datang)
    - Implementasi SendInvoiceNotification: kirim notifikasi ke kepala_yayasan saat invoice dibuat
    - Implementasi SendPaymentNotification: kirim notifikasi ke supplier saat pembayaran dilakukan
    - Gunakan pola existing: push ke Firebase + simpan in-app notification, graceful degradation jika FCM gagal
    - _Requirements: 1.7, 2.8, 5.7, 5.8, 7.1, 7.2, 9.4, 9.8_

  - [ ]* 20.2 Tulis unit test untuk notification helpers
    - Test bahwa notifikasi dibuat dengan type, title, message yang benar
    - Test bahwa notifikasi dikirim ke user yang tepat
    - _Requirements: 7.1, 7.2, 7.3_

- [x] 21. Endpoint Restriction untuk Role Supplier
  - [x] 21.1 Tambahkan middleware restriction untuk role supplier di `backend/internal/router/router.go`
    - Pastikan role "supplier" hanya bisa mengakses: supplier-products (CRUD), purchase-orders (GET only), invoices (CRUD), notifications, supplier/dashboard, supplier/payments
    - Endpoint lain harus mengembalikan HTTP 403 Forbidden untuk role supplier
    - _Requirements: 10.4, 10.6_

  - [ ]* 21.2 Tulis property test untuk supplier endpoint restriction
    - **Property 27: Supplier endpoint restriction untuk role supplier**
    - **Validates: Requirements 10.4, 10.6**

- [x] 22. Final checkpoint — Pastikan semua test lulus dan integrasi end-to-end berfungsi
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Task yang ditandai `*` bersifat opsional dan dapat dilewati untuk MVP lebih cepat
- Setiap task mereferensikan requirements spesifik untuk traceability
- Checkpoint memastikan validasi inkremental di setiap fase
- Property tests memvalidasi correctness properties universal dari design document
- Unit tests memvalidasi contoh spesifik dan edge cases
- Urutan implementasi: Model → Migration → Middleware/Auth → Services → Handlers → Router → Notifikasi → Integrasi
