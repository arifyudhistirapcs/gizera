# Field "Diterima Oleh" di Goods Receipt Note (GRN)

## Overview

Field "Diterima Oleh" (`received_by`) di GRN diisi **otomatis** dari user yang sedang login saat membuat GRN. Ini untuk audit trail - mencatat siapa yang menerima barang dari supplier.

## Cara Kerja

### 1. Saat Membuat GRN (Frontend)

User **TIDAK** perlu mengisi field "Diterima Oleh" di form. Form hanya berisi:
- Purchase Order (pilih PO yang akan diterima)
- Tanggal Penerimaan
- Foto Invoice/Nota
- Rating Kualitas Barang (1-5 bintang)
- Item yang Diterima (jumlah, tanggal kadaluarsa)
- Catatan (opsional)

### 2. Backend Processing

Saat API `POST /api/v1/goods-receipts` dipanggil:

1. **Middleware JWT** mengekstrak `user_id` dari token
2. **Handler** mengambil `user_id` dari context:
   ```go
   userID, _ := c.Get("user_id")
   ```
3. **Service** mengisi field `ReceivedBy`:
   ```go
   grn.ReceivedBy = userID
   grn.ReceiptDate = time.Now()
   ```

### 3. Response ke Frontend

Backend mengembalikan GRN dengan field tambahan `received_by_name`:

```json
{
  "success": true,
  "goods_receipt": {
    "id": 1,
    "grn_number": "GRN-20260305-0001",
    "received_by": 4,
    "received_by_name": "Test User",
    ...
  }
}
```

## Database Schema

```sql
CREATE TABLE goods_receipts (
  id INTEGER PRIMARY KEY,
  grn_number VARCHAR(50) UNIQUE NOT NULL,
  po_id INTEGER NOT NULL,
  receipt_date DATETIME NOT NULL,
  received_by INTEGER NOT NULL,  -- Auto-filled from JWT token
  notes TEXT,
  quality_rating REAL DEFAULT 0,
  invoice_photo VARCHAR(500),
  created_at DATETIME,
  FOREIGN KEY (received_by) REFERENCES users(id)
);
```

## Frontend Display

### Tabel List GRN

```vue
<a-table-column title="Diterima Oleh" dataIndex="received_by_name" />
```

Data ditampilkan dari field `received_by_name` yang digenerate backend.

### Detail Modal

```vue
<a-descriptions-item label="Diterima Oleh">
  {{ selectedGRN.received_by_name }}
</a-descriptions-item>
```

## Security & Audit Trail

### Keuntungan Auto-Fill:

1. **Tidak bisa dimanipulasi** - User tidak bisa mengisi nama orang lain
2. **Audit trail akurat** - Selalu tercatat siapa yang benar-benar menerima barang
3. **Sesuai dengan JWT token** - Terikat dengan session login user

### Validasi:

- JWT token harus valid
- User harus memiliki role yang berhak membuat GRN (pengadaan, kepala_sppg, dll)
- `received_by` selalu terisi (NOT NULL di database)

## Testing

### Test Case 1: Create GRN dengan User Login

1. Login sebagai user (contoh: kepala.sppg@sppg.com)
2. Buat GRN baru
3. Verify: Field `received_by` terisi dengan ID user yang login
4. Verify: Field `received_by_name` menampilkan nama user

### Test Case 2: View GRN List

1. Buka halaman Penerimaan Barang
2. Verify: Kolom "Diterima Oleh" menampilkan nama user
3. Verify: Tidak ada kolom yang kosong

### Test Case 3: View GRN Detail

1. Klik "Detail" pada salah satu GRN
2. Verify: Field "Diterima Oleh" menampilkan nama lengkap user

## Troubleshooting

### Kolom "Diterima Oleh" Kosong

**Kemungkinan Penyebab:**

1. **Backend tidak mengembalikan `received_by_name`**
   - Solusi: Pastikan handler menambahkan field ini ke response
   - Check: Handler `GetAllGoodsReceipts` dan `GetGoodsReceipt`

2. **Data `Receiver` tidak di-preload**
   - Solusi: Pastikan service melakukan `Preload("Receiver")`
   - Check: `goods_receipt_service.go` line 168

3. **User tidak ada di database**
   - Solusi: Pastikan `received_by` merujuk ke user yang valid
   - Check: Foreign key constraint

### Field `received_by_name` Tidak Muncul di Response

**Solusi:**
Handler harus menggunakan map untuk menambahkan field custom:

```go
response := map[string]interface{}{
    "id": grn.ID,
    "grn_number": grn.GRNNumber,
    // ... fields lainnya
    "received_by_name": grn.Receiver.FullName,
}
```

## API Documentation

### Create GRN

```http
POST /api/v1/goods-receipts
Authorization: Bearer {token}
Content-Type: application/json

{
  "po_id": 1,
  "quality_rating": 4.5,
  "notes": "Barang diterima dalam kondisi baik",
  "items": [
    {
      "ingredient_id": 1,
      "received_quantity": 100,
      "expiry_date": "2026-12-31"
    }
  ]
}
```

**Note:** Field `received_by` TIDAK perlu dikirim - akan diisi otomatis dari token.

### Get All GRNs

```http
GET /api/v1/goods-receipts
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "goods_receipts": [
    {
      "id": 1,
      "grn_number": "GRN-20260305-0001",
      "received_by": 4,
      "received_by_name": "Test User",
      "receiver": {
        "id": 4,
        "full_name": "Test User",
        "email": "test@sppg.com"
      },
      ...
    }
  ]
}
```

## Summary

- ✅ Field "Diterima Oleh" diisi **otomatis** dari user yang login
- ✅ User **TIDAK** perlu mengisi field ini di form
- ✅ Backend mengambil `user_id` dari JWT token
- ✅ Frontend menampilkan `received_by_name` di tabel dan detail
- ✅ Audit trail akurat dan tidak bisa dimanipulasi
