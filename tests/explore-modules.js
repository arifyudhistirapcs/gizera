#!/usr/bin/env node

/**
 * Module Explorer
 * 
 * Script untuk explore semua modul dan fitur-fiturnya
 * Akan menghasilkan report tentang elemen UI yang ditemukan di setiap halaman
 */

const { chromium } = require('playwright');
const fs = require('fs');
const path = require('path');

// Configuration
const BASE_URL = 'http://localhost:5173';
const CREDENTIALS = {
  username: 'kepala.sppg@sppg.com',
  password: 'password123'
};

// Modules to explore
const MODULES = [
  { name: 'Dashboard', path: '/dashboard' },
  { name: 'Monitoring Aktivitas', path: '/monitoring-activity' },
  { name: 'Ulasan & Rating', path: '/reviews' },
  { name: 'Display/KDS', path: '/kds' },
  { name: 'Menu Perencanaan', path: '/menu-planning' },
  { name: 'Menu Manajemen', path: '/menu-management' },
  { name: 'Menu Komponen', path: '/menu-components' },
  { name: 'Supply Chain - Supplier', path: '/suppliers' },
  { name: 'Supply Chain - Purchase Order', path: '/purchase-orders' },
  { name: 'Supply Chain - Penerimaan Barang', path: '/goods-receipt' },
  { name: 'Supply Chain - Bahan Baku', path: '/raw-materials' },
  { name: 'Logistik - Data Sekolah', path: '/schools' },
  { name: 'Logistik - Tugas Pengiriman', path: '/delivery-tasks' },
  { name: 'SDM - Data Karyawan', path: '/employees' },
  { name: 'SDM - Laporan Absensi', path: '/attendance-reports' },
  { name: 'SDM - Konfigurasi Absensi', path: '/attendance-config' },
  { name: 'Keuangan - Aset Dapur', path: '/kitchen-assets' },
  { name: 'Keuangan - Arus Kas', path: '/cash-flow' },
  { name: 'Keuangan - Laporan', path: '/financial-reports' },
  { name: 'Sistem - Audit Trail', path: '/audit-trail' },
  { name: 'Sistem - Konfigurasi', path: '/system-config' }
];

async function login(page) {
  console.log('🔐 Logging in...');
  
  await page.goto(BASE_URL);
  await page.waitForLoadState('networkidle');
  
  // Check if already logged in
  const isLoginPage = await page.locator('input[type="password"]').count() > 0;
  
  if (isLoginPage) {
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    await usernameInput.fill(CREDENTIALS.username);
    
    const passwordInput = page.locator('input[type="password"]').first();
    await passwordInput.fill(CREDENTIALS.password);
    
    const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
    await loginButton.click();
    
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    console.log('✓ Logged in successfully\n');
  } else {
    console.log('✓ Already logged in\n');
  }
}

async function exploreModule(page, module) {
  console.log(`\n${'='.repeat(80)}`);
  console.log(`📋 Exploring: ${module.name}`);
  console.log(`🔗 Path: ${module.path}`);
  console.log('='.repeat(80));
  
  const report = {
    name: module.name,
    path: module.path,
    accessible: false,
    elements: {
      headings: [],
      buttons: [],
      inputs: [],
      tables: false,
      forms: false,
      modals: false,
      filters: [],
      cards: 0,
      tabs: []
    },
    screenshots: []
  };
  
  try {
    // Navigate to module
    await page.goto(`${BASE_URL}${module.path}`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // Check if page is accessible
    const url = page.url();
    if (url.includes('login') || url.includes('404') || url.includes('403')) {
      console.log('❌ Page not accessible (redirected or error)');
      report.accessible = false;
      return report;
    }
    
    report.accessible = true;
    console.log('✓ Page accessible');
    
    // Take screenshot
    const screenshotPath = `screenshots/explore-${module.name.toLowerCase().replace(/[^a-z0-9]/g, '-')}.png`;
    await page.screenshot({ path: screenshotPath, fullPage: true });
    report.screenshots.push(screenshotPath);
    console.log(`📸 Screenshot saved: ${screenshotPath}`);
    
    // Find headings
    const headings = await page.locator('h1, h2, h3, [class*="title"], [class*="heading"]').allTextContents();
    report.elements.headings = headings.filter(h => h.trim().length > 0).slice(0, 5);
    console.log(`\n📝 Headings found: ${report.elements.headings.length}`);
    report.elements.headings.forEach(h => console.log(`   - ${h.trim()}`));
    
    // Find buttons
    const buttons = await page.locator('button').allTextContents();
    report.elements.buttons = buttons.filter(b => b.trim().length > 0 && b.trim().length < 50).slice(0, 10);
    console.log(`\n🔘 Buttons found: ${report.elements.buttons.length}`);
    report.elements.buttons.forEach(b => console.log(`   - ${b.trim()}`));
    
    // Find input fields
    const inputs = await page.locator('input[type="text"], input[type="email"], input[type="number"], input[type="date"], textarea').count();
    report.elements.inputs = inputs;
    console.log(`\n📝 Input fields: ${inputs}`);
    
    // Check for tables
    const tables = await page.locator('.ant-table, table').count();
    report.elements.tables = tables > 0;
    console.log(`📊 Tables: ${tables > 0 ? 'Yes' : 'No'} (${tables} found)`);
    
    // Check for forms
    const forms = await page.locator('form').count();
    report.elements.forms = forms > 0;
    console.log(`📋 Forms: ${forms > 0 ? 'Yes' : 'No'} (${forms} found)`);
    
    // Check for filters
    const filterElements = await page.locator('[class*="filter"], [placeholder*="filter" i], [placeholder*="cari" i], button:has-text("Filter")').allTextContents();
    report.elements.filters = filterElements.filter(f => f.trim().length > 0).slice(0, 5);
    console.log(`🔍 Filters found: ${report.elements.filters.length}`);
    report.elements.filters.forEach(f => console.log(`   - ${f.trim()}`));
    
    // Count cards
    const cards = await page.locator('[class*="card"], .ant-card, .ant-statistic').count();
    report.elements.cards = cards;
    console.log(`🃏 Cards: ${cards}`);
    
    // Find tabs
    const tabs = await page.locator('.ant-tabs-tab, [role="tab"]').allTextContents();
    report.elements.tabs = tabs.filter(t => t.trim().length > 0);
    console.log(`📑 Tabs found: ${report.elements.tabs.length}`);
    report.elements.tabs.forEach(t => console.log(`   - ${t.trim()}`));
    
    // Check for date pickers
    const datePickers = await page.locator('.ant-picker, input[type="date"]').count();
    console.log(`📅 Date pickers: ${datePickers}`);
    
    // Check for dropdowns/selects
    const selects = await page.locator('.ant-select, select').count();
    console.log(`📋 Dropdowns: ${selects}`);
    
    // Check for modals (if any are open)
    const modals = await page.locator('.ant-modal, .modal').count();
    report.elements.modals = modals > 0;
    console.log(`🪟 Modals: ${modals > 0 ? 'Yes' : 'No'} (${modals} found)`);
    
  } catch (error) {
    console.log(`❌ Error exploring module: ${error.message}`);
    report.error = error.message;
  }
  
  return report;
}

async function main() {
  console.log('🚀 Starting Module Explorer\n');
  
  const browser = await chromium.launch({ 
    headless: false,
    slowMo: 100 
  });
  
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 }
  });
  
  const page = await context.newPage();
  
  // Login first
  await login(page);
  
  // Explore all modules
  const reports = [];
  
  for (const module of MODULES) {
    const report = await exploreModule(page, module);
    reports.push(report);
    
    // Wait a bit between modules
    await page.waitForTimeout(1000);
  }
  
  // Save reports to JSON
  const reportPath = path.join(__dirname, 'exploration-report.json');
  fs.writeFileSync(reportPath, JSON.stringify(reports, null, 2));
  console.log(`\n\n✅ Exploration complete!`);
  console.log(`📄 Report saved to: ${reportPath}`);
  
  // Generate summary
  console.log(`\n${'='.repeat(80)}`);
  console.log('📊 SUMMARY');
  console.log('='.repeat(80));
  
  const accessible = reports.filter(r => r.accessible).length;
  const notAccessible = reports.filter(r => !r.accessible).length;
  
  console.log(`\n✓ Accessible modules: ${accessible}/${MODULES.length}`);
  console.log(`✗ Not accessible: ${notAccessible}/${MODULES.length}`);
  
  console.log('\n📋 Accessible Modules:');
  reports.filter(r => r.accessible).forEach(r => {
    console.log(`   ✓ ${r.name} (${r.path})`);
  });
  
  if (notAccessible > 0) {
    console.log('\n❌ Not Accessible Modules:');
    reports.filter(r => !r.accessible).forEach(r => {
      console.log(`   ✗ ${r.name} (${r.path})`);
    });
  }
  
  await browser.close();
}

main().catch(console.error);
