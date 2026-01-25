-- Create patient_allergies table
CREATE TABLE IF NOT EXISTS patient_allergies (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    patient_id BIGINT UNSIGNED NOT NULL,
    
    -- Allergy Information
    allergen VARCHAR(100) NOT NULL,
    allergen_type VARCHAR(20) NOT NULL,
    reaction TEXT,
    severity VARCHAR(20) NOT NULL,
    diagnosed_date DATE,
    notes TEXT,
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_patient_allergies_patient_id (patient_id),
    INDEX idx_patient_allergies_allergen_type (allergen_type),
    INDEX idx_patient_allergies_severity (severity),
    INDEX idx_patient_allergies_is_active (is_active),
    INDEX idx_patient_allergies_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
