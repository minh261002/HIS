-- Seed common ICD-10 codes for testing
INSERT INTO icd10_codes (code, description, category, is_active, created_at, updated_at) VALUES
-- Infectious diseases (A00-B99)
('A00.0', 'Cholera do vi khuẩn Vibrio cholerae 01, biovar cholerae', 'Infectious', TRUE, NOW(), NOW()),
('A09', 'Tiêu chảy và viêm dạ dày ruột', 'Infectious', TRUE, NOW(), NOW()),
('A15.0', 'Lao phổi', 'Infectious', TRUE, NOW(), NOW()),

-- Endocrine, nutritional and metabolic diseases (E00-E89)
('E11.9', 'Đái tháo đường type 2 không có biến chứng', 'Endocrine', TRUE, NOW(), NOW()),
('E11.2', 'Đái tháo đường type 2 với biến chứng thận', 'Endocrine', TRUE, NOW(), NOW()),
('E78.5', 'Tăng lipid máu không xác định', 'Endocrine', TRUE, NOW(), NOW()),
('E66.9', 'Béo phì không xác định', 'Endocrine', TRUE, NOW(), NOW()),

-- Circulatory system (I00-I99)
('I10', 'Tăng huyết áp nguyên phát (cao huyết áp)', 'Circulatory', TRUE, NOW(), NOW()),
('I25.1', 'Bệnh tim do xơ vữa động mạch', 'Circulatory', TRUE, NOW(), NOW()),
('I20.0', 'Đau thắt ngực không ổn định', 'Circulatory', TRUE, NOW(), NOW()),
('I50.9', 'Suy tim không xác định', 'Circulatory', TRUE, NOW(), NOW()),

-- Respiratory system (J00-J99)
('J00', 'Viêm mũi họng cấp (cảm lạnh thông thường)', 'Respiratory', TRUE, NOW(), NOW()),
('J06.9', 'Nhiễm trùng đường hô hấp trên cấp tính không xác định', 'Respiratory', TRUE, NOW(), NOW()),
('J18.9', 'Viêm phổi không xác định', 'Respiratory', TRUE, NOW(), NOW()),
('J45.9', 'Hen phế quản không xác định', 'Respiratory', TRUE, NOW(), NOW()),
('J44.9', 'Bệnh phổi tắc nghẽn mạn tính không xác định', 'Respiratory', TRUE, NOW(), NOW()),

-- Digestive system (K00-K95)
('K29.7', 'Viêm dạ dày không xác định', 'Digestive', TRUE, NOW(), NOW()),
('K21.9', 'Bệnh trào ngược dạ dày thực quản không viêm thực quản', 'Digestive', TRUE, NOW(), NOW()),
('K30', 'Khó tiêu', 'Digestive', TRUE, NOW(), NOW()),
('K59.0', 'Táo bón', 'Digestive', TRUE, NOW(), NOW()),

-- Musculoskeletal system (M00-M99)
('M54.5', 'Đau lưng dưới', 'Musculoskeletal', TRUE, NOW(), NOW()),
('M25.5', 'Đau khớp', 'Musculoskeletal', TRUE, NOW(), NOW()),
('M79.3', 'Viêm mô mỡ dưới da không xác định', 'Musculoskeletal', TRUE, NOW(), NOW()),

-- Nervous system (G00-G99)
('G43.9', 'Đau nửa đầu không xác định', 'Nervous', TRUE, NOW(), NOW()),
('G44.2', 'Đau đầu căng thẳng', 'Nervous', TRUE, NOW(), NOW()),

-- Symptoms and signs (R00-R99)
('R50.9', 'Sốt không xác định', 'Symptoms', TRUE, NOW(), NOW()),
('R51', 'Đau đầu', 'Symptoms', TRUE, NOW(), NOW()),
('R10.4', 'Đau bụng không xác định', 'Symptoms', TRUE, NOW(), NOW()),
('R05', 'Ho', 'Symptoms', TRUE, NOW(), NOW());
