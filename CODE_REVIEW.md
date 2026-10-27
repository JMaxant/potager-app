# Code Review Report

**Scope**: PR #19 "Feat: go app skeleton" (closes #5), branch `feat/go_skeleton_5` against `main`, 9 commits, 10 files, +507/-2.
**Stack**: Go 1.26 (`go.mod`), stdlib `net/http`, `pgx/v5`. Conventions from `references/golang.md`.
**Workload**: a request-serving HTTP service behind Caddy, single instance, one operator.
**Checks run**: `gofmt -l` (clean), `go vet ./...` (clean), `go test ./...` (14 passed).

## Summary
- **Status**: Approved, with minor comments (none blocks the merge)
- **Overall severity**: Low

## Strengths
- Graceful shutdown is correct: `signal.NotifyContext`, then `Shutdown` on a **fresh** `context.Background()` with its own 10 s budget, because the signal context is already cancelled at that point. The `errs` channel is buffered (size 1), so the serving goroutine cannot leak.
- `ReadHeaderTimeout` is set (Slowloris), and the database connection attempt is bounded by a 10 s context.
- `web.Store` is a one-method interface declared on the consumer side, as the project's own design doc demands. `*store.Store` satisfies it implicitly.
- Errors are wrapped with `%w` and keep their chain; `errors.Join` reports every configuration problem at once.
- `loadConfig` takes a `getenv` function instead of reading `os.Getenv`: pure and trivially testable.
- No secret is logged: only `listen_addr` is. The health endpoint answers `ok` / `health check failed` and leaks nothing about the database.
- Tests follow the table-driven convention with `t.Run`, use a fake store, and `TestServer_Handler` pins the 405 / 404 behaviour of the router.
- The configuration section in `docs/design.md` matches the code (strict, once at start-up, no default for `LISTEN_ADDR`).

## Observations and Improvements

### `internal/web/server.go`: `statusRecorder`
- **Issue**: the wrapper embeds `http.ResponseWriter` but hides the optional interfaces of the real writer (`http.Flusher`, `http.Hijacker`, ...), and it has no `Unwrap() http.ResponseWriter` method. It also lives in `server.go` while its only user is `middleware.go`.
- **Impact**: Maintainability. Nothing breaks today. It will the day a handler streams or flushes (HTMX partials, SSE), because `http.NewResponseController(w).Flush()` cannot reach the real writer through the wrapper.
- **Suggestion**: add `func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }` and move the type next to `logRequests`.

### `cmd/potager/main.go`: `http.Server` settings
- **Issue**: only `ReadHeaderTimeout` is set. `IdleTimeout` and `WriteTimeout` are left at zero (no limit).
- **Impact**: Security / Scalability, low. Caddy sits in front and absorbs most of it, but idle keep-alive connections are never reclaimed.
- **Suggestion**: set an `IdleTimeout` (60 s is a common value) now; choose `WriteTimeout` once real handlers exist.

### `cmd/potager/main.go`: `ListenAndServe` result
- **Issue**: after `Shutdown`, `ListenAndServe` always returns `http.ErrServerClosed`, which the goroutine wraps and pushes to `errs` like a real failure.
- **Impact**: Maintainability. Harmless today: the `select` has already left and the channel is buffered. It becomes a trap if someone later reads `errs` after the shutdown.
- **Suggestion**: `if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { ... }`.

### `cmd/potager/main.go`: logger wiring
- **Issue**: the logger is both injected (`web.NewServer(logger, db)`) and installed as the global default (`slog.SetDefault`).
- **Impact**: Maintainability. `docs/design.md` says dependencies are wired by hand, no global state; `SetDefault` is a quiet exception.
- **Suggestion**: drop `slog.SetDefault` unless something needs it (a library logging through the default logger), or note the exception in the design doc.

### `internal/web/middleware.go`: no test
- **Issue**: `logRequests` has no test. Default status 200 when the handler never calls `WriteHeader`, the first-`WriteHeader`-wins rule, and the log fields are all unverified.
- **Impact**: Maintainability. This is the only code in the PR that changes behaviour without a safety net.
- **Suggestion**: a table-driven test with a `slog` handler writing to a `bytes.Buffer`, and three cases: no explicit status, explicit 404, double `WriteHeader`. Suitable as a follow-up exercise.

### `cmd/potager/main_test.go`: assertions
- **Issue**: error cases only check that the message contains the variable name, and never that the *other* variables are absent.
- **Impact**: Maintainability. Already tracked in issue #20; no action in this PR.

### `docs/design.md`: `LOG_LEVEL` values
- **Issue**: the table says `debug`, `info`, `warn`, `error`. `slog.Level.UnmarshalText` also accepts upper case and offsets (`INFO+2`, `debug-1`).
- **Impact**: Documentation accuracy, cosmetic.
- **Suggestion**: write "case-insensitive; offsets such as `warn+1` are also accepted", or leave it if the strict list is the intended contract.

### Style nits
- A blank line separates a call from its `if err != nil` in `store.go`, and in `main.go` (lines around `store.New`, `Shutdown`, `parseLogLevel`). Idiomatic Go keeps them adjacent.
- `fmt.Errorf("config:\n %w", err)` embeds a newline and a space in the message; fine for a multi-line `errors.Join`, but a comment would explain the intent.
- Loop variable is `tt` in one test and `test` in the other; pick one.
- `duration` is logged as a `time.Duration`, which the JSON handler renders in nanoseconds. Log `.Milliseconds()` if the logs are read by humans.

## Questions / Clarifications
- Should `/health` requests appear in the request log? Caddy or a Compose healthcheck will call it every few seconds and drown real traffic at `info` level.
- Is `IdleTimeout` deliberately left unset until deployment (issue #16), or just not thought of yet?
- `internal/store` has no test (it needs a database). Is that intended to wait for the CI task (#4)?
