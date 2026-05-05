export interface APIResource {
  id: number
  api_path: string
  api_method: string
  api_group: string
  api_desc: string
  api_status: number
  created_at: string
  updated_at: string
}

export interface APIListQuery {
  page: number
  page_size: number
  keyword?: string
  api_group?: string
  api_method?: string
  api_status?: number
}

export interface APIListResult {
  total: number
  items: APIResource[]
  page: number
  page_size: number
}

// HTTP_METHODS mirrors apimodel.ValidHTTPMethods on the backend.
export const HTTP_METHODS = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'] as const
export type HTTPMethod = (typeof HTTP_METHODS)[number]

// APIStatusActive / APIStatusDisabled mirror apimodel.StatusActive / StatusDisabled.
export const APIStatusActive = 1
export const APIStatusDisabled = 2

export interface SaveAPIPayload {
  api_path: string
  api_method: string
  api_group: string
  api_desc: string
  api_status?: number
}

export interface SkipRule {
  id: number
  api_path: string
  api_method: string
  skip_reason: string
  created_at: string
}

export interface SkipRuleListQuery {
  page: number
  page_size: number
  keyword?: string
  api_method?: string
}

export interface SkipRuleListResult {
  total: number
  items: SkipRule[]
  page: number
  page_size: number
}

export interface SaveSkipRulePayload {
  api_path: string
  api_method: string
  skip_reason: string
}
