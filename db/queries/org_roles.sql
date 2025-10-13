-- name: AssignRoleToOrg :exec
INSERT INTO organisation_roles (org_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
RETURNING org_id, role_id, created_at;
