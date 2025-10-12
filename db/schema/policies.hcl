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
  foreign_key "policies_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_policies_org" {
    columns = [column.org_id]
  }
  index "idx_policies_type" {
    columns = [column.policy_type]
  }
}

table "policy_applications" {
  schema = schema.public
  column "policy_id" {
    null = false
    type = uuid
  }
  column "app_id" {
    null = false
    type = uuid
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.policy_id, column.app_id]
  }
  foreign_key "policy_applications_policy_id_fkey" {
    columns     = [column.policy_id]
    ref_columns = [table.policies.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "policy_applications_app_id_fkey" {
    columns     = [column.app_id]
    ref_columns = [table.applications.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_policy_applications_app" {
    columns = [column.app_id]
  }
  index "idx_policy_applications_policy" {
    columns = [column.policy_id]
  }
}
