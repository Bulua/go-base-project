package apiservice

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	apimodel "gobaseproject/server/internal/model/api"
	auditmodel "gobaseproject/server/internal/model/audit"
	"gobaseproject/server/pkg/routereg"
)

const (
	defaultPageSize = 20
	maxPageSize     = 200
)

// validMethodSet is derived from apimodel.ValidHTTPMethods for O(1) lookup.
var validMethodSet = func() map[string]bool {
	m := make(map[string]bool, len(apimodel.ValidHTTPMethods))
	for _, v := range apimodel.ValidHTTPMethods {
		m[v] = true
	}
	return m
}()

type repository interface {
	CountAPIs(ctx context.Context, q apimodel.APIListQuery) (int64, error)
	ListAPIs(ctx context.Context, q apimodel.APIListQuery) ([]apimodel.APIResource, error)
	GetAPIByID(ctx context.Context, id uint64) (*apimodel.APIResource, error)
	ExistsAPIPathMethod(ctx context.Context, path, method string, excludeID uint64) (bool, error)
	CreateAPI(ctx context.Context, p apimodel.SaveAPIPayload) (uint64, error)
	UpdateAPI(ctx context.Context, id uint64, p apimodel.SaveAPIPayload) error
	DeleteAPI(ctx context.Context, id uint64) error
	HasPolicyForAPI(ctx context.Context, apiID uint64) (bool, error)
	ListAllAPIGroups(ctx context.Context) ([]string, error)

	CountSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) (int64, error)
	ListSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) ([]apimodel.SkipRule, error)
	GetSkipRuleByID(ctx context.Context, id uint64) (*apimodel.SkipRule, error)
	ExistsSkipRule(ctx context.Context, path, method string, excludeID uint64) (bool, error)
	CreateSkipRule(ctx context.Context, p apimodel.SaveSkipRulePayload) (uint64, error)
	DeleteSkipRule(ctx context.Context, id uint64) error

	UpsertRoutes(ctx context.Context, routes []routereg.Route) error
	GrantSuperAdmin(ctx context.Context) error
}

type AuditRepository interface {
	InsertOperationLog(ctx context.Context, log auditmodel.OperationLog) error
}

type Service struct {
	repo  repository
	audit AuditRepository
}

func NewService(repo repository, audit AuditRepository) *Service {
	return &Service{repo: repo, audit: audit}
}

// ── API Resources ──────────────────────────────────────────────────────────

func (s *Service) ListAPIs(ctx context.Context, q apimodel.APIListQuery) (*apimodel.APIListResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > maxPageSize {
		q.PageSize = defaultPageSize
	}
	total, err := s.repo.CountAPIs(ctx, q)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListAPIs(ctx, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []apimodel.APIResource{}
	}
	return &apimodel.APIListResult{Total: total, Items: items, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *Service) ListAllAPIGroups(ctx context.Context) ([]string, error) {
	return s.repo.ListAllAPIGroups(ctx)
}

func (s *Service) CreateAPI(ctx context.Context, p apimodel.SaveAPIPayload, actor apimodel.ActorContext) (*apimodel.APIResource, error) {
	if err := validateAPI(p); err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/apis", http.StatusBadRequest, err)
		return nil, err
	}
	taken, err := s.repo.ExistsAPIPathMethod(ctx, p.APIPath, p.APIMethod, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/apis", http.StatusConflict, apimodel.ErrAPIPathTaken)
		return nil, apimodel.ErrAPIPathTaken
	}
	id, err := s.repo.CreateAPI(ctx, p)
	if err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/apis", http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPost, "/api/v1/apis", http.StatusOK, fmt.Sprintf("created api_id=%d %s:%s", id, p.APIMethod, p.APIPath))
	return s.repo.GetAPIByID(ctx, id)
}

func (s *Service) UpdateAPI(ctx context.Context, id uint64, p apimodel.SaveAPIPayload, actor apimodel.ActorContext) (*apimodel.APIResource, error) {
	path := fmt.Sprintf("/api/v1/apis/%d", id)
	if err := validateAPI(p); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, path, http.StatusBadRequest, err)
		return nil, err
	}
	existing, err := s.repo.GetAPIByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		s.logFailure(ctx, actor, http.MethodPut, path, http.StatusNotFound, apimodel.ErrAPINotFound)
		return nil, apimodel.ErrAPINotFound
	}
	taken, err := s.repo.ExistsAPIPathMethod(ctx, p.APIPath, p.APIMethod, id)
	if err != nil {
		return nil, err
	}
	if taken {
		s.logFailure(ctx, actor, http.MethodPut, path, http.StatusConflict, apimodel.ErrAPIPathTaken)
		return nil, apimodel.ErrAPIPathTaken
	}
	if err := s.repo.UpdateAPI(ctx, id, p); err != nil {
		s.logFailure(ctx, actor, http.MethodPut, path, http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPut, path, http.StatusOK, fmt.Sprintf("updated %s:%s", p.APIMethod, p.APIPath))
	return s.repo.GetAPIByID(ctx, id)
}

func (s *Service) DeleteAPI(ctx context.Context, id uint64, actor apimodel.ActorContext) error {
	path := fmt.Sprintf("/api/v1/apis/%d", id)
	existing, err := s.repo.GetAPIByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		s.logFailure(ctx, actor, http.MethodDelete, path, http.StatusNotFound, apimodel.ErrAPINotFound)
		return apimodel.ErrAPINotFound
	}
	hasPolicies, err := s.repo.HasPolicyForAPI(ctx, id)
	if err != nil {
		return err
	}
	if hasPolicies {
		s.logFailure(ctx, actor, http.MethodDelete, path, http.StatusConflict, apimodel.ErrAPIHasPolicies)
		return apimodel.ErrAPIHasPolicies
	}
	if err := s.repo.DeleteAPI(ctx, id); err != nil {
		s.logFailure(ctx, actor, http.MethodDelete, path, http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodDelete, path, http.StatusOK, fmt.Sprintf("deleted %s:%s", existing.APIMethod, existing.APIPath))
	return nil
}

// ── Skip Rules ─────────────────────────────────────────────────────────────

func (s *Service) ListSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) (*apimodel.SkipRuleListResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > maxPageSize {
		q.PageSize = defaultPageSize
	}
	total, err := s.repo.CountSkipRules(ctx, q)
	if err != nil {
		return nil, err
	}
	items, err := s.repo.ListSkipRules(ctx, q)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []apimodel.SkipRule{}
	}
	return &apimodel.SkipRuleListResult{Total: total, Items: items, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *Service) CreateSkipRule(ctx context.Context, p apimodel.SaveSkipRulePayload, actor apimodel.ActorContext) (*apimodel.SkipRule, error) {
	if strings.TrimSpace(p.APIPath) == "" {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/api-skip-rules", http.StatusBadRequest, apimodel.ErrSkipPathEmpty)
		return nil, apimodel.ErrSkipPathEmpty
	}
	if !validMethodSet[strings.ToUpper(p.APIMethod)] {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/api-skip-rules", http.StatusBadRequest, apimodel.ErrAPIMethodInvalid)
		return nil, apimodel.ErrAPIMethodInvalid
	}
	p.APIMethod = strings.ToUpper(p.APIMethod)
	taken, err := s.repo.ExistsSkipRule(ctx, p.APIPath, p.APIMethod, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/api-skip-rules", http.StatusConflict, apimodel.ErrSkipPathTaken)
		return nil, apimodel.ErrSkipPathTaken
	}
	id, err := s.repo.CreateSkipRule(ctx, p)
	if err != nil {
		s.logFailure(ctx, actor, http.MethodPost, "/api/v1/api-skip-rules", http.StatusInternalServerError, err)
		return nil, err
	}
	s.logSuccess(ctx, actor, http.MethodPost, "/api/v1/api-skip-rules", http.StatusOK, fmt.Sprintf("created skip_rule_id=%d %s:%s", id, p.APIMethod, p.APIPath))
	return s.repo.GetSkipRuleByID(ctx, id)
}

func (s *Service) DeleteSkipRule(ctx context.Context, id uint64, actor apimodel.ActorContext) error {
	path := fmt.Sprintf("/api/v1/api-skip-rules/%d", id)
	existing, err := s.repo.GetSkipRuleByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		s.logFailure(ctx, actor, http.MethodDelete, path, http.StatusNotFound, apimodel.ErrSkipNotFound)
		return apimodel.ErrSkipNotFound
	}
	if err := s.repo.DeleteSkipRule(ctx, id); err != nil {
		s.logFailure(ctx, actor, http.MethodDelete, path, http.StatusInternalServerError, err)
		return err
	}
	s.logSuccess(ctx, actor, http.MethodDelete, path, http.StatusOK, fmt.Sprintf("deleted %s:%s", existing.APIMethod, existing.APIPath))
	return nil
}

// SyncRoutes upserts all registered routes into gbp_api_resources and ensures
// super_admin (role_id=1) has allow policies for every active API.
func (s *Service) SyncRoutes(ctx context.Context, routes []routereg.Route) error {
	if err := s.repo.UpsertRoutes(ctx, routes); err != nil {
		return err
	}
	return s.repo.GrantSuperAdmin(ctx)
}

func validateAPI(p apimodel.SaveAPIPayload) error {
	if strings.TrimSpace(p.APIPath) == "" {
		return apimodel.ErrAPIPathEmpty
	}
	if !validMethodSet[strings.ToUpper(p.APIMethod)] {
		return apimodel.ErrAPIMethodInvalid
	}
	return nil
}

// ── audit helpers ──────────────────────────────────────────────────────────

func (s *Service) logSuccess(ctx context.Context, actor apimodel.ActorContext, method, path string, status int, note string) {
	s.writeAudit(ctx, actor, method, path, status, "", note)
}

func (s *Service) logFailure(ctx context.Context, actor apimodel.ActorContext, method, path string, status int, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	s.writeAudit(ctx, actor, method, path, status, msg, "")
}

func (s *Service) writeAudit(ctx context.Context, actor apimodel.ActorContext, method, path string, status int, errMsg, body string) {
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
