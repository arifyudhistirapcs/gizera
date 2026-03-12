const { chromium } = require('playwright');

async function checkSearchInput() {
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

    console.log('\n=== Checking Search Input ===\n');

    // Find all inputs
    const inputs = await page.locator('input').all();
    console.log(`Found ${inputs.length} inputs`);

    for (let i = 0; i < inputs.length; i++) {
      const input = inputs[i];
      const placeholder = await input.getAttribute('placeholder');
      const type = await input.getAttribute('type');
      const className = await input.getAttribute('class');
      
      console.log(`\nInput ${i + 1}:`);
      console.log(`  Type: ${type}`);
      console.log(`  Placeholder: ${placeholder}`);
      console.log(`  Class: ${className}`);
    }

    // Try to find the search input and use it
    console.log('\n=== Testing Search ===\n');
    
    const searchInput = page.locator('input[placeholder*="Cari"]').first();
    if (await searchInput.count() > 0) {
      console.log('Found search input with "Cari" placeholder');
      
      // Type in search
      await searchInput.fill('API Test Supplier');
      await page.waitForTimeout(2000);
      
      // Check table content
      const tableContent = await page.locator('table tbody').textContent();
      console.log('\nTable content after search:');
      console.log(tableContent.substring(0, 500));
      
      // Clear search
      await searchInput.clear();
      await page.waitForTimeout(2000);
    } else {
      console.log('Search input with "Cari" placeholder NOT found');
    }

    console.log('\n=== Check Complete ===');
    await page.waitForTimeout(5000);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await browser.close();
  }
}

checkSearchInput();
