-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
ALTER TABLE venues DROP COLUMN "open_time";
ALTER TABLE venues DROP COLUMN "close_time";
ALTER TABLE venues ADD COLUMN "open_range" JSON;
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
