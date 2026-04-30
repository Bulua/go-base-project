package userservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	auditmodel "gobaseproject/server/internal/model/audit"
	usermodel "gobaseproject/server/internal/model/user"
	authservice "gobaseproject/server/internal/service/auth"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var loginNamePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]{2,31}$`)

type Repository interface {
	Count(ctx context.Context, q usermodel.ListQuery) (int64, error)
	List(ctx context.Context, q usermodel.ListQuery) ([]usermodel.User, error)
	GetByID(ctx context.Context, id uint64) (*usermodel.User, error)
	GetByLoginName(ctx context.Context, loginName string) (*usermodel.User, error)
	Create(ctx context.Context, userUUID string, req usermodel.CreateRequest, passwordHash string) (uint64, error)
	Update(ctx context.Context, id uint64, req usermodel.UpdateRequest) error
	SoftDelete(ctx context.Context, id uint64) error
	UpdateStatus(ctx context.Context, id uint64, status int) error
	UpdatePassword(ctx context.Context, id uint64, passwordHash string, mustChange bool) error
	ReplaceRoles(ctx context.Context, userID uint64, roleIDs []uint64) error
	ListActiveRoles(ctx context.Context) ([]usermodel.Role, error)
}

type AuditRepository interface {
	InsertOperationLog(ctx context.Context, log auditmodel.OperationLog) error
}

type Service struct {
	repo  Repository
	audit AuditRepository
}

func NewService(repo Repository, audit AuditRepository) *Service {
	return &Service{repo: repo, audit: audit}
}

func (s *Service) List(ctx context.Context, q usermodel.ListQuery) (usermodel.ListResult, error) {
	q = normalizeListQuery(q)
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return usermodel.ListResult{}, err
	}
	if total == 0 {
		return usermodel.ListResult{Items: []usermodel.User{}, Page: q.Page, Size: q.PageSize}, nil
	}
	items, err := s.repo.List(ctx, q)
	if err != nil {
		return usermodel.ListResult{}, err
	}
	return usermodel.ListResult{
		Total: total,
		Items: items,
		Page:  q.Page,
		Size:  q.PageSize,
	}, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*usermodel.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListRoles(ctx context.Context) ([]usermodel.Role, error) {
	return s.repo.ListActiveRoles(ctx)
}

func (s *Service) Create(ctx context.Context, req usermodel.CreateRequest, actor usermodel.ActorContext) (*usermodel.User, error) {
	req.LoginName = strings.TrimSpace(req.LoginName)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if !loginNamePattern.MatchString(req.LoginName) {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/users", http.StatusBadRequest, usermodel.ErrLoginNameInvalid)
		return nil, usermodel.ErrLoginNameInvalid
	}
	if err := validatePassword(req.Password); err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/users", http.StatusBadRequest, err)
		return nil, err
	}
	if req.UserStatus != usermodel.StatusActive && req.UserStatus != usermodel.StatusFrozen {
		req.UserStatus = usermodel.StatusActive
	}
	if req.DisplayName == "" {
		req.DisplayName = req.LoginName
	}

	if existing, err := s.repo.GetByLoginName(ctx, req.LoginName); err != nil && !errors.Is(err, usermodel.ErrUserNotFound) {
		return nil, err
	} else if existing != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/users", http.StatusConflict, usermodel.ErrLoginNameTaken)
		return nil, usermodel.ErrLoginNameTaken
	}

	// 选了主角色就必须把它写进 gbp_user_roles，否则用户登录后没有任何菜单/按钮权限。
	req.RoleIDs = mergePrimaryIntoRoles(req.RoleIDs, req.PrimaryRoleID)

	hash, err := authservice.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	uuid, err := generateUUIDv4()
	if err != nil {
		return nil, err
	}
	id, err := s.repo.Create(ctx, uuid, req, hash)
	if err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/users", http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPost, "/api/v1/users", http.StatusOK, fmt.Sprintf("created user_id=%d", id))
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint64, req usermodel.UpdateRequest, actor usermodel.ActorContext) (*usermodel.User, error) {
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.DisplayName == "" {
		req.DisplayName = target.DisplayName
	}
	if err := s.repo.Update(ctx, id, req); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", id), http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d", id), http.StatusOK, "")
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint64, actor usermodel.ActorContext) error {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if isProtectedAdmin(target) {
		s.logFailure(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", id), http.StatusForbidden, usermodel.ErrAdminProtected)
		return usermodel.ErrAdminProtected
	}
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		s.logFailure(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/users/%d", id), http.StatusOK, "")
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, id uint64, status int, actor usermodel.ActorContext) error {
	if status != usermodel.StatusActive && status != usermodel.StatusFrozen {
		return usermodel.ErrInvalidStatus
	}
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if isProtectedAdmin(target) && status == usermodel.StatusFrozen {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/status", id), http.StatusForbidden, usermodel.ErrAdminProtected)
		return usermodel.ErrAdminProtected
	}
	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/status", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/status", id), http.StatusOK, fmt.Sprintf("status=%d", status))
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, id uint64, req usermodel.ResetPasswordRequest, actor usermodel.ActorContext) error {
	if err := validatePassword(req.Password); err != nil {
		return err
	}
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	hash, err := authservice.HashPassword(req.Password)
	if err != nil {
		return err
	}
	mustChange := true
	if req.MustChangePassword != nil {
		mustChange = *req.MustChangePassword
	}
	if err := s.repo.UpdatePassword(ctx, id, hash, mustChange); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/password", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/password", id), http.StatusOK, "")
	return nil
}

func (s *Service) AssignRoles(ctx context.Context, id uint64, roleIDs []uint64, actor usermodel.ActorContext) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.ReplaceRoles(ctx, id, dedupeRoles(roleIDs)); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/roles", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/users/%d/roles", id), http.StatusOK, fmt.Sprintf("roles=%v", roleIDs))
	return nil
}

// ── helpers ─────────────────────────────────────────────────────────────

func normalizeListQuery(q usermodel.ListQuery) usermodel.ListQuery {
	if q.Page <= 0 {
		q.Page = defaultPage
	}
	if q.PageSize <= 0 {
		q.PageSize = defaultPageSize
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
	q.Keyword = strings.TrimSpace(q.Keyword)
	return q
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return usermodel.ErrPasswordWeak
	}
	return nil
}

func isProtectedAdmin(u *usermodel.User) bool {
	return u != nil && (u.ID == 1 || strings.EqualFold(u.LoginName, usermodel.BuiltinAdminLogin))
}

func dedupeRoles(ids []uint64) []uint64 {
	seen := map[uint64]struct{}{}
	out := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func mergePrimaryIntoRoles(ids []uint64, primary *uint64) []uint64 {
	merged := dedupeRoles(ids)
	if primary == nil || *primary == 0 {
		return merged
	}
	for _, id := range merged {
		if id == *primary {
			return merged
		}
	}
	return append(merged, *primary)
}

func generateUUIDv4() (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	buf[6] = (buf[6] & 0x0f) | 0x40 // version 4
	buf[8] = (buf[8] & 0x3f) | 0x80 // variant 10
	hexs := hex.EncodeToString(buf[:])
	return hexs[0:8] + "-" + hexs[8:12] + "-" + hexs[12:16] + "-" + hexs[16:20] + "-" + hexs[20:32], nil
}

func (s *Service) logSuccess(ctx context.Context, actor usermodel.ActorContext, method, path string, status int, note string) {
	s.writeAudit(ctx, actor, method, path, status, "", note)
}

func (s *Service) logFailure(ctx context.Context, actor usermodel.ActorContext, method, path string, status int, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	s.writeAudit(ctx, actor, method, path, status, msg, "")
}

func (s *Service) writeAudit(ctx context.Context, actor usermodel.ActorContext, method, path string, status int, errMsg, body string) {
	if s.audit == nil {
		return
	}
	logCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var userID *uint64
	if actor.UserID != 0 {
		v := actor.UserID
		userID = &v
	}
	_ = s.audit.InsertOperationLog(logCtx, auditmodel.OperationLog{
		UserID:        userID,
		SourceIP:      actor.SourceIP,
		RequestMethod: method,
		RequestPath:   path,
		StatusCode:    status,
		UserAgent:     actor.UserAgent,
		ErrorMessage:  errMsg,
		ResponseBody:  body,
	})
}
