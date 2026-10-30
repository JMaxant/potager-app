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

## Done so far

`Date` (`DateFromTime`, `DateFromInt`, `String`, `addDays`, `before`), `Season`.

## Current step: `TestNetEffect`

File: `internal/potager/neteffect.go` and `neteffect_test.go`. Rule: the table
in `docs/design.md` (*Net soil effect*). No dates involved.

New concept: **enumerations in Go** — a named type plus constants.
`type Purpose string` with `const PurposeFood Purpose = "food"`. String-based
rather than `iota`, because the values mirror the Postgres columns
(`plantings.purpose`, `species.nitrogen_demand`) and read well in test output.

Types (names to confirm):

- `Purpose`: `PurposeFood`, `PurposeCover`.
- `Demand`: `DemandHeavy`, `DemandModerate`, `DemandLight`.
- `Effect`: `EffectEnriching`, `EffectNeutral`, `EffectSlightlyDepleting`, `EffectDepleting`.

Signature: `func NetEffect(purpose Purpose, fixesNitrogen bool, demand Demand) Effect`.
`fixesNitrogen` comes from `families.nitrogen_fixing`. One species at a time;
a mix of several varieties comes later.

Cases (one per table row, plus the "any" that matters):

1. **(model)** Engrais vert fixateur (féverole) → **enriching**.
2. Engrais vert non fixateur, gourmand (moutarde) → **neutral**.
3. Vivrier fixateur (pois récolté), même gourmand → **neutral**.
4. Vivrier non fixateur, gourmand (tomate) → **depleting**.
5. Vivrier non fixateur, moyen → **slightly depleting**.
6. Vivrier non fixateur, faible → **slightly depleting**.

Model case:

```go
func TestNetEffect(t *testing.T) {
	tests := []struct {
		name          string
		purpose       Purpose
		fixesNitrogen bool
		demand        Demand
		want          Effect
	}{
		{
			name:          "a nitrogen-fixing cover crop enriches the bed",
			purpose:       PurposeCover,
			fixesNitrogen: true,
			demand:        DemandModerate,
			want:          EffectEnriching,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NetEffect(tt.purpose, tt.fixesNitrogen, tt.demand)
			if got != tt.want {
				t.Errorf("NetEffect(%v, %v, %v) = %v, want %v",
					tt.purpose, tt.fixesNitrogen, tt.demand, got, tt.want)
			}
		})
	}
}
```

**Done when:**

- `go test ./internal/potager/ -run TestNetEffect -v`: **6 cases, all PASS**.
- **Sabotage 1:** a fixing species always enriches, whatever its purpose. Only
  case 3 fails: the row a single "nitrogen" enum would get wrong.
- **Sabotage 2:** `moderate` treated like `heavy`. Only case 5 fails.

## Where `NetEffect` stands (to pick up here)

Done: the three types, `NetEffect` as a flat `switch` (one `case` per table
row, `default` panics), 6 table cases passing.

To do, in order:

- [ ] **`TestNetEffect`**: cases 2 and 3 back to `DemandHeavy` (moutarde, pois
      récolté), as planned. With `DemandModerate` they would not catch a
      `switch demand` back in the cover branch.
- [ ] Remove the `"an invalid value is passed to NetEffect"` case from the
      table: it has no `want` and makes the whole test panic.
      Check: `go test ./internal/potager/ -run 'TestNetEffect$' -v` → **6 PASS**.
- [ ] **`TestNetEffectPanics`**, a separate test, `defer` + `recover`
      (Go blog *Defer, Panic, and Recover*; spec *Handling panics*):
  1. **(model)** `purpose` inconnu (`"orchard"`) → panic.
  2. `purpose` vide `""` (colonne `NULL` mal lue) → panic.
  3. Vivrier non fixateur, `demand` vide `""` → panic. **Fails at first**: the
     `case purpose == PurposeFood:` line catches every demand.

  ```go
  func TestNetEffectPanics(t *testing.T) {
  	tests := []struct {
  		name          string
  		purpose       Purpose
  		fixesNitrogen bool
  		demand        Demand
  	}{
  		{
  			name:          "an unknown purpose panics",
  			purpose:       Purpose("orchard"),
  			fixesNitrogen: false,
  			demand:        DemandHeavy,
  		},
  	}
  	for _, tt := range tests {
  		t.Run(tt.name, func(t *testing.T) {
  			defer func() {
  				if r := recover(); r == nil {
  					t.Errorf("NetEffect(%q, %v, %q) did not panic",
  						tt.purpose, tt.fixesNitrogen, tt.demand)
  				}
  			}()
  			NetEffect(tt.purpose, tt.fixesNitrogen, tt.demand)
  		})
  	}
  }
  ```

- [ ] **Green**: write the last food line with explicit values,
      `demand == DemandModerate || demand == DemandLight`.
      Check: `go test ./internal/potager/ -run TestNetEffectPanics -v` → **3 PASS**.
- [ ] Panic message: lower case, with the values through `fmt.Sprintf` and
      `%q`, e.g. `netEffect: unknown purpose "orchard" or demand ""`.
- [ ] **Sabotage**: `default` returns `EffectNeutral` instead of panicking →
      cases 1 and 2 fail with "did not panic". Restore.
- [ ] The two sabotages planned above for the table (fixing always enriches;
      moderate treated as heavy).
- [ ] Commit (`refs #8`, with the `#`).

## Next step: `ExpectedEnd`

Order changed on 2026-10-30: `ExpectedEnd` before `Status`, as planned above.

Rule (`step-1-spec.md`, *Expected end*; `design.md`; AC 5 and 6):

1. A stored `expected_end_at` always wins, food or cover. Shown "(saisie)".
2. Cover crop without a stored value: no expected end.
3. Food crop: start + maturity days, the start being the date `maturity_from`
   names (`sowing` or `transplant`).
4. No fallback: a `transplant` species without a transplant date has no
   expected end, and the screen asks for that date.
5. Days: `maturity_days_override` if set, else `varieties.days_to_maturity`
   (nullable). No `MAX` in step 1: one variety.

### Decision: optional values use the zero value

| Value | Absent = | Why |
| --- | --- | --- |
| `Date` | `Date{}`, checked with `IsZero()` | `Date` is comparable, no `nil` to dereference, same idiom as `time.Time.IsZero`. 1 January of year 1 is the sentinel; no garden reaches it. |
| maturity days (`int`) | `0` | A zero-day maturity is meaningless, so `0` cannot be real data. |

Checked: `DateFromTime(time.Time{}) == Date{}` (Paris in year 1 is +00:09:21,
same day), and `DateFromInt(1, time.January, 1) == Date{}`.

The zero value only works where zero cannot be real data: it would not work for
a `bool` such as `fixesNitrogen`.

### Steps

- [x] `Date.IsZero()`, test-first: `TestDateIsZero`, 4 cases, in `date_test.go`.
- [ ] `MaturityFrom` enum (`sowing`, `transplant`), like `Purpose`.
- [ ] Signature and result shape (computed date, stored date, nothing, and
      whether to ask for the transplant date).
- [ ] `TestExpectedEnd`, then the implementation.
