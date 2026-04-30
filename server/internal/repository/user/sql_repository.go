package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	usermodel "gobaseproject/server/internal/model/user"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Count(ctx context.Context, q usermodel.ListQuery) (int64, error) {
	where, args := buildListWhere(q)
	query := "SELECT COUNT(DISTINCT u.id) FROM gbp_users u "
	if q.RoleID > 0 {
		query += "JOIN gbp_user_roles ur ON ur.user_id = u.id "
	}
	query += where
	var total int64
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *SQLRepository) List(ctx context.Context, q usermodel.ListQuery) ([]usermodel.User, error) {
	where, args := buildListWhere(q)
	query := `SELECT DISTINCT u.id, u.user_uuid, u.login_name, u.display_name, u.avatar_url,
       u.primary_role_id, u.phone_number, u.email_address, u.user_status,
       u.must_change_password, u.last_login_at, u.remark, u.created_at, u.updated_at
FROM gbp_users u `
	if q.RoleID > 0 {
		query += "JOIN gbp_user_roles ur ON ur.user_id = u.id "
	}
	query += where + " ORDER BY u.id DESC LIMIT ? OFFSET ?"
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []usermodel.User{}
	ids := []uint64{}
	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
		ids = append(ids, u.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return users, nil
	}
	rolesByUser, err := r.listRolesForUsers(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range users {
		users[i].Roles = rolesByUser[users[i].ID]
	}
	return users, nil
}

func (r *SQLRepository) GetByID(ctx context.Context, id uint64) (*usermodel.User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_uuid, login_name, display_name, avatar_url, primary_role_id,
       phone_number, email_address, user_status, must_change_password,
       last_login_at, remark, created_at, updated_at
FROM gbp_users
WHERE id = ? AND deleted_at IS NULL`, id)
	u, err := scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usermodel.ErrUserNotFound
		}
		return nil, err
	}
	roles, err := r.listRolesForUsers(ctx, []uint64{id})
	if err != nil {
		return nil, err
	}
	u.Roles = roles[id]
	return &u, nil
}

func (r *SQLRepository) GetByLoginName(ctx context.Context, loginName string) (*usermodel.User, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, user_uuid, login_name, display_name, avatar_url, primary_role_id,
       phone_number, email_address, user_status, must_change_password,
       last_login_at, remark, created_at, updated_at
FROM gbp_users
WHERE login_name = ? AND deleted_at IS NULL`, loginName)
	u, err := scanUserRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, usermodel.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *SQLRepository) Create(ctx context.Context, userUUID string, req usermodel.CreateRequest, passwordHash string) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if err := archiveDeletedLoginNameTx(ctx, tx, req.LoginName); err != nil {
		return 0, err
	}

	res, err := tx.ExecContext(ctx, `
INSERT INTO gbp_users
  (user_uuid, login_name, password_hash, display_name, avatar_url, primary_role_id,
   phone_number, email_address, user_status, must_change_password, remark)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userUUID,
		req.LoginName,
		passwordHash,
		req.DisplayName,
		nullStringPtr(req.AvatarURL),
		nullableUint64(req.PrimaryRoleID),
		nullStringPtr(req.PhoneNumber),
		nullStringPtr(req.EmailAddress),
		req.UserStatus,
		false,
		nullStringPtr(req.Remark),
	)
	if err != nil {
		return 0, err
	}
	insertedID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	id := uint64(insertedID)

	if err := replaceUserRolesTx(ctx, tx, id, req.RoleIDs); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLRepository) Update(ctx context.Context, id uint64, req usermodel.UpdateRequest) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_users
SET display_name = ?, avatar_url = ?, primary_role_id = ?,
    phone_number = ?, email_address = ?, remark = ?
WHERE id = ? AND deleted_at IS NULL`,
		req.DisplayName,
		nullStringPtr(req.AvatarURL),
		nullableUint64(req.PrimaryRoleID),
		nullStringPtr(req.PhoneNumber),
		nullStringPtr(req.EmailAddress),
		nullStringPtr(req.Remark),
		id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usermodel.ErrUserNotFound
	}
	return nil
}

func (r *SQLRepository) SoftDelete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_users
SET login_name = CONCAT(login_name, '#deleted#', id), deleted_at = NOW(3)
WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usermodel.ErrUserNotFound
	}
	return nil
}

func (r *SQLRepository) UpdateStatus(ctx context.Context, id uint64, status int) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_users SET user_status = ? WHERE id = ? AND deleted_at IS NULL`, status, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usermodel.ErrUserNotFound
	}
	return nil
}

func (r *SQLRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string, mustChange bool) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_users SET password_hash = ?, must_change_password = ?
WHERE id = ? AND deleted_at IS NULL`, passwordHash, mustChange, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return usermodel.ErrUserNotFound
	}
	return nil
}

func (r *SQLRepository) ReplaceRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := replaceUserRolesTx(ctx, tx, userID, roleIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *SQLRepository) ListActiveRoles(ctx context.Context) ([]usermodel.Role, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, role_code, role_name FROM gbp_roles
WHERE deleted_at IS NULL AND role_status = 1
ORDER BY sort_no ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []usermodel.Role{}
	for rows.Next() {
		var role usermodel.Role
		if err := rows.Scan(&role.ID, &role.RoleCode, &role.RoleName); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

// ── helpers ────────────────────────────────────────────────────────────────

func buildListWhere(q usermodel.ListQuery) (string, []interface{}) {
	conditions := []string{"u.deleted_at IS NULL"}
	args := []interface{}{}
	if q.Keyword != "" {
		conditions = append(conditions, "(u.login_name LIKE ? OR u.display_name LIKE ? OR u.email_address LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like, like)
	}
	if q.UserStatus > 0 {
		conditions = append(conditions, "u.user_status = ?")
		args = append(args, q.UserStatus)
	}
	if q.RoleID > 0 {
		conditions = append(conditions, "ur.role_id = ?")
		args = append(args, q.RoleID)
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanUserRow(row rowScanner) (usermodel.User, error) {
	var u usermodel.User
	var avatarURL, phoneNumber, emailAddress, remark sql.NullString
	var primaryRoleID sql.NullInt64
	var lastLoginAt sql.NullTime
	if err := row.Scan(
		&u.ID, &u.UserUUID, &u.LoginName, &u.DisplayName, &avatarURL,
		&primaryRoleID, &phoneNumber, &emailAddress, &u.UserStatus,
		&u.MustChangePassword, &lastLoginAt, &remark, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return u, err
	}
	if avatarURL.Valid {
		u.AvatarURL = &avatarURL.String
	}
	if primaryRoleID.Valid {
		v := uint64(primaryRoleID.Int64)
		u.PrimaryRoleID = &v
	}
	if phoneNumber.Valid {
		u.PhoneNumber = &phoneNumber.String
	}
	if emailAddress.Valid {
		u.EmailAddress = &emailAddress.String
	}
	if remark.Valid {
		u.Remark = &remark.String
	}
	if lastLoginAt.Valid {
		t := lastLoginAt.Time
		u.LastLoginAt = &t
	}
	return u, nil
}

func (r *SQLRepository) listRolesForUsers(ctx context.Context, userIDs []uint64) (map[uint64][]usermodel.Role, error) {
	if len(userIDs) == 0 {
		return map[uint64][]usermodel.Role{}, nil
	}
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`
SELECT ur.user_id, r.id, r.role_code, r.role_name
FROM gbp_user_roles ur
JOIN gbp_roles r ON r.id = ur.role_id
WHERE ur.user_id IN (%s) AND r.deleted_at IS NULL
ORDER BY r.sort_no ASC, r.id ASC`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[uint64][]usermodel.Role{}
	for rows.Next() {
		var userID uint64
		var role usermodel.Role
		if err := rows.Scan(&userID, &role.ID, &role.RoleCode, &role.RoleName); err != nil {
			return nil, err
		}
		out[userID] = append(out[userID], role)
	}
	return out, rows.Err()
}

func replaceUserRolesTx(ctx context.Context, tx *sql.Tx, userID uint64, roleIDs []uint64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM gbp_user_roles WHERE user_id = ?`, userID); err != nil {
		return err
	}
	if len(roleIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(roleIDs))
	args := make([]interface{}, 0, len(roleIDs)*2)
	for i, rid := range roleIDs {
		placeholders[i] = "(?, ?)"
		args = append(args, userID, rid)
	}
	query := "INSERT INTO gbp_user_roles (user_id, role_id) VALUES " + strings.Join(placeholders, ",")
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

func archiveDeletedLoginNameTx(ctx context.Context, tx *sql.Tx, loginName string) error {
	_, err := tx.ExecContext(ctx, `
UPDATE gbp_users
SET login_name = CONCAT(login_name, '#deleted#', id)
WHERE login_name = ? AND deleted_at IS NOT NULL`, loginName)
	return err
}

func nullStringPtr(v *string) interface{} {
	if v == nil || *v == "" {
		return nil
	}
	return *v
}

func nullableUint64(v *uint64) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
