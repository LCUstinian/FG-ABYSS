-- migrations/004_audit.sql
CREATE TABLE audit_events (
    id           TEXT PRIMARY KEY,
    webshell_id  TEXT,
    action       TEXT NOT NULL,
    detail       TEXT NOT NULL,
    created_at   INTEGER NOT NULL
);
CREATE INDEX idx_audit_webshell ON audit_events(webshell_id, created_at);
