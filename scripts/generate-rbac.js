const XLSX = require('../web/node_modules/xlsx');
const wb = XLSX.utils.book_new();

// Sheet 1: Role Overview
const roles = [
  ['Role', 'Deskripsi', 'Platform', 'Default Landing Page', 'Jumlah Fitur'],
  ['superadmin', 'Administrator sistem tertinggi, kelola semua organisasi', 'Web', '/yayasan', '15+'],
  ['admin_bgn', 'Admin Badan Gizi Nasional, monitoring nasional', 'Web + PWA', '/dashboard-bgn', '8'],
  ['kepala_yayasan', 'Kepala Yayasan, monitoring multi-SPPG + audit', 'Web + PWA', '/dashboard-yayasan', '12'],
  ['kepala_sppg', 'Kepala SPPG, kelola operasional dapur harian', 'Web + PWA', '/dashboard', '25+'],
  ['ahli_gizi', 'Ahli Gizi, kelola resep dan menu', 'Web + PWA', '/menu-planning', '8'],
  ['pengadaan', 'Staff Pengadaan, kelola supplier dan PO', 'Web', '/purchase-orders', '7'],
  ['akuntan', 'Akuntan, kelola keuangan dan laporan', 'Web', '/financial-reports', '10'],
  ['chef', 'Chef/Juru Masak, KDS cooking', 'Web', '/kds/cooking', '4'],
  ['packing', 'Staff Packing, KDS packing', 'Web', '/kds/cooking', '3'],
  ['driver', 'Driver pengiriman makanan', 'PWA', '/tasks', '6'],
  ['asisten_lapangan', 'Asisten Lapangan, bantu pengiriman', 'PWA', '/tasks', '6'],
  ['kebersihan', 'Staff Kebersihan, pencucian ompreng', 'Web', '/kds/cleaning', '3'],
  ['sekolah', 'Pihak Sekolah, monitoring penerimaan', 'PWA', '/school-monitoring', '3'],
];
const ws1 = XLSX.utils.aoa_to_sheet(roles);
ws1['!cols'] = [{wch:18},{wch:52},{wch:12},{wch:25},{wch:14}];
XLSX.utils.book_append_sheet(wb, ws1, 'Role Overview');

// Sheet 2: RBAC Matrix
const roleNames = ['superadmin','admin_bgn','kepala_yayasan','kepala_sppg','ahli_gizi','pengadaan','akuntan','chef','packing','driver','asisten_lapangan','kebersihan','sekolah'];
const features = [
  ['Dashboard Kepala SPPG',       '-','-','-','V','-','-','-','-','-','-','-','-','-'],
  ['Dashboard Kepala Yayasan',    '-','-','V','-','-','-','-','-','-','-','-','-','-'],
  ['Dashboard BGN',               'V','V','-','-','-','-','-','-','-','-','-','-','-'],
  ['Manajemen Yayasan',           'V','V','-','-','-','-','-','-','-','-','-','-','-'],
  ['Manajemen SPPG',              'V','V','-','-','-','-','-','-','-','-','-','-','-'],
  ['Manajemen User',              'V','-','V','V','-','-','-','-','-','-','-','-','-'],
  ['Manajemen Resep',             '-','-','-','V','V','-','-','-','-','-','-','-','-'],
  ['Barang Setengah Jadi',        '-','-','-','V','V','-','-','V','-','-','-','-','-'],
  ['Perencanaan Menu',            '-','-','-','V','V','-','-','-','-','-','-','-','-'],
  ['KDS Cooking',                 '-','-','-','V','V','-','-','V','-','-','-','-','-'],
  ['KDS Packing',                 '-','-','-','V','V','-','-','V','V','-','-','-','-'],
  ['KDS Pencucian Ompreng',       '-','-','-','V','-','-','-','-','-','-','-','V','-'],
  ['Manajemen Supplier',          '-','-','-','V','-','V','-','-','-','-','-','-','-'],
  ['Purchase Order',              '-','-','-','V','-','V','-','-','-','-','-','-','-'],
  ['Penerimaan Barang (GRN)',     '-','-','-','V','-','V','-','-','-','-','-','-','-'],
  ['Inventori Bahan Baku',        '-','-','-','V','-','V','V','-','-','-','-','-','-'],
  ['Stok Opname',                 '-','-','-','V (approve)','-','V (input)','-','-','-','-','-','-','-'],
  ['Manajemen Sekolah',           '-','-','-','V','-','-','-','-','-','V (view)','V (view)','-','-'],
  ['Tugas Pengiriman',            '-','-','-','V (assign)','-','-','-','-','-','V (exec)','V (exec)','-','-'],
  ['Tugas Pengambilan (Pickup)',   '-','-','V (view)','V (assign)','-','-','-','-','-','V (exec)','V (exec)','-','-'],
  ['e-POD',                       '-','-','-','-','-','-','-','-','-','V','V','-','-'],
  ['Pelacakan Ompreng',           '-','-','-','V','-','-','-','-','-','V','V','-','-'],
  ['Monitoring Pengiriman',       '-','-','V','V','V','V','V','V','V','V','V','-','-'],
  ['Monitoring Aktivitas',        '-','-','V','V','V','V','V','V','V','V','V','-','-'],
  ['Manajemen Karyawan',          '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Laporan Absensi',             '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Konfigurasi Absensi',         '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Aset Dapur',                  '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Arus Kas',                    '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Laporan Keuangan',            '-','-','-','V','-','-','V','-','-','-','-','-','-'],
  ['Ulasan & Rating',             '-','-','V','V','-','-','V','-','-','-','-','-','-'],
  ['Risk Assessment (Audit SOP)', 'V (template)','-','V (audit)','-','-','-','-','-','-','-','-','-','-'],
  ['Audit Trail',                 'V','V','V','V','-','-','-','-','-','-','-','-','-'],
  ['Konfigurasi Sistem',          'V','-','-','V','-','-','-','-','-','-','-','-','-'],
  ['Absensi (Check-in/out) PWA',  '-','-','-','-','-','-','-','-','-','V','V','-','-'],
  ['Review Sekolah (PWA)',        '-','-','-','-','-','-','-','-','-','V','V','-','V'],
  ['School Monitoring (PWA)',     '-','-','-','-','-','-','-','-','-','-','-','-','V'],
  ['Profil',                      'V','V','V','V','V','V','V','V','V','V','V','V','V'],
  ['Notifikasi',                  'V','V','V','V','V','V','V','V','V','V','V','V','V'],
];

const header2 = ['Fitur / Modul', ...roleNames];
const ws2 = XLSX.utils.aoa_to_sheet([header2, ...features]);
ws2['!cols'] = [{wch:32}, ...roleNames.map(() => ({wch:15}))];
XLSX.utils.book_append_sheet(wb, ws2, 'RBAC Matrix');

// Sheet 3: Platform Access
const platformData = [
  ['Role', 'Web Dashboard', 'PWA Mobile', 'Keterangan'],
  ['superadmin', 'Full Access', 'Dashboard BGN', 'Akses semua modul web + dashboard BGN di PWA'],
  ['admin_bgn', 'Dashboard + Org', 'Dashboard BGN', 'Dashboard BGN, manajemen yayasan/SPPG'],
  ['kepala_yayasan', 'Dashboard + Audit', 'Dashboard + Audit', 'Dashboard yayasan, risk assessment, review di web & PWA'],
  ['kepala_sppg', 'Full Operasional', 'Dashboard + Monitoring', 'Semua modul operasional di web, dashboard + monitoring di PWA'],
  ['ahli_gizi', 'Menu & KDS', 'Menu Planning', 'Resep, menu planning, KDS di web; menu planning di PWA'],
  ['pengadaan', 'Supply Chain', 'Tidak ada akses', 'Supplier, PO, GRN, inventori - hanya web'],
  ['akuntan', 'Keuangan + SDM', 'Tidak ada akses', 'Keuangan, karyawan, absensi - hanya web'],
  ['chef', 'KDS Cooking', 'Tidak ada akses', 'KDS cooking & packing - hanya web'],
  ['packing', 'KDS Packing', 'Tidak ada akses', 'KDS packing - hanya web'],
  ['driver', 'Tidak ada akses', 'Tugas + Absensi', 'Tugas pengiriman, e-POD, absensi - hanya PWA'],
  ['asisten_lapangan', 'Tidak ada akses', 'Tugas + Absensi', 'Tugas pengiriman, e-POD, absensi - hanya PWA'],
  ['kebersihan', 'KDS Cleaning', 'Tidak ada akses', 'Pencucian ompreng - hanya web'],
  ['sekolah', 'Tidak ada akses', 'School Monitoring', 'Monitoring penerimaan makanan - hanya PWA'],
];
const ws3 = XLSX.utils.aoa_to_sheet(platformData);
ws3['!cols'] = [{wch:18},{wch:22},{wch:22},{wch:58}];
XLSX.utils.book_append_sheet(wb, ws3, 'Platform Access');

// Sheet 4: Role Relationships
const relData = [
  ['Role', 'Hubungan dengan Role Lain', 'Hierarki', 'Alur Kerja'],
  ['superadmin', 'Membuat & mengelola semua user termasuk admin_bgn dan kepala_yayasan. Mengelola SOP template untuk risk assessment.', 'Level 1 (Tertinggi)', 'Superadmin -> Buat Yayasan -> Assign Kepala Yayasan -> Kepala Yayasan buat SPPG'],
  ['admin_bgn', 'Setara superadmin untuk scope BGN. Melihat dashboard agregat semua yayasan & SPPG.', 'Level 1', 'Admin BGN -> Monitor semua Yayasan -> Drill-down ke SPPG'],
  ['kepala_yayasan', 'Mengawasi semua SPPG di bawah yayasannya. Melakukan risk assessment/audit ke SPPG.', 'Level 2', 'Kepala Yayasan -> Audit SPPG -> Beri skor -> Submit laporan'],
  ['kepala_sppg', 'Mengelola semua operasional 1 SPPG. Approve stok opname & menu plan.', 'Level 3', 'Kepala SPPG -> Approve menu -> Monitor produksi -> Assign driver -> Review hasil'],
  ['ahli_gizi', 'Membuat resep & menu plan. Data digunakan chef untuk KDS.', 'Level 4', 'Ahli Gizi -> Buat resep -> Buat menu plan -> Kepala SPPG approve -> Chef eksekusi'],
  ['pengadaan', 'Membuat PO. Menerima barang (GRN) yang update inventori.', 'Level 4', 'Pengadaan -> Buat PO -> Kepala SPPG approve -> Terima barang (GRN) -> Update stok'],
  ['akuntan', 'Mengelola keuangan & SDM. Melihat laporan dari data operasional.', 'Level 4', 'Akuntan -> Input arus kas -> Kelola karyawan -> Generate laporan keuangan'],
  ['chef', 'Mengeksekusi menu via KDS cooking. Status trigger KDS packing.', 'Level 5', 'Chef -> Lihat menu hari ini di KDS -> Masak -> Update status selesai -> Trigger packing'],
  ['packing', 'Mengemas makanan setelah chef selesai. Status trigger tugas pengiriman.', 'Level 5', 'Packing -> Lihat order di KDS -> Kemas per sekolah -> Update status -> Trigger delivery'],
  ['driver', 'Menerima tugas pengiriman. Mengisi e-POD saat serah terima.', 'Level 5', 'Driver -> Terima tugas -> Kirim ke sekolah -> Isi e-POD -> Sekolah review'],
  ['asisten_lapangan', 'Membantu driver. Bisa mengisi e-POD dan pickup task.', 'Level 5', 'Asisten -> Bantu driver -> Isi e-POD -> Pickup ompreng dari sekolah'],
  ['kebersihan', 'Mencuci ompreng setelah dikembalikan dari sekolah.', 'Level 5', 'Kebersihan -> Lihat ompreng pending -> Mulai cuci -> Selesai cuci -> Update status'],
  ['sekolah', 'Pihak eksternal penerima makanan. Memberikan review/rating.', 'Eksternal', 'Sekolah -> Monitor status pengiriman -> Terima makanan -> Beri review & rating'],
];
const ws4 = XLSX.utils.aoa_to_sheet(relData);
ws4['!cols'] = [{wch:18},{wch:65},{wch:18},{wch:65}];
XLSX.utils.book_append_sheet(wb, ws4, 'Role Relationships');

// Sheet 5: API Endpoint Access
const apiData = [
  ['Endpoint Group', 'Method', 'Roles yang Diizinkan', 'Keterangan'],
  ['/auth/login', 'POST', 'Public', 'Login semua role'],
  ['/auth/me', 'GET', 'Semua (Authenticated)', 'Info user yang sedang login'],
  ['/recipes', 'GET/POST/PUT/DELETE', 'Semua (Authenticated)', 'CRUD resep - frontend restrict ke ahli_gizi, kepala_sppg'],
  ['/ingredients', 'GET/POST', 'Semua (Authenticated)', 'CRUD bahan baku'],
  ['/semi-finished', 'GET/POST/PUT/DELETE', 'Semua (Authenticated)', 'Barang setengah jadi'],
  ['/menu-plans', 'GET/POST/PUT', 'Semua (Authenticated)', 'Menu planning - frontend restrict'],
  ['/menu-plans/:id/approve', 'POST', 'Semua (Authenticated)', 'Approve menu - frontend restrict ke kepala_sppg'],
  ['/kds/cooking', 'GET/PUT/POST', 'Semua (Authenticated)', 'KDS cooking - frontend restrict ke chef'],
  ['/kds/packing', 'GET/PUT/POST', 'Semua (Authenticated)', 'KDS packing - frontend restrict ke packing'],
  ['/suppliers', 'GET/POST/PUT', 'Semua (Authenticated)', 'Supplier - frontend restrict ke pengadaan'],
  ['/purchase-orders', 'GET/POST/PUT', 'Semua (Authenticated)', 'PO - frontend restrict ke pengadaan'],
  ['/purchase-orders/:id/approve', 'POST', 'Semua (Authenticated)', 'Approve PO - frontend restrict ke kepala_sppg'],
  ['/goods-receipts', 'GET/POST', 'Semua (Authenticated)', 'GRN - frontend restrict ke pengadaan'],
  ['/inventory', 'GET/POST', 'Semua (Authenticated)', 'Inventori bahan baku'],
  ['/stok-opname/forms/:id/approve', 'POST', 'kepala_sppg (backend enforced)', 'Hanya kepala SPPG bisa approve stok opname'],
  ['/stok-opname/forms/:id/reject', 'POST', 'kepala_sppg (backend enforced)', 'Hanya kepala SPPG bisa reject stok opname'],
  ['/schools', 'GET/POST/PUT/DELETE', 'Semua (Authenticated)', 'CRUD sekolah'],
  ['/delivery-tasks', 'GET/POST/PUT/DELETE', 'Semua (Authenticated)', 'Tugas pengiriman'],
  ['/pickup-tasks', 'GET/POST/PUT/DELETE', 'kepala_sppg, kepala_yayasan, asisten_lapangan, driver', 'Backend enforced role check'],
  ['/epod', 'GET/POST', 'Semua (Authenticated)', 'e-POD - frontend restrict ke driver, asisten'],
  ['/reviews', 'GET/POST', 'Semua (Authenticated)', 'Ulasan & rating'],
  ['/employees', 'GET/POST/PUT', 'Semua (Authenticated)', 'Karyawan - frontend restrict'],
  ['/attendance', 'GET/POST', 'Semua (Authenticated)', 'Absensi check-in/out'],
  ['/assets', 'GET/POST/PUT/DELETE', 'Semua (Authenticated)', 'Aset dapur - frontend restrict ke akuntan'],
  ['/cash-flow', 'GET/POST', 'Semua (Authenticated)', 'Arus kas - frontend restrict ke akuntan'],
  ['/financial-reports', 'GET/POST', 'Semua (Authenticated)', 'Laporan keuangan'],
  ['/dashboard/kepala-sppg', 'GET', 'Semua (Authenticated)', 'Dashboard SPPG'],
  ['/dashboard/kepala-yayasan', 'GET', 'Permission: dashboard_yayasan', 'Dashboard agregat yayasan'],
  ['/dashboard/admin-bgn', 'GET', 'Permission: dashboard_bgn', 'Dashboard agregat BGN'],
  ['/organizations/yayasan', 'GET/POST/PUT/PATCH', 'superadmin, admin_bgn', 'Backend enforced - CRUD yayasan'],
  ['/organizations/sppg', 'GET/POST/PUT/PATCH', 'superadmin, admin_bgn', 'Backend enforced - CRUD SPPG'],
  ['/users', 'GET/POST/PUT/PATCH', 'Permission: user_provisioning', 'Backend enforced - CRUD user'],
  ['/monitoring', 'GET/PUT', 'Semua kecuali kebersihan', 'Backend enforced - monitoring pengiriman'],
  ['/cleaning', 'GET/POST', 'kebersihan, kepala_sppg, kepala_yayasan', 'Backend enforced - pencucian ompreng'],
  ['/activity-tracker', 'GET/PUT/POST', 'kepala_sppg, kepala_yayasan, akuntan', 'Backend enforced - pelacakan aktivitas'],
  ['/risk-assessment/sop-categories (write)', 'POST/PUT', 'superadmin', 'Backend enforced - kelola template SOP'],
  ['/risk-assessment/sop-checklist-items (write)', 'POST/PUT/PATCH', 'superadmin', 'Backend enforced - kelola checklist SOP'],
  ['/risk-assessment/forms', 'GET/POST/PUT/DELETE', 'kepala_yayasan, superadmin', 'Backend enforced - form audit'],
  ['/risk-assessment/stats', 'GET', 'kepala_yayasan, superadmin', 'Backend enforced - statistik audit'],
  ['/audit-trail', 'GET', 'Semua (Authenticated)', 'Log audit - frontend restrict'],
  ['/system-config', 'GET/POST/DELETE', 'Semua + IP Whitelist', 'Konfigurasi sistem'],
  ['/notifications', 'GET/PUT/DELETE', 'Semua (Authenticated)', 'Notifikasi'],
];
const ws5 = XLSX.utils.aoa_to_sheet(apiData);
ws5['!cols'] = [{wch:40},{wch:22},{wch:42},{wch:52}];
XLSX.utils.book_append_sheet(wb, ws5, 'API Endpoint Access');

XLSX.writeFile(wb, 'docs/RBAC-Dapur-Sehat.xlsx');
console.log('Created: docs/RBAC-Dapur-Sehat.xlsx');
