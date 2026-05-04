// migrate applies SQL migration files from deploy/mysql/init/ in order.
// Usage: go run ./cmd/migrate [path-to-init-dir]
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gobaseproject/server/internal/infra/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	initDir := "../deploy/mysql/init"
	if len(os.Args) > 1 {
		initDir = os.Args[1]
	}

	db, err := sql.Open("mysql", cfg.Database.DSN()+"&multiStatements=true")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	entries, err := os.ReadDir(initDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read dir %s: %v\n", initDir, err)
		os.Exit(1)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, filepath.Join(initDir, e.Name()))
		}
	}
	sort.Strings(files)

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", f, err)
			os.Exit(1)
		}
		fmt.Printf("applying %s ...\n", filepath.Base(f))
		if _, err := db.Exec(string(content)); err != nil {
			fmt.Fprintf(os.Stderr, "exec %s: %v\n", f, err)
			os.Exit(1)
		}
		fmt.Printf("  done\n")
	}
	fmt.Println("migration complete")
}
