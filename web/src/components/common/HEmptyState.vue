<template>
  <div class="h-empty-state">
    <div class="h-empty-state__visual">
      <LottiePlayer
        v-if="lottie"
        :src="lottie"
        :autoplay="true"
        :loop="true"
        width="160px"
        height="160px"
      />
      <img
        v-else-if="illustration"
        :src="illustration"
        :alt="title || 'Empty state'"
        class="h-empty-state__illustration"
      />
    </div>
    <h3 v-if="title" class="h-empty-state__title">{{ title }}</h3>
    <p v-if="description" class="h-empty-state__description">{{ description }}</p>
    <div v-if="$slots.default" class="h-empty-state__action">
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
.h-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.h-empty-state__visual {
  margin-bottom: 16px;
}

.h-empty-state__illustration {
  width: 160px;
  height: auto;
  max-height: 160px;
  object-fit: contain;
}

.h-empty-state__title {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 600;
  color: #303030;
  font-family: 'Urbanist', sans-serif;
}

.h-empty-state__description {
  margin: 0;
  font-size: 14px;
  color: #6B6B6B;
  max-width: 320px;
  line-height: 1.5;
  font-family: 'Urbanist', sans-serif;
}
</style>