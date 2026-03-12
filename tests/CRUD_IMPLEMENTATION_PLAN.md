# CRUD Testing Implementation Plan

**Date**: 2025-01-15  
**Status**: In Progress  
**Purpose**: Implement actual CRUD operations testing for all modules

---

## 🎯 Objective

Replace UI-only tests with comprehensive CRUD tests that:
- ✅ CREATE: Actually fill forms and submit data
- ✅ READ: Verify data is displayed correctly in lists and details
- ✅ UPDATE: Edit existing records and verify changes
- ✅ DELETE: Remove records and verify deletion
- ✅ VALIDATE: Test form validation with invalid/empty data
- ✅ SEARCH: Test search and filter functionality

---

## 🔧 Implementation Approach

### Phase 1: Create CRUD Helper Utility ✅
**File**: `tests/utils/crud-helper.js`

**Features**:
- Generate unique timestamps for test data
- Fill text inputs, textareas, dropdowns
- Submit and cancel forms
- Wait for success/error messages
- Verify data in tables
- Click row actions (Edit, Delete, Detail)
- Confirm deletion dialogs
- Search and filter data
- Take screenshots for debugging

### Phase 2: Create CRUD Test Suites
**Pattern**: `tests/test-suites/{module}-crud.spec.js`

**Test Structure**:
1. **CRUD-001: CREATE** - Add new record with all required fields
2. **CRUD-002: READ** - View record details
3. **CRUD-003: UPDATE** - Edit record and verify changes
4. **CRUD-004: DELETE** - Remove record and verify deletion
5. **CRUD-005: VALIDATION** - Test form validation
6. **CRUD-006: SEARCH** - Search for records

### Phase 3: Execute and Document
- Run CRUD tests for each module
- Document results with screenshots
- Upload to Google Sheets
- Fix any bugs found

---

## 📋 Modules Priority

### Priority 1: Core CRUD Modules (Start Here)

#### 1. Supply Chain - Supplier ✅ IMPLEMENTED
**File**: `tests/test-suites/supply-chain-supplier-crud.spec.js`
**Fields**:
- Nama (Name) - Required
- Kontak (Contact) - Required
- Email - Optional
- Alamat (Address) - Optional

**Status**: CRUD test suite created

#### 2. Logistik - Data Sekolah
**File**: `tests/test-suites/logistik-data-sekolah-crud.spec.js`
**Fields**:
- Nama Sekolah (School Name) - Required
- Alamat (Address) - Required
- Jumlah Siswa (Student Count) - Required
- Latitude - Optional
- Longitude - Optional

**Status**: Pending

#### 3. SDM - Data Karyawan
**File**: `tests/test-suites/sdm-data-karyawan-crud.spec.js`
**Fields**:
- Nama (Name) - Required
- Email - Required
- No. Telepon (Phone) - Required
- Role - Required (Dropdown)
- Status - Required (Dropdown)

**Status**: Pending

#### 4. Keuangan - Aset Dapur
**File**: `tests/test-suites/keuangan-aset-dapur-crud.spec.js`
**Fields**:
- Nama Aset (Asset Name) - Required
- Kategori (Category) - Required (Dropdown)
- Nilai (Value) - Required
- Kondisi (Condition) - Required (Dropdown)
- Tanggal Pembelian (Purchase Date) - Optional

**Status**: Pending

#### 5. Menu Manajemen (Resep)
**File**: `tests/test-suites/menu-manajemen-crud.spec.js`
**Fields**:
- Nama Resep (Recipe Name) - Required
- Kategori (Category) - Required (Dropdown)
- Porsi (Portions) - Required
- Bahan (Ingredients) - Required (Multiple)
- Instruksi (Instructions) - Optional

**Status**: Pending

### Priority 2: Transaction Modules

#### 6. Supply Chain - Purchase Order
**File**: `tests/test-suites/supply-chain-purchase-order-crud.spec.js`
**Fields**:
- Supplier - Required (Dropdown)
- Tanggal PO (PO Date) - Required
- Items - Required (Multiple)
- Catatan (Notes) - Optional

**Status**: Pending

#### 7. Menu Perencanaan
**File**: `tests/test-suites/menu-perencanaan-crud.spec.js`
**Fields**:
- Tanggal (Date) - Required
- Menu - Required (Dropdown)
- Sekolah (School) - Required (Dropdown)
- Jumlah Porsi (Portions) - Required

**Status**: Pending

#### 8. Logistik - Tugas Pengiriman
**File**: `tests/test-suites/logistik-tugas-pengiriman-crud.spec.js`
**Fields**:
- Tanggal (Date) - Required
- Sekolah (School) - Required (Dropdown)
- Driver - Required (Dropdown)
- Menu - Required (Dropdown)
- Status - Required (Dropdown)

**Status**: Pending

#### 9. Keuangan - Arus Kas
**File**: `tests/test-suites/keuangan-arus-kas-crud.spec.js`
**Fields**:
- Tanggal (Date) - Required
- Tipe (Type) - Required (Dropdown: Pemasukan/Pengeluaran)
- Kategori (Category) - Required (Dropdown)
- Jumlah (Amount) - Required
- Deskripsi (Description) - Optional

**Status**: Pending

#### 10. Supply Chain - Penerimaan Barang (GRN)
**File**: `tests/test-suites/supply-chain-penerimaan-barang-crud.spec.js`
**Fields**:
- Purchase Order - Required (Dropdown)
- Tanggal Terima (Received Date) - Required
- Items - Required (Multiple with quantities)
- Catatan (Notes) - Optional

**Status**: Pending

### Priority 3: Inventory & Stock Modules

#### 11. Supply Chain - Bahan Baku
**File**: `tests/test-suites/supply-chain-bahan-baku-crud.spec.js`
**Operations**:
- Initialize stock
- Adjust stock levels
- View stock history

**Status**: Pending

#### 12. Menu Komponen (Barang Setengah Jadi)
**File**: `tests/test-suites/menu-komponen-crud.spec.js`
**Fields**:
- Nama Komponen (Component Name) - Required
- Kategori (Category) - Required
- Bahan (Ingredients) - Required (Multiple)
- Satuan (Unit) - Required

**Status**: Pending

---

## 🔍 Field Discovery Process

For each module, we need to:

1. **Navigate to the module**
2. **Click "Tambah" button**
3. **Inspect the form** to identify:
   - All input fields (text, number, date, etc.)
   - All dropdowns/selects
   - All textareas
   - Required vs optional fields (marked with *)
   - Field labels and placeholders
4. **Document field names and types**
5. **Create test data** for each field
6. **Implement CRUD tests**

---

## 📝 Test Data Strategy

### Naming Convention
Use timestamps to ensure uniqueness:
```javascript
const timestamp = Date.now();
const testName = `Test Item ${timestamp}`;
```

### Realistic Data
- Names: "Test Supplier 1234567890"
- Emails: "test1234567890@example.com"
- Phones: "081234567890"
- Addresses: "Jl. Test No. 123, Jakarta"
- Numbers: Use realistic values (e.g., 100 for student count)

### Cleanup Strategy
- Delete test data after each test
- Use unique names to avoid conflicts
- Handle cases where deletion might fail

---

## 🚀 Execution Plan

### Step 1: Implement CRUD Helper ✅
- Created `tests/utils/crud-helper.js`
- Tested with Supplier module

### Step 2: Create First CRUD Test Suite ✅
- Created `tests/test-suites/supply-chain-supplier-crud.spec.js`
- Implemented all 6 CRUD test cases

### Step 3: Run and Verify First Module
```bash
cd tests
npx playwright test supply-chain-supplier-crud --headed
```

### Step 4: Create CRUD Tests for Remaining Modules
- Use Supplier test as template
- Adapt for each module's specific fields
- Test one module at a time

### Step 5: Document Results
- Create detailed test results CSV
- Include screenshots of each operation
- Document any bugs found

### Step 6: Upload to Google Sheets
- Update upload script for CRUD results
- Create new sheet: "CRUD Test Results"
- Include all test data and screenshots

---

## 📊 Expected Results

### Coverage
- **Modules with CRUD**: 12 modules
- **Tests per module**: 6 tests (Create, Read, Update, Delete, Validation, Search)
- **Total CRUD tests**: 72 tests

### Quality Metrics
- **Pass Rate Target**: >90%
- **Bug Discovery**: Expected to find 5-10 bugs
- **Execution Time**: ~10-15 minutes per module

---

## 🐛 Common Issues to Watch For

### Form Issues
- Required fields not marked
- Validation not working
- Success messages not showing
- Forms not clearing after submit

### Data Issues
- Duplicate data not prevented
- Data not appearing in table after creation
- Data not updating after edit
- Data not removed after deletion

### UI Issues
- Buttons not clickable
- Modals not closing
- Dropdowns not working
- Search not filtering correctly

---

## 📸 Screenshot Strategy

Take screenshots at key points:
1. **Before**: Empty form
2. **Filled**: Form with test data
3. **After**: Success message or updated table
4. **Errors**: Validation errors or failures

---

## ✅ Success Criteria

A module's CRUD tests are complete when:
- ✅ All 6 test cases implemented
- ✅ Tests run successfully (>90% pass rate)
- ✅ Screenshots captured for all operations
- ✅ Results documented in CSV
- ✅ Results uploaded to Google Sheets
- ✅ Any bugs found are documented

---

## 🎯 Next Steps

1. **Run Supplier CRUD tests** to verify implementation
2. **Create CRUD tests for Logistik - Data Sekolah**
3. **Create CRUD tests for SDM - Data Karyawan**
4. **Continue with remaining Priority 1 modules**
5. **Document all results**
6. **Upload to Google Sheets**

---

**Status**: Ready to execute Supplier CRUD tests

**Command to run**:
```bash
cd tests
npx playwright test supply-chain-supplier-crud --headed
```

