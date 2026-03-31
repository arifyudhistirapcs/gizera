# Bugfix Requirements Document

## Introduction

This bugfix addresses a status value mismatch between backend and frontend that prevents Kepala SPPG from approving menu planning items. The backend creates menu plans with status 'draft', but the frontend PWA only displays approve/reject buttons when status is 'pending', resulting in newly created draft menus being inaccessible for approval.

## Bug Analysis

### Current Behavior (Defect)

1.1 WHEN a menu plan has status 'draft' THEN the system does not display approve/reject buttons in MenuPlanningView.vue

1.2 WHEN a menu plan has status 'draft' THEN the system does not display approve/reject buttons in MenuWeekCard.vue

1.3 WHEN Kepala SPPG views newly created menu plans THEN the system prevents approval workflow due to missing action buttons

### Expected Behavior (Correct)

2.1 WHEN a menu plan has status 'draft' THEN the system SHALL display approve/reject buttons in MenuPlanningView.vue

2.2 WHEN a menu plan has status 'draft' THEN the system SHALL display approve/reject buttons in MenuWeekCard.vue

2.3 WHEN Kepala SPPG clicks approve on a draft menu THEN the system SHALL change status to 'approved'

2.4 WHEN Kepala SPPG clicks reject on a draft menu THEN the system SHALL change status to 'rejected'

### Unchanged Behavior (Regression Prevention)

3.1 WHEN a menu plan has status 'approved' THEN the system SHALL CONTINUE TO display it as approved without action buttons

3.2 WHEN a menu plan has status 'rejected' THEN the system SHALL CONTINUE TO display it as rejected without action buttons

3.3 WHEN a menu plan has status 'pending' (if any exist) THEN the system SHALL CONTINUE TO display approve/reject buttons

3.4 WHEN the menuPlanning store maps status from backend THEN the system SHALL CONTINUE TO correctly map all status values

3.5 WHEN backend creates menu plans with 'draft' status THEN the system SHALL CONTINUE TO use 'draft' as the initial status
