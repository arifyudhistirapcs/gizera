const { chromium } = require('playwright');
const path = require('path');

async function exploreSupplierForm() {
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

    // Click Tambah Supplier
    const tambahButton = page.locator('button:has-text("Tambah Supplier")').first();
    await tambahButton.click();
    await page.waitForTimeout(2000);

    // Take screenshot
    await page.screenshot({ path: 'tests/screenshots/supplier-form-exploration.png', fullPage: true });

    // Explore form elements
    console.log('\n=== FORM EXPLORATION ===\n');

    // Find all inputs
    const inputs = await page.locator('input').all();
    console.log(`Found ${inputs.length} input elements:`);
    for (let i = 0; i < inputs.length; i++) {
      const input = inputs[i];
      const type = await input.getAttribute('type');
      const placeholder = await input.getAttribute('placeholder');
      const name = await input.getAttribute('name');
      const id = await input.getAttribute('id');
      const ariaLabel = await input.getAttribute('aria-label');
      
      console.log(`\nInput ${i + 1}:`);
      console.log(`  Type: ${type}`);
      console.log(`  Placeholder: ${placeholder}`);
      console.log(`  Name: ${name}`);
      console.log(`  ID: ${id}`);
      console.log(`  Aria-label: ${ariaLabel}`);
    }

    // Find all textareas
    const textareas = await page.locator('textarea').all();
    console.log(`\n\nFound ${textareas.length} textarea elements:`);
    for (let i = 0; i < textareas.length; i++) {
      const textarea = textareas[i];
      const placeholder = await textarea.getAttribute('placeholder');
      const name = await textarea.getAttribute('name');
      const id = await textarea.getAttribute('id');
      
      console.log(`\nTextarea ${i + 1}:`);
      console.log(`  Placeholder: ${placeholder}`);
      console.log(`  Name: ${name}`);
      console.log(`  ID: ${id}`);
    }

    // Find all selects
    const selects = await page.locator('.ant-select, select').all();
    console.log(`\n\nFound ${selects.length} select/dropdown elements`);

    // Find all labels
    const labels = await page.locator('label').all();
    console.log(`\n\nFound ${labels.length} label elements:`);
    for (let i = 0; i < labels.length; i++) {
      const label = labels[i];
      const text = await label.textContent();
      const forAttr = await label.getAttribute('for');
      
      console.log(`\nLabel ${i + 1}:`);
      console.log(`  Text: ${text}`);
      console.log(`  For: ${forAttr}`);
    }

    // Find all buttons
    const buttons = await page.locator('button').all();
    console.log(`\n\nFound ${buttons.length} button elements:`);
    for (let i = 0; i < buttons.length; i++) {
      const button = buttons[i];
      const text = await button.textContent();
      const type = await button.getAttribute('type');
      
      console.log(`\nButton ${i + 1}:`);
      console.log(`  Text: ${text}`);
      console.log(`  Type: ${type}`);
    }

    console.log('\n\n=== EXPLORATION COMPLETE ===\n');
    console.log('Screenshot saved to: tests/screenshots/supplier-form-exploration.png');
    
    // Wait for user to inspect
    await page.waitForTimeout(10000);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    await browser.close();
  }
}

exploreSupplierForm();
