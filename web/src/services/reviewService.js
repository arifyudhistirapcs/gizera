import api from './api'

const reviewService = {
  // Get all reviews with filters
  getReviews: async (params = {}) => {
    const response = await api.get('/reviews', { params })
    return response.data
  },

  // Get review by ID
  getReview: async (id) => {
    const response = await api.get(`/reviews/${id}`)
    return response.data
  },

  // Get review by delivery record ID
  getReviewByDeliveryRecord: async (deliveryRecordId) => {
    const response = await api.get('/reviews/by-delivery', {
      params: { delivery_record_id: deliveryRecordId }
    })
    return response.data
  },

  // Check if review exists for delivery record
  checkReviewExists: async (deliveryRecordId) => {
    const response = await api.get('/reviews/check', {
      params: { delivery_record_id: deliveryRecordId }
    })
    return response.data
  },

  // Get review summary/statistics
  getSummary: async (params = {}) => {
    const response = await api.get('/reviews/summary', { params })
    return response.data
  }
}

export default reviewService
