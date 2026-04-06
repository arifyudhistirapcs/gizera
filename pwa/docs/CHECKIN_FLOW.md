# Check-in Flow Documentation

## Overview
Aplikasi PWA menggunakan validasi otomatis untuk check-in karyawan dengan dua metode: WiFi dan GPS.

## Flow Check-in

### 1. User Klik Tombol "Check In"

### 2. Load Authorized Networks
```javascript
GET /api/v1/wifi-config?active_only=true

Response:
{
  "success": true,
  "data": [
    {
      "id": 1,
      "ssid": "SPPG-Office",
      "bssid": "00:11:22:33:44:55",
      "location": "Kantor Pusat",
      "is_active": true
    }
  ]
}
```

### 3. Validasi WiFi (Step 1)
- Cek apakah device terhubung ke WiFi menggunakan Network Information API
- Jika terhubung WiFi, lanjut ke check-in
- Jika tidak, lanjut ke validasi GPS

### 4. Check-in Request (WiFi Valid)
```javascript
POST /api/v1/hrm/attendance/check-in

Request Body:
{
  "ssid": "SPPG-Office",
  "bssid": "00:11:22:33:44:55"
}

Response (Success):
{
  "success": true,
  "message": "Check-in berhasil",
  "data": {
    "id": 123,
    "employee_id": 1,
    "check_in_time": "2024-02-24T10:00:00Z",
    "ssid": "SPPG-Office",
    "bssid": "00:11:22:33:44:55"
  }
}

Response (Error - Invalid WiFi):
{
  "success": false,
  "error_code": "INVALID_WIFI",
  "message": "Anda harus terhubung ke Wi-Fi kantor untuk absen"
}
```

### 5. Validasi GPS (Step 2 - jika WiFi gagal)
- Request GPS location dari browser
- Validasi lokasi (saat ini permisif, akan ditambahkan GPS boundaries)
- Jika valid, kirim check-in request dengan SSID/BSSID dari authorized network pertama

### 6. Check-in Request (GPS Valid)
```javascript
POST /api/v1/hrm/attendance/check-in

Request Body:
{
  "ssid": "SPPG-Office",  // dari authorized network
  "bssid": "00:11:22:33:44:55"
}
```

### 7. Error Handling
Jika kedua validasi gagal:
```
"Check-in gagal. Anda tidak berada di area kantor. 
Pastikan terhubung ke Wi-Fi kantor atau berada di lokasi kantor."
```

## Backend Validation

Backend melakukan validasi di `AttendanceService.CheckIn()`:
1. Cek apakah SSID dan BSSID ada di tabel `wifi_configs` dengan `is_active = true`
2. Cek apakah employee sudah check-in hari ini
3. Jika valid, simpan attendance record

## Console Logs untuk Debugging

PWA akan menampilkan console logs:
```javascript
console.log('Loaded authorized networks:', authorizedNetworks)
console.log('Check-in request:', checkInData)
console.log('Check-in request (GPS):', checkInData)
```

## Testing

1. Buka browser DevTools (F12)
2. Buka tab Console
3. Klik tombol "Check In"
4. Lihat request yang dikirim:
   - SSID dan BSSID harus sesuai dengan data di web admin
   - Jika WiFi tidak terdeteksi, akan coba GPS
   - Jika keduanya gagal, akan tampil error

## Konfigurasi WiFi di Web Admin

1. Login ke web admin
2. Menu: SDM → Konfigurasi Wi-Fi
3. Tambah jaringan WiFi dengan:
   - SSID: Nama jaringan WiFi kantor
   - BSSID: MAC address router (format: 00:11:22:33:44:55)
   - Lokasi: Nama lokasi kantor
   - Status: Aktif

## Notes

- Browser tidak bisa mendeteksi SSID/BSSID secara langsung karena security restrictions
- Network Information API hanya bisa deteksi tipe koneksi (wifi/cellular/unknown)
- Untuk production, sebaiknya tambahkan GPS boundaries di model WiFiConfig
