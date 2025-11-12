-- Insert into storages
INSERT INTO storages (id, manufacturer, model, capacity, created_at, updated_at) VALUES
(101, 'BigBattery, Inc.', 'FETHS-48051-G1-08-18K-01 [240V]', 4096, NOW(), NOW()),
(102, 'BigBattery, Inc.', 'FETHS-48051-G1-08-18K-02 [208V]', 4096, NOW(), NOW()),
(103, 'BigBattery, Inc.', 'FETHS-48051-G1-08-18K-02 [240V]', 4096, NOW(), NOW()),
(104, 'BigBattery, Inc.', 'FETHS-48051-G1-08-18K-03 [208V]', 4096, NOW(), NOW()),
(105, 'BigBattery, Inc.', 'FETHS-48051-G1-08-18K-03 [240V]', 4096, NOW(), NOW()),
(106, 'BigBattery, Inc.', 'FETHS-48051-G1-08-INV020-01 [208V]', 4096, NOW(), NOW()),
(107, 'BigBattery, Inc.', 'FETHS-48051-G1-08-INV020-01 [240V]', 4096, NOW(), NOW()),
(108, 'BigBattery, Inc.', 'FETHS-48051-G1-08-INV020-02 [208V]', 4096, NOW(), NOW()),
(109, 'BigBattery, Inc.', 'FETHS-48051-G1-08-INV020-02 [240V]', 4096, NOW(), NOW()),
(110, 'BigBattery, Inc.', 'FETHS-48051-G1-08-INV020-03 [208V]', 4096, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert into inverters
INSERT INTO inverters (id, manufacturer, model, capacity, created_at, updated_at) VALUES
(5557, 'ABB', 'PVI-3.0-OUTD-S-US-A [208V]', 3, NOW(), NOW()),
(5558, 'ABB', 'PVI-3.0-OUTD-S-US-A [240V]', 3, NOW(), NOW()),
(5559, 'ABB', 'PVI-3.0-OUTD-S-US-A [277V]', 3, NOW(), NOW()),
(5560, 'ABB', 'PVI-3.0-OUTD-S-US-Z-A [208V]', 3, NOW(), NOW()),
(6035, 'Enphase Energy Inc.', 'IQ7A-72-E-ACM-US-NM [240V]', 0.349, NOW(), NOW()),
(6036, 'Enphase Energy Inc.', 'IQ7A-72-E-US [240V]', 0.349, NOW(), NOW()),
(6037, 'Enphase Energy Inc.', 'IQ7A-72-M-US [240V]', 0.349, NOW(), NOW()),
(6091, 'Fortress Power LLC', 'FP-ENVY-8K [208V]', 7.85995, NOW(), NOW()),
(6092, 'Fortress Power LLC', 'FP-ENVY-8K [240V]', 7.86266, NOW(), NOW()),
(6105, 'Fronius International GmbH', 'Fronius Primo 3.8-1 208-240 [240V]', 3.8, NOW(), NOW()),
(6106, 'Fronius International GmbH', 'Primo GEN24 3.8 208-240 [240V]', 3.802, NOW(), NOW()),
(6107, 'Fronius International GmbH', 'Primo GEN24 3.8 208-240 Plus [240V]', 3.802, NOW(), NOW()),
(6128, 'Fronius International GmbH', 'Fronius Primo 8.2-1 208-240 [208V]', 7.9, NOW(), NOW()),
(6129, 'Fronius International GmbH', 'Fronius Primo 8.2-1 208-240 [240V]', 8.2, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert into panels
INSERT INTO panels (id, manufacturer, model, wattage, longside, shortside, created_at, updated_at) VALUES
(289, 'Anji Technology', 'AJP-M660-230', 230, 1.632, 0.995, NOW(), NOW()),
(290, 'Anji Technology', 'AJP-M660-235', 235, 1.632, 0.995, NOW(), NOW()),
(291, 'Anji Technology', 'AJP-M660-240', 240, 1.632, 0.995, NOW(), NOW()),
(292, 'Anji Technology', 'AJP-M660-245', 245, 1.632, 0.995, NOW(), NOW()),
(293, 'Anji Technology', 'AJP-M660-250', 250, 1.632, 0.995, NOW(), NOW()),
(294, 'Apollo Renewables Inc', 'Apollo-M660BH-350BB', 350, NULL, NULL, NOW(), NOW()),
(295, 'Apollo Renewables Inc', 'Apollo-M660BH-355BB', 355, NULL, NULL, NOW(), NOW()),
(296, 'Apollo Renewables Inc', 'Apollo-M660BH-360BB', 360, NULL, NULL, NOW(), NOW()),
(297, 'Apollo Renewables Inc', 'Apollo-M660BH-365BB', 365, NULL, NULL, NOW(), NOW()),
(298, 'Apollo Renewables Inc', 'Apollo-M660BH-370BB', 370, NULL, NULL, NOW(), NOW()),
(299, 'Apollo Renewables Inc', 'Apollo-M754BH-385BB', 385, NULL, NULL, NOW(), NOW()),
(300, 'Apollo Renewables Inc', 'Apollo-M754BH-390BB', 390, NULL, NULL, NOW(), NOW()),
(301, 'Apollo Renewables Inc', 'Apollo-M754BH-395BB', 395, NULL, NULL, NOW(), NOW()),
(320, 'APOS Energy', 'AP180', 180, 1.341, 1, NOW(), NOW()),
(321, 'APOS Energy', 'AS180', 180, 1.341, 1, NOW(), NOW()),
(322, 'APOS Energy', 'AP185', 185, 1.341, 1, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
