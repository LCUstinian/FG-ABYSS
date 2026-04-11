// 插件类型定义
export interface Plugin {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  enabled: boolean;
  signature: string;
  path?: string; // 本地插件路径
  isBuiltIn: boolean; // 是否为内置插件
}

// 插件类型枚举
export enum PluginType {
  COMMAND_EXECUTION = 'command_execution',
  FILE_MANAGEMENT = 'file_management',
  DATABASE_MANAGEMENT = 'database_management',
  CUSTOM = 'custom'
}

// 插件状态
export enum PluginStatus {
  INSTALLED = 'installed',
  ENABLED = 'enabled',
  DISABLED = 'disabled',
  ERROR = 'error'
}

// 插件操作请求
export interface PluginActionRequest {
  pluginId: string;
  action: 'enable' | 'disable' | 'uninstall' | 'install';
  path?: string; // 安装时需要
}

// 插件加载结果
export interface PluginLoadResult {
  success: boolean;
  message: string;
  plugin?: Plugin;
}

// 插件签名验证结果
export interface PluginSignatureResult {
  valid: boolean;
  message: string;
  pluginId?: string;
}
