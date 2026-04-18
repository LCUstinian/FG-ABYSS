use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct C2Profile {
    pub name:            String,
    pub request_param:   String,
    pub extra_headers:   HashMap<String, String>,
    pub user_agent:      String,
    pub request_wrapper: Option<String>,
    pub response_prefix: String,
    pub jitter_ms:       (u64, u64),
    pub padding_range:   Option<(usize, usize)>,
}

impl Default for C2Profile {
    fn default() -> Self {
        Self {
            name:            "default".into(),
            request_param:   "pass".into(),
            extra_headers:   HashMap::new(),
            user_agent:      "Mozilla/5.0".into(),
            request_wrapper: None,
            response_prefix: String::new(),
            jitter_ms:       (0, 0),
            padding_range:   None,
        }
    }
}

impl C2Profile {
    pub fn by_name(name: &str) -> Self {
        match name {
            "cdn-callback" => Self {
                name:          "cdn-callback".into(),
                request_param: "data".into(),
                extra_headers: [("Content-Type".into(), "application/octet-stream".into())].into(),
                ..Default::default()
            },
            "api-json" => Self {
                name:            "api-json".into(),
                request_param:   "data".into(),
                request_wrapper: Some(r#"{"code":0,"data":"{DATA}"}"#.into()),
                ..Default::default()
            },
            "form-submit" => Self {
                name:          "form-submit".into(),
                request_param: "_token".into(),
                ..Default::default()
            },
            _ => Self::default(),
        }
    }
}
