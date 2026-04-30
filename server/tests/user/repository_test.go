package user_test

import (
	"context"
	"regexp"
	"testing"

	usermodel "gobaseproject/server/internal/model/user"
	userrepo "gobaseproject/server/internal/repository/user"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRepositorySoftDeleteArchivesLoginName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open sqlmock: %v", err)
	}
	defer db.Close()
	repo := userrepo.NewSQLRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`
UPDATE gbp_users
SET login_name = CONCAT(login_name, '#deleted#', id), deleted_at = NOW(3)
WHERE id = ? AND deleted_at IS NULL`)).
		WithArgs(uint64(12)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.SoftDelete(context.Background(), 12); err != nil {
		t.Fatalf("soft delete: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}

func TestRepositoryCreateArchivesDeletedLoginNameBeforeInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("open sqlmock: %v", err)
	}
	defer db.Close()
	repo := userrepo.NewSQLRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`
UPDATE gbp_users
SET login_name = CONCAT(login_name, '#deleted#', id)
WHERE login_name = ? AND deleted_at IS NOT NULL`)).
		WithArgs("test").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`
INSERT INTO gbp_users
  (user_uuid, login_name, password_hash, display_name, avatar_url, primary_role_id,
   phone_number, email_address, user_status, must_change_password, remark)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
		WithArgs(
			"uuid-1",
			"test",
			"hash",
			"Test User",
			nil,
			nil,
			nil,
			nil,
			usermodel.StatusActive,
			false,
			nil,
		).
		WillReturnResult(sqlmock.NewResult(88, 1))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM gbp_user_roles WHERE user_id = ?`)).
		WithArgs(uint64(88)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	id, err := repo.Create(context.Background(), "uuid-1", usermodel.CreateRequest{
		LoginName:   "test",
		DisplayName: "Test User",
		UserStatus:  usermodel.StatusActive,
	}, "hash")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if id != 88 {
		t.Fatalf("expected inserted id 88, got %d", id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expectations: %v", err)
	}
}
