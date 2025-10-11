table "users" {
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
  column "sub" {
    null = true
    type = text
  }
  column "email" {
    null = false
    type = sql("citext")
  }
  column "name" {
    null = true
    type = text
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
  foreign_key "users_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_users_org_email" {
    columns = [column.org_id, column.email]
  }
  index "idx_users_sub" {
    columns = [column.sub]
  }
  unique "users_org_id_email_key" {
    columns = [column.org_id, column.email]
  }
  unique "users_sub_key" {
    columns = [column.sub]
  }
}