# Google Sheets API Setup Guide

Panduan lengkap untuk setup automated upload test results ke Google Sheets menggunakan Google Sheets API.

---

## 📋 Prerequisites

- Node.js installed
- Google Account
- Access to Google Cloud Console
- Google Sheets yang sudah dibuat

---

## 🔧 Setup Steps

### Step 1: Enable Google Sheets API

1. **Buka Google Cloud Console**
   - Go to: https://console.cloud.google.com/

2. **Create New Project** (atau gunakan existing project)
   - Click "Select a project" di top bar
   - Click "NEW PROJECT"
   - Project name: `ERP SPPG Testing` (atau nama lain)
   - Click "CREATE"

3. **Enable Google Sheets API**
   - Go to: https://console.cloud.google.com/apis/library
   - Search for "Google Sheets API"
   - Click "Google Sheets API"
   - Click "ENABLE"

### Step 2: Create Service Account

1. **Go to Credentials Page**
   - Go to: https://console.cloud.google.com/apis/credentials
   - Click "CREATE CREDENTIALS"
   - Select "Service Account"

2. **Fill Service Account Details**
   - Service account name: `test-results-uploader`
   - Service account ID: `test-results-uploader` (auto-generated)
   - Description: `Service account for uploading test results to Google Sheets`
   - Click "CREATE AND CONTINUE"

3. **Grant Permissions** (Optional)
   - Skip this step (click "CONTINUE")
   - Skip "Grant users access" (click "DONE")

4. **Create Key**
   - Click on the service account you just created
   - Go to "KEYS" tab
   - Click "ADD KEY" → "Create new key"
   - Select "JSON"
   - Click "CREATE"
   - **File will download automatically** → Save as `google-credentials.json`

### Step 3: Share Google Sheets with Service Account

1. **Copy Service Account Email**
   - From the downloaded JSON file, copy the `client_email` value
   - Example: `test-results-uploader@erp-sppg-testing.iam.gserviceaccount.com`

2. **Share Google Sheets**
   - Open your Google Sheets: https://docs.google.com/spreadsheets/d/1KwhjYaURuqgzg1lITtTeSCPLZNym7fTjxe9HsdcLf6c/edit
   - Click "Share" button (top right)
   - Paste the service account email
   - Set permission to "Editor"
   - **UNCHECK** "Notify people"
   - Click "Share"

### Step 4: Install Dependencies

```bash
cd tests
npm install googleapis
```

### Step 5: Configure Credentials

1. **Move credentials file**
   ```bash
   mv ~/Downloads/google-credentials.json tests/config/google-credentials.json
   ```

2. **Verify config file**
   - File `tests/config/sheets-config.json` sudah dibuat dengan spreadsheet ID Anda
   - Jika perlu edit, buka file dan update:
   ```json
   {
     "spreadsheetId": "1KwhjYaURuqgzg1lITtTeSCPLZNym7fTjxe9HsdcLf6c",
     "testResultsSheetName": "Authentication Test Results",
     "testSummarySheetName": "Test Summary"
   }
   ```

### Step 6: Test Upload

```bash
cd tests
node upload-to-sheets.js
```

Expected output:
```
🚀 Starting Google Sheets upload...

✓ Google Sheets API initialized successfully

📤 Uploading test results to sheet: Authentication Test Results
✓ Read 8 rows from CSV
✓ Created sheet: Authentication Test Results
✓ Cleared sheet: Authentication Test Results
✓ Uploaded 8 rows to sheet: Authentication Test Results
✓ Formatted headers for sheet: Authentication Test Results
✓ Formatted status column for sheet: Authentication Test Results
✅ Successfully uploaded test results to: Authentication Test Results

📤 Uploading test summary to sheet: Test Summary
✓ Read 15 rows from CSV
✓ Created sheet: Test Summary
✓ Cleared sheet: Test Summary
✓ Uploaded 15 rows to sheet: Test Summary
✓ Formatted headers for sheet: Test Summary
✅ Successfully uploaded test summary to: Test Summary

✅ All data uploaded successfully!

📊 View your results at:
https://docs.google.com/spreadsheets/d/1KwhjYaURuqgzg1lITtTeSCPLZNym7fTjxe9HsdcLf6c/edit
```

---

## 📁 File Structure

```
tests/
├── config/
│   ├── google-credentials.json    # Service account credentials (DO NOT COMMIT!)
│   └── sheets-config.json          # Spreadsheet configuration
├── utils/
│   └── google-sheets-uploader.js   # Uploader class
├── test-results/
│   ├── authentication-test-results.csv
│   └── test-summary.csv
└── upload-to-sheets.js             # Main upload script
```

---

## 🔒 Security Best Practices

### 1. Add to .gitignore

**IMPORTANT**: Never commit credentials to git!

Add to `.gitignore`:
```
tests/config/google-credentials.json
```

### 2. Restrict Service Account Permissions

- Only grant "Editor" access to specific spreadsheets
- Don't grant project-wide permissions
- Rotate credentials periodically

### 3. Environment Variables (Optional)

For CI/CD, use environment variables instead of files:

```bash
export GOOGLE_CREDENTIALS='{"type":"service_account",...}'
export SPREADSHEET_ID='1KwhjYaURuqgzg1lITtTeSCPLZNym7fTjxe9HsdcLf6c'
```

Update script to read from env:
```javascript
const credentials = process.env.GOOGLE_CREDENTIALS 
  ? JSON.parse(process.env.GOOGLE_CREDENTIALS)
  : JSON.parse(fs.readFileSync(CREDENTIALS_PATH, 'utf8'));
```

---

## 🎨 Formatting Features

The uploader automatically applies:

1. **Header Formatting**
   - Bold text
   - Blue background (#3399DD)
   - White text color
   - Frozen first row

2. **Status Column Formatting**
   - PASS → Green background (#B7E1B7)
   - FAIL → Red background (#F2B2B2)
   - SKIP → Yellow background (#FFF2B2)

3. **Auto-sizing**
   - Columns auto-adjust to content width

---

## 🔄 Automated Upload in CI/CD

### GitHub Actions Example

Create `.github/workflows/test-and-upload.yml`:

```yaml
name: Test and Upload Results

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
    
    - name: Install dependencies
      run: |
        cd tests
        npm install
    
    - name: Run tests
      run: |
        cd tests
        npx playwright test authentication --headed
    
    - name: Upload to Google Sheets
      env:
        GOOGLE_CREDENTIALS: ${{ secrets.GOOGLE_CREDENTIALS }}
        SPREADSHEET_ID: ${{ secrets.SPREADSHEET_ID }}
      run: |
        cd tests
        node upload-to-sheets.js
```

Add secrets in GitHub:
- `GOOGLE_CREDENTIALS`: Content of google-credentials.json
- `SPREADSHEET_ID`: Your spreadsheet ID

---

## 🐛 Troubleshooting

### Error: "Failed to initialize Google Sheets API"

**Cause**: Credentials file not found or invalid

**Solution**:
1. Check file exists: `tests/config/google-credentials.json`
2. Verify JSON format is valid
3. Re-download credentials from Google Cloud Console

### Error: "The caller does not have permission"

**Cause**: Service account doesn't have access to spreadsheet

**Solution**:
1. Copy service account email from credentials file
2. Share Google Sheets with that email
3. Grant "Editor" permission

### Error: "Spreadsheet not found"

**Cause**: Invalid spreadsheet ID

**Solution**:
1. Check spreadsheet ID in `sheets-config.json`
2. Verify spreadsheet exists and is accessible
3. Make sure URL format is correct

### Error: "API has not been used in project"

**Cause**: Google Sheets API not enabled

**Solution**:
1. Go to Google Cloud Console
2. Enable Google Sheets API
3. Wait a few minutes for propagation

---

## 📚 Additional Resources

- [Google Sheets API Documentation](https://developers.google.com/sheets/api)
- [Service Account Authentication](https://cloud.google.com/iam/docs/service-accounts)
- [googleapis npm package](https://www.npmjs.com/package/googleapis)

---

## 💡 Tips

1. **Test with a copy first**: Create a copy of your spreadsheet for testing
2. **Use descriptive sheet names**: Include date/time in sheet names for history
3. **Backup data**: Export sheets regularly as backup
4. **Monitor API usage**: Check quota usage in Google Cloud Console
5. **Version control**: Keep track of script changes in git

---

## 🎯 Next Steps

After successful setup:

1. ✅ Run automated upload after each test run
2. ✅ Create charts and dashboards in Google Sheets
3. ✅ Set up scheduled reports
4. ✅ Integrate with CI/CD pipeline
5. ✅ Share results with team members

---

## 📞 Need Help?

If you encounter issues:
1. Check this documentation
2. Review error messages carefully
3. Verify all setup steps completed
4. Check Google Cloud Console for API status
5. Test with a simple spreadsheet first

Happy testing! 🚀
