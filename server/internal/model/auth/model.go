package authmodel

import "time"

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type User struct {
	ID                 uint64  `json:"id"`
	UserUUID           string  `json:"user_uuid"`
	LoginName          string  `json:"login_name"`
	DisplayName        string  `json:"display_name"`
	AvatarURL          *string `json:"avatar_url,omitempty"`
	PrimaryRoleID      *uint64 `json:"primary_role_id,omitempty"`
	PhoneNumber        *string `json:"phone_number,omitempty"`
	EmailAddress       *string `json:"email_address,omitempty"`
	Remark             *string `json:"remark,omitempty"`
	PasswordHash       string  `json:"-"`
	UserStatus         int     `json:"user_status"`
	MustChangePassword bool    `json:"must_change_password"`
	Roles              []Role  `json:"roles,omitempty"`
}

type Role struct {
	ID           uint64 `json:"id"`
	RoleCode     string `json:"role_code"`
	RoleName     string `json:"role_name"`
	DefaultRoute string `json:"default_route"`
}

type Menu struct {
	ID             uint64  `json:"id"`
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
	Children       []Menu  `json:"children,omitempty"`
}

type Action struct {
	MenuID     uint64 `json:"menu_id"`
	ActionID   uint64 `json:"action_id"`
	ActionCode string `json:"action_code"`
	ActionName string `json:"action_name"`
}

type LoginAudit struct {
	UserID       *uint64
	LoginName    string
	SourceIP     string
	LoginSuccess bool
	FailReason   string
	UserAgent    string
}

type BlockedToken struct {
	TokenHash string
	ExpiresAt time.Time
	Reason    string
}

type LoginRequest struct {
	LoginName string `json:"login_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RequestMeta struct {
	SourceIP  string
	UserAgent string
}

type Session struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	TokenType        string    `json:"token_type"`
	ExpiresAt        time.Time `json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	User             User      `json:"user"`
}
