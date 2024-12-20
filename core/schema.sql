DROP TABLE IF EXISTS membership;
DROP TABLE IF EXISTS person;
DROP TABLE IF EXISTS external_org;
DROP TABLE IF EXISTS org;

-- Organization: an organization of people.
CREATE TABLE IF NOT EXISTS org
(
    id      INTEGER PRIMARY KEY,
    created INTEGER NOT NULL,
    updated INTEGER NOT NULL,

    -- ## Parent Organization ID
    -- Root organization has the ID 1
    -- ID 0 is not used
    parent  INTEGER,

    -- ## Sequence ID to order the organization within the parent
    sid  INTEGER,

    -- ## Source of data for this object
    -- 0: Local (default) - data stored in this PORGS instance
    -- 1: Peer - data stored in a peer PORGS instance
    source INTEGER NOT NULL DEFAULT 0,

    -- ## Organization Type
    type    INTEGER NOT NULL DEFAULT 1,

    external_id TEXT,
    external_sid TEXT,

    name    TEXT    NOT NULL,

    -- ## Translations
    trlx TEXT,

    -- ## Custom Properties
    -- Schema for the custom properties is determined by the organization type (type field)
    cxp TEXT,

    -- ## Translations for Custom Properties
    cxp_trlx TEXT,

    FOREIGN KEY (parent) REFERENCES org (id)
);

CREATE TABLE IF NOT EXISTS org_type
(
    id          INTEGER PRIMARY KEY,
    created     INTEGER NOT NULL,
    updated     INTEGER NOT NULL,
    name        TEXT    NOT NULL,
    description TEXT,
    cxp         TEXT
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
