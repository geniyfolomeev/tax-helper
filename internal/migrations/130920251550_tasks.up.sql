CREATE TABLE tasks (
   id SERIAL PRIMARY KEY,
   entrepreneur_id BIGINT NOT NULL REFERENCES entrepreneurs(id) ON DELETE CASCADE,
   status TEXT NOT NULL,
   type TEXT NOT NULL,
   run_at TIMESTAMP NOT NULL
);
