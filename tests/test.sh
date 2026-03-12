#!/bin/bash

# Playwright Web Testing - Test Runner Script

echo "🎭 Playwright Web Testing System"
echo "================================"
echo ""

# Check if PWA is running
echo "Checking if PWA is running..."
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo "✓ PWA is running at http://localhost:5173"
else
    echo "✗ PWA is not running!"
    echo "  Please start PWA first: cd pwa && npm run dev"
    exit 1
fi

echo ""

# Parse command line arguments
if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    echo "Usage:"
    echo "  ./test.sh                    Run all tests in headed mode"
    echo "  ./test.sh auth               Run authentication tests"
    echo "  ./test.sh --debug            Run in debug mode"
    echo "  ./test.sh --headless         Run in headless mode"
    echo "  ./test.sh --report           Show HTML report"
    echo ""
    echo "Examples:"
    echo "  ./test.sh                    # Run all tests"
    echo "  ./test.sh auth               # Run auth tests only"
    echo "  ./test.sh --debug            # Debug mode"
    exit 0
fi

if [ "$1" == "--report" ]; then
    echo "Opening HTML report..."
    npx playwright show-report test-results/html-report
    exit 0
fi

# Determine test mode
if [ "$1" == "--debug" ]; then
    echo "Running in DEBUG mode..."
    npx playwright test --debug
elif [ "$1" == "--headless" ]; then
    echo "Running in HEADLESS mode..."
    npx playwright test
elif [ "$1" == "auth" ] || [ "$1" == "authentication" ]; then
    echo "Running AUTHENTICATION tests in HEADED mode..."
    npx playwright test authentication --headed
elif [ -z "$1" ]; then
    echo "Running ALL tests in HEADED mode..."
    npx playwright test --headed
else
    echo "Running tests matching: $1"
    npx playwright test "$1" --headed
fi

echo ""
echo "================================"
echo "Test execution completed!"
echo ""
echo "View results:"
echo "  - HTML Report: npx playwright show-report test-results/html-report"
echo "  - Screenshots: ls screenshots/"
echo "  - JSON Results: cat test-results/test-results.json"
echo ""
