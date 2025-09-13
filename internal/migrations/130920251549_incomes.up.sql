CREATE TABLE incomes (
    id SERIAL PRIMARY KEY,
    entrepreneur_id BIGINT NOT NULL REFERENCES entrepreneurs(id) ON DELETE CASCADE,
    amount NUMERIC NOT NULL,
    source_amount NUMERIC NOT NULL,
    source_currency CHAR(3) NOT NULL,
    date TIMESTAMP NOT NULL
);