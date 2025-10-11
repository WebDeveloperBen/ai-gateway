table "organisation_keys" {
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
  column "key_hash" {
    null = false
    type = text
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "last_used" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "organisation_keys_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_org_keys_org" {
    columns = [column.org_id]
  }
}