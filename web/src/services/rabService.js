import api from './api'

const rabService = {
  // Get all RABs with optional filters
  getRABList(params = {}) {
    return api.get('/rab', { params })
  },

  // Get single RAB by ID
  getRABDetail(id) {
    return api.get(`/rab/${id}`)
  },

  // Update existing RAB
  updateRAB(id, data) {
    return api.put(`/rab/${id}`, data)
  },

  // Approve RAB at SPPG level
  approveSPPG(id) {
    return api.post(`/rab/${id}/approve-sppg`)
  },

  // Approve RAB at Yayasan level
  approveYayasan(id) {
    return api.post(`/rab/${id}/approve-yayasan`)
  },

  // Reject RAB with revision notes
  rejectRAB(id, data) {
    return api.post(`/rab/${id}/reject`, data)
  },

  // Resubmit RAB after revision
  resubmitRAB(id) {
    return api.post(`/rab/${id}/resubmit`)
  },

  // Get RAB comparison (planned vs actual)
  getRABComparison(id) {
    return api.get(`/rab/${id}/comparison`)
  },

  // Get RAB PO tracking
  getRABPOTracking(id) {
    return api.get(`/rab/${id}/po-tracking`)
  }
}

export default rabService
