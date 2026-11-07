-- +goose Up
CREATE TABLE families
(
    id                      bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                    text CHECK ( name <> '' ) UNIQUE                NOT NULL,
    rotation_interval_years smallint CHECK ( rotation_interval_years >= 0 ) NOT NULL,
    nitrogen_fixing         BOOL        DEFAULT false                       NOT NULL,
    created_at              timestamptz DEFAULT now()                       NOT NULL
);

CREATE TABLE species
(
    id                    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    family_id             bigint REFERENCES families (id) ON DELETE RESTRICT NOT NULL,
    name                  text CHECK ( name <> '' ) UNIQUE                   NOT NULL,
    suitable_as_cover     BOOL DEFAULT false                                 NOT NULL,
    termination_lead_days smallint,
    CONSTRAINT is_cover_crop CHECK (
        termination_lead_days IS NULL                     -- Nothing is provided.
            OR
        (suitable_as_cover AND termination_lead_days > 0) -- A value is provided: valid if it is a cover crop and its termination delay is > 0.
        )
);
-- +goose Down
DROP TABLE species;
DROP TABLE families;
