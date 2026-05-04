import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import {
  getActions,
  getProfile,
  getRoutes,
  login,
  logout,
} from '@/api/auth'
import {
  clearStoredSession,
  getStoredSession,
  setStoredSession,
} from '@/api/auth/session'
import router from '@/router'
import { clearDynamicRoutes, registerDynamicRoutes } from '@/router/dynamic'
import type { AuthAction, AuthSession, CurrentUser, MenuRoute } from '@/types/auth'

export interface FlatMenu extends MenuRoute {
  depth: number
}

export const useAuthStore = defineStore('auth', () => {
  const session = ref<AuthSession | null>(getStoredSession())
  const currentUser = ref<CurrentUser | null>(session.value?.user ?? null)
  const menuRoutes = ref<MenuRoute[]>([])
  const actions = ref<AuthAction[]>([])
  const loginLoading = ref(false)
  const workspaceLoading = ref(false)

  const isAuthenticated = computed(() => Boolean(session.value?.access_token))
  const roleNames = computed(() =>
    currentUser.value?.roles?.map((role) => role.role_name).join(' / ') || '未分配角色',
  )
  const flatMenus = computed(() => flattenMenus(menuRoutes.value))
  // actionCodes: Set of action_code for quick loose check (e.g. 'add')
  const actionCodes = computed(() => new Set(actions.value.map((a) => a.action_code)))
  // actionKeys: Set of "menuId:actionCode" for precise check (e.g. '11:add')
  const actionKeys = computed(
    () => new Set(actions.value.map((a) => `${a.menu_id}:${a.action_code}`)),
  )

  async function loginWithPassword(loginName: string, password: string) {
    loginLoading.value = true
    try {
      const nextSession = await login({
        login_name: loginName.trim(),
        password,
      })
      session.value = nextSession
      currentUser.value = nextSession.user ?? null
      setStoredSession(nextSession)
      await loadWorkspace()
      ElMessage.success('登录成功')
      const target =
        (router.currentRoute.value.query.redirect as string | undefined) ?? '/'
      await router.push(target)
    } catch (error) {
      ElMessage.error(error instanceof Error ? error.message : '登录失败')
      throw error
    } finally {
      loginLoading.value = false
    }
  }

  async function logoutCurrentUser() {
    try {
      if (session.value) {
        await logout()
      }
    } catch {
      // 本地登录态仍然需要清理，避免过期 Token 卡在页面里。
    } finally {
      clearLocalAuth()
      await router.push('/login')
    }
  }

  async function loadWorkspace() {
    workspaceLoading.value = true
    try {
      const [profile, routes, actionList] = await Promise.all([
        getProfile(),
        getRoutes(),
        getActions(),
      ])
      currentUser.value = profile
      menuRoutes.value = routes
      actions.value = actionList
      if (session.value) {
        setStoredSession({ ...session.value, user: profile })
      }
      registerDynamicRoutes(routes)
    } catch (error) {
      clearLocalAuth()
      ElMessage.error(error instanceof Error ? error.message : '登录状态已失效')
    } finally {
      workspaceLoading.value = false
    }
  }

  function clearLocalAuth() {
    clearStoredSession()
    session.value = null
    currentUser.value = null
    menuRoutes.value = []
    actions.value = []
    clearDynamicRoutes()
  }

  return {
    actions,
    actionCodes,
    actionKeys,
    currentUser,
    flatMenus,
    isAuthenticated,
    loadWorkspace,
    loginLoading,
    loginWithPassword,
    logoutCurrentUser,
    menuRoutes,
    roleNames,
    session,
    workspaceLoading,
  }
})

function flattenMenus(menus: MenuRoute[], depth = 0): FlatMenu[] {
  return menus.flatMap((menu) => [
    { ...menu, depth },
    ...flattenMenus(menu.children ?? [], depth + 1),
  ])
}
