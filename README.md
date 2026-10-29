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

- `task init` creates `.env` from `.env.example` (an existing `.env` is never overwritten) and installs the Git hooks
- `task up` starts the Postgres service and waits until it is healthy
- `task run` starts the server; in another terminal, `curl localhost:8081/health` answers `ok` once the database is reachable
- `task --list` shows every command; `task check` runs what CI runs

### Git hooks

[Lefthook](https://lefthook.dev/) runs two hooks, configured in [`lefthook.yaml`](lefthook.yaml)
and installed by `task init` (or `task hooks` alone):

- `pre-commit` refuses staged Go files that `gofmt` would change
- `pre-push` runs `task lint`, then `task test`, and stops at the first failure

The hooks are a convenience, not a gate: `git commit --no-verify` and
`git push --no-verify` skip them. CI runs the same checks and has the last word.

Two things to know:

- The installed hook script embeds the path of the pinned Lefthook in the Go
  build cache. After `go clean -cache` or a Lefthook version bump, run
  `task hooks` again.
- If a `lefthook` is already in your `PATH`, the hook uses it instead of the
  pinned version (Lefthook is pinned in `tools.mod` like the other tools).

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
