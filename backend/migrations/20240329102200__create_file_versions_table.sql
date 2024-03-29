CREATE TABLE file_versions (
    file_version_id VARCHAR(80) PRIMARY KEY,
    file_id VARCHAR(80) NOT NULL,
    size INT NOT NULL,
    extension VARCHAR(30) NOT NULL,
    user_id VARCHAR(80) NOT NULL,    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (file_id) REFERENCES files(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);