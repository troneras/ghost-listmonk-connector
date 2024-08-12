CREATE TABLE son_execution_action_logs (
    id VARCHAR(36) PRIMARY KEY,
    son_execution_log_id VARCHAR(36) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    action_status VARCHAR(10) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    error_message TEXT,
    FOREIGN KEY (son_execution_log_id) REFERENCES son_execution_logs(id) ON DELETE CASCADE,
    CONSTRAINT chk_action_status CHECK (action_status IN ('success', 'failure'))
);