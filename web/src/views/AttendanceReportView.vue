<template>
  <div class="attendance-report">
    <a-page-header
      title="Laporan Absensi"
      sub-title="Laporan kehadiran karyawan dengan filter tanggal dan karyawan"
    >
      <template #extra>
        <a-space>
          <a-button 
            type="primary" 
            :loading="exportingExcel"
            @click="exportToExcel"
            :disabled="!hasData"
          >
            <template #icon><FileExcelOutlined /></template>
            Export Excel
          </a-button>
          <a-button 
            :loading="exportingPDF"
            @click="exportToPDF"
            :disabled="!hasData"
          >
            <template #icon><FilePdfOutlined /></template>
            Export PDF
          </a-button>
        </a-space>
      </template>
    </a-page-header>

    <a-card>
      <a-space direction="vertical" style="width: 100%" :size="16">
        <!-- Filter Section -->
        <a-card size="small" title="Filter Laporan">
          <a-form layout="inline" :model="filters" @finish="handleSearch">
            <a-form-item label="Periode">
              <a-range-picker
                v-model:value="dateRange"
                format="DD/MM/YYYY"
                placeholder="['Tanggal Mulai', 'Tanggal Akhir']"
                style="width: 280px"
              />
            </a-form-item>
            
            <a-form-item label="Karyawan">
              <a-select
                v-model:value="filters.employeeId"
                placeholder="Semua Karyawan"
                style="width: 200px"
                show-search
                option-filter-prop="children"
                allow-clear
                :filter-option="filterOption"
              >
                <a-select-option 
                  v-for="employee in employees" 
                  :key="employee.id" 
                  :value="employee.id"
                >
                  {{ employee.full_name }} ({{ employee.nik }})
                </a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item>
              <a-button type="primary" html-type="submit" :loading="loading">
                <template #icon><SearchOutlined /></template>
                Cari
              </a-button>
            </a-form-item>

            <a-form-item>
              <a-button @click="resetFilters">
                <template #icon><ClearOutlined /></template>
                Reset
              </a-button>
            </a-form-item>
          </a-form>
        </a-card>

        <!-- Summary Statistics -->
        <a-row :gutter="16" v-if="hasData">
          <a-col :span="6">
            <a-card size="small">
              <a-statistic
                title="Total Karyawan"
                :value="reportData.length"
                :value-style="{ color: '#1890ff' }"
              />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic
                title="Total Hari Kerja"
                :value="totalWorkDays"
                :value-style="{ color: '#52c41a' }"
              />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic
                title="Total Jam Kerja"
                :value="totalWorkHours"
                :precision="1"
                :value-style="{ color: '#722ed1' }"
                suffix="jam"
              />
            </a-card>
          </a-col>
          <a-col :span="6">
            <a-card size="small">
              <a-statistic
                title="Rata-rata Jam/Hari"
                :value="averageHoursPerDay"
                :precision="1"
                :value-style="{ color: '#fa8c16' }"
                suffix="jam"
              />
            </a-card>
          </a-col>
        </a-row>

        <!-- Report Table -->
        <a-table
          :columns="columns"
          :data-source="reportData"
          :loading="loading"
          :pagination="pagination"
          @change="handleTableChange"
          row-key="employee_id"
          :scroll="{ x: 800 }"
          :expandedRowKeys="expandedRowKeys"
          @expand="onExpand"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'full_name'">
              <span style="font-weight: 500;">{{ record.full_name }}</span>
            </template>
            <template v-else-if="column.key === 'total_hours'">
              {{ formatHours(record.total_hours) }}
            </template>
            <template v-else-if="column.key === 'average_hours'">
              {{ formatHours(record.average_hours) }}
            </template>
            <template v-else-if="column.key === 'attendance_rate'">
              <a-progress 
                :percent="calculateAttendanceRate(record.total_days)" 
                size="small"
                :status="getAttendanceStatus(record.total_days)"
              />
            </template>
          </template>

          <!-- Expandable Row Content -->
          <template #expandedRowRender="{ record }">
            <div style="padding: 16px; background: #fafafa;">
              <a-spin :spinning="expandedRowLoading[record.employee_id]">
                <div v-if="expandedRowData[record.employee_id] && expandedRowData[record.employee_id].length > 0">
                  <h4 style="margin-bottom: 12px; color: #5A4372;">
                    <ClockCircleOutlined /> Detail Kehadiran - {{ record.full_name }}
                  </h4>
                  <a-table
                    :columns="detailColumns"
                    :data-source="expandedRowData[record.employee_id]"
                    :pagination="false"
                    size="small"
                    :scroll="{ x: 600 }"
                  >
                    <template #bodyCell="{ column, record: detailRecord }">
                      <template v-if="column.key === 'date'">
                        {{ formatDate(detailRecord.date) }}
                      </template>
                      <template v-else-if="column.key === 'check_in'">
                        <a-tag color="green">
                          <LoginOutlined /> {{ formatTime(detailRecord.check_in) }}
                        </a-tag>
                      </template>
                      <template v-else-if="column.key === 'check_out'">
                        <a-tag v-if="detailRecord.check_out" color="red">
                          <LogoutOutlined /> {{ formatTime(detailRecord.check_out) }}
                        </a-tag>
                        <a-tag v-else color="orange">
                          Belum Check Out
                        </a-tag>
                      </template>
                      <template v-else-if="column.key === 'work_hours'">
                        <span style="font-weight: 500;">
                          {{ formatHours(detailRecord.work_hours) }}
                        </span>
                      </template>
                      <template v-else-if="column.key === 'status'">
                        <a-tag :color="getStatusColor(detailRecord)">
                          {{ getStatusText(detailRecord) }}
                        </a-tag>
                      </template>
                    </template>
                  </a-table>
                </div>
                <a-empty 
                  v-else-if="!expandedRowLoading[record.employee_id]"
                  description="Tidak ada data detail"
                  :image="Empty.PRESENTED_IMAGE_SIMPLE"
                />
              </a-spin>
            </div>
          </template>
        </a-table>

        <!-- Empty State -->
        <a-empty 
          v-if="!loading && !hasData" 
          description="Tidak ada data absensi untuk periode yang dipilih"
        >
          <a-button type="primary" @click="resetFilters">Reset Filter</a-button>
        </a-empty>
      </a-space>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { message, Empty } from 'ant-design-vue'
import { 
  FileExcelOutlined, 
  FilePdfOutlined, 
  SearchOutlined, 
  ClearOutlined,
  ClockCircleOutlined,
  LoginOutlined,
  LogoutOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import attendanceService from '@/services/attendanceService'
import employeeService from '@/services/employeeService'

const loading = ref(false)
const exportingExcel = ref(false)
const exportingPDF = ref(false)
const reportData = ref([])
const employees = ref([])
const dateRange = ref([])
const expandedRowKeys = ref([])
const expandedRowData = ref({})
const expandedRowLoading = ref({})

const filters = reactive({
  employeeId: undefined
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showQuickJumper: true,
  showTotal: (total, range) => `${range[0]}-${range[1]} dari ${total} data`
})

const columns = [
  {
    title: 'Nama Karyawan',
    key: 'full_name',
    dataIndex: 'full_name',
    sorter: true,
    width: 200
  },
  {
    title: 'Posisi',
    dataIndex: 'position',
    key: 'position',
    width: 150
  },
  {
    title: 'Total Hari',
    dataIndex: 'total_days',
    key: 'total_days',
    sorter: true,
    width: 100,
    align: 'center'
  },
  {
    title: 'Total Jam',
    key: 'total_hours',
    sorter: true,
    width: 120,
    align: 'right'
  },
  {
    title: 'Rata-rata Jam/Hari',
    key: 'average_hours',
    sorter: true,
    width: 140,
    align: 'right'
  },
  {
    title: 'Tingkat Kehadiran',
    key: 'attendance_rate',
    width: 150,
    align: 'center'
  }
]

const detailColumns = [
  {
    title: 'Tanggal',
    key: 'date',
    width: 120
  },
  {
    title: 'Check In',
    key: 'check_in',
    width: 120
  },
  {
    title: 'Check Out',
    key: 'check_out',
    width: 120
  },
  {
    title: 'Jam Kerja',
    key: 'work_hours',
    width: 100,
    align: 'right'
  },
  {
    title: 'Status',
    key: 'status',
    width: 120,
    align: 'center'
  }
]

const hasData = computed(() => reportData.value.length > 0)

const totalWorkDays = computed(() => {
  return reportData.value.reduce((sum, item) => sum + item.total_days, 0)
})

const totalWorkHours = computed(() => {
  return reportData.value.reduce((sum, item) => sum + item.total_hours, 0)
})

const averageHoursPerDay = computed(() => {
  if (totalWorkDays.value === 0) return 0
  return totalWorkHours.value / totalWorkDays.value
})

const fetchEmployees = async () => {
  try {
    const response = await employeeService.getEmployees({ is_active: true })
    employees.value = response.data || []
  } catch (error) {
    console.error('Gagal memuat data karyawan:', error)
  }
}

const fetchReport = async () => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    message.warning('Silakan pilih periode tanggal terlebih dahulu')
    return
  }

  loading.value = true
  try {
    const params = {
      start_date: dateRange.value[0].format('YYYY-MM-DD'),
      end_date: dateRange.value[1].format('YYYY-MM-DD')
    }

    if (filters.employeeId) {
      params.employee_id = filters.employeeId
    }

    console.log('[AttendanceReport] Fetching report with params:', params)
    const response = await attendanceService.getAttendanceReport(params)
    console.log('[AttendanceReport] Response:', response)
    
    // Handle response structure
    if (response && response.success && response.data) {
      reportData.value = response.data
    } else if (response && response.data) {
      reportData.value = response.data
    } else {
      reportData.value = []
    }
    
    pagination.total = reportData.value.length
    
    if (reportData.value.length === 0) {
      console.warn('[AttendanceReport] No data returned from API')
    } else {
      console.log('[AttendanceReport] Loaded', reportData.value.length, 'records')
    }
  } catch (error) {
    message.error('Gagal memuat laporan absensi')
    console.error('[AttendanceReport] Error:', error)
    console.error('[AttendanceReport] Error response:', error.response?.data)
    reportData.value = []
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.current = 1
  fetchReport()
}

const handleTableChange = (pag, filters, sorter) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
}

const onExpand = async (expanded, record) => {
  if (expanded) {
    // Add to expanded keys
    if (!expandedRowKeys.value.includes(record.employee_id)) {
      expandedRowKeys.value.push(record.employee_id)
    }
    
    // Load detail data if not already loaded
    if (!expandedRowData.value[record.employee_id]) {
      await loadExpandedRowData(record)
    }
  } else {
    // Remove from expanded keys
    expandedRowKeys.value = expandedRowKeys.value.filter(key => key !== record.employee_id)
  }
}

const loadExpandedRowData = async (record) => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    message.warning('Periode tanggal tidak valid')
    return
  }

  expandedRowLoading.value[record.employee_id] = true

  try {
    console.log('[loadExpandedRowData] Loading details for employee:', record.employee_id)
    const response = await attendanceService.getAttendanceByDateRange(
      record.employee_id,
      dateRange.value[0].format('YYYY-MM-DD'),
      dateRange.value[1].format('YYYY-MM-DD')
    )
    
    console.log('[loadExpandedRowData] Response:', response)
    
    if (response && response.success && response.data) {
      expandedRowData.value[record.employee_id] = response.data
    } else if (response && response.data) {
      expandedRowData.value[record.employee_id] = response.data
    } else {
      expandedRowData.value[record.employee_id] = []
    }
    
    console.log('[loadExpandedRowData] Loaded', expandedRowData.value[record.employee_id].length, 'records')
  } catch (error) {
    message.error('Gagal memuat detail absensi')
    console.error('[loadExpandedRowData] Error:', error)
    expandedRowData.value[record.employee_id] = []
  } finally {
    expandedRowLoading.value[record.employee_id] = false
  }
}

const resetFilters = () => {
  dateRange.value = []
  filters.employeeId = undefined
  reportData.value = []
  pagination.current = 1
  expandedRowKeys.value = []
  expandedRowData.value = {}
}

const exportToExcel = async () => {
  if (!hasData.value) {
    message.warning('Tidak ada data untuk diekspor')
    return
  }

  exportingExcel.value = true
  try {
    const params = {
      start_date: dateRange.value[0].format('YYYY-MM-DD'),
      end_date: dateRange.value[1].format('YYYY-MM-DD'),
      format: 'excel'
    }

    if (filters.employeeId) {
      params.employee_id = filters.employeeId
    }

    const response = await attendanceService.exportToExcel(params)
    
    // Create download link
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `laporan-absensi-${params.start_date}-${params.end_date}.xlsx`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)

    message.success('Laporan Excel berhasil diunduh')
  } catch (error) {
    message.error('Gagal mengekspor ke Excel')
    console.error(error)
  } finally {
    exportingExcel.value = false
  }
}

const exportToPDF = async () => {
  if (!hasData.value) {
    message.warning('Tidak ada data untuk diekspor')
    return
  }

  exportingPDF.value = true
  try {
    const params = {
      start_date: dateRange.value[0].format('YYYY-MM-DD'),
      end_date: dateRange.value[1].format('YYYY-MM-DD'),
      format: 'pdf'
    }

    if (filters.employeeId) {
      params.employee_id = filters.employeeId
    }

    const response = await attendanceService.exportToPDF(params)
    
    // Create download link
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `laporan-absensi-${params.start_date}-${params.end_date}.pdf`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)

    message.success('Laporan PDF berhasil diunduh')
  } catch (error) {
    message.error('Gagal mengekspor ke PDF')
    console.error(error)
  } finally {
    exportingPDF.value = false
  }
}

const filterOption = (input, option) => {
  return option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
}

const formatHours = (hours) => {
  if (!hours) return '0.0 jam'
  return `${parseFloat(hours).toFixed(1)} jam`
}

const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('DD/MM/YYYY')
}

const formatTime = (time) => {
  if (!time) return '-'
  return dayjs(time).format('HH:mm')
}

const calculateAttendanceRate = (totalDays) => {
  try {
    if (!dateRange.value || dateRange.value.length !== 2) return 0
    if (!dateRange.value[0] || !dateRange.value[1]) return 0
    
    const workingDays = getWorkingDays(dateRange.value[0], dateRange.value[1])
    if (workingDays === 0) return 0
    
    return Math.round((totalDays / workingDays) * 100)
  } catch (error) {
    console.error('[calculateAttendanceRate] Error:', error)
    return 0
  }
}

const getWorkingDays = (startDate, endDate) => {
  try {
    if (!startDate || !endDate) return 0
    
    let count = 0
    let current = startDate.clone()
    
    while (current.isSameOrBefore(endDate)) {
      // Skip weekends (Saturday = 6, Sunday = 0)
      if (current.day() !== 0 && current.day() !== 6) {
        count++
      }
      current = current.add(1, 'day')
    }
    
    return count
  } catch (error) {
    console.error('[getWorkingDays] Error:', error)
    return 0
  }
}

const getAttendanceStatus = (totalDays) => {
  const rate = calculateAttendanceRate(totalDays)
  if (rate >= 90) return 'success'
  if (rate >= 75) return 'normal'
  return 'exception'
}

const getStatusColor = (record) => {
  if (!record.check_out) return 'orange'
  if (record.work_hours >= 8) return 'green'
  if (record.work_hours >= 6) return 'blue'
  return 'red'
}

const getStatusText = (record) => {
  if (!record.check_out) return 'Belum Check Out'
  if (record.work_hours >= 8) return 'Lengkap'
  if (record.work_hours >= 6) return 'Cukup'
  return 'Kurang'
}

onMounted(async () => {
  await fetchEmployees()
  
  // Set default date range to current month
  const now = dayjs()
  dateRange.value = [
    now.startOf('month'),
    now.endOf('month')
  ]
  
  // Auto-fetch report with default date range
  console.log('[AttendanceReport] Auto-fetching report on mount')
  await fetchReport()
})
</script>

<style scoped>
.attendance-report {
  padding: 24px;
}

:deep(.ant-table-thead > tr > th) {
  background-color: #fafafa;
  font-weight: 600;
}

:deep(.ant-statistic-content) {
  font-size: 20px;
}

:deep(.ant-progress-line) {
  margin-right: 8px;
}

/* Expandable row styling */
:deep(.ant-table-expanded-row > td) {
  padding: 0 !important;
  background: #f5f5f5;
}

:deep(.ant-table-expanded-row .ant-table) {
  margin: 0;
}

:deep(.ant-table-expanded-row .ant-table-thead > tr > th) {
  background-color: #fff;
  border-bottom: 1px solid #f0f0f0;
}

/* Fix expand icon alignment */
:deep(.ant-table-row-expand-icon-cell) {
  vertical-align: middle !important;
}

:deep(.ant-table-row-expand-icon) {
  color: #5A4372;
  vertical-align: middle;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

:deep(.ant-table-row-expand-icon:hover) {
  color: #7B5E9D;
}

/* Ensure table cells are vertically centered */
:deep(.ant-table-tbody > tr > td) {
  vertical-align: middle !important;
}

/* Detail table styling */
:deep(.ant-table-small .ant-table-tbody > tr > td) {
  padding: 8px;
}
</style>