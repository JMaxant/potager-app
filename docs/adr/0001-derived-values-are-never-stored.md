# Status, season, expected end and soil effect are derived, never stored

Only dates and the gardener's intent are stored. Everything else is computed on read from the dates and the catalogue. A stored value goes stale silently when the catalogue is corrected, whereas a derived one is fixed retroactively. The cost is that the rules must be cheap and pure, which is why they live in `internal/potager` with no database access. The one exception is `expected_end_at`, a deliberate override entered by hand that always wins.
