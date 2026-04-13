pub mod models;

use rusqlite::{Connection, Result};
use std::path::Path;

/// 初始化数据库连接
pub fn init_db(path: &Path) -> Result<Connection> {
    let conn = Connection::open(path)?;
    
    // 启用外键支持
    conn.execute("PRAGMA foreign_keys = ON", [])?;
    
    Ok(conn)
}

/// 运行数据库迁移
pub fn run_migrations(conn: &Connection) -> Result<()> {
    let migrations = [
        include_str!("../../../../migrations/001_create_tables.sql"),
    ];
    
    for migration in migrations.iter() {
        conn.execute_batch(migration)?;
    }
    
    Ok(())
}

/// 获取数据库路径
pub fn get_database_path() -> std::result::Result<std::path::PathBuf, std::io::Error> {
    let config_dir = dirs::config_dir()
        .ok_or_else(|| std::io::Error::new(
            std::io::ErrorKind::NotFound,
            "无法获取配置目录"
        ))?;
    
    let db_dir = config_dir.join("fg-abyss");
    std::fs::create_dir_all(&db_dir)?;
    
    Ok(db_dir.join("fg-abyss.db"))
}
