# Requirements Document

## Introduction

Dokumen ini mendefinisikan kebutuhan untuk rombak total UI aplikasi PWA ERP SPPG dan penambahan modul-modul baru (view only). Redesign mengacu pada style referensi HR Attendee app yang modern dan clean, dengan tetap mempertahankan color palette purple (#5A4372) yang sudah ada, seluruh business logic, dan flow yang sudah berjalan. Saat ini PWA memiliki 2 modul: Absensi (semua karyawan) dan Tugas (driver & asisten lapangan). Akan ditambahkan 3 modul baru: Dashboard, Monitoring Aktivitas, dan Perencanaan Menu, serta fitur pendukung view-only pada modul yang sudah ada.

## Glossary

- **PWA**: Progressive Web Application ERP SPPG yang dibangun dengan Vue 3, Vant UI, dan Pinia
- **Attendance_Module**: Modul absensi yang menangani check-in/check-out karyawan dengan validasi Wi-Fi/GPS
- **Task_Module**: Modul tugas pengiriman dan pengambilan untuk driver dan asisten lapangan
- **Dashboard_Module**: Modul baru untuk menampilkan ringkasan data operasional (view only)
- **Activity_Monitoring_Module**: Modul baru untuk menampilkan monitoring aktivitas (view only)
- **Menu_Planning_Module**: Modul baru untuk menampilkan perencanaan menu dan approval menu mingguan
- **Bottom_Navigation**: Komponen navigasi utama di bagian bawah layar mengikuti style HR Attendee app
- **Card_Component**: Komponen kartu dengan rounded corners, shadow, dan layout modern
- **Kepala_SPPG**: Role kepala SPPG/yayasan yang memiliki akses ke modul Dashboard, Monitoring Aktivitas, dan Perencanaan Menu
- **Ahli_Gizi**: Role ahli gizi yang memiliki akses ke modul Perencanaan Menu (view only)
- **Driver**: Role driver yang memiliki akses ke modul Tugas
- **Asisten_Lapangan**: Role asisten lapangan yang memiliki akses ke modul Tugas
- **Karyawan**: Semua role karyawan yang memiliki akses ke modul Absensi
- **Color_Palette**: Skema warna yang sudah ada: primary #5A4372, accent #3D2B53, bg #F8FDEA, text #322837
- **Swipe_Action**: Gesture geser horizontal untuk melakukan aksi (check-in/check-out) mengikuti referensi HR Attendee
- **Summary_Card**: Kartu ringkasan dengan ikon, label, dan nilai numerik
- **Date_Selector**: Komponen pemilih tanggal di bagian atas halaman
- **Activity_Log**: Daftar log aktivitas dengan timestamp dan status
- **Approval_Action**: Aksi approve/reject pada menu mingguan oleh Kepala_SPPG
- **Skeleton_Loading**: Placeholder animasi saat data sedang dimuat
- **Pull_To_Refresh**: Gesture tarik ke bawah untuk memperbarui data
- **Sekolah**: Role sekolah yang memiliki akun login sendiri, terhubung ke satu sekolah spesifik di database, dan memiliki akses ke modul Monitoring Sekolah (view only)
- **School_Monitoring_Module**: Modul baru khusus role Sekolah untuk melihat menu hari ini dan status/progres pengiriman ke sekolah mereka (view only)

## Requirements

### Requirement 1: Redesign Halaman Login

**User Story:** Sebagai karyawan, saya ingin halaman login yang modern dan clean mengikuti style HR Attendee app, sehingga pengalaman pertama menggunakan aplikasi terasa profesional.

#### Acceptance Criteria

1. THE PWA SHALL menampilkan halaman login dengan layout single-column centered yang memiliki background gradient menggunakan Color_Palette primary (#5A4372) dan accent (#3D2B53)
2. THE PWA SHALL menampilkan logo aplikasi, judul "ERP SPPG", dan subtitle di bagian atas halaman login
3. THE PWA SHALL menampilkan form login dengan field NIK/Email dan Password di dalam Card_Component dengan border-radius 16px dan shadow
4. THE PWA SHALL menampilkan tombol login full-width dengan tinggi minimal 48px, border-radius 12px, dan warna primary #5A4372
5. WHEN karyawan menekan tombol login dengan kredensial valid, THE PWA SHALL mengautentikasi dan mengarahkan karyawan ke halaman utama sesuai role
6. IF karyawan memasukkan kredensial tidak valid, THEN THE PWA SHALL menampilkan pesan error yang deskriptif tanpa mereset field NIK/Email
7. WHILE proses login berlangsung, THE PWA SHALL menampilkan loading indicator pada tombol login dan menonaktifkan interaksi form

### Requirement 2: Redesign Layout Utama dan Bottom Navigation

**User Story:** Sebagai karyawan, saya ingin navigasi yang mudah diakses dengan bottom navigation bar mengikuti style HR Attendee app, sehingga saya dapat berpindah antar modul dengan cepat.

#### Acceptance Criteria

1. THE PWA SHALL menampilkan Bottom_Navigation dengan tinggi 56px, background putih, shadow atas, dan ikon berukuran 22px di setiap halaman setelah login
2. THE PWA SHALL menampilkan item Bottom_Navigation sesuai role karyawan yang sedang login
3. WHEN karyawan dengan role Driver atau Asisten_Lapangan login, THE PWA SHALL menampilkan item Bottom_Navigation: Tugas, Absensi, dan Profil
4. WHEN karyawan dengan role Kepala_SPPG login, THE PWA SHALL menampilkan item Bottom_Navigation: Dashboard, Monitoring, Menu, Absensi, dan Profil
5. WHEN karyawan dengan role Ahli_Gizi login, THE PWA SHALL menampilkan item Bottom_Navigation: Menu, Absensi, dan Profil
6. WHEN karyawan dengan role Sekolah login, THE PWA SHALL menampilkan item Bottom_Navigation: Monitoring dan Profil
7. WHEN karyawan dengan role selain yang disebutkan login, THE PWA SHALL menampilkan item Bottom_Navigation: Absensi dan Profil
8. THE PWA SHALL menandai item Bottom_Navigation yang aktif dengan warna primary #5A4372 dan ikon filled
9. WHEN karyawan menekan item Bottom_Navigation, THE PWA SHALL menavigasi ke halaman terkait dengan transisi halus dalam waktu kurang dari 300ms

### Requirement 3: Redesign Modul Absensi

**User Story:** Sebagai karyawan, saya ingin tampilan absensi yang modern dengan card-based layout, date selector, dan swipe action mengikuti style HR Attendee app, sehingga proses check-in/check-out lebih intuitif.

#### Acceptance Criteria

1. THE Attendance_Module SHALL menampilkan header dengan nama karyawan, tanggal hari ini, dan Date_Selector untuk melihat riwayat absensi
2. THE Attendance_Module SHALL menampilkan status absensi hari ini dalam Card_Component dengan layout dua kolom: kartu Check-In (hijau) dan kartu Check-Out (merah) yang menampilkan waktu masing-masing
3. THE Attendance_Module SHALL menampilkan Summary_Card untuk total hari kerja, jam kerja, dan status kehadiran bulan berjalan
4. THE Attendance_Module SHALL menampilkan Swipe_Action untuk melakukan check-in (geser ke kanan) dan check-out (geser ke kiri) sebagai alternatif tombol
5. WHEN karyawan melakukan swipe check-in, THE Attendance_Module SHALL menjalankan proses validasi Wi-Fi/GPS dan check-in yang sama dengan flow check-in tombol yang sudah ada
6. WHEN karyawan melakukan swipe check-out, THE Attendance_Module SHALL menjalankan proses konfirmasi dan check-out yang sama dengan flow check-out tombol yang sudah ada
7. THE Attendance_Module SHALL menampilkan Activity_Log berupa daftar riwayat absensi 7 hari terakhir dengan format kartu yang menampilkan tanggal, waktu masuk, waktu keluar, dan durasi kerja
8. WHILE data absensi sedang dimuat, THE Attendance_Module SHALL menampilkan Skeleton_Loading pada setiap Card_Component
9. THE Attendance_Module SHALL mendukung Pull_To_Refresh untuk memperbarui data absensi
10. THE Attendance_Module SHALL mempertahankan seluruh business logic check-in/check-out yang sudah ada termasuk validasi Wi-Fi, GPS, dan offline support

### Requirement 4: Redesign Modul Tugas

**User Story:** Sebagai driver atau asisten lapangan, saya ingin tampilan daftar tugas yang modern dengan card-based layout, tab selection untuk memisahkan pengiriman dan pengambilan, dan informasi yang lebih terstruktur, sehingga saya dapat mengelola tugas pengiriman dan pengambilan dengan lebih efisien.

#### Acceptance Criteria

1. THE Task_Module SHALL menampilkan header dengan judul halaman, jumlah tugas hari ini, dan tombol refresh
2. THE Task_Module SHALL menampilkan Summary_Card di bagian atas yang menunjukkan jumlah tugas berdasarkan status: Menunggu, Dalam Perjalanan, dan Selesai
3. THE Task_Module SHALL menampilkan daftar tugas dalam Card_Component dengan layout yang menampilkan nama sekolah, alamat, jenis tugas (pengiriman/pengambilan), status, dan urutan rute
4. THE Task_Module SHALL menampilkan tag jenis tugas dengan warna berbeda: hijau untuk pengiriman makanan dan oranye untuk pengambilan ompreng
5. WHEN driver menekan kartu tugas pengiriman, THE Task_Module SHALL menavigasi ke halaman detail tugas pengiriman
6. WHEN driver menekan kartu tugas pengambilan, THE Task_Module SHALL menavigasi ke halaman detail tugas pengambilan
7. THE Task_Module SHALL mendukung Pull_To_Refresh untuk memperbarui daftar tugas
8. WHILE daftar tugas sedang dimuat, THE Task_Module SHALL menampilkan Skeleton_Loading
9. THE Task_Module SHALL mempertahankan seluruh business logic tugas yang sudah ada termasuk offline caching, sync, dan e-POD submission
10. THE Task_Module SHALL menampilkan tab selection dengan dua tab: "Pengiriman" dan "Pengambilan" di bagian atas daftar tugas, dengan style pill-shaped (border-radius 20px) mengikuti referensi HR Attendee app
11. THE Task_Module SHALL menampilkan tab aktif dengan background Color_Palette primary (#5A4372) dan teks putih, serta tab tidak aktif dengan background transparan dan teks Color_Palette text-secondary (#74788C)
12. WHEN driver memilih tab "Pengiriman", THE Task_Module SHALL memfilter daftar tugas untuk hanya menampilkan tugas pengiriman makanan
13. WHEN driver memilih tab "Pengambilan", THE Task_Module SHALL memfilter daftar tugas untuk hanya menampilkan tugas pengambilan ompreng
14. WHEN driver berpindah tab, THE Task_Module SHALL memperbarui Summary_Card agar menampilkan jumlah status (Menunggu, Dalam Perjalanan, Selesai) hanya dari tugas yang sesuai dengan tab aktif

### Requirement 5: Redesign Halaman Detail Tugas dan e-POD

**User Story:** Sebagai driver, saya ingin halaman detail tugas dan form e-POD yang lebih rapi dan mudah digunakan, sehingga proses pengiriman dan dokumentasi lebih lancar.

#### Acceptance Criteria

1. THE Task_Module SHALL menampilkan halaman detail tugas dengan Card_Component terpisah untuk: informasi status, informasi sekolah, informasi menu/porsi, dan aksi
2. THE Task_Module SHALL menampilkan progress indicator visual yang menunjukkan tahapan tugas saat ini
3. THE Task_Module SHALL menampilkan tombol navigasi Google Maps dengan style Card_Component
4. THE Task_Module SHALL menampilkan form e-POD dengan layout Card_Component yang memisahkan setiap bagian: lokasi GPS, foto bukti, tanda tangan, dan nama penerima
5. WHILE form e-POD sedang disubmit, THE Task_Module SHALL menampilkan progress indicator dan menonaktifkan tombol submit
6. THE Task_Module SHALL mempertahankan seluruh business logic detail tugas dan e-POD yang sudah ada termasuk upload foto, tanda tangan digital, dan offline sync

### Requirement 6: Redesign Halaman Profil

**User Story:** Sebagai karyawan, saya ingin halaman profil yang modern dengan tab-based layout mengikuti style HR Attendee app, sehingga informasi akun saya tertata rapi.

#### Acceptance Criteria

1. THE PWA SHALL menampilkan halaman profil dengan banner gradient Color_Palette di bagian atas dan avatar karyawan yang overlapping
2. THE PWA SHALL menampilkan nama lengkap dan NIK karyawan di bawah avatar
3. THE PWA SHALL menampilkan informasi akun dalam Card_Component dengan tab Personal (NIK, nama, email, role) dan tab Pengaturan (versi aplikasi, mode offline)
4. THE PWA SHALL menampilkan tombol logout full-width dengan konfirmasi dialog
5. THE PWA SHALL mempertahankan seluruh business logic profil dan logout yang sudah ada

### Requirement 7: Modul Dashboard (View Only)

**User Story:** Sebagai Kepala SPPG/Yayasan, saya ingin melihat ringkasan data operasional dalam dashboard yang minimalis di PWA, sehingga saya dapat memantau kinerja operasional dari perangkat mobile.

#### Acceptance Criteria

1. THE Dashboard_Module SHALL hanya dapat diakses oleh karyawan dengan role Kepala_SPPG
2. THE Dashboard_Module SHALL menampilkan Summary_Card untuk metrik utama: total karyawan hadir hari ini, total tugas pengiriman hari ini, total tugas selesai, dan total sekolah yang dilayani
3. THE Dashboard_Module SHALL menampilkan grafik ringkasan kehadiran 7 hari terakhir dalam Card_Component
4. THE Dashboard_Module SHALL menampilkan daftar tugas pengiriman terbaru dengan status dalam Card_Component
5. THE Dashboard_Module SHALL mengambil data dari API backend yang sudah ada tanpa modifikasi endpoint
6. WHILE data dashboard sedang dimuat, THE Dashboard_Module SHALL menampilkan Skeleton_Loading pada setiap Card_Component
7. THE Dashboard_Module SHALL mendukung Pull_To_Refresh untuk memperbarui seluruh data dashboard
8. IF pengambilan data dashboard gagal, THEN THE Dashboard_Module SHALL menampilkan pesan error dan tombol retry
9. THE Dashboard_Module SHALL bersifat read-only tanpa aksi modifikasi data

### Requirement 8: Modul Monitoring Aktivitas (View Only)

**User Story:** Sebagai Kepala SPPG/Yayasan, saya ingin melihat monitoring aktivitas operasional dalam tampilan minimalis di PWA, sehingga saya dapat memantau aktivitas harian dari perangkat mobile.

#### Acceptance Criteria

1. THE Activity_Monitoring_Module SHALL hanya dapat diakses oleh karyawan dengan role Kepala_SPPG
2. THE Activity_Monitoring_Module SHALL menampilkan Date_Selector untuk memfilter aktivitas berdasarkan tanggal
3. THE Activity_Monitoring_Module SHALL menampilkan daftar aktivitas dalam Card_Component dengan informasi: nama karyawan, jenis aktivitas, waktu, dan status
4. THE Activity_Monitoring_Module SHALL menampilkan filter berdasarkan jenis aktivitas (absensi, pengiriman, pengambilan) menggunakan tab atau chip selector
5. THE Activity_Monitoring_Module SHALL mengambil data dari API backend yang sudah ada tanpa modifikasi endpoint
6. WHILE data monitoring sedang dimuat, THE Activity_Monitoring_Module SHALL menampilkan Skeleton_Loading
7. THE Activity_Monitoring_Module SHALL mendukung Pull_To_Refresh untuk memperbarui data monitoring
8. IF pengambilan data monitoring gagal, THEN THE Activity_Monitoring_Module SHALL menampilkan pesan error dan tombol retry
9. THE Activity_Monitoring_Module SHALL bersifat read-only tanpa aksi modifikasi data
10. THE Activity_Monitoring_Module SHALL mendukung infinite scroll untuk memuat data aktivitas yang lebih banyak

### Requirement 9: Modul Perencanaan Menu (View Only + Approval)

**User Story:** Sebagai Kepala SPPG/Yayasan atau Ahli Gizi, saya ingin melihat perencanaan menu dalam tampilan minimalis di PWA, dan sebagai Kepala SPPG saya ingin dapat melakukan approval menu mingguan, sehingga proses perencanaan menu dapat dipantau dan disetujui dari perangkat mobile.

#### Acceptance Criteria

1. THE Menu_Planning_Module SHALL dapat diakses oleh karyawan dengan role Kepala_SPPG dan Ahli_Gizi
2. THE Menu_Planning_Module SHALL menampilkan daftar rencana menu mingguan dalam Card_Component dengan informasi: periode minggu, status approval, dan jumlah menu
3. THE Menu_Planning_Module SHALL menampilkan detail menu harian dalam Card_Component yang menampilkan nama menu, komponen menu, dan jumlah porsi per hari
4. THE Menu_Planning_Module SHALL menampilkan filter berdasarkan minggu menggunakan Date_Selector
5. WHEN karyawan dengan role Kepala_SPPG melihat menu mingguan yang belum diapprove, THE Menu_Planning_Module SHALL menampilkan tombol Approve dan Reject
6. WHEN Kepala_SPPG menekan tombol Approve pada menu mingguan, THE Menu_Planning_Module SHALL mengirim request approval ke API backend dan memperbarui status menu menjadi approved
7. WHEN Kepala_SPPG menekan tombol Reject pada menu mingguan, THE Menu_Planning_Module SHALL menampilkan dialog input alasan penolakan dan mengirim request rejection ke API backend
8. WHILE karyawan dengan role Ahli_Gizi mengakses Menu_Planning_Module, THE Menu_Planning_Module SHALL menampilkan data menu dalam mode read-only tanpa tombol approval
9. THE Menu_Planning_Module SHALL mengambil data dari API backend yang sudah ada tanpa modifikasi endpoint
10. WHILE data menu sedang dimuat, THE Menu_Planning_Module SHALL menampilkan Skeleton_Loading
11. THE Menu_Planning_Module SHALL mendukung Pull_To_Refresh untuk memperbarui data menu
12. IF pengambilan data menu gagal, THEN THE Menu_Planning_Module SHALL menampilkan pesan error dan tombol retry
13. IF proses approval atau rejection gagal, THEN THE Menu_Planning_Module SHALL menampilkan pesan error dan mempertahankan status menu sebelumnya

### Requirement 10: Design System dan Komponen UI

**User Story:** Sebagai developer, saya ingin design system yang konsisten mengikuti style HR Attendee app dengan color palette yang sudah ada, sehingga seluruh halaman PWA memiliki tampilan yang seragam dan modern.

#### Acceptance Criteria

1. THE PWA SHALL menggunakan CSS variables untuk seluruh Color_Palette: primary #5A4372, accent #3D2B53, background #F8FDEA, text-primary #322837, text-secondary #74788C, success #05CD99, warning #FFB547, error #EE5D50
2. THE PWA SHALL menggunakan font family 'DM Sans' sebagai font utama dengan fallback ke system fonts
3. THE PWA SHALL menggunakan border-radius 16px untuk Card_Component, 12px untuk tombol, dan 8px untuk elemen kecil
4. THE PWA SHALL menggunakan shadow card (0px 18px 40px rgba(112, 144, 176, 0.12)) untuk seluruh Card_Component
5. THE PWA SHALL menggunakan minimum touch target 44x44px untuk seluruh elemen interaktif
6. THE PWA SHALL menampilkan Skeleton_Loading sebagai placeholder saat data sedang dimuat pada setiap halaman
7. THE PWA SHALL menggunakan transisi halus (200ms ease-in-out) untuk perubahan state pada komponen interaktif
8. THE PWA SHALL memastikan kontras warna teks terhadap background memenuhi rasio minimal 4.5:1 untuk teks normal

### Requirement 11: Fitur Pendukung View Only pada Modul Absensi

**User Story:** Sebagai karyawan, saya ingin melihat statistik kehadiran dan ringkasan absensi bulanan di modul absensi, sehingga saya dapat memantau rekap kehadiran saya sendiri.

#### Acceptance Criteria

1. THE Attendance_Module SHALL menampilkan Summary_Card statistik kehadiran bulan berjalan: total hari hadir, total hari tidak hadir, total terlambat, dan rata-rata jam kerja
2. THE Attendance_Module SHALL menampilkan kalender mini yang menandai hari hadir (hijau), tidak hadir (merah), dan terlambat (kuning) pada bulan berjalan
3. WHEN karyawan memilih tanggal pada kalender mini, THE Attendance_Module SHALL menampilkan detail absensi pada tanggal tersebut
4. THE Attendance_Module SHALL mengambil data statistik dari API backend yang sudah ada

### Requirement 12: Fitur Pendukung View Only pada Modul Tugas

**User Story:** Sebagai driver atau asisten lapangan, saya ingin melihat ringkasan performa tugas dan riwayat tugas sebelumnya, sehingga saya dapat memantau kinerja pengiriman saya.

#### Acceptance Criteria

1. THE Task_Module SHALL menampilkan Summary_Card performa tugas: total tugas hari ini, tugas selesai, tugas pending, dan persentase penyelesaian
2. THE Task_Module SHALL menampilkan tab Riwayat yang menampilkan daftar tugas 7 hari terakhir dalam Card_Component
3. WHEN driver memilih tanggal pada tab Riwayat, THE Task_Module SHALL menampilkan daftar tugas pada tanggal tersebut
4. THE Task_Module SHALL mengambil data riwayat dari API backend yang sudah ada

### Requirement 13: Routing dan Akses Kontrol Modul Baru

**User Story:** Sebagai developer, saya ingin routing dan akses kontrol yang tepat untuk modul-modul baru, sehingga setiap role hanya dapat mengakses modul yang sesuai.

#### Acceptance Criteria

1. THE PWA SHALL mendaftarkan route /dashboard untuk Dashboard_Module dengan akses terbatas pada role Kepala_SPPG
2. THE PWA SHALL mendaftarkan route /monitoring untuk Activity_Monitoring_Module dengan akses terbatas pada role Kepala_SPPG
3. THE PWA SHALL mendaftarkan route /menu-planning untuk Menu_Planning_Module dengan akses terbatas pada role Kepala_SPPG dan Ahli_Gizi
4. THE PWA SHALL mendaftarkan route /school-monitoring untuk School_Monitoring_Module dengan akses terbatas pada role Sekolah
5. WHEN karyawan dengan role tidak sesuai mengakses route modul baru, THE PWA SHALL mengarahkan karyawan ke halaman default sesuai role
6. WHEN karyawan dengan role Kepala_SPPG login, THE PWA SHALL mengarahkan ke halaman /dashboard sebagai halaman default
7. WHEN karyawan dengan role Ahli_Gizi login, THE PWA SHALL mengarahkan ke halaman /menu-planning sebagai halaman default
8. WHEN pengguna dengan role Sekolah login, THE PWA SHALL mengarahkan ke halaman /school-monitoring sebagai halaman default
9. THE PWA SHALL mempertahankan seluruh routing dan akses kontrol yang sudah ada untuk modul Absensi dan Tugas


### Requirement 14: Modul Monitoring Sekolah (View Only)

**User Story:** Sebagai pihak sekolah, saya ingin memantau menu yang dikirim hari ini dan status/progres pengiriman ke sekolah saya melalui PWA, sehingga saya mengetahui apa yang akan diterima dan kapan pengiriman tiba.

#### Acceptance Criteria

1. THE School_Monitoring_Module SHALL hanya dapat diakses oleh pengguna dengan role Sekolah
2. THE School_Monitoring_Module SHALL menampilkan data hanya untuk sekolah yang terhubung dengan akun Sekolah yang sedang login (berdasarkan schoolId)
3. THE School_Monitoring_Module SHALL menampilkan informasi menu hari ini untuk sekolah tersebut dalam Card_Component, meliputi: nama menu, komponen menu, dan jumlah porsi
4. THE School_Monitoring_Module SHALL menampilkan status pengiriman terkini ke sekolah tersebut dalam Card_Component, meliputi: status (menunggu/dalam perjalanan/sampai/selesai), nama driver, dan estimasi waktu jika tersedia
5. THE School_Monitoring_Module SHALL menampilkan status pengambilan ompreng jika ada tugas pengambilan terjadwal untuk sekolah tersebut
6. THE School_Monitoring_Module SHALL bersifat read-only tanpa aksi modifikasi data apapun
7. THE School_Monitoring_Module SHALL mengambil data dari API backend yang sudah ada tanpa modifikasi endpoint
8. WHILE data monitoring sekolah sedang dimuat, THE School_Monitoring_Module SHALL menampilkan Skeleton_Loading pada setiap Card_Component
9. THE School_Monitoring_Module SHALL mendukung Pull_To_Refresh untuk memperbarui seluruh data monitoring sekolah
10. IF pengambilan data monitoring sekolah gagal, THEN THE School_Monitoring_Module SHALL menampilkan pesan error dan tombol retry
