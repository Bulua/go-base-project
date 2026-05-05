package apihandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gobaseproject/server/internal/apperror"
	apimodel "gobaseproject/server/internal/model/api"
	authservice "gobaseproject/server/internal/service/auth"
	"gobaseproject/server/pkg/response"
	"gobaseproject/server/pkg/routereg"
)

type apiService interface {
	ListAPIs(ctx context.Context, q apimodel.APIListQuery) (*apimodel.APIListResult, error)
	ListAllAPIGroups(ctx context.Context) ([]string, error)
	CreateAPI(ctx context.Context, p apimodel.SaveAPIPayload, actor apimodel.ActorContext) (*apimodel.APIResource, error)
	UpdateAPI(ctx context.Context, id uint64, p apimodel.SaveAPIPayload, actor apimodel.ActorContext) (*apimodel.APIResource, error)
	DeleteAPI(ctx context.Context, id uint64, actor apimodel.ActorContext) error

	ListSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) (*apimodel.SkipRuleListResult, error)
	CreateSkipRule(ctx context.Context, p apimodel.SaveSkipRulePayload, actor apimodel.ActorContext) (*apimodel.SkipRule, error)
	DeleteSkipRule(ctx context.Context, id uint64, actor apimodel.ActorContext) error
}

type TokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	svc    apiService
	tokens TokenParser
}

func NewHandler(svc apiService, tokens TokenParser) *Handler {
	return &Handler{svc: svc, tokens: tokens}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/apis/groups", h.withAuth(h.listGroups))
	mux.HandleFunc("GET /api/v1/apis", h.withAuth(h.listAPIs))
	mux.HandleFunc("POST /api/v1/apis", h.withAuth(h.createAPI))
	mux.HandleFunc("PUT /api/v1/apis/{id}", h.withAuth(h.updateAPI))
	mux.HandleFunc("DELETE /api/v1/apis/{id}", h.withAuth(h.deleteAPI))

	mux.HandleFunc("GET /api/v1/api-skip-rules", h.withAuth(h.listSkipRules))
	mux.HandleFunc("POST /api/v1/api-skip-rules", h.withAuth(h.createSkipRule))
	mux.HandleFunc("DELETE /api/v1/api-skip-rules/{id}", h.withAuth(h.deleteSkipRule))

	routereg.Add("GET",    "/api/v1/apis/groups",        "api", "API分组列表")
	routereg.Add("GET",    "/api/v1/apis",               "api", "API资源列表")
	routereg.Add("POST",   "/api/v1/apis",               "api", "创建API资源")
	routereg.Add("PUT",    "/api/v1/apis/{id}",          "api", "修改API资源")
	routereg.Add("DELETE", "/api/v1/apis/{id}",          "api", "删除API资源")
	routereg.Add("GET",    "/api/v1/api-skip-rules",     "api", "白名单列表")
	routereg.Add("POST",   "/api/v1/api-skip-rules",     "api", "创建白名单")
	routereg.Add("DELETE", "/api/v1/api-skip-rules/{id}","api", "删除白名单")
}

// GET /api/v1/apis/groups
func (h *Handler) listGroups(w http.ResponseWriter, r *http.Request, _ apimodel.ActorContext) {
	groups, err := h.svc.ListAllAPIGroups(r.Context())
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	if groups == nil {
		groups = []string{}
	}
	response.OK(w, r, groups)
}

// GET /api/v1/apis
func (h *Handler) listAPIs(w http.ResponseWriter, r *http.Request, _ apimodel.ActorContext) {
	q := apimodel.APIListQuery{
		Page:      parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:  parseIntDefault(r.URL.Query().Get("page_size"), 10),
		Keyword:   r.URL.Query().Get("keyword"),
		APIGroup:  r.URL.Query().Get("api_group"),
		APIMethod: r.URL.Query().Get("api_method"),
		APIStatus: parseIntDefault(r.URL.Query().Get("api_status"), 0),
	}
	result, err := h.svc.ListAPIs(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// POST /api/v1/apis
func (h *Handler) createAPI(w http.ResponseWriter, r *http.Request, actor apimodel.ActorContext) {
	var body apimodel.CreateAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	p := apimodel.SaveAPIPayload{
		APIPath:   strings.TrimSpace(body.APIPath),
		APIMethod: strings.ToUpper(strings.TrimSpace(body.APIMethod)),
		APIGroup:  strings.TrimSpace(body.APIGroup),
		APIDesc:   strings.TrimSpace(body.APIDesc),
		APIStatus: apimodel.StatusActive,
	}
	resource, err := h.svc.CreateAPI(r.Context(), p, actor)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, resource)
}

// PUT /api/v1/apis/{id}
func (h *Handler) updateAPI(w http.ResponseWriter, r *http.Request, actor apimodel.ActorContext) {
	id, err := pathUint64(r, "id")
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidParams)
		return
	}
	var body apimodel.UpdateAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if body.APIStatus == 0 {
		body.APIStatus = apimodel.StatusActive
	}
	p := apimodel.SaveAPIPayload{
		APIPath:   strings.TrimSpace(body.APIPath),
		APIMethod: strings.ToUpper(strings.TrimSpace(body.APIMethod)),
		APIGroup:  strings.TrimSpace(body.APIGroup),
		APIDesc:   strings.TrimSpace(body.APIDesc),
		APIStatus: body.APIStatus,
	}
	resource, err := h.svc.UpdateAPI(r.Context(), id, p, actor)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, resource)
}

// DELETE /api/v1/apis/{id}
func (h *Handler) deleteAPI(w http.ResponseWriter, r *http.Request, actor apimodel.ActorContext) {
	id, err := pathUint64(r, "id")
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidParams)
		return
	}
	if err := h.svc.DeleteAPI(r.Context(), id, actor); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// GET /api/v1/api-skip-rules
func (h *Handler) listSkipRules(w http.ResponseWriter, r *http.Request, _ apimodel.ActorContext) {
	q := apimodel.SkipRuleListQuery{
		Page:      parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:  parseIntDefault(r.URL.Query().Get("page_size"), 10),
		Keyword:   r.URL.Query().Get("keyword"),
		APIMethod: r.URL.Query().Get("api_method"),
	}
	result, err := h.svc.ListSkipRules(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// POST /api/v1/api-skip-rules
func (h *Handler) createSkipRule(w http.ResponseWriter, r *http.Request, actor apimodel.ActorContext) {
	var body apimodel.CreateSkipRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	p := apimodel.SaveSkipRulePayload{
		APIPath:    strings.TrimSpace(body.APIPath),
		APIMethod:  strings.ToUpper(strings.TrimSpace(body.APIMethod)),
		SkipReason: strings.TrimSpace(body.SkipReason),
	}
	rule, err := h.svc.CreateSkipRule(r.Context(), p, actor)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, rule)
}

// DELETE /api/v1/api-skip-rules/{id}
func (h *Handler) deleteSkipRule(w http.ResponseWriter, r *http.Request, actor apimodel.ActorContext) {
	id, err := pathUint64(r, "id")
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidParams)
		return
	}
	if err := h.svc.DeleteSkipRule(r.Context(), id, actor); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Middleware ────────────────────────────────────────────────────────────

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, apimodel.ActorContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			apperror.WriteDefinition(w, r, apperror.MissingAuthToken)
			return
		}
		claims, err := h.tokens.ParseAccessToken(r.Context(), token)
		if err != nil {
			apperror.WriteDefinition(w, r, apperror.InvalidAuthToken)
			return
		}
		actor := apimodel.ActorContext{
			UserID:    claims.UserID,
			LoginName: claims.LoginName,
			SourceIP:  clientIP(r),
			UserAgent: r.UserAgent(),
		}
		next(w, r, actor)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

func pathUint64(r *http.Request, key string) (uint64, error) {
	return strconv.ParseUint(r.PathValue(key), 10, 64)
}

func parseIntDefault(raw string, def int) int {
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(strings.ToLower(header), strings.ToLower(prefix)) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func clientIP(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	if realIP := strings.TrimSpace(r.Header.Get("X-Real-IP")); realIP != "" {
		return realIP
	}
	host := r.RemoteAddr
	if idx := strings.LastIndex(host, ":"); idx > -1 {
		return host[:idx]
	}
	return host
}

