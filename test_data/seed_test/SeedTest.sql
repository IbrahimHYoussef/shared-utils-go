-- +seed Up
-- +seed StatementBegin
INSERT INTO test_table (id, name) VALUES
    (1, 'alice'),
    (2, 'bob'),
    (3, 'carol');
-- +seed StatementEnd

-- +seed Down
-- +seed StatementBegin
DELETE FROM test_table WHERE id IN (1, 2, 3);
-- +seed StatementEnd
