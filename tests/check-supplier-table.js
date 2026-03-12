const { chromium } = require('playwright');

async function checkSupplierTable() {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

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

    console.log('\n=== Analyzing Supplier Page ===\n');

    // Check for pagination
    const pagination = await page.locator('.ant-pagination, [class*="pagination"]').count();
    console.log(`Pagination elements: ${pagination}`);

    if (pagination > 0) {
      const paginationText = await page.locator('.ant-pagination').textContent();
      console.log(`Pagination text: ${paginationText}`);
    }

    // Check for filters
    const filters = await page.locator('.ant-select, select, input[type="search"]').count();
    console.log(`Filter elements: ${filters}`);

    // Check all tables
    const tables = await page.locator('table').all();
    console.log(`\nTotal tables: ${tables.length}`);

    for (let i = 0; i < tables.length; i++) {
      const table = tables[i];
      const rows = await table.locator('tbody tr').count();
      console.log(`\nTable ${i + 1}:`);
      console.log(`  Rows: ${rows}`);
      
      if (rows > 0) {
        const headers = await table.locator('thead th').allTextContents();
        console.log(`  Headers: ${headers.join(', ')}`);
        
        // Get first row content
        const firstRow = await table.locator('tbody tr').first().textContent();
        console.log(`  First row preview: ${firstRow.substring(0, 100)}...`);
      }
    }

    // Check for "Show All" or similar buttons
    const showAllButton = await page.locator('button:has-text("Tampilkan Semua"), button:has-text("Show All")').count();
    console.log(`\nShow All button: ${showAllButton > 0 ? 'Found' : 'Not found'}`);

    // Check for status filter
    const statusFilter = await page.locator('.ant-select:has-text("Status"), select:has-text("Status")').count();
    console.log(`Status filter: ${statusFilter > 0 ? 'Found' : 'Not found'}`);

    console.log('\n=== Analysis Complete ===');
    await page.waitForTimeout(10000);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await browser.close();
  }
}

checkSupplierTable();
