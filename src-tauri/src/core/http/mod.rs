use reqwest::{Client, Proxy, Response};
use std::collections::HashMap;
use std::time::Duration;
use tokio::sync::Mutex;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum HttpError {
    #[error("HTTP 请求失败：{0}")]
    RequestFailed(#[from] reqwest::Error),
    #[error("无效的 URL: {0}")]
    InvalidUrl(String),
    #[error("请求超时")]
    Timeout,
    #[error("无效的响应状态码：{0}")]
    InvalidStatusCode(u16),
    #[error("代理配置错误：{0}")]
    ProxyError(String),
}

pub type Result<T> = std::result::Result<T, HttpError>;

/// HTTP 客户端配置
#[derive(Debug, Clone)]
pub struct HttpClientConfig {
    /// 请求超时时间（秒）
    pub timeout: u64,
    /// 连接超时时间（秒）
    pub connect_timeout: u64,
    /// 最大重试次数
    pub max_retries: u32,
    /// 重试延迟（毫秒）
    pub retry_delay_ms: u64,
    /// 速率限制（每秒请求数）
    pub rate_limit: Option<u32>,
    /// 代理 URL
    pub proxy_url: Option<String>,
    /// User-Agent
    pub user_agent: String,
    /// 是否验证 SSL 证书
    pub verify_ssl: bool,
}

impl Default for HttpClientConfig {
    fn default() -> Self {
        Self {
            timeout: 30,
            connect_timeout: 10,
            max_retries: 3,
            retry_delay_ms: 1000,
            rate_limit: Some(10),
            proxy_url: None,
            user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36".to_string(),
            verify_ssl: false,
        }
    }
}

/// HTTP 客户端
pub struct HttpClient {
    client: Client,
    config: HttpClientConfig,
    cookies: Mutex<HashMap<String, String>>,
}

impl HttpClient {
    /// 创建新的 HTTP 客户端
    pub fn new(config: HttpClientConfig) -> Result<Self> {
        let mut client_builder = Client::builder()
            .timeout(Duration::from_secs(config.timeout))
            .connect_timeout(Duration::from_secs(config.connect_timeout))
            .user_agent(&config.user_agent)
            .danger_accept_invalid_certs(!config.verify_ssl);

        // 配置代理
        if let Some(proxy_url) = &config.proxy_url {
            let proxy = Proxy::all(proxy_url)
                .map_err(|e| HttpError::ProxyError(e.to_string()))?;
            client_builder = client_builder.proxy(proxy);
        }

        let client = client_builder
            .build()
            .map_err(HttpError::from)?;

        Ok(Self {
            client,
            config,
            cookies: Mutex::new(HashMap::new()),
        })
    }

    /// 发送 GET 请求
    pub async fn get(&self, url: &str, headers: Option<HashMap<String, String>>) -> Result<HttpResponse> {
        self.request("GET", url, None, headers).await
    }

    /// 发送 POST 请求
    pub async fn post(&self, url: &str, body: Option<Vec<u8>>, headers: Option<HashMap<String, String>>) -> Result<HttpResponse> {
        self.request("POST", url, body, headers).await
    }

    /// 发送通用请求
    pub async fn request(
        &self,
        method: &str,
        url: &str,
        body: Option<Vec<u8>>,
        headers: Option<HashMap<String, String>>,
    ) -> Result<HttpResponse> {
        let mut retries = 0;
        let mut last_error: Option<HttpError> = None;

        while retries <= self.config.max_retries {
            // 应用速率限制
            if let Some(rate_limit) = self.config.rate_limit {
                tokio::time::sleep(Duration::from_millis(1000 / rate_limit as u64)).await;
            }

            let result = self.execute_request(method, url, &body, &headers).await;

            match result {
                Ok(response) => return Ok(response),
                Err(e) => {
                    last_error = Some(e);
                    retries += 1;
                    if retries <= self.config.max_retries {
                        tokio::time::sleep(Duration::from_millis(self.config.retry_delay_ms)).await;
                    }
                }
            }
        }

        Err(last_error.unwrap())
    }

    /// 执行单个请求
    async fn execute_request(
        &self,
        method: &str,
        url: &str,
        body: &Option<Vec<u8>>,
        headers: &Option<HashMap<String, String>>,
    ) -> Result<HttpResponse> {
        let mut request_builder = match method.to_uppercase().as_str() {
            "GET" => self.client.get(url),
            "POST" => self.client.post(url),
            "PUT" => self.client.put(url),
            "DELETE" => self.client.delete(url),
            "HEAD" => self.client.head(url),
            "PATCH" => self.client.patch(url),
            _ => return Err(HttpError::InvalidStatusCode(405)),
        };

        // 添加自定义 headers
        if let Some(headers_map) = headers {
            for (key, value) in headers_map {
                request_builder = request_builder.header(key, value);
            }
        }

        // 添加 Cookies
        let cookies = self.cookies.lock().await;
        if !cookies.is_empty() {
            let cookie_header = cookies
                .iter()
                .map(|(k, v)| format!("{}={}", k, v))
                .collect::<Vec<_>>()
                .join("; ");
            drop(cookies);
            request_builder = request_builder.header("Cookie", cookie_header);
        }

        // 添加请求体
        if let Some(body_data) = body {
            request_builder = request_builder.body(body_data.clone());
        }

        let response: Response = request_builder.send().await?;
        let status = response.status().as_u16();
        let response_headers: HashMap<String, String> = response
            .headers()
            .iter()
            .map(|(k, v)| (k.to_string(), v.to_str().unwrap_or("").to_string()))
            .collect();

        // 更新 Cookies
        self.update_cookies_from_response(&response).await;

        let body_bytes: Vec<u8> = response.bytes().await?.to_vec();

        Ok(HttpResponse {
            status_code: status,
            headers: response_headers,
            body: body_bytes,
        })
    }

    /// 从响应中更新 Cookies
    async fn update_cookies_from_response(&self, response: &Response) {
        let mut cookies = self.cookies.lock().await;
        for cookie in response.cookies() {
            let name = cookie.name();
            let value = cookie.value();
            cookies.insert(name.to_string(), value.to_string());
        }
    }

    /// 清除所有 Cookies
    pub async fn clear_cookies(&self) {
        self.cookies.lock().await.clear();
    }

    /// 获取当前配置
    pub fn config(&self) -> &HttpClientConfig {
        &self.config
    }
}

/// HTTP 响应
#[derive(Debug, Clone)]
pub struct HttpResponse {
    pub status_code: u16,
    pub headers: HashMap<String, String>,
    pub body: Vec<u8>,
}

impl HttpResponse {
    /// 检查响应是否成功
    pub fn is_success(&self) -> bool {
        self.status_code >= 200 && self.status_code < 300
    }

    /// 获取响应体字符串
    pub fn text(&self) -> Result<String> {
        String::from_utf8(self.body.clone())
            .map_err(|e| HttpError::InvalidUrl(format!("UTF-8 解码失败：{}", e)))
    }

    /// 检查是否包含特定状态码
    pub fn has_status(&self, status: u16) -> bool {
        self.status_code == status
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_config() {
        let config = HttpClientConfig::default();
        assert_eq!(config.timeout, 30);
        assert_eq!(config.connect_timeout, 10);
        assert_eq!(config.max_retries, 3);
        assert!(!config.verify_ssl);
    }

    #[test]
    fn test_http_response() {
        let response = HttpResponse {
            status_code: 200,
            headers: HashMap::new(),
            body: b"Hello, World!".to_vec(),
        };
        assert!(response.is_success());
        assert_eq!(response.text().unwrap(), "Hello, World!");
    }
}
