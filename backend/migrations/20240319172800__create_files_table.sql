CREATE TABLE files (
    file_id VARCHAR(80) PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    size INT NOT NULL,
    extension VARCHAR(10) NOT NULL,
    user_id VARCHAR(80) NOT NULL,
    directory_id VARCHAR(80) NOT NULL,
    full_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (directory_id) REFERENCES directories(directory_id)
);