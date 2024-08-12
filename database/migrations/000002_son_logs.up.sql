-- 000002_son_execution_logs.up.sql
CREATE TABLE son_execution_logs (
    id VARCHAR(36) PRIMARY KEY,
    son_id VARCHAR(36) NOT NULL,
    webhook_log_id VARCHAR(36) NOT NULL,
    execution_status VARCHAR(10) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    error_message TEXT,
    FOREIGN KEY (son_id) REFERENCES sons(id) ON DELETE CASCADE,
    FOREIGN KEY (webhook_log_id) REFERENCES webhook_logs(id) ON DELETE CASCADE,
    CONSTRAINT chk_execution_status CHECK (execution_status IN ('success', 'failure'))
);
