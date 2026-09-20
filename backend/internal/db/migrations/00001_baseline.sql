-- +goose Up
-- Intentional no-op: proves the migration pipeline (goose, this
-- directory, and the tooling in TASK-2) applies cleanly against a
-- fresh database before any product schema exists. FR issues that
-- need persistent storage add real migrations after this one.
SELECT 1;

-- +goose Down
SELECT 1;
