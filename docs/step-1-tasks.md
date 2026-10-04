# Step 1 — task breakdown

Tasks are vertical slices: each one goes from the database to the screen and
leaves something that works. The owner writes the code; see
[CLAUDE.md](../CLAUDE.md). Acceptance numbers refer to
[step-1-spec.md](step-1-spec.md). The `scaffold` branch holds an earlier
implementation to compare against once a task is done, not before.

| # | Task | Done when | Go ground covered |
| --- | --- | --- | --- |
| 1 | **Project skeleton.** Module, `Taskfile`, Postgres in Compose, `/healthz`, graceful shutdown, JSON logs. | `task run` serves `/healthz`; `go vet` clean. | `context`, `os/signal`, `net/http`, `slog` |
| 2 | **Schema and catalogue.** Migrations for the model in `design.md`, seeded families, species and varieties. | `task migrate` on an empty database yields the full catalogue. | goose, Postgres types and constraints |
| 3 | **Domain rules.** `Date`, `ExpectedEnd`, `Season`, `Status` (four states), `NetEffect`. Test-first. | Every rule in *Derived values* has a table-driven test, including the no-fallback case for transplanted species. | table-driven tests, optional values, pure functions |
| 4 | **Plots: list and create.** `GET /`, `GET /plots/new`, `POST /plots`. First sqlc queries and first templ components. | AC 1 (create, list), form validation, redirect after POST. | sqlc, templ, form parsing, post-redirect-get |
| 5 | **Plots: detail, edit, archive, restore.** | AC 1 and 2, including the archive confirmation counting plantings in progress. | routing with path values, consumer-side interfaces |
| 6 | **Planting form, food crop.** Species, then variety through HTMX, with a no-JS fallback. Validation messages. | AC 3, 7, 14; entered values survive a refusal. | HTMX, partials rendered standalone and inline, transactions |
| 7 | **Derived display.** Expected end, status, the three sections of the plot detail, history by season. Injectable clock. | AC 5, 6, 8, 12. | `time`, injected dependencies, `httptest` |
| 8 | **Cover crops.** Purpose switch, species filter, seed rate, *terminé* wording. | AC 4 and 6. | conditional templates |
| 9 | **Edit and delete a planting.** Confirmation screen. | AC 13. | |
| 10 | **Plot list without N+1; mobile pass.** | AC 9 and 10. | query shape, CSS only |
| 11 | **Deploy.** Caddy with basic auth, `migrate` service, CI quality gates on `main`. | AC 11, 15, 16, then 17: one real sowing recorded from a phone. | Docker, Compose, GitHub Actions |

## Order and dependencies

Tasks 2 and 3 only need task 1 and can be done in either order. 4 needs 2. 5 needs 4. 6 needs 3 and 5.
7 needs 3 and 6. 8 needs 6. 9 needs 6. 10 needs 7. 11 needs everything else, though
the Compose and CI parts can start as soon as task 1 exists.

## Design iterations

After task 5 and after task 8, stop for a short review of the spec against what
was learnt, before going on. Anything that changes a decision gets an ADR or an
edit to `design.md`.
