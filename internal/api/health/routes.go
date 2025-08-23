package health

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterPublicRoutes(api huma.API) {
	// Register GET /greeting/{name}
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Health check for the service",
		Description: "Health check endpoint for the api service",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, input *struct{}) (*GetResponse, error) {
		return &GetResponse{
			Body: Health{
				Message: "All systems online and healthy",
				Status:  http.StatusOK,
			},
		}, nil
	})
}
