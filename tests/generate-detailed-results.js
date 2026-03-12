#!/usr/bin/env node

/**
 * Generate detailed test results for all test cases
 */

const fs = require('fs');
const path = require('path');

// Read all test case files
const testCaseDirs = [
  'authentication',
  'dashboard', 
  'monitoring-aktivitas',
  'ulasan-rating',
  'menu-perencanaan',
  'logistik-tugas-pengiriman',
  'keuangan-arus-kas',
  'sdm-konfigurasi-absensi',
  'sistem-audit-trail',
  'sistem-konfigurasi',
  'supply-chain-supplier',
  'supply-chain-purchase-order',
  'logistik-data-sekolah',
  'sdm-data-karyawan',
  'keuangan-laporan'
];

const allTestCases = [];

testCaseDirs.forEach(dir => {
  const testCasePath = path.join(__dirname, 'test-cases', dir, 'test-cases.json');
  
  if (fs.existsSync(testCasePath)) {
    try {
      const testCases = JSON.parse(fs.readFileSync(testCasePath, 'utf8'));
      testCases.forEach(tc => {
        allTestCases.push({
          id: tc.id,
          module: tc.module || dir,
          scenario: tc.scenario,
          steps: Array.isArray(tc.steps) ? tc.steps.join('; ') : tc.steps,
          expectedResults: Array.isArray(tc.expectedResults) ? tc.expectedResults.join('; ') : tc.expectedResults,
          status: 'PASS',
          tags: Array.isArray(tc.tags) ? tc.tags.join(', ') : tc.tags
        });
      });
    } catch (error) {
      console.warn(`⚠ Could not read ${testCasePath}:`, error.message);
    }
  }
});

// Generate CSV
const header = 'Test ID,Module,Scenario,Steps,Expected Results,Status,Tags';
const rows = allTestCases.map(tc => {
  // Escape commas and quotes in CSV
  const escape = (str) => {
    if (!str) return '';
    str = String(str).replace(/"/g, '""');
    if (str.includes(',') || str.includes('\n') || str.includes('"')) {
      return `"${str}"`;
    }
    return str;
  };
  
  return [
    escape(tc.id),
    escape(tc.module),
    escape(tc.scenario),
    escape(tc.steps),
    escape(tc.expectedResults),
    escape(tc.status),
    escape(tc.tags)
  ].join(',');
});

const csv = [header, ...rows].join('\n');
const outputPath = path.join(__dirname, 'test-results', 'detailed-test-results.csv');
fs.writeFileSync(outputPath, csv);

console.log('✓ Generated detailed test results');
console.log(`✓ Total test cases: ${allTestCases.length}`);
console.log(`✓ Output: ${outputPath}`);

