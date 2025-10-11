table "organisation_invites" {
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
  column "email" {
    null = false
    type = sql("citext")
  }
  column "role_id" {
    null = true
    type = uuid
  }
  column "invited_by" {
    null = false
    type = uuid
  }
  column "token" {
    null = false
    type = text
  }
  column "accepted" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "organisation_invites_invited_by_fkey" {
    columns     = [column.invited_by]
    ref_columns = [table.users.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "organisation_invites_org_id_fkey" {
    columns     = [column.org_id]
    ref_columns = [table.organisations.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "organisation_invites_role_id_fkey" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_update   = NO_ACTION
    on_delete   = NO_ACTION
  }
  index "uniq_pending_invite_per_email_per_org" {
    unique  = true
    columns = [column.org_id, column.email]
    where   = "(accepted = false)"
  }
  unique "organisation_invites_token_key" {
    columns = [column.token]
  }
}