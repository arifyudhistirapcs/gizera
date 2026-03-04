import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DateSelector from '@/components/mobile/DateSelector.vue'

describe('DateSelector', () => {
  const defaultProps = {
    modelValue: new Date(2025, 0, 6) // Senin, 6 Jan 2025
  }

  const mountSelector = (props = {}) => {
    return mount(DateSelector, {
      props: { ...defaultProps, ...props },
      global: {
        stubs: {
          'van-icon': true,
          'van-calendar': true
        }
      }
    })
  }

  it('renders the date selector display', () => {
    const wrapper = mountSelector()
    expect(wrapper.find('.date-selector').exists()).toBe(true)
    expect(wrapper.find('.date-selector__display').exists()).toBe(true)
  })

  it('displays formatted date in Indonesian format', () => {
    const wrapper = mountSelector()
    expect(wrapper.find('.date-selector__text').text()).toBe('Senin, 6 Jan 2025')
  })

  it('formats different dates correctly', () => {
    const wrapper = mountSelector({ modelValue: new Date(2025, 6, 20) }) // Minggu, 20 Jul 2025
    expect(wrapper.find('.date-selector__text').text()).toBe('Minggu, 20 Jul 2025')
  })

  it('opens calendar popup on tap', async () => {
    const wrapper = mountSelector()
    const calendar = wrapper.findComponent({ name: 'van-calendar' })
    expect(calendar.attributes('show')).toBe('false')

    await wrapper.find('.date-selector').trigger('click')
    expect(wrapper.findComponent({ name: 'van-calendar' }).attributes('show')).toBe('true')
  })

  it('emits update:modelValue when date is confirmed', async () => {
    const wrapper = mountSelector()
    const newDate = new Date(2025, 1, 14)

    const calendar = wrapper.findComponent({ name: 'van-calendar' })
    await calendar.vm.$emit('confirm', newDate)

    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([newDate])
  })

  it('closes calendar after date confirmation', async () => {
    const wrapper = mountSelector()

    // Open calendar
    await wrapper.find('.date-selector').trigger('click')

    // Confirm a date
    const calendar = wrapper.findComponent({ name: 'van-calendar' })
    await calendar.vm.$emit('confirm', new Date(2025, 1, 14))

    // Calendar should be closed
    expect(wrapper.findComponent({ name: 'van-calendar' }).attributes('show')).toBe('false')
  })

  it('passes minDate and maxDate to van-calendar', () => {
    const minDate = new Date(2024, 0, 1)
    const maxDate = new Date(2025, 11, 31)
    const wrapper = mountSelector({ minDate, maxDate })

    const calendar = wrapper.findComponent({ name: 'van-calendar' })
    expect(calendar.attributes('min-date')).toBeDefined()
    expect(calendar.attributes('max-date')).toBeDefined()
  })

  it('supports v-model binding pattern', async () => {
    const wrapper = mountSelector()
    const newDate = new Date(2025, 5, 15)

    const calendar = wrapper.findComponent({ name: 'van-calendar' })
    await calendar.vm.$emit('confirm', newDate)

    const emitted = wrapper.emitted('update:modelValue')
    expect(emitted).toHaveLength(1)
    expect(emitted[0][0]).toEqual(newDate)
  })

  it('renders calendar and arrow icons', () => {
    const wrapper = mountSelector()
    const icons = wrapper.findAll('van-icon-stub')
    expect(icons.length).toBe(2)
    expect(icons[0].attributes('name')).toBe('calendar-o')
    expect(icons[1].attributes('name')).toBe('arrow-down')
  })
})
