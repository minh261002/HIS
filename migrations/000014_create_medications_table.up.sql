-- Create medications table
CREATE TABLE IF NOT EXISTS medications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Medication Information
    name VARCHAR(200) NOT NULL,
    generic_name VARCHAR(200),
    dosage_form VARCHAR(20) NOT NULL,
    strength VARCHAR(50),
    unit VARCHAR(20),
    manufacturer VARCHAR(200),
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_medications_name (name),
    INDEX idx_medications_generic_name (generic_name),
    INDEX idx_medications_dosage_form (dosage_form),
    INDEX idx_medications_is_active (is_active),
    INDEX idx_medications_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
