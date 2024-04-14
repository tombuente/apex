INSERT INTO logistics.item_categories (id, name)
VALUES (1, 'None')
ON CONFLICT DO NOTHING;

INSERT INTO logistics.addresses (zip, city, street, country)
VALUES ('44444', 'a city', 'a street', 'a country')

INSERT INTO logistics.addresses (name, sku, category_id, gross_price, net_price)
VALUES ('an item', 1, 499.99, 899.99)
