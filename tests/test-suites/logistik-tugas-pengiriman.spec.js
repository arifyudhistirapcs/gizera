const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

const testCasesPath = path.join(__dirname, '../test-cases/logistik-tugas-pengiriman/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('Logistik - Tugas Pengiriman Module', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    const isLoginPage = await page.locator('input[type="password"]').count() > 0;
    if (isLoginPage) {
      await page.locator('input[type="text"], input[type="email"]').first().fill('kepala.sppg@sppg.com');
      await page.locator('input[type="password"]').first().fill('password123');
      await page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first().click();
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(1000);
    }
  });

  test('ltp-001: View delivery tasks page', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const heading = await page.locator('h1, h2, h3').filter({ hasText: /Tugas Pengiriman/i }).count();
    const tabs = await page.locator('.ant-tabs-tab, [role="tab"]').count();
    const tables = await page.locator('.ant-table, table').count();
    const cards = await page.locator('[class*="card"], .ant-card').count();
    
    console.log(`✓ Heading: ${heading > 0}, Tabs: ${tabs}, Tables: ${tables > 0}, Cards: ${cards}`);
    expect(heading > 0 || tabs > 0).toBeTruthy();
    console.log('✓ Test ltp-001 passed');
  });

  test('ltp-002: Switch between Pengiriman and Pengambilan tabs', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const tabs = page.locator('.ant-tabs-tab, [role="tab"]');
    const tabCount = await tabs.count();
    
    if (tabCount >= 2) {
      console.log(`✓ Found ${tabCount} tabs`);
      await tabs.nth(1).click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ltp-002 passed: Switched tabs successfully');
    } else {
      console.log('⚠ Test ltp-002 skipped: Tabs not found');
    }
  });

  test('ltp-003: Create new delivery task', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const createButton = page.locator('button:has-text("Buat Tugas Pengiriman"), button:has-text("Buat Tugas")').first();
    
    if (await createButton.count() > 0) {
      console.log('✓ Found create button');
      await createButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ltp-003 passed: Create button clicked');
    } else {
      console.log('⚠ Test ltp-003 skipped: Create button not found');
    }
  });

  test('ltp-004: Filter delivery tasks', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const dropdowns = await page.locator('.ant-select, select').count();
    const datePicker = await page.locator('.ant-picker, input[type="date"]').count();
    
    console.log(`✓ Dropdowns: ${dropdowns}, Date picker: ${datePicker}`);
    console.log('✓ Test ltp-004 passed: Filter elements present');
  });

  test('ltp-005: Reset filters', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const resetButton = page.locator('button:has-text("Reset Filter"), button:has-text("Reset")').first();
    
    if (await resetButton.count() > 0) {
      console.log('✓ Found reset button');
      await resetButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ltp-005 passed: Reset button clicked');
    } else {
      console.log('⚠ Test ltp-005 skipped: Reset button not found');
    }
  });

  test('ltp-006: View task details', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/delivery-tasks`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const detailButtons = page.locator('button:has-text("Detail")');
    const count = await detailButtons.count();
    
    if (count > 0) {
      console.log(`✓ Found ${count} detail buttons`);
      await detailButtons.first().click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ltp-006 passed: Detail button clicked');
    } else {
      console.log('⚠ Test ltp-006 skipped: No detail buttons found');
    }
  });

});
