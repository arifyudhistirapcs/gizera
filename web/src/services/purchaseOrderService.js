import api from './api'

const purchaseOrderService = {
  // Get all purchase orders with optional filters
  getPurchaseOrders(params = {}) {
    return api.get('/purchase-orders', { params })
  },

  // Get single purchase order by ID
  getPurchaseOrder(id) {
    return api.get(`/purchase-orders/${id}`)
  },

  // Create new purchase order
  createPurchaseOrder(data) {
    return api.post('/purchase-orders', data)
  },

  // Update existing purchase order
  updatePurchaseOrder(id, data) {
    return api.put(`/purchase-orders/${id}`, data)
  },

  // Delete purchase order
  deletePurchaseOrder(id) {
    return api.delete(`/purchase-orders/${id}`)
  },

  // Approve purchase order
  approvePurchaseOrder(id) {
    return api.post(`/purchase-orders/${id}/approve`)
  },

  // Create batch POs from RAB (kepala_yayasan)
  createBatchFromRAB(data) {
    return api.post('/purchase-orders/batch-from-rab', data)
  },

  // Supplier confirms PO as-is
  confirmBySupplier(id) {
    return api.post(`/purchase-orders/${id}/confirm`)
  },

  // Supplier marks PO as shipping
  markAsShipping(id) {
    return api.post(`/purchase-orders/${id}/shipping`)
  },

  // Supplier requests revision (changed items/notes)
  requestRevision(id, data) {
    return api.post(`/purchase-orders/${id}/request-revision`, data)
  },

  // Kepala yayasan accepts supplier's revision
  acceptRevision(id) {
    return api.post(`/purchase-orders/${id}/accept-revision`)
  },

  // Kepala yayasan revises PO and sends back to supplier
  revisePO(id, data) {
    return api.post(`/purchase-orders/${id}/revise`, data)
  }
}

export default purchaseOrderService
