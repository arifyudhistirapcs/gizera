<template>
  <div ref="mapContainer" :style="{ height: height + 'px', width: '100%' }"></div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

const props = defineProps({
  markers: {
    type: Array,
    default: () => []
  },
  height: {
    type: Number,
    default: 400
  }
})

const mapContainer = ref(null)
let map = null
let markerLayer = null

const COLORS = {
  Yayasan: '#722ed1',
  SPPG: '#52c41a',
  Sekolah: '#1890ff',
  Supplier: '#fa8c16'
}

const DEFAULT_CENTER = [-2.5, 118]
const DEFAULT_ZOOM = 5

const initMap = () => {
  if (!mapContainer.value) return
  map = L.map(mapContainer.value).setView(DEFAULT_CENTER, DEFAULT_ZOOM)
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    maxZoom: 19
  }).addTo(map)
  markerLayer = L.layerGroup().addTo(map)
  renderMarkers()
}

const renderMarkers = () => {
  if (!map || !markerLayer) return
  markerLayer.clearLayers()

  const validMarkers = props.markers.filter(m => m.latitude && m.longitude)
  if (validMarkers.length === 0) {
    map.setView(DEFAULT_CENTER, DEFAULT_ZOOM)
    return
  }

  const bounds = []
  validMarkers.forEach(m => {
    const color = COLORS[m.type] || '#1890ff'
    const latlng = [m.latitude, m.longitude]
    bounds.push(latlng)

    const circle = L.circleMarker(latlng, {
      radius: 8,
      fillColor: color,
      color: '#fff',
      weight: 2,
      opacity: 1,
      fillOpacity: 0.85
    })

    const badgeColor = COLORS[m.type] || '#1890ff'
    const statsHtml = m.stats ? (() => {
      if (m.type === 'Yayasan') {
        return `
          <table style="margin-top:8px;font-size:12px;width:100%;border-collapse:collapse">
            <tr><td style="color:#888;padding:2px 8px 2px 0">Jumlah SPPG</td><td style="text-align:right;font-weight:600">${m.stats.sppgCount ?? '-'}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Total Porsi</td><td style="text-align:right;font-weight:600">${(m.stats.totalPortions ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Pengeluaran</td><td style="text-align:right;font-weight:600">Rp ${(m.stats.totalSpent ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Rating</td><td style="text-align:right;font-weight:600">${(m.stats.rating ?? 0).toFixed(1)} / 5</td></tr>
          </table>`
      } else if (m.type === 'Sekolah') {
        return `
          <table style="margin-top:8px;font-size:12px;width:100%;border-collapse:collapse">
            <tr><td style="color:#888;padding:2px 8px 2px 0">Porsi Kecil</td><td style="text-align:right;font-weight:600">${(m.stats.portionsSmall ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Porsi Besar</td><td style="text-align:right;font-weight:600">${(m.stats.portionsLarge ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Total Distribusi</td><td style="text-align:right;font-weight:600">${(m.stats.totalDelivered ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Rating Feedback</td><td style="text-align:right;font-weight:600">${(m.stats.rating ?? 0).toFixed(1)} / 5</td></tr>
          </table>`
      } else {
        return `
          <table style="margin-top:8px;font-size:12px;width:100%;border-collapse:collapse">
            <tr><td style="color:#888;padding:2px 8px 2px 0">Total Porsi</td><td style="text-align:right;font-weight:600">${(m.stats.totalPortions ?? 0).toLocaleString('id-ID')}</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Delivery Rate</td><td style="text-align:right;font-weight:600">${(m.stats.deliveryRate ?? 0).toFixed(1)}%</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Penyerapan</td><td style="text-align:right;font-weight:600">${(m.stats.budgetAbsorption ?? 0).toFixed(1)}%</td></tr>
            <tr><td style="color:#888;padding:2px 8px 2px 0">Rating</td><td style="text-align:right;font-weight:600">${(m.stats.rating ?? 0).toFixed(1)} / 5</td></tr>
          </table>`
      }
    })() : ''
    const popup = `
      <div style="min-width:200px">
        <strong style="font-size:14px">${m.name}</strong><br/>
        <span style="background:${badgeColor};color:#fff;padding:1px 8px;border-radius:10px;font-size:11px">${m.type}</span>
        <span style="margin-left:6px;color:#888;font-size:12px">${m.kode}</span>
        ${m.details ? `<div style="margin-top:6px;font-size:12px;color:#555">${m.details}</div>` : ''}
        ${statsHtml}
      </div>
    `
    circle.bindPopup(popup)
    circle.addTo(markerLayer)
  })

  if (bounds.length > 0) {
    map.fitBounds(bounds, { padding: [30, 30], maxZoom: 12 })
  }
}

watch(() => props.markers, renderMarkers, { deep: true })

onMounted(initMap)

onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>
