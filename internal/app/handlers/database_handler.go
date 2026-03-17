package handlers

import (
	"context"

	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// DatabaseHandler 数据库管理处理器
type DatabaseHandler struct {
	service *services.DatabaseService
}

// NewDatabaseHandler 创建数据库处理器
func NewDatabaseHandler() *DatabaseHandler {
	return &DatabaseHandler{
		service: services.NewDatabaseService(),
	}
}

// ConnectRequest 连接请求
type ConnectRequest struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Charset   string `json:"charset"`
	SSLMode   bool   `json:"ssl_mode"`
}

// ConnectResponse 连接响应
type ConnectResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Connect 连接数据库
func (h *DatabaseHandler) Connect(ctx context.Context, req *ConnectRequest) (*ConnectResponse, error) {
	conn := &entity.DatabaseConnection{
		ID:       req.ID,
		Name:     req.Name,
		Type:     entity.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Database: req.Database,
		Charset:  req.Charset,
		SSLMode:  req.SSLMode,
	}
	
	if err := h.service.Connect(conn); err != nil {
		return &ConnectResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &ConnectResponse{
		Success: true,
		Message: "Connected successfully",
	}, nil
}

// DisconnectRequest 断开连接请求
type DisconnectRequest struct {
	ConnID int64 `json:"conn_id"`
}

// Disconnect 断开数据库连接
func (h *DatabaseHandler) Disconnect(ctx context.Context, req *DisconnectRequest) (*ConnectResponse, error) {
	if err := h.service.Disconnect(req.ConnID); err != nil {
		return &ConnectResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &ConnectResponse{
		Success: true,
		Message: "Disconnected successfully",
	}, nil
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Charset  string `json:"charset"`
	SSLMode  bool   `json:"ssl_mode"`
}

// TestConnectionResponse 测试连接响应
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TestConnection 测试数据库连接
func (h *DatabaseHandler) TestConnection(ctx context.Context, req *TestConnectionRequest) (*TestConnectionResponse, error) {
	conn := &entity.DatabaseConnection{
		Type:     entity.DatabaseType(req.Type),
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Database: req.Database,
		Charset:  req.Charset,
		SSLMode:  req.SSLMode,
	}
	
	success, message := h.service.TestConnection(conn)
	
	return &TestConnectionResponse{
		Success: success,
		Message: message,
	}, nil
}

// ExecuteQueryRequest 执行查询请求
type ExecuteQueryRequest struct {
	ConnID  int64  `json:"conn_id"`
	SQL     string `json:"sql"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
	Timeout int    `json:"timeout"`
}

// ExecuteQueryResponse 执行查询响应
type ExecuteQueryResponse struct {
	Columns  []string               `json:"columns"`
	Rows     []map[string]interface{} `json:"rows"`
	Affected int64                  `json:"affected"`
	Duration string                 `json:"duration"`
	Success  bool                   `json:"success"`
	Message  string                 `json:"message"`
}

// ExecuteQuery 执行 SQL 查询
func (h *DatabaseHandler) ExecuteQuery(ctx context.Context, req *ExecuteQueryRequest) (*ExecuteQueryResponse, error) {
	query := entity.DatabaseQuery{
		SQL:     req.SQL,
		Limit:   req.Limit,
		Offset:  req.Offset,
		Timeout: req.Timeout,
	}
	
	result, err := h.service.ExecuteQuery(req.ConnID, query)
	if err != nil {
		return &ExecuteQueryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &ExecuteQueryResponse{
		Columns:  result.Columns,
		Rows:     result.Rows,
		Affected: result.Affected,
		Duration: result.Duration,
		Success:  result.Success,
		Message:  result.Message,
	}, nil
}

// ExecuteUpdateRequest 执行更新请求
type ExecuteUpdateRequest struct {
	ConnID int64  `json:"conn_id"`
	SQL    string `json:"sql"`
}

// ExecuteUpdate 执行 SQL 更新
func (h *DatabaseHandler) ExecuteUpdate(ctx context.Context, req *ExecuteUpdateRequest) (*ExecuteQueryResponse, error) {
	result, err := h.service.ExecuteUpdate(req.ConnID, req.SQL)
	if err != nil {
		return &ExecuteQueryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	
	return &ExecuteQueryResponse{
		Affected: result.Affected,
		Duration: result.Duration,
		Success:  result.Success,
		Message:  result.Message,
	}, nil
}

// GetTablesRequest 获取表请求
type GetTablesRequest struct {
	ConnID   int64  `json:"conn_id"`
	ConnType string `json:"conn_type"`
}

// GetTablesResponse 获取表响应
type GetTablesResponse struct {
	Tables []entity.DatabaseTable `json:"tables"`
}

// GetTables 获取数据库表列表
func (h *DatabaseHandler) GetTables(ctx context.Context, req *GetTablesRequest) (*GetTablesResponse, error) {
	tables, err := h.service.GetTables(req.ConnID, entity.DatabaseType(req.ConnType))
	if err != nil {
		return &GetTablesResponse{
			Tables: []entity.DatabaseTable{},
		}, nil
	}
	
	return &GetTablesResponse{
		Tables: tables,
	}, nil
}

// GetTableColumnsRequest 获取表列请求
type GetTableColumnsRequest struct {
	ConnID    int64  `json:"conn_id"`
	ConnType  string `json:"conn_type"`
	TableName string `json:"table_name"`
}

// GetTableColumnsResponse 获取表列响应
type GetTableColumnsResponse struct {
	Columns []entity.TableColumn `json:"columns"`
}

// GetTableColumns 获取表列信息
func (h *DatabaseHandler) GetTableColumns(ctx context.Context, req *GetTableColumnsRequest) (*GetTableColumnsResponse, error) {
	columns, err := h.service.GetTableColumns(req.ConnID, entity.DatabaseType(req.ConnType), req.TableName)
	if err != nil {
		return &GetTableColumnsResponse{
			Columns: []entity.TableColumn{},
		}, nil
	}
	
	return &GetTableColumnsResponse{
		Columns: columns,
	}, nil
}

// GetDatabasesRequest 获取数据库列表请求
type GetDatabasesRequest struct {
	ConnID int64 `json:"conn_id"`
}

// GetDatabasesResponse 获取数据库列表响应
type GetDatabasesResponse struct {
	Databases []string `json:"databases"`
}

// GetDatabases 获取数据库列表（用于 MySQL/PostgreSQL）
func (h *DatabaseHandler) GetDatabases(ctx context.Context, req *GetDatabasesRequest) (*GetDatabasesResponse, error) {
	// TODO: 实现获取数据库列表
	return &GetDatabasesResponse{
		Databases: []string{},
	}, nil
}
