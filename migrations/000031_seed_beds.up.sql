-- Seed sample beds

-- Internal Medicine Ward A (20 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('A101', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A102', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A103', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A104', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A105', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A106', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A107', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A108', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A109', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A110', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A111', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A112', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A113', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A114', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A115', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A116', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A117', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A118', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A119', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('A120', 'INTERNAL_MEDICINE', 'Ward A', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW());

-- Surgery Ward B (15 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('B101', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B102', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B103', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B104', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B105', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B106', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B107', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B108', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B109', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B110', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B111', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B112', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B113', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B114', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('B115', 'SURGERY', 'Ward B', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW());

-- Pediatrics Ward P (15 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('P101', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P102', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P103', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P104', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P105', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P106', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P107', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P108', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P109', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P110', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P111', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P112', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P113', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P114', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW()),
('P115', 'PEDIATRICS', 'Ward P', 'STANDARD', 'AVAILABLE', TRUE, NOW(), NOW());

-- ICU (10 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('ICU-01', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-02', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-03', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-04', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-05', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-06', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-07', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-08', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-09', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW()),
('ICU-10', 'ICU', 'ICU', 'ICU', 'AVAILABLE', TRUE, NOW(), NOW());

-- Isolation Ward (5 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('ISO-01', 'INTERNAL_MEDICINE', 'Isolation', 'ISOLATION', 'AVAILABLE', TRUE, NOW(), NOW()),
('ISO-02', 'INTERNAL_MEDICINE', 'Isolation', 'ISOLATION', 'AVAILABLE', TRUE, NOW(), NOW()),
('ISO-03', 'INTERNAL_MEDICINE', 'Isolation', 'ISOLATION', 'AVAILABLE', TRUE, NOW(), NOW()),
('ISO-04', 'INTERNAL_MEDICINE', 'Isolation', 'ISOLATION', 'AVAILABLE', TRUE, NOW(), NOW()),
('ISO-05', 'INTERNAL_MEDICINE', 'Isolation', 'ISOLATION', 'AVAILABLE', TRUE, NOW(), NOW());

-- VIP Ward (5 beds)
INSERT INTO beds (bed_number, department, ward, bed_type, status, is_active, created_at, updated_at) VALUES
('VIP-01', 'INTERNAL_MEDICINE', 'VIP', 'VIP', 'AVAILABLE', TRUE, NOW(), NOW()),
('VIP-02', 'INTERNAL_MEDICINE', 'VIP', 'VIP', 'AVAILABLE', TRUE, NOW(), NOW()),
('VIP-03', 'SURGERY', 'VIP', 'VIP', 'AVAILABLE', TRUE, NOW(), NOW()),
('VIP-04', 'SURGERY', 'VIP', 'VIP', 'AVAILABLE', TRUE, NOW(), NOW()),
('VIP-05', 'OBSTETRICS', 'VIP', 'VIP', 'AVAILABLE', TRUE, NOW(), NOW());
