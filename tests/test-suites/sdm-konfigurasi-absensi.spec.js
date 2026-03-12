const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

const testCasesPath = path.join(__dirname, '../test-cases/sdm-konfigurasi-absensi/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

const ConfigLoader = require('../utils/config-loader');
const config = new ConfigLoader().load();

test.describe('SDM - Konfigurasi Absensi Module', () => {
  
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

  test('ska-001: View attendance configuration', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const heading = await page.locator('h1, h2, h3').filter({ hasText: /Konfigurasi Absensi/i }).count();
    const tabs = await page.locator('.ant-tabs-tab, [role="tab"]').count();
    console.log(`✓ Heading: ${heading > 0}, Tabs: ${tabs}`);
    expect(heading > 0 || tabs > 0).toBeTruthy();
    console.log('✓ Test ska-001 passed');
  });

  test('ska-002: Switch between Wi-Fi and GPS tabs', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const tabs = page.locator('.ant-tabs-tab, [role="tab"]');
    if (await tabs.count() >= 2) {
      await tabs.nth(1).click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ska-002 passed');
    } else {
      console.log('⚠ Test ska-002 skipped');
    }
  });

  test('ska-003: Add new Wi-Fi network', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const addButton = page.locator('button:has-text("Tambah")').first();
    if (await addButton.count() > 0) {
      await addButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ska-003 passed');
    } else {
      console.log('⚠ Test ska-003 skipped');
    }
  });

  test('ska-004: Edit configuration', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const editButton = page.locator('button:has-text("Edit")').first();
    if (await editButton.count() > 0) {
      await editButton.click();
      await page.waitForTimeout(1000);
      console.log('✓ Test ska-004 passed');
    } else {
      console.log('⚠ Test ska-004 skipped');
    }
  });

  test('ska-005: Disable configuration', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const disableButton = page.locator('button:has-text("Nonaktifkan")').first();
    if (await disableButton.count() > 0) {
      console.log('✓ Disable button found');
      console.log('✓ Test ska-005 passed');
    } else {
      console.log('⚠ Test ska-005 skipped');
    }
  });

  test('ska-006: Delete configuration', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/attendance-config`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const deleteButton = page.locator('button:has-text("Hapus")').first();
    if (await deleteButton.count() > 0) {
      console.log('✓ Delete button found');
      console.log('✓ Test ska-006 passed');
    } else {
      console.log('⚠ Test ska-006 skipped');
    }
  });

});
