<template>
  <div class="sppg-list">
    <a-page-header title="Manajemen SPPG" sub-title="Kelola data Satuan Pelayanan Pemenuhan Gizi">
      <template #extra>
        <a-button type="primary" @click="openModal()">
          <template #icon><PlusOutlined /></template>
          Tambah SPPG
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <a-input-search
          v-model:value="searchText"
          placeholder="Cari nama atau kode SPPG..."
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
            <template v-if="column.key === 'yayasan'">
              {{ record.yayasan?.nama || '-' }}
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
                  :title="`${record.is_active ? 'Nonaktifkan' : 'Aktifkan'} SPPG ini?`"
                  @confirm="toggleStatus(record)"
                >
                  <a-button size="small" :danger="record.is_active">
                    {{ record.is_active ? 'Nonaktifkan' : 'Aktifkan' }}
                  </a-button>
                </a-popconfirm>
                <a-button size="small" type="dashed" @click="openTransferModal(record)">
                  Transfer
                </a-button>
              </a-space>
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:open="modalVisible"
      :title="editingRecord ? 'Edit SPPG' : 'Tambah SPPG'"
      :confirm-loading="saving"
      @ok="handleSave"
      @cancel="closeModal"
      width="650px"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item label="Nama" required>
          <a-input v-model:value="form.nama" placeholder="Nama SPPG" />
        </a-form-item>
        <a-form-item label="Yayasan" required>
          <a-select
            v-model:value="form.yayasan_id"
            placeholder="Pilih Yayasan"
            :options="yayasanOptions"
            show-search
            :filter-option="filterOption"
          />
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
        <a-form-item label="Lokasi">
          <LocationPicker
            v-model:latitude="form.latitude"
            v-model:longitude="form.longitude"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Transfer Modal -->
    <a-modal
      v-model:open="transferModalVisible"
      title="Transfer SPPG ke Yayasan Lain"
      :confirm-loading="saving"
      @ok="handleTransfer"
      @cancel="transferModalVisible = false"
    >
      <p>Pindahkan <strong>{{ transferRecord?.nama }}</strong> ke Yayasan:</p>
      <a-select
        v-model:value="transferYayasanId"
        placeholder="Pilih Yayasan tujuan"
        :options="yayasanOptions"
        show-search
        :filter-option="filterOption"
        style="width: 100%"
      />
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
const yayasanList = ref([])
const searchText = ref('')
const modalVisible = ref(false)
const transferModalVisible = ref(false)
const editingRecord = ref(null)
const transferRecord = ref(null)
const transferYayasanId = ref(null)

const form = ref({
  nama: '',
  yayasan_id: null,
  alamat: '',
  nomor_telepon: '',
  email: '',
  latitude: 0,
  longitude: 0
})

const columns = [
  { title: 'Kode', dataIndex: 'kode', key: 'kode', width: 130 },
  { title: 'Nama', dataIndex: 'nama', key: 'nama', sorter: (a, b) => a.nama.localeCompare(b.nama) },
  { title: 'Yayasan Induk', key: 'yayasan' },
  { title: 'Status', key: 'is_active', width: 100, align: 'center' },
  { title: 'Aksi', key: 'action', width: 250, align: 'center' }
]

const yayasanOptions = computed(() =>
  yayasanList.value.map(y => ({ value: y.id, label: `${y.kode} - ${y.nama}` }))
)

const filteredData = computed(() => {
  if (!searchText.value) return dataSource.value
  const q = searchText.value.toLowerCase()
  return dataSource.value.filter(
    r => r.nama?.toLowerCase().includes(q) || r.kode?.toLowerCase().includes(q)
  )
})

const filterOption = (input, option) =>
  option.label.toLowerCase().includes(input.toLowerCase())

const fetchData = async () => {
  loading.value = true
  try {
    const [sppgRes, yayasanRes] = await Promise.all([
      organizationService.getSPPGList(),
      organizationService.getYayasanList()
    ])
    dataSource.value = sppgRes.data?.data || sppgRes.data?.sppgs || []
    yayasanList.value = yayasanRes.data?.data || yayasanRes.data?.yayasans || []
  } catch (e) {
    message.error('Gagal memuat data SPPG')
  } finally {
    loading.value = false
  }
}

const openModal = (record = null) => {
  editingRecord.value = record
  if (record) {
    form.value = { nama: record.nama, yayasan_id: record.yayasan_id, alamat: record.alamat, nomor_telepon: record.nomor_telepon, email: record.email, latitude: record.latitude || 0, longitude: record.longitude || 0 }
  } else {
    form.value = { nama: '', yayasan_id: null, alamat: '', nomor_telepon: '', email: '', latitude: 0, longitude: 0 }
  }
  modalVisible.value = true
}

const closeModal = () => {
  modalVisible.value = false
  editingRecord.value = null
}

const handleSave = async () => {
  if (!form.value.nama || !form.value.yayasan_id) {
    message.warning('Nama dan Yayasan wajib diisi')
    return
  }
  saving.value = true
  try {
    if (editingRecord.value) {
      await organizationService.updateSPPG(editingRecord.value.id, form.value)
      message.success('SPPG berhasil diperbarui')
    } else {
      await organizationService.createSPPG(form.value)
      message.success('SPPG berhasil ditambahkan')
    }
    closeModal()
    fetchData()
  } catch (e) {
    message.error(e.response?.data?.message || 'Gagal menyimpan SPPG')
  } finally {
    saving.value = false
  }
}

const toggleStatus = async (record) => {
  try {
    await organizationService.setSPPGStatus(record.id, !record.is_active)
    message.success(`SPPG berhasil ${record.is_active ? 'dinonaktifkan' : 'diaktifkan'}`)
    fetchData()
  } catch (e) {
    message.error('Gagal mengubah status SPPG')
  }
}

const openTransferModal = (record) => {
  transferRecord.value = record
  transferYayasanId.value = null
  transferModalVisible.value = true
}

const handleTransfer = async () => {
  if (!transferYayasanId.value) {
    message.warning('Pilih Yayasan tujuan')
    return
  }
  saving.value = true
  try {
    await organizationService.transferSPPG(transferRecord.value.id, transferYayasanId.value)
    message.success('SPPG berhasil dipindahkan')
    transferModalVisible.value = false
    fetchData()
  } catch (e) {
    message.error(e.response?.data?.message || 'Gagal memindahkan SPPG')
  } finally {
    saving.value = false
  }
}

onMounted(fetchData)
</script>
