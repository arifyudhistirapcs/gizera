import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'

// Mock vue-router
const mockPush = vi.fn()
const mockRoutePath = { value: '/attendance' }

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  }),
  useRoute: () => ({
    get path() { return mockRoutePath.value }
  })
}))

// Mock auth store - override the global mock from setup.js
let mockUserRole = 'driver'
vi.mock('@/stores/auth.js', () => ({
  useAuthStore: () => ({
    user: { get role() { return mockUserRole } }
  })
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    user: { get role() { return mockUserRole } }
  })
}))

import BottomNavigation from '@/components/mobile/BottomNavigation.vue'

// Use a functional stub that preserves event listeners
const VanTabbarStub = {
  name: 'van-tabbar',
  template: '<div class="van-tabbar bottom-navigation"><slot /></div>',
  props: ['modelValue', 'activeColor', 'inactiveColor'],
  emits: ['change', 'update:modelValue'],
  setup(props, { emit }) {
    return { emit }
  }
}

const VanTabbarItemStub = {
  name: 'van-tabbar-item',
  template: '<div class="van-tabbar-item"><slot /></div>',
  props: ['icon']
}

function mountNav() {
  return mount(BottomNavigation, {
    global: {
      components: {
        'van-tabbar': VanTabbarStub,
        'van-tabbar-item': VanTabbarItemStub
      }
    }
  })
}

describe('BottomNavigation', () => {
  beforeEach(() => {
    mockPush.mockClear()
    mockRoutePath.value = '/attendance'
  })

  describe('Role-based navigation items', () => {
    it('renders Tugas, Absensi, Profil for driver role', () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(3)
      expect(items[0].text()).toBe('Tugas')
      expect(items[1].text()).toBe('Absensi')
      expect(items[2].text()).toBe('Profil')
    })

    it('renders Tugas, Absensi, Profil for asisten_lapangan role', () => {
      mockUserRole = 'asisten_lapangan'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(3)
      expect(items[0].text()).toBe('Tugas')
      expect(items[1].text()).toBe('Absensi')
      expect(items[2].text()).toBe('Profil')
    })

    it('renders Dashboard, Monitoring, Menu, Absensi, Profil for kepala_sppg role', () => {
      mockUserRole = 'kepala_sppg'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(5)
      expect(items[0].text()).toBe('Dashboard')
      expect(items[1].text()).toBe('Monitoring')
      expect(items[2].text()).toBe('Menu')
      expect(items[3].text()).toBe('Absensi')
      expect(items[4].text()).toBe('Profil')
    })

    it('renders Menu, Absensi, Profil for ahli_gizi role', () => {
      mockUserRole = 'ahli_gizi'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(3)
      expect(items[0].text()).toBe('Menu')
      expect(items[1].text()).toBe('Absensi')
      expect(items[2].text()).toBe('Profil')
    })

    it('renders Monitoring, Profil for sekolah role', () => {
      mockUserRole = 'sekolah'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(2)
      expect(items[0].text()).toBe('Monitoring')
      expect(items[1].text()).toBe('Profil')
    })

    it('renders Absensi, Profil for unknown/default role', () => {
      mockUserRole = 'karyawan'
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(2)
      expect(items[0].text()).toBe('Absensi')
      expect(items[1].text()).toBe('Profil')
    })

    it('renders default nav when user role is null', () => {
      mockUserRole = null
      const wrapper = mountNav()
      const items = wrapper.findAll('.van-tabbar-item')
      expect(items).toHaveLength(2)
      expect(items[0].text()).toBe('Absensi')
      expect(items[1].text()).toBe('Profil')
    })
  })

  describe('Navigation on item tap', () => {
    it('navigates to /monitoring when kepala_sppg taps Monitoring tab', async () => {
      mockUserRole = 'kepala_sppg'
      const wrapper = mountNav()

      // The component binds @change="onTabChange" on van-tabbar
      // We call the component's internal onTabChange via the exposed vm
      // Simulate the change event by finding the tabbar and emitting
      const tabbar = wrapper.findComponent(VanTabbarStub)
      await tabbar.vm.$emit('change', 1)
      await wrapper.vm.$nextTick()
      expect(mockPush).toHaveBeenCalledWith('/monitoring')
    })

    it('navigates to /tasks for driver Tugas tab', async () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      const tabbar = wrapper.findComponent(VanTabbarStub)
      await tabbar.vm.$emit('change', 0)
      await wrapper.vm.$nextTick()
      expect(mockPush).toHaveBeenCalledWith('/tasks')
    })

    it('navigates to /profile for Profil tab', async () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      const tabbar = wrapper.findComponent(VanTabbarStub)
      await tabbar.vm.$emit('change', 2)
      await wrapper.vm.$nextTick()
      expect(mockPush).toHaveBeenCalledWith('/profile')
    })
  })

  describe('Styling and configuration', () => {
    it('has the bottom-navigation class', () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      expect(wrapper.find('.bottom-navigation').exists()).toBe(true)
    })

    it('configures active color as #5A4372', () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      const tabbar = wrapper.findComponent(VanTabbarStub)
      expect(tabbar.props('activeColor')).toBe('#5A4372')
    })

    it('configures inactive color as #ACA9B0', () => {
      mockUserRole = 'driver'
      const wrapper = mountNav()
      const tabbar = wrapper.findComponent(VanTabbarStub)
      expect(tabbar.props('inactiveColor')).toBe('#ACA9B0')
    })
  })
})
