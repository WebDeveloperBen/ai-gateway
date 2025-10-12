-- name: GetPolicy :one
SELECT * FROM policies
WHERE id = $1 LIMIT 1;

-- name: ListPolicies :many
SELECT * FROM policies
WHERE app_id = $1
ORDER BY policy_type
LIMIT $2 OFFSET $3;

-- name: ListEnabledPolicies :many
SELECT * FROM policies
WHERE app_id = $1 AND enabled = true
ORDER BY policy_type
LIMIT $2 OFFSET $3;

-- name: GetPoliciesByType :many
SELECT * FROM policies
WHERE app_id = $1 AND policy_type = $2
ORDER BY created_at;

-- name: CreatePolicy :one
INSERT INTO policies (
  org_id, app_id, policy_type, config, enabled
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdatePolicy :one
UPDATE policies
SET policy_type = $2,
    config = $3,
    enabled = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePolicy :exec
DELETE FROM policies
WHERE id = $1;

-- name: EnablePolicy :exec
UPDATE policies
SET enabled = true,
    updated_at = now()
WHERE id = $1;

-- name: DisablePolicy :exec
UPDATE policies
SET enabled = false,
    updated_at = now()
WHERE id = $1;
