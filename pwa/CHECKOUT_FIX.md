# Check-out Fix Documentation

## Problem

Check-out gagal dengan error 400 (Bad Request).

### Root Cause

PWA mengirim data yang salah ke backend:

**PWA mengirim:**
```javascript
{
  employee_id: 4,
  attendance_id: 35,
  check_out_time: "2026-03-04T13:11:00.000Z"
}
```

**Backend mengharapkan:**
```go
type CheckInRequest struct {
    SSID  string `json:"ssid" binding:"required"`
    BSSID string `json:"bssid" binding:"required"`
}
```

Backend CheckOut handler menggunakan struct yang sama dengan CheckIn, hanya memerlukan `ssid` dan `bssid` untuk validasi WiFi/IP.

## Solution

Update `attendanceService.checkOut()` untuk mengirim format yang benar.

### Files Modified:

#### 1. `pwa/src/services/attendanceService.js`

**Before:**
```javascript
async checkOut() {
  const checkOutData = {
    employee_id: authStore.user.id,
    attendance_id: this.currentAttendance.id,
    check_out_time: new Date().toISOString()
  }
  const response = await attendanceAPI.checkOut(checkOutData)
  // ...
}
```

**After:**
```javascript
async checkOut() {
  // Backend expects same format as check-in: { ssid, bssid }
  // Use dummy values for IP-based validation
  const checkOutData = {
    ssid: 'AUTO-DETECT',
    bssid: '00:00:00:00:00:00'
  }
  const response = await attendanceAPI.checkOut(checkOutData)
  // ...
}
```

#### 2. `pwa/src/views/AttendanceView.vue`

**Updated `performCheckOut()`:**
- Update both `currentAttendance.value` and `attendanceService.currentAttendance`
- Call `initializeAttendance()` after successful checkout
- Better error handling with specific messages
- Use `showNotify` instead of `showToast` for consistency

## Backend Flow

### CheckOut Handler (`backend/internal/handlers/hrm_handler.go`):

```go
func (h *HRMHandler) CheckOut(c *gin.Context) {
    var req CheckInRequest  // Expects: { ssid, bssid }
    
    // Get employee from JWT token
    userID, _ := c.Get("user_id")
    employee, _ := h.employeeService.GetEmployeeByUserID(userID.(uint))
    
    // Validate WiFi/IP and perform checkout
    attendance, err := h.attendanceService.CheckOut(employee.ID, req.SSID, req.BSSID)
    
    // Return attendance with check_out time and work_hours
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Check-out berhasil",
        "data":    attendance,
    })
}
```

### Key Points:

1. **Employee ID**: Diambil dari JWT token, tidak perlu dikirim di request body
2. **Attendance ID**: Backend otomatis cari attendance hari ini yang belum checkout
3. **Check-out Time**: Backend otomatis set ke `time.Now()`
4. **Work Hours**: Backend otomatis hitung dari selisih check-in dan check-out
5. **SSID/BSSID**: Untuk validasi WiFi/IP (gunakan dummy value untuk IP-based)

## Testing

### Test Flow:

1. **Check-in**:
   ```
   POST /api/v1/attendance/check-in
   Body: { ssid: "AUTO-DETECT", bssid: "00:00:00:00:00:00" }
   Response: { success: true, data: { id: 35, check_in: "...", check_out: null } }
   ```

2. **Check-out**:
   ```
   POST /api/v1/attendance/check-out
   Body: { ssid: "AUTO-DETECT", bssid: "00:00:00:00:00:00" }
   Response: { success: true, data: { id: 35, check_in: "...", check_out: "...", work_hours: 8.5 } }
   ```

### Expected Behavior:

- ✅ Check-in: Tombol berubah dari "Check In" → "Check Out"
- ✅ Check-out: Tombol berubah dari "Check Out" → "Check In"
- ✅ Status menampilkan jam masuk dan keluar
- ✅ Menampilkan total jam kerja
- ✅ Riwayat ter-update dengan data terbaru

## Error Handling

### Common Errors:

1. **400 Bad Request**:
   - Cause: Request body format salah
   - Solution: Pastikan mengirim `{ ssid, bssid }`

2. **400 NOT_CHECKED_IN**:
   - Cause: Belum check-in hari ini
   - Solution: Check-in dulu sebelum checkout

3. **409 ALREADY_CHECKED_OUT**:
   - Cause: Sudah checkout hari ini
   - Solution: Reload data untuk update UI

4. **403 INVALID_WIFI**:
   - Cause: IP tidak valid (jika tidak di localhost)
   - Solution: Pastikan IP dalam range yang dikonfigurasi

## Additional Improvements

### Error 409 Handling:

Ketika user sudah check-in dan mencoba check-in lagi (error 409), sekarang sistem akan:
1. Menampilkan pesan "Anda sudah melakukan check-in hari ini"
2. Otomatis reload data untuk update UI
3. Tombol berubah ke "Check Out"

**Code:**
```javascript
if (error.response.status === 409) {
  errorMessage = 'Anda sudah melakukan check-in hari ini'
  // Reload data to update UI even if already checked in
  await initializeAttendance()
}
```

## Summary

Perbaikan ini memastikan:
- ✅ Check-out menggunakan format request yang benar
- ✅ UI ter-update setelah check-in/check-out
- ✅ Error handling yang lebih baik
- ✅ Konsistensi dengan IP-based validation
- ✅ State management yang proper (update both local and service state)
