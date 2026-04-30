package rolemodel

import (
	"errors"
	"time"
)

const (
	StatusActive   = 1
	StatusDisabled = 2

	// BuiltinRoleID 是 seed 写入的内置最高权限角色（id=1）。
	// 不依赖 role_code 字符串，避免管理员改名后保护逻辑失效。
	BuiltinRoleID uint64 = 1

	SubjectTypeRole = "role"
	ResourceTypeAPI = "api"
	EffectAllow     = "allow"
)

var (
	ErrRoleNotFound    = errors.New("role not found")
	ErrRoleCodeTaken   = errors.New("role code already exists")
	ErrRoleCodeInvalid = errors.New("role code invalid")
	ErrInvalidStatus   = errors.New("invalid role status")
	ErrBuiltinProtect  = errors.New("built-in role is protected")
	ErrRoleHasUsers    = errors.New("role still has users assigned")
	ErrRoleHasChildren = errors.New("role still has child roles")
	ErrParentLoop      = errors.New("parent role would create a cycle")
)

type Role struct {
	ID           uint64    `json:"id"`
	RoleCode     string    `json:"role_code"`
	RoleName     string    `json:"role_name"`
	ParentRoleID uint64    `json:"parent_role_id"`
	DefaultRoute string    `json:"default_route"`
	SortNo       int       `json:"sort_no"`
	RoleStatus   int       `json:"role_status"`
	Remark       *string   `json:"remark,omitempty"`
	UserCount    int64     `json:"user_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Children     []Role    `json:"children,omitempty"`
}

type ListQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	RoleStatus int
}

type ListResult struct {
	Total int64  `json:"total"`
	Items []Role `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"page_size"`
}

type CreateRequest struct {
	RoleCode     string  `json:"role_code"`
	RoleName     string  `json:"role_name"`
	ParentRoleID uint64  `json:"parent_role_id"`
	DefaultRoute string  `json:"default_route"`
	SortNo       int     `json:"sort_no"`
	RoleStatus   int     `json:"role_status"`
	Remark       *string `json:"remark,omitempty"`
}

type UpdateRequest struct {
	RoleCode     string  `json:"role_code"`
	RoleName     string  `json:"role_name"`
	ParentRoleID uint64  `json:"parent_role_id"`
	DefaultRoute string  `json:"default_route"`
	SortNo       int     `json:"sort_no"`
	RoleStatus   int     `json:"role_status"`
	Remark       *string `json:"remark,omitempty"`
}

type AssignIDsRequest struct {
	IDs []uint64 `json:"ids"`
}

// MenuOption is a node in the menu tree picker.
type MenuOption struct {
	ID        uint64       `json:"id"`
	ParentID  uint64       `json:"parent_id"`
	MenuTitle string       `json:"menu_title"`
	MenuType  int          `json:"menu_type"`
	SortNo    int          `json:"sort_no"`
	Children  []MenuOption `json:"children,omitempty"`
}

// ActionOption groups action buttons by their owning menu.
type ActionOption struct {
	ID         uint64 `json:"id"`
	MenuID     uint64 `json:"menu_id"`
	MenuTitle  string `json:"menu_title"`
	ActionCode string `json:"action_code"`
	ActionName string `json:"action_name"`
	SortNo     int    `json:"sort_no"`
}

// APIOption groups API resources by their group.
type APIOption struct {
	ID         uint64 `json:"id"`
	APIPath    string `json:"api_path"`
	APIMethod  string `json:"api_method"`
	APIGroup   string `json:"api_group"`
	APIDesc    string `json:"api_desc"`
	APIStatus  int    `json:"api_status"`
}

type ActorContext struct {
	UserID    uint64
	LoginName string
	SourceIP  string
	UserAgent string
}
