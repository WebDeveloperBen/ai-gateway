package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/insurgence-ai/llm-gateway/internal/kv"
	"github.com/insurgence-ai/llm-gateway/internal/model"
)

type Registry struct {
	kv  kv.Store
	ctx context.Context
}

func NewRegistry(ctx context.Context, store kv.Store) *Registry {
	return &Registry{kv: store, ctx: ctx}
}

func (r *Registry) Add(md model.ModelDeployment, ttl time.Duration) error {
	b, err := json.Marshal(md)
	if err != nil {
		return err
	}
	return r.kv.Set(r.ctx, r.key(md.Model, md.Tenant), string(b), ttl)
}

func (r *Registry) Update(md model.ModelDeployment, ttl time.Duration) error {
	return r.Add(md, ttl)
}

func (r *Registry) Remove(model, tenant string) error {
	return r.kv.Del(r.ctx, r.key(model, tenant))
}

func (r *Registry) Get(mod, tenant string) (model.ModelDeployment, bool, error) {
	val, err := r.kv.Get(r.ctx, r.key(mod, tenant))
	if err != nil || val == "" {
		return model.ModelDeployment{}, false, err
	}
	var md model.ModelDeployment
	if err := json.Unmarshal([]byte(val), &md); err != nil {
		return model.ModelDeployment{}, false, err
	}
	return md, true, nil
}

func (r *Registry) All(pattern string) ([]model.ModelDeployment, error) {
	keys, err := r.kv.Keys(r.ctx, pattern)
	if err != nil {
		return nil, err
	}
	var all []model.ModelDeployment
	for _, k := range keys {
		v, _ := r.kv.Get(r.ctx, k)
		if v == "" {
			continue
		}
		var md model.ModelDeployment
		if json.Unmarshal([]byte(v), &md) == nil {
			all = append(all, md)
		}
	}
	return all, nil
}

func (r *Registry) key(model, tenant string) string {
	return fmt.Sprintf("modelreg:%s:%s", tenant, model)
}

// EnsureRegistryPopulated loads all model deployments from the registry, and if empty, calls the provided loader to seed the registry with initial deployments.
// Returns the final set of all model deployments available in the registry.
func EnsureRegistryPopulated(reg *Registry, loadFn func() []model.ModelDeployment) []model.ModelDeployment {
	all, err := reg.All("modelreg:*")
	if err != nil {
		log.Fatal("registry read failed: ", err)
	}
	if len(all) == 0 {
		modelDeployments := loadFn()
		for _, md := range modelDeployments {
			if err := reg.Add(md, 0); err != nil {
				log.Printf("failed to add model to registry: %+v", err)
			}
		}
		all = modelDeployments
	}
	log.Printf("Loaded %d active models from registry", len(all))
	return all
}
