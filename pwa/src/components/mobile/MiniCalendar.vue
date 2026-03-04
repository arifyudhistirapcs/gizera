<template>
  <div class="mini-calendar">
    <div class="mini-calendar__header">
      <button class="mini-calendar__nav" @click="prevMonth" aria-label="Bulan sebelumnya">
        <van-icon name="arrow-left" size="16" />
      </button>
      <span class="mini-calendar__title">{{ monthYearLabel }}</span>
      <button class="mini-calendar__nav" @click="nextMonth" aria-label="Bulan berikutnya">
        <van-icon name="arrow" size="16" />
      </button>
    </div>

    <div class="mini-calendar__weekdays">
      <span v-for="day in weekdays" :key="day" class="mini-calendar__weekday">{{ day }}</span>
    </div>

    <div class="mini-calendar__grid">
      <div
        v-for="(cell, index) in calendarCells"
        :key="index"
        class="mini-calendar__cell"
        :class="{
          'mini-calendar__cell--empty': !cell.day,
          'mini-calendar__cell--selected': cell.isSelected,
          'mini-calendar__cell--today': cell.isToday
        }"
        @click="cell.day ? onSelectDate(cell.dateStr) : null"
      >
        <span v-if="cell.day" class="mini-calendar__day">{{ cell.day }}</span>
        <span
          v-if="cell.day && cell.status"
          class="mini-calendar__dot"
          :class="`mini-calendar__dot--${cell.status}`"
        ></span>
      </div>
    </div>

    <div class="mini-calendar__legend">
      <span class="mini-calendar__legend-item">
        <span class="mini-calendar__dot mini-calendar__dot--present"></span>
        Hadir
      </span>
      <span class="mini-calendar__legend-item">
        <span class="mini-calendar__dot mini-calendar__dot--absent"></span>
        Tidak Hadir
      </span>
      <span class="mini-calendar__legend-item">
        <span class="mini-calendar__dot mini-calendar__dot--late"></span>
        Terlambat
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  attendanceData: {
    type: Array,
    default: () => []
  },
  selectedDate: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['select-date'])

const weekdays = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab']
const monthNames = [
  'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]

const today = new Date()
const currentMonth = ref(today.getMonth())
const currentYear = ref(today.getFullYear())

const monthYearLabel = computed(() => {
  return `${monthNames[currentMonth.value]} ${currentYear.value}`
})

const attendanceMap = computed(() => {
  const map = {}
  for (const entry of props.attendanceData) {
    map[entry.date] = entry.status
  }
  return map
})

const todayStr = computed(() => {
  return formatDate(today.getFullYear(), today.getMonth(), today.getDate())
})

const calendarCells = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1).getDay()
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  const cells = []

  // Empty cells for days before the 1st
  for (let i = 0; i < firstDay; i++) {
    cells.push({ day: null, dateStr: '', status: null, isSelected: false, isToday: false })
  }

  // Day cells
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = formatDate(year, month, d)
    cells.push({
      day: d,
      dateStr,
      status: attendanceMap.value[dateStr] || null,
      isSelected: dateStr === props.selectedDate,
      isToday: dateStr === todayStr.value
    })
  }

  return cells
})

function formatDate(year, month, day) {
  const m = String(month + 1).padStart(2, '0')
  const d = String(day).padStart(2, '0')
  return `${year}-${m}-${d}`
}

function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

function onSelectDate(dateStr) {
  emit('select-date', dateStr)
}
</script>

<style scoped>
.mini-calendar {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
}

.mini-calendar__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--h-spacing-md);
}

.mini-calendar__title {
  font-size: 15px;
  font-weight: 700;
  color: var(--h-text-primary);
}

.mini-calendar__nav {
  width: 32px;
  height: 32px;
  min-height: 44px;
  min-width: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: var(--h-bg-light, #F4F7FE);
  border-radius: var(--h-radius-sm);
  color: var(--h-text-secondary);
  cursor: pointer;
  transition: background var(--h-transition-base);
  -webkit-tap-highlight-color: transparent;
}

.mini-calendar__nav:active {
  background: var(--h-primary-lighter);
}

.mini-calendar__weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: var(--h-spacing-xs);
}

.mini-calendar__weekday {
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  color: var(--h-text-light);
  padding: var(--h-spacing-xs) 0;
}

.mini-calendar__grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
}

.mini-calendar__cell {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 6px 0;
  min-height: 44px;
  cursor: pointer;
  border-radius: var(--h-radius-sm);
  transition: background var(--h-transition-base);
  -webkit-tap-highlight-color: transparent;
}

.mini-calendar__cell--empty {
  cursor: default;
}

.mini-calendar__cell:not(.mini-calendar__cell--empty):active {
  background: var(--h-primary-lighter);
}

.mini-calendar__cell--selected {
  background: var(--h-primary);
}

.mini-calendar__cell--selected .mini-calendar__day {
  color: #ffffff;
}

.mini-calendar__cell--today:not(.mini-calendar__cell--selected) {
  background: var(--h-primary-lighter);
}

.mini-calendar__day {
  font-size: 13px;
  font-weight: 500;
  color: var(--h-text-primary);
  line-height: 1;
}

.mini-calendar__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-top: 4px;
  display: inline-block;
}

.mini-calendar__dot--present {
  background-color: var(--h-success);
}

.mini-calendar__dot--absent {
  background-color: var(--h-error);
}

.mini-calendar__dot--late {
  background-color: var(--h-warning);
}

.mini-calendar__legend {
  display: flex;
  justify-content: center;
  gap: var(--h-spacing-lg);
  margin-top: var(--h-spacing-md);
  padding-top: var(--h-spacing-md);
  border-top: 1px solid var(--h-border-color, #E9EDF7);
}

.mini-calendar__legend-item {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-xs);
  font-size: 11px;
  color: var(--h-text-secondary);
}

.mini-calendar__legend-item .mini-calendar__dot {
  margin-top: 0;
}
</style>
