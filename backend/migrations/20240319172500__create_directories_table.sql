CREATE TABLE directories (
    directory_id VARCHAR(80) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    full_path TEXT NOT NULL,
    user_id VARCHAR(80) NOT NULL,
    parent_directory_id VARCHAR(80),
    level int NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

