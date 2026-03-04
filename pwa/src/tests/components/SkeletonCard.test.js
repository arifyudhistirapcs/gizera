import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

describe('SkeletonCard', () => {
  it('renders with default props', () => {
    const wrapper = mount(SkeletonCard, {
      global: { stubs: { 'van-skeleton': true } }
    })
    expect(wrapper.find('.skeleton-card').exists()).toBe(true)
    const skeleton = wrapper.find('van-skeleton-stub')
    expect(skeleton.exists()).toBe(true)
    expect(skeleton.attributes('row')).toBe('3')
    expect(skeleton.attributes('avatar')).toBe('false')
    expect(skeleton.attributes('title')).toBe('false')
    expect(skeleton.attributes('round')).toBeDefined()
    expect(skeleton.attributes('loading')).toBeDefined()
  })

  it('passes custom rows prop to van-skeleton', () => {
    const wrapper = mount(SkeletonCard, {
      props: { rows: 5 },
      global: { stubs: { 'van-skeleton': true } }
    })
    const skeleton = wrapper.find('van-skeleton-stub')
    expect(skeleton.attributes('row')).toBe('5')
  })

  it('passes avatar prop to van-skeleton', () => {
    const wrapper = mount(SkeletonCard, {
      props: { avatar: true },
      global: { stubs: { 'van-skeleton': true } }
    })
    const skeleton = wrapper.find('van-skeleton-stub')
    expect(skeleton.attributes('avatar')).toBe('true')
  })

  it('passes title prop to van-skeleton', () => {
    const wrapper = mount(SkeletonCard, {
      props: { title: true },
      global: { stubs: { 'van-skeleton': true } }
    })
    const skeleton = wrapper.find('van-skeleton-stub')
    expect(skeleton.attributes('title')).toBe('true')
  })

  it('applies card styling class', () => {
    const wrapper = mount(SkeletonCard, {
      global: { stubs: { 'van-skeleton': true } }
    })
    const card = wrapper.find('.skeleton-card')
    expect(card.exists()).toBe(true)
  })
})
