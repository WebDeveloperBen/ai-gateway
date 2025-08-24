// Package testkit implements testing utilities to keep tests DRY
package testkit

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

// SetupPublicTestAPI creates a test Huma API instance with the provided group registration.
func SetupPublicTestAPI(
	t *testing.T,
	register func(grp *huma.Group),
) humatest.TestAPI {
	_, api := humatest.New(t)
	group := huma.NewGroup(api, "/api")
	register(group)
	return api
}
