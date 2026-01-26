-- Create lab_test_templates table
CREATE TABLE IF NOT EXISTS lab_test_templates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Template Information
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(20) NOT NULL,
    description TEXT,
    sample_type VARCHAR(20) NOT NULL,
    preparation_instructions TEXT,
    turnaround_time_hours INT,
    price DECIMAL(10,2),
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_lab_test_templates_code (code),
    INDEX idx_lab_test_templates_name (name),
    INDEX idx_lab_test_templates_category (category),
    INDEX idx_lab_test_templates_is_active (is_active),
    INDEX idx_lab_test_templates_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
