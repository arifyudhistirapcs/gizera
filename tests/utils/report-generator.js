const fs = require('fs');
const path = require('path');

class ReportGenerator {
  constructor() {
    this.reportsDir = path.join(__dirname, '../test-results');
    
    // Ensure reports directory exists
    if (!fs.existsSync(this.reportsDir)) {
      fs.mkdirSync(this.reportsDir, { recursive: true });
    }
  }

  /**
   * Generate test report from test results
   */
  async generateTestReport(testResults) {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const reportId = `report-${timestamp}`;

    const statistics = this.calculateStatistics(testResults);

    const report = {
      id: reportId,
      timestamp: new Date().toISOString(),
      totalTests: statistics.totalTests,
      passed: statistics.passed,
      failed: statistics.failed,
      blocked: statistics.blocked,
      notRun: statistics.notRun,
      executionTime: statistics.totalExecutionTime,
      modules: statistics.moduleResults,
      bugReports: []
    };

    console.log(`\n📊 Test Report Generated: ${reportId}`);
    console.log(`   Total Tests: ${report.totalTests}`);
    console.log(`   Passed: ${report.passed}`);
    console.log(`   Failed: ${report.failed}`);
    console.log(`   Blocked: ${report.blocked}`);
    console.log(`   Execution Time: ${report.executionTime}ms`);

    return report;
  }

  /**
   * Calculate statistics from test results
   */
  calculateStatistics(testResults) {
    const stats = {
      totalTests: 0,
      passed: 0,
      failed: 0,
      blocked: 0,
      notRun: 0,
      totalExecutionTime: 0,
      moduleResults: []
    };

    // If testResults is an array of module results
    if (Array.isArray(testResults)) {
      testResults.forEach(moduleResult => {
        stats.totalTests += moduleResult.totalTests || 0;
        stats.passed += moduleResult.passed || 0;
        stats.failed += moduleResult.failed || 0;
        stats.blocked += moduleResult.blocked || 0;
        
        const moduleExecutionTime = moduleResult.results?.reduce((sum, r) => sum + (r.executionTime || 0), 0) || 0;
        stats.totalExecutionTime += moduleExecutionTime;

        stats.moduleResults.push({
          module: moduleResult.module,
          totalTests: moduleResult.totalTests || 0,
          passed: moduleResult.passed || 0,
          failed: moduleResult.failed || 0,
          blocked: moduleResult.blocked || 0,
          executionTime: moduleExecutionTime
        });
      });
    }
    // If testResults is an object with moduleResults
    else if (testResults.moduleResults) {
      return this.calculateStatistics(testResults.moduleResults);
    }

    return stats;
  }


  /**
   * Link bug reports to test results
   */
  async linkBugReports(testResults, bugReports) {
    const linkedBugReports = [];

    bugReports.forEach(bugReport => {
      // Find matching test result
      const matchingResult = this.findMatchingTestResult(testResults, bugReport.testCaseId);
      
      if (matchingResult) {
        linkedBugReports.push(bugReport.id);
      }
    });

    return linkedBugReports;
  }

  /**
   * Find matching test result for a test case ID
   */
  findMatchingTestResult(testResults, testCaseId) {
    if (Array.isArray(testResults)) {
      for (const moduleResult of testResults) {
        if (moduleResult.results) {
          const match = moduleResult.results.find(r => r.id === testCaseId);
          if (match) return match;
        }
      }
    }
    return null;
  }

  /**
   * Save report to file
   */
  async saveReport(report) {
    try {
      const filename = `${report.id}.json`;
      const filepath = path.join(this.reportsDir, filename);
      
      fs.writeFileSync(filepath, JSON.stringify(report, null, 2), 'utf8');
      console.log(`   Report saved: ${filepath}`);
      
      return filepath;
    } catch (error) {
      throw new Error(`Failed to save report: ${error.message}`);
    }
  }

  /**
   * Format report in different formats
   */
  async formatReport(report, format = 'json') {
    if (format === 'json') {
      return JSON.stringify(report, null, 2);
    }
    
    if (format === 'text' || format === 'human') {
      return this.formatHumanReadable(report);
    }

    throw new Error(`Unsupported format: ${format}`);
  }

  /**
   * Format report in human-readable text format
   */
  formatHumanReadable(report) {
    let output = '';
    
    output += `\n${'='.repeat(80)}\n`;
    output += `TEST EXECUTION REPORT\n`;
    output += `${'='.repeat(80)}\n`;
    output += `Report ID: ${report.id}\n`;
    output += `Timestamp: ${report.timestamp}\n`;
    output += `\n`;
    output += `SUMMARY:\n`;
    output += `  Total Tests: ${report.totalTests}\n`;
    output += `  Passed: ${report.passed} (${this.percentage(report.passed, report.totalTests)}%)\n`;
    output += `  Failed: ${report.failed} (${this.percentage(report.failed, report.totalTests)}%)\n`;
    output += `  Blocked: ${report.blocked} (${this.percentage(report.blocked, report.totalTests)}%)\n`;
    output += `  Not Run: ${report.notRun} (${this.percentage(report.notRun, report.totalTests)}%)\n`;
    output += `  Execution Time: ${report.executionTime}ms (${(report.executionTime / 1000).toFixed(2)}s)\n`;
    output += `\n`;
    
    if (report.modules && report.modules.length > 0) {
      output += `MODULE RESULTS:\n`;
      report.modules.forEach(module => {
        output += `\n  ${module.module}:\n`;
        output += `    Total: ${module.totalTests}\n`;
        output += `    Passed: ${module.passed}\n`;
        output += `    Failed: ${module.failed}\n`;
        output += `    Blocked: ${module.blocked}\n`;
        output += `    Time: ${module.executionTime}ms\n`;
      });
      output += `\n`;
    }
    
    if (report.bugReports && report.bugReports.length > 0) {
      output += `BUG REPORTS:\n`;
      report.bugReports.forEach(bugId => {
        output += `  - ${bugId}\n`;
      });
      output += `\n`;
    }
    
    output += `${'='.repeat(80)}\n`;
    
    return output;
  }

  /**
   * Calculate percentage
   */
  percentage(value, total) {
    if (total === 0) return 0;
    return ((value / total) * 100).toFixed(1);
  }
}

module.exports = ReportGenerator;
