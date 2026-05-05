package filehandler

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"gobaseproject/server/internal/apperror"
	authservice "gobaseproject/server/internal/service/auth"
	filemodel "gobaseproject/server/internal/model/file"
	fileservice "gobaseproject/server/internal/service/file"
	"gobaseproject/server/pkg/response"
	"gobaseproject/server/pkg/routereg"
)

const maxMultipartMemory = 32 << 20 // 32 MB in memory, rest on disk

type tokenParser interface {
	ParseAccessToken(ctx context.Context, token string) (*authservice.Claims, error)
}

type Handler struct {
	svc    *fileservice.Service
	tokens tokenParser
}

func NewHandler(svc *fileservice.Service, tokens tokenParser) *Handler {
	return &Handler{svc: svc, tokens: tokens}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/files", h.upload)
	mux.HandleFunc("GET /api/v1/files", h.list)
	mux.HandleFunc("DELETE /api/v1/files/{id}", h.delete)
	mux.HandleFunc("GET /api/v1/files/{id}/raw", h.serve)

	routereg.Add("POST",   "/api/v1/files",          "file", "上传文件")
	routereg.Add("GET",    "/api/v1/files",           "file", "文件列表")
	routereg.Add("DELETE", "/api/v1/files/{id}",      "file", "删除文件")
	routereg.Add("GET",    "/api/v1/files/{id}/raw",  "file", "下载/预览文件")
}

// POST /api/v1/files — multipart upload, field name "file"
func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxMultipartMemory); err != nil {
		apperror.Write(w, r, fmt.Errorf("解析请求失败: %w", err))
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		apperror.Write(w, r, fmt.Errorf("获取文件失败: %w", err))
		return
	}
	file.Close()

	uploaderID := h.extractUserID(r)
	rec, err := h.svc.Upload(r.Context(), header, uploaderID)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, rec)
}

// GET /api/v1/files
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := filemodel.FileListQuery{
		Page:         parseIntDefault(r.URL.Query().Get("page"), 1),
		PageSize:     parseIntDefault(r.URL.Query().Get("page_size"), 10),
		Keyword:      r.URL.Query().Get("keyword"),
		MimeCategory: r.URL.Query().Get("mime_category"),
		StartDate:    r.URL.Query().Get("start_date"),
		EndDate:      r.URL.Query().Get("end_date"),
	}
	result, err := h.svc.List(r.Context(), q)
	if err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, result)
}

// DELETE /api/v1/files/{id}
func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		apperror.Write(w, r, filemodel.ErrFileNotFound)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		apperror.Write(w, r, err)
		return
	}
	response.OK(w, r, map[string]bool{"success": true})
}

// GET /api/v1/files/{id}/raw — serve file content
func (h *Handler) serve(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		http.NotFound(w, r)
		return
	}
	rec, err := h.svc.GetByID(r.Context(), id)
	if err != nil || rec == nil {
		http.NotFound(w, r)
		return
	}

	fullPath := h.svc.StoragePath(rec.StorageKey)

	// Safety: ensure resolved path is under the upload directory
	uploadDir := h.svc.UploadDir()
	cleanPath := filepath.Clean(fullPath)
	cleanDir := filepath.Clean(uploadDir)
	if !strings.HasPrefix(cleanPath, cleanDir+string(filepath.Separator)) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	contentType := rec.MimeType
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(rec.StorageKey))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	disposition := "inline"
	if !strings.HasPrefix(contentType, "image/") && !strings.HasPrefix(contentType, "text/") {
		disposition = fmt.Sprintf(`attachment; filename="%s"`, escapeHeader(rec.OriginalName))
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", disposition)
	http.ServeFile(w, r, fullPath)
}

// ── helpers ────────────────────────────────────────────────────────────────

func (h *Handler) extractUserID(r *http.Request) *uint64 {
	token := bearerToken(r)
	if token == "" {
		return nil
	}
	claims, err := h.tokens.ParseAccessToken(r.Context(), token)
	if err != nil {
		return nil
	}
	v := claims.UserID
	return &v
}

func parseID(r *http.Request) (uint64, bool) {
	// The path is /api/v1/files/{id}/raw or /api/v1/files/{id}
	// r.PathValue("id") works for Go 1.22 ServeMux
	raw := r.PathValue("id")
	if raw == "" {
		// fallback: extract second-to-last segment
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		for i, p := range parts {
			if p == "files" && i+1 < len(parts) {
				raw = parts[i+1]
				break
			}
		}
	}
	n, err := strconv.ParseUint(raw, 10, 64)
	return n, err == nil && n > 0
}

func parseIntDefault(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil && n > 0 {
		return n
	}
	return def
}

func bearerToken(r *http.Request) string {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	const prefix = "Bearer "
	if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

func escapeHeader(s string) string {
	return strings.NewReplacer(`"`, `\"`, `\`, `\\`).Replace(s)
}

