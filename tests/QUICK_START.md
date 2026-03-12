# Quick Start Guide - Playwright Web Testing

## Prerequisites

1. **PWA must be running** at `http://localhost:5173`
2. **Backend must be accessible** at configured URL
3. **Node.js** installed (v16+)

## Installation (Already Done)

```bash
cd tests
npm install
npm run install:browsers
```

## Running Tests

### Option 1: Run with Custom Test Runner (Recommended for QA)

This will run tests with headed browser and generate bug reports:

```bash
# Run all modules
node run-tests.js

# Run specific module
node run-tests.js --module authentication
```

### Option 2: Run with Playwright Test Runner

This uses Playwright's built-in test runner:

```bash
# Run all test suites
npx playwright test

# Run specific test suite
npx playwright test authentication

# Run in headed mode (see browser)
npx playwright test --headed

# Run in debug mode
npx playwright test --debug

# Run specific test
npx playwright test -g "auth-001"
```

## First Time Setup

### 1. Update Test Credentials

Edit `test-suites/authentication.spec.js` and replace:
```javascript
await usernameInput.fill('admin@example.com'); // Your actual username
await passwordInput.fill('password123');        // Your actual password
```

### 2. Verify Configuration

Check `config/test.config.json`:
```json
{
  "pwaBaseUrl": "http://localhost:5173",  // Your PWA URL
  "backendBaseUrl": "https://your-backend-url/api/v1"
}
```

### 3. Start Your Application

```bash
# Terminal 1: Start backend
cd backend
go run cmd/server/main.go

# Terminal 2: Start PWA
cd pwa
npm run dev

# Terminal 3: Run tests
cd tests
npx playwright test --headed
```

## Quick Test Run

Test authentication module only:

```bash
npx playwright test authentication --headed
```

You should see:
- ✅ Chrome browser opens
- ✅ Tests run visually
- ✅ Results in terminal
- ✅ Screenshots for failures in `screenshots/`
- ✅ HTML report generated

## View Test Results

### HTML Report

```bash
npx playwright show-report test-results/html-report
```

### JSON Results

```bash
cat test-results/test-results.json | json_pp
```

### Screenshots

Failed tests automatically capture screenshots:
```bash
ls -la screenshots/
```

## Debugging Failed Tests

### 1. Run in Debug Mode

```bash
npx playwright test authentication --debug
```

This opens Playwright Inspector where you can:
- Step through test
- Inspect elements
- See console logs
- Pause/resume execution

### 2. Check Screenshots

```bash
open screenshots/  # macOS
xdg-open screenshots/  # Linux
explorer screenshots\  # Windows
```

### 3. Check Bug Reports

```bash
cat bug-reports/*.json | json_pp
```

## Common Issues

### Issue: Browser doesn't launch

**Solution:**
```bash
npx playwright install chromium
```

### Issue: Tests fail with "element not found"

**Solution:**
1. Check if PWA is running: `curl http://localhost:5173`
2. Run in headed mode to see what's happening: `npx playwright test --headed`
3. Update selectors in test suite if UI changed

### Issue: Login credentials don't work

**Solution:**
1. Update credentials in `test-suites/authentication.spec.js`
2. Or create test user in your application
3. Verify backend is accessible

### Issue: Tests timeout

**Solution:**
1. Increase timeout in `config/test.config.json`:
```json
{
  "timeouts": {
    "default": 60000,
    "navigation": 120000
  }
}
```

## Next Steps

### 1. Customize Authentication Tests

Edit `test-suites/authentication.spec.js`:
- Update selectors to match your actual UI
- Add more test scenarios
- Update test credentials

### 2. Create More Test Suites

Copy `test-suites/authentication.spec.js` as template:
```bash
cp test-suites/authentication.spec.js test-suites/dashboard.spec.js
```

Then customize for dashboard module.

### 3. Run Tests in CI/CD

Add to your CI pipeline:
```yaml
- name: Run E2E Tests
  run: |
    cd tests
    npx playwright test --reporter=json
```

## Tips for Effective Testing

1. **Run tests in headed mode first** to see what's happening
2. **Use Playwright Inspector** for debugging: `--debug` flag
3. **Check screenshots** when tests fail
4. **Update selectors** if UI changes
5. **Keep test data separate** from production data
6. **Run tests regularly** to catch regressions early

## Example Test Session

```bash
# 1. Start your app
cd pwa && npm run dev

# 2. In another terminal, run tests
cd tests
npx playwright test authentication --headed

# 3. Watch tests run in browser
# 4. Check results in terminal
# 5. View HTML report
npx playwright show-report test-results/html-report
```

## Getting Help

- Check `README.md` for detailed documentation
- Review test cases in `test-cases/authentication/test-cases.json`
- Look at example test suite in `test-suites/authentication.spec.js`
- Check Playwright docs: https://playwright.dev

## Success Indicators

You'll know it's working when:
- ✅ Browser opens automatically
- ✅ You see login page
- ✅ Tests interact with UI
- ✅ Results show in terminal
- ✅ Screenshots captured on failures
- ✅ HTML report generated

Happy Testing! 🎭
