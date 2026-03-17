package encoder

// NoneEncoder 无编码器
type NoneEncoder struct{}

// Name 编码器名称
func (e *NoneEncoder) Name() string {
	return string(EncoderTypeNone)
}

// Encode 编码数据（不做任何处理）
func (e *NoneEncoder) Encode(data []byte) ([]byte, error) {
	return data, nil
}

// Decode 解码数据（不做任何处理）
func (e *NoneEncoder) Decode(data []byte) ([]byte, error) {
	return data, nil
}
