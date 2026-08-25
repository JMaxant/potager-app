-- +goose Up
-- +goose StatementBegin

INSERT INTO families (name, latin_name, rotation_interval_years, nitrogen_fixing) VALUES
    ('Solanacées',       'Solanaceae',     4, false),
    ('Brassicacées',     'Brassicaceae',   4, false),
    ('Fabacées',         'Fabaceae',       3, true),
    ('Apiacées',         'Apiaceae',       3, false),
    ('Alliacées',        'Amaryllidaceae', 3, false),
    ('Cucurbitacées',    'Cucurbitaceae',  3, false),
    ('Astéracées',       'Asteraceae',     2, false),
    ('Chénopodiacées',   'Amaranthaceae',  3, false),
    ('Poacées',          'Poaceae',        2, false),
    ('Polygonacées',     'Polygonaceae',   1, false),
    ('Boraginacées',     'Boraginaceae',   1, false);

INSERT INTO species (family_id, name, nitrogen_demand, maturity_from, spacing_cm, suitable_as_cover, termination_lead_days)
SELECT f.id, s.name, s.demand::nitrogen_demand, s.from_date::maturity_from, s.spacing, s.cover, s.lead
FROM (VALUES
    -- Solanacées
    ('Solanacées',     'Tomate',            'heavy',    'transplant',  60, false, NULL),
    ('Solanacées',     'Aubergine',         'heavy',    'transplant',  60, false, NULL),
    ('Solanacées',     'Poivron',           'heavy',    'transplant',  50, false, NULL),
    ('Solanacées',     'Pomme de terre',    'heavy',    'sowing',      35, false, NULL),
    -- Brassicacées
    ('Brassicacées',   'Chou pommé',        'heavy',    'transplant',  50, false, NULL),
    ('Brassicacées',   'Chou-fleur',        'heavy',    'transplant',  60, false, NULL),
    ('Brassicacées',   'Brocoli',           'heavy',    'transplant',  50, false, NULL),
    ('Brassicacées',   'Navet',             'moderate', 'sowing',      15, false, NULL),
    ('Brassicacées',   'Radis',             'light',    'sowing',       5, false, NULL),
    ('Brassicacées',   'Roquette',          'light',    'sowing',      10, false, NULL),
    ('Brassicacées',   'Moutarde blanche',  'light',    'sowing',       5, true,    14),
    ('Brassicacées',   'Radis fourrager',   'light',    'sowing',       5, true,    14),
    -- Fabacées
    ('Fabacées',       'Haricot',           'light',    'sowing',      10, false, NULL),
    ('Fabacées',       'Pois',              'light',    'sowing',       5, false, NULL),
    ('Fabacées',       'Fève',              'light',    'sowing',      15, true,    21),
    ('Fabacées',       'Vesce commune',     'light',    'sowing',       3, true,    21),
    ('Fabacées',       'Trèfle incarnat',   'light',    'sowing',       3, true,    21),
    ('Fabacées',       'Luzerne',           'light',    'sowing',       3, true,    28),
    ('Fabacées',       'Féverole',          'light',    'sowing',      10, true,    21),
    -- Apiacées
    ('Apiacées',       'Carotte',           'moderate', 'sowing',       5, false, NULL),
    ('Apiacées',       'Panais',            'moderate', 'sowing',      15, false, NULL),
    ('Apiacées',       'Céleri',            'heavy',    'transplant',  30, false, NULL),
    ('Apiacées',       'Persil',            'light',    'sowing',      15, false, NULL),
    -- Alliacées
    ('Alliacées',      'Oignon',            'moderate', 'transplant',  12, false, NULL),
    ('Alliacées',      'Échalote',          'moderate', 'sowing',      15, false, NULL),
    ('Alliacées',      'Ail',               'light',    'sowing',      12, false, NULL),
    ('Alliacées',      'Poireau',           'heavy',    'transplant',  15, false, NULL),
    -- Cucurbitacées
    ('Cucurbitacées',  'Courgette',         'heavy',    'transplant',  90, false, NULL),
    ('Cucurbitacées',  'Courge',            'heavy',    'transplant', 100, false, NULL),
    ('Cucurbitacées',  'Concombre',         'heavy',    'transplant',  50, false, NULL),
    -- Astéracées
    ('Astéracées',     'Laitue',            'moderate', 'transplant',  25, false, NULL),
    -- Chénopodiacées
    ('Chénopodiacées', 'Betterave',         'moderate', 'sowing',      10, false, NULL),
    ('Chénopodiacées', 'Épinard',           'moderate', 'sowing',      10, false, NULL),
    ('Chénopodiacées', 'Blette',            'moderate', 'sowing',      30, false, NULL),
    -- Poacées
    ('Poacées',        'Maïs doux',         'heavy',    'transplant',  40, false, NULL),
    ('Poacées',        'Seigle fourrager',  'light',    'sowing',       3, true,    28),
    -- Polygonacées
    ('Polygonacées',   'Sarrasin',          'light',    'sowing',       5, true,    14),
    -- Boraginacées
    ('Boraginacées',   'Phacélie',          'light',    'sowing',       3, true,    21)
) AS s (family, name, demand, from_date, spacing, cover, lead)
JOIN families f ON f.name = s.family;

INSERT INTO varieties (species_id, name, days_to_maturity)
SELECT sp.id, v.name, v.days
FROM (VALUES
    ('Tomate',           'Marmande',            80),
    ('Tomate',           'Cœur de bœuf',        85),
    ('Tomate',           'Cerise Gardener''s Delight', 65),
    ('Aubergine',        'Violette de Barbentane', 85),
    ('Poivron',          'Doux long des Landes', 80),
    ('Pomme de terre',   'Charlotte',          100),
    ('Pomme de terre',   'Bintje',             120),
    ('Chou pommé',       'Cœur de bœuf des Vertus', 90),
    ('Chou-fleur',       'Merveille de toutes saisons', 100),
    ('Brocoli',          'Calabrais',           80),
    ('Navet',            'Boule d''or',         60),
    ('Radis',            'De 18 jours',         25),
    ('Radis',            'Gaudry',              30),
    ('Roquette',         'Cultivée',            40),
    ('Haricot',          'Contender',           55),
    ('Haricot',          'Coco de Paimpol',     90),
    ('Pois',             'Merveille de Kelvedon', 70),
    ('Fève',             'Aguadulce',           120),
    ('Carotte',          'Nantaise',            80),
    ('Carotte',          'De Colmar',          110),
    ('Panais',           'Demi-long de Guernesey', 130),
    ('Céleri',           'Vert d''Elne',       120),
    ('Persil',           'Géant d''Italie',     70),
    ('Oignon',           'Jaune paille des Vertus', 120),
    ('Échalote',         'Grise',              120),
    ('Ail',              'Blanc de Lomagne',   240),
    ('Poireau',          'Bleu de Solaise',    150),
    ('Courgette',        'Verte des Maraîchers', 55),
    ('Courge',           'Butternut',          110),
    ('Courge',           'Potimarron',         100),
    ('Concombre',        'Marketmore',          60),
    ('Laitue',           'Batavia Reine des Glaces', 60),
    ('Laitue',           'Merveille des quatre saisons', 65),
    ('Betterave',        'Crapaudine',          95),
    ('Épinard',          'Géant d''hiver',      50),
    ('Blette',           'Verte à carde blanche', 60),
    ('Maïs doux',        'Golden Bantam',       90),
    -- Cover crops: no days_to_maturity, their end date is driven backwards from
    -- when the bed is next needed.
    ('Moutarde blanche', 'Commune',           NULL),
    ('Radis fourrager',  'Structurator',      NULL),
    ('Vesce commune',    'De printemps',      NULL),
    ('Trèfle incarnat',  'Commun',            NULL),
    ('Luzerne',          'Commune',           NULL),
    ('Féverole',         'De printemps',      NULL),
    ('Seigle fourrager', 'Commun',            NULL),
    ('Sarrasin',         'Commun',            NULL),
    ('Phacélie',         'Commune',           NULL)
) AS v (species, name, days)
JOIN species sp ON sp.name = v.species;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM varieties;
DELETE FROM species;
DELETE FROM families;
-- +goose StatementEnd
