package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"
	"time"

	"fg-abyss/internal/app/services"
)

// PluginLoader 插件加载器
type PluginLoader struct {
	pluginDir      string
	dataDir        string
	loadedPlugins  map[string]Plugin
	externalPlugins map[string]*ExternalPluginInfo
	mu             sync.RWMutex
	api            *PluginAPIImpl
}

// ExternalPluginInfo 外置插件信息
type ExternalPluginInfo struct {
	Path       string         `json:"path"`
	Metadata   *PluginMetadata `json:"metadata"`
	Plugin     *plugin.Plugin  `json:"-"`
	IsLoaded   bool           `json:"is_loaded"`
	IsEnabled  bool           `json:"is_enabled"`
	Error      string         `json:"error,omitempty"`
}

// PluginAPIImpl 插件 API 实现
type PluginAPIImpl struct {
	appVersion      string
	pluginDir       string
	dataDir         string
	config          map[string]interface{}
	commands        map[string]CommandHandler
	menus           []*Menu
	events          map[string][]EventHandler
	eventMu         sync.RWMutex
	permissions     map[string][]Permission
	permissionMu    sync.RWMutex
	connectionService *services.ConnectionService
	fileService       *services.FileService
	databaseService   *services.DatabaseService
}

// NewPluginLoader 创建插件加载器
func NewPluginLoader(pluginDir, dataDir, appVersion string) *PluginLoader {
	return &PluginLoader{
		pluginDir:       pluginDir,
		dataDir:         dataDir,
		loadedPlugins:   make(map[string]Plugin),
		externalPlugins: make(map[string]*ExternalPluginInfo),
		api: &PluginAPIImpl{
			appVersion:        appVersion,
			pluginDir:         pluginDir,
			dataDir:           dataDir,
			config:            make(map[string]interface{}),
			commands:          make(map[string]CommandHandler),
			menus:             make([]*Menu, 0),
			events:            make(map[string][]EventHandler),
			permissions:       make(map[string][]Permission),
			connectionService: services.NewConnectionService(),
			fileService:       services.NewFileService(services.NewConnectionService()),
			databaseService:   services.NewDatabaseService(),
		},
	}
}

// Initialize 初始化插件加载器
func (l *PluginLoader) Initialize(ctx context.Context) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 创建插件目录
	if err := os.MkdirAll(l.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// 创建数据目录
	if err := os.MkdirAll(l.dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// 加载内置插件
	if err := l.loadBuiltinPlugins(ctx); err != nil {
		return fmt.Errorf("failed to load builtin plugins: %w", err)
	}

	// 加载外置插件
	if err := l.loadExternalPlugins(ctx); err != nil {
		return fmt.Errorf("failed to load external plugins: %w", err)
	}

	return nil
}

// loadBuiltinPlugins 加载内置插件
func (l *PluginLoader) loadBuiltinPlugins(ctx context.Context) error {
	// 注册内置插件
	builtinPlugins := []Plugin{
		NewCommandExecutorPlugin(),
		NewFileManagerPlugin(),
		NewDatabaseManagerPlugin(),
	}

	for _, p := range builtinPlugins {
		if err := l.initializePlugin(ctx, p); err != nil {
			return fmt.Errorf("failed to initialize builtin plugin %s: %w", p.GetMetadata().ID, err)
		}
		l.loadedPlugins[p.GetMetadata().ID] = p
	}

	return nil
}

// loadExternalPlugins 加载外置插件
func (l *PluginLoader) loadExternalPlugins(ctx context.Context) error {
	// 扫描插件目录
	entries, err := os.ReadDir(l.pluginDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(l.pluginDir, entry.Name())
		
		// 检查是否是外置插件（.so 或 .dll 文件）
		pluginFile, err := l.findPluginFile(pluginPath)
		if err != nil {
			continue // 跳过无效的插件目录
		}

		// 加载插件
		info, err := l.loadPluginFile(ctx, pluginFile)
		if err != nil {
			info.Error = err.Error()
			l.externalPlugins[entry.Name()] = info
			continue
		}

		l.externalPlugins[entry.Name()] = info
	}

	return nil
}

// findPluginFile 查找插件文件
func (l *PluginLoader) findPluginFile(pluginPath string) (string, error) {
	// 支持的插件文件扩展名
	extensions := []string{".so", ".dll", ".dylib"}
	
	for _, ext := range extensions {
		files, err := filepath.Glob(filepath.Join(pluginPath, "*"+ext))
		if err != nil {
			continue
		}
		if len(files) > 0 {
			return files[0], nil
		}
	}

	return "", fmt.Errorf("no plugin file found in %s", pluginPath)
}

// loadPluginFile 加载插件文件
func (l *PluginLoader) loadPluginFile(ctx context.Context, pluginPath string) (*ExternalPluginInfo, error) {
	// 打开插件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// 查找元数据符号
	metadataSym, err := p.Lookup("PluginMetadata")
	if err != nil {
		return nil, fmt.Errorf("failed to find PluginMetadata: %w", err)
	}

	metadata, ok := metadataSym.(*PluginMetadata)
	if !ok {
		return nil, fmt.Errorf("invalid PluginMetadata type")
	}

	// 查找插件创建函数
	newPluginSym, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("failed to find NewPlugin: %w", err)
	}

	newPluginFunc, ok := newPluginSym.(func() Plugin)
	if !ok {
		return nil, fmt.Errorf("invalid NewPlugin type")
	}

	// 创建插件实例
	pInstance := newPluginFunc()

	info := &ExternalPluginInfo{
		Path:     pluginPath,
		Metadata: metadata,
		Plugin:   p,
		IsLoaded: false,
		IsEnabled: false,
	}

	// 初始化插件
	if err := l.initializePlugin(ctx, pInstance); err != nil {
		return info, fmt.Errorf("failed to initialize plugin: %w", err)
	}

	info.IsLoaded = true
	l.loadedPlugins[metadata.ID] = pInstance

	return info, nil
}

// initializePlugin 初始化插件
func (l *PluginLoader) initializePlugin(ctx context.Context, p Plugin) error {
	if err := p.Initialize(ctx, l.api); err != nil {
		return err
	}

	// 启动插件
	if err := p.Start(ctx); err != nil {
		return err
	}

	return nil
}

// InstallPlugin 安装插件
func (l *PluginLoader) InstallPlugin(ctx context.Context, pluginPath string) (*ExternalPluginInfo, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 验证插件文件
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin file not found: %s", pluginPath)
	}

	// 创建插件目录
	pluginName := filepath.Base(pluginPath)
	pluginDir := filepath.Join(l.pluginDir, pluginName[:len(pluginName)-len(filepath.Ext(pluginName))])
	
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// 复制插件文件
	destPath := filepath.Join(pluginDir, pluginName)
	if err := copyFile(pluginPath, destPath); err != nil {
		return nil, fmt.Errorf("failed to copy plugin file: %w", err)
	}

	// 加载插件
	info, err := l.loadPluginFile(ctx, destPath)
	if err != nil {
		os.RemoveAll(pluginDir)
		return nil, err
	}

	l.externalPlugins[pluginDir] = info

	return info, nil
}

// UninstallPlugin 卸载插件
func (l *PluginLoader) UninstallPlugin(ctx context.Context, pluginID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 查找插件
	var pluginDir string
	var info *ExternalPluginInfo
	for dir, i := range l.externalPlugins {
		if i.Metadata.ID == pluginID {
			pluginDir = dir
			info = i
			break
		}
	}

	if info == nil {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}

	// 停止并卸载插件
	if p, exists := l.loadedPlugins[pluginID]; exists {
		if err := p.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop plugin: %w", err)
		}
		if err := p.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown plugin: %w", err)
		}
		delete(l.loadedPlugins, pluginID)
	}

	// 关闭插件文件
	if info.Plugin != nil {
		// Go 的 plugin 包不支持关闭，依赖 GC
	}

	// 删除插件目录
	if err := os.RemoveAll(pluginDir); err != nil {
		return fmt.Errorf("failed to remove plugin directory: %w", err)
	}

	delete(l.externalPlugins, pluginDir)

	return nil
}

// EnablePlugin 启用插件
func (l *PluginLoader) EnablePlugin(ctx context.Context, pluginID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	p, exists := l.loadedPlugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	if err := p.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugin: %w", err)
	}

	// 更新外置插件状态
	for _, info := range l.externalPlugins {
		if info.Metadata.ID == pluginID {
			info.IsEnabled = true
			break
		}
	}

	return nil
}

// DisablePlugin 禁用插件
func (l *PluginLoader) DisablePlugin(ctx context.Context, pluginID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	p, exists := l.loadedPlugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not loaded: %s", pluginID)
	}

	if err := p.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop plugin: %w", err)
	}

	// 更新外置插件状态
	for _, info := range l.externalPlugins {
		if info.Metadata.ID == pluginID {
			info.IsEnabled = false
			break
		}
	}

	return nil
}

// GetPluginList 获取插件列表
func (l *PluginLoader) GetPluginList() []*PluginInfo {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var result []*PluginInfo

	// 添加内置插件
	for id, p := range l.loadedPlugins {
		if p.GetMetadata().Type == PluginTypeBuiltin {
			result = append(result, &PluginInfo{
				ID:        id,
				Metadata:  p.GetMetadata(),
				Status:    p.GetStatus(),
				IsEnabled: p.GetStatus() == PluginStatusEnabled,
			})
		}
	}

	// 添加外置插件
	for _, info := range l.externalPlugins {
		status := PluginStatusUnloaded
		if info.IsLoaded {
			if info.IsEnabled {
				status = PluginStatusEnabled
			} else {
				status = PluginStatusDisabled
			}
		}

		result = append(result, &PluginInfo{
			ID:        info.Metadata.ID,
			Metadata:  info.Metadata,
			Status:    status,
			IsEnabled: info.IsEnabled,
			Error:     info.Error,
		})
	}

	return result
}

// PluginInfo 插件信息（用于 API 返回）
type PluginInfo struct {
	ID        string         `json:"id"`
	Metadata  *PluginMetadata `json:"metadata"`
	Status    PluginStatus   `json:"status"`
	IsEnabled bool           `json:"is_enabled"`
	Error     string         `json:"error,omitempty"`
}

// GetAPI 获取插件 API
func (l *PluginLoader) GetAPI() *PluginAPIImpl {
	return l.api
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// ========== PluginAPIImpl 实现 ==========

// GetAppVersion 获取应用版本
func (api *PluginAPIImpl) GetAppVersion() string {
	return api.appVersion
}

// GetPluginDir 获取插件目录
func (api *PluginAPIImpl) GetPluginDir() string {
	return api.pluginDir
}

// GetDataDir 获取数据目录
func (api *PluginAPIImpl) GetDataDir() string {
	return api.dataDir
}

// GetConfig 获取配置
func (api *PluginAPIImpl) GetConfig(key string) (interface{}, error) {
	if v, ok := api.config[key]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("config key not found: %s", key)
}

// SetConfig 设置配置
func (api *PluginAPIImpl) SetConfig(key string, value interface{}) error {
	api.config[key] = value
	return nil
}

// RegisterCommand 注册命令
func (api *PluginAPIImpl) RegisterCommand(name string, handler CommandHandler) error {
	api.commands[name] = handler
	return nil
}

// RegisterMenu 注册菜单
func (api *PluginAPIImpl) RegisterMenu(parentID string, menu *Menu) error {
	api.menus = append(api.menus, menu)
	return nil
}

// EmitEvent 触发事件
func (api *PluginAPIImpl) EmitEvent(event string, data interface{}) error {
	api.eventMu.RLock()
	defer api.eventMu.RUnlock()

	if handlers, ok := api.events[event]; ok {
		for _, handler := range handlers {
			if err := handler(event, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// SubscribeEvent 订阅事件
func (api *PluginAPIImpl) SubscribeEvent(event string, handler EventHandler) error {
	api.eventMu.Lock()
	defer api.eventMu.Unlock()

	api.events[event] = append(api.events[event], handler)
	return nil
}

// Log 日志记录
func (api *PluginAPIImpl) Log(level LogLevel, message string, fields map[string]interface{}) {
	// 这里应该集成到应用的日志系统
	logData := map[string]interface{}{
		"level":   level,
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}
	
	for k, v := range fields {
		logData[k] = v
	}

	jsonData, _ := json.Marshal(logData)
	fmt.Printf("[PLUGIN LOG] %s\n", string(jsonData))
}

// HasPermission 检查权限
func (api *PluginAPIImpl) HasPermission(permission Permission) bool {
	api.permissionMu.RLock()
	defer api.permissionMu.RUnlock()

	// TODO: 实现权限检查逻辑
	return true
}

// CallPlugin 调用其他插件
func (api *PluginAPIImpl) CallPlugin(pluginID string, method string, params interface{}) (interface{}, error) {
	// TODO: 实现插件间调用
	return nil, fmt.Errorf("plugin call not implemented")
}

// GrantPermission 授予权限
func (api *PluginAPIImpl) GrantPermission(pluginID string, permissions []Permission) {
	api.permissionMu.Lock()
	defer api.permissionMu.Unlock()
	api.permissions[pluginID] = permissions
}

// RevokePermission 撤销权限
func (api *PluginAPIImpl) RevokePermission(pluginID string) {
	api.permissionMu.Lock()
	defer api.permissionMu.Unlock()
	delete(api.permissions, pluginID)
}

// GetCommands 获取所有注册的命令
func (api *PluginAPIImpl) GetCommands() map[string]CommandHandler {
	return api.commands
}

// GetMenus 获取所有注册的菜单
func (api *PluginAPIImpl) GetMenus() []*Menu {
	return api.menus
}

// GetConnectionService 获取连接服务
func (api *PluginAPIImpl) GetConnectionService() *services.ConnectionService {
	return api.connectionService
}

// GetFileService 获取文件服务
func (api *PluginAPIImpl) GetFileService() *services.FileService {
	return api.fileService
}

// GetDatabaseService 获取数据库服务
func (api *PluginAPIImpl) GetDatabaseService() *services.DatabaseService {
	return api.databaseService
}
