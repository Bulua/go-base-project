<script setup lang="ts">
import { reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Key, Lock, User } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()

const loginForm = reactive({
  login_name: 'admin',
  password: '123456',
})

async function handleLogin() {
  if (!loginForm.login_name.trim() || !loginForm.password) {
    ElMessage.warning('请输入账号和密码')
    return
  }
  await authStore.loginWithPassword(loginForm.login_name, loginForm.password)
}
</script>

<template>
  <main class="login-page">
    <section class="login-panel">
      <div class="login-brand">
        <span class="brand-mark">G</span>
        <div>
          <strong>GoBaseProject</strong>
          <small>后台认证入口</small>
        </div>
      </div>

      <el-form label-position="top" @submit.prevent>
        <el-form-item label="账号">
          <el-input
            v-model="loginForm.login_name"
            autocomplete="username"
            size="large"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="密码">
          <el-input
            v-model="loginForm.password"
            autocomplete="current-password"
            show-password
            size="large"
            type="password"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-button
          class="login-button"
          :icon="Key"
          :loading="authStore.loginLoading"
          size="large"
          type="primary"
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form>
    </section>
  </main>
</template>
