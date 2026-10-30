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
- [ ] **`compare` checks equality**, but `Status` needs the *order*. Decision:
      `before(other Date) bool`, strict, lower case like `addDays`. `==` is
      enough for equality.
- [ ] Optional: rename `DateFromInt` to `NewDate` or `DateOf` (the name used by
      `civil`).
- [ ] Before deploying (T11): `import _ "time/tzdata"`, because a minimal Docker
      image has no tzdata.

## Current step: `TestBefore`

File: `internal/potager/date_test.go`. `Status` uses it twice: *start date in
the future* is `today.before(start)`, *expected end passed* is
`end.before(today)`. It is **strict**: on the expected end day itself the
planting is still « En cours ».

Cases:

1. **(model)** La veille est avant : 28/07/2026 avant 29/07/2026 → **true**.
2. Le lendemain n'est pas avant : 30/07/2026 avant 29/07/2026 → **false**.
3. Même jour → **false** (le jour de la récolte prévue, la plantation est encore « En cours »).
4. Changement d'année : 31/12/2025 avant 01/01/2026 → **true**.

Model case:

```go
func TestBefore(t *testing.T) {
	tests := []struct {
		name  string
		d     Date
		other Date
		want  bool
	}{
		{
			name:  "the day before is before",
			d:     DateFromInt(2026, time.July, 28),
			other: DateFromInt(2026, time.July, 29),
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.before(tt.other)
			if got != tt.want {
				t.Errorf("%v.before(%v) = %v, want %v", tt.d, tt.other, got, tt.want)
			}
		})
	}
}
```

Red first: replace the commented `compare` with a `before` stub that panics.

**Done when:**

- `go test ./internal/potager/ -run TestBefore -v`: **4 cases, all PASS**.
- **Sabotage 1:** `!d.t.After(other.t)` (≤ instead of <). Only case 3 fails.
- **Sabotage 2:** `d.t.Day() < other.t.Day()`. Only case 4 fails. Explain why.

## Next step

`Season`, the year of `ended_at`.
