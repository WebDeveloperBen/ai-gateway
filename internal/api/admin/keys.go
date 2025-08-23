package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/insurgence-ai/llm-gateway/internal/admin/models"
	"github.com/insurgence-ai/llm-gateway/internal/admin/services"
)

type Server struct {
	Keys services.KeysService
}

func NewServer(keys services.KeysService) *Server {
	return &Server{Keys: keys}
}

func (s *Server) RegisterRoutes(grp *huma.Group) {
	// POST /admin/keys
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-mint-key",
		Method:        http.MethodPost,
		Path:          "/admin/keys",
		Summary:       "Mint API key",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		Tenant   string         `json:"tenant" required:"true"`
		App      string         `json:"app" required:"true"`
		TTL      *time.Duration `json:"ttl,omitempty"`
		Prefix   string         `json:"prefix,omitempty"`
		Metadata map[string]any `json:"metadata,omitempty"`
	}) (*struct {
		Token string        `json:"token"`
		Key   models.APIKey `json:"key"`
	}, error,
	) {
		var ttl time.Duration
		if in.TTL != nil {
			ttl = *in.TTL
		}
		out, err := s.Keys.MintKey(ctx, models.MintKeyRequest{
			Tenant:   in.Tenant,
			App:      in.App,
			TTL:      ttl,
			Prefix:   in.Prefix,
			Metadata: in.Metadata,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("mint failed")
		}
		return &struct {
			Token string        `json:"token"`
			Key   models.APIKey `json:"key"`
		}{Token: out.Token, Key: out.Key}, nil
	})

	// POST /admin/keys/{key_id}/revoke
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-revoke-key",
		Method:        http.MethodPost,
		Path:          "/admin/keys/{key_id}/revoke",
		Summary:       "Revoke API key",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"Admin"},
	}, func(ctx context.Context, in *struct {
		KeyID string `path:"key_id" required:"true"`
	},
	) (*struct{}, error) {
		if err := s.Keys.RevokeKey(ctx, in.KeyID); err != nil {
			return nil, huma.Error404NotFound("not found")
		}
		return &struct{}{}, nil
	})
}
