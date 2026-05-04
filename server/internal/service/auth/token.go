package authservice

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	authmodel "gobaseproject/server/internal/model/auth"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	TokenIssuer     string
	Clock           func() time.Time
}

type TokenManager struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
	clock           func() time.Time
}

type Claims struct {
	UserID    uint64   `json:"uid"`
	LoginName string   `json:"login_name"`
	TokenType string   `json:"typ"`
	RoleIDs   []uint64 `json:"role_ids,omitempty"`
	jwt.RegisteredClaims
}

func NewTokenManager(config TokenConfig) TokenManager {
	clock := config.Clock
	if clock == nil {
		clock = time.Now
	}
	if config.AccessTokenTTL <= 0 {
		config.AccessTokenTTL = 2 * time.Hour
	}
	if config.RefreshTokenTTL <= 0 {
		config.RefreshTokenTTL = 7 * 24 * time.Hour
	}
	if config.TokenIssuer == "" {
		config.TokenIssuer = "gobaseproject"
	}
	return TokenManager{
		secret:          []byte(config.Secret),
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.TokenIssuer,
		clock:           clock,
	}
}

func (m TokenManager) CreateToken(userID uint64, loginName string, tokenType string, roleIDs []uint64) (string, time.Time, error) {
	if len(m.secret) == 0 {
		return "", time.Time{}, errors.New("jwt secret is required")
	}
	now := m.clock().UTC()
	ttl := m.accessTokenTTL
	if tokenType == authmodel.TokenTypeRefresh {
		ttl = m.refreshTokenTTL
	}
	expiresAt := now.Add(ttl)
	claims := Claims{
		UserID:    userID,
		LoginName: loginName,
		TokenType: tokenType,
		RoleIDs:   roleIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	value, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return value, expiresAt, nil
}

func (m TokenManager) Parse(tokenValue string, expectedType string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, authmodel.ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithTimeFunc(m.clock))
	if err != nil || token == nil || !token.Valid {
		return nil, authmodel.ErrInvalidToken
	}
	if expectedType != "" && claims.TokenType != expectedType {
		return nil, authmodel.ErrInvalidToken
	}
	return claims, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
