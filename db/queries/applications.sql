-- name: GetApplication :one
SELECT * FROM applications
WHERE id = $1 LIMIT 1;

-- name: GetApplicationByName :one
SELECT * FROM applications
WHERE org_id = $1 AND name = $2 LIMIT 1;

-- name: ListApplications :many
SELECT * FROM applications
WHERE org_id = $1
ORDER BY name;

-- name: CreateApplication :one
INSERT INTO applications (
  org_id, name, description
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateApplication :one
UPDATE applications
SET name = $2,
    description = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = $1;
