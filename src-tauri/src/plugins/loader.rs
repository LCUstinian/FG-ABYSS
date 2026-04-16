use crate::plugins::types::{Plugin, PluginError, PluginMetadata, Result};
use std::fs;
use std::path::{Path, PathBuf};

/// 插件加载器
pub struct PluginLoader {
    /// 插件目录
    plugins_dir: PathBuf,
    /// 临时目录
    temp_dir: PathBuf,
}

impl PluginLoader {
    pub fn new(plugins_dir: PathBuf) -> Self {
        let temp_dir = plugins_dir.join(".temp");
        
        // 确保目录存在
        if !plugins_dir.exists() {
            fs::create_dir_all(&plugins_dir).expect("创建插件目录失败");
        }
        if !temp_dir.exists() {
            fs::create_dir_all(&temp_dir).expect("创建临时目录失败");
        }

        Self {
            plugins_dir,
            temp_dir,
        }
    }

    /// 加载插件
    pub fn load_plugin(&self, plugin_path: &Path) -> Result<Plugin> {
        // 验证插件路径
        if !plugin_path.exists() {
            return Err(PluginError::NotFound(format!(
                "插件路径不存在：{:?}",
                plugin_path
            )));
        }

        // 读取插件元数据
        let metadata_path = plugin_path.join("plugin.json");
        if !metadata_path.exists() {
            return Err(PluginError::ValidationFailed(
                "插件元数据文件不存在".to_string(),
            ));
        }

        let metadata_content = fs::read_to_string(&metadata_path)
            .map_err(|e| PluginError::LoadFailed(format!("读取元数据失败：{}", e)))?;

        let metadata: PluginMetadata = serde_json::from_str(&metadata_content)
            .map_err(|e| PluginError::ValidationFailed(format!("解析元数据失败：{}", e)))?;

        // 创建插件实例
        let plugin = Plugin::new(plugin_path.to_path_buf(), metadata);

        // 验证插件
        plugin.validate(env!("CARGO_PKG_VERSION"))?;

        Ok(plugin)
    }

    /// 卸载插件
    pub fn unload_plugin(&self, plugin_id: &str) -> Result<()> {
        let plugin_dir = self.plugins_dir.join(plugin_id);
        
        if !plugin_dir.exists() {
            return Err(PluginError::NotFound(format!(
                "插件不存在：{}",
                plugin_id
            )));
        }

        // 移动到临时目录
        let temp_path = self.temp_dir.join(format!("{}_{}", plugin_id, chrono::Utc::now().timestamp()));
        fs::rename(&plugin_dir, &temp_path)
            .map_err(|e| PluginError::LoadFailed(format!("移动插件失败：{}", e)))?;

        // 删除临时目录
        fs::remove_dir_all(&temp_path)
            .map_err(|e| PluginError::LoadFailed(format!("删除插件失败：{}", e)))?;

        Ok(())
    }

    /// 扫描插件目录
    pub fn scan_plugins(&self) -> Result<Vec<Plugin>> {
        let mut plugins = Vec::new();

        let entries = fs::read_dir(&self.plugins_dir)
            .map_err(|e| PluginError::LoadFailed(format!("读取插件目录失败：{}", e)))?;

        for entry in entries.flatten() {
            let path = entry.path();
            if path.is_dir() {
                match self.load_plugin(&path) {
                    Ok(plugin) => plugins.push(plugin),
                    Err(e) => {
                        log::warn!("加载插件 {:?} 失败：{}", path, e);
                    }
                }
            }
        }

        Ok(plugins)
    }

    /// 安装插件（从 ZIP 文件）
    pub fn install_plugin(&self, zip_path: &Path) -> Result<String> {
        // 验证 ZIP 文件
        if !zip_path.exists() {
            return Err(PluginError::NotFound(format!(
                "插件文件不存在：{:?}",
                zip_path
            )));
        }

        // 解压到临时目录
        let temp_plugin_dir = self.temp_dir.join(format!("plugin_{}", chrono::Utc::now().timestamp()));
        self.extract_plugin(zip_path, &temp_plugin_dir)?;

        // 加载并验证插件
        let plugin = self.load_plugin(&temp_plugin_dir)?;
        let plugin_id = plugin.metadata.id.clone();

        // 移动到插件目录
        let final_path = self.plugins_dir.join(&plugin_id);
        if final_path.exists() {
            return Err(PluginError::AlreadyExists(plugin_id.clone()));
        }

        fs::rename(&temp_plugin_dir, &final_path)
            .map_err(|e| PluginError::LoadFailed(format!("安装插件失败：{}", e)))?;

        Ok(plugin_id)
    }

    /// 解压插件
    fn extract_plugin(&self, zip_path: &Path, dest_dir: &Path) -> Result<()> {
        // 使用 zip 解压库（需要添加依赖）
        // 这里简化实现
        fs::create_dir_all(dest_dir)
            .map_err(|e| PluginError::LoadFailed(format!("创建目录失败：{}", e)))?;

        // TODO: 实际实现 ZIP 解压
        log::info!("解压插件：{:?} -> {:?}", zip_path, dest_dir);

        Ok(())
    }

    /// 获取插件目录
    pub fn plugins_dir(&self) -> &Path {
        &self.plugins_dir
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::fs;

    #[test]
    fn test_plugin_loader_creation() {
        let temp_dir = std::env::temp_dir().join("fg_abyss_test_plugins");
        let _ = fs::remove_dir_all(&temp_dir); // 清理旧数据
        
        let loader = PluginLoader::new(temp_dir.clone());
        assert!(loader.plugins_dir().exists());
        
        // 清理
        let _ = fs::remove_dir_all(&temp_dir);
    }
}
