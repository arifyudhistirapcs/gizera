#!/usr/bin/env node

/**
 * Complete Upload Script
 * Generates and uploads all test results to Google Sheets
 */

const { execSync } = require('child_process');

console.log('🚀 Starting complete upload process...\n');

try {
  // Step 1: Generate summary and module results
  console.log('📊 Step 1: Generating summary and module results...');
  execSync('node generate-all-results.js', { stdio: 'inherit', cwd: __dirname });
  console.log('✓ Summary generated\n');

  // Step 2: Generate detailed test results
  console.log('📋 Step 2: Generating detailed test results...');
  execSync('node generate-detailed-results.js', { stdio: 'inherit', cwd: __dirname });
  console.log('✓ Detailed results generated\n');

  // Step 3: Upload all results
  console.log('📤 Step 3: Uploading to Google Sheets...');
  execSync('node upload-all-results.js', { stdio: 'inherit', cwd: __dirname });
  console.log('✓ Summary and module results uploaded\n');

  // Step 4: Upload detailed results
  console.log('📤 Step 4: Uploading detailed results...');
  execSync('node upload-detailed-results.js', { stdio: 'inherit', cwd: __dirname });
  console.log('✓ Detailed results uploaded\n');

  console.log('✅ Complete upload process finished successfully!\n');
  console.log('📊 Your Google Sheets now contains:');
  console.log('   1. Test Summary - Overall statistics');
  console.log('   2. Module Results - Module breakdown');
  console.log('   3. Detailed Test Results - All test cases\n');
  
  console.log('🔗 View at: https://docs.google.com/spreadsheets/d/1UI329CBX5MnQ_-qfplE37JXAOFKmYVaA11HAb6ie_Pc/edit\n');

} catch (error) {
  console.error('\n❌ Upload process failed:', error.message);
  process.exit(1);
}

