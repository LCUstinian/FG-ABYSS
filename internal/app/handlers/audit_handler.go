package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
)

// AuditHandler 审计日志处理器
type AuditHandler struct {
	service *services.AuditService
}

// NewAuditHandler 创建审计处理器
func NewAuditHandler() *AuditHandler {
	config := &services.AuditConfig{
		MaxSize:  10000,
		FilePath: "data/audit_logs.json",
		AutoSave: true,
	}

	return &AuditHandler{
		service: services.NewAuditService(config),
	}
}

// LogRequest 日志请求
type LogRequest struct {
	Level     string `json:"level"`
	Operation string `json:"operation"`
	User      string `json:"user"`
	Resource  string `json:"resource"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Success   bool   `json:"success"`
	Duration  int64  `json:"duration"`
}

// Log 记录审计日志
func (h *AuditHandler) Log(ctx context.Context, req *LogRequest) error {
	level := services.AuditLevel(req.Level)
	operation := services.AuditOperation(req.Operation)

	h.service.Log(level, operation, req.User, req.Resource, req.Action, req.Details, req.Success, req.Duration)
	return nil
}

// GetLogsRequest 获取日志请求
type GetLogsRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// LogEntry 日志条目
type LogEntry struct {
	ID        int64  `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Operation string `json:"operation"`
	User      string `json:"user"`
	Resource  string `json:"resource"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address"`
	Details   string `json:"details"`
	Success   bool   `json:"success"`
	Duration  int64  `json:"duration"`
}

// GetLogsResponse 获取日志响应
type GetLogsResponse struct {
	Logs []LogEntry `json:"logs"`
}

// GetLogs 获取日志列表
func (h *AuditHandler) GetLogs(ctx context.Context, req *GetLogsRequest) (*GetLogsResponse, error) {
	logs := h.service.GetLogs(req.Limit, req.Offset)

	entries := make([]LogEntry, len(logs))
	for i, log := range logs {
		entries[i] = LogEntry{
			ID:        log.ID,
			Timestamp: log.Timestamp.Format("2006-01-02 15:04:05"),
			Level:     string(log.Level),
			Operation: string(log.Operation),
			User:      log.User,
			Resource:  log.Resource,
			Action:    log.Action,
			Details:   log.Details,
			Success:   log.Success,
			Duration:  log.Duration,
		}
	}

	return &GetLogsResponse{
		Logs: entries,
	}, nil
}

// GetLogsByUserRequest 按用户获取日志请求
type GetLogsByUserRequest struct {
	User  string `json:"user"`
	Limit int    `json:"limit"`
}

// GetLogsByUser 按用户获取日志
func (h *AuditHandler) GetLogsByUser(ctx context.Context, req *GetLogsByUserRequest) (*GetLogsResponse, error) {
	logs := h.service.GetLogsByUser(req.User, req.Limit)

	entries := make([]LogEntry, len(logs))
	for i, log := range logs {
		entries[i] = LogEntry{
			ID:        log.ID,
			Timestamp: log.Timestamp.Format("2006-01-02 15:04:05"),
			Level:     string(log.Level),
			Operation: string(log.Operation),
			User:      log.User,
			Resource:  log.Resource,
			Action:    log.Action,
			Details:   log.Details,
			Success:   log.Success,
			Duration:  log.Duration,
		}
	}

	return &GetLogsResponse{
		Logs: entries,
	}, nil
}

// GetStatisticsResponse 获取统计响应
type GetStatisticsResponse struct {
	TotalLogs      int `json:"total_logs"`
	InfoCount      int `json:"info_count"`
	WarningCount   int `json:"warning_count"`
	ErrorCount     int `json:"error_count"`
	CriticalCount  int `json:"critical_count"`
	SuccessCount   int `json:"success_count"`
	FailureCount   int `json:"failure_count"`
}

// GetStatistics 获取统计信息
func (h *AuditHandler) GetStatistics(ctx context.Context) (*GetStatisticsResponse, error) {
	stats := h.service.GetStatistics()

	return &GetStatisticsResponse{
		TotalLogs:      stats.TotalLogs,
		InfoCount:      stats.InfoCount,
		WarningCount:   stats.WarningCount,
		ErrorCount:     stats.ErrorCount,
		CriticalCount:  stats.CriticalCount,
		SuccessCount:   stats.SuccessCount,
		FailureCount:   stats.FailureCount,
	}, nil
}

// SearchLogsRequest 搜索日志请求
type SearchLogsRequest struct {
	Keyword string `json:"keyword"`
	Limit   int    `json:"limit"`
}

// SearchLogs 搜索日志
func (h *AuditHandler) SearchLogs(ctx context.Context, req *SearchLogsRequest) (*GetLogsResponse, error) {
	logs := h.service.SearchLogs(req.Keyword, req.Limit)

	entries := make([]LogEntry, len(logs))
	for i, log := range logs {
		entries[i] = LogEntry{
			ID:        log.ID,
			Timestamp: log.Timestamp.Format("2006-01-02 15:04:05"),
			Level:     string(log.Level),
			Operation: string(log.Operation),
			User:      log.User,
			Resource:  log.Resource,
			Action:    log.Action,
			Details:   log.Details,
			Success:   log.Success,
			Duration:  log.Duration,
		}
	}

	return &GetLogsResponse{
		Logs: entries,
	}, nil
}

// ExportLogsRequest 导出日志请求
type ExportLogsRequest struct {
	FilePath string `json:"file_path"`
}

// ExportLogs 导出日志
func (h *AuditHandler) ExportLogs(ctx context.Context, req *ExportLogsRequest) error {
	return h.service.ExportToJSON(req.FilePath)
}

// ClearLogs 清空日志
func (h *AuditHandler) ClearLogs(ctx context.Context) error {
	h.service.Clear()
	return nil
}
