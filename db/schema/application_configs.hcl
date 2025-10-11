table "application_configs" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("uuid_generate_v4()")
  }
  column "app_id" {
    null = false
    type = uuid
  }
  column "org_id" {
    null    = false
    type    = uuid
    default = sql("app_current_org()")
  }
  column "environment" {
    null = false
    type = text
  }
  column "config" {
    null = false
    type = jsonb
    default = sql("'{}'::jsonb")
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
  foreign_key "application_configs_app_id_fkey" {
    columns     = [column.app_id]
    ref_columns = [table.applications.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "application_configs_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_app_configs_app_env" {
    columns = [column.app_id, column.environment]
    unique  = true
  }
  index "idx_app_configs_org" {
    columns = [column.org_id]
  }
}
