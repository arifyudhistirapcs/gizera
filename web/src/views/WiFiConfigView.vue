<template>
  <div class="wifi-config">
    <a-page-header
      title="Konfigurasi Wi-Fi"
      sub-title="Kelola jaringan Wi-Fi yang diotorisasi untuk absensi karyawan"
    >
      <template #extra>
        <a-button type="primary" @click="showCreateModal">
          <template #icon><PlusOutlined /></template>
          Tambah Jaringan Wi-Fi
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Search and Filter -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari SSID atau lokasi..."
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </a-input>
          </a-col>
          <a-col :span="6">
            <a-select
              v-model:value="filterStatus"
              placeholder="Status"
              style="width: 100%"
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <a-select-option value="active">Aktif</a-select-option>
              <a-select-option value="inactive">Tidak Aktif</a-select-option>
            </a-select>
          </a-col>
        </a-row>

        <!-- Statistics Cards -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-card size="small">
              <a-statistic
                title="Total Jaringan"
                :value="stats.total || 0"
                :value-style="{ color: '#1890ff' }"
              />
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card size="small">
              <a-statistic
                title="Jaringan Aktif"
                :value="stats.active || 0"
                :value-style="{ color: '#52c41a' }"
              />
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card size="small">
              <a-statistic
                title="Tidak Aktif"
                :value="stats.inactive || 0"
                :value-style="{ color: '#ff4d4f' }"
              />
            </a-card>
          </a-col>
        </a-row>

        <!-- Table -->
        <a-table
          :columns="columns"
          :data-source="wifiConfigs"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'ssid'">
              <a-typography-text copyable>{{ record.ssid }}</a-typography-text>
            </template>
            <template v-else-if="column.key === 'ip_range'">
              <a-typography-text copyable code>{{ record.ip_range || '-' }}</a-typography-text>
            </template>
            <template v-else-if="column.key === 'is_active'">
              <a-tag :color="record.is_active ? 'green' : 'red'">
                {{ record.is_active ? 'Aktif' : 'Tidak Aktif' }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDateTime(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'actions'">
              <a-space>
                <a-button type="link" size="small" @click="editWiFiConfig(record)">
                  Edit
                </a-button>
                <a-popconfirm
                  :title="record.is_active ? 'Yakin ingin menonaktifkan jaringan Wi-Fi ini?' : 'Yakin ingin mengaktifkan jaringan Wi-Fi ini?'"
                  ok-text="Ya"
                  cancel-text="Tidak"
                  @confirm="toggleWiFiConfigStatus(record)"
                >
                  <a-button type="link" size="small" :danger="record.is_active">
                    {{ record.is_active ? 'Nonaktifkan' : 'Aktifkan' }}
                  </a-button>
                </a-popconfirm>
                <a-popconfirm
                  title="Yakin ingin menghapus jaringan Wi-Fi ini?"
                  ok-text="Ya"
                  cancel-text="Tidak"
                  @confirm="deleteWiFiConfig(record)"
                >
                  <a-button type="link" size="small" danger>
                    Hapus
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingConfig ? 'Edit Jaringan Wi-Fi' : 'Tambah Jaringan Wi-Fi'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="600px"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="SSID (Nama Jaringan)" name="ssid">
          <a-input 
            v-model:value="formData.ssid" 
            placeholder="Masukkan nama jaringan Wi-Fi (SSID)"
            :maxlength="100"
          />
          <div style="font-size: 12px; color: #666; margin-top: 4px;">
            Contoh: SPPG-Office, Kantor-WiFi
          </div>
        </a-form-item>

        <a-form-item label="IP Range (CIDR Notation)" name="ip_range">
          <a-input 
            v-model:value="formData.ip_range" 
            placeholder="Masukkan IP Range dalam format CIDR"
            :maxlength="18"
          />
          <div style="font-size: 12px; color: #666; margin-top: 4px;">
            Format: 192.168.1.0/24 (untuk jaringan 192.168.1.0 - 192.168.1.255)
          </div>
        </a-form-item>

        <a-form-item label="IP Spesifik (Opsional)" name="allowed_ips">
          <a-textarea 
            v-model:value="formData.allowed_ips" 
            placeholder="Masukkan IP address spesifik yang diizinkan (satu per baris)"
            :rows="3"
          />
          <div style="font-size: 12px; color: #666; margin-top: 4px;">
            Contoh:<br/>
            192.168.1.100<br/>
            192.168.1.101<br/>
            Kosongkan jika menggunakan IP Range saja
          </div>
        </a-form-item>

        <a-form-item label="Lokasi" name="location">
          <a-input 
            v-model:value="formData.location" 
            placeholder="Masukkan lokasi jaringan Wi-Fi"
            :maxlength="200"
          />
          <div style="font-size: 12px; color: #666; margin-top: 4px;">
            Contoh: Kantor Pusat SPPG, Ruang Meeting Lt. 2
          </div>
        </a-form-item>

        <a-form-item label="Status" name="is_active">
          <a-switch 
            v-model:checked="formData.is_active" 
            checked-children="Aktif" 
            un-checked-children="Tidak Aktif" 
          />
          <div style="font-size: 12px; color: #666; margin-top: 4px;">
            Hanya jaringan aktif yang dapat digunakan untuk absensi
          </div>
        </a-form-item>
      </a-form>

      <a-alert
        message="Informasi Penting"
        description="Sistem akan memvalidasi check-in berdasarkan IP address karyawan. Pastikan IP Range yang dimasukkan sesuai dengan jaringan kantor Anda."
        type="info"
        show-icon
        style="margin-top: 16px"
      />
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import wifiConfigService from '@/services/wifiConfigService'

const loading = ref(false)
const submitting = ref(false)
const modalVisible = ref(false)
const editingConfig = ref(null)
const wifiConfigs = ref([])
const searchText = ref('')
const filterStatus = ref(undefined)
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  ssid: '',
  ip_range: '',
  allowed_ips: '',
  location: '',
  is_active: true
})

const rules = {
  ssid: [
    { required: true, message: 'SSID wajib diisi' },
    { min: 1, max: 100, message: 'SSID harus antara 1-100 karakter' }
  ],
  ip_range: [
    { required: true, message: 'IP Range wajib diisi' },
    { 
      pattern: /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/, 
      message: 'Format IP Range tidak valid (contoh: 192.168.1.0/24)' 
    }
  ],
  location: [
    { required: true, message: 'Lokasi wajib diisi' },
    { max: 200, message: 'Lokasi maksimal 200 karakter' }
  ]
}

const columns = [
  {
    title: 'SSID',
    key: 'ssid',
    dataIndex: 'ssid',
    sorter: true,
    width: 200
  },
  {
    title: 'IP Range',
    key: 'ip_range',
    dataIndex: 'ip_range',
    width: 180
  },
  {
    title: 'Lokasi',
    dataIndex: 'location',
    key: 'location'
  },
  {
    title: 'Status',
    key: 'is_active',
    width: 100
  },
  {
    title: 'Dibuat Pada',
    key: 'created_at',
    width: 160
  },
  {
    title: 'Aksi',
    key: 'actions',
    width: 200
  }
]

const stats = computed(() => {
  const total = wifiConfigs.value.length
  const active = wifiConfigs.value.filter(config => config.is_active).length
  const inactive = total - active
  
  return { total, active, inactive }
})

const fetchWiFiConfigs = async () => {
  loading.value = true
  try {
    const response = await wifiConfigService.getWiFiConfigs()
    wifiConfigs.value = response.data || []
    pagination.total = wifiConfigs.value.length
    
    // Apply client-side filtering
    applyFilters()
  } catch (error) {
    message.error('Gagal memuat konfigurasi Wi-Fi')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  let filtered = [...wifiConfigs.value]
  
  // Apply search filter
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    filtered = filtered.filter(config => 
      config.ssid.toLowerCase().includes(search) ||
      config.location.toLowerCase().includes(search)
    )
  }
  
  // Apply status filter
  if (filterStatus.value === 'active') {
    filtered = filtered.filter(config => config.is_active)
  } else if (filterStatus.value === 'inactive') {
    filtered = filtered.filter(config => !config.is_active)
  }
  
  // Update pagination
  pagination.total = filtered.length
  
  // Apply pagination
  const start = (pagination.current - 1) * pagination.pageSize
  const end = start + pagination.pageSize
  wifiConfigs.value = filtered.slice(start, end)
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchWiFiConfigs()
}

const handleSearch = () => {
  pagination.current = 1
  fetchWiFiConfigs()
}

const showCreateModal = () => {
  editingConfig.value = null
  resetForm()
  modalVisible.value = true
}

const editWiFiConfig = (config) => {
  editingConfig.value = config
  
  // Convert allowed_ips array to string (one per line)
  const allowedIpsStr = Array.isArray(config.allowed_ips) 
    ? config.allowed_ips.join('\n') 
    : ''
  
  Object.assign(formData, {
    ssid: config.ssid,
    ip_range: config.ip_range || '',
    allowed_ips: allowedIpsStr,
    location: config.location,
    is_active: config.is_active
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    // Convert allowed_ips string to array
    const allowedIpsArray = formData.allowed_ips
      ? formData.allowed_ips.split('\n').map(ip => ip.trim()).filter(ip => ip)
      : []

    const submitData = {
      ssid: formData.ssid,
      ip_range: formData.ip_range,
      allowed_ips: allowedIpsArray,
      location: formData.location,
      is_active: formData.is_active,
      // Keep bssid for backward compatibility, use dummy value
      bssid: '00:00:00:00:00:00'
    }

    if (editingConfig.value) {
      await wifiConfigService.updateWiFiConfig(editingConfig.value.id, submitData)
      message.success('Konfigurasi Wi-Fi berhasil diperbarui')
    } else {
      await wifiConfigService.createWiFiConfig(submitData)
      message.success('Konfigurasi Wi-Fi berhasil ditambahkan')
    }

    modalVisible.value = false
    fetchWiFiConfigs()
  } catch (error) {
    if (error.errorFields) {
      return
    }
    
    const errorMessage = error.response?.data?.message || 'Gagal menyimpan konfigurasi Wi-Fi'
    message.error(errorMessage)
    console.error(error)
  } finally {
    submitting.value = false
  }
}

const toggleWiFiConfigStatus = async (config) => {
  try {
    await wifiConfigService.updateWiFiConfig(config.id, { 
      is_active: !config.is_active 
    })
    message.success(`Jaringan Wi-Fi berhasil ${config.is_active ? 'dinonaktifkan' : 'diaktifkan'}`)
    fetchWiFiConfigs()
  } catch (error) {
    message.error('Gagal mengubah status jaringan Wi-Fi')
    console.error(error)
  }
}

const deleteWiFiConfig = async (config) => {
  try {
    await wifiConfigService.deleteWiFiConfig(config.id)
    message.success('Jaringan Wi-Fi berhasil dihapus')
    fetchWiFiConfigs()
  } catch (error) {
    message.error('Gagal menghapus jaringan Wi-Fi')
    console.error(error)
  }
}

const handleCancel = () => {
  modalVisible.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    ssid: '',
    ip_range: '',
    allowed_ips: '',
    location: '',
    is_active: true
  })
  formRef.value?.resetFields()
}

const formatDateTime = (date) => {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY HH:mm')
}

onMounted(() => {
  fetchWiFiConfigs()
})
</script>

<style scoped>
.wifi-config {
  padding: 24px;
}
</style>