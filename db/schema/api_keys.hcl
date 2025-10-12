table "api_keys" {
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
  column "user_id" {
    null = false
    type = uuid
  }
  column "key_prefix" {
    null = false
    type = text
  }
  column "secret_phc" {
    null = false
    type = text
  }
  column "status" {
    null    = false
    type    = text
    default = "active"
  }
  column "last_four" {
    null = false
    type = text
  }
  column "expires_at" {
    null = true
    type = timestamptz
  }
  column "last_used_at" {
    null = true
    type = timestamptz
  }
  column "metadata" {
    null    = false
    type    = jsonb
    default = sql("'{}'::jsonb")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "api_keys_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "api_keys_app_id_fkey" {
    columns     = [column.app_id]
    ref_columns = [table.applications.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "api_keys_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_api_keys_org_app" {
    columns = [column.org_id, column.app_id]
  }
  index "idx_api_keys_prefix" {
    columns = [column.key_prefix]
  }
}