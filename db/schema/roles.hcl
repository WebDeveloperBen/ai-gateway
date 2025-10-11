table "roles" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("uuid_generate_v4()")
  }
  column "name" {
    null = false
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
  unique "roles_name_key" {
    columns = [column.name]
  }
}