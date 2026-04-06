<template>
  <div>
    <!-- Role-specific dashboard redirect for Kepala SPPG -->
    <div v-if="authStore.user?.role === 'kepala_sppg'" class="redirect-container">
      <LottiePlayer src="/lottie/loading-cooking.json" width="150px" height="150px" />
      <p class="redirect-text">Mengarahkan ke Dashboard Kepala SPPG...</p>
    </div>

    <!-- General dashboard for other roles -->
    <template v-else>
      <div class="h-card welcome-card">
        <div class="welcome-card__content">
          <div class="welcome-card__text">
            <h2 class="welcome-title">Selamat Datang, {{ userName }}!</h2>
            <p class="welcome-subtitle">Anda login sebagai {{ roleLabel }}</p>
          </div>
          <img v-if="welcomeIllustration" :src="welcomeIllustration" alt="Welcome" class="welcome-card__illustration" />
        </div>
      </div>

      <!-- Overview Highlight Widget -->
      <div class="overview-widget">
        <h3 class="overview-widget__title">Overview</h3>
        <div class="overview-widget__cards">
          <div class="overview-widget__card">
            <span class="overview-widget__label">Total Resep</span>
            <span class="overview-widget__value">0</span>
          </div>
          <div class="overview-widget__card">
            <span class="overview-widget__label">Menu Aktif</span>
            <span class="overview-widget__value">0</span>
          </div>
          <div class="overview-widget__card">
            <span class="overview-widget__label">Pengiriman</span>
            <span class="overview-widget__value">0</span>
          </div>
          <div class="overview-widget__card">
            <span class="overview-widget__label">Stok Menipis</span>
            <span class="overview-widget__value">0</span>
          </div>
        </div>
      </div>

      <div class="stats-row">
        <HStatCard
          :icon="BookOutlined"
          icon-bg="#FDEAE7"
          label="Total Resep"
          value="0"
          :loading="false"
        />
        <HStatCard
          :icon="CalendarOutlined"
          icon-bg="#D1FAE5"
          label="Menu Aktif"
          value="0"
          :loading="false"
        />
        <HStatCard
          :icon="CarOutlined"
          icon-bg="#DBEAFE"
          label="Pengiriman Hari Ini"
          value="0"
          :loading="false"
        />
        <HStatCard
          :icon="WarningOutlined"
          icon-bg="#FEF3C7"
          label="Stok Menipis"
          value="0"
          :loading="false"
        />
      </div>

      <div class="content-row">
        <div class="h-card content-card">
          <h3 class="card-title">Status Produksi Hari Ini</h3>
          <HEmptyState lottie="/lottie/empty-box.json" description="Data akan ditampilkan setelah modul KDS diimplementasikan" />
        </div>
        <div class="h-card content-card">
          <h3 class="card-title">Aktivitas Terbaru</h3>
          <HEmptyState lottie="/lottie/empty-box.json" description="Data akan ditampilkan setelah modul audit trail diimplementasikan" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { h, computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  BookOutlined,
  CalendarOutlined,
  CarOutlined,
  WarningOutlined
} from '@ant-design/icons-vue'
import HStatCard from '@/components/horizon/HStatCard.vue'
import HEmptyState from '@/components/common/HEmptyState.vue'
import LottiePlayer from '@/components/common/LottiePlayer.vue'
import welcomeChefSvg from '@/assets/illustrations/welcome-chef.svg'
import welcomeDeliverySvg from '@/assets/illustrations/welcome-delivery.svg'
import welcomeAnalyticsSvg from '@/assets/illustrations/welcome-analytics.svg'

const router = useRouter()
const authStore = useAuthStore()
const redirecting = ref(false)

const userName = computed(() => {
  return authStore.user?.fullName || authStore.user?.email || 'User'
})

const roleLabel = computed(() => {
  const roleLabels = {
    'kepala_sppg': 'Kepala SPPG',
    'kepala_yayasan': 'Kepala Yayasan',
    'akuntan': 'Akuntan',
    'ahli_gizi': 'Ahli Gizi',
    'pengadaan': 'Staff Pengadaan',
    'chef': 'Chef',
    'packing': 'Staff Packing',
    'driver': 'Driver',
    'asisten': 'Asisten Lapangan'
  }
  return roleLabels[authStore.user?.role] || 'User'
})

const welcomeIllustration = computed(() => {
  const role = authStore.user?.role
  const illustrationMap = {
    'ahli_gizi': welcomeChefSvg,
    'chef': welcomeChefSvg,
    'driver': welcomeDeliverySvg,
    'akuntan': welcomeAnalyticsSvg
  }
  return illustrationMap[role] || null
})

const goToKepalaSSPGDashboard = () => {
  redirecting.value = true
  router.push('/dashboard/kepala-sppg')
}

onMounted(() => {
  if (authStore.user?.role === 'kepala_sppg') {
    setTimeout(() => {
      goToKepalaSSPGDashboard()
    }, 2000)
  }
})
</script>

<style scoped>
.redirect-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 16px;
}

.redirect-text {
  font-size: 16px;
  font-weight: 500;
  color: #6B6B6B;
  margin: 0;
}

.welcome-card {
  text-align: center;
  padding: 20px;
}

.welcome-card__content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
}

.welcome-card__text {
  text-align: left;
}

.welcome-card__illustration {
  width: 80px;
  height: 80px;
  object-fit: contain;
  flex-shrink: 0;
}

.welcome-title {
  font-size: var(--h-text-2xl, 24px);
  font-weight: 600;
  color: var(--h-text-primary, #303030);
  margin: 0 0 var(--h-spacing-2, 8px) 0;
}

.welcome-subtitle {
  font-size: var(--h-text-sm, 14px);
  color: var(--h-text-secondary, #6B6B6B);
  margin: 0;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

@media (max-width: 1024px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-row {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}

.content-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

@media (max-width: 1024px) {
  .content-row {
    grid-template-columns: 1fr;
  }
}

.content-card {
  padding: var(--h-spacing-6, 24px);
}

.card-title {
  font-size: var(--h-text-lg, 18px);
  font-weight: 600;
  color: var(--h-text-primary, #303030);
  margin: 0 0 var(--h-spacing-4, 16px) 0;
}

.h-button {
  background: #C94A3A;
  border: none;
  border-radius: 6px;
  height: 44px;
  font-weight: var(--h-font-semibold, 600);
}

/* Overview Highlight Widget */
.overview-widget {
  background: #CCE2C8;
  border-radius: 8px;
  padding: 20px;
}

.overview-widget__title {
  font-size: 18px;
  font-weight: 600;
  color: #303030;
  margin: 0 0 16px 0;
}

.overview-widget__cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.overview-widget__card {
  background: rgba(255, 255, 255, 0.7);
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.overview-widget__label {
  font-size: 12px;
  color: #6B6B6B;
  font-weight: 500;
}

.overview-widget__value {
  font-size: 24px;
  font-weight: 600;
  color: #303030;
}

@media (max-width: 768px) {
  .overview-widget__cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
