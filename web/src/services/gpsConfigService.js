import api from './api'

const gpsConfigService = {
  // Get all GPS configurations
  async getGPSConfigs(activeOnly = false) {
    const params = activeOnly ? { active_only: 'true' } : {}
    const response = await api.get('/gps-config', { params })
    return response.data
  },

  // Create a new GPS configuration
  async createGPSConfig(data) {
    const response = await api.post('/gps-config', data)
    return response.data
  },

  // Update a GPS configuration
  async updateGPSConfig(id, data) {
    const response = await api.put(`/gps-config/${id}`, data)
    return response.data
  },

  // Delete a GPS configuration
  async deleteGPSConfig(id) {
    const response = await api.delete(`/gps-config/${id}`)
    return response.data
  }
}

export default gpsConfigService
