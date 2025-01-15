-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."session_rules" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "rule_text" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "session_rules_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."play_sessions"("id") ON DELETE CASCADE
);
-- Indexes
CREATE INDEX idx_session_rules_session ON public.session_rules (session_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."session_rules";
-- +goose StatementEnd
