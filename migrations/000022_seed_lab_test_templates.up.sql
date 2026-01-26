-- Seed common lab test templates with parameters

-- Complete Blood Count (CBC)
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('CBC', 'Complete Blood Count', 'HEMATOLOGY', 'Đếm tế bào máu toàn phần', 'BLOOD', 'Không cần nhịn ăn', 2, 150000, TRUE, NOW(), NOW());

SET @cbc_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@cbc_id, 'WBC (Bạch cầu)', '10^9/L', 4.0, 11.0, '4.0-11.0', 1, NOW(), NOW()),
(@cbc_id, 'RBC (Hồng cầu)', '10^12/L', 4.0, 5.5, '4.0-5.5', 2, NOW(), NOW()),
(@cbc_id, 'Hemoglobin', 'g/dL', 12.0, 17.5, '12.0-17.5', 3, NOW(), NOW()),
(@cbc_id, 'Hematocrit', '%', 36, 50, '36-50', 4, NOW(), NOW()),
(@cbc_id, 'Platelets (Tiểu cầu)', '10^9/L', 150, 400, '150-400', 5, NOW(), NOW());

-- Fasting Blood Sugar (FBS)
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('FBS', 'Fasting Blood Sugar', 'BIOCHEMISTRY', 'Đường huyết lúc đói', 'BLOOD', 'Nhịn ăn 8-12 giờ', 1, 50000, TRUE, NOW(), NOW());

SET @fbs_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@fbs_id, 'Glucose', 'mg/dL', 70, 100, '70-100', 1, NOW(), NOW());

-- Lipid Panel
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('LIPID', 'Lipid Panel', 'BIOCHEMISTRY', 'Bộ xét nghiệm mỡ máu', 'BLOOD', 'Nhịn ăn 12 giờ', 2, 200000, TRUE, NOW(), NOW());

SET @lipid_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@lipid_id, 'Total Cholesterol', 'mg/dL', 0, 200, '<200', 1, NOW(), NOW()),
(@lipid_id, 'LDL Cholesterol', 'mg/dL', 0, 100, '<100', 2, NOW(), NOW()),
(@lipid_id, 'HDL Cholesterol', 'mg/dL', 40, 999, '>40', 3, NOW(), NOW()),
(@lipid_id, 'Triglycerides', 'mg/dL', 0, 150, '<150', 4, NOW(), NOW());

-- Liver Function Test (LFT)
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('LFT', 'Liver Function Test', 'BIOCHEMISTRY', 'Chức năng gan', 'BLOOD', 'Nhịn ăn 8 giờ', 2, 250000, TRUE, NOW(), NOW());

SET @lft_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@lft_id, 'ALT (SGPT)', 'U/L', 7, 56, '7-56', 1, NOW(), NOW()),
(@lft_id, 'AST (SGOT)', 'U/L', 10, 40, '10-40', 2, NOW(), NOW()),
(@lft_id, 'Alkaline Phosphatase', 'U/L', 44, 147, '44-147', 3, NOW(), NOW()),
(@lft_id, 'Total Bilirubin', 'mg/dL', 0.1, 1.2, '0.1-1.2', 4, NOW(), NOW());

-- Kidney Function Test (RFT)
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('RFT', 'Renal Function Test', 'BIOCHEMISTRY', 'Chức năng thận', 'BLOOD', 'Không cần nhịn ăn', 2, 180000, TRUE, NOW(), NOW());

SET @rft_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@rft_id, 'Creatinine', 'mg/dL', 0.6, 1.2, '0.6-1.2', 1, NOW(), NOW()),
(@rft_id, 'Urea', 'mg/dL', 15, 40, '15-40', 2, NOW(), NOW()),
(@rft_id, 'Uric Acid', 'mg/dL', 3.5, 7.2, '3.5-7.2', 3, NOW(), NOW());

-- Urinalysis
INSERT INTO lab_test_templates (code, name, category, description, sample_type, preparation_instructions, turnaround_time_hours, price, is_active, created_at, updated_at) 
VALUES ('UA', 'Urinalysis', 'UROLOGY', 'Xét nghiệm nước tiểu', 'URINE', 'Lấy mẫu nước tiểu giữa dòng', 1, 80000, TRUE, NOW(), NOW());

SET @ua_id = LAST_INSERT_ID();

INSERT INTO lab_test_template_parameters (template_id, parameter_name, unit, normal_range_min, normal_range_max, normal_range_text, display_order, created_at, updated_at) VALUES
(@ua_id, 'pH', '', 4.5, 8.0, '4.5-8.0', 1, NOW(), NOW()),
(@ua_id, 'Protein', 'mg/dL', 0, 10, 'Negative', 2, NOW(), NOW()),
(@ua_id, 'Glucose', 'mg/dL', 0, 15, 'Negative', 3, NOW(), NOW()),
(@ua_id, 'WBC', 'cells/HPF', 0, 5, '0-5', 4, NOW(), NOW()),
(@ua_id, 'RBC', 'cells/HPF', 0, 3, '0-3', 5, NOW(), NOW());
