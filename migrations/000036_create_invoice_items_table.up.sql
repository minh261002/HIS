-- Create invoice_items table
CREATE TABLE IF NOT EXISTS invoice_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    
    -- Foreign Key
    invoice_id BIGINT UNSIGNED NOT NULL,
    
    -- Item Details
    item_type VARCHAR(20) NOT NULL,
    item_id BIGINT UNSIGNED,
    description TEXT NOT NULL,
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_invoice_items_invoice_id (invoice_id),
    INDEX idx_invoice_items_item_id (item_id),
    INDEX idx_invoice_items_deleted_at (deleted_at),
    
    -- Foreign Key
    FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
