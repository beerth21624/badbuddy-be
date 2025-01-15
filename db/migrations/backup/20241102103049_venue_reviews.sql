-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."venue_reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "rating" int4 NOT NULL,
    "comment" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "venue_reviews_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "public"."venues"("id") ON DELETE CASCADE,
    CONSTRAINT "venue_reviews_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE
);
-- Indexes
CREATE UNIQUE INDEX unique_user_venue_review ON public.venue_reviews (venue_id, user_id);
CREATE INDEX idx_venue_reviews_venue_id ON public.venue_reviews (venue_id);
CREATE INDEX idx_venue_reviews_user_id ON public.venue_reviews (user_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."venue_reviews";
-- +goose StatementEnd
