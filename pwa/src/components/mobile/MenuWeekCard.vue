<template>
  <div class="menu-week-card" @click="$emit('click')">
    <div class="menu-week-card__header">
      <h3 class="menu-week-card__period">{{ weekPeriod }}</h3>
      <span
        :class="['menu-week-card__status', `menu-week-card__status--${approvalStatus}`]"
      >
        {{ statusLabel }}
      </span>
    </div>
    <div class="menu-week-card__body">
      <div class="menu-week-card__info">
        <van-icon name="notes-o" size="16" color="var(--h-text-secondary)" />
        <span class="menu-week-card__count">{{ menuCount }} menu</span>
      </div>
    </div>
    <div v-if="canApprove && approvalStatus === 'pending'" class="menu-week-card__actions">
      <van-button
        type="primary"
        size="small"
        round
        class="menu-week-card__btn menu-week-card__btn--approve"
        @click.stop="$emit('approve')"
      >
        Approve
      </van-button>
      <van-button
        type="danger"
        size="small"
        round
        class="menu-week-card__btn menu-week-card__btn--reject"
        @click.stop="$emit('reject')"
      >
        Reject
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  weekPeriod: {
    type: String,
    required: true
  },
  approvalStatus: {
    type: String,
    required: true,
    validator: (val) => ['pending', 'approved', 'rejected'].includes(val)
  },
  menuCount: {
    type: Number,
    required: true
  },
  canApprove: {
    type: Boolean,
    default: false
  }
})

defineEmits(['approve', 'reject', 'click'])

const statusLabel = computed(() => {
  const labels = {
    pending: 'Pending',
    approved: 'Approved',
    rejected: 'Rejected'
  }
  return labels[props.approvalStatus] || props.approvalStatus
})
</script>

<style scoped>
.menu-week-card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
  cursor: pointer;
  transition: transform var(--h-transition-base);
}

.menu-week-card:active {
  transform: scale(0.98);
}

.menu-week-card__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--h-spacing-sm);
}

.menu-week-card__period {
  font-size: 15px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0;
  line-height: 1.4;
}

.menu-week-card__status {
  font-size: 12px;
  font-weight: 500;
  padding: 2px 10px;
  border-radius: var(--h-radius-full);
}

.menu-week-card__status--pending {
  color: var(--h-warning);
  background: rgba(255, 181, 71, 0.1);
}

.menu-week-card__status--approved {
  color: var(--h-success);
  background: rgba(5, 205, 153, 0.1);
}

.menu-week-card__status--rejected {
  color: var(--h-error);
  background: rgba(238, 93, 80, 0.1);
}

.menu-week-card__body {
  margin-bottom: var(--h-spacing-sm);
}

.menu-week-card__info {
  display: flex;
  align-items: center;
  gap: var(--h-spacing-xs);
}

.menu-week-card__count {
  font-size: 13px;
  color: var(--h-text-secondary);
  line-height: 1.4;
}

.menu-week-card__actions {
  display: flex;
  gap: var(--h-spacing-sm);
  padding-top: var(--h-spacing-sm);
  border-top: 1px solid var(--h-border-light);
}

.menu-week-card__btn {
  flex: 1;
  font-weight: 600;
  font-size: 13px;
}

.menu-week-card__btn--approve {
  background: var(--h-success) !important;
  border-color: var(--h-success) !important;
}

.menu-week-card__btn--reject {
  background: var(--h-error) !important;
  border-color: var(--h-error) !important;
}
</style>
