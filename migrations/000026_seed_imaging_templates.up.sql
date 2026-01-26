-- Seed common imaging templates

-- X-Ray Templates
INSERT INTO imaging_templates (code, name, modality, body_part, description, template_content, price, is_active, created_at, updated_at) VALUES
('XRAY_CHEST', 'Chest X-Ray', 'XRAY', 'CHEST', 'Chụp X-quang ngực', 'FINDINGS:\n- Heart size: Normal\n- Lung fields: Clear\n- Costophrenic angles: Sharp\n- Bony thorax: Intact\n\nIMPRESSION:\nNo acute cardiopulmonary abnormality.', 150000, TRUE, NOW(), NOW()),
('XRAY_ABD', 'Abdomen X-Ray', 'XRAY', 'ABDOMEN', 'Chụp X-quang bụng', 'FINDINGS:\n- Bowel gas pattern: Normal\n- No free air\n- No abnormal calcifications\n\nIMPRESSION:\nNormal abdominal radiograph.', 150000, TRUE, NOW(), NOW()),
('XRAY_SPINE', 'Spine X-Ray', 'XRAY', 'SPINE', 'Chụp X-quang cột sống', 'FINDINGS:\n- Vertebral alignment: Normal\n- Disc spaces: Preserved\n- No fracture or dislocation\n\nIMPRESSION:\nNo acute abnormality.', 180000, TRUE, NOW(), NOW());

-- CT Scan Templates
INSERT INTO imaging_templates (code, name, modality, body_part, description, template_content, price, is_active, created_at, updated_at) VALUES
('CT_HEAD', 'CT Scan Head', 'CT', 'HEAD', 'Chụp CT sọ não', 'FINDINGS:\n- Brain parenchyma: Normal attenuation\n- Ventricles: Normal size and configuration\n- No intracranial hemorrhage\n- No mass effect or midline shift\n\nIMPRESSION:\nNo acute intracranial abnormality.', 800000, TRUE, NOW(), NOW()),
('CT_CHEST', 'CT Scan Chest', 'CT', 'CHEST', 'Chụp CT lồng ngực', 'FINDINGS:\n- Lungs: Clear, no nodules or masses\n- Mediastinum: Normal\n- Heart: Normal size\n- No pleural effusion\n\nIMPRESSION:\nNo acute abnormality.', 900000, TRUE, NOW(), NOW()),
('CT_ABD', 'CT Scan Abdomen', 'CT', 'ABDOMEN', 'Chụp CT bụng', 'FINDINGS:\n- Liver: Normal size and attenuation\n- Spleen, pancreas, kidneys: Normal\n- No free fluid\n- Bowel: Normal\n\nIMPRESSION:\nNo acute abdominal abnormality.', 1000000, TRUE, NOW(), NOW());

-- Ultrasound Templates
INSERT INTO imaging_templates (code, name, modality, body_part, description, template_content, price, is_active, created_at, updated_at) VALUES
('US_ABD', 'Abdominal Ultrasound', 'ULTRASOUND', 'ABDOMEN', 'Siêu âm bụng', 'FINDINGS:\n- Liver: Normal size and echogenicity\n- Gallbladder: No stones\n- Spleen: Normal\n- Kidneys: Normal size, no hydronephrosis\n- Pancreas: Visualized portions normal\n\nIMPRESSION:\nNormal abdominal ultrasound.', 250000, TRUE, NOW(), NOW()),
('US_OB', 'Obstetric Ultrasound', 'ULTRASOUND', 'PELVIS', 'Siêu âm sản khoa', 'FINDINGS:\n- Single intrauterine pregnancy\n- Fetal cardiac activity: Present\n- Gestational age: [GA] weeks\n- Amniotic fluid: Normal\n- Placenta: [Location]\n\nIMPRESSION:\nViable intrauterine pregnancy.', 300000, TRUE, NOW(), NOW()),
('US_THYROID', 'Thyroid Ultrasound', 'ULTRASOUND', 'HEAD', 'Siêu âm tuyến giáp', 'FINDINGS:\n- Right lobe: Normal size and echogenicity\n- Left lobe: Normal size and echogenicity\n- Isthmus: Normal\n- No focal lesions\n\nIMPRESSION:\nNormal thyroid ultrasound.', 200000, TRUE, NOW(), NOW());

-- MRI Templates
INSERT INTO imaging_templates (code, name, modality, body_part, description, template_content, price, is_active, created_at, updated_at) VALUES
('MRI_BRAIN', 'Brain MRI', 'MRI', 'HEAD', 'Chụp MRI não', 'FINDINGS:\n- Brain parenchyma: Normal signal intensity\n- Ventricles: Normal size\n- No abnormal enhancement\n- No mass or hemorrhage\n\nIMPRESSION:\nNo acute intracranial abnormality.', 2500000, TRUE, NOW(), NOW()),
('MRI_SPINE', 'Spine MRI', 'MRI', 'SPINE', 'Chụp MRI cột sống', 'FINDINGS:\n- Vertebral alignment: Normal\n- Spinal cord: Normal signal\n- Disc spaces: Preserved\n- No significant stenosis\n\nIMPRESSION:\nNo significant abnormality.', 2800000, TRUE, NOW(), NOW());

-- Mammography
INSERT INTO imaging_templates (code, name, modality, body_part, description, template_content, price, is_active, created_at, updated_at) VALUES
('MAMMO', 'Mammography', 'MAMMOGRAPHY', 'CHEST', 'Chụp nhũ ảnh', 'FINDINGS:\n- Breast composition: [Type]\n- Right breast: No suspicious masses or calcifications\n- Left breast: No suspicious masses or calcifications\n\nIMPRESSION:\nBI-RADS Category [1-6]', 400000, TRUE, NOW(), NOW());
