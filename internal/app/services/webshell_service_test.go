package services

import (
	"errors"
	"testing"

	"fg-abyss/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWebShellRepository 模拟 WebShell 仓库
type MockWebShellRepository struct {
	mock.Mock
}

func (m *MockWebShellRepository) FindAll() ([]entity.WebShell, error) {
	args := m.Called()
	return args.Get(0).([]entity.WebShell), args.Error(1)
}

func (m *MockWebShellRepository) FindByID(id string) (*entity.WebShell, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.WebShell), args.Error(1)
}

func (m *MockWebShellRepository) FindByProjectID(projectID string) ([]entity.WebShell, error) {
	args := m.Called(projectID)
	return args.Get(0).([]entity.WebShell), args.Error(1)
}

func (m *MockWebShellRepository) FindByProjectIDPaginated(projectID string, page int, pageSize int, searchQuery string, sortField string, sortDir string) ([]entity.WebShell, int64, error) {
	args := m.Called(projectID, page, pageSize, searchQuery, sortField, sortDir)
	return args.Get(0).([]entity.WebShell), args.Get(1).(int64), args.Error(2)
}

func (m *MockWebShellRepository) FindDeletedByProjectID(projectID string) ([]entity.WebShell, error) {
	args := m.Called(projectID)
	return args.Get(0).([]entity.WebShell), args.Error(1)
}

func (m *MockWebShellRepository) Save(webshell *entity.WebShell) error {
	args := m.Called(webshell)
	return args.Error(0)
}

func (m *MockWebShellRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWebShellRepository) DeleteSoft(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockWebShellRepository) Recover(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestWebShellService_GetAll(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	expectedWebShells := []entity.WebShell{
		{ID: "1", Url: "http://example.com/1", ProjectID: "project-1"},
		{ID: "2", Url: "http://example.com/2", ProjectID: "project-1"},
	}

	mockRepo.On("FindByProjectID", "project-1").Return(expectedWebShells, nil)

	webshells, err := service.GetAll("project-1")

	assert.NoError(t, err)
	assert.Len(t, webshells, 2)
	assert.Equal(t, expectedWebShells, webshells)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_GetByID(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	expectedWebShell := &entity.WebShell{
		ID:        "1",
		ProjectID: "project-1",
		Url:       "http://example.com/shell.php",
		Payload:   "php",
		Cryption:  "base64",
		Encoding:  "UTF-8",
		Status:    "active",
	}

	mockRepo.On("FindByID", "1").Return(expectedWebShell, nil)

	webshell, err := service.GetByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, webshell)
	assert.Equal(t, expectedWebShell, webshell)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("FindByID", "not-found").Return((*entity.WebShell)(nil), errors.New("WebShell 不存在"))

	webshell, err := service.GetByID("not-found")

	assert.Error(t, err)
	assert.Nil(t, webshell)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_GetPaginated(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	expectedWebShells := []entity.WebShell{
		{ID: "1", Url: "http://example.com/1"},
		{ID: "2", Url: "http://example.com/2"},
	}
	totalCount := int64(10)

	mockRepo.On("FindByProjectIDPaginated", "project-1", 1, 10, "", "id", "asc").
		Return(expectedWebShells, totalCount, nil)

	webshells, total, err := service.GetPaginated("project-1", 1, 10, "", "id", "asc")

	assert.NoError(t, err)
	assert.Len(t, webshells, 2)
	assert.Equal(t, totalCount, total)
	assert.Equal(t, expectedWebShells, webshells)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Create_Success(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	webshell := &entity.WebShell{
		ID:        "1",
		ProjectID: "project-1",
		Url:       "http://example.com/shell.php",
		Payload:   "php",
		Cryption:  "base64",
	}

	mockRepo.On("Save", webshell).Return(nil)

	result, err := service.Create(webshell)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, webshell, result)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Create_SaveError(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	webshell := &entity.WebShell{
		ID:        "1",
		ProjectID: "project-1",
		Url:       "http://example.com/shell.php",
	}

	mockRepo.On("Save", webshell).Return(errors.New("数据库错误"))

	result, err := service.Create(webshell)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "数据库错误", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Update_Success(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	existingWebShell := &entity.WebShell{
		ID:        "1",
		ProjectID: "project-1",
		Url:       "http://example.com/old.php",
	}

	updatedWebShell := &entity.WebShell{
		ID:        "1",
		ProjectID: "project-1",
		Url:       "http://example.com/new.php",
		Payload:   "php",
		Cryption:  "aes",
	}

	mockRepo.On("FindByID", "1").Return(existingWebShell, nil)
	mockRepo.On("Save", updatedWebShell).Return(nil)

	result, err := service.Update(updatedWebShell)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedWebShell, result)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Update_NotFound(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	webshell := &entity.WebShell{
		ID:        "not-found",
		ProjectID: "project-1",
		Url:       "http://example.com/shell.php",
	}

	mockRepo.On("FindByID", "not-found").Return((*entity.WebShell)(nil), errors.New("WebShell 不存在"))

	result, err := service.Update(webshell)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "WebShell 不存在", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Delete_Success(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("DeleteSoft", "1").Return(nil)

	err := service.Delete("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Delete_Error(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("DeleteSoft", "not-found").Return(errors.New("删除失败"))

	err := service.Delete("not-found")

	assert.Error(t, err)
	assert.Equal(t, "删除失败", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Recover_Success(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("Recover", "1").Return(nil)

	err := service.Recover("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_Recover_Error(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("Recover", "not-found").Return(errors.New("恢复失败"))

	err := service.Recover("not-found")

	assert.Error(t, err)
	assert.Equal(t, "恢复失败", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_GetDeleted(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	expectedWebShells := []entity.WebShell{
		{ID: "1", Url: "http://example.com/deleted1", ProjectID: "project-1"},
		{ID: "2", Url: "http://example.com/deleted2", ProjectID: "project-1"},
	}

	mockRepo.On("FindDeletedByProjectID", "project-1").Return(expectedWebShells, nil)

	webshells, err := service.GetDeleted("project-1")

	assert.NoError(t, err)
	assert.Len(t, webshells, 2)
	assert.Equal(t, expectedWebShells, webshells)
	mockRepo.AssertExpectations(t)
}

func TestWebShellService_GetDeleted_Empty(t *testing.T) {
	mockRepo := new(MockWebShellRepository)
	service := NewWebShellService(mockRepo)

	mockRepo.On("FindDeletedByProjectID", "project-1").Return([]entity.WebShell{}, nil)

	webshells, err := service.GetDeleted("project-1")

	assert.NoError(t, err)
	assert.Empty(t, webshells)
	mockRepo.AssertExpectations(t)
}
