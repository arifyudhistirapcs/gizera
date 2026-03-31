# Fitur Expand/Collapse Detail Kehadiran

## Overview

Fitur baru yang memungkinkan user melihat detail log kehadiran setiap karyawan langsung di tabel laporan tanpa perlu membuka modal terpisah. Menggunakan expandable rows dari Ant Design Vue.

## Fitur

### 1. Expandable Rows
- Setiap baris di tabel laporan memiliki icon expand (+) di sebelah kiri
- Klik icon untuk expand/collapse detail kehadiran
- Multiple rows dapat di-expand sekaligus

### 2. Detail Kehadiran
Saat row di-expand, akan muncul tabel detail yang menampilkan:
- **Tanggal**: Tanggal kehadiran (format: DD/MM/YYYY)
- **Check In**: Waktu check-in dengan icon dan tag hijau
- **Check Out**: Waktu check-out dengan icon dan tag merah (atau "Belum Check Out" jika belum)
- **Jam Kerja**: Total jam kerja (format: X.X jam)
- **Status**: Status kehadiran dengan color coding
  - **Lengkap** (Hijau): ≥8 jam
  - **Cukup** (Biru): 6-8 jam
  - **Kurang** (Merah): <6 jam
  - **Belum Check Out** (Orange): Sudah check-in tapi belum check-out

### 3. Lazy Loading
- Detail data hanya dimuat saat row di-expand (tidak semua data dimuat sekaligus)
- Menampilkan loading spinner saat memuat data
- Data di-cache setelah dimuat pertama kali

### 4. Visual Design
- Background abu-abu terang untuk area expanded
- Tabel detail dengan styling yang konsisten
- Icon yang jelas untuk check-in/check-out
- Color coding untuk status

## Cara Menggunakan

### 1. Buka Laporan Absensi
```
Web Admin → SDM → Laporan Absensi
```

### 2. Filter Data
- Pilih periode tanggal
- (Opsional) Pilih karyawan tertentu
- Klik "Cari"

### 3. Expand Detail
- Klik icon **+** di sebelah kiri nama karyawan
- Tabel detail akan muncul di bawah row tersebut
- Loading spinner akan muncul saat memuat data

### 4. Collapse Detail
- Klik icon **-** untuk menutup detail
- Data tetap ter-cache untuk akses cepat berikutnya

### 5. Multiple Expand
- Anda bisa expand multiple rows sekaligus
- Setiap row memiliki data detail sendiri

## Technical Implementation

### Components Used
- `a-table` dengan prop `:expandedRowKeys` dan `@expand`
- `expandedRowRender` slot untuk custom content
- Icons: `ClockCircleOutlined`, `LoginOutlined`, `LogoutOutlined`

### State Management
```javascript
const expandedRowKeys = ref([])        // Array of expanded employee IDs
const expandedRowData = ref({})        // Object mapping employee_id to detail data
const expandedRowLoading = ref({})     // Object mapping employee_id to loading state
```

### API Calls
```javascript
// Called when row is expanded
GET /api/v1/attendance/by-date-range
Query Parameters:
  - employee_id: integer (required)
  - start_date: YYYY-MM-DD (required)
  - end_date: YYYY-MM-DD (required)

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
      "work_hours": 0.002329273888888889
    }
  ]
}
```

### Functions

#### onExpand
```javascript
const onExpand = async (expanded, record) => {
  if (expanded) {
    // Add to expanded keys
    expandedRowKeys.value.push(record.employee_id)
    
    // Load detail data if not cached
    if (!expandedRowData.value[record.employee_id]) {
      await loadExpandedRowData(record)
    }
  } else {
    // Remove from expanded keys
    expandedRowKeys.value = expandedRowKeys.value.filter(
      key => key !== record.employee_id
    )
  }
}
```

#### loadExpandedRowData
```javascript
const loadExpandedRowData = async (record) => {
  expandedRowLoading.value[record.employee_id] = true
  
  try {
    const response = await attendanceService.getAttendanceByDateRange(
      record.employee_id,
      dateRange.value[0].format('YYYY-MM-DD'),
      dateRange.value[1].format('YYYY-MM-DD')
    )
    
    expandedRowData.value[record.employee_id] = response.data || []
  } catch (error) {
    message.error('Gagal memuat detail absensi')
    expandedRowData.value[record.employee_id] = []
  } finally {
    expandedRowLoading.value[record.employee_id] = false
  }
}
```

## Styling

### Expandable Row
```css
/* Remove padding from expanded row cell */
:deep(.ant-table-expanded-row > td) {
  padding: 0 !important;
  background: #f5f5f5;
}

/* Expand icon color */
:deep(.ant-table-row-expand-icon) {
  color: #5A4372;
}

:deep(.ant-table-row-expand-icon:hover) {
  color: #7B5E9D;
}
```

### Detail Table
```css
/* Detail table styling */
:deep(.ant-table-small .ant-table-tbody > tr > td) {
  padding: 8px;
}

:deep(.ant-table-expanded-row .ant-table-thead > tr > th) {
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
}
```

## Performance Optimization

### 1. Lazy Loading
- Detail data hanya dimuat saat row di-expand
- Tidak memuat semua detail data di awal

### 2. Caching
- Data yang sudah dimuat disimpan di `expandedRowData`
- Tidak perlu fetch ulang saat expand/collapse/expand lagi

### 3. Conditional Rendering
- Loading spinner hanya muncul saat memuat data
- Empty state hanya muncul jika tidak ada data

## User Experience

### Before (Modal)
1. Klik nama karyawan
2. Modal muncul (overlay)
3. Lihat detail
4. Tutup modal
5. Scroll untuk cari karyawan lain
6. Repeat

### After (Expandable)
1. Klik icon expand
2. Detail muncul inline
3. Lihat detail
4. Scroll ke karyawan lain
5. Expand lagi (bisa multiple)
6. Bandingkan data side-by-side

### Advantages
- ✅ Tidak perlu buka/tutup modal
- ✅ Bisa expand multiple rows sekaligus
- ✅ Lebih cepat untuk membandingkan data
- ✅ Tidak kehilangan context (tetap di halaman yang sama)
- ✅ Lebih intuitive dan modern

## Browser Compatibility

Tested on:
- Chrome 120+
- Firefox 120+
- Safari 17+
- Edge 120+

## Troubleshooting

### Detail tidak muncul saat expand
**Solusi:**
1. Cek console untuk error
2. Pastikan periode tanggal sudah dipilih
3. Cek Network tab untuk melihat API response

### Loading terus-menerus
**Solusi:**
1. Cek backend log untuk error
2. Pastikan endpoint `/attendance/by-date-range` berfungsi
3. Test dengan cURL:
```bash
TOKEN="your-token"
curl -X GET "http://localhost:8080/api/v1/attendance/by-date-range?employee_id=3&start_date=2026-03-01&end_date=2026-03-31" \
  -H "Authorization: Bearer $TOKEN"
```

### Data tidak ter-cache
**Solusi:**
1. Cek `expandedRowData` di Vue DevTools
2. Pastikan `employee_id` konsisten
3. Clear cache dengan reset filter

## Future Enhancements

### Possible Improvements
1. **Export Detail**: Export detail kehadiran per karyawan
2. **Filter Detail**: Filter detail berdasarkan status (Lengkap/Kurang/dll)
3. **Sort Detail**: Sort detail berdasarkan tanggal/jam kerja
4. **Pagination Detail**: Jika data terlalu banyak
5. **Summary Row**: Tampilkan summary di bawah detail table
6. **Edit Inline**: Edit check-in/check-out langsung dari detail table
7. **Bulk Actions**: Select multiple detail rows untuk bulk operations

### Performance Improvements
1. Virtual scrolling untuk detail table yang panjang
2. Debounce untuk expand/collapse yang cepat
3. Prefetch data untuk row berikutnya

## Changelog

### Version 1.0.0 (2026-03-04)
- ✅ Initial implementation of expandable rows
- ✅ Lazy loading of detail data
- ✅ Caching mechanism
- ✅ Visual design with icons and color coding
- ✅ Loading states and error handling
- ✅ Removed old modal implementation

## Related Files

- `web/src/views/AttendanceReportView.vue` - Main component
- `web/src/services/attendanceService.js` - API service
- `backend/internal/handlers/hrm_handler.go` - Backend handler
- `backend/internal/services/attendance_service.go` - Backend service

## Screenshots

### Collapsed State
```
┌─────────────────────────────────────────────────────────┐
│ [+] Test User    │ Kepala SPPG │ 1 │ 8.5 jam │ 8.5 jam │
└─────────────────────────────────────────────────────────┘
```

### Expanded State
```
┌─────────────────────────────────────────────────────────┐
│ [-] Test User    │ Kepala SPPG │ 1 │ 8.5 jam │ 8.5 jam │
├─────────────────────────────────────────────────────────┤
│   Detail Kehadiran - Test User                          │
│   ┌───────────┬──────────┬──────────┬─────────┬────────┐│
│   │ Tanggal   │ Check In │ Check Out│ Jam Kerja│ Status ││
│   ├───────────┼──────────┼──────────┼─────────┼────────┤│
│   │ 04/03/2026│ 08:00    │ 17:00    │ 8.5 jam │ Lengkap││
│   └───────────┴──────────┴──────────┴─────────┴────────┘│
└─────────────────────────────────────────────────────────┘
```

## Support

Jika ada masalah atau pertanyaan:
1. Cek console browser untuk error
2. Cek Network tab untuk API response
3. Cek backend log untuk server error
4. Refer to `ATTENDANCE_REPORT_GUIDE.md` untuk troubleshooting umum
