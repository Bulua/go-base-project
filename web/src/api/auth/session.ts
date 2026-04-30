import type { AuthSession } from '@/types/auth'

export const AUTH_SESSION_KEY = 'gobaseproject.auth.session'

export function getStoredSession(): AuthSession | null {
  const raw = localStorage.getItem(AUTH_SESSION_KEY)
  if (!raw) {
    return null
  }
  try {
    const session = JSON.parse(raw) as Partial<AuthSession>
    if (!session.access_token || !session.refresh_token) {
      clearStoredSession()
      return null
    }
    return session as AuthSession
  } catch {
    clearStoredSession()
    return null
  }
}

export function setStoredSession(session: AuthSession): void {
  localStorage.setItem(AUTH_SESSION_KEY, JSON.stringify(session))
}

export function clearStoredSession(): void {
  localStorage.removeItem(AUTH_SESSION_KEY)
}
