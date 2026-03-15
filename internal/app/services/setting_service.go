package services

import (
	"encoding/json"
	"errors"
	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
	"strconv"
	"sync"
)

// SettingService 设置服务
type SettingService struct {
	settingRepo   repository.SettingRepository
	settingsCache map[string]*entity.Setting
	cacheMutex    sync.RWMutex
}

// NewSettingService 创建设置服务实例
func NewSettingService(settingRepo repository.SettingRepository) *SettingService {
	return &SettingService{
		settingRepo:   settingRepo,
		settingsCache: make(map[string]*entity.Setting),
	}
}

// Initialize 初始化设置（启动时从数据库加载到内存）
func (s *SettingService) Initialize() error {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return err
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for _, setting := range settings {
		settingCopy := setting
		s.settingsCache[setting.Key] = &settingCopy
	}

	return nil
}

// InitializeDefaults 初始化默认设置
func (s *SettingService) InitializeDefaults(defaults []entity.Setting) error {
	if err := s.settingRepo.InitializeDefaults(defaults); err != nil {
		return err
	}
	return s.Initialize()
}

// GetAllSettings 获取所有设置
func (s *SettingService) GetAllSettings() ([]entity.Setting, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	settings := make([]entity.Setting, 0, len(s.settingsCache))
	for _, setting := range s.settingsCache {
		settings = append(settings, *setting)
	}
	return settings, nil
}

// GetSettingsByCategory 根据分类获取设置
func (s *SettingService) GetSettingsByCategory(category string) ([]entity.Setting, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	settings := make([]entity.Setting, 0)
	for _, setting := range s.settingsCache {
		if setting.Category == category {
			settings = append(settings, *setting)
		}
	}
	return settings, nil
}

// GetSettingByKey 根据键名获取设置
func (s *SettingService) GetSettingByKey(key string) (*entity.Setting, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	setting, exists := s.settingsCache[key]
	if !exists {
		return nil, errors.New("设置不存在")
	}

	settingCopy := *setting
	return &settingCopy, nil
}

// GetStringValue 获取字符串类型设置值
func (s *SettingService) GetStringValue(key string) (string, error) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	setting, exists := s.settingsCache[key]
	if !exists {
		return "", errors.New("设置不存在")
	}
	return setting.Value, nil
}

// GetBoolValue 获取布尔类型设置值
func (s *SettingService) GetBoolValue(key string) (bool, error) {
	value, err := s.GetStringValue(key)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

// GetIntValue 获取整数类型设置值
func (s *SettingService) GetIntValue(key string) (int, error) {
	value, err := s.GetStringValue(key)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

// GetJSONValue 获取 JSON 类型设置值
func (s *SettingService) GetJSONValue(key string, result interface{}) error {
	value, err := s.GetStringValue(key)
	if err != nil {
		return err
	}
	if value == "" {
		return errors.New("设置值为空")
	}
	return json.Unmarshal([]byte(value), result)
}

// UpdateSetting 更新设置值（同时更新数据库和内存缓存）
func (s *SettingService) UpdateSetting(key string, value string) error {
	// 更新数据库
	if err := s.settingRepo.UpdateValue(key, value); err != nil {
		return err
	}

	// 更新内存缓存
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	if setting, exists := s.settingsCache[key]; exists {
		setting.Value = value
	}

	return nil
}

// UpdateSettings 批量更新设置
func (s *SettingService) UpdateSettings(updates map[string]string) error {
	// 批量保存到数据库
	settings := make([]entity.Setting, 0, len(updates))
	for key, value := range updates {
		settings = append(settings, entity.Setting{
			Key:   key,
			Value: value,
		})
	}

	if err := s.settingRepo.SaveAll(settings); err != nil {
		return err
	}

	// 更新内存缓存
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	for key, value := range updates {
		if setting, exists := s.settingsCache[key]; exists {
			setting.Value = value
		}
	}

	return nil
}

// GetDefaultSettings 获取所有默认设置（基于现有设置界面）
func (s *SettingService) GetDefaultSettings() []entity.Setting {
	return []entity.Setting{
		// 外观设置
		{
			Category:    "appearance",
			Key:         "theme",
			Value:       "dark",
			Type:        "string",
			Name:        "主题",
			Description: "应用程序主题",
		},
		{
			Category:    "appearance",
			Key:         "language",
			Value:       "zh-CN",
			Type:        "string",
			Name:        "语言",
			Description: "界面显示语言",
		},
		{
			Category:    "appearance",
			Key:         "accentColor",
			Value:       "#18a058",
			Type:        "color",
			Name:        "强调色",
			Description: "界面强调色",
		},
		// 代理设置
		{
			Category:    "proxy",
			Key:         "proxyEnabled",
			Value:       "false",
			Type:        "bool",
			Name:        "启用代理",
			Description: "是否启用代理",
		},
		{
			Category:    "proxy",
			Key:         "proxyAddress",
			Value:       "",
			Type:        "string",
			Name:        "代理地址",
			Description: "代理服务器地址",
		},
		{
			Category:    "proxy",
			Key:         "proxyPort",
			Value:       "8080",
			Type:        "int",
			Name:        "代理端口",
			Description: "代理服务器端口",
		},
		// 网络设置
		{
			Category:    "network",
			Key:         "timeout",
			Value:       "30",
			Type:        "int",
			Name:        "超时时间",
			Description: "网络请求超时时间（秒）",
		},
		{
			Category:    "network",
			Key:         "retryCount",
			Value:       "3",
			Type:        "int",
			Name:        "重试次数",
			Description: "网络请求失败重试次数",
		},
	}
}
