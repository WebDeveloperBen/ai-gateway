-- Template for org-scoped table:

-- +goose StatementBegin
DO $$
BEGIN
ALTER TABLE {{table}} ENABLE ROW LEVEL SECURITY;
ALTER TABLE {{table}}
  ALTER COLUMN org_id SET DEFAULT current_setting('app.current_org')::uuid;
CREATE POLICY org_isolation_{{table}} ON {{table}}
  USING (org_id = current_setting('app.current_org')::uuid)
  WITH CHECK (org_id = current_setting('app.current_org')::uuid);
END
$$;
-- +goose StatementEnd
