package keys

import (
	"context"
	"net/http"

	"github.com/WebDeveloperBen/ai-gateway/internal/exceptions"
	"github.com/danielgtaylor/huma/v2"
)

type KeyService struct {
	Keys KeysService
}

func NewRouter(keys KeysService) *KeyService {
	return &KeyService{Keys: keys}
}

func (s *KeyService) RegisterRoutes(grp *huma.Group) {
	// POST /keys
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-mint-key",
		Method:        http.MethodPost,
		Path:          "/keys",
		Summary:       "Mint API key",
		Description:   "Creates a new API key for accessing the gateway with specified tenant and application.",
		DefaultStatus: http.StatusCreated,
		Tags:          []string{"API Keys"},
	}, exceptions.Handle(func(ctx context.Context, in *MintKeyRequest) (*MintKeyResponse, error) {
		out, err := s.Keys.MintKey(ctx, in.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("mint failed")
		}
		return &out, nil
	}))

	// POST /keys/{key_id}/revoke
	huma.Register(grp, huma.Operation{
		OperationID:   "admin-revoke-key",
		Method:        http.MethodPost,
		Path:          "/keys/{key_id}/revoke",
		Summary:       "Revoke API key",
		Description:   "Revokes an API key, making it permanently unusable for authentication.",
		DefaultStatus: http.StatusNoContent,
		Tags:          []string{"API Keys"},
	}, exceptions.Handle(func(ctx context.Context, in *RevokeKeyRequest) (*struct{}, error) {
		if err := s.Keys.RevokeKey(ctx, in.KeyID); err != nil {
			return nil, huma.Error404NotFound("not found")
		}
		return &struct{}{}, nil
	}))
}
