package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"gobaseproject/server/internal/apperror"
	authmodel "gobaseproject/server/internal/model/auth"
	authservice "gobaseproject/server/internal/service/auth"
)

type jwtParser interface {
	Parse(tokenValue string, expectedType string) (*authservice.Claims, error)
}

// Permission enforces API-level access control based on gbp_permission_policies.
// Skip rules (gbp_api_skip_rules) are checked first; matching requests bypass auth entirely.
// All other requests require a valid JWT whose embedded role IDs have an allow policy.
func Permission(db *sql.DB, tokens jwtParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			if isSkipped(r.Context(), db, r.Method, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			token := bearerToken(r)
			if token == "" {
				apperror.WriteDefinition(w, r, apperror.MissingAuthToken)
				return
			}

			claims, err := tokens.Parse(token, authmodel.TokenTypeAccess)
			if err != nil {
				apperror.WriteDefinition(w, r, apperror.InvalidAuthToken)
				return
			}

			if len(claims.RoleIDs) == 0 {
				apperror.WriteDefinition(w, r, apperror.Forbidden)
				return
			}

			resourceKey := r.Method + ":" + normalizePath(r.URL.Path)
			if !hasPermission(r.Context(), db, claims.RoleIDs, resourceKey) {
				apperror.WriteDefinition(w, r, apperror.Forbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isSkipped(ctx context.Context, db *sql.DB, method, path string) bool {
	rows, err := db.QueryContext(ctx,
		`SELECT api_path FROM gbp_api_skip_rules WHERE api_method = ? AND deleted_at IS NULL`,
		method,
	)
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var pattern string
		if rows.Scan(&pattern) != nil {
			continue
		}
		if pattern == path || matchPathPattern(pattern, path) {
			return true
		}
	}
	return false
}

// matchPathPattern checks whether path matches a pattern that may contain
// {placeholder} segments (e.g. /api/v1/dictionaries/{dictCode}/items).
func matchPathPattern(pattern, path string) bool {
	ps := strings.Split(pattern, "/")
	ss := strings.Split(path, "/")
	if len(ps) != len(ss) {
		return false
	}
	for i := range ps {
		if strings.HasPrefix(ps[i], "{") && strings.HasSuffix(ps[i], "}") {
			continue
		}
		if ps[i] != ss[i] {
			return false
		}
	}
	return true
}

func hasPermission(ctx context.Context, db *sql.DB, roleIDs []uint64, resourceKey string) bool {
	if len(roleIDs) == 0 {
		return false
	}
	placeholders := make([]string, len(roleIDs))
	args := make([]interface{}, 0, len(roleIDs)+1)
	for i, id := range roleIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	args = append(args, resourceKey)
	query := fmt.Sprintf(
		`SELECT COUNT(1) FROM gbp_permission_policies
		 WHERE subject_type = 'role' AND subject_id IN (%s)
		   AND resource_type = 'api' AND resource_key = ?
		   AND effect = 'allow' AND policy_status = 1 AND deleted_at IS NULL`,
		strings.Join(placeholders, ","),
	)
	var n int
	err := db.QueryRowContext(ctx, query, args...).Scan(&n)
	return err == nil && n > 0
}

// normalizePath replaces all-numeric path segments with {id} so that stored
// patterns like /api/v1/users/{id} match incoming paths like /api/v1/users/42.
func normalizePath(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if isNumeric(seg) {
			segments[i] = "{id}"
		}
	}
	return strings.Join(segments, "/")
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}
