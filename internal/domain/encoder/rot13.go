package encoder

import (
	"strings"
)

// ROT13Encoder ROT13 编码器
type ROT13Encoder struct{}

// Name 编码器名称
func (e *ROT13Encoder) Name() string {
	return string(EncoderTypeROT13)
}

// Encode 使用 ROT13 编码数据
func (e *ROT13Encoder) Encode(data []byte) ([]byte, error) {
	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'M':
			return r + 13
		case r >= 'M' && r <= 'Z':
			return r - 13
		case r >= 'a' && r <= 'm':
			return r + 13
		case r >= 'm' && r <= 'z':
			return r - 13
		}
		return r
	}
	
	result := strings.Map(rot13, string(data))
	return []byte(result), nil
}

// Decode 使用 ROT13 解码数据（ROT13 的逆操作与正操作相同）
func (e *ROT13Encoder) Decode(data []byte) ([]byte, error) {
	return e.Encode(data)
}
