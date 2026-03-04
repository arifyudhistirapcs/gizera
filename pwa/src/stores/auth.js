import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const schoolId = computed(() => user.value?.schoolId ?? null)

  function setAuth(userData, authToken) {
    user.value = userData
    token.value = authToken
    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function clearAuth() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function login(credentials) {
    try {
      const response = await authAPI.login(credentials)
      if (response.data.success) {
        setAuth(response.data.data.user, response.data.data.token)
        return { success: true, data: response.data.data }
      }
      return { success: false, message: response.data.message }
    } catch (error) {
      return { 
        success: false, 
        message: error.response?.data?.message || 'Login gagal' 
      }
    }
  }

  async function logout() {
    try {
      await authAPI.logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuth()
    }
  }

  // Initialize user from localStorage on store creation
  const storedUser = localStorage.getItem('user')
  if (storedUser) {
    try {
      user.value = JSON.parse(storedUser)
    } catch (error) {
      console.error('Error parsing stored user:', error)
      localStorage.removeItem('user')
    }
  }

  return {
    user,
    token,
    isAuthenticated,
    schoolId,
    setAuth,
    clearAuth,
    login,
    logout
  }
})
