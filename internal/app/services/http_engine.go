package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"fg-abyss/internal/domain/encoder"
)

// HTTPRequestEngine HTTP 请求引擎
type HTTPRequestEngine struct {
	client    *http.Client
	transport *http.Transport
}

// HTTPRequestConfig HTTP 请求配置
type HTTPRequestConfig struct {
	// URL 请求地址
	URL string
	// Method 请求方法（POST/GET）
	Method string
	// Headers 请求头
	Headers map[string]string
	// Cookies Cookie 列表
	Cookies map[string]string
	// Body 请求体
	Body []byte
	// Timeout 超时时间（秒）
	Timeout int
	// ProxyType 代理类型（none, http, https, socks5）
	ProxyType string
	// ProxyAddress 代理地址
	ProxyAddress string
	// SSLVerify 是否验证 SSL 证书
	SSLVerify bool
}

// HTTPResponse HTTP 响应
type HTTPResponse struct {
	// StatusCode HTTP 状态码
	StatusCode int
	// Headers 响应头
	Headers http.Header
	// Body 响应体
	Body []byte
	// Duration 请求耗时（毫秒）
	Duration int64
	// Error 错误信息
	Error string
}

// NewHTTPRequestEngine 创建 HTTP 请求引擎
func NewHTTPRequestEngine() *HTTPRequestEngine {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // 默认验证 SSL
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 默认 30 秒超时
	}

	return &HTTPRequestEngine{
		client:    client,
		transport: transport,
	}
}

// Send 发送 HTTP 请求
func (e *HTTPRequestEngine) Send(ctx context.Context, config *HTTPRequestConfig) (*HTTPResponse, error) {
	if config.URL == "" {
		return nil, errors.New("URL cannot be empty")
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, strings.NewReader(string(config.Body)))
	if err != nil {
		return nil, err
	}

	// 设置默认请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	// 设置自定义请求头
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// 设置 Cookies
	for key, value := range config.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}

	// 配置代理
	if err := e.configureProxy(config); err != nil {
		return nil, err
	}

	// 配置 SSL 验证
	e.transport.TLSClientConfig.InsecureSkipVerify = !config.SSLVerify

	// 设置超时
	if config.Timeout > 0 {
		e.client.Timeout = time.Duration(config.Timeout) * time.Second
	}

	// 发送请求
	startTime := time.Now()
	resp, err := e.client.Do(req)
	duration := time.Since(startTime).Milliseconds()

	if err != nil {
		return &HTTPResponse{
			StatusCode: 0,
			Duration:   duration,
			Error:      err.Error(),
		}, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &HTTPResponse{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Duration:   duration,
			Error:      err.Error(),
		}, err
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
		Duration:   duration,
	}, nil
}

// configureProxy 配置代理
func (e *HTTPRequestEngine) configureProxy(config *HTTPRequestConfig) error {
	if config.ProxyType == "none" || config.ProxyAddress == "" {
		e.transport.Proxy = nil
		return nil
	}

	proxyURL, err := url.Parse(config.ProxyAddress)
	if err != nil {
		return err
	}

	e.transport.Proxy = http.ProxyURL(proxyURL)
	return nil
}

// SetProxy 动态设置代理
func (e *HTTPRequestEngine) SetProxy(proxyType, proxyAddress string) error {
	if proxyType == "none" {
		e.transport.Proxy = nil
		return nil
	}

	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return err
	}

	e.transport.Proxy = http.ProxyURL(proxyURL)
	return nil
}

// SetSSLVerify 设置是否验证 SSL
func (e *HTTPRequestEngine) SetSSLVerify(verify bool) {
	e.transport.TLSClientConfig.InsecureSkipVerify = !verify
}

// SetTimeout 设置超时时间
func (e *HTTPRequestEngine) SetTimeout(seconds int) {
	e.client.Timeout = time.Duration(seconds) * time.Second
}

// MarshalJSON 序列化为 JSON
func (r *HTTPResponse) MarshalJSON() ([]byte, error) {
	type Alias HTTPResponse
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	})
}

// EncodeBody 编码请求体
func EncodeBody(data []byte, encoderType string, key string) ([]byte, error) {
	enc, err := encoder.GetEncoder(encoder.EncoderType(encoderType), key)
	if err != nil {
		return nil, err
	}
	return enc.Encode(data)
}

// DecodeBody 解码响应体
func DecodeBody(data []byte, encoderType string, key string) ([]byte, error) {
	enc, err := encoder.GetEncoder(encoder.EncoderType(encoderType), key)
	if err != nil {
		return nil, err
	}
	return enc.Decode(data)
}
