# Plot occupancy is not an exclusion constraint

Undersowing and intercropping deliberately put two plantings in one plot at the same time, so Postgres is not allowed to refuse overlap. The GiST index on `(plot_id, occupancy)` exists for query speed only, and overlap becomes an advisory notice from step 2. A reader seeing a `daterange` with a GiST index will expect `EXCLUDE USING gist`; this is the reason it is absent.
