-- name: FindUserBySubOrEmail :one
SELECT * FROM users
WHERE sub = $1 OR email = $2
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (org_id, sub, email, name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id, org_id)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;
