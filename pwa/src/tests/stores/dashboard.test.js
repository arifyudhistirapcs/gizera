import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useDashboardStore } from '@/stores/dashboard'

// Mock the api module
vi.mock('@/services/api', () => ({
  default: {
    get: vi.fn()
  }
}))

import api from '@/services/api'

describe('Dashboard Store', () => {
  let store

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useDashboardStore()
    vi.clearAllMocks()
  })

  it('has correct initial state', () => {
    expect(store.summary).toEqual({
      totalHadir: 0,
      totalPengiriman: 0,
      totalSelesai: 0,
      totalSekolah: 0
    })
    expect(store.attendanceChart).toEqual([])
    expect(store.recentTasks).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('sets loading to true while fetching and false after', async () => {
    api.get.mockResolvedValue({ data: { success: true, data: {} } })

    const promise = store.fetchDashboardData()
    expect(store.loading).toBe(true)

    await promise
    expect(store.loading).toBe(false)
  })

  it('fetches dashboard data successfully', async () => {
    api.get
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: { totalHadir: 45, totalPengiriman: 12, totalSelesai: 8, totalSekolah: 15 }
        }
      })
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: [
            { date: '2025-01-01', count: 40 },
            { date: '2025-01-02', count: 42 }
          ]
        }
      })
      .mockResolvedValueOnce({
        data: {
          success: true,
          data: [
            { id: '1', schoolName: 'SD Negeri 1', status: 'completed', taskType: 'delivery' }
          ]
        }
      })

    await store.fetchDashboardData()

    expect(store.summary).toEqual({
      totalHadir: 45,
      totalPengiriman: 12,
      totalSelesai: 8,
      totalSekolah: 15
    })
    expect(store.attendanceChart).toHaveLength(2)
    expect(store.recentTasks).toHaveLength(1)
    expect(store.error).toBeNull()
  })

  it('sets error message on API failure', async () => {
    api.get.mockRejectedValue({
      response: { data: { message: 'Server error' } }
    })

    await store.fetchDashboardData()

    expect(store.error).toBe('Server error')
    expect(store.loading).toBe(false)
  })

  it('sets fallback error message when no response message', async () => {
    api.get.mockRejectedValue(new Error('Network Error'))

    await store.fetchDashboardData()

    expect(store.error).toBe('Gagal memuat data dashboard. Silakan coba lagi.')
    expect(store.loading).toBe(false)
  })

  it('clears error on successful retry', async () => {
    // First call fails
    api.get.mockRejectedValue(new Error('fail'))
    await store.fetchDashboardData()
    expect(store.error).not.toBeNull()

    // Retry succeeds
    api.get.mockResolvedValue({ data: { success: true, data: {} } })
    await store.retry()

    expect(store.error).toBeNull()
  })

  it('calls all three API endpoints', async () => {
    api.get.mockResolvedValue({ data: { success: true, data: {} } })

    await store.fetchDashboardData()

    expect(api.get).toHaveBeenCalledWith('/dashboard/summary')
    expect(api.get).toHaveBeenCalledWith('/dashboard/attendance-chart')
    expect(api.get).toHaveBeenCalledWith('/dashboard/recent-tasks')
  })

  it('handles missing data fields with defaults', async () => {
    api.get
      .mockResolvedValueOnce({
        data: { success: true, data: { totalHadir: 10 } }
      })
      .mockResolvedValueOnce({
        data: { success: true, data: null }
      })
      .mockResolvedValueOnce({
        data: { success: true, data: null }
      })

    await store.fetchDashboardData()

    expect(store.summary.totalHadir).toBe(10)
    expect(store.summary.totalPengiriman).toBe(0)
    expect(store.summary.totalSelesai).toBe(0)
    expect(store.summary.totalSekolah).toBe(0)
    expect(store.attendanceChart).toEqual([])
    expect(store.recentTasks).toEqual([])
  })
})
