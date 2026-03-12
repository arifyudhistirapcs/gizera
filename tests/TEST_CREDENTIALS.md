# Test Credentials Configuration

## Overview
Test suite memerlukan credentials yang valid untuk menjalankan authentication tests. Credentials default perlu diganti dengan data yang sesuai dari database Anda.

## Current Test Credentials

File yang perlu diupdate: `tests/test-suites/authentication.spec.js`

### Default Credentials (Perlu Diganti)
```javascript
Username: admin
Password: admin123
```

## Cara Update Credentials

### Option 1: Update Langsung di Test File

1. Buka file `tests/test-suites/authentication.spec.js`
2. Cari semua baris dengan comment `// TODO: Replace with actual test credentials`
3. Ganti `'admin'` dan `'admin123'` dengan credentials yang valid dari database

**Contoh:**
```javascript
// Sebelum:
await usernameInput.fill('admin'); // TODO: Replace with actual test credentials
await passwordInput.fill('admin123'); // TODO: Replace with actual test credentials

// Sesudah (contoh):
await usernameInput.fill('testuser@sppg.com');
await passwordInput.fill('SecurePass123!');
```

### Option 2: Buat Test User di Database

Alternatif yang lebih baik adalah membuat dedicated test user di database:

1. Login ke database Anda
2. Buat user baru dengan credentials:
   - NIK/Email: `admin` atau `testuser@sppg.com`
   - Password: `admin123` atau password pilihan Anda
   - Role: `kepala_sppg` (untuk akses penuh ke semua fitur)

**SQL Example (sesuaikan dengan schema Anda):**
```sql
INSERT INTO users (nik, email, full_name, password, role, is_active)
VALUES (
  'TEST001',
  'testuser@sppg.com',
  'Test User',
  '$2a$10$...', -- hashed password untuk 'admin123'
  'kepala_sppg',
  true
);
```

### Option 3: Environment Variables (Recommended)

Untuk keamanan yang lebih baik, gunakan environment variables:

1. Buat file `.env` di folder `tests/`:
```env
TEST_USERNAME=admin
TEST_PASSWORD=admin123
```

2. Update `tests/test-suites/authentication.spec.js` untuk membaca dari env:
```javascript
const testUsername = process.env.TEST_USERNAME || 'admin';
const testPassword = process.env.TEST_PASSWORD || 'admin123';

await usernameInput.fill(testUsername);
await passwordInput.fill(testPassword);
```

3. Tambahkan `.env` ke `.gitignore` agar credentials tidak ter-commit

## Verifikasi Credentials

Setelah update credentials, jalankan test untuk verifikasi:

```bash
cd tests
npx playwright test authentication --headed
```

Test yang harus pass:
- ✅ auth-001: User login with valid credentials
- ✅ auth-005: User logout successfully  
- ✅ auth-006: Session persistence after page refresh

## Troubleshooting

### Test auth-001 gagal: "Logout button tidak ditemukan"
- Pastikan credentials yang digunakan valid dan bisa login
- Cek apakah user memiliki role yang sesuai (minimal role yang bisa akses dashboard)
- Cek console browser untuk error messages

### Test auth-002/auth-003 gagal: "Error message tidak terdeteksi"
- Ini normal jika backend tidak mengembalikan error message
- Test akan di-skip atau perlu adjustment pada selector

### Test auth-006 gagal: "Session tidak persist"
- Pastikan auth store menyimpan token ke localStorage
- Cek browser console untuk error saat restore auth state
- Verifikasi router guard membaca token dari localStorage

## Best Practices

1. **Jangan commit credentials** ke repository
2. **Gunakan dedicated test user** dengan role terbatas jika memungkinkan
3. **Rotate test credentials** secara berkala
4. **Document credentials** di tempat yang aman (password manager, secret management)
5. **Use environment variables** untuk CI/CD pipeline

## Contact

Jika ada pertanyaan tentang test credentials, hubungi:
- DevOps team untuk production test credentials
- Database admin untuk membuat test user
