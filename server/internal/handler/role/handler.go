package rolehandler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	rolemodel "gobaseproject/server/internal/model/role"
	authservice "gobaseproject/server/internal/service/auth"
	roleservice "gobaseproject/server/internal/service/role"
	"gobaseproject/server/pkg/response"
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
		response.Error(w, r, http.StatusInternalServerError, 500, "list roles failed: "+err.Error())
		return
	}
	response.OK(w, r, result)
}

func (h *Handler) tree(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	tree, err := h.service.Tree(r.Context())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "load role tree failed: "+err.Error())
		return
	}
	response.OK(w, r, tree)
}

func (h *Handler) resources(w http.ResponseWriter, r *http.Request, _ rolemodel.ActorContext) {
	data, err := h.service.ResourceCatalog(r.Context())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, 500, "load resources failed: "+err.Error())
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
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
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
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
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
		response.Error(w, r, http.StatusBadRequest, 400, "invalid request body")
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
			response.Error(w, r, http.StatusUnauthorized, 401, "missing authorization token")
			return
		}
		claims, err := h.tokens.ParseAccessToken(r.Context(), token)
		if err != nil {
			response.Error(w, r, http.StatusUnauthorized, 401, "invalid or expired token")
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
		response.Error(w, r, http.StatusBadRequest, 400, "invalid role id")
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
	switch {
	case errors.Is(err, rolemodel.ErrRoleNotFound):
		response.Error(w, r, http.StatusNotFound, 404, "role not found")
	case errors.Is(err, rolemodel.ErrRoleCodeTaken):
		response.Error(w, r, http.StatusConflict, 409, "role code already exists")
	case errors.Is(err, rolemodel.ErrRoleCodeInvalid):
		response.Error(w, r, http.StatusBadRequest, 400, "role code must start with a lowercase letter and be 3-64 chars of [a-z0-9_]")
	case errors.Is(err, rolemodel.ErrInvalidStatus):
		response.Error(w, r, http.StatusBadRequest, 400, "invalid role status")
	case errors.Is(err, rolemodel.ErrBuiltinProtect):
		response.Error(w, r, http.StatusForbidden, 403, "built-in role is protected")
	case errors.Is(err, rolemodel.ErrRoleHasUsers):
		response.Error(w, r, http.StatusConflict, 409, "role still has users assigned, deassign them first")
	case errors.Is(err, rolemodel.ErrRoleHasChildren):
		response.Error(w, r, http.StatusConflict, 409, "role still has child roles, remove them first")
	case errors.Is(err, rolemodel.ErrParentLoop):
		response.Error(w, r, http.StatusBadRequest, 400, "parent role would create a cycle")
	default:
		response.Error(w, r, http.StatusInternalServerError, 500, "internal server error: "+err.Error())
	}
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
