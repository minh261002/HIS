-- Create visits table
CREATE TABLE IF NOT EXISTS visits (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    visit_code VARCHAR(20) NOT NULL UNIQUE,
    
    -- Foreign Keys
    appointment_id BIGINT UNSIGNED,
    patient_id BIGINT UNSIGNED NOT NULL,
    doctor_id BIGINT UNSIGNED NOT NULL,
    
    -- Visit Details
    visit_date DATE NOT NULL,
    visit_time TIME NOT NULL,
    visit_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'WAITING',
    
    -- Chief Complaint & Symptoms
    chief_complaint TEXT NOT NULL,
    symptoms TEXT,
    
    -- Vital Signs
    temperature DECIMAL(4,1),
    blood_pressure_systolic INT,
    blood_pressure_diastolic INT,
    heart_rate INT,
    respiratory_rate INT,
    oxygen_saturation INT,
    weight DECIMAL(5,2),
    height DECIMAL(5,2),
    bmi DECIMAL(5,2),
    
    -- Clinical Documentation
    physical_examination TEXT,
    clinical_notes TEXT,
    treatment_plan TEXT,
    follow_up_instructions TEXT,
    next_visit_date DATE,
    
    -- Audit fields
    created_by BIGINT UNSIGNED NOT NULL,
    updated_by BIGINT UNSIGNED,
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_visits_visit_code (visit_code),
    INDEX idx_visits_appointment_id (appointment_id),
    INDEX idx_visits_patient_id (patient_id),
    INDEX idx_visits_doctor_id (doctor_id),
    INDEX idx_visits_visit_date (visit_date),
    INDEX idx_visits_status (status),
    INDEX idx_visits_deleted_at (deleted_at),
    
    -- Foreign Keys
    FOREIGN KEY (appointment_id) REFERENCES appointments(id),
    FOREIGN KEY (patient_id) REFERENCES patients(id),
    FOREIGN KEY (doctor_id) REFERENCES users(id),
    FOREIGN KEY (created_by) REFERENCES users(id),
    FOREIGN KEY (updated_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
