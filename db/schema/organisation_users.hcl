table "organisation_users" {
  schema = schema.public
  column "user_id" {
    null = false
    type = uuid
  }
  column "org_id" {
    null    = false
    type    = uuid
    default = sql("app_current_org()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.user_id, column.org_id]
  }
  foreign_key "organisation_users_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "organisation_users_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_org_users_org" {
    columns = [column.org_id]
  }
  index "idx_org_users_user" {
    columns = [column.user_id]
  }
}