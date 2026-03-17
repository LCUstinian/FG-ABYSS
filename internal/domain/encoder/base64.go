package encoder

import (
	"encoding/base64"
)

// Base64Encoder Base64 编码器
type Base64Encoder struct{}

// Name 编码器名称
func (e *Base64Encoder) Name() string {
	return string(EncoderTypeBase64)
}

// Encode 编码数据为 Base64
func (e *Base64Encoder) Encode(data []byte) ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	return []byte(encoded), nil
}

// Decode 从 Base64 解码数据
func (e *Base64Encoder) Decode(data []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, ErrDecodeFailed
	}
	return decoded, nil
}
