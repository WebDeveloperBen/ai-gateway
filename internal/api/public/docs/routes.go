// Package docs implements a scalar api docs interface
package docs

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes creates a Scalar documentation route
func RegisterRoutes(router http.Handler, api huma.API) {
	r, ok := router.(interface {
		Get(pattern string, handlerFn http.HandlerFunc)
	})
	if !ok {
		panic("router must support GET method")
	}

	// Custom OpenAPI endpoint with tag grouping
	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		openAPI := api.OpenAPI()

		// Add x-tagGroups extension for hierarchical organization
		// This is supported by Scalar v1.19.3+
		if openAPI.Extensions == nil {
			openAPI.Extensions = make(map[string]any)
		}

		// Collect all existing tags from the spec
		allTags := make(map[string]bool)
		for _, pathItem := range openAPI.Paths {
			ops := []*huma.Operation{pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Patch, pathItem.Delete}
			for _, op := range ops {
				if op != nil {
					for _, tag := range op.Tags {
						allTags[tag] = true
					}
				}
			}
		}

		// Automatically organize tags into groups based on the API path structure
		// Convention: /api/providers/* → Providers, /api/v1/admin/* → Admin, everything else → Public
		tagToGroup := make(map[string]string)

		for path, pathItem := range openAPI.Paths {
			ops := []*huma.Operation{pathItem.Get, pathItem.Post, pathItem.Put, pathItem.Patch, pathItem.Delete}
			for _, op := range ops {
				if op != nil && len(op.Tags) > 0 {
					tag := op.Tags[0]

					// Determine group based on path prefix convention
					if strings.HasPrefix(path, "/api/providers/") {
						tagToGroup[tag] = "Providers"
					} else if strings.HasPrefix(path, "/api/v1/admin/") {
						tagToGroup[tag] = "Admin"
					} else {
						tagToGroup[tag] = "Public"
					}
				}
			}
		}

		// Group tags by their assigned group
		groups := make(map[string][]string)
		for tag, group := range tagToGroup {
			groups[group] = append(groups[group], tag)
		}

		// Build x-tagGroups in a consistent order
		tagGroups := []map[string]any{}
		for _, groupName := range []string{"Providers", "Admin", "Public"} {
			if tags, exists := groups[groupName]; exists && len(tags) > 0 {
				tagGroups = append(tagGroups, map[string]any{
					"name": groupName,
					"tags": tags,
				})
			}
		}

		openAPI.Extensions["x-tagGroups"] = tagGroups

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(openAPI)
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
	<head>
		<title>API Reference</title>
		<meta charset="utf-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1" />
	</head>
	<body>
		<script
			id="api-reference"
			data-url="/openapi.json">
		</script>

		<script>
			var configuration = {
				theme: "deep-space",
				tagsSorter: "alpha",
				operationsSorter: "alpha"
			};
			document.getElementById('api-reference').dataset.configuration = JSON.stringify(configuration);
		</script>

		<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
	</body>
</html>`))
	})
}
