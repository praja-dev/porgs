DROP TABLE IF EXISTS membership;
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS external_org;
DROP TABLE IF EXISTS org;

CREATE TABLE IF NOT EXISTS org
(
    id      INTEGER PRIMARY KEY,
    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,
    name    TEXT    NOT NULL,
    type    TEXT    NOT NULL DEFAULT 'praja',
    storage INTEGER NOT NULL DEFAULT 0, -- 0: local, 1: external
    parent  INTEGER NOT NULL DEFAULT 1, -- The root organization has id 1
    FOREIGN KEY (parent) REFERENCES org (id)
);

CREATE TABLE IF NOT EXISTS external_org
(
    id      INTEGER PRIMARY KEY,
    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,
    url     TEXT    NOT NULL,
    key     TEXT    NOT NULL,
    org     INTEGER NOT NULL,
    FOREIGN KEY (org) REFERENCES org (id)
);

CREATE TABLE IF NOT EXISTS person
(
    id          INTEGER PRIMARY KEY,
    created     INTEGER NOT NULL,
    updated     INTEGER NOT NULL,
    type        INTEGER NOT NULL DEFAULT 0, -- 1: human, 2: company, 3: system
    status      INTEGER NOT NULL DEFAULT 0, -- 0: inactive, 1: active, 2: paused, 3: stopped, 4: deleted
    pname       TEXT,                       -- Preferred name
    name        TEXT,                       -- JSON: {title, given, middle, family, suffix}
    addresses   TEXT,                       -- JSON: []{type, primary, street, locality, region, code, country}
    phones      TEXT,                       -- JSON: []{type, primary, value}
    emails      TEXT,                       -- JSON: []{type, primary, value}
    dob         INTEGER,                    -- Date of birth as Unix timestamp
    gender      TEXT,                       -- M (Male), F (Female), O (Other)
    marital     TEXT,                       -- M (Married), S (Single), D (Divorced), W (Widowed),
    religion    TEXT,
    race        TEXT,
    citizenship TEXT
);

CREATE TABLE IF NOT EXISTS membership
(
    id      INTEGER PRIMARY KEY,
    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,
    org     INTEGER NOT NULL,
    person  INTEGER NOT NULL,
    status  INTEGER NOT NULL DEFAULT 0, -- 0: inactive, 1: active, 2: paused, 3: stopped, 4: deleted
    roles    TEXT,                      -- JSON: []string
    fields   TEXT,                      -- JSON: []{name, value}
    FOREIGN KEY (org)    REFERENCES org    (id),
    FOREIGN KEY (person) REFERENCES person (id)
);
