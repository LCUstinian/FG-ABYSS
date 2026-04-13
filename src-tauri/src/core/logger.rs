use log::{info, LevelFilter};
use env_logger::Builder;

pub fn init() {
    Builder::new()
        .filter(None, LevelFilter::Info)
        .format(|buf, record| {
            use std::io::Write;
            writeln!(
                buf,
                "[{} {}] {}: {}",
                chrono::Local::now().format("%Y-%m-%d %H:%M:%S"),
                std::thread::current().name().unwrap_or("main"),
                record.level(),
                record.args()
            )
        })
        .init();
    
    info!("Logger initialized");
}
