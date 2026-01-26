-- Create inventory table
CREATE TABLE IF NOT EXISTS inventory (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key
    medication_id BIGINT UNSIGNED NOT NULL,
    
    -- Inventory Details
    batch_number VARCHAR(50) NOT NULL,
    expiry_date DATE NOT NULL,
    quantity INT NOT NULL,
    unit VARCHAR(20) NOT NULL,
    cost_price DECIMAL(10,2),
    supplier VARCHAR(200),
    received_date DATE NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_inventory_medication_id (medication_id),
    INDEX idx_inventory_batch_number (batch_number),
    INDEX idx_inventory_expiry_date (expiry_date),
    INDEX idx_inventory_deleted_at (deleted_at),
    
    -- Foreign Key
    FOREIGN KEY (medication_id) REFERENCES medications(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
