const { defineConfig, devices } = require('@playwright/test');
const ConfigLoader = require('../utils/config-loader');

// Load test configuration
const configLoader = new ConfigLoader();
const testConfig = configLoader.load();

module.exports = defineConfig({
  testDir: '../test-suites',
  
  // Maximum time one test can run
  timeout: testConfig.timeouts.default || 30000,
  
  // Test execution settings
  fullyParallel: false, // Run tests sequentially for better isolation
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // Single worker for headed mode observation
  
  // Reporter configuration
  reporter: [
    ['list'],
    ['json', { outputFile: '../test-results/test-results.json' }],
    ['html', { outputFolder: '../test-results/html-report', open: 'never' }]
  ],
  
  // Shared settings for all projects
  use: {
    // Base URL for navigation
    baseURL: testConfig.pwaBaseUrl,
    
    // Browser context options
    trace: 'on-first-retry',
    screenshot: testConfig.screenshots.onFailure ? 'only-on-failure' : 'off',
    video: 'retain-on-failure',
    
    // Navigation timeout
    navigationTimeout: testConfig.timeouts.navigation || 60000,
    
    // Action timeout
    actionTimeout: testConfig.timeouts.action || 10000,
  },

  // Configure projects for different browsers
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        headless: testConfig.browser.headless,
        slowMo: testConfig.browser.slowMo || 0,
        viewport: testConfig.browser.viewport || { width: 1920, height: 1080 },
        launchOptions: {
          args: ['--start-maximized'],
          slowMo: testConfig.browser.slowMo || 0,
        }
      },
    },
  ],

  // Output folder for test artifacts
  outputDir: '../test-results/artifacts',
});
