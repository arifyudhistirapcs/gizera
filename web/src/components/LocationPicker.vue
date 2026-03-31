<template>
  <div class="location-picker">
    <!-- Search bar -->
    <div class="location-picker__search">
      <a-input-search
        v-model:value="searchQuery"
        placeholder="Cari alamat atau nama tempat..."
        :loading="searching"
        @search="handleSearch"
        allow-clear
        style="margin-bottom: 8px"
      />
      <!-- Search results -->
      <div v-if="searchResults.length" class="location-picker__results">
        <div
          v-for="result in searchResults"
          :key="result.place_id"
          class="location-picker__result-item"
          @click="selectSearchResult(result)"
        >
          <EnvironmentOutlined style="margin-right: 6px; color: #5A4372" />
          <span>{{ result.display_name }}</span>
        </div>
      </div>
    </div>

    <!-- Map -->
    <div ref="mapContainer" class="location-picker__map"></div>

    <!-- Controls -->
    <div class="location-picker__controls">
      <a-button size="small" @click="getMyLocation" :loading="gettingLocation">
        <template #icon><AimOutlined /></template>
        Lokasi Saya
      </a-button>
      <span class="location-picker__coords">
        {{ latitude?.toFixed(6) || '-' }}, {{ longitude?.toFixed(6) || '-' }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import { AimOutlined, EnvironmentOutlined } from '@ant-design/icons-vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps({
  latitude: { type: Number, default: 0 },
  longitude: { type: Number, default: 0 }
})

const emit = defineEmits(['update:latitude', 'update:longitude'])

const mapContainer = ref(null)
const searchQuery = ref('')
const searching = ref(false)
const gettingLocation = ref(false)
const searchResults = ref([])

let map = null
let marker = null

const initMap = () => {
  if (!mapContainer.value) return

  const lat = props.latitude || -6.2088
  const lng = props.longitude || 106.8456
  const zoom = (props.latitude && props.longitude) ? 15 : 5

  map = L.map(mapContainer.value).setView([lat, lng], zoom)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap',
    maxZoom: 19
  }).addTo(map)

  // Place marker if coords exist
  if (props.latitude && props.longitude) {
    placeMarker(props.latitude, props.longitude)
  }

  // Click on map to place marker
  map.on('click', (e) => {
    placeMarker(e.latlng.lat, e.latlng.lng)
    emit('update:latitude', e.latlng.lat)
    emit('update:longitude', e.latlng.lng)
    searchResults.value = []
  })
}

const placeMarker = (lat, lng) => {
  if (marker) {
    marker.setLatLng([lat, lng])
  } else {
    marker = L.marker([lat, lng], {
      draggable: true
    }).addTo(map)

    marker.on('dragend', () => {
      const pos = marker.getLatLng()
      emit('update:latitude', pos.lat)
      emit('update:longitude', pos.lng)
    })
  }
  map.setView([lat, lng], Math.max(map.getZoom(), 13))
}

const getMyLocation = () => {
  if (!navigator.geolocation) {
    message.error('Browser tidak mendukung geolokasi')
    return
  }
  gettingLocation.value = true
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      const lat = pos.coords.latitude
      const lng = pos.coords.longitude
      placeMarker(lat, lng)
      emit('update:latitude', lat)
      emit('update:longitude', lng)
      gettingLocation.value = false
      message.success('Lokasi berhasil diambil')
    },
    (err) => {
      gettingLocation.value = false
      message.error('Gagal mengambil lokasi: ' + err.message)
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

const handleSearch = async () => {
  if (!searchQuery.value.trim()) return
  searching.value = true
  searchResults.value = []
  try {
    const q = encodeURIComponent(searchQuery.value)
    const res = await fetch(
      `https://nominatim.openstreetmap.org/search?format=json&q=${q}&limit=5&countrycodes=id`
    )
    searchResults.value = await res.json()
  } catch (e) {
    message.error('Gagal mencari lokasi')
  } finally {
    searching.value = false
  }
}

const selectSearchResult = (result) => {
  const lat = parseFloat(result.lat)
  const lng = parseFloat(result.lon)
  placeMarker(lat, lng)
  emit('update:latitude', lat)
  emit('update:longitude', lng)
  searchResults.value = []
  searchQuery.value = result.display_name
}

watch(() => [props.latitude, props.longitude], ([lat, lng]) => {
  if (map && lat && lng) {
    placeMarker(lat, lng)
  }
})

onMounted(initMap)
onUnmounted(() => { if (map) { map.remove(); map = null } })
</script>

<style scoped>
.location-picker__search {
  position: relative;
}
.location-picker__results {
  position: absolute;
  z-index: 1000;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  max-height: 200px;
  overflow-y: auto;
  width: 100%;
  box-shadow: 0 2px 8px rgba(0,0,0,0.15);
}
.location-picker__result-item {
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: flex-start;
}
.location-picker__result-item:hover {
  background: #f4f7fe;
}
.location-picker__result-item:last-child {
  border-bottom: none;
}
.location-picker__map {
  height: 250px;
  width: 100%;
  border-radius: 8px;
  border: 1px solid #d9d9d9;
  z-index: 0;
}
.location-picker__controls {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 8px;
}
.location-picker__coords {
  font-size: 12px;
  color: #888;
  font-family: monospace;
}
</style>
