const fs = require('fs');
const path = require('path');

class TestCaseLoader {
  constructor() {
    this.testCasesDir = path.join(__dirname, '../test-cases');
  }

  /**
   * Load test cases for a specific module
   */
  loadTestCases(module) {
    const filePath = path.join(this.testCasesDir, module, 'test-cases.json');
    
    if (!fs.existsSync(filePath)) {
      throw new Error(`Test cases file not found for module: ${module}`);
    }

    try {
      const fileContent = fs.readFileSync(filePath, 'utf8');
      const testCases = JSON.parse(fileContent);
      
      // Validate test cases
      this.validateTestCases(testCases);
      
      return testCases;
    } catch (error) {
      throw new Error(`Failed to load test cases for ${module}: ${error.message}`);
    }
  }

  /**
   * Validate test case structure
   */
  validateTestCases(testCases) {
    if (!Array.isArray(testCases)) {
      throw new Error('Test cases must be an array');
    }

    testCases.forEach((testCase, index) => {
      const errors = [];

      // Check required fields
      if (!testCase.id) errors.push('id is required');
      if (!testCase.module) errors.push('module is required');
      if (!testCase.scenario || testCase.scenario.trim() === '') {
        errors.push('scenario is required and must not be empty');
      }
      if (!Array.isArray(testCase.steps) || testCase.steps.length === 0) {
        errors.push('steps must be a non-empty array');
      }
      if (!Array.isArray(testCase.expectedResults) || testCase.expectedResults.length === 0) {
        errors.push('expectedResults must be a non-empty array');
      }
      if (!testCase.hasOwnProperty('actualResults')) {
        errors.push('actualResults field is required (can be empty array)');
      }
      if (!testCase.status) {
        errors.push('status is required');
      }
      
      // Validate status value
      const validStatuses = ['not_run', 'pass', 'fail', 'blocked'];
      if (testCase.status && !validStatuses.includes(testCase.status)) {
        errors.push(`status must be one of: ${validStatuses.join(', ')}`);
      }

      if (errors.length > 0) {
        throw new Error(`Test case at index ${index} validation failed:\n${errors.join('\n')}`);
      }
    });
  }


  /**
   * Update test case status and results
   */
  updateTestCase(module, testCaseId, updates) {
    const filePath = path.join(this.testCasesDir, module, 'test-cases.json');
    
    if (!fs.existsSync(filePath)) {
      throw new Error(`Test cases file not found for module: ${module}`);
    }

    try {
      const fileContent = fs.readFileSync(filePath, 'utf8');
      const testCases = JSON.parse(fileContent);
      
      // Find and update the test case
      const testCaseIndex = testCases.findIndex(tc => tc.id === testCaseId);
      if (testCaseIndex === -1) {
        throw new Error(`Test case not found: ${testCaseId}`);
      }

      // Apply updates
      testCases[testCaseIndex] = {
        ...testCases[testCaseIndex],
        ...updates,
        lastExecuted: new Date().toISOString(),
      };

      // Write back to file
      fs.writeFileSync(filePath, JSON.stringify(testCases, null, 2), 'utf8');
      
      return testCases[testCaseIndex];
    } catch (error) {
      throw new Error(`Failed to update test case ${testCaseId}: ${error.message}`);
    }
  }

  /**
   * Get all test cases across all modules
   */
  loadAllTestCases() {
    const modules = fs.readdirSync(this.testCasesDir, { withFileTypes: true })
      .filter(dirent => dirent.isDirectory())
      .map(dirent => dirent.name);

    const allTestCases = {};
    
    modules.forEach(module => {
      try {
        allTestCases[module] = this.loadTestCases(module);
      } catch (error) {
        console.warn(`Skipping module ${module}: ${error.message}`);
      }
    });

    return allTestCases;
  }

  /**
   * Create a new test case file for a module
   */
  createTestCaseFile(module, testCases = []) {
    const moduleDir = path.join(this.testCasesDir, module);
    const filePath = path.join(moduleDir, 'test-cases.json');

    // Create module directory if it doesn't exist
    if (!fs.existsSync(moduleDir)) {
      fs.mkdirSync(moduleDir, { recursive: true });
    }

    // Validate test cases before writing
    if (testCases.length > 0) {
      this.validateTestCases(testCases);
    }

    // Write test cases to file
    fs.writeFileSync(filePath, JSON.stringify(testCases, null, 2), 'utf8');
    console.log(`Test case file created for module: ${module}`);
  }
}

module.exports = TestCaseLoader;
