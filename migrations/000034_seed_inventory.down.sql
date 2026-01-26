-- Remove seeded inventory
DELETE FROM inventory WHERE batch_number IN (
    'PAR-2024-001', 'PAR-2024-002',
    'AMX-2024-001', 'AMX-2024-002',
    'IBU-2024-001',
    'OME-2024-001',
    'MET-2024-001',
    'AML-2024-001',
    'CET-2024-001',
    'SAL-2024-001',
    'INS-2024-001',
    'VIT-2024-001'
);
