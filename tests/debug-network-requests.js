const { chromium } = require('playwright');
const CRUDHelper = require('./utils/crud-helper');

async function debugNetworkRequests() {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();
  const crudHelper = new CRUDHelper(page);

  // Listen to all network requests
  page.on('request', request => {
    if (request.url().includes('/suppliers')) {
      console.log(`>> ${request.method()} ${request.url()}`);
      if (request.method() === 'POST') {
        console.log(`   Body: ${request.postData()}`);
      }
    }
  });

  page.on('response', async response => {
    if (response.url().includes('/suppliers')) {
      console.log(`<< ${response.status()} ${response.url()}`);
      try {
        const body = await response.json();
        console.log(`   Response: ${JSON.stringify(body, null, 2)}`);
      } catch (e) {
        // Not JSON
      }
    }
  });

  try {
    // Login
    await page.goto('http://localhost:5173');
    await page.waitForLoadState('networkidle');
    
    const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
    await usernameInput.fill('kepala.sppg@sppg.com');
    
    const passwordInput = page.locator('input[type="password"]').first();
    await passwordInput.fill('password123');
    
    const loginButton = page.locator('button:has-text("Login"), button:has-text("Masuk"), button[type="submit"]').first();
    await loginButton.click();
    
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    // Navigate to suppliers
    console.log('\n=== Navigating to Suppliers Page ===\n');
    await page.goto('http://localhost:5173/suppliers');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const testName = `Network Test ${Date.now()}`;
    console.log(`\n=== Creating supplier: ${testName} ===\n`);

    // Create supplier
    const tambahButton = page.locator('button:has-text("Tambah Supplier")').first();
    await tambahButton.click();
    await page.waitForTimeout(1000);
    
    await crudHelper.fillInput('form_item_name', testName);
    await crudHelper.fillInput('form_item_product_category', 'Sayuran');
    await crudHelper.fillInput('form_item_contact_person', 'Network Test Contact');
    await crudHelper.fillInput('form_item_phone_number', '081234567890');
    await crudHelper.fillInput('form_item_email', `networktest${Date.now()}@test.com`);
    await crudHelper.fillTextarea('form_item_address', 'Jl. Network Test No. 123');
    
    console.log('\n=== Submitting Form ===\n');
    await crudHelper.clickSubmit('OK');
    await page.waitForTimeout(3000);

    console.log('\n=== Checking Table ===\n');
    
    // Check if data is in table
    const tableContent = await page.locator('table tbody').textContent();
    console.log(`Table contains "${testName}": ${tableContent.includes(testName)}`);
    
    // Try search
    console.log('\n=== Trying Search ===\n');
    const searchInput = page.locator('input[placeholder*="Cari"]').first();
    await searchInput.fill(testName);
    await page.waitForTimeout(2000);
    
    const tableAfterSearch = await page.locator('table tbody').textContent();
    console.log(`Table after search contains "${testName}": ${tableAfterSearch.includes(testName)}`);
    
    console.log('\n=== Debug Complete ===');
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await browser.close();
  }
}

debugNetworkRequests();
