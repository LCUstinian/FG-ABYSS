package services

import (
	"errors"
)

var (
	// ErrTemplateNotFound 模板不存在
	ErrTemplateNotFound = errors.New("template not found")
	// ErrInvalidConfig 配置无效
	ErrInvalidConfig = errors.New("invalid configuration")
	// ErrGenerateFailed 生成失败
	ErrGenerateFailed = errors.New("generate failed")
	// ErrObfuscationFailed 混淆失败
	ErrObfuscationFailed = errors.New("obfuscation failed")
	// ErrEncryptionFailed 加密失败
	ErrEncryptionFailed = errors.New("encryption failed")
)
