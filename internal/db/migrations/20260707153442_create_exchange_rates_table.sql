-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS exchange_rates
(
    id                 INTEGER PRIMARY KEY,
    base_currency_id   INT REFERENCES currencies (id) ON DELETE CASCADE,
    target_currency_id INT REFERENCES currencies (id) ON DELETE CASCADE,
    rate               REAL NOT NULL,

    UNIQUE (base_currency_id, target_currency_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchange_rates;
-- +goose StatementEnd
