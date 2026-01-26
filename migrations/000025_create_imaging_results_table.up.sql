-- Create imaging_results table
CREATE TABLE IF NOT EXISTS imaging_results (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key (one-to-one with request)
    request_id BIGINT UNSIGNED NOT NULL UNIQUE,
    radiologist_id BIGINT UNSIGNED NOT NULL,
    
    -- Result Information
    findings TEXT NOT NULL,
    impression TEXT NOT NULL,
    dicom_files JSON,
    report_date TIMESTAMP NOT NULL,
    is_critical BOOLEAN DEFAULT FALSE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_imaging_results_request_id (request_id),
    INDEX idx_imaging_results_radiologist_id (radiologist_id),
    INDEX idx_imaging_results_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (request_id) REFERENCES imaging_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (radiologist_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
