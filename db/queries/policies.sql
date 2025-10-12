-- name: GetPolicy :one
SELECT * FROM policies
WHERE id = $1 LIMIT 1;

-- name: ListPolicies :many
SELECT p.* FROM policies p
JOIN policy_applications pa ON p.id = pa.policy_id
WHERE pa.app_id = $1
ORDER BY p.policy_type
LIMIT $2 OFFSET $3;

-- name: ListEnabledPolicies :many
SELECT p.* FROM policies p
JOIN policy_applications pa ON p.id = pa.policy_id
WHERE pa.app_id = $1 AND p.enabled = true
ORDER BY p.policy_type
LIMIT $2 OFFSET $3;

-- name: GetPoliciesByType :many
SELECT p.* FROM policies p
JOIN policy_applications pa ON p.id = pa.policy_id
WHERE pa.app_id = $1 AND p.policy_type = $2
ORDER BY p.created_at;

-- name: CreatePolicy :one
INSERT INTO policies (
  org_id, policy_type, config, enabled
) VALUES (
  $1, $2, $3, $4
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

-- name: AttachPolicyToApp :exec
INSERT INTO policy_applications (policy_id, app_id)
VALUES ($1, $2)
ON CONFLICT (policy_id, app_id) DO NOTHING;

-- name: DetachPolicyFromApp :exec
DELETE FROM policy_applications
WHERE policy_id = $1 AND app_id = $2;

-- name: GetAppsForPolicy :many
SELECT a.* FROM applications a
JOIN policy_applications pa ON a.id = pa.app_id
WHERE pa.policy_id = $1
ORDER BY a.name;

-- name: GetPoliciesForApp :many
SELECT p.* FROM policies p
JOIN policy_applications pa ON p.id = pa.policy_id
WHERE pa.app_id = $1
ORDER BY p.policy_type;
