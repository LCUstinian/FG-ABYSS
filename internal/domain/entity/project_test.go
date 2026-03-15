package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name        string
		project     *Project
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效项目",
			project: &Project{
				Name:        "测试项目",
				Description: "测试描述",
				Status:      0,
			},
			expectError: false,
		},
		{
			name: "名称为空",
			project: &Project{
				Name:        "",
				Description: "测试描述",
			},
			expectError: true,
			errorMsg:    "项目名称不能为空",
		},
		{
			name: "只有空格",
			project: &Project{
				Name:        "   ",
				Description: "测试描述",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProject_ValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "name",
		Message: "项目名称不能为空",
	}

	assert.Equal(t, "name", err.Field)
	assert.Equal(t, "项目名称不能为空", err.Message)
	assert.Equal(t, "项目名称不能为空", err.Error())
}
