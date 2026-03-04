<template>
  <SkeletonCard v-if="loading" :rows="2" />
  <div v-else class="summary-card">
    <div
      class="summary-card__icon"
      :style="{ backgroundColor: iconBgColor }"
    >
      <van-icon :name="icon" :color="iconColor" size="20" />
    </div>
    <div class="summary-card__content">
      <span class="summary-card__label">{{ label }}</span>
      <span class="summary-card__value">{{ value }}</span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import SkeletonCard from './SkeletonCard.vue'

const props = defineProps({
  icon: {
    type: String,
    required: true
  },
  iconColor: {
    type: String,
    default: 'var(--h-primary)'
  },
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const iconBgColor = computed(() => {
  return 'var(--h-primary-lighter)'
})
</script>

<style scoped>
.summary-card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--h-spacing-md);
}

.summary-card__icon {
  width: 40px;
  height: 40px;
  border-radius: var(--h-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.summary-card__content {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.summary-card__label {
  font-size: 12px;
  color: var(--h-text-secondary);
  line-height: 1.4;
}

.summary-card__value {
  font-size: 20px;
  font-weight: 700;
  color: var(--h-text-primary);
  line-height: 1.3;
}
</style>
