-- name: GetModel :one
SELECT * FROM models
WHERE id = $1 LIMIT 1;

-- name: GetModelByProviderAndName :one
SELECT * FROM models
WHERE org_id = $1 AND provider = $2 AND model_name = $3 LIMIT 1;

-- name: ListModels :many
SELECT * FROM models
WHERE org_id = $1
ORDER BY provider, model_name
LIMIT $2 OFFSET $3;

-- name: ListEnabledModels :many
SELECT * FROM models
WHERE org_id = $1 AND enabled = true
ORDER BY provider, model_name
LIMIT $2 OFFSET $3;

-- name: CreateModel :one
INSERT INTO models (
  org_id, provider, model_name, deployment_name, endpoint_url,
  auth_type, auth_config, metadata, enabled
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateModel :one
UPDATE models
SET provider = $2,
    model_name = $3,
    deployment_name = $4,
    endpoint_url = $5,
    auth_type = $6,
    auth_config = $7,
    metadata = $8,
    enabled = $9,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteModel :exec
DELETE FROM models
WHERE id = $1;

-- name: EnableModel :exec
UPDATE models
SET enabled = true,
    updated_at = now()
WHERE id = $1;

-- name: DisableModel :exec
UPDATE models
SET enabled = false,
    updated_at = now()
WHERE id = $1;
