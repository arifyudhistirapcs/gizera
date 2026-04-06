<template>
  <div class="gps-config-tab">
    <a-space direction="vertical" style="width: 100%" :size="16">
      <!-- Header with Add Button -->
      <a-row justify="space-between" align="middle">
        <a-col>
          <a-space>
            <a-input
              v-model:value="searchText"
              placeholder="Cari nama lokasi..."
              allow-clear
              style="width: 250px"
            >
              <template #prefix>
                <SearchOutlined />
              </template>
            </a-input>
            <a-select
              v-model:value="filterStatus"
              placeholder="Status"
              style="width: 150px"
              allow-clear
            >
              <a-select-option value="active">Aktif</a-select-option>
              <a-select-option value="inactive">Tidak Aktif</a-select-option>
            </a-select>
          </a-space>
        </a-col>
        <a-col>
          <a-button type="primary" @click="showCreateModal">
            <template #icon><PlusOutlined /></template>
            Tambah Lokasi GPS
          </a-button>
        </a-col>
      </a-row>

      <!-- Statistics Cards -->
      <a-row :gutter="16">
        <a-col :span="8">
          <a-card size="small">
            <a-statistic
              title="Total Lokasi"
              :value="stats.total || 0"
              :value-style="{ color: '#1890ff' }"
            />
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card size="small">
            <a-statistic
              title="Lokasi Aktif"
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

      <!-- Info Alert -->
      <a-alert
        message="Informasi GPS"
        description="Karyawan harus berada dalam radius yang ditentukan dari koordinat lokasi untuk dapat melakukan check-in. Pastikan koordinat dan radius sudah sesuai dengan area kantor."
        type="info"
        show-icon
      />

      <!-- Table -->
      <a-table
        :columns="columns"
        :data-source="filteredConfigs"
        :loading="loading"
        :pagination="pagination"
        @change="handleTableChange"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-space>
              <EnvironmentOutlined style="color: #303030" />
              <span style="font-weight: 500;">{{ record.name }}</span>
            </a-space>
          </template>
          <template v-else-if="column.key === 'coordinates'">
            <a-typography-text copyable :content="`${record.latitude}, ${record.longitude}`">
              {{ record.latitude.toFixed(6) }}, {{ record.longitude.toFixed(6) }}
            </a-typography-text>
          </template>
          <template v-else-if="column.key === 'radius'">
            <a-tag color="blue">{{ record.radius }} meter</a-tag>
          </template>
          <template v-else-if="column.key === 'is_active'">
            <a-tag :color="record.is_active ? 'green' : 'red'">
              {{ record.is_active ? 'Aktif' : 'Tidak Aktif' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" size="small" @click="viewOnMap(record)">
                <template #icon><AimOutlined /></template>
                Lihat
              </a-button>
              <a-button type="link" size="small" @click="editConfig(record)">
                Edit
              </a-button>
              <a-popconfirm
                :title="record.is_active ? 'Nonaktifkan lokasi ini?' : 'Aktifkan lokasi ini?'"
                @confirm="toggleStatus(record)"
              >
                <a-button type="link" size="small" :danger="record.is_active">
                  {{ record.is_active ? 'Nonaktifkan' : 'Aktifkan' }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm
                title="Hapus lokasi GPS ini?"
                @confirm="deleteConfig(record)"
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

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingConfig ? 'Edit Lokasi GPS' : 'Tambah Lokasi GPS'"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="700px"
    >
      <a-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        layout="vertical"
      >
        <a-form-item label="Nama Lokasi" name="name">
          <a-input 
            v-model:value="formData.name" 
            placeholder="Contoh: Kantor Pusat SPPG"
          />
        </a-form-item>

        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Latitude" name="latitude">
              <a-input-number 
                v-model:value="formData.latitude" 
                placeholder="-6.123456"
                :precision="6"
                :step="0.000001"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Longitude" name="longitude">
              <a-input-number 
                v-model:value="formData.longitude" 
                placeholder="106.123456"
                :precision="6"
                :step="0.000001"
                style="width: 100%"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item>
          <a-button type="dashed" block @click="getCurrentLocation" :loading="gettingLocation">
            <template #icon><AimOutlined /></template>
            Gunakan Lokasi Saat Ini
          </a-button>
        </a-form-item>

        <a-form-item label="Radius (meter)" name="radius">
          <a-slider 
            v-model:value="formData.radius" 
            :min="10" 
            :max="500" 
            :marks="radiusMarks"
          />
          <div class="radius-display">
            <a-tag color="blue">{{ formData.radius }} meter</a-tag>
            <span class="radius-hint">Area yang diizinkan untuk check-in</span>
          </div>
        </a-form-item>

        <a-form-item label="Alamat (Opsional)" name="address">
          <a-textarea 
            v-model:value="formData.address" 
            placeholder="Alamat lengkap lokasi"
            :rows="2"
          />
        </a-form-item>

        <a-form-item label="Deskripsi (Opsional)" name="description">
          <a-textarea 
            v-model:value="formData.description" 
            placeholder="Keterangan tambahan"
            :rows="2"
          />
        </a-form-item>

        <a-form-item label="Status" name="is_active">
          <a-switch 
            v-model:checked="formData.is_active" 
            checked-children="Aktif" 
            un-checked-children="Tidak Aktif" 
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Map Modal -->
    <a-modal
      v-model:open="mapModalVisible"
      :title="`Lokasi: ${selectedConfig?.name || ''}`"
      :footer="null"
      width="800px"
    >
      <div v-if="selectedConfig" class="map-info">
        <a-descriptions :column="2" size="small" bordered>
          <a-descriptions-item label="Koordinat">
            {{ selectedConfig.latitude.toFixed(6) }}, {{ selectedConfig.longitude.toFixed(6) }}
          </a-descriptions-item>
          <a-descriptions-item label="Radius">
            {{ selectedConfig.radius }} meter
          </a-descriptions-item>
          <a-descriptions-item label="Alamat" :span="2">
            {{ selectedConfig.address || '-' }}
          </a-descriptions-item>
        </a-descriptions>
        
        <div class="map-actions">
          <a-button 
            type="primary" 
            @click="openInGoogleMaps(selectedConfig)"
          >
            <template #icon><GlobalOutlined /></template>
            Buka di Google Maps
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { 
  PlusOutlined, 
  SearchOutlined, 
  EnvironmentOutlined, 
  AimOutlined,
  GlobalOutlined
} from '@ant-design/icons-vue'
import gpsConfigService from '@/services/gpsConfigService'

const loading = ref(false)
const submitting = ref(false)
const gettingLocation = ref(false)
const modalVisible = ref(false)
const mapModalVisible = ref(false)
const editingConfig = ref(null)
const selectedConfig = ref(null)
const gpsConfigs = ref([])
const searchText = ref('')
const filterStatus = ref(undefined)
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  name: '',
  latitude: null,
  longitude: null,
  radius: 100,
  address: '',
  description: '',
  is_active: true
})

const radiusMarks = {
  10: '10m',
  50: '50m',
  100: '100m',
  200: '200m',
  300: '300m',
  500: '500m'
}

const rules = {
  name: [{ required: true, message: 'Nama lokasi wajib diisi' }],
  latitude: [
    { required: true, message: 'Latitude wajib diisi' },
    { type: 'number', min: -90, max: 90, message: 'Latitude harus antara -90 dan 90' }
  ],
  longitude: [
    { required: true, message: 'Longitude wajib diisi' },
    { type: 'number', min: -180, max: 180, message: 'Longitude harus antara -180 dan 180' }
  ],
  radius: [
    { required: true, message: 'Radius wajib diisi' },
    { type: 'number', min: 10, max: 500, message: 'Radius harus antara 10 dan 500 meter' }
  ]
}

const columns = [
  { title: 'Nama Lokasi', key: 'name', width: 200 },
  { title: 'Koordinat', key: 'coordinates', width: 220 },
  { title: 'Radius', key: 'radius', width: 100 },
  { title: 'Status', key: 'is_active', width: 100 },
  { title: 'Aksi', key: 'actions', width: 250 }
]

const stats = computed(() => {
  const total = gpsConfigs.value.length
  const active = gpsConfigs.value.filter(c => c.is_active).length
  return { total, active, inactive: total - active }
})

const filteredConfigs = computed(() => {
  let result = [...gpsConfigs.value]
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(c => 
      c.name.toLowerCase().includes(search) ||
      c.address?.toLowerCase().includes(search)
    )
  }
  
  if (filterStatus.value === 'active') {
    result = result.filter(c => c.is_active)
  } else if (filterStatus.value === 'inactive') {
    result = result.filter(c => !c.is_active)
  }
  
  return result
})

const fetchConfigs = async () => {
  loading.value = true
  try {
    const response = await gpsConfigService.getGPSConfigs()
    gpsConfigs.value = response.data || []
  } catch (error) {
    message.error('Gagal memuat data')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const showCreateModal = () => {
  editingConfig.value = null
  resetForm()
  modalVisible.value = true
}

const editConfig = (config) => {
  editingConfig.value = config
  Object.assign(formData, {
    name: config.name,
    latitude: config.latitude,
    longitude: config.longitude,
    radius: config.radius,
    address: config.address || '',
    description: config.description || '',
    is_active: config.is_active
  })
  modalVisible.value = true
}

const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    message.error('Browser tidak mendukung geolocation')
    return
  }

  gettingLocation.value = true
  navigator.geolocation.getCurrentPosition(
    (position) => {
      formData.latitude = position.coords.latitude
      formData.longitude = position.coords.longitude
      message.success('Lokasi berhasil didapatkan')
      gettingLocation.value = false
    },
    (error) => {
      let errorMsg = 'Gagal mendapatkan lokasi'
      if (error.code === 1) errorMsg = 'Akses lokasi ditolak'
      else if (error.code === 2) errorMsg = 'Lokasi tidak tersedia'
      else if (error.code === 3) errorMsg = 'Timeout mendapatkan lokasi'
      message.error(errorMsg)
      gettingLocation.value = false
    },
    { enableHighAccuracy: true, timeout: 10000 }
  )
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    const submitData = {
      name: formData.name,
      latitude: formData.latitude,
      longitude: formData.longitude,
      radius: formData.radius,
      address: formData.address,
      description: formData.description,
      is_active: formData.is_active
    }

    if (editingConfig.value) {
      await gpsConfigService.updateGPSConfig(editingConfig.value.id, submitData)
      message.success('Lokasi berhasil diperbarui')
    } else {
      await gpsConfigService.createGPSConfig(submitData)
      message.success('Lokasi berhasil ditambahkan')
    }

    modalVisible.value = false
    fetchConfigs()
  } catch (error) {
    if (!error.errorFields) {
      message.error(error.response?.data?.message || 'Gagal menyimpan')
    }
  } finally {
    submitting.value = false
  }
}

const toggleStatus = async (config) => {
  try {
    await gpsConfigService.updateGPSConfig(config.id, { is_active: !config.is_active })
    message.success('Status berhasil diubah')
    fetchConfigs()
  } catch (error) {
    message.error('Gagal mengubah status')
  }
}

const deleteConfig = async (config) => {
  try {
    await gpsConfigService.deleteGPSConfig(config.id)
    message.success('Lokasi berhasil dihapus')
    fetchConfigs()
  } catch (error) {
    message.error('Gagal menghapus')
  }
}

const viewOnMap = (config) => {
  selectedConfig.value = config
  mapModalVisible.value = true
}

const openInGoogleMaps = (config) => {
  const url = `https://www.google.com/maps?q=${config.latitude},${config.longitude}`
  window.open(url, '_blank')
}

const handleCancel = () => {
  modalVisible.value = false
  resetForm()
}

const resetForm = () => {
  Object.assign(formData, {
    name: '',
    latitude: null,
    longitude: null,
    radius: 100,
    address: '',
    description: '',
    is_active: true
  })
  formRef.value?.resetFields()
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
}

onMounted(fetchConfigs)
</script>

<style scoped>
.gps-config-tab {
  padding: 16px 0;
}

.radius-display {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
}

.radius-hint {
  font-size: 12px;
  color: #666;
}

.map-info {
  padding: 16px 0;
}

.map-actions {
  margin-top: 16px;
  text-align: center;
}
</style>
