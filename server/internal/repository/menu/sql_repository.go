package menurepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	menumodel "gobaseproject/server/internal/model/menu"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Count(ctx context.Context, q menumodel.ListQuery) (int64, error) {
	where, args := buildWhere(q)
	var total int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_menus "+where, args...).Scan(&total)
	return total, err
}

func (r *SQLRepository) List(ctx context.Context, q menumodel.ListQuery) ([]menumodel.Menu, error) {
	where, args := buildWhere(q)
	query := `SELECT id, parent_id, menu_type, route_path, route_name, component_path, redirect_path,
	       menu_title, menu_icon, sort_no, is_hidden, is_keep_alive, is_affix,
	       active_route, transition_name, external_url, menu_status, created_at, updated_at
FROM gbp_menus ` + where + ` ORDER BY parent_id ASC, sort_no ASC, id ASC LIMIT ? OFFSET ?`
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMenuRows(rows)
}

func (r *SQLRepository) ListAll(ctx context.Context) ([]menumodel.Menu, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, parent_id, menu_type, route_path, route_name, component_path, redirect_path,
       menu_title, menu_icon, sort_no, is_hidden, is_keep_alive, is_affix,
       active_route, transition_name, external_url, menu_status, created_at, updated_at
FROM gbp_menus
WHERE deleted_at IS NULL
ORDER BY parent_id ASC, sort_no ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMenuRows(rows)
}

func (r *SQLRepository) GetByID(ctx context.Context, id uint64) (*menumodel.Menu, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, parent_id, menu_type, route_path, route_name, component_path, redirect_path,
       menu_title, menu_icon, sort_no, is_hidden, is_keep_alive, is_affix,
       active_route, transition_name, external_url, menu_status, created_at, updated_at
FROM gbp_menus
WHERE id = ? AND deleted_at IS NULL`, id)
	menu, err := scanMenuRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, menumodel.ErrMenuNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *SQLRepository) CountChildren(ctx context.Context, parentID uint64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM gbp_menus WHERE parent_id = ? AND deleted_at IS NULL", parentID).Scan(&count)
	return count, err
}

func (r *SQLRepository) Create(ctx context.Context, req menumodel.SaveRequest) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_menus
  (parent_id, menu_type, route_path, route_name, component_path, redirect_path,
   menu_title, menu_icon, sort_no, is_hidden, is_keep_alive, is_affix,
   active_route, transition_name, external_url, menu_status)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		req.ParentID,
		req.MenuType,
		nullStr(req.RoutePath),
		nullStr(req.RouteName),
		nullStr(req.ComponentPath),
		nullStr(req.RedirectPath),
		req.MenuTitle,
		nullStr(req.MenuIcon),
		req.SortNo,
		req.IsHidden,
		req.IsKeepAlive,
		req.IsAffix,
		nullStr(req.ActiveRoute),
		nullStr(req.TransitionName),
		nullStr(req.ExternalURL),
		req.MenuStatus,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) Update(ctx context.Context, id uint64, req menumodel.SaveRequest) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_menus SET
  parent_id = ?, menu_type = ?, route_path = ?, route_name = ?, component_path = ?,
  redirect_path = ?, menu_title = ?, menu_icon = ?, sort_no = ?,
  is_hidden = ?, is_keep_alive = ?, is_affix = ?,
  active_route = ?, transition_name = ?, external_url = ?, menu_status = ?
WHERE id = ? AND deleted_at IS NULL`,
		req.ParentID,
		req.MenuType,
		nullStr(req.RoutePath),
		nullStr(req.RouteName),
		nullStr(req.ComponentPath),
		nullStr(req.RedirectPath),
		req.MenuTitle,
		nullStr(req.MenuIcon),
		req.SortNo,
		req.IsHidden,
		req.IsKeepAlive,
		req.IsAffix,
		nullStr(req.ActiveRoute),
		nullStr(req.TransitionName),
		nullStr(req.ExternalURL),
		req.MenuStatus,
		id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return menumodel.ErrMenuNotFound
	}
	return nil
}

func (r *SQLRepository) SoftDelete(ctx context.Context, id uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"UPDATE gbp_menus SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return menumodel.ErrMenuNotFound
	}

	// cascade: soft-delete menu actions
	if _, err := tx.ExecContext(ctx,
		"UPDATE gbp_menu_actions SET deleted_at = NOW(3) WHERE menu_id = ? AND deleted_at IS NULL", id); err != nil {
		return err
	}

	return tx.Commit()
}

// ── Menu Actions ─────────────────────────────────────────────────────────────

func (r *SQLRepository) ListActions(ctx context.Context, menuID uint64) ([]menumodel.MenuAction, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, menu_id, action_code, action_name, action_desc, sort_no, action_status, created_at, updated_at
FROM gbp_menu_actions
WHERE menu_id = ? AND deleted_at IS NULL
ORDER BY sort_no ASC, id ASC`, menuID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanActionRows(rows)
}

func (r *SQLRepository) GetActionByID(ctx context.Context, id uint64) (*menumodel.MenuAction, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, menu_id, action_code, action_name, action_desc, sort_no, action_status, created_at, updated_at
FROM gbp_menu_actions WHERE id = ? AND deleted_at IS NULL`, id)
	a, err := scanActionRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, menumodel.ErrActionNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *SQLRepository) CountActionByCode(ctx context.Context, menuID uint64, code string, excludeID uint64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(1) FROM gbp_menu_actions
WHERE menu_id = ? AND action_code = ? AND deleted_at IS NULL AND id != ?`,
		menuID, code, excludeID).Scan(&count)
	return count, err
}

func (r *SQLRepository) CreateAction(ctx context.Context, menuID uint64, req menumodel.SaveActionRequest) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_menu_actions (menu_id, action_code, action_name, action_desc, sort_no, action_status)
VALUES (?, ?, ?, ?, ?, ?)`,
		menuID, req.ActionCode, req.ActionName, nullStr(req.ActionDesc), req.SortNo, req.ActionStatus)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) UpdateAction(ctx context.Context, id uint64, req menumodel.SaveActionRequest) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_menu_actions
SET action_code = ?, action_name = ?, action_desc = ?, sort_no = ?, action_status = ?
WHERE id = ? AND deleted_at IS NULL`,
		req.ActionCode, req.ActionName, nullStr(req.ActionDesc), req.SortNo, req.ActionStatus, id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return menumodel.ErrActionNotFound
	}
	return nil
}

func (r *SQLRepository) DeleteAction(ctx context.Context, id uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		"UPDATE gbp_menu_actions SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return menumodel.ErrActionNotFound
	}
	// cascade: remove role assignments for this action
	if _, err := tx.ExecContext(ctx,
		"DELETE FROM gbp_role_actions WHERE action_id = ?", id); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *SQLRepository) AssignToRole(ctx context.Context, roleID, menuID uint64) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT IGNORE INTO gbp_role_menus (role_id, menu_id) VALUES (?, ?)", roleID, menuID)
	return err
}

// ── helpers ──────────────────────────────────────────────────────────────────

func buildWhere(q menumodel.ListQuery) (string, []interface{}) {
	conds := []string{"deleted_at IS NULL"}
	args := []interface{}{}
	if q.Keyword != "" {
		conds = append(conds, "(menu_title LIKE ? OR route_path LIKE ? OR component_path LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like, like)
	}
	if q.MenuStatus > 0 {
		conds = append(conds, "menu_status = ?")
		args = append(args, q.MenuStatus)
	}
	if q.MenuType > 0 {
		conds = append(conds, "menu_type = ?")
		args = append(args, q.MenuType)
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanMenuRow(row rowScanner) (menumodel.Menu, error) {
	var m menumodel.Menu
	var routePath, routeName, componentPath, redirectPath, menuIcon sql.NullString
	var activeRoute, transitionName, externalURL sql.NullString
	err := row.Scan(
		&m.ID, &m.ParentID, &m.MenuType,
		&routePath, &routeName, &componentPath, &redirectPath,
		&m.MenuTitle, &menuIcon,
		&m.SortNo, &m.IsHidden, &m.IsKeepAlive, &m.IsAffix,
		&activeRoute, &transitionName, &externalURL,
		&m.MenuStatus, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return m, err
	}
	m.RoutePath = strPtr(routePath)
	m.RouteName = strPtr(routeName)
	m.ComponentPath = strPtr(componentPath)
	m.RedirectPath = strPtr(redirectPath)
	m.MenuIcon = strPtr(menuIcon)
	m.ActiveRoute = strPtr(activeRoute)
	m.TransitionName = strPtr(transitionName)
	m.ExternalURL = strPtr(externalURL)
	return m, nil
}

func scanMenuRows(rows *sql.Rows) ([]menumodel.Menu, error) {
	menus := []menumodel.Menu{}
	for rows.Next() {
		m, err := scanMenuRow(rows)
		if err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	return menus, rows.Err()
}

func strPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func scanActionRow(row rowScanner) (menumodel.MenuAction, error) {
	var a menumodel.MenuAction
	var desc sql.NullString
	err := row.Scan(&a.ID, &a.MenuID, &a.ActionCode, &a.ActionName, &desc,
		&a.SortNo, &a.ActionStatus, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return a, err
	}
	if desc.Valid {
		a.ActionDesc = &desc.String
	}
	return a, nil
}

func scanActionRows(rows *sql.Rows) ([]menumodel.MenuAction, error) {
	actions := []menumodel.MenuAction{}
	for rows.Next() {
		a, err := scanActionRow(rows)
		if err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, rows.Err()
}

func nullStr(v *string) interface{} {
	if v == nil || *v == "" {
		return nil
	}
	return *v
}
