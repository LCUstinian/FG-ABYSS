package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
)

// FileHandler 文件管理处理器（Wails 绑定用）
type FileHandler struct {
	service *services.FileService
}

// NewFileHandler 创建文件管理处理器
func NewFileHandler() *FileHandler {
	return &FileHandler{
		service: services.NewFileService(services.NewConnectionService()),
	}
}

// ListFiles 列出目录内容
// @summary 列出目录内容
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 目录路径
// @return 文件列表
func (h *FileHandler) ListFiles(
	ctx context.Context,
	webshellID string,
	path string,
) (*services.FileListResponse, error) {
	req := &services.FileListRequest{
		WebShellID: webshellID,
		Path:       path,
	}
	
	return h.service.ListFiles(ctx, req)
}

// ReadFile 读取文件
// @summary 读取文件内容
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @param limit 读取限制（字节）
// @return 文件内容
func (h *FileHandler) ReadFile(
	ctx context.Context,
	webshellID string,
	path string,
	limit int,
) (string, error) {
	req := &services.FileReadRequest{
		WebShellID: webshellID,
		Path:       path,
		Limit:      limit,
	}
	
	return h.service.ReadFile(ctx, req)
}

// WriteFile 写入文件
// @summary 写入文件内容
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @param content 文件内容
// @param isBase64 是否为 Base64 编码
// @return 错误信息
func (h *FileHandler) WriteFile(
	ctx context.Context,
	webshellID string,
	path string,
	content string,
	isBase64 bool,
) error {
	req := &services.FileWriteRequest{
		WebShellID: webshellID,
		Path:       path,
		Content:    content,
		IsBase64:   isBase64,
	}
	
	return h.service.WriteFile(ctx, req)
}

// UploadFile 上传文件
// @summary 上传文件
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 上传路径
// @param filename 文件名
// @param content 文件内容（Base64）
// @return 错误信息
func (h *FileHandler) UploadFile(
	ctx context.Context,
	webshellID string,
	path string,
	filename string,
	content string,
) error {
	req := &services.FileUploadRequest{
		WebShellID: webshellID,
		Path:       path,
		Filename:   filename,
		Content:    content,
	}
	
	return h.service.UploadFile(ctx, req)
}

// DownloadFile 下载文件
// @summary 下载文件
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @return 文件内容（字节）
func (h *FileHandler) DownloadFile(
	ctx context.Context,
	webshellID string,
	path string,
) ([]byte, error) {
	req := &services.FileDownloadRequest{
		WebShellID: webshellID,
		Path:       path,
	}
	
	return h.service.DownloadFile(ctx, req)
}

// DeleteFile 删除文件
// @summary 删除文件
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @return 错误信息
func (h *FileHandler) DeleteFile(
	ctx context.Context,
	webshellID string,
	path string,
) error {
	return h.service.DeleteFile(ctx, webshellID, path)
}

// CreateDirectory 创建目录
// @summary 创建目录
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 目录路径
// @return 错误信息
func (h *FileHandler) CreateDirectory(
	ctx context.Context,
	webshellID string,
	path string,
) error {
	return h.service.CreateDirectory(ctx, webshellID, path)
}

// RenameFile 重命名文件
// @summary 重命名文件/目录
// @param ctx 上下文
// @param webshellID WebShell ID
// @param oldPath 原路径
// @param newPath 新路径
// @return 错误信息
func (h *FileHandler) RenameFile(
	ctx context.Context,
	webshellID string,
	oldPath string,
	newPath string,
) error {
	return h.service.RenameFile(ctx, webshellID, oldPath, newPath)
}

// SearchFiles 搜索文件
// @summary 搜索文件
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 搜索路径
// @param pattern 搜索模式
// @return 文件列表
func (h *FileHandler) SearchFiles(
	ctx context.Context,
	webshellID string,
	path string,
	pattern string,
) ([]services.FileInfo, error) {
	return h.service.SearchFiles(ctx, webshellID, path, pattern)
}

// GetFilePermissions 获取文件权限
// @summary 获取文件权限
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @return 权限字符串
func (h *FileHandler) GetFilePermissions(
	ctx context.Context,
	webshellID string,
	path string,
) (string, error) {
	return h.service.GetFilePermissions(ctx, webshellID, path)
}

// SetFilePermissions 设置文件权限
// @summary 设置文件权限
// @param ctx 上下文
// @param webshellID WebShell ID
// @param path 文件路径
// @param permissions 权限（如 755）
// @return 错误信息
func (h *FileHandler) SetFilePermissions(
	ctx context.Context,
	webshellID string,
	path string,
	permissions string,
) error {
	return h.service.SetFilePermissions(ctx, webshellID, path, permissions)
}

// GetDiskUsage 获取磁盘使用情况
// @summary 获取磁盘使用情况
// @param ctx 上下文
// @param webshellID WebShell ID
// @return 磁盘使用信息
func (h *FileHandler) GetDiskUsage(
	ctx context.Context,
	webshellID string,
) (map[string]interface{}, error) {
	return h.service.GetDiskUsage(ctx, webshellID)
}
