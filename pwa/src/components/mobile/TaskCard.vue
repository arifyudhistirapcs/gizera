<template>
  <div class="task-card" @click="$emit('click')">
    <div class="task-card__header">
      <span class="task-card__order">{{ routeOrder }}</span>
      <van-tag
        :class="['task-card__type-tag', `task-card__type-tag--${taskType}`]"
        round
        size="medium"
      >
        {{ taskType === 'delivery' ? 'Pengiriman' : 'Pengambilan' }}
      </van-tag>
    </div>
    <div class="task-card__body">
      <h3 class="task-card__school">{{ schoolName }}</h3>
      <p class="task-card__address">{{ address }}</p>
    </div>
    <div class="task-card__footer">
      <span :class="['task-card__status', `task-card__status--${status}`]">
        {{ statusLabel }}
      </span>
      <van-icon name="arrow" class="task-card__arrow" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  schoolName: {
    type: String,
    required: true
  },
  address: {
    type: String,
    required: true
  },
  taskType: {
    type: String,
    required: true,
    validator: (val) => ['delivery', 'pickup'].includes(val)
  },
  status: {
    type: String,
    required: true,
    validator: (val) => ['pending', 'in_progress', 'arrived', 'completed'].includes(val)
  },
  routeOrder: {
    type: Number,
    required: true
  }
})

defineEmits(['click'])

const statusLabel = computed(() => {
  const labels = {
    pending: 'Menunggu',
    in_progress: 'Dalam Perjalanan',
    arrived: 'Tiba di Lokasi',
    completed: 'Selesai'
  }
  return labels[props.status] || props.status
})
</script>

<style scoped>
.task-card {
  background: var(--h-bg-card);
  border-radius: 8px;
  border: 1px solid var(--h-border-color);
  padding: var(--h-spacing-lg);
  cursor: pointer;
  transition: transform var(--h-transition-base);
}

.task-card:active {
  transform: scale(0.98);
}

.task-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--h-spacing-sm);
}

.task-card__order {
  width: 28px;
  height: 28px;
  border-radius: var(--h-radius-full);
  background: #F0F0F0;
  color: #303030;
  font-weight: 700;
  font-size: 13px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.task-card__type-tag {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border: none;
  min-height: auto;
  min-width: auto;
  border-radius: 4px !important;
}

.task-card__type-tag--delivery {
  background: #DCF0D8 !important;
  color: #303030 !important;
}

.task-card__type-tag--pickup {
  background: #FEF3C7 !important;
  color: #303030 !important;
}

.task-card__body {
  margin-bottom: var(--h-spacing-sm);
}

.task-card__school {
  font-size: 15px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 2px 0;
  line-height: 1.4;
}

.task-card__address {
  font-size: 13px;
  color: var(--h-text-secondary);
  margin: 0;
  line-height: 1.4;
}

.task-card__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.task-card__status {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 8px;
  border-radius: var(--h-radius-sm);
}

.task-card__status--pending {
  color: var(--h-warning);
  background: rgba(255, 181, 71, 0.1);
}

.task-card__status--in_progress {
  color: #303030;
  background: #F0F0F0;
}

.task-card__status--arrived {
  color: #1989fa;
  background: rgba(25, 137, 250, 0.1);
}

.task-card__status--completed {
  color: var(--h-success);
  background: rgba(5, 205, 153, 0.1);
}

.task-card__arrow {
  color: var(--h-text-light);
  font-size: 14px;
}
</style>
