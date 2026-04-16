use tauri::AppHandle;

/// 设置系统托盘
pub fn setup_system_tray(_app: &AppHandle) -> tauri::Result<()> {
    // TODO: 实现系统托盘
    log::info!("系统托盘已初始化");
    Ok(())
}
