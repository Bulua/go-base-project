package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	apihandler "gobaseproject/server/internal/handler/api"
	audithandler "gobaseproject/server/internal/handler/audit"
	authhandler "gobaseproject/server/internal/handler/auth"
	dicthandler "gobaseproject/server/internal/handler/dict"
	menuhandler "gobaseproject/server/internal/handler/menu"
	rolehandler "gobaseproject/server/internal/handler/role"
	userhandler "gobaseproject/server/internal/handler/user"
	"gobaseproject/server/internal/infra/config"
	"gobaseproject/server/internal/middleware"
	apirepo "gobaseproject/server/internal/repository/api"
	auditrepo "gobaseproject/server/internal/repository/audit"
	authrepo "gobaseproject/server/internal/repository/auth"
	dictrepo "gobaseproject/server/internal/repository/dict"
	menurepo "gobaseproject/server/internal/repository/menu"
	rolerepo "gobaseproject/server/internal/repository/role"
	userrepo "gobaseproject/server/internal/repository/user"
	apiservice "gobaseproject/server/internal/service/api"
	auditservice "gobaseproject/server/internal/service/audit"
	authservice "gobaseproject/server/internal/service/auth"
	dictservice "gobaseproject/server/internal/service/dict"
	menuservice "gobaseproject/server/internal/service/menu"
	roleservice "gobaseproject/server/internal/service/role"
	userservice "gobaseproject/server/internal/service/user"
	"gobaseproject/server/pkg/response"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}

	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}
	db, err := sql.Open("mysql", cfg.Database.DSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	tokenManager := authservice.NewTokenManager(authservice.TokenConfig{
		Secret:          cfg.JWT.Secret,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL(),
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL(),
		TokenIssuer:     cfg.App.Name,
	})

	authService := authservice.NewService(authrepo.NewSQLRepository(db), tokenManager)

	auditRepository := auditrepo.NewSQLRepository(db)
	auditSvc := auditservice.NewService(auditRepository)
	apiSvc := apiservice.NewService(apirepo.NewSQLRepository(db), auditRepository)
	userService := userservice.NewService(userrepo.NewSQLRepository(db), auditRepository)
	roleService := roleservice.NewService(rolerepo.NewSQLRepository(db), auditRepository)
	menuService := menuservice.NewService(menurepo.NewSQLRepository(db))
	dictService := dictservice.NewService(dictrepo.NewSQLRepository(db))

	mux := http.NewServeMux()
	registerSystemRoutes(mux, cfg)
	authhandler.NewHandler(authService).RegisterRoutes(mux)
	userhandler.NewHandler(userService, authService).RegisterRoutes(mux)
	rolehandler.NewHandler(roleService, authService).RegisterRoutes(mux)
	menuhandler.NewHandler(menuService, authService, cfg.CodeGen.WebSrcRoot).RegisterRoutes(mux)
	dicthandler.NewHandler(dictService, authService).RegisterRoutes(mux)
	audithandler.NewHandler(auditSvc).RegisterRoutes(mux)
	apihandler.NewHandler(apiSvc, authService).RegisterRoutes(mux)

	server := &http.Server{
		Addr: ":" + cfg.App.Port,
		Handler: withCORS(
			middleware.Permission(db, tokenManager)(
				middleware.OperationLog(auditRepository, tokenManager)(mux),
			),
		),
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("GoBaseProject server listening on :%s\n", cfg.App.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

func registerSystemRoutes(mux *http.ServeMux, cfg config.Config) {
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, response.Body{
			Code:    0,
			Message: "ok",
			Data: map[string]string{
				"app": cfg.App.Name,
				"env": cfg.App.Env,
				"now": time.Now().Format(time.RFC3339),
			},
		})
	})

	mux.HandleFunc("/api/v1/system/runtime", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, response.Body{
			Code:    0,
			Message: "ok",
			Data: map[string]string{
				"database_host": cfg.Database.Host,
				"database_name": cfg.Database.Name,
				"redis_addr":    cfg.Redis.Addr(),
			},
		})
	})
}

func healthcheck() int {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + cfg.App.Port + "/api/v1/health")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "unexpected status: %s\n", resp.Status)
		return 1
	}
	return 0
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Trace-ID")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
