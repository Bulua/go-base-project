import axios from 'axios'
import { getStoredSession } from '@/api/auth/session'
import type { ApiResponse } from '@/types/auth'

export const request = axios.create({
  baseURL: '/',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const session = getStoredSession()
  if (session?.access_token) {
    config.headers.Authorization = `Bearer ${session.access_token}`
  }
  return config
})

export async function unwrap<T>(promise: Promise<{ data: ApiResponse<T> }>): Promise<T> {
  const response = await promise
  if (response.data.code !== 0) {
    throw new Error(response.data.message || '请求失败')
  }
  return response.data.data
}
