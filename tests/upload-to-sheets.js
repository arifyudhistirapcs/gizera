#!/usr/bin/env node

/**
 * Upload Test Results to Google Sheets
 * 
 * Usage:
 *   node upload-to-sheets.js
 * 
 * Requirements:
 *   - Google Sheets API credentials in tests/config/google-credentials.json
 *   - Spreadsheet ID in tests/config/sheets-config.json
 */

const path = require('path');
const GoogleSheetsUploader = require('./utils/google-sheets-uploader');

// Configuration
const CREDENTIALS_PATH = path.join(__dirname, 'config', 'google-credentials.json');
const CONFIG_PATH = path.join(__dirname, 'config', 'sheets-config.json');
const TEST_RESULTS_CSV = path.join(__dirname, 'test-results', 'monitoring-aktivitas-test-results-clean.csv');
const TEST_SUMMARY_CSV = path.join(__dirname, 'test-results', 'test-summary.csv');

async function main() {
  try {
    console.log('🚀 Starting Google Sheets upload...\n');

    // Load configuration
    const fs = require('fs');
    let config;
    
    try {
      config = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
    } catch (error) {
      console.error('❌ Failed to load sheets-config.json');
      console.error('Please create tests/config/sheets-config.json with:');
      console.error(JSON.stringify({
        spreadsheetId: 'YOUR_SPREADSHEET_ID',
        testResultsSheetName: 'Authentication Test Results',
        testSummarySheetName: 'Test Summary'
      }, null, 2));
      process.exit(1);
    }

    // Check if credentials exist
    if (!fs.existsSync(CREDENTIALS_PATH)) {
      console.error('❌ Google credentials not found!');
      console.error('Please create tests/config/google-credentials.json');
      console.error('See tests/GOOGLE_SHEETS_SETUP.md for instructions');
      process.exit(1);
    }

    // Check if CSV files exist
    if (!fs.existsSync(TEST_RESULTS_CSV)) {
      console.error('❌ Test results CSV not found:', TEST_RESULTS_CSV);
      process.exit(1);
    }

    if (!fs.existsSync(TEST_SUMMARY_CSV)) {
      console.error('❌ Test summary CSV not found:', TEST_SUMMARY_CSV);
      process.exit(1);
    }

    // Initialize uploader
    const uploader = new GoogleSheetsUploader(CREDENTIALS_PATH, config.spreadsheetId);
    await uploader.initialize();

    // Upload test results (append mode for multi-module support)
    await uploader.uploadTestResults(
      TEST_RESULTS_CSV,
      config.testResultsSheetName || 'Test Results',
      true // append mode
    );

    // Upload test summary
    await uploader.uploadTestSummary(
      TEST_SUMMARY_CSV,
      config.testSummarySheetName || 'Test Summary'
    );

    console.log('✅ All data uploaded successfully!');
    console.log(`\n📊 View your results at:`);
    console.log(`https://docs.google.com/spreadsheets/d/${config.spreadsheetId}/edit\n`);

  } catch (error) {
    console.error('\n❌ Upload failed:', error.message);
    console.error('\nFor help, see: tests/GOOGLE_SHEETS_SETUP.md');
    process.exit(1);
  }
}

// Run main function
main();
