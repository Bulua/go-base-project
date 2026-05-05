<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as ElIcons from '@element-plus/icons-vue'
import {
  ArrowDown,
  ArrowRight,
  Bell,
  HomeFilled,
  Moon,
  Setting,
  Sunny,
  SwitchButton,
  Expand,
  Fold,
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/modules/auth'
import { useTheme } from '@/composables/common/useTheme'
import { useWatermark } from '@/composables/common/useWatermark'
import SettingsDrawer from '@/components/layout/SettingsDrawer.vue'
import TabBar from '@/components/layout/TabBar.vue'
import { useTabsStore } from '@/store/modules/tabs'
import type { MenuRoute } from '@/types/auth'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const { isDark, toggleDark } = useTheme()

const collapsed = ref(false)
const openGroups = ref<Set<number>>(new Set())

const iconMap = ElIcons as Record<string, unknown>
function resolveIcon(name?: string) {
  if (!name) return null
  return iconMap[name] ?? null
}

function joinPath(parent: string, child: string): string {
  if (!child) return parent || '/'
  if (child.startsWith('/')) return child
  const sep = parent.endsWith('/') ? '' : '/'
  return `${parent || ''}${sep}${child}`.replace(/\/+/g, '/')
}

interface DisplayMenu extends MenuRoute {
  fullPath: string
  hasChildren: boolean
}

function decorate(menus: MenuRoute[], parentPath = ''): DisplayMenu[] {
  return menus
    .filter((m) => !m.is_hidden)
    .map((m) => {
      const fullPath = joinPath(parentPath, m.route_path ?? '')
      const children = m.children ? decorate(m.children, fullPath) : []
      return { ...m, fullPath, hasChildren: children.length > 0, children: children as MenuRoute[] }
    })
}

const decoratedMenus = computed(() => decorate(authStore.menuRoutes))

const activePath = computed(() => route.path)

function isItemActive(item: DisplayMenu): boolean {
  if (item.fullPath === activePath.value) return true
  if (item.hasChildren) {
    return (item.children as DisplayMenu[]).some((c) => isItemActive(c))
  }
  return false
}

function toggleGroup(id: number) {
  if (openGroups.value.has(id)) openGroups.value.delete(id)
  else openGroups.value.add(id)
}

function isGroupOpen(item: DisplayMenu): boolean {
  // Auto-open if a descendant is active
  if (isItemActive(item) && item.fullPath !== activePath.value) return true
  return openGroups.value.has(item.id)
}

async function navigate(item: DisplayMenu) {
  if (item.hasChildren) {
    toggleGroup(item.id)
    return
  }
  if (item.fullPath && item.fullPath !== activePath.value) {
    await router.push(item.fullPath)
  }
}

const breadcrumbs = computed(() => {
  const trail: { title: string; path: string }[] = []
  const findPath = (menus: DisplayMenu[], parents: DisplayMenu[]): boolean => {
    for (const m of menus) {
      if (m.fullPath === activePath.value) {
        const chain = [...parents, m]
        for (const c of chain) {
          trail.push({ title: c.menu_title, path: c.fullPath })
        }
        return true
      }
      if (m.hasChildren && findPath(m.children as DisplayMenu[], [...parents, m])) {
        return true
      }
    }
    return false
  }
  findPath(decoratedMenus.value, [])
  if (trail.length === 0) {
    trail.push({ title: (route.meta.title as string) || '页面', path: activePath.value })
  }
  return trail
})

const username = computed(
  () => authStore.currentUser?.display_name || authStore.currentUser?.login_name || 'Admin',
)

const userInitial = computed(() => username.value.charAt(0).toUpperCase())

useWatermark(username)

const settingsVisible = ref(false)
const tabsStore = useTabsStore()

async function handleLogout() {
  tabsStore.reset()
  await authStore.logoutCurrentUser()
}
</script>

<template>
  <div class="bp-app" v-loading="authStore.workspaceLoading">
    <!-- Sidebar -->
    <aside
      class="bp-sidebar"
      :class="{ collapsed }"
      :style="{ width: collapsed ? 'var(--sidebar-col)' : 'var(--sidebar-w)' }"
    >
      <div class="bp-sidebar-logo">
        <span class="mark">B</span>
        <span v-if="!collapsed" class="name">BaseProject</span>
      </div>

      <nav class="bp-sidebar-nav">
        <template v-for="item in decoratedMenus" :key="item.id">
          <div
            class="bp-nav-item"
            :class="{ active: isItemActive(item) && !item.hasChildren }"
            :title="collapsed ? item.menu_title : undefined"
            @click="navigate(item)"
          >
            <el-icon class="icon" v-if="resolveIcon(item.menu_icon)">
              <component :is="resolveIcon(item.menu_icon)" />
            </el-icon>
            <span v-if="!collapsed" class="label">{{ item.menu_title }}</span>
            <el-icon
              v-if="!collapsed && item.hasChildren"
              class="chev"
              :class="{ open: isGroupOpen(item) }"
            >
              <ArrowDown />
            </el-icon>
          </div>

          <template v-if="item.hasChildren && isGroupOpen(item) && !collapsed">
            <div
              v-for="child in (item.children as DisplayMenu[])"
              :key="child.id"
              class="bp-nav-subitem"
              :class="{ active: isItemActive(child) }"
              @click="navigate(child)"
            >
              <el-icon class="icon" v-if="resolveIcon(child.menu_icon)">
                <component :is="resolveIcon(child.menu_icon)" />
              </el-icon>
              <span class="label">{{ child.menu_title }}</span>
            </div>
          </template>
        </template>
      </nav>

      <div class="bp-sidebar-footer">
        <button
          class="bp-sidebar-collapse-btn"
          :title="collapsed ? '展开' : '收起'"
          @click="collapsed = !collapsed"
        >
          <el-icon>
            <component :is="collapsed ? Expand : Fold" />
          </el-icon>
        </button>
      </div>
    </aside>

    <div class="bp-main-area">
      <header class="bp-navbar">
        <div class="bp-breadcrumb">
          <el-icon class="home"><HomeFilled /></el-icon>
          <template v-for="(c, i) in breadcrumbs" :key="c.path + i">
            <el-icon v-if="i > 0" class="sep"><ArrowRight /></el-icon>
            <span class="crumb" :class="{ last: i === breadcrumbs.length - 1 }">
              {{ c.title }}
            </span>
          </template>
        </div>

        <div class="bp-navbar-actions">
          <button class="bp-icon-btn" :title="isDark ? '切换浅色' : '切换深色'" @click="toggleDark">
            <el-icon>
              <component :is="isDark ? Sunny : Moon" />
            </el-icon>
          </button>

          <button class="bp-icon-btn" title="外观设置" @click="settingsVisible = true">
            <el-icon><Setting /></el-icon>
          </button>

          <el-badge :value="3" :max="9">
            <button class="bp-icon-btn" title="通知">
              <el-icon><Bell /></el-icon>
            </button>
          </el-badge>

          <el-dropdown trigger="click" @command="(c: string) => c === 'logout' && handleLogout()">
            <div class="bp-user-trigger">
              <div class="bp-user-avatar">{{ userInitial }}</div>
              <span class="bp-user-name">{{ username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>{{ authStore.roleNames }}</el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <TabBar />

      <main class="bp-content">
        <router-view v-slot="{ Component }">
          <keep-alive :include="tabsStore.cachedNames">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </main>
    </div>

    <SettingsDrawer v-model="settingsVisible" />
  </div>
</template>

<style scoped>
</style>

