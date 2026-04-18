use std::time::Instant;
use zeroize::Zeroizing;

pub struct WebshellSession {
    pub session_id:    String,
    pub session_key:   Zeroizing<[u8; 32]>,
    pub init_time:     Instant,
    pub response_mark: (String, String),
    pub timeout_secs:  u64,
}

impl WebshellSession {
    pub fn is_expired(&self) -> bool {
        self.init_time.elapsed().as_secs() > self.timeout_secs
    }
}
