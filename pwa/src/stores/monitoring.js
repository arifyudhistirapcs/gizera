import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useMonitoringStore = defineStore('monitoring', () => {
  const activities = ref([])
  const selectedDate = ref(new Date())
  const filterType = ref('all') // 'all' | 'attendance' | 'delivery' | 'pickup'
  const loading = ref(false)
  const hasMore = ref(true)
  const page = ref(1)
  const error = ref(null)

  const PAGE_LIMIT = 20

  function _formatDate(date) {
    const d = date instanceof Date ? date : new Date(date)
    return d.toISOString().split('T')[0]
  }

  async function fetchActivities() {
    loading.value = true
    error.value = null
    page.value = 1
    hasMore.value = true

    try {
      const params = {
        date: _formatDate(selectedDate.value),
        page: 1,
        limit: PAGE_LIMIT
      }
      if (filterType.value !== 'all') {
        params.filter = filterType.value
      }

      const res = await api.get('/monitoring/activities', { params })

      if (res.data.success) {
        activities.value = res.data.data ?? []
        hasMore.value = activities.value.length >= PAGE_LIMIT
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat data monitoring. Silakan coba lagi.'
    } finally {
      loading.value = false
    }
  }

  async function loadMore() {
    if (!hasMore.value || loading.value) return

    loading.value = true
    error.value = null

    try {
      const nextPage = page.value + 1
      const params = {
        date: _formatDate(selectedDate.value),
        page: nextPage,
        limit: PAGE_LIMIT
      }
      if (filterType.value !== 'all') {
        params.filter = filterType.value
      }

      const res = await api.get('/monitoring/activities', { params })

      if (res.data.success) {
        const newItems = res.data.data ?? []
        activities.value = [...activities.value, ...newItems]
        page.value = nextPage
        hasMore.value = newItems.length >= PAGE_LIMIT
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat data monitoring. Silakan coba lagi.'
    } finally {
      loading.value = false
    }
  }

  function setFilter(type) {
    filterType.value = type
    return fetchActivities()
  }

  function setDate(date) {
    selectedDate.value = date
    return fetchActivities()
  }

  function retry() {
    return fetchActivities()
  }

  return {
    activities,
    selectedDate,
    filterType,
    loading,
    hasMore,
    page,
    error,
    fetchActivities,
    loadMore,
    setFilter,
    setDate,
    retry
  }
})
