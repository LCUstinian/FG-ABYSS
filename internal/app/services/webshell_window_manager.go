package services

import (
	"fmt"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// WebShellWindowManager WebShell 窗口管理器
type WebShellWindowManager struct {
	windows map[string]*application.Window
	mu      sync.RWMutex
	app     *application.App
}

// NewWebShellWindowManager 创建窗口管理器
func NewWebShellWindowManager(app *application.App) *WebShellWindowManager {
	return &WebShellWindowManager{
		windows: make(map[string]*application.Window),
		app:     app,
	}
}

// WindowInfo 窗口信息
type WindowInfo struct {
	WebShellID string `json:"webshell_id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	IsOpen     bool   `json:"is_open"`
}

// OpenWindow 打开 WebShell 窗口
func (m *WebShellWindowManager) OpenWindow(webshellID, name, url string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查窗口是否已存在
	if _, exists := m.windows[webshellID]; exists {
		return fmt.Errorf("window for WebShell %s already exists", webshellID)
	}

	// 创建新窗口
	// 注意：Wails v3 的窗口 API 可能有所不同，这里使用伪代码
	// 实际实现需要根据 Wails v3 的具体 API 调整
	
	// windowURL := fmt.Sprintf("index.html#/webshell-control?id=%s", webshellID)
	
	// 这里需要使用 Wails v3 的窗口创建 API
	// window, err := m.app.NewWindow(application.WindowOptions{
	// 	Title:  fmt.Sprintf("WebShell Control - %s", name),
	// 	URL:    windowURL,
	// 	Width:  1200,
	// 	Height: 800,
	// })
	
	// if err != nil {
	// 	return err
	// }
	
	// m.windows[webshellID] = window
	
	// TODO: 暂时返回 nil，实际实现需要 Wails v3 的窗口 API 支持
	return nil
}

// CloseWindow 关闭 WebShell 窗口
func (m *WebShellWindowManager) CloseWindow(webshellID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.windows[webshellID]
	if !exists {
		return fmt.Errorf("window for WebShell %s not found", webshellID)
	}

	// 关闭窗口
	// window.Close()
	
	delete(m.windows, webshellID)
	return nil
}

// GetWindow 获取窗口
func (m *WebShellWindowManager) GetWindow(webshellID string) (*application.Window, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	window, exists := m.windows[webshellID]
	if !exists {
		return nil, fmt.Errorf("window for WebShell %s not found", webshellID)
	}

	return window, nil
}

// IsWindowOpen 检查窗口是否已打开
func (m *WebShellWindowManager) IsWindowOpen(webshellID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.windows[webshellID]
	return exists
}

// GetAllWindows 获取所有窗口信息
func (m *WebShellWindowManager) GetAllWindows() []WindowInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	infos := make([]WindowInfo, 0, len(m.windows))
	for id := range m.windows {
		infos = append(infos, WindowInfo{
			WebShellID: id,
			IsOpen:     true,
		})
	}

	return infos
}

// CloseAllWindows 关闭所有窗口
func (m *WebShellWindowManager) CloseAllWindows() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id := range m.windows {
		// window.Close()
		delete(m.windows, id)
	}
}
