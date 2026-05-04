import { request, unwrap } from '@/api/request'
import type {
  LoginLogQuery,
  LoginLogResult,
  OperationLogQuery,
  OperationLogResult,
  OperationLogRecord,
  CleanupResult,
} from '@/types/audit'

export function getLoginLogs(params: LoginLogQuery): Promise<LoginLogResult> {
  return unwrap<LoginLogResult>(request.get('/api/v1/audit/login-logs', { params }))
}

export function cleanupLoginLogs(days = 90): Promise<CleanupResult> {
  return unwrap<CleanupResult>(
    request.delete('/api/v1/audit/login-logs/cleanup', { params: { days } }),
  )
}

export function getOperationLogs(params: OperationLogQuery): Promise<OperationLogResult> {
  return unwrap<OperationLogResult>(request.get('/api/v1/audit/operation-logs', { params }))
}

export function getOperationLog(id: number): Promise<OperationLogRecord> {
  return unwrap<OperationLogRecord>(request.get(`/api/v1/audit/operation-logs/${id}`))
}

export function cleanupOperationLogs(days = 90): Promise<CleanupResult> {
  return unwrap<CleanupResult>(
    request.delete('/api/v1/audit/operation-logs/cleanup', { params: { days } }),
  )
}
