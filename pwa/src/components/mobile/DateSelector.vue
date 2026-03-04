<template>
  <div class="date-selector" @click="showCalendar = true">
    <div class="date-selector__display">
      <van-icon name="calendar-o" size="18" color="var(--h-primary)" />
      <span class="date-selector__text">{{ formattedDate }}</span>
      <van-icon name="arrow-down" size="14" color="var(--h-text-secondary)" />
    </div>
  </div>

  <van-calendar
    v-model:show="showCalendar"
    :default-date="modelValue"
    :min-date="minDate"
    :max-date="maxDate"
    @confirm="onConfirm"
    color="var(--h-primary)"
    :show-confirm="false"
  />
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Date,
    required: true
  },
  minDate: {
    type: Date,
    default: () => new Date(new Date().getFullYear() - 1, 0, 1)
  },
  maxDate: {
    type: Date,
    default: () => new Date(new Date().getFullYear() + 1, 11, 31)
  }
})

const emit = defineEmits(['update:modelValue'])

const showCalendar = ref(false)

const formattedDate = computed(() => {
  const date = props.modelValue
  const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu']
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun', 'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des']

  const dayName = days[date.getDay()]
  const day = date.getDate()
  const month = months[date.getMonth()]
  const year = date.getFullYear()

  return `${dayName}, ${day} ${month} ${year}`
})

function onConfirm(date) {
  emit('update:modelValue', date)
  showCalendar.value = false
}
</script>

<style scoped>
.date-selector {
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}

.date-selector__display {
  display: inline-flex;
  align-items: center;
  gap: var(--h-spacing-sm);
  padding: var(--h-spacing-sm) var(--h-spacing-md);
  background: var(--h-bg-card);
  border-radius: var(--h-radius-md);
  box-shadow: var(--h-shadow-sm);
  transition: box-shadow var(--h-transition-base);
  min-height: 44px;
}

.date-selector__display:active {
  box-shadow: var(--h-shadow-md);
}

.date-selector__text {
  font-size: 14px;
  font-weight: 600;
  color: var(--h-text-primary);
  white-space: nowrap;
}
</style>
