import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import ActivityLogItem from '@/components/mobile/ActivityLogItem.vue'

describe('ActivityLogItem', () => {
  const defaultProps = {
    employeeName: 'Budi Santoso',
    activityType: 'attendance',
    timestamp: '08:00',
    status: 'Hadir'
  }

  const mountItem = (props = {}) => {
    return mount(ActivityLogItem, {
      props: { ...defaultProps, ...props },
      global: {
        stubs: {
          'van-icon': true
        }
      }
    })
  }

  it('renders the activity log item card', () => {
    const wrapper = mountItem()
    expect(wrapper.find('.activity-log-item').exists()).toBe(true)
  })

  it('displays the employee name', () => {
    const wrapper = mountItem({ employeeName: 'Andi Wijaya' })
    expect(wrapper.find('.activity-log-item__name').text()).toBe('Andi Wijaya')
  })

  it('displays attendance activity type as Absensi', () => {
    const wrapper = mountItem({ activityType: 'attendance' })
    expect(wrapper.find('.activity-log-item__type').text()).toBe('Absensi')
  })

  it('displays delivery activity type as Pengiriman', () => {
    const wrapper = mountItem({ activityType: 'delivery' })
    expect(wrapper.find('.activity-log-item__type').text()).toBe('Pengiriman')
  })

  it('displays pickup activity type as Pengambilan', () => {
    const wrapper = mountItem({ activityType: 'pickup' })
    expect(wrapper.find('.activity-log-item__type').text()).toBe('Pengambilan')
  })

  it('displays the timestamp badge', () => {
    const wrapper = mountItem({ timestamp: '09:30' })
    expect(wrapper.find('.activity-log-item__timestamp').text()).toBe('09:30')
  })

  it('displays the status text', () => {
    const wrapper = mountItem({ status: 'Selesai' })
    expect(wrapper.find('.activity-log-item__status').text()).toBe('Selesai')
  })

  it('applies success status class for Hadir', () => {
    const wrapper = mountItem({ status: 'Hadir' })
    expect(wrapper.find('.activity-log-item__status').classes()).toContain('activity-log-item__status--success')
  })

  it('applies warning status class for Terlambat', () => {
    const wrapper = mountItem({ status: 'Terlambat' })
    expect(wrapper.find('.activity-log-item__status').classes()).toContain('activity-log-item__status--warning')
  })

  it('applies error status class for Tidak Hadir', () => {
    const wrapper = mountItem({ status: 'Tidak Hadir' })
    expect(wrapper.find('.activity-log-item__status').classes()).toContain('activity-log-item__status--error')
  })

  it('applies primary status class for Dalam Perjalanan', () => {
    const wrapper = mountItem({ status: 'Dalam Perjalanan' })
    expect(wrapper.find('.activity-log-item__status').classes()).toContain('activity-log-item__status--primary')
  })

  it('renders correct icon class for attendance type', () => {
    const wrapper = mountItem({ activityType: 'attendance' })
    const icon = wrapper.find('.activity-log-item__type-icon')
    expect(icon.classes()).toContain('activity-log-item__type-icon--attendance')
  })

  it('renders correct icon class for delivery type', () => {
    const wrapper = mountItem({ activityType: 'delivery' })
    const icon = wrapper.find('.activity-log-item__type-icon')
    expect(icon.classes()).toContain('activity-log-item__type-icon--delivery')
  })

  it('renders correct icon class for pickup type', () => {
    const wrapper = mountItem({ activityType: 'pickup' })
    const icon = wrapper.find('.activity-log-item__type-icon')
    expect(icon.classes()).toContain('activity-log-item__type-icon--pickup')
  })
})
