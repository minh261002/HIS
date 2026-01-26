-- Create imaging_templates table
CREATE TABLE IF NOT EXISTS imaging_templates (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Template Information
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    modality VARCHAR(20) NOT NULL,
    body_part VARCHAR(20) NOT NULL,
    description TEXT,
    template_content TEXT,
    price DECIMAL(10,2),
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_imaging_templates_code (code),
    INDEX idx_imaging_templates_name (name),
    INDEX idx_imaging_templates_modality (modality),
    INDEX idx_imaging_templates_is_active (is_active),
    INDEX idx_imaging_templates_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
