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

- [ ] **`LoadLocation` on every call** in `DateFromTime`: load the zone **once**
      into a package-level variable (`init()`, or `var paris = mustLoad(...)` in
      the spirit of `template.Must` / `regexp.MustCompile`). A missing tzdata then
      fails at start-up rather than on the first request. Case 5 of the test
      below needs this variable.
- [ ] **`compare` checks equality**, but `Status` needs the *order*. Replace it
      with `Before(other Date) bool` or `Compare(other Date) int` (-1/0/+1, as
      in `cmp.Compare` and `time.Time.Compare`). `==` is enough for equality.
- [ ] `addDays`: a single line, `d.t.AddDate(0, 0, days)` wrapped in a `Date`.
      Write it **after** its test.
- [ ] Optional: rename `DateFromInt` to `NewDate` or `DateOf` (the name used by
      `civil`).
- [ ] Before deploying (T11): `import _ "time/tzdata"`, because a minimal Docker
      image has no tzdata.

## Current step: `TestDateFromTime`

File: `internal/potager/date_test.go`.

Cases (rule: *« `Date` is a calendar date without time or zone »* + ADR 0003):

1. **(model)** 22h30 UTC le 28/03/2026 = 23h30 à Paris (heure d'hiver) → **28/03**.
2. 23h30 UTC le 28/03/2026 = 00h30 à Paris → **29/03** (le jour bascule).
3. 22h30 UTC le 15/07/2026 = 00h30 à Paris (heure d'été, +2) → **16/07**.
4. Jour du changement d'heure, 29/03/2026 : 01h30 UTC = 03h30 à Paris → **29/03**.
5. Entrée déjà à l'heure de Paris : 9h00 le 12/10/2025 (le semis d'ail de la spec) → **12/10/2025**.

Model case:

```go
package potager

import (
	"testing"
	"time"
)

func TestDateFromTime(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want Date
	}{
		{
			name: "late evening in winter stays on the same day in Paris",
			in:   time.Date(2026, time.March, 28, 22, 30, 0, 0, time.UTC), // 23:30 CET
			want: DateFromInt(2026, time.March, 28),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DateFromTime(tt.in)
			if got != tt.want {
				t.Errorf("DateFromTime(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
```

The comparison uses `!=`, not `Equal`, on purpose: it also checks that the
constructors produce comparable `Date` values.

**Done when:**

- `go test ./internal/potager/ -run TestDateFromTime -v`: **5 cases, all PASS**.
- **Sabotage:** remove `.In(loc)` from `DateFromTime`. Cases 2 and 3 must fail,
  cases 1 and 4 still pass. Explain why, then restore.

Then send Claude the test file and the explanation of the sabotage.

## Next step

Test `addDays` and the comparison (`Before` or `Compare`), then move on to
`Season`.
