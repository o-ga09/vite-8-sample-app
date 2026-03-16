package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func moduleRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current directory: %v", err)
	}

	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

	return string(content)
}

func TestIssue1RequiredFilesExist(t *testing.T) {
	root := moduleRoot(t)

	required := []string{
		".github/workflows/ci.yml",
		"bobgen.yaml",
		"queries/reports.sql",
		"db/migrations/001_create_enums.sql",
		"db/migrations/002_create_workspaces.sql",
		"db/migrations/003_create_members.sql",
		"db/migrations/004_create_accounts.sql",
		"db/migrations/005_create_categories.sql",
		"db/migrations/006_create_transactions.sql",
		"db/migrations/007_create_triggers.sql",
	}

	for _, rel := range required {
		fullPath := filepath.Join(root, rel)
		if _, err := os.Stat(fullPath); err != nil {
			t.Fatalf("required file %s does not exist: %v", rel, err)
		}
	}
}

func TestIssue1MisePortsAndTasks(t *testing.T) {
	root := moduleRoot(t)
	mise := readFile(t, filepath.Join(root, "mise.toml"))

	expected := []string{
		"POSTGRES_PORT = \"15432\"",
		"GRAFANA_PORT = \"13000\"",
		"LOKI_PORT = \"13100\"",
		"APP_PORT = \"18080\"",
		"WEBAPP_PORT = \"13001\"",
		"[tasks.\"migrate-down\"]",
		"[tasks.lint]",
		"[tasks.test]",
		"[tasks.build]",
		"[tasks.clean]",
	}

	for _, e := range expected {
		if !strings.Contains(mise, e) {
			t.Fatalf("mise.toml does not contain required entry: %s", e)
		}
	}
}

func TestIssue1ComposeDefaultPorts(t *testing.T) {
	root := moduleRoot(t)
	compose := readFile(t, filepath.Join(root, "compose.yaml"))

	expected := []string{
		"${POSTGRES_PORT:-15432}:5432",
		"${GRAFANA_PORT:-13000}:3000",
		"${LOKI_PORT:-13100}:3100",
		"postgres:",
		"grafana:",
		"loki:",
		"promtail:",
	}

	for _, e := range expected {
		if !strings.Contains(compose, e) {
			t.Fatalf("compose.yaml does not contain required entry: %s", e)
		}
	}
}
