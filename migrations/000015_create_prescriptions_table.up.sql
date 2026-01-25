-- Create prescriptions table
CREATE TABLE IF NOT EXISTS prescriptions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    prescription_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    visit_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    doctor_id BIGINT UNSIGNED NOT NULL,
    diagnosis_id BIGINT UNSIGNED,
    
    -- Prescription Details
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    prescribed_date DATE NOT NULL,
    notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_prescriptions_prescription_code (prescription_code),
    INDEX idx_prescriptions_visit_id (visit_id),
    INDEX idx_prescriptions_patient_id (patient_id),
    INDEX idx_prescriptions_doctor_id (doctor_id),
    INDEX idx_prescriptions_diagnosis_id (diagnosis_id),
    INDEX idx_prescriptions_status (status),
    INDEX idx_prescriptions_prescribed_date (prescribed_date),
    INDEX idx_prescriptions_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    FOREIGN KEY (diagnosis_id) REFERENCES diagnoses(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
