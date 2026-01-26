-- Create nursing_notes table
CREATE TABLE IF NOT EXISTS nursing_notes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Keys
    admission_id BIGINT UNSIGNED NOT NULL,
    nurse_id BIGINT UNSIGNED NOT NULL,
    
    -- Note Details
    note_date TIMESTAMP NOT NULL,
    vital_signs JSON,
    observations TEXT NOT NULL,
    interventions TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_nursing_notes_admission_id (admission_id),
    INDEX idx_nursing_notes_nurse_id (nurse_id),
    INDEX idx_nursing_notes_note_date (note_date),
    INDEX idx_nursing_notes_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (admission_id) REFERENCES admissions(id) ON DELETE CASCADE,
    FOREIGN KEY (nurse_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
