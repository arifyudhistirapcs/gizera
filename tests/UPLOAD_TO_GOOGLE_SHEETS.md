# Upload Test Results to Google Sheets

## Files Ready for Upload

Saya sudah membuatkan 2 file CSV yang siap diupload ke Google Sheets:

1. **Test Results Detail**: `tests/test-results/authentication-test-results.csv`
   - Berisi detail semua test cases dan hasilnya
   - Format: Test ID, Module, Scenario, Steps, Expected Results, Actual Results, Status, Execution Time, Tags, Notes

2. **Test Summary**: `tests/test-results/test-summary.csv`
   - Berisi summary metrics dari test run
   - Format: Metric, Value

## Cara Upload ke Google Sheets

### Option 1: Manual Upload (Recommended)

1. **Buka Google Sheets** yang sudah Anda share:
   https://docs.google.com/spreadsheets/d/1KwhjYaURuqgzg1lITtTeSCPLZNym7fTjxe9HsdcLf6c/edit?usp=sharing

2. **Upload Test Results Detail**:
   - Klik tab/sheet baru atau pilih sheet yang ada
   - Klik `File` → `Import`
   - Pilih `Upload` tab
   - Drag & drop file `tests/test-results/authentication-test-results.csv`
   - Pilih `Replace current sheet` atau `Insert new sheet(s)`
   - Klik `Import data`

3. **Upload Test Summary**:
   - Ulangi langkah yang sama untuk file `tests/test-results/test-summary.csv`
   - Atau paste di sheet yang berbeda

### Option 2: Copy-Paste

1. **Buka file CSV** dengan text editor atau Excel
2. **Copy semua content**
3. **Paste ke Google Sheets** di cell A1

### Option 3: Using Google Sheets API (Advanced)

Jika Anda ingin automated upload, saya bisa buatkan script Node.js yang menggunakan Google Sheets API. Untuk ini Anda perlu:
1. Enable Google Sheets API di Google Cloud Console
2. Create service account dan download credentials
3. Share Google Sheets dengan service account email

## Format Data

### Test Results Detail Columns:
- **Test ID**: Unique identifier (auth-001, auth-002, etc.)
- **Module**: Module name (Authentication)
- **Scenario**: Test scenario description
- **Steps**: Numbered steps to execute the test
- **Expected Results**: What should happen
- **Actual Results**: What actually happened
- **Status**: PASS/FAIL/SKIPPED
- **Execution Time**: Time taken to run the test
- **Tags**: Test categories (critical, smoke, high, etc.)
- **Notes**: Additional information or observations

### Test Summary Columns:
- **Metric**: Name of the metric
- **Value**: Value of the metric

## Formatting Tips

Setelah upload ke Google Sheets, Anda bisa:

1. **Format Headers**:
   - Bold the first row
   - Add background color (e.g., light blue)
   - Freeze the first row

2. **Format Status Column**:
   - PASS → Green background
   - FAIL → Red background
   - SKIPPED → Yellow background

3. **Add Conditional Formatting**:
   - Highlight rows based on status
   - Color code execution times (slow tests in red)

4. **Create Charts**:
   - Pie chart for Pass/Fail distribution
   - Bar chart for execution times
   - Timeline chart for test execution

## Sample Google Sheets Formula

### Calculate Pass Rate:
```
=COUNTIF(G:G,"PASS")/COUNTA(G:G)*100
```

### Average Execution Time:
```
=AVERAGE(H:H)
```

### Count by Status:
```
=COUNTIF(G:G,"PASS")
=COUNTIF(G:G,"FAIL")
=COUNTIF(G:G,"SKIPPED")
```

## Next Steps

After uploading:
1. ✅ Verify all data imported correctly
2. ✅ Apply formatting for better readability
3. ✅ Create charts/graphs for visualization
4. ✅ Share with team members
5. ✅ Set up automated updates (optional)

## Need Help?

Jika Anda ingin saya buatkan:
- Script untuk automated upload
- Custom formatting
- Additional reports
- Integration dengan CI/CD

Silakan beritahu saya!
