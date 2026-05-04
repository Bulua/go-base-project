package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateVueFile creates a Vue3 template file at webSrcRoot/<componentPath>
// if the file does not already exist.
// Returns (filePath, created, error).
// If webSrcRoot is empty or the views directory does not exist, it is a no-op.
func GenerateVueFile(webSrcRoot, componentPath, menuTitle string) (string, bool, error) {
	if webSrcRoot == "" || componentPath == "" {
		return "", false, nil
	}

	// normalize: component_path may start with "views/" or "/views/"
	rel := strings.TrimPrefix(componentPath, "/")
	if !strings.HasPrefix(rel, "views/") {
		rel = "views/" + rel
	}
	rel = strings.TrimSuffix(rel, ".vue") + ".vue"

	absPath := filepath.Join(webSrcRoot, rel)

	// skip if file already exists
	if _, err := os.Stat(absPath); err == nil {
		return absPath, false, nil
	}

	// skip silently if the views root does not exist (e.g. Docker)
	viewsRoot := filepath.Join(webSrcRoot, "views")
	if _, err := os.Stat(viewsRoot); err != nil {
		return "", false, nil
	}

	if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
		return "", false, fmt.Errorf("codegen mkdir: %w", err)
	}

	content := vueTemplate(menuTitle)
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return "", false, fmt.Errorf("codegen write: %w", err)
	}

	return absPath, true, nil
}

func vueTemplate(title string) string {
	if title == "" {
		title = "新页面"
	}
	return fmt.Sprintf(`<script setup lang="ts">
// TODO: implement %s page
</script>

<template>
  <div style="padding: 16px">
    <h2>%s</h2>
    <p style="color: #909399">该页面待实现。</p>
  </div>
</template>
`, title, title)
}
