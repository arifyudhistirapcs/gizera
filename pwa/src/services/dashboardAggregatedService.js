import api from './api'

const dashboardAggregatedService = {
  // ==================== Kepala Yayasan Dashboard ====================

  /**
   * Get aggregated dashboard for Kepala Yayasan
   * @param {Object} params - { start_date, end_date, sppg_id }
   */
  getKepalaYayasanDashboard(params = {}) {
    return api.get('/dashboard/kepala-yayasan', { params })
  },

  /**
   * Export Kepala Yayasan dashboard data
   * @param {Object} params - { start_date, end_date, sppg_id }
   */
  exportKepalaYayasanDashboard(params = {}) {
    return api.get('/dashboard/kepala-yayasan/export', { params })
  },

  // ==================== Admin BGN Dashboard ====================

  /**
   * Get aggregated dashboard for Admin BGN
   * @param {Object} params - { start_date, end_date, yayasan_id, sppg_id }
   */
  getAdminBGNDashboard(params = {}) {
    return api.get('/dashboard/admin-bgn', { params })
  },

  /**
   * Export Admin BGN dashboard data
   * @param {Object} params - { start_date, end_date, yayasan_id, sppg_id }
   */
  exportAdminBGNDashboard(params = {}) {
    return api.get('/dashboard/admin-bgn/export', { params })
  }
}

export default dashboardAggregatedService
