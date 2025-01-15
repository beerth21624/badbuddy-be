-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."venue_images" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "url" text NOT NULL,
    "description" text,
    "order_index" int4 NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id"),
    CONSTRAINT "venue_images_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "public"."venues"("id") ON DELETE CASCADE
);
-- Indexes
CREATE UNIQUE INDEX unique_venue_image_order ON public.venue_images (venue_id, order_index);
CREATE INDEX idx_venue_images_venue_id ON public.venue_images (venue_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."venue_images";
-- +goose StatementEnd
