<template>
  <div class="wifi-config-tab">
    <a-space direction="vertical" style="width: 100%" :size="16">
      <!-- Header with Add Button -->
      <a-row justify="space-between" align="middle">
        <a-col>
          <a-space>
            <a-input
              v-model:value="searchText"
              placeholder="Cari SSID atau lokasi..."
              @change="handleSearch"
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
              @change="handleSearch"
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
            Tambah Jaringan Wi-Fi
          </a-button>
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

      <!-- Info Alert -->
      <a-alert
        message="Informasi Wi-Fi / IP"
        description="Karyawan harus terhubung ke jaringan Wi-Fi kantor atau menggunakan IP address yang terdaftar untuk dapat melakukan check-in. Pastikan IP Range sudah sesuai dengan jaringan kantor."
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
              <a-button type="link" size="small" @click="editConfig(record)">
                Edit
              </a-button>
              <a-popconfirm
                :title="record.is_active ? 'Nonaktifkan jaringan ini?' : 'Aktifkan jaringan ini?'"
                @confirm="toggleStatus(record)"
              >
                <a-button type="link" size="small" :danger="record.is_active">
                  {{ record.is_active ? 'Nonaktifkan' : 'Aktifkan' }}
                </a-button>
              </a-popconfirm>
              <a-popconfirm
                title="Hapus jaringan Wi-Fi ini?"
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
            placeholder="Masukkan nama jaringan Wi-Fi"
          />
        </a-form-item>

        <a-form-item label="IP Range (CIDR)" name="ip_range">
          <a-input 
            v-model:value="formData.ip_range" 
            placeholder="Contoh: 192.168.1.0/24"
          />
          <div class="form-hint">
            Format CIDR untuk range IP yang diizinkan
          </div>
        </a-form-item>

        <a-form-item label="IP Spesifik (Opsional)" name="allowed_ips">
          <a-textarea 
            v-model:value="formData.allowed_ips" 
            placeholder="Satu IP per baris"
            :rows="3"
          />
        </a-form-item>

        <a-form-item label="Lokasi" name="location">
          <a-input 
            v-model:value="formData.location" 
            placeholder="Contoh: Kantor Pusat"
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, SearchOutlined } from '@ant-design/icons-vue'
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
  ssid: [{ required: true, message: 'SSID wajib diisi' }],
  ip_range: [
    { required: true, message: 'IP Range wajib diisi' },
    { pattern: /^(\d{1,3}\.){3}\d{1,3}\/\d{1,2}$/, message: 'Format tidak valid' }
  ],
  location: [{ required: true, message: 'Lokasi wajib diisi' }]
}

const columns = [
  { title: 'SSID', key: 'ssid', dataIndex: 'ssid', width: 180 },
  { title: 'IP Range', key: 'ip_range', dataIndex: 'ip_range', width: 160 },
  { title: 'Lokasi', dataIndex: 'location', key: 'location' },
  { title: 'Status', key: 'is_active', width: 100 },
  { title: 'Dibuat', key: 'created_at', width: 150 },
  { title: 'Aksi', key: 'actions', width: 200 }
]

const stats = computed(() => {
  const total = wifiConfigs.value.length
  const active = wifiConfigs.value.filter(c => c.is_active).length
  return { total, active, inactive: total - active }
})

const filteredConfigs = computed(() => {
  let result = [...wifiConfigs.value]
  
  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter(c => 
      c.ssid.toLowerCase().includes(search) ||
      c.location?.toLowerCase().includes(search)
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
    const response = await wifiConfigService.getWiFiConfigs()
    wifiConfigs.value = response.data || []
  } catch (error) {
    message.error('Gagal memuat data')
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
    ssid: config.ssid,
    ip_range: config.ip_range || '',
    allowed_ips: Array.isArray(config.allowed_ips) ? config.allowed_ips.join('\n') : '',
    location: config.location,
    is_active: config.is_active
  })
  modalVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    const submitData = {
      ssid: formData.ssid,
      ip_range: formData.ip_range,
      allowed_ips: formData.allowed_ips ? formData.allowed_ips.split('\n').filter(ip => ip.trim()) : [],
      location: formData.location,
      is_active: formData.is_active,
      bssid: '00:00:00:00:00:00'
    }

    if (editingConfig.value) {
      await wifiConfigService.updateWiFiConfig(editingConfig.value.id, submitData)
      message.success('Berhasil diperbarui')
    } else {
      await wifiConfigService.createWiFiConfig(submitData)
      message.success('Berhasil ditambahkan')
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
    await wifiConfigService.updateWiFiConfig(config.id, { is_active: !config.is_active })
    message.success('Status berhasil diubah')
    fetchConfigs()
  } catch (error) {
    message.error('Gagal mengubah status')
  }
}

const deleteConfig = async (config) => {
  try {
    await wifiConfigService.deleteWiFiConfig(config.id)
    message.success('Berhasil dihapus')
    fetchConfigs()
  } catch (error) {
    message.error('Gagal menghapus')
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

const handleSearch = () => {}
const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
}

const formatDateTime = (date) => date ? dayjs(date).format('DD/MM/YYYY HH:mm') : '-'

onMounted(fetchConfigs)
</script>

<style scoped>
.wifi-config-tab {
  padding: 16px 0;
}
.form-hint {
  font-size: 12px;
  color: #666;
  margin-top: 4px;
}
</style>
