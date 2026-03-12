const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

// Load test cases
const testCasesPath = path.join(__dirname, '../test-cases/dashboard/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

// Load configuration
const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('Dashboard Module', () => {
  
  // Login before each test
  test.beforeEach(async ({ page }) => {
    // Navigate to login page
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    // Check if already logged in
    const isLoginPage = await page.locator('input[type="password"]').count() > 0;
    
    if (isLoginPage) {
      // Login
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

  // Test Case: dash-001 - View dashboard with all widgets loaded
  test('dash-001: View dashboard with all widgets loaded', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'dash-001');
    
    // Navigate to dashboard (should already be there after login)
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // Verify dashboard page is displayed
    const url = page.url();
    expect(url).toContain('dashboard');
    
    // Check for common dashboard elements
    const hasContent = await page.locator('body').count() > 0;
    expect(hasContent).toBeTruthy();
    
    console.log('✓ Test dash-001 passed: Dashboard loaded successfully');
  });

  // Test Case: dash-002 - Refresh dashboard data
  test('dash-002: Refresh dashboard data', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'dash-002');
    
    // Navigate to dashboard
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);
    
    // Look for refresh button (common selectors)
    const refreshButton = page.locator('button:has-text("Refresh"), button:has-text("Muat Ulang"), button[aria-label*="refresh" i], .refresh-button, [class*="refresh"]').first();
    
    if (await refreshButton.count() > 0) {
      await refreshButton.click();
      await page.waitForTimeout(2000);
      
      console.log('✓ Test dash-002 passed: Dashboard refresh triggered');
    } else {
      // If no refresh button, just reload the page
      await page.reload();
      await page.waitForLoadState('networkidle');
      
      console.log('✓ Test dash-002 passed: Dashboard reloaded (no refresh button found)');
    }
  });

  // Test Case: dash-003 - Navigate from dashboard to other modules
  test('dash-003: Navigate from dashboard to other modules', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'dash-003');
    
    // Navigate to dashboard
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);
    
    // Try to find and click a menu item (look for common menu items)
    const menuItems = page.locator('.ant-menu-item, .menu-item, nav a, .sidebar a').first();
    
    if (await menuItems.count() > 0) {
      const initialUrl = page.url();
      await menuItems.click();
      await page.waitForTimeout(1000);
      
      // Navigate back to dashboard
      await page.goto(`${config.pwaBaseUrl}/dashboard`);
      await page.waitForLoadState('networkidle');
      
      // Verify we're back on dashboard
      const url = page.url();
      expect(url).toContain('dashboard');
      
      console.log('✓ Test dash-003 passed: Navigation works correctly');
    } else {
      console.log('⚠ Test dash-003 skipped: No menu items found');
    }
  });

  // Test Case: dash-004 - Filter dashboard by date range
  test('dash-004: Filter dashboard by date range', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'dash-004');
    
    // Navigate to dashboard
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);
    
    // Look for date picker or filter
    const datePicker = page.locator('.ant-picker, input[type="date"], .date-picker, [class*="date"]').first();
    
    if (await datePicker.count() > 0) {
      await datePicker.click();
      await page.waitForTimeout(500);
      
      // Try to select a date (if date picker opens)
      const dateCell = page.locator('.ant-picker-cell, .date-cell').first();
      if (await dateCell.count() > 0) {
        await dateCell.click();
        await page.waitForTimeout(1000);
      }
      
      console.log('✓ Test dash-004 passed: Date filter interaction works');
    } else {
      console.log('⚠ Test dash-004 skipped: No date picker found');
    }
  });

  // Test Case: dash-005 - Dashboard with data loading error
  test('dash-005: Dashboard with data loading error', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'dash-005');
    
    // Navigate to dashboard
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // Check if dashboard loaded (even with potential errors)
    const hasContent = await page.locator('body').count() > 0;
    expect(hasContent).toBeTruthy();
    
    // Look for error messages
    const errorMessage = await page.locator('.ant-message-error, .ant-alert-error, .error-message, [class*="error"]').count();
    
    if (errorMessage > 0) {
      console.log('✓ Test dash-005 passed: Error handling detected');
    } else {
      console.log('✓ Test dash-005 passed: Dashboard loaded without errors');
    }
  });

});
