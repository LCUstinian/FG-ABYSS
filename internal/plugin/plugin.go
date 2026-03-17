package plugin

import (
	"context"
	"time"
)

// PluginType 插件类型
type PluginType string

const (
	// PluginTypeBuiltin 内置插件
	PluginTypeBuiltin PluginType = "builtin"
	// PluginTypeExternal 外置插件
	PluginTypeExternal PluginType = "external"
)

// PluginStatus 插件状态
type PluginStatus string

const (
	// PluginStatusLoaded 已加载
	PluginStatusLoaded PluginStatus = "loaded"
	// PluginStatusUnloaded 未加载
	PluginStatusUnloaded PluginStatus = "unloaded"
	// PluginStatusEnabled 已启用
	PluginStatusEnabled PluginStatus = "enabled"
	// PluginStatusDisabled 已禁用
	PluginStatusDisabled PluginStatus = "disabled"
	// PluginStatusError 错误状态
	PluginStatusError PluginStatus = "error"
)

// PluginMetadata 插件元数据
type PluginMetadata struct {
	// ID 插件唯一标识
	ID string `json:"id"`
	// Name 插件名称
	Name string `json:"name"`
	// Version 版本号 (语义化版本：MAJOR.MINOR.PATCH)
	Version string `json:"version"`
	// Description 插件描述
	Description string `json:"description"`
	// Author 作者信息
	Author string `json:"author"`
	// Type 插件类型
	Type PluginType `json:"type"`
	// Category 插件分类
	Category string `json:"category"`
	// Tags 标签
	Tags []string `json:"tags,omitempty"`
	// Homepage 主页 URL
	Homepage string `json:"homepage,omitempty"`
	// Repository 代码仓库 URL
	Repository string `json:"repository,omitempty"`
	// License 许可证
	License string `json:"license,omitempty"`
	// MinAppVersion 最低应用版本要求
	MinAppVersion string `json:"min_app_version,omitempty"`
	// Dependencies 依赖插件列表
	Dependencies []PluginDependency `json:"dependencies,omitempty"`
	// Permissions 权限声明
	Permissions []Permission `json:"permissions,omitempty"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// PluginDependency 插件依赖
type PluginDependency struct {
	// ID 依赖插件 ID
	ID string `json:"id"`
	// Version 版本约束 (如："^1.0.0", ">=2.0.0", "~1.2.0")
	Version string `json:"version"`
	// Required 是否必需
	Required bool `json:"required"`
}

// Permission 权限定义
type Permission string

const (
	// PermissionFileSystem 文件系统访问
	PermissionFileSystem Permission = "fs:readwrite"
	// PermissionNetwork 网络访问
	PermissionNetwork Permission = "network"
	// PermissionDatabase 数据库访问
	PermissionDatabase Permission = "database"
	// PermissionExecute 命令执行
	PermissionExecute Permission = "execute"
	// PermissionClipboard 剪贴板访问
	PermissionClipboard Permission = "clipboard"
	// PermissionSettings 设置访问
	PermissionSettings Permission = "settings"
	// PermissionPlugin 插件管理
	PermissionPlugin Permission = "plugin"
)

// PluginLifecycle 插件生命周期接口
type PluginLifecycle interface {
	// Initialize 初始化插件
	Initialize(ctx context.Context, api PluginAPI) error
	// Start 启动插件
	Start(ctx context.Context) error
	// Stop 停止插件
	Stop(ctx context.Context) error
	// Shutdown 关闭插件（清理资源）
	Shutdown(ctx context.Context) error
}

// PluginFunction 插件功能接口
type PluginFunction interface {
	// GetMetadata 获取插件元数据
	GetMetadata() *PluginMetadata
	// GetStatus 获取插件状态
	GetStatus() PluginStatus
	// SetStatus 设置插件状态
	SetStatus(status PluginStatus)
}

// Plugin 插件完整接口
type Plugin interface {
	PluginLifecycle
	PluginFunction
}

// PluginAPI 插件 API 接口（提供给插件使用的 API）
type PluginAPI interface {
	// GetAppVersion 获取应用版本
	GetAppVersion() string
	// GetPluginDir 获取插件目录
	GetPluginDir() string
	// GetDataDir 获取数据目录
	GetDataDir() string
	// GetConfig 获取配置
	GetConfig(key string) (interface{}, error)
	// SetConfig 设置配置
	SetConfig(key string, value interface{}) error
	// RegisterCommand 注册命令
	RegisterCommand(name string, handler CommandHandler) error
	// RegisterMenu 注册菜单
	RegisterMenu(parentID string, menu *Menu) error
	// EmitEvent 触发事件
	EmitEvent(event string, data interface{}) error
	// SubscribeEvent 订阅事件
	SubscribeEvent(event string, handler EventHandler) error
	// Log 日志记录
	Log(level LogLevel, message string, fields map[string]interface{})
	// HasPermission 检查权限
	HasPermission(permission Permission) bool
	// CallPlugin 调用其他插件
	CallPlugin(pluginID string, method string, params interface{}) (interface{}, error)
}

// CommandHandler 命令处理器
type CommandHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// EventHandler 事件处理器
type EventHandler func(event string, data interface{}) error

// Menu 菜单定义
type Menu struct {
	ID       string     `json:"id"`
	Label    string     `json:"label"`
	Icon     string     `json:"icon,omitempty"`
	Action   string     `json:"action"`
	Children []*Menu    `json:"children,omitempty"`
	Enabled  bool       `json:"enabled"`
	Visible  bool       `json:"visible"`
	Shortcut string     `json:"shortcut,omitempty"`
}

// LogLevel 日志级别
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// BasePlugin 插件基础实现
type BasePlugin struct {
	metadata *PluginMetadata
	status   PluginStatus
	api      PluginAPI
}

// NewBasePlugin 创建基础插件
func NewBasePlugin(metadata *PluginMetadata) *BasePlugin {
	return &BasePlugin{
		metadata: metadata,
		status:   PluginStatusUnloaded,
	}
}

// GetMetadata 获取插件元数据
func (p *BasePlugin) GetMetadata() *PluginMetadata {
	return p.metadata
}

// GetStatus 获取插件状态
func (p *BasePlugin) GetStatus() PluginStatus {
	return p.status
}

// SetStatus 设置插件状态
func (p *BasePlugin) SetStatus(status PluginStatus) {
	p.status = status
}

// Initialize 初始化（默认实现）
func (p *BasePlugin) Initialize(ctx context.Context, api PluginAPI) error {
	p.api = api
	p.status = PluginStatusLoaded
	return nil
}

// Start 启动（默认实现）
func (p *BasePlugin) Start(ctx context.Context) error {
	p.status = PluginStatusEnabled
	return nil
}

// Stop 停止（默认实现）
func (p *BasePlugin) Stop(ctx context.Context) error {
	p.status = PluginStatusDisabled
	return nil
}

// Shutdown 关闭（默认实现）
func (p *BasePlugin) Shutdown(ctx context.Context) error {
	p.status = PluginStatusUnloaded
	p.api = nil
	return nil
}
