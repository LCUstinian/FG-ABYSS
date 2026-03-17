package encoder

// Encoder 编码器接口
type Encoder interface {
	// Encode 编码数据
	Encode(data []byte) ([]byte, error)
	// Decode 解码数据
	Decode(data []byte) ([]byte, error)
	// Name 编码器名称
	Name() string
}

// EncoderType 编码器类型
type EncoderType string

const (
	// EncoderTypeNone 无编码
	EncoderTypeNone EncoderType = "none"
	// EncoderTypeBase64 Base64 编码
	EncoderTypeBase64 EncoderType = "base64"
	// EncoderTypeROT13 ROT13 编码
	EncoderTypeROT13 EncoderType = "rot13"
	// EncoderTypeXOR XOR 编码
	EncoderTypeXOR EncoderType = "xor"
)

// GetEncoder 获取编码器实例
func GetEncoder(encoderType EncoderType, key string) (Encoder, error) {
	switch encoderType {
	case EncoderTypeNone:
		return &NoneEncoder{}, nil
	case EncoderTypeBase64:
		return &Base64Encoder{}, nil
	case EncoderTypeROT13:
		return &ROT13Encoder{}, nil
	case EncoderTypeXOR:
		return NewXOREncoder(key)
	default:
		return nil, ErrUnknownEncoder
	}
}
