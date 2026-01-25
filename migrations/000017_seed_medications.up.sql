-- Seed common medications
INSERT INTO medications (name, generic_name, dosage_form, strength, unit, manufacturer, is_active, created_at, updated_at) VALUES
-- Pain relievers
('Paracetamol', 'Acetaminophen', 'TABLET', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Ibuprofen', 'Ibuprofen', 'TABLET', '400mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Aspirin', 'Acetylsalicylic acid', 'TABLET', '100mg', 'mg', 'Various', TRUE, NOW(), NOW()),

-- Antibiotics
('Amoxicillin', 'Amoxicillin', 'CAPSULE', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Azithromycin', 'Azithromycin', 'TABLET', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Ciprofloxacin', 'Ciprofloxacin', 'TABLET', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Cephalexin', 'Cephalexin', 'CAPSULE', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),

-- Diabetes medications
('Metformin', 'Metformin', 'TABLET', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Glibenclamide', 'Glibenclamide', 'TABLET', '5mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Insulin Glargine', 'Insulin Glargine', 'INJECTION', '100IU/ml', 'IU', 'Various', TRUE, NOW(), NOW()),

-- Hypertension medications
('Amlodipine', 'Amlodipine', 'TABLET', '5mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Losartan', 'Losartan', 'TABLET', '50mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Enalapril', 'Enalapril', 'TABLET', '5mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Hydrochlorothiazide', 'Hydrochlorothiazide', 'TABLET', '25mg', 'mg', 'Various', TRUE, NOW(), NOW()),

-- Gastrointestinal
('Omeprazole', 'Omeprazole', 'CAPSULE', '20mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Ranitidine', 'Ranitidine', 'TABLET', '150mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Loperamide', 'Loperamide', 'TABLET', '2mg', 'mg', 'Various', TRUE, NOW(), NOW()),

-- Respiratory
('Salbutamol', 'Salbutamol', 'INHALER', '100mcg', 'mcg', 'Various', TRUE, NOW(), NOW()),
('Cetirizine', 'Cetirizine', 'TABLET', '10mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Loratadine', 'Loratadine', 'TABLET', '10mg', 'mg', 'Various', TRUE, NOW(), NOW()),

-- Vitamins
('Vitamin C', 'Ascorbic acid', 'TABLET', '500mg', 'mg', 'Various', TRUE, NOW(), NOW()),
('Vitamin D3', 'Cholecalciferol', 'TABLET', '1000IU', 'IU', 'Various', TRUE, NOW(), NOW()),
('Multivitamin', 'Multivitamin', 'TABLET', 'Various', 'Various', 'Various', TRUE, NOW(), NOW());
