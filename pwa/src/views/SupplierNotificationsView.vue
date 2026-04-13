<template>
  <div class="supplier-notifications-page">
    <!-- Header -->
    <van-nav-bar title="Notifikasi" fixed>
      <template #right>
        <span v-if="unreadCount > 0" class="mark-all" @click="markAllRead">
          Tandai semua dibaca
        </span>
      </template>
    </van-nav-bar>

    <div class="page-content">
      <!-- Loading -->
      <van-loading v-if="loading" type="spinner" vertical class="page-loading">
        Memuat notifikasi...
      </van-loading>

      <!-- Empty -->
      <van-empty
        v-else-if="notifications.length === 0"
        image="search"
        description="Tidak ada notifikasi"
      />

      <!-- Notification List -->
      <van-cell-group v-else inset>
        <van-cell
          v-for="notif in notifications"
          :key="notif.id"
          :title="notif.title || 'Notifikasi'"
          :label="notif.message"
          clickable
          :class="{ 'notif-unread': !notif.is_read }"
          @click="onNotifClick(notif)"
        >
          <template #value>
            <div class="notif-meta">
              <span class="notif-time">{{ formatTime(notif.created_at) }}</span>
              <van-badge v-if="!notif.is_read" dot />
            </div>
          </template>
        </van-cell>
      </van-cell-group>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import api from '@/services/api'

const router = useRouter()
const loading = ref(false)
const notifications = ref([])
const unreadCount = ref(0)

function formatTime(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const diffMs = now - d
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return 'Baru saja'
  if (diffMin < 60) return `${diffMin} menit lalu`
  const diffHour = Math.floor(diffMin / 60)
  if (diffHour < 24) return `${diffHour} jam lalu`
  const diffDay = Math.floor(diffHour / 24)
  if (diffDay < 7) return `${diffDay} hari lalu`
  return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short' })
}

async function loadNotifications() {
  loading.value = true
  try {
    const [notifRes, countRes] = await Promise.all([
      api.get('/notifications'),
      api.get('/notifications/unread-count')
    ])
    notifications.value = notifRes.data?.data ?? notifRes.data ?? []
    unreadCount.value = countRes.data?.data?.count ?? countRes.data?.count ?? 0
  } catch (e) {
    showToast(e.response?.data?.message || 'Gagal memuat notifikasi')
  } finally {
    loading.value = false
  }
}

async function onNotifClick(notif) {
  // Mark as read
  if (!notif.is_read) {
    try {
      await api.put(`/notifications/${notif.id}/read`)
      notif.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (e) {
      // Silent fail
    }
  }

  // Navigate to related page based on notification type
  if (notif.reference_type === 'purchase_order' && notif.reference_id) {
    router.push(`/supplier-po/${notif.reference_id}`)
  } else if (notif.reference_type === 'invoice' && notif.reference_id) {
    router.push('/supplier-invoices')
  } else if (notif.reference_type === 'payment') {
    router.push('/supplier-payments')
  }
}

async function markAllRead() {
  try {
    await api.put('/notifications/read-all')
    notifications.value.forEach(n => { n.is_read = true })
    unreadCount.value = 0
    showToast('Semua notifikasi ditandai dibaca')
  } catch (e) {
    showToast('Gagal menandai notifikasi')
  }
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.supplier-notifications-page {
  min-height: 100vh;
  background: #F7F8FA;
  padding-top: 46px;
  padding-bottom: 88px;
}

.page-content {
  padding: 16px 0;
}

.page-loading {
  padding: 60px 0;
}

.mark-all {
  font-size: 12px;
  color: #1989fa;
  cursor: pointer;
}

.notif-unread {
  background: #f0f7ff;
}

.notif-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.notif-time {
  font-size: 11px;
  color: #969799;
  white-space: nowrap;
}
</style>
