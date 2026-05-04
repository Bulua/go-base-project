export interface LoginLogRecord {
  id: number
  user_id: number | null
  login_name: string | null
  source_ip: string | null
  login_success: boolean
  fail_reason: string | null
  user_agent: string | null
  created_at: string
}

export interface OperationLogRecord {
  id: number
  user_id: number | null
  source_ip: string | null
  request_method: string | null
  request_path: string | null
  status_code: number | null
  cost_ms: number | null
  user_agent: string | null
  error_message: string | null
  request_body: string | null
  response_body: string | null
  created_at: string
}

export interface LoginLogQuery {
  page: number
  page_size: number
  keyword?: string
  login_success?: number
  start_date?: string
  end_date?: string
}

export interface OperationLogQuery {
  page: number
  page_size: number
  keyword?: string
  method?: string
  status_code?: number
  start_date?: string
  end_date?: string
}

export interface LoginLogResult {
  total: number
  items: LoginLogRecord[]
  page: number
  page_size: number
}

export interface OperationLogResult {
  total: number
  items: OperationLogRecord[]
  page: number
  page_size: number
}

export interface CleanupResult {
  deleted: number
}
