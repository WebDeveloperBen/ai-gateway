-- Applied after Atlas creates the tables

-- Triggers (idempotent creation)
DROP TRIGGER IF EXISTS organisations_set_updated_at ON organisations;
CREATE TRIGGER organisations_set_updated_at
BEFORE UPDATE ON organisations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS users_set_updated_at ON users;
CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS applications_set_updated_at ON applications;
CREATE TRIGGER applications_set_updated_at
BEFORE UPDATE ON applications
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS application_configs_set_updated_at ON application_configs;
CREATE TRIGGER application_configs_set_updated_at
BEFORE UPDATE ON application_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS models_set_updated_at ON models;
CREATE TRIGGER models_set_updated_at
BEFORE UPDATE ON models
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS policies_set_updated_at ON policies;
CREATE TRIGGER policies_set_updated_at
BEFORE UPDATE ON policies
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Row Level Security Policies (safe creation - only create if not exists)
DO $$
BEGIN
    -- organisations
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'organisations' AND policyname = 'org_isolation_organisations') THEN
        ALTER TABLE organisations ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_organisations ON organisations USING (id = app_current_org());
    END IF;

    -- organisation_roles
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'organisation_roles' AND policyname = 'org_isolation_organisation_roles') THEN
        ALTER TABLE organisation_roles ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_organisation_roles ON organisation_roles
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- users
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'users' AND policyname = 'org_isolation_users') THEN
        ALTER TABLE users ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_users ON users
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- user_roles
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'user_roles' AND policyname = 'org_isolation_user_roles') THEN
        ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_user_roles ON user_roles
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- organisation_users
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'organisation_users' AND policyname = 'org_isolation_organisation_users') THEN
        ALTER TABLE organisation_users ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_organisation_users ON organisation_users
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- organisation_invites
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'organisation_invites' AND policyname = 'org_isolation_invites') THEN
        ALTER TABLE organisation_invites ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_invites ON organisation_invites
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- api_keys
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'api_keys' AND policyname = 'org_isolation_api_keys') THEN
        ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_api_keys ON api_keys
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- organisation_keys
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'organisation_keys' AND policyname = 'org_isolation_org_keys') THEN
        ALTER TABLE organisation_keys ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_org_keys ON organisation_keys
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- applications
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'applications' AND policyname = 'org_isolation_applications') THEN
        ALTER TABLE applications ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_applications ON applications
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- application_configs
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'application_configs' AND policyname = 'org_isolation_application_configs') THEN
        ALTER TABLE application_configs ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_application_configs ON application_configs
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- models
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'models' AND policyname = 'org_isolation_models') THEN
        ALTER TABLE models ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_models ON models
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- policies
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'policies' AND policyname = 'org_isolation_policies') THEN
        ALTER TABLE policies ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_policies ON policies
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;

    -- usage_metrics
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'usage_metrics' AND policyname = 'org_isolation_usage_metrics') THEN
        ALTER TABLE usage_metrics ENABLE ROW LEVEL SECURITY;
        CREATE POLICY org_isolation_usage_metrics ON usage_metrics
          USING (org_id = app_current_org()) WITH CHECK (org_id = app_current_org());
    END IF;
END $$;

-- Insert seed data for roles
INSERT INTO roles (name, description)
VALUES
  ('owner', 'Organisation owner'),
  ('admin', 'Admin user'),
  ('member', 'Standard member')
ON CONFLICT (name) DO NOTHING;