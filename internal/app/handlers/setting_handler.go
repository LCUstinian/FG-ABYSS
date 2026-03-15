package handlers

import (
	"fg-abyss/internal/app/services"
	"fg-abyss/internal/domain/entity"
)

// SettingHandler 设置处理器
type SettingHandler struct {
	settingService *services.SettingService
}

// NewSettingHandler 创建设置处理器实例
func NewSettingHandler(settingService *services.SettingService) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
	}
}

// GetAllSettings 获取所有设置
func (h *SettingHandler) GetAllSettings() ([]entity.Setting, error) {
	return h.settingService.GetAllSettings()
}

// GetSettingsByCategory 根据分类获取设置
func (h *SettingHandler) GetSettingsByCategory(category string) ([]entity.Setting, error) {
	return h.settingService.GetSettingsByCategory(category)
}

// GetSettingByKey 根据键名获取设置
func (h *SettingHandler) GetSettingByKey(key string) (*entity.Setting, error) {
	return h.settingService.GetSettingByKey(key)
}

// UpdateSetting 更新设置值
func (h *SettingHandler) UpdateSetting(key string, value string) error {
	return h.settingService.UpdateSetting(key, value)
}

// UpdateSettings 批量更新设置
func (h *SettingHandler) UpdateSettings(updates map[string]string) error {
	return h.settingService.UpdateSettings(updates)
}

// GetStringValue 获取字符串类型设置值
func (h *SettingHandler) GetStringValue(key string) (string, error) {
	return h.settingService.GetStringValue(key)
}

// GetBoolValue 获取布尔类型设置值
func (h *SettingHandler) GetBoolValue(key string) (bool, error) {
	return h.settingService.GetBoolValue(key)
}

// GetIntValue 获取整数类型设置值
func (h *SettingHandler) GetIntValue(key string) (int, error) {
	return h.settingService.GetIntValue(key)
}
