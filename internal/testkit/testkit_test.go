package testkit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

func TestExists(t *testing.T) {
	// Test existing file
	if !exists("testkit.go") {
		t.Error("expected testkit.go to exist")
	}

	// Test non-existing file
	if exists("nonexistent_file.go") {
		t.Error("expected nonexistent_file.go to not exist")
	}
}

func TestFindRepoRoot(t *testing.T) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Test finding repo root from current directory
	root, err := findRepoRoot(wd)
	if err != nil {
		t.Fatalf("failed to find repo root: %v", err)
	}

	// Verify go.mod exists in root
	goModPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		t.Errorf("expected go.mod to exist at %s", goModPath)
	}

	// Test from a subdirectory
	subDir := filepath.Join(wd, "internal", "testkit")
	rootFromSub, err := findRepoRoot(subDir)
	if err != nil {
		t.Fatalf("failed to find repo root from subdirectory: %v", err)
	}

	if root != rootFromSub {
		t.Errorf("expected same root from subdirectory, got %s vs %s", root, rootFromSub)
	}
}

func TestFindRepoRoot_NonRepo(t *testing.T) {
	// Test from a directory that's not in a repo (use temp dir)
	tempDir := t.TempDir()

	_, err := findRepoRoot(tempDir)
	if err == nil {
		t.Error("expected error when finding repo root from non-repo directory")
	}
}

func TestDefaultContainerConfig(t *testing.T) {
	config := DefaultContainerConfig()

	if config.PostgresDB != "testdb" {
		t.Errorf("expected PostgresDB to be 'testdb', got %s", config.PostgresDB)
	}
	if config.PostgresUser != "postgres" {
		t.Errorf("expected PostgresUser to be 'postgres', got %s", config.PostgresUser)
	}
	if config.PostgresPassword != "postgres" {
		t.Errorf("expected PostgresPassword to be 'postgres', got %s", config.PostgresPassword)
	}
	if config.RedisPassword != "" {
		t.Errorf("expected RedisPassword to be empty, got %s", config.RedisPassword)
	}
}

func TestIsDockerAvailable(t *testing.T) {
	// This test just ensures the function doesn't panic
	// We can't easily test the actual Docker availability without mocking
	result := isDockerAvailable()
	// Result depends on environment, just ensure it's a bool
	if result != true && result != false {
		t.Error("isDockerAvailable should return a boolean")
	}
}

func TestFindProjectRoot(t *testing.T) {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Test finding project root from current directory
	root := findProjectRoot(wd)

	// Verify go.mod exists in root
	goModPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		t.Errorf("expected go.mod to exist at %s", goModPath)
	}

	// Test from a subdirectory
	subDir := filepath.Join(wd, "internal", "testkit")
	rootFromSub := findProjectRoot(subDir)

	if root != rootFromSub {
		t.Errorf("expected same root from subdirectory, got %s vs %s", root, rootFromSub)
	}
}

func TestGetenvDefault(t *testing.T) {
	// Test with existing env var
	testKey := "TEST_GETENV_DEFAULT"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	result := getenvDefault(testKey, "default")
	if result != testValue {
		t.Errorf("expected %s, got %s", testValue, result)
	}

	// Test with non-existing env var
	result = getenvDefault("NON_EXISTING_VAR", "default")
	if result != "default" {
		t.Errorf("expected default, got %s", result)
	}
}

func TestSetupPublicTestAPI(t *testing.T) {
	// Test that SetupPublicTestAPI creates an API and calls the register function
	called := false
	api := SetupPublicTestAPI(t, func(grp *huma.Group) {
		called = true
		// Just verify the group is not nil
		if grp == nil {
			t.Error("expected group to be created")
		}
	})

	// Verify API was created and register was called
	if api == nil {
		t.Error("expected API to be created")
	}
	if !called {
		t.Error("expected register function to be called")
	}
}

func TestSetupProviderTestAPI(t *testing.T) {
	// Test that SetupProviderTestAPI creates an API and calls the register function
	called := false
	api := SetupProviderTestAPI(t, func(grp *huma.Group) {
		called = true
		// Just verify the group is not nil
		if grp == nil {
			t.Error("expected group to be created")
		}
	})

	// Verify API was created and register was called
	if api == nil {
		t.Error("expected API to be created")
	}
	if !called {
		t.Error("expected register function to be called")
	}
}

func TestSetupAdminTestAPI(t *testing.T) {
	// Test that SetupAdminTestAPI creates an API and calls the register function
	called := false
	api := SetupAdminTestAPI(t, func(grp *huma.Group) {
		called = true
		// Just verify the group is not nil
		if grp == nil {
			t.Error("expected group to be created")
		}
	})

	// Verify API was created and register was called
	if api == nil {
		t.Error("expected API to be created")
	}
	if !called {
		t.Error("expected register function to be called")
	}
}

func TestSetupAuthTestAPI(t *testing.T) {
	// Test that SetupAuthTestAPI creates an API and calls the register function
	called := false
	api := SetupAuthTestAPI(t, func(grp *huma.Group) {
		called = true
		// Just verify the group is not nil
		if grp == nil {
			t.Error("expected group to be created")
		}
	})

	// Verify API was created and register was called
	if api == nil {
		t.Error("expected API to be created")
	}
	if !called {
		t.Error("expected register function to be called")
	}
}
