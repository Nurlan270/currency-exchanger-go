-- +goose Up
-- +goose StatementBegin
INSERT INTO currencies (full_name, code, sign)
VALUES ('US Dollar', 'USD', '$'),
       ('Euro', 'EUR', '€'),
       ('Russian Ruble', 'RUB', '₽'),
       ('British Pound', 'GBP', '£'),
       ('Japanese Yen', 'JPY', '¥'),
       ('Canadian Dollar', 'CAD', '$'),
       ('Australian Dollar', 'AUD', '$'),
       ('Swiss Franc', 'CHF', 'CHF'),
       ('Chinese Yuan', 'CNY', '¥'),
       ('Indian Rupee', 'INR', '₹'),
       ('Brazilian Real', 'BRL', 'R$');

INSERT INTO exchange_rates (rate, base_currency_id, target_currency_id)
VALUES (0.85, 1, 2),
       (74.5, 1, 3),
       (0.75, 1, 4),
       (110.0, 1, 5),
       (1.25, 1, 6),
       (1.35, 1, 7),
       (0.92, 1, 8),
       (6.45, 1, 9),
       (74.0, 1, 10),
       (5.2, 1, 11);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM currencies;
DELETE
FROM exchange_rates;
-- +goose StatementEnd
