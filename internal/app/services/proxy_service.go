package services

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// ProxyType 代理类型
type ProxyType string

const (
	ProxyHTTP   ProxyType = "http"
	ProxyHTTPS  ProxyType = "https"
	ProxySOCKS5 ProxyType = "socks5"
)

// ProxyConfig 代理配置
type ProxyConfig struct {
	Type     ProxyType `json:"type"`
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	Timeout  int       `json:"timeout,omitempty"`
}

// ProxyService 代理服务
type ProxyService struct {
	currentProxy *ProxyConfig
	httpClient   *http.Client
}

// NewProxyService 创建代理服务
func NewProxyService() *ProxyService {
	return &ProxyService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetProxy 设置代理
func (s *ProxyService) SetProxy(config *ProxyConfig) error {
	if config == nil {
		// 移除代理
		s.currentProxy = nil
		s.httpClient.Transport = &http.Transport{
			IdleConnTimeout: 30 * time.Second,
		}
		return nil
	}

	// 验证配置
	if config.Host == "" || config.Port <= 0 {
		return errors.New("invalid proxy configuration")
	}

	s.currentProxy = config

	// 创建代理传输
	transport, err := s.createTransport(config)
	if err != nil {
		return err
	}

	timeout := 30 * time.Second
	if config.Timeout > 0 {
		timeout = time.Duration(config.Timeout) * time.Second
	}

	s.httpClient = &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return nil
}

// createTransport 创建代理传输
func (s *ProxyService) createTransport(config *ProxyConfig) (*http.Transport, error) {
	proxyURL, err := s.buildProxyURL(config)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 注意：生产环境应该验证证书
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return transport, nil
}

// buildProxyURL 构建代理 URL
func (s *ProxyService) buildProxyURL(config *ProxyConfig) (*url.URL, error) {
	var scheme string
	switch config.Type {
	case ProxyHTTP:
		scheme = "http"
	case ProxyHTTPS:
		scheme = "https"
	case ProxySOCKS5:
		scheme = "socks5"
	default:
		return nil, fmt.Errorf("unsupported proxy type: %s", config.Type)
	}

	// 构建用户信息
	var user *url.Userinfo
	if config.Username != "" {
		if config.Password != "" {
			user = url.UserPassword(config.Username, config.Password)
		} else {
			user = url.User(config.Username)
		}
	}

	// 构建 URL
	proxyURL := &url.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:   user,
	}

	return proxyURL, nil
}

// TestProxy 测试代理连接
func (s *ProxyService) TestProxy(config *ProxyConfig) (bool, string, error) {
	// 临时设置代理
	originalProxy := s.currentProxy
	s.SetProxy(config)
	defer func() {
		if originalProxy != nil {
			s.SetProxy(originalProxy)
		} else {
			s.SetProxy(nil)
		}
	}()

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.google.com", nil)
	if err != nil {
		return false, "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Sprintf("Connection failed: %v", err), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, "Proxy connection successful", nil
	}

	return false, fmt.Sprintf("Unexpected status code: %d", resp.StatusCode), nil
}

// GetHTTPClient 获取配置了代理的 HTTP 客户端
func (s *ProxyService) GetHTTPClient() *http.Client {
	return s.httpClient
}

// GetCurrentProxy 获取当前代理配置
func (s *ProxyService) GetCurrentProxy() *ProxyConfig {
	return s.currentProxy
}

// IsProxyEnabled 检查是否启用了代理
func (s *ProxyService) IsProxyEnabled() bool {
	return s.currentProxy != nil
}

// ClearProxy 清除代理配置
func (s *ProxyService) ClearProxy() {
	s.SetProxy(nil)
}

// DialWithProxy 使用代理建立连接
func (s *ProxyService) DialWithProxy(network, address string) (net.Conn, error) {
	if s.currentProxy == nil {
		return net.DialTimeout(network, address, 30*time.Second)
	}

	// 对于 SOCKS5，使用专门的拨号器
	if s.currentProxy.Type == ProxySOCKS5 {
		proxyDialer, err := s.createSOCKS5Dialer()
		if err != nil {
			return nil, err
		}
		return proxyDialer.Dial(network, address)
	}

	// 对于 HTTP/HTTPS 代理，使用 net.Dial 然后建立隧道
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.currentProxy.Host, s.currentProxy.Port), 30*time.Second)
	if err != nil {
		return nil, err
	}

	// 如果是 HTTPS，需要建立 CONNECT 隧道
	if network == "tcp" && s.currentProxy.Type == ProxyHTTPS {
		// TODO: 实现 CONNECT 方法
	}

	return conn, nil
}

// createSOCKS5Dialer 创建 SOCKS5 拨号器
func (s *ProxyService) createSOCKS5Dialer() (proxy.Dialer, error) {
	if s.currentProxy == nil || s.currentProxy.Type != ProxySOCKS5 {
		return nil, errors.New("not a SOCKS5 proxy")
	}

	auth := &proxy.Auth{}
	if s.currentProxy.Username != "" {
		auth.User = s.currentProxy.Username
		auth.Password = s.currentProxy.Password
	}

	dialer, err := proxy.SOCKS5("tcp", 
		fmt.Sprintf("%s:%d", s.currentProxy.Host, s.currentProxy.Port), 
		auth, 
		proxy.Direct)
	if err != nil {
		return nil, err
	}

	return dialer, nil
}

// GetProxyStats 获取代理统计信息
func (s *ProxyService) GetProxyStats() map[string]interface{} {
	stats := make(map[string]interface{})
	
	stats["enabled"] = s.IsProxyEnabled()
	
	if s.currentProxy != nil {
		stats["type"] = string(s.currentProxy.Type)
		stats["host"] = s.currentProxy.Host
		stats["port"] = s.currentProxy.Port
		stats["has_auth"] = s.currentProxy.Username != ""
	} else {
		stats["type"] = nil
		stats["host"] = nil
		stats["port"] = nil
		stats["has_auth"] = false
	}

	return stats
}

// ValidateProxyConfig 验证代理配置
func (s *ProxyService) ValidateProxyConfig(config *ProxyConfig) error {
	if config.Host == "" {
		return errors.New("proxy host is required")
	}

	if config.Port <= 0 || config.Port > 65535 {
		return errors.New("proxy port must be between 1 and 65535")
	}

	switch config.Type {
	case ProxyHTTP, ProxyHTTPS, ProxySOCKS5:
		// 有效
	default:
		return fmt.Errorf("unsupported proxy type: %s", config.Type)
	}

	if config.Timeout < 0 {
		return errors.New("timeout cannot be negative")
	}

	return nil
}
