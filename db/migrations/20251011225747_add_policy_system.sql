-- +goose Up
-- +goose StatementBegin
CREATE TABLE "public"."applications" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "name" text NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "applications_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE UNIQUE INDEX "idx_applications_name_org" ON "public"."applications" ("org_id", "name");
CREATE INDEX "idx_applications_org" ON "public"."applications" ("org_id");

-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "public"."application_configs" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "app_id" uuid NOT NULL,
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "environment" text NOT NULL,
  "config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "application_configs_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "application_configs_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE UNIQUE INDEX "idx_app_configs_app_env" ON "public"."application_configs" ("app_id", "environment");
CREATE INDEX "idx_app_configs_org" ON "public"."application_configs" ("org_id");

-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "public"."models" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "provider" text NOT NULL,
  "model_name" text NOT NULL,
  "deployment_name" text NULL,
  "endpoint_url" text NOT NULL,
  "auth_type" text NOT NULL,
  "auth_config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "metadata" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "enabled" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "models_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_models_org_enabled" ON "public"."models" ("org_id", "enabled");
CREATE INDEX "idx_models_provider_model" ON "public"."models" ("org_id", "provider", "model_name");

-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "public"."policies" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
  "app_id" uuid NOT NULL,
  "policy_type" text NOT NULL,
  "config" jsonb NOT NULL DEFAULT '{}'::jsonb,
  "enabled" boolean NOT NULL DEFAULT true,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "policies_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "policies_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX "idx_policies_app_enabled" ON "public"."policies" ("app_id", "enabled");
CREATE INDEX "idx_policies_org" ON "public"."policies" ("org_id");
CREATE INDEX "idx_policies_type" ON "public"."policies" ("policy_type");

-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE "public"."usage_metrics" (
  "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
  "org_id" uuid NOT NULL DEFAULT app_current_org(),
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
  CONSTRAINT "usage_metrics_org_id_fkey" FOREIGN KEY ("org_id") REFERENCES "public"."organisations" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "usage_metrics_app_id_fkey" FOREIGN KEY ("app_id") REFERENCES "public"."applications" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "usage_metrics_api_key_id_fkey" FOREIGN KEY ("api_key_id") REFERENCES "public"."api_keys" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "usage_metrics_model_id_fkey" FOREIGN KEY ("model_id") REFERENCES "public"."models" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);

CREATE INDEX "idx_usage_metrics_org_timestamp" ON "public"."usage_metrics" ("org_id", "timestamp");
CREATE INDEX "idx_usage_metrics_app_timestamp" ON "public"."usage_metrics" ("app_id", "timestamp");
CREATE INDEX "idx_usage_metrics_timestamp" ON "public"."usage_metrics" ("timestamp");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."usage_metrics";
DROP TABLE IF EXISTS "public"."policies";
DROP TABLE IF EXISTS "public"."models";
DROP TABLE IF EXISTS "public"."application_configs";
DROP TABLE IF EXISTS "public"."applications";
-- +goose StatementEnd
