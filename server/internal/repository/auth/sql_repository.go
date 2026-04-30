package authrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	authmodel "gobaseproject/server/internal/model/auth"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) FindUserByLogin(ctx context.Context, loginName string) (*authmodel.User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_uuid, login_name, password_hash, display_name, avatar_url, primary_role_id,
       phone_number, email_address, user_status, must_change_password
FROM gbp_users
WHERE login_name = ? AND deleted_at IS NULL
LIMIT 1`, loginName)

	var user authmodel.User
	var avatarURL, phoneNumber, emailAddress sql.NullString
	var primaryRoleID sql.NullInt64
	if err := row.Scan(
		&user.ID,
		&user.UserUUID,
		&user.LoginName,
		&user.PasswordHash,
		&user.DisplayName,
		&avatarURL,
		&primaryRoleID,
		&phoneNumber,
		&emailAddress,
		&user.UserStatus,
		&user.MustChangePassword,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authmodel.ErrUserNotFound
		}
		return nil, err
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	if primaryRoleID.Valid {
		value := uint64(primaryRoleID.Int64)
		user.PrimaryRoleID = &value
	}
	if phoneNumber.Valid {
		user.PhoneNumber = &phoneNumber.String
	}
	if emailAddress.Valid {
		user.EmailAddress = &emailAddress.String
	}
	return &user, nil
}

func (r *SQLRepository) ListRolesByUserID(ctx context.Context, userID uint64) ([]authmodel.Role, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT r.id, r.role_code, r.role_name, r.default_route
FROM gbp_roles r
JOIN gbp_user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = ? AND r.deleted_at IS NULL AND r.role_status = 1
ORDER BY r.sort_no ASC, r.id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []authmodel.Role
	for rows.Next() {
		var role authmodel.Role
		if err := rows.Scan(&role.ID, &role.RoleCode, &role.RoleName, &role.DefaultRoute); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *SQLRepository) InsertLoginAudit(ctx context.Context, audit authmodel.LoginAudit) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_login_audit_logs (user_id, login_name, source_ip, login_success, fail_reason, user_agent)
VALUES (?, ?, ?, ?, ?, ?)`,
		nullableUint64(audit.UserID),
		nullString(audit.LoginName),
		nullString(audit.SourceIP),
		audit.LoginSuccess,
		nullString(audit.FailReason),
		nullString(audit.UserAgent),
	)
	return err
}

func (r *SQLRepository) InsertTokenBlocklist(ctx context.Context, tokenHash string, expiresAt time.Time, reason string) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_jwt_blocklist (token_hash, expires_at, logout_reason)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE expires_at = VALUES(expires_at), logout_reason = VALUES(logout_reason)`,
		tokenHash,
		expiresAt,
		nullString(reason),
	)
	return err
}

func (r *SQLRepository) IsTokenBlocked(ctx context.Context, tokenHash string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
SELECT COUNT(1)
FROM gbp_jwt_blocklist
WHERE token_hash = ? AND expires_at > NOW(3) AND deleted_at IS NULL`, tokenHash).Scan(&count)
	return count > 0, err
}

func (r *SQLRepository) ListMenusByUserID(ctx context.Context, userID uint64) ([]authmodel.Menu, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT DISTINCT m.id, m.parent_id, m.menu_type, m.route_path, m.route_name, m.component_path,
       m.redirect_path, m.menu_title, m.menu_icon, m.sort_no, m.is_hidden, m.is_keep_alive,
       m.is_affix, m.active_route, m.transition_name, m.external_url
FROM gbp_menus m
JOIN gbp_role_menus rm ON rm.menu_id = m.id
JOIN gbp_user_roles ur ON ur.role_id = rm.role_id
WHERE ur.user_id = ? AND m.deleted_at IS NULL AND m.menu_status = 1
ORDER BY m.parent_id ASC, m.sort_no ASC, m.id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []authmodel.Menu
	for rows.Next() {
		var menu authmodel.Menu
		var routePath, routeName, componentPath, redirectPath, menuIcon sql.NullString
		var activeRoute, transitionName, externalURL sql.NullString
		if err := rows.Scan(
			&menu.ID,
			&menu.ParentID,
			&menu.MenuType,
			&routePath,
			&routeName,
			&componentPath,
			&redirectPath,
			&menu.MenuTitle,
			&menuIcon,
			&menu.SortNo,
			&menu.IsHidden,
			&menu.IsKeepAlive,
			&menu.IsAffix,
			&activeRoute,
			&transitionName,
			&externalURL,
		); err != nil {
			return nil, err
		}
		menu.RoutePath = stringPtr(routePath)
		menu.RouteName = stringPtr(routeName)
		menu.ComponentPath = stringPtr(componentPath)
		menu.RedirectPath = stringPtr(redirectPath)
		menu.MenuIcon = stringPtr(menuIcon)
		menu.ActiveRoute = stringPtr(activeRoute)
		menu.TransitionName = stringPtr(transitionName)
		menu.ExternalURL = stringPtr(externalURL)
		menus = append(menus, menu)
	}
	return menus, rows.Err()
}

func (r *SQLRepository) ListActionsByUserID(ctx context.Context, userID uint64) ([]authmodel.Action, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT DISTINCT ma.menu_id, ma.id, ma.action_code, ma.action_name
FROM gbp_menu_actions ma
JOIN gbp_role_actions ra ON ra.action_id = ma.id
JOIN gbp_user_roles ur ON ur.role_id = ra.role_id
WHERE ur.user_id = ? AND ma.deleted_at IS NULL AND ma.action_status = 1
ORDER BY ma.menu_id ASC, ma.id ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []authmodel.Action
	for rows.Next() {
		var action authmodel.Action
		if err := rows.Scan(&action.MenuID, &action.ActionID, &action.ActionCode, &action.ActionName); err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, rows.Err()
}

func nullString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}

func nullableUint64(value *uint64) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func stringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}
