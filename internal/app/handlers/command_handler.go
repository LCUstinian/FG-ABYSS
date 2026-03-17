package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
)

// CommandHandler 命令执行处理器（Wails 绑定用）
type CommandHandler struct {
	service *services.CommandService
}

// NewCommandHandler 创建命令执行处理器
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		service: services.NewCommandService(services.NewConnectionService()),
	}
}

// ExecuteCommand 执行命令
// @summary 执行命令
// @param ctx 上下文
// @param webshellID WebShell ID
// @param command 要执行的命令
// @param timeout 超时时间（秒）
// @return 命令执行结果
func (h *CommandHandler) ExecuteCommand(
	ctx context.Context,
	webshellID string,
	command string,
	timeout int,
) (*CommandResult, error) {
	req := &services.CommandRequest{
		WebShellID: webshellID,
		Command:    command,
		Timeout:    timeout,
	}
	
	resp, err := h.service.Execute(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &CommandResult{
		Output:   resp.Output,
		Error:    resp.Error,
		ExitCode: resp.ExitCode,
		Duration: resp.Duration,
		Success:  resp.Success,
	}, nil
}

// ExecuteCommandWithResult 执行命令并返回结构化结果
// @summary 执行命令并返回结构化结果
// @param ctx 上下文
// @param webshellID WebShell ID
// @param command 要执行的命令
// @return 结构化结果
func (h *CommandHandler) ExecuteCommandWithResult(
	ctx context.Context,
	webshellID string,
	command string,
) (map[string]interface{}, error) {
	return h.service.ExecuteWithResult(ctx, webshellID, command)
}

// GetSystemInfo 获取系统信息
// @summary 获取系统信息
// @param ctx 上下文
// @param webshellID WebShell ID
// @return 系统信息
func (h *CommandHandler) GetSystemInfo(
	ctx context.Context,
	webshellID string,
) (map[string]interface{}, error) {
	return h.service.GetSystemInfo(ctx, webshellID)
}

// GetCurrentUser 获取当前用户
// @summary 获取当前用户
// @param ctx 上下文
// @param webshellID WebShell ID
// @return 用户名
func (h *CommandHandler) GetCurrentUser(
	ctx context.Context,
	webshellID string,
) (string, error) {
	return h.service.GetCurrentUser(ctx, webshellID)
}

// GetWorkingDirectory 获取工作目录
// @summary 获取工作目录
// @param ctx 上下文
// @param webshellID WebShell ID
// @return 工作目录路径
func (h *CommandHandler) GetWorkingDirectory(
	ctx context.Context,
	webshellID string,
) (string, error) {
	return h.service.GetWorkingDirectory(ctx, webshellID)
}

// KillProcess 终止进程
// @summary 终止进程
// @param ctx 上下文
// @param webshellID WebShell ID
// @param pid 进程 ID
// @return 错误信息
func (h *CommandHandler) KillProcess(
	ctx context.Context,
	webshellID string,
	pid int,
) error {
	return h.service.KillProcess(ctx, webshellID, pid)
}

// CommandResult 命令执行结果
type CommandResult struct {
	Output   string `json:"output"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
	Duration int64  `json:"duration"`
	Success  bool   `json:"success"`
}
