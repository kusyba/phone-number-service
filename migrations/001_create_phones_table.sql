-- +goose Up
CREATE TABLE IF NOT EXISTS phones (
    id SERIAL PRIMARY KEY,
    number VARCHAR(20) UNIQUE NOT NULL,
    country VARCHAR(100),
    region VARCHAR(100),
    provider VARCHAR(100),
    source VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_phones_number ON phones(number);
CREATE INDEX IF NOT EXISTS idx_phones_country ON phones(country);
CREATE INDEX IF NOT EXISTS idx_phones_region ON phones(region);
CREATE INDEX IF NOT EXISTS idx_phones_provider ON phones(provider);

-- +goose Down
DROP INDEX IF EXISTS idx_phones_number;
DROP INDEX IF EXISTS idx_phones_country;
DROP INDEX IF EXISTS idx_phones_region;
DROP INDEX IF EXISTS idx_phones_provider;
DROP TABLE IF EXISTS phones;
