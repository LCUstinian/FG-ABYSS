package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebShell_Validate(t *testing.T) {
	tests := []struct {
		name        string
		webshell    *WebShell
		expectError bool
		errorMsg    string
	}{
		{
			name: "有效 WebShell",
			webshell: &WebShell{
				ID:        "test-id",
				ProjectID: "project-123",
				Url:       "http://example.com/shell.php",
				Payload:   "php",
				Cryption:  "base64",
				Encoding:  "UTF-8",
				ProxyType: "http",
				Remark:    "测试 WebShell",
				Status:    "active",
			},
			expectError: false,
		},
		{
			name: "URL 为空",
			webshell: &WebShell{
				ID:        "test-id",
				ProjectID: "project-123",
				Url:       "",
			},
			expectError: true,
			errorMsg:    "URL 不能为空",
		},
		{
			name: "项目 ID 为空",
			webshell: &WebShell{
				ID:  "test-id",
				Url: "http://example.com/shell.php",
			},
			expectError: true,
			errorMsg:    "项目 ID 不能为空",
		},
		{
			name: "URL 和项目 ID 都为空",
			webshell: &WebShell{
				ID: "test-id",
			},
			expectError: true,
			errorMsg:    "URL 不能为空",
		},
		{
			name: "最小有效 WebShell",
			webshell: &WebShell{
				ID:        "test-id",
				ProjectID: "project-123",
				Url:       "http://example.com",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.webshell.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
