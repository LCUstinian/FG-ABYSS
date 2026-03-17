package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptionModule 提供加密算法支持
type EncryptionModule struct{}

// NewEncryptionModule 创建加密模块实例
func NewEncryptionModule() *EncryptionModule {
	return &EncryptionModule{}
}

// EncryptAES 使用 AES-256-CBC 加密数据
func (em *EncryptionModule) EncryptAES(plaintext string, password string) (string, error) {
	// 从密码派生密钥
	key := pbkdf2.Key([]byte(password), []byte("fg-abyss-salt"), 10000, 32, nil)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 生成随机 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// PKCS7 填充
	plaintextBytes := []byte(plaintext)
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	if padding > 0 {
		paddingBytes := make([]byte, padding)
		for i := range paddingBytes {
			paddingBytes[i] = byte(padding)
		}
		plaintextBytes = append(plaintextBytes, paddingBytes...)
	}

	// CBC 模式加密
	ciphertext := make([]byte, len(plaintextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// 拼接 IV 和密文
	result := append(iv, ciphertext...)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(result), nil
}

// DecryptAES 使用 AES-256-CBC 解密数据
func (em *EncryptionModule) DecryptAES(ciphertext string, password string) (string, error) {
	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// 分离 IV 和密文
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	ciphertextBytes := data[aes.BlockSize:]

	// 从密码派生密钥
	key := pbkdf2.Key([]byte(password), []byte("fg-abyss-salt"), 10000, 32, nil)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)

	// 移除 PKCS7 填充
	padding := int(ciphertextBytes[len(ciphertextBytes)-1])
	if padding > len(ciphertextBytes) || padding > aes.BlockSize {
		return "", errors.New("invalid padding")
	}

	// 验证填充
	for i := 0; i < padding; i++ {
		if ciphertextBytes[len(ciphertextBytes)-1-i] != byte(padding) {
			return "", errors.New("invalid padding")
		}
	}

	return string(ciphertextBytes[:len(ciphertextBytes)-padding]), nil
}

// EncryptXOR 简单的 XOR 加密（用于快速混淆）
func (em *EncryptionModule) EncryptXOR(data string, key byte) string {
	bytes := []byte(data)
	for i := range bytes {
		bytes[i] ^= key
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

// DecryptXOR 解密 XOR 加密的数据
func (em *EncryptionModule) DecryptXOR(encoded string, key byte) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	for i := range bytes {
		bytes[i] ^= key
	}
	return string(bytes), nil
}

// GenerateKey 生成随机密钥
func (em *EncryptionModule) GenerateKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// HashPassword 使用 PBKDF2 哈希密码
func (em *EncryptionModule) HashPassword(password string, salt string) string {
	key := pbkdf2.Key([]byte(password), []byte(salt), 10000, 32, nil)
	return base64.StdEncoding.EncodeToString(key)
}

// VerifyPassword 验证密码哈希
func (em *EncryptionModule) VerifyPassword(password string, salt string, hash string) bool {
	computedHash := em.HashPassword(password, salt)
	return computedHash == hash
}

// EncryptROT13 简单的 ROT13 编码（仅用于混淆）
func (em *EncryptionModule) EncryptROT13(input string) string {
	result := make([]byte, len(input))
	for i, b := range []byte(input) {
		switch {
		case b >= 'A' && b <= 'Z':
			result[i] = 'A' + (b-'A'+13)%26
		case b >= 'a' && b <= 'z':
			result[i] = 'a' + (b-'a'+13)%26
		default:
			result[i] = b
		}
	}
	return string(result)
}

// DecryptROT13 解密 ROT13 编码（ROT13 是对称的）
func (em *EncryptionModule) DecryptROT13(input string) string {
	return em.EncryptROT13(input)
}
