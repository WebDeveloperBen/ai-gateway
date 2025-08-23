// Package docs implements a scalar api docs interface
package docs

import "net/http"

// RegisterRoutes creates a Scalar documentation route
func RegisterRoutes(router http.Handler) {
	r, ok := router.(interface {
		Get(pattern string, handlerFn http.HandlerFunc)
	})
	if !ok {
		panic("router must support GET method")
	}
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
				theme: "deep-space"
			};
			document.getElementById('api-reference').dataset.configuration = JSON.stringify(configuration);
		</script>

		<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
	</body>
</html>`))
	})
}
