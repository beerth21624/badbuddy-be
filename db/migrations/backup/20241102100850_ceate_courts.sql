-- +goose Up
-- Enable the uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."courts" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "price_per_hour" float4,
    "status" varchar,
    "updated_at" timestamptz,
    CONSTRAINT "courts_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "public"."venues"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- Indexes
CREATE UNIQUE INDEX courts_venue_id_name_key ON public.courts (venue_id, name);
CREATE INDEX idx_courts_venue ON public.courts (venue_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."courts";
-- +goose StatementEnd
