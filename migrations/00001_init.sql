-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TYPE nitrogen_demand AS ENUM ('heavy', 'moderate', 'light');
CREATE TYPE maturity_from AS ENUM ('sowing', 'transplant');
CREATE TYPE plot_type AS ENUM ('open_ground', 'raised_bed', 'greenhouse', 'coldframe');
CREATE TYPE planting_purpose AS ENUM ('food', 'cover');

-- Botanical family. Carries the rotation interval (how many years before this
-- family should return to the same plot) and nitrogen fixation, which is a
-- family-level trait: every Fabaceae nodulates, nothing else here does.
CREATE TABLE families (
    id                      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                    text    NOT NULL UNIQUE,
    latin_name              text    NOT NULL UNIQUE,
    rotation_interval_years int     NOT NULL CHECK (rotation_interval_years > 0),
    nitrogen_fixing         boolean NOT NULL DEFAULT false
);

-- A species is the catalogue level that carries agronomic behaviour.
-- termination_lead_days: how long before the bed is needed a cover crop of this
-- species must be terminated. Rye is allelopathic and decomposes slowly, so it
-- wants far more lead time than mustard.
CREATE TABLE species (
    id                    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    family_id             bigint          NOT NULL REFERENCES families (id) ON DELETE RESTRICT,
    name                  text            NOT NULL UNIQUE,
    nitrogen_demand       nitrogen_demand NOT NULL,
    maturity_from         maturity_from   NOT NULL,
    spacing_cm            int             CHECK (spacing_cm > 0),
    suitable_as_cover     boolean         NOT NULL DEFAULT false,
    termination_lead_days int             CHECK (termination_lead_days >= 0)
);

CREATE INDEX species_family_id_idx ON species (family_id);

-- days_to_maturity is null for species only ever grown as a cover crop: their
-- end date is driven backwards from when the bed is needed, not forwards from
-- sowing.
CREATE TABLE varieties (
    id               bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    species_id       bigint NOT NULL REFERENCES species (id) ON DELETE RESTRICT,
    name             text   NOT NULL,
    days_to_maturity int    CHECK (days_to_maturity > 0),
    UNIQUE (species_id, name)
);

CREATE INDEX varieties_species_id_idx ON varieties (species_id);

CREATE TABLE plots (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name        text      NOT NULL,
    type        plot_type NOT NULL DEFAULT 'open_ground',
    length_cm   int       CHECK (length_cm > 0),
    width_cm    int       CHECK (width_cm > 0),
    position    int       NOT NULL DEFAULT 0,
    archived_at timestamptz
);

CREATE INDEX plots_position_idx ON plots (position, id);

-- A planting is one act of sowing: one plot, one set of dates, one intent.
-- expected_end_at is a nullable stored override. When it is null and purpose is
-- 'food' the expected end is computed on read from the components' maturity;
-- when purpose is 'cover' nothing is computed and the form merely suggests a
-- date derived backwards from the plot's next planned planting.
--
-- user_id is unused until step 3 (auth) but present from the start so that
-- adding it later is not a migration across every row.
CREATE TABLE plantings (
    id                     bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    plot_id                bigint           NOT NULL REFERENCES plots (id) ON DELETE RESTRICT,
    purpose                planting_purpose NOT NULL DEFAULT 'food',
    sown_at                date,
    transplanted_at        date,
    ended_at               date,
    expected_end_at        date,
    maturity_days_override int CHECK (maturity_days_override > 0),
    notes                  text NOT NULL DEFAULT '',
    archived_at            timestamptz,
    user_id                bigint,
    created_at             timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT plantings_needs_a_start CHECK (sown_at IS NOT NULL OR transplanted_at IS NOT NULL),

    -- Occupancy of the plot over time. Rotation history, the season overview and
    -- the termination suggestion all read this. Deliberately NOT an exclusion
    -- constraint: undersowing and intercropping mean two plantings legitimately
    -- share a plot and a period.
    occupancy daterange GENERATED ALWAYS AS (
        CASE
            WHEN COALESCE(transplanted_at, sown_at) IS NULL THEN NULL
            ELSE daterange(COALESCE(transplanted_at, sown_at), COALESCE(ended_at, expected_end_at), '[]')
        END
    ) STORED
);

CREATE INDEX plantings_plot_occupancy_idx ON plantings USING gist (plot_id, occupancy);
CREATE INDEX plantings_plot_id_idx ON plantings (plot_id);

-- A planting holds one or more varieties: a cover crop mix is one sowing of
-- several species, not several sowings.
CREATE TABLE planting_varieties (
    planting_id    bigint NOT NULL REFERENCES plantings (id) ON DELETE CASCADE,
    variety_id     bigint NOT NULL REFERENCES varieties (id) ON DELETE RESTRICT,
    quantity       int     CHECK (quantity > 0),
    seed_rate_g_m2 numeric(6, 2) CHECK (seed_rate_g_m2 > 0),
    PRIMARY KEY (planting_id, variety_id)
);

CREATE INDEX planting_varieties_variety_id_idx ON planting_varieties (variety_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS planting_varieties;
DROP TABLE IF EXISTS plantings;
DROP TABLE IF EXISTS plots;
DROP TABLE IF EXISTS varieties;
DROP TABLE IF EXISTS species;
DROP TABLE IF EXISTS families;
DROP TYPE IF EXISTS planting_purpose;
DROP TYPE IF EXISTS plot_type;
DROP TYPE IF EXISTS maturity_from;
DROP TYPE IF EXISTS nitrogen_demand;
-- +goose StatementEnd
