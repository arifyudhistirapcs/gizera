const fs = require('fs');
const path = require('path');
const readline = require('readline');

class BugReporter {
  constructor(config) {
    this.config = config;
    this.bugReportsDir = path.join(__dirname, '../bug-reports');
    
    // Ensure bug reports directory exists
    if (!fs.existsSync(this.bugReportsDir)) {
      fs.mkdirSync(this.bugReportsDir, { recursive: true });
    }
  }

  /**
   * Create a bug report from a failed test case
   */
  async createBugReport(testCase, failureDetails) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const bugId = `bug-${timestamp}`;

    const bugReport = {
      id: bugId,
      testCaseId: testCase.id,
      module: testCase.module,
      severity: this.determineSeverity(testCase),
      title: `${testCase.scenario} - Test Failed`,
      description: `Test case "${testCase.scenario}" failed during execution`,
      reproductionSteps: testCase.steps,
      expectedBehavior: testCase.expectedResults.join('; '),
      actualBehavior: failureDetails.actualResults ? failureDetails.actualResults.join('; ') : 'Test execution failed',
      screenshots: failureDetails.screenshot ? [failureDetails.screenshot] : [],
      errorMessages: failureDetails.error ? [failureDetails.error] : [],
      environment: {
        pwaUrl: this.config.pwaBaseUrl,
        backendUrl: this.config.backendBaseUrl,
        browser: 'Chrome (Playwright)',
        timestamp: new Date().toISOString()
      },
      status: 'open',
      fixApplied: false,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    console.log(`\n🐛 Bug Report Created: ${bugId}`);
    console.log(`   Module: ${bugReport.module}`);
    console.log(`   Severity: ${bugReport.severity}`);
    console.log(`   Title: ${bugReport.title}`);

    return bugReport;
  }

  /**
   * Determine bug severity based on test case tags or default
   */
  determineSeverity(testCase) {
    if (testCase.tags) {
      if (testCase.tags.includes('critical')) return 'critical';
      if (testCase.tags.includes('high')) return 'high';
      if (testCase.tags.includes('low')) return 'low';
    }
    return 'medium'; // default
  }

  /**
   * Save bug report to file
   */
  async saveBugReport(bugReport) {
    try {
      const filename = `${bugReport.id}.json`;
      const filepath = path.join(this.bugReportsDir, filename);
      
      fs.writeFileSync(filepath, JSON.stringify(bugReport, null, 2), 'utf8');
      console.log(`   Bug report saved: ${filepath}`);
      
      return filepath;
    } catch (error) {
      throw new Error(`Failed to save bug report: ${error.message}`);
    }
  }


  /**
   * Load bug reports with optional filtering
   */
  async loadBugReports(filter = {}) {
    try {
      const files = fs.readdirSync(this.bugReportsDir)
        .filter(file => file.endsWith('.json'));

      const bugReports = files.map(file => {
        const filepath = path.join(this.bugReportsDir, file);
        const content = fs.readFileSync(filepath, 'utf8');
        return JSON.parse(content);
      });

      // Apply filters
      let filtered = bugReports;
      
      if (filter.module) {
        filtered = filtered.filter(bug => bug.module === filter.module);
      }
      if (filter.status) {
        filtered = filtered.filter(bug => bug.status === filter.status);
      }
      if (filter.severity) {
        filtered = filtered.filter(bug => bug.severity === filter.severity);
      }

      return filtered;
    } catch (error) {
      throw new Error(`Failed to load bug reports: ${error.message}`);
    }
  }

  /**
   * Update bug report status
   */
  async updateBugStatus(bugId, status) {
    try {
      const filepath = path.join(this.bugReportsDir, `${bugId}.json`);
      
      if (!fs.existsSync(filepath)) {
        throw new Error(`Bug report not found: ${bugId}`);
      }

      const content = fs.readFileSync(filepath, 'utf8');
      const bugReport = JSON.parse(content);
      
      bugReport.status = status;
      bugReport.updatedAt = new Date().toISOString();
      
      fs.writeFileSync(filepath, JSON.stringify(bugReport, null, 2), 'utf8');
      console.log(`Bug report ${bugId} status updated to: ${status}`);
      
      return bugReport;
    } catch (error) {
      throw new Error(`Failed to update bug status: ${error.message}`);
    }
  }

  /**
   * Request user confirmation for bug fix
   */
  async requestUserConfirmation(bugReport) {
    console.log(`\n${'='.repeat(80)}`);
    console.log(`BUG REPORT: ${bugReport.id}`);
    console.log(`${'='.repeat(80)}`);
    console.log(`Module: ${bugReport.module}`);
    console.log(`Severity: ${bugReport.severity}`);
    console.log(`Title: ${bugReport.title}`);
    console.log(`\nDescription:`);
    console.log(`  ${bugReport.description}`);
    console.log(`\nExpected Behavior:`);
    console.log(`  ${bugReport.expectedBehavior}`);
    console.log(`\nActual Behavior:`);
    console.log(`  ${bugReport.actualBehavior}`);
    
    if (bugReport.errorMessages.length > 0) {
      console.log(`\nError Messages:`);
      bugReport.errorMessages.forEach(msg => console.log(`  - ${msg}`));
    }
    
    if (bugReport.screenshots.length > 0) {
      console.log(`\nScreenshots:`);
      bugReport.screenshots.forEach(screenshot => console.log(`  - ${screenshot}`));
    }
    
    console.log(`\n${'='.repeat(80)}`);
    console.log(`Would you like to proceed with fixing this bug?`);
    console.log(`${'='.repeat(80)}\n`);

    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout
    });

    return new Promise((resolve) => {
      rl.question('Enter your choice (yes/no): ', (answer) => {
        rl.close();
        const confirmed = answer.toLowerCase().trim() === 'yes' || answer.toLowerCase().trim() === 'y';
        console.log(confirmed ? '\n✓ User confirmed bug fix\n' : '\n✗ User declined bug fix\n');
        resolve(confirmed);
      });
    });
  }

  /**
   * Present bug report to user (without confirmation)
   */
  async presentBugReport(bugReport) {
    console.log(`\n${'='.repeat(80)}`);
    console.log(`BUG REPORT: ${bugReport.id}`);
    console.log(`${'='.repeat(80)}`);
    console.log(`Module: ${bugReport.module}`);
    console.log(`Severity: ${bugReport.severity}`);
    console.log(`Title: ${bugReport.title}`);
    console.log(`Description: ${bugReport.description}`);
    console.log(`Expected: ${bugReport.expectedBehavior}`);
    console.log(`Actual: ${bugReport.actualBehavior}`);
    console.log(`${'='.repeat(80)}\n`);
  }
}

module.exports = BugReporter;
