-- Create lab_test_requests table
CREATE TABLE IF NOT EXISTS lab_test_requests (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    request_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    visit_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    doctor_id BIGINT UNSIGNED NOT NULL,
    template_id BIGINT UNSIGNED NOT NULL,
    
    -- Request Details
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    priority VARCHAR(20) NOT NULL DEFAULT 'ROUTINE',
    requested_date TIMESTAMP NOT NULL,
    sample_collected_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    clinical_notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_lab_test_requests_request_code (request_code),
    INDEX idx_lab_test_requests_visit_id (visit_id),
    INDEX idx_lab_test_requests_patient_id (patient_id),
    INDEX idx_lab_test_requests_doctor_id (doctor_id),
    INDEX idx_lab_test_requests_template_id (template_id),
    INDEX idx_lab_test_requests_status (status),
    INDEX idx_lab_test_requests_requested_date (requested_date),
    INDEX idx_lab_test_requests_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    FOREIGN KEY (template_id) REFERENCES lab_test_templates(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
