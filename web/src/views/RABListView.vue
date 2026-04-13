<template>
  <div class="rab-list">
    <a-page-header
      title="Daftar RAB"
      sub-title="Rencana Anggaran Belanja dari menu plan"
    />

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Filters -->
        <a-row :gutter="16">
          <a-col :span="8">
            <a-input
              v-model:value="searchText"
              placeholder="Cari nomor RAB..."
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
              placeholder="Filter Status"
              style="width: 100%"
              @change="handleSearch"
              allow-clear
              size="large"
            >
              <a-select-option value="draft">Draft</a-select-option>
              <a-select-option value="approved_sppg">Approved SPPG</a-select-option>
              <a-select-option value="approved_yayasan">Approved Yayasan</a-select-option>
              <a-select-option value="revision_requested">Revisi</a-select-option>
              <a-select-option value="completed">Selesai</a-select-option>
            </a-select>
          </a-col>
        </a-row>

        <!-- Table -->
        <a-table
          :columns="columns"
          :data-source="rabList"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="id"
          :custom-row="customRow"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'status'">
              <a-tag :color="getStatusColor(record.status)">
                {{ getStatusLabel(record.status) }}
              </a-tag>
            </template>
            <template v-else-if="column.key === 'total_amount'">
              {{ formatRupiah(record.total_amount) }}
            </template>
            <template v-else-if="column.key === 'created_at'">
              {{ formatDate(record.created_at) }}
            </template>
            <template v-else-if="column.key === 'menu_plan'">
              {{ record.menu_plan ? `${formatDate(record.menu_plan.week_start)} - ${formatDate(record.menu_plan.week_end)}` : '-' }}
            </template>
            <template v-else-if="column.key === 'sppg'">
              {{ record.sppg?.nama || record.sppg?.name || '-' }}
            </template>
          </template>
        </a-table>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import rabService from '@/services/rabService'

const router = useRouter()

const loading = ref(false)
const rabList = ref([])
const searchText = ref('')
const filterStatus = ref(undefined)

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const columns = [
  { title: 'Nomor RAB', dataIndex: 'rab_number', key: 'rab_number' },
  { title: 'Menu Plan', key: 'menu_plan' },
  { title: 'SPPG', key: 'sppg' },
  { title: 'Status', key: 'status', width: 160 },
  { title: 'Total', key: 'total_amount' },
  { title: 'Tanggal', key: 'created_at' }
]

const fetchRABList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.current,
      per_page: pagination.pageSize,
      status: filterStatus.value || undefined,
      search: searchText.value || undefined
    }
    const response = await rabService.getRABList(params)
    rabList.value = response.data.rabs || response.data.data || []
    pagination.total = response.data.total || 0
  } catch (error) {
    message.error('Gagal memuat data RAB')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleTableChange = (pag) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchRABList()
}

const handleSearch = () => {
  pagination.current = 1
  fetchRABList()
}

const customRow = (record) => ({
  onClick: () => {
    router.push(`/rab/${record.id}`)
  },
  style: { cursor: 'pointer' }
})

const getStatusColor = (status) => {
  const colors = {
    draft: 'default',
    approved_sppg: 'blue',
    approved_yayasan: 'green',
    revision_requested: 'orange',
    completed: 'purple'
  }
  return colors[status] || 'default'
}

const getStatusLabel = (status) => {
  const labels = {
    draft: 'Draft',
    approved_sppg: 'Approved SPPG',
    approved_yayasan: 'Approved Yayasan',
    revision_requested: 'Revisi',
    completed: 'Selesai'
  }
  return labels[status] || status
}

const formatRupiah = (value) => {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(value || 0)
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' })
}

onMounted(() => {
  fetchRABList()
})
</script>

<style scoped>
.rab-list {
  padding: 24px;
}
</style>
