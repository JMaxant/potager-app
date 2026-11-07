-- +goose Up
CREATE TABLE families
(
    id                      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                    text CHECK ( name <> '' ) UNIQUE                NOT NULL,
    rotation_interval_years smallint CHECK ( rotation_interval_years >= 0 ) NOT NULL,
    nitrogen_fixing         BOOL        DEFAULT false                       NOT NULL,
    created_at              timestamptz DEFAULT now()                       NOT NULL
);

CREATE TYPE nitrogen_demand AS ENUM ('heavy', 'moderate', 'light');
CREATE TYPE maturity_from AS ENUM ('sowing', 'transplant');
CREATE TABLE species
(
    id                    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    family_id             bigint REFERENCES families (id) ON DELETE RESTRICT NOT NULL,
    name                  text CHECK ( name <> '' ) UNIQUE                   NOT NULL,
    suitable_as_cover     BOOL DEFAULT false                                 NOT NULL,
    termination_lead_days smallint,
    nitrogen_demand       nitrogen_demand                                    NOT NULL,
    maturity_from         maturity_from                                      NOT NULL,
    CONSTRAINT termination_lead_days_cover_only CHECK (
        termination_lead_days IS NULL                     -- Nothing is provided.
            OR
        (suitable_as_cover AND termination_lead_days > 0) -- A value is provided: valid if it is a cover crop and its termination delay is > 0.
        )
);

CREATE TABLE varieties
(
    id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    species_id       bigint REFERENCES species (id) ON DELETE RESTRICT NOT NULL,
    name             text CHECK ( name <> '' )                         NOT NULL,
    days_to_maturity smallint CHECK ( days_to_maturity > 0 ),
    UNIQUE (species_id, name)
);

CREATE TYPE plot_type AS ENUM ('open_ground', 'raised_bed', 'greenhouse', 'coldframe');
CREATE TABLE plots
(
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        varchar(80) CHECK ( name <> '' ) NOT NULL,
    length_cm   smallint CHECK ( length_cm BETWEEN 1 AND 10000),
    width_cm    smallint CHECK ( width_cm BETWEEN 1 AND 10000 ),
    plot_type   plot_type                        NOT NULL DEFAULT 'open_ground',
    position    smallint                         NOT NULL,
    archived_at timestamptz
);

CREATE TYPE planting_purpose AS ENUM ('food', 'cover');
CREATE TABLE plantings
(
    id                     bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    plot_id                bigint REFERENCES plots (id) ON DELETE RESTRICT NOT NULL,
    purpose                planting_purpose                                NOT NULL,
    sown_at                date,
    sown_in_place          bool                                            NOT NULL,
    transplanted_at        date,
    ended_at               date,
    expected_end_at        date,
    maturity_days_override smallint CHECK ( maturity_days_override > 0 ),
    notes                  text CHECK ( notes <> '' )
);
-- +goose Down
DROP TABLE plantings;
DROP TYPE planting_purpose;
DROP TABLE plots;
DROP TYPE plot_type;
DROP TABLE varieties;
DROP TABLE species;
DROP TYPE nitrogen_demand;
DROP TYPE maturity_from;
DROP TABLE families;
