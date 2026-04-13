import api from './api'

const supplierPortalService = {
  // Get supplier dashboard summary
  getDashboard() {
    return api.get('/supplier/dashboard')
  },

  // Get supplier payment history
  getPayments(params = {}) {
    return api.get('/supplier/payments', { params })
  },

  // Get supplier products
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
  },

  // Get invoices
  getInvoices(params = {}) {
    return api.get('/invoices', { params })
  },

  // Create new invoice
  createInvoice(data) {
    return api.post('/invoices', data)
  },

  // Get purchase orders
  getPurchaseOrders(params = {}) {
    return api.get('/purchase-orders', { params })
  }
}

export default supplierPortalService
