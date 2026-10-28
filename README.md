# potager-app

A kitchen-garden manager: which crops are growing in which plot, when they went
in, when they come out, and whether the succession makes agronomic sense.

The interesting part is not the CRUD. It is that the app knows **moutarde is a
brassica** — so sowing it as a cover crop blocks cabbage on that bed for four
years — and that it knows **a cover crop returns the nitrogen it took** while a
harvested pea crop carries its own away in the seed.

The interface is in French; the code and documentation are in English. This is
also a Go learning project: the owner writes the code (see [CLAUDE.md](CLAUDE.md)).

## Stack

Go 1.26 (stdlib `net/http`, no framework) · templ + HTMX · Postgres via sqlc ·
goose migrations · Docker Compose + Caddy.

## Setup

### Prerequisites

- Go 1.26
- Docker Compose (plugin compose v2)
- [Task](https://taskfile.dev/installation/)

goose, sqlc, templ and staticcheck need no installation. Their versions are
pinned in [`tools.mod`](tools.mod) and `task` runs them with `go tool`, so the
first task that needs one downloads it (network required, a little slower once).

#### Optional

- `postgresql-client` for troubleshooting

### Getting started

- `cp .env.example .env`
- `task up` starts the Postgres service and waits until it is healthy
- `task run` starts the server; in another terminal, `curl localhost:8081/health` answers `ok` once the database is reachable
- `task --list` shows every command; `task check` runs what CI runs

## Documents

- [CONTEXT.md](CONTEXT.md) — the vocabulary
- [CONTRIBUTING.md](CONTRIBUTING.md) — branches, pull requests, commits
- [docs/design.md](docs/design.md) — principles, model, derived logic, sequence
- [docs/step-1-spec.md](docs/step-1-spec.md) — functional specification of step 1
- [docs/step-1-tasks.md](docs/step-1-tasks.md) — task breakdown of step 1
- [docs/adr/](docs/adr/) — decisions that are hard to reverse and surprising

## Status

Specification only. Step 1 of 6 is specified and broken into tasks; no
application code on `main` yet. The branch `scaffold` holds an earlier
implementation, kept as a reference.
