-- name: CreateUsageMetric :one
INSERT INTO usage_metrics (
  org_id, app_id, api_key_id, model_id, provider, model_name,
  prompt_tokens, completion_tokens, total_tokens,
  request_size_bytes, response_size_bytes, timestamp
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: GetUsageMetricsByApp :many
SELECT * FROM usage_metrics
WHERE app_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp DESC
LIMIT $4 OFFSET $5;

-- name: GetUsageMetricsByOrg :many
SELECT * FROM usage_metrics
WHERE org_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp DESC;

-- name: GetUsageMetricsByAPIKey :many
SELECT * FROM usage_metrics
WHERE api_key_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp DESC;

-- name: SumTokensByApp :one
SELECT
  COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
  COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
  COALESCE(SUM(total_tokens), 0) as total_tokens,
  COUNT(*) as request_count
FROM usage_metrics
WHERE app_id = $1
  AND timestamp >= $2
  AND timestamp <= $3;

-- name: SumTokensByOrg :one
SELECT
  COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
  COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
  COALESCE(SUM(total_tokens), 0) as total_tokens,
  COUNT(*) as request_count
FROM usage_metrics
WHERE org_id = $1
  AND timestamp >= $2
  AND timestamp <= $3;

-- name: GetUsageByModel :many
SELECT
  model_name,
  provider,
  COALESCE(SUM(prompt_tokens), 0) as total_prompt_tokens,
  COALESCE(SUM(completion_tokens), 0) as total_completion_tokens,
  COALESCE(SUM(total_tokens), 0) as total_tokens,
  COUNT(*) as request_count
FROM usage_metrics
WHERE app_id = $1
  AND timestamp >= $2
  AND timestamp <= $3
GROUP BY model_name, provider
ORDER BY total_tokens DESC
LIMIT $4 OFFSET $5;
