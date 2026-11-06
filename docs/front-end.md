# potager-app — front-end

Server-rendered HTML from templ components, plain CSS, HTMX on two
interactions. Nothing is fetched from a third-party host at runtime, and nothing
is built beyond `templ generate`.

## Static files

```
internal/web/static/
  htmx.min.js         HTMX 4, vendored ([ADR 0009](adr/0009-htmx-is-vendored.md))
  styles/
    main.css          the only file carrying @import
    base/             reset.css, tokens.css, elements.css
    components/       one file per styled partial: field.css, badge.css, …
    layout/           shell.css (header, nav, flash messages)
    utils.css
```

**Everything under `static/` is public, nothing else is.** The directory is
embedded with `//go:embed static` and served under `/static/` by
`http.FileServerFS`. Stylesheets do not sit beside the `.templ` files: a pattern
reaching into `templates/` would one day be widened and publish the template
sources and generated Go.

**No asset pipeline.** Go has no equivalent of a library system: `go:embed`
copies files as they are, so they are served exactly as written, unminified.
Caddy compresses them (#16).

## Stylesheet

- `main.css` declares the cascade layers first, then imports every other file
  into one of them:
  `@layer reset, tokens, base, components, layout, utils;` then
  `@import url("components/field.css") layer(components);` and so on. Layer
  order, not file order or selector specificity, decides which rule wins, so
  `utils` beats a component wherever it is imported. An `@import` after any
  other rule is ignored by the browser, silently.
- The browser follows the imports itself: there is no bundler to inline them.
  That costs one extra round trip on first load and is accepted.
- **Tokens** (colours, spacing, type scale, `--size-touch-target: 44px`) live as
  custom properties on `:root` in `base/tokens.css`. To change a colour, that is
  the file. A component may define a local custom property for geometry nobody
  else uses.
- Native nesting, no preprocessor. Phone first: base rules are the single
  column, `@media (width >= 640px)` adds the wider layout.

## Naming: BEM, ABEM-style states, no prefixes

| Kind | Form | Question it answers | Example |
| --- | --- | --- | --- |
| Block | `.name` | which component? | `.planting-row`, `.field` |
| Element | `.block__part` | which part of it? | `.planting-row__name`, `.field__error` |
| Modifier | `.block--variant` | which version of it? (fixed at render) | `.badge--finished`, `.field--inline` |
| State | `.-state` on a block or element | what is it doing now? (changes at runtime) | `.planting-row.-removing` |

A state that ARIA already expresses is styled through the attribute, not a
class: `.field [aria-invalid="true"]`, `[aria-current]`. The attribute is needed
for assistive technology anyway, and one source of truth cannot disagree with
itself. Class names are English, like the rest of the code.

## Components

The split is the one in [`step-1-spec.md`](step-1-spec.md#components) (*Components*): one `.templ` file
per screen, shared partials for what repeats. Three conventions on top:

1. **One CSS file per styled partial**, same name, same BEM block:
   `field.templ` → `components/field.css` → `.field`. Screens have no
   stylesheet of their own; they assemble components. `layout.templ` maps to
   `layout/shell.css`.
2. **A partial is an exported templ component, rendered two ways.** Inline from
   a page (`@VarietyOptions(varieties)` in `planting_form.templ`), and
   standalone from a handler (`templates.VarietyOptions(varieties).Render(ctx, w)`,
   no layout). Same component, same HTML: that is how "no fragment is reachable
   only through HTMX" holds.
3. **A route always returns the same shape.** Screen routes return full pages,
   fragment routes (`/varieties/options`) return fragments. No handler branches
   on the `HX-Request` header, so no response needs `Vary: HX-Request` and the
   cache below never mixes a page and a fragment.

## Caching

Files in an `embed.FS` have a zero modification time, so `http.FileServer`
sends no `Last-Modified` and no `ETag`, and the browser re-downloads every file
on every page. A middleware in front of the file server fixes that:

- at start-up, walk the embedded files and compute a content hash for each
  (`io/fs.WalkDir`, `crypto/sha256`);
- on each request, set `ETag` to that hash and `Cache-Control: no-cache`, then
  hand over to the file server.

`no-cache` means "revalidate before use", not "do not store": the browser keeps
its copy and asks; `http.FileServer` already answers `304 Not Modified` when
the request's `If-None-Match` matches an `ETag` set before it runs. Hashed file
names with `immutable` were rejected: the `@import` URLs in `main.css` are fixed
strings and could not follow renamed files without rewriting the CSS. Check
in #16 that Caddy's compression keeps the `ETag` usable.

## Development loop

`task dev` runs templ's watch mode, already pinned in `tools.mod`: no other
tool. It regenerates on change, restarts `go run ./cmd/potager` (`--cmd`), and
serves a proxy on port 7331 that injects a reload script, so the browser
refreshes itself (`--proxy`, pointed at `LISTEN_ADDR`). `--watch-pattern`
extends the default with `.css`: stylesheets are embedded, so a CSS change
needs a rebuild, one or two seconds.

Watch mode writes development-only Go that reads template text from
`*_templ.txt` files. Those files are ignored by Git, and the generated Go must
be regenerated with `task generate` before a commit; `templ generate --check`
in `task check` refuses it otherwise.

If the rebuild on every CSS change becomes a nuisance: the file server takes an
`fs.FS`, so development can pass `os.DirFS("internal/web/static")` where
production passes the `embed.FS`. Not done until needed.

## Versions

The templ CLI (`tools.mod`) and the templ library (`go.mod`, added with the
first component) are chosen together ([ADR 0006](adr/0006-dev-tools-are-pinned-in-tools-mod.md)). HTMX is pinned by the file in
the repository and its version recorded in ADR 0009.
