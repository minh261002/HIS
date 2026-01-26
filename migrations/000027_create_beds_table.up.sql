-- Create beds table
CREATE TABLE IF NOT EXISTS beds (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Bed Information
    bed_number VARCHAR(20) NOT NULL UNIQUE,
    department VARCHAR(30) NOT NULL,
    ward VARCHAR(50) NOT NULL,
    bed_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    is_active BOOLEAN DEFAULT TRUE,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_beds_bed_number (bed_number),
    INDEX idx_beds_department (department),
    INDEX idx_beds_bed_type (bed_type),
    INDEX idx_beds_status (status),
    INDEX idx_beds_is_active (is_active),
    INDEX idx_beds_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
