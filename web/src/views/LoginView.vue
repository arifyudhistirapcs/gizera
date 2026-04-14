<template>
  <div class="login-page">
    <!-- Left: Form -->
    <div class="login-form-side">
      <div class="login-form-wrap">
        <!-- Mobile logo -->
        <div class="mobile-logo">
          <img src="/pose-logo.svg" alt="POSe" />
        </div>

        <div class="form-header">
          <h1>Selamat Datang 👋</h1>
          <p>Masuk untuk mengelola operasional dapur</p>
        </div>

        <a-form
          :model="formState"
          :rules="rules"
          @finish="handleLogin"
          layout="vertical"
          class="login-form"
        >
          <a-form-item label="NIK / Email" name="identifier" :validate-status="error ? 'error' : ''">
            <a-input
              v-model:value="formState.identifier"
              placeholder="Masukkan NIK atau Email"
              size="large"
              :disabled="loading"
            >
              <template #prefix>
                <UserOutlined style="color: #8C8C8C" />
              </template>
            </a-input>
          </a-form-item>

          <a-form-item label="Password" name="password" :validate-status="error ? 'error' : ''" :help="error">
            <a-input-password
              v-model:value="formState.password"
              placeholder="Masukkan password"
              size="large"
              :disabled="loading"
            >
              <template #prefix>
                <LockOutlined style="color: #8C8C8C" />
              </template>
            </a-input-password>
          </a-form-item>

          <div class="form-options">
            <a-checkbox v-model:checked="formState.rememberMe">Ingat saya</a-checkbox>
          </div>

          <a-button type="primary" html-type="submit" block size="large" class="login-btn" :loading="loading">
            Masuk
          </a-button>
        </a-form>

        <div class="form-footer">
          <p>© 2026 POSe</p>
        </div>
      </div>
    </div>

    <!-- Right: Branding -->
    <div class="login-brand-side">
      <div class="brand-circles">
        <div class="circle circle-1"></div>
        <div class="circle circle-2"></div>
        <div class="circle circle-3"></div>
      </div>
      <div class="brand-content">
        <img src="/pose-logo-white.svg" alt="POSe" class="brand-logo" />
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
  identifier: [{ required: true, message: 'NIK atau Email harus diisi', trigger: 'blur' }],
  password: [
    { required: true, message: 'Password harus diisi', trigger: 'blur' },
    { min: 6, message: 'Password minimal 6 karakter', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  loading.value = true
  error.value = null
  try {
    await authStore.login({ identifier: formState.identifier, password: formState.password })
    message.success('Login berhasil!')
    const role = authStore.user?.role || 'default'
    const roleRoutes = {
      superadmin: '/yayasan',
      admin_bgn: '/dashboard-bgn',
      kepala_yayasan: '/dashboard-yayasan',
      kepala_sppg: '/dashboard/kepala-sppg',
      ahli_gizi: '/menu-planning',
      pengadaan: '/goods-receipts',
      akuntan: '/financial-reports',
      chef: '/kds/cooking',
      packing: '/kds/cooking',
      supplier: '/supplier-dashboard'
    }
    window.location.href = roleRoutes[role] || '/dashboard'
  } catch (err) {
    error.value = err.response?.data?.message || 'Login gagal. Periksa NIK/Email dan password Anda.'
    message.error(error.value)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  min-height: 100vh;
}

/* ===== LEFT: FORM ===== */
.login-form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
  background: #fff;
}

.login-form-wrap {
  width: 100%;
  max-width: 420px;
}

.mobile-logo {
  display: none;
  margin-bottom: 40px;
  text-align: center;
}

.mobile-logo img {
  height: 56px;
  width: auto;
  display: block;
  margin: 0 auto 4px;
}

.form-header {
  margin-bottom: 36px;
}

.form-header h1 {
  font-size: 32px;
  font-weight: 800;
  color: #1a1a1a;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.form-header p {
  font-size: 15px;
  color: #8C8C8C;
  margin: 0;
}

/* Form fields */
.login-form :deep(.ant-form-item) {
  margin-bottom: 24px;
}

.login-form :deep(.ant-form-item-label) {
  padding-bottom: 8px;
}

.login-form :deep(.ant-form-item-label > label) {
  font-size: 14px;
  font-weight: 600;
  color: #303030;
  letter-spacing: 0;
  text-transform: none;
}

.login-form :deep(.ant-input-affix-wrapper) {
  border-radius: 12px;
  border: 2px solid #EBEBEB;
  background: #FAFAFA;
  padding: 0 16px;
  height: 52px;
  transition: all 0.25s ease;
}

.login-form :deep(.ant-input-affix-wrapper:hover) {
  border-color: #D0D0D0;
  background: #fff;
}

.login-form :deep(.ant-input-affix-wrapper-focused),
.login-form :deep(.ant-input-affix-wrapper:focus) {
  border-color: #F82C17;
  background: #fff;
  box-shadow: 0 0 0 4px rgba(248, 44, 23, 0.06);
}

.login-form :deep(.ant-input) {
  font-size: 15px;
  background: transparent;
}

.login-form :deep(.ant-input::placeholder) {
  color: #BFBFBF;
}

.login-form :deep(.ant-input-prefix) {
  margin-right: 12px;
  font-size: 18px;
}

.login-form :deep(.ant-input-suffix) {
  font-size: 16px;
  color: #BFBFBF;
}

.login-form :deep(.ant-input-suffix:hover) {
  color: #666;
}

.form-options {
  margin-bottom: 28px;
}

.form-options :deep(.ant-checkbox-wrapper) {
  font-size: 14px;
  color: #999;
}

.form-options :deep(.ant-checkbox-checked .ant-checkbox-inner) {
  background-color: #F82C17;
  border-color: #F82C17;
}

.login-btn {
  height: 52px !important;
  border-radius: 12px !important;
  font-size: 16px !important;
  font-weight: 700 !important;
  background: #F82C17 !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba(248, 44, 23, 0.2);
  transition: all 0.25s ease;
  letter-spacing: 0.3px;
}

.login-btn:hover {
  background: #E02510 !important;
  box-shadow: 0 6px 24px rgba(248, 44, 23, 0.3) !important;
  transform: translateY(-1px);
}

.login-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(248, 44, 23, 0.2) !important;
}

.form-footer {
  text-align: center;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid #F0F0F0;
}

.form-footer p {
  font-size: 13px;
  color: #BFBFBF;
  margin: 0;
}

/* ===== RIGHT: BRANDING ===== */
.login-brand-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F82C17;
  position: relative;
  overflow: hidden;
}

.brand-circles .circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
}

.circle-1 {
  width: 500px;
  height: 500px;
  top: -150px;
  right: -120px;
}

.circle-2 {
  width: 350px;
  height: 350px;
  bottom: -100px;
  left: -80px;
}

.circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 15%;
  background: rgba(255, 255, 255, 0.04) !important;
}

.brand-content {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.brand-logo {
  width: 50%;
  min-width: 280px;
  max-width: 400px;
  height: auto;
}

/* ===== RESPONSIVE ===== */
@media (max-width: 1024px) {
  .login-brand-side { display: none; }
  .mobile-logo { display: block; }
}

@media (max-width: 768px) {
  .login-form-side { padding: 24px; }
  .form-header h1 { font-size: 26px; }
}

@media (max-width: 480px) {
  .login-form-side { padding: 20px; }
  .form-header h1 { font-size: 24px; }
  .mobile-logo img { height: 44px; }
}

/* ===== DARK MODE ===== */
.dark .login-form-side { background: #1A1A1A; }
.dark .form-header h1 { color: #F7F8FA; }
.dark .form-header p { color: #999; }
.dark .login-form :deep(.ant-form-item-label > label) { color: #ccc; }
.dark .login-form :deep(.ant-input-affix-wrapper) { background: #252525; border-color: #3A3A3A; }
.dark .login-form :deep(.ant-input-affix-wrapper:hover) { border-color: #555; background: #2A2A2A; }
.dark .login-form :deep(.ant-input-affix-wrapper-focused) { border-color: #F82C17; background: #2A2A2A; }
.dark .login-form :deep(.ant-input) { color: #F7F8FA; }
.dark .login-form :deep(.ant-input::placeholder) { color: #666; }
.dark .form-footer { border-top-color: #333; }
</style>
