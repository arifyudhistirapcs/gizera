import api from './api'

const aggregatedDashboardService = {
  // ==================== Kepala Yayasan Dashboard ====================

  getKepalaYayasanDashboard(params = {}) {
    return api.get('/dashboard/kepala-yayasan', { params })
  },

  exportKepalaYayasanDashboard(params = {}) {
    return api.get('/dashboard/kepala-yayasan/export', { params })
  },

  // ==================== Admin BGN Dashboard ====================

  getAdminBGNDashboard(params = {}) {
    return api.get('/dashboard/admin-bgn', { params })
  },

  exportAdminBGNDashboard(params = {}) {
    return api.get('/dashboard/admin-bgn/export', { params })
  }
}

export default aggregatedDashboardService
