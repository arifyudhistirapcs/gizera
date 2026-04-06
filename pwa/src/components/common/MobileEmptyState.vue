<template>
  <div class="mobile-empty-state">
    <div class="mobile-empty-state__visual">
      <LottiePlayer
        v-if="lottie"
        :src="lottie"
        :autoplay="true"
        :loop="true"
        width="120px"
        height="120px"
      />
      <img
        v-else-if="illustration"
        :src="illustration"
        :alt="title || 'Empty state'"
        class="mobile-empty-state__illustration"
      />
    </div>
    <h3 v-if="title" class="mobile-empty-state__title">{{ title }}</h3>
    <p v-if="description" class="mobile-empty-state__description">{{ description }}</p>
    <div v-if="$slots.default" class="mobile-empty-state__action">
      <slot />
    </div>
  </div>
</template>

<script setup>
import LottiePlayer from './LottiePlayer.vue'

defineProps({
  illustration: { type: String, default: '' },
  lottie: { type: String, default: '' },
  title: { type: String, default: '' },
  description: { type: String, default: '' }
})
</script>

<style scoped>
.mobile-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
}

.mobile-empty-state__visual {
  margin-bottom: 12px;
}

.mobile-empty-state__illustration {
  width: 120px;
  height: auto;
  max-height: 120px;
  object-fit: contain;
}

.mobile-empty-state__title {
  margin: 0 0 6px;
  font-size: 16px;
  font-weight: 600;
  color: #303030;
  font-family: 'Urbanist', sans-serif;
}

.mobile-empty-state__description {
  margin: 0;
  font-size: 13px;
  color: #6B6B6B;
  max-width: 260px;
  line-height: 1.5;
  font-family: 'Urbanist', sans-serif;
}
</style>