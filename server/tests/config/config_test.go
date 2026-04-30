package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gobaseproject/server/internal/infra/config"
)

func TestLoadReadsYAMLConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
app:
  name: TestProject
  env: local
  port: 18080
database:
  host: 127.0.0.1
  port: 3306
  name: test_db
  user: test_user
  password: test_password
  charset: utf8mb4
redis:
  host: 127.0.0.1
  port: 6379
  password: redis_password
  db: 2
jwt:
  secret: yaml-secret
  access_token_ttl_minutes: 30
  refresh_token_ttl_minutes: 1440
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.App.Port != "18080" {
		t.Fatalf("expected app port 18080, got %q", cfg.App.Port)
	}
	if cfg.Database.Host != "127.0.0.1" {
		t.Fatalf("expected database host from yaml, got %q", cfg.Database.Host)
	}
	if cfg.JWT.Secret != "yaml-secret" {
		t.Fatalf("expected jwt secret from yaml, got %q", cfg.JWT.Secret)
	}
	if cfg.Redis.Addr() != "127.0.0.1:6379" {
		t.Fatalf("expected redis addr, got %q", cfg.Redis.Addr())
	}
}

func TestLoadLetsEnvironmentOverrideYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
app:
  name: TestProject
  env: local
  port: 18080
database:
  host: yaml-host
  port: 3306
  name: yaml_db
  user: yaml_user
  password: yaml_password
  charset: utf8mb4
redis:
  host: yaml-redis
  port: 6379
  password: yaml_redis_password
  db: 0
jwt:
  secret: yaml-secret
  access_token_ttl_minutes: 30
  refresh_token_ttl_minutes: 1440
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("APP_PORT", "19090")
	t.Setenv("DB_HOST", "env-host")
	t.Setenv("JWT_SECRET", "env-secret")

	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.App.Port != "19090" {
		t.Fatalf("expected env app port, got %q", cfg.App.Port)
	}
	if cfg.Database.Host != "env-host" {
		t.Fatalf("expected env database host, got %q", cfg.Database.Host)
	}
	if cfg.JWT.Secret != "env-secret" {
		t.Fatalf("expected env jwt secret, got %q", cfg.JWT.Secret)
	}
}

func TestLoadFailsWhenConfigFileIsMissing(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if !errors.Is(err, config.ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestLoadCanUseEnvironmentWhenConfigFileIsMissing(t *testing.T) {
	t.Setenv("APP_NAME", "EnvProject")
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "18080")
	t.Setenv("DB_HOST", "env-db")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_NAME", "env_db")
	t.Setenv("DB_USER", "env_user")
	t.Setenv("DB_PASSWORD", "env_password")
	t.Setenv("DB_CHARSET", "utf8mb4")
	t.Setenv("REDIS_HOST", "env-redis")
	t.Setenv("REDIS_PORT_IN_CONTAINER", "6379")
	t.Setenv("JWT_SECRET", "env-secret")
	t.Setenv("JWT_ACCESS_TTL_MINUTES", "30")
	t.Setenv("JWT_REFRESH_TTL_MINUTES", "1440")

	cfg, err := config.Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("load env-only config: %v", err)
	}

	if cfg.App.Name != "EnvProject" {
		t.Fatalf("expected app name from env, got %q", cfg.App.Name)
	}
	if cfg.Database.Host != "env-db" {
		t.Fatalf("expected database host from env, got %q", cfg.Database.Host)
	}
	if cfg.JWT.AccessTokenTTLMinutes != 30 {
		t.Fatalf("expected jwt ttl from env, got %d", cfg.JWT.AccessTokenTTLMinutes)
	}
}

func TestLoadFailsWhenRequiredDatabaseConfigIsMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`
app:
  name: TestProject
  env: local
  port: 18080
jwt:
  secret: yaml-secret
  access_token_ttl_minutes: 30
  refresh_token_ttl_minutes: 1440
`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := config.Load(path)
	if !errors.Is(err, config.ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}
