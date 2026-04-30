package roleservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	auditmodel "gobaseproject/server/internal/model/audit"
	rolemodel "gobaseproject/server/internal/model/role"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100

	defaultRouteFallback = "dashboard"
)

var roleCodePattern = regexp.MustCompile(`^[a-z][a-z0-9_]{2,63}$`)

type Repository interface {
	Count(ctx context.Context, q rolemodel.ListQuery) (int64, error)
	List(ctx context.Context, q rolemodel.ListQuery) ([]rolemodel.Role, error)
	ListAll(ctx context.Context) ([]rolemodel.Role, error)
	GetByID(ctx context.Context, id uint64) (*rolemodel.Role, error)
	GetByCode(ctx context.Context, code string) (*rolemodel.Role, error)
	Create(ctx context.Context, req rolemodel.CreateRequest) (uint64, error)
	Update(ctx context.Context, id uint64, req rolemodel.UpdateRequest) error
	SoftDelete(ctx context.Context, id uint64) error
	CountUsers(ctx context.Context, roleID uint64) (int64, error)
	CountChildren(ctx context.Context, roleID uint64) (int64, error)
	ListMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error)
	ListActionIDs(ctx context.Context, roleID uint64) ([]uint64, error)
	ListAPIIDs(ctx context.Context, roleID uint64) ([]uint64, error)
	ListDataScopeIDs(ctx context.Context, roleID uint64) ([]uint64, error)
	ReplaceMenus(ctx context.Context, roleID uint64, ids []uint64) error
	ReplaceActions(ctx context.Context, roleID uint64, ids []uint64) error
	ReplaceAPIs(ctx context.Context, roleID uint64, ids []uint64) error
	ReplaceDataScopes(ctx context.Context, roleID uint64, ids []uint64) error
	ListAllMenus(ctx context.Context) ([]rolemodel.MenuOption, error)
	ListAllActions(ctx context.Context) ([]rolemodel.ActionOption, error)
	ListAllAPIs(ctx context.Context) ([]rolemodel.APIOption, error)
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

// ── Read ─────────────────────────────────────────────────────────────────

func (s *Service) List(ctx context.Context, q rolemodel.ListQuery) (rolemodel.ListResult, error) {
	q = normalizeListQuery(q)
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return rolemodel.ListResult{}, err
	}
	if total == 0 {
		return rolemodel.ListResult{Items: []rolemodel.Role{}, Page: q.Page, Size: q.PageSize}, nil
	}
	items, err := s.repo.List(ctx, q)
	if err != nil {
		return rolemodel.ListResult{}, err
	}
	return rolemodel.ListResult{Total: total, Items: items, Page: q.Page, Size: q.PageSize}, nil
}

func (s *Service) Tree(ctx context.Context) ([]rolemodel.Role, error) {
	roles, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return buildRoleTree(roles), nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*rolemodel.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListMenuIDs(ctx, roleID)
}

func (s *Service) ListActionIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListActionIDs(ctx, roleID)
}

func (s *Service) ListAPIIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListAPIIDs(ctx, roleID)
}

func (s *Service) ListDataScopeIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	if _, err := s.repo.GetByID(ctx, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListDataScopeIDs(ctx, roleID)
}

func (s *Service) ResourceCatalog(ctx context.Context) (map[string]interface{}, error) {
	menus, err := s.repo.ListAllMenus(ctx)
	if err != nil {
		return nil, err
	}
	actions, err := s.repo.ListAllActions(ctx)
	if err != nil {
		return nil, err
	}
	apis, err := s.repo.ListAllAPIs(ctx)
	if err != nil {
		return nil, err
	}
	roles, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"menus":   buildMenuTree(menus),
		"actions": actions,
		"apis":    apis,
		"roles":   roles,
	}, nil
}

// ── Mutations ────────────────────────────────────────────────────────────

func (s *Service) Create(ctx context.Context, req rolemodel.CreateRequest, actor rolemodel.ActorContext) (*rolemodel.Role, error) {
	req.RoleCode = strings.TrimSpace(strings.ToLower(req.RoleCode))
	req.RoleName = strings.TrimSpace(req.RoleName)
	if !roleCodePattern.MatchString(req.RoleCode) {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/roles", http.StatusBadRequest, rolemodel.ErrRoleCodeInvalid)
		return nil, rolemodel.ErrRoleCodeInvalid
	}
	if req.RoleName == "" {
		req.RoleName = req.RoleCode
	}
	if req.DefaultRoute == "" {
		req.DefaultRoute = defaultRouteFallback
	}
	if req.RoleStatus != rolemodel.StatusActive && req.RoleStatus != rolemodel.StatusDisabled {
		req.RoleStatus = rolemodel.StatusActive
	}

	if existing, err := s.repo.GetByCode(ctx, req.RoleCode); err != nil && !errors.Is(err, rolemodel.ErrRoleNotFound) {
		return nil, err
	} else if existing != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/roles", http.StatusConflict, rolemodel.ErrRoleCodeTaken)
		return nil, rolemodel.ErrRoleCodeTaken
	}

	if req.ParentRoleID > 0 {
		if _, err := s.repo.GetByID(ctx, req.ParentRoleID); err != nil {
			return nil, err
		}
	}

	id, err := s.repo.Create(ctx, req)
	if err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/roles", http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPost, "/api/v1/roles", http.StatusOK, fmt.Sprintf("created role_id=%d", id))
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id uint64, req rolemodel.UpdateRequest, actor rolemodel.ActorContext) (*rolemodel.Role, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	req.RoleCode = strings.TrimSpace(strings.ToLower(req.RoleCode))
	if req.RoleCode == "" {
		req.RoleCode = target.RoleCode
	}
	req.RoleName = strings.TrimSpace(req.RoleName)
	if req.RoleName == "" {
		req.RoleName = target.RoleName
	}
	if req.DefaultRoute == "" {
		req.DefaultRoute = target.DefaultRoute
	}
	if req.RoleStatus != rolemodel.StatusActive && req.RoleStatus != rolemodel.StatusDisabled {
		req.RoleStatus = target.RoleStatus
	}
	if req.ParentRoleID > 0 {
		if req.ParentRoleID == id {
			return nil, rolemodel.ErrParentLoop
		}
		if descendantOf(ctx, s.repo, req.ParentRoleID, id) {
			return nil, rolemodel.ErrParentLoop
		}
	}
	if isBuiltinRole(target) && req.RoleStatus == rolemodel.StatusDisabled {
		// 内置角色（id=1）允许改名/改 code，但禁止禁用——否则当前管理员可能直接锁掉自己。
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusForbidden, rolemodel.ErrBuiltinProtect)
		return nil, rolemodel.ErrBuiltinProtect
	}
	if req.RoleCode != target.RoleCode {
		if !roleCodePattern.MatchString(req.RoleCode) {
			return nil, rolemodel.ErrRoleCodeInvalid
		}
		if existing, err := s.repo.GetByCode(ctx, req.RoleCode); err != nil && !errors.Is(err, rolemodel.ErrRoleNotFound) {
			return nil, err
		} else if existing != nil && existing.ID != id {
			s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusConflict, rolemodel.ErrRoleCodeTaken)
			return nil, rolemodel.ErrRoleCodeTaken
		}
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusOK, "")
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id uint64, actor rolemodel.ActorContext) error {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if isBuiltinRole(target) {
		s.logFailure(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusForbidden, rolemodel.ErrBuiltinProtect)
		return rolemodel.ErrBuiltinProtect
	}
	// 用户/子角色不再硬阻塞——前端会做警告确认。软删后：
	//   - gbp_user_roles 中关联记录残留，但 ListRolesByUserID 通过 r.deleted_at IS NULL 过滤跳过，用户自动失去该角色
	//   - 子角色 parent_role_id 指向已软删的父，相当于变成顶级角色（前端会显示"—"）
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		s.logFailure(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodDelete, fmt.Sprintf("/api/v1/roles/%d", id), http.StatusOK, "")
	return nil
}

func (s *Service) AssignMenus(ctx context.Context, id uint64, menuIDs []uint64, actor rolemodel.ActorContext) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.ReplaceMenus(ctx, id, dedupe(menuIDs)); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/menus", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/menus", id), http.StatusOK, fmt.Sprintf("count=%d", len(menuIDs)))
	return nil
}

func (s *Service) AssignActions(ctx context.Context, id uint64, actionIDs []uint64, actor rolemodel.ActorContext) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.ReplaceActions(ctx, id, dedupe(actionIDs)); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/actions", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/actions", id), http.StatusOK, fmt.Sprintf("count=%d", len(actionIDs)))
	return nil
}

func (s *Service) AssignAPIs(ctx context.Context, id uint64, apiIDs []uint64, actor rolemodel.ActorContext) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.ReplaceAPIs(ctx, id, dedupe(apiIDs)); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/apis", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/apis", id), http.StatusOK, fmt.Sprintf("count=%d", len(apiIDs)))
	return nil
}

func (s *Service) AssignDataScopes(ctx context.Context, id uint64, visibleRoleIDs []uint64, actor rolemodel.ActorContext) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return err
	}
	if err := s.repo.ReplaceDataScopes(ctx, id, dedupe(visibleRoleIDs)); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/data-scopes", id), http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodPut, fmt.Sprintf("/api/v1/roles/%d/data-scopes", id), http.StatusOK, fmt.Sprintf("count=%d", len(visibleRoleIDs)))
	return nil
}

// ── helpers ──────────────────────────────────────────────────────────────

func normalizeListQuery(q rolemodel.ListQuery) rolemodel.ListQuery {
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

func isBuiltinRole(role *rolemodel.Role) bool {
	return role != nil && role.ID == rolemodel.BuiltinRoleID
}

func dedupe(ids []uint64) []uint64 {
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

func buildRoleTree(roles []rolemodel.Role) []rolemodel.Role {
	byParent := map[uint64][]rolemodel.Role{}
	for _, r := range roles {
		byParent[r.ParentRoleID] = append(byParent[r.ParentRoleID], r)
	}
	var attach func(parent uint64) []rolemodel.Role
	attach = func(parent uint64) []rolemodel.Role {
		items := byParent[parent]
		for i := range items {
			items[i].Children = attach(items[i].ID)
		}
		return items
	}
	return attach(0)
}

func buildMenuTree(menus []rolemodel.MenuOption) []rolemodel.MenuOption {
	byParent := map[uint64][]rolemodel.MenuOption{}
	for _, m := range menus {
		byParent[m.ParentID] = append(byParent[m.ParentID], m)
	}
	var attach func(parent uint64) []rolemodel.MenuOption
	attach = func(parent uint64) []rolemodel.MenuOption {
		items := byParent[parent]
		for i := range items {
			items[i].Children = attach(items[i].ID)
		}
		return items
	}
	return attach(0)
}

// descendantOf reports whether candidate is in the descendant chain of root.
// Used to prevent setting parent_role_id to a role that's already this role's child.
func descendantOf(ctx context.Context, repo Repository, candidate uint64, root uint64) bool {
	cur := candidate
	for depth := 0; depth < 32 && cur != 0; depth++ {
		if cur == root {
			return true
		}
		role, err := repo.GetByID(ctx, cur)
		if err != nil {
			return false
		}
		cur = role.ParentRoleID
	}
	return false
}

func (s *Service) logSuccess(ctx context.Context, actor rolemodel.ActorContext, method, path string, status int, note string) {
	s.writeAudit(ctx, actor, method, path, status, "", note)
}

func (s *Service) logFailure(ctx context.Context, actor rolemodel.ActorContext, method, path string, status int, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	s.writeAudit(ctx, actor, method, path, status, msg, "")
}

func (s *Service) writeAudit(ctx context.Context, actor rolemodel.ActorContext, method, path string, status int, errMsg, body string) {
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
