<script setup lang="ts">
import { useAuthStore } from '@/store/modules/auth'

const authStore = useAuthStore()
</script>

<template>
  <div class="bp-placeholder">
    <h1>仪表盘 · 工作台</h1>
    <p>欢迎回来，{{ authStore.currentUser?.display_name }}。该页面将在后续迭代中接入数据看板。</p>
  </div>
</template>
