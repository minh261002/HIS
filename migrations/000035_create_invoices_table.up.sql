-- Create invoices table
CREATE TABLE IF NOT EXISTS invoices (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    invoice_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    visit_id BIGINT UNSIGNED,
    patient_id BIGINT UNSIGNED NOT NULL,
    
    -- Invoice Details
    invoice_date TIMESTAMP NOT NULL,
    due_date TIMESTAMP NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_invoices_invoice_code (invoice_code),
    INDEX idx_invoices_visit_id (visit_id),
    INDEX idx_invoices_patient_id (patient_id),
    INDEX idx_invoices_status (status),
    INDEX idx_invoices_invoice_date (invoice_date),
    INDEX idx_invoices_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (visit_id) REFERENCES visits(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
