-- name: GetApplicationConfig :one
SELECT * FROM application_configs
WHERE id = $1 LIMIT 1;

-- name: GetApplicationConfigByEnv :one
SELECT * FROM application_configs
WHERE app_id = $1 AND environment = $2 LIMIT 1;

-- name: ListApplicationConfigs :many
SELECT * FROM application_configs
WHERE app_id = $1
ORDER BY environment;

-- name: CreateApplicationConfig :one
INSERT INTO application_configs (
  app_id, org_id, environment, config
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateApplicationConfig :one
UPDATE application_configs
SET config = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteApplicationConfig :exec
DELETE FROM application_configs
WHERE id = $1;
