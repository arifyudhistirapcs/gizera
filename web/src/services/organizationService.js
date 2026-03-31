import api from './api'

const organizationService = {
  // ==================== Yayasan ====================

  getYayasanList(params = {}) {
    return api.get('/organizations/yayasan', { params })
  },

  getYayasan(id) {
    return api.get(`/organizations/yayasan/${id}`)
  },

  createYayasan(data) {
    return api.post('/organizations/yayasan', data)
  },

  updateYayasan(id, data) {
    return api.put(`/organizations/yayasan/${id}`, data)
  },

  setYayasanStatus(id, isActive) {
    return api.patch(`/organizations/yayasan/${id}/status`, { is_active: isActive })
  },

  // ==================== SPPG ====================

  getSPPGList(params = {}) {
    return api.get('/organizations/sppg', { params })
  },

  getSPPG(id) {
    return api.get(`/organizations/sppg/${id}`)
  },

  createSPPG(data) {
    return api.post('/organizations/sppg', data)
  },

  updateSPPG(id, data) {
    return api.put(`/organizations/sppg/${id}`, data)
  },

  setSPPGStatus(id, isActive) {
    return api.patch(`/organizations/sppg/${id}/status`, { is_active: isActive })
  },

  transferSPPG(id, yayasanId) {
    return api.put(`/organizations/sppg/${id}/transfer`, { yayasan_id: yayasanId })
  },

  // ==================== User Provisioning ====================

  getUsers(params = {}) {
    return api.get('/users', { params })
  },

  getUser(id) {
    return api.get(`/users/${id}`)
  },

  createUser(data) {
    return api.post('/users', data)
  },

  updateUser(id, data) {
    return api.put(`/users/${id}`, data)
  },

  setUserStatus(id, isActive) {
    return api.patch(`/users/${id}/status`, { is_active: isActive })
  }
}

export default organizationService
