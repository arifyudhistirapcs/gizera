<template>
  <div class="attendance-config">
    <a-page-header
      title="Konfigurasi Absensi"
      sub-title="Kelola konfigurasi Wi-Fi dan GPS untuk validasi absensi karyawan"
    />

    <a-card>
      <a-tabs v-model:activeKey="activeTab" type="card">
        <!-- Wi-Fi Tab -->
        <a-tab-pane key="wifi" tab="Wi-Fi / IP">
          <template #tab>
            <span>
              <WifiOutlined />
              Wi-Fi / IP
            </span>
          </template>
          
          <WiFiConfigTab />
        </a-tab-pane>

        <!-- GPS Tab -->
        <a-tab-pane key="gps" tab="GPS / Lokasi">
          <template #tab>
            <span>
              <EnvironmentOutlined />
              GPS / Lokasi
            </span>
          </template>
          
          <GPSConfigTab />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { WifiOutlined, EnvironmentOutlined } from '@ant-design/icons-vue'
import WiFiConfigTab from '@/components/attendance/WiFiConfigTab.vue'
import GPSConfigTab from '@/components/attendance/GPSConfigTab.vue'

const activeTab = ref('wifi')
</script>

<style scoped>
.attendance-config {
  padding: 24px;
}

:deep(.ant-tabs-card > .ant-tabs-nav .ant-tabs-tab) {
  padding: 8px 16px;
}

:deep(.ant-tabs-card > .ant-tabs-nav .ant-tabs-tab-active) {
  background-color: #fff;
  border-bottom-color: #fff;
}
</style>
