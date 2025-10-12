// Package testkit implements testing utilities to keep tests DRY
package testkit

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/joho/godotenv"
)

// SetupPublicTestAPI creates a test Huma API instance with the provided group registration.
func SetupPublicTestAPI(
	t *testing.T,
	register func(grp *huma.Group),
) humatest.TestAPI {
	_ = godotenv.Load(".env", "../.env", "../../.env", "../../../.env")

	// Load .env from the repo root before constructing deps that read env.
	LoadDotenvFromRepoRoot(t)

	_, api := humatest.New(t)
	group := huma.NewGroup(api, "/api")
	register(group)
	return api
}

// SetupProviderTestAPI creates a test Huma API instance for provider tests with /api/providers prefix.
func SetupProviderTestAPI(
	t *testing.T,
	register func(grp *huma.Group),
) humatest.TestAPI {
	_ = godotenv.Load(".env", "../.env", "../../.env", "../../../.env")

	// Load .env from the repo root before constructing deps that read env.
	LoadDotenvFromRepoRoot(t)

	_, api := humatest.New(t)
	group := huma.NewGroup(api, "/api/providers")
	register(group)
	return api
}

// SetupAdminTestAPI creates a test Huma API instance for admin tests with /api prefix.
func SetupAdminTestAPI(
	t *testing.T,
	register func(grp *huma.Group),
) humatest.TestAPI {
	_ = godotenv.Load(".env", "../.env", "../../.env", "../../../.env")

	// Load .env from the repo root before constructing deps that read env.
	LoadDotenvFromRepoRoot(t)

	_, api := humatest.New(t)
	group := huma.NewGroup(api, "/api")
	register(group)
	return api
}

// SetupAuthTestAPI creates a test Huma API instance for auth tests with /auth prefix.
func SetupAuthTestAPI(
	t *testing.T,
	register func(grp *huma.Group),
) humatest.TestAPI {
	_ = godotenv.Load(".env", "../.env", "../../.env", "../../../.env")

	// Load .env from the repo root before constructing deps that read env.
	LoadDotenvFromRepoRoot(t)

	_, api := humatest.New(t)
	group := huma.NewGroup(api, "/auth")
	register(group)
	return api
}
