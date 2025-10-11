table "usage_metrics" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("uuid_generate_v4()")
  }
  column "org_id" {
    null    = false
    type    = uuid
    default = sql("app_current_org()")
  }
  column "app_id" {
    null = false
    type = uuid
  }
  column "api_key_id" {
    null = false
    type = uuid
  }
  column "model_id" {
    null = true
    type = uuid
  }
  column "provider" {
    null = false
    type = text
  }
  column "model_name" {
    null = false
    type = text
  }
  column "prompt_tokens" {
    null = false
    type = integer
    default = 0
  }
  column "completion_tokens" {
    null = false
    type = integer
    default = 0
  }
  column "total_tokens" {
    null = false
    type = integer
    default = 0
  }
  column "request_size_bytes" {
    null = false
    type = integer
    default = 0
  }
  column "response_size_bytes" {
    null = false
    type = integer
    default = 0
  }
  column "timestamp" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "usage_metrics_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "usage_metrics_app_id_fkey" {
    columns     = [column.app_id]
    ref_columns = [table.applications.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "usage_metrics_api_key_id_fkey" {
    columns     = [column.api_key_id]
    ref_columns = [table.api_keys.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "usage_metrics_model_id_fkey" {
    columns     = [column.model_id]
    ref_columns = [table.models.column.id]
    on_update   = NO_ACTION
    on_delete   = SET_NULL
  }
  index "idx_usage_metrics_org_timestamp" {
    columns = [column.org_id, column.timestamp]
  }
  index "idx_usage_metrics_app_timestamp" {
    columns = [column.app_id, column.timestamp]
  }
  index "idx_usage_metrics_timestamp" {
    columns = [column.timestamp]
  }
}
