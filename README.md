# potager-app

A kitchen-garden manager: which crops are growing in which plot, when they went
in, when they come out, and whether the succession makes agronomic sense.

The interesting part is not the CRUD. It is that the app knows **moutarde is a
brassica** — so sowing it as green manure blocks cabbage on that bed for four
years — and that it knows **a cover crop returns the nitrogen it took** while a
harvested pea crop carries its own away in the seed.

The interface is in French; the code and documentation are in English.

## Stack

Go 1.25 (stdlib `net/http`, no framework) · templ + HTMX · Postgres via sqlc ·
goose migrations · Docker Compose.

See [docs/design.md](docs/design.md) for the design and the reasoning behind it.

## Running it locally

Requires Go 1.25, Docker, and the following tools:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/a-h/templ/cmd/templ@latest
```

Then:

```sh
cp .env.example .env
task up        # start Postgres
task migrate   # apply migrations, including the seeded catalogue
task run       # http://localhost:8080
```

`task check` runs everything CI runs.

## Status

Step 1 of 6: schema, seeded catalogue and domain rules are in place; the screens
are not built yet. The sequence is in the design document.
