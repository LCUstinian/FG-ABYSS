-- migrations/002_plugins.sql
CREATE TABLE plugins (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    version     TEXT NOT NULL,
    enabled     INTEGER NOT NULL DEFAULT 1,
    config      TEXT NOT NULL DEFAULT '{}',
    source      TEXT NOT NULL DEFAULT 'builtin',
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL
);
