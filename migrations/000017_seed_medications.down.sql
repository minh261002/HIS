-- Remove seeded medications
DELETE FROM medications WHERE name IN (
    'Paracetamol', 'Ibuprofen', 'Aspirin',
    'Amoxicillin', 'Azithromycin', 'Ciprofloxacin', 'Cephalexin',
    'Metformin', 'Glibenclamide', 'Insulin Glargine',
    'Amlodipine', 'Losartan', 'Enalapril', 'Hydrochlorothiazide',
    'Omeprazole', 'Ranitidine', 'Loperamide',
    'Salbutamol', 'Cetirizine', 'Loratadine',
    'Vitamin C', 'Vitamin D3', 'Multivitamin'
);
