package dicthandler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"gobaseproject/server/internal/apperror"
	dictmodel "gobaseproject/server/internal/model/dict"
	authservice "gobaseproject/server/internal/service/auth"
	dictservice "gobaseproject/server/internal/service/dict"
	"gobaseproject/server/pkg/response"
)

type TokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	service *dictservice.Service
	tokens  TokenParser
}

func NewHandler(service *dictservice.Service, tokens TokenParser) *Handler {
	return &Handler{service: service, tokens: tokens}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// specific sub-paths before /{id}
	mux.HandleFunc("GET /api/v1/dictionaries/{dictCode}/items", h.withAuth(h.listItemsByCode))
	mux.HandleFunc("GET /api/v1/dictionaries", h.withAuth(h.list))
	mux.HandleFunc("POST /api/v1/dictionaries", h.withAuth(h.create))
	mux.HandleFunc("PUT /api/v1/dictionaries/{id}", h.withAuth(h.update))
	mux.HandleFunc("DELETE /api/v1/dictionaries/{id}", h.withAuth(h.delete))
	mux.HandleFunc("POST /api/v1/dictionaries/{id}/items", h.withAuth(h.createItem))
	mux.HandleFunc("PUT /api/v1/dictionary-items/{id}", h.withAuth(h.updateItem))
	mux.HandleFunc("DELETE /api/v1/dictionary-items/{id}", h.withAuth(h.deleteItem))
}

// ── Endpoints ─────────────────────────────────────────────────────────────

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := dictmodel.ListQuery{
		Page:       parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:   parseIntDefault(r.URL.Query().Get("page_size"), 20),
		Keyword:    r.URL.Query().Get("keyword"),
		DictStatus: parseIntDefault(r.URL.Query().Get("dict_status"), 0),
	}
	result, err := h.service.List(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var p struct {
		DictName   string  `json:"dict_name"`
		DictCode   string  `json:"dict_code"`
		DictStatus int     `json:"dict_status"`
		ParentID   uint64  `json:"parent_id"`
		Remark     *string `json:"remark"`
	}
	if err := response.ReadJSON(r, &p); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	d, err := h.service.Create(r.Context(), dictmodel.SaveDictPayload{
		DictName: p.DictName, DictCode: p.DictCode,
		DictStatus: p.DictStatus, ParentID: p.ParentID, Remark: p.Remark,
	})
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, d)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var p struct {
		DictName   string  `json:"dict_name"`
		DictCode   string  `json:"dict_code"`
		DictStatus int     `json:"dict_status"`
		ParentID   uint64  `json:"parent_id"`
		Remark     *string `json:"remark"`
	}
	if err := response.ReadJSON(r, &p); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	d, err := h.service.Update(r.Context(), id, dictmodel.SaveDictPayload{
		DictName: p.DictName, DictCode: p.DictCode,
		DictStatus: p.DictStatus, ParentID: p.ParentID, Remark: p.Remark,
	})
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, d)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

func (h *Handler) listItemsByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("dictCode")
	if code == "" {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	items, err := h.service.ListItemsByCode(r.Context(), code)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, items)
}

func (h *Handler) createItem(w http.ResponseWriter, r *http.Request) {
	dictID, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var p struct {
		ItemLabel  string  `json:"item_label"`
		ItemValue  string  `json:"item_value"`
		ItemExtra  *string `json:"item_extra"`
		ItemStatus int     `json:"item_status"`
		SortNo     int     `json:"sort_no"`
		ParentID   uint64  `json:"parent_id"`
	}
	if err := response.ReadJSON(r, &p); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	item, err := h.service.CreateItem(r.Context(), dictID, dictmodel.SaveItemPayload{
		ItemLabel: p.ItemLabel, ItemValue: p.ItemValue, ItemExtra: p.ItemExtra,
		ItemStatus: p.ItemStatus, SortNo: p.SortNo, ParentID: p.ParentID,
	})
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, item)
}

func (h *Handler) updateItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	var p struct {
		ItemLabel  string  `json:"item_label"`
		ItemValue  string  `json:"item_value"`
		ItemExtra  *string `json:"item_extra"`
		ItemStatus int     `json:"item_status"`
		SortNo     int     `json:"sort_no"`
		ParentID   uint64  `json:"parent_id"`
	}
	if err := response.ReadJSON(r, &p); err != nil {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return
	}
	item, err := h.service.UpdateItem(r.Context(), id, dictmodel.SaveItemPayload{
		ItemLabel: p.ItemLabel, ItemValue: p.ItemValue, ItemExtra: p.ItemExtra,
		ItemStatus: p.ItemStatus, SortNo: p.SortNo, ParentID: p.ParentID,
	})
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, item)
}

func (h *Handler) deleteItem(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r, "id")
	if !ok {
		return
	}
	if err := h.service.DeleteItem(r.Context(), id); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// ── Auth middleware ────────────────────────────────────────────────────────

func (h *Handler) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			apperror.WriteDefinition(w, r, apperror.MissingAuthToken)
			return
		}
		if _, err := h.tokens.ParseAccessToken(r.Context(), token); err != nil {
			apperror.WriteDefinition(w, r, apperror.InvalidAuthToken)
			return
		}
		next(w, r)
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────

func parseID(w http.ResponseWriter, r *http.Request, key string) (uint64, bool) {
	id, err := strconv.ParseUint(r.PathValue(key), 10, 64)
	if err != nil || id == 0 {
		apperror.WriteDefinition(w, r, apperror.InvalidRequestBody)
		return 0, false
	}
	return id, true
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

func bearerToken(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "
	if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}
