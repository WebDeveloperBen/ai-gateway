-- Applied after Atlas creates the tables

-- Triggers
CREATE TRIGGER organisations_set_updated_at
BEFORE UPDATE ON organisations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER applications_set_updated_at
BEFORE UPDATE ON applications
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER application_configs_set_updated_at
BEFORE UPDATE ON application_configs
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER models_set_updated_at
BEFORE UPDATE ON models
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER policies_set_updated_at
BEFORE UPDATE ON policies
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Row Level Security Policies
ALTER TABLE organisations ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_organisations ON organisations
  USING (id = app_current_org());

ALTER TABLE organisation_roles ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_organisation_roles ON organisation_roles
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_users ON users
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_user_roles ON user_roles
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE organisation_users ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_organisation_users ON organisation_users
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE organisation_invites ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_invites ON organisation_invites
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_api_keys ON api_keys
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE organisation_keys ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_org_keys ON organisation_keys
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE applications ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_applications ON applications
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE application_configs ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_application_configs ON application_configs
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE models ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_models ON models
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE policies ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_policies ON policies
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

ALTER TABLE usage_metrics ENABLE ROW LEVEL SECURITY;
CREATE POLICY org_isolation_usage_metrics ON usage_metrics
  USING (org_id = app_current_org())
  WITH CHECK (org_id = app_current_org());

-- Insert seed data for roles
INSERT INTO roles (name, description)
VALUES
  ('owner', 'Organisation owner'),
  ('admin', 'Admin user'),
  ('member', 'Standard member')
ON CONFLICT (name) DO NOTHING;