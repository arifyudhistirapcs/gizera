# Implementation Plan: Risk Assessment (Audit Kepatuhan SOP)

## Overview

Implementasi fitur Risk Assessment secara inkremental: dimulai dari data models dan migrasi, lalu service layer (score calculation, snapshot capture), handler + routes, seed data, lalu frontend (web monitoring dashboard dan PWA form). Setiap task membangun di atas task sebelumnya. Property-based tests menggunakan `pgregory.net/rapid`.

## Tasks

- [x] 1. Define data models and register for auto-migration
  - [x] 1.1 Create risk assessment models file
    - Create `backend/internal/models/risk_assessment.go` with all 6 structs: `SOPCategory`, `SOPChecklistItem`, `RiskAssessmentForm`, `RiskAssessmentItem`, `RiskAssessmentCategoryScore`, `SPPGOperationalSnapshot`
    - Include GORM tags, JSON tags, and validation tags as specified in design
    - Add `TableName()` methods for each struct
    - _Requirements: 1.1, 1.2, 2.3, 7.1, 7.2, 8.2_

  - [x] 1.2 Register models in AllModels() for auto-migration
    - Add all 6 new model structs to `backend/internal/models/models.go` `AllModels()` function
    - Place them under a new comment section `// Risk Assessment`
    - _Requirements: 1.1_

- [x] 2. Implement RiskAssessmentService with score calculation and form lifecycle
  - [x] 2.1 Create RiskAssessmentService with form CRUD operations
    - Create `backend/internal/services/risk_assessment_service.go`
    - Implement `NewRiskAssessmentService(db *gorm.DB, snapshotService *SnapshotService)`
    - Implement `CreateForm(sppgID, yayasanID, createdByUserID uint)` — validates SPPG ownership, fetches active checklist items, creates form with items (snapshot item_nama/category_nama), calls SnapshotService
    - Implement `GetForm(formID, yayasanID uint)` — loads form with Items, CategoryScores, Snapshot, applies tenant filter
    - Implement `ListForms(filter FormFilter)` — paginated list with tenant filtering, ordered by created_at DESC
    - Implement `UpdateDraft(formID, yayasanID uint, items []UpdateItemRequest)` — validates form is draft, validates scores 1-5, updates items
    - Implement `GetStats(sppgIDs []uint)` — aggregated stats per SPPG (count, avg score, trend)
    - Define `FormFilter` struct with fields: YayasanID, SPPGID, Status, RiskLevel, DateFrom, DateTo, Page, PageSize
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.6, 3.6, 4.1, 4.4, 4.5, 5.1, 5.2, 6.1, 6.2_

  - [x] 2.2 Implement SubmitForm with score calculation and penalty logic
    - Implement `SubmitForm(formID, yayasanID uint)` in the same service file
    - Validate all items have non-nil ComplianceScore; return error listing unscored items
    - Implement `calculateCategoryScores(items []RiskAssessmentItem)` — group items by SOPCategoryID, compute average per category, determine risk level per category
    - Implement `calculateOverallScore(items []RiskAssessmentItem)` — arithmetic mean of all compliance scores
    - Implement `applyOperationalPenalty(sopScore float64, snapshot *SPPGOperationalSnapshot)` — apply -0.5 per trigger (review < 3.0, budget absorption < 50%, on-time delivery < 70%), max -1.5, clamp result to [1.0, 5.0]
    - Implement `determineRiskLevel(score float64)` — "rendah" [4.0-5.0], "sedang" [2.5-3.9], "tinggi" [1.0-2.4]
    - Set form status to "submitted", record SubmittedAt timestamp
    - Save category scores to `risk_assessment_category_scores` table
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 7.1, 7.2, 7.4, 8.5_

  - [ ]* 2.3 Write property tests for score calculation (Properties 9, 10, 11, 12)
    - **Property 9: Category Score Calculation Accuracy** — generate random items with scores 1-5 across multiple categories, verify average per category matches arithmetic mean
    - **Validates: Requirements 7.1, 7.2**
    - **Property 10: Overall SOP Score Calculation** — generate random items, verify overall score equals arithmetic mean of all scores
    - **Validates: Requirements 3.3**
    - **Property 11: Risk Level Mapping Determinism** — generate random scores in [1.0, 5.0], verify deterministic risk level assignment per threshold
    - **Validates: Requirements 3.4, 7.4**
    - **Property 12: Operational Penalty Calculation** — generate random snapshot values, verify penalty is 0.5 per trigger, max 1.5, final score clamped to [1.0, 5.0]
    - **Validates: Requirements 8.5**
    - Create test file `backend/internal/services/risk_assessment_service_test.go` using `pgregory.net/rapid`

  - [ ]* 2.4 Write property test for form lifecycle (Property 13)
    - **Property 13: Form Status Lifecycle** — generate forms, verify draft→submitted transition, verify submitted forms reject updates
    - **Validates: Requirements 2.6, 3.5, 3.6**

  - [ ]* 2.5 Write property test for compliance score validation (Property 7)
    - **Property 7: Compliance Score Validation** — generate random integers, verify only [1,5] accepted, values outside range rejected
    - **Validates: Requirements 2.3**

  - [ ]* 2.6 Write property test for submit requires complete scores (Property 8)
    - **Property 8: Submit Requires Complete Scores** — generate forms with some nil scores, verify submission rejected and form stays draft
    - **Validates: Requirements 3.1, 3.2**

- [x] 3. Implement SnapshotService for capturing SPPG operational data
  - [x] 3.1 Create SnapshotService
    - Create `backend/internal/services/snapshot_service.go`
    - Implement `NewSnapshotService(db *gorm.DB, aggregatedService *AggregatedDashboardService)`
    - Implement `CaptureSnapshot(sppgID uint)` that:
      - Uses `aggregatedService.getAggregatedReview()` for review metrics
      - Uses `aggregatedService.getAggregatedFinancial()` for financial metrics
      - Uses `aggregatedService.getAggregatedDelivery()` for delivery metrics
      - Uses `aggregatedService.getAggregatedProduction()` for production metrics
      - Queries inventory alerts directly (items below min threshold)
      - Queries active employee count and attendance rate directly
      - Sets SnapshotPeriodStart to first day of current month, SnapshotPeriodEnd and CapturedAt to now
    - Return populated `SPPGOperationalSnapshot` struct
    - _Requirements: 8.1, 8.2, 8.3_

  - [ ]* 3.2 Write property test for snapshot completeness (Property 15)
    - **Property 15: Operational Snapshot Completeness** — verify snapshot has valid timestamps and all metric fields populated when SPPG has data
    - **Validates: Requirements 8.1, 8.2**

  - [ ]* 3.3 Write property test for snapshot immutability (Property 14)
    - **Property 14: Operational Snapshot Immutability** — verify snapshot values don't change when underlying SPPG data changes after capture
    - **Validates: Requirements 8.3**

- [x] 4. Checkpoint - Ensure all backend service tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [x] 5. Implement SOP template management in RiskAssessmentService
  - [x] 5.1 Add SOP category and checklist item CRUD methods
    - Add to `risk_assessment_service.go`:
    - `GetSOPCategories()` — list all categories ordered by Urutan
    - `CreateSOPCategory(input)` — create with auto-assigned Urutan
    - `UpdateSOPCategory(id, input)` — update nama, deskripsi, urutan
    - `GetSOPChecklistItems(categoryID *uint)` — list items, optionally filtered by category, ordered by Urutan
    - `CreateSOPChecklistItem(input)` — create with auto-assigned Urutan (max existing + 1 in category)
    - `UpdateSOPChecklistItem(id, input)` — update without affecting existing forms
    - `SetSOPChecklistItemStatus(id, isActive bool)` — activate/deactivate item
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

  - [ ]* 5.2 Write property test for auto-assigned display order (Property 18)
    - **Property 18: Auto-Assigned Display Order** — generate random items added to a category, verify each new item gets urutan > all existing items in same category
    - **Validates: Requirements 1.3**

  - [ ]* 5.3 Write property tests for form creation and checklist items (Properties 2, 3)
    - **Property 2: Form Creation Contains All Active Checklist Items** — generate sets of active/inactive items, verify form contains exactly active items
    - **Validates: Requirements 1.5, 2.2**
    - **Property 3: Checklist Item Snapshot Immutability** — create form, update source item nama, verify form item_nama unchanged
    - **Validates: Requirements 1.4**

- [x] 6. Implement RiskAssessmentHandler with all endpoints
  - [x] 6.1 Create RiskAssessmentHandler
    - Create `backend/internal/handlers/risk_assessment_handler.go`
    - Implement `NewRiskAssessmentHandler(service *RiskAssessmentService)`
    - Implement SOP template endpoints: `GetSOPCategories`, `CreateSOPCategory`, `UpdateSOPCategory`, `GetSOPChecklistItems`, `CreateSOPChecklistItem`, `UpdateSOPChecklistItem`, `SetSOPChecklistItemStatus`
    - Implement form endpoints: `CreateForm`, `GetForms`, `GetForm`, `UpdateDraft`, `SubmitForm`, `UploadEvidence`, `GetStats`
    - Extract yayasan_id from tenant middleware context for all form operations
    - Use consistent API response format: `{"success": bool, "data": ..., "message": ..., "error_code": ...}`
    - Handle file upload for evidence photos (save to `./uploads/risk-assessment/`)
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

  - [ ]* 6.2 Write property test for API response format consistency (Property 16)
    - **Property 16: API Response Format Consistency** — verify all endpoints return standard response format with success, data/error_code, message fields
    - **Validates: Requirements 6.3, 6.4**

  - [ ]* 6.3 Write property test for form list ordering (Property 17)
    - **Property 17: Form List Ordering** — create forms with various timestamps, verify list returns newest first
    - **Validates: Requirements 4.1**

- [x] 7. Register routes and wire services in router.go
  - [x] 7.1 Wire services and register risk assessment routes
    - In `backend/internal/router/router.go`:
    - Instantiate `SnapshotService` with existing `aggregatedDashboardService`
    - Instantiate `RiskAssessmentService` with db and snapshotService
    - Instantiate `RiskAssessmentHandler` with riskAssessmentService
    - Register all routes under `/risk-assessment` group as specified in design:
      - SOP template write routes with `middleware.RequireRole("superadmin")`
      - SOP template read routes with `middleware.RequireRole("kepala_yayasan", "superadmin")`
      - Form routes with `middleware.RequireRole("kepala_yayasan", "superadmin")`
      - Stats route with `middleware.RequireRole("kepala_yayasan", "superadmin")`
    - _Requirements: 5.5, 6.1_

  - [ ]* 7.2 Write property tests for tenant isolation (Properties 4, 5, 6)
    - **Property 4: Tenant Isolation for Kepala Yayasan** — verify all returned forms belong to user's Yayasan SPPGs
    - **Validates: Requirements 2.1, 4.1, 5.2**
    - **Property 5: Cross-Tenant Access Returns 404** — verify accessing forms outside tenant scope returns 404
    - **Validates: Requirements 5.4**
    - **Property 6: SPPG Ownership Validation on Form Creation** — verify form creation rejected for SPPGs outside user's Yayasan
    - **Validates: Requirements 6.2**

- [x] 8. Create seed data for default SOP categories
  - [x] 8.1 Add SOP seed data to seed command
    - Create `backend/cmd/seed/seed_sop_categories.go` or add to existing seed main.go
    - Seed 7 default SOP categories with checklist items based on SOP Dapur MBG document:
      1. Higienitas Dapur dan Sanitasi
      2. Standar Persiapan Makanan
      3. Penyimpanan dan Kontrol Suhu
      4. Prosedur Pengiriman
      5. Kebersihan Staf dan APD
      6. Pemeliharaan Peralatan
      7. Dokumentasi dan Pencatatan
    - Each category should have 3-5 representative checklist items
    - Use upsert logic to avoid duplicates on re-run
    - _Requirements: 1.6_

  - [ ]* 8.2 Write unit tests for seed data verification
    - Verify all 7 categories are created with correct names
    - Verify each category has associated checklist items
    - Verify urutan values are sequential
    - _Requirements: 1.6_

- [x] 9. Checkpoint - Ensure all backend tests pass and API endpoints work
  - Ensure all tests pass, ask the user if questions arise.

- [x] 10. Implement Web App monitoring dashboard
  - [x] 10.1 Create risk assessment API service and types
    - Create `web/src/api/riskAssessment.ts` with API client functions for all risk assessment endpoints
    - Create `web/src/types/riskAssessment.ts` with TypeScript interfaces matching backend models
    - _Requirements: 4.1, 6.1_

  - [x] 10.2 Create RiskAssessmentListView page
    - Create `web/src/views/RiskAssessmentListView.vue`
    - Display table of risk assessment forms with columns: SPPG name, audit date, overall risk score, risk level (color-coded badge), status
    - Implement filters: SPPG dropdown, date range picker, risk level selector
    - Display summary statistics per SPPG: total audits, average score, risk level trend
    - Sorted by newest first
    - Use Ant Design Vue components (a-table, a-select, a-date-picker, a-tag, a-statistic)
    - _Requirements: 4.1, 4.2, 4.4, 4.5_

  - [x] 10.3 Create RiskAssessmentDetailView page
    - Create `web/src/views/RiskAssessmentDetailView.vue`
    - Display full form detail: header info (SPPG, date, score, risk level, status)
    - Display items grouped by SOP category with compliance scores, notes, evidence photos
    - Display category scores in visual format (progress bars or radar chart)
    - Display SPPG Operational Snapshot section with color-coded metrics (green/yellow/red based on thresholds from design)
    - Use Ant Design Vue components (a-descriptions, a-collapse, a-progress, a-image, a-tag)
    - _Requirements: 4.3, 7.3, 8.4, 8.6_

  - [x] 10.4 Register routes in web app router
    - Add risk assessment routes to `web/src/router/index.ts`:
      - `/risk-assessment` → RiskAssessmentListView
      - `/risk-assessment/:id` → RiskAssessmentDetailView
    - Add navigation menu item for Risk Assessment (visible to kepala_yayasan and superadmin roles)
    - _Requirements: 4.1_

- [x] 11. Implement PWA App form creation and filling
  - [x] 11.1 Create risk assessment API service and types for PWA
    - Create `pwa/src/api/riskAssessment.ts` with API client functions
    - Create `pwa/src/types/riskAssessment.ts` with TypeScript interfaces
    - _Requirements: 2.1, 6.1_

  - [x] 11.2 Create SPPG selection and form creation page
    - Create `pwa/src/views/RiskAssessmentSelectView.vue`
    - Display list of SPPGs under user's Yayasan
    - On SPPG selection, call POST /risk-assessment/forms to create new form
    - Navigate to form filling page on success
    - Use Vant UI components
    - _Requirements: 2.1, 2.2_

  - [x] 11.3 Create form filling page with score input and evidence upload
    - Create `pwa/src/views/RiskAssessmentFormView.vue`
    - Display checklist items grouped by SOP category (collapsible sections)
    - For each item: score selector (1-5 radio/stepper), notes text input, photo upload button
    - Show progress indicator (X of Y items scored)
    - Auto-save draft on score/note changes (debounced)
    - Submit button with validation (all items must have scores)
    - Display operational snapshot section (read-only)
    - Use Vant UI components (van-collapse, van-stepper, van-field, van-uploader, van-button, van-progress)
    - _Requirements: 2.2, 2.3, 2.4, 2.5, 2.6, 3.1, 3.2, 8.4_

  - [x] 11.4 Implement offline support with IndexedDB
    - Create `pwa/src/utils/riskAssessmentOffline.ts`
    - Implement IndexedDB store for `risk_assessment_drafts` and `risk_assessment_cache`
    - Save form data to IndexedDB on each auto-save
    - Detect offline status and switch to local storage
    - Implement background sync: on reconnect, push all `pending_sync` forms to API
    - Handle sync conflicts: server timestamp wins, skip duplicates
    - Show sync status indicator in UI
    - _Requirements: 2.6, 2.7_

  - [x] 11.5 Register routes in PWA router
    - Add risk assessment routes to PWA router:
      - `/risk-assessment` → RiskAssessmentSelectView
      - `/risk-assessment/:id` → RiskAssessmentFormView
    - Add navigation entry visible to kepala_yayasan role
    - _Requirements: 2.1_

- [ ] 12. Write property test for data round-trip (Property 1)
  - [ ]* 12.1 Write property test for risk assessment data round-trip
    - **Property 1: Risk Assessment Data Round-Trip** — generate random valid compliance scores (1-5), catatan strings, and evidence URLs; save and retrieve; verify exact match
    - **Validates: Requirements 1.1, 2.3, 2.4, 2.5**
    - Add to `backend/internal/services/risk_assessment_service_test.go`

- [x] 13. Final checkpoint - Ensure all tests pass and features are integrated
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property-based tests use `pgregory.net/rapid` library with minimum 100 iterations
- Checkpoints ensure incremental validation
- Backend tasks (1-9) should be completed before frontend tasks (10-11)
- The design specifies Go for backend and TypeScript/Vue for frontend — both are used accordingly
