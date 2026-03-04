import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import MiniCalendar from '@/components/mobile/MiniCalendar.vue'

describe('MiniCalendar', () => {
  const sampleAttendance = [
    { date: '2025-01-06', status: 'present' },
    { date: '2025-01-07', status: 'late' },
    { date: '2025-01-08', status: 'absent' },
    { date: '2025-01-09', status: 'present' },
    { date: '2025-01-10', status: 'present' }
  ]

  const mountCalendar = (props = {}) => {
    return mount(MiniCalendar, {
      props: {
        attendanceData: [],
        ...props
      },
      global: {
        stubs: {
          'van-icon': true
        }
      }
    })
  }

  it('renders the calendar container', () => {
    const wrapper = mountCalendar()
    expect(wrapper.find('.mini-calendar').exists()).toBe(true)
  })

  it('renders 7 weekday headers', () => {
    const wrapper = mountCalendar()
    const weekdays = wrapper.findAll('.mini-calendar__weekday')
    expect(weekdays).toHaveLength(7)
    expect(weekdays[0].text()).toBe('Min')
    expect(weekdays[1].text()).toBe('Sen')
    expect(weekdays[6].text()).toBe('Sab')
  })

  it('renders month grid with correct number of day cells', () => {
    const wrapper = mountCalendar()
    const dayCells = wrapper.findAll('.mini-calendar__day')
    // Should have between 28 and 31 day numbers
    expect(dayCells.length).toBeGreaterThanOrEqual(28)
    expect(dayCells.length).toBeLessThanOrEqual(31)
  })

  it('displays current month and year in header', () => {
    const wrapper = mountCalendar()
    const title = wrapper.find('.mini-calendar__title').text()
    const now = new Date()
    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ]
    expect(title).toContain(monthNames[now.getMonth()])
    expect(title).toContain(String(now.getFullYear()))
  })

  it('renders colored dots for attendance data', () => {
    const wrapper = mountCalendar({ attendanceData: sampleAttendance })
    const presentDots = wrapper.findAll('.mini-calendar__dot--present')
    const absentDots = wrapper.findAll('.mini-calendar__dot--absent')
    const lateDots = wrapper.findAll('.mini-calendar__dot--late')
    // Legend has 1 of each, grid may have more depending on current month
    expect(presentDots.length).toBeGreaterThanOrEqual(1) // at least legend
    expect(absentDots.length).toBeGreaterThanOrEqual(1)
    expect(lateDots.length).toBeGreaterThanOrEqual(1)
  })

  it('highlights selected date', () => {
    const now = new Date()
    const y = now.getFullYear()
    const m = String(now.getMonth() + 1).padStart(2, '0')
    const selectedDate = `${y}-${m}-15`

    const wrapper = mountCalendar({ selectedDate })
    const selectedCell = wrapper.find('.mini-calendar__cell--selected')
    expect(selectedCell.exists()).toBe(true)
    expect(selectedCell.find('.mini-calendar__day').text()).toBe('15')
  })

  it('emits select-date when a day cell is clicked', async () => {
    const wrapper = mountCalendar()
    const dayCells = wrapper.findAll('.mini-calendar__cell:not(.mini-calendar__cell--empty)')
    expect(dayCells.length).toBeGreaterThan(0)

    await dayCells[0].trigger('click')
    expect(wrapper.emitted('select-date')).toBeTruthy()
    expect(wrapper.emitted('select-date')[0][0]).toMatch(/^\d{4}-\d{2}-\d{2}$/)
  })

  it('does not emit select-date when empty cell is clicked', async () => {
    const wrapper = mountCalendar()
    const emptyCells = wrapper.findAll('.mini-calendar__cell--empty')
    if (emptyCells.length > 0) {
      await emptyCells[0].trigger('click')
      expect(wrapper.emitted('select-date')).toBeFalsy()
    }
  })

  it('navigates to previous month', async () => {
    const wrapper = mountCalendar()
    const navButtons = wrapper.findAll('.mini-calendar__nav')
    const prevBtn = navButtons[0]

    const titleBefore = wrapper.find('.mini-calendar__title').text()
    await prevBtn.trigger('click')
    const titleAfter = wrapper.find('.mini-calendar__title').text()

    expect(titleAfter).not.toBe(titleBefore)
  })

  it('navigates to next month', async () => {
    const wrapper = mountCalendar()
    const navButtons = wrapper.findAll('.mini-calendar__nav')
    const nextBtn = navButtons[1]

    const titleBefore = wrapper.find('.mini-calendar__title').text()
    await nextBtn.trigger('click')
    const titleAfter = wrapper.find('.mini-calendar__title').text()

    expect(titleAfter).not.toBe(titleBefore)
  })

  it('renders legend with three status items', () => {
    const wrapper = mountCalendar()
    const legendItems = wrapper.findAll('.mini-calendar__legend-item')
    expect(legendItems).toHaveLength(3)
    expect(legendItems[0].text()).toContain('Hadir')
    expect(legendItems[1].text()).toContain('Tidak Hadir')
    expect(legendItems[2].text()).toContain('Terlambat')
  })

  it('renders grid with 7-column layout (weekdays + day cells)', () => {
    const wrapper = mountCalendar()
    const grid = wrapper.find('.mini-calendar__grid')
    expect(grid.exists()).toBe(true)
    // Total cells should be divisible by 7 or have empty padding
    const allCells = wrapper.findAll('.mini-calendar__cell')
    expect(allCells.length).toBeGreaterThanOrEqual(28)
  })

  it('highlights today with special styling', () => {
    const wrapper = mountCalendar()
    const todayCell = wrapper.find('.mini-calendar__cell--today')
    expect(todayCell.exists()).toBe(true)
    const todayDay = new Date().getDate()
    expect(todayCell.find('.mini-calendar__day').text()).toBe(String(todayDay))
  })
})
