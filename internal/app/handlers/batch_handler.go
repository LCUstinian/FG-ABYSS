package handlers

import (
	"context"
	"strconv"

	"fg-abyss/internal/app/services"
)

// BatchHandler 批量操作处理器
type BatchHandler struct {
	service *services.BatchService
}

// NewBatchHandler 创建批量操作处理器
func NewBatchHandler() *BatchHandler {
	return &BatchHandler{
		service: services.NewBatchService(),
	}
}

// BatchOperationResponse 批量操作响应
type BatchOperationResponse struct {
	Total    int      `json:"total"`
	Success  int      `json:"success"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
	Skipped  int      `json:"skipped"`
	Duration string   `json:"duration"`
}

// ImportRequest 导入请求
type ImportRequest struct {
	// Data 导入数据
	Data string `json:"data"`
	// Format 导入格式（json/csv/xml）
	Format string `json:"format"`
	// FileType 文件类型（webshell/project）
	FileType string `json:"file_type"`
}

// Import 批量导入
func (h *BatchHandler) Import(ctx context.Context, req *ImportRequest) (*BatchOperationResponse, error) {
	format := services.ImportFormat(req.Format)
	
	result, err := h.service.ImportWebShells([]byte(req.Data), format)
	if err != nil {
		return &BatchOperationResponse{
			Errors: []string{err.Error()},
		}, nil
	}
	
	return &BatchOperationResponse{
		Total:    result.Total,
		Success:  result.Success,
		Failed:   result.Failed,
		Errors:   result.Errors,
		Skipped:  result.Skipped,
		Duration: result.Duration,
	}, nil
}

// ExportRequest 导出请求
type ExportRequest struct {
	// IDs 导出的 ID 列表
	IDs []int64 `json:"ids"`
	// Format 导出格式（json/csv/xml）
	Format string `json:"format"`
	// Filename 输出文件名
	Filename string `json:"filename"`
	// FileType 文件类型
	FileType string `json:"file_type"`
}

// ExportResponse 导出响应
type ExportResponse struct {
	// Success 是否成功
	Success bool `json:"success"`
	// Filename 文件名
	Filename string `json:"filename"`
	// Data 文件内容（Base64 编码）
	Data string `json:"data"`
	// Message 消息
	Message string `json:"message"`
}

// Export 批量导出
func (h *BatchHandler) Export(ctx context.Context, req *ExportRequest) (*ExportResponse, error) {
	// TODO: 根据 IDs 从数据库获取 WebShell 列表
	// 将 int64 IDs 转换为 string IDs
	webshellIDs := make([]string, len(req.IDs))
	for i, id := range req.IDs {
		webshellIDs[i] = strconv.FormatInt(id, 10)
	}
	
	format := services.ImportFormat(req.Format)
	
	err := h.service.ExportWebShells(webshellIDs, format, req.Filename)
	if err != nil {
		return &ExportResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &ExportResponse{
		Success:  true,
		Filename: req.Filename,
		Message:  "Export successful",
	}, nil
}

// BatchDeleteRequest 批量删除请求
type BatchDeleteRequest struct {
	IDs []string `json:"ids"`
}

// BatchDelete 批量删除
func (h *BatchHandler) BatchDelete(ctx context.Context, req *BatchDeleteRequest) (*BatchOperationResponse, error) {
	result, err := h.service.BatchDeleteWebShells(req.IDs)
	if err != nil {
		return &BatchOperationResponse{
			Errors: []string{err.Error()},
		}, nil
	}
	
	return &BatchOperationResponse{
		Total:    result.Total,
		Success:  result.Success,
		Failed:   result.Failed,
		Errors:   result.Errors,
		Skipped:  result.Skipped,
		Duration: result.Duration,
	}, nil
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

// BatchUpdateStatus 批量更新状态
func (h *BatchHandler) BatchUpdateStatus(ctx context.Context, req *BatchUpdateStatusRequest) (*BatchOperationResponse, error) {
	result, err := h.service.BatchUpdateStatus(req.IDs, req.Status)
	if err != nil {
		return &BatchOperationResponse{
			Errors: []string{err.Error()},
		}, nil
	}
	
	return &BatchOperationResponse{
		Total:    result.Total,
		Success:  result.Success,
		Failed:   result.Failed,
		Errors:   result.Errors,
		Skipped:  result.Skipped,
		Duration: result.Duration,
	}, nil
}

// BatchTestRequest 批量测试请求
type BatchTestRequest struct {
	IDs []string `json:"ids"`
}

// BatchTest 批量测试连接
func (h *BatchHandler) BatchTest(ctx context.Context, req *BatchTestRequest) (*BatchOperationResponse, error) {
	// TODO: 根据 IDs 获取 WebShell 列表
	
	result, err := h.service.BatchTestConnection(req.IDs)
	if err != nil {
		return &BatchOperationResponse{
			Errors: []string{err.Error()},
		}, nil
	}
	
	return &BatchOperationResponse{
		Total:    result.Total,
		Success:  result.Success,
		Failed:   result.Failed,
		Errors:   result.Errors,
		Skipped:  result.Skipped,
		Duration: result.Duration,
	}, nil
}

// GetImportTemplateRequest 获取导入模板请求
type GetImportTemplateRequest struct {
	Format string `json:"format"`
}

// GetImportTemplateResponse 获取导入模板响应
type GetImportTemplateResponse struct {
	// Template 模板内容
	Template string `json:"template"`
	// Format 格式
	Format string `json:"format"`
}

// GetImportTemplate 获取导入模板
func (h *BatchHandler) GetImportTemplate(ctx context.Context, req *GetImportTemplateRequest) (*GetImportTemplateResponse, error) {
	format := services.ImportFormat(req.Format)
	
	data, err := h.service.GenerateImportTemplate(format)
	if err != nil {
		return &GetImportTemplateResponse{
			Template: "",
			Format:   req.Format,
		}, nil
	}
	
	return &GetImportTemplateResponse{
		Template: string(data),
		Format:   req.Format,
	}, nil
}

// ValidateImportRequest 验证导入数据请求
type ValidateImportRequest struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}

// ValidateImport 验证导入数据
func (h *BatchHandler) ValidateImport(ctx context.Context, req *ValidateImportRequest) (*BatchOperationResponse, error) {
	format := services.ImportFormat(req.Format)
	
	result, err := h.service.ValidateImportData([]byte(req.Data), format)
	if err != nil {
		return &BatchOperationResponse{
			Errors: []string{err.Error()},
		}, nil
	}
	
	return &BatchOperationResponse{
		Total:    result.Total,
		Success:  result.Success,
		Failed:   result.Failed,
		Errors:   result.Errors,
		Skipped:  result.Skipped,
		Duration: result.Duration,
	}, nil
}
