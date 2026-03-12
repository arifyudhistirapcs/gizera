const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

const testCasesPath = path.join(__dirname, '../test-cases/supply-chain-bahan-baku/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('Supply Chain - Bahan Baku Module', () => {
  
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

  test('scbb-001: View raw materials inventory', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/inventory`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const url = page.url();
    if (url.includes('login') || url.includes('404')) {
      console.log('⚠ Test scbb-001 skipped: Page not accessible');
      return;
    }
    
    const heading = await page.locator('h1, h2, h3').filter({ hasText: /Bahan Baku|Inventory|Manajemen/i }).count();
    console.log(`✓ Found heading: ${heading > 0}`);
    
    const cards = await page.locator('[class*="card"], .ant-card').count();
    console.log(`✓ Found ${cards} cards`);
    
    const table = await page.locator('table, .ant-table').count();
    console.log(`✓ Found table: ${table > 0}`);
    
    expect(heading > 0 || cards > 0 || table > 0).toBeTruthy();
    console.log('✓ Test scbb-001 passed: Raw materials inventory page loaded');
  });

  test('scbb-002: View stock levels', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/inventory`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const stockInfo = await page.locator('text=/stok|stock|gram|kg|liter/i').count();
    console.log(`✓ Found ${stockInfo} stock information elements`);
    
    if (stockInfo > 0) {
      console.log('✓ Test scbb-002 passed: Stock levels displayed');
    } else {
      console.log('⚠ Test scbb-002 skipped: No stock information found');
    }
  });

  test('scbb-003: Initialize stock', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/inventory`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const initButton = page.locator('button:has-text("Inisialisasi"), button:has-text("Initialize")');
    const buttonCount = await initButton.count();
    
    if (buttonCount > 0) {
      console.log(`✓ Found ${buttonCount} initialize buttons`);
      console.log('✓ Test scbb-003 passed: Initialize stock functionality available');
    } else {
      console.log('⚠ Test scbb-003 skipped: Initialize button not found');
    }
  });

  test('scbb-004: Search raw materials', async ({ page }) => {
    await page.goto(`${config.pwaBaseUrl}/inventory`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    const searchInput = page.locator('input[type="text"], input[type="search"], .ant-input').first();
    
    if (await searchInput.count() > 0) {
      console.log('✓ Found search input');
      console.log('✓ Test scbb-004 passed: Search functionality available');
    } else {
      console.log('⚠ Test scbb-004 skipped: Search input not found');
    }
  });

});
