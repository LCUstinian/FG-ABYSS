-- migrations/001_init.sql
CREATE TABLE projects (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  INTEGER NOT NULL,
    updated_at  INTEGER NOT NULL,
    deleted_at  INTEGER
);

CREATE TABLE webshells (
    id                TEXT PRIMARY KEY,
    name              TEXT NOT NULL,
    url               TEXT NOT NULL,
    password          TEXT NOT NULL,
    payload_type      TEXT NOT NULL,
    project_id        TEXT REFERENCES projects(id),
    status            TEXT NOT NULL DEFAULT 'inactive'
                          CHECK (status IN ('inactive', 'active', 'needs_redeploy')),
    tags              TEXT NOT NULL DEFAULT '[]',
    custom_headers    TEXT NOT NULL DEFAULT '{}',
    cookies           TEXT NOT NULL DEFAULT '{}',
    proxy_override    TEXT,
    http_method       TEXT NOT NULL DEFAULT 'post',
    c2_profile        TEXT NOT NULL DEFAULT 'default',
    crypto_chain      TEXT NOT NULL DEFAULT '[]',
    fingerprint       TEXT,
    notes             TEXT,
    last_connected_at INTEGER,
    created_at        INTEGER NOT NULL,
    updated_at        INTEGER NOT NULL,
    deleted_at        INTEGER
);
CREATE INDEX idx_webshells_active ON webshells(deleted_at, project_id);
CREATE INDEX idx_webshells_status ON webshells(status, deleted_at);
