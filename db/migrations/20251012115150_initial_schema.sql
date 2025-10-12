-- +goose Up
-- create "organisations" table
CREATE TABLE "public"."organisations" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- create "applications" table
CREATE TABLE "public"."applications" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "applications_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_applications_name_org" to table: "applications"
CREATE UNIQUE INDEX "idx_applications_name_org" ON "public"."applications" ("org_id", "name");
-- create index "idx_applications_org" to table: "applications"
CREATE INDEX "idx_applications_org" ON "public"."applications" ("org_id");
-- create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "sub" text NULL,
  "email" public.citext NOT NULL,
  "name" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_org_id_email_key" UNIQUE ("org_id", "email"),
  CONSTRAINT "users_sub_key" UNIQUE ("sub"),
  CONSTRAINT "users_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_users_org_email" to table: "users"
CREATE INDEX "idx_users_org_email" ON "public"."users" ("org_id", "email");
-- create index "idx_users_sub" to table: "users"
CREATE INDEX "idx_users_sub" ON "public"."users" ("sub");
-- create "api_keys" table
CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "app_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "key_prefix" text NOT NULL,
  "secret_phc" text NOT NULL,
  "status" text NOT NULL DEFAULT 'active',
  "last_four" text NOT NULL,
  "expires_at" timestamptz NULL,
  "last_used_at" timestamptz NULL,
  "metadata" jsonb NOT NULL DEFAULT '{}',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "api_keys_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "api_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "api_keys_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_api_keys_org_app" to table: "api_keys"
CREATE INDEX "idx_api_keys_org_app" ON "public"."api_keys" ("org_id", "app_id");
-- create index "idx_api_keys_prefix" to table: "api_keys"
CREATE INDEX "idx_api_keys_prefix" ON "public"."api_keys" ("key_prefix");
-- create "application_configs" table
CREATE TABLE "public"."application_configs" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "app_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "environment" text NOT NULL,
  "config" jsonb NOT NULL DEFAULT '{}',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "application_configs_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "application_configs_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_app_configs_app_env" to table: "application_configs"
CREATE UNIQUE INDEX "idx_app_configs_app_env" ON "public"."application_configs" ("app_id", "environment");
-- create index "idx_app_configs_org" to table: "application_configs"
CREATE INDEX "idx_app_configs_org" ON "public"."application_configs" ("org_id");
-- create "models" table
CREATE TABLE "public"."models" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "provider" text NOT NULL,
  "model_name" text NOT NULL,
  "deployment_name" text NULL,
  "endpoint_url" text NOT NULL,
  "auth_type" text NOT NULL,
  "auth_config" jsonb NOT NULL DEFAULT '{}',
  "metadata" jsonb NOT NULL DEFAULT '{}',
  "enabled" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "models_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_models_org_enabled" to table: "models"
CREATE INDEX "idx_models_org_enabled" ON "public"."models" ("org_id", "enabled");
-- create index "idx_models_provider_model" to table: "models"
CREATE INDEX "idx_models_provider_model" ON "public"."models" ("org_id", "provider", "model_name");
-- create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "roles_name_key" UNIQUE ("name")
);
-- create "organisation_invites" table
CREATE TABLE "public"."organisation_invites" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "email" public.citext NOT NULL,
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
-- create index "uniq_pending_invite_per_email_per_org" to table: "organisation_invites"
CREATE UNIQUE INDEX "uniq_pending_invite_per_email_per_org" ON "public"."organisation_invites" ("org_id", "email") WHERE (accepted = false);
-- create "organisation_keys" table
CREATE TABLE "public"."organisation_keys" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "key_hash" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "last_used" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "organisation_keys_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_org_keys_org" to table: "organisation_keys"
CREATE INDEX "idx_org_keys_org" ON "public"."organisation_keys" ("org_id");
-- create "organisation_roles" table
CREATE TABLE "public"."organisation_roles" (
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "role_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("org_id", "role_id"),
  CONSTRAINT "organisation_roles_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "organisation_users" table
CREATE TABLE "public"."organisation_users" (
  "user_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "org_id"),
  CONSTRAINT "organisation_users_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "organisation_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_org_users_org" to table: "organisation_users"
CREATE INDEX "idx_org_users_org" ON "public"."organisation_users" ("org_id");
-- create index "idx_org_users_user" to table: "organisation_users"
CREATE INDEX "idx_org_users_user" ON "public"."organisation_users" ("user_id");
-- create "policies" table
CREATE TABLE "public"."policies" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "policy_type" text NOT NULL,
  "config" jsonb NOT NULL DEFAULT '{}',
  "enabled" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "policies_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_policies_org" to table: "policies"
CREATE INDEX "idx_policies_org" ON "public"."policies" ("org_id");
-- create index "idx_policies_type" to table: "policies"
CREATE INDEX "idx_policies_type" ON "public"."policies" ("policy_type");
-- create "policy_applications" table
CREATE TABLE "public"."policy_applications" (
  "policy_id" uuid NOT NULL,
  "app_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("policy_id", "app_id"),
  CONSTRAINT "policy_applications_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "policy_applications_policy_id_fkey" FOREIGN KEY ("policy_id") REFERENCES "public"."policies" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_policy_applications_app" to table: "policy_applications"
CREATE INDEX "idx_policy_applications_app" ON "public"."policy_applications" ("app_id");
-- create index "idx_policy_applications_policy" to table: "policy_applications"
CREATE INDEX "idx_policy_applications_policy" ON "public"."policy_applications" ("policy_id");
-- create "usage_metrics" table
CREATE TABLE "public"."usage_metrics" (
  "id" uuid NOT NULL DEFAULT public.uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "app_id" uuid NOT NULL,
  "api_key_id" uuid NOT NULL,
  "model_id" uuid NULL,
  "provider" text NOT NULL,
  "model_name" text NOT NULL,
  "prompt_tokens" integer NOT NULL DEFAULT 0,
  "completion_tokens" integer NOT NULL DEFAULT 0,
  "total_tokens" integer NOT NULL DEFAULT 0,
  "request_size_bytes" integer NOT NULL DEFAULT 0,
  "response_size_bytes" integer NOT NULL DEFAULT 0,
  "timestamp" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "usage_metrics_api_key_id_fkey" FOREIGN KEY ("api_key_id") REFERENCES "public"."api_keys" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "usage_metrics_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "usage_metrics_model_id_fkey" FOREIGN KEY ("model_id") REFERENCES "public"."models" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "usage_metrics_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_usage_metrics_app_timestamp" to table: "usage_metrics"
CREATE INDEX "idx_usage_metrics_app_timestamp" ON "public"."usage_metrics" ("app_id", "timestamp");
-- create index "idx_usage_metrics_org_timestamp" to table: "usage_metrics"
CREATE INDEX "idx_usage_metrics_org_timestamp" ON "public"."usage_metrics" ("org_id", "timestamp");
-- create index "idx_usage_metrics_timestamp" to table: "usage_metrics"
CREATE INDEX "idx_usage_metrics_timestamp" ON "public"."usage_metrics" ("timestamp");
-- create "user_roles" table
CREATE TABLE "public"."user_roles" (
  "user_id" uuid NOT NULL,
  "role_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT public.app_current_org(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id", "role_id", "org_id"),
  CONSTRAINT "user_roles_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_roles_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_user_roles_org_role" to table: "user_roles"
CREATE INDEX "idx_user_roles_org_role" ON "public"."user_roles" ("org_id", "role_id");
-- create index "idx_user_roles_org_user" to table: "user_roles"
CREATE INDEX "idx_user_roles_org_user" ON "public"."user_roles" ("org_id", "user_id");

-- +goose Down
-- reverse: create index "idx_user_roles_org_user" to table: "user_roles"
DROP INDEX "public"."idx_user_roles_org_user";
-- reverse: create index "idx_user_roles_org_role" to table: "user_roles"
DROP INDEX "public"."idx_user_roles_org_role";
-- reverse: create "user_roles" table
DROP TABLE "public"."user_roles";
-- reverse: create index "idx_usage_metrics_timestamp" to table: "usage_metrics"
DROP INDEX "public"."idx_usage_metrics_timestamp";
-- reverse: create index "idx_usage_metrics_org_timestamp" to table: "usage_metrics"
DROP INDEX "public"."idx_usage_metrics_org_timestamp";
-- reverse: create index "idx_usage_metrics_app_timestamp" to table: "usage_metrics"
DROP INDEX "public"."idx_usage_metrics_app_timestamp";
-- reverse: create "usage_metrics" table
DROP TABLE "public"."usage_metrics";
-- reverse: create index "idx_policy_applications_policy" to table: "policy_applications"
DROP INDEX "public"."idx_policy_applications_policy";
-- reverse: create index "idx_policy_applications_app" to table: "policy_applications"
DROP INDEX "public"."idx_policy_applications_app";
-- reverse: create "policy_applications" table
DROP TABLE "public"."policy_applications";
-- reverse: create index "idx_policies_type" to table: "policies"
DROP INDEX "public"."idx_policies_type";
-- reverse: create index "idx_policies_org" to table: "policies"
DROP INDEX "public"."idx_policies_org";
-- reverse: create "policies" table
DROP TABLE "public"."policies";
-- reverse: create index "idx_org_users_user" to table: "organisation_users"
DROP INDEX "public"."idx_org_users_user";
-- reverse: create index "idx_org_users_org" to table: "organisation_users"
DROP INDEX "public"."idx_org_users_org";
-- reverse: create "organisation_users" table
DROP TABLE "public"."organisation_users";
-- reverse: create "organisation_roles" table
DROP TABLE "public"."organisation_roles";
-- reverse: create index "idx_org_keys_org" to table: "organisation_keys"
DROP INDEX "public"."idx_org_keys_org";
-- reverse: create "organisation_keys" table
DROP TABLE "public"."organisation_keys";
-- reverse: create index "uniq_pending_invite_per_email_per_org" to table: "organisation_invites"
DROP INDEX "public"."uniq_pending_invite_per_email_per_org";
-- reverse: create "organisation_invites" table
DROP TABLE "public"."organisation_invites";
-- reverse: create "roles" table
DROP TABLE "public"."roles";
-- reverse: create index "idx_models_provider_model" to table: "models"
DROP INDEX "public"."idx_models_provider_model";
-- reverse: create index "idx_models_org_enabled" to table: "models"
DROP INDEX "public"."idx_models_org_enabled";
-- reverse: create "models" table
DROP TABLE "public"."models";
-- reverse: create index "idx_app_configs_org" to table: "application_configs"
DROP INDEX "public"."idx_app_configs_org";
-- reverse: create index "idx_app_configs_app_env" to table: "application_configs"
DROP INDEX "public"."idx_app_configs_app_env";
-- reverse: create "application_configs" table
DROP TABLE "public"."application_configs";
-- reverse: create index "idx_api_keys_prefix" to table: "api_keys"
DROP INDEX "public"."idx_api_keys_prefix";
-- reverse: create index "idx_api_keys_org_app" to table: "api_keys"
DROP INDEX "public"."idx_api_keys_org_app";
-- reverse: create "api_keys" table
DROP TABLE "public"."api_keys";
-- reverse: create index "idx_users_sub" to table: "users"
DROP INDEX "public"."idx_users_sub";
-- reverse: create index "idx_users_org_email" to table: "users"
DROP INDEX "public"."idx_users_org_email";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_applications_org" to table: "applications"
DROP INDEX "public"."idx_applications_org";
-- reverse: create index "idx_applications_name_org" to table: "applications"
DROP INDEX "public"."idx_applications_name_org";
-- reverse: create "applications" table
DROP TABLE "public"."applications";
-- reverse: create "organisations" table
DROP TABLE "public"."organisations";
