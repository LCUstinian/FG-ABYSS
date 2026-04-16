use crate::plugins::types::{Plugin, PluginConfig, PluginError, PluginStatus, Result};
use crate::plugins::loader::PluginLoader;
use crate::plugins::sandbox::{PluginSandbox, ResourceLimits};
use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

/// 插件运行时
pub struct PluginRuntime {
    /// 已加载的插件
    plugins: Arc<RwLock<HashMap<String, Plugin>>>,
    /// 插件配置
    configs: Arc<RwLock<HashMap<String, PluginConfig>>>,
    /// 插件加载器
    loader: PluginLoader,
    /// 插件沙箱
    sandbox: Arc<RwLock<PluginSandbox>>,
}

impl PluginRuntime {
    pub fn new(plugins_dir: std::path::PathBuf) -> Self {
        let loader = PluginLoader::new(plugins_dir);
        let sandbox = Arc::new(RwLock::new(PluginSandbox::with_limits(ResourceLimits::default())));

        Self {
            plugins: Arc::new(RwLock::new(HashMap::new())),
            configs: Arc::new(RwLock::new(HashMap::new())),
            loader,
            sandbox,
        }
    }

    /// 初始化运行时
    pub async fn initialize(&self) -> Result<()> {
        log::info!("初始化插件运行时...");

        // 扫描并加载插件
        let plugins = self.loader.scan_plugins()?;
        
        let mut plugins_map = self.plugins.write().await;
        for plugin in plugins {
            let plugin_id = plugin.metadata.id.clone();
            log::info!("发现插件：{} v{}", plugin.metadata.name, plugin.metadata.version);
            plugins_map.insert(plugin_id, plugin);
        }

        log::info!("插件运行时初始化完成，共加载 {} 个插件", plugins_map.len());

        Ok(())
    }

    /// 加载插件
    pub async fn load_plugin(&self, plugin_path: &std::path::Path) -> Result<String> {
        let plugin = self.loader.load_plugin(plugin_path)?;
        let plugin_id = plugin.metadata.id.clone();

        // 验证插件
        plugin.validate(env!("CARGO_PKG_VERSION"))?;

        // 初始化沙箱
        let mut sandbox = self.sandbox.write().await;
        sandbox.initialize(&plugin_id)?;

        // 添加到插件列表
        let mut plugins = self.plugins.write().await;
        plugins.insert(plugin_id.clone(), plugin);

        // 启动沙箱
        sandbox.start(&plugin_id)?;
        drop(sandbox);

        log::info!("加载插件：{}", plugin_id);

        Ok(plugin_id)
    }

    /// 卸载插件
    pub async fn unload_plugin(&self, plugin_id: &str) -> Result<()> {
        // 终止沙箱
        let mut sandbox = self.sandbox.write().await;
        sandbox.terminate(plugin_id)?;
        drop(sandbox);

        // 从列表中移除
        let mut plugins = self.plugins.write().await;
        plugins.remove(plugin_id);

        // 从文件系统删除
        self.loader.unload_plugin(plugin_id)?;

        log::info!("卸载插件：{}", plugin_id);

        Ok(())
    }

    /// 启用插件
    pub async fn enable_plugin(&self, plugin_id: &str) -> Result<()> {
        let mut plugins = self.plugins.write().await;
        let plugin = plugins.get_mut(plugin_id)
            .ok_or_else(|| PluginError::NotFound(plugin_id.to_string()))?;

        plugin.status = PluginStatus::Active;
        
        // 启动沙箱
        let mut sandbox = self.sandbox.write().await;
        sandbox.start(plugin_id)?;
        drop(sandbox);

        // 更新配置
        let mut configs = self.configs.write().await;
        if let Some(config) = configs.get_mut(plugin_id) {
            config.enabled = true;
        }

        log::info!("启用插件：{}", plugin_id);

        Ok(())
    }

    /// 禁用插件
    pub async fn disable_plugin(&self, plugin_id: &str) -> Result<()> {
        let mut plugins = self.plugins.write().await;
        let plugin = plugins.get_mut(plugin_id)
            .ok_or_else(|| PluginError::NotFound(plugin_id.to_string()))?;

        plugin.status = PluginStatus::Inactive;
        
        // 暂停沙箱
        let mut sandbox = self.sandbox.write().await;
        sandbox.suspend(plugin_id)?;
        drop(sandbox);

        // 更新配置
        let mut configs = self.configs.write().await;
        if let Some(config) = configs.get_mut(plugin_id) {
            config.enabled = false;
        }

        log::info!("禁用插件：{}", plugin_id);

        Ok(())
    }

    /// 获取所有插件
    pub async fn get_all_plugins(&self) -> Vec<Plugin> {
        let plugins = self.plugins.read().await;
        plugins.values().cloned().collect()
    }

    /// 获取活跃的插件
    pub async fn get_active_plugins(&self) -> Vec<Plugin> {
        let plugins = self.plugins.read().await;
        plugins.values()
            .filter(|p| p.is_active())
            .cloned()
            .collect()
    }

    /// 获取插件
    pub async fn get_plugin(&self, plugin_id: &str) -> Option<Plugin> {
        let plugins = self.plugins.read().await;
        plugins.get(plugin_id).cloned()
    }

    /// 执行插件
    pub async fn execute_plugin(&self, plugin_id: &str, code: &str) -> Result<String> {
        let plugins = self.plugins.read().await;
        let plugin = plugins.get(plugin_id)
            .ok_or_else(|| PluginError::NotFound(plugin_id.to_string()))?;

        if !plugin.is_active() {
            return Err(PluginError::ExecutionFailed("插件未激活".to_string()));
        }

        // 在沙箱中执行
        let sandbox = self.sandbox.read().await;
        sandbox.execute(plugin_id, code)
    }

    /// 获取插件配置
    pub async fn get_plugin_config(&self, plugin_id: &str) -> PluginConfig {
        let configs = self.configs.read().await;
        configs.get(plugin_id).cloned().unwrap_or_default()
    }

    /// 更新插件配置
    pub async fn update_plugin_config(&self, plugin_id: &str, config: PluginConfig) -> Result<()> {
        let mut configs = self.configs.write().await;
        configs.insert(plugin_id.to_string(), config);
        
        log::info!("更新插件配置：{}", plugin_id);

        Ok(())
    }

    /// 安装插件
    pub async fn install_plugin(&self, zip_path: &std::path::Path) -> Result<String> {
        let plugin_id = self.loader.install_plugin(zip_path)?;
        log::info!("安装插件：{}", plugin_id);
        Ok(plugin_id)
    }

    /// 获取运行时信息
    pub async fn get_runtime_info(&self) -> RuntimeInfo {
        let plugins = self.plugins.read().await;
        let active_count = plugins.values().filter(|p| p.is_active()).count();

        RuntimeInfo {
            total_plugins: plugins.len(),
            active_plugins: active_count,
            inactive_plugins: plugins.len() - active_count,
        }
    }
}

/// 运行时信息
#[derive(Debug, Clone)]
pub struct RuntimeInfo {
    pub total_plugins: usize,
    pub active_plugins: usize,
    pub inactive_plugins: usize,
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::path::PathBuf;

    #[tokio::test]
    async fn test_runtime_initialization() {
        let temp_dir = std::env::temp_dir().join("fg_abyss_test_runtime");
        let _ = std::fs::remove_dir_all(&temp_dir);
        
        let runtime = PluginRuntime::new(temp_dir.clone());
        runtime.initialize().await.unwrap();
        
        let info = runtime.get_runtime_info().await;
        assert_eq!(info.total_plugins, 0);
        assert_eq!(info.active_plugins, 0);
        
        // 清理
        let _ = std::fs::remove_dir_all(&temp_dir);
    }
}
