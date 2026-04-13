-- 数据库迁移脚本 001
-- 创建所有基础表

-- 项目表
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER
);

-- WebShell 表
CREATE TABLE IF NOT EXISTS webshells (
    id TEXT PRIMARY KEY,
    project_id TEXT,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    password TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    encryption TEXT DEFAULT 'aes-256-gcm',
    status TEXT DEFAULT 'unknown',
    last_connected_at INTEGER,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

-- 载荷配置表
CREATE TABLE IF NOT EXISTS payload_configs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    password TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    encryption TEXT DEFAULT 'aes-256-gcm',
    obfuscation_level INTEGER DEFAULT 1,
    tags TEXT,
    "group" TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    deleted_at INTEGER
);

-- 载荷模板表
CREATE TABLE IF NOT EXISTS payload_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    payload_type TEXT NOT NULL,
    code TEXT NOT NULL,
    is_built_in INTEGER DEFAULT 0,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- 载荷历史表
CREATE TABLE IF NOT EXISTS payload_history (
    id TEXT PRIMARY KEY,
    config_id TEXT NOT NULL,
    generated_code TEXT NOT NULL,
    generated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (config_id) REFERENCES payload_configs(id)
);

-- 插件表
CREATE TABLE IF NOT EXISTS plugins (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    enabled INTEGER DEFAULT 1,
    config TEXT,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_projects_deleted ON projects(deleted_at);
CREATE INDEX IF NOT EXISTS idx_webshells_project ON webshells(project_id);
CREATE INDEX IF NOT EXISTS idx_webshells_deleted ON webshells(deleted_at);
CREATE INDEX IF NOT EXISTS idx_payload_configs_deleted ON payload_configs(deleted_at);
