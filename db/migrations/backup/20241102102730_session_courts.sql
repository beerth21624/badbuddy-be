-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."session_courts" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "court_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "session_courts_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."play_sessions"("id") ON DELETE CASCADE,
    CONSTRAINT "session_courts_court_id_fkey" FOREIGN KEY ("court_id") REFERENCES "public"."courts"("id") ON DELETE CASCADE
);

-- Indexes
CREATE UNIQUE INDEX session_courts_session_id_court_id_key ON public.session_courts (session_id, court_id);
CREATE INDEX idx_session_courts_session ON public.session_courts (session_id);
CREATE INDEX idx_session_courts_court ON public.session_courts (court_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."session_courts";
-- +goose StatementEnd
