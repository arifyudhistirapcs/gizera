---
inclusion: fileMatch
fileMatchPattern: "web/src/**/*.vue,web/src/**/*.js,pwa/src/**/*.vue,pwa/src/**/*.js"
---

# Vue 3 Patterns

## Component Structure
```vue
<template>
  <!-- template -->
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
// imports, reactive state, computed, methods, lifecycle
</script>

<style scoped>
/* scoped styles */
</style>
```

## API Service Pattern
```js
import api from './api'

const myService = {
  getItems(params = {}) {
    return api.get('/endpoint', { params })
  },
  createItem(data) {
    return api.post('/endpoint', data)
  }
}
export default myService
```

## Response Parsing
Backend returns `{ success, data, message }`. Parse consistently:
```js
const response = await myService.getItems()
const items = response.data?.data || []
```

## State Pattern
```js
const loading = ref(false)
const error = ref(null)
const data = ref([])

async function fetchData() {
  loading.value = true
  error.value = null
  try {
    const res = await myService.getItems()
    data.value = res.data?.data || []
  } catch (err) {
    error.value = err.response?.data?.message || 'Gagal memuat data'
  } finally {
    loading.value = false
  }
}
```

## Web (Ant Design Vue)
- Use `a-table`, `a-form`, `a-select`, `a-tag`, `a-card`, `a-descriptions`
- Color-code risk levels: rendah=green, sedang=orange, tinggi=red
- Use `a-page-header` for page titles

## PWA (Vant UI)
- Use `van-nav-bar`, `van-cell-group`, `van-cell`, `van-button`, `van-tag`
- Use `van-loading` for loading states, `showToast`/`showFailToast` for feedback
- Background color: `#F8FDEA`, card border-radius: `16px`
- Bottom nav uses floating center button pattern
