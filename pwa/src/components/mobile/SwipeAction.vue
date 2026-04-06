<template>
  <div class="swipe-action" :class="{ 'swipe-action--disabled': loading || (!canCheckIn && !canCheckOut) }">
    <div class="swipe-action__track" ref="trackRef">
      <!-- Directional arrows -->
      <div class="swipe-action__arrows swipe-action__arrows--left" v-if="canCheckOut">
        <van-icon name="arrow-left" color="rgba(255,255,255,0.5)" size="16" />
        <van-icon name="arrow-left" color="rgba(255,255,255,0.3)" size="16" />
      </div>
      <div class="swipe-action__label">
        <span v-if="loading">
          <van-loading size="20" color="#fff" />
        </span>
        <span v-else-if="activeDirection === 'right' && canCheckIn" class="swipe-action__label-text">
          Check In
        </span>
        <span v-else-if="activeDirection === 'left' && canCheckOut" class="swipe-action__label-text">
          Check Out
        </span>
        <span v-else class="swipe-action__label-text">
          {{ canCheckIn && canCheckOut ? 'Geser untuk absen' : canCheckIn ? 'Geser kanan → Check In' : canCheckOut ? '← Geser kiri Check Out' : '' }}
        </span>
      </div>
      <div class="swipe-action__arrows swipe-action__arrows--right" v-if="canCheckIn">
        <van-icon name="arrow" color="rgba(255,255,255,0.3)" size="16" />
        <van-icon name="arrow" color="rgba(255,255,255,0.5)" size="16" />
      </div>

      <!-- Slider thumb -->
      <div
        class="swipe-action__thumb"
        ref="thumbRef"
        :style="thumbStyle"
        @touchstart.prevent="onTouchStart"
        @touchmove.prevent="onTouchMove"
        @touchend.prevent="onTouchEnd"
      >
        <van-loading v-if="loading" size="20" color="var(--h-primary)" />
        <van-icon v-else-if="activeDirection === 'right'" name="arrow" color="var(--h-primary)" size="20" />
        <van-icon v-else-if="activeDirection === 'left'" name="arrow-left" color="var(--h-primary)" size="20" />
        <van-icon v-else name="exchange" color="var(--h-primary)" size="20" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  canCheckIn: {
    type: Boolean,
    default: false
  },
  canCheckOut: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['check-in', 'check-out'])

const trackRef = ref(null)
const thumbRef = ref(null)

const THUMB_SIZE = 48
const SWIPE_THRESHOLD = 0.6 // 60% of track width

// Touch state
const isDragging = ref(false)
const startX = ref(0)
const currentOffset = ref(0)

const activeDirection = computed(() => {
  // For check-out only (thumb starts from right), positive offset means swiping left
  if (props.canCheckOut && !props.canCheckIn) {
    if (currentOffset.value > 10) return 'left'
    return null
  }
  
  // For check-in (thumb starts from left), positive offset means swiping right
  if (currentOffset.value > 10) return 'right'
  if (currentOffset.value < -10) return 'left'
  return null
})

const thumbStyle = computed(() => {
  const trackEl = trackRef.value
  if (!trackEl) {
    // Default position: left for check-in, right for check-out only
    const defaultLeft = props.canCheckOut && !props.canCheckIn ? 'auto' : '4px'
    const defaultRight = props.canCheckOut && !props.canCheckIn ? '4px' : 'auto'
    return { left: defaultLeft, right: defaultRight, transform: 'translateX(0px)' }
  }

  const trackWidth = trackEl.offsetWidth
  const maxOffset = trackWidth - THUMB_SIZE - 8 // 8px padding

  let offset = currentOffset.value

  // Clamp offset to track bounds (always positive for both directions now)
  offset = Math.max(0, Math.min(offset, maxOffset))

  // Position thumb: left for check-in, right for check-out only
  const defaultLeft = props.canCheckOut && !props.canCheckIn ? 'auto' : '4px'
  const defaultRight = props.canCheckOut && !props.canCheckIn ? '4px' : 'auto'
  
  // For check-out (starts from right), translate in negative direction (left)
  const translateDirection = props.canCheckOut && !props.canCheckIn ? -1 : 1

  return {
    left: defaultLeft,
    right: defaultRight,
    transform: `translateX(${offset * translateDirection}px)`,
    transition: isDragging.value ? 'none' : 'transform 300ms ease-out'
  }
})

function onTouchStart(e) {
  if (props.loading || (!props.canCheckIn && !props.canCheckOut)) return

  isDragging.value = true
  startX.value = e.touches[0].clientX
  currentOffset.value = 0
}

function onTouchMove(e) {
  if (!isDragging.value) return

  let deltaX = e.touches[0].clientX - startX.value
  
  // For check-out only (thumb starts from right), invert the delta
  // so swiping left (negative deltaX) becomes positive offset
  if (props.canCheckOut && !props.canCheckIn) {
    deltaX = -deltaX
  }
  
  currentOffset.value = deltaX
}

function onTouchEnd() {
  if (!isDragging.value) return

  isDragging.value = false

  const trackEl = trackRef.value
  if (!trackEl) {
    currentOffset.value = 0
    return
  }

  const trackWidth = trackEl.offsetWidth
  const threshold = (trackWidth - THUMB_SIZE) * SWIPE_THRESHOLD

  // For check-in: swipe right (positive offset)
  if (currentOffset.value > threshold && props.canCheckIn) {
    emit('check-in')
  } 
  // For check-out: swipe left from right position (positive offset after inversion)
  else if (currentOffset.value > threshold && props.canCheckOut) {
    emit('check-out')
  }

  // Snap back
  currentOffset.value = 0
}
</script>

<style scoped>
.swipe-action {
  width: 100%;
  padding: 0 var(--h-spacing-lg);
  box-sizing: border-box;
}

.swipe-action--disabled {
  opacity: 0.5;
  pointer-events: none;
}

.swipe-action__track {
  position: relative;
  height: 56px;
  background: #303030;
  border-radius: var(--h-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 0 4px;
  user-select: none;
  -webkit-user-select: none;
  touch-action: none;
}

.swipe-action__label {
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
  pointer-events: none;
}

.swipe-action__label-text {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}

.swipe-action__arrows {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 2px;
  z-index: 1;
  pointer-events: none;
}

.swipe-action__arrows--left {
  left: 16px;
}

.swipe-action__arrows--right {
  right: 16px;
}

.swipe-action__thumb {
  position: absolute;
  width: 48px;
  height: 48px;
  min-height: 48px;
  min-width: 48px;
  background: #FFFFFF;
  border-radius: var(--h-radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--h-shadow-md);
  cursor: grab;
  z-index: 2;
  touch-action: none;
}

.swipe-action__thumb:active {
  cursor: grabbing;
}
</style>
