import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'
import { useAuthStore } from '@/stores/auth'

export const useSchoolMonitoringStore = defineStore('schoolMonitoring', () => {
  const todayMenu = ref(null)        // { menuName, components, portions }
  const deliveryStatus = ref(null)   // { status, driverName, estimatedTime }
  const pickupStatus = ref(null)     // { status, driverName, estimatedTime } or null
  const loading = ref(false)
  const error = ref(null)

  const schoolId = computed(() => {
    const authStore = useAuthStore()
    return authStore.user?.schoolId ?? null
  })

  async function fetchSchoolData() {
    if (!schoolId.value) {
      error.value = 'School ID tidak ditemukan. Silakan login ulang.'
      return
    }

    loading.value = true
    error.value = null

    try {
      const res = await api.get(`/school-monitoring/${schoolId.value}/today`)

      if (res.data.success) {
        const data = res.data.data
        todayMenu.value = data.menu ?? null
        deliveryStatus.value = data.delivery ?? null
        pickupStatus.value = data.pickup ?? null
      }
    } catch (err) {
      error.value = err.response?.data?.message || 'Gagal memuat data monitoring sekolah. Silakan coba lagi.'
    } finally {
      loading.value = false
    }
  }

  function retry() {
    return fetchSchoolData()
  }

  return {
    todayMenu,
    deliveryStatus,
    pickupStatus,
    loading,
    error,
    schoolId,
    fetchSchoolData,
    retry
  }
})
