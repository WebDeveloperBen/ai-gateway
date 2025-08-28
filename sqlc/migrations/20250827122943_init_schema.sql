-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

-- =====================
-- Postgres Helper Functions
-- =====================

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app_current_org() RETURNS uuid
LANGUAGE sql STABLE AS $$
  SELECT COALESCE(
    NULLIF(current_setting('app.current_org', true), ''),
    '00000000-0000-0000-0000-000000000000'
  )::uuid;
$$;
-- +goose StatementEnd

-- =====================
-- Organisations
-- =====================
CREATE TABLE organisations (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now()
);
CREATE TRIGGER organisations_set_updated_at
BEFORE UPDATE ON organisations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

ALTER TABLE organisations ENABLE ROW LEVEL SECURITY;

CREATE POLICY org_isolation_organisations ON organisations
  USING (id = app_current_org());

-- =====================
-- Global Roles (definitions reused by all orgs)
-- =====================
CREATE TABLE roles (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text NOT NULL UNIQUE,
    description text,
    created_at  timestamptz NOT NULL DEFAULT now()
);

-- =====================
-- Organisation ⇄ Role assignments
-- =====================
CREATE TABLE organisation_roles (
    org_id      uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    role_id     uuid NOT NULL REFERENCES roles(id),
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (org_id, role_id)
);
ALTER TABLE organisation_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE organisation_roles
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_organisation_roles ON organisation_roles
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

-- =====================
-- Users (each user has a "home" org via org_id)
-- =====================
CREATE TABLE users (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    sub        text UNIQUE,                    -- unique when non-null (Postgres allows multiple NULLs)
    email      citext NOT NULL,               -- case-insensitive email
    name       text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (org_id, email)
);

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE users
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_users ON users
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

CREATE INDEX idx_users_sub ON users(sub);
CREATE INDEX idx_users_org_email ON users(org_id, email);

-- =====================
-- User ↔ Role (scoped per org)
-- =====================
CREATE TABLE user_roles (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    org_id  uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id, org_id)
);
ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_roles
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_user_roles ON user_roles
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

CREATE INDEX idx_user_roles_org_user ON user_roles(org_id, user_id);
CREATE INDEX idx_user_roles_org_role ON user_roles(org_id, role_id);

-- =====================
-- User ↔ Organisation membership (multi-org)
-- =====================
CREATE TABLE organisation_users (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id  uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, org_id)
);
ALTER TABLE organisation_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE organisation_users
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_organisation_users ON organisation_users
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

CREATE INDEX idx_org_users_org ON organisation_users(org_id);
CREATE INDEX idx_org_users_user ON organisation_users(user_id);

-- =====================
-- Organisation Invites
-- =====================
CREATE TABLE organisation_invites (
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    email        citext NOT NULL,
    role_id      uuid REFERENCES roles(id),
    invited_by   uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token        text NOT NULL UNIQUE,
    accepted     boolean NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL
);
ALTER TABLE organisation_invites ENABLE ROW LEVEL SECURITY;
ALTER TABLE organisation_invites
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_invites ON organisation_invites
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

-- prevent multiple pending invites for same email in an org
CREATE UNIQUE INDEX uniq_pending_invite_per_email_per_org
ON organisation_invites(org_id, email)
WHERE accepted = false;

-- =====================
-- API Keys (user-scoped)
-- =====================
CREATE TABLE api_keys (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_hash   text NOT NULL,   -- store hash, not raw key
    created_at timestamptz NOT NULL DEFAULT now(),
    last_used  timestamptz
);
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_api_keys ON api_keys
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

CREATE INDEX idx_api_keys_org_user ON api_keys(org_id, user_id);

-- =====================
-- Organisation Keys (org-scoped)
-- =====================
CREATE TABLE organisation_keys (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    key_hash   text NOT NULL,  -- store hash, not raw key
    created_at timestamptz NOT NULL DEFAULT now(),
    last_used  timestamptz
);
ALTER TABLE organisation_keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE organisation_keys
  ALTER COLUMN org_id SET DEFAULT app_current_org();

CREATE POLICY org_isolation_org_keys ON organisation_keys
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

CREATE INDEX idx_org_keys_org ON organisation_keys(org_id);

-- +goose Down
DROP TABLE IF EXISTS organisation_keys CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS organisation_invites CASCADE;
DROP TABLE IF EXISTS organisation_users CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS organisations_roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS organisations CASCADE;

DROP FUNCTION IF EXISTS set_updated_at();
DROP FUNCTION IF EXISTS app_current_org();
