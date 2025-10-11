package gateway

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/WebDeveloperBen/ai-gateway/internal/drivers/kv"
	"github.com/WebDeveloperBen/ai-gateway/internal/model"
)

type Registry struct {
	kv  kv.KvStore
	ctx context.Context
}

func NewRegistry(ctx context.Context, store kv.KvStore) *Registry {
	return &Registry{kv: store, ctx: ctx}
}

/* ---------- small helpers to centralize key/pattern usage ---------- */

func (r *Registry) key(tenant, model string) string {
	return kv.KeyModel(tenant, model)
}

func (r *Registry) patternAll() string {
	return kv.PatternAll()
}

func (r *Registry) patternTenantAll(tenant string) string {
	return kv.PatternTenantAll(tenant)
}

/* ---------------------------- CRUD --------------------------------- */

func (r *Registry) Add(md model.ModelDeployment, ttl time.Duration) error {
	b, err := json.Marshal(md)
	if err != nil {
		return err
	}
	return r.kv.Set(r.ctx, r.key(md.Tenant, md.Model), string(b), ttl)
}

func (r *Registry) Update(md model.ModelDeployment, ttl time.Duration) error {
	return r.Add(md, ttl)
}

func (r *Registry) Remove(modelName, tenant string) error {
	return r.kv.Del(r.ctx, r.key(tenant, modelName))
}

func (r *Registry) Get(mod, tenant string) (model.ModelDeployment, bool, error) {
	val, err := r.kv.Get(r.ctx, r.key(tenant, mod))
	if err != nil {
		return model.ModelDeployment{}, false, err
	}
	if val == "" {
		return model.ModelDeployment{}, false, nil
	}
	var md model.ModelDeployment
	if err := json.Unmarshal([]byte(val), &md); err != nil {
		return model.ModelDeployment{}, false, err
	}
	return md, true, nil
}

/* ------------------------- listing / scans -------------------------- */

// All returns all deployments matching a Redis MATCH pattern using ScanGetAll.
// Prefer using pattern helpers: kv.PatternAll() or kv.PatternTenantAll(tenant).
func (r *Registry) All(pattern string) ([]model.ModelDeployment, error) {
	kvs, err := r.kv.ScanGetAll(r.ctx, pattern, 1024)
	if err != nil {
		return nil, err
	}
	out := make([]model.ModelDeployment, 0, len(kvs))
	for _, v := range kvs {
		var md model.ModelDeployment
		if err := json.Unmarshal([]byte(v), &md); err == nil {
			out = append(out, md)
		}
	}
	return out, nil
}

// DeploymentsForModel returns all deployments for a given model and (optional) tenant.
func (r *Registry) DeploymentsForModel(mod, tenant string) ([]model.ModelDeployment, error) {
	// Fast path: exact key if both provided
	if tenant != "" && mod != "" {
		md, ok, err := r.Get(mod, tenant)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, nil
		}
		return []model.ModelDeployment{md}, nil
	}

	// Otherwise scan the tightest namespace
	pattern := r.patternAll()
	if tenant != "" {
		pattern = r.patternTenantAll(tenant)
	}
	kvs, err := r.kv.ScanGetAll(r.ctx, pattern, 1024)
	if err != nil {
		return nil, err
	}

	var result []model.ModelDeployment
	for _, raw := range kvs {
		var d model.ModelDeployment
		if err := json.Unmarshal([]byte(raw), &d); err != nil {
			continue
		}
		if (mod == "" || d.Model == mod) && (tenant == "" || d.Tenant == tenant) {
			result = append(result, d)
		}
	}
	return result, nil
}

/* --------------------------- bootstrap ------------------------------ */

// EnsureRegistryPopulated loads all deployments; if empty, seeds via loadFn.
// Returns the final set of deployments.
func EnsureRegistryPopulated(reg *Registry, loadFn func() []model.ModelDeployment) []model.ModelDeployment {
	all, err := reg.All(reg.patternAll())
	if err != nil {
		log.Fatal("registry read failed: ", err)
	}
	if len(all) == 0 {
		seeds := loadFn()
		for _, md := range seeds {
			if err := reg.Add(md, 0); err != nil {
				log.Printf("failed to add model to registry: %+v", err)
			}
		}
		// re-read to reflect actual stored state
		all, err = reg.All(reg.patternAll())
		if err != nil {
			log.Fatal("registry read after seed failed: ", err)
		}
	}
	log.Printf("Loaded %d active models from registry", len(all))
	return all
}
