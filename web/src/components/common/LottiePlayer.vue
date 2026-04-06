<template>
  <div
    ref="containerRef"
    class="lottie-player"
    :style="{ width, height }"
  >
    <div ref="animationRef" class="lottie-player__canvas" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import lottie from 'lottie-web'

const props = defineProps({
  src: { type: String, required: true },
  autoplay: { type: Boolean, default: true },
  loop: { type: Boolean, default: true },
  width: { type: String, default: '200px' },
  height: { type: String, default: '200px' }
})

const containerRef = ref(null)
const animationRef = ref(null)
const isVisible = ref(false)
let animInstance = null
let observer = null

const prefersReducedMotion = computed(() => {
  if (typeof window === 'undefined') return false
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
})

async function loadAndPlay() {
  if (!animationRef.value || animInstance) return

  try {
    // Dynamic import for Vite — resolve the asset path
    let animationData = null

    // Try fetching the src directly (works for absolute paths)
    const response = await fetch(props.src)
    if (response.ok) {
      animationData = await response.json()
    }

    if (!animationData) return

    animInstance = lottie.loadAnimation({
      container: animationRef.value,
      renderer: 'svg',
      loop: props.loop,
      autoplay: props.autoplay && !prefersReducedMotion.value,
      animationData
    })

    // If reduced motion, go to first frame and stop
    if (prefersReducedMotion.value) {
      animInstance.goToAndStop(0, true)
    }
  } catch (e) {
    console.warn('[LottiePlayer] Failed to load animation:', props.src, e)
  }
}

onMounted(() => {
  if (!containerRef.value) return

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting && !isVisible.value) {
          isVisible.value = true
          loadAndPlay()
        }
      })
    },
    { threshold: 0.1 }
  )
  observer.observe(containerRef.value)
})

onUnmounted(() => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
  if (animInstance) {
    animInstance.destroy()
    animInstance = null
  }
})
</script>

<style scoped>
.lottie-player {
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.lottie-player__canvas {
  width: 100%;
  height: 100%;
}
</style>
