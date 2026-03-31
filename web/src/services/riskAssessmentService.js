import api from './api'

const riskAssessmentService = {
  // SOP Category operations
  getSOPCategories() {
    return api.get('/risk-assessment/sop-categories')
  },

  createSOPCategory(data) {
    return api.post('/risk-assessment/sop-categories', data)
  },

  updateSOPCategory(id, data) {
    return api.put(`/risk-assessment/sop-categories/${id}`, data)
  },

  // SOP Checklist Item operations
  getSOPChecklistItems(params = {}) {
    return api.get('/risk-assessment/sop-checklist-items', { params })
  },

  createSOPChecklistItem(data) {
    return api.post('/risk-assessment/sop-checklist-items', data)
  },

  updateSOPChecklistItem(id, data) {
    return api.put(`/risk-assessment/sop-checklist-items/${id}`, data)
  },

  setSOPChecklistItemStatus(id, isActive) {
    return api.patch(`/risk-assessment/sop-checklist-items/${id}/status`, { is_active: isActive })
  },

  // Form operations
  createForm(data) {
    return api.post('/risk-assessment/forms', data)
  },

  getForms(params = {}) {
    return api.get('/risk-assessment/forms', { params })
  },

  getForm(id) {
    return api.get(`/risk-assessment/forms/${id}`)
  },

  updateDraft(id, data) {
    return api.put(`/risk-assessment/forms/${id}`, data)
  },

  submitForm(id) {
    return api.post(`/risk-assessment/forms/${id}/submit`)
  },

  uploadEvidence(id, formData) {
    return api.post(`/risk-assessment/forms/${id}/evidence`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  // Stats
  getStats(sppgIds) {
    return api.get('/risk-assessment/stats', {
      params: { sppg_ids: sppgIds.join(',') }
    })
  },

  // SPPG list for risk assessment
  getSPPGList() {
    return api.get('/risk-assessment/sppg-list')
  }
}

export default riskAssessmentService
