package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/insurgence-ai/llm-gateway/internal/kv"
	"github.com/insurgence-ai/llm-gateway/internal/model/models"
)

type Registry struct {
	kv  kv.Store
	ctx context.Context
}

func NewRegistry(ctx context.Context, store kv.Store) *Registry {
	return &Registry{kv: store, ctx: ctx}
}

func (r *Registry) Add(md models.ModelDeployment, ttl time.Duration) error {
	b, err := json.Marshal(md)
	if err != nil {
		return err
	}
	return r.kv.Set(r.ctx, r.key(md.Model, md.Tenant), string(b), ttl)
}

func (r *Registry) Update(md models.ModelDeployment, ttl time.Duration) error {
	return r.Add(md, ttl)
}

func (r *Registry) Remove(model, tenant string) error {
	return r.kv.Del(r.ctx, r.key(model, tenant))
}

func (r *Registry) Get(model, tenant string) (models.ModelDeployment, bool, error) {
	val, err := r.kv.Get(r.ctx, r.key(model, tenant))
	if err != nil || val == "" {
		return models.ModelDeployment{}, false, err
	}
	var md models.ModelDeployment
	if err := json.Unmarshal([]byte(val), &md); err != nil {
		return models.ModelDeployment{}, false, err
	}
	return md, true, nil
}

func (r *Registry) All(pattern string) ([]models.ModelDeployment, error) {
	keys, err := r.kv.Keys(r.ctx, pattern)
	if err != nil {
		return nil, err
	}
	var all []models.ModelDeployment
	for _, k := range keys {
		v, _ := r.kv.Get(r.ctx, k)
		if v == "" {
			continue
		}
		var md models.ModelDeployment
		if json.Unmarshal([]byte(v), &md) == nil {
			all = append(all, md)
		}
	}
	return all, nil
}

func (r *Registry) key(model, tenant string) string {
	return fmt.Sprintf("modelreg:%s:%s", tenant, model)
}
