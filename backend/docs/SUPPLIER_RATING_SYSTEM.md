# Sistem Rating dan Metrik Supplier

## Overview

Sistem ini secara otomatis menghitung dan memperbarui metrik performa supplier berdasarkan data penerimaan barang (Goods Receipt Note / GRN).

## Metrik yang Dihitung

### 1. Pengiriman Tepat Waktu (On-Time Delivery Rate)

**Cara Kerja:**
- Dihitung otomatis saat barang diterima (GRN dibuat)
- Membandingkan `receipt_date` dengan `expected_delivery` dari Purchase Order
- Formula: `(jumlah pengiriman tepat waktu / total pengiriman selesai) × 100%`

**Kriteria Tepat Waktu:**
- Barang diterima pada atau sebelum tanggal `expected_delivery`

**Update:**
- Otomatis setiap kali GRN dibuat
- Disimpan di field `suppliers.on_time_delivery`

### 2. Rating Kualitas (Quality Rating)

**Cara Kerja:**
- Diisi saat membuat GRN (form penerimaan barang)
- Field `quality_rating` di tabel `goods_receipts` (nilai 0-5)
- Rata-rata dari semua rating GRN untuk supplier tersebut

**Cara Mengisi:**
- Saat membuat GRN, isi field `quality_rating` dengan nilai 0-5
- 0 = Tidak ada rating
- 1 = Sangat buruk
- 2 = Buruk
- 3 = Cukup
- 4 = Baik
- 5 = Sangat baik

**Update:**
- Otomatis dihitung rata-rata dari semua GRN
- Disimpan di field `suppliers.quality_rating`

## API Endpoints

### 1. Create Goods Receipt (dengan rating)

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

**Response:**
```json
{
  "success": true,
  "message": "Goods receipt berhasil dibuat dan metrik supplier diperbarui",
  "goods_receipt": {
    "id": 1,
    "grn_number": "GRN-2026-001",
    "po_id": 1,
    "quality_rating": 4.5,
    ...
  }
}
```

### 2. Get Supplier Performance

```http
GET /api/v1/suppliers/{id}/performance
Authorization: Bearer {token}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "supplier_id": 1,
    "supplier_name": "PT Menyala Abangku",
    "total_orders": 10,
    "completed_orders": 8,
    "on_time_deliveries": 7,
    "on_time_rate": 87.5,
    "quality_rating": 4.2,
    "total_amount": 50000000,
    "transactions": [...]
  }
}
```

## Database Schema

### goods_receipts Table

```sql
CREATE TABLE goods_receipts (
  id INTEGER PRIMARY KEY,
  grn_number VARCHAR(50) UNIQUE NOT NULL,
  po_id INTEGER NOT NULL,
  receipt_date DATETIME NOT NULL,
  received_by INTEGER NOT NULL,
  notes TEXT,
  quality_rating REAL DEFAULT 0,  -- NEW FIELD
  created_at DATETIME,
  FOREIGN KEY (po_id) REFERENCES purchase_orders(id)
);
```

### suppliers Table

```sql
CREATE TABLE suppliers (
  id INTEGER PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  ...
  on_time_delivery REAL DEFAULT 0,  -- Auto-calculated
  quality_rating REAL DEFAULT 0,     -- Auto-calculated average
  ...
);
```

## Workflow

### Saat Membuat GRN:

1. User mengisi form GRN dengan rating kualitas (0-5)
2. System menyimpan GRN dengan `quality_rating`
3. System otomatis:
   - Menghitung on-time delivery (membandingkan tanggal terima vs expected)
   - Menghitung rata-rata quality rating dari semua GRN
   - Update `suppliers.on_time_delivery` dan `suppliers.quality_rating`

### Saat Melihat Detail Supplier:

1. System menampilkan metrik terkini:
   - Pengiriman Tepat Waktu (%)
   - Rating Kualitas (rata-rata dari GRN)
2. Riwayat transaksi (10 terakhir)

## Frontend Integration

### Form GRN (Penerimaan Barang)

Tambahkan field rating di form:

```vue
<a-form-item label="Rating Kualitas Barang">
  <a-rate v-model:value="formData.quality_rating" allow-half />
  <a-typography-text type="secondary">
    Berikan rating untuk kualitas barang yang diterima (1-5)
  </a-typography-text>
</a-form-item>
```

### Detail Supplier

Metrik ditampilkan otomatis dengan penjelasan:

```vue
<a-statistic title="Pengiriman Tepat Waktu" :value="supplier.on_time_delivery" suffix="%" />
<a-typography-text type="secondary">
  Dihitung otomatis saat barang diterima
</a-typography-text>

<a-statistic title="Rating Kualitas">
  <template #formatter>
    <a-rate :value="supplier.quality_rating" disabled allow-half />
  </template>
</a-statistic>
<a-typography-text type="secondary">
  Rata-rata dari rating saat penerimaan barang (GRN)
</a-typography-text>
```

## Migration

File: `backend/migrations/add_quality_rating_to_goods_receipts.sql`

```sql
ALTER TABLE goods_receipts ADD COLUMN quality_rating REAL DEFAULT 0;
UPDATE goods_receipts SET quality_rating = 0 WHERE quality_rating IS NULL;
```

## Testing

### Test Scenario 1: On-Time Delivery

1. Buat PO dengan `expected_delivery = 2026-03-10`
2. Buat GRN dengan `receipt_date = 2026-03-09` (tepat waktu)
3. Verify: `on_time_rate` meningkat

### Test Scenario 2: Late Delivery

1. Buat PO dengan `expected_delivery = 2026-03-10`
2. Buat GRN dengan `receipt_date = 2026-03-12` (terlambat)
3. Verify: `on_time_rate` menurun

### Test Scenario 3: Quality Rating

1. Buat GRN dengan `quality_rating = 5`
2. Buat GRN lagi dengan `quality_rating = 3`
3. Verify: `quality_rating` supplier = 4.0 (rata-rata)

## Notes

- Rating 0 tidak dihitung dalam rata-rata (dianggap tidak ada rating)
- On-time delivery hanya dihitung untuk PO dengan status "received"
- Metrik diupdate secara asynchronous untuk tidak memblokir response GRN
- Supplier performance dapat di-refresh manual dengan memanggil `GET /suppliers/{id}/performance`

## Future Enhancements

1. Rating per kategori (kualitas, kemasan, dokumentasi)
2. Trend analysis (performa per bulan)
3. Alert untuk supplier dengan performa rendah
4. Weighted average rating (GRN terbaru lebih berpengaruh)
