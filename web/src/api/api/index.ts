import { request, unwrap } from '@/api/request'
import type {
  APIListQuery,
  APIListResult,
  APIResource,
  SaveAPIPayload,
  SkipRule,
  SkipRuleListQuery,
  SkipRuleListResult,
  SaveSkipRulePayload,
} from '@/types/api'

export function getAPIGroups(): Promise<string[]> {
  return unwrap<string[]>(request.get('/api/v1/apis/groups'))
}

export function listAPIs(params: APIListQuery): Promise<APIListResult> {
  return unwrap<APIListResult>(request.get('/api/v1/apis', { params }))
}

export function createAPI(payload: SaveAPIPayload): Promise<APIResource> {
  return unwrap<APIResource>(request.post('/api/v1/apis', payload))
}

export function updateAPI(id: number, payload: SaveAPIPayload): Promise<APIResource> {
  return unwrap<APIResource>(request.put(`/api/v1/apis/${id}`, payload))
}

export function deleteAPI(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/apis/${id}`))
}

export function listSkipRules(params: SkipRuleListQuery): Promise<SkipRuleListResult> {
  return unwrap<SkipRuleListResult>(request.get('/api/v1/api-skip-rules', { params }))
}

export function createSkipRule(payload: SaveSkipRulePayload): Promise<SkipRule> {
  return unwrap<SkipRule>(request.post('/api/v1/api-skip-rules', payload))
}

export function deleteSkipRule(id: number): Promise<{ success: boolean }> {
  return unwrap<{ success: boolean }>(request.delete(`/api/v1/api-skip-rules/${id}`))
}
