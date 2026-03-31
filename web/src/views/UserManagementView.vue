<template>
  <div class="user-management">
    <a-page-header title="Manajemen User" sub-title="Kelola akun pengguna sistem">
      <template #extra>
        <a-button type="primary" @click="openModal()">
          <template #icon><PlusOutlined /></template>
          Tambah User
        </a-button>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input-search
              v-model:value="searchText"
              placeholder="Cari NIK, nama, atau email..."
              allow-clear
            />
          </a-col>
          <a-col :span="5">
            <a-select
              v-model:value="filterRole"
              placeholder="Filter Peran"
              allow-clear
              style="width: 100%"
              :options="allRoleOptions"
            />
          </a-col>
          <a-col :span="5">
            <a-select
              v-model:value="filterStatus"
              placeholder="Filter Status"
              allow-clear
              style="width: 100%"
              :options="[{ value: true, label: 'Aktif' }, { value: false, label: 'Nonaktif' }]"
            />
          </a-col>
        </a-row>

        <a-table
          :columns="columns"
          :data-source="filteredData"
          :loading="loading"
          :pagination="{ pageSize: 10, showSizeChanger: true }"
          row-key="id"
          size="middle"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'role'">
              <a-tag :color="getRoleColor(record.role)">{{ getRoleLabel(record.role) }}</a-tag>
            </template>
            <template v-if="column.key === 'sppg'">
              {{ record.sppg?.nama || '-' }}
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
                  :title="`${record.is_active ? 'Nonaktifkan' : 'Aktifkan'} user ini?`"
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
      :title="editingRecord ? 'Edit User' : 'Tambah User'"
      :confirm-loading="saving"
      @ok="handleSave"
      @cancel="closeModal"
      width="600px"
    >
      <a-form :model="form" layout="vertical">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="NIK" required>
              <a-input v-model:value="form.nik" placeholder="NIK" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Nama Lengkap" required>
              <a-input v-model:value="form.full_name" placeholder="Nama Lengkap" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="Email" required>
              <a-input v-model:value="form.email" placeholder="Email" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="No. Telepon">
              <a-input v-model:value="form.phone_number" placeholder="No. Telepon" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item v-if="!editingRecord" label="Password" required>
          <a-input-password v-model:value="form.password" placeholder="Password" />
        </a-form-item>
        <a-form-item label="Peran" required>
          <a-select
            v-model:value="form.role"
            placeholder="Pilih Peran"
            :options="allowedRoleOptions"
            @change="handleRoleChange"
          />
        </a-form-item>
        <a-form-item v-if="showYayasanField" label="Yayasan" required>
          <a-select
            v-model:value="form.yayasan_id"
            placeholder="Pilih Yayasan"
            :options="yayasanOptions"
            show-search
            :filter-option="filterOption"
            @change="handleYayasanChange"
          />
        </a-form-item>
        <a-form-item v-if="showSPPGField" label="SPPG" required>
          <a-select
            v-model:value="form.sppg_id"
            placeholder="Pilih SPPG"
            :options="sppgOptions"
            show-search
            :filter-option="filterOption"
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
import { useAuthStore } from '@/stores/auth'
import organizationService from '@/services/organizationService'

const authStore = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const dataSource = ref([])
const yayasanList = ref([])
const sppgList = ref([])
const searchText = ref('')
const filterRole = ref(null)
const filterStatus = ref(null)
const modalVisible = ref(false)
const editingRecord = ref(null)

const form = ref({
  nik: '', full_name: '', email: '', phone_number: '', password: '',
  role: null, sppg_id: null, yayasan_id: null
})

const columns = [
  { title: 'NIK', dataIndex: 'nik', key: 'nik', width: 120 },
  { title: 'Nama', dataIndex: 'full_name', key: 'full_name', sorter: (a, b) => (a.full_name || '').localeCompare(b.full_name || '') },
  { title: 'Email', dataIndex: 'email', key: 'email' },
  { title: 'Peran', key: 'role', width: 140 },
  { title: 'SPPG', key: 'sppg', width: 150 },
  { title: 'Yayasan', key: 'yayasan', width: 150 },
  { title: 'Status', key: 'is_active', width: 90, align: 'center' },
  { title: 'Aksi', key: 'action', width: 170, align: 'center' }
]

const ROLE_MAP = {
  superadmin: { label: 'Superadmin', color: 'purple' },
  admin_bgn: { label: 'Admin BGN', color: 'magenta' },
  kepala_yayasan: { label: 'Kepala Yayasan', color: 'geekblue' },
  kepala_sppg: { label: 'Kepala SPPG', color: 'blue' },
  akuntan: { label: 'Akuntan', color: 'cyan' },
  ahli_gizi: { label: 'Ahli Gizi', color: 'green' },
  pengadaan: { label: 'Pengadaan', color: 'lime' },
  chef: { label: 'Chef', color: 'orange' },
  packing: { label: 'Packing', color: 'gold' },
  driver: { label: 'Driver', color: 'volcano' },
  asisten_lapangan: { label: 'Asisten Lapangan', color: 'red' },
  kebersihan: { label: 'Kebersihan', color: 'default' }
}

const OPERATIONAL_ROLES = ['akuntan', 'ahli_gizi', 'pengadaan', 'chef', 'packing', 'driver', 'asisten_lapangan', 'kebersihan']

const getRoleLabel = (role) => ROLE_MAP[role]?.label || role
const getRoleColor = (role) => ROLE_MAP[role]?.color || 'default'

const allRoleOptions = computed(() =>
  Object.entries(ROLE_MAP).map(([value, { label }]) => ({ value, label }))
)

// Roles the current user is allowed to create
const allowedRoleOptions = computed(() => {
  const r = authStore.role
  let allowed = []
  if (r === 'superadmin') {
    allowed = Object.keys(ROLE_MAP)
  } else if (r === 'kepala_yayasan') {
    allowed = ['kepala_sppg', ...OPERATIONAL_ROLES]
  } else if (r === 'kepala_sppg') {
    allowed = OPERATIONAL_ROLES
  }
  return allowed.map(v => ({ value: v, label: ROLE_MAP[v]?.label || v }))
})

// Show Yayasan dropdown for kepala_yayasan role or when superadmin creates non-system roles
const showYayasanField = computed(() => {
  const r = form.value.role
  if (!r) return false
  return r === 'kepala_yayasan' || (authStore.isSuperadmin && !['superadmin', 'admin_bgn'].includes(r))
})

// Show SPPG dropdown for SPPG-level roles
const showSPPGField = computed(() => {
  const r = form.value.role
  if (!r) return false
  return ['kepala_sppg', ...OPERATIONAL_ROLES].includes(r)
})

const yayasanOptions = computed(() => {
  // Kepala Yayasan can only see their own Yayasan
  if (authStore.isKepalaYayasan) {
    return yayasanList.value
      .filter(y => y.id === authStore.yayasanId)
      .map(y => ({ value: y.id, label: `${y.kode} - ${y.nama}` }))
  }
  return yayasanList.value.map(y => ({ value: y.id, label: `${y.kode} - ${y.nama}` }))
})

const sppgOptions = computed(() => {
  let list = sppgList.value
  // Scope by Yayasan if selected or if user is Kepala Yayasan
  if (form.value.yayasan_id) {
    list = list.filter(s => s.yayasan_id === form.value.yayasan_id)
  } else if (authStore.isKepalaYayasan) {
    list = list.filter(s => s.yayasan_id === authStore.yayasanId)
  }
  return list.map(s => ({ value: s.id, label: `${s.kode} - ${s.nama}` }))
})

const filteredData = computed(() => {
  let data = dataSource.value
  if (searchText.value) {
    const q = searchText.value.toLowerCase()
    data = data.filter(r =>
      r.nik?.toLowerCase().includes(q) ||
      r.full_name?.toLowerCase().includes(q) ||
      r.email?.toLowerCase().includes(q)
    )
  }
  if (filterRole.value) data = data.filter(r => r.role === filterRole.value)
  if (filterStatus.value !== null && filterStatus.value !== undefined) {
    data = data.filter(r => r.is_active === filterStatus.value)
  }
  return data
})

const filterOption = (input, option) =>
  option.label.toLowerCase().includes(input.toLowerCase())

const handleRoleChange = () => {
  // Reset dependent fields
  if (!showYayasanField.value) form.value.yayasan_id = null
  if (!showSPPGField.value) form.value.sppg_id = null
}

const handleYayasanChange = () => {
  form.value.sppg_id = null
}

const fetchData = async () => {
  loading.value = true
  try {
    const [usersRes, yayasanRes, sppgRes] = await Promise.all([
      organizationService.getUsers(),
      organizationService.getYayasanList(),
      organizationService.getSPPGList()
    ])
    dataSource.value = usersRes.data?.data || usersRes.data?.users || []
    yayasanList.value = yayasanRes.data?.data || yayasanRes.data?.yayasans || []
    sppgList.value = sppgRes.data?.data || sppgRes.data?.sppgs || []
  } catch (e) {
    message.error('Gagal memuat data user')
  } finally {
    loading.value = false
  }
}

const openModal = (record = null) => {
  editingRecord.value = record
  if (record) {
    form.value = {
      nik: record.nik, full_name: record.full_name, email: record.email,
      phone_number: record.phone_number, password: '',
      role: record.role, sppg_id: record.sppg_id, yayasan_id: record.yayasan_id
    }
  } else {
    form.value = { nik: '', full_name: '', email: '', phone_number: '', password: '', role: null, sppg_id: null, yayasan_id: null }
  }
  modalVisible.value = true
}

const closeModal = () => {
  modalVisible.value = false
  editingRecord.value = null
}

const handleSave = async () => {
  if (!form.value.nik || !form.value.full_name || !form.value.email || !form.value.role) {
    message.warning('NIK, Nama, Email, dan Peran wajib diisi')
    return
  }
  if (!editingRecord.value && !form.value.password) {
    message.warning('Password wajib diisi untuk user baru')
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (editingRecord.value && !payload.password) delete payload.password
    if (editingRecord.value) {
      await organizationService.updateUser(editingRecord.value.id, payload)
      message.success('User berhasil diperbarui')
    } else {
      await organizationService.createUser(payload)
      message.success('User berhasil ditambahkan')
    }
    closeModal()
    fetchData()
  } catch (e) {
    message.error(e.response?.data?.message || 'Gagal menyimpan user')
  } finally {
    saving.value = false
  }
}

const toggleStatus = async (record) => {
  try {
    await organizationService.setUserStatus(record.id, !record.is_active)
    message.success(`User berhasil ${record.is_active ? 'dinonaktifkan' : 'diaktifkan'}`)
    fetchData()
  } catch (e) {
    message.error('Gagal mengubah status user')
  }
}

onMounted(fetchData)
</script>
