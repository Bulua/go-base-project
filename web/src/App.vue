<script setup lang="ts">
// Theme is initialised on import so that <html data-theme> and --hue
// reflect persisted settings before the first paint.
import '@/composables/common/useTheme'
</script>

<template>
  <router-view />
</template>
