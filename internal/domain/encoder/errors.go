package encoder

import "errors"

var (
	// ErrUnknownEncoder 未知的编码器类型
	ErrUnknownEncoder = errors.New("unknown encoder type")
	// ErrInvalidKey 无效的密钥
	ErrInvalidKey = errors.New("invalid encryption key")
	// ErrDecodeFailed 解码失败
	ErrDecodeFailed = errors.New("decode failed")
)
