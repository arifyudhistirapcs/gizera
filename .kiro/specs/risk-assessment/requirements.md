# Requirements Document

## Introduction

Fitur Risk Assessment memungkinkan peran kepala_yayasan untuk melakukan audit kepatuhan SOP (Standard Operating Procedure) terhadap SPPG-SPPG yang dikelolanya. Audit dilakukan melalui formulir penilaian risiko yang diisi saat kunjungan lapangan menggunakan PWA mobile app, dan hasilnya dapat dipantau melalui web app. Formulir audit didasarkan pada dokumen SOP Dapur MBG yang mencakup area seperti higienitas dapur, standar persiapan makanan, penyimpanan, pengiriman, kebersihan staf, pemeliharaan peralatan, dan dokumentasi.

## Glossary

- **Risk_Assessment_Form**: Formulir audit yang diisi oleh Kepala Yayasan untuk menilai kepatuhan SPPG terhadap SOP
- **SOP_Category**: Kategori utama dalam dokumen SOP Dapur MBG (misalnya: Higienitas Dapur, Persiapan Makanan, Penyimpanan, Pengiriman, Kebersihan Staf, Peralatan, Dokumentasi)
- **SOP_Checklist_Item**: Item individual dalam SOP yang dinilai kepatuhannya (misalnya: "Lantai dapur bersih dan kering")
- **Compliance_Score**: Nilai kepatuhan per item checklist, menggunakan skala numerik (1-5) di mana 1 = Tidak Patuh, 5 = Sangat Patuh
- **Overall_Risk_Score**: Skor risiko keseluruhan dari satu formulir audit, dihitung dari rata-rata semua Compliance_Score
- **Risk_Level**: Klasifikasi tingkat risiko berdasarkan Overall_Risk_Score: Rendah (4.0-5.0), Sedang (2.5-3.9), Tinggi (1.0-2.4)
- **Kepala_Yayasan**: Peran pengguna yang bertanggung jawab mengawasi SPPG di bawah yayasannya, memiliki akses read-only ke data operasional dan akses penuh ke fitur Risk Assessment
- **SPPG**: Satuan Pelayanan Pemenuhan Gizi, unit operasional dapur yang menjadi objek audit
- **Yayasan**: Organisasi induk yang mengelola satu atau lebih SPPG
- **PWA_App**: Progressive Web Application untuk operasi lapangan
- **Web_App**: Aplikasi web utama dengan Ant Design Vue untuk monitoring dan manajemen
- **Assessment_Evidence**: Foto atau catatan pendukung yang dilampirkan pada item checklist sebagai bukti audit
- **SPPG_Operational_Snapshot**: Data operasional SPPG yang di-capture otomatis saat formulir risk assessment dibuat, sebagai konteks tambahan untuk penilaian

## Requirements

### Requirement 1: Manajemen Template SOP

**User Story:** Sebagai superadmin, saya ingin mengelola template SOP checklist, sehingga item-item audit selalu sesuai dengan dokumen SOP terbaru.

#### Acceptance Criteria

1. THE System SHALL menyimpan SOP_Checklist_Item dengan atribut: nama item, deskripsi, SOP_Category, dan urutan tampil
2. THE System SHALL mengelompokkan SOP_Checklist_Item berdasarkan SOP_Category
3. WHEN superadmin menambahkan SOP_Checklist_Item baru, THE System SHALL menyimpan item tersebut dan menetapkan urutan tampil secara otomatis
4. WHEN superadmin mengubah SOP_Checklist_Item, THE System SHALL memperbarui item tersebut tanpa memengaruhi Risk_Assessment_Form yang sudah diisi sebelumnya
5. WHEN superadmin menonaktifkan SOP_Checklist_Item, THE System SHALL menandai item tersebut sebagai tidak aktif dan mengecualikannya dari formulir audit baru
6. THE System SHALL menyediakan daftar SOP_Category default berdasarkan dokumen SOP Dapur MBG: Higienitas Dapur dan Sanitasi, Standar Persiapan Makanan, Penyimpanan dan Kontrol Suhu, Prosedur Pengiriman, Kebersihan Staf dan APD, Pemeliharaan Peralatan, Dokumentasi dan Pencatatan

### Requirement 2: Pembuatan Risk Assessment Form di PWA

**User Story:** Sebagai kepala_yayasan, saya ingin membuat dan mengisi formulir risk assessment saat kunjungan lapangan ke SPPG, sehingga saya dapat mendokumentasikan hasil audit secara langsung.

#### Acceptance Criteria

1. WHEN Kepala_Yayasan membuka fitur Risk Assessment di PWA_App, THE System SHALL menampilkan daftar SPPG yang berada di bawah Yayasan Kepala_Yayasan tersebut
2. WHEN Kepala_Yayasan memilih SPPG untuk diaudit, THE System SHALL membuat Risk_Assessment_Form baru dengan semua SOP_Checklist_Item aktif yang dikelompokkan berdasarkan SOP_Category
3. WHEN Kepala_Yayasan mengisi Compliance_Score untuk sebuah SOP_Checklist_Item, THE System SHALL menyimpan skor tersebut dengan nilai antara 1 sampai 5
4. WHEN Kepala_Yayasan menambahkan catatan pada SOP_Checklist_Item, THE System SHALL menyimpan catatan teks tersebut bersama item yang bersangkutan
5. WHEN Kepala_Yayasan melampirkan foto sebagai Assessment_Evidence pada SOP_Checklist_Item, THE System SHALL mengunggah dan menyimpan foto tersebut terkait dengan item yang bersangkutan
6. WHILE Risk_Assessment_Form belum disubmit, THE System SHALL menyimpan progres pengisian secara otomatis sebagai draft
7. IF koneksi jaringan terputus saat pengisian formulir, THEN THE System SHALL menyimpan data formulir secara lokal di perangkat dan melakukan sinkronisasi otomatis saat koneksi tersedia kembali

### Requirement 3: Submit dan Kalkulasi Skor

**User Story:** Sebagai kepala_yayasan, saya ingin mengirimkan formulir risk assessment yang sudah lengkap, sehingga hasil audit tercatat secara resmi.

#### Acceptance Criteria

1. WHEN Kepala_Yayasan menekan tombol submit pada Risk_Assessment_Form, THE System SHALL memvalidasi bahwa semua SOP_Checklist_Item memiliki Compliance_Score
2. IF terdapat SOP_Checklist_Item yang belum memiliki Compliance_Score saat submit, THEN THE System SHALL menampilkan pesan error yang menyebutkan item mana saja yang belum dinilai
3. WHEN Risk_Assessment_Form berhasil disubmit, THE System SHALL menghitung Overall_Risk_Score sebagai rata-rata dari semua Compliance_Score dalam formulir tersebut
4. WHEN Overall_Risk_Score telah dihitung, THE System SHALL menetapkan Risk_Level berdasarkan skor: Rendah untuk 4.0-5.0, Sedang untuk 2.5-3.9, Tinggi untuk 1.0-2.4
5. WHEN Risk_Assessment_Form berhasil disubmit, THE System SHALL mencatat tanggal submit dan mengubah status formulir dari draft menjadi submitted
6. WHEN Risk_Assessment_Form sudah berstatus submitted, THE System SHALL mencegah perubahan pada formulir tersebut

### Requirement 4: Monitoring Risk Assessment di Web App

**User Story:** Sebagai kepala_yayasan, saya ingin melihat riwayat dan ringkasan hasil risk assessment di web app, sehingga saya dapat memantau tren kepatuhan SPPG dari waktu ke waktu.

#### Acceptance Criteria

1. WHEN Kepala_Yayasan membuka halaman Risk Assessment di Web_App, THE System SHALL menampilkan daftar semua Risk_Assessment_Form milik SPPG di bawah Yayasan tersebut, diurutkan berdasarkan tanggal terbaru
2. THE System SHALL menampilkan informasi berikut pada setiap item daftar: nama SPPG, tanggal audit, Overall_Risk_Score, Risk_Level, dan status formulir
3. WHEN Kepala_Yayasan memilih Risk_Assessment_Form dari daftar, THE System SHALL menampilkan detail lengkap formulir termasuk semua SOP_Checklist_Item beserta Compliance_Score, catatan, dan Assessment_Evidence
4. THE System SHALL menyediakan filter berdasarkan SPPG, rentang tanggal, dan Risk_Level pada halaman daftar Risk Assessment
5. THE System SHALL menampilkan ringkasan statistik per SPPG yang mencakup: jumlah audit yang dilakukan, rata-rata Overall_Risk_Score, dan tren Risk_Level dari waktu ke waktu

### Requirement 5: Tenant Isolation untuk Risk Assessment

**User Story:** Sebagai pengelola sistem, saya ingin memastikan data risk assessment terisolasi sesuai hierarki tenant, sehingga setiap yayasan hanya dapat mengakses data audit SPPG miliknya.

#### Acceptance Criteria

1. THE System SHALL menyimpan yayasan_id dan sppg_id pada setiap Risk_Assessment_Form
2. WHEN Kepala_Yayasan mengakses data Risk Assessment, THE System SHALL hanya menampilkan formulir dari SPPG yang berada di bawah Yayasan Kepala_Yayasan tersebut
3. WHEN superadmin mengakses data Risk Assessment, THE System SHALL menampilkan semua formulir dari seluruh SPPG
4. IF pengguna mencoba mengakses Risk_Assessment_Form milik SPPG di luar cakupan tenant-nya, THEN THE System SHALL mengembalikan respons 404 Not Found
5. THE System SHALL menerapkan tenant filtering melalui middleware yang sudah ada secara konsisten pada semua endpoint Risk Assessment

### Requirement 6: API Risk Assessment

**User Story:** Sebagai pengembang, saya ingin API yang konsisten untuk operasi CRUD risk assessment, sehingga web app dan PWA app dapat mengakses data dengan cara yang seragam.

#### Acceptance Criteria

1. THE System SHALL menyediakan endpoint REST API untuk operasi berikut: membuat formulir baru, mengambil daftar formulir, mengambil detail formulir, memperbarui draft formulir, dan mengirimkan formulir
2. WHEN API menerima request pembuatan Risk_Assessment_Form, THE System SHALL memvalidasi bahwa SPPG yang dipilih berada di bawah Yayasan pengguna yang membuat request
3. WHEN API menerima request dengan data tidak valid, THE System SHALL mengembalikan respons error dengan kode error dan pesan yang deskriptif
4. THE System SHALL menggunakan format respons API yang konsisten dengan endpoint lain dalam sistem: field success, data, message, dan error_code
5. WHEN Risk_Assessment_Form disubmit melalui API, THE System SHALL mencatat aksi tersebut dalam audit trail

### Requirement 7: Skor per Kategori SOP

**User Story:** Sebagai kepala_yayasan, saya ingin melihat skor kepatuhan per kategori SOP, sehingga saya dapat mengidentifikasi area spesifik yang memerlukan perbaikan di setiap SPPG.

#### Acceptance Criteria

1. WHEN Risk_Assessment_Form disubmit, THE System SHALL menghitung skor rata-rata per SOP_Category dari Compliance_Score item-item dalam kategori tersebut
2. THE System SHALL menyimpan skor per SOP_Category bersama Risk_Assessment_Form
3. WHEN detail Risk_Assessment_Form ditampilkan, THE System SHALL menampilkan skor per SOP_Category dalam format visual yang mudah dipahami
4. THE System SHALL menetapkan Risk_Level per SOP_Category menggunakan skala yang sama dengan Overall_Risk_Score

### Requirement 8: Snapshot Data Operasional SPPG

**User Story:** Sebagai kepala_yayasan, saya ingin formulir risk assessment secara otomatis merekam data operasional terkini dari SPPG yang diaudit, sehingga saya memiliki konteks lengkap tentang kondisi SPPG saat audit dilakukan dan data tersebut dapat memengaruhi penilaian keseluruhan.

#### Acceptance Criteria

1. WHEN Risk_Assessment_Form baru dibuat untuk sebuah SPPG, THE System SHALL secara otomatis mengambil dan menyimpan snapshot data operasional SPPG tersebut pada saat itu
2. THE SPPG_Operational_Snapshot SHALL mencakup data berikut:
   - Rata-rata rating review dari sekolah (overall, menu, layanan)
   - Jumlah total ulasan yang diterima
   - Total pemasukan dan pengeluaran keuangan bulan berjalan
   - Penyerapan anggaran bulan berjalan (persentase)
   - Jumlah pengiriman selesai dan tingkat ketepatan waktu (on-time delivery rate)
   - Jumlah porsi yang diproduksi bulan berjalan
   - Jumlah item stok yang berada di bawah threshold minimum (stok kritis)
   - Jumlah karyawan aktif
   - Tingkat kehadiran karyawan bulan berjalan (persentase)
3. THE System SHALL menyimpan SPPG_Operational_Snapshot sebagai bagian dari Risk_Assessment_Form sehingga data tersebut tidak berubah meskipun data operasional SPPG berubah di kemudian hari
4. WHEN detail Risk_Assessment_Form ditampilkan di Web_App atau PWA_App, THE System SHALL menampilkan SPPG_Operational_Snapshot dalam section terpisah dengan label "Data Operasional SPPG Saat Audit"
5. WHEN menghitung Overall_Risk_Score, THE System SHALL memperhitungkan data operasional sebagai faktor tambahan: SPPG dengan rating review rendah (<3.0), penyerapan anggaran rendah (<50%), atau on-time delivery rate rendah (<70%) SHALL mendapat penalti pengurangan skor sebesar 0.5 poin pada Overall_Risk_Score
6. THE System SHALL menampilkan indikator visual (warna merah/kuning/hijau) pada setiap metrik operasional berdasarkan threshold: hijau (baik), kuning (perlu perhatian), merah (kritis)
