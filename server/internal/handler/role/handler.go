package rolehandler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"gobaseproject/server/internal/apperror"
	rolemodel "gobaseproject/server/internal/model/role"
	authservice "gobaseproject/server/internal/service/auth"
	roleservice "gobaseproject/server/internal/service/role"
	"gobaseproject/server/pkg/response"
	"gobaseproject/server/pkg/routereg"
)

type TokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	service *roleservice.Service
	tokens  TokenParser
}

func NewHandler(service *roleservice.Service, tokens TokenParser) *Handler {
	return &Handler{service: service, tokens: tokens}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// 资源目录端点放在 /:id 之前，避免被路径参数吞掉
	mux.HandleFunc("GET /api/v1/roles/tree", h.withAuth(h.tree))
	mux.HandleFunc("GET /api/v1/roles/resources", h.withAuth(h.resources))
	mux.HandleFunc("GET /api/v1/roles", h.withAuth(h.list))
	mux.HandleFunc("POST /api/v1/roles", h.withAuth(h.create))
	mux.HandleFunc("GET /api/v1/roles/{id}", h.withAuth(h.get))
	mux.HandleFunc("PUT /api/v1/roles/{id}", h.withAuth(h.update))
	mux.HandleFunc("DELETE /api/v1/roles/{id}", h.withAuth(h.delete))
	mux.HandleFunc("GET /api/v1/roles/{id}/menus", h.withAuth(h.getMenus))
	mux.HandleFunc("PUT /api/v1/roles/{id}/menus", h.withAuth(h.assignMenus))
	mux.HandleFunc("GET /api/v1/roles/{id}/actions", h.withAuth(h.getActions))
	mux.HandleFunc("PUT /api/v1/roles/{id}/actions", h.withAuth(h.assignActions))
	mux.HandleFunc("GET /api/v1/roles/{id}/apis", h.withAuth(h.getAPIs))
	mux.HandleFunc("PUT /api/v1/roles/{id}/apis", h.withAuth(h.assignAPIs))
	mux.HandleFunc("GET /api/v1/roles/{id}/data-scopes", h.withAuth(h.getDataScopes))
	mux.HandleFunc("PUT /api/v1/roles/{id}/data-scopes", h.withAuth(h.assignDataScopes))

	routereg.Add("GET",    "/api/v1/roles/tree",              "role", "角色树")
	routereg.Add("GET",    "/api/v1/roles/resources",         "role", "权限资源目录")
	routereg.Add("GET",    "/api/v1/roles",                   "role", "角色列表")
	routereg.Add("POST",   "/api/v1/roles",                   "role", "创建角色")
	routereg.Add("GET",    "/api/v1/roles/{id}",              "role", "角色详情")
	routereg.Add("PUT",    "/api/v1/roles/{id}",              "role", "修改角色")
	routereg.Add("DELETE", "/api/v1/roles/{id}",              "role", "删除角色")
	routereg.Add("GET",    "/api/v1/roles/{id}/menus",        "role", "角色菜单")
	routereg.Add("PUT",    "/api/v1/roles/{id}/menus",        "role", "分配菜单")
	routereg.Add("GET",    "/api/v1/roles/{id}/actions",      "role", "角色按钮权限")
	routereg.Add("PUT",    "/api/v1/roles/{id}/actions",      "role", "分配按钮权限")
	routereg.Add("GET",    "/api/v1/roles/{id}/apis",         "role", "角色API权限")
	routereg.Add("PUT",    "/api/v1/roles/{id}/apis",         "role", "分配API权限")
	routereg.Add("GET",    "/api/v1/roles/{id}/data-scopes",  "role", "角色数据范围")
	routereg.Add("PUT",    "/api/v1/roles/{id}/data-scopes",  "role", "分配数据范围")
}

// ── Endpoints ─────────────────────────────────────────────────────────────

func (h *Handler) list(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	q := rolemodel.ListQuery{
		Page:       parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:   parseIntDefault(r.URL.Query().Get("page_size"), 20),
		Keyword:    r.URL.Query().Get("keyword"),
		RoleStatus: parseIntDefault(r.URL.Query().Get("role_status"), 0),
	}
	result, err := h.service.List(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

func (h *Handler) tree(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	tree, err := h.service.Tree(r.Context())
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, tree)
}

func (h *Handler) resources(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	data, err := h.service.ResourceCatalog(r.Context())
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, data)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	role, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, role)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	var req rolemodel.CreateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	role, err := h.service.Create(r.Context(), req, actor)
	if err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, role)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req rolemodel.UpdateRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	role, err := h.service.Update(r.Context(), id, req, actor)
	if err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, role)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id, actor); err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Assignment GET ────────────────────────────────────────────────────────

func (h *Handler) getMenus(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	h.assignedIDs(w, r, h.service.ListMenuIDs)
}
func (h *Handler) getActions(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	h.assignedIDs(w, r, h.service.ListActionIDs)
}
func (h *Handler) getAPIs(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	h.assignedIDs(w, r, h.service.ListAPIIDs)
}
func (h *Handler) getDataScopes(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	h.assignedIDs(w, r, h.service.ListDataScopeIDs)
}

func (h *Handler) assignedIDs(w http.ResponseWriter, r *http.Request, fn func(context.Context, uint64) ([]uint64, error)) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	ids, err := fn(r.Context(), id)
	if err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, map[string]interface{}{"ids": ids})
}

// ── Assignment PUT ────────────────────────────────────────────────────────

func (h *Handler) assignMenus(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	h.assignIDs(w, r, actor, h.service.AssignMenus)
}
func (h *Handler) assignActions(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	h.assignIDs(w, r, actor, h.service.AssignActions)
}
func (h *Handler) assignAPIs(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	h.assignIDs(w, r, actor, h.service.AssignAPIs)
}
func (h *Handler) assignDataScopes(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext) {
	h.assignIDs(w, r, actor, h.service.AssignDataScopes)
}

func (h *Handler) assignIDs(w http.ResponseWriter, r *http.Request, actor rolemodel.ActorContext, fn func(context.Context, uint64, []uint64, rolemodel.ActorContext) error) {
	id, ok := parseIDFromPath(w, r)
	if !ok {
		return
	}
	var req rolemodel.AssignIDsRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if err := fn(r.Context(), id, req.IDs, actor); err != nil {
		writeRoleError(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Middleware ────────────────────────────────────────────────────────────

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, rolemodel.ActorContext)) http.HandlerFunc {
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
		actor := rolemodel.ActorContext{
			UserID:    claims.UserID,
			LoginName: claims.LoginName,
			SourceIP:  clientIP(r),
			UserAgent: r.UserAgent(),
		}
		next(w, r, actor)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

func parseIDFromPath(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	raw := r.PathValue("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidRoleID)
		return 0, false
	}
	return id, true
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

func writeRoleError(w http.ResponseWriter, r *http.Request, err error) {
	apperror.Write(w, r, err)
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	prefix := "Bearer "
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
