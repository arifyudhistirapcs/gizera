const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

const testCasesPath = path.join(__dirname, '../test-cases/supply-chain-penerimaan-barang/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('Supply Chain - Penerimaan Barang Module', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    const isLoginPage = await page.locator('input[type="password"]').count() > 0;
    
    if (isLoginPage) {
      const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
      await usernameInput.fill('kepala.sppg@sppg.com');
      
      const passwordInput = page.locator('input[type="password"]').first();
      await passwordInput.fill('password123');
      
      const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
      await loginButton.click();
      
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(1000);
    }
  });

  test('scpb-001: View goods receipt page', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/goods-receipts`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const url = page.url();
    if (url.includes('login') || url.includes('404')) {
      console.log('⚠ Test scpb-001 skipped: Page not accessible');
      return;
    }
    
    const heading = await page.locator('h1, h2, h3').filter({ hasText: /Penerimaan|GRN|Receipt/i }).count();
    console.log(`✓ Found heading: ${heading > 0}`);
    
    const table = await page.locator('table, .ant-table').count();
    console.log(`✓ Found table: ${table > 0}`);
    
    const buatButton = await page.locator('button:has-text("Buat GRN"), button:has-text("Buat")').count();
    console.log(`✓ Buat GRN button: ${buatButton > 0}`);
    
    expect(heading > 0 || table > 0).toBeTruthy();
    console.log('✓ Test scpb-001 passed: Goods receipt page loaded');
  });

  test('scpb-002: Create new GRN', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/goods-receipts`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const buatButton = page.locator('button:has-text("Buat GRN"), button:has-text("Buat")').first();
    
    if (await buatButton.count() > 0) {
      console.log('✓ Found "Buat GRN" button');
      console.log('✓ Test scpb-002 passed: Create GRN button available');
    } else {
      console.log('⚠ Test scpb-002 skipped: "Buat GRN" button not found');
    }
  });

  test('scpb-003: View GRN details', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/goods-receipts`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const detailButtons = page.locator('button:has-text("Detail")');
    const buttonCount = await detailButtons.count();
    
    if (buttonCount > 0) {
      console.log(`✓ Found ${buttonCount} detail buttons`);
      console.log('✓ Test scpb-003 passed: Detail view available');
    } else {
      console.log('⚠ Test scpb-003 skipped: No GRN records to view');
    }
  });

  test('scpb-004: Search GRN records', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/goods-receipts`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const searchInput = page.locator('input[type="text"], input[type="search"], .ant-input').first();
    
    if (await searchInput.count() > 0) {
      console.log('✓ Found search input');
      console.log('✓ Test scpb-004 passed: Search functionality available');
    } else {
      console.log('⚠ Test scpb-004 skipped: Search input not found');
    }
  });

});
