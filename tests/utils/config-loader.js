const fs = require('fs');
const path = require('path');

class ConfigLoader {
  constructor() {
    this.config = null;
  }

  /**
   * Load configuration from JSON file and environment variables
   * Environment variables take precedence over JSON config
   */
  load(configPath = null) {
    const defaultPath = path.join(__dirname, '../config/test.config.json');
    const filePath = configPath || defaultPath;

    // Check if config file exists
    if (!fs.existsSync(filePath)) {
      throw new Error(`Configuration file not found: ${filePath}`);
    }

    // Read and parse JSON config
    try {
      const fileContent = fs.readFileSync(filePath, 'utf8');
      this.config = JSON.parse(fileContent);
    } catch (error) {
      throw new Error(`Failed to parse configuration file: ${error.message}`);
    }

    // Override with environment variables if present
    if (process.env.PWA_BASE_URL) {
      this.config.pwaBaseUrl = process.env.PWA_BASE_URL;
    }
    if (process.env.BACKEND_BASE_URL) {
      this.config.backendBaseUrl = process.env.BACKEND_BASE_URL;
    }
    if (process.env.HEADLESS !== undefined) {
      this.config.browser.headless = process.env.HEADLESS === 'true';
    }

    // Validate required fields
    this.validate();

    return this.config;
  }

  /**
   * Validate configuration has all required fields
   */
  validate() {
    const errors = [];

    // Check required fields
    if (!this.config.pwaBaseUrl) {
      errors.push('pwaBaseUrl is required');
    }
    if (!this.config.backendBaseUrl) {
      errors.push('backendBaseUrl is required');
    }

    // Validate URLs
    if (this.config.pwaBaseUrl && !this.isValidUrl(this.config.pwaBaseUrl)) {
      errors.push(`pwaBaseUrl is not a valid URL: ${this.config.pwaBaseUrl}`);
    }
    if (this.config.backendBaseUrl && !this.isValidUrl(this.config.backendBaseUrl)) {
      errors.push(`backendBaseUrl is not a valid URL: ${this.config.backendBaseUrl}`);
    }


    // Validate timeouts
    if (this.config.timeouts) {
      if (this.config.timeouts.default && this.config.timeouts.default <= 0) {
        errors.push('timeouts.default must be greater than 0');
      }
      if (this.config.timeouts.navigation && this.config.timeouts.navigation <= 0) {
        errors.push('timeouts.navigation must be greater than 0');
      }
      if (this.config.timeouts.action && this.config.timeouts.action <= 0) {
        errors.push('timeouts.action must be greater than 0');
      }
    }

    if (errors.length > 0) {
      throw new Error(`Configuration validation failed:\n${errors.join('\n')}`);
    }
  }

  /**
   * Check if a string is a valid URL
   */
  isValidUrl(urlString) {
    try {
      new URL(urlString);
      return true;
    } catch (error) {
      return false;
    }
  }

  /**
   * Get the loaded configuration
   */
  getConfig() {
    if (!this.config) {
      throw new Error('Configuration not loaded. Call load() first.');
    }
    return this.config;
  }
}

module.exports = ConfigLoader;
