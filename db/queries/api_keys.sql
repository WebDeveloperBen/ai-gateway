-- name: InsertAPIKey :one
INSERT INTO api_keys (
  org_id,
  app_id,
  user_id,
  key_prefix,
  secret_phc,
  status,
  last_four,
  expires_at,
  metadata
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetAPIKeyByID :one
SELECT * FROM api_keys
WHERE id = $1;

-- name: GetAPIKeyByPrefix :one
SELECT * FROM api_keys
WHERE key_prefix = $1;

-- name: GetSecretPHCByPrefix :one
SELECT secret_phc FROM api_keys
WHERE key_prefix = $1;

-- name: UpdateAPIKeyStatus :exec
UPDATE api_keys
SET status = $2
WHERE key_prefix = $1;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used_at = now()
WHERE key_prefix = $1;

-- name: ListAPIKeysByOrgID :many
SELECT * FROM api_keys
WHERE org_id = $1
ORDER BY created_at DESC;

-- name: ListAPIKeysByAppID :many
SELECT * FROM api_keys
WHERE app_id = $1
ORDER BY created_at DESC;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys
WHERE id = $1;
