-- name: CreateAPIKey :one
INSERT INTO api_keys (org_id, user_id, key_hash)
VALUES ($1, $2, $3)
RETURNING id, org_id, user_id, key_hash, created_at, last_used;

-- name: GetAPIKey :one
SELECT id, org_id, user_id, key_hash, created_at, last_used
FROM api_keys
WHERE id = $1;

-- name: GetAPIKeyByHash :one
SELECT id, org_id, user_id, key_hash, created_at, last_used
FROM api_keys
WHERE key_hash = $1;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used = now()
WHERE id = $1;