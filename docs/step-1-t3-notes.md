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
- [ ] **`compare` checks equality**, but `Status` needs the *order*. Replace it
      with `Before(other Date) bool` or `Compare(other Date) int` (-1/0/+1, as
      in `cmp.Compare` and `time.Time.Compare`). `==` is enough for equality.
- [ ] `addDays`: a single line, `d.t.AddDate(0, 0, days)` wrapped in a `Date`.
      Write it **after** its test.
- [ ] Optional: rename `DateFromInt` to `NewDate` or `DateOf` (the name used by
      `civil`).
- [ ] Before deploying (T11): `import _ "time/tzdata"`, because a minimal Docker
      image has no tzdata.

## Current step: `TestAddDays`

File: `internal/potager/date_test.go`. `addDays` is how `ExpectedEnd` counts
maturity from a start date (`docs/design.md`), and how a termination date is
counted back from the next planting (negative days).

Cases:

1. **(model)** AC 5 : tomate repiquée le 10/05/2026, variété à 80 jours → **29/07/2026**.
2. + 0 jour → **même date**.
3. Changement d'année : 15/12/2025 + 20 → **04/01/2026**.
4. Année bissextile : 28/02/2028 + 1 → **29/02/2028**.
5. Année non bissextile : 28/02/2026 + 1 → **01/03/2026**.
6. Jours négatifs (date de destruction comptée à rebours) : 10/05/2026 − 21 → **19/04/2026**.

Model case:

```go
func TestAddDays(t *testing.T) {
	tests := []struct {
		name string
		in   Date
		days int
		want Date
	}{
		{
			name: "80-day tomato transplanted on 10 May is ready on 29 July (AC 5)",
			in:   DateFromInt(2026, time.May, 10),
			days: 80,
			want: DateFromInt(2026, time.July, 29),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.addDays(tt.days)
			if got != tt.want {
				t.Errorf("%v.addDays(%d) = %v, want %v", tt.in, tt.days, got, tt.want)
			}
		})
	}
}
```

Red first: uncomment the `addDays` stub (`panic("todo")`) so the test compiles
and fails, then write the one-line body.

**Done when:**

- `go test ./internal/potager/ -run TestAddDays -v`: **6 cases, all PASS**.
- **Sabotage 1:** `AddDate(0, days, 0)` (months instead of days). Every case but
  case 2 fails. Explain why case 2 survives.
- **Sabotage 2:** `d.t.Add(time.Duration(days) * 24 * time.Hour)`. Nothing
  fails. Explain why, and what would break if `t` were midnight in Paris.

## Next step

Test the comparison (`Before` or `Compare`), then move on to `Season`.
