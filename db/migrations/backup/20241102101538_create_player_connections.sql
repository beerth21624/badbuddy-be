-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."player_connections" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "requestor_id" uuid NOT NULL,
    "requestee_id" uuid NOT NULL,
    "status" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "player_connections_requestor_id_fkey" FOREIGN KEY ("requestor_id") REFERENCES "public"."users"("id") ON DELETE CASCADE,
    CONSTRAINT "player_connections_requestee_id_fkey" FOREIGN KEY ("requestee_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."player_connections";
-- +goose StatementEnd
