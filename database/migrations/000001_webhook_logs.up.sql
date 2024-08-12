CREATE TABLE webhook_logs (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    headers TEXT NOT NULL,
    body TEXT NOT NULL,
    status_code INT NOT NULL,
    response_body TEXT,
    duration INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);