-- +goose Up
ALTER TABLE events ADD COLUMN notify_sent BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE events DROP COLUMN notify_sent;
