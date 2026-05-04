package dictrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	dictmodel "gobaseproject/server/internal/model/dict"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

// ── Dictionary ────────────────────────────────────────────────────────────

func (r *SQLRepository) Count(ctx context.Context, q dictmodel.ListQuery) (int64, error) {
	where, args := buildWhere(q)
	var n int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM gbp_dictionaries d "+where, args...).Scan(&n)
	return n, err
}

func (r *SQLRepository) List(ctx context.Context, q dictmodel.ListQuery) ([]dictmodel.Dictionary, error) {
	where, args := buildWhere(q)
	query := `
SELECT d.id, d.dict_name, d.dict_code, d.dict_status, d.parent_id, d.remark,
       COUNT(i.id) AS item_count, d.created_at, d.updated_at
FROM gbp_dictionaries d
LEFT JOIN gbp_dictionary_items i ON i.dict_id = d.id AND i.deleted_at IS NULL
` + where + `
GROUP BY d.id
ORDER BY d.id ASC
LIMIT ? OFFSET ?`
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []dictmodel.Dictionary{}
	for rows.Next() {
		d, err := scanDict(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *SQLRepository) GetByID(ctx context.Context, id uint64) (*dictmodel.Dictionary, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT d.id, d.dict_name, d.dict_code, d.dict_status, d.parent_id, d.remark,
       COUNT(i.id) AS item_count, d.created_at, d.updated_at
FROM gbp_dictionaries d
LEFT JOIN gbp_dictionary_items i ON i.dict_id = d.id AND i.deleted_at IS NULL
WHERE d.id = ? AND d.deleted_at IS NULL
GROUP BY d.id`, id)
	d, err := scanDict(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dictmodel.ErrDictNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *SQLRepository) GetByCode(ctx context.Context, code string) (*dictmodel.Dictionary, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT d.id, d.dict_name, d.dict_code, d.dict_status, d.parent_id, d.remark,
       COUNT(i.id) AS item_count, d.created_at, d.updated_at
FROM gbp_dictionaries d
LEFT JOIN gbp_dictionary_items i ON i.dict_id = d.id AND i.deleted_at IS NULL
WHERE d.dict_code = ? AND d.deleted_at IS NULL
GROUP BY d.id`, code)
	d, err := scanDict(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dictmodel.ErrDictNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *SQLRepository) CountByCode(ctx context.Context, code string, excludeID uint64) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM gbp_dictionaries WHERE dict_code = ? AND id != ? AND deleted_at IS NULL`,
		code, excludeID,
	).Scan(&n)
	return n, err
}

func (r *SQLRepository) Create(ctx context.Context, p dictmodel.SaveDictPayload) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_dictionaries (dict_name, dict_code, dict_status, parent_id, remark)
VALUES (?, ?, ?, ?, ?)`,
		p.DictName, p.DictCode, p.DictStatus, p.ParentID, nullStr(p.Remark),
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) Update(ctx context.Context, id uint64, p dictmodel.SaveDictPayload) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_dictionaries
SET dict_name = ?, dict_code = ?, dict_status = ?, parent_id = ?, remark = ?
WHERE id = ? AND deleted_at IS NULL`,
		p.DictName, p.DictCode, p.DictStatus, p.ParentID, nullStr(p.Remark), id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return dictmodel.ErrDictNotFound
	}
	return nil
}

func (r *SQLRepository) SoftDelete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE gbp_dictionaries SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return dictmodel.ErrDictNotFound
	}
	return nil
}

// ── DictItem ──────────────────────────────────────────────────────────────

func (r *SQLRepository) ListItems(ctx context.Context, dictID uint64) ([]dictmodel.DictItem, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT id, dict_id, item_label, item_value, item_extra, item_status,
       sort_no, parent_id, tree_level, tree_path, created_at, updated_at
FROM gbp_dictionary_items
WHERE dict_id = ? AND deleted_at IS NULL
ORDER BY sort_no ASC, id ASC`, dictID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []dictmodel.DictItem{}
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *SQLRepository) GetItemByID(ctx context.Context, id uint64) (*dictmodel.DictItem, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT id, dict_id, item_label, item_value, item_extra, item_status,
       sort_no, parent_id, tree_level, tree_path, created_at, updated_at
FROM gbp_dictionary_items WHERE id = ? AND deleted_at IS NULL`, id)
	item, err := scanItem(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, dictmodel.ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *SQLRepository) CountItemByValue(ctx context.Context, dictID uint64, value string, excludeID uint64) (int64, error) {
	var n int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM gbp_dictionary_items WHERE dict_id = ? AND item_value = ? AND id != ? AND deleted_at IS NULL`,
		dictID, value, excludeID,
	).Scan(&n)
	return n, err
}

func (r *SQLRepository) CreateItem(ctx context.Context, dictID uint64, p dictmodel.SaveItemPayload) (uint64, error) {
	res, err := r.db.ExecContext(ctx, `
INSERT INTO gbp_dictionary_items
  (dict_id, item_label, item_value, item_extra, item_status, sort_no, parent_id)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		dictID, p.ItemLabel, p.ItemValue, nullStr(p.ItemExtra),
		p.ItemStatus, p.SortNo, p.ParentID,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return uint64(id), err
}

func (r *SQLRepository) UpdateItem(ctx context.Context, id uint64, p dictmodel.SaveItemPayload) error {
	res, err := r.db.ExecContext(ctx, `
UPDATE gbp_dictionary_items
SET item_label = ?, item_value = ?, item_extra = ?, item_status = ?, sort_no = ?, parent_id = ?
WHERE id = ? AND deleted_at IS NULL`,
		p.ItemLabel, p.ItemValue, nullStr(p.ItemExtra),
		p.ItemStatus, p.SortNo, p.ParentID, id,
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return dictmodel.ErrItemNotFound
	}
	return nil
}

func (r *SQLRepository) SoftDeleteItem(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE gbp_dictionary_items SET deleted_at = NOW(3) WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return dictmodel.ErrItemNotFound
	}
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanDict(row rowScanner) (dictmodel.Dictionary, error) {
	var d dictmodel.Dictionary
	var remark sql.NullString
	err := row.Scan(
		&d.ID, &d.DictName, &d.DictCode, &d.DictStatus, &d.ParentID,
		&remark, &d.ItemCount, &d.CreatedAt, &d.UpdatedAt,
	)
	if remark.Valid {
		d.Remark = &remark.String
	}
	return d, err
}

func scanItem(row rowScanner) (dictmodel.DictItem, error) {
	var item dictmodel.DictItem
	var extra, treePath sql.NullString
	err := row.Scan(
		&item.ID, &item.DictID, &item.ItemLabel, &item.ItemValue,
		&extra, &item.ItemStatus, &item.SortNo, &item.ParentID,
		&item.TreeLevel, &treePath, &item.CreatedAt, &item.UpdatedAt,
	)
	if extra.Valid {
		item.ItemExtra = &extra.String
	}
	if treePath.Valid {
		item.TreePath = &treePath.String
	}
	return item, err
}

func buildWhere(q dictmodel.ListQuery) (string, []interface{}) {
	conds := []string{"d.deleted_at IS NULL"}
	args := []interface{}{}
	if q.Keyword != "" {
		conds = append(conds, "(d.dict_name LIKE ? OR d.dict_code LIKE ?)")
		like := "%" + q.Keyword + "%"
		args = append(args, like, like)
	}
	if q.DictStatus > 0 {
		conds = append(conds, "d.dict_status = ?")
		args = append(args, q.DictStatus)
	}
	return "WHERE " + strings.Join(conds, " AND "), args
}

func nullStr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}
