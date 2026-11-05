# T3 — domain rules: working notes

Working notes to pick up task 3 (issue #8) on another machine. Delete this file
once T3 is merged. Branch: `feature/domain_rules_8`.

## Plan

Scope from `docs/step-1-tasks.md` (widened by #28 / ADR 0008):

1. **`Date`** — done.
2. **`Season`** — done. The year of `ended_at`.
3. **`NetEffect`** — done. The table in `docs/design.md`.
4. **`ExpectedEnd`** — done, panics included (`f1cd089`).
5. **`StartDate`** — done.
6. **Occupancy** — done, `OccupiedFrom` (ADR 0008).
7. **`Status`** — done, with `InPlotBy` (`b62ae4e`).
8. **Area** — **next**, added to T3 on 2026-11-05 (it is in the spec's
   *Derived values*). Plan below.

For each rule: the test plan in French plus one model case from Claude, then the
owner writes the remaining cases, then the implementation, then the sabotages.

## Decisions taken

### `Date`

| Decision                                                     | Why                                                                                                                                                                                                                                                                                                                                   |
|--------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `struct{ t time.Time }`, unexported field, **not** embedding | Embedding promotes every `time.Time` method (`Hour`, `Add`, `In`, `Before(time.Time)`…): the "no time, no zone" rule would only be a convention. A private field exposes only the methods we write. Embedding suits a decorator that passes everything through (`statusRecorder`); a private field suits a value type that restricts. |
| `t` is always **midnight UTC**                               | With `time.Date(y, m, d, 0, 0, 0, 0, time.UTC)`, two `Date` values for the same day are bit-for-bit identical, so `==` works. Checked: midnight in `Europe/Paris` gives `==` false, because every `LoadLocation` returns a new `*time.Location` pointer.                                                                              |
| Paris appears **in one place only**: `DateFromTime`          | The zone only matters to turn an instant into a calendar day (ADR 0003). `DateFromInt` needs no zone.                                                                                                                                                                                                                                 |
| `DateFromTime` = `In(paris).Date()`, then `DateFromInt`      | A single function builds the struct, so the "midnight UTC" rule lives in one place.                                                                                                                                                                                                                                                   |
| No `Year()` for now                                          | `Season` lives in `potager` and can read `d.t.Year()`: `t` is private to the package, not to the type. Add it when `web` needs it.                                                                                                                                                                                                    |
| `addDays` / `before` in lower case                           | Fine while the rules live in `potager`. The capital letter decides visibility: it is a scope decision, not a style one.                                                                                                                                                                                                               |

### Enumerations

`type Purpose string` plus constants, string-based rather than `iota`: the
values mirror the Postgres columns and read well in test output.

A Go "enum" is **not a closed set**: the compiler checks the type, never the
value. `MaturityFrom("anything")` compiles, and so does an untyped literal
(`MaturityFrom: "transplanting"` in a test). Hence:

- every `switch` on an enum lists the known values and **panics in `default`**,
  with `fmt.Sprintf("Func: unknown field %q", v)`;
- a value is checked **only on the path that reads it**: `NetEffect` checks
  `demand` for food crops only, `ExpectedEnd` checks `MaturityFrom` for food
  crops only, after the stored value;
- the real check belongs **at the boundary** (see "With the store" below); the
  domain panic is a safety belt for a programmer error.

### Optional values use the zero value

| Value                 | Absent =                          | Why                                                                                                                                         |
|-----------------------|-----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|
| `Date`                | `Date{}`, checked with `IsZero()` | `Date` is comparable, no `nil` to dereference, same idiom as `time.Time.IsZero`. 1 January of year 1 is the sentinel; no garden reaches it. |
| maturity days (`int`) | `0`                               | A zero-day maturity is meaningless, so `0` cannot be real data.                                                                             |

The zero value only works where zero cannot be real data. It would not work for
an optional `bool`; `SownInPlace` is fine because the column is always set
(defaulted from the species in the form), so `false` is real data, not absence.

### Rules on a `Planting` are methods

Since `b874909`: `p.StartDate()`, `p.OccupiedFrom()`, `p.ExpectedEnd()`, and
`p.Status(today)` to come. Value receivers. A rule that reads only a planting's
fields is a property of that value, like `t.IsZero()` or `d.addDays(n)`; the
spec groups them under *Derived values*. A rule that combines independent
inputs stays a function: `NetEffect(purpose, fixesNitrogen, demand)`. Methods
and their tests live in `planting.go` / `planting_test.go`, tests named
`TestPlanting_Method` (checked by `go vet`). `expectedend.go` keeps the
`MaturityFrom` and `ExpectedEndSource` enums.

`Planting` fields: `MaturityFrom` and `DaysToMaturity` are catalogue facts
(species, variety); every other field comes from `plantings`, including
`MaturityDaysOverride`, entered by the gardener. With mixes (step 2), only the
catalogue facts become a list.

### `ExpectedEnd`

`func (p Planting) ExpectedEnd() (Date, ExpectedEndSource)`. Order of the guard
clauses:

1. `ExpectedEndAt` set → `stored`, nothing else is read.
2. `switch p.Purpose`: cover → `none`; food → go on; else panic.
3. `switch p.MaturityFrom` picks the start date; transplant species with no
   transplant date → `needs transplant` (no fallback, AC 5); else panic.
4. Days: the override, else `DaysToMaturity`.
5. No start date (a bought plant of a sown species) or no days → `none`.
6. `start.addDays(days)` → `computed`.

Tests: `TestPlanting_ExpectedEnd` 13 cases, `TestPlanting_ExpectedEndPanics` 4
cases. Two cases of `TestPlanting_ExpectedEnd` assert that a panic does **not** happen (unknown purpose
with a stored value, cover crop without `MaturityFrom`): they pin where the
checks sit.

### `StartDate`

`func (p Planting) StartDate() Date`: `SownAt` when set, else `TransplantedAt`;
`Date{}` when neither is set. `TestPlanting_StartDate` 5 cases.

Decided on 2026-11-05: no real minimum with `before`. A transplant before the
sowing is **invalid data**, not a case to compute: a plant is always sown
before it is planted out. The rule is enforced at the boundary (form validation,
and a `CHECK (transplanted_at >= sown_at)` in the migration), so `StartDate`
assumes it and says so in its doc comment. Neither an `error` (every caller
would handle a case that cannot happen) nor a panic (one bad row would turn the
plot page into a 500).

Sabotage kept: the scaffold's `COALESCE(transplanted_at, sown_at)` makes
case 1 fail, the bug ADR 0008 rejects.

### `OccupiedFrom`

`func (p Planting) OccupiedFrom() Date`: `SownAt` when `SownInPlace`, else
`TransplantedAt`, so a tray seedling not yet planted out occupies nothing
(`Date{}`). Only the start: the end is `ended_at`, read by `Status`.
`SownInPlace` is recorded, never inferred from `MaturityFrom` (ADR 0008).
Invalid data is not tested: "sown in place" excludes a transplant date in
validation. `TestPlanting_OccupiedFrom` 5 cases.

Sabotages: the scaffold's `COALESCE` makes the tray-seedling case fail;
inferring from `MaturityFrom` makes the "sowing species sown in a tray" case
fail.

### `Status` and `InPlotBy`

`func (p Planting) Status(today Date) Status`: the spec's table, read top-down,
first match wins:

1. `EndedAt` set → `ended`. Only the gardener ends a planting: an expected end
   that has passed makes it `due`, never `ended`.
2. `!p.InPlotBy(today)` → `planned`.
3. `ExpectedEnd()` set and strictly before today → `due`. On the expected day
   itself the planting is still `growing` (decided 2026-11-05).
4. Otherwise `growing`.

`Status` combines `InPlotBy` and `ExpectedEnd`; it never reads `MaturityFrom`,
`Purpose` or the raw dates. A first version that re-derived occupancy from
`MaturityFrom` passed 9 cases and was wrong on five situations: the reason for
cases 10 to 12.

- `Status` is a `string` enum (`planned`, `growing`, `due`, `ended`) although
  it is never stored: the zero value `""` is no valid state, so a forgotten
  value or a stub shows up; readable test output; same as the other enums. An
  `iota` enum would make `0` a valid state unless it starts at `iota + 1`.
- It is a state, not a label: `web` picks *Récoltée* / *Terminée*,
  *À récolter* / *À terminer* from `Purpose`.
- `func (p Planting) InPlotBy(day Date) bool`: occupancy started on or before
  `day`; false while the seedling waits in its tray. It handles the zero
  `OccupiedFrom()` once, so callers never compare `Date{}` (year 1, before
  every real day). It ignores `EndedAt`. `web` reuses it for the *À venir*
  section. A method rather than `inPlotBy(occupiedFrom, day Date)`: with two
  dates, a caller could pass `StartDate()` and the compiler would accept it.

Tests: `TestPlanting_Status` 12 cases, `TestPlanting_InPlotBy` 4 cases.
Sabotages: `StartDate()` for occupancy (tray tomato fails), `due` checked
before `ended` (case 8, today after the expected end, fails), `due` on the day
itself (case 6 fails), the `IsZero` check dropped (tray tomato fails).

## Next: Area

Rule (`step-1-spec.md`, *Area*): `length_cm × width_cm`, displayed in m² with
one decimal. Both dimensions are optional (> 0, ≤ 100 m); the area shows only
when both are set.

- New type `Plot` in `plot.go`, with only the fields the rule needs for now:
  `LengthCm` and `WidthCm` (`int`, `0` = absent, as for maturity days).
- `func (p Plot) AreaCm2() int`: the area in **cm²**, `0` when a dimension is
  missing. The domain keeps an exact integer; `web` turns 36 000 cm² into
  `3,6 m²` (÷ 10 000, one decimal, comma). No `float64` in the domain: no
  rounding question in the tests, `==` works.

### `TestPlot_AreaCm2` plan

1. **(model)** Planche de 3,00 × 1,20 m (300 × 120 cm) → 36 000 cm² (3,6 m²).
2. Longueur seule → 0.
3. Largeur seule → 0.
4. Aucune dimension → 0.
5. Le maximum, 100 × 100 m → 100 000 000 cm².
6. Le minimum, 1 × 1 cm → 1 cm².

**Done when:** `go test ./internal/potager/ -run TestPlot_AreaCm2 -v` →
**6 cases PASS**. Sabotages: `length × length` (case 1 fails, the plot is not
square); converting to m² with an integer division by 10 000 (case 1 gives 3,
case 6 gives 0).

## Before closing T3

- [ ] Doc comments of `Purpose`, `Demand`, `Effect` (`neteffect.go`): domain
      words, none documented.
- [ ] Doc comment of `DateFromTime` (`date.go`): "returns a Date in
      Europe/Paris timezone" contradicts the `Date` decision (no zone). It
      returns the calendar day in Paris of that instant.
- [ ] A story test (or `Example`) following one tomato through a season: tray
      in March → planned, planted out in May → growing, end of July → due,
      harvested → ended.
- [ ] Move what outlives T3 (the "With the store" list below) into the tasks
      it belongs to, then delete this file and open the PR.

### Closed on 2026-11-05

- *`transplant` species, no days, no transplant date → `needs transplant`*:
  kept. The screen says « Indiquez la date de plantation ou saisissez une fin
  prévue », so it stays useful without days.
- *`Season(endedAt Date)` as a function*: kept. It takes a date, not a
  planting, so the "methods on `Planting`" rule does not apply; `web` groups
  with `Season(p.EndedAt)`.
- *`TestNetEffectPanics` case for `""`*: dropped. The `default` covers it, and
  `TestPlanting_ExpectedEndPanics` already documents an empty purpose.
- *Rename `DateFromInt`*: dropped. The name is clear and used in every test.
- Done: `ExpectedEndSource` and `Planting` doc comments, `InPlotBy(day Date)`,
  panic prefix `Planting.ExpectedEnd:` (the rule became a method).

## With the store (T4 or later)

- [ ] Migration for `plantings`: `CHECK (transplanted_at >= sown_at)`, the
      invariant `StartDate` relies on (see `docs/design.md`).
- [ ] `ParsePurpose`, `ParseMaturityFrom`, `ParseDemand`
      (`func(string) (X, error)`, like `strconv.ParseBool`): the boundary where
      a column becomes an enum. An `error` there, not a panic: bad data at the
      boundary is expected, not a bug.
- [ ] Optional: the `exhaustive` linter (golangci-lint) to catch a `switch` on
      an enum that forgets a constant.
- [ ] Before deploying (T11): `import _ "time/tzdata"`, because a minimal Docker
      image has no tzdata.
