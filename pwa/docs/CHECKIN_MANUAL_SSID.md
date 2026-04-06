# Check-in Flow dengan Manual SSID Input

## Overview
User diminta input SSID WiFi secara manual karena browser tidak bisa mendeteksi SSID otomatis.

## Flow Check-in Baru

### 1. User Klik "Check In"

### 2. Load Authorized Networks
```javascript
GET /api/v1/wifi-config?active_only=true
```

### 3. Tampilkan Dialog Input SSID
Dialog menampilkan:
- Daftar WiFi yang diotorisasi
- Input field untuk SSID
- Tombol "Lanjut" atau "Gunakan GPS"

### 4a. User Input SSID (Path 1)
- User ketik nama WiFi yang sedang terhubung
- Sistem validasi apakah SSID ada di authorized networks
- Jika valid → Check-in dengan SSID tersebut
- Jika tidak valid → Lanjut ke GPS validation

### 4b. User Pilih "Gunakan GPS" (Path 2)
- Skip input SSID
- Langsung ke GPS validation

### 5. Check-in Request
```javascript
POST /api/v1/attendance/check-in
{
  "ssid": "SPPG-Office",  // SSID yang diinput user
  "bssid": "00:11:22:33:44:55"  // BSSID dari database
}
```

### 6. GPS Validation (Fallback)
Jika SSID tidak valid atau user skip:
- Request GPS location
- Validasi lokasi
- Check-in dengan SSID dari authorized network pertama

## Kenapa Perlu Manual Input?

Browser tidak bisa akses SSID WiFi karena:
1. **Security Restrictions**: Browser API tidak expose SSID
2. **Privacy**: Mencegah website tracking lokasi via WiFi
3. **Platform Limitations**: Hanya native app yang bisa akses WiFi info

## Alternatif Solusi

### Untuk Production:
1. **Native App**: Build sebagai native app (React Native/Flutter)
2. **QR Code**: Scan QR code di kantor yang berisi SSID
3. **NFC**: Tap NFC tag di kantor
4. **Geofencing**: Pure GPS-based validation
5. **Bluetooth Beacon**: Detect beacon di kantor

### Untuk Development/Testing:
- Manual input SSID (current implementation)
- GPS validation sebagai fallback

## Testing

1. **Refresh PWA** di browser
2. **Klik "Check In"**
3. **Dialog muncul** dengan daftar WiFi authorized
4. **Input SSID**: Ketik "SPPG-Office"
5. **Klik OK**
6. **Check-in berhasil** dengan SSID yang diinput

## Console Logs

```javascript
console.log('Check-in request (Manual SSID):', {
  ssid: "SPPG-Office",
  bssid: "00:11:22:33:44:55"
})
```

## Error Handling

- **SSID tidak valid**: "WiFi tidak diotorisasi. Mencoba validasi GPS..."
- **GPS gagal**: "Check-in gagal. WiFi tidak valid dan GPS tidak tersedia"
- **Already checked in**: "Anda sudah melakukan check-in hari ini"
- **Network error**: "Terjadi kesalahan saat check-in"
