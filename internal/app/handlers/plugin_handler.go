package handlers

import (
	"context"

	"fg-abyss/internal/plugin"
)

// PluginHandler 插件管理处理器
type PluginHandler struct {
	loader *plugin.PluginLoader
}

// NewPluginHandler 创建插件处理器
func NewPluginHandler(loader *plugin.PluginLoader) *PluginHandler {
	return &PluginHandler{
		loader: loader,
	}
}

// PluginInfoResponse 插件信息响应
type PluginInfoResponse struct {
	ID         string              `json:"id"`
	Metadata   *plugin.PluginMetadata `json:"metadata"`
	Status     plugin.PluginStatus `json:"status"`
	IsEnabled  bool                `json:"is_enabled"`
	CanDisable bool                `json:"can_disable"` // 是否可以禁用（内置插件通常不可禁用）
	Error      string              `json:"error,omitempty"`
}

// PluginListResponse 插件列表响应
type PluginListResponse struct {
	Total   int                    `json:"total"`
	Builtin int                    `json:"builtin"`
	External int                   `json:"external"`
	Plugins []*PluginInfoResponse `json:"plugins"`
}

// InstallPluginRequest 安装插件请求
type InstallPluginRequest struct {
	PluginPath string `json:"plugin_path"`
}

// InstallPluginResponse 安装插件响应
type InstallPluginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Plugin  *PluginInfoResponse `json:"plugin,omitempty"`
}

// UninstallPluginRequest 卸载插件请求
type UninstallPluginRequest struct {
	PluginID string `json:"plugin_id"`
}

// UninstallPluginResponse 卸载插件响应
type UninstallPluginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// EnablePluginRequest 启用插件请求
type EnablePluginRequest struct {
	PluginID string `json:"plugin_id"`
}

// EnablePluginResponse 启用插件响应
type EnablePluginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DisablePluginRequest 禁用插件请求
type DisablePluginRequest struct {
	PluginID string `json:"plugin_id"`
}

// DisablePluginResponse 禁用插件响应
type DisablePluginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetPluginList 获取插件列表
func (h *PluginHandler) GetPluginList(ctx context.Context) (*PluginListResponse, error) {
	pluginList := h.loader.GetPluginList()
	
	response := &PluginListResponse{
		Plugins: make([]*PluginInfoResponse, 0, len(pluginList)),
	}

	for _, p := range pluginList {
		canDisable := p.Metadata.Type == plugin.PluginTypeExternal
		
		response.Plugins = append(response.Plugins, &PluginInfoResponse{
			ID:         p.ID,
			Metadata:   p.Metadata,
			Status:     p.Status,
			IsEnabled:  p.IsEnabled,
			CanDisable: canDisable,
			Error:      p.Error,
		})

		response.Total++
		if p.Metadata.Type == plugin.PluginTypeBuiltin {
			response.Builtin++
		} else {
			response.External++
		}
	}

	return response, nil
}

// InstallPlugin 安装插件
func (h *PluginHandler) InstallPlugin(ctx context.Context, req *InstallPluginRequest) (*InstallPluginResponse, error) {
	info, err := h.loader.InstallPlugin(ctx, req.PluginPath)
	if err != nil {
		return &InstallPluginResponse{
			Success: false,
			Message: "安装失败：" + err.Error(),
		}, nil
	}

	return &InstallPluginResponse{
		Success: true,
		Message: "插件安装成功",
		Plugin: &PluginInfoResponse{
			ID:         info.Metadata.ID,
			Metadata:   info.Metadata,
			Status:     plugin.PluginStatusEnabled,
			IsEnabled:  info.IsEnabled,
			CanDisable: true,
		},
	}, nil
}

// UninstallPlugin 卸载插件
func (h *PluginHandler) UninstallPlugin(ctx context.Context, req *UninstallPluginRequest) (*UninstallPluginResponse, error) {
	err := h.loader.UninstallPlugin(ctx, req.PluginID)
	if err != nil {
		return &UninstallPluginResponse{
			Success: false,
			Message: "卸载失败：" + err.Error(),
		}, nil
	}

	return &UninstallPluginResponse{
		Success: true,
		Message: "插件卸载成功",
	}, nil
}

// EnablePlugin 启用插件
func (h *PluginHandler) EnablePlugin(ctx context.Context, req *EnablePluginRequest) (*EnablePluginResponse, error) {
	err := h.loader.EnablePlugin(ctx, req.PluginID)
	if err != nil {
		return &EnablePluginResponse{
			Success: false,
			Message: "启用失败：" + err.Error(),
		}, nil
	}

	return &EnablePluginResponse{
		Success: true,
		Message: "插件已启用",
	}, nil
}

// DisablePlugin 禁用插件
func (h *PluginHandler) DisablePlugin(ctx context.Context, req *DisablePluginRequest) (*DisablePluginResponse, error) {
	err := h.loader.DisablePlugin(ctx, req.PluginID)
	if err != nil {
		return &DisablePluginResponse{
			Success: false,
			Message: "禁用失败：" + err.Error(),
		}, nil
	}

	return &DisablePluginResponse{
		Success: true,
		Message: "插件已禁用",
	}, nil
}

// GetPluginDetail 获取插件详情
func (h *PluginHandler) GetPluginDetail(ctx context.Context, pluginID string) (*PluginInfoResponse, error) {
	pluginList := h.loader.GetPluginList()
	
	for _, p := range pluginList {
		if p.ID == pluginID {
			canDisable := p.Metadata.Type == plugin.PluginTypeExternal
			
			return &PluginInfoResponse{
				ID:         p.ID,
				Metadata:   p.Metadata,
				Status:     p.Status,
				IsEnabled:  p.IsEnabled,
				CanDisable: canDisable,
				Error:      p.Error,
			}, nil
		}
	}

	return nil, nil
}

// CheckUpdate 检查插件更新
func (h *PluginHandler) CheckUpdate(ctx context.Context, pluginID string) (*PluginUpdateInfo, error) {
	// TODO: 实现插件更新检查
	// 这里可以预留插件商店接口
	return &PluginUpdateInfo{
		HasUpdate: false,
		Message:   "当前已是最新版本",
	}, nil
}

// PluginUpdateInfo 插件更新信息
type PluginUpdateInfo struct {
	HasUpdate     bool   `json:"has_update"`
	LatestVersion string `json:"latest_version,omitempty"`
	CurrentVersion string `json:"current_version"`
	Message       string `json:"message"`
	DownloadURL   string `json:"download_url,omitempty"`
}

// UpdatePlugin 更新插件
func (h *PluginHandler) UpdatePlugin(ctx context.Context, pluginID string) (*InstallPluginResponse, error) {
	// TODO: 实现插件更新
	return &InstallPluginResponse{
		Success: false,
		Message: "插件更新功能尚未实现",
	}, nil
}
