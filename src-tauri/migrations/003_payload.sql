-- migrations/003_payload.sql
CREATE TABLE payloads (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    config       TEXT NOT NULL,
    created_at   INTEGER NOT NULL,
    deleted_at   INTEGER
);

CREATE TABLE payload_history (
    id               TEXT PRIMARY KEY,
    payload_id       TEXT REFERENCES payloads(id),
    webshell_id      TEXT,
    code             TEXT NOT NULL,
    template_version TEXT NOT NULL DEFAULT 'v1',
    created_at       INTEGER NOT NULL
);
CREATE INDEX idx_payload_history_webshell ON payload_history(webshell_id, created_at);
