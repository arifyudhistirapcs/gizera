#!/usr/bin/env node

/**
 * Upload Detailed Test Results to Google Sheets
 */

const path = require('path');
const fs = require('fs');
const GoogleSheetsUploader = require('./utils/google-sheets-uploader');

const CREDENTIALS_PATH = path.join(__dirname, 'config', 'google-credentials.json');
const CONFIG_PATH = path.join(__dirname, 'config', 'sheets-config.json');
const DETAILED_CSV = path.join(__dirname, 'test-results', 'detailed-test-results.csv');

async function main() {
  try {
    console.log('🚀 Uploading detailed test results...\n');

    const config = JSON.parse(fs.readFileSync(CONFIG_PATH, 'utf8'));

    if (!fs.existsSync(CREDENTIALS_PATH)) {
      console.error('❌ Google credentials not found!');
      process.exit(1);
    }

    if (!fs.existsSync(DETAILED_CSV)) {
      console.error('❌ Detailed results CSV not found. Run: node generate-detailed-results.js');
      process.exit(1);
    }

    console.log('📡 Connecting to Google Sheets...');
    const uploader = new GoogleSheetsUploader(CREDENTIALS_PATH, config.spreadsheetId);
    await uploader.initialize();
    console.log('✓ Connected\n');

    console.log('📋 Uploading detailed test results...');
    await uploader.uploadTestResults(
      DETAILED_CSV,
      'Detailed Test Results',
      false
    );
    console.log('✓ Detailed results uploaded\n');

    console.log('✅ Upload complete!');
    console.log(`\n📊 View at: https://docs.google.com/spreadsheets/d/${config.spreadsheetId}/edit\n`);

  } catch (error) {
    console.error('\n❌ Upload failed:', error.message);
    process.exit(1);
  }
}

main();

