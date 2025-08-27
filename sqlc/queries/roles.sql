-- name: FindRoleByOrgAndName :one
SELECT 
  id,
  org_id,
  name,
  description,
  created_at
FROM roles
WHERE org_id = $1 AND name = $2
LIMIT 1;

-- name: CreateRole :one
INSERT INTO roles (org_id, name, description)
VALUES ($1, $2, $3)
RETURNING id, org_id, name, description, created_at;
