-- Create lab_test_template_parameters table
CREATE TABLE IF NOT EXISTS lab_test_template_parameters (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key
    template_id BIGINT UNSIGNED NOT NULL,
    
    -- Parameter Information
    parameter_name VARCHAR(100) NOT NULL,
    unit VARCHAR(50),
    normal_range_min DECIMAL(10,2),
    normal_range_max DECIMAL(10,2),
    normal_range_text VARCHAR(100),
    display_order INT DEFAULT 0,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_lab_test_template_parameters_template_id (template_id),
    INDEX idx_lab_test_template_parameters_deleted_at (deleted_at),
    
    -- Foreign Key
    FOREIGN KEY (template_id) REFERENCES lab_test_templates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
