# T3 — domain rules: working notes

Working notes to pick up task 3 (issue #8) on another machine. Delete this file
once T3 is merged. Branch: `feature/domain_rules_8`.

## Plan

The rules build on one another, in this order:

1. **`Date`** — in progress (see below).
2. **`Season`** — the year of `ended_at`.
3. **`NetEffect`** — the table in `docs/design.md`, no dates involved.
4. **`ExpectedEnd`** — stored value, `maturity_from`, cover crop.
5. **`Status`** — four states, depends on `ExpectedEnd` and on a "today"
   passed as a parameter (ADR 0003).

For each rule: the test plan in French plus one model case from Claude, then the
owner writes the remaining cases, then the implementation.

## Decisions taken on `Date`

| Decision | Why |
| --- | --- |
| `struct{ t time.Time }`, unexported field, **not** embedding | Embedding promotes every `time.Time` method (`Hour`, `Add`, `In`, `Before(time.Time)`…): the "no time, no zone" rule would only be a convention. A private field exposes only the methods we write. Embedding suits a decorator that passes everything through (`statusRecorder`); a private field suits a value type that restricts. |
| `t` is always **midnight UTC** | With `time.Date(y, m, d, 0, 0, 0, 0, time.UTC)`, two `Date` values for the same day are bit-for-bit identical, so `==` works. Checked: midnight in `Europe/Paris` gives `==` false, because every `LoadLocation` returns a new `*time.Location` pointer. |
| Paris appears **in one place only**: `DateFromTime` | The zone only matters to turn an instant into a calendar day (ADR 0003). `DateFromInt` needs no zone. |
| `DateFromTime` = `In(paris).Date()`, then `DateFromInt` | A single function builds the struct, so the "midnight UTC" rule lives in one place. |
| No `Year()` for now | `Season` lives in `potager` and can read `d.t.Year()`: `t` is private to the package, not to the type. Add it when `web` needs it. |
| `addDays` / `compare` in lower case | Fine while `ExpectedEnd` and `Status` live in `potager`. The capital letter decides visibility: it is a scope decision, not a style one. |

## Still open in `internal/potager/date.go`

- [x] Zone loaded once into the package variable `locParis` (`init()`).
- [x] `TestDateFromTime`: 5 cases, both sabotages done (`.In` removed, `.UTC()`).
- [x] `addDays`: `TestAddDays`, 6 cases. Negative days kept (backward
      termination in `design.md`); positive maturity is checked at input.
- [x] `before(other Date) bool`, strict, lower case like `addDays`: `TestBefore`,
      4 cases. No `after`: `b.before(a)` says it.
- [ ] Optional: rename `DateFromInt` to `NewDate` or `DateOf` (the name used by
      `civil`).
- [ ] Before deploying (T11): `import _ "time/tzdata"`, because a minimal Docker
      image has no tzdata.

## Current step: `TestSeason`

File: `internal/potager/season_test.go` (or `date_test.go`). Rule: *« Season is
the year of `ended_at` »*, used only for closed plantings (history section).

Signature: `func Season(endedAt Date) int`. No `Planting` type yet; the caller
only asks for the season of a closed planting, so there is no optional value to
handle here. The optional `ended_at` comes with `Status`.

Cases:

1. **(model)** AC 8 : ail semé le 12/10/2025, récolté le 04/07/2026 → **2026** (seule la fin compte).
2. Fin le 01/01/2026 → **2026**.
3. Fin le 31/12/2025 → **2025**.

Model case:

```go
func TestSeason(t *testing.T) {
	tests := []struct {
		name    string
		endedAt Date
		want    int
	}{
		{
			name:    "garlic sown in October 2025 and lifted in July 2026 is season 2026 (AC 8)",
			endedAt: DateFromInt(2026, time.July, 4),
			want:    2026,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Season(tt.endedAt)
			if got != tt.want {
				t.Errorf("Season(%v) = %d, want %d", tt.endedAt, got, tt.want)
			}
		})
	}
}
```

**Done when:**

- `go test ./internal/potager/ -run TestSeason -v`: **3 cases, all PASS**.
- **Sabotage:** `endedAt.addDays(-1).t.Year()`. Only case 2 fails. Explain why.

## Next step

`NetEffect`, the table in `docs/design.md`.
