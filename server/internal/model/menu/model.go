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
	ErrMenuNotFound       = errors.New("menu not found")
	ErrMenuTitleEmpty     = errors.New("menu title is required")
	ErrInvalidStatus      = errors.New("invalid menu status")
	ErrInvalidType        = errors.New("invalid menu type")
	ErrParentNotFound     = errors.New("parent menu not found")
	ErrHasChildren        = errors.New("menu still has children")
	ErrParamNotFound      = errors.New("route param not found")
	ErrParamKeyEmpty      = errors.New("param key is required")
	ErrParamModeBad       = errors.New("param mode must be query or path")
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
	MenuStatus     int          `json:"menu_status"`
	Params         []RouteParam `json:"params,omitempty"`
	Children       []Menu       `json:"children,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type RouteParam struct {
	ID         uint64    `json:"id"`
	MenuID     uint64    `json:"menu_id"`
	ParamMode  string    `json:"param_mode"`
	ParamKey   string    `json:"param_key"`
	ParamValue *string   `json:"param_value,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
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

type CreateParamRequest struct {
	ParamMode  string  `json:"param_mode"`
	ParamKey   string  `json:"param_key"`
	ParamValue *string `json:"param_value,omitempty"`
}

type ActorContext struct {
	UserID    uint64
	LoginName string
	SourceIP  string
	UserAgent string
}
