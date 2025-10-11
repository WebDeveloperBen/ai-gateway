table "organisation_roles" {
  schema = schema.public
  column "org_id" {
    null    = false
    type    = uuid
    default = sql("app_current_org()")
  }
  column "role_id" {
    null = false
    type = uuid
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.org_id, column.role_id]
  }
  foreign_key "organisation_roles_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "organisation_roles_role_id_fkey" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
}