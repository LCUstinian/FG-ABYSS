package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// AuditLevel 审计级别
type AuditLevel string

const (
	AuditInfo     AuditLevel = "INFO"
	AuditWarning  AuditLevel = "WARNING"
	AuditError    AuditLevel = "ERROR"
	AuditCritical AuditLevel = "CRITICAL"
)

// AuditOperation 审计操作类型
type AuditOperation string

const (
	OpLogin          AuditOperation = "LOGIN"
	OpLogout         AuditOperation = "LOGOUT"
	OpCreate         AuditOperation = "CREATE"
	OpUpdate         AuditOperation = "UPDATE"
	OpDelete         AuditOperation = "DELETE"
	OpQuery          AuditOperation = "QUERY"
	OpExecute        AuditOperation = "EXECUTE"
	OpUpload         AuditOperation = "UPLOAD"
	OpDownload       AuditOperation = "DOWNLOAD"
	OpConnect        AuditOperation = "CONNECT"
	OpDisconnect     AuditOperation = "DISCONNECT"
	OpEncrypt        AuditOperation = "ENCRYPT"
	OpDecrypt        AuditOperation = "DECRYPT"
	OpGenerate       AuditOperation = "GENERATE"
	OpImport         AuditOperation = "IMPORT"
	OpExport         AuditOperation = "EXPORT"
)

// AuditLog 审计日志条目
type AuditLog struct {
	ID         int64          `json:"id"`
	Timestamp  time.Time      `json:"timestamp"`
	Level      AuditLevel     `json:"level"`
	Operation  AuditOperation `json:"operation"`
	User       string         `json:"user"`
	Resource   string         `json:"resource"`
	Action     string         `json:"action"`
	IPAddress  string         `json:"ip_address"`
	UserAgent  string         `json:"user_agent"`
	Details    string         `json:"details,omitempty"`
	Success    bool           `json:"success"`
	Duration   int64          `json:"duration_ms"`
}

// AuditService 审计日志服务
type AuditService struct {
	logs       []*AuditLog
	mu         sync.RWMutex
	maxSize    int
	filePath   string
	autoSave   bool
	saveTicker *time.Ticker
	stopChan   chan struct{}
}

// AuditConfig 审计配置
type AuditConfig struct {
	MaxSize  int    `json:"max_size"`  // 最大日志条数
	FilePath string `json:"file_path"` // 日志文件路径
	AutoSave bool   `json:"auto_save"` // 是否自动保存
}

// NewAuditService 创建审计服务
func NewAuditService(config *AuditConfig) *AuditService {
	service := &AuditService{
		logs:     make([]*AuditLog, 0),
		maxSize:  config.MaxSize,
		filePath: config.FilePath,
		autoSave: config.AutoSave,
		stopChan: make(chan struct{}),
	}

	// 确保目录存在
	if config.AutoSave && config.FilePath != "" {
		dir := filepath.Dir(config.FilePath)
		os.MkdirAll(dir, 0755)
	}

	return service
}

// Log 记录审计日志
func (a *AuditService) Log(level AuditLevel, operation AuditOperation, user string, resource string, action string, details string, success bool, duration int64) {
	a.mu.Lock()
	defer a.mu.Unlock()

	log := &AuditLog{
		ID:        time.Now().UnixNano(),
		Timestamp: time.Now(),
		Level:     level,
		Operation: operation,
		User:      user,
		Resource:  resource,
		Action:    action,
		Details:   details,
		Success:   success,
		Duration:  duration,
	}

	a.logs = append(a.logs, log)

	// 限制日志大小
	if len(a.logs) > a.maxSize {
		a.logs = a.logs[len(a.logs)-a.maxSize:]
	}

	// 自动保存
	if a.autoSave {
		go a.saveToFile()
	}
}

// Info 记录信息级别日志
func (a *AuditService) Info(operation AuditOperation, user string, resource string, action string) {
	a.Log(AuditInfo, operation, user, resource, action, "", true, 0)
}

// InfoWithDetails 记录信息级别日志（带详情）
func (a *AuditService) InfoWithDetails(operation AuditOperation, user string, resource string, action string, details string) {
	a.Log(AuditInfo, operation, user, resource, action, details, true, 0)
}

// Warning 记录警告级别日志
func (a *AuditService) Warning(operation AuditOperation, user string, resource string, action string, details string) {
	a.Log(AuditWarning, operation, user, resource, action, details, true, 0)
}

// Error 记录错误级别日志
func (a *AuditService) Error(operation AuditOperation, user string, resource string, action string, err error, duration int64) {
	details := ""
	if err != nil {
		details = err.Error()
	}
	a.Log(AuditError, operation, user, resource, action, details, false, duration)
}

// GetLogs 获取日志
func (a *AuditService) GetLogs(limit int, offset int) []*AuditLog {
	a.mu.RLock()
	defer a.mu.RUnlock()

	start := len(a.logs) - offset - limit
	if start < 0 {
		start = 0
	}

	end := len(a.logs) - offset
	if end > len(a.logs) {
		end = len(a.logs)
	}

	if start >= end {
		return []*AuditLog{}
	}

	result := make([]*AuditLog, end-start)
	copy(result, a.logs[start:end])
	return result
}

// GetLogsByUser 按用户获取日志
func (a *AuditService) GetLogsByUser(user string, limit int) []*AuditLog {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(a.logs) - 1; i >= 0 && count < limit; i-- {
		if a.logs[i].User == user {
			result = append(result, a.logs[i])
			count++
		}
	}
	return result
}

// GetLogsByOperation 按操作类型获取日志
func (a *AuditService) GetLogsByOperation(operation AuditOperation, limit int) []*AuditLog {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(a.logs) - 1; i >= 0 && count < limit; i-- {
		if a.logs[i].Operation == operation {
			result = append(result, a.logs[i])
			count++
		}
	}
	return result
}

// GetLogsByLevel 按级别获取日志
func (a *AuditService) GetLogsByLevel(level AuditLevel, limit int) []*AuditLog {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(a.logs) - 1; i >= 0 && count < limit; i-- {
		if a.logs[i].Level == level {
			result = append(result, a.logs[i])
			count++
		}
	}
	return result
}

// Clear 清空日志
func (a *AuditService) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logs = make([]*AuditLog, 0)
}

// Count 获取日志数量
func (a *AuditService) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.logs)
}

// saveToFile 保存到文件
func (a *AuditService) saveToFile() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	data, err := json.MarshalIndent(a.logs, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal audit logs: %v\n", err)
		return
	}

	err = os.WriteFile(a.filePath, data, 0644)
	if err != nil {
		fmt.Printf("Failed to save audit logs: %v\n", err)
	}
}

// ExportToJSON 导出为 JSON
func (a *AuditService) ExportToJSON(filePath string) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	data, err := json.MarshalIndent(a.logs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// ImportFromJSON 从 JSON 导入
func (a *AuditService) ImportFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var logs []*AuditLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	a.logs = append(a.logs, logs...)

	// 限制大小
	if len(a.logs) > a.maxSize {
		a.logs = a.logs[len(a.logs)-a.maxSize:]
	}

	return nil
}

// StartAutoSave 启动自动保存
func (a *AuditService) StartAutoSave(interval time.Duration) {
	if !a.autoSave {
		return
	}

	a.saveTicker = time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-a.saveTicker.C:
				a.saveToFile()
			case <-a.stopChan:
				a.saveTicker.Stop()
				return
			}
		}
	}()
}

// Stop 停止服务
func (a *AuditService) Stop() {
	close(a.stopChan)
	if a.saveTicker != nil {
		a.saveTicker.Stop()
	}
	// 最后保存一次
	a.saveToFile()
}

// GetStatistics 获取统计信息
func (a *AuditService) GetStatistics() *AuditStatistics {
	a.mu.RLock()
	defer a.mu.RUnlock()

	stats := &AuditStatistics{
		TotalLogs: len(a.logs),
	}

	for _, log := range a.logs {
		switch log.Level {
		case AuditInfo:
			stats.InfoCount++
		case AuditWarning:
			stats.WarningCount++
		case AuditError:
			stats.ErrorCount++
		case AuditCritical:
			stats.CriticalCount++
		}

		if log.Success {
			stats.SuccessCount++
		} else {
			stats.FailureCount++
		}
	}

	return stats
}

// AuditStatistics 审计统计
type AuditStatistics struct {
	TotalLogs      int `json:"total_logs"`
	InfoCount      int `json:"info_count"`
	WarningCount   int `json:"warning_count"`
	ErrorCount     int `json:"error_count"`
	CriticalCount  int `json:"critical_count"`
	SuccessCount   int `json:"success_count"`
	FailureCount   int `json:"failure_count"`
}

// SearchLogs 搜索日志
func (a *AuditService) SearchLogs(keyword string, limit int) []*AuditLog {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var result []*AuditLog
	count := 0
	for i := len(a.logs) - 1; i >= 0 && count < limit; i-- {
		log := a.logs[i]
		if containsKeyword(log, keyword) {
			result = append(result, log)
			count++
		}
	}
	return result
}

func containsKeyword(log *AuditLog, keyword string) bool {
	return strings.Contains(log.User, keyword) ||
		strings.Contains(log.Resource, keyword) ||
		strings.Contains(log.Action, keyword) ||
		strings.Contains(log.Details, keyword)
}
