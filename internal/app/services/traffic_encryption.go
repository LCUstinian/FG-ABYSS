package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// TrafficEncryption 流量加密服务
type TrafficEncryption struct {
	key       []byte
	iv        []byte
	signature []byte
}

// EncryptionConfig 加密配置
type EncryptionConfig struct {
	Key       string `json:"key"`
	IV        string `json:"iv"`
	Signature string `json:"signature"`
}

// NewTrafficEncryption 创建流量加密服务
func NewTrafficEncryption(config *EncryptionConfig) (*TrafficEncryption, error) {
	key, err := hex.DecodeString(config.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	iv, err := hex.DecodeString(config.IV)
	if err != nil {
		return nil, fmt.Errorf("failed to decode IV: %w", err)
	}

	signature, err := hex.DecodeString(config.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}

	return &TrafficEncryption{
		key:       key,
		iv:        iv,
		signature: signature,
	}, nil
}

// EncryptData 加密数据
func (e *TrafficEncryption) EncryptData(plaintext []byte) ([]byte, error) {
	// 创建 AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	// 使用 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// DecryptData 解密数据
func (e *TrafficEncryption) DecryptData(ciphertext []byte) ([]byte, error) {
	// 创建 AES cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	// 使用 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 验证 nonce 长度
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// 分离 nonce 和密文
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GenerateSignature 生成请求签名
func (e *TrafficEncryption) GenerateSignature(data []byte, timestamp int64) string {
	// 组合数据和 timestamp
	message := fmt.Sprintf("%s%d", string(data), timestamp)

	// 使用 HMAC-SHA256 生成签名
	h := hmac.New(sha256.New, e.signature)
	h.Write([]byte(message))
	signature := h.Sum(nil)

	return hex.EncodeToString(signature)
}

// VerifySignature 验证请求签名
func (e *TrafficEncryption) VerifySignature(data []byte, timestamp int64, signature string) bool {
	expectedSignature := e.GenerateSignature(data, timestamp)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// GenerateTimestamp 生成时间戳
func (e *TrafficEncryption) GenerateTimestamp() int64 {
	return time.Now().UnixNano() / 1000000 // 毫秒级时间戳
}

// VerifyTimestamp 验证时间戳（防止重放攻击）
func (e *TrafficEncryption) VerifyTimestamp(timestamp int64, window time.Duration) bool {
	now := e.GenerateTimestamp()
	diff := now - timestamp

	// 检查时间戳是否在有效窗口内
	return diff >= 0 && diff <= int64(window.Milliseconds())
}

// EncryptRequest 加密请求
func (e *TrafficEncryption) EncryptRequest(data []byte) (*EncryptedRequest, error) {
	// 生成时间戳
	timestamp := e.GenerateTimestamp()

	// 生成签名
	signature := e.GenerateSignature(data, timestamp)

	// 加密数据
	encryptedData, err := e.EncryptData(data)
	if err != nil {
		return nil, err
	}

	return &EncryptedRequest{
		Data:      base64.StdEncoding.EncodeToString(encryptedData),
		Timestamp: timestamp,
		Signature: signature,
	}, nil
}

// DecryptRequest 解密请求
func (e *TrafficEncryption) DecryptRequest(req *EncryptedRequest, window time.Duration) ([]byte, error) {
	// 验证时间戳
	if !e.VerifyTimestamp(req.Timestamp, window) {
		return nil, errors.New("timestamp verification failed")
	}

	// 解码数据
	encryptedData, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	// 验证签名
	dataBytes, _ := base64.StdEncoding.DecodeString(req.Data)
	if !e.VerifySignature(dataBytes, req.Timestamp, req.Signature) {
		return nil, errors.New("signature verification failed")
	}

	// 解密数据
	plaintext, err := e.DecryptData(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// EncryptedRequest 加密请求
type EncryptedRequest struct {
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// GenerateKey 生成随机密钥
func GenerateKey() (string, error) {
	key := make([]byte, 32) // 256 位
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// GenerateIV 生成随机 IV
func GenerateIV() (string, error) {
	iv := make([]byte, 12) // GCM 推荐的 nonce 大小
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	return hex.EncodeToString(iv), nil
}

// GenerateSignature 生成随机签名密钥
func GenerateSignature() (string, error) {
	signature := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, signature); err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

// GenerateEncryptionConfig 生成完整的加密配置
func GenerateEncryptionConfig() (*EncryptionConfig, error) {
	key, err := GenerateKey()
	if err != nil {
		return nil, err
	}

	iv, err := GenerateIV()
	if err != nil {
		return nil, err
	}

	signature, err := GenerateSignature()
	if err != nil {
		return nil, err
	}

	return &EncryptionConfig{
		Key:       key,
		IV:        iv,
		Signature: signature,
	}, nil
}

// ValidateEncryptionConfig 验证加密配置
func ValidateEncryptionConfig(config *EncryptionConfig) error {
	if config.Key == "" {
		return errors.New("key is required")
	}

	if config.IV == "" {
		return errors.New("IV is required")
	}

	if config.Signature == "" {
		return errors.New("signature is required")
	}

	// 验证密钥长度（256 位 = 64 个十六进制字符）
	keyBytes, err := hex.DecodeString(config.Key)
	if err != nil {
		return fmt.Errorf("invalid key format: %w", err)
	}
	if len(keyBytes) != 32 {
		return errors.New("key must be 256 bits (64 hex characters)")
	}

	// 验证 IV 长度（GCM nonce = 12 字节 = 24 个十六进制字符）
	ivBytes, err := hex.DecodeString(config.IV)
	if err != nil {
		return fmt.Errorf("invalid IV format: %w", err)
	}
	if len(ivBytes) != 12 {
		return errors.New("IV must be 96 bits (24 hex characters)")
	}

	return nil
}
