package plugin

import (
	"context"
	"runtime"
	"time"

	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// CommandExecutorPlugin 命令执行插件
type CommandExecutorPlugin struct {
	*BasePlugin
	api *PluginAPIImpl
}

// NewCommandExecutorPlugin 创建命令执行插件
func NewCommandExecutorPlugin() *CommandExecutorPlugin {
	return &CommandExecutorPlugin{
		BasePlugin: NewBasePlugin(&PluginMetadata{
			ID:          "builtin.command_executor",
			Name:        "命令执行",
			Version:     "1.0.0",
			Description: "提供远程命令执行功能",
			Author:      "FG-ABYSS Team",
			Type:        PluginTypeBuiltin,
			Category:    "core",
			Tags:        []string{"command", "execute", "shell"},
			Permissions: []Permission{PermissionExecute},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}),
	}
}

// Initialize 初始化插件
func (p *CommandExecutorPlugin) Initialize(ctx context.Context, api PluginAPI) error {
	p.api = api.(*PluginAPIImpl)
	return p.BasePlugin.Initialize(ctx, api)
}

// ExecuteCommand 执行命令
func (p *CommandExecutorPlugin) ExecuteCommand(ctx context.Context, webshellID, command string, timeout int) (string, error) {
	if !p.api.HasPermission(PermissionExecute) {
		return "", ErrPermissionDenied
	}

	cmdService := services.NewCommandService(p.api.GetConnectionService())
	
	req := &services.CommandRequest{
		WebShellID: webshellID,
		Command:    command,
		Timeout:    timeout,
	}

	resp, err := cmdService.Execute(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Output, nil
}

// GetSystemInfo 获取系统信息
func (p *CommandExecutorPlugin) GetSystemInfo(ctx context.Context, webshellID string) (map[string]interface{}, error) {
	cmdService := services.NewCommandService(p.api.GetConnectionService())
	return cmdService.GetSystemInfo(ctx, webshellID)
}

// GetCurrentUser 获取当前用户
func (p *CommandExecutorPlugin) GetCurrentUser(ctx context.Context, webshellID string) (string, error) {
	cmdService := services.NewCommandService(p.api.GetConnectionService())
	return cmdService.GetCurrentUser(ctx, webshellID)
}

// GetWorkingDirectory 获取工作目录
func (p *CommandExecutorPlugin) GetWorkingDirectory(ctx context.Context, webshellID string) (string, error) {
	cmdService := services.NewCommandService(p.api.GetConnectionService())
	return cmdService.GetWorkingDirectory(ctx, webshellID)
}

// FileManagerPlugin 文件管理插件
type FileManagerPlugin struct {
	*BasePlugin
	api *PluginAPIImpl
}

// NewFileManagerPlugin 创建文件管理插件
func NewFileManagerPlugin() *FileManagerPlugin {
	return &FileManagerPlugin{
		BasePlugin: NewBasePlugin(&PluginMetadata{
			ID:          "builtin.file_manager",
			Name:        "文件管理",
			Version:     "1.0.0",
			Description: "提供远程文件管理功能",
			Author:      "FG-ABYSS Team",
			Type:        PluginTypeBuiltin,
			Category:    "core",
			Tags:        []string{"file", "manager", "upload", "download"},
			Permissions: []Permission{PermissionFileSystem},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}),
	}
}

// Initialize 初始化插件
func (p *FileManagerPlugin) Initialize(ctx context.Context, api PluginAPI) error {
	p.api = api.(*PluginAPIImpl)
	return p.BasePlugin.Initialize(ctx, api)
}

// ListFiles 列出目录内容
func (p *FileManagerPlugin) ListFiles(ctx context.Context, webshellID, path string) (*services.FileListResponse, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return nil, ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	
	req := &services.FileListRequest{
		WebShellID: webshellID,
		Path:       path,
	}

	return fileService.ListFiles(ctx, req)
}

// ReadFile 读取文件
func (p *FileManagerPlugin) ReadFile(ctx context.Context, webshellID, path string, limit int) (string, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return "", ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	
	req := &services.FileReadRequest{
		WebShellID: webshellID,
		Path:       path,
		Limit:      limit,
	}

	return fileService.ReadFile(ctx, req)
}

// WriteFile 写入文件
func (p *FileManagerPlugin) WriteFile(ctx context.Context, webshellID, path, content string, isBase64 bool) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	
	req := &services.FileWriteRequest{
		WebShellID: webshellID,
		Path:       path,
		Content:    content,
		IsBase64:   isBase64,
	}

	return fileService.WriteFile(ctx, req)
}

// UploadFile 上传文件
func (p *FileManagerPlugin) UploadFile(ctx context.Context, webshellID, path, filename, content string) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	
	req := &services.FileUploadRequest{
		WebShellID: webshellID,
		Path:       path,
		Filename:   filename,
		Content:    content,
	}

	return fileService.UploadFile(ctx, req)
}

// DownloadFile 下载文件
func (p *FileManagerPlugin) DownloadFile(ctx context.Context, webshellID, path string) ([]byte, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return nil, ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	
	req := &services.FileDownloadRequest{
		WebShellID: webshellID,
		Path:       path,
	}

	return fileService.DownloadFile(ctx, req)
}

// DeleteFile 删除文件
func (p *FileManagerPlugin) DeleteFile(ctx context.Context, webshellID, path string) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.DeleteFile(ctx, webshellID, path)
}

// CreateDirectory 创建目录
func (p *FileManagerPlugin) CreateDirectory(ctx context.Context, webshellID, path string) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.CreateDirectory(ctx, webshellID, path)
}

// RenameFile 重命名文件
func (p *FileManagerPlugin) RenameFile(ctx context.Context, webshellID, oldPath, newPath string) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.RenameFile(ctx, webshellID, oldPath, newPath)
}

// GetFilePermissions 获取文件权限
func (p *FileManagerPlugin) GetFilePermissions(ctx context.Context, webshellID, path string) (string, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return "", ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.GetFilePermissions(ctx, webshellID, path)
}

// SetFilePermissions 设置文件权限
func (p *FileManagerPlugin) SetFilePermissions(ctx context.Context, webshellID, path, permissions string) error {
	if !p.api.HasPermission(PermissionFileSystem) {
		return ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.SetFilePermissions(ctx, webshellID, path, permissions)
}

// SearchFiles 搜索文件
func (p *FileManagerPlugin) SearchFiles(ctx context.Context, webshellID, path, pattern string) ([]services.FileInfo, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return nil, ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.SearchFiles(ctx, webshellID, path, pattern)
}

// GetDiskUsage 获取磁盘使用情况
func (p *FileManagerPlugin) GetDiskUsage(ctx context.Context, webshellID string) (map[string]interface{}, error) {
	if !p.api.HasPermission(PermissionFileSystem) {
		return nil, ErrPermissionDenied
	}

	fileService := services.NewFileService(p.api.GetConnectionService())
	return fileService.GetDiskUsage(ctx, webshellID)
}

// DatabaseManagerPlugin 数据库管理插件
type DatabaseManagerPlugin struct {
	*BasePlugin
	api *PluginAPIImpl
}

// NewDatabaseManagerPlugin 创建数据库管理插件
func NewDatabaseManagerPlugin() *DatabaseManagerPlugin {
	return &DatabaseManagerPlugin{
		BasePlugin: NewBasePlugin(&PluginMetadata{
			ID:          "builtin.database_manager",
			Name:        "数据库管理",
			Version:     "1.0.0",
			Description: "提供数据库连接和管理功能",
			Author:      "FG-ABYSS Team",
			Type:        PluginTypeBuiltin,
			Category:    "core",
			Tags:        []string{"database", "mysql", "postgresql", "sqlite", "mssql"},
			Permissions: []Permission{PermissionDatabase},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}),
	}
}

// Initialize 初始化插件
func (p *DatabaseManagerPlugin) Initialize(ctx context.Context, api PluginAPI) error {
	p.api = api.(*PluginAPIImpl)
	return p.BasePlugin.Initialize(ctx, api)
}

// ConnectDatabase 连接数据库
func (p *DatabaseManagerPlugin) ConnectDatabase(ctx context.Context, conn *entity.DatabaseConnection) error {
	if !p.api.HasPermission(PermissionDatabase) {
		return ErrPermissionDenied
	}

	dbService := services.NewDatabaseService()
	return dbService.Connect(conn)
}

// ExecuteQuery 执行查询
func (p *DatabaseManagerPlugin) ExecuteQuery(ctx context.Context, connectionID int64, query string) (*entity.DatabaseQueryResult, error) {
	if !p.api.HasPermission(PermissionDatabase) {
		return nil, ErrPermissionDenied
	}

	dbService := services.NewDatabaseService()
	
	q := entity.DatabaseQuery{
		SQL: query,
	}

	return dbService.ExecuteQuery(connectionID, q)
}

// GetTables 获取表列表
func (p *DatabaseManagerPlugin) GetTables(ctx context.Context, connectionID int64, dbType entity.DatabaseType) ([]string, error) {
	if !p.api.HasPermission(PermissionDatabase) {
		return nil, ErrPermissionDenied
	}

	dbService := services.NewDatabaseService()
	tables, err := dbService.GetTables(connectionID, dbType)
	if err != nil {
		return nil, err
	}
	
	// 转换为字符串数组
	result := make([]string, len(tables))
	for i, t := range tables {
		result[i] = t.Name
	}
	return result, nil
}

// GetTableData 获取表数据
func (p *DatabaseManagerPlugin) GetTableData(ctx context.Context, connectionID int64, tableName string, limit, offset int) (*entity.DatabaseQueryResult, error) {
	if !p.api.HasPermission(PermissionDatabase) {
		return nil, ErrPermissionDenied
	}

	query := "SELECT * FROM " + tableName
	if limit > 0 {
		query += " LIMIT ?"
	}
	if offset > 0 {
		query += " OFFSET ?"
	}

	dbService := services.NewDatabaseService()
	
	q := entity.DatabaseQuery{
		SQL:    query,
		Limit:  limit,
		Offset: offset,
	}

	return dbService.ExecuteQuery(connectionID, q)
}

// CloseConnection 关闭连接
func (p *DatabaseManagerPlugin) CloseConnection(ctx context.Context, connectionID int64) error {
	if !p.api.HasPermission(PermissionDatabase) {
		return ErrPermissionDenied
	}

	dbService := services.NewDatabaseService()
	return dbService.Disconnect(connectionID)
}

// ErrPermissionDenied 权限拒绝错误
var ErrPermissionDenied = &PermissionError{"permission denied"}

// PermissionError 权限错误
type PermissionError struct {
	Message string
}

func (e *PermissionError) Error() string {
	return e.Message
}

// GetOSInfo 获取操作系统信息（辅助函数）
func GetOSInfo() map[string]interface{} {
	return map[string]interface{}{
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"num_cpu":  runtime.NumCPU(),
		"go_version": runtime.Version(),
	}
}
