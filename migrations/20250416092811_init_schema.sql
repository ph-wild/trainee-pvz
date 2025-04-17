-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT CHECK (role IN ('employee', 'moderator')) NOT NULL
);

CREATE TABLE pvz (
    id UUID PRIMARY KEY,
    registration_date TIMESTAMPTZ NOT NULL,
    city TEXT CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань')) NOT NULL
);

CREATE TABLE receptions (
    id UUID PRIMARY KEY,
    datetime TIMESTAMPTZ NOT NULL,
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    status TEXT CHECK (status IN ('in_progress', 'close')) NOT NULL
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    datetime TIMESTAMPTZ NOT NULL,
    type TEXT CHECK (type IN ('электроника', 'одежда', 'обувь')) NOT NULL,
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS receptions;
DROP TABLE IF EXISTS pvz;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd