package testkit

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

func LoadDotenvFromRepoRoot(t *testing.T, names ...string) string {
	t.Helper()

	if len(names) == 0 {
		names = []string{".env", ".env.e2e", ".env.local"}
	}

	start, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	root, err := findRepoRoot(start)
	if err != nil {
		t.Fatalf("find repo root from %q: %v", start, err)
	}

	for _, n := range names {
		p := filepath.Join(root, n)
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err != nil {
				t.Fatalf("failed to load %s: %v", p, err)
			}
			t.Logf("Loaded dotenv: %s (wd=%s)", p, start)
			return p
		}
	}
	t.Fatalf("no dotenv file found at %s (looked for %v)", root, names)
	return ""
}

func findRepoRoot(dir string) (string, error) {
	for {
		if exists(filepath.Join(dir, "go.mod")) || exists(filepath.Join(dir, ".git")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("reached filesystem root")
		}
		dir = parent
	}
}
func exists(p string) bool { _, err := os.Stat(p); return err == nil }
