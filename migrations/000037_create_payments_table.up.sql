-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    payment_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    invoice_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    
    -- Payment Details
    payment_method VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    transaction_id VARCHAR(100),
    gateway_response JSON,
    notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_payments_payment_code (payment_code),
    INDEX idx_payments_invoice_id (invoice_id),
    INDEX idx_payments_patient_id (patient_id),
    INDEX idx_payments_payment_method (payment_method),
    INDEX idx_payments_status (status),
    INDEX idx_payments_payment_date (payment_date),
    INDEX idx_payments_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (invoice_id) REFERENCES invoices(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
