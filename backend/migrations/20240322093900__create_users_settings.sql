CREATE TABLE user_settings (
    user_setting_id VARCHAR(80) PRIMARY KEY,
    user_id VARCHAR(80) NOT NULL UNIQUE,
    storage_vendor VARCHAR(255),
    storage_credentials JSON,
    storage_informations JSON,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);