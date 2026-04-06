---
inclusion: auto
---

# Coding Standards

## Go Backend

### Naming
- `camelCase` unexported, `PascalCase` exported, `UPPER_SNAKE_CASE` constants
- Package names: short, lowercase, singular

### Functions
- `context.Context` as first param, `error` as last return
- Max 50 lines per function, early returns, single responsibility

### Error Handling
- Sentinel errors: `var ErrNotFound = errors.New("...")`
- Wrap with context: `fmt.Errorf("creating form: %w", err)`
- Never swallow errors silently

### GORM Models
- Use `gorm:"primaryKey"`, `json:"field_name"`, `validate:"required"`
- Pointers for nullable: `*string`, `*time.Time`, `*uint`
- `json:"-"` for passwords, tokens
- `TableName()` method for every model
- Add to `AllModels()` in `models.go` for auto-migration

### Handlers
- Bind JSON → validate → call service → return standard response
- Extract `yayasan_id` from context for tenant-scoped operations
- Map service errors to HTTP status codes consistently

### Testing
- Table-driven tests with `testify`
- Property-based tests with `pgregory.net/rapid`
- Test file: `*_test.go` in same package

## Vue 3 (Web + PWA)

### Components
- `<script setup>` only (Composition API)
- `ref()`, `computed()`, `watch()` for reactivity
- Ant Design Vue components for web, Vant UI for PWA

### API Calls
- Always through `services/*.js` using shared Axios instance
- Handle loading, error, empty states in every view
- Web uses Vite proxy (`/api` → backend), PWA uses direct URL from `.env`

### State
- Pinia for auth and shared state only
- Local `ref()` for component state
- Never duplicate server state in Pinia

### Naming
- PascalCase for components and views
- camelCase for functions, variables, composables
- Files: `PascalCase.vue` for views, `camelCase.js` for services
