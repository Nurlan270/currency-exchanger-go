-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currencies
(
    id        INTEGER PRIMARY KEY,
    full_name VARCHAR NOT NULL,
    code      VARCHAR NOT NULL UNIQUE,
    sign      VARCHAR NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd
