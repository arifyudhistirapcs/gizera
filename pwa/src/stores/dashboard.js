import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useDashboardStore = defineStore('dashboard', () => {
  const summary = ref({
    totalHadir: 0,
    totalPengiriman: 0,
    totalSelesai: 0,
    totalSekolah: 0
  })
  const attendanceChart = ref([]) // Array of { date, count } for 7 days
  const recentTasks = ref([])     // Array of { id, schoolName, status, taskType }
  const loading = ref(false)
  const error = ref(null)

  async function fetchDashboardData() {
    loading.value = true
    error.value = null

    try {
      const [summaryRes, chartRes, tasksRes] = await Promise.all([
        api.get('/dashboard/summary'),
        api.get('/dashboard/attendance-chart'),
        api.get('/dashboard/recent-tasks')
      ])

      if (summaryRes.data.success) {
        const s = summaryRes.data.data
        summary.value = {
          totalHadir: s.totalHadir ?? 0,
          totalPengiriman: s.totalPengiriman ?? 0,
          totalSelesai: s.totalSelesai ?? 0,
          totalSekolah: s.totalSekolah ?? 0
        }
      }

      if (chartRes.data.success) {
        attendanceChart.value = chartRes.data.data ?? []
      }

      if (tasksRes.data.success) {
        recentTasks.value = tasksRes.data.data ?? []
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat data dashboard. Silakan coba lagi.'
    } finally {
      loading.value = false
    }
  }

  function retry() {
    return fetchDashboardData()
  }

  return {
    summary,
    attendanceChart,
    recentTasks,
    loading,
    error,
    fetchDashboardData,
    retry
  }
})
