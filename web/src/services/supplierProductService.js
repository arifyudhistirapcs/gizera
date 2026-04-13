import api from './api'

const supplierProductService = {
  // Get all supplier products with optional filters
  getProducts(params = {}) {
    return api.get('/supplier-products', { params })
  },

  // Create new supplier product
  createProduct(data) {
    return api.post('/supplier-products', data)
  },

  // Update existing supplier product
  updateProduct(id, data) {
    return api.put(`/supplier-products/${id}`, data)
  },

  // Delete supplier product
  deleteProduct(id) {
    return api.delete(`/supplier-products/${id}`)
  }
}

export default supplierProductService
