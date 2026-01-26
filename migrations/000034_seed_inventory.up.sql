-- Seed sample inventory for existing medications

-- Paracetamol 500mg (medication_id: 1)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(1, 'PAR-2024-001', '2026-12-31', 1000, 'TABLET', 0.50, 'ABC Pharma', '2024-01-15', NOW(), NOW()),
(1, 'PAR-2024-002', '2027-06-30', 500, 'TABLET', 0.55, 'ABC Pharma', '2024-06-01', NOW(), NOW());

-- Amoxicillin 500mg (medication_id: 2)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(2, 'AMX-2024-001', '2026-06-30', 500, 'CAPSULE', 2.00, 'XYZ Pharma', '2024-02-01', NOW(), NOW()),
(2, 'AMX-2024-002', '2026-12-31', 300, 'CAPSULE', 2.10, 'XYZ Pharma', '2024-07-15', NOW(), NOW());

-- Ibuprofen 400mg (medication_id: 3)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(3, 'IBU-2024-001', '2026-09-30', 800, 'TABLET', 1.20, 'ABC Pharma', '2024-03-10', NOW(), NOW());

-- Omeprazole 20mg (medication_id: 4)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(4, 'OME-2024-001', '2026-11-30', 400, 'CAPSULE', 3.50, 'MediCare', '2024-04-05', NOW(), NOW());

-- Metformin 500mg (medication_id: 5)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(5, 'MET-2024-001', '2027-03-31', 600, 'TABLET', 1.80, 'DiabCare', '2024-05-20', NOW(), NOW());

-- Amlodipine 5mg (medication_id: 6)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(6, 'AML-2024-001', '2026-08-31', 350, 'TABLET', 2.50, 'CardioMed', '2024-06-10', NOW(), NOW());

-- Cetirizine 10mg (medication_id: 7)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(7, 'CET-2024-001', '2026-10-31', 450, 'TABLET', 1.50, 'AllergyFree', '2024-07-01', NOW(), NOW());

-- Salbutamol Inhaler (medication_id: 8)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(8, 'SAL-2024-001', '2026-05-31', 100, 'BOTTLE', 45.00, 'RespiroCare', '2024-08-15', NOW(), NOW());

-- Insulin Glargine (medication_id: 9)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(9, 'INS-2024-001', '2026-04-30', 80, 'VIAL', 250.00, 'DiabCare', '2024-09-01', NOW(), NOW());

-- Vitamin B Complex (medication_id: 10)
INSERT INTO inventory (medication_id, batch_number, expiry_date, quantity, unit, cost_price, supplier, received_date, created_at, updated_at) VALUES
(10, 'VIT-2024-001', '2027-02-28', 700, 'TABLET', 0.80, 'VitaHealth', '2024-10-10', NOW(), NOW());
