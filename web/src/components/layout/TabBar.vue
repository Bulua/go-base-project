<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close } from '@element-plus/icons-vue'
import { useTabsStore } from '@/store/modules/tabs'

const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()

const activePath = computed(() => route.path)

watch(
  () => route.path,
  () => {
    const name = (route.name as string) || route.path
    const title = (route.meta?.title as string) || name
    tabsStore.add({
      path: route.path,
      name,
      title,
      keepAlive: !!route.meta?.keepAlive,
      affix: !!route.meta?.affix,
    })
  },
  { immediate: true },
)

function switchTab(path: string) {
  if (path !== route.path) router.push(path)
}

function closeTab(e: MouseEvent, path: string) {
  e.stopPropagation()
  const nextPath = tabsStore.close(path)
  if (path === route.path && nextPath !== path) {
    router.push(nextPath ?? '/')
  }
}

// ── Context menu ─────────────────────────────────────────────────────────────
const ctxVisible = ref(false)
const ctxX = ref(0)
const ctxY = ref(0)
const ctxPath = ref('')

const ctxTabAffix = computed(
  () => tabsStore.list.find((t) => t.path === ctxPath.value)?.affix ?? false,
)

function showCtx(e: MouseEvent, path: string) {
  e.preventDefault()
  ctxPath.value = path
  ctxX.value = e.clientX
  ctxY.value = e.clientY
  ctxVisible.value = true
}

function hideCtx() {
  ctxVisible.value = false
}

function ctxCloseCurrent() {
  if (ctxTabAffix.value) return
  closeTab(new MouseEvent('click'), ctxPath.value)
  hideCtx()
}

function ctxCloseOthers() {
  tabsStore.closeOthers(ctxPath.value)
  if (route.path !== ctxPath.value) router.push(ctxPath.value)
  hideCtx()
}

function ctxCloseAll() {
  tabsStore.closeAll()
  const first = tabsStore.list[0]
  router.push(first?.path ?? '/')
  hideCtx()
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') hideCtx()
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="tab-bar">
    <div class="tab-list">
      <div
        v-for="tab in tabsStore.list"
        :key="tab.path"
        class="tab-item"
        :class="{ active: tab.path === activePath, affix: tab.affix }"
        @click="switchTab(tab.path)"
        @contextmenu.prevent="showCtx($event, tab.path)"
      >
        <span class="tab-title">{{ tab.title }}</span>
        <el-icon v-if="!tab.affix" class="tab-close" @click="closeTab($event, tab.path)">
          <Close />
        </el-icon>
      </div>
    </div>
  </div>

  <!-- Right-click context menu -->
  <teleport to="body">
    <div v-if="ctxVisible" class="tab-ctx-mask" @click="hideCtx" @contextmenu.prevent />
    <ul v-if="ctxVisible" class="tab-ctx-menu" :style="{ left: ctxX + 'px', top: ctxY + 'px' }">
      <li :class="{ disabled: ctxTabAffix }" @click="ctxCloseCurrent">关闭当前</li>
      <li @click="ctxCloseOthers">关闭其他</li>
      <li @click="ctxCloseAll">关闭全部</li>
    </ul>
  </teleport>
</template>

<style scoped>
.tab-bar {
  height: 36px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: stretch;
  flex-shrink: 0;
  overflow: hidden;
}

.tab-list {
  display: flex;
  align-items: stretch;
  overflow-x: auto;
  scrollbar-width: none;
  flex: 1;
}

.tab-list::-webkit-scrollbar {
  display: none;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 12px;
  font-size: 13px;
  color: var(--el-text-color-secondary);
  cursor: pointer;
  white-space: nowrap;
  border-right: 1px solid var(--el-border-color-lighter);
  transition: background 0.15s, color 0.15s;
  position: relative;
  flex-shrink: 0;
}

.tab-item:hover {
  background: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
}

.tab-item.active {
  background: var(--el-bg-color-page);
  color: var(--el-color-primary);
  font-weight: 500;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--el-color-primary);
}

.tab-item.affix {
  padding-right: 16px;
}

.tab-title {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-close {
  font-size: 12px;
  border-radius: 50%;
  padding: 1px;
  flex-shrink: 0;
  color: var(--el-text-color-placeholder);
  transition: background 0.15s, color 0.15s;
}

.tab-close:hover {
  background: var(--el-fill-color-dark);
  color: var(--el-text-color-primary);
}

/* Context menu (not scoped - rendered via teleport) */
</style>

<style>
.tab-ctx-mask {
  position: fixed;
  inset: 0;
  z-index: 9998;
}

.tab-ctx-menu {
  position: fixed;
  z-index: 9999;
  background: var(--el-bg-color-overlay);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  box-shadow: var(--el-box-shadow-light);
  padding: 4px 0;
  list-style: none;
  margin: 0;
  min-width: 120px;
}

.tab-ctx-menu li {
  padding: 7px 16px;
  font-size: 13px;
  color: var(--el-text-color-regular);
  cursor: pointer;
  transition: background 0.15s;
}

.tab-ctx-menu li:hover {
  background: var(--el-fill-color-light);
}

.tab-ctx-menu li.disabled {
  color: var(--el-text-color-placeholder);
  cursor: not-allowed;
}

.tab-ctx-menu li.disabled:hover {
  background: transparent;
}
</style>
