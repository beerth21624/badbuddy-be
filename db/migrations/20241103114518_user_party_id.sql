-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE users ADD COLUMN venue_id UUID REFERENCES venues(id);

CREATE INDEX idx_users_venue_id ON users(venue_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
