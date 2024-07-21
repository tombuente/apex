INSERT INTO accounting.currencies (id, iso, name)
VALUES
    (1, 'USD', 'U.S. dollar'),
    (2, 'EUR', 'Euro'),
    (3, 'JPY', 'Japanese yen'),
    (4, 'GBP', 'Pound sterling'),
    (5, 'CNY', 'Renminbi'),
    (6, 'AUD', 'Australian dollar'),
    (7, 'CAD', 'Canadian dollar'),
    (8, 'CHF', 'Swiss franc'),
    (9, 'HKD', 'Hong Kong dollar'),
    (10, 'SGD', 'Singapore dollar'),
    (11, 'SEK', 'Swedish krona'),
    (12, 'KRW', 'South Korean won'),
    (13, 'NOK', 'Norwegian krone'),
    (14, 'NZD', 'New Zealand dollar'),
    (15, 'INR', 'Indian rupee'),
    (16, 'MXN', 'Mexican peso'),
    (17, 'TWD', 'New Taiwan dollar'),
    (18, 'ZAR', 'South African rand'),
    (19, 'BRL', 'Brazilian real'),
    (20, 'DKK', 'Danish krone'),
    (21, 'PLN', 'Polish złoty'),
    (22, 'THB', 'Thai baht'),
    (23, 'ILS', 'Israeli new shekel'),
    (24, 'IDR', 'Indonesian rupiah'),
    (25, 'CZK', 'Czech koruna'),
    (26, 'AED', 'UAE dirham'),
    (27, 'TRY', 'Turkish lira'),
    (28, 'HUF', 'Hungarian forint'),
    (29, 'CLP', 'Chilean peso'),
    (30, 'SAR', 'Saudi riyal'),
    (31, 'PHP', 'Philippine peso'),
    (32, 'MYR', 'Malaysian ringgit'),
    (33, 'COP', 'Colombian peso'),
    (34, 'RUB', 'Russian ruble'),
    (35, 'RON', 'Romanian leu'),
    (36, 'PEN', 'Peruvian sol'),
    (37, 'BHD', 'Bahraini dinar'),
    (38, 'BGN', 'Bulgarian lev'),
    (39, 'ARS', 'Argentine peso')
ON CONFLICT DO NOTHING;

INSERT INTO accounting.document_position_types (id, description)
VALUES
    (1, 'Debit'),
    (2, 'Credit')
ON CONFLICT DO NOTHING;
