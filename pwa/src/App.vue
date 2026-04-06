<template>
  <router-view />
</template>

<script setup>
import '@/styles/horizon-mobile.css'
</script>

<style>
body {
  margin: 0;
  padding: 0;
  font-family: 'Urbanist', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  min-height: 100vh;
}
</style>
