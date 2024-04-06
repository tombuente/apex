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
