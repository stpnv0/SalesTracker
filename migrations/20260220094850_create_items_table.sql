-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type        VARCHAR(10)    NOT NULL CHECK (type IN ('income', 'expense')),
    amount      NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    category    VARCHAR(100)   NOT NULL,
    description TEXT           NOT NULL DEFAULT '',
    date        DATE           NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX idx_items_date ON items (date);
CREATE INDEX idx_items_category ON items (category);
CREATE INDEX idx_items_type ON items (type);

-- +goose Down
DROP TABLE IF EXISTS items;
