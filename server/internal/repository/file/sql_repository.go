package filerepo

import (
	"context"
	"database/sql"
	"strings"

	filemodel "gobaseproject/server/internal/model/file"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Insert(ctx context.Context, f filemodel.FileRecord) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_files (original_name, storage_key, file_size, mime_type, uploader_id)
VALUES (?, ?, ?, ?, ?)`,
		f.OriginalName, f.StorageKey, f.FileSize, f.MimeType, f.UploaderID)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) Count(ctx context.Context, q filemodel.FileListQuery) (int64, error) {
	where, args := buildWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_files "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) List(ctx context.Context, q filemodel.FileListQuery) ([]filemodel.FileRecord, error) {
	where, args := buildWhere(q)
	offset := (q.Page - 1) * q.PageSize
	args = append(args, q.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, `
SELECT id, original_name, storage_key, file_size, mime_type, uploader_id, created_at
FROM gbp_files `+where+`
ORDER BY created_at DESC
LIMIT ? OFFSET ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []filemodel.FileRecord
	for rows.Next() {
		var f filemodel.FileRecord
		if err := rows.Scan(&f.ID, &f.OriginalName, &f.StorageKey, &f.FileSize, &f.MimeType, &f.UploaderID, &f.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, f)
	}
	return items, rows.Err()
}

func (r *SQLRepository) GetByID(ctx context.Context, id uint64) (*filemodel.FileRecord, error) {
	var f filemodel.FileRecord
	err := r.db.QueryRowContext(ctx, `
SELECT id, original_name, storage_key, file_size, mime_type, uploader_id, created_at
FROM gbp_files WHERE id = ? AND deleted_at IS NULL`, id).Scan(
		&f.ID, &f.OriginalName, &f.StorageKey, &f.FileSize, &f.MimeType, &f.UploaderID, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &f, err
}

func (r *SQLRepository) Delete(ctx context.Context, id uint64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE gbp_files SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	return err
}

func buildWhere(q filemodel.FileListQuery) (string, []interface{}) {
	conds := []string{"deleted_at IS NULL"}
	var args []interface{}
	if q.Keyword != "" {
		conds = append(conds, "original_name LIKE ?")
		args = append(args, "%"+escapeLike(q.Keyword)+"%")
	}
	switch q.MimeCategory {
	case "image":
		conds = append(conds, "mime_type LIKE 'image/%'")
	case "video":
		conds = append(conds, "mime_type LIKE 'video/%'")
	case "audio":
		conds = append(conds, "mime_type LIKE 'audio/%'")
	case "text":
		conds = append(conds, "mime_type LIKE 'text/%'")
	case "pdf":
		conds = append(conds, "mime_type LIKE '%pdf%'")
	case "other":
		conds = append(conds, "mime_type NOT LIKE 'image/%' AND mime_type NOT LIKE 'video/%' AND mime_type NOT LIKE 'audio/%' AND mime_type NOT LIKE 'text/%' AND mime_type NOT LIKE '%pdf%'")
	}
	if q.StartDate != "" {
		conds = append(conds, "DATE(created_at) >= ?")
		args = append(args, q.StartDate)
	}
	if q.EndDate != "" {
		conds = append(conds, "DATE(created_at) <= ?")
		args = append(args, q.EndDate)
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
