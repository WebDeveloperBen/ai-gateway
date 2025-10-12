package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/gateway"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
	"github.com/stretchr/testify/require"
)

func setupRegistry(t *testing.T) (*gateway.Registry, func()) {
	ctx := context.Background()
	kvStore := kv.NewMemoryStore()
	reg := gateway.NewRegistry(ctx, kvStore)
	return reg, func() {
		// Cleanup if needed
	}
}

func TestRegistry_AddAndGet(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	md := model.ModelDeployment{
		Tenant:   "tenant1",
		Model:    "gpt-4",
		Provider: "openai",
		Meta:     map[string]string{"endpoint": "https://api.openai.com"},
	}

	// Add deployment
	err := reg.Add(md, 0)
	require.NoError(t, err)

	// Get deployment
	retrieved, exists, err := reg.Get("gpt-4", "tenant1")
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, md, retrieved)
}

func TestRegistry_GetNonExistent(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	retrieved, exists, err := reg.Get("non-existent", "tenant1")
	require.NoError(t, err)
	require.False(t, exists)
	require.Equal(t, model.ModelDeployment{}, retrieved)
}

func TestRegistry_Update(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	md := model.ModelDeployment{
		Tenant:   "tenant2",
		Model:    "gpt-3.5",
		Provider: "openai",
		Meta:     map[string]string{"token": "sk-old"},
	}

	// Add initial
	err := reg.Add(md, 0)
	require.NoError(t, err)

	// Update
	md.Meta["token"] = "sk-new"
	err = reg.Update(md, 0)
	require.NoError(t, err)

	// Verify update
	retrieved, exists, err := reg.Get("gpt-3.5", "tenant2")
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, "sk-new", retrieved.Meta["token"])
}

func TestRegistry_Remove(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	md := model.ModelDeployment{
		Tenant:   "tenant3",
		Model:    "claude",
		Provider: "anthropic",
		Meta:     map[string]string{"endpoint": "https://api.anthropic.com"},
	}

	// Add
	err := reg.Add(md, 0)
	require.NoError(t, err)

	// Verify exists
	_, exists, err := reg.Get("claude", "tenant3")
	require.NoError(t, err)
	require.True(t, exists)

	// Remove
	err = reg.Remove("claude", "tenant3")
	require.NoError(t, err)

	// Verify gone
	_, exists, err = reg.Get("claude", "tenant3")
	require.NoError(t, err)
	require.False(t, exists)
}

func TestRegistry_AllWithPattern(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	// Add multiple deployments
	deployments := []model.ModelDeployment{
		{Tenant: "tenant1", Model: "gpt-4", Provider: "openai", Meta: map[string]string{}},
		{Tenant: "tenant1", Model: "claude", Provider: "anthropic", Meta: map[string]string{}},
		{Tenant: "tenant2", Model: "gpt-4", Provider: "azure", Meta: map[string]string{}},
	}

	for _, md := range deployments {
		err := reg.Add(md, 0)
		require.NoError(t, err)
	}

	// Get all
	all, err := reg.All("modelreg:*")
	require.NoError(t, err)
	require.Len(t, all, 3)

	// Check that all deployments are present
	models := make(map[string]bool)
	for _, md := range all {
		key := md.Tenant + ":" + md.Model
		models[key] = true
	}
	require.True(t, models["tenant1:gpt-4"])
	require.True(t, models["tenant1:claude"])
	require.True(t, models["tenant2:gpt-4"])
}

func TestRegistry_DeploymentsForModel(t *testing.T) {
	reg, cleanup := setupRegistry(t)
	defer cleanup()

	// Add deployments
	deployments := []model.ModelDeployment{
		{Tenant: "tenant1", Model: "gpt-4", Provider: "openai", Meta: map[string]string{"endpoint": "https://api1.openai.com"}},
		{Tenant: "tenant1", Model: "gpt-4-turbo", Provider: "azure", Meta: map[string]string{"endpoint": "https://api2.azure.com"}},
		{Tenant: "tenant2", Model: "gpt-4", Provider: "openai", Meta: map[string]string{"endpoint": "https://api3.openai.com"}},
		{Tenant: "tenant1", Model: "claude", Provider: "anthropic", Meta: map[string]string{"endpoint": "https://api.anthropic.com"}},
	}

	for _, md := range deployments {
		err := reg.Add(md, time.Hour)
		require.NoError(t, err)
	}

	// Get all gpt-4 deployments
	results, err := reg.DeploymentsForModel("gpt-4", "")
	require.NoError(t, err)
	require.Len(t, results, 2) // tenant1/gpt-4 and tenant2/gpt-4

	// Get gpt-4 for tenant1 only
	results, err = reg.DeploymentsForModel("gpt-4", "tenant1")
	require.NoError(t, err)
	require.Len(t, results, 1)

	// Verify tenant1 results
	require.Equal(t, "tenant1", results[0].Tenant)
	require.Equal(t, "gpt-4", results[0].Model)
	require.Equal(t, "https://api1.openai.com", results[0].Meta["endpoint"])

	// Get claude for tenant1
	results, err = reg.DeploymentsForModel("claude", "tenant1")
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, "claude", results[0].Model)
	require.Equal(t, "tenant1", results[0].Tenant)
}
