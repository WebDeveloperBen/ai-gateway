-- name: CreateOrg :one
INSERT INTO organisations (name)
VALUES ($1)
RETURNING *;

-- name: FindOrgByID :one
SELECT * FROM organisations
WHERE id = $1;
