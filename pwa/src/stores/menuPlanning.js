import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useMenuPlanningStore = defineStore('menuPlanning', () => {
  const weeklyPlans = ref([])       // Array of { id, weekPeriod, approvalStatus, menuCount, menus: [...] }
  const selectedWeek = ref(null)    // Currently selected week for filtering (e.g. '2025-W02')
  const loading = ref(false)
  const approvalLoading = ref(false) // Separate loading for approve/reject actions
  const error = ref(null)

  async function fetchMenuPlans() {
    loading.value = true
    error.value = null

    try {
      const params = {}
      if (selectedWeek.value) {
        params.week = selectedWeek.value
      }

      const res = await api.get('/menu-planning/plans', { params })

      if (res.data.success) {
        weeklyPlans.value = res.data.data ?? []
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat data perencanaan menu. Silakan coba lagi.'
    } finally {
      loading.value = false
    }
  }

  async function approveMenu(id) {
    approvalLoading.value = true
    error.value = null

    try {
      await api.post(`/menu-planning/plans/${id}/approve`)
      // Re-fetch to get updated status from server (no optimistic update)
      await fetchMenuPlans()
    } catch (err) {
      // Preserve previous status — no state mutation on failure
      error.value = err.response?.data?.message || 'Gagal menyetujui menu. Silakan coba lagi.'
    } finally {
      approvalLoading.value = false
    }
  }

  async function rejectMenu(id, reason) {
    approvalLoading.value = true
    error.value = null

    try {
      await api.post(`/menu-planning/plans/${id}/reject`, { reason })
      // Re-fetch to get updated status from server (no optimistic update)
      await fetchMenuPlans()
    } catch (err) {
      // Preserve previous status — no state mutation on failure
      error.value = err.response?.data?.message || 'Gagal menolak menu. Silakan coba lagi.'
    } finally {
      approvalLoading.value = false
    }
  }

  function setWeek(week) {
    selectedWeek.value = week
    return fetchMenuPlans()
  }

  function retry() {
    return fetchMenuPlans()
  }

  return {
    weeklyPlans,
    selectedWeek,
    loading,
    approvalLoading,
    error,
    fetchMenuPlans,
    approveMenu,
    rejectMenu,
    setWeek,
    retry
  }
})
