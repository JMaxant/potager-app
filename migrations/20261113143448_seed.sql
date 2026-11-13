-- +goose Up
INSERT INTO families (name, rotation_interval_years, nitrogen_fixing)
VALUES ('Fabacées', 3, true),
       ('Solanacées', 4, false),
       ('Apiacées', 3, false),
       ('Poacées', 2, false),
       ('Brassicacées', 4, false),
       ('Amaryllidacées', 4, false),
       ('Astéracées', 2, false),
       ('Boraginacées', 0, false);


INSERT INTO species (family_id, name, suitable_as_cover, termination_lead_days, nitrogen_demand, maturity_from)
VALUES ((SELECT id FROM families WHERE name = 'Solanacées'), 'Tomate', false, NULL, 'heavy', 'transplant'),
       ((SELECT id FROM families WHERE name = 'Apiacées'), 'Carotte', false, NULL, 'light', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Amaryllidacées'), 'Ail', false, NULL, 'light', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Astéracées'), 'Laitue', false, NULL, 'moderate', 'transplant'),
       ((SELECT id FROM families WHERE name = 'Fabacées'), 'Fève', true, 21, 'light', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Brassicacées'), 'Moutarde', true, 15, 'heavy', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Brassicacées'), 'Radis', true, 15, 'light', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Fabacées'), 'Trèfle', true, 21, 'light', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Poacées'), 'Seigle', true, 42, 'moderate', 'sowing'),
       ((SELECT id FROM families WHERE name = 'Boraginacées'), 'Phacélie', true, 15, 'moderate', 'sowing');

INSERT INTO varieties (name, species_id, days_to_maturity)
VALUES ('Roma', (SELECT id FROM species WHERE name = 'Tomate'), 75),
       ('Rose de Berne', (SELECT id FROM species WHERE name = 'Tomate'), 80),
       ('Nantaise 2', (SELECT id FROM species WHERE name = 'Carotte'), 90),
       ('Violet de Cadours', (SELECT id FROM species WHERE name = 'Ail'), 240),
       ('Batavia blonde de Paris', (SELECT id FROM species WHERE name = 'Laitue'), 50),
       ('Aguadulce', (SELECT id FROM species WHERE name = 'Fève'), 180),
       ('Féverole d''hiver', (SELECT id FROM species WHERE name = 'Fève'), NULL),
       ('Moutarde blanche', (SELECT id FROM species WHERE name = 'Moutarde'), NULL),
       ('Radis de 18 jours', (SELECT id FROM species WHERE name = 'Radis'), 21),
       ('Radis fourrager', (SELECT id FROM species WHERE name = 'Radis'), NULL),
       ('Trèfle incarnat', (SELECT id FROM species WHERE name = 'Trèfle'), NULL),
       ('Seigle d''hiver', (SELECT id FROM species WHERE name = 'Seigle'), NULL),
       ('Phacélie', (SELECT id FROM species WHERE name = 'Phacélie'), NULL);

-- +goose Down
DELETE
FROM varieties
WHERE name IN ('Roma', 'Rose de Berne', 'Nantaise 2', 'Violet de Cadours', 'Batavia blonde de Paris', 'Aguadulce',
               'Féverole d''hiver', 'Moutarde blanche', 'Radis de 18 jours', 'Trèfle incarnat', 'Seigle d''hiver',
               'Radis fourrager', 'Phacélie');
DELETE
FROM species
WHERE name IN ('Tomate', 'Carotte', 'Ail', 'Laitue', 'Fève', 'Moutarde', 'Radis', 'Trèfle', 'Seigle', 'Phacélie');
DELETE
FROM families
WHERE name IN
      ('Fabacées', 'Solanacées', 'Apiacées', 'Poacées', 'Brassicacées', 'Amaryllidacées', 'Astéracées', 'Boraginacées');
