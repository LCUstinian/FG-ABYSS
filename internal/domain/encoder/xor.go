package encoder

// XOREncoder XOR 编码器
type XOREncoder struct {
	key []byte
}

// NewXOREncoder 创建新的 XOR 编码器
func NewXOREncoder(key string) (*XOREncoder, error) {
	if key == "" {
		return nil, ErrInvalidKey
	}
	return &XOREncoder{key: []byte(key)}, nil
}

// Name 编码器名称
func (e *XOREncoder) Name() string {
	return string(EncoderTypeXOR)
}

// Encode 使用 XOR 编码数据
func (e *XOREncoder) Encode(data []byte) ([]byte, error) {
	return e.xor(data), nil
}

// Decode 使用 XOR 解码数据（XOR 的逆操作与正操作相同）
func (e *XOREncoder) Decode(data []byte) ([]byte, error) {
	return e.xor(data), nil
}

// xor 执行 XOR 操作
func (e *XOREncoder) xor(data []byte) []byte {
	result := make([]byte, len(data))
	keyLen := len(e.key)
	
	for i := 0; i < len(data); i++ {
		result[i] = data[i] ^ e.key[i%keyLen]
	}
	
	return result
}
