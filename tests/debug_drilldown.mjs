import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch({ headless: false, slowMo: 300 });
  const context = await browser.newContext();
  const page = await context.newPage();

  // Listen to network
  page.on('response', async response => {
    const url = response.url();
    if (url.includes('/auth/') || url.includes('/dashboard/') || url.includes('/organizations/')) {
      console.log(`${response.request().method()} ${url} → ${response.status()}`);
    }
  });

  page.on('console', msg => {
    if (msg.type() === 'error' || msg.type() === 'warning') {
      console.log(`[BROWSER ${msg.type()}] ${msg.text()}`);
    }
  });

  // Login first
  await page.goto('http://localhost:5173/login');
  await page.evaluate(() => localStorage.clear());
  await page.reload();
  await page.waitForLoadState('networkidle');

  await page.locator('input').first().fill('SA001');
  await page.locator('input[type="password"]').fill('superadmin123');
  await page.locator('button:has-text("Masuk")').click();
  await page.waitForTimeout(2000);

  console.log('\n--- After login, URL:', page.url());

  // Navigate to Dashboard BGN
  await page.goto('http://localhost:5173/dashboard-bgn');
  await page.waitForTimeout(3000);
  console.log('--- Dashboard BGN loaded, URL:', page.url());

  // Check if Detail button exists
  const detailBtn = page.locator('button:has-text("Detail"), a:has-text("Detail")');
  const count = await detailBtn.count();
  console.log(`--- Found ${count} Detail button(s)`);

  if (count > 0) {
    // Check if button is clickable
    const box = await detailBtn.first().boundingBox();
    console.log('--- Button bounding box:', box);
    
    const isVisible = await detailBtn.first().isVisible();
    const isEnabled = await detailBtn.first().isEnabled();
    console.log(`--- Visible: ${isVisible}, Enabled: ${isEnabled}`);

    // Try clicking
    console.log('--- Clicking Detail button...');
    await detailBtn.first().click({ force: true });
    await page.waitForTimeout(2000);
    console.log('--- After click, URL:', page.url());

    // Check if alert/indicator appeared
    const alert = await page.locator('.ant-alert').isVisible().catch(() => false);
    console.log('--- Alert visible after click:', alert);

    // Check filter dropdown value
    const filterVal = await page.evaluate(() => {
      // Check if yayasan filter changed
      const selects = document.querySelectorAll('.ant-select-selection-item');
      return Array.from(selects).map(s => s.textContent);
    });
    console.log('--- Select values:', filterVal);
  }

  await page.screenshot({ path: 'tests/debug_drilldown_result.png', fullPage: true });
  console.log('\nScreenshot saved');
  
  await page.waitForTimeout(5000);
  await browser.close();
})();
