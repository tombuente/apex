INSERT INTO accounting.currencies (id, iso, name)
VALUES
    (1, 'EUR', 'Euro'),
    (2, 'USD', 'US dollar')
    (3, 'GBP', 'Pound sterling')
ON CONFLICT DO NOTHING;
