# HTMX 4 is vendored and embedded in the binary

`htmx.min.js` is copied into `internal/web/static/` and served from the binary's embedded filesystem, like the stylesheets. A CDN was rejected: the application must work offline and behind a firewall, with nothing fetched from a third-party host at runtime, and the file is small enough that caching across sites is no argument. Upgrading is a deliberate act: replace the file, record the version here, read the release notes.

The version is 4 (4.0.1 at the time of writing), although npm still tags 2.x as `latest` and publishes 4 under `next`. The project starts with no HTMX code to migrate, uses it in two places only, and HTMX 4 is what the owner wants to learn; the cost is fewer tutorials and answers online, most of which describe 2.x. Read the 2 → 4 migration notes before relying on an example written for 2.x.

Dropping HTMX entirely stays possible at any time: every screen works without JavaScript, and no fragment is reachable only through HTMX.
