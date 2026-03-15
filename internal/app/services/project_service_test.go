package services

import (
	"errors"
	"testing"

	"fg-abyss/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectRepository 模拟项目仓库
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) FindAll() ([]entity.Project, error) {
	args := m.Called()
	return args.Get(0).([]entity.Project), args.Error(1)
}

func (m *MockProjectRepository) FindByID(id string) (*entity.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Project), args.Error(1)
}

func (m *MockProjectRepository) FindByName(name string) (*entity.Project, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Project), args.Error(1)
}

func (m *MockProjectRepository) Save(project *entity.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProjectRepository) DeleteSoft(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProjectService_GetAll(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	expectedProjects := []entity.Project{
		{ID: "1", Name: "项目 1"},
		{ID: "2", Name: "项目 2"},
	}

	mockRepo.On("FindAll").Return(expectedProjects, nil)

	projects, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, projects, 2)
	assert.Equal(t, expectedProjects, projects)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_GetByID(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	expectedProject := &entity.Project{
		ID:          "1",
		Name:        "测试项目",
		Description: "测试描述",
	}

	mockRepo.On("FindByID", "1").Return(expectedProject, nil)

	project, err := service.GetByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, expectedProject, project)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	mockRepo.On("FindByID", "not-found").Return((*entity.Project)(nil), errors.New("项目不存在"))

	project, err := service.GetByID("not-found")

	assert.Error(t, err)
	assert.Nil(t, project)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Create_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	// 模拟名称检查返回不存在
	mockRepo.On("FindByName", "新项目").Return((*entity.Project)(nil), nil)
	// 模拟保存成功
	mockRepo.On("Save", mock.AnythingOfType("*entity.Project")).Return(nil)

	project, err := service.Create("新项目", "新描述")

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "新项目", project.Name)
	assert.Equal(t, "新描述", project.Description)
	assert.Equal(t, 0, project.Status)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Create_NameExists(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	existingProject := &entity.Project{
		ID:   "existing-id",
		Name: "已存在的名称",
	}

	mockRepo.On("FindByName", "已存在的名称").Return(existingProject, nil)

	project, err := service.Create("已存在的名称", "描述")

	assert.Error(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "项目名称已存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Create_SaveError(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	// 模拟名称检查返回不存在
	mockRepo.On("FindByName", "新项目").Return((*entity.Project)(nil), nil)
	// 模拟保存失败
	mockRepo.On("Save", mock.AnythingOfType("*entity.Project")).Return(errors.New("数据库错误"))

	project, err := service.Create("新项目", "新描述")

	assert.Error(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "数据库错误", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Update_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	existingProject := &entity.Project{
		ID:          "1",
		Name:        "旧名称",
		Description: "旧描述",
		Status:      0,
	}

	updatedProject := &entity.Project{
		ID:          "1",
		Name:        "新名称",
		Description: "新描述",
		Status:      1,
	}

	mockRepo.On("FindByID", "1").Return(existingProject, nil)
	mockRepo.On("FindByName", "新名称").Return((*entity.Project)(nil), nil)
	mockRepo.On("Save", updatedProject).Return(nil)

	project, err := service.Update("1", "新名称", "新描述", 1)

	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, "新名称", project.Name)
	assert.Equal(t, "新描述", project.Description)
	assert.Equal(t, 1, project.Status)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	mockRepo.On("FindByID", "not-found").Return((*entity.Project)(nil), errors.New("项目不存在"))

	project, err := service.Update("not-found", "名称", "描述", 0)

	assert.Error(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "项目不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Update_NameExists(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	existingProject := &entity.Project{
		ID:   "1",
		Name: "旧名称",
	}

	duplicateProject := &entity.Project{
		ID:   "2",
		Name: "新名称",
	}

	mockRepo.On("FindByID", "1").Return(existingProject, nil)
	mockRepo.On("FindByName", "新名称").Return(duplicateProject, nil)

	project, err := service.Update("1", "新名称", "描述", 0)

	assert.Error(t, err)
	assert.Nil(t, project)
	assert.Equal(t, "项目名称已存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Delete_Success(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	project := &entity.Project{
		ID:   "1",
		Name: "普通项目",
	}

	mockRepo.On("FindByID", "1").Return(project, nil)
	mockRepo.On("DeleteSoft", "1").Return(nil)

	err := service.Delete("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Delete_NotFound(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	mockRepo.On("FindByID", "not-found").Return((*entity.Project)(nil), errors.New("项目不存在"))

	err := service.Delete("not-found")

	assert.Error(t, err)
	assert.Equal(t, "项目不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProjectService_Delete_DefaultProject(t *testing.T) {
	mockRepo := new(MockProjectRepository)
	service := NewProjectService(mockRepo)

	project := &entity.Project{
		ID:   "default-id",
		Name: "默认项目",
	}

	mockRepo.On("FindByID", "default-id").Return(project, nil)

	err := service.Delete("default-id")

	assert.Error(t, err)
	assert.Equal(t, "默认项目不能删除", err.Error())
	mockRepo.AssertExpectations(t)
}
