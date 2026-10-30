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

`Date` (`DateFromTime`, `DateFromInt`, `String`, `IsZero`, `addDays`, `before`), `Season`, `NetEffect`. In progress: `ExpectedEnd` (see "Pick up here" below).

## `NetEffect` plan (done, kept for reference)

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

## Where `NetEffect` stands: done

`NetEffect` as a flat `switch`, one `case` per table row. Two panics with their
own message (`fmt.Sprintf` + `%q`): `NetEffect: unknown demand` for a food crop
whose demand matched no row, `NetEffect: unknown purpose` in the `default`.
`TestNetEffect` 6 cases, `TestNetEffectPanics` 3 cases (unknown purpose, unknown
demand, empty demand), sabotages done, committed.

Left over, small:

- [ ] `NetEffect` doc comment: "May panic if" → "It panics if" (the behaviour
      is guaranteed, as in `regexp.MustCompile`).
- [ ] Optional: a `TestNetEffectPanics` case for an empty purpose `""` (a
      `NULL` column read badly). Already covered by the `default`; the case
      would document it.

## Current step: `ExpectedEnd`

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
- [x] `MaturityFrom` enum (`sowing`, `transplant`), like `Purpose`, in
      `expectedend.go`.
- [x] Signature and result shape, decided below.
- [x] `Planting` in `planting.go` (own file: a central domain noun, used by
      `ExpectedEnd`, then `Status`, the store and `web`). Exported fields: the
      store fills them from another package.
- [x] `ExpectedEndSource` enum in `expectedend.go`, four values. A `string`
      like the other enums, although it is never stored: readable test output.
- [x] `CONTEXT.md`: *Catalogue* and *Expected end* added to the glossary.

**Pick up here (stopped on 2026-10-30, evening).** Uncommitted:
`planting.go`, `ExpectedEndSource`, `CONTEXT.md`, these notes.

- [ ] **Doc comments**, before the commit:
  - `ExpectedEndSource`: "tells where the source of expected end date" does
    not parse. E.g. "ExpectedEndSource tells where the date returned by
    ExpectedEnd comes from."
  - `ExpectedEndSourceNeedsTransplant`: a sentence after the name, e.g.
    "ExpectedEndSourceNeedsTransplant means there is no date and the screen
    asks for the transplant date (AC 5)." Its name gives the cause, not what
    the screen must do; without it, `web` could treat it like `None`.
  - `MaturityFrom`: "the step" → "the date of the planting"; add the
    no-fallback rule ("a transplant species without a transplant date has no
    expected end").
  - Optional, `Planting`: "Maturity days fields" → name them, `MaturityFrom`
    and `DaysToMaturity`.
- [ ] Commit (`refs #8`).
- [ ] **Stub**: `ExpectedEnd(p Planting) (Date, ExpectedEndSource)` returning
      `Date{}, ExpectedEndSourceNone`, only so the test compiles.
- [ ] **`TestExpectedEnd`** in `expectedend_test.go`: the model case below,
      then cases 2 to 10 of the plan.
      Check: `go test ./internal/potager/ -run TestExpectedEnd -v` → **10
      cases, 2 PASS** (8 and 10, the only ones expecting "none"). Any other
      green case does not test what it claims.
- [ ] Implementation, then the three sabotages listed under the plan.
- [ ] Decide the open edge case (`sowing` species with only a transplant
      date), write it in the spec, add the case.
- [ ] Panics: an unknown `Purpose` or `MaturityFrom`, like `TestNetEffectPanics`.
- [ ] Then `Status`: it receives the expected end computed here.

Model case:

```go
func TestExpectedEnd(t *testing.T) {
	tests := []struct {
		name       string
		planting   Planting
		wantDate   Date
		wantSource ExpectedEndSource
	}{
		{
			name: "an 80-day tomato transplanted on 10/05/2026 is expected on 29/07/2026 (AC 5)",
			planting: Planting{
				Purpose:        PurposeFood,
				MaturityFrom:   MaturityFromTransplant,
				TransplantedAt: DateFromInt(2026, time.May, 10),
				DaysToMaturity: 80,
			},
			wantDate:   DateFromInt(2026, time.July, 29),
			wantSource: ExpectedEndSourceComputed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, gotSource := ExpectedEnd(tt.planting)
			if gotDate != tt.wantDate || gotSource != tt.wantSource {
				t.Errorf("ExpectedEnd(%+v) = (%v, %q), want (%v, %q)",
					tt.planting, gotDate, gotSource, tt.wantDate, tt.wantSource)
			}
		})
	}
}
```

Fields left out of a struct literal take their zero value, i.e. "absent":
each case lists only what matters to it. `%+v` prints field names and calls
`Date.String()` on nested dates.

### Decision: signature

`func ExpectedEnd(p Planting) (Date, ExpectedEndSource)`

- **Input: `Planting`**, the domain noun (`CONTEXT.md`). Five fields come from
  `plantings`, two are catalogue facts borrowed from the species and the
  variety. In step 1 a planting has one variety, so those two sit flat in the
  struct; with mixes (step 2) they become a list of components and a `MAX`.
  A struct rather than seven arguments: three are `Date` and two are `int`, and
  swapping `SownAt` and `TransplantedAt` would compile silently.
- **Output: `(Date, ExpectedEndSource)`**: the date (`Date{}` when there is
  none) and where it comes from. An enum rather than two `bool`s: four values,
  no impossible combination such as "stored and needs a transplant date".

Names (all confirmed and declared):

| `Planting` field | Column |
| --- | --- |
| `Purpose` | `plantings.purpose` |
| `MaturityFrom` | `species.maturity_from` |
| `SownAt` | `plantings.sown_at` |
| `TransplantedAt` | `plantings.transplanted_at` |
| `ExpectedEndAt` | `plantings.expected_end_at` (entered by hand) |
| `DaysToMaturity` | `varieties.days_to_maturity` |
| `MaturityDaysOverride` | `plantings.maturity_days_override` |

| `ExpectedEndSource` | Meaning | Screen |
| --- | --- | --- |
| `ExpectedEndSourceComputed` | start + days | « Récolte prévue : 29/07/2026 » |
| `ExpectedEndSourceStored` | entered by hand | « … (saisie) » |
| `ExpectedEndSourceNone` | nothing to show | nothing |
| `ExpectedEndSourceNeedsTransplant` | transplant species, no transplant date | « Indiquez la date de plantation… » |

### `TestExpectedEnd` plan

1. **(model)** AC 5: tomate repiquée le 10/05/2026, 80 jours → 29/07/2026, computed.
2. AC 5: la même tomate, seulement semée → `Date{}`, needs transplant.
3. Espèce `sowing` (radis) semée le 01/04/2026, 30 jours → 01/05/2026, computed.
4. Espèce `transplant` semée **et** repiquée → compte depuis la plantation.
5. Override: tomate repiquée le 10/05/2026, variété 80, override 90 → 08/08/2026, computed.
6. Saisie sur une culture vivrière: 01/07/2026 saisie l'emporte sur le 29/07 calculé → stored.
7. Saisie sans date de plantation (espèce `transplant`, seulement semée) → stored: rule 1 comes first.
8. AC 6: engrais vert sans saisie → `Date{}`, none.
9. Engrais vert avec saisie → stored.
10. Culture vivrière, variété sans durée, pas d'override → `Date{}`, none.

Edge case left open (spec silent): a `sowing` species with only a transplant
date. Validation lets it through; decide before adding a case, and write the
answer in the spec.

Sabotages:

- fall back to the sowing date → case 2 fails;
- ignore the override → case 5 fails;
- check the stored value after the computation → case 6 fails.
