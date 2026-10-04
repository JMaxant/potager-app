# "Today" is a parameter, never `CURRENT_DATE`

Queries that depend on the current date receive it as an argument, computed in Go in the `Europe/Paris` zone from an injectable clock. `CURRENT_DATE` follows the database session zone, which is UTC on the VPS, so the database and the application would disagree around midnight, and integration tests could not pin the date. The price is one extra parameter on a handful of queries.
