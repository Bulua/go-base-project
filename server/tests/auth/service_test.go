package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	authmodel "gobaseproject/server/internal/model/auth"
	authservice "gobaseproject/server/internal/service/auth"
)

func TestServiceLoginReturnsTokensForAdminCredentials(t *testing.T) {
	ctx := context.Background()
	passwordHash, err := authservice.HashPassword("123456")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	repo := &fakeRepository{
		user: &authmodel.User{
			ID:           1,
			UserUUID:     "00000000-0000-0000-0000-000000000001",
			LoginName:    "admin",
			DisplayName:  "Super Admin",
			PasswordHash: passwordHash,
			UserStatus:   1,
		},
		roles: []authmodel.Role{{ID: 1, RoleCode: "super_admin", RoleName: "Super Admin"}},
	}
	service := authservice.NewService(repo, authservice.NewTokenManager(authservice.TokenConfig{
		Secret:          "unit-test-secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: 24 * time.Hour,
		TokenIssuer:     "gobaseproject-test",
		Clock:           func() time.Time { return time.Unix(1000, 0).UTC() },
	}))

	session, err := service.Login(ctx, authmodel.LoginRequest{
		LoginName: "admin",
		Password:  "123456",
	}, authmodel.RequestMeta{SourceIP: "127.0.0.1", UserAgent: "unit-test"})
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	if session.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if session.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if session.User.LoginName != "admin" {
		t.Fatalf("expected admin user, got %q", session.User.LoginName)
	}
	if len(session.User.Roles) != 1 || session.User.Roles[0].RoleCode != "super_admin" {
		t.Fatalf("expected super_admin role, got %#v", session.User.Roles)
	}
	if len(repo.loginAudits) != 1 || !repo.loginAudits[0].LoginSuccess {
		t.Fatalf("expected successful login audit, got %#v", repo.loginAudits)
	}
}

func TestServiceLoginRejectsWrongPasswordAndAuditsFailure(t *testing.T) {
	ctx := context.Background()
	passwordHash, err := authservice.HashPassword("123456")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	repo := &fakeRepository{
		user: &authmodel.User{
			ID:           1,
			LoginName:    "admin",
			DisplayName:  "Super Admin",
			PasswordHash: passwordHash,
			UserStatus:   1,
		},
	}
	service := authservice.NewService(repo, authservice.NewTokenManager(authservice.TokenConfig{
		Secret:          "unit-test-secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: 24 * time.Hour,
		TokenIssuer:     "gobaseproject-test",
	}))

	_, err = service.Login(ctx, authmodel.LoginRequest{
		LoginName: "admin",
		Password:  "bad-password",
	}, authmodel.RequestMeta{SourceIP: "127.0.0.1"})
	if !errors.Is(err, authmodel.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	if len(repo.loginAudits) != 1 {
		t.Fatalf("expected one login audit, got %d", len(repo.loginAudits))
	}
	audit := repo.loginAudits[0]
	if audit.LoginSuccess {
		t.Fatal("expected failed login audit")
	}
	if audit.FailReason == "" {
		t.Fatal("expected fail reason")
	}
}

func TestServiceLogoutBlocksPresentedToken(t *testing.T) {
	ctx := context.Background()
	repo := &fakeRepository{}
	clock := func() time.Time { return time.Unix(1000, 0).UTC() }
	manager := authservice.NewTokenManager(authservice.TokenConfig{
		Secret:          "unit-test-secret",
		AccessTokenTTL:  time.Hour,
		RefreshTokenTTL: 24 * time.Hour,
		TokenIssuer:     "gobaseproject-test",
		Clock:           clock,
	})
	service := authservice.NewService(repo, manager)

	token, expiresAt, err := manager.CreateToken(1, "admin", authmodel.TokenTypeAccess)
	if err != nil {
		t.Fatalf("create token: %v", err)
	}
	if err := service.Logout(ctx, token, "manual"); err != nil {
		t.Fatalf("logout: %v", err)
	}

	tokenHash := authservice.HashToken(token)
	blocked, ok := repo.blockedTokens[tokenHash]
	if !ok {
		t.Fatal("expected token hash to be inserted into blocklist")
	}
	if !blocked.ExpiresAt.Equal(expiresAt) {
		t.Fatalf("expected expires_at %s, got %s", expiresAt, blocked.ExpiresAt)
	}
	if blocked.Reason != "manual" {
		t.Fatalf("expected reason manual, got %q", blocked.Reason)
	}
}

type fakeRepository struct {
	user          *authmodel.User
	roles         []authmodel.Role
	loginAudits   []authmodel.LoginAudit
	blockedTokens map[string]authmodel.BlockedToken
}

func (r *fakeRepository) FindUserByLogin(ctx context.Context, loginName string) (*authmodel.User, error) {
	if r.user == nil || r.user.LoginName != loginName {
		return nil, authmodel.ErrUserNotFound
	}
	user := *r.user
	return &user, nil
}

func (r *fakeRepository) ListRolesByUserID(ctx context.Context, userID uint64) ([]authmodel.Role, error) {
	return append([]authmodel.Role(nil), r.roles...), nil
}

func (r *fakeRepository) InsertLoginAudit(ctx context.Context, audit authmodel.LoginAudit) error {
	r.loginAudits = append(r.loginAudits, audit)
	return nil
}

func (r *fakeRepository) InsertTokenBlocklist(ctx context.Context, tokenHash string, expiresAt time.Time, reason string) error {
	if r.blockedTokens == nil {
		r.blockedTokens = map[string]authmodel.BlockedToken{}
	}
	r.blockedTokens[tokenHash] = authmodel.BlockedToken{TokenHash: tokenHash, ExpiresAt: expiresAt, Reason: reason}
	return nil
}

func (r *fakeRepository) IsTokenBlocked(ctx context.Context, tokenHash string) (bool, error) {
	_, ok := r.blockedTokens[tokenHash]
	return ok, nil
}

func (r *fakeRepository) ListMenusByUserID(ctx context.Context, userID uint64) ([]authmodel.Menu, error) {
	return nil, nil
}

func (r *fakeRepository) ListActionsByUserID(ctx context.Context, userID uint64) ([]authmodel.Action, error) {
	return nil, nil
}
