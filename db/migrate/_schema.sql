CREATE TABLE people
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE INDEX people_name_idx ON people (name);
