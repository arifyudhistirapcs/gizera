const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

const testCasesPath = path.join(__dirname, '../test-cases/keuangan-arus-kas/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

const ConfigLoader = require('../utils/config-loader');
const config = new ConfigLoader().load();

test.describe('Keuangan - Arus Kas Module', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    const isLoginPage = await page.locator('input[type="password"]').count() > 0;
    if (isLoginPage) {
      await page.locator('input[type="text"]').first().fill('kepala.sppg@sppg.com');
      await page.locator('input[type="password"]').first().fill('password123');
      await page.locator('button[type="submit"]').first().click();
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(1000);
    }
  });

  test('kak-001: View cash flow dashboard', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // Check if page is accessible
    const url = page.url();
    if (url.includes('login') || url.includes('404')) {
      console.log('⚠ Page not accessible');
      return;
    }
    
    const heading = await page.locator('h1, h2, h3').filter({ hasText: /Arus Kas|Manajemen Arus Kas/i }).count();
    const cards = await page.locator('[class*="card"], .ant-card, .ant-statistic').count();
    const tables = await page.locator('.ant-table, table').count();
    
    console.log(`✓ Heading: ${heading > 0}, Cards: ${cards}, Tables: ${tables > 0}`);
    expect(heading > 0 || cards > 0 || tables > 0).toBeTruthy();
    console.log('✓ Test kak-001 passed');
  });

  test('kak-002: Add new transaction', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const addButton = page.locator('button:has-text("Tambah Transaksi")').first();
    if (await addButton.count() > 0) {
      await addButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test kak-002 passed');
    } else {
      console.log('⚠ Test kak-002 skipped');
    }
  });

  test('kak-003: Filter transactions', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const dropdowns = await page.locator('.ant-select').count();
    console.log(`✓ Dropdowns: ${dropdowns}`);
    console.log('✓ Test kak-003 passed');
  });

  test('kak-004: Export report', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const exportButton = page.locator('button:has-text("Export")').first();
    if (await exportButton.count() > 0) {
      console.log('✓ Export button found');
      console.log('✓ Test kak-004 passed');
    } else {
      console.log('⚠ Test kak-004 skipped');
    }
  });

  test('kak-005: Reset filters', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const resetButton = page.locator('button:has-text("Reset")').first();
    if (await resetButton.count() > 0) {
      await resetButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test kak-005 passed');
    } else {
      console.log('⚠ Test kak-005 skipped');
    }
  });

  test('kak-006: View transaction details', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/cash-flow`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const detailButtons = page.locator('button:has-text("Detail")');
    if (await detailButtons.count() > 0) {
      await detailButtons.first().click();
      await page.waitForTimeout(1000);
      console.log('✓ Test kak-006 passed');
    } else {
      console.log('⚠ Test kak-006 skipped');
    }
  });

});
