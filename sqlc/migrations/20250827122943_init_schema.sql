-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================
-- Organisations
-- =====================
CREATE TABLE organisations (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE organisations ENABLE ROW LEVEL SECURITY;

CREATE POLICY org_isolation_organisations ON organisations
    USING (id = current_setting('app.current_org')::uuid);

-- =====================
-- Roles
-- =====================
CREATE TABLE roles (
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id      uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    name        text NOT NULL,
    description text,
    created_at  timestamptz NOT NULL DEFAULT now(),
    UNIQUE (org_id, name)
);

ALTER TABLE roles ENABLE ROW LEVEL SECURITY;

ALTER TABLE roles
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_roles ON roles
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

-- =====================
-- Users
-- =====================
CREATE TABLE users (
    id         uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id     uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    sub        text UNIQUE,
    email      text NOT NULL,
    name       text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (org_id, email)
);

ALTER TABLE users ENABLE ROW LEVEL SECURITY;

ALTER TABLE users
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_users ON users
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

-- =====================
-- User ↔ Role join
-- =====================
CREATE TABLE user_roles (
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id uuid NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    org_id  uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;

ALTER TABLE user_roles
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_user_roles ON user_roles
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

-- =====================
-- Organisation Invites
-- =====================
CREATE TABLE organisation_invites (
    id           uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id       uuid NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
    email        text NOT NULL,
    role_id      uuid REFERENCES roles(id),
    invited_by   uuid NOT NULL REFERENCES users(id),
    token        text NOT NULL,
    accepted     boolean NOT NULL DEFAULT false,
    created_at   timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL
);

ALTER TABLE organisation_invites ENABLE ROW LEVEL SECURITY;

ALTER TABLE organisation_invites
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_invites ON organisation_invites
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

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
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_api_keys ON api_keys
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

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
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;

CREATE POLICY org_isolation_org_keys ON organisation_keys
    USING (org_id = current_setting('app.current_org')::uuid)
    WITH CHECK (org_id = current_setting('app.current_org')::uuid);

-- +goose Down
DROP TABLE IF EXISTS organisation_keys CASCADE;
DROP TABLE IF EXISTS api_keys CASCADE;
DROP TABLE IF EXISTS organisation_invites CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS organisations CASCADE;
