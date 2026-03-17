package services

import (
	"errors"
)

var (
	// ErrConnectionNotFound 连接不存在
	ErrConnectionNotFound = errors.New("connection not found")
	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("connection failed")
	// ErrRequestTimeout 请求超时
	ErrRequestTimeout = errors.New("request timeout")
	// ErrInvalidResponse 响应无效
	ErrInvalidResponse = errors.New("invalid response")
)
