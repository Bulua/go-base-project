package menuhandler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"gobaseproject/server/internal/apperror"
	menumodel "gobaseproject/server/internal/model/menu"
	authservice "gobaseproject/server/internal/service/auth"
	menuservice "gobaseproject/server/internal/service/menu"
	"gobaseproject/server/pkg/codegen"
	"gobaseproject/server/pkg/response"
	"gobaseproject/server/pkg/routereg"
)

type TokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	service    *menuservice.Service
	tokens     TokenParser
	webSrcRoot string
}

func NewHandler(service *menuservice.Service, tokens TokenParser, webSrcRoot string) *Handler {
	return &Handler{service: service, tokens: tokens, webSrcRoot: webSrcRoot}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// specific sub-paths before /{id}
	mux.HandleFunc("GET /api/v1/menus/tree", h.withAuth(h.tree))
	mux.HandleFunc("PUT /api/v1/menu-actions/{id}", h.withAuth(h.updateAction))
	mux.HandleFunc("DELETE /api/v1/menu-actions/{id}", h.withAuth(h.deleteAction))
	mux.HandleFunc("GET /api/v1/menus", h.withAuth(h.list))
	mux.HandleFunc("POST /api/v1/menus", h.withAuth(h.create))
	mux.HandleFunc("GET /api/v1/menus/{id}", h.withAuth(h.get))
	mux.HandleFunc("PUT /api/v1/menus/{id}", h.withAuth(h.update))
	mux.HandleFunc("DELETE /api/v1/menus/{id}", h.withAuth(h.delete))
	mux.HandleFunc("GET /api/v1/menus/{id}/actions", h.withAuth(h.listActions))
	mux.HandleFunc("POST /api/v1/menus/{id}/actions", h.withAuth(h.createAction))

	routereg.Add("GET",    "/api/v1/menus/tree",          "menu", "菜单树")
	routereg.Add("GET",    "/api/v1/menus",               "menu", "菜单列表")
	routereg.Add("POST",   "/api/v1/menus",               "menu", "创建菜单")
	routereg.Add("GET",    "/api/v1/menus/{id}",          "menu", "菜单详情")
	routereg.Add("PUT",    "/api/v1/menus/{id}",          "menu", "修改菜单")
	routereg.Add("DELETE", "/api/v1/menus/{id}",          "menu", "删除菜单")
	routereg.Add("GET",    "/api/v1/menus/{id}/actions",  "menu", "菜单按钮权限")
	routereg.Add("POST",   "/api/v1/menus/{id}/actions",  "menu", "新增按钮权限")
	routereg.Add("PUT",    "/api/v1/menu-actions/{id}",   "menu", "修改按钮权限")
	routereg.Add("DELETE", "/api/v1/menu-actions/{id}",   "menu", "删除按钮权限")
}

// ── Endpoints ─────────────────────────────────────────────────────────────

func (h *Handler) list(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	q := menumodel.ListQuery{
		Page:       parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:   parseIntDefault(r.URL.Query().Get("page_size"), 20),
		Keyword:    r.URL.Query().Get("keyword"),
		MenuStatus: parseIntDefault(r.URL.Query().Get("menu_status"), 0),
		MenuType:   parseIntDefault(r.URL.Query().Get("menu_type"), 0),
	}
	result, err := h.service.List(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

func (h *Handler) tree(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	tree, err := h.service.Tree(r.Context())
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, tree)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	menu, err := h.service.Get(r.Context(), id)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, menu)
}

type createResult struct {
	*menumodel.Menu
	CodeGenerated bool   `json:"code_generated"`
	CodePath      string `json:"code_path,omitempty"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	var req menumodel.SaveRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if req.MenuStatus == 0 {
		req.MenuStatus = menumodel.StatusActive
	}
	menu, err := h.service.Create(r.Context(), req)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}

	result := createResult{Menu: menu}
	if req.ComponentPath != nil && *req.ComponentPath != "" {
		path, created, _ := codegen.GenerateVueFile(h.webSrcRoot, *req.ComponentPath, req.MenuTitle)
		result.CodeGenerated = created
		if created {
			result.CodePath = path
		}
	}
	response.OK(w, r, result)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	var req menumodel.SaveRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if req.MenuStatus == 0 {
		req.MenuStatus = menumodel.StatusActive
	}
	menu, err := h.service.Update(r.Context(), id, req)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, menu)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) listActions(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	actions, err := h.service.ListActions(r.Context(), id)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, actions)
}

func (h *Handler) createAction(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	var req menumodel.SaveActionRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if req.ActionStatus == 0 {
		req.ActionStatus = menumodel.StatusActive
	}
	action, err := h.service.CreateAction(r.Context(), id, req)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, action)
}

func (h *Handler) updateAction(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	var req menumodel.SaveActionRequest
	if err := response.ReadJSON(r, &req); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	if req.ActionStatus == 0 {
		req.ActionStatus = menumodel.StatusActive
	}
	action, err := h.service.UpdateAction(r.Context(), id, req)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, action)
}

func (h *Handler) deleteAction(w http.ResponseWriter, r *http.Request, _ menumodel.ActorContext) {
	id, ok := parseIDFromPath(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteAction(r.Context(), id); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Middleware ────────────────────────────────────────────────────────────

func (h *Handler) withAuth(next func(http.ResponseWriter, *http.Request, menumodel.ActorContext)) http.HandlerFunc {
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
		actor := menumodel.ActorContext{
			UserID:    claims.UserID,
			LoginName: claims.LoginName,
			SourceIP:  clientIP(r),
			UserAgent: r.UserAgent(),
		}
		next(w, r, actor)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

func parseIDFromPath(w http.ResponseWriter, r *http.Request, key string) (uint64, bool) {
	raw := r.PathValue(key)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidMenuID)
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
