import { request, unwrap } from '@/api/request'
import type { AuthAction, AuthSession, CurrentUser, MenuRoute } from '@/types/auth'

export interface LoginPayload {
  login_name: string
  password: string
}

export function login(payload: LoginPayload): Promise<AuthSession> {
  return unwrap<AuthSession>(request.post('/api/v1/auth/login', payload))
}

export function logout(): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.post('/api/v1/auth/logout'))
}

export function refresh(refreshToken: string): Promise<AuthSession> {
  return unwrap<AuthSession>(request.post('/api/v1/auth/refresh', { refresh_token: refreshToken }))
}

export function getProfile(): Promise<CurrentUser> {
  return unwrap<CurrentUser>(request.get('/api/v1/auth/profile'))
}

export function getRoutes(): Promise<MenuRoute[]> {
  return unwrap<MenuRoute[]>(request.get('/api/v1/auth/routes'))
}

export function getActions(): Promise<AuthAction[]> {
  return unwrap<AuthAction[]>(request.get('/api/v1/auth/actions'))
}
