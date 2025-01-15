-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."play_sessions" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "start_time" timestamptz NOT NULL,
    "end_time" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "status" "public"."session_status_enum" NOT NULL DEFAULT 'open',
    PRIMARY KEY ("id"),
    CONSTRAINT "play_sessions_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "public"."venues"("id") ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."play_sessions";
-- +goose StatementEnd
