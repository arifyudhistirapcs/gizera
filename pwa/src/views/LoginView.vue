<template>
  <div class="login-container">
    <!-- Gradient Header -->
    <div class="login-hero">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <img src="@/logo/dapur-sehat-light.png" alt="Dapur Sehat" class="hero-logo" />
      </div>
      <!-- Curved bottom -->
      <svg class="hero-curve" viewBox="0 0 400 40" preserveAspectRatio="none">
        <path d="M0,40 L0,15 Q200,-15 400,15 L400,40 Z" fill="#F7F8FA"/>
      </svg>
    </div>

    <!-- Login Form -->
    <div class="login-body">
      <van-form @submit.prevent="handleLogin" class="login-form">
        <div class="login-card">
          <h2 class="card-title">Masuk ke Akun</h2>
          <p class="card-subtitle">Silakan masukkan kredensial Anda</p>

          <!-- Error Alert -->
          <div v-if="errorMessage" class="error-alert">
            <van-icon name="warning-o" />
            <span>{{ errorMessage }}</span>
          </div>
          
          <van-field
            v-model="formData.identifier"
            name="identifier"
            label="NIK/Email"
            placeholder="Masukkan NIK atau Email"
            :rules="[{ required: true, message: 'NIK/Email wajib diisi' }]"
            left-icon="user-o"
            clearable
            :disabled="loading"
            @focus="errorMessage = ''"
          />
          <van-field
            v-model="formData.password"
            :type="showPassword ? 'text' : 'password'"
            name="password"
            label="Password"
            placeholder="Masukkan password"
            :rules="[{ required: true, message: 'Password wajib diisi' }]"
            left-icon="lock"
            :right-icon="showPassword ? 'eye-o' : 'closed-eye'"
            :disabled="loading"
            @click-right-icon="showPassword = !showPassword"
            @focus="errorMessage = ''"
          />

          <div class="login-btn-wrapper">
            <van-button
              block
              type="primary"
              native-type="button"
              :loading="loading"
              loading-text="Memproses..."
              :disabled="loading"
              class="login-btn"
              @click="handleLogin"
            >
              Login
            </van-button>
          </div>
        </div>
      </van-form>

      <div class="login-footer">
        <p>© 2026 Dapur Sehat</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authAPI } from '@/services/api'
import { showToast, showNotify, showDialog } from 'vant'
import LottiePlayer from '@/components/common/LottiePlayer.vue'


const router = useRouter()
const authStore = useAuthStore()

const formData = ref({
  identifier: '',
  password: ''
})

const loading = ref(false)
const showPassword = ref(false)
const errorMessage = ref('')

const handleLogin = async () => {
  if (!formData.value.identifier || !formData.value.password) {
    errorMessage.value = 'NIK/Email dan Password wajib diisi'
    return
  }

  loading.value = true
  errorMessage.value = ''
  
  try {
    const response = await authAPI.login({
      identifier: formData.value.identifier,
      password: formData.value.password
    })

    if (response.data.success) {
      const userData = {
        id: response.data.user.id,
        nik: response.data.user.nik,
        email: response.data.user.email,
        name: response.data.user.full_name,
        username: response.data.user.nik,
        role: response.data.user.role
      }
      
      authStore.setAuth(userData, response.data.token)
      
      showToast({
        type: 'success',
        message: 'Login berhasil!',
        duration: 2000
      })

      setTimeout(() => {
        router.push('/')
      }, 500)
    } else {
      errorMessage.value = response.data.message || 'Login gagal'
      showDialog({
        title: 'Login Gagal',
        message: errorMessage.value,
        confirmButtonText: 'OK',
        confirmButtonColor: '#C94A3A'
      })
    }
  } catch (error) {
    console.error('Login error:', error)
    
    let errMsg = 'Terjadi kesalahan saat login'
    
    if (error.response) {
      if (error.response.status === 401) {
        errMsg = 'NIK/Email atau password salah'
      } else if (error.response.status === 403) {
        errMsg = 'Akun Anda tidak aktif. Silakan hubungi administrator'
      } else if (error.response.data?.message) {
        errMsg = error.response.data.message
      }
    } else if (error.request) {
      errMsg = 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.'
    }
    
    errorMessage.value = errMsg
    showDialog({
      title: 'Login Gagal',
      message: errMsg,
      confirmButtonText: 'OK',
      confirmButtonColor: '#C94A3A'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: #F7F8FA;
  display: flex;
  flex-direction: column;
}

/* ========== HERO GRADIENT HEADER ========== */
.login-hero {
  position: relative;
  min-height: 280px;
  overflow: hidden;
  flex-shrink: 0;
}

.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #C94A3A 0%, #D4553E 45%, #1E8A6E 100%);
}

.hero-content {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px 56px;
}

.hero-logo {
  width: 340px;
  max-width: 80vw;
  height: auto;
  filter: brightness(0) invert(1);
  margin-bottom: 8px;
}

.hero-curve {
  position: absolute;
  bottom: -1px;
  left: 0;
  width: 100%;
  height: 40px;
  z-index: 3;
}

/* ========== LOGIN BODY ========== */
.login-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 20px 24px;
  margin-top: -12px;
  position: relative;
  z-index: 4;
}

.login-form {
  width: 100%;
  max-width: 400px;
}

.login-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  padding: 28px 24px;
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.card-title {
  font-size: 22px;
  font-weight: 700;
  color: #303030;
  margin: 0 0 4px;
}

.card-subtitle {
  font-size: 14px;
  color: #8C8C8C;
  margin: 0 0 24px;
}

/* Error */
.error-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(201, 74, 58, 0.08);
  border: 1px solid rgba(201, 74, 58, 0.2);
  border-radius: 10px;
  margin-bottom: 16px;
  color: #C94A3A;
  font-size: 14px;
  font-weight: 500;
}

.error-alert .van-icon {
  font-size: 18px;
  flex-shrink: 0;
}

/* Form fields */
.login-card :deep(.van-field) {
  padding: 14px 0;
}

.login-card :deep(.van-cell::after) {
  left: 0;
  right: 0;
  border-color: #E8E8E8;
}

.login-card :deep(.van-field__label) {
  color: #303030;
  font-weight: 600;
  font-size: 14px;
}

.login-card :deep(.van-field__control) {
  color: #303030;
  font-size: 15px;
}

.login-card :deep(.van-field__control::placeholder) {
  color: #BFBFBF;
}

.login-card :deep(.van-field__left-icon) {
  color: #1E8A6E;
  font-size: 18px;
}

.login-card :deep(.van-field__right-icon) {
  color: #8C8C8C;
}

/* Button */
.login-btn-wrapper {
  margin-top: 28px;
}

.login-btn {
  height: 50px !important;
  font-size: 16px !important;
  font-weight: 700 !important;
  border-radius: 12px !important;
  background: linear-gradient(135deg, #C94A3A 0%, #1E8A6E 100%) !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba(201, 74, 58, 0.3);
  transition: all 0.2s ease;
  letter-spacing: 0.3px;
}

.login-btn:active {
  transform: scale(0.97);
  box-shadow: 0 2px 8px rgba(201, 74, 58, 0.2);
}

/* Footer */
.login-footer {
  margin-top: auto;
  text-align: center;
  color: #BFBFBF;
  font-size: 12px;
  padding: 20px 0;
}

.login-footer p {
  margin: 0;
}
</style>
