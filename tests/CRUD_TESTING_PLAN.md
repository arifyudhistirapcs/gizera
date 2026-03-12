# Comprehensive CRUD Testing Plan

**Date**: 2025-01-15  
**Purpose**: Test actual CRUD operations, not just UI presence

---

## 🎯 Objective

Test ACTUAL data operations:
- ✅ CREATE: Add new records with real data
- ✅ READ: Verify data is displayed correctly
- ✅ UPDATE: Edit existing records
- ✅ DELETE: Remove records and verify deletion

---

## 📋 Modules Requiring CRUD Testing

### Priority 1: Core CRUD Modules

#### 1. Supply Chain - Supplier
**CRUD Operations:**
- CREATE: Add new supplier with name, contact, address
- READ: View supplier list and details
- UPDATE: Edit supplier information
- DELETE: Remove supplier

#### 2. Supply Chain - Purchase Order
**CRUD Operations:**
- CREATE: Create new PO with items and quantities
- READ: View PO list and details
- UPDATE: Modify PO items or status
- DELETE: Cancel/delete PO

#### 3. Logistik - Data Sekolah
**CRUD Operations:**
- CREATE: Add new school with name, address, coordinates
- READ: View school list and on map
- UPDATE: Edit school information
- DELETE: Remove school

#### 4. SDM - Data Karyawan
**CRUD Operations:**
- CREATE: Add new employee with role, contact
- READ: View employee list
- UPDATE: Edit employee details
- DELETE: Remove employee

#### 5. Menu Manajemen (Resep)
**CRUD Operations:**
- CREATE: Create new recipe with ingredients
- READ: View recipe list and details
- UPDATE: Edit recipe components
- DELETE: Remove recipe

#### 6. Menu Komponen (Barang Setengah Jadi)
**CRUD Operations:**
- CREATE: Add new semi-finished component
- READ: View component list
- UPDATE: Edit component details
- DELETE: Remove component

#### 7. Supply Chain - Bahan Baku
**CRUD Operations:**
- CREATE: Initialize new raw material stock
- READ: View inventory levels
- UPDATE: Adjust stock levels
- DELETE: Remove material (if applicable)

#### 8. Keuangan - Aset Dapur
**CRUD Operations:**
- CREATE: Add new kitchen asset
- READ: View asset list
- UPDATE: Update asset condition/value
- DELETE: Remove asset

### Priority 2: Transaction/Planning Modules

#### 9. Menu Perencanaan
**Operations:**
- CREATE: Plan menu for specific day
- READ: View weekly menu plan
- UPDATE: Modify planned menu
- DELETE: Remove menu from day

#### 10. Logistik - Tugas Pengiriman
**Operations:**
- CREATE: Create new delivery task
- READ: View delivery tasks
- UPDATE: Update task status
- DELETE: Cancel delivery task

#### 11. Keuangan - Arus Kas
**Operations:**
- CREATE: Add new transaction
- READ: View transaction list
- UPDATE: Edit transaction details
- DELETE: Remove transaction

#### 12. Supply Chain - Penerimaan Barang (GRN)
**Operations:**
- CREATE: Create new GRN
- READ: View GRN list
- UPDATE: Update received quantities
- DELETE: Cancel GRN

---

## 🔧 Testing Approach

### Phase 1: Explore with Playwright
For each module:
1. Navigate to the module
2. Click "Tambah" button
3. Identify all form fields
4. Document required vs optional fields
5. Document validation rules

### Phase 2: Create Test Data
- Use realistic test data
- Include edge cases
- Test validation (empty fields, invalid data)

### Phase 3: Execute CRUD Flow
1. **CREATE**: Fill form and submit
2. **VERIFY**: Check success message and data appears in list
3. **READ**: Open detail view and verify all fields
4. **UPDATE**: Edit data and save
5. **VERIFY**: Check updated data is correct
6. **DELETE**: Remove record
7. **VERIFY**: Check record is gone

### Phase 4: Cleanup
- Remove all test data
- Reset to original state

---

## 📝 Test Data Examples

### Supplier Test Data
```javascript
{
  name: "Test Supplier " + timestamp,
  contact: "081234567890",
  address: "Jl. Test No. 123",
  email: "test@supplier.com"
}
```

### School Test Data
```javascript
{
  name: "SD Test " + timestamp,
  address: "Jl. Pendidikan No. 456",
  students: 100,
  latitude: -6.200000,
  longitude: 106.816666
}
```

### Employee Test Data
```javascript
{
  name: "Test Employee " + timestamp,
  email: "test.employee@test.com",
  phone: "081234567890",
  role: "Driver",
  status: "Active"
}
```

---

## ⚠️ Important Notes

1. **Use Timestamps**: Add timestamp to names to avoid duplicates
2. **Cleanup**: Always delete test data after testing
3. **Verification**: Always verify data after each operation
4. **Error Handling**: Test validation errors
5. **Real Operations**: Actually submit forms, don't just check UI

---

## 🚀 Next Steps

1. Create exploration script to identify all form fields
2. Create CRUD test suites for each module
3. Execute tests with real data operations
4. Document results
5. Upload to Google Sheets

---

**Status**: Planning Complete - Ready for Implementation

