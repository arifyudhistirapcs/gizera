<template>
  <div class="delivery-timeline">
    <a-timeline>
      <a-timeline-item
        v-for="stage in stages"
        :key="stage.status"
        :color="getStageColor(stage.status)"
      >
        <template #dot>
          <check-circle-outlined v-if="isCompleted(stage.status)" style="font-size: 16px" />
          <sync-outlined v-else-if="isInProgress(stage.status)" style="font-size: 16px" />
          <clock-circle-outlined v-else style="font-size: 16px" />
        </template>
        
        <div class="timeline-content">
          <div class="stage-header">
            <span class="stage-title">{{ stage.title }}</span>
            <a-tooltip v-if="stage.status === 'sudah_diterima_pihak_sekolah' && isCompletedOrInProgress(stage.status)" title="Lihat e-POD">
              <eye-outlined 
                class="epod-icon" 
                @click="$emit('viewEpod')"
              />
            </a-tooltip>
          </div>
          <div class="stage-description">{{ stage.description }}</div>
          <div v-if="getStageTimestamp(stage.status)" class="stage-timestamp">
            {{ formatTimestamp(getStageTimestamp(stage.status)) }}
          </div>
        </div>
      </a-timeline-item>
    </a-timeline>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  CheckCircleOutlined,
  SyncOutlined,
  ClockCircleOutlined,
  EyeOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/id'
import timezone from 'dayjs/plugin/timezone'
import utc from 'dayjs/plugin/utc'

dayjs.extend(utc)
dayjs.extend(timezone)
dayjs.locale('id')

const props = defineProps({
  currentStatus: {
    type: String,
    required: true
  },
  activityLog: {
    type: Array,
    default: () => []
  }
})

defineEmits(['viewEpod'])

// Define all 16 lifecycle stages
const stages = [
  {
    status: 'sedang_dimasak',
    title: 'Sedang Dimasak',
    description: 'Menu sedang dalam proses memasak'
  },
  {
    status: 'selesai_dimasak',
    title: 'Selesai Dimasak',
    description: 'Proses memasak telah selesai'
  },
  {
    status: 'siap_dipacking',
    title: 'Siap Dipacking',
    description: 'Menu siap untuk dikemas'
  },
  {
    status: 'selesai_dipacking',
    title: 'Selesai Dipacking',
    description: 'Menu telah dikemas'
  },
  {
    status: 'siap_dikirim',
    title: 'Siap Dikirim',
    description: 'Menu siap untuk dikirim'
  },
  {
    status: 'diperjalanan',
    title: 'Diperjalanan',
    description: 'Driver sedang dalam perjalanan ke sekolah'
  },
  {
    status: 'sudah_sampai_sekolah',
    title: 'Sudah Sampai Sekolah',
    description: 'Driver telah tiba di sekolah'
  },
  {
    status: 'sudah_diterima_pihak_sekolah',
    title: 'Sudah Diterima',
    description: 'Menu telah diterima oleh pihak sekolah'
  },
  {
    status: 'driver_menuju_lokasi_pengambilan',
    title: 'Driver Menuju Lokasi Pengambilan',
    description: 'Driver menuju lokasi untuk mengambil ompreng'
  },
  {
    status: 'driver_tiba_di_lokasi_pengambilan',
    title: 'Driver Tiba di Lokasi',
    description: 'Driver telah sampai di lokasi untuk pengambilan'
  },
  {
    status: 'driver_kembali_ke_sppg',
    title: 'Driver Kembali ke SPPG',
    description: 'Driver dalam perjalanan kembali ke SPPG'
  },
  {
    status: 'driver_tiba_di_sppg',
    title: 'Driver Tiba di SPPG',
    description: 'Driver telah tiba di SPPG dengan ompreng'
  },
  {
    status: 'ompreng_siap_dicuci',
    title: 'Ompreng Siap Dicuci',
    description: 'Ompreng siap untuk proses pencucian'
  },
  {
    status: 'ompreng_proses_pencucian',
    title: 'Proses Pencucian',
    description: 'Ompreng sedang dalam proses pencucian'
  },
  {
    status: 'ompreng_selesai_dicuci',
    title: 'Selesai Dicuci',
    description: 'Ompreng telah selesai dicuci'
  }
]

// Get the index of current status
const currentStatusIndex = computed(() => {
  return stages.findIndex(stage => stage.status === props.currentStatus)
})

// Check if a stage is completed
const isCompleted = (status) => {
  const stageIndex = stages.findIndex(stage => stage.status === status)
  // Mark stage as completed if it's before current stage
  // OR if it's the last stage (index 14 for 15 stages) and current stage has reached it
  return stageIndex < currentStatusIndex.value || (stageIndex === 14 && currentStatusIndex.value >= 14)
}

// Check if a stage is in progress
const isInProgress = (status) => {
  return status === props.currentStatus
}

// Check if a stage is completed or in progress
const isCompletedOrInProgress = (status) => {
  return isCompleted(status) || isInProgress(status)
}

// Get stage color based on status
const getStageColor = (status) => {
  if (isCompleted(status)) {
    return 'green'
  } else if (isInProgress(status)) {
    return 'blue'
  } else {
    return 'gray'
  }
}

// Get timestamp for a stage from activity log
const getStageTimestamp = (status) => {
  const activity = props.activityLog.find(log => log.to_status === status)
  return activity?.transitioned_at
}

// Format timestamp without timezone conversion (display as-is from backend)
const formatTimestamp = (timestamp) => {
  if (!timestamp) return ''
  // Parse as UTC and display without timezone conversion
  return dayjs.utc(timestamp).format('DD MMM YYYY, HH:mm') + ' WIB'
}
</script>

<style scoped>
.delivery-timeline {
  padding: 16px 0;
}

.timeline-content {
  padding-bottom: 16px;
}

.stage-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stage-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(0, 0, 0, 0.85);
}

.epod-icon {
  font-size: 16px;
  color: #5A4372;
  cursor: pointer;
  transition: color 0.2s;
}

.epod-icon:hover {
  color: #7B5A9A;
}

.stage-description {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.65);
  margin-bottom: 4px;
  margin-top: 4px;
}

.stage-timestamp {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
  font-style: italic;
}

:deep(.ant-timeline-item-tail) {
  border-left-width: 2px;
}

:deep(.ant-timeline-item-head) {
  background-color: transparent;
}
</style>
