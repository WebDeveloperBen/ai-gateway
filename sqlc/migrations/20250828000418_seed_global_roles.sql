-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name, description)
VALUES
  ('owner', 'Organisation owner'),
  ('admin', 'Admin user'),
  ('member', 'Standard member')
ON CONFLICT ( name) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
