// 示例外置插件 - Hello World
// 编译命令：go build -buildmode=plugin -o hello.so hello_plugin.go

package main

import (
	"context"
	"fg-abyss/internal/plugin"
	"time"
)

// PluginMetadata 插件元数据（必须导出）
var PluginMetadata = &plugin.PluginMetadata{
	ID:          "example.hello_world",
	Name:        "Hello World",
	Version:     "1.0.0",
	Description: "示例插件 - Hello World",
	Author:      "Your Name",
	Type:        plugin.PluginTypeExternal,
	Category:    "example",
	Tags:        []string{"example", "demo"},
	Homepage:    "https://github.com/yourname/hello-plugin",
	License:     "MIT",
	Permissions: []plugin.Permission{},
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

// HelloPlugin Hello World 插件
type HelloPlugin struct {
	*plugin.BasePlugin
	api plugin.PluginAPI
}

// NewPlugin 创建插件实例（必须导出）
func NewPlugin() plugin.Plugin {
	return &HelloPlugin{
		BasePlugin: plugin.NewBasePlugin(PluginMetadata),
	}
}

// Initialize 初始化插件
func (p *HelloPlugin) Initialize(ctx context.Context, api plugin.PluginAPI) error {
	p.api = api
	
	// 注册命令
	if err := api.RegisterCommand("hello", p.handleHelloCommand); err != nil {
		return err
	}

	// 订阅事件
	if err := api.SubscribeEvent("app.started", p.handleAppStarted); err != nil {
		return err
	}

	// 记录日志
	api.Log(plugin.LogLevelInfo, "Hello World 插件已初始化", nil)

	return p.BasePlugin.Initialize(ctx, api)
}

// Start 启动插件
func (p *HelloPlugin) Start(ctx context.Context) error {
	p.api.Log(plugin.LogLevelInfo, "Hello World 插件已启动", nil)
	return p.BasePlugin.Start(ctx)
}

// Stop 停止插件
func (p *HelloPlugin) Stop(ctx context.Context) error {
	p.api.Log(plugin.LogLevelInfo, "Hello World 插件已停止", nil)
	return p.BasePlugin.Stop(ctx)
}

// Shutdown 关闭插件
func (p *HelloPlugin) Shutdown(ctx context.Context) error {
	p.api.Log(plugin.LogLevelInfo, "Hello World 插件已关闭", nil)
	return p.BasePlugin.Shutdown(ctx)
}

// handleHelloCommand 处理 Hello 命令
func (p *HelloPlugin) handleHelloCommand(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	name := ""
	if n, ok := params["name"].(string); ok {
		name = n
	}

	if name == "" {
		name = "World"
	}

	message := "Hello, " + name + "! from Hello World plugin"
	
	p.api.Log(plugin.LogLevelInfo, "Hello command executed", map[string]interface{}{
		"message": message,
	})

	return map[string]interface{}{
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}, nil
}

// handleAppStarted 处理应用启动事件
func (p *HelloPlugin) handleAppStarted(event string, data interface{}) error {
	p.api.Log(plugin.LogLevelInfo, "Application started event received", nil)
	return nil
}
