import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authService from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const loading = ref(false)
  const error = ref(null)

  // Tenant & RBAC fields from /auth/me
  const sppgId = ref(null)
  const yayasanId = ref(null)
  const role = ref(null)
  const modules = ref([])
  const permissions = ref([])

  const isAuthenticated = computed(() => !!token.value)

  // Role getters
  const isSuperadmin = computed(() => role.value === 'superadmin')
  const isAdminBGN = computed(() => role.value === 'admin_bgn')
  const isKepalaYayasan = computed(() => role.value === 'kepala_yayasan')

  // Initialize user from localStorage if token exists
  const initAuth = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser && token.value) {
      try {
        const parsed = JSON.parse(storedUser)
        user.value = parsed
        _syncTenantFields(parsed)
      } catch (e) {
        clearAuth()
      }
    }
  }

  // Sync tenant/RBAC fields from user object
  function _syncTenantFields(userData) {
    if (!userData) {
      sppgId.value = null
      yayasanId.value = null
      role.value = null
      modules.value = []
      permissions.value = []
      return
    }
    sppgId.value = userData.sppg_id ?? null
    yayasanId.value = userData.yayasan_id ?? null
    role.value = userData.role ?? null
    modules.value = Array.isArray(userData.modules) ? userData.modules : []
    permissions.value = Array.isArray(userData.permissions) ? userData.permissions : []
  }

  const login = async (credentials) => {
    loading.value = true
    error.value = null
    try {
      const response = await authService.login(credentials)
      setAuth(response.user, response.token)
      // Fetch full user data (modules, permissions, sppg, yayasan) from /auth/me
      try {
        await getCurrentUser()
      } catch (e) {
        console.warn('Failed to fetch full user data after login:', e)
      }
      return response
    } catch (err) {
      error.value = err.response?.data?.message || 'Login gagal. Silakan coba lagi.'
      throw err
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    loading.value = true
    try {
      await authService.logout()
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      clearAuth()
      loading.value = false
    }
  }

  const refreshToken = async () => {
    try {
      const response = await authService.refreshToken()
      setAuth(response.user, response.token)
      return response
    } catch (err) {
      clearAuth()
      throw err
    }
  }

  const getCurrentUser = async () => {
    try {
      const response = await authService.getCurrentUser()
      user.value = response.user
      _syncTenantFields(response.user)
      localStorage.setItem('user', JSON.stringify(response.user))
      return response.user
    } catch (err) {
      clearAuth()
      throw err
    }
  }

  function setAuth(userData, authToken) {
    user.value = userData
    token.value = authToken
    _syncTenantFields(userData)
    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function clearAuth() {
    user.value = null
    token.value = null
    _syncTenantFields(null)
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // Initialize on store creation
  initAuth()

  return {
    user,
    token,
    loading,
    error,
    sppgId,
    yayasanId,
    role,
    modules,
    permissions,
    isAuthenticated,
    isSuperadmin,
    isAdminBGN,
    isKepalaYayasan,
    login,
    logout,
    refreshToken,
    getCurrentUser,
    setAuth,
    clearAuth
  }
})
