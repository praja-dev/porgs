DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS user;

CREATE TABLE IF NOT EXISTS user
(
    username TEXT PRIMARY KEY,
    password TEXT,                       -- Argon2 hashed
    salt     TEXT,                       -- Salt for hashing algorithm
    status   INTEGER NOT NULL DEFAULT 0, -- 0: disabled, 1: enabled, 2: locked, 3: archived, 4: deleted
    email    TEXT                        -- Not unique
) WITHOUT ROWID;

CREATE TABLE IF NOT EXISTS session
(
    id       TEXT PRIMARY KEY,
    created  INTEGER NOT NULL,
    updated  INTEGER NOT NULL,
    username TEXT    NOT NULL,
    FOREIGN KEY (username) REFERENCES user (username)
) WITHOUT ROWID;
