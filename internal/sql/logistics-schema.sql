CREATE TABLE IF NOT EXISTS logistics_items (
    id          INTEGER PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    sku         VARCHAR(255) UNIQUE NOT NULL,
    category_id INTEGER NOT NULL,
    gross_price INTEGER DEFAULT 0 NOT NULL,
    net_price   INTEGER DEFAULT 0 NOT NULL,
    FOREIGN KEY(category_id) REFERENCES logistics_item_categories(id)
);

CREATE TABLE IF NOT EXISTS logistics_item_categories (
    id   INTEGER PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS logistics_addresses (
    id            INTEGER PRIMARY KEY,
    zip           INTEGER NOT NULL,
    city          VARCHAR(255) NOT NULL,
    street        VARCHAR(255) NOT NULL,
    street_number VARCHAR(255) NOT NULL,
    country       VARCHAR(255) NOT NULL,
    longitude     REAL NOT NULL,
    latitude      REAL NOT NULL
);

CREATE TABLE IF NOT EXISTS logistics_plants (
    id         INTEGER PRIMARY KEY,
    name       VARCHAR(255) UNIQUE NOT NULL,
    address_id INTEGER NOT NULL,
    FOREIGN KEY(address_id) REFERENCES logistics_addresses(id)
);
