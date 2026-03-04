# Panduan Laporan Absensi

## Masalah: Data Absensi Tidak Muncul di Laporan

### Penyebab Umum

1. **Filter Tanggal Tidak Sesuai**
   - Pastikan periode tanggal mencakup tanggal check-in Anda
   - Format tanggal: DD/MM/YYYY (contoh: 01/03/2026 - 31/03/2026)

2. **Data Belum Ter-refresh**
   - Setelah check-in/check-out, refresh halaman laporan
   - Klik tombol "Cari" untuk reload data

3. **Token Expired**
   - Jika sudah login lama, token mungkin expired
   - Logout dan login kembali

### Cara Menggunakan Laporan Absensi

#### 1. Akses Halaman Laporan
- Login ke web admin: http://localhost:5173
- Menu: **SDM** → **Laporan Absensi**

#### 2. Filter Data
- **Periode**: Pilih tanggal mulai dan tanggal akhir
- **Karyawan**: (Opsional) Pilih karyawan tertentu atau kosongkan untuk semua karyawan
- Klik tombol **Cari** untuk menampilkan data

#### 3. Melihat Data
Tabel akan menampilkan:
- Nama Karyawan
- Posisi
- Total Hari (berapa hari masuk kerja)
- Total Jam (total jam kerja)
- Rata-rata Jam/Hari
- Tingkat Kehadiran (persentase)

#### 4. Detail Absensi
- Klik **nama karyawan** di tabel untuk melihat detail harian
- Modal akan muncul dengan data:
  - Tanggal
  - Jam Check In
  - Jam Check Out
  - Total Jam Kerja
  - Status (Lengkap/Cukup/Kurang/Belum Check Out)

#### 5. Export Data
- **Export Excel**: Klik tombol "Export Excel" untuk download file .xlsx
- **Export PDF**: Klik tombol "Export PDF" untuk download file .pdf

### Troubleshooting

#### Data Tidak Muncul Setelah Check-in/Check-out

**Solusi 1: Refresh Halaman**
```
1. Buka halaman Laporan Absensi
2. Tekan F5 atau Ctrl+R (Cmd+R di Mac)
3. Klik tombol "Cari" lagi
```

**Solusi 2: Cek Filter Tanggal**
```
1. Pastikan periode mencakup hari ini
2. Contoh: Jika hari ini 04/03/2026, pastikan filter adalah:
   - Dari: 01/03/2026
   - Sampai: 31/03/2026 (atau lebih)
```

**Solusi 3: Cek Data di Database**
```bash
# Cek data absensi hari ini
psql -U arifyudhistira -d erp_sppg -c "
  SELECT 
    e.full_name, 
    a.date, 
    a.check_in, 
    a.check_out, 
    a.work_hours 
  FROM attendances a 
  JOIN employees e ON e.id = a.employee_id 
  WHERE DATE(a.date) = CURRENT_DATE 
  ORDER BY a.check_in DESC;
"
```

**Solusi 4: Cek Browser Console**
```
1. Buka Developer Tools (F12)
2. Tab "Console"
3. Cari error atau warning
4. Cari log "[AttendanceReport]" untuk debug info
```

**Solusi 5: Cek Network Request**
```
1. Buka Developer Tools (F12)
2. Tab "Network"
3. Klik tombol "Cari" di halaman laporan
4. Cari request ke "/api/v1/attendance/report"
5. Cek:
   - Status: Harus 200 OK
   - Response: Harus ada data
   - Request URL: Cek parameter start_date dan end_date
```

#### Error "Tidak ada data absensi untuk periode yang dipilih"

**Penyebab:**
- Tidak ada data check-in/check-out dalam periode tersebut
- Filter karyawan tidak match
- Timezone issue (jarang terjadi)

**Solusi:**
```
1. Reset filter: Klik tombol "Reset"
2. Pilih periode yang lebih luas (misalnya 1 bulan penuh)
3. Kosongkan filter karyawan
4. Klik "Cari" lagi
```

#### Data Muncul di PWA tapi Tidak di Web

**Penyebab:**
- Cache browser
- Token berbeda
- Build web belum ter-update

**Solusi:**
```bash
# 1. Rebuild web
cd web
npm run build

# 2. Clear browser cache
# - Chrome: Ctrl+Shift+Delete → Clear cache
# - Firefox: Ctrl+Shift+Delete → Clear cache

# 3. Hard refresh
# - Chrome/Firefox: Ctrl+Shift+R (Cmd+Shift+R di Mac)
```

### Fitur Laporan

#### 1. Summary Statistics
Di bagian atas tabel, ada 4 kartu statistik:
- **Total Karyawan**: Jumlah karyawan yang ada di laporan
- **Total Hari Kerja**: Total hari kerja semua karyawan
- **Total Jam Kerja**: Total jam kerja semua karyawan
- **Rata-rata Jam/Hari**: Rata-rata jam kerja per hari

#### 2. Tingkat Kehadiran
Progress bar menunjukkan persentase kehadiran:
- **Hijau (≥90%)**: Sangat baik
- **Biru (75-89%)**: Baik
- **Merah (<75%)**: Perlu perhatian

Perhitungan:
```
Tingkat Kehadiran = (Total Hari Masuk / Total Hari Kerja) × 100%
```

Catatan: Hari kerja tidak termasuk Sabtu dan Minggu

#### 3. Status Absensi Harian
Dalam detail absensi, status ditentukan berdasarkan jam kerja:
- **Lengkap** (Hijau): ≥8 jam
- **Cukup** (Biru): 6-8 jam
- **Kurang** (Merah): <6 jam
- **Belum Check Out** (Orange): Sudah check-in tapi belum check-out

### API Endpoints

#### Get Attendance Report
```
GET /api/v1/attendance/report
Query Parameters:
  - start_date: YYYY-MM-DD (required)
  - end_date: YYYY-MM-DD (required)
  - employee_id: integer (optional)

Response:
{
  "success": true,
  "data": [
    {
      "employee_id": 3,
      "full_name": "Test User",
      "position": "Kepala SPPG",
      "total_days": 1,
      "total_hours": 8.5,
      "average_hours": 8.5
    }
  ]
}
```

#### Get Attendance by Date Range
```
GET /api/v1/attendance/by-date-range
Query Parameters:
  - start_date: YYYY-MM-DD (required)
  - end_date: YYYY-MM-DD (required)
  - employee_id: integer (optional, uses JWT if not provided)

Response:
{
  "success": true,
  "data": [
    {
      "id": 37,
      "employee_id": 3,
      "date": "2026-03-04T13:23:43.407083+07:00",
      "check_in": "2026-03-04T13:23:43.407083+07:00",
      "check_out": "2026-03-04T13:23:51.792469+07:00",
      "work_hours": 0.002329273888888889,
      "ssid": "AUTO-DETECT",
      "bssid": "00:00:00:00:00:00"
    }
  ]
}
```

### Testing

#### Test dengan cURL
```bash
# 1. Login untuk mendapatkan token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"TEST001","password":"password123"}' \
  2>/dev/null | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

# 2. Get attendance report
curl -X GET "http://localhost:8080/api/v1/attendance/report?start_date=2026-03-01&end_date=2026-03-31" \
  -H "Authorization: Bearer $TOKEN" \
  2>/dev/null | python3 -m json.tool

# 3. Get detailed attendance
curl -X GET "http://localhost:8080/api/v1/attendance/by-date-range?start_date=2026-03-01&end_date=2026-03-31" \
  -H "Authorization: Bearer $TOKEN" \
  2>/dev/null | python3 -m json.tool
```

### Database Queries

#### Cek Data Absensi
```sql
-- Semua absensi hari ini
SELECT 
  e.nik,
  e.full_name,
  a.check_in,
  a.check_out,
  a.work_hours
FROM attendances a
JOIN employees e ON e.id = a.employee_id
WHERE DATE(a.date) = CURRENT_DATE
ORDER BY a.check_in DESC;

-- Absensi dalam periode tertentu
SELECT 
  e.full_name,
  COUNT(*) as total_days,
  SUM(a.work_hours) as total_hours,
  AVG(a.work_hours) as avg_hours
FROM attendances a
JOIN employees e ON e.id = a.employee_id
WHERE DATE(a.date) >= '2026-03-01' 
  AND DATE(a.date) <= '2026-03-31'
GROUP BY e.id, e.full_name
ORDER BY e.full_name;
```

### Changelog

#### 2026-03-04: Fix Timezone Issue
- **Problem**: Data tidak muncul karena timezone mismatch
- **Solution**: Menggunakan `DATE()` function untuk compare tanggal saja
- **Changes**:
  - Backend: `attendance_service.go` - Updated `GetAttendanceReport()` query
  - Frontend: `AttendanceReportView.vue` - Added console logging
  - Frontend: Auto-fetch report on mount

**Before:**
```sql
WHERE attendances.date >= ? AND attendances.date < ?
```

**After:**
```sql
WHERE DATE(attendances.date) >= ? AND DATE(attendances.date) <= ?
```

### Tips

1. **Gunakan Filter Periode yang Tepat**
   - Untuk laporan harian: Pilih 1 hari saja
   - Untuk laporan mingguan: Pilih 7 hari
   - Untuk laporan bulanan: Pilih 1 bulan penuh

2. **Export Data Secara Berkala**
   - Export data setiap akhir bulan untuk arsip
   - Simpan file Excel/PDF sebagai backup

3. **Monitor Tingkat Kehadiran**
   - Cek karyawan dengan tingkat kehadiran rendah (<75%)
   - Follow up dengan karyawan yang sering terlambat

4. **Gunakan Detail View**
   - Klik nama karyawan untuk melihat pola kehadiran
   - Identifikasi masalah (sering lupa check-out, jam kerja kurang, dll)

### Support

Jika masih ada masalah:
1. Cek log backend: `backend/app.log`
2. Cek browser console untuk error
3. Restart backend dan web server
4. Clear browser cache dan cookies
