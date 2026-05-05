package apirepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"

	apimodel "gobaseproject/server/internal/model/api"
	"gobaseproject/server/pkg/routereg"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// ── API Resources ──────────────────────────────────────────────────────────

func (r *SQLRepository) CountAPIs(ctx context.Context, q apimodel.APIListQuery) (int64, error) {
	where, args := buildAPIWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_api_resources a "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) ListAPIs(ctx context.Context, q apimodel.APIListQuery) ([]apimodel.APIResource, error) {
	where, args := buildAPIWhere(q)
	offset := (q.Page - 1) * q.PageSize
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT a.id, a.api_path, a.api_method, COALESCE(a.api_group,''), COALESCE(a.api_desc,''), a.api_status, a.created_at, a.updated_at
FROM gbp_api_resources a `+where+`
ORDER BY a.api_group ASC, a.api_path ASC, a.api_method ASC
LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []apimodel.APIResource
	for rows.Next() {
		var a apimodel.APIResource
		if err := rows.Scan(&a.ID, &a.APIPath, &a.APIMethod, &a.APIGroup, &a.APIDesc, &a.APIStatus, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, a)
	}
	return items, rows.Err()
}

func (r *SQLRepository) GetAPIByID(ctx context.Context, id uint64) (*apimodel.APIResource, error) {
	var a apimodel.APIResource
	err := r.db.QueryRowContext(ctx, `
SELECT id, api_path, api_method, COALESCE(api_group,''), COALESCE(api_desc,''), api_status, created_at, updated_at
FROM gbp_api_resources WHERE id = ? AND deleted_at IS NULL`, id).Scan(
		&a.ID, &a.APIPath, &a.APIMethod, &a.APIGroup, &a.APIDesc, &a.APIStatus, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func (r *SQLRepository) ExistsAPIPathMethod(ctx context.Context, path, method string, excludeID uint64) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM gbp_api_resources WHERE api_path = ? AND api_method = ? AND deleted_at IS NULL AND id != ?`,
		path, method, excludeID).Scan(&n)
	return n > 0, err
}

func (r *SQLRepository) CreateAPI(ctx context.Context, p apimodel.SaveAPIPayload) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_api_resources (api_path, api_method, api_group, api_desc, api_status)
VALUES (?, ?, ?, ?, ?)`,
		p.APIPath, p.APIMethod, p.APIGroup, p.APIDesc, p.APIStatus)
	if err != nil {
		if isDuplicateKey(err) {
			return 0, apimodel.ErrAPIPathTaken
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) UpdateAPI(ctx context.Context, id uint64, p apimodel.SaveAPIPayload) error {
	_, err := r.db.ExecContext(ctx, `
UPDATE gbp_api_resources SET api_path=?, api_method=?, api_group=?, api_desc=?, api_status=? WHERE id=? AND deleted_at IS NULL`,
		p.APIPath, p.APIMethod, p.APIGroup, p.APIDesc, p.APIStatus, id)
	return err
}

func (r *SQLRepository) DeleteAPI(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE gbp_api_resources SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	return err
}

func (r *SQLRepository) ListAllAPIGroups(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT api_group FROM gbp_api_resources WHERE api_group IS NOT NULL AND api_group != '' AND deleted_at IS NULL ORDER BY api_group`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var groups []string
	for rows.Next() {
		var g string
		if rows.Scan(&g) == nil {
			groups = append(groups, g)
		}
	}
	return groups, rows.Err()
}

func buildAPIWhere(q apimodel.APIListQuery) (string, []interface{}) {
	var conds []string
	var args []interface{}
	conds = append(conds, "a.deleted_at IS NULL")
	if q.Keyword != "" {
		conds = append(conds, "(a.api_path LIKE ? OR a.api_desc LIKE ?)")
		like := "%" + escapeLike(q.Keyword) + "%"
		args = append(args, like, like)
	}
	if q.APIGroup != "" {
		conds = append(conds, "a.api_group = ?")
		args = append(args, q.APIGroup)
	}
	if q.APIMethod != "" {
		conds = append(conds, "a.api_method = ?")
		args = append(args, strings.ToUpper(q.APIMethod))
	}
	if q.APIStatus > 0 {
		conds = append(conds, "a.api_status = ?")
		args = append(args, q.APIStatus)
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

// ── Skip Rules ─────────────────────────────────────────────────────────────

func (r *SQLRepository) CountSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) (int64, error) {
	where, args := buildSkipWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_api_skip_rules s "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) ListSkipRules(ctx context.Context, q apimodel.SkipRuleListQuery) ([]apimodel.SkipRule, error) {
	where, args := buildSkipWhere(q)
	offset := (q.Page - 1) * q.PageSize
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT s.id, s.api_path, s.api_method, COALESCE(s.skip_reason,''), s.created_at
FROM gbp_api_skip_rules s `+where+`
ORDER BY s.api_method ASC, s.api_path ASC
LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []apimodel.SkipRule
	for rows.Next() {
		var s apimodel.SkipRule
		if err := rows.Scan(&s.ID, &s.APIPath, &s.APIMethod, &s.SkipReason, &s.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

func (r *SQLRepository) ExistsSkipRule(ctx context.Context, path, method string, excludeID uint64) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM gbp_api_skip_rules WHERE api_path = ? AND api_method = ? AND deleted_at IS NULL AND id != ?`,
		path, method, excludeID).Scan(&n)
	return n > 0, err
}

func (r *SQLRepository) CreateSkipRule(ctx context.Context, p apimodel.SaveSkipRulePayload) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_api_skip_rules (api_path, api_method, skip_reason) VALUES (?, ?, ?)`,
		p.APIPath, p.APIMethod, p.SkipReason)
	if err != nil {
		if isDuplicateKey(err) {
			return 0, apimodel.ErrSkipPathTaken
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) HasPolicyForAPI(ctx context.Context, apiID uint64) (bool, error) {
	var n int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM gbp_permission_policies WHERE resource_type = 'api' AND resource_id = ? AND deleted_at IS NULL`,
		apiID).Scan(&n)
	return n > 0, err
}

func isDuplicateKey(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *SQLRepository) GetSkipRuleByID(ctx context.Context, id uint64) (*apimodel.SkipRule, error) {
	var s apimodel.SkipRule
	err := r.db.QueryRowContext(ctx, `
SELECT id, api_path, api_method, COALESCE(skip_reason,''), created_at
FROM gbp_api_skip_rules WHERE id = ? AND deleted_at IS NULL`, id).Scan(
		&s.ID, &s.APIPath, &s.APIMethod, &s.SkipReason, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &s, err
}

func (r *SQLRepository) DeleteSkipRule(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE gbp_api_skip_rules SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	return err
}

func buildSkipWhere(q apimodel.SkipRuleListQuery) (string, []interface{}) {
	var conds []string
	var args []interface{}
	conds = append(conds, "s.deleted_at IS NULL")
	if q.Keyword != "" {
		conds = append(conds, "(s.api_path LIKE ? OR s.skip_reason LIKE ?)")
		like := "%" + escapeLike(q.Keyword) + "%"
		args = append(args, like, like)
	}
	if q.APIMethod != "" {
		conds = append(conds, "s.api_method = ?")
		args = append(args, strings.ToUpper(q.APIMethod))
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// ── Route sync ─────────────────────────────────────────────────────────────

// UpsertRoutes inserts or updates API resource records from the route registry.
// Existing rows are matched on (api_path, api_method); group/desc are refreshed
// and deleted_at is cleared so soft-deleted entries come back automatically.
func (r *SQLRepository) UpsertRoutes(ctx context.Context, routes []routereg.Route) error {
	if len(routes) == 0 {
		return nil
	}
	placeholders := make([]string, len(routes))
	args := make([]interface{}, 0, len(routes)*4)
	for i, rt := range routes {
		placeholders[i] = "(?, ?, ?, ?, 1)"
		args = append(args, rt.Path, rt.Method, rt.Group, rt.Desc)
	}
	query := `INSERT INTO gbp_api_resources (api_path, api_method, api_group, api_desc, api_status) VALUES ` +
		strings.Join(placeholders, ",") +
		` ON DUPLICATE KEY UPDATE api_group=VALUES(api_group), api_desc=VALUES(api_desc), api_status=1, deleted_at=NULL`
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

// GrantSuperAdmin ensures role_id=1 has an allow policy for every active API resource.
func (r *SQLRepository) GrantSuperAdmin(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_permission_policies
  (subject_type, subject_id, resource_type, resource_id, resource_key, action, effect, policy_status)
SELECT 'role', 1, 'api', a.id, CONCAT(a.api_method, ':', a.api_path), a.api_method, 'allow', 1
FROM gbp_api_resources a
WHERE a.deleted_at IS NULL
ON DUPLICATE KEY UPDATE deleted_at=NULL, policy_status=1`)
	return err
}
