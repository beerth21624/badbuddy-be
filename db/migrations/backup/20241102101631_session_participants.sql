-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."session_participants" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "status" "public"."participant_status_enum" NOT NULL DEFAULT 'pending',
    "joined_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "session_participants_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "public"."play_sessions"("id") ON DELETE CASCADE,
    CONSTRAINT "session_participants_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."session_participants";
-- +goose StatementEnd
