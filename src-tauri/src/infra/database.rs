use rusqlite::{Connection, Result};
use std::path::Path;

pub fn init_db(path: &Path) -> Result<Connection> {
    let conn = Connection::open(path)?;
    
    conn.execute(
        "CREATE TABLE IF NOT EXISTS projects (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            description TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )",
        [],
    )?;
    
    Ok(conn)
}
