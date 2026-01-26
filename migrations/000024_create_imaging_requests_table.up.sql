-- Create imaging_requests table
CREATE TABLE IF NOT EXISTS imaging_requests (
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
    scheduled_date TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    clinical_indication TEXT,
    special_instructions TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_imaging_requests_request_code (request_code),
    INDEX idx_imaging_requests_visit_id (visit_id),
    INDEX idx_imaging_requests_patient_id (patient_id),
    INDEX idx_imaging_requests_doctor_id (doctor_id),
    INDEX idx_imaging_requests_template_id (template_id),
    INDEX idx_imaging_requests_status (status),
    INDEX idx_imaging_requests_requested_date (requested_date),
    INDEX idx_imaging_requests_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id) ON DELETE CASCADE,
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    FOREIGN KEY (template_id) REFERENCES imaging_templates(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
