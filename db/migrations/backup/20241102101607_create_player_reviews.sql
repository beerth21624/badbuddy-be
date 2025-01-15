-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."player_reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "reviewer_id" uuid NOT NULL,
    "reviewee_id" uuid NOT NULL,
    "rating" int NOT NULL CHECK (rating BETWEEN 1 AND 5),
    "comment" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "player_reviews_reviewer_id_fkey" FOREIGN KEY ("reviewer_id") REFERENCES "public"."users"("id") ON DELETE CASCADE,
    CONSTRAINT "player_reviews_reviewee_id_fkey" FOREIGN KEY ("reviewee_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."player_reviews";
-- +goose StatementEnd
