package authservice

import (
	"context"
	"strings"
	"time"

	authmodel "gobaseproject/server/internal/model/auth"
)

type Repository interface {
	FindUserByLogin(ctx context.Context, loginName string) (*authmodel.User, error)
	ListRolesByUserID(ctx context.Context, userID uint64) ([]authmodel.Role, error)
	InsertLoginAudit(ctx context.Context, audit authmodel.LoginAudit) error
	UpdateLastLoginAt(ctx context.Context, userID uint64, lastLoginAt time.Time) error
	InsertTokenBlocklist(ctx context.Context, tokenHash string, expiresAt time.Time, reason string) error
	IsTokenBlocked(ctx context.Context, tokenHash string) (bool, error)
	ListMenusByUserID(ctx context.Context, userID uint64) ([]authmodel.Menu, error)
	ListActionsByUserID(ctx context.Context, userID uint64) ([]authmodel.Action, error)
}

type Service struct {
	repo   Repository
	tokens TokenManager
}

func NewService(repo Repository, tokens TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, req authmodel.LoginRequest, meta authmodel.RequestMeta) (authmodel.Session, error) {
	loginName := strings.TrimSpace(req.LoginName)
	if loginName == "" {
		loginName = strings.TrimSpace(req.Username)
	}
	password := req.Password
	user, err := s.repo.FindUserByLogin(ctx, loginName)
	if err != nil {
		_ = s.repo.InsertLoginAudit(ctx, authmodel.LoginAudit{
			LoginName:    loginName,
			SourceIP:     meta.SourceIP,
			UserAgent:    meta.UserAgent,
			LoginSuccess: false,
			FailReason:   "invalid credentials",
		})
		return authmodel.Session{}, authmodel.ErrInvalidCredentials
	}

	if user.UserStatus != 1 {
		userID := user.ID
		_ = s.repo.InsertLoginAudit(ctx, authmodel.LoginAudit{
			UserID:       &userID,
			LoginName:    loginName,
			SourceIP:     meta.SourceIP,
			UserAgent:    meta.UserAgent,
			LoginSuccess: false,
			FailReason:   "user disabled",
		})
		return authmodel.Session{}, authmodel.ErrUserDisabled
	}

	if !CheckPassword(password, user.PasswordHash) {
		userID := user.ID
		_ = s.repo.InsertLoginAudit(ctx, authmodel.LoginAudit{
			UserID:       &userID,
			LoginName:    loginName,
			SourceIP:     meta.SourceIP,
			UserAgent:    meta.UserAgent,
			LoginSuccess: false,
			FailReason:   "invalid credentials",
		})
		return authmodel.Session{}, authmodel.ErrInvalidCredentials
	}

	now := s.tokens.clock().UTC()
	if err := s.repo.UpdateLastLoginAt(ctx, user.ID, now); err != nil {
		return authmodel.Session{}, err
	}

	roles, err := s.repo.ListRolesByUserID(ctx, user.ID)
	if err != nil {
		return authmodel.Session{}, err
	}
	user.Roles = roles
	accessToken, expiresAt, err := s.tokens.CreateToken(user.ID, user.LoginName, authmodel.TokenTypeAccess)
	if err != nil {
		return authmodel.Session{}, err
	}
	refreshToken, refreshExpiresAt, err := s.tokens.CreateToken(user.ID, user.LoginName, authmodel.TokenTypeRefresh)
	if err != nil {
		return authmodel.Session{}, err
	}
	userID := user.ID
	_ = s.repo.InsertLoginAudit(ctx, authmodel.LoginAudit{
		UserID:       &userID,
		LoginName:    loginName,
		SourceIP:     meta.SourceIP,
		UserAgent:    meta.UserAgent,
		LoginSuccess: true,
	})
	return authmodel.Session{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		User:             *user,
	}, nil
}

func (s *Service) Logout(ctx context.Context, tokenValue string, reason string) error {
	claims, err := s.tokens.Parse(tokenValue, "")
	if err != nil {
		return err
	}
	if claims.ExpiresAt == nil {
		return authmodel.ErrInvalidToken
	}
	return s.repo.InsertTokenBlocklist(ctx, HashToken(tokenValue), claims.ExpiresAt.Time, reason)
}

func (s *Service) ParseAccessToken(ctx context.Context, tokenValue string) (*Claims, error) {
	if blocked, err := s.repo.IsTokenBlocked(ctx, HashToken(tokenValue)); err != nil {
		return nil, err
	} else if blocked {
		return nil, authmodel.ErrTokenBlocked
	}
	return s.tokens.Parse(tokenValue, authmodel.TokenTypeAccess)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (authmodel.Session, error) {
	if blocked, err := s.repo.IsTokenBlocked(ctx, HashToken(refreshToken)); err != nil {
		return authmodel.Session{}, err
	} else if blocked {
		return authmodel.Session{}, authmodel.ErrTokenBlocked
	}
	claims, err := s.tokens.Parse(refreshToken, authmodel.TokenTypeRefresh)
	if err != nil {
		return authmodel.Session{}, err
	}
	user, err := s.repo.FindUserByLogin(ctx, claims.LoginName)
	if err != nil {
		return authmodel.Session{}, authmodel.ErrInvalidCredentials
	}
	roles, err := s.repo.ListRolesByUserID(ctx, user.ID)
	if err != nil {
		return authmodel.Session{}, err
	}
	user.Roles = roles
	accessToken, expiresAt, err := s.tokens.CreateToken(user.ID, user.LoginName, authmodel.TokenTypeAccess)
	if err != nil {
		return authmodel.Session{}, err
	}
	nextRefreshToken, refreshExpiresAt, err := s.tokens.CreateToken(user.ID, user.LoginName, authmodel.TokenTypeRefresh)
	if err != nil {
		return authmodel.Session{}, err
	}
	return authmodel.Session{
		AccessToken:      accessToken,
		RefreshToken:     nextRefreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: refreshExpiresAt,
		User:             *user,
	}, nil
}

func (s *Service) Profile(ctx context.Context, userID uint64, loginName string) (authmodel.User, error) {
	user, err := s.repo.FindUserByLogin(ctx, loginName)
	if err != nil {
		return authmodel.User{}, err
	}
	if user.ID != userID {
		return authmodel.User{}, authmodel.ErrInvalidToken
	}
	roles, err := s.repo.ListRolesByUserID(ctx, user.ID)
	if err != nil {
		return authmodel.User{}, err
	}
	user.Roles = roles
	return *user, nil
}

func (s *Service) Routes(ctx context.Context, userID uint64) ([]authmodel.Menu, error) {
	menus, err := s.repo.ListMenusByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return BuildMenuTree(menus), nil
}

func (s *Service) Actions(ctx context.Context, userID uint64) ([]authmodel.Action, error) {
	return s.repo.ListActionsByUserID(ctx, userID)
}

func BuildMenuTree(menus []authmodel.Menu) []authmodel.Menu {
	childrenByParent := map[uint64][]authmodel.Menu{}
	for _, menu := range menus {
		childrenByParent[menu.ParentID] = append(childrenByParent[menu.ParentID], menu)
	}
	var attach func(parentID uint64) []authmodel.Menu
	attach = func(parentID uint64) []authmodel.Menu {
		items := childrenByParent[parentID]
		for i := range items {
			items[i].Children = attach(items[i].ID)
		}
		return items
	}
	return attach(0)
}
