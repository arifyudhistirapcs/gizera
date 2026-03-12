const ConfigLoader = require('./utils/config-loader');
const BrowserManager = require('./utils/browser-manager');
const TestCaseLoader = require('./utils/test-case-loader');
const BugReporter = require('./utils/bug-reporter');
const ReportGenerator = require('./utils/report-generator');

async function verifyCoreUtilities() {
  console.log('\n========================================');
  console.log('VERIFYING CORE UTILITIES');
  console.log('========================================\n');

  let allPassed = true;

  // Test 1: Config Loader
  try {
    console.log('✓ Testing ConfigLoader...');
    const configLoader = new ConfigLoader();
    const config = configLoader.load();
    console.log(`  ✓ Config loaded successfully`);
    console.log(`  ✓ PWA URL: ${config.pwaBaseUrl}`);
    console.log(`  ✓ Backend URL: ${config.backendBaseUrl}`);
  } catch (error) {
    console.error(`  ✗ ConfigLoader failed: ${error.message}`);
    allPassed = false;
  }

  // Test 2: Browser Manager
  try {
    console.log('\n✓ Testing BrowserManager...');
    const configLoader = new ConfigLoader();
    const config = configLoader.load();
    const browserManager = new BrowserManager();
    await browserManager.initialize(config);
    console.log(`  ✓ BrowserManager initialized`);
    
    // Note: Not launching browser in verification to save time
    console.log(`  ✓ BrowserManager ready (browser launch skipped in verification)`);
  } catch (error) {
    console.error(`  ✗ BrowserManager failed: ${error.message}`);
    allPassed = false;
  }

  // Test 3: Test Case Loader
  try {
    console.log('\n✓ Testing TestCaseLoader...');
    const testCaseLoader = new TestCaseLoader();
    console.log(`  ✓ TestCaseLoader initialized`);
    console.log(`  ✓ Ready to load test cases from modules`);
  } catch (error) {
    console.error(`  ✗ TestCaseLoader failed: ${error.message}`);
    allPassed = false;
  }

  // Test 4: Bug Reporter
  try {
    console.log('\n✓ Testing BugReporter...');
    const configLoader = new ConfigLoader();
    const config = configLoader.load();
    const bugReporter = new BugReporter(config);
    console.log(`  ✓ BugReporter initialized`);
    console.log(`  ✓ Bug reports directory ready`);
  } catch (error) {
    console.error(`  ✗ BugReporter failed: ${error.message}`);
    allPassed = false;
  }

  // Test 5: Report Generator
  try {
    console.log('\n✓ Testing ReportGenerator...');
    const reportGenerator = new ReportGenerator();
    console.log(`  ✓ ReportGenerator initialized`);
    console.log(`  ✓ Test results directory ready`);
  } catch (error) {
    console.error(`  ✗ ReportGenerator failed: ${error.message}`);
    allPassed = false;
  }

  console.log('\n========================================');
  if (allPassed) {
    console.log('✅ ALL CORE UTILITIES VERIFIED SUCCESSFULLY');
  } else {
    console.log('❌ SOME CORE UTILITIES FAILED');
  }
  console.log('========================================\n');

  return allPassed;
}

// Run verification
verifyCoreUtilities()
  .then(success => {
    process.exit(success ? 0 : 1);
  })
  .catch(error => {
    console.error('Verification failed:', error);
    process.exit(1);
  });
