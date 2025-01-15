-- +goose Up
-- Enable the uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TYPE "public"."player_level_enum" AS ENUM ('beginner', 'intermediate', 'advanced');
CREATE TYPE "public"."session_status_enum" AS ENUM ('open', 'full', 'completed', 'cancelled');
CREATE TYPE "public"."participant_status_enum" AS ENUM ('confirmed', 'pending', 'cancelled');
CREATE TYPE "public"."user_status_enum" AS ENUM ('active', 'inactive');
CREATE TYPE "public"."venue_status" AS ENUM ('active', 'inactive', 'pending', 'suspended');

CREATE TABLE "public"."venues" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(255) NOT NULL,
    "address" text NOT NULL,
    "location" varchar(255) NOT NULL,
    "description" text,
    "phone" varchar(20),
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    "owner_id" uuid,
    "email" varchar(255),
    "open_time" time NOT NULL DEFAULT '00:00:00'::time without time zone,
    "close_time" time NOT NULL DEFAULT '00:00:00'::time without time zone,
    "status" "public"."venue_status" NOT NULL DEFAULT 'pending'::venue_status,
    "rating" numeric(3,2) DEFAULT 0.00,
    "total_reviews" int4 DEFAULT 0,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "image_urls" text,
    "search_vector" tsvector,
    PRIMARY KEY ("id")
);

-- Indexes
CREATE INDEX idx_venues_location ON public.venues (location);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."venues";
SELECT 'down SQL query';
DROP TYPE IF EXISTS "public"."player_level_enum";
DROP TYPE IF EXISTS "public"."session_status_enum";
DROP TYPE IF EXISTS "public"."participant_status_enum";
DROP TYPE IF EXISTS "public"."user_status_enum";
DROP TYPE IF EXISTS "public"."venue_status";
-- +goose StatementEnd
