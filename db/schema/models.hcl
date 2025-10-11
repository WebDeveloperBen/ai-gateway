table "models" {
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
  column "provider" {
    null = false
    type = text
  }
  column "model_name" {
    null = false
    type = text
  }
  column "deployment_name" {
    null = true
    type = text
  }
  column "endpoint_url" {
    null = false
    type = text
  }
  column "auth_type" {
    null = false
    type = text
  }
  column "auth_config" {
    null = false
    type = jsonb
    default = sql("'{}'::jsonb")
  }
  column "metadata" {
    null = false
    type = jsonb
    default = sql("'{}'::jsonb")
  }
  column "enabled" {
    null    = false
    type    = boolean
    default = true
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "models_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_models_org_enabled" {
    columns = [column.org_id, column.enabled]
  }
  index "idx_models_provider_model" {
    columns = [column.org_id, column.provider, column.model_name]
  }
}
