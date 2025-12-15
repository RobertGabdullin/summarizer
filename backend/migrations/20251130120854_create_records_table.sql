-- +goose Up
-- +goose StatementBegin
CREATE TABLE records(
    id SERIAL PRIMARY KEY,
    prompt TEXT,
    participants INTEGER,
    status TEXT NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd
