const { chromium } = require('playwright');
const CRUDHelper = require('./utils/crud-helper');

async function debugSupplierOperations() {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();
  const crudHelper = new CRUDHelper(page);

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
    await page.goto('http://localhost:5173/suppliers');
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const testName = `Debug Supplier ${Date.now()}`;
    console.log(`\n=== Creating supplier: ${testName} ===`);

    // Create supplier
    const tambahButton = page.locator('button:has-text("Tambah Supplier")').first();
    await tambahButton.click();
    await page.waitForTimeout(1000);
    
    await crudHelper.fillInput('form_item_name', testName);
    await crudHelper.fillInput('form_item_product_category', 'Sayuran');
    await crudHelper.fillInput('form_item_contact_person', 'Debug Contact');
    await crudHelper.fillInput('form_item_phone_number', '081234567890');
    await crudHelper.fillInput('form_item_email', `debug${Date.now()}@test.com`);
    await crudHelper.fillTextarea('form_item_address', 'Jl. Debug No. 123');
    
    await crudHelper.clickSubmit('OK');
    await page.waitForTimeout(3000);

    console.log('\n=== Testing Detail View ===');
    // Click Detail button
    const detailButton = page.locator(`table tr:has-text("${testName}") button:has-text("Detail")`).first();
    await detailButton.click();
    await page.waitForTimeout(2000);

    // Check what's in the modal
    const modals = await page.locator('.ant-modal').all();
    console.log(`Found ${modals.length} modals`);
    
    for (let i = 0; i < modals.length; i++) {
      const modal = modals[i];
      const isVisible = await modal.isVisible();
      console.log(`\nModal ${i + 1}:`);
      console.log(`  Visible: ${isVisible}`);
      
      if (isVisible) {
        const content = await modal.textContent();
        console.log(`  Content preview: ${content.substring(0, 200)}...`);
        console.log(`  Contains test name: ${content.includes(testName)}`);
      }
    }

    // Close modal
    await page.keyboard.press('Escape');
    await page.waitForTimeout(1000);

    console.log('\n=== Testing Edit ===');
    // Click Edit button
    const editButton = page.locator(`table tr:has-text("${testName}") button:has-text("Edit")`).first();
    await editButton.click();
    await page.waitForTimeout(2000);

    // Update name
    const updatedName = `${testName} (Updated)`;
    const nameInput = page.locator('#form_item_name');
    await nameInput.clear();
    await nameInput.fill(updatedName);
    
    await crudHelper.clickSubmit('OK');
    await page.waitForTimeout(3000);

    console.log('\n=== Checking if update worked ===');
    // Check all tables
    const tables = await page.locator('table tbody').all();
    console.log(`Found ${tables.length} tables with tbody`);
    
    for (let i = 0; i < tables.length; i++) {
      const table = tables[i];
      const content = await table.textContent();
      console.log(`\nTable ${i + 1}:`);
      console.log(`  Contains original name: ${content.includes(testName)}`);
      console.log(`  Contains updated name: ${content.includes(updatedName)}`);
      if (content.includes(testName) || content.includes(updatedName)) {
        console.log(`  Content preview: ${content.substring(0, 300)}...`);
      }
    }

    // Reload and check again
    await page.reload();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const tableContentAfterReload = await page.locator('table').textContent();
    console.log(`\nAfter reload:`);
    console.log(`Table contains original name: ${tableContentAfterReload.includes(testName)}`);
    console.log(`Table contains updated name: ${tableContentAfterReload.includes(updatedName)}`);

    console.log('\n=== Testing Delete ===');
    // Find the supplier (use whichever name is in the table)
    const nameToDelete = tableContentAfterReload.includes(updatedName) ? updatedName : testName;
    console.log(`Deleting: ${nameToDelete}`);

    const deleteButton = page.locator(`table tr:has-text("${nameToDelete}") button:has-text("Hapus")`).first();
    await deleteButton.click();
    await page.waitForTimeout(1000);

    // Confirm
    const confirmButton = page.locator('button:has-text("Ya"), button.ant-btn-primary').first();
    await confirmButton.click();
    await page.waitForTimeout(3000);

    console.log('\n=== Checking if delete worked ===');
    const tableAfterDelete = await page.locator('table').textContent();
    console.log(`Table still contains name: ${tableAfterDelete.includes(nameToDelete)}`);

    // Check if row has inactive status
    const row = page.locator(`table tr:has-text("${nameToDelete}")`);
    if (await row.count() > 0) {
      const rowText = await row.textContent();
      console.log(`Row text: ${rowText}`);
      console.log(`Contains "Tidak Aktif": ${rowText.includes('Tidak Aktif')}`);
      console.log(`Contains "Inactive": ${rowText.includes('Inactive')}`);
    }

    // Reload and check
    await page.reload();
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);

    const tableAfterReload2 = await page.locator('table').textContent();
    console.log(`\nAfter reload:`);
    console.log(`Table still contains name: ${tableAfterReload2.includes(nameToDelete)}`);

    console.log('\n=== Debug Complete ===');
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await browser.close();
  }
}

debugSupplierOperations();
