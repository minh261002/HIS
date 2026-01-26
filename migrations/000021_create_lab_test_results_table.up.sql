-- Create lab_test_results table
CREATE TABLE IF NOT EXISTS lab_test_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key
    request_id BIGINT UNSIGNED NOT NULL,
    
    -- Result Information
    parameter_name VARCHAR(100) NOT NULL,
    value VARCHAR(100),
    unit VARCHAR(50),
    normal_range_text VARCHAR(100),
    is_abnormal BOOLEAN DEFAULT FALSE,
    remarks TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_lab_test_results_request_id (request_id),
    INDEX idx_lab_test_results_deleted_at (deleted_at),
    
    -- Foreign Key
    FOREIGN KEY (request_id) REFERENCES lab_test_requests(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
