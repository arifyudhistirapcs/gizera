const { chromium } = require('@playwright/test');
const ConfigLoader = require('./utils/config-loader');

const configLoader = new ConfigLoader();
const config = configLoader.load();

const modules = [
  { name: 'Logistik - Data Sekolah', url: '/schools', addButton: 'Tambah Sekolah' },
  { name: 'SDM - Data Karyawan', url: '/employees', addButton: 'Tambah Karyawan' },
  { name: 'Keuangan - Aset Dapur', url: '/kitchen-assets', addButton: 'Tambah Aset' },
  { name: 'Menu Manajemen', url: '/recipes', addButton: 'Tambah Resep' },
  { name: 'Supply Chain - Purchase Order', url: '/purchase-orders', addButton: 'Tambah PO' },
  { name: 'Menu Perencanaan', url: '/menu-planning', addButton: 'Tambah Menu' },
  { name: 'Logistik - Tugas Pengiriman', url: '/delivery-tasks', addButton: 'Tambah Tugas' },
  { name: 'Keuangan - Arus Kas', url: '/cash-flow', addButton: 'Tambah Transaksi' },
];

async function discoverModule(page, module) {
  console.log(`\n${'='.repeat(60)}`);
  console.log(`Discovering: ${module.name}`);
  console.log(`${'='.repeat(60)}`);
  
  try {
    await page.goto(`${config.pwaBaseUrl}${module.url}`);
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(2000);
    
    // Check if page loaded
    const pageTitle = await page.title();
    console.log(`Page title: ${pageTitle}`);
    
    // Look for add button
    const addButton = page.locator(`button:has-text("${module.addButton}")`).first();
    const buttonExists = await addButton.count() > 0;
    
    if (!buttonExists) {
      console.log(`⚠ Add button "${module.addButton}" not found`);
      console.log(`Trying alternative button texts...`);
      
      const altButtons = await page.locator('button').all();
      console.log(`Found ${altButtons.length} buttons on page`);
      
      for (let i = 0; i < Math.min(altButtons.length, 10); i++) {
        const text = await altButtons[i].textContent();
        if (text && text.includes('Tambah')) {
          console.log(`  - Button ${i + 1}: "${text}"`);
        }
      }
      return;
    }
    
    console.log(`✓ Found add button: "${module.addButton}"`);
    
    // Click add button
    await addButton.click();
    await page.waitForTimeout(1500);
    
    // Discover form fields
    console.log(`\nForm Fields:`);
    
    // Find all inputs
    const inputs = await page.locator('input[id^="form_item"]').all();
    console.log(`\nText Inputs (${inputs.length}):`);
    for (const input of inputs) {
      const id = await input.getAttribute('id');
      const type = await input.getAttribute('type');
      const placeholder = await input.getAttribute('placeholder');
      console.log(`  - ID: ${id}, Type: ${type}, Placeholder: ${placeholder || 'none'}`);
    }
    
    // Find all textareas
    const textareas = await page.locator('textarea[id^="form_item"]').all();
    console.log(`\nTextareas (${textareas.length}):`);
    for (const textarea of textareas) {
      const id = await textarea.getAttribute('id');
      const placeholder = await textarea.getAttribute('placeholder');
      console.log(`  - ID: ${id}, Placeholder: ${placeholder || 'none'}`);
    }
    
    // Find all selects/dropdowns
    const selects = await page.locator('.ant-select[id^="form_item"]').all();
    console.log(`\nDropdowns (${selects.length}):`);
    for (const select of selects) {
      const id = await select.getAttribute('id');
      console.log(`  - ID: ${id}`);
    }
    
    // Close modal
    await page.keyboard.press('Escape');
    await page.waitForTimeout(500);
    
  } catch (error) {
    console.log(`✗ Error: ${error.message}`);
  }
}

async function main() {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();
  
  // Login
  console.log('Logging in...');
  await page.goto(config.pwaBaseUrl);
  await page.waitForLoadState('networkidle');
  
  const usernameInput = page.locator('input[type="text"], input[type="email"]').first();
  await usernameInput.fill('kepala.sppg@sppg.com');
  
  const passwordInput = page.locator('input[type="password"]').first();
  await passwordInput.fill('password123');
  
  const loginButton = page.locator('button[type="submit"]').first();
  await loginButton.click();
  
  await page.waitForLoadState('networkidle');
  await page.waitForTimeout(2000);
  console.log('✓ Logged in successfully\n');
  
  // Discover each module
  for (const module of modules) {
    await discoverModule(page, module);
    await page.waitForTimeout(1000);
  }
  
  console.log(`\n${'='.repeat(60)}`);
  console.log('Discovery complete!');
  console.log(`${'='.repeat(60)}\n`);
  
  await browser.close();
}

main().catch(console.error);
