package apimodel

import (
	"errors"
	"time"
)

// ValidHTTPMethods is the single source of truth for allowed HTTP methods.
// The service, handler, and frontend should all derive from this list.
var ValidHTTPMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

// ── API Resources ──────────────────────────────────────────────────────────

type APIResource struct {
	ID        uint64    `json:"id"`
	APIPath   string    `json:"api_path"`
	APIMethod string    `json:"api_method"`
	APIGroup  string    `json:"api_group"`
	APIDesc   string    `json:"api_desc"`
	APIStatus int       `json:"api_status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIListQuery struct {
	Page      int
	PageSize  int
	Keyword   string
	APIGroup  string
	APIMethod string
	APIStatus int
}

type APIListResult struct {
	Total    int64         `json:"total"`
	Items    []APIResource `json:"items"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// CreateAPIRequest is the decoded request body for POST /api/v1/apis.
type CreateAPIRequest struct {
	APIPath   string `json:"api_path"`
	APIMethod string `json:"api_method"`
	APIGroup  string `json:"api_group"`
	APIDesc   string `json:"api_desc"`
}

// UpdateAPIRequest is the decoded request body for PUT /api/v1/apis/{id}.
type UpdateAPIRequest struct {
	APIPath   string `json:"api_path"`
	APIMethod string `json:"api_method"`
	APIGroup  string `json:"api_group"`
	APIDesc   string `json:"api_desc"`
	APIStatus int    `json:"api_status"`
}

type SaveAPIPayload struct {
	APIPath   string
	APIMethod string
	APIGroup  string
	APIDesc   string
	APIStatus int
}

// ── Skip Rules ─────────────────────────────────────────────────────────────

type SkipRule struct {
	ID         uint64    `json:"id"`
	APIPath    string    `json:"api_path"`
	APIMethod  string    `json:"api_method"`
	SkipReason string    `json:"skip_reason"`
	CreatedAt  time.Time `json:"created_at"`
}

type SkipRuleListQuery struct {
	Page      int
	PageSize  int
	Keyword   string
	APIMethod string
}

type SkipRuleListResult struct {
	Total    int64      `json:"total"`
	Items    []SkipRule `json:"items"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// CreateSkipRuleRequest is the decoded request body for POST /api/v1/api-skip-rules.
type CreateSkipRuleRequest struct {
	APIPath    string `json:"api_path"`
	APIMethod  string `json:"api_method"`
	SkipReason string `json:"skip_reason"`
}

type SaveSkipRulePayload struct {
	APIPath    string
	APIMethod  string
	SkipReason string
}

// ── Actor ──────────────────────────────────────────────────────────────────

type ActorContext struct {
	UserID    uint64
	LoginName string
	SourceIP  string
	UserAgent string
}

// ── Errors ─────────────────────────────────────────────────────────────────

const (
	StatusActive   = 1
	StatusDisabled = 2
)

var (
	ErrAPINotFound      = errors.New("api: not found")
	ErrAPIPathEmpty     = errors.New("api: path is empty")
	ErrAPIMethodInvalid = errors.New("api: method is invalid")
	ErrAPIPathTaken     = errors.New("api: path+method already exists")
	ErrAPIHasPolicies   = errors.New("api: has associated permission policies")
	ErrSkipNotFound     = errors.New("skip rule: not found")
	ErrSkipPathEmpty    = errors.New("skip rule: path is empty")
	ErrSkipPathTaken    = errors.New("skip rule: path+method already exists")
)
