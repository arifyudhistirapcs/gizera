const TestCaseLoader = require('./test-case-loader');
const BrowserManager = require('./browser-manager');

class TestExecutor {
  constructor(browserManager, testCaseLoader) {
    this.browserManager = browserManager;
    this.testCaseLoader = testCaseLoader || new TestCaseLoader();
    this.results = [];
  }

  /**
   * Load test cases for a specific module
   */
  async loadTestCases(module) {
    return this.testCaseLoader.loadTestCases(module);
  }

  /**
   * Execute a single test case
   */
  async executeTestCase(testCase, page) {
    const startTime = Date.now();
    const result = {
      id: testCase.id,
      module: testCase.module,
      scenario: testCase.scenario,
      status: 'not_run',
      actualResults: [],
      error: null,
      executionTime: 0,
      screenshot: null
    };

    try {
      console.log(`\nExecuting test case: ${testCase.id} - ${testCase.scenario}`);
      
      // Execute each step
      for (let i = 0; i < testCase.steps.length; i++) {
        const step = testCase.steps[i];
        console.log(`  Step ${i + 1}: ${step}`);
        
        try {
          // This is a placeholder - actual step execution would be implemented
          // based on the specific step description
          await this.executeStep(step, page);
          
          // Capture actual result for this step
          const actualResult = await this.captureActualResult(page, testCase.expectedResults[i]);
          result.actualResults.push(actualResult);
          
        } catch (stepError) {
          console.error(`  Step ${i + 1} failed: ${stepError.message}`);
          result.error = stepError.message;
          result.status = 'fail';
          
          // Capture screenshot on failure
          const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
          const screenshotName = `${timestamp}-${testCase.id}.png`;
          result.screenshot = await this.browserManager.captureScreenshot(screenshotName);
          
          break;
        }
      }

      // Determine test status if not already failed
      if (result.status !== 'fail') {
        result.status = this.compareResults(testCase.expectedResults, result.actualResults) ? 'pass' : 'fail';
      }

      result.executionTime = Date.now() - startTime;
      console.log(`  Result: ${result.status} (${result.executionTime}ms)`);

    } catch (error) {
      result.status = 'fail';
      result.error = error.message;
      result.executionTime = Date.now() - startTime;
      console.error(`Test case failed: ${error.message}`);
    }

    // Update test case in storage
    await this.updateTestStatus(testCase.module, testCase.id, result.status, result.actualResults, result.executionTime);

    return result;
  }


  /**
   * Execute a single step (placeholder - to be implemented based on step description)
   */
  async executeStep(step, page) {
    // This is a simplified implementation
    // In a real scenario, you would parse the step description and execute appropriate actions
    
    const stepLower = step.toLowerCase();
    
    // Navigation steps
    if (stepLower.includes('navigate to') || stepLower.includes('go to')) {
      // Extract URL or page name from step
      // For now, just wait a bit
      await page.waitForTimeout(500);
    }
    // Click steps
    else if (stepLower.includes('click')) {
      // Would need to extract selector from step
      await page.waitForTimeout(500);
    }
    // Input steps
    else if (stepLower.includes('enter') || stepLower.includes('type')) {
      // Would need to extract selector and value from step
      await page.waitForTimeout(500);
    }
    // Wait steps
    else if (stepLower.includes('wait')) {
      await page.waitForTimeout(1000);
    }
    else {
      // Generic wait for other steps
      await page.waitForTimeout(500);
    }
  }

  /**
   * Capture actual result from page
   */
  async captureActualResult(page, expectedResult) {
    try {
      // This is a placeholder - actual implementation would depend on what to verify
      const url = page.url();
      const title = await page.title();
      
      return `Page: ${title}, URL: ${url}`;
    } catch (error) {
      return `Error capturing result: ${error.message}`;
    }
  }

  /**
   * Compare expected and actual results
   */
  compareResults(expectedResults, actualResults) {
    if (expectedResults.length !== actualResults.length) {
      return false;
    }
    
    // Simple comparison - in real scenario would be more sophisticated
    return actualResults.length > 0;
  }

  /**
   * Execute all test cases in a module
   */
  async executeTestSuite(module) {
    console.log(`\n========================================`);
    console.log(`Executing test suite for module: ${module}`);
    console.log(`========================================`);

    const testCases = await this.loadTestCases(module);
    const page = this.browserManager.getPage();
    
    if (!page) {
      throw new Error('Browser page not available. Ensure browser is initialized.');
    }

    const suiteResults = {
      module: module,
      totalTests: testCases.length,
      passed: 0,
      failed: 0,
      blocked: 0,
      results: []
    };

    for (const testCase of testCases) {
      const result = await this.executeTestCase(testCase, page);
      suiteResults.results.push(result);
      
      if (result.status === 'pass') suiteResults.passed++;
      else if (result.status === 'fail') suiteResults.failed++;
      else if (result.status === 'blocked') suiteResults.blocked++;
    }

    console.log(`\n========================================`);
    console.log(`Module: ${module} - Summary`);
    console.log(`Total: ${suiteResults.totalTests}, Passed: ${suiteResults.passed}, Failed: ${suiteResults.failed}, Blocked: ${suiteResults.blocked}`);
    console.log(`========================================\n`);

    return suiteResults;
  }

  /**
   * Execute all tests across all modules
   */
  async executeAllTests() {
    const allTestCases = this.testCaseLoader.loadAllTestCases();
    const modules = Object.keys(allTestCases);
    
    const allResults = {
      totalModules: modules.length,
      moduleResults: []
    };

    for (const module of modules) {
      try {
        const suiteResult = await this.executeTestSuite(module);
        allResults.moduleResults.push(suiteResult);
      } catch (error) {
        console.error(`Failed to execute test suite for ${module}: ${error.message}`);
      }
    }

    return allResults;
  }

  /**
   * Update test case status in storage
   */
  async updateTestStatus(module, testCaseId, status, actualResults, executionTime) {
    try {
      this.testCaseLoader.updateTestCase(module, testCaseId, {
        status: status,
        actualResults: actualResults,
        executionTime: executionTime
      });
    } catch (error) {
      console.error(`Failed to update test status: ${error.message}`);
    }
  }

  /**
   * Capture test failure details
   */
  async captureTestFailure(testCase, error) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const screenshotName = `${timestamp}-${testCase.id}-failure.png`;
    
    const failureDetails = {
      testCaseId: testCase.id,
      module: testCase.module,
      scenario: testCase.scenario,
      error: error.message,
      timestamp: new Date().toISOString(),
      screenshot: null
    };

    try {
      failureDetails.screenshot = await this.browserManager.captureScreenshot(screenshotName);
    } catch (screenshotError) {
      console.error(`Failed to capture failure screenshot: ${screenshotError.message}`);
    }

    return failureDetails;
  }
}

module.exports = TestExecutor;
