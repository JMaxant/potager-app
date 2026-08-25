# potager-app — design

A kitchen-garden manager: which crops are growing in which plot, when they went
in, when they come out, and whether the succession makes agronomic sense.

## Principles

- **One step at a time.** Each step deploys and gets used for real before the
  next one starts.
- **Advisory, never blocking.** Rotation, nitrogen and overlap checks warn; no
  rule ever refuses a save. There is the plan, and then there is the weather.
- **Dates are facts, everything else is derived.** No stored status, no stored
  season, no stored soil effect. A derived value can be corrected retroactively
  by fixing the catalogue; a stored one silently goes stale.
- **French UI, English code.** Domain nouns are English (`Plot`, `Species`,
  `Variety`, `Planting`). French appears only in templ templates, in a single
  `labels.go`, and in catalogue data. Docs and comments are English.

## Stack

| Concern | Choice |
| --- | --- |
| Backend | Go 1.25, stdlib `net/http`, dependencies wired by hand in `main.go`. No framework, no DI container. |
| Templates | `templ` + HTMX, used minimally |
| Data | Postgres, `sqlc` over hand-written SQL, `goose` migrations. No ORM. |
| Operations | `log/slog`, `context` throughout, graceful shutdown, env config, wrapped errors, `/healthz` |
| Deploy | Docker Compose + Caddy on a small VPS, by hand until it is boring |
| CI | GitHub Actions: vet, staticcheck, tests against a real Postgres, `sqlc diff`, `templ generate --check` |

`sqlc diff` and `templ generate --check` are in CI on purpose: generated code that
has drifted from its source fails the build instead of rotting quietly.

### Layout

```
cmd/potager/      entrypoint
internal/potager/ domain types and agronomic rules — imports no net/http, no database/sql
internal/store/   sqlc output and queries
internal/web/     handlers, templ components, middleware
migrations/       goose migrations, including the seeded catalogue
```

Interfaces are declared where they are consumed (`web.Store`), not beside the
implementation. That is the Go convention and the one most Java-shaped Go
codebases get backwards.

## Model

```
families ──< species ──< varieties ──< planting_varieties >── plantings >── plots
```

- **`families`** — `rotation_interval_years` (years before this family should
  return to a plot) and `nitrogen_fixing`. Fixation is a family-level botanical
  trait: every Fabaceae nodulates, nothing else in the catalogue does.
- **`species`** — `nitrogen_demand` (heavy/moderate/light), `maturity_from`
  (sowing/transplant), `spacing_cm`, `suitable_as_cover`,
  `termination_lead_days`.
- **`varieties`** — `days_to_maturity`, null for species only ever sown as green
  manure.
- **`plots`** — name, type, dimensions, position, `archived_at`. No coordinates:
  a map can be derived from dimensions and ordering if one is ever built.
- **`plantings`** — one plot, one `purpose`, one set of dates. Occupancy is a
  generated `daterange`.
- **`planting_varieties`** — components of a planting, with `quantity` and
  `seed_rate_g_m2`.

### Why a planting holds many varieties

A cover crop mix — trèfle + seigle + radis fourrager — is **one sowing**: one
date, one incorporation, one decision. Modelling it as three overlapping
plantings would mean three rows that must agree forever, and three passes through
the form on a phone in the garden.

Registering the mix as its own species would be worse: a "Mélange hivernal" has
no botanical family, so the radis fourrager in it would silently vanish from
brassica rotation checks. Joining through `planting_varieties` to reach families
costs one join and keeps that correct.

### Why occupancy is not an exclusion constraint

Undersowing (trèfle beneath a standing crop) and intercropping (radis between
carrot rows) both put two plantings in one plot at one time, deliberately. The
GiST index is for query speed only. Overlap produces an advisory notice, never a
rejection.

Note that Postgres canonicalises discrete `daterange` values to `[)`, so
`upper(occupancy)` is one day past the recorded end date.

### Purpose lives on the planting, not the species

Féverole, moutarde and sarrasin are dual-use: you eat them or you turn them in.
Intent belongs where intent is formed. `species.suitable_as_cover` exists only to
sort the picker and is never the source of truth.

## Derived logic

**Net soil effect** — never stored, a pure function of intent, appetite and
botany:

| purpose | nitrogen-fixing | demand | effect |
| --- | --- | --- | --- |
| cover | yes | any | enriching |
| cover | no | any | neutral |
| food | yes | any | neutral |
| food | no | heavy | depleting |
| food | no | moderate/light | slightly depleting |

The two rows that matter are the ones that a single "nitrogen" enum would have
got wrong. A **cover crop is terminated before it fruits**, so everything it took
up returns to the bed on incorporation — moutarde is a hungry plant with a
neutral net effect. A **legume grown for food exports its nitrogen in the seed**,
so a harvested pea crop enriches far less than its nodules suggest.

History is judged on net effect; a candidate crop is judged on demand.

**Expected end** — a stored `expected_end_at` always wins. Otherwise a food
planting is `start + MAX(days_to_maturity)` across its components, counted from
the date `maturity_from` names. A cover crop gets nothing computed: its
termination is driven **backwards** from when the bed is next needed
(`next planting start − MAX(termination_lead_days)`), because that is how the
decision is actually made. Seigle is allelopathic and decomposes slowly, so in a
mix it sets the lead time for the whole sowing.

**Season** — the year a planting leaves the ground. Garlic sown in October 2025
and lifted in July 2026 is one season, 2026, not two.

**Status** — derived from which dates are filled in.

**Rotation neutrality** — a family with no food-purpose species (phacélie,
sarrasin) is free to sow anywhere. Derived, not flagged, so it stays true as the
catalogue grows.

## Sequence

1. **Schema, seeded catalogue, CRUD.** Plot CRUD, planting CRUD (one variety row
   in the form), expected harvest date. Four screens: plot list, plot detail,
   plot form, planting form. Deployed by hand and used for real.
2. **Rotation.** Family interval warnings, nitrogen succession, overlap notices,
   the multi-variety picker, the backward-derived termination suggestion.
3. **Auth.** Users, session cookies, argon2id. `user_id` already exists on
   plantings; plots need one too.
4. **Season overview and archiving UI.**
5. **Catalogue CRUD and reusable cover-crop mixes.**
6. Plot map, spacing advice, whatever the garden asks for.

## Deliberately not doing

Offline/PWA, i18n, a JSON API, an ORM, a DI container, ports-and-adapters,
stored coordinates, planned-vs-actual date pairs, per-variety spacing overrides,
automated deploys before the first manual one, a table-driven rule engine.

## Known constraints

- Until step 5, adding a species means editing a migration. That constraint is
  what keeps step 1 short. If it chafes within a week, pull step 5 forward.
- The `daterange` occupancy column is load-bearing: rotation, the season overview
  and termination suggestions all read it.
