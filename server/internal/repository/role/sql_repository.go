package rolerepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	rolemodel "gobaseproject/server/internal/model/role"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// ── List / Get ───────────────────────────────────────────────────────────

func (r *SQLRepository) Count(ctx context.Context, q rolemodel.ListQuery) (int64, error) {
	where, args := buildListWhere(q)
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_roles "+where, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SQLRepository) List(ctx context.Context, q rolemodel.ListQuery) ([]rolemodel.Role, error) {
	where, args := buildListWhere(q)
	query := `SELECT id, role_code, role_name, parent_role_id, default_route, sort_no, role_status, remark, created_at, updated_at
FROM gbp_roles ` + where + ` ORDER BY sort_no ASC, id ASC LIMIT ? OFFSET ?`
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []rolemodel.Role{}
	ids := []uint64{}
	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
		ids = append(ids, role.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return roles, nil
	}
	counts, err := r.userCounts(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range roles {
		roles[i].UserCount = counts[roles[i].ID]
	}
	return roles, nil
}

func (r *SQLRepository) ListAll(ctx context.Context) ([]rolemodel.Role, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, role_code, role_name, parent_role_id, default_route, sort_no, role_status, remark, created_at, updated_at
FROM gbp_roles WHERE deleted_at IS NULL
ORDER BY sort_no ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []rolemodel.Role{}
	for rows.Next() {
		role, err := scanRoleRow(rows)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *SQLRepository) GetByID(ctx context.Context, id uint64) (*rolemodel.Role, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, role_code, role_name, parent_role_id, default_route, sort_no, role_status, remark, created_at, updated_at
FROM gbp_roles WHERE id = ? AND deleted_at IS NULL`, id)
	role, err := scanRoleRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rolemodel.ErrRoleNotFound
		}
		return nil, err
	}
	counts, err := r.userCounts(ctx, []uint64{id})
	if err != nil {
		return nil, err
	}
	role.UserCount = counts[id]
	return &role, nil
}

func (r *SQLRepository) GetByCode(ctx context.Context, code string) (*rolemodel.Role, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, role_code, role_name, parent_role_id, default_route, sort_no, role_status, remark, created_at, updated_at
FROM gbp_roles WHERE role_code = ? AND deleted_at IS NULL`, code)
	role, err := scanRoleRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, rolemodel.ErrRoleNotFound
		}
		return nil, err
	}
	return &role, nil
}

// ── Create / Update / Delete ─────────────────────────────────────────────

func (r *SQLRepository) Create(ctx context.Context, req rolemodel.CreateRequest) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_roles
  (role_code, role_name, parent_role_id, default_route, sort_no, role_status, remark)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		req.RoleCode,
		req.RoleName,
		req.ParentRoleID,
		req.DefaultRoute,
		req.SortNo,
		req.RoleStatus,
		nullStringPtr(req.Remark),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func (r *SQLRepository) Update(ctx context.Context, id uint64, req rolemodel.UpdateRequest) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_roles SET role_code = ?, role_name = ?, parent_role_id = ?, default_route = ?,
  sort_no = ?, role_status = ?, remark = ?
WHERE id = ? AND deleted_at IS NULL`,
		req.RoleCode,
		req.RoleName,
		req.ParentRoleID,
		req.DefaultRoute,
		req.SortNo,
		req.RoleStatus,
		nullStringPtr(req.Remark),
		id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return rolemodel.ErrRoleNotFound
	}
	return nil
}

func (r *SQLRepository) SoftDelete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_roles SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return rolemodel.ErrRoleNotFound
	}
	return nil
}

func (r *SQLRepository) CountUsers(ctx context.Context, roleID uint64) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(1) FROM gbp_user_roles ur
JOIN gbp_users u ON u.id = ur.user_id
WHERE ur.role_id = ? AND u.deleted_at IS NULL`, roleID).Scan(&n)
	return n, err
}

func (r *SQLRepository) CountChildren(ctx context.Context, roleID uint64) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(1) FROM gbp_roles WHERE parent_role_id = ? AND deleted_at IS NULL`, roleID).Scan(&n)
	return n, err
}

// ── Assignment reads ─────────────────────────────────────────────────────

func (r *SQLRepository) ListMenuIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return r.queryIDs(ctx, `SELECT menu_id FROM gbp_role_menus WHERE role_id = ? ORDER BY menu_id ASC`, roleID)
}

func (r *SQLRepository) ListActionIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return r.queryIDs(ctx, `SELECT action_id FROM gbp_role_actions WHERE role_id = ? ORDER BY action_id ASC`, roleID)
}

func (r *SQLRepository) ListAPIIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return r.queryIDs(ctx, `
SELECT resource_id FROM gbp_permission_policies
WHERE subject_type = 'role' AND subject_id = ? AND resource_type = 'api'
  AND resource_id IS NOT NULL AND deleted_at IS NULL
ORDER BY resource_id ASC`, roleID)
}

func (r *SQLRepository) ListDataScopeIDs(ctx context.Context, roleID uint64) ([]uint64, error) {
	return r.queryIDs(ctx, `SELECT visible_role_id FROM gbp_role_data_scopes WHERE role_id = ? ORDER BY visible_role_id ASC`, roleID)
}

// ── Assignment writes ────────────────────────────────────────────────────

func (r *SQLRepository) ReplaceMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	return r.replaceLinks(ctx,
		`DELETE FROM gbp_role_menus WHERE role_id = ?`,
		"INSERT INTO gbp_role_menus (role_id, menu_id) VALUES ",
		roleID, menuIDs,
	)
}

func (r *SQLRepository) ReplaceActions(ctx context.Context, roleID uint64, actionIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM gbp_role_actions WHERE role_id = ?`, roleID); err != nil {
		return err
	}
	if len(actionIDs) > 0 {
		// gbp_role_actions has (role_id, menu_id, action_id) PK; we need menu_id from gbp_menu_actions.
		placeholders := make([]string, len(actionIDs))
		args := make([]interface{}, 0, len(actionIDs)*1+1)
		args = append(args, roleID)
		for i, id := range actionIDs {
			placeholders[i] = "?"
			args = append(args, id)
		}
		query := fmt.Sprintf(`
INSERT INTO gbp_role_actions (role_id, menu_id, action_id)
SELECT ?, menu_id, id FROM gbp_menu_actions
WHERE id IN (%s) AND deleted_at IS NULL`, strings.Join(placeholders, ","))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SQLRepository) ReplaceAPIs(ctx context.Context, roleID uint64, apiIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `
DELETE FROM gbp_permission_policies
WHERE subject_type = 'role' AND subject_id = ? AND resource_type = 'api'`, roleID); err != nil {
		return err
	}
	for _, apiID := range apiIDs {
		// resource_key = METHOD:/path; action mirrors the method.
		var apiPath, apiMethod string
		if err := tx.QueryRowContext(ctx, `SELECT api_path, api_method FROM gbp_api_resources WHERE id = ? AND deleted_at IS NULL`, apiID).Scan(&apiPath, &apiMethod); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		key := apiMethod + ":" + apiPath
		if _, err := tx.ExecContext(ctx, `
INSERT INTO gbp_permission_policies
  (subject_type, subject_id, resource_type, resource_id, resource_key, action, effect, policy_status)
VALUES ('role', ?, 'api', ?, ?, ?, 'allow', 1)
ON DUPLICATE KEY UPDATE deleted_at = NULL, resource_id = VALUES(resource_id), policy_status = 1`,
			roleID, apiID, key, apiMethod); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SQLRepository) ReplaceDataScopes(ctx context.Context, roleID uint64, visibleRoleIDs []uint64) error {
	return r.replaceLinks(ctx,
		`DELETE FROM gbp_role_data_scopes WHERE role_id = ?`,
		"INSERT INTO gbp_role_data_scopes (role_id, visible_role_id) VALUES ",
		roleID, visibleRoleIDs,
	)
}

// ── Resource catalogs (used by role-permissions UI) ──────────────────────

func (r *SQLRepository) ListAllMenus(ctx context.Context) ([]rolemodel.MenuOption, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, parent_id, menu_title, menu_type, sort_no FROM gbp_menus
WHERE deleted_at IS NULL AND menu_status = 1
ORDER BY parent_id ASC, sort_no ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []rolemodel.MenuOption{}
	for rows.Next() {
		var m rolemodel.MenuOption
		if err := rows.Scan(&m.ID, &m.ParentID, &m.MenuTitle, &m.MenuType, &m.SortNo); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *SQLRepository) ListAllActions(ctx context.Context) ([]rolemodel.ActionOption, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT a.id, a.menu_id, m.menu_title, a.action_code, a.action_name, a.sort_no
FROM gbp_menu_actions a
JOIN gbp_menus m ON m.id = a.menu_id
WHERE a.deleted_at IS NULL AND a.action_status = 1 AND m.deleted_at IS NULL
ORDER BY a.menu_id ASC, a.sort_no ASC, a.id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []rolemodel.ActionOption{}
	for rows.Next() {
		var a rolemodel.ActionOption
		if err := rows.Scan(&a.ID, &a.MenuID, &a.MenuTitle, &a.ActionCode, &a.ActionName, &a.SortNo); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *SQLRepository) ListAllAPIs(ctx context.Context) ([]rolemodel.APIOption, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, api_path, api_method, COALESCE(api_group, ''), COALESCE(api_desc, ''), api_status
FROM gbp_api_resources WHERE deleted_at IS NULL
ORDER BY api_group ASC, api_path ASC, api_method ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []rolemodel.APIOption{}
	for rows.Next() {
		var a rolemodel.APIOption
		if err := rows.Scan(&a.ID, &a.APIPath, &a.APIMethod, &a.APIGroup, &a.APIDesc, &a.APIStatus); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// ── Helpers ──────────────────────────────────────────────────────────────

func buildListWhere(q rolemodel.ListQuery) (string, []interface{}) {
	conditions := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	if q.Keyword != "" {
		conditions = append(conditions, "(role_code LIKE ? OR role_name LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like)
	}
	if q.RoleStatus > 0 {
		conditions = append(conditions, "role_status = ?")
		args = append(args, q.RoleStatus)
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanRoleRow(row rowScanner) (rolemodel.Role, error) {
	var role rolemodel.Role
	var remark sql.NullString
	if err := row.Scan(
		&role.ID, &role.RoleCode, &role.RoleName, &role.ParentRoleID, &role.DefaultRoute,
		&role.SortNo, &role.RoleStatus, &remark, &role.CreatedAt, &role.UpdatedAt,
	); err != nil {
		return role, err
	}
	if remark.Valid {
		role.Remark = &remark.String
	}
	return role, nil
}

func (r *SQLRepository) userCounts(ctx context.Context, roleIDs []uint64) (map[uint64]int64, error) {
	out := map[uint64]int64{}
	if len(roleIDs) == 0 {
		return out, nil
	}
	placeholders := make([]string, len(roleIDs))
	args := make([]interface{}, len(roleIDs))
	for i, id := range roleIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`
SELECT ur.role_id, COUNT(1)
FROM gbp_user_roles ur
JOIN gbp_users u ON u.id = ur.user_id
WHERE ur.role_id IN (%s) AND u.deleted_at IS NULL
GROUP BY ur.role_id`, strings.Join(placeholders, ","))
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var roleID uint64
		var count int64
		if err := rows.Scan(&roleID, &count); err != nil {
			return nil, err
		}
		out[roleID] = count
	}
	return out, rows.Err()
}

func (r *SQLRepository) queryIDs(ctx context.Context, query string, args ...interface{}) ([]uint64, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []uint64{}
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (r *SQLRepository) replaceLinks(ctx context.Context, deleteSQL, insertSQLPrefix string, roleID uint64, ids []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, deleteSQL, roleID); err != nil {
		return err
	}
	if len(ids) == 0 {
		return tx.Commit()
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids)*2)
	for i, id := range ids {
		placeholders[i] = "(?, ?)"
		args = append(args, roleID, id)
	}
	query := insertSQLPrefix + strings.Join(placeholders, ",")
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	return tx.Commit()
}

func nullStringPtr(v *string) interface{} {
	if v == nil || *v == "" {
		return nil
	}
	return *v
}
