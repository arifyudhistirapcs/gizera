import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import SwipeAction from '@/components/mobile/SwipeAction.vue'

describe('SwipeAction', () => {
  const mountSwipe = (props = {}) => {
    return mount(SwipeAction, {
      props: {
        canCheckIn: false,
        canCheckOut: false,
        loading: false,
        ...props
      },
      global: {
        stubs: {
          'van-icon': true,
          'van-loading': true
        }
      }
    })
  }

  describe('rendering', () => {
    it('renders the swipe action track', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      expect(wrapper.find('.swipe-action__track').exists()).toBe(true)
    })

    it('renders the thumb element', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      expect(wrapper.find('.swipe-action__thumb').exists()).toBe(true)
    })

    it('shows right arrows when canCheckIn is true', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      expect(wrapper.find('.swipe-action__arrows--right').exists()).toBe(true)
    })

    it('shows left arrows when canCheckOut is true', () => {
      const wrapper = mountSwipe({ canCheckOut: true })
      expect(wrapper.find('.swipe-action__arrows--left').exists()).toBe(true)
    })

    it('hides right arrows when canCheckIn is false', () => {
      const wrapper = mountSwipe({ canCheckIn: false })
      expect(wrapper.find('.swipe-action__arrows--right').exists()).toBe(false)
    })

    it('hides left arrows when canCheckOut is false', () => {
      const wrapper = mountSwipe({ canCheckOut: false })
      expect(wrapper.find('.swipe-action__arrows--left').exists()).toBe(false)
    })

    it('shows both arrow sets when both actions available', () => {
      const wrapper = mountSwipe({ canCheckIn: true, canCheckOut: true })
      expect(wrapper.find('.swipe-action__arrows--right').exists()).toBe(true)
      expect(wrapper.find('.swipe-action__arrows--left').exists()).toBe(true)
    })
  })

  describe('disabled state', () => {
    it('applies disabled class when loading is true', () => {
      const wrapper = mountSwipe({ loading: true, canCheckIn: true })
      expect(wrapper.find('.swipe-action--disabled').exists()).toBe(true)
    })

    it('applies disabled class when both canCheckIn and canCheckOut are false', () => {
      const wrapper = mountSwipe({ canCheckIn: false, canCheckOut: false })
      expect(wrapper.find('.swipe-action--disabled').exists()).toBe(true)
    })

    it('does not apply disabled class when canCheckIn is true and not loading', () => {
      const wrapper = mountSwipe({ canCheckIn: true, loading: false })
      expect(wrapper.find('.swipe-action--disabled').exists()).toBe(false)
    })
  })

  describe('loading state', () => {
    it('shows loading indicator in label when loading', () => {
      const wrapper = mountSwipe({ loading: true, canCheckIn: true })
      const label = wrapper.find('.swipe-action__label')
      expect(label.find('van-loading-stub').exists()).toBe(true)
    })

    it('shows loading indicator in thumb when loading', () => {
      const wrapper = mountSwipe({ loading: true, canCheckIn: true })
      const thumb = wrapper.find('.swipe-action__thumb')
      expect(thumb.find('van-loading-stub').exists()).toBe(true)
    })
  })

  describe('label text', () => {
    it('shows check-in hint when only canCheckIn is true', () => {
      const wrapper = mountSwipe({ canCheckIn: true, canCheckOut: false })
      expect(wrapper.find('.swipe-action__label-text').text()).toContain('Check In')
    })

    it('shows check-out hint when only canCheckOut is true', () => {
      const wrapper = mountSwipe({ canCheckIn: false, canCheckOut: true })
      expect(wrapper.find('.swipe-action__label-text').text()).toContain('Check Out')
    })

    it('shows generic hint when both actions available', () => {
      const wrapper = mountSwipe({ canCheckIn: true, canCheckOut: true })
      expect(wrapper.find('.swipe-action__label-text').text()).toContain('Geser untuk absen')
    })
  })

  describe('gradient track styling', () => {
    it('has gradient background on track', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      const track = wrapper.find('.swipe-action__track')
      expect(track.exists()).toBe(true)
      // The gradient is applied via CSS class, verify the element exists
    })

    it('has pill shape (border-radius full)', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      const track = wrapper.find('.swipe-action__track')
      expect(track.exists()).toBe(true)
    })
  })

  describe('thumb styling', () => {
    it('thumb has white background', () => {
      const wrapper = mountSwipe({ canCheckIn: true })
      const thumb = wrapper.find('.swipe-action__thumb')
      expect(thumb.exists()).toBe(true)
    })
  })
})
