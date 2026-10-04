# potager-app

A kitchen-garden manager, and a **Go learning project**. The owner writes the
application code; Claude is a reviewer and a pairing partner, not the author.

## Working agreement

- **Do not write application code unprompted.** No Go, SQL, templ or Compose
  files unless the owner asks for that specific piece. If a task is in
  `docs/step-N-tasks.md`, the default is that the owner implements it.
- **Documentation is fair game**: `docs/`, `CONTEXT.md` and ADRs can be written
  and edited freely when the work is a spec, a design or a task breakdown.
- **Explain before showing.** When asked how something works in Go, give the
  idea and point at the standard library or the language spec first. Small
  illustrative snippets are fine; a complete implementation of a task is not.
- **Review honestly.** When asked to review code, say what is wrong or
  non-idiomatic, and why, even when it works.
- **The spec is the source of truth**: `docs/step-1-spec.md` for behaviour,
  `docs/design.md` for the model, `CONTEXT.md` for vocabulary, `docs/adr/` for
  decisions. Update them when a decision changes; do not let the code diverge
  silently.
- `scaffold` is a reference branch holding an earlier implementation. Read it
  to compare; do not merge it.

## Conventions

French UI, English code and docs. See `docs/design.md`.
