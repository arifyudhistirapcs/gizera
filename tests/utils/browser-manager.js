const { chromium } = require('playwright');
const path = require('path');
const fs = require('fs');

class BrowserManager {
  constructor() {
    this.browser = null;
    this.context = null;
    this.page = null;
    this.config = null;
  }

  /**
   * Initialize browser manager with configuration
   */
  async initialize(config) {
    this.config = config;
  }

  /**
   * Launch Chrome browser in headed or headless mode
   */
  async launchBrowser(headless = false) {
    try {
      this.browser = await chromium.launch({
        headless: headless,
        slowMo: this.config?.browser?.slowMo || 0,
        args: ['--start-maximized']
      });
      console.log(`Browser launched in ${headless ? 'headless' : 'headed'} mode`);
      return this.browser;
    } catch (error) {
      throw new Error(`Failed to launch browser: ${error.message}`);
    }
  }

  /**
   * Create a new browser context with clean state
   */
  async createContext(options = {}) {
    if (!this.browser) {
      throw new Error('Browser not launched. Call launchBrowser() first.');
    }

    try {
      const contextOptions = {
        viewport: this.config?.browser?.viewport || { width: 1920, height: 1080 },
        ...options
      };

      this.context = await this.browser.newContext(contextOptions);
      this.page = await this.context.newPage();
      
      // Set default timeouts
      if (this.config?.timeouts) {
        this.page.setDefaultTimeout(this.config.timeouts.default || 30000);
        this.page.setDefaultNavigationTimeout(this.config.timeouts.navigation || 60000);
      }

      console.log('Browser context created with clean state');
      return this.context;
    } catch (error) {
      throw new Error(`Failed to create browser context: ${error.message}`);
    }
  }


  /**
   * Navigate to a specific URL
   */
  async navigateToPage(url) {
    if (!this.page) {
      throw new Error('Page not created. Call createContext() first.');
    }

    try {
      console.log(`Navigating to: ${url}`);
      await this.page.goto(url, {
        waitUntil: 'networkidle',
        timeout: this.config?.timeouts?.navigation || 60000
      });
      console.log(`Successfully navigated to: ${url}`);
    } catch (error) {
      if (error.message.includes('net::ERR')) {
        throw new Error(`Navigation failed - URL unreachable: ${url}. ${error.message}`);
      }
      throw new Error(`Navigation error: ${error.message}`);
    }
  }

  /**
   * Capture screenshot and save to file
   */
  async captureScreenshot(filename) {
    if (!this.page) {
      throw new Error('Page not created. Call createContext() first.');
    }

    try {
      const screenshotPath = this.config?.screenshots?.path || 'screenshots';
      const fullPath = path.join(process.cwd(), screenshotPath);
      
      // Ensure screenshot directory exists
      if (!fs.existsSync(fullPath)) {
        fs.mkdirSync(fullPath, { recursive: true });
      }

      const filepath = path.join(fullPath, filename);
      await this.page.screenshot({ path: filepath, fullPage: true });
      console.log(`Screenshot saved: ${filepath}`);
      return filepath;
    } catch (error) {
      console.error(`Failed to capture screenshot: ${error.message}`);
      throw new Error(`Screenshot capture failed: ${error.message}`);
    }
  }

  /**
   * Clean up browser resources
   */
  async cleanup() {
    try {
      if (this.page) {
        await this.page.close();
        this.page = null;
      }
      if (this.context) {
        await this.context.close();
        this.context = null;
      }
      if (this.browser) {
        await this.browser.close();
        this.browser = null;
      }
      console.log('Browser resources cleaned up');
    } catch (error) {
      console.error(`Cleanup error: ${error.message}`);
    }
  }

  /**
   * Reset browser state (clear cookies, storage, etc.)
   */
  async resetState() {
    if (!this.context) {
      throw new Error('Context not created. Call createContext() first.');
    }

    try {
      await this.context.clearCookies();
      await this.page.evaluate(() => {
        localStorage.clear();
        sessionStorage.clear();
      });
      console.log('Browser state reset (cookies and storage cleared)');
    } catch (error) {
      throw new Error(`Failed to reset state: ${error.message}`);
    }
  }

  /**
   * Get current page instance
   */
  getPage() {
    return this.page;
  }

  /**
   * Get current context instance
   */
  getContext() {
    return this.context;
  }

  /**
   * Get current browser instance
   */
  getBrowser() {
    return this.browser;
  }
}

module.exports = BrowserManager;
