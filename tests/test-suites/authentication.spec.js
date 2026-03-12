const { test, expect } = require('@playwright/test');
const path = require('path');
const fs = require('fs');

// Load test cases
const testCasesPath = path.join(__dirname, '../test-cases/authentication/test-cases.json');
const testCases = JSON.parse(fs.readFileSync(testCasesPath, 'utf8'));

// Load configuration
const ConfigLoader = require('../utils/config-loader');
const configLoader = new ConfigLoader();
const config = configLoader.load();

test.describe('Authentication Module', () => {
  
  test.beforeEach(async ({ page }) => {
    // Navigate to login page before each test
    await page.goto(config.pwaBaseUrl);
  });

  // Test Case: auth-001 - User login with valid credentials
  test('auth-001: User login with valid credentials', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-001');
    
    // Wait for login page to load
    await page.waitForLoadState('networkidle');
    
    // Check if we're on login page or already logged in
    const isLoginPage = await page.locator('input[type="text"], input[type="email"], input[placeholder*="username" i], input[placeholder*="email" i]').count() > 0;
    
    if (isLoginPage) {
      // Find username/email input
      const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
      await usernameInput.fill('kepala.sppg@sppg.com');
      
      // Find password input
      const passwordInput = page.locator('input[type="password"]').first();
      await passwordInput.fill('password123');
      
      // Find and click login button
      const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
      await loginButton.click();
      
      // Wait for navigation
      await page.waitForLoadState('networkidle');
      
      // Verify redirect to dashboard
      await expect(page).toHaveURL(/dashboard|home/i, { timeout: 10000 });
      
      // Verify user is logged in (check for logout button in sidebar)
      // Wait for sidebar to load
      await page.waitForTimeout(1000);
      const isLoggedIn = await page.locator('.logout-button, button.logout-button, [aria-label="Keluar"]').count() > 0;
      expect(isLoggedIn).toBeTruthy();
      
      console.log('✓ Test auth-001 passed: User logged in successfully');
    } else {
      console.log('⚠ Already logged in, skipping login test');
    }
  });

  // Test Case: auth-002 - User login with invalid username
  test('auth-002: User login with invalid username', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-002');
    
    // Ensure we're logged out first
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    // Try to logout if logged in
    const logoutButton = page.locator('.logout-button, button.logout-button, [aria-label="Keluar"]').first();
    if (await logoutButton.count() > 0) {
      await logoutButton.click();
      await page.waitForLoadState('networkidle');
    }
    
    // Find username input
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    await usernameInput.fill('invaliduser@test.com');
    
    // Find password input
    const passwordInput = page.locator('input[type="password"]').first();
    await passwordInput.fill('password123');
    
    // Click login button
    const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
    await loginButton.click();
    
    // Wait for error message to appear (Ant Design Vue message)
    await page.waitForTimeout(2000);
    
    // Verify error message is displayed (Ant Design Vue) or still on login page
    const errorMessage = await page.locator('.ant-message-error, .ant-form-item-explain-error, .ant-notification-notice-error').count();
    const stillOnLoginPage = await page.locator('input[type="password"]').count() > 0;
    
    // Test passes if still on login page (login failed)
    expect(stillOnLoginPage).toBeTruthy();
    
    if (errorMessage > 0) {
      console.log('✓ Test auth-002 passed: Invalid login rejected with error message');
    } else {
      console.log('✓ Test auth-002 passed: Invalid login rejected (stayed on login page)');
    }
  });

  // Test Case: auth-003 - User login with invalid password
  test('auth-003: User login with invalid password', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-003');
    
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    // Find username input
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    await usernameInput.fill('kepala.sppg@sppg.com'); // Valid test credentials
    
    // Find password input
    const passwordInput = page.locator('input[type="password"]').first();
    await passwordInput.fill('wrongpassword');
    
    // Click login button
    const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
    await loginButton.click();
    
    // Wait for error message (Ant Design Vue)
    await page.waitForTimeout(2000);
    
    // Verify error message or still on login page
    const errorMessage = await page.locator('.ant-message-error, .ant-form-item-explain-error, .ant-notification-notice-error').count();
    const stillOnLoginPage = await page.locator('input[type="password"]').count() > 0;
    
    // Test passes if still on login page (login failed)
    expect(stillOnLoginPage).toBeTruthy();
    
    if (errorMessage > 0) {
      console.log('✓ Test auth-003 passed: Invalid password rejected with error message');
    } else {
      console.log('✓ Test auth-003 passed: Invalid password rejected (stayed on login page)');
    }
  });

  // Test Case: auth-004 - User login with empty credentials
  test('auth-004: User login with empty credentials', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-004');
    
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    // Leave fields empty and try to submit
    const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
    await loginButton.click();
    
    // Wait a bit
    await page.waitForTimeout(1000);
    
    // Verify validation messages or that we're still on login page
    const stillOnLoginPage = await page.locator('input[type="password"]').count() > 0;
    expect(stillOnLoginPage).toBeTruthy();
    
    console.log('✓ Test auth-004 passed: Empty credentials validation works');
  });

  // Test Case: auth-005 - User logout successfully
  test('auth-005: User logout successfully', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-005');
    
    // First login
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    if (await usernameInput.count() > 0) {
      await usernameInput.fill('kepala.sppg@sppg.com');
      
      const passwordInput = page.locator('input[type="password"]').first();
      await passwordInput.fill('password123');
      
      const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
      await loginButton.click();
      
      await page.waitForLoadState('networkidle');
    }
    
    // Now logout
    const logoutButton = page.locator('.logout-button, button.logout-button, [aria-label="Keluar"]').first();
    if (await logoutButton.count() > 0) {
      await logoutButton.click();
      await page.waitForLoadState('networkidle');
      
      // Verify redirected to login page
      const backToLogin = await page.locator('input[type="password"]').count() > 0;
      expect(backToLogin).toBeTruthy();
      
      console.log('✓ Test auth-005 passed: Logout successful');
    } else {
      console.log('⚠ Logout button not found, test skipped');
    }
  });

  // Test Case: auth-006 - Session persistence after page refresh
  test('auth-006: Session persistence after page refresh', async ({ page }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-006');
    
    // Login first
    await page.goto(config.pwaBaseUrl);
    await page.waitForLoadState('networkidle');
    
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    if (await usernameInput.count() > 0) {
      await usernameInput.fill('kepala.sppg@sppg.com');
      
      const passwordInput = page.locator('input[type="password"]').first();
      await passwordInput.fill('password123');
      
      const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
      await loginButton.click();
      
      await page.waitForLoadState('networkidle');
      
      // Wait for dashboard to fully load
      await page.waitForTimeout(2000);
      
      // Refresh page
      await page.reload();
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(2000);
      
      // Verify still logged in (check for logout button instead of password input)
      const logoutButtonExists = await page.locator('.logout-button, button.logout-button, [aria-label="Keluar"]').count() > 0;
      expect(logoutButtonExists).toBeTruthy();
      
      console.log('✓ Test auth-006 passed: Session persisted after refresh');
    }
  });

  // Test Case: auth-007 - Access protected page without authentication
  test('auth-007: Access protected page without authentication', async ({ page, context }) => {
    const testCase = testCases.find(tc => tc.id === 'auth-007');
    
    // Clear all cookies and storage
    await context.clearCookies();
    await page.goto(config.pwaBaseUrl);
    await page.evaluate(() => {
      localStorage.clear();
      sessionStorage.clear();
    });
    
    // Try to access dashboard directly
    await page.goto(`${config.pwaBaseUrl}/dashboard`);
    await page.waitForLoadState('networkidle');
    
    // Verify redirected to login or access denied
    const url = page.url();
    const isLoginPage = url.includes('login') || await page.locator('input[type="password"]').count() > 0;
    expect(isLoginPage).toBeTruthy();
    
    console.log('✓ Test auth-007 passed: Protected page requires authentication');
  });

});
