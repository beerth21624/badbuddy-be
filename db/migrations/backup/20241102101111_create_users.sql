-- +goose Up
-- Enable the uuid-ossp extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" varchar(255) NOT NULL,
    "password" text NOT NULL,
    "first_name" varchar(100) NOT NULL,
    "last_name" varchar(100) NOT NULL,
    "phone" varchar(20),
    "play_level" "public"."player_level_enum" NOT NULL,
    "location" varchar(255),
    "bio" text,
    "avatar_url" text,
    "status" "public"."user_status_enum" DEFAULT 'active',
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_active_at" timestamptz,
    PRIMARY KEY ("id")
);

-- Indexes
CREATE UNIQUE INDEX users_email_key ON public.users (email);
CREATE INDEX idx_users_location_level ON public.users (location, play_level);
CREATE INDEX idx_users_status ON public.users (status);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "public"."users";
-- +goose StatementEnd
