-- Create dispensing table
CREATE TABLE IF NOT EXISTS dispensing (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    dispensing_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    prescription_id BIGINT UNSIGNED NOT NULL,
    prescription_item_id BIGINT UNSIGNED NOT NULL,
    medication_id BIGINT UNSIGNED NOT NULL,
    inventory_id BIGINT UNSIGNED NOT NULL,
    patient_id BIGINT UNSIGNED NOT NULL,
    pharmacist_id BIGINT UNSIGNED NOT NULL,
    
    -- Dispensing Details
    quantity_dispensed INT NOT NULL,
    batch_number VARCHAR(50) NOT NULL,
    dispensed_date TIMESTAMP NOT NULL,
    notes TEXT,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_dispensing_dispensing_code (dispensing_code),
    INDEX idx_dispensing_prescription_id (prescription_id),
    INDEX idx_dispensing_prescription_item_id (prescription_item_id),
    INDEX idx_dispensing_medication_id (medication_id),
    INDEX idx_dispensing_inventory_id (inventory_id),
    INDEX idx_dispensing_patient_id (patient_id),
    INDEX idx_dispensing_pharmacist_id (pharmacist_id),
    INDEX idx_dispensing_dispensed_date (dispensed_date),
    INDEX idx_dispensing_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (prescription_id) REFERENCES prescriptions(id),
    FOREIGN KEY (prescription_item_id) REFERENCES prescription_items(id),
    FOREIGN KEY (medication_id) REFERENCES medications(id),
    FOREIGN KEY (inventory_id) REFERENCES inventory(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (pharmacist_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
