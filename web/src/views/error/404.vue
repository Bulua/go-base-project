<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
  <div class="bp-placeholder" style="margin: 80px auto; max-width: 480px; text-align: center">
    <h1>404</h1>
    <p>页面不存在或您没有访问权限。</p>
    <el-button type="primary" style="margin-top: 16px" @click="router.push('/')">返回首页</el-button>
  </div>
</template>
