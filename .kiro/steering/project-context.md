---
inclusion: auto
---

# Gizera ERP SPPG — Project Context

## Tech Stack
- Backend: Go (Gin + GORM + PostgreSQL) di `backend/`
- Web Dashboard: Vue 3 + Ant Design Vue + Vite di `web/`
- PWA Mobile: Vue 3 + Vant UI + Vite di `pwa/`
- Database: PostgreSQL
- Real-time: Firebase Realtime Database
- Maps: Leaflet + OpenStreetMap
- Offline: Dexie (IndexedDB) di PWA

## Architecture
- Backend: handler (`internal/handlers/`) → service (`internal/services/`) → model (`internal/models/`)
- Middleware: auth, tenant isolation, CORS, rate limiting, security (`internal/middleware/`)
- Web: Vue 3 SFC `<script setup>`, services/api.js (Axios), Pinia stores, Vite proxy `/api` → localhost:8080
- PWA: Vue 3 SFC `<script setup>`, services/api.js (Axios direct), Dexie for offline, Vant UI

## Multi-Tenant Hierarchy
BGN → Yayasan → SPPG (dapur)

## User Roles (13)
superadmin, admin_bgn, kepala_yayasan, kepala_sppg, akuntan, ahli_gizi, pengadaan, chef, packing, driver, asisten_lapangan, kebersihan, sekolah

## API Response Format
```json
{ "success": true, "data": {...}, "message": "..." }
{ "success": false, "error_code": "...", "message": "..." }
```

## Key Conventions
- Go: handler validates input → service handles business logic → model is GORM struct
- All tenant-scoped queries filter by `yayasan_id` or `sppg_id`
- GORM models use `json:"-"` for sensitive fields, pointers for nullable fields
- Vue components use `<script setup>` with Composition API
- API services in `web/src/services/` and `pwa/src/services/` follow same pattern
- Error handling: Go returns sentinel errors, handlers map to HTTP status + error_code
