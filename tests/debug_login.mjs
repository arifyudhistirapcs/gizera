import { chromium } from 'playwright';

(async () => {
  const browser = await chromium.launch({ headless: false, slowMo: 300 });
  const context = await browser.newContext();
  const page = await context.newPage();

  // Listen to network responses for auth endpoints
  page.on('response', async response => {
    const url = response.url();
    if (url.includes('/auth/login') || url.includes('/auth/me')) {
      console.log(`\n=== ${response.request().method()} ${url} ===`);
      console.log(`Status: ${response.status()}`);
      try {
        const body = await response.json();
        console.log('Response:', JSON.stringify(body, null, 2));
      } catch (e) {
        console.log('(non-JSON or already consumed)');
      }
    }
  });

  // Go to login page fresh
  await page.goto('http://localhost:5173/login');
  await page.evaluate(() => { localStorage.clear(); });
  await page.reload();
  await page.waitForLoadState('networkidle');

  // Fill login form using Ant Design Vue input selectors
  console.log('\n--- Filling login form ---');
  const identifierInput = page.locator('input').first();
  await identifierInput.fill('SA001');
  
  const passwordInput = page.locator('input[type="password"]');
  await passwordInput.fill('superadmin123');
  
  // Click login button
  console.log('--- Clicking login ---');
  await page.locator('button[type="submit"], button:has-text("Masuk")').first().click();
  
  // Wait for response
  await page.waitForTimeout(3000);
  
  // Check current URL
  const currentUrl = page.url();
  console.log(`\n--- Current URL: ${currentUrl} ---`);
  
  // Check localStorage
  const storageData = await page.evaluate(() => {
    return {
      token: localStorage.getItem('token') ? 'exists (' + localStorage.getItem('token').substring(0, 20) + '...)' : 'MISSING',
      user: localStorage.getItem('user')
    };
  });
  console.log('Token:', storageData.token);
  console.log('User data:', storageData.user);

  // Take screenshot
  await page.screenshot({ path: 'tests/debug_login_result.png', fullPage: true });
  console.log('\nScreenshot saved');

  await page.waitForTimeout(5000);
  await browser.close();
})();
