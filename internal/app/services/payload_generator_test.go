package services

import (
	"fg-abyss/internal/domain/entity"
	"strings"
	"testing"
)

func TestPayloadGenerator_Generate(t *testing.T) {
	generator := NewPayloadGenerator()
	
	tests := []struct {
		name     string
		config   *entity.PayloadConfig
		wantErr  bool
		contains string
	}{
		{
			name: "PHP Basic Payload",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypePHP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "test123",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  false,
			contains: "test123",
		},
		{
			name: "PHP Full Payload",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypePHP,
				Function:         entity.PayloadFunctionFull,
				Password:         "mypassword",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  false,
			contains: "mypassword",
		},
		{
			name: "ASP Basic Payload",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypeASP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "asppass",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  false,
			contains: "asppass",
		},
		{
			name: "ASPX Basic Payload",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypeASPX,
				Function:         entity.PayloadFunctionBasic,
				Password:         "aspxpass",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  false,
			contains: "aspxpass",
		},
		{
			name: "JSP Basic Payload",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypeJSP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "jsppass",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  false,
			contains: "jsppass",
		},
		{
			name: "Missing Password",
			config: &entity.PayloadConfig{
				Type:             entity.PayloadTypePHP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "",
				Encoder:          "none",
				ObfuscationLevel: entity.ObfuscationNone,
			},
			wantErr:  true,
			contains: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.Generate(tt.config)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if !result.Success {
					t.Errorf("Generate() result.Success = false, want true")
				}
				
				if tt.contains != "" && !strings.Contains(result.Content, tt.contains) {
					t.Errorf("Generate() content does not contain '%s': got %s", tt.contains, result.Content)
				}
				
				if result.Size == 0 {
					t.Errorf("Generate() result.Size = 0, want non-zero")
				}
			}
		})
	}
}

func TestPayloadGenerator_TemplateRendering(t *testing.T) {
	generator := NewPayloadGenerator()
	
	// Test PHP Basic template rendering
	config := &entity.PayloadConfig{
		Type:             entity.PayloadTypePHP,
		Function:         entity.PayloadFunctionBasic,
		Password:         "custompass",
		Encoder:          "none",
		ObfuscationLevel: entity.ObfuscationNone,
	}
	
	result, err := generator.Generate(config)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	
	// Check that password was properly substituted
	if !strings.Contains(result.Content, "custompass") {
		t.Errorf("Template did not substitute password. Got: %s", result.Content)
	}
	
	// Check that template variables were not left as placeholders
	if strings.Contains(result.Content, "{{.Password}}") || strings.Contains(result.Content, "{{.password}}") {
		t.Errorf("Template variables were not rendered. Got: %s", result.Content)
	}
}

func TestPayloadGenerator_Obfuscation(t *testing.T) {
	generator := NewPayloadGenerator()
	
	tests := []struct {
		name   string
		level  entity.ObfuscationLevel
		expect string
	}{
		{
			name:   "No Obfuscation",
			level:  entity.ObfuscationNone,
			expect: "<?php",
		},
		{
			name:   "Low Obfuscation",
			level:  entity.ObfuscationLow,
			expect: "<?php",
		},
		{
			name:   "Medium Obfuscation",
			level:  entity.ObfuscationMedium,
			expect: "<?php",
		},
		{
			name:   "High Obfuscation",
			level:  entity.ObfuscationHigh,
			expect: "<?php",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &entity.PayloadConfig{
				Type:             entity.PayloadTypePHP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "test123",
				Encoder:          "none",
				ObfuscationLevel: tt.level,
			}
			
			result, err := generator.Generate(config)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			
			if !strings.Contains(result.Content, tt.expect) {
				t.Errorf("Obfuscation %s: expected to contain %s, got %s", tt.level, tt.expect, result.Content)
			}
		})
	}
}

func TestPayloadGenerator_Encoding(t *testing.T) {
	generator := NewPayloadGenerator()
	
	tests := []struct {
		name    string
		encoder string
	}{
		{name: "Base64", encoder: "base64"},
		{name: "ROT13", encoder: "rot13"},
		{name: "URL", encoder: "urlencode"},
		{name: "Hex", encoder: "hex"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &entity.PayloadConfig{
				Type:             entity.PayloadTypePHP,
				Function:         entity.PayloadFunctionBasic,
				Password:         "test123",
				Encoder:          tt.encoder,
				ObfuscationLevel: entity.ObfuscationNone,
			}
			
			result, err := generator.Generate(config)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			
			if result.Size == 0 {
				t.Errorf("Encoded payload size is 0")
			}
		})
	}
}
