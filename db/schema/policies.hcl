table "policies" {
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
  column "policy_type" {
    null = false
    type = text
  }
  column "config" {
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
  foreign_key "policies_app_id_fkey" {
    columns     = [column.app_id]
    ref_columns = [table.applications.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "policies_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_policies_app_enabled" {
    columns = [column.app_id, column.enabled]
  }
  index "idx_policies_org" {
    columns = [column.org_id]
  }
  index "idx_policies_type" {
    columns = [column.policy_type]
  }
}
