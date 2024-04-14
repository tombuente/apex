CREATE SCHEMA IF NOT EXISTS logistics;
GRANT ALL ON SCHEMA logistics TO postgres;

CREATE TABLE IF NOT EXISTS logistics.item_categories (
    id   SERIAL       PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS logistics.addresses (
    id        SERIAL       PRIMARY KEY,
    zip       VARCHAR(255) NOT NULL,
    city      VARCHAR(255) NOT NULL,
    street    VARCHAR(255) NOT NULL,
    country   VARCHAR(255) NOT NULL,
    latitude  NUMERIC(8,6) NOT NULL DEFAULT 0,
    longitude NUMERIC(9,6) NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS logistics.items (
    id          SERIAL       PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    sku         VARCHAR(255) NOT NULL UNIQUE,
    category_id INTEGER      NOT NULL REFERENCES logistics.item_categories(id),
    gross_price INTEGER      NOT NULL DEFAULT 0,
    net_price   INTEGER      NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS logistics.plants (
    id         SERIAL       PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    address_id INTEGER      NOT NULL REFERENCES logistics.addresses(id)
);
