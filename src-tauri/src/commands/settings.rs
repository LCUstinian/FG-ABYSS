use tauri::State;
use crate::Result;
use crate::state::{AppState, check_locked};
use crate::infra::config::Config;

#[tauri::command]
#[specta::specta]
pub async fn get_settings(state: State<'_, AppState>) -> Result<Config> {
    Ok(state.settings_service.get().await)
}

#[derive(Debug, serde::Deserialize, specta::Type)]
#[serde(rename_all = "camelCase")]
pub struct UpdateSettingsInput {
    pub theme:        Option<String>,
    pub language:     Option<String>,
    pub accent_color: Option<String>,
}

#[tauri::command]
#[specta::specta]
pub async fn update_settings(
    state: State<'_, AppState>,
    input: UpdateSettingsInput,
) -> Result<Config> {
    check_locked(&state)?;
    let current = state.settings_service.get().await;
    let theme   = input.theme.unwrap_or(current.appearance.theme);
    let lang    = input.language.unwrap_or(current.appearance.language);
    let accent  = input.accent_color.unwrap_or(current.appearance.accent_color);
    state.settings_service.update_appearance(theme, lang, accent).await?;
    Ok(state.settings_service.get().await)
}

#[tauri::command]
#[specta::specta]
pub async fn unlock(
    state:    State<'_, AppState>,
    password: String,
) -> Result<bool> {
    let cfg = state.settings_service.get().await;
    if !cfg.security.master_password_enabled {
        return Ok(true);
    }
    use argon2::{Argon2, PasswordHash, PasswordVerifier};
    let hash = cfg.security.master_password_hash;
    let ok = PasswordHash::new(&hash)
        .map(|h| Argon2::default().verify_password(password.as_bytes(), &h).is_ok())
        .unwrap_or(false);
    if ok {
        state.is_locked.store(false, std::sync::atomic::Ordering::Relaxed);
    }
    Ok(ok)
}

#[tauri::command]
#[specta::specta]
pub async fn set_master_password(
    state:        State<'_, AppState>,
    old_password: Option<String>,
    new_password: String,
) -> Result<()> {
    check_locked(&state)?;
    use argon2::{Argon2, PasswordHash, PasswordHasher, PasswordVerifier};
    use argon2::password_hash::{rand_core::OsRng, SaltString};

    let cfg = state.settings_service.get().await;
    if cfg.security.master_password_enabled {
        match old_password.as_deref() {
            Some(old) => {
                let hash_str = &cfg.security.master_password_hash;
                let ok = PasswordHash::new(hash_str)
                    .map(|h| Argon2::default().verify_password(old.as_bytes(), &h).is_ok())
                    .unwrap_or(false);
                if !ok {
                    return Err(crate::AppError::InvalidInput("old password incorrect".into()));
                }
            }
            None => return Err(crate::AppError::InvalidInput("old password required".into())),
        }
    }
    let salt = SaltString::generate(&mut OsRng);
    let new_hash = Argon2::default()
        .hash_password(new_password.as_bytes(), &salt)
        .map(|h| h.to_string())
        .map_err(|e| crate::AppError::Crypto(e.to_string()))?;

    let path = state.paths.config.clone();
    let mut full_cfg = state.settings_service.get().await;
    full_cfg.security.master_password_hash    = new_hash;
    full_cfg.security.master_password_enabled = true;
    crate::infra::config::save(&path, &full_cfg)?;
    state.is_locked.store(true, std::sync::atomic::Ordering::Relaxed);
    Ok(())
}
