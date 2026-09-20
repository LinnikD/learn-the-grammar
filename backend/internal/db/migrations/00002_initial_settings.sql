-- +goose Up
CREATE TABLE users (
    id                    uuid PRIMARY KEY,
    level                 text NOT NULL DEFAULT 'A1' CHECK (level IN ('A1', 'A2', 'B1')),
    onboarding_completed_at timestamptz NULL,
    created_at            timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE topics (
    id     uuid PRIMARY KEY,
    name   text NOT NULL,
    active boolean NOT NULL DEFAULT true
);

-- Single stub user every visitor is mapped to until real per-visitor
-- accounts exist (FR-29, DEBT-1: see backend/internal/db/stub.go).
INSERT INTO users (id) VALUES ('00000000-0000-0000-0000-000000000001');

-- +goose Down
DROP TABLE topics;
DROP TABLE users;
