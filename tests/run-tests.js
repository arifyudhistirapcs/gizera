#!/usr/bin/env node

const ConfigLoader = require('./utils/config-loader');
const BrowserManager = require('./utils/browser-manager');
const TestCaseLoader = require('./utils/test-case-loader');
const TestExecutor = require('./utils/test-executor');
const BugReporter = require('./utils/bug-reporter');
const ReportGenerator = require('./utils/report-generator');

class TestOrchestrator {
  constructor() {
    this.configLoader = new ConfigLoader();
    this.config = null;
    this.browserManager = null;
    this.testCaseLoader = null;
    this.testExecutor = null;
    this.bugReporter = null;
    this.reportGenerator = null;
    this.bugReports = [];
  }

  async initialize() {
    console.log('\n🚀 Initializing Playwright Web Testing System...\n');
    
    // Load configuration
    this.config = this.configLoader.load();
    console.log('✓ Configuration loaded');
    console.log(`  PWA URL: ${this.config.pwaBaseUrl}`);
    console.log(`  Backend URL: ${this.config.backendBaseUrl}`);
    console.log(`  Headed Mode: ${!this.config.browser.headless}`);
    
    // Initialize components
    this.browserManager = new BrowserManager();
    await this.browserManager.initialize(this.config);
    console.log('✓ Browser Manager initialized');
    
    this.testCaseLoader = new TestCaseLoader();
    console.log('✓ Test Case Loader initialized');
    
    this.bugReporter = new BugReporter(this.config);
    console.log('✓ Bug Reporter initialized');
    
    this.reportGenerator = new ReportGenerator();
    console.log('✓ Report Generator initialized');
    
    console.log('\n');
  }

  async runTests(modules = null) {
    try {
      // Launch browser
      console.log('🌐 Launching Chrome browser in headed mode...\n');
      await this.browserManager.launchBrowser(this.config.browser.headless);
      await this.browserManager.createContext();
      
      // Initialize test executor
      this.testExecutor = new TestExecutor(this.browserManager, this.testCaseLoader);
      
      // Determine which modules to test
      const modulesToTest = modules || this.config.modules;
      console.log(`📋 Testing ${modulesToTest.length} modules:\n`);
      modulesToTest.forEach((mod, idx) => {
        console.log(`   ${idx + 1}. ${mod}`);
      });
      console.log('\n');
      
      // Execute tests for each module
      const allResults = [];
      for (const module of modulesToTest) {
        try {
          const result = await this.testExecutor.executeTestSuite(module);
          allResults.push(result);
          
          // Create bug reports for failed tests
          await this.createBugReportsForFailures(result);
          
        } catch (error) {
          console.error(`❌ Failed to execute tests for ${module}: ${error.message}\n`);
        }
      }
      
      // Generate test report
      const report = await this.reportGenerator.generateTestReport(allResults);
      
      // Link bug reports
      report.bugReports = await this.reportGenerator.linkBugReports(allResults, this.bugReports);
      
      // Save report
      await this.reportGenerator.saveReport(report);
      
      // Display summary
      this.displaySummary(report);
      
      // Handle bug reports if any
      if (this.bugReports.length > 0) {
        await this.handleBugReports();
      }
      
      return report;
      
    } finally {
      // Cleanup
      await this.browserManager.cleanup();
      console.log('\n✓ Browser closed\n');
    }
  }


  async createBugReportsForFailures(moduleResult) {
    for (const testResult of moduleResult.results) {
      if (testResult.status === 'fail') {
        // Find the original test case
        const testCases = this.testCaseLoader.loadTestCases(moduleResult.module);
        const testCase = testCases.find(tc => tc.id === testResult.id);
        
        if (testCase) {
          const bugReport = await this.bugReporter.createBugReport(testCase, testResult);
          await this.bugReporter.saveBugReport(bugReport);
          this.bugReports.push(bugReport);
        }
      }
    }
  }

  async handleBugReports() {
    console.log(`\n${'='.repeat(80)}`);
    console.log(`🐛 ${this.bugReports.length} BUG REPORT(S) GENERATED`);
    console.log(`${'='.repeat(80)}\n`);
    
    for (const bugReport of this.bugReports) {
      await this.bugReporter.presentBugReport(bugReport);
    }
    
    console.log(`\nBug reports have been saved to: bug-reports/\n`);
    console.log(`To fix bugs, review the reports and run bug fix workflow.\n`);
  }

  displaySummary(report) {
    const formatted = this.reportGenerator.formatHumanReadable(report);
    console.log(formatted);
  }

  async runModule(moduleName) {
    return await this.runTests([moduleName]);
  }
}

// CLI Interface
async function main() {
  const args = process.argv.slice(2);
  const orchestrator = new TestOrchestrator();
  
  try {
    await orchestrator.initialize();
    
    if (args.length === 0) {
      // Run all tests
      console.log('Running all tests...\n');
      await orchestrator.runTests();
    } else if (args[0] === '--module' && args[1]) {
      // Run specific module
      console.log(`Running tests for module: ${args[1]}\n`);
      await orchestrator.runModule(args[1]);
    } else if (args[0] === '--help' || args[0] === '-h') {
      console.log(`
Playwright Web Testing System

Usage:
  node run-tests.js                    Run all tests
  node run-tests.js --module <name>    Run tests for specific module
  node run-tests.js --help             Show this help message

Examples:
  node run-tests.js
  node run-tests.js --module authentication
  node run-tests.js --module dashboard

Available modules:
  - authentication
  - dashboard
  - monitoring-aktivitas
  - ulasan-rating
  - display-kds
  - menu-perencanaan
  - menu-manajemen
  - menu-komponen
  - supply-chain-supplier
  - supply-chain-purchase-order
  - supply-chain-penerimaan-barang
  - supply-chain-bahan-baku
  - logistik-data-sekolah
  - logistik-tugas-pengiriman
  - sdm-data-karyawan
  - sdm-laporan-absensi
  - sdm-konfigurasi-absensi
  - keuangan-aset-dapur
  - keuangan-arus-kas
  - keuangan-laporan
  - sistem-audit-trail
  - sistem-konfigurasi
      `);
    } else {
      console.error('Invalid arguments. Use --help for usage information.');
      process.exit(1);
    }
    
    process.exit(0);
  } catch (error) {
    console.error(`\n❌ Error: ${error.message}\n`);
    console.error(error.stack);
    process.exit(1);
  }
}

// Run if called directly
if (require.main === module) {
  main();
}

module.exports = TestOrchestrator;
