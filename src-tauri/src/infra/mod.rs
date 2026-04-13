pub mod database;

use rusqlite::Result;

/// 初始化数据库层
pub fn init() -> Result<()> {
    let db_path = database::get_database_path()
        .map_err(|e| rusqlite::Error::SqliteFailure(
            rusqlite::ffi::Error::new(1),
            Some(format!("获取数据库路径失败：{}", e))
        ))?;
    
    let conn = database::init_db(&db_path)?;
    database::run_migrations(&conn)?;
    
    Ok(())
}
