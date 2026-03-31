import api from './api'

const riskAssessmentService = {
  // ==================== SOP Categories ====================

  /**
   * Get all SOP categories ordered by urutan
   * @returns {Promise} API response with categories array
   */
  getSOPCategories() {
    return api.get('/risk-assessment/sop-categories')
  },

  // ==================== SOP Checklist Items ====================

  /**
   * Get SOP checklist items, optionally filtered by category
   * @param {Object} params - { category_id }
   * @returns {Promise} API response with checklist items array
   */
  getSOPChecklistItems(params = {}) {
    return api.get('/risk-assessment/sop-checklist-items', { params })
  },

  // ==================== Form Operations ====================

  /**
   * Create a new risk assessment form for a SPPG
   * @param {Object} data - { sppg_id }
   * @returns {Promise} API response with created form
   */
  createForm(data) {
    return api.post('/risk-assessment/forms', data)
  },

  /**
   * Get paginated list of risk assessment forms
   * @param {Object} params - { sppg_id, status, risk_level, date_from, date_to, page, page_size }
   * @returns {Promise} API response with forms array and pagination
   */
  getForms(params = {}) {
    return api.get('/risk-assessment/forms', { params })
  },

  /**
   * Get a single risk assessment form with items, category scores, and snapshot
   * @param {number} id - Form ID
   * @returns {Promise} API response with full form detail
   */
  getForm(id) {
    return api.get(`/risk-assessment/forms/${id}`)
  },

  /**
   * Update a draft form's items (scores, notes)
   * @param {number} id - Form ID
   * @param {Object} data - { items: [{ id, compliance_score, catatan }] }
   * @returns {Promise} API response
   */
  updateDraft(id, data) {
    return api.put(`/risk-assessment/forms/${id}`, data)
  },

  /**
   * Submit a completed form for score calculation
   * @param {number} id - Form ID
   * @returns {Promise} API response with calculated scores and risk level
   */
  submitForm(id) {
    return api.post(`/risk-assessment/forms/${id}/submit`)
  },

  deleteForm(id) {
    return api.delete(`/risk-assessment/forms/${id}`)
  },

  /**
   * Upload evidence photo for a form item
   * @param {number} id - Form ID
   * @param {FormData} formData - Multipart form data with file and item_id
   * @returns {Promise} API response with evidence URL
   */
  uploadEvidence(id, formData) {
    return api.post(`/risk-assessment/forms/${id}/evidence`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      timeout: 30000
    })
  },

  // ==================== SPPG List ====================

  /**
   * Get list of SPPGs under the user's Yayasan
   * @returns {Promise} API response with SPPG array
   */
  getSPPGList() {
    return api.get('/risk-assessment/sppg-list')
  },

  // ==================== Stats ====================

  /**
   * Get risk assessment statistics per SPPG
   * @param {number[]} sppgIds - Array of SPPG IDs
   * @returns {Promise} API response with stats
   */
  getStats(sppgIds) {
    return api.get('/risk-assessment/stats', {
      params: { sppg_ids: sppgIds.join(',') }
    })
  }
}

export default riskAssessmentService
