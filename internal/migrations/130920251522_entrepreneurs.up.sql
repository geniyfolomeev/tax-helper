CREATE TABLE entrepreneurs (
   id BIGINT PRIMARY KEY,
   status TEXT NOT NULL,
   registered_at TIMESTAMP NOT NULL,
   last_sent_at TIMESTAMP,
   year_total_amount NUMERIC NOT NULL
);