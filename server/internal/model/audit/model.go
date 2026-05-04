package auditmodel

import "time"

// ── Write models (used by middleware / auth service) ──────────────────────

type OperationLog struct {
	UserID        *uint64
	SourceIP      string
	RequestMethod string
	RequestPath   string
	StatusCode    int
	CostMs        int64
	UserAgent     string
	ErrorMessage  string
	RequestBody   string
	ResponseBody  string
}

// ── Read models (used by audit query API) ────────────────────────────────

type LoginLogRecord struct {
	ID           uint64     `json:"id"`
	UserID       *uint64    `json:"user_id"`
	LoginName    *string    `json:"login_name"`
	SourceIP     *string    `json:"source_ip"`
	LoginSuccess bool       `json:"login_success"`
	FailReason   *string    `json:"fail_reason"`
	UserAgent    *string    `json:"user_agent"`
	CreatedAt    time.Time  `json:"created_at"`
}

type OperationLogRecord struct {
	ID            uint64     `json:"id"`
	UserID        *uint64    `json:"user_id"`
	SourceIP      *string    `json:"source_ip"`
	RequestMethod *string    `json:"request_method"`
	RequestPath   *string    `json:"request_path"`
	StatusCode    *int       `json:"status_code"`
	CostMs        *int64     `json:"cost_ms"`
	UserAgent     *string    `json:"user_agent"`
	ErrorMessage  *string    `json:"error_message"`
	RequestBody   *string    `json:"request_body"`
	ResponseBody  *string    `json:"response_body"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ── Query types ───────────────────────────────────────────────────────────

type LoginLogQuery struct {
	Page        int
	PageSize    int
	Keyword     string
	LoginSuccess int // 0=all 1=success 2=fail
	StartDate   string
	EndDate     string
}

type OperationLogQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	Method     string
	StatusCode int
	StartDate  string
	EndDate    string
}

type LoginLogResult struct {
	Total    int64            `json:"total"`
	Items    []LoginLogRecord `json:"items"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

type OperationLogResult struct {
	Total    int64               `json:"total"`
	Items    []OperationLogRecord `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

type CleanupResult struct {
	Deleted int64 `json:"deleted"`
}
