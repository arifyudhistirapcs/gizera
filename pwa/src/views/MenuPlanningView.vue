<template>
  <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
    <div class="menu-planning-page">
      <!-- NavBar -->
      <van-nav-bar title="Perencanaan Menu" />

      <!-- Date/Week Selector -->
      <div class="date-section">
        <DateSelector v-model="selectedDate" @update:modelValue="onWeekChange" />
      </div>

      <!-- Loading State -->
      <template v-if="menuStore.loading">
        <div class="plans-list">
          <SkeletonCard v-for="n in 3" :key="n" :rows="3" />
        </div>
      </template>

      <!-- Error State -->
      <div v-else-if="menuStore.error && !menuStore.weeklyPlans.length" class="error-state">
        <van-icon name="warning-o" size="48" color="var(--h-error)" />
        <p class="error-state__message">{{ menuStore.error }}</p>
        <van-button type="primary" size="normal" @click="menuStore.retry()">
          Coba Lagi
        </van-button>
      </div>

      <!-- Weekly Plans List -->
      <template v-else>
        <div v-if="menuStore.weeklyPlans.length" class="plans-list">
          <div v-for="plan in menuStore.weeklyPlans" :key="plan.id">
            <MenuWeekCard
              :weekPeriod="plan.weekPeriod"
              :approvalStatus="plan.approvalStatus"
              :menuCount="plan.menuCount"
              :canApprove="isKepalaSPPG && plan.approvalStatus === 'pending'"
              @approve="handleApprove(plan.id)"
              @reject="showRejectDialog(plan.id)"
              @click="toggleDetail(plan.id)"
            />

            <!-- Expanded Detail: Daily Menus -->
            <transition name="expand">
              <div v-if="expandedPlanId === plan.id && plan.menus && plan.menus.length" class="daily-menus">
                <div v-for="menu in plan.menus" :key="menu.day" class="menu-day-card h-card">
                  <h4 class="menu-day-card__day">{{ menu.day }}</h4>
                  <p class="menu-day-card__name">{{ menu.menuName }}</p>
                  <p class="menu-day-card__components">
                    <van-icon name="label-o" size="14" color="var(--h-text-secondary)" />
                    Komponen: {{ Array.isArray(menu.components) ? menu.components.join(', ') : menu.components }}
                  </p>
                  <p class="menu-day-card__portions">
                    <van-icon name="friends-o" size="14" color="var(--h-text-secondary)" />
                    Porsi: {{ menu.portions }}
                  </p>
                </div>
              </div>
            </transition>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="!menuStore.weeklyPlans.length" class="empty-state">
          <van-icon name="notes-o" size="48" color="var(--h-text-light)" />
          <p class="empty-state__text">Tidak ada rencana menu ditemukan</p>
        </div>
      </template>

      <!-- Approval Loading Overlay -->
      <div v-if="menuStore.approvalLoading" class="approval-overlay">
        <van-loading type="spinner" color="var(--h-primary)" />
      </div>
    </div>
  </van-pull-refresh>

  <!-- Reject Dialog -->
  <van-dialog
    v-model:show="showReject"
    title="Tolak Menu"
    show-cancel-button
    confirm-button-text="Tolak"
    cancel-button-text="Batal"
    @confirm="handleReject"
    @cancel="closeRejectDialog"
  >
    <div class="reject-dialog__content">
      <van-field
        v-model="rejectReason"
        placeholder="Alasan penolakan"
        type="textarea"
        rows="3"
        autosize
      />
    </div>
  </van-dialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useMenuPlanningStore } from '@/stores/menuPlanning'
import { useAuthStore } from '@/stores/auth'
import DateSelector from '@/components/mobile/DateSelector.vue'
import MenuWeekCard from '@/components/mobile/MenuWeekCard.vue'
import SkeletonCard from '@/components/mobile/SkeletonCard.vue'

const menuStore = useMenuPlanningStore()
const authStore = useAuthStore()

const refreshing = ref(false)
const selectedDate = ref(new Date())
const expandedPlanId = ref(null)
const showReject = ref(false)
const rejectReason = ref('')
const rejectTargetId = ref(null)

// Role check — only Kepala_SPPG can approve/reject
const isKepalaSPPG = computed(() => authStore.user?.role?.toLowerCase() === 'kepala_sppg')

function onWeekChange(date) {
  // Convert date to ISO week string for the store
  const d = new Date(date)
  const year = d.getFullYear()
  const oneJan = new Date(year, 0, 1)
  const weekNum = Math.ceil(((d - oneJan) / 86400000 + oneJan.getDay() + 1) / 7)
  const weekStr = `${year}-W${String(weekNum).padStart(2, '0')}`
  menuStore.setWeek(weekStr)
}

function toggleDetail(planId) {
  expandedPlanId.value = expandedPlanId.value === planId ? null : planId
}

async function handleApprove(planId) {
  await menuStore.approveMenu(planId)
}

function showRejectDialog(planId) {
  rejectTargetId.value = planId
  rejectReason.value = ''
  showReject.value = true
}

function closeRejectDialog() {
  showReject.value = false
  rejectReason.value = ''
  rejectTargetId.value = null
}

async function handleReject() {
  if (rejectTargetId.value) {
    await menuStore.rejectMenu(rejectTargetId.value, rejectReason.value)
  }
  closeRejectDialog()
}

async function onRefresh() {
  await menuStore.fetchMenuPlans()
  refreshing.value = false
}

onMounted(() => {
  menuStore.fetchMenuPlans()
})
</script>

<style scoped>
.menu-planning-page {
  padding: var(--h-spacing-lg);
  padding-bottom: 80px;
  min-height: 100vh;
  position: relative;
}

/* Date Section */
.date-section {
  margin-bottom: var(--h-spacing-lg);
}

/* Plans List */
.plans-list {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-md);
}

/* Daily Menus (expanded detail) */
.daily-menus {
  display: flex;
  flex-direction: column;
  gap: var(--h-spacing-sm);
  margin-top: var(--h-spacing-sm);
  padding-left: var(--h-spacing-md);
}

.menu-day-card {
  background: var(--h-bg-card);
  border-radius: var(--h-radius-lg);
  box-shadow: var(--h-shadow-card);
  padding: var(--h-spacing-md);
  border-left: 3px solid var(--h-primary);
}

.menu-day-card__day {
  font-size: 14px;
  font-weight: 600;
  color: var(--h-primary);
  margin: 0 0 4px 0;
}

.menu-day-card__name {
  font-size: 15px;
  font-weight: 600;
  color: var(--h-text-primary);
  margin: 0 0 8px 0;
}

.menu-day-card__components,
.menu-day-card__portions {
  font-size: 13px;
  color: var(--h-text-secondary);
  margin: 0 0 4px 0;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Expand Transition */
.expand-enter-active,
.expand-leave-active {
  transition: all var(--h-transition-base);
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 1000px;
}

/* Approval Loading Overlay */
.approval-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.6);
  z-index: 10;
  border-radius: var(--h-radius-lg);
}

/* Error State */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px var(--h-spacing-xl);
  text-align: center;
}

.error-state__message {
  font-size: 14px;
  color: var(--h-text-secondary);
  margin: var(--h-spacing-lg) 0;
  line-height: 1.5;
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px var(--h-spacing-xl);
  text-align: center;
}

.empty-state__text {
  font-size: 14px;
  color: var(--h-text-light);
  margin: var(--h-spacing-md) 0 0 0;
}

/* Reject Dialog */
.reject-dialog__content {
  padding: var(--h-spacing-md);
}
</style>
