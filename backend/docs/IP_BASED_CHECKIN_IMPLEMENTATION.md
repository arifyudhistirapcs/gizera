# Implementasi Check-in Berbasis IP Address

## Ringkasan Perubahan

Sistem check-in karyawan sekarang menggunakan validasi berbasis IP address, bukan lagi BSSID/MAC address WiFi. Perubahan ini dilakukan karena browser tidak dapat mengakses informasi SSID/BSSID WiFi karena alasan keamanan.

## Perubahan yang Dilakukan

### 1. Backend

#### Database Schema
Tabel `wi_fi_configs` sekarang memiliki kolom tambahan:
- `ip_range` (VARCHAR 100): Range IP dalam format CIDR (contoh: 192.168.1.0/24)
- `allowed_ips` (TEXT[]): Array IP address spesifik yang diizinkan

Kolom lama (`ss_id`, `bss_id`) masih ada untuk backward compatibility tapi tidak lagi digunakan untuk validasi.

#### Service Layer (`backend/internal/services/attendance_service.go`)

**Method `ValidateIP(ipAddress string)`:**
```go
// Validasi IP address terhadap konfigurasi WiFi yang aktif
// 1. Development mode: Accept localhost (127.0.0.1, ::1, localhost)
// 2. Check IP Range (CIDR notation)
// 3. Check Allowed IPs (specific IP list)
```

**Method `isIPInRange(ipAddress, cidr string)`:**
```go
// Validasi apakah IP address berada dalam range CIDR
// Mendukung /24 dan /16 network
```

#### Handler Layer (`backend/internal/handlers/hrm_handler.go`)

**Method `CheckIn()`:**
```go
// Flow check-in:
// 1. Get client IP dari c.ClientIP()
// 2. Try IP validation first
// 3. Jika IP valid, gunakan WiFi config dari matched network
// 4. Jika IP tidak valid, fallback ke SSID/BSSID validation (untuk backward compatibility)
// 5. Return validated_by info dengan method dan client_ip
```

### 2. Frontend PWA

#### AttendanceView.vue
**Method `performAutoCheckIn()`:**
- Mengirim dummy SSID/BSSID ke backend
- Backend otomatis validasi berdasarkan IP address
- Menampilkan pesan sukses/error dengan informasi IP address
- Tidak perlu input manual SSID dari user

**Request payload:**
```javascript
{
  ssid: 'AUTO-DETECT',
  bssid: '00:00:00:00:00:00'
}
```

**Response sukses:**
```javascript
{
  success: true,
  message: 'Check-in berhasil',
  data: { /* attendance record */ },
  validated_by: {
    method: 'ip_validation',
    client_ip: '::1'  // atau IP address sebenarnya
  }
}
```

### 3. Frontend Web Admin

#### WiFiConfigView.vue
Form konfigurasi WiFi sekarang memiliki:
- **SSID**: Nama jaringan (untuk identifikasi)
- **IP Range**: Range IP dalam format CIDR (contoh: 192.168.1.0/24)
- **IP Spesifik**: Daftar IP address spesifik yang diizinkan (opsional)
- **Lokasi**: Lokasi fisik jaringan
- **Status**: Aktif/Tidak Aktif

Field BSSID/MAC Address sudah dihapus dari form.

## Cara Kerja

### Flow Check-in:

```
1. User buka PWA → Login → Klik "Check In"
   ↓
2. PWA kirim request ke backend dengan dummy SSID/BSSID
   ↓
3. Backend ambil client IP dari HTTP request headers
   ↓
4. Backend validasi IP:
   - Cek apakah IP dalam range yang dikonfigurasi (CIDR)
   - Atau cek apakah IP ada dalam daftar allowed_ips
   ↓
5. Jika valid:
   - Buat attendance record
   - Return success dengan validated_by info
   ↓
6. Jika tidak valid:
   - Return 403 error dengan client IP
```

### Development Mode:

Untuk memudahkan testing di localhost, sistem otomatis menerima IP berikut:
- `127.0.0.1` (IPv4 localhost)
- `::1` (IPv6 localhost)
- `localhost` (hostname)

**PENTING:** Hapus atau disable fitur ini di production!

## Konfigurasi

### Contoh Konfigurasi WiFi:

#### Menggunakan IP Range (CIDR):
```
SSID: SPPG-Office
IP Range: 192.168.1.0/24
Lokasi: Kantor Pusat
Status: Aktif
```
Ini akan menerima semua IP dari 192.168.1.0 sampai 192.168.1.255

#### Menggunakan IP Spesifik:
```
SSID: SPPG-Office
IP Range: (kosong)
Allowed IPs:
  - 192.168.1.100
  - 192.168.1.101
  - 192.168.1.102
Lokasi: Kantor Pusat
Status: Aktif
```
Ini hanya akan menerima 3 IP address tersebut.

#### Kombinasi:
```
SSID: SPPG-Office
IP Range: 192.168.1.0/24
Allowed IPs:
  - 10.0.0.50
  - 10.0.0.51
Lokasi: Kantor Pusat
Status: Aktif
```
Ini akan menerima semua IP dari range 192.168.1.x DAN juga 2 IP spesifik dari network lain.

## Testing

### 1. Test dari Localhost (Development):

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier": "TEST001", "password": "password123"}'

# Check-in
curl -X POST http://localhost:8080/api/v1/attendance/check-in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"ssid": "AUTO-DETECT", "bssid": "00:00:00:00:00:00"}'
```

### 2. Test dari Browser:

1. Buka http://localhost:5174/
2. Login dengan NIK: TEST001, password: password123
3. Klik tab "Absensi"
4. Klik tombol "Check In"
5. Lihat console browser untuk melihat request/response

### 3. Test dari Network Sebenarnya:

1. Deploy PWA ke device di jaringan kantor (192.168.1.x)
2. Pastikan IP range sudah dikonfigurasi di database
3. Test check-in dari device tersebut
4. Verifikasi IP address di response

## Database Migration

Jika Anda sudah memiliki data WiFi config lama, update dengan:

```sql
-- Update existing config dengan IP range
UPDATE wi_fi_configs 
SET ip_range = '192.168.1.0/24' 
WHERE id = 1;

-- Atau tambah IP spesifik
UPDATE wi_fi_configs 
SET allowed_ips = ARRAY['192.168.1.100', '192.168.1.101'] 
WHERE id = 1;
```

## Security Considerations

### 1. IP Spoofing
Client IP diambil dari `c.ClientIP()` yang menggunakan X-Forwarded-For header. Pastikan:
- Reverse proxy (nginx/apache) dikonfigurasi dengan benar
- Hanya trust X-Forwarded-For dari proxy yang valid
- Jangan expose backend langsung ke internet

### 2. Localhost Exception
Development mode menerima localhost IP. Di production:
```go
// REMOVE atau COMMENT OUT bagian ini:
if ipAddress == "127.0.0.1" || ipAddress == "::1" || ipAddress == "localhost" {
    if len(configs) > 0 {
        return true, &configs[0], nil
    }
}
```

### 3. CIDR Validation
Implementasi saat ini basic (hanya support /24 dan /16). Untuk production, pertimbangkan:
- Gunakan library CIDR yang proper (contoh: `github.com/apparentlymart/go-cidr`)
- Validasi input CIDR di frontend dan backend
- Support IPv6 jika diperlukan

### 4. Multiple Locations
Jika ada multiple kantor dengan network berbeda:
- Buat WiFi config terpisah untuk setiap lokasi
- Set IP range yang berbeda untuk setiap lokasi
- User bisa check-in dari lokasi manapun yang terdaftar

## Troubleshooting

### Check-in Gagal dengan Error 403

**Kemungkinan penyebab:**
1. IP address tidak dalam range yang dikonfigurasi
2. WiFi config tidak aktif
3. Tidak ada WiFi config yang terdaftar

**Solusi:**
```sql
-- Cek konfigurasi WiFi
SELECT id, ss_id, ip_range, allowed_ips, is_active 
FROM wi_fi_configs;

-- Cek IP address dari error message
-- Update IP range jika perlu
UPDATE wi_fi_configs 
SET ip_range = 'YOUR_IP_RANGE' 
WHERE id = 1;
```

### IP Address Salah di Response

**Kemungkinan penyebab:**
1. Behind reverse proxy yang tidak dikonfigurasi dengan benar
2. Multiple proxy layers

**Solusi:**
- Cek konfigurasi nginx/apache
- Pastikan X-Forwarded-For header di-set dengan benar
- Test dengan `curl -H "X-Forwarded-For: 192.168.1.100" ...`

### Check-in Berhasil dari IP yang Tidak Seharusnya

**Kemungkinan penyebab:**
1. Localhost exception masih aktif di production
2. IP range terlalu luas

**Solusi:**
- Disable localhost exception
- Review dan perbaiki IP range configuration
- Gunakan IP spesifik jika perlu kontrol lebih ketat

## Files Modified

### Backend:
1. `backend/internal/models/hr.go` - Added IPRange and AllowedIPs fields
2. `backend/internal/services/attendance_service.go` - Added ValidateIP() and isIPInRange()
3. `backend/internal/handlers/hrm_handler.go` - Updated CheckIn() to use IP validation

### Frontend PWA:
1. `pwa/src/views/AttendanceView.vue` - Updated performAutoCheckIn() to use IP-based validation

### Frontend Web:
1. `web/src/views/WiFiConfigView.vue` - Updated form to use IP Range instead of BSSID

## Next Steps

### Untuk Production:
1. ✅ Remove localhost exception dari ValidateIP()
2. ✅ Configure actual office IP ranges
3. ✅ Test from real network devices
4. ✅ Consider using proper CIDR library
5. ✅ Add IPv6 support if needed
6. ✅ Setup proper reverse proxy configuration
7. ✅ Add monitoring/logging for failed check-in attempts

### Untuk Enhancement:
1. Add IP geolocation untuk additional validation
2. Add time-based restrictions (hanya bisa check-in di jam kerja)
3. Add notification untuk suspicious check-in attempts
4. Add admin dashboard untuk monitoring check-in patterns
5. Add API untuk test IP validation tanpa create attendance record

## Conclusion

Implementasi check-in berbasis IP address berhasil mengatasi keterbatasan browser yang tidak bisa mengakses informasi WiFi SSID/BSSID. Sistem sekarang lebih reliable dan mudah dikonfigurasi untuk multiple locations.
