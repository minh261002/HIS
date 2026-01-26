-- Remove seeded lab test templates and parameters
DELETE FROM lab_test_templates WHERE code IN ('CBC', 'FBS', 'LIPID', 'LFT', 'RFT', 'UA');
