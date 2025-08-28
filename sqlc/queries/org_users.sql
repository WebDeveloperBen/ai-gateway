-- name: EnsureOrgMembership :exec
INSERT INTO organisation_users (org_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;
