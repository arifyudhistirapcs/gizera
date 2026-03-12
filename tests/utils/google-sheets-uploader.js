/**
 * Google Sheets Uploader
 * 
 * Script untuk automated upload test results ke Google Sheets
 * menggunakan Google Sheets API v4
 */

const { google } = require('googleapis');
const fs = require('fs');
const path = require('path');

class GoogleSheetsUploader {
  constructor(credentialsPath, spreadsheetId) {
    this.credentialsPath = credentialsPath;
    this.spreadsheetId = spreadsheetId;
    this.auth = null;
    this.sheets = null;
  }

  /**
   * Initialize Google Sheets API client
   */
  async initialize() {
    try {
      // Load credentials
      const credentials = JSON.parse(fs.readFileSync(this.credentialsPath, 'utf8'));
      
      // Create auth client
      this.auth = new google.auth.GoogleAuth({
        credentials: credentials,
        scopes: ['https://www.googleapis.com/auth/spreadsheets'],
      });

      // Create sheets client
      this.sheets = google.sheets({ version: 'v4', auth: this.auth });
      
      console.log('✓ Google Sheets API initialized successfully');
      return true;
    } catch (error) {
      console.error('✗ Failed to initialize Google Sheets API:', error.message);
      throw error;
    }
  }

  /**
   * Read CSV file and convert to 2D array
   */
  readCSV(filePath) {
    const content = fs.readFileSync(filePath, 'utf8');
    const lines = content.split('\n').filter(line => line.trim());
    return lines.map(line => {
      // Simple CSV parser (handles basic cases)
      const values = [];
      let current = '';
      let inQuotes = false;
      
      for (let i = 0; i < line.length; i++) {
        const char = line[i];
        
        if (char === '"') {
          inQuotes = !inQuotes;
        } else if (char === ',' && !inQuotes) {
          values.push(current.trim());
          current = '';
        } else {
          current += char;
        }
      }
      values.push(current.trim());
      
      return values;
    });
  }

  /**
   * Create a new sheet in the spreadsheet
   */
  async createSheet(sheetName) {
    try {
      const request = {
        spreadsheetId: this.spreadsheetId,
        resource: {
          requests: [{
            addSheet: {
              properties: {
                title: sheetName
              }
            }
          }]
        }
      };

      await this.sheets.spreadsheets.batchUpdate(request);
      console.log(`✓ Created sheet: ${sheetName}`);
      return true;
    } catch (error) {
      if (error.message.includes('already exists')) {
        console.log(`⚠ Sheet "${sheetName}" already exists, will update it`);
        return true;
      }
      console.error(`✗ Failed to create sheet: ${error.message}`);
      throw error;
    }
  }

  /**
   * Clear existing data in a sheet
   */
  async clearSheet(sheetName) {
    try {
      await this.sheets.spreadsheets.values.clear({
        spreadsheetId: this.spreadsheetId,
        range: `${sheetName}!A1:Z1000`,
      });
      console.log(`✓ Cleared sheet: ${sheetName}`);
    } catch (error) {
      console.error(`✗ Failed to clear sheet: ${error.message}`);
      throw error;
    }
  }

  /**
   * Upload data to a sheet
   */
  async uploadData(sheetName, data) {
    try {
      const request = {
        spreadsheetId: this.spreadsheetId,
        range: `${sheetName}!A1`,
        valueInputOption: 'RAW',
        resource: {
          values: data
        }
      };

      await this.sheets.spreadsheets.values.update(request);
      console.log(`✓ Uploaded ${data.length} rows to sheet: ${sheetName}`);
      return true;
    } catch (error) {
      console.error(`✗ Failed to upload data: ${error.message}`);
      throw error;
    }
  }

  /**
   * Append data to a sheet
   */
  async appendData(sheetName, data) {
    try {
      const request = {
        spreadsheetId: this.spreadsheetId,
        range: `${sheetName}!A:A`,
        valueInputOption: 'RAW',
        insertDataOption: 'INSERT_ROWS',
        resource: {
          values: data
        }
      };

      await this.sheets.spreadsheets.values.append(request);
      console.log(`✓ Appended ${data.length} rows to sheet: ${sheetName}`);
      return true;
    } catch (error) {
      console.error(`✗ Failed to append data: ${error.message}`);
      throw error;
    }
  }

  /**
   * Format sheet headers (bold, background color)
   */
  async formatHeaders(sheetName) {
    try {
      // Get sheet ID
      const spreadsheet = await this.sheets.spreadsheets.get({
        spreadsheetId: this.spreadsheetId
      });
      
      const sheet = spreadsheet.data.sheets.find(s => s.properties.title === sheetName);
      if (!sheet) {
        console.log(`⚠ Sheet "${sheetName}" not found for formatting`);
        return;
      }
      
      const sheetId = sheet.properties.sheetId;

      // Format header row
      const requests = [
        {
          repeatCell: {
            range: {
              sheetId: sheetId,
              startRowIndex: 0,
              endRowIndex: 1
            },
            cell: {
              userEnteredFormat: {
                backgroundColor: {
                  red: 0.2,
                  green: 0.6,
                  blue: 0.86
                },
                textFormat: {
                  bold: true,
                  foregroundColor: {
                    red: 1,
                    green: 1,
                    blue: 1
                  }
                }
              }
            },
            fields: 'userEnteredFormat(backgroundColor,textFormat)'
          }
        },
        {
          updateSheetProperties: {
            properties: {
              sheetId: sheetId,
              gridProperties: {
                frozenRowCount: 1
              }
            },
            fields: 'gridProperties.frozenRowCount'
          }
        }
      ];

      await this.sheets.spreadsheets.batchUpdate({
        spreadsheetId: this.spreadsheetId,
        resource: { requests }
      });

      console.log(`✓ Formatted headers for sheet: ${sheetName}`);
    } catch (error) {
      console.error(`✗ Failed to format headers: ${error.message}`);
      // Don't throw, formatting is optional
    }
  }

  /**
   * Format status column with colors
   */
  async formatStatusColumn(sheetName, statusColumnIndex, rowCount) {
    try {
      const spreadsheet = await this.sheets.spreadsheets.get({
        spreadsheetId: this.spreadsheetId
      });
      
      const sheet = spreadsheet.data.sheets.find(s => s.properties.title === sheetName);
      if (!sheet) return;
      
      const sheetId = sheet.properties.sheetId;

      // Add conditional formatting for status column
      const requests = [
        {
          addConditionalFormatRule: {
            rule: {
              ranges: [{
                sheetId: sheetId,
                startRowIndex: 1,
                endRowIndex: rowCount,
                startColumnIndex: statusColumnIndex,
                endColumnIndex: statusColumnIndex + 1
              }],
              booleanRule: {
                condition: {
                  type: 'TEXT_CONTAINS',
                  values: [{ userEnteredValue: 'PASS' }]
                },
                format: {
                  backgroundColor: {
                    red: 0.7,
                    green: 0.9,
                    blue: 0.7
                  }
                }
              }
            },
            index: 0
          }
        },
        {
          addConditionalFormatRule: {
            rule: {
              ranges: [{
                sheetId: sheetId,
                startRowIndex: 1,
                endRowIndex: rowCount,
                startColumnIndex: statusColumnIndex,
                endColumnIndex: statusColumnIndex + 1
              }],
              booleanRule: {
                condition: {
                  type: 'TEXT_CONTAINS',
                  values: [{ userEnteredValue: 'FAIL' }]
                },
                format: {
                  backgroundColor: {
                    red: 0.95,
                    green: 0.7,
                    blue: 0.7
                  }
                }
              }
            },
            index: 1
          }
        },
        {
          addConditionalFormatRule: {
            rule: {
              ranges: [{
                sheetId: sheetId,
                startRowIndex: 1,
                endRowIndex: rowCount,
                startColumnIndex: statusColumnIndex,
                endColumnIndex: statusColumnIndex + 1
              }],
              booleanRule: {
                condition: {
                  type: 'TEXT_CONTAINS',
                  values: [{ userEnteredValue: 'SKIP' }]
                },
                format: {
                  backgroundColor: {
                    red: 1,
                    green: 0.95,
                    blue: 0.7
                  }
                }
              }
            },
            index: 2
          }
        }
      ];

      await this.sheets.spreadsheets.batchUpdate({
        spreadsheetId: this.spreadsheetId,
        resource: { requests }
      });

      console.log(`✓ Formatted status column for sheet: ${sheetName}`);
    } catch (error) {
      console.error(`✗ Failed to format status column: ${error.message}`);
    }
  }

  /**
   * Auto-resize columns to fit content
   */
  async autoResizeColumns(sheetName) {
    try {
      const spreadsheet = await this.sheets.spreadsheets.get({
        spreadsheetId: this.spreadsheetId
      });
      
      const sheet = spreadsheet.data.sheets.find(s => s.properties.title === sheetName);
      if (!sheet) return;
      
      const sheetId = sheet.properties.sheetId;

      // Auto-resize all columns
      const requests = [{
        autoResizeDimensions: {
          dimensions: {
            sheetId: sheetId,
            dimension: 'COLUMNS',
            startIndex: 0,
            endIndex: 20
          }
        }
      }];

      await this.sheets.spreadsheets.batchUpdate({
        spreadsheetId: this.spreadsheetId,
        resource: { requests }
      });

      console.log(`✓ Auto-resized columns for sheet: ${sheetName}`);
    } catch (error) {
      console.error(`✗ Failed to auto-resize columns: ${error.message}`);
    }
  }

  /**
   * Upload test results from CSV file (append mode)
   */
  async uploadTestResults(csvPath, sheetName, appendMode = false) {
    console.log(`\n📤 Uploading test results to sheet: ${sheetName}`);
    
    // Read CSV data
    const data = this.readCSV(csvPath);
    console.log(`✓ Read ${data.length} rows from CSV`);

    // Create sheet if not exists
    await this.createSheet(sheetName);
    
    if (!appendMode) {
      // Clear sheet if not in append mode
      await this.clearSheet(sheetName);
      // Upload data with headers
      await this.uploadData(sheetName, data);
    } else {
      // Append mode: skip header row and append to existing data
      const dataWithoutHeader = data.slice(1);
      await this.appendData(sheetName, dataWithoutHeader);
    }

    // Format sheet (only if not append mode)
    if (!appendMode) {
      await this.formatHeaders(sheetName);
    }
    
    // Format status column (column G = index 6)
    await this.formatStatusColumn(sheetName, 6, data.length + 100); // Add buffer for future appends
    
    // Auto-resize columns
    await this.autoResizeColumns(sheetName);

    console.log(`✅ Successfully uploaded test results to: ${sheetName}\n`);
  }

  /**
   * Upload test summary from CSV file
   */
  async uploadTestSummary(csvPath, sheetName) {
    console.log(`\n📤 Uploading test summary to sheet: ${sheetName}`);
    
    // Read CSV data
    const data = this.readCSV(csvPath);
    console.log(`✓ Read ${data.length} rows from CSV`);

    // Create or clear sheet
    await this.createSheet(sheetName);
    await this.clearSheet(sheetName);

    // Upload data
    await this.uploadData(sheetName, data);

    // Format sheet
    await this.formatHeaders(sheetName);
    
    // Auto-resize columns
    await this.autoResizeColumns(sheetName);

    console.log(`✅ Successfully uploaded test summary to: ${sheetName}\n`);
  }
}

module.exports = GoogleSheetsUploader;
