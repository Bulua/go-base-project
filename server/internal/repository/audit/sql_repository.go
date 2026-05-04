package auditrepo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	auditmodel "gobaseproject/server/internal/model/audit"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// ── Write ─────────────────────────────────────────────────────────────────

func (r *SQLRepository) InsertOperationLog(ctx context.Context, log auditmodel.OperationLog) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_operation_audit_logs
  (user_id, source_ip, request_method, request_path, status_code, cost_ms, user_agent, error_message, request_body, response_body)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		nullableUint64(log.UserID),
		nullString(log.SourceIP),
		nullString(log.RequestMethod),
		nullString(log.RequestPath),
		nullInt(log.StatusCode),
		nullInt64(log.CostMs),
		nullString(log.UserAgent),
		nullString(log.ErrorMessage),
		nullString(log.RequestBody),
		nullString(log.ResponseBody),
	)
	return err
}

// ── Login logs ────────────────────────────────────────────────────────────

func (r *SQLRepository) CountLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) (int64, error) {
	where, args := buildLoginWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM gbp_login_audit_logs l "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) ListLoginLogs(ctx context.Context, q auditmodel.LoginLogQuery) ([]auditmodel.LoginLogRecord, error) {
	where, args := buildLoginWhere(q)
	offset := (q.Page - 1) * q.PageSize
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT l.id, l.user_id, l.login_name, l.source_ip, l.login_success, l.fail_reason, l.user_agent, l.created_at
FROM gbp_login_audit_logs l `+where+`
ORDER BY l.created_at DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []auditmodel.LoginLogRecord
	for rows.Next() {
		var rec auditmodel.LoginLogRecord
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.LoginName, &rec.SourceIP,
			&rec.LoginSuccess, &rec.FailReason, &rec.UserAgent, &rec.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, rec)
	}
	return items, rows.Err()
}

func (r *SQLRepository) CleanupLoginLogs(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM gbp_login_audit_logs WHERE created_at < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func buildLoginWhere(q auditmodel.LoginLogQuery) (string, []interface{}) {
	var conds []string
	var args []interface{}
	conds = append(conds, "l.deleted_at IS NULL")
	if q.Keyword != "" {
		conds = append(conds, "(l.login_name LIKE ? OR l.source_ip LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like)
	}
	if q.LoginSuccess == 1 {
		conds = append(conds, "l.login_success = 1")
	} else if q.LoginSuccess == 2 {
		conds = append(conds, "l.login_success = 0")
	}
	if q.StartDate != "" {
		conds = append(conds, "l.created_at >= ?")
		args = append(args, q.StartDate+" 00:00:00")
	}
	if q.EndDate != "" {
		conds = append(conds, "l.created_at <= ?")
		args = append(args, q.EndDate+" 23:59:59")
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

// ── Operation logs ────────────────────────────────────────────────────────

func (r *SQLRepository) CountOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) (int64, error) {
	where, args := buildOperationWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM gbp_operation_audit_logs l "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) ListOperationLogs(ctx context.Context, q auditmodel.OperationLogQuery) ([]auditmodel.OperationLogRecord, error) {
	where, args := buildOperationWhere(q)
	offset := (q.Page - 1) * q.PageSize
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT l.id, l.user_id, l.source_ip, l.request_method, l.request_path,
       l.status_code, l.cost_ms, l.user_agent, l.error_message, l.request_body, l.response_body, l.created_at
FROM gbp_operation_audit_logs l `+where+`
ORDER BY l.created_at DESC LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []auditmodel.OperationLogRecord
	for rows.Next() {
		var rec auditmodel.OperationLogRecord
		err := rows.Scan(&rec.ID, &rec.UserID, &rec.SourceIP, &rec.RequestMethod, &rec.RequestPath,
			&rec.StatusCode, &rec.CostMs, &rec.UserAgent, &rec.ErrorMessage,
			&rec.RequestBody, &rec.ResponseBody, &rec.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, rec)
	}
	return items, rows.Err()
}

func (r *SQLRepository) GetOperationLogByID(ctx context.Context, id uint64) (*auditmodel.OperationLogRecord, error) {
	var rec auditmodel.OperationLogRecord
	err := r.db.QueryRowContext(ctx, `
SELECT id, user_id, source_ip, request_method, request_path,
       status_code, cost_ms, user_agent, error_message, request_body, response_body, created_at
FROM gbp_operation_audit_logs
WHERE id = ? AND deleted_at IS NULL`, id).Scan(
		&rec.ID, &rec.UserID, &rec.SourceIP, &rec.RequestMethod, &rec.RequestPath,
		&rec.StatusCode, &rec.CostMs, &rec.UserAgent, &rec.ErrorMessage,
		&rec.RequestBody, &rec.ResponseBody, &rec.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rec, err
}

func (r *SQLRepository) CleanupOperationLogs(ctx context.Context, days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM gbp_operation_audit_logs WHERE created_at < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func buildOperationWhere(q auditmodel.OperationLogQuery) (string, []interface{}) {
	var conds []string
	var args []interface{}
	conds = append(conds, "l.deleted_at IS NULL")
	if q.Keyword != "" {
		conds = append(conds, "(l.request_path LIKE ? OR l.source_ip LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like)
	}
	if q.Method != "" {
		conds = append(conds, "l.request_method = ?")
		args = append(args, strings.ToUpper(q.Method))
	}
	if q.StatusCode > 0 {
		conds = append(conds, "l.status_code = ?")
		args = append(args, q.StatusCode)
	}
	if q.StartDate != "" {
		conds = append(conds, "l.created_at >= ?")
		args = append(args, q.StartDate+" 00:00:00")
	}
	if q.EndDate != "" {
		conds = append(conds, "l.created_at <= ?")
		args = append(args, q.EndDate+" 23:59:59")
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

// ── helpers ───────────────────────────────────────────────────────────────

func nullString(v string) interface{} {
	if v == "" {
		return nil
	}
	return v
}

func nullInt(v int) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func nullInt64(v int64) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

func nullableUint64(v *uint64) interface{} {
	if v == nil {
		return nil
	}
	return *v
}

