# Quick Start: Upload Test Results to Google Sheets

Panduan cepat untuk upload test results ke Google Sheets dalam 5 menit!

---

## ⚡ Quick Steps

### 1. Enable Google Sheets API (2 minutes)

1. Go to: https://console.cloud.google.com/
2. Create new project: `ERP SPPG Testing`
3. Enable Google Sheets API: https://console.cloud.google.com/apis/library/sheets.googleapis.com
4. Click "ENABLE"

### 2. Create Service Account (2 minutes)

1. Go to: https://console.cloud.google.com/apis/credentials
2. Click "CREATE CREDENTIALS" → "Service Account"
3. Name: `test-results-uploader`
4. Click "CREATE AND CONTINUE" → "DONE"
5. Click on the service account
6. Go to "KEYS" tab → "ADD KEY" → "Create new key" → "JSON"
7. Save downloaded file as `tests/config/google-credentials.json`

### 3. Share Google Sheets (1 minute)

1. Open downloaded JSON file
2. Copy the `client_email` value (looks like: `xxx@xxx.iam.gserviceaccount.com`)
3. Open your Google Sheets
4. Click "Share"
5. Paste the email
6. Set permission to "Editor"
7. **UNCHECK** "Notify people"
8. Click "Share"

### 4. Run Upload Script

```bash
cd tests
node upload-to-sheets.js
```

Done! ✅

---

## 📋 Checklist

Before running the script, make sure:

- [ ] Google Sheets API enabled in Google Cloud Console
- [ ] Service account created
- [ ] Credentials file saved as `tests/config/google-credentials.json`
- [ ] Google Sheets shared with service account email
- [ ] `googleapis` package installed (`npm install googleapis`)
- [ ] Test results CSV files exist in `tests/test-results/`

---

## 🎯 Expected Output

```
🚀 Starting Google Sheets upload...

✓ Google Sheets API initialized successfully

📤 Uploading test results to sheet: Authentication Test Results
✓ Read 8 rows from CSV
✓ Uploaded 8 rows to sheet: Authentication Test Results
✅ Successfully uploaded test results

📤 Uploading test summary to sheet: Test Summary
✓ Read 15 rows from CSV
✓ Uploaded 15 rows to sheet: Test Summary
✅ Successfully uploaded test summary

✅ All data uploaded successfully!

📊 View your results at:
https://docs.google.com/spreadsheets/d/YOUR_SPREADSHEET_ID/edit
```

---

## ❌ Common Errors

### "Failed to load sheets-config.json"
**Fix**: File already created at `tests/config/sheets-config.json` with your spreadsheet ID

### "Google credentials not found"
**Fix**: Save credentials as `tests/config/google-credentials.json`

### "The caller does not have permission"
**Fix**: Share Google Sheets with service account email (from credentials JSON)

### "API has not been used in project"
**Fix**: Enable Google Sheets API in Google Cloud Console

---

## 📚 Full Documentation

For detailed setup and troubleshooting, see: `tests/GOOGLE_SHEETS_SETUP.md`

---

## 🔒 Security Note

**IMPORTANT**: Add to `.gitignore`:
```
tests/config/google-credentials.json
```

Never commit credentials to git!

---

Happy testing! 🚀
