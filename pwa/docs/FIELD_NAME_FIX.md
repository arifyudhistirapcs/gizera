# Field Name Compatibility Fix

## Problem

Backend API mengirim field dengan nama:
- `check_in` (bukan `check_in_time`)
- `check_out` (bukan `check_out_time`)

Tapi PWA frontend menggunakan:
- `check_in_time`
- `check_out_time`

Ini menyebabkan:
1. ❌ Tombol Check In/Check Out tidak berubah status
2. ❌ Status absensi hari ini tidak muncul
3. ❌ Riwayat absensi tidak menampilkan jam masuk/keluar

## Solution

Update semua referensi di PWA untuk support kedua format field name.

### Files Modified:

#### 1. `pwa/src/services/attendanceService.js`

**Method `canCheckIn()`:**
```javascript
canCheckIn() {
  // Support both field names
  return !this.currentAttendance || 
         (this.currentAttendance && (this.currentAttendance.check_out || this.currentAttendance.check_out_time))
}
```

**Method `canCheckOut()`:**
```javascript
canCheckOut() {
  // Support both field names
  return this.currentAttendance && 
         (this.currentAttendance.check_in || this.currentAttendance.check_in_time) && 
         !(this.currentAttendance.check_out || this.currentAttendance.check_out_time)
}
```

#### 2. `pwa/src/views/AttendanceView.vue`

**Template - Status Absensi Hari Ini:**
```vue
<div v-if="currentAttendance.check_in || currentAttendance.check_in_time" class="check-in">
  <van-icon name="play-circle-o" color="#07c160" />
  <span>Masuk: {{ formatTime(currentAttendance.check_in || currentAttendance.check_in_time) }}</span>
</div>
<div v-if="currentAttendance.check_out || currentAttendance.check_out_time" class="check-out">
  <van-icon name="stop-circle-o" color="#ee0a24" />
  <span>Keluar: {{ formatTime(currentAttendance.check_out || currentAttendance.check_out_time) }}</span>
</div>
```

**Function `getAttendanceLabel()`:**
```javascript
function getAttendanceLabel(record) {
  const checkIn = (record.check_in || record.check_in_time) 
    ? formatTime(record.check_in || record.check_in_time) 
    : '-'
  const checkOut = (record.check_out || record.check_out_time) 
    ? formatTime(record.check_out || record.check_out_time) 
    : '-'
  return `Masuk: ${checkIn} | Keluar: ${checkOut}`
}
```

**Function `performAutoCheckIn()` - Update service state:**
```javascript
if (response.data.success) {
  // Update both local and service attendance
  currentAttendance.value = response.data.data
  attendanceService.currentAttendance = response.data.data
  
  // ... rest of code
}
```

## Backend API Response Format

### `/api/v1/attendance/today`
```json
{
  "success": true,
  "data": {
    "id": 33,
    "employee_id": 3,
    "date": "2026-03-04T12:55:50.251694+07:00",
    "check_in": "2026-03-04T12:55:50.251695+07:00",
    "check_out": null,
    "work_hours": 0,
    "ssid": "SPPG-Office",
    "bssid": "00:11:22:33:44:55"
  }
}
```

### `/api/v1/attendance/by-date-range`
```json
{
  "success": true,
  "data": [
    {
      "id": 33,
      "employee_id": 3,
      "date": "2026-03-04T12:55:50.251694+07:00",
      "check_in": "2026-03-04T12:55:50.251695+07:00",
      "check_out": null,
      "work_hours": 0
    }
  ]
}
```

## Testing

### Before Fix:
- ❌ Tombol tetap "Check In" setelah check-in
- ❌ Status hari ini: "Belum ada absensi hari ini"
- ❌ Riwayat: "Masuk: - | Keluar: -"

### After Fix:
- ✅ Tombol berubah menjadi "Check Out" setelah check-in
- ✅ Status hari ini: "Masuk: 12:55"
- ✅ Riwayat: "Masuk: 12:55 | Keluar: -"

## Why This Approach?

Menggunakan fallback (`field1 || field2`) memastikan kompatibilitas dengan:
1. Backend API yang mengirim `check_in`/`check_out`
2. Kode lama yang mungkin menggunakan `check_in_time`/`check_out_time`
3. Future changes tanpa breaking existing code

## Alternative Solution (Not Recommended)

Mengubah backend untuk mengirim `check_in_time` dan `check_out_time` akan:
- ❌ Breaking change untuk web admin
- ❌ Perlu update database schema
- ❌ Perlu update semua API consumers

Lebih baik frontend yang adaptif terhadap backend API.
