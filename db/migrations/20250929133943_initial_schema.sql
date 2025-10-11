-- +goose Up
-- Create "organisations" table
CREATE TABLE "public"."organisations" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);

-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "sub" text NULL,
  "email" citext NOT NULL,
  "name" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_org_id_email_key" UNIQUE ("org_id", "email"),
  CONSTRAINT "users_sub_key" UNIQUE ("sub"),
  CONSTRAINT "users_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_users_org_email" ON "public"."users" ("org_id", "email");
CREATE INDEX "idx_users_sub" ON "public"."users" ("sub");

-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "roles_name_key" UNIQUE ("name")
);

-- Create "api_keys" table
CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "user_id" uuid NOT NULL,
  "key_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "last_used" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "api_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "api_keys_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_api_keys_org_user" ON "public"."api_keys" ("org_id", "user_id");

-- Create "organisation_keys" table
CREATE TABLE "public"."organisation_keys" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "key_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "last_used" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "organisation_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_org_keys_org" ON "public"."organisation_keys" ("org_id");

-- Create "organisation_roles" table
CREATE TABLE "public"."organisation_roles" (
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "role_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("org_id", "role_id"),
  CONSTRAINT "organisation_roles_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- Create "organisation_users" table
CREATE TABLE "public"."organisation_users" (
  "user_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "org_id"),
  CONSTRAINT "organisation_users_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_org_users_org" ON "public"."organisation_users" ("org_id");
CREATE INDEX "idx_org_users_user" ON "public"."organisation_users" ("user_id");

-- Create "organisation_invites" table
CREATE TABLE "public"."organisation_invites" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "email" citext NOT NULL,
  "role_id" uuid NULL,
  "invited_by" uuid NOT NULL,
  "token" text NOT NULL,
  "accepted" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "expires_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "organisation_invites_token_key" UNIQUE ("token"),
  CONSTRAINT "organisation_invites_invited_by_fkey" FOREIGN KEY ("invited_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_invites_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_invites_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
CREATE UNIQUE INDEX "uniq_pending_invite_per_email_per_org" ON "public"."organisation_invites" ("org_id", "email") WHERE (accepted = false);

-- Create "user_roles" table
CREATE TABLE "public"."user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "role_id", "org_id"),
  CONSTRAINT "user_roles_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
CREATE INDEX "idx_user_roles_org_role" ON "public"."user_roles" ("org_id", "role_id");
CREATE INDEX "idx_user_roles_org_user" ON "public"."user_roles" ("org_id", "user_id");

-- +goose Down
DROP TABLE IF EXISTS "user_roles";
DROP TABLE IF EXISTS "organisation_invites";
DROP TABLE IF EXISTS "organisation_users";
DROP TABLE IF EXISTS "organisation_roles";
DROP TABLE IF EXISTS "organisation_keys";
DROP TABLE IF EXISTS "api_keys";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "organisations";
