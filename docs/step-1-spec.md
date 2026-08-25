# Step 1 — functional specification

Step 1 makes the app usable for its core purpose: recording what is planted
where and when, and seeing when it is expected to come out. It is the smallest
slice worth deploying.

This spec is written against the model in [design.md](design.md), which is
already migrated and seeded. Nothing here changes the schema.

## Scope

**In:** plot CRUD, planting CRUD with a single variety, derived expected end
date, derived status, four screens, deployed by hand and used for one real
sowing.

**Out, and named so it stays out:** rotation and nitrogen warnings, overlap
notices, multi-variety mixes, the backward-derived termination suggestion, auth,
the season overview, catalogue editing, the plot map, spacing advice.

Two consequences of that list worth stating plainly:

- **The catalogue is read-only.** Adding a species means editing
  `migrations/00002_seed_catalogue.sql` and running a new migration. This is the
  constraint that keeps step 1 short.
- **A planting holds exactly one variety.** The `planting_varieties` table
  already supports several; the step 1 form writes one row. No schema change is
  needed when step 2 opens it up.

## Conventions

**URLs are English**, matching the code rather than the interface. French URLs
would mean percent-encoded accents in paths (`/variétés`) for no benefit. The
interface itself is French throughout.

**Dates** display as `15/03/2026` and are entered with `<input type="date">`.
**Dimensions** are entered in metres (`3,00`) and stored as centimetres. Decimal
separator is a comma in display, a period in form values.

**Forms are ordinary HTML.** POST, validate, redirect on success; re-render with
errors on failure. HTMX is used in exactly two places (below). Every screen works
with JavaScript disabled.

## Routes

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/` | Plot list |
| GET | `/plots/new` | Plot form, empty |
| POST | `/plots` | Create plot |
| GET | `/plots/{id}` | Plot detail |
| GET | `/plots/{id}/edit` | Plot form, populated |
| POST | `/plots/{id}` | Update plot |
| POST | `/plots/{id}/archive` | Archive plot |
| POST | `/plots/{id}/restore` | Restore plot |
| GET | `/plantings/new?plot_id={id}` | Planting form, empty |
| POST | `/plantings` | Create planting |
| GET | `/plantings/{id}/edit` | Planting form, populated |
| POST | `/plantings/{id}` | Update planting |
| POST | `/plantings/{id}/archive` | Archive planting *(HTMX)* |
| GET | `/varieties/options?species_id={id}` | Variety `<select>` options *(HTMX)* |
| GET | `/healthz` | Liveness + database check |

HTML forms cannot issue PUT or DELETE, so updates and archives are POST. There
are no hard deletes in step 1.

### On archiving

The design puts the archiving *UI* in step 4, but archiving without a way to see
or undo it is a trap: a mis-click makes a plot vanish permanently. Step 1
therefore includes `archive`, `restore`, and an **"Afficher les archivées"**
checkbox on the plot list. The richer history views stay in step 4. This is a
deliberate, small pull-forward.

## Screens

### 1. Plot list — `GET /`

Title **« Mes parcelles »**. Plots ordered by `position, id`, archived ones
hidden unless the checkbox is ticked.

Each plot shows:

- Name, and its type as a badge (*Pleine terre*, *Bac surélevé*, *Serre*,
  *Châssis*)
- Dimensions when both are set: `3,00 × 1,20 m`, followed by the area `3,6 m²`
- The plantings **currently occupying it** — `occupancy @> CURRENT_DATE` — as
  variety name, status, and expected end date
- A count of upcoming plantings when there are any: *« 2 cultures à venir »*

Actions: **« Nouvelle parcelle »**, and per plot **« Ajouter une culture »**.

Empty state: *« Aucune parcelle pour l'instant. »* with the create button as the
only call to action.

### 2. Plot detail — `GET /plots/{id}`

Heading is the plot name, with type, dimensions and area beneath it. Actions:
**« Modifier »**, **« Archiver »**, **« Ajouter une culture »**.

Three sections, each hidden when empty:

- **« Cultures en cours »** — `occupancy @> CURRENT_DATE`, ordered by start date
- **« À venir »** — `lower(occupancy) > CURRENT_DATE`, ordered by start date
- **« Historique »** — everything else, grouped by season year descending, the
  heading being *« Saison 2025 »*

Each planting row shows variety (with species in parentheses when they differ),
purpose badge for cover crops (**« Engrais vert »**), the relevant dates, status,
and an archive control.

### 3. Plot form — `GET /plots/new`, `GET /plots/{id}/edit`

| Field | Label | Type | Rules |
| --- | --- | --- | --- |
| `name` | Nom | text | required, 1–80 characters |
| `type` | Type | select | required, defaults to *Pleine terre* |
| `length_cm` | Longueur (m) | number, step 0,01 | optional, > 0, ≤ 100 m |
| `width_cm` | Largeur (m) | number, step 0,01 | optional, > 0, ≤ 100 m |
| `position` | Ordre d'affichage | number | required, defaults to `max(position) + 1` |

Names are not required to be unique — two beds may legitimately share a name —
but the form warns when the name already exists: *« Une parcelle porte déjà ce
nom. »*

### 4. Planting form — `GET /plantings/new?plot_id={id}`, `GET /plantings/{id}/edit`

| Field | Label | Type | Rules |
| --- | --- | --- | --- |
| `plot_id` | Parcelle | select | required, preselected from the query parameter |
| `purpose` | Type de culture | radio | required: *Culture vivrière* / *Engrais vert*, defaults to vivrière |
| `species_id` | Espèce | select, `<optgroup>` by family | required |
| `variety_id` | Variété | select, loaded by HTMX | required |
| `quantity` | Quantité | number | optional, > 0 |
| `seed_rate_g_m2` | Dose de semis (g/m²) | number, step 0,1 | optional, > 0, shown only for *Engrais vert* |
| `sown_at` | Semé le | date | see date rules |
| `transplanted_at` | Planté le | date | see date rules |
| `ended_at` | Récolté le / Enfoui le | date | optional |
| `expected_end_at` | Fin prévue | date | optional override |
| `maturity_days_override` | Durée réelle (jours) | number | optional, inside « Options avancées » |
| `notes` | Notes | textarea | optional |

**Selecting *Engrais vert* filters the species list** to those with
`suitable_as_cover`, and reveals the seed-rate field. Selecting a species
reloads the variety options.

**Labels change with purpose**: a food planting is *« Récolté le »* and *« Récolte
prévue »*; a cover crop is *« Enfoui le »* and *« Destruction prévue »*.

#### Validation

| Rule | Message |
| --- | --- |
| at least one of `sown_at` / `transplanted_at` | Indiquez au moins une date de semis ou de plantation. |
| `transplanted_at` ≥ `sown_at` | La date de plantation ne peut pas précéder le semis. |
| `ended_at` ≥ start date | La date de fin ne peut pas précéder le début de la culture. |
| `expected_end_at` ≥ start date | La fin prévue ne peut pas précéder le début de la culture. |
| `variety_id` belongs to `species_id` | Cette variété n'appartient pas à l'espèce choisie. |
| dates within 1900–2100 | Cette date semble erronée. |

Errors render beside their field, with a summary at the top of the form. Entered
values are always preserved on re-render.

## Derived values

Computed by `internal/potager`, never stored, displayed wherever a planting
appears.

**Expected end** — `ExpectedEnd`. A food planting with a maturity figure shows
*« Récolte prévue : 29/07/2026 »*. A cover crop shows nothing unless
`expected_end_at` was set by hand. A stored value always wins and is marked as
such: *« Récolte prévue : 01/07/2026 (saisie) »*.

**Status** — three states, derived from dates:

| Condition | Food label | Cover label |
| --- | --- | --- |
| `ended_at` set | Récoltée | Enfouie |
| start date in the future | Prévue | Prévue |
| otherwise | En cours | En cours |

**Season** — the year of `ended_at`, falling back to the expected end. Used only
to group the history section.

**Area** — `length_cm × width_cm`, displayed in m² with one decimal.

## Queries

To be written in `internal/store/queries/` and generated with sqlc.

`plots.sql` — `ListPlots` (with an `include_archived` flag), `GetPlot`,
`CreatePlot`, `UpdatePlot`, `ArchivePlot`, `RestorePlot`, `NextPlotPosition`,
`CountPlotsByName`.

`plantings.sql` — `ListCurrentPlantingsByPlot` (list screen, all plots at once,
to avoid N+1), `ListPlantingsForPlot` (detail screen, every planting with its
variety, species and family joined), `GetPlanting`, `CreatePlanting`,
`UpdatePlanting`, `ArchivePlanting`, `CountUpcomingByPlot`.

`catalogue.sql` — `ListSpecies` (ordered by family then name, with an optional
cover-only filter), `ListVarietiesBySpecies`, `GetVariety`.

`CreatePlanting` and `UpdatePlanting` write `plantings` and
`planting_varieties` in one transaction.

## Components

`internal/web/templates/` — `layout.templ` (shell, nav, flash messages),
`plots_list.templ`, `plot_detail.templ`, `plot_form.templ`,
`planting_form.templ`, and the shared partials `planting_row.templ`,
`variety_options.templ`, `field.templ` (label + control + error), `badge.templ`.

`variety_options.templ` and `planting_row.templ` are the two rendered both
standalone for HTMX and inline in a full page — the rule from the design: no
fragment is reachable only through HTMX.

## Non-functional

- Single column below 640 px; the sidebar becomes a top bar. Touch targets ≥ 44 px,
  because the phone is used standing in the garden.
- No JavaScript beyond HTMX itself. Both HTMX endpoints degrade: without it the
  variety list shows every variety, and archiving performs a normal POST and
  redirect.
- Every form control has a `<label for>`; errors are linked with
  `aria-describedby`; the error summary receives focus on failed submission.
- `lang="fr"` on `<html>`.
- Requests logged with method, path, status and duration; database errors return
  500 with a French message and the detail in the log, never on the page.

## Acceptance criteria

1. A plot can be created, edited, listed, archived and restored.
2. Archived plots are hidden by default and appear when the checkbox is ticked.
3. A planting can be created against a plot, choosing species then variety, with
   the variety list following the species without a page reload — and with one
   too, when JavaScript is off.
4. Choosing *Engrais vert* restricts the species list to cover-crop species and
   reveals the seed-rate field.
5. A tomato transplanted on 10/05/2026 with an 80-day variety shows *« Récolte
   prévue : 29/07/2026 »*; the same variety dated only by its sowing counts from
   the sowing instead.
6. A cover crop shows no expected end unless one was entered by hand.
7. A planting with neither sowing nor planting date is refused with the message
   above, and the entered values survive the refusal.
8. Garlic sown 12/10/2025 and lifted 04/07/2026 appears under **Saison 2026** in
   the history section.
9. The plot list issues a constant number of queries regardless of how many plots
   exist.
10. Every screen is usable at 375 px wide.
11. `task check` passes: vet, staticcheck, tests, `sqlc diff`, `templ generate --check`.
12. The app is deployed and one real sowing has been recorded through it.

## Tests

- **Domain** — done: `ExpectedEnd`, `StartDate`, `Season`, `NetEffect`.
- **Store** — one integration test per query against a real Postgres, including
  the transactional write of a planting plus its variety row, and the
  archived-plot filter.
- **Handlers** — `httptest` coverage of: create plot happy path and redirect;
  planting validation failure preserving input; the variety-options fragment;
  archive and restore. Assert on status codes and a few DOM markers, not on
  full-page snapshots.

## Definition of done

`task check` green, deployed to the VPS, and one real planting recorded from a
phone. Step 2 begins only after that.
