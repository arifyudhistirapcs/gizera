<template>
  <div class="horizon-login">
    <!-- Left Side: Login Form -->
    <div class="horizon-login__form-side">
      <div class="horizon-login__form-container">
        <!-- Mobile Logo (visible only when branding side is hidden) -->
        <div class="horizon-login__mobile-logo">
          <img src="/gizera-light.png" alt="Dapur Sehat" class="horizon-login__mobile-logo-img" />
        </div>

        <!-- Title -->
        <div class="horizon-login__header">
          <h1 class="horizon-login__title">Masuk ke Dapur Sehat</h1>
          <p class="horizon-login__subtitle">Kelola operasional dapur dengan mudah</p>
        </div>

        <!-- Form -->
        <a-form
          :model="formState"
          :rules="rules"
          @finish="handleLogin"
          layout="vertical"
          class="horizon-login__form"
        >
          <!-- NIK / Email Input -->
          <a-form-item
            label="NIK / Email"
            name="identifier"
            :validate-status="error ? 'error' : ''"
          >
            <a-input
              v-model:value="formState.identifier"
              placeholder="Masukkan NIK atau Email"
              class="horizon-login__input"
              :disabled="loading"
            >
              <template #prefix>
                <UserOutlined class="horizon-login__input-icon" />
              </template>
            </a-input>
          </a-form-item>

          <!-- Password Input -->
          <a-form-item
            label="Password"
            name="password"
            :validate-status="error ? 'error' : ''"
            :help="error"
          >
            <a-input-password
              v-model:value="formState.password"
              placeholder="Masukkan password"
              class="horizon-login__input"
              :disabled="loading"
            >
              <template #prefix>
                <LockOutlined class="horizon-login__input-icon" />
              </template>
            </a-input-password>
          </a-form-item>

          <!-- Keep me logged in -->
          <div class="horizon-login__options">
            <a-checkbox v-model:checked="formState.rememberMe">
              Keep me logged in
            </a-checkbox>
          </div>

          <!-- Submit Button -->
          <a-button
            type="primary"
            html-type="submit"
            class="horizon-login__submit"
            :loading="loading"
          >
            Masuk
          </a-button>
        </a-form>

        <!-- Footer -->
        <div class="horizon-login__footer">
          <p class="horizon-login__footer-text">
            Dapur Sehat &mdash; Sistem Manajemen Operasional
          </p>
        </div>
      </div>
    </div>

    <!-- Right Side: Branding -->
    <div class="horizon-login__brand-side">
      <!-- Decorative background circles -->
      <div class="horizon-login__brand-bg-circle horizon-login__brand-bg-circle--top" />
      <div class="horizon-login__brand-bg-circle horizon-login__brand-bg-circle--bottom" />

      <div class="horizon-login__brand-content">
        <!-- Dapur Sehat Logo -->
        <div class="horizon-login__brand-logo-wrap">
          <img src="/gizera-light.png" alt="Dapur Sehat" class="horizon-login__brand-gizera-logo" />
        </div>

        <!-- Lottie Animation as Centerpiece -->
        <div class="horizon-login__brand-lottie">
          <LottiePlayer src="/lottie/data-flow.json" width="300px" height="300px" :loop="true" />
        </div>

        <!-- SVG Illustration -->
        <div class="horizon-login__brand-illustration-wrap">
          <img src="@/assets/illustrations/login-branding.svg" alt="Dapur Sehat Branding" class="horizon-login__brand-illustration" />
        </div>

        <!-- Tagline -->
        <h2 class="horizon-login__brand-title">Dapur Sehat</h2>
        <p class="horizon-login__brand-subtitle">
          Sistem Manajemen Operasional Terpadu
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import LottiePlayer from '@/components/common/LottiePlayer.vue'

const router = useRouter()
const authStore = useAuthStore()

const formState = reactive({
  identifier: '',
  password: '',
  rememberMe: false
})

const loading = ref(false)
const error = ref(null)

const rules = {
  identifier: [
    { required: true, message: 'NIK atau Email harus diisi', trigger: 'blur' }
  ],
  password: [
    { required: true, message: 'Password harus diisi', trigger: 'blur' },
    { min: 6, message: 'Password minimal 6 karakter', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  loading.value = true
  error.value = null

  try {
    await authStore.login({
      identifier: formState.identifier,
      password: formState.password
    })

    message.success('Login berhasil!')
    
    // Redirect based on user role
    const user = authStore.user
    const role = user?.role || 'default'
    
    const roleRoutes = {
      superadmin: '/yayasan',
      admin_bgn: '/dashboard-bgn',
      kepala_yayasan: '/dashboard-yayasan',
      kepala_sppg: '/dashboard/kepala-sppg',
      ahli_gizi: '/menu-planning',
      pengadaan: '/purchase-orders',
      akuntan: '/financial-reports',
      chef: '/kds/cooking',
      packing: '/kds/cooking'
    }
    
    const target = roleRoutes[role] || '/dashboard'
    
    // Use window.location for reliable redirect after login
    window.location.href = target
  } catch (err) {
    console.error('Login error:', err)
    error.value = err.response?.data?.message || 'Login gagal. Periksa NIK/Email dan password Anda.'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* ========================================
   MAIN CONTAINER - Split Screen Layout
   ======================================== */

.horizon-login {
  display: flex;
  min-height: 100vh;
  background-color: var(--h-bg-secondary, #FFFFFF);
}

/* ========================================
   LEFT SIDE - Login Form
   ======================================== */

.horizon-login__form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--h-spacing-12, 48px);
  background-color: var(--h-bg-secondary, #FFFFFF);
}

.horizon-login__form-container {
  width: 100%;
  max-width: 480px;
}

/* Mobile Logo — hidden on desktop, shown when branding side is hidden */
.horizon-login__mobile-logo {
  display: none;
  text-align: center;
  margin-bottom: 32px;
}

.horizon-login__mobile-logo-img {
  height: 48px;
  width: auto;
  object-fit: contain;
}

/* Header */
.horizon-login__header {
  margin-bottom: var(--h-spacing-8, 32px);
}

.horizon-login__title {
  font-size: 34px;
  font-weight: 700;
  color: var(--h-text-primary, #303030);
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.horizon-login__subtitle {
  font-size: 15px;
  color: var(--h-text-secondary, #6B6B6B);
  margin: 0;
  line-height: 1.5;
}

/* Form */
.horizon-login__form {
  margin-bottom: var(--h-spacing-6, 24px);
}

.horizon-login__form :deep(.ant-form-item-label > label) {
  font-size: 14px;
  font-weight: 500;
  color: var(--h-text-primary, #303030);
}

.horizon-login__input {
  height: 48px;
  border-radius: 12px;
  font-size: 15px;
  border: 1.5px solid var(--h-border-color, #D8D8DB);
  transition: all 200ms ease;
}

.horizon-login__input:hover {
  border-color: #C94A3A;
}

.horizon-login__input:focus,
.horizon-login__input:focus-within {
  border-color: #C94A3A;
  box-shadow: 0 0 0 3px rgba(201, 74, 58, 0.1);
}

.horizon-login__input :deep(.ant-input) {
  font-size: 15px;
}

.horizon-login__input-icon {
  color: var(--h-text-secondary, #6B6B6B);
  font-size: 18px;
}

/* Options (Keep me logged in) */
.horizon-login__options {
  margin-bottom: 24px;
}

.horizon-login__options :deep(.ant-checkbox-wrapper) {
  font-size: 14px;
  color: var(--h-text-secondary, #6B6B6B);
}

/* Submit Button */
.horizon-login__submit {
  width: 100%;
  height: 52px;
  border-radius: 14px;
  font-size: 16px;
  font-weight: 600;
  background: #C94A3A;
  border: none;
  transition: all 200ms ease;
  box-shadow: 0 4px 14px rgba(201, 74, 58, 0.3);
}

.horizon-login__submit:hover {
  background: #B5412F;
  box-shadow: 0 6px 20px rgba(201, 74, 58, 0.4);
  transform: translateY(-1px);
}

.horizon-login__submit:active {
  background: #9C3828;
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(201, 74, 58, 0.3);
}

/* Footer */
.horizon-login__footer {
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid var(--h-border-color, #E8E8EB);
}

.horizon-login__footer-text {
  font-size: 13px;
  color: var(--h-text-secondary, #9B9B9B);
  margin: 0;
  letter-spacing: 0.02em;
}

/* ========================================
   RIGHT SIDE - Branding
   ======================================== */

.horizon-login__brand-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(160deg, #C94A3A 0%, #A33D30 40%, #1E8A6E 100%);
  padding: 48px;
  position: relative;
  overflow: hidden;
}

/* Decorative background circles */
.horizon-login__brand-bg-circle {
  position: absolute;
  border-radius: 50%;
  pointer-events: none;
}

.horizon-login__brand-bg-circle--top {
  width: 400px;
  height: 400px;
  top: -120px;
  right: -100px;
  background: rgba(255, 255, 255, 0.06);
}

.horizon-login__brand-bg-circle--bottom {
  width: 300px;
  height: 300px;
  bottom: -80px;
  left: -60px;
  background: rgba(255, 255, 255, 0.04);
}

.horizon-login__brand-content {
  text-align: center;
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

/* Gizera Logo */
.horizon-login__brand-logo-wrap {
  margin-bottom: 8px;
}

.horizon-login__brand-gizera-logo {
  height: 72px;
  width: auto;
  object-fit: contain;
  filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.15));
}

/* Lottie Centerpiece */
.horizon-login__brand-lottie {
  display: flex;
  align-items: center;
  justify-content: center;
  filter: drop-shadow(0 8px 24px rgba(0, 0, 0, 0.1));
}

/* SVG Illustration */
.horizon-login__brand-illustration-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.horizon-login__brand-illustration {
  width: 300px;
  height: auto;
  max-height: 300px;
  object-fit: contain;
  filter: drop-shadow(0 4px 16px rgba(0, 0, 0, 0.1));
}

/* Brand Title */
.horizon-login__brand-title {
  font-size: 36px;
  font-weight: 700;
  color: #FFFFFF;
  margin: 0;
  line-height: 1.2;
  letter-spacing: -0.01em;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.horizon-login__brand-subtitle {
  font-size: 17px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
  line-height: 1.5;
  font-weight: 400;
  letter-spacing: 0.01em;
}

/* ========================================
   RESPONSIVE - Tablet & Mobile
   ======================================== */

/* Tablet (< 1024px): hide branding, show mobile logo */
@media (max-width: 1024px) {
  .horizon-login__brand-side {
    display: none;
  }

  .horizon-login__form-side {
    flex: 1;
  }

  .horizon-login__mobile-logo {
    display: block;
  }
}

/* Mobile (< 768px) */
@media (max-width: 768px) {
  .horizon-login__form-side {
    padding: 24px;
  }

  .horizon-login__title {
    font-size: 28px;
  }

  .horizon-login__form-container {
    max-width: 100%;
  }

  .horizon-login__submit {
    height: 48px;
    border-radius: 12px;
  }
}

/* Small Mobile (< 480px) */
@media (max-width: 480px) {
  .horizon-login__form-side {
    padding: 16px;
  }

  .horizon-login__title {
    font-size: 24px;
  }

  .horizon-login__input,
  .horizon-login__submit {
    height: 44px;
  }

  .horizon-login__mobile-logo-img {
    height: 40px;
  }
}

/* ========================================
   DARK MODE SUPPORT
   ======================================== */

.dark .horizon-login__form-side {
  background-color: #1A1A1A;
}

.dark .horizon-login__title {
  color: var(--h-text-primary-dark, #F7F8FA);
}

.dark .horizon-login__subtitle,
.dark .horizon-login__footer-text {
  color: var(--h-text-secondary-dark, #D8D8DB);
}

.dark .horizon-login__input {
  background-color: #252525;
  border-color: #404040;
  color: #F7F8FA;
}

.dark .horizon-login__input::placeholder {
  color: #D8D8DB;
}

.dark .horizon-login__footer {
  border-top-color: #404040;
}
</style>
