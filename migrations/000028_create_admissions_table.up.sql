-- Create admissions table
CREATE TABLE IF NOT EXISTS admissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    admission_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    visit_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    doctor_id BIGINT UNSIGNED NOT NULL,
    
    -- Admission Details
    admission_date TIMESTAMP NOT NULL,
    discharge_date TIMESTAMP NULL,
    admission_diagnosis TEXT NOT NULL,
    discharge_diagnosis TEXT,
    discharge_summary TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'ADMITTED',
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_admissions_admission_code (admission_code),
    INDEX idx_admissions_visit_id (visit_id),
    INDEX idx_admissions_patient_id (patient_id),
    INDEX idx_admissions_doctor_id (doctor_id),
    INDEX idx_admissions_status (status),
    INDEX idx_admissions_admission_date (admission_date),
    INDEX idx_admissions_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
