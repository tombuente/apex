INSERT INTO accounting.accounts (description)
VALUES
    ('Cash Account'),
    ('Revenue Account');

INSERT INTO accounting.documents (id, date, posting_date, reference, description, currency_id)
VALUES
    (1, '2024-01-01', '2024-01-01', 'DOC1-REF', 'DOC1', 1);

INSERT INTO accounting.document_positions (document_id, account_id, description, type_id, amount)
VALUES
    (1, 1, 'Cash Sale', 1, 100.00),
    (1, 2, 'Revenue Recognition', 2, 100.00);
