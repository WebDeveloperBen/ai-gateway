-- name: FindRoleByName :one
SELECT 
  id,
  name,
  description,
  created_at
FROM roles
WHERE name = $1
LIMIT 1;

-- name: FindRoleByID :one
SELECT
  id,
  name,
  description,
  created_at
FROM roles
WHERE id = $1
LIMIT 1;

-- name: CreateRole :one
INSERT INTO roles (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: ListRoles :many
SELECT 
  id,
  name,
  description,
  created_at
FROM roles;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3
WHERE id = $1
RETURNING id, name, description, created_at;
