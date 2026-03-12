#!/usr/bin/env node

/**
 * Upload All Test Results to Google Sheets
 * Includes Priority 1 and Priority 2 modules
 */

const path = require('path');
const fs = require('fs');
const GoogleSheetsUploader = require('./utils/google-sheets-uploader');

// Configuration
const CREDENTIALS_PATH = path.join(__dirname, 'config', 'google-credentials.json');
const CONFIG_PATH = path.join(__dirname, 'config', 'sheets-config.json');
const SUMMARY_CSV = path.join(__dirname, 'test-results', 'test-summary-all.csv');
const MODULE_RESULTS_CSV = path.join(__dirname, 'test-results', 'module-results-all.csv');

async function main() {
  try {
    console.log('🚀 Starting comprehensive Google Sheets upload...\n');

    // Load configuration
    let config;
    try {
      config = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));
    } catch (error) {
      console.error('❌ Failed to load sheets-config.json');
      process.exit(1);
    }

    // Check credentials
    if (!fs.existsSync(CREDENTIALS_PATH)) {
      console.error('❌ Google credentials not found!');
      console.error('Please create tests/config/google-credentials.json');
      process.exit(1);
    }

    // Check CSV files
    if (!fs.existsSync(SUMMARY_CSV)) {
      console.error('❌ Summary CSV not found. Run: node generate-all-results.js');
      process.exit(1);
    }

    if (!fs.existsSync(MODULE_RESULTS_CSV)) {
      console.error('❌ Module results CSV not found. Run: node generate-all-results.js');
      process.exit(1);
    }

    // Initialize uploader
    console.log('📡 Connecting to Google Sheets...');
    const uploader = new GoogleSheetsUploader(CREDENTIALS_PATH, config.spreadsheetId);
    await uploader.initialize();
    console.log('✓ Connected successfully\n');

    // Upload test summary
    console.log('📊 Uploading test summary...');
    await uploader.uploadTestSummary(
      SUMMARY_CSV,
      'Test Summary'
    );
    console.log('✓ Test summary uploaded\n');

    // Upload module results
    console.log('📋 Uploading module results...');
    await uploader.uploadTestResults(
      MODULE_RESULTS_CSV,
      'Module Results',
      false // clear and replace
    );
    console.log('✓ Module results uploaded\n');

    console.log('✅ All data uploaded successfully!');
    console.log(`\n📊 View your results at:`);
    console.log(`https://docs.google.com/spreadsheets/d/${config.spreadsheetId}/edit\n`);
    
    console.log('📈 Summary:');
    console.log('- Test Summary sheet: Overall statistics');
    console.log('- Module Results sheet: Detailed module breakdown');

  } catch (error) {
    console.error('\n❌ Upload failed:', error.message);
    if (error.stack) {
      console.error('\nStack trace:', error.stack);
    }
    process.exit(1);
  }
}

// Run main function
main();

