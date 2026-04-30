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
  let response: { data: ApiResponse<T> }
  try {
    response = await promise
  } catch (error) {
    throw new Error(extractErrorMessage(error))
  }
  if (response.data.code !== 0) {
    throw new Error(response.data.message || '请求失败')
  }
  return response.data.data
}

function extractErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data
    if (isApiErrorBody(data) && data.message) {
      return data.message
    }
  }
  if (error instanceof Error && error.message) {
    return error.message
  }
  return '请求失败'
}

function isApiErrorBody(value: unknown): value is { message?: string } {
  return typeof value === 'object' && value !== null && 'message' in value
}
