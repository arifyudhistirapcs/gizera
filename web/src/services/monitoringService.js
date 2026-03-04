import api from './api'

/**
 * Get delivery records for a specific date with optional filters
 * @param {string} date - Date in YYYY-MM-DD format
 * @param {object} filters - Optional filters (school_id, status, driver_id)
 * @returns {Promise} API response
 */
export const getDeliveryRecords = async (date, filters = {}) => {
  try {
    const params = { date, ...filters }
    const response = await api.get('/monitoring/deliveries', { params })
    return response.data
  } catch (error) {
    console.error('Error fetching delivery records:', error)
    throw error
  }
}

/**
 * Get delivery record detail by ID
 * @param {number} id - Delivery record ID
 * @returns {Promise} API response
 */
export const getDeliveryDetail = async (id) => {
  try {
    const response = await api.get(`/monitoring/deliveries/${id}`)
    return response.data
  } catch (error) {
    console.error('Error fetching delivery detail:', error)
    throw error
  }
}

/**
 * Update delivery status
 * @param {number} id - Delivery record ID
 * @param {string} status - New status
 * @param {string} notes - Optional notes
 * @returns {Promise} API response
 */
export const updateDeliveryStatus = async (id, status, notes = '') => {
  try {
    const response = await api.put(`/monitoring/deliveries/${id}/status`, {
      status,
      notes
    })
    return response.data
  } catch (error) {
    console.error('Error updating delivery status:', error)
    throw error
  }
}

/**
 * Get activity log for a delivery record
 * @param {number} id - Delivery record ID
 * @returns {Promise} API response
 */
export const getActivityLog = async (id) => {
  try {
    const response = await api.get(`/monitoring/deliveries/${id}/activity`)
    return response.data
  } catch (error) {
    console.error('Error fetching activity log:', error)
    throw error
  }
}

/**
 * Get daily summary statistics
 * @param {string} date - Date in YYYY-MM-DD format
 * @returns {Promise} API response
 */
export const getDailySummary = async (date) => {
  try {
    const response = await api.get('/monitoring/summary', {
      params: { date }
    })
    return response.data
  } catch (error) {
    console.error('Error fetching daily summary:', error)
    throw error
  }
}

/**
 * Get e-POD data by delivery task ID or delivery record ID
 * @param {number} id - Delivery task ID or delivery record ID
 * @param {string} type - 'task' or 'record' (default: 'record')
 * @returns {Promise} API response
 */
export const getEPODByDeliveryTask = async (id, type = 'record') => {
  try {
    const paramKey = type === 'task' ? 'delivery_task_id' : 'delivery_record_id'
    const response = await api.get('/epod', {
      params: { [paramKey]: id }
    })
    return response.data
  } catch (error) {
    console.error('Error fetching e-POD:', error)
    throw error
  }
}
