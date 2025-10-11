table "user_roles" {
  schema = schema.public
  column "user_id" {
    null = false
    type = uuid
  }
  column "role_id" {
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
    columns = [column.user_id, column.role_id, column.org_id]
  }
  foreign_key "user_roles_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "user_roles_role_id_fkey" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "user_roles_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_user_roles_org_role" {
    columns = [column.org_id, column.role_id]
  }
  index "idx_user_roles_org_user" {
    columns = [column.org_id, column.user_id]
  }
}