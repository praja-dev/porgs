-- TODO
CREATE TABLE IF NOT EXISTS task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    summary TEXT NOT NULL,
    description TEXT,
    status INTEGER NOT NULL DEFAULT 0
);
