package menumodel

import (
	"errors"
	"time"
)

const (
	StatusActive   = 1
	StatusDisabled = 2

	TypeDirectory = 1
	TypePage      = 2
	TypeHidden    = 3
	TypeExternal  = 4
)

var (
	ErrMenuNotFound    = errors.New("menu not found")
	ErrMenuTitleEmpty  = errors.New("menu title is required")
	ErrInvalidStatus   = errors.New("invalid menu status")
	ErrInvalidType     = errors.New("invalid menu type")
	ErrParentNotFound  = errors.New("parent menu not found")
	ErrHasChildren     = errors.New("menu still has children")
	ErrActionNotFound  = errors.New("action not found")
	ErrActionCodeEmpty = errors.New("action code is required")
	ErrActionNameEmpty = errors.New("action name is required")
	ErrActionCodeTaken = errors.New("action code already exists for this menu")
)

type Menu struct {
	ID             uint64       `json:"id"`
	ParentID       uint64       `json:"parent_id"`
	MenuType       int          `json:"menu_type"`
	RoutePath      *string      `json:"route_path,omitempty"`
	RouteName      *string      `json:"route_name,omitempty"`
	ComponentPath  *string      `json:"component_path,omitempty"`
	RedirectPath   *string      `json:"redirect_path,omitempty"`
	MenuTitle      string       `json:"menu_title"`
	MenuIcon       *string      `json:"menu_icon,omitempty"`
	SortNo         int          `json:"sort_no"`
	IsHidden       bool         `json:"is_hidden"`
	IsKeepAlive    bool         `json:"is_keep_alive"`
	IsAffix        bool         `json:"is_affix"`
	ActiveRoute    *string      `json:"active_route,omitempty"`
	TransitionName *string      `json:"transition_name,omitempty"`
	ExternalURL    *string      `json:"external_url,omitempty"`
	MenuStatus     int    `json:"menu_status"`
	Children       []Menu `json:"children,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type ListQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	MenuStatus int
	MenuType   int
}

type ListResult struct {
	Total int64  `json:"total"`
	Items []Menu `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"page_size"`
}

type SaveRequest struct {
	ParentID       uint64  `json:"parent_id"`
	MenuType       int     `json:"menu_type"`
	RoutePath      *string `json:"route_path,omitempty"`
	RouteName      *string `json:"route_name,omitempty"`
	ComponentPath  *string `json:"component_path,omitempty"`
	RedirectPath   *string `json:"redirect_path,omitempty"`
	MenuTitle      string  `json:"menu_title"`
	MenuIcon       *string `json:"menu_icon,omitempty"`
	SortNo         int     `json:"sort_no"`
	IsHidden       bool    `json:"is_hidden"`
	IsKeepAlive    bool    `json:"is_keep_alive"`
	IsAffix        bool    `json:"is_affix"`
	ActiveRoute    *string `json:"active_route,omitempty"`
	TransitionName *string `json:"transition_name,omitempty"`
	ExternalURL    *string `json:"external_url,omitempty"`
	MenuStatus     int     `json:"menu_status"`
}

type MenuAction struct {
	ID           uint64    `json:"id"`
	MenuID       uint64    `json:"menu_id"`
	ActionCode   string    `json:"action_code"`
	ActionName   string    `json:"action_name"`
	ActionDesc   *string   `json:"action_desc,omitempty"`
	SortNo       int       `json:"sort_no"`
	ActionStatus int       `json:"action_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SaveActionRequest struct {
	ActionCode   string  `json:"action_code"`
	ActionName   string  `json:"action_name"`
	ActionDesc   *string `json:"action_desc,omitempty"`
	SortNo       int     `json:"sort_no"`
	ActionStatus int     `json:"action_status"`
}

type ActorContext struct {
	UserID    uint64
	LoginName string
	SourceIP  string
	UserAgent string
}
