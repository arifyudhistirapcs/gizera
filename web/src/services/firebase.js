import { initializeApp } from 'firebase/app'
import { getDatabase } from 'firebase/database'
import { getStorage } from 'firebase/storage'

const firebaseConfig = {
  apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
  authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
  databaseURL: import.meta.env.VITE_FIREBASE_DATABASE_URL,
  projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
  storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
  appId: import.meta.env.VITE_FIREBASE_APP_ID
}

const app = initializeApp(firebaseConfig)
const database = getDatabase(app)
const storage = getStorage(app)

/**
 * Get tenant context from auth store.
 * Returns { sppg_id, yayasan_id, role } from the current user.
 */
function getTenantContext() {
  try {
    const storedUser = localStorage.getItem('user')
    if (storedUser) {
      const user = JSON.parse(storedUser)
      return {
        sppg_id: user.sppg_id || null,
        yayasan_id: user.yayasan_id || null,
        role: user.role || null
      }
    }
  } catch (e) {
    console.warn('[Firebase] Failed to get tenant context:', e)
  }
  return { sppg_id: null, yayasan_id: null, role: null }
}

/**
 * Build tenant-aware Firebase paths based on user role.
 *
 * SPPG-level roles: paths include /{sppg_id}/
 * kepala_yayasan: dashboard path uses yayasan_id
 * admin_bgn/superadmin: dashboard path uses /bgn/
 */
const firebasePaths = {
  /**
   * KDS Cooking path: /kds/cooking/{sppg_id}/{date}
   */
  kdsCooking(date) {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      console.warn('[Firebase] No sppg_id for KDS cooking path, falling back to legacy path')
      return `/kds/cooking/${date}`
    }
    return `/kds/cooking/${sppg_id}/${date}`
  },

  /**
   * KDS Packing path: /kds/packing/{sppg_id}/{date}
   */
  kdsPacking(date) {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      console.warn('[Firebase] No sppg_id for KDS packing path, falling back to legacy path')
      return `/kds/packing/${date}`
    }
    return `/kds/packing/${sppg_id}/${date}`
  },

  /**
   * Cleaning path: /cleaning/{sppg_id}/pending
   */
  cleaningPending() {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      console.warn('[Firebase] No sppg_id for cleaning path, falling back to legacy path')
      return '/cleaning/pending'
    }
    return `/cleaning/${sppg_id}/pending`
  },

  /**
   * Monitoring path: /monitoring/{sppg_id}/{date}
   */
  monitoring(date) {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      return `/monitoring/${date}`
    }
    return `/monitoring/${sppg_id}/${date}`
  },

  /**
   * Dashboard Kepala SPPG: /dashboard/kepala_sppg/{sppg_id}
   */
  dashboardKepalaSPPG() {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      return '/dashboard/kepala_sppg'
    }
    return `/dashboard/kepala_sppg/${sppg_id}`
  },

  /**
   * Dashboard Kepala Yayasan: /dashboard/kepala_yayasan/{yayasan_id}
   */
  dashboardKepalaYayasan() {
    const { yayasan_id } = getTenantContext()
    if (!yayasan_id) {
      return '/dashboard/kepala_yayasan'
    }
    return `/dashboard/kepala_yayasan/${yayasan_id}`
  },

  /**
   * Dashboard BGN (admin_bgn): /dashboard/bgn
   */
  dashboardBGN() {
    return '/dashboard/bgn'
  },

  /**
   * Notifications path: /notifications/{sppg_id}/...
   */
  notifications(subPath) {
    const { sppg_id } = getTenantContext()
    if (!sppg_id) {
      return `/notifications/${subPath}`
    }
    return `/notifications/${sppg_id}/${subPath}`
  },

  /**
   * Get the appropriate dashboard path based on user role.
   */
  dashboardForCurrentRole() {
    const { role } = getTenantContext()
    switch (role) {
      case 'admin_bgn':
      case 'superadmin':
        return this.dashboardBGN()
      case 'kepala_yayasan':
        return this.dashboardKepalaYayasan()
      default:
        return this.dashboardKepalaSPPG()
    }
  }
}

export { app, database, storage, firebasePaths, getTenantContext }
