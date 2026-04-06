<template>
  <div class="login-container">
    <div class="login-header">
      <div class="logo">
        <img src="@/logo/gizera-dark-mobile2.png" alt="GIZERA Logo" class="logo-image" />
      </div>
      <img src="@/assets/illustrations/login-header.svg" alt="Login illustration" class="login-header__illustration" />
    </div>

    <van-form @submit.prevent="handleLogin" class="login-form">
      <div class="login-card">
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
      <p>© 2024 SPPG. All rights reserved.</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { authAPI } from '@/services/api'
import { showToast, showNotify, showDialog } from 'vant'

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
  // Validate form first
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
      // Set auth data - backend returns token and user directly
      const userData = {
        id: response.data.user.id,
        nik: response.data.user.nik,
        email: response.data.user.email,
        name: response.data.user.full_name,
        username: response.data.user.nik, // Use NIK as username
        role: response.data.user.role
      }
      
      authStore.setAuth(userData, response.data.token)
      
      showToast({
        type: 'success',
        message: 'Login berhasil!',
        duration: 2000
      })

      // Redirect to role-based default page
      setTimeout(() => {
        router.push('/')
      }, 500)
    } else {
      errorMessage.value = response.data.message || 'Login gagal'
      showDialog({
        title: 'Login Gagal',
        message: errorMessage.value,
        confirmButtonText: 'OK',
        confirmButtonColor: '#EE5D50'
      })
    }
  } catch (error) {
    console.error('Login error:', error)
    
    let errMsg = 'Terjadi kesalahan saat login'
    
    if (error.response) {
      // Server responded with error
      if (error.response.status === 401) {
        errMsg = 'NIK/Email atau password salah'
      } else if (error.response.status === 403) {
        errMsg = 'Akun Anda tidak aktif. Silakan hubungi administrator'
      } else if (error.response.data?.message) {
        errMsg = error.response.data.message
      }
    } else if (error.request) {
      // Request made but no response
      errMsg = 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.'
    }
    
    errorMessage.value = errMsg
    showDialog({
      title: 'Login Gagal',
      message: errMsg,
      confirmButtonText: 'OK',
      confirmButtonColor: '#EE5D50'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  background: #E8EDE5;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--h-spacing-xl);
  position: relative;
  overflow: hidden;
}

.login-header {
  text-align: center;
  color: #303030;
  margin-top: 60px;
  margin-bottom: 40px;
  position: relative;
  z-index: 1;
}

.logo {
  margin-bottom: var(--h-spacing-xl);
}

.logo-image {
  width: 400px;
  max-width: 90vw;
  height: auto;
}

.login-header__illustration {
  width: 200px;
  max-width: 60vw;
  height: auto;
  margin: 0 auto;
  display: block;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

.login-header h1 {
  font-size: 28px;
  font-weight: 600;
  margin: 10px 0;
  font-family: var(--h-font-family);
}

.login-header p {
  font-size: 14px;
  opacity: 0.9;
  margin: 5px 0;
  font-weight: 400;
}

.login-form {
  width: 100%;
  max-width: 400px;
  position: relative;
  z-index: 1;
}

.login-card {
  background: var(--h-bg-card);
  border-radius: 8px;
  border: 1px solid #D8D8DB;
  box-shadow: none;
  padding: var(--h-spacing-2xl);
  overflow: hidden;
}

.error-alert {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(238, 93, 80, 0.1);
  border: 1px solid rgba(238, 93, 80, 0.3);
  border-radius: var(--h-radius-md);
  margin-bottom: 16px;
  color: #EE5D50;
  font-size: 14px;
  font-weight: 500;
}

.error-alert .van-icon {
  font-size: 18px;
  flex-shrink: 0;
}

.login-card :deep(.van-field) {
  padding: 12px 0;
}

.login-card :deep(.van-cell::after) {
  left: 0;
  right: 0;
  border-color: var(--h-border-color);
}

.login-card :deep(.van-field__label) {
  color: var(--h-text-primary);
  font-weight: 500;
  font-size: 14px;
}

.login-card :deep(.van-field__control) {
  color: var(--h-text-primary);
}

.login-card :deep(.van-field__control::placeholder) {
  color: var(--h-text-light);
}

.login-card :deep(.van-field__left-icon) {
  color: var(--h-text-secondary);
}

.login-card :deep(.van-field__right-icon) {
  color: var(--h-text-secondary);
}

.login-btn-wrapper {
  margin-top: var(--h-spacing-2xl);
}

.login-btn {
  height: 48px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  border-radius: 6px !important;
  background: #C94A3A !important;
  border-color: #C94A3A !important;
  box-shadow: none;
  transition: all var(--h-transition-base);
}

.login-btn:active {
  transform: scale(0.98);
  background: #A33D30 !important;
  border-color: #A33D30 !important;
}

.login-footer {
  margin-top: auto;
  text-align: center;
  color: #6B6B6B;
  opacity: 0.8;
  font-size: 12px;
  padding: var(--h-spacing-xl) 0;
  position: relative;
  z-index: 1;
}
</style>
