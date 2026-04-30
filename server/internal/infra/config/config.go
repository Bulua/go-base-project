package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "configs/config.yaml"

var (
	ErrConfigNotFound = errors.New("config file not found")
	ErrInvalidConfig  = errors.New("invalid config")
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Env  string `yaml:"env"`
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret                 string `yaml:"secret"`
	AccessTokenTTLMinutes  int    `yaml:"access_token_ttl_minutes"`
	RefreshTokenTTLMinutes int    `yaml:"refresh_token_ttl_minutes"`
}

func Load(path string) (Config, error) {
	var cfg Config
	if path == "" {
		path = DefaultPath
	}
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			applyEnv(&cfg)
			if err := validate(cfg); err == nil {
				return cfg, nil
			}
			return Config{}, fmt.Errorf("%w: %s", ErrConfigNotFound, path)
		}
		return Config{}, err
	} else if err := yaml.Unmarshal(content, &cfg); err != nil {
		return Config{}, err
	}
	applyEnv(&cfg)
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.Charset,
	)
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c JWTConfig) AccessTokenTTL() time.Duration {
	return time.Duration(c.AccessTokenTTLMinutes) * time.Minute
}

func (c JWTConfig) RefreshTokenTTL() time.Duration {
	return time.Duration(c.RefreshTokenTTLMinutes) * time.Minute
}

func applyEnv(cfg *Config) {
	setString(&cfg.App.Name, "APP_NAME")
	setString(&cfg.App.Env, "APP_ENV")
	setString(&cfg.App.Port, "APP_PORT")

	setString(&cfg.Database.Host, "DB_HOST")
	setString(&cfg.Database.Port, "DB_PORT")
	setString(&cfg.Database.Name, "DB_NAME")
	setString(&cfg.Database.User, "DB_USER")
	setString(&cfg.Database.Password, "DB_PASSWORD")
	setString(&cfg.Database.Charset, "DB_CHARSET")

	setString(&cfg.Redis.Host, "REDIS_HOST")
	setString(&cfg.Redis.Port, "REDIS_PORT_IN_CONTAINER")
	setString(&cfg.Redis.Password, "REDIS_PASSWORD")
	setInt(&cfg.Redis.DB, "REDIS_DB")

	setString(&cfg.JWT.Secret, "JWT_SECRET")
	setInt(&cfg.JWT.AccessTokenTTLMinutes, "JWT_ACCESS_TTL_MINUTES")
	setInt(&cfg.JWT.RefreshTokenTTLMinutes, "JWT_REFRESH_TTL_MINUTES")
}

func validate(cfg Config) error {
	required := map[string]string{
		"app.name":          cfg.App.Name,
		"app.env":           cfg.App.Env,
		"app.port":          cfg.App.Port,
		"database.host":     cfg.Database.Host,
		"database.port":     cfg.Database.Port,
		"database.name":     cfg.Database.Name,
		"database.user":     cfg.Database.User,
		"database.password": cfg.Database.Password,
		"database.charset":  cfg.Database.Charset,
		"redis.host":        cfg.Redis.Host,
		"redis.port":        cfg.Redis.Port,
		"jwt.secret":        cfg.JWT.Secret,
	}
	var missing []string
	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, key)
		}
	}
	if cfg.JWT.AccessTokenTTLMinutes <= 0 {
		missing = append(missing, "jwt.access_token_ttl_minutes")
	}
	if cfg.JWT.RefreshTokenTTLMinutes <= 0 {
		missing = append(missing, "jwt.refresh_token_ttl_minutes")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing %s", ErrInvalidConfig, strings.Join(missing, ", "))
	}
	return nil
}

func setString(target *string, key string) {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		*target = value
	}
}

func setInt(target *int, key string) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return
	}
	value, err := strconv.Atoi(raw)
	if err == nil {
		*target = value
	}
}
