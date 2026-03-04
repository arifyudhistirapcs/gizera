<template>
  <div class="profile-page">
    <!-- Gradient Banner -->
    <div class="profile-banner">
      <div class="banner-decoration"></div>
      <div class="avatar-wrapper">
        <div class="avatar-circle">
          <van-icon name="user-o" size="40" color="var(--h-primary)" />
        </div>
      </div>
    </div>

    <!-- User Info -->
    <div class="user-info">
      <h2 class="user-name">{{ user?.name || 'User' }}</h2>
      <p class="user-nik">NIK: {{ user?.nik || '-' }}</p>
    </div>

    <!-- Tabs: Personal & Pengaturan -->
    <div class="profile-tabs-wrapper">
      <van-tabs v-model:active="activeTab" shrink animated>
        <van-tab title="Personal">
          <div class="tab-content">
            <div class="info-card">
              <van-cell-group :border="false">
                <van-cell title="NIK" :value="user?.nik || '-'" />
                <van-cell title="Nama Lengkap" :value="user?.name || '-'" />
                <van-cell title="Email" :value="user?.email || '-'" />
                <van-cell title="Role" :value="user?.role || '-'" />
              </van-cell-group>
            </div>
          </div>
        </van-tab>
        <van-tab title="Pengaturan">
          <div class="tab-content">
            <div class="info-card">
              <van-cell-group :border="false">
                <van-cell title="Versi Aplikasi" value="1.0.0" />
                <van-cell title="Mode Offline" is-link @click="showOfflineInfo" />
              </van-cell-group>
            </div>
          </div>
        </van-tab>
      </van-tabs>
    </div>

    <!-- Logout Button -->
    <div class="logout-section">
      <van-button
        block
        type="danger"
        icon="cross"
        :loading="loading"
        loading-text="Logging out..."
        class="logout-btn"
        @click="handleLogout"
      >
        Logout
      </van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { showDialog, showNotify } from 'vant'

const router = useRouter()
const authStore = useAuthStore()

const activeTab = ref(0)
const loading = ref(false)

const user = computed(() => authStore.user)

// Check if user can access tasks page (driver or asisten_lapangan only)
const canAccessTasks = computed(() => {
  const allowedRoles = ['driver', 'asisten_lapangan']
  const userRole = authStore.user?.role?.toLowerCase()
  return allowedRoles.includes(userRole)
})

const handleLogout = async () => {
  showDialog({
    title: 'Konfirmasi Logout',
    message: 'Apakah Anda yakin ingin keluar?',
    showCancelButton: true,
    confirmButtonText: 'Ya, Keluar',
    cancelButtonText: 'Batal'
  }).then(async () => {
    loading.value = true
    
    try {
      await authStore.logout()
      
      showNotify({
        type: 'success',
        message: 'Logout berhasil'
      })
      
      router.push('/login')
    } catch (error) {
      console.error('Logout error:', error)
      showNotify({
        type: 'warning',
        message: 'Logout berhasil (offline mode)'
      })
      router.push('/login')
    } finally {
      loading.value = false
    }
  }).catch(() => {
    // User cancelled
  })
}

const showOfflineInfo = () => {
  showDialog({
    title: 'Mode Offline',
    message: 'Aplikasi ini dapat bekerja secara offline. Data akan disinkronkan otomatis saat koneksi tersedia.',
    confirmButtonText: 'Mengerti'
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: var(--h-bg-primary);
}

/* Gradient Banner */
.profile-banner {
  background: linear-gradient(135deg, var(--h-primary) 0%, var(--h-accent) 100%);
  height: 180px;
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  overflow: hidden;
}

.banner-decoration {
  position: absolute;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  top: -80px;
  right: -60px;
}

.banner-decoration::after {
  content: '';
  position: absolute;
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  bottom: -100px;
  left: -160px;
}

/* Avatar overlapping the banner */
.avatar-wrapper {
  position: absolute;
  bottom: -40px;
  z-index: 2;
}

.avatar-circle {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: var(--h-bg-card);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--h-shadow-card);
  border: 3px solid var(--h-bg-card);
}

/* User Info below avatar */
.user-info {
  text-align: center;
  padding-top: 50px;
  padding-bottom: var(--h-spacing-lg);
}

.user-name {
  font-family: var(--h-font-family);
  font-size: 20px;
  font-weight: 700;
  color: var(--h-text-primary);
  margin: 0 0 4px;
}

.user-nik {
  font-family: var(--h-font-family);
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: 0;
}

/* Tabs */
.profile-tabs-wrapper {
  padding: 0 var(--h-spacing-lg);
}

.profile-tabs-wrapper :deep(.van-tabs__nav) {
  background: transparent;
}

.profile-tabs-wrapper :deep(.van-tabs__line) {
  background: var(--h-primary);
  height: 3px;
  border-radius: 2px;
}

.profile-tabs-wrapper :deep(.van-tab) {
  color: var(--h-text-secondary);
  font-weight: 500;
  font-family: var(--h-font-family);
}

.profile-tabs-wrapper :deep(.van-tab--active) {
  color: var(--h-primary);
  font-weight: 600;
}

.tab-content {
  padding-top: var(--h-spacing-lg);
}

.info-card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  overflow: hidden;
}

.info-card :deep(.van-cell) {
  background: transparent;
}

.info-card :deep(.van-cell__title) {
  color: var(--h-text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.info-card :deep(.van-cell__value) {
  color: var(--h-text-primary);
  font-weight: 500;
}

/* Logout */
.logout-section {
  padding: var(--h-spacing-2xl) var(--h-spacing-lg);
}

.logout-btn {
  height: 50px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  border-radius: var(--h-radius-md) !important;
  box-shadow: 0px 4px 12px rgba(238, 93, 80, 0.3);
  transition: all var(--h-transition-base);
}

.logout-btn:active {
  transform: scale(0.98);
}
</style>
