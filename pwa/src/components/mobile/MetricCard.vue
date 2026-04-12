<template>
  <div class="metric-card h-card">
    <div class="metric-card__icon" :style="{ backgroundColor: iconColor }">
      <van-icon :name="icon" size="20" color="#303030" />
    </div>
    <div class="metric-card__content">
      <span class="metric-card__label">{{ label }}</span>
      <span class="metric-card__value">{{ value }}</span>
      <span v-if="trend" class="metric-card__trend" :class="trendClass">
        <van-icon :name="trendIcon" size="10" />
        {{ trend }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  icon: {
    type: String,
    required: true
  },
  iconColor: {
    type: String,
    default: '#F0F0F0'
  },
  label: {
    type: String,
    required: true
  },
  value: {
    type: [String, Number],
    required: true
  },
  trend: {
    type: String,
    default: ''
  },
  trendUp: {
    type: Boolean,
    default: false
  },
  trendDown: {
    type: Boolean,
    default: false
  }
})

const trendClass = computed(() => {
  if (props.trendUp) return 'trend-up'
  if (props.trendDown) return 'trend-down'
  return ''
})

const trendIcon = computed(() => {
  if (props.trendUp) return 'arrow-up'
  if (props.trendDown) return 'arrow-down'
  return ''
})
</script>

<style scoped>
.metric-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  min-height: 80px;
  border: none;
  border-radius: 16px;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.metric-card__icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.metric-card__content {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-width: 0;
}

.metric-card__label {
  font-size: 11px;
  color: #8C8C8C;
  font-weight: 500;
  line-height: 1.2;
}

.metric-card__value {
  font-size: 20px;
  font-weight: 700;
  color: #303030;
  line-height: 1;
}

.metric-card__trend {
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 2px;
  margin-top: 2px;
}

.metric-card__trend.trend-up {
  color: #764AF1;
}

.metric-card__trend.trend-down {
  color: #F82C17;
}
</style>
