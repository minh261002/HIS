-- Create diagnoses table
CREATE TABLE IF NOT EXISTS diagnoses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Keys
    visit_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    icd10_code_id BIGINT UNSIGNED NOT NULL,
    diagnosed_by BIGINT UNSIGNED NOT NULL,
    
    -- Diagnosis Details
    diagnosis_type VARCHAR(20) NOT NULL,
    diagnosis_status VARCHAR(20) NOT NULL,
    clinical_notes TEXT,
    diagnosed_at TIMESTAMP NOT NULL,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_diagnoses_visit_id (visit_id),
    INDEX idx_diagnoses_patient_id (patient_id),
    INDEX idx_diagnoses_icd10_code_id (icd10_code_id),
    INDEX idx_diagnoses_diagnosed_by (diagnosed_by),
    INDEX idx_diagnoses_diagnosis_type (diagnosis_type),
    INDEX idx_diagnoses_diagnosis_status (diagnosis_status),
    INDEX idx_diagnoses_diagnosed_at (diagnosed_at),
    INDEX idx_diagnoses_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (icd10_code_id) REFERENCES icd10_codes(id),
    FOREIGN KEY (diagnosed_by) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
