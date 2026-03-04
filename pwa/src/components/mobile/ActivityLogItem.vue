<template>
  <div class="activity-log-item">
    <div class="activity-log-item__left">
      <div :class="['activity-log-item__type-icon', `activity-log-item__type-icon--${activityType}`]">
        <van-icon :name="typeIcon" size="18" color="#ffffff" />
      </div>
    </div>
    <div class="activity-log-item__content">
      <h4 class="activity-log-item__name">{{ employeeName }}</h4>
      <p class="activity-log-item__type">{{ typeLabel }}</p>
    </div>
    <div class="activity-log-item__right">
      <span class="activity-log-item__timestamp">{{ timestamp }}</span>
      <span :class="['activity-log-item__status', statusClass]">
        {{ status }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  employeeName: {
    type: String,
    required: true
  },
  activityType: {
    type: String,
    required: true,
    validator: (val) => ['attendance', 'delivery', 'pickup'].includes(val)
  },
  timestamp: {
    type: String,
    required: true
  },
  status: {
    type: String,
    required: true
  }
})

const typeLabel = computed(() => {
  const labels = {
    attendance: 'Absensi',
    delivery: 'Pengiriman',
    pickup: 'Pengambilan'
  }
  return labels[props.activityType] || props.activityType
})

const typeIcon = computed(() => {
  const icons = {
    attendance: 'clock-o',
    delivery: 'logistics',
    pickup: 'revoke'
  }
  return icons[props.activityType] || 'info-o'
})

const statusClass = computed(() => {
  const s = props.status.toLowerCase()
  if (s === 'hadir' || s === 'selesai' || s === 'completed' || s === 'present') {
    return 'activity-log-item__status--success'
  }
  if (s === 'terlambat' || s === 'late' || s === 'pending' || s === 'menunggu') {
    return 'activity-log-item__status--warning'
  }
  if (s === 'tidak hadir' || s === 'absent' || s === 'gagal' || s === 'failed') {
    return 'activity-log-item__status--error'
  }
  if (s === 'dalam perjalanan' || s === 'in_progress') {
    return 'activity-log-item__status--primary'
  }
  return 'activity-log-item__status--default'
})
</script>

<style scoped>
.activity-log-item {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--h-spacing-md);
}

.activity-log-item__left {
  flex-shrink: 0;
}

.activity-log-item__type-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--h-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
}

.activity-log-item__type-icon--attendance {
  background: var(--h-primary);
}

.activity-log-item__type-icon--delivery {
  background: var(--h-success);
}

.activity-log-item__type-icon--pickup {
  background: var(--h-warning);
}

.activity-log-item__content {
  flex: 1;
  min-width: 0;
}

.activity-log-item__name {
  font-size: 14px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 2px 0;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.activity-log-item__type {
  font-size: 12px;
  color: var(--h-text-secondary);
  margin: 0;
  line-height: 1.4;
}

.activity-log-item__right {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.activity-log-item__timestamp {
  font-size: 12px;
  font-weight: 500;
  color: var(--h-text-secondary);
  background: var(--h-bg-light);
  padding: 2px 8px;
  border-radius: var(--h-radius-sm);
}

.activity-log-item__status {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: var(--h-radius-sm);
}

.activity-log-item__status--success {
  color: var(--h-success);
  background: rgba(5, 205, 153, 0.1);
}

.activity-log-item__status--warning {
  color: var(--h-warning);
  background: rgba(255, 181, 71, 0.1);
}

.activity-log-item__status--error {
  color: var(--h-error);
  background: rgba(238, 93, 80, 0.1);
}

.activity-log-item__status--primary {
  color: var(--h-primary);
  background: var(--h-primary-lighter);
}

.activity-log-item__status--default {
  color: var(--h-text-secondary);
  background: var(--h-bg-light);
}
</style>
