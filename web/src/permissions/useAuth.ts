import { useAuthStore } from '@/store/modules/auth'

/**
 * Check whether the current user has a specific button action permission.
 * Usage:
 *   const { hasAction } = useAuth()
 *   hasAction('add')           // any menu has 'add'
 *   hasAction(11, 'add')       // menu 11 has 'add'
 */
export function useAuth() {
  const authStore = useAuthStore()

  function hasAction(codeOrMenuId: string | number, code?: string): boolean {
    if (typeof codeOrMenuId === 'string') {
      return authStore.actionCodes.has(codeOrMenuId)
    }
    return authStore.actionKeys.has(`${codeOrMenuId}:${code}`)
  }

  return { hasAction }
}
