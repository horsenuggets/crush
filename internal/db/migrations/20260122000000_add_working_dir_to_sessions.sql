-- +goose Up
-- +goose StatementBegin
ALTER TABLE sessions ADD COLUMN working_dir TEXT;
CREATE INDEX IF NOT EXISTS idx_sessions_working_dir ON sessions (working_dir);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_sessions_working_dir;
-- SQLite doesn't support DROP COLUMN, so we skip that in down migration
-- +goose StatementEnd
