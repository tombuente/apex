INSERT INTO accounting.currencies (id, iso, name)
VALUES
    (1, 'EUR', 'Euro'),
    (2, 'USD', 'US dollar')
    (3, 'GBP', 'Pound sterling')
ON CONFLICT DO NOTHING;

INSERT INTO accounting.document_position_types (id, description)
VALUES
    (1, 'Debit'),
    (2, 'Credit')
ON CONFLICT DO NOTHING;
