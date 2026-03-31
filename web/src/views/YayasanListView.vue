<template>
  <div class="yayasan-list">
    <a-page-header title="Manajemen Yayasan" sub-title="Kelola data Yayasan">
      <template #extra>
        <a-button type="primary" @click="openModal()">
          <template #icon><PlusOutlined /></template>
          Tambah Yayasan
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <a-input-search
          v-model:value="searchText"
          placeholder="Cari nama atau kode Yayasan..."
          style="max-width: 400px"
          allow-clear
          @search="fetchData"
        />

        <a-table
          :columns="columns"
          :data-source="filteredData"
          :loading="loading"
          :pagination="{ pageSize: 10, showSizeChanger: true }"
          row-key="id"
          size="middle"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'kode'">
              <a-tag color="blue">{{ record.kode }}</a-tag>
            </template>
            <template v-if="column.key === 'sppg_count'">
              {{ record.sppgs ? record.sppgs.length : 0 }}
            </template>
            <template v-if="column.key === 'is_active'">
              <a-tag :color="record.is_active ? 'green' : 'red'">
                {{ record.is_active ? 'Aktif' : 'Nonaktif' }}
              </a-tag>
            </template>
            <template v-if="column.key === 'action'">
              <a-space>
                <a-button size="small" @click="openModal(record)">Edit</a-button>
                <a-popconfirm
                  :title="`${record.is_active ? 'Nonaktifkan' : 'Aktifkan'} Yayasan ini?`"
                  @confirm="toggleStatus(record)"
                >
                  <a-button size="small" :danger="record.is_active">
                    {{ record.is_active ? 'Nonaktifkan' : 'Aktifkan' }}
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
      :title="editingRecord ? 'Edit Yayasan' : 'Tambah Yayasan'"
      :confirm-loading="saving"
      @ok="handleSave"
      @cancel="closeModal"
      width="650px"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="Nama" required>
          <a-input v-model:value="form.nama" placeholder="Nama Yayasan" />
        </a-form-item>
        <a-form-item label="Alamat">
          <a-textarea v-model:value="form.alamat" placeholder="Alamat" :rows="2" />
        </a-form-item>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Nomor Telepon">
              <a-input v-model:value="form.nomor_telepon" placeholder="Nomor Telepon" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Email">
              <a-input v-model:value="form.email" placeholder="Email" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Penanggung Jawab">
              <a-input v-model:value="form.penanggung_jawab" placeholder="Nama Penanggung Jawab" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="NPWP">
              <a-input v-model:value="form.npwp" placeholder="NPWP" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="Lokasi">
          <LocationPicker
            v-model:latitude="form.latitude"
            v-model:longitude="form.longitude"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import organizationService from '@/services/organizationService'
import LocationPicker from '@/components/LocationPicker.vue'

const loading = ref(false)
const saving = ref(false)
const dataSource = ref([])
const searchText = ref('')
const modalVisible = ref(false)
const editingRecord = ref(null)

const form = ref({
  nama: '',
  alamat: '',
  nomor_telepon: '',
  email: '',
  penanggung_jawab: '',
  npwp: '',
  latitude: 0,
  longitude: 0
})

const columns = [
  { title: 'Kode', dataIndex: 'kode', key: 'kode', width: 120 },
  { title: 'Nama', dataIndex: 'nama', key: 'nama', sorter: (a, b) => a.nama.localeCompare(b.nama) },
  { title: 'Penanggung Jawab', dataIndex: 'penanggung_jawab', key: 'penanggung_jawab' },
  { title: 'Jumlah SPPG', key: 'sppg_count', width: 120, align: 'center' },
  { title: 'Status', key: 'is_active', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 180, align: 'center' }
]

const filteredData = computed(() => {
  if (!searchText.value) return dataSource.value
  const q = searchText.value.toLowerCase()
  return dataSource.value.filter(
    r => r.nama?.toLowerCase().includes(q) || r.kode?.toLowerCase().includes(q)
  )
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await organizationService.getYayasanList()
    dataSource.value = res.data?.data || res.data?.yayasans || []
  } catch (e) {
    message.error('Gagal memuat data Yayasan')
  } finally {
    loading.value = false
  }
}

const openModal = (record = null) => {
  editingRecord.value = record
  if (record) {
    form.value = { nama: record.nama, alamat: record.alamat, nomor_telepon: record.nomor_telepon, email: record.email, penanggung_jawab: record.penanggung_jawab, npwp: record.npwp, latitude: record.latitude || 0, longitude: record.longitude || 0 }
  } else {
    form.value = { nama: '', alamat: '', nomor_telepon: '', email: '', penanggung_jawab: '', npwp: '', latitude: 0, longitude: 0 }
  }
  modalVisible.value = true
}

const closeModal = () => {
  modalVisible.value = false
  editingRecord.value = null
}

const handleSave = async () => {
  if (!form.value.nama) {
    message.warning('Nama Yayasan wajib diisi')
    return
  }
  saving.value = true
  try {
    if (editingRecord.value) {
      await organizationService.updateYayasan(editingRecord.value.id, form.value)
      message.success('Yayasan berhasil diperbarui')
    } else {
      await organizationService.createYayasan(form.value)
      message.success('Yayasan berhasil ditambahkan')
    }
    closeModal()
    fetchData()
  } catch (e) {
    message.error(e.response?.data?.message || 'Gagal menyimpan Yayasan')
  } finally {
    saving.value = false
  }
}

const toggleStatus = async (record) => {
  try {
    await organizationService.setYayasanStatus(record.id, !record.is_active)
    message.success(`Yayasan berhasil ${record.is_active ? 'dinonaktifkan' : 'diaktifkan'}`)
    fetchData()
  } catch (e) {
    message.error('Gagal mengubah status Yayasan')
  }
}

onMounted(fetchData)
</script>
