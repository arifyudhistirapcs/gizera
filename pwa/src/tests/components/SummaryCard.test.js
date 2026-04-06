import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SummaryCard from '@/components/mobile/SummaryCard.vue'

describe('SummaryCard', () => {
  const defaultProps = {
    icon: 'friends-o',
    iconColor: '#F0F0F0',
    label: 'Total Hadir',
    value: '24',
    loading: false
  }

  const mountCard = (props = {}) => {
    return mount(SummaryCard, {
      props: { ...defaultProps, ...props },
      global: {
        stubs: {
          'van-icon': true,
          'van-skeleton': true
        }
      }
    })
  }

  it('renders stat card when loading is false', () => {
    const wrapper = mountCard()
    expect(wrapper.find('.summary-card').exists()).toBe(true)
    expect(wrapper.find('.skeleton-card').exists()).toBe(false)
  })

  it('renders SkeletonCard when loading is true', () => {
    const wrapper = mountCard({ loading: true })
    expect(wrapper.find('.summary-card').exists()).toBe(false)
    expect(wrapper.find('.skeleton-card').exists()).toBe(true)
  })

  it('displays the label text', () => {
    const wrapper = mountCard({ label: 'Tugas Selesai' })
    expect(wrapper.find('.summary-card__label').text()).toBe('Tugas Selesai')
  })

  it('displays string value', () => {
    const wrapper = mountCard({ value: '95%' })
    expect(wrapper.find('.summary-card__value').text()).toBe('95%')
  })

  it('displays numeric value', () => {
    const wrapper = mountCard({ value: 42 })
    expect(wrapper.find('.summary-card__value').text()).toBe('42')
  })

  it('renders icon with correct props', () => {
    const wrapper = mountCard({ icon: 'chart-trending-o', iconColor: '#05CD99' })
    const icon = wrapper.find('van-icon-stub')
    expect(icon.exists()).toBe(true)
    expect(icon.attributes('name')).toBe('chart-trending-o')
    expect(icon.attributes('color')).toBe('#05CD99')
  })

  it('renders 40x40 icon circle', () => {
    const wrapper = mountCard()
    const iconCircle = wrapper.find('.summary-card__icon')
    expect(iconCircle.exists()).toBe(true)
  })

  it('uses default iconColor when not provided', () => {
    const wrapper = mount(SummaryCard, {
      props: {
        icon: 'friends-o',
        label: 'Test',
        value: '1',
        loading: false
      },
      global: {
        stubs: {
          'van-icon': true,
          'van-skeleton': true
        }
      }
    })
    const icon = wrapper.find('van-icon-stub')
    expect(icon.attributes('color')).toBe('var(--h-primary)')
  })
})
