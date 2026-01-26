-- Create insurance_claims table
CREATE TABLE IF NOT EXISTS insurance_claims (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    claim_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    invoice_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    
    -- Claim Details
    insurance_provider VARCHAR(200) NOT NULL,
    policy_number VARCHAR(100) NOT NULL,
    claim_amount DECIMAL(10,2) NOT NULL,
    approved_amount DECIMAL(10,2),
    claim_date TIMESTAMP NOT NULL,
    approval_date TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'SUBMITTED',
    rejection_reason TEXT,
    notes TEXT,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_insurance_claims_claim_code (claim_code),
    INDEX idx_insurance_claims_invoice_id (invoice_id),
    INDEX idx_insurance_claims_patient_id (patient_id),
    INDEX idx_insurance_claims_status (status),
    INDEX idx_insurance_claims_claim_date (claim_date),
    INDEX idx_insurance_claims_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (invoice_id) REFERENCES invoices(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
