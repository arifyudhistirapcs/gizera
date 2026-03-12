# Playwright Web Testing System

Comprehensive QA automation framework for testing PWA frontend with headed browser mode using Playwright.

## Features

- ✅ Headed browser testing with Google Chrome for visual observation
- ✅ 22 application modules covered with comprehensive test cases
- ✅ Automatic bug report generation for failed tests
- ✅ Test result reporting with statistics and summaries
- ✅ Module-based test organization
- ✅ Screenshot capture on test failures
- ✅ Configurable test execution parameters

## Installation

### Prerequisites

- Node.js (v16 or higher)
- npm or yarn

### Setup

1. Navigate to the tests directory:
```bash
cd tests
```

2. Install dependencies:
```bash
npm install
```

3. Install Playwright browsers:
```bash
npm run install:browsers
```

## Configuration

### Test Configuration File

Edit `config/test.config.json` to configure:

```json
{
  "pwaBaseUrl": "http://localhost:5173",
  "backendBaseUrl": "https://your-backend-url/api/v1",
  "browser": {
    "headless": false,
    "slowMo": 100,
    "viewport": { "width": 1920, "height": 1080 }
  },
  "timeouts": {
    "default": 30000,
    "navigation": 60000,
    "action": 10000
  },
  "screenshots": {
    "onFailure": true,
    "path": "screenshots"
  }
}
```

### Environment Variables

You can override configuration using environment variables:

```bash
export PWA_BASE_URL=http://localhost:5173
export BACKEND_BASE_URL=https://your-backend-url/api/v1
export HEADLESS=false
```

## Running Tests

### Run All Tests

```bash
node run-tests.js
```

### Run Specific Module

```bash
node run-tests.js --module authentication
node run-tests.js --module dashboard
node run-tests.js --module menu-perencanaan
```

### Available Modules

- `authentication` - User login/logout
- `dashboard` - Main dashboard
- `monitoring-aktivitas` - Activity monitoring
- `ulasan-rating` - Reviews and ratings
- `display-kds` - Kitchen Display System
- `menu-perencanaan` - Menu planning
- `menu-manajemen` - Menu management
- `menu-komponen` - Component management
- `supply-chain-supplier` - Supplier management
- `supply-chain-purchase-order` - Purchase orders
- `supply-chain-penerimaan-barang` - Goods receipt
- `supply-chain-bahan-baku` - Raw materials
- `logistik-data-sekolah` - School data
- `logistik-tugas-pengiriman` - Delivery tasks
- `sdm-data-karyawan` - Employee data
- `sdm-laporan-absensi` - Attendance reports
- `sdm-konfigurasi-absensi` - Attendance configuration
- `keuangan-aset-dapur` - Kitchen assets
- `keuangan-arus-kas` - Cash flow
- `keuangan-laporan` - Financial reports
- `sistem-audit-trail` - Audit trail
- `sistem-konfigurasi` - System configuration

## Test Case Structure

Test cases are stored in JSON format under `test-cases/<module>/test-cases.json`:

```json
{
  "id": "auth-001",
  "module": "authentication",
  "scenario": "User login with valid credentials",
  "steps": [
    "Navigate to login page",
    "Enter valid username",
    "Enter valid password",
    "Click login button"
  ],
  "expectedResults": [
    "User is redirected to dashboard",
    "User profile is displayed in header"
  ],
  "actualResults": [],
  "status": "not_run",
  "lastExecuted": null,
  "executionTime": null,
  "tags": ["critical", "smoke"]
}
```

## Test Results

### Test Reports

Test reports are saved to `test-results/` directory:
- JSON format: `report-<timestamp>.json`
- Contains statistics, module results, and bug report references

### Bug Reports

Bug reports are automatically generated for failed tests and saved to `bug-reports/`:
- Format: `bug-<timestamp>.json`
- Includes test case details, reproduction steps, screenshots, and error messages

### Screenshots

Screenshots are captured on test failures and saved to `screenshots/`:
- Format: `<timestamp>-<test-id>.png`

## Directory Structure

```
tests/
├── config/
│   ├── playwright.config.js      # Playwright configuration
│   └── test.config.json          # Test execution settings
├── test-cases/
│   ├── authentication/           # Auth module tests
│   ├── dashboard/                # Dashboard tests
│   └── ...                       # Other modules
├── test-suites/                  # Playwright test suites (future)
├── bug-reports/                  # Generated bug reports
├── test-results/                 # Test execution reports
├── screenshots/                  # Failure screenshots
├── utils/
│   ├── browser-manager.js        # Browser lifecycle management
│   ├── test-executor.js          # Test execution engine
│   ├── bug-reporter.js           # Bug documentation
│   ├── report-generator.js       # Test report generation
│   ├── config-loader.js          # Configuration loader
│   └── test-case-loader.js       # Test case management
└── run-tests.js                  # Main test runner

```

## Adding New Test Cases

1. Navigate to the module directory: `test-cases/<module>/`
2. Edit `test-cases.json`
3. Add new test case following the structure above
4. Run tests for that module

## Troubleshooting

### Browser Not Launching

If Chrome doesn't launch:
```bash
npx playwright install chromium
```

### Configuration Errors

Verify your `test.config.json` has valid URLs and settings.

### Test Failures

Check:
1. PWA is running at configured URL
2. Backend is accessible
3. Test data is available
4. Screenshots in `screenshots/` directory for visual debugging

## Best Practices

1. **Run tests with PWA and backend running**
2. **Review bug reports** after test execution
3. **Keep test cases updated** when features change
4. **Use headed mode** for debugging and observation
5. **Check screenshots** for visual verification of failures

## Support

For issues or questions, review:
- Bug reports in `bug-reports/`
- Test results in `test-results/`
- Screenshots in `screenshots/`
- Configuration in `config/test.config.json`
