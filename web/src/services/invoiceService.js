import api from './api'

const invoiceService = {
  // Get all invoices with optional filters
  getInvoices(params = {}) {
    return api.get('/invoices', { params })
  },

  // Create new invoice
  createInvoice(data) {
    return api.post('/invoices', data)
  },

  // Get single invoice by ID
  getInvoiceDetail(id) {
    return api.get(`/invoices/${id}`)
  },

  // Pay invoice
  payInvoice(id, data) {
    return api.post(`/invoices/${id}/pay`, data)
  },

  // Upload payment proof
  uploadPaymentProof(id, formData) {
    return api.post(`/invoices/${id}/upload-proof`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}

export default invoiceService
