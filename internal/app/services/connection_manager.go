package services

import (
	"context"
	"sync"
	"time"

	"fg-abyss/internal/domain/entity"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*entity.Connection
	mu          sync.RWMutex
	httpEngine  *HTTPRequestEngine
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*entity.Connection),
		httpEngine:  NewHTTPRequestEngine(),
	}
}

// CreateConnection 创建新连接
func (m *ConnectionManager) CreateConnection(webshellID string, config *entity.ConnectionConfig) (*entity.Connection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查是否已存在连接
	if conn, exists := m.connections[webshellID]; exists {
		return conn, nil
	}

	// 创建新连接
	conn := entity.NewConnection(webshellID, config)
	m.connections[webshellID] = conn

	return conn, nil
}

// GetConnection 获取连接
func (m *ConnectionManager) GetConnection(webshellID string) (*entity.Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.connections[webshellID]
	return conn, exists
}

// RemoveConnection 移除连接
func (m *ConnectionManager) RemoveConnection(webshellID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.connections, webshellID)
}

// UpdateConnectionStatus 更新连接状态
func (m *ConnectionManager) UpdateConnectionStatus(webshellID string, status entity.ConnectionStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	conn, exists := m.connections[webshellID]
	if !exists {
		return ErrConnectionNotFound
	}

	conn.UpdateStatus(status)
	return nil
}

// GetHTTPEngine 获取 HTTP 引擎
func (m *ConnectionManager) GetHTTPEngine() *HTTPRequestEngine {
	return m.httpEngine
}

// GetAllConnections 获取所有连接
func (m *ConnectionManager) GetAllConnections() []*entity.Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connections := make([]*entity.Connection, 0, len(m.connections))
	for _, conn := range m.connections {
		connections = append(connections, conn)
	}

	return connections
}

// GetConnectionCount 获取连接数量
func (m *ConnectionManager) GetConnectionCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// CleanupInactiveConnections 清理不活跃的连接
func (m *ConnectionManager) CleanupInactiveConnections(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, conn := range m.connections {
		if now.Sub(conn.LastActiveTime) > timeout {
			delete(m.connections, id)
		}
	}
}

// StartHeartbeat 启动心跳检测
func (m *ConnectionManager) StartHeartbeat(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.CleanupInactiveConnections(interval * 2)
		}
	}
}
