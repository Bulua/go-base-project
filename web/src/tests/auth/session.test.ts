import { beforeEach, describe, expect, it } from 'vitest'
import {
  AUTH_SESSION_KEY,
  clearStoredSession,
  getStoredSession,
  setStoredSession,
} from '@/api/auth/session'

describe('auth session storage', () => {
  beforeEach(() => {
    const store = new Map<string, string>()
    const storage = {
      getItem: (key: string) => store.get(key) ?? null,
      setItem: (key: string, value: string) => store.set(key, value),
      removeItem: (key: string) => store.delete(key),
      clear: () => store.clear(),
      key: (index: number) => Array.from(store.keys())[index] ?? null,
      get length() {
        return store.size
      },
    } as Storage
    Object.defineProperty(globalThis, 'localStorage', { value: storage, configurable: true })
    Object.defineProperty(window, 'localStorage', { value: storage, configurable: true })
    window.localStorage.clear()
  })

  it('stores and restores token session data', () => {
    setStoredSession({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      token_type: 'Bearer',
      expires_at: '2026-04-30T10:00:00Z',
      refresh_expires_at: '2026-05-07T10:00:00Z',
    })

    expect(getStoredSession()?.access_token).toBe('access-token')
    expect(window.localStorage.getItem(AUTH_SESSION_KEY)).toContain('refresh-token')
  })

  it('clears malformed session data', () => {
    window.localStorage.setItem(AUTH_SESSION_KEY, '{bad-json')

    expect(getStoredSession()).toBeNull()
    expect(window.localStorage.getItem(AUTH_SESSION_KEY)).toBeNull()
  })

  it('removes session data on logout', () => {
    setStoredSession({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      token_type: 'Bearer',
      expires_at: '2026-04-30T10:00:00Z',
      refresh_expires_at: '2026-05-07T10:00:00Z',
    })

    clearStoredSession()

    expect(getStoredSession()).toBeNull()
  })
})
