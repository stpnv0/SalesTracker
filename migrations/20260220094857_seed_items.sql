-- +goose Up
INSERT INTO items (type, amount, category, description, date) VALUES
    -- Январь
    ('income', 120000.00, 'salary',      'seed__salary_2026_01',     '2026-01-10'),
    ('expense',  3500.00, 'food',        'seed__groceries_01',       '2026-01-11'),
    ('expense',  1200.00, 'transport',   'seed__metro_01',           '2026-01-11'),
    ('expense',  8900.00, 'housing',     'seed__utilities_01',       '2026-01-15'),
    ('expense',  2100.00, 'health',      'seed__pharmacy_01',        '2026-01-18'),
    ('income',   15000.00, 'freelance',  'seed__freelance_01',       '2026-01-20'),
    ('expense',  5400.00, 'entertainment','seed__cinema_01',         '2026-01-22'),
    ('expense',  1999.99, 'education',   'seed__course_01',          '2026-01-25'),

    -- Февраль
    ('income', 120000.00, 'salary',      'seed__salary_2026_02',     '2026-02-10'),
    ('expense',  4200.00, 'food',        'seed__groceries_02',       '2026-02-10'),
    ('expense',  1600.00, 'transport',   'seed__taxi_02',            '2026-02-12'),
    ('expense', 10500.00, 'housing',     'seed__rent_part_02',       '2026-02-15'),
    ('income',    200.00, 'investments', 'seed__dividends_usd_02',   '2026-02-16'),
    ('expense',   35.00,  'subscriptions','seed__subscription_usd',  '2026-02-18');

-- +goose Down
DELETE FROM items WHERE description LIKE 'seed__%';