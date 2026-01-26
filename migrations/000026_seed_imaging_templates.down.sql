-- Remove seeded imaging templates
DELETE FROM imaging_templates WHERE code IN (
    'XRAY_CHEST', 'XRAY_ABD', 'XRAY_SPINE',
    'CT_HEAD', 'CT_CHEST', 'CT_ABD',
    'US_ABD', 'US_OB', 'US_THYROID',
    'MRI_BRAIN', 'MRI_SPINE',
    'MAMMO'
);
