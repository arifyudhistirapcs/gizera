import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TaskCard from '@/components/mobile/TaskCard.vue'

describe('TaskCard', () => {
  const defaultProps = {
    schoolName: 'SD Negeri 1',
    address: 'Jl. Merdeka 10',
    taskType: 'delivery',
    status: 'pending',
    routeOrder: 1
  }

  const mountCard = (props = {}) => {
    return mount(TaskCard, {
      props: { ...defaultProps, ...props },
      global: {
        stubs: {
          'van-tag': true,
          'van-icon': true
        }
      }
    })
  }

  it('renders the task card', () => {
    const wrapper = mountCard()
    expect(wrapper.find('.task-card').exists()).toBe(true)
  })

  it('displays the school name', () => {
    const wrapper = mountCard({ schoolName: 'SMP Negeri 3' })
    expect(wrapper.find('.task-card__school').text()).toBe('SMP Negeri 3')
  })

  it('displays the address', () => {
    const wrapper = mountCard({ address: 'Jl. Sudirman 5' })
    expect(wrapper.find('.task-card__address').text()).toBe('Jl. Sudirman 5')
  })

  it('displays the route order number', () => {
    const wrapper = mountCard({ routeOrder: 3 })
    expect(wrapper.find('.task-card__order').text()).toBe('3')
  })

  it('renders delivery tag with correct class', () => {
    const wrapper = mountCard({ taskType: 'delivery' })
    const tag = wrapper.find('.task-card__type-tag')
    expect(tag.exists()).toBe(true)
    expect(tag.classes()).toContain('task-card__type-tag--delivery')
  })

  it('renders pickup tag with correct class', () => {
    const wrapper = mountCard({ taskType: 'pickup' })
    const tag = wrapper.find('.task-card__type-tag')
    expect(tag.exists()).toBe(true)
    expect(tag.classes()).toContain('task-card__type-tag--pickup')
  })

  it('displays pending status label', () => {
    const wrapper = mountCard({ status: 'pending' })
    const status = wrapper.find('.task-card__status')
    expect(status.text()).toBe('Menunggu')
    expect(status.classes()).toContain('task-card__status--pending')
  })

  it('displays in_progress status label', () => {
    const wrapper = mountCard({ status: 'in_progress' })
    const status = wrapper.find('.task-card__status')
    expect(status.text()).toBe('Dalam Perjalanan')
    expect(status.classes()).toContain('task-card__status--in_progress')
  })

  it('displays completed status label', () => {
    const wrapper = mountCard({ status: 'completed' })
    const status = wrapper.find('.task-card__status')
    expect(status.text()).toBe('Selesai')
    expect(status.classes()).toContain('task-card__status--completed')
  })

  it('emits click event when card is clicked', async () => {
    const wrapper = mountCard()
    await wrapper.find('.task-card').trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')).toHaveLength(1)
  })
})
