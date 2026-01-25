-- Create icd10_codes table
CREATE TABLE IF NOT EXISTS icd10_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- ICD-10 Information
    code VARCHAR(10) NOT NULL UNIQUE,
    description VARCHAR(500) NOT NULL,
    category VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_icd10_codes_code (code),
    INDEX idx_icd10_codes_category (category),
    INDEX idx_icd10_codes_is_active (is_active),
    INDEX idx_icd10_codes_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
