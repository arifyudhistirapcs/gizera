# Implementation Plan: Multi-Tenancy & Hierarki Organisasi

## Overview

Transformasi Sistem ERP SPPG dari single-tenant menjadi multi-tenant dengan hierarki organisasi bertingkat (Superadmin → BGN → Yayasan → SPPG → Karyawan). Implementasi mencakup entitas organisasi baru, peran baru, tenant isolation via GORM scopes, dashboard agregasi, Firebase tenant-aware, migrasi data, dan dukungan PWA mobile.

## Tasks

- [x] 1. Data models dan migrasi database
  - [x] 1.1 Buat model Yayasan dan SPPG di `backend/internal/models/organization.go`
    - Definisikan struct `Yayasan` dan `SPPG` sesuai desain (kode unik, relasi FK, validasi)
    - Tambahkan ke `AllModels()` di `backend/internal/models/models.go`
    - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 2.4_

  - [x] 1.2 Perbarui model User di `backend/internal/models/user.go`
    - Tambahkan kolom `SPPGID *uint`, `YayasanID *uint`, `CreatedBy *uint` pada struct User
    - Perbarui validasi Role untuk menyertakan `superadmin` dan `admin_bgn`
    - Tambahkan relasi GORM `SPPG` dan `Yayasan` pada User
    - _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7_

  - [x] 1.3 Tambahkan kolom `sppg_id` pada semua model operasional
    - Tambahkan `SPPGID *uint gorm:"index"` pada 26 model operasional (Recipe, Ingredient, MenuPlan, Supplier, PurchaseOrder, GoodsReceipt, InventoryItem, School, DeliveryTask, Employee, Attendance, dll sesuai desain)
    - Pastikan tabel `audit_trails` dan `system_configs` TIDAK ditambahkan `sppg_id`
    - _Requirements: 6.1, 6.6_

  - [x] 1.4 Buat model dashboard agregasi di `backend/internal/models/dashboard.go`
    - Definisikan struct: `KepalaYayasanAggregatedDashboard`, `SPPGSummary`, `AggregatedProduction`, `AggregatedDelivery`, `AggregatedFinancial`, `AggregatedReview`, `AdminBGNDashboard`, `YayasanSummary`
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 10.1, 10.2, 10.3, 10.4, 10.5_

  - [x] 1.5 Buat fungsi migrasi multi-tenant di `backend/internal/database/migrate.go`
    - Tambahkan fungsi `MigrateMultiTenant(db *gorm.DB) error`
    - Buat Yayasan default ("Yayasan Default", kode "YYS-0001") dan SPPG default ("SPPG Default", kode "SPPG-0001")
    - Isi `sppg_id` pada semua record operasional existing dengan ID SPPG default
    - Isi `sppg_id` dan `yayasan_id` pada record User existing yang memiliki peran SPPG
    - Buat akun superadmin default
    - Tambahkan indeks pada kolom `sppg_id`
    - Validasi tidak ada record dengan `sppg_id = NULL` pada tabel operasional
    - Implementasikan rollback jika migrasi gagal
    - _Requirements: 12.1, 12.2, 12.3, 12.4, 12.5, 12.6, 12.7_

  - [ ]* 1.6 Write property test untuk migrasi data completeness
    - **Property 11: Migrasi Data Completeness**
    - Verifikasi semua record operasional memiliki `sppg_id` valid setelah migrasi
    - **Validates: Requirements 12.2, 12.5**

  - [ ]* 1.7 Write property test untuk kode unik Yayasan dan SPPG
    - **Property 10: Kode Unik Yayasan dan SPPG**
    - Verifikasi tidak ada duplikasi kode Yayasan atau SPPG di seluruh sistem
    - **Validates: Requirements 1.3, 2.4**

  - [ ]* 1.8 Write property test untuk Yayasan-SPPG relationship integrity
    - **Property 7: Yayasan-SPPG Relationship Integrity**
    - Verifikasi setiap SPPG terhubung ke tepat satu Yayasan yang valid
    - **Validates: Requirements 2.3**

- [x] 2. Checkpoint — Pastikan semua model dan migrasi berjalan
  - Ensure all tests pass, ask the user if questions arise.

- [x] 3. Auth service dan tenant middleware
  - [x] 3.1 Perbarui JWT claims di `backend/internal/services/auth_service.go`
    - Tambahkan `SPPGID *uint` dan `YayasanID *uint` pada struct `JWTClaims`
    - Perbarui `Login()` untuk query User dengan preload SPPG/Yayasan dan menyertakan tenant info di token
    - Perbarui `GenerateToken()` untuk menerima dan menyertakan `sppg_id` dan `yayasan_id`
    - Perbarui `ValidateToken()` untuk mengekstrak tenant claims
    - _Requirements: 8.8, 7.1_

  - [x] 3.2 Perbarui auth middleware di `backend/internal/middleware/auth.go`
    - Perbarui `JWTAuth()` untuk set `sppg_id` dan `yayasan_id` di Gin context dari JWT claims
    - Perbarui `RequireRole()` untuk menyertakan `superadmin` dan `admin_bgn`
    - Perbarui `PermissionChecker` untuk menambahkan permission baru (organization management, dashboard BGN, dashboard Yayasan, user provisioning)
    - _Requirements: 15.1, 15.2, 5.1, 4.1, 17.1, 17.2, 17.3, 17.4, 17.5_

  - [x] 3.3 Buat tenant middleware di `backend/internal/middleware/tenant.go`
    - Implementasikan `TenantMiddleware(db *gorm.DB) gin.HandlerFunc` yang mengekstrak tenant context dari Gin context
    - Implementasikan `TenantScope(c *gin.Context) func(db *gorm.DB) *gorm.DB` sebagai GORM scope
    - Untuk peran SPPG: filter `WHERE sppg_id = ?`
    - Untuk kepala_yayasan: filter `WHERE sppg_id IN (SELECT id FROM sppgs WHERE yayasan_id = ?)`
    - Untuk admin_bgn/superadmin: bypass filter (dengan optional query param filtering)
    - Implementasikan auto-inject `sppg_id` pada INSERT untuk peran SPPG
    - Implementasikan fail-closed: tolak request jika tenant context gagal diekstrak
    - Implementasikan read-only enforcement untuk kepala_yayasan dan admin_bgn pada data operasional
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 4.4, 5.4, 6.2, 6.3, 6.4, 6.5, 14.1, 14.5, 14.6_

  - [ ]* 3.4 Write property test untuk isolasi data tenant
    - **Property 1: Isolasi Data Tenant**
    - Verifikasi pengguna SPPG hanya melihat data dengan `sppg_id` miliknya
    - **Validates: Requirements 6.2, 6.3, 6.5**

  - [ ]* 3.5 Write property test untuk auto-inject sppg_id
    - **Property 2: Auto-Inject SPPG_ID pada Insert**
    - Verifikasi record baru dari pengguna SPPG otomatis terisi `sppg_id`
    - **Validates: Requirements 6.4, 7.3**

  - [ ]* 3.6 Write property test untuk scope Kepala Yayasan
    - **Property 3: Scope Kepala Yayasan**
    - Verifikasi Kepala Yayasan hanya melihat data dari SPPG di bawah Yayasan-nya
    - **Validates: Requirements 4.2, 7.4**

  - [ ]* 3.7 Write property test untuk read-only enforcement Kepala Yayasan
    - **Property 4: Read-Only Enforcement untuk Kepala Yayasan**
    - Verifikasi write request pada data operasional oleh Kepala Yayasan ditolak
    - **Validates: Requirements 4.4**

  - [ ]* 3.8 Write property test untuk read-only enforcement Admin BGN
    - **Property 5: Read-Only Enforcement untuk Admin BGN**
    - Verifikasi write request pada data operasional oleh Admin BGN ditolak (kecuali Yayasan/SPPG CRUD)
    - **Validates: Requirements 5.4, 5.5**

  - [ ]* 3.9 Write property test untuk fail-closed tenant middleware
    - **Property 9: Fail-Closed Tenant Middleware**
    - Verifikasi request tanpa tenant context yang valid ditolak
    - **Validates: Requirements 7.6, 14.6**

  - [ ]* 3.10 Write property test untuk cross-tenant access prevention
    - **Property 13: Cross-Tenant Access Prevention**
    - Verifikasi akses resource milik tenant lain mengembalikan 404 (bukan 403)
    - **Validates: Requirements 14.2**

- [x] 4. Checkpoint — Pastikan auth dan tenant middleware berfungsi
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Organization CRUD service dan handler
  - [x] 5.1 Buat organization service di `backend/internal/services/organization_service.go`
    - Implementasikan `YayasanService`: Create, GetAll, GetByID, Update, SetStatus (aktif/nonaktif)
    - Auto-generate kode unik format "YYS-XXXX" pada Create
    - Validasi uniqueness kode, email, NPWP sebelum simpan
    - Cegah pembuatan SPPG baru di bawah Yayasan yang nonaktif
    - Implementasikan `SPPGService`: Create, GetAll, GetByID, Update, SetStatus, Transfer (pindah Yayasan)
    - Auto-generate kode unik format "SPPG-XXXX" pada Create
    - Validasi Yayasan_ID wajib dan valid pada Create SPPG
    - Catat perubahan dalam audit trail
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 3.2, 3.3_

  - [x] 5.2 Buat organization handler di `backend/internal/handlers/organization_handler.go`
    - Implementasikan endpoint Yayasan: POST, GET list, GET detail, PUT, PATCH status
    - Implementasikan endpoint SPPG: POST, GET list, GET detail, PUT, PATCH status, PUT transfer
    - Terapkan role check: hanya superadmin dan admin_bgn yang boleh akses
    - _Requirements: 1.1, 1.4, 1.5, 2.1, 2.6, 3.4, 3.5, 15.3, 15.4_

  - [ ]* 5.3 Write unit tests untuk organization service
    - Test CRUD Yayasan dan SPPG
    - Test validasi uniqueness kode, email, NPWP
    - Test pencegahan SPPG baru di bawah Yayasan nonaktif
    - Test transfer SPPG antar Yayasan
    - _Requirements: 1.1, 1.3, 1.6, 2.1, 2.4, 2.5, 2.6_

- [x] 6. User provisioning service
  - [x] 6.1 Buat provisioning service di `backend/internal/services/provisioning_service.go`
    - Implementasikan `CreateUser()` dengan delegated provisioning rules
    - Superadmin: boleh buat semua peran
    - Admin BGN: tidak boleh buat akun (tolak request)
    - Kepala Yayasan: boleh buat kepala_sppg dan peran operasional untuk SPPG di bawah Yayasan-nya
    - Kepala SPPG: boleh buat peran operasional untuk SPPG-nya
    - Auto-fill `sppg_id` dan `yayasan_id` sesuai konteks pembuat
    - Catat `created_by` pada record User baru
    - Implementasikan `GetUsers()`, `GetUserByID()`, `UpdateUser()`, `SetUserStatus()` dengan tenant scoping
    - _Requirements: 16.1, 16.2, 16.3, 16.4, 16.5, 16.6, 16.7_

  - [x] 6.2 Buat provisioning handler di `backend/internal/handlers/provisioning_handler.go`
    - Implementasikan endpoint: POST /users, GET /users, GET /users/:id, PUT /users/:id, PATCH /users/:id/status
    - Terapkan role-based access control sesuai aturan provisioning
    - _Requirements: 16.1, 16.2, 16.3, 16.4, 16.5_

  - [ ]* 6.3 Write property test untuk delegated provisioning boundary
    - **Property 6: Delegated Provisioning Boundary**
    - Verifikasi setiap pembuat hanya bisa membuat peran dalam cakupan wewenangnya
    - **Validates: Requirements 16.1, 16.2, 16.3, 16.4, 16.5**

  - [ ]* 6.4 Write property test untuk user-tenant consistency
    - **Property 8: User-Tenant Consistency**
    - Verifikasi konsistensi `sppg_id`/`yayasan_id` berdasarkan peran user
    - **Validates: Requirements 8.3, 8.4, 8.5, 8.6**

- [x] 7. Checkpoint — Pastikan organization CRUD dan provisioning berfungsi
  - Ensure all tests pass, ask the user if questions arise.

- [x] 8. Dashboard agregasi backend
  - [x] 8.1 Buat dashboard agregasi service di `backend/internal/services/aggregated_dashboard_service.go`
    - Implementasikan `GetKepalaYayasanDashboard()`: agregasi produksi, pengiriman, keuangan, review dari semua SPPG di bawah Yayasan
    - Implementasikan `GetAdminBGNDashboard()`: agregasi nasional dari seluruh SPPG di semua Yayasan
    - Implementasikan drill-down: filter per SPPG untuk Kepala Yayasan, filter per Yayasan/SPPG untuk Admin BGN
    - Implementasikan filter berdasarkan rentang tanggal
    - Implementasikan daftar SPPG/Yayasan dengan ringkasan performa
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 9.7, 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 10.7, 10.8_

  - [x] 8.2 Perbarui dashboard handler di `backend/internal/handlers/dashboard_handler.go`
    - Tambahkan endpoint `GET /api/v1/dashboard/kepala-yayasan` dengan query params: start_date, end_date, sppg_id
    - Tambahkan endpoint `GET /api/v1/dashboard/admin-bgn` dengan query params: start_date, end_date, yayasan_id, sppg_id
    - Tambahkan endpoint export untuk kedua dashboard
    - Terapkan role check: kepala_yayasan untuk dashboard Yayasan, admin_bgn/superadmin untuk dashboard BGN
    - _Requirements: 9.6, 9.7, 9.8, 10.6, 10.7, 10.8, 10.9_

  - [ ]* 8.3 Write property test untuk dashboard agregasi accuracy
    - **Property 12: Dashboard Agregasi Accuracy**
    - Verifikasi metrik agregat = penjumlahan metrik dari SPPG individual
    - **Validates: Requirements 9.1, 9.2, 9.3, 10.1, 10.2, 10.3**

- [x] 9. Firebase tenant-aware
  - [x] 9.1 Perbarui Firebase service di `backend/internal/firebase/firebase.go`
    - Ubah semua path Firebase untuk menyertakan `sppg_id` segment: `/kds/cooking/{sppg_id}/...`, `/kds/packing/{sppg_id}/...`, dll
    - Tambahkan path baru: `/dashboard/kepala_yayasan/{yayasan_id}/...`, `/dashboard/bgn/...`
    - Perbarui semua service yang menulis ke Firebase (KDS, monitoring, cleaning, dashboard) untuk menggunakan path tenant-aware
    - _Requirements: 11.1, 11.2, 11.6_

  - [x] 9.2 Perbarui Firebase client di `web/src/services/firebase.js` dan `pwa/src/services/firebase.js`
    - Ubah listener path untuk menyertakan `sppg_id` dari auth store
    - Untuk kepala_yayasan: listen ke path semua SPPG di bawah Yayasan
    - Untuk admin_bgn: listen ke path agregat
    - _Requirements: 11.3, 11.4, 11.5_

- [x] 10. Wiring router dan integrasi tenant middleware
  - [x] 10.1 Perbarui router di `backend/internal/router/router.go`
    - Tambahkan tenant middleware pada protected routes setelah JWT auth
    - Daftarkan route baru: `/api/v1/organizations/yayasan/*`, `/api/v1/organizations/sppg/*`
    - Daftarkan route baru: `/api/v1/users/*` untuk provisioning
    - Daftarkan route baru: `/api/v1/dashboard/kepala-yayasan`, `/api/v1/dashboard/admin-bgn`
    - Perbarui `GetMe` endpoint untuk menyertakan `modules` dan `permissions` berdasarkan peran
    - Terapkan role-based route guards pada semua route baru
    - _Requirements: 13.1, 13.2, 13.3, 13.4, 13.5, 17.1, 17.2, 17.3, 17.4, 17.5, 17.6_

  - [x] 10.2 Perbarui semua handler existing untuk menggunakan tenant-scoped DB
    - Modifikasi handler yang mengakses data operasional untuk menggunakan `TenantScope` dari context
    - Pastikan semua query SELECT, INSERT, UPDATE, DELETE melewati tenant scope
    - Pastikan backward compatibility: pengguna SPPG tidak perlu ubah request format
    - _Requirements: 13.1, 13.2, 14.1, 14.4, 14.5_

  - [ ]* 10.3 Write property test untuk backward compatibility
    - **Property 14: Backward Compatibility**
    - Verifikasi endpoint API v1 existing tetap berfungsi untuk pengguna SPPG tanpa perubahan request
    - **Validates: Requirements 13.1, 13.2**

- [x] 11. Checkpoint — Pastikan backend terintegrasi penuh
  - Ensure all tests pass, ask the user if questions arise.

- [x] 12. Web frontend — Auth, navigation, dan module visibility
  - [x] 12.1 Perbarui auth store di `web/src/stores/auth.js`
    - Simpan `sppg_id`, `yayasan_id`, `role`, `modules`, `permissions` dari response `/auth/me`
    - Tambahkan getter untuk cek peran: `isSuperadmin`, `isAdminBGN`, `isKepalaYayasan`
    - _Requirements: 8.8, 17.1, 17.2, 17.3_

  - [x] 12.2 Perbarui router dan navigation di `web/src/router/index.js`
    - Tambahkan route baru: `/yayasan`, `/sppg`, `/users`, `/dashboard-bgn`, `/dashboard-yayasan`
    - Terapkan route guards berdasarkan peran dan module visibility
    - Redirect setelah login sesuai peran: Superadmin → Yayasan, Admin BGN → Dashboard BGN, Kepala Yayasan → Dashboard Yayasan
    - _Requirements: 17.6_

  - [x] 12.3 Perbarui sidebar/navigation di `web/src/layouts/MainLayout.vue` atau `HorizonLayout.vue`
    - Tampilkan menu items berdasarkan `modules` dari auth store
    - Tambahkan menu: Manajemen Yayasan, Manajemen SPPG, Manajemen User, Dashboard BGN, Dashboard Yayasan
    - Sembunyikan modul operasional dari Superadmin, Admin BGN, Kepala Yayasan
    - _Requirements: 17.1, 17.2, 17.3, 17.4, 17.5_

  - [x] 12.4 Perbarui permissions utility di `web/src/utils/permissions.js` dan `web/src/composables/usePermissions.js`
    - Tambahkan permission checks untuk peran baru
    - Tambahkan helper functions untuk module visibility
    - _Requirements: 17.1, 17.2, 17.3, 17.4, 17.5_

- [x] 13. Web frontend — Organization management views
  - [x] 13.1 Buat API service di `web/src/services/organizationService.js`
    - Implementasikan fungsi API untuk CRUD Yayasan dan SPPG
    - Implementasikan fungsi API untuk user provisioning
    - _Requirements: 1.1, 2.1_

  - [x] 13.2 Buat view `web/src/views/YayasanListView.vue`
    - Tabel daftar Yayasan dengan kolom: kode, nama, penanggung jawab, jumlah SPPG, status
    - Form modal untuk create/edit Yayasan
    - Tombol aktifkan/nonaktifkan
    - _Requirements: 1.1, 1.4, 1.5, 3.2_

  - [x] 13.3 Buat view `web/src/views/SPPGListView.vue`
    - Tabel daftar SPPG dengan kolom: kode, nama, Yayasan induk, status
    - Form modal untuk create/edit SPPG dengan dropdown Yayasan
    - Tombol aktifkan/nonaktifkan dan transfer Yayasan
    - _Requirements: 2.1, 2.2, 2.6, 3.3_

  - [x] 13.4 Buat view `web/src/views/UserManagementView.vue`
    - Tabel daftar user dengan kolom: NIK, nama, email, peran, SPPG, Yayasan, status
    - Form modal untuk create/edit user dengan dropdown peran (filtered sesuai wewenang pembuat)
    - Dropdown SPPG dan Yayasan yang di-scope sesuai pembuat
    - Tombol aktifkan/nonaktifkan user
    - _Requirements: 16.1, 16.3, 16.4, 16.5, 16.6_

- [x] 14. Web frontend — Dashboard agregasi views
  - [x] 14.1 Buat view `web/src/views/DashboardBGNView.vue`
    - Kartu metrik agregat nasional: produksi, pengiriman, keuangan, review
    - Tabel daftar Yayasan dengan ringkasan performa
    - Filter: Yayasan, SPPG, rentang tanggal
    - Drill-down ke detail Yayasan dan SPPG
    - Tombol export
    - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6, 10.7, 10.8, 10.9_

  - [x] 14.2 Perbarui view `web/src/views/DashboardKepalaYayasanView.vue`
    - Kartu metrik agregat Yayasan: produksi, pengiriman, keuangan, review
    - Tabel daftar SPPG di bawah Yayasan dengan ringkasan performa
    - Filter: SPPG, rentang tanggal
    - Drill-down ke detail SPPG
    - Tombol export
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 9.5, 9.6, 9.7, 9.8_

  - [x] 14.3 Buat dashboard service di `web/src/services/aggregatedDashboardService.js`
    - Implementasikan API calls untuk dashboard Kepala Yayasan dan Admin BGN
    - _Requirements: 9.1, 10.1_

- [x] 15. Checkpoint — Pastikan web frontend terintegrasi
  - Ensure all tests pass, ask the user if questions arise.

- [x] 16. PWA mobile — Dukungan Admin BGN dan Kepala Yayasan
  - [x] 16.1 Perbarui auth store di `pwa/src/stores/auth.js`
    - Simpan `sppg_id`, `yayasan_id`, `role` dari response login
    - Tambahkan getter untuk cek peran: `isAdminBGN`, `isKepalaYayasan`
    - _Requirements: 18.1_

  - [x] 16.2 Perbarui router di `pwa/src/router/index.js`
    - Tambahkan route baru: `/pwa/dashboard-yayasan`, `/pwa/dashboard-yayasan/:sppg_id`, `/pwa/dashboard-bgn`, `/pwa/dashboard-bgn/:yayasan_id`, `/pwa/dashboard-bgn/:yayasan_id/:sppg_id`
    - Terapkan route guards berdasarkan peran
    - Redirect setelah login sesuai peran
    - _Requirements: 18.1, 18.2, 18.4_

  - [x] 16.3 Buat PWA API service di `pwa/src/services/dashboardAggregatedService.js`
    - Implementasikan API calls untuk dashboard Kepala Yayasan dan Admin BGN
    - _Requirements: 18.2, 18.4_

  - [x] 16.4 Buat view `pwa/src/views/DashboardYayasanView.vue`
    - Daftar SPPG di bawah Yayasan dengan ringkasan performa (mobile-optimized)
    - Kartu metrik agregat
    - Navigasi drill-down ke detail SPPG
    - Indikator SPPG/Yayasan yang sedang dilihat
    - _Requirements: 18.2, 18.3, 18.5, 18.7_

  - [x] 16.5 Buat view `pwa/src/views/DashboardBGNView.vue`
    - Daftar Yayasan dengan ringkasan performa (mobile-optimized)
    - Kartu metrik agregat nasional
    - Navigasi drill-down: Yayasan → SPPG
    - Indikator Yayasan/SPPG yang sedang dilihat
    - _Requirements: 18.4, 18.5, 18.7_

  - [x] 16.6 Implementasikan offline support di `pwa/src/services/db.js`
    - Cache data SPPG yang sudah diakses ke IndexedDB
    - Tampilkan data terakhir yang di-cache saat offline
    - Sync otomatis saat online kembali
    - _Requirements: 18.6_

  - [x] 16.7 Perbarui Firebase listener di `pwa/src/services/firebase.js`
    - Ubah listener path untuk menyertakan `sppg_id`
    - Untuk kepala_yayasan: listen ke path semua SPPG di bawah Yayasan
    - Untuk admin_bgn: listen ke path agregat
    - _Requirements: 11.3, 11.4, 11.5_

- [x] 17. Checkpoint — Pastikan PWA mobile terintegrasi
  - Ensure all tests pass, ask the user if questions arise.

- [x] 18. Audit trail dan keamanan lintas tenant
  - [x] 18.1 Perbarui audit middleware di `backend/internal/middleware/audit.go`
    - Catat `sppg_id` dan `yayasan_id` pada setiap audit trail entry
    - Catat percobaan akses lintas tenant dengan level "warning"
    - _Requirements: 14.3, 15.7, 16.7_

  - [x] 18.2 Perbarui audit handler di `backend/internal/handlers/audit_handler.go`
    - Scope audit trail query berdasarkan peran: superadmin melihat semua, admin_bgn lintas tenant, kepala_yayasan scope Yayasan, kepala_sppg scope SPPG
    - _Requirements: 17.1, 17.2, 17.3_

  - [x] 18.3 Perbarui auth handler `GetMe` di `backend/internal/handlers/auth_handler.go`
    - Perluas response `/auth/me` dengan field `modules`, `permissions`, `sppg`, `yayasan`
    - Kembalikan daftar modul yang visible berdasarkan peran
    - _Requirements: 17.1, 17.2, 17.3, 17.4, 17.5_

- [x] 19. Final checkpoint — Pastikan seluruh sistem terintegrasi
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- Backend menggunakan Go (Golang) + Gin + GORM + PostgreSQL
- Web menggunakan Vue 3 + JavaScript + Ant Design Vue
- PWA menggunakan Vue 3 + JavaScript + Vant UI
- Semua 14 correctness properties dari design.md tercakup dalam property test tasks
