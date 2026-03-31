# Dokumen Persyaratan: Multi-Tenancy & Hierarki Organisasi

## Pendahuluan

Fitur Multi-Tenancy & Hierarki Organisasi mentransformasi Sistem ERP SPPG dari aplikasi single-tenant menjadi platform multi-tenant dengan struktur organisasi bertingkat. Fitur ini memungkinkan satu instalasi sistem melayani banyak SPPG yang dikelola oleh berbagai Yayasan, semuanya di bawah pengawasan BGN (Badan Gizi Nasional).

Hierarki organisasi yang diimplementasikan:

```
Superadmin (administrator teknis — full provisioning seluruh sistem)
  └── BGN (Badan Gizi Nasional) — tingkat tertinggi, mengawasi semua Yayasan dan SPPG
        └── Yayasan (Yayasan/Lembaga) — 1 Yayasan mengelola banyak SPPG
              └── SPPG (unit operasional) — tingkat operasional yang sudah ada
                    └── Karyawan (Kepala SPPG, Chef, Driver, dll.)
```

Delegated provisioning:
- **Superadmin**: Dapat membuat dan mengelola semua entitas (Yayasan, SPPG) dan semua akun pengguna di semua level
- **Kepala Yayasan**: Dapat membuat akun Kepala SPPG dan karyawan untuk SPPG di bawah Yayasan-nya
- **Kepala SPPG**: Dapat membuat akun karyawan operasional di SPPG-nya (seperti saat ini)

Fitur ini menjadi fondasi untuk fitur Risk Assessment di masa depan dan memastikan isolasi data antar tenant sambil memungkinkan agregasi data untuk pelaporan tingkat manajerial.

## Glosarium

- **BGN**: Badan Gizi Nasional — otoritas tingkat tertinggi yang mengawasi seluruh program MBG secara nasional
- **Yayasan**: Yayasan/lembaga yang mengelola satu atau lebih SPPG dalam program MBG
- **SPPG**: Satuan Pelayanan Pemenuhan Gizi — unit operasional yang menjalankan produksi dan distribusi makanan
- **Tenant**: Konteks isolasi data; dalam sistem ini, setiap SPPG merupakan satu tenant
- **Tenant_Context**: Informasi SPPG aktif yang melekat pada sesi pengguna untuk memfilter data secara otomatis
- **Tenant_Middleware**: Komponen middleware yang menyisipkan filter tenant pada setiap query database
- **System**: Sistem ERP SPPG secara keseluruhan (backend, web, dan PWA)
- **Backend**: Server aplikasi Golang yang mengelola logika bisnis dan database
- **Web_App**: Aplikasi web desktop untuk staff kantor dan manajemen
- **RBAC**: Role-Based Access Control — kontrol akses berbasis peran
- **Kepala_Yayasan**: Peran manajerial yang mengawasi semua SPPG di bawah satu Yayasan (read-only)
- **Admin_BGN**: Peran pengawasan tertinggi yang mengawasi semua Yayasan dan SPPG (read-only)
- **SPPG_ID**: Foreign key yang menghubungkan setiap record data operasional ke SPPG tertentu
- **Yayasan_ID**: Foreign key yang menghubungkan setiap SPPG ke Yayasan tertentu
- **Superadmin**: Peran administrator teknis dengan akses penuh ke seluruh sistem termasuk provisioning entitas dan akun
- **Delegated_Provisioning**: Mekanisme di mana pembuatan akun pengguna didelegasikan sesuai hierarki organisasi
- **Aggregated_Dashboard**: Dashboard yang menampilkan metrik gabungan dari beberapa SPPG
- **Drill_Down**: Kemampuan navigasi dari data agregat ke data detail per SPPG
- **Firebase**: Platform Google untuk real-time database dan sinkronisasi data

## Persyaratan

### Persyaratan 1: Entitas Organisasi — Yayasan

**User Story:** Sebagai administrator sistem, saya ingin mendaftarkan dan mengelola data Yayasan agar setiap SPPG memiliki induk organisasi yang jelas dan terstruktur.

#### Acceptance Criteria

1. WHEN seorang Superadmin atau Admin_BGN membuat Yayasan baru, THE System SHALL menyimpan nama, alamat, nomor telepon, email, nama penanggung jawab, dan NPWP Yayasan
2. THE System SHALL menghasilkan kode unik Yayasan secara otomatis dengan format "YYS-XXXX"
3. THE System SHALL memvalidasi bahwa kode Yayasan, email, dan NPWP bersifat unik sebelum menyimpan
4. WHEN seorang Superadmin atau Admin_BGN memperbarui data Yayasan, THE System SHALL mencatat perubahan dalam Audit Trail
5. THE System SHALL memungkinkan Superadmin atau Admin_BGN untuk mengaktifkan atau menonaktifkan Yayasan
6. WHEN sebuah Yayasan dinonaktifkan, THE System SHALL mencegah pembuatan SPPG baru di bawah Yayasan tersebut tetapi tetap mempertahankan data historis

### Persyaratan 2: Entitas Organisasi — SPPG sebagai Tenant

**User Story:** Sebagai administrator sistem, saya ingin setiap SPPG terhubung ke satu Yayasan agar hierarki organisasi terjaga dan data operasional terisolasi per SPPG.

#### Acceptance Criteria

1. WHEN seorang Superadmin atau Admin_BGN membuat SPPG baru, THE System SHALL mewajibkan pemilihan satu Yayasan sebagai induk organisasi
2. THE System SHALL menyimpan data SPPG meliputi nama, alamat, nomor telepon, email, kode SPPG unik, dan Yayasan_ID
3. THE System SHALL memvalidasi bahwa setiap SPPG terhubung ke tepat satu Yayasan
4. THE System SHALL memvalidasi bahwa kode SPPG bersifat unik di seluruh sistem
5. WHEN sebuah Yayasan dinonaktifkan, THE System SHALL tetap mempertahankan hubungan SPPG yang sudah ada dengan Yayasan tersebut
6. THE System SHALL memungkinkan Superadmin atau Admin_BGN untuk memindahkan SPPG dari satu Yayasan ke Yayasan lain dengan mencatat perubahan dalam Audit Trail

### Persyaratan 3: Entitas Organisasi — BGN sebagai Otoritas Tertinggi

**User Story:** Sebagai pejabat BGN, saya ingin memiliki akses pengawasan ke seluruh data Yayasan dan SPPG agar saya dapat memantau program MBG secara nasional.

#### Acceptance Criteria

1. THE System SHALL menyediakan entitas BGN sebagai tingkat tertinggi dalam hierarki organisasi
2. THE System SHALL memungkinkan Admin_BGN untuk melihat daftar semua Yayasan beserta jumlah SPPG yang dikelola masing-masing
3. THE System SHALL memungkinkan Admin_BGN untuk melihat daftar semua SPPG beserta Yayasan induknya
4. WHEN seorang Admin_BGN mengakses data operasional, THE System SHALL menampilkan data dari seluruh SPPG tanpa batasan tenant
5. THE System SHALL menyediakan filter berdasarkan Yayasan dan SPPG pada semua tampilan data untuk Admin_BGN

### Persyaratan 4: Peran Baru — Kepala Yayasan

**User Story:** Sebagai Kepala Yayasan, saya ingin memantau performa operasional semua SPPG di bawah Yayasan saya agar saya dapat mengidentifikasi masalah dan memastikan program berjalan efisien.

#### Acceptance Criteria

1. THE System SHALL mendukung peran "kepala_yayasan" dalam RBAC dengan akses read-only ke data operasional
2. WHEN seorang Kepala_Yayasan login, THE System SHALL membatasi akses data hanya ke SPPG yang berada di bawah Yayasan yang sama
3. THE System SHALL memungkinkan Kepala_Yayasan untuk melihat dashboard, laporan keuangan, data pengiriman, dan metrik performa untuk semua SPPG di bawah Yayasan yang bersangkutan
4. THE System SHALL mencegah Kepala_Yayasan dari melakukan operasi tulis (create, update, delete) pada data operasional SPPG seperti resep, menu, KDS, inventaris, dan purchase order
5. WHEN seorang Kepala_Yayasan mengakses data, THE System SHALL menampilkan indikator SPPG asal pada setiap record data
6. THE System SHALL memungkinkan Kepala_Yayasan untuk memfilter data berdasarkan SPPG tertentu di bawah Yayasan yang bersangkutan

### Persyaratan 5: Peran Baru — Admin BGN

**User Story:** Sebagai pejabat BGN, saya ingin memiliki akses pengawasan lintas Yayasan agar saya dapat memantau seluruh program MBG secara nasional.

#### Acceptance Criteria

1. THE System SHALL mendukung peran "admin_bgn" dalam RBAC dengan akses read-only ke seluruh data operasional
2. WHEN seorang Admin_BGN login, THE System SHALL memberikan akses ke data dari seluruh Yayasan dan seluruh SPPG
3. THE System SHALL memungkinkan Admin_BGN untuk melihat dashboard, laporan keuangan, data pengiriman, dan metrik performa untuk semua SPPG di seluruh Yayasan
4. THE System SHALL mencegah Admin_BGN dari melakukan operasi tulis pada data operasional SPPG kecuali untuk pengelolaan entitas Yayasan dan SPPG
5. THE System SHALL mencegah Admin_BGN dari melakukan user provisioning (pembuatan akun pengguna); provisioning hanya dapat dilakukan oleh Superadmin dan Kepala Yayasan sesuai hierarki
6. WHEN seorang Admin_BGN mengakses data, THE System SHALL menampilkan indikator Yayasan dan SPPG asal pada setiap record data
7. THE System SHALL memungkinkan Admin_BGN untuk memfilter data berdasarkan Yayasan dan SPPG tertentu


### Persyaratan 6: Isolasi Data Tenant pada Tingkat SPPG

**User Story:** Sebagai pengguna SPPG, saya ingin data operasional SPPG saya terisolasi dari SPPG lain agar kerahasiaan dan integritas data terjaga.

#### Acceptance Criteria

1. THE System SHALL menambahkan kolom SPPG_ID sebagai foreign key pada semua tabel data operasional yang ada (Recipe, MenuPlan, Supplier, PurchaseOrder, GoodsReceipt, InventoryItem, School, DeliveryTask, Employee, Attendance, KitchenAsset, CashFlowEntry, dan tabel terkait lainnya)
2. WHEN seorang pengguna tingkat SPPG (Kepala SPPG, Akuntan, Ahli Gizi, Pengadaan, Chef, Packing, Driver, Asisten Lapangan, Kebersihan) mengakses data, THE System SHALL secara otomatis memfilter data berdasarkan SPPG_ID pengguna tersebut
3. THE System SHALL mencegah pengguna tingkat SPPG dari mengakses, membuat, atau memodifikasi data milik SPPG lain
4. WHEN seorang pengguna tingkat SPPG membuat record baru, THE System SHALL secara otomatis menyisipkan SPPG_ID pengguna tersebut ke dalam record
5. IF seorang pengguna tingkat SPPG mencoba mengakses data dengan SPPG_ID yang berbeda, THEN THE System SHALL menolak permintaan dan mengembalikan kode error "FORBIDDEN"
6. THE System SHALL menambahkan indeks database pada kolom SPPG_ID di semua tabel yang relevan untuk menjaga performa query

### Persyaratan 7: Tenant Middleware untuk Penyaringan Data Otomatis

**User Story:** Sebagai pengembang sistem, saya ingin middleware yang secara otomatis menyaring data berdasarkan konteks tenant agar isolasi data terjamin tanpa perlu modifikasi manual di setiap endpoint.

#### Acceptance Criteria

1. THE Tenant_Middleware SHALL mengekstrak SPPG_ID dari JWT token pengguna yang terautentikasi pada setiap permintaan API
2. WHEN pengguna tingkat SPPG mengirim permintaan API, THE Tenant_Middleware SHALL menyisipkan filter `WHERE sppg_id = ?` pada semua query database secara otomatis
3. WHEN pengguna tingkat SPPG mengirim permintaan POST, THE Tenant_Middleware SHALL menyisipkan SPPG_ID ke dalam data yang akan disimpan
4. WHEN seorang Kepala_Yayasan mengirim permintaan API, THE Tenant_Middleware SHALL menyisipkan filter `WHERE sppg_id IN (?)` dengan daftar SPPG_ID di bawah Yayasan yang bersangkutan
5. WHEN seorang Admin_BGN mengirim permintaan API, THE Tenant_Middleware SHALL melewatkan filter tenant sehingga semua data dapat diakses
6. IF JWT token tidak mengandung informasi SPPG_ID yang valid untuk pengguna tingkat SPPG, THEN THE Tenant_Middleware SHALL menolak permintaan dan mengembalikan kode error "UNAUTHORIZED"

### Persyaratan 8: Perubahan Model User untuk Multi-Tenancy

**User Story:** Sebagai administrator sistem, saya ingin model User mendukung informasi tenant agar setiap pengguna terhubung ke SPPG dan Yayasan yang tepat.

#### Acceptance Criteria

1. THE System SHALL menambahkan kolom SPPG_ID (nullable) pada tabel User untuk menghubungkan pengguna ke SPPG tertentu
2. THE System SHALL menambahkan kolom Yayasan_ID (nullable) pada tabel User untuk menghubungkan pengguna ke Yayasan tertentu
3. WHEN seorang pengguna memiliki peran tingkat SPPG, THE System SHALL mewajibkan SPPG_ID terisi pada record User
4. WHEN seorang pengguna memiliki peran "kepala_yayasan", THE System SHALL mewajibkan Yayasan_ID terisi pada record User
5. WHEN seorang pengguna memiliki peran "admin_bgn", THE System SHALL mengizinkan SPPG_ID dan Yayasan_ID bernilai null
6. WHEN seorang pengguna memiliki peran "superadmin", THE System SHALL mengizinkan SPPG_ID dan Yayasan_ID bernilai null
7. THE System SHALL memperbarui validasi peran pada model User untuk menyertakan "superadmin" dan "admin_bgn" sebagai peran yang valid
8. THE System SHALL menyertakan informasi SPPG_ID dan Yayasan_ID dalam payload JWT token saat login

### Persyaratan 9: Dashboard Agregasi untuk Kepala Yayasan

**User Story:** Sebagai Kepala Yayasan, saya ingin melihat dashboard dengan metrik agregat dari semua SPPG di bawah Yayasan saya agar saya dapat memantau performa program secara menyeluruh.

#### Acceptance Criteria

1. WHEN seorang Kepala_Yayasan mengakses dashboard, THE System SHALL menampilkan metrik produksi agregat (total porsi diproduksi, tingkat penyelesaian) dari semua SPPG di bawah Yayasan yang bersangkutan
2. WHEN seorang Kepala_Yayasan mengakses dashboard, THE System SHALL menampilkan metrik pengiriman agregat (total pengiriman selesai, tingkat ketepatan waktu) dari semua SPPG di bawah Yayasan yang bersangkutan
3. WHEN seorang Kepala_Yayasan mengakses dashboard, THE System SHALL menampilkan metrik keuangan agregat (total pengeluaran, penyerapan anggaran) dari semua SPPG di bawah Yayasan yang bersangkutan
4. WHEN seorang Kepala_Yayasan mengakses dashboard, THE System SHALL menampilkan monitoring ulasan (delivery review) agregat dari semua SPPG di bawah Yayasan, termasuk rata-rata rating, jumlah ulasan, dan tren kepuasan penerima manfaat
5. THE System SHALL menampilkan daftar SPPG di bawah Yayasan dengan ringkasan performa masing-masing
6. WHEN seorang Kepala_Yayasan memilih SPPG tertentu pada dashboard, THE System SHALL menampilkan detail operasional SPPG tersebut (drill-down)
7. THE System SHALL memungkinkan Kepala_Yayasan untuk memfilter data dashboard berdasarkan rentang tanggal
8. THE System SHALL memungkinkan ekspor data dashboard Kepala_Yayasan dalam format yang dapat digunakan untuk pelaporan

### Persyaratan 10: Dashboard Agregasi untuk BGN

**User Story:** Sebagai pejabat BGN, saya ingin melihat dashboard dengan metrik agregat dari seluruh Yayasan dan SPPG agar saya dapat memantau program MBG secara nasional.

#### Acceptance Criteria

1. WHEN seorang Admin_BGN mengakses dashboard, THE System SHALL menampilkan metrik produksi agregat dari seluruh SPPG di semua Yayasan
2. WHEN seorang Admin_BGN mengakses dashboard, THE System SHALL menampilkan metrik pengiriman agregat dari seluruh SPPG di semua Yayasan
3. WHEN seorang Admin_BGN mengakses dashboard, THE System SHALL menampilkan metrik keuangan agregat dari seluruh SPPG di semua Yayasan
4. WHEN seorang Admin_BGN mengakses dashboard, THE System SHALL menampilkan monitoring ulasan (delivery review) agregat dari seluruh SPPG di semua Yayasan, termasuk rata-rata rating, jumlah ulasan, dan tren kepuasan penerima manfaat
5. THE System SHALL menampilkan daftar semua Yayasan dengan ringkasan performa masing-masing (jumlah SPPG, total porsi, total pengeluaran)
6. WHEN seorang Admin_BGN memilih Yayasan tertentu pada dashboard, THE System SHALL menampilkan detail performa Yayasan tersebut beserta daftar SPPG-nya (drill-down tingkat Yayasan)
7. WHEN seorang Admin_BGN memilih SPPG tertentu pada dashboard, THE System SHALL menampilkan detail operasional SPPG tersebut (drill-down tingkat SPPG)
8. THE System SHALL memungkinkan Admin_BGN untuk memfilter data dashboard berdasarkan Yayasan, SPPG, dan rentang tanggal
9. THE System SHALL memungkinkan ekspor data dashboard BGN dalam format yang dapat digunakan untuk pelaporan nasional

### Persyaratan 11: Sinkronisasi Firebase Tenant-Aware

**User Story:** Sebagai pengguna sistem, saya ingin data real-time Firebase terisolasi per SPPG agar saya hanya menerima update yang relevan dengan SPPG saya.

#### Acceptance Criteria

1. THE System SHALL mengorganisasi data Firebase dengan path yang menyertakan SPPG_ID (contoh: `/kds/cooking/{sppg_id}/{date}/{recipe_id}`)
2. WHEN data berubah di Backend, THE System SHALL melakukan push update ke path Firebase yang sesuai dengan SPPG_ID terkait
3. WHEN klien Web_App atau PWA_App terhubung ke Firebase, THE System SHALL hanya mendengarkan path Firebase yang sesuai dengan SPPG_ID pengguna
4. WHEN seorang Kepala_Yayasan terhubung ke Firebase, THE System SHALL mendengarkan path Firebase untuk semua SPPG di bawah Yayasan yang bersangkutan
5. WHEN seorang Admin_BGN terhubung ke Firebase, THE System SHALL mendengarkan path Firebase tingkat agregat
6. THE System SHALL memperbarui struktur path Firebase yang ada untuk menyertakan segmen SPPG_ID tanpa mengganggu fungsionalitas real-time

### Persyaratan 12: Migrasi Data Single-Tenant ke Multi-Tenant

**User Story:** Sebagai administrator sistem, saya ingin data yang sudah ada bermigrasi dengan aman ke struktur multi-tenant agar tidak ada data yang hilang dan sistem tetap berfungsi normal.

#### Acceptance Criteria

1. WHEN migrasi dijalankan, THE System SHALL membuat satu Yayasan default dan satu SPPG default untuk menampung data yang sudah ada
2. WHEN migrasi dijalankan, THE System SHALL mengisi kolom SPPG_ID pada semua record data operasional yang ada dengan ID SPPG default
3. WHEN migrasi dijalankan, THE System SHALL mengisi kolom Yayasan_ID pada record SPPG default dengan ID Yayasan default
4. WHEN migrasi dijalankan, THE System SHALL mengisi kolom SPPG_ID pada semua record User yang memiliki peran tingkat SPPG dengan ID SPPG default
5. THE System SHALL memvalidasi bahwa semua record data operasional memiliki SPPG_ID yang valid setelah migrasi selesai
6. THE System SHALL menyediakan laporan migrasi yang menunjukkan jumlah record yang berhasil dimigrasi per tabel
7. IF migrasi gagal pada salah satu tabel, THEN THE System SHALL melakukan rollback seluruh perubahan dan melaporkan error

### Persyaratan 13: Kompatibilitas Mundur API

**User Story:** Sebagai pengguna SPPG yang sudah ada, saya ingin endpoint API yang sudah ada tetap berfungsi normal setelah multi-tenancy diaktifkan agar pekerjaan saya tidak terganggu.

#### Acceptance Criteria

1. THE System SHALL memastikan semua endpoint API v1 yang ada tetap berfungsi untuk pengguna tingkat SPPG tanpa perubahan pada request format
2. WHEN pengguna tingkat SPPG mengakses endpoint API yang ada, THE System SHALL secara otomatis menerapkan filter tenant berdasarkan JWT token tanpa memerlukan parameter SPPG_ID eksplisit dalam request
3. THE System SHALL menambahkan endpoint API baru untuk pengelolaan Yayasan dan SPPG di bawah prefix `/api/v1/organizations/`
4. THE System SHALL menambahkan endpoint API baru untuk dashboard agregasi di bawah prefix `/api/v1/dashboard/`
5. WHEN pengguna tingkat SPPG mengakses endpoint dashboard yang ada (`/api/v1/dashboard/kepala-sppg`), THE System SHALL mengembalikan data yang sama seperti sebelumnya, hanya untuk SPPG pengguna tersebut
6. THE System SHALL mendokumentasikan semua perubahan API dalam changelog

### Persyaratan 14: Keamanan dan Validasi Akses Lintas Tenant

**User Story:** Sebagai administrator keamanan, saya ingin memastikan tidak ada kebocoran data antar tenant agar kerahasiaan data setiap SPPG terjaga.

#### Acceptance Criteria

1. THE System SHALL memvalidasi bahwa setiap permintaan API mengandung Tenant_Context yang valid sebelum mengeksekusi query database
2. WHEN seorang pengguna mencoba mengakses resource dengan ID yang bukan milik tenant-nya, THE System SHALL mengembalikan response "404 Not Found" alih-alih "403 Forbidden" untuk mencegah enumerasi data
3. THE System SHALL mencatat setiap percobaan akses lintas tenant dalam Audit Trail dengan level "warning"
4. THE System SHALL memvalidasi bahwa operasi bulk (import, export) hanya memproses data dalam scope tenant pengguna
5. THE System SHALL memastikan bahwa query JOIN antar tabel tetap menerapkan filter SPPG_ID pada semua tabel yang terlibat
6. IF terjadi kegagalan pada Tenant_Middleware, THEN THE System SHALL menolak permintaan secara default (fail-closed) alih-alih mengizinkan akses tanpa filter

### Persyaratan 15: Peran Baru — Superadmin

**User Story:** Sebagai Superadmin, saya ingin memiliki akses penuh ke seluruh sistem termasuk provisioning entitas organisasi dan akun pengguna di semua level agar saya dapat mengelola infrastruktur sistem secara menyeluruh.

#### Acceptance Criteria

1. THE System SHALL mendukung peran "superadmin" dalam RBAC dengan akses penuh (read dan write) ke seluruh sistem
2. WHEN seorang Superadmin login, THE System SHALL memberikan akses tanpa batasan tenant ke seluruh data dan fitur
3. THE System SHALL memungkinkan Superadmin untuk membuat, memperbarui, dan menonaktifkan entitas Yayasan
4. THE System SHALL memungkinkan Superadmin untuk membuat, memperbarui, dan menonaktifkan entitas SPPG
5. THE System SHALL memungkinkan Superadmin untuk membuat akun pengguna di semua level (Admin BGN, Kepala Yayasan, Kepala SPPG, dan semua peran karyawan operasional)
6. THE System SHALL memungkinkan Superadmin untuk mengubah peran dan afiliasi tenant (SPPG_ID, Yayasan_ID) pada akun pengguna yang sudah ada
7. THE System SHALL mencatat semua aksi Superadmin dalam Audit Trail

### Persyaratan 16: User Provisioning Berjenjang (Delegated Provisioning)

**User Story:** Sebagai pengelola di setiap tingkat hierarki, saya ingin dapat membuat akun pengguna sesuai cakupan tanggung jawab saya agar proses onboarding karyawan baru efisien dan terkontrol.

#### Acceptance Criteria

1. WHEN seorang Superadmin membuat akun pengguna, THE System SHALL mengizinkan pemilihan peran apa pun termasuk "superadmin", "admin_bgn", "kepala_yayasan", "kepala_sppg", dan semua peran karyawan operasional
2. WHEN seorang Admin_BGN membuat akun pengguna, THE System SHALL menolak permintaan karena Admin_BGN tidak memiliki wewenang user provisioning
3. WHEN seorang Kepala_Yayasan membuat akun pengguna, THE System SHALL hanya mengizinkan pemilihan peran "kepala_sppg" dan peran karyawan operasional untuk SPPG yang berada di bawah Yayasan yang bersangkutan
4. WHEN seorang Kepala_SPPG membuat akun pengguna, THE System SHALL hanya mengizinkan pemilihan peran karyawan operasional (Akuntan, Ahli Gizi, Pengadaan, Chef, Packing, Driver, Asisten Lapangan, Kebersihan) untuk SPPG yang bersangkutan
5. IF seorang pengguna mencoba membuat akun dengan peran di luar cakupan wewenangnya, THEN THE System SHALL menolak permintaan dan mengembalikan pesan error yang jelas
6. WHEN akun pengguna baru dibuat, THE System SHALL secara otomatis mengisi SPPG_ID dan/atau Yayasan_ID sesuai konteks pembuat akun
7. THE System SHALL mencatat siapa yang membuat setiap akun pengguna dalam Audit Trail

### Persyaratan 17: Matriks Visibilitas Modul per Peran

**User Story:** Sebagai pengguna sistem, saya ingin hanya melihat modul dan fitur yang relevan dengan peran saya agar antarmuka tidak membingungkan dan navigasi lebih efisien.

#### Acceptance Criteria

1. WHEN seorang Superadmin login, THE System SHALL menampilkan modul: Manajemen Yayasan, Manajemen SPPG, Manajemen User (semua level), System Configuration, Audit Trail (seluruh sistem), dan semua modul yang tersedia untuk Admin_BGN
2. WHEN seorang Admin_BGN login, THE System SHALL menampilkan modul: Dashboard BGN (agregasi nasional), Daftar Yayasan, Daftar SPPG, Manajemen Yayasan (CRUD), Manajemen SPPG (CRUD), Laporan Keuangan (agregat, read-only), Data Pengiriman/Logistik (agregat, read-only), Data Inventaris (agregat, read-only), Monitoring Ulasan/Review (agregat, read-only), dan Audit Trail (lintas tenant)
3. WHEN seorang Kepala_Yayasan login, THE System SHALL menampilkan modul: Dashboard Yayasan (agregasi SPPG di bawahnya), Daftar SPPG di bawah Yayasan, User Provisioning (Kepala SPPG dan karyawan), Laporan Keuangan (agregat SPPG-nya, read-only), Data Pengiriman/Logistik (agregat, read-only), Data Inventaris (agregat, read-only), Monitoring Ulasan/Review (agregat, read-only), dan Audit Trail (scope Yayasan)
4. THE System SHALL menyembunyikan modul operasional harian (KDS, Menu Planning editor, PO creation, Cooking status, Packing, Absensi) dari peran Superadmin, Admin_BGN, dan Kepala_Yayasan
5. THE System SHALL mempertahankan visibilitas modul yang sudah ada untuk peran tingkat SPPG (Kepala SPPG, Akuntan, Ahli Gizi, Pengadaan, Chef, Packing, Driver, Asisten Lapangan, Kebersihan) tanpa perubahan
6. THE System SHALL mengarahkan pengguna ke halaman default yang sesuai dengan perannya setelah login (Superadmin → Manajemen Yayasan, Admin_BGN → Dashboard BGN, Kepala_Yayasan → Dashboard Yayasan)

### Persyaratan 18: Akses PWA Mobile untuk Admin BGN dan Kepala Yayasan

**User Story:** Sebagai Kepala Yayasan atau pejabat BGN, saya ingin mengakses data monitoring melalui aplikasi mobile (PWA) saat melakukan kunjungan lapangan ke SPPG agar saya dapat langsung melihat informasi operasional SPPG yang sedang dikunjungi.

#### Acceptance Criteria

1. THE PWA_App SHALL mendukung login untuk peran "admin_bgn" dan "kepala_yayasan" selain peran operasional yang sudah ada (Driver, Asisten Lapangan)
2. WHEN seorang Kepala_Yayasan login ke PWA_App, THE System SHALL menampilkan daftar SPPG di bawah Yayasan-nya dengan ringkasan performa masing-masing
3. WHEN seorang Kepala_Yayasan memilih SPPG tertentu di PWA_App, THE System SHALL menampilkan detail operasional SPPG tersebut meliputi: status produksi hari ini, status pengiriman, ringkasan keuangan, dan ulasan terbaru
4. WHEN seorang Admin_BGN login ke PWA_App, THE System SHALL menampilkan daftar Yayasan dengan ringkasan performa, dan memungkinkan drill-down ke SPPG tertentu
5. THE PWA_App SHALL menyediakan tampilan dashboard yang dioptimalkan untuk layar mobile bagi peran Admin_BGN dan Kepala_Yayasan
6. THE PWA_App SHALL mendukung kemampuan offline untuk data SPPG yang sudah pernah diakses, sehingga informasi tetap tersedia saat koneksi internet terbatas di lokasi kunjungan
7. THE PWA_App SHALL menampilkan indikator SPPG dan Yayasan yang sedang dilihat secara jelas pada setiap halaman monitoring
