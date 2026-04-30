package usermodel

import (
	"errors"
	"time"
)

const (
	StatusActive = 1
	StatusFrozen = 2

	BuiltinAdminLogin = "admin"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrLoginNameTaken    = errors.New("login name already exists")
	ErrLoginNameInvalid  = errors.New("login name invalid")
	ErrPasswordWeak      = errors.New("password too weak")
	ErrInvalidStatus     = errors.New("invalid user status")
	ErrAdminProtected    = errors.New("admin user is protected")
	ErrRoleAssignmentBad = errors.New("invalid role assignment")
)

type Role struct {
	ID       uint64 `json:"id"`
	RoleCode string `json:"role_code"`
	RoleName string `json:"role_name"`
}

type User struct {
	ID                 uint64    `json:"id"`
	UserUUID           string    `json:"user_uuid"`
	LoginName          string    `json:"login_name"`
	DisplayName        string    `json:"display_name"`
	AvatarURL          *string   `json:"avatar_url,omitempty"`
	PrimaryRoleID      *uint64   `json:"primary_role_id,omitempty"`
	PhoneNumber        *string   `json:"phone_number,omitempty"`
	EmailAddress       *string   `json:"email_address,omitempty"`
	UserStatus         int       `json:"user_status"`
	MustChangePassword bool      `json:"must_change_password"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	Remark             *string   `json:"remark,omitempty"`
	Roles              []Role    `json:"roles,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ListQuery struct {
	Page       int
	PageSize   int
	Keyword    string
	UserStatus int
	RoleID     uint64
}

type ListResult struct {
	Total int64  `json:"total"`
	Items []User `json:"items"`
	Page  int    `json:"page"`
	Size  int    `json:"page_size"`
}

type CreateRequest struct {
	LoginName     string   `json:"login_name"`
	Password      string   `json:"password"`
	DisplayName   string   `json:"display_name"`
	AvatarURL     *string  `json:"avatar_url,omitempty"`
	PrimaryRoleID *uint64  `json:"primary_role_id,omitempty"`
	PhoneNumber   *string  `json:"phone_number,omitempty"`
	EmailAddress  *string  `json:"email_address,omitempty"`
	UserStatus    int      `json:"user_status"`
	Remark        *string  `json:"remark,omitempty"`
	RoleIDs       []uint64 `json:"role_ids,omitempty"`
}

type UpdateRequest struct {
	DisplayName   string  `json:"display_name"`
	AvatarURL     *string `json:"avatar_url,omitempty"`
	PrimaryRoleID *uint64 `json:"primary_role_id,omitempty"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	EmailAddress  *string `json:"email_address,omitempty"`
	Remark        *string `json:"remark,omitempty"`
}

type UpdateStatusRequest struct {
	UserStatus int `json:"user_status"`
}

type ResetPasswordRequest struct {
	Password           string `json:"password"`
	MustChangePassword *bool  `json:"must_change_password,omitempty"`
}

type AssignRolesRequest struct {
	RoleIDs []uint64 `json:"role_ids"`
}

type ActorContext struct {
	UserID    uint64
	LoginName string
	SourceIP  string
	UserAgent string
}
