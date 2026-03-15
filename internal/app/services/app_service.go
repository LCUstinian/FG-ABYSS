package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"gorm.io/gorm"

	"fg-abyss/internal/domain/entity"
	"fg-abyss/internal/domain/repository"
)

// AppService 应用服务
type AppService struct {
	db             *gorm.DB
	projectRepo    repository.ProjectRepository
	webshellRepo   repository.WebShellRepository
	startTime      time.Time
	process        *process.Process
	processID      int32
}

// NewAppService 创建应用服务实例
func NewAppService(
	db *gorm.DB,
	projectRepo repository.ProjectRepository,
	webshellRepo repository.WebShellRepository,
) *AppService {
	// 获取当前进程信息
	pid := os.Getpid()
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Printf("Warning: Failed to get process info: %v\n", err)
	}

	return &AppService{
		db:           db,
		projectRepo:  projectRepo,
		webshellRepo: webshellRepo,
		startTime:    time.Now(),
		process:      proc,
		processID:    int32(pid),
	}
}

// GetSystemStatus 获取系统状态
func (s *AppService) GetSystemStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// 获取内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		status["memory"] = map[string]interface{}{
			"total": memInfo.Total,
			"used":  memInfo.Used,
			"free":  memInfo.Free,
		}
	}

	// 获取进程信息
	if s.process != nil {
		if memPercent, err := s.process.MemoryPercent(); err == nil {
			status["processMemoryPercent"] = memPercent
		}
	}

	status["startTime"] = s.startTime
	status["uptime"] = time.Since(s.startTime).String()

	return status, nil
}

// GetProjects 获取所有项目
func (s *AppService) GetProjects() ([]entity.Project, error) {
	return s.projectRepo.FindAll()
}

// GetProjectByID 根据 ID 获取项目
func (s *AppService) GetProjectByID(id string) (*entity.Project, error) {
	return s.projectRepo.FindByID(id)
}

// CreateProject 创建项目
func (s *AppService) CreateProject(name, description string) (*entity.Project, error) {
	project := &entity.Project{
		Name:        name,
		Description: description,
		Status:      0, // Active
	}

	if err := s.projectRepo.Save(project); err != nil {
		return nil, err
	}

	return project, nil
}

// UpdateProject 更新项目
func (s *AppService) UpdateProject(id, name, description string, status int) (*entity.Project, error) {
	project, err := s.projectRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	project.Name = name
	project.Description = description
	project.Status = status

	if err := s.projectRepo.Save(project); err != nil {
		return nil, err
	}

	return project, nil
}

// DeleteProject 删除项目
func (s *AppService) DeleteProject(id string) error {
	return s.projectRepo.DeleteSoft(id)
}

// GetWebShells 获取 WebShell 列表
func (s *AppService) GetWebShells(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error) {
	return s.webshellRepo.FindByProjectIDPaginated(projectID, page, pageSize, searchQuery, sortField, sortDir)
}

// GetWebShellByID 根据 ID 获取 WebShell
func (s *AppService) GetWebShellByID(id string) (*entity.WebShell, error) {
	return s.webshellRepo.FindByID(id)
}

// CreateWebShell 创建 WebShell
func (s *AppService) CreateWebShell(
	projectID, url, payload, cryption, encoding, proxyType, remark, status string,
) (*entity.WebShell, error) {
	// 验证项目是否存在
	_, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	webshell := &entity.WebShell{
		ProjectID: projectID,
		Url:       url,
		Payload:   payload,
		Cryption:  cryption,
		Encoding:  encoding,
		ProxyType: proxyType,
		Remark:    remark,
		Status:    status,
	}

	if err := s.webshellRepo.Save(webshell); err != nil {
		return nil, err
	}

	return webshell, nil
}

// UpdateWebShell 更新 WebShell
func (s *AppService) UpdateWebShell(
	id, projectID, url, payload, cryption, encoding, proxyType, remark, status string,
) (*entity.WebShell, error) {
	webshell, err := s.webshellRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	webshell.ProjectID = projectID
	webshell.Url = url
	webshell.Payload = payload
	webshell.Cryption = cryption
	webshell.Encoding = encoding
	webshell.ProxyType = proxyType
	webshell.Remark = remark
	webshell.Status = status

	if err := s.webshellRepo.Save(webshell); err != nil {
		return nil, err
	}

	return webshell, nil
}

// DeleteWebShell 删除 WebShell
func (s *AppService) DeleteWebShell(id string) error {
	return s.webshellRepo.DeleteSoft(id)
}

// RecoverWebShell 恢复已删除的 WebShell
func (s *AppService) RecoverWebShell(id string) error {
	return s.webshellRepo.Recover(id)
}

// camelToSnake 驼峰转下划线（辅助函数）
func camelToSnake(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToLower(result)
}
