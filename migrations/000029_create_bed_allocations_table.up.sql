-- Create bed_allocations table
CREATE TABLE IF NOT EXISTS bed_allocations (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Keys
    admission_id BIGINT UNSIGNED NOT NULL,
    bed_id BIGINT UNSIGNED NOT NULL,
    
    -- Allocation Details
    allocated_date TIMESTAMP NOT NULL,
    released_date TIMESTAMP NULL,
    is_current BOOLEAN DEFAULT TRUE,
    notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_bed_allocations_admission_id (admission_id),
    INDEX idx_bed_allocations_bed_id (bed_id),
    INDEX idx_bed_allocations_is_current (is_current),
    INDEX idx_bed_allocations_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (admission_id) REFERENCES admissions(id) ON DELETE CASCADE,
    FOREIGN KEY (bed_id) REFERENCES beds(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
