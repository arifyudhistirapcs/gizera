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
  border-radius: 8px;
  border: 1px solid var(--h-border-color);
  padding: var(--h-spacing-lg);
}
</style>
