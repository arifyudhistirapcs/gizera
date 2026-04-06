---
inclusion: fileMatch
fileMatchPattern: "backend/**/*_test.go,web/src/**/*.test.js,pwa/src/**/*.test.js"
---

# Testing Guide

## Go Backend Tests
- Use `testify` for assertions: `require.NoError`, `assert.Equal`
- Table-driven tests for multiple scenarios
- Property-based tests with `pgregory.net/rapid` (min 100 iterations)
- Use SQLite in-memory for unit tests: `gorm.Open(sqlite.Open(":memory:"))`
- Always test: happy path, error cases, edge cases, tenant isolation

## Vue Tests (Vitest)
- Run: `npm run test` (web) or `npm run test` (pwa)
- Use `@vue/test-utils` for component testing
- Mock API calls with `vi.mock()`

## Verification Commands
- Backend compile: `cd backend && go build ./internal/...`
- Backend vet: `cd backend && go vet ./internal/...`
- Backend tests: `cd backend && go test ./internal/... -v`
- Web tests: `cd web && npm run test`
- PWA tests: `cd pwa && npm run test`
