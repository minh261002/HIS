-- Create patient_medical_history table
CREATE TABLE IF NOT EXISTS patient_medical_history (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    patient_id BIGINT UNSIGNED NOT NULL,
    
    -- Condition Information
    condition_name VARCHAR(200) NOT NULL,
    condition_type VARCHAR(20) NOT NULL,
    diagnosis_date DATE,
    status VARCHAR(20) NOT NULL,
    treatment TEXT,
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
    INDEX idx_patient_medical_history_patient_id (patient_id),
    INDEX idx_patient_medical_history_condition_type (condition_type),
    INDEX idx_patient_medical_history_status (status),
    INDEX idx_patient_medical_history_diagnosis_date (diagnosis_date),
    INDEX idx_patient_medical_history_is_active (is_active),
    INDEX idx_patient_medical_history_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
