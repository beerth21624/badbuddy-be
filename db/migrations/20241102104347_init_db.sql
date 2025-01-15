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

-- Create Users table
CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "email" varchar(255) NOT NULL,
    "password" text NOT NULL,
    "first_name" varchar(100) NOT NULL,
    "last_name" varchar(100) NOT NULL,
    "phone" varchar(20),
    "play_level" player_level_enum NOT NULL,
    "location" varchar(255),
    "bio" text,
    "avatar_url" text,
    "status" user_status_enum DEFAULT 'active',
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_active_at" timestamptz,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users USING btree (email);
CREATE INDEX IF NOT EXISTS idx_users_location_level ON users USING btree (location, play_level);
CREATE INDEX IF NOT EXISTS idx_users_status ON users USING btree (status);

-- Create Venues table
CREATE TABLE IF NOT EXISTS "venues" (
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
    "open_time" time NOT NULL DEFAULT '00:00:00',
    "close_time" time NOT NULL DEFAULT '00:00:00',
    "status" venue_status NOT NULL DEFAULT 'pending',
    "rating" numeric(3,2) DEFAULT 0.00,
    "total_reviews" int4 DEFAULT 0,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "image_urls" text,
    "search_vector" tsvector,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS idx_venues_location ON venues USING btree (location);

-- Create Courts table
CREATE TABLE IF NOT EXISTS "courts" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "name" varchar(100) NOT NULL,
    "description" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamptz,
    "price_per_hour" float4,
    "status" varchar,
    "updated_at" timestamptz,
    CONSTRAINT "courts_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "venues"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS courts_venue_id_name_key ON courts USING btree (venue_id, name);
CREATE INDEX IF NOT EXISTS idx_courts_venue ON courts USING btree (venue_id);

-- Create Play Sessions table
CREATE TABLE IF NOT EXISTS "play_sessions" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "host_id" uuid NOT NULL,
    "venue_id" uuid NOT NULL,
    "title" varchar(255) NOT NULL,
    "description" text,
    "player_level" player_level_enum NOT NULL,
    "session_date" date NOT NULL,
    "start_time" time NOT NULL,
    "end_time" time NOT NULL,
    "max_participants" int4 NOT NULL DEFAULT 4,
    "cost_per_person" numeric(10,2) NOT NULL DEFAULT 0,
    "allow_cancellation" bool NOT NULL DEFAULT true,
    "cancellation_deadline_hours" int4 DEFAULT 2,
    "status" session_status_enum DEFAULT 'open',
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "search_vector" tsvector,
    CONSTRAINT "play_sessions_host_id_fkey" FOREIGN KEY ("host_id") REFERENCES "users"("id"),
    CONSTRAINT "play_sessions_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "venues"("id"),
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS idx_sessions_search ON play_sessions USING gin (search_vector);
CREATE INDEX IF NOT EXISTS idx_sessions_date_status ON play_sessions USING btree (session_date, status);
CREATE INDEX IF NOT EXISTS idx_sessions_host ON play_sessions USING btree (host_id);
CREATE INDEX IF NOT EXISTS idx_sessions_venue ON play_sessions USING btree (venue_id);
CREATE INDEX IF NOT EXISTS idx_sessions_level ON play_sessions USING btree (player_level);

-- Create Player Connections table
CREATE TABLE IF NOT EXISTS "player_connections" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "player1_id" uuid NOT NULL,
    "player2_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "player_connections_player1_id_fkey" FOREIGN KEY ("player1_id") REFERENCES "users"("id"),
    CONSTRAINT "player_connections_player2_id_fkey" FOREIGN KEY ("player2_id") REFERENCES "users"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS player_connections_player1_id_player2_id_key ON player_connections USING btree (player1_id, player2_id);

-- Create Player Reviews table
CREATE TABLE IF NOT EXISTS "player_reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "reviewer_id" uuid NOT NULL,
    "reviewed_id" uuid NOT NULL,
    "session_id" uuid NOT NULL,
    "rating" int4 NOT NULL,
    "comment" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "player_reviews_reviewer_id_fkey" FOREIGN KEY ("reviewer_id") REFERENCES "users"("id"),
    CONSTRAINT "player_reviews_reviewed_id_fkey" FOREIGN KEY ("reviewed_id") REFERENCES "users"("id"),
    CONSTRAINT "player_reviews_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "play_sessions"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS player_reviews_reviewer_id_reviewed_id_session_id_key ON player_reviews USING btree (reviewer_id, reviewed_id, session_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviewed ON player_reviews USING btree (reviewed_id);
CREATE INDEX IF NOT EXISTS idx_reviews_session ON player_reviews USING btree (session_id);

-- Create Session Courts table
CREATE TABLE IF NOT EXISTS "session_courts" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "court_id" uuid NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "session_courts_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "play_sessions"("id") ON DELETE CASCADE,
    CONSTRAINT "session_courts_court_id_fkey" FOREIGN KEY ("court_id") REFERENCES "courts"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS session_courts_session_id_court_id_key ON session_courts USING btree (session_id, court_id);

-- Create Session Participants table
CREATE TABLE IF NOT EXISTS "session_participants" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "status" participant_status_enum DEFAULT 'pending',
    "joined_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "cancelled_at" timestamptz,
    CONSTRAINT "session_participants_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "play_sessions"("id") ON DELETE CASCADE,
    CONSTRAINT "session_participants_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS session_participants_session_id_user_id_key ON session_participants USING btree (session_id, user_id);
CREATE INDEX IF NOT EXISTS idx_participants_session ON session_participants USING btree (session_id);
CREATE INDEX IF NOT EXISTS idx_participants_user ON session_participants USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_participants_status ON session_participants USING btree (status);

-- Create Session Rules table
CREATE TABLE IF NOT EXISTS "session_rules" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "session_id" uuid NOT NULL,
    "rule_text" text NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "session_rules_session_id_fkey" FOREIGN KEY ("session_id") REFERENCES "play_sessions"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- Create Venue Images table
CREATE TABLE IF NOT EXISTS "venue_images" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "url" text NOT NULL,
    "description" text,
    "order_index" int4 NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "venue_images_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "venues"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_venue_image_order ON venue_images USING btree (venue_id, order_index);
CREATE INDEX IF NOT EXISTS idx_venue_images_venue_id ON venue_images USING btree (venue_id);

-- Create Venue Reviews table
CREATE TABLE IF NOT EXISTS "venue_reviews" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "venue_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "rating" int4 NOT NULL,
    "comment" text,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "venue_reviews_venue_id_fkey" FOREIGN KEY ("venue_id") REFERENCES "venues"("id"),
    CONSTRAINT "venue_reviews_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_user_venue_review ON venue_reviews USING btree (venue_id, user_id);
CREATE INDEX IF NOT EXISTS idx_venue_reviews_venue_id ON venue_reviews USING btree (venue_id);
CREATE INDEX IF NOT EXISTS idx_venue_reviews_user_id ON venue_reviews USING btree (user_id);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS "session_rules" CASCADE;
DROP TABLE IF EXISTS "session_participants" CASCADE;
DROP TABLE IF EXISTS "session_courts" CASCADE;
DROP TABLE IF EXISTS "player_reviews" CASCADE;
DROP TABLE IF EXISTS "player_connections" CASCADE;
DROP TABLE IF EXISTS "play_sessions" CASCADE;
DROP TABLE IF EXISTS "venue_reviews" CASCADE;
DROP TABLE IF EXISTS "venue_images" CASCADE;
DROP TABLE IF EXISTS "courts" CASCADE;
DROP TABLE IF EXISTS "venues" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;

-- Drop types
DROP TYPE IF EXISTS player_level_enum CASCADE;
DROP TYPE IF EXISTS session_status_enum CASCADE;
DROP TYPE IF EXISTS participant_status_enum CASCADE;
DROP TYPE IF EXISTS venue_status CASCADE;
DROP TYPE IF EXISTS user_status_enum CASCADE;
-- +goose StatementEnd
