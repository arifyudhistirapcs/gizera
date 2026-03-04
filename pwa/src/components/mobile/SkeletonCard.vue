<template>
  <div class="skeleton-card">
    <van-skeleton
      :row="rows"
      :avatar="avatar"
      :title="title"
      avatar-size="40px"
      round
      loading
    />
  </div>
</template>

<script setup>
defineProps({
  rows: {
    type: Number,
    default: 3
  },
  avatar: {
    type: Boolean,
    default: false
  },
  title: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
.skeleton-card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-lg);
}
</style>
