/**
 * CRUD Helper Utility
 * Provides reusable functions for Create, Read, Update, Delete operations
 */

class CRUDHelper {
  constructor(page) {
    this.page = page;
  }

  /**
   * Generate timestamp for unique test data
   */
  getTimestamp() {
    return Date.now();
  }

  /**
   * Wait for success message
   */
  async waitForSuccessMessage(timeout = 5000) {
    try {
      await this.page.waitForSelector(
        '.ant-message-success, .ant-notification-success, [class*="success"]',
        { timeout }
      );
      const message = await this.page.locator('.ant-message-success, .ant-notification-success').first().textContent();
      console.log(`✓ Success message: ${message}`);
      return true;
    } catch (error) {
      console.log('⚠ No success message found');
      return false;
    }
  }

  /**
   * Wait for error message
   */
  async waitForErrorMessage(timeout = 5000) {
    try {
      await this.page.waitForSelector(
        '.ant-message-error, .ant-notification-error, .ant-form-item-explain-error',
        { timeout }
      );
      const message = await this.page.locator('.ant-message-error, .ant-notification-error').first().textContent();
      console.log(`✓ Error message: ${message}`);
      return true;
    } catch (error) {
      console.log('⚠ No error message found');
      return false;
    }
  }

  /**
   * Fill text input by label, ID, or placeholder
   */
  async fillInput(labelOrIdOrPlaceholder, value) {
    try {
      // Try by ID first (most reliable)
      const byId = this.page.locator(`#${labelOrIdOrPlaceholder}`);
      if (await byId.count() > 0) {
        await byId.fill(value);
        console.log(`✓ Filled input by ID "${labelOrIdOrPlaceholder}": ${value}`);
        return true;
      }

      // Try by label (look for label with "for" attribute)
      const labelElement = this.page.locator(`label:has-text("${labelOrIdOrPlaceholder}")`);
      if (await labelElement.count() > 0) {
        const forAttr = await labelElement.first().getAttribute('for');
        if (forAttr) {
          const input = this.page.locator(`#${forAttr}`);
          if (await input.count() > 0) {
            await input.fill(value);
            console.log(`✓ Filled input "${labelOrIdOrPlaceholder}": ${value}`);
            return true;
          }
        }
      }

      // Try by placeholder
      const byPlaceholder = this.page.locator(`input[placeholder*="${labelOrIdOrPlaceholder}"]`);
      if (await byPlaceholder.count() > 0) {
        await byPlaceholder.first().fill(value);
        console.log(`✓ Filled input by placeholder "${labelOrIdOrPlaceholder}": ${value}`);
        return true;
      }

      console.log(`⚠ Input "${labelOrIdOrPlaceholder}" not found`);
      return false;
    } catch (error) {
      console.log(`✗ Error filling input "${labelOrIdOrPlaceholder}": ${error.message}`);
      return false;
    }
  }

  /**
   * Fill textarea by label, ID, or placeholder
   */
  async fillTextarea(labelOrIdOrPlaceholder, value) {
    try {
      // Try by ID first
      const byId = this.page.locator(`#${labelOrIdOrPlaceholder}`);
      if (await byId.count() > 0) {
        await byId.fill(value);
        console.log(`✓ Filled textarea by ID "${labelOrIdOrPlaceholder}": ${value}`);
        return true;
      }

      // Try by label
      const labelElement = this.page.locator(`label:has-text("${labelOrIdOrPlaceholder}")`);
      if (await labelElement.count() > 0) {
        const forAttr = await labelElement.first().getAttribute('for');
        if (forAttr) {
          const textarea = this.page.locator(`#${forAttr}`);
          if (await textarea.count() > 0) {
            await textarea.fill(value);
            console.log(`✓ Filled textarea "${labelOrIdOrPlaceholder}": ${value}`);
            return true;
          }
        }
      }

      // Try by placeholder
      const byPlaceholder = this.page.locator(`textarea[placeholder*="${labelOrIdOrPlaceholder}"]`);
      if (await byPlaceholder.count() > 0) {
        await byPlaceholder.first().fill(value);
        console.log(`✓ Filled textarea by placeholder "${labelOrIdOrPlaceholder}": ${value}`);
        return true;
      }

      console.log(`⚠ Textarea "${labelOrIdOrPlaceholder}" not found`);
      return false;
    } catch (error) {
      console.log(`✗ Error filling textarea "${labelOrIdOrPlaceholder}": ${error.message}`);
      return false;
    }
  }

  /**
   * Select dropdown option by label
   */
  async selectDropdown(label, optionText) {
    try {
      // Find the select by label
      const selectWrapper = this.page.locator(`label:has-text("${label}") + * .ant-select, label:has-text("${label}") ~ * .ant-select`);
      
      if (await selectWrapper.count() > 0) {
        await selectWrapper.first().click();
        await this.page.waitForTimeout(500);
        
        // Click the option
        const option = this.page.locator(`.ant-select-dropdown .ant-select-item:has-text("${optionText}")`);
        if (await option.count() > 0) {
          await option.first().click();
          console.log(`✓ Selected "${optionText}" in dropdown "${label}"`);
          return true;
        }
      }
      
      console.log(`⚠ Dropdown "${label}" or option "${optionText}" not found`);
      return false;
    } catch (error) {
      console.log(`✗ Error selecting dropdown "${label}": ${error.message}`);
      return false;
    }
  }

  /**
   * Click submit button
   */
  async clickSubmit(buttonText = 'OK') {
    try {
      const submitButton = this.page.locator(`button:has-text("${buttonText}"), button[type="submit"]`);
      if (await submitButton.count() > 0) {
        await submitButton.first().click();
        console.log(`✓ Clicked submit button "${buttonText}"`);
        await this.page.waitForTimeout(1000);
        return true;
      }
      console.log(`⚠ Submit button "${buttonText}" not found`);
      return false;
    } catch (error) {
      console.log(`✗ Error clicking submit: ${error.message}`);
      return false;
    }
  }

  /**
   * Click cancel button
   */
  async clickCancel(buttonText = 'Batal') {
    try {
      const cancelButton = this.page.locator(`button:has-text("${buttonText}")`);
      if (await cancelButton.count() > 0) {
        await cancelButton.first().click();
        console.log(`✓ Clicked cancel button "${buttonText}"`);
        await this.page.waitForTimeout(500);
        return true;
      }
      return false;
    } catch (error) {
      return false;
    }
  }

  /**
   * Verify data exists in table (checks main data table, handles pagination)
   */
  async verifyDataInTable(searchText) {
    try {
      // First check if data is visible in current view
      const tableCell = this.page.locator('table tbody td').filter({ hasText: searchText });
      let exists = await tableCell.count() > 0;
      
      if (exists) {
        console.log(`✓ Data "${searchText}" found in table (current page)`);
        return true;
      }
      
      // If not found on current page, navigate to last page (new data appears there)
      console.log(`  Data not on current page, checking last page...`);
      const paginationItems = this.page.locator('.ant-pagination-item');
      const paginationCount = await paginationItems.count();
      
      if (paginationCount > 0) {
        // Click the last page
        const lastPageButton = paginationItems.last();
        await lastPageButton.click();
        await this.page.waitForTimeout(2000); // Wait for page to load
        
        // Check if data exists on last page
        const tableCellLastPage = this.page.locator('table tbody td').filter({ hasText: searchText });
        exists = await tableCellLastPage.count() > 0;
        
        if (exists) {
          console.log(`✓ Data "${searchText}" found on last page`);
          return true;
        }
      }
      
      // If still not found, try search as fallback
      console.log(`  Data not on last page, trying search...`);
      const searchInput = this.page.locator('input[placeholder*="Cari"]').first();
      if (await searchInput.count() > 0) {
        // Clear any existing search first
        await searchInput.clear();
        await this.page.waitForTimeout(500);
        
        // Search for our data
        await searchInput.fill(searchText);
        await this.page.waitForTimeout(2000); // Wait for search to filter
        
        // Check how many rows we have
        const rowCount = await this.page.locator('table tbody tr').count();
        console.log(`  Found ${rowCount} rows after search`);
        
        if (rowCount > 0) {
          // Check all rows for our text (exact match preferred)
          const rows = await this.page.locator('table tbody tr').all();
          for (let i = 0; i < rows.length; i++) {
            const rowText = await rows[i].textContent();
            if (rowText.includes(searchText)) {
              console.log(`✓ Data "${searchText}" found via search in row ${i + 1}`);
              // Clear search before returning
              await searchInput.clear();
              await this.page.waitForTimeout(500);
              return true;
            }
          }
          
          // Log first row for debugging
          const firstRowText = await rows[0].textContent();
          console.log(`  First row after search: ${firstRowText.substring(0, 100)}...`);
        }
        
        // Clear search
        await searchInput.clear();
        await this.page.waitForTimeout(500);
      }
      
      console.log(`✗ Data "${searchText}" NOT found in table`);
      return false;
    } catch (error) {
      console.log(`✗ Error verifying data: ${error.message}`);
      return false;
    }
  }

  /**
   * Click action button in table row (Edit, Delete, Detail)
   */
  async clickRowAction(searchText, actionButtonText) {
    try {
      // First check if row is visible on current page
      let row = this.page.locator(`table tbody tr:has-text("${searchText}")`);
      
      if (await row.count() === 0) {
        // Not on current page, try last page
        console.log(`  Row not on current page, checking last page...`);
        const paginationItems = this.page.locator('.ant-pagination-item');
        const paginationCount = await paginationItems.count();
        
        if (paginationCount > 0) {
          const lastPageButton = paginationItems.last();
          await lastPageButton.click();
          await this.page.waitForTimeout(2000);
          
          // Try to find row again
          row = this.page.locator(`table tbody tr:has-text("${searchText}")`);
        }
      }
      
      if (await row.count() > 0) {
        // Find the action button in that row
        const actionButton = row.locator(`button:has-text("${actionButtonText}")`);
        
        if (await actionButton.count() > 0) {
          await actionButton.first().click();
          console.log(`✓ Clicked "${actionButtonText}" for "${searchText}"`);
          await this.page.waitForTimeout(1000);
          return true;
        }
      }
      
      console.log(`✗ Could not find "${actionButtonText}" button for "${searchText}"`);
      return false;
    } catch (error) {
      console.log(`✗ Error clicking row action: ${error.message}`);
      return false;
    }
  }

  /**
   * Confirm deletion dialog
   */
  async confirmDelete(confirmButtonText = 'Ya') {
    try {
      // Wait for confirmation dialog
      await this.page.waitForTimeout(500);
      
      const confirmButton = this.page.locator(
        `.ant-modal button:has-text("${confirmButtonText}"), ` +
        `.ant-popconfirm button:has-text("${confirmButtonText}"), ` +
        `button.ant-btn-primary:has-text("${confirmButtonText}")`
      );
      
      if (await confirmButton.count() > 0) {
        await confirmButton.first().click();
        console.log(`✓ Confirmed deletion`);
        await this.page.waitForTimeout(1000);
        return true;
      }
      
      console.log(`⚠ Confirm button not found`);
      return false;
    } catch (error) {
      console.log(`✗ Error confirming delete: ${error.message}`);
      return false;
    }
  }

  /**
   * Search for data
   */
  async search(searchText) {
    try {
      const searchInput = this.page.locator('input[type="text"], input[type="search"], .ant-input').first();
      
      if (await searchInput.count() > 0) {
        await searchInput.fill(searchText);
        await this.page.waitForTimeout(1000);
        console.log(`✓ Searched for: ${searchText}`);
        return true;
      }
      
      return false;
    } catch (error) {
      console.log(`✗ Error searching: ${error.message}`);
      return false;
    }
  }

  /**
   * Clear search
   */
  async clearSearch() {
    try {
      const searchInput = this.page.locator('input[type="text"], input[type="search"], .ant-input').first();
      
      if (await searchInput.count() > 0) {
        await searchInput.clear();
        await this.page.waitForTimeout(1000);
        console.log(`✓ Cleared search`);
        return true;
      }
      
      return false;
    } catch (error) {
      return false;
    }
  }

  /**
   * Take screenshot for debugging
   */
  async takeScreenshot(name) {
    try {
      await this.page.screenshot({ path: `tests/screenshots/crud-${name}-${Date.now()}.png` });
      console.log(`✓ Screenshot saved: crud-${name}`);
    } catch (error) {
      console.log(`⚠ Could not save screenshot: ${error.message}`);
    }
  }
}

module.exports = CRUDHelper;
