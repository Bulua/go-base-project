package auditrepo

import (
	"context"
	"database/sql"

	auditmodel "gobaseproject/server/internal/model/audit"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

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
