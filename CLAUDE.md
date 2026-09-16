# Ginko

Alpha-stage Go package for building lightweight desktop apps in pure Go: a
local HTTP server renders HTML into a native WebView window, with HTMX
handling the Ajax interactions. Single package, `github.com/serge-hulne/ginko`.

## Layout

- `lib.go` — the entire library (~120 lines): server bootstrap, WebView
  launch, and two HTML helpers (`ButtonHTMX`, `HeadHTMX`).
- `assets/w3.css`, `assets/htmx.min.js` — vendored third-party assets,
  embedded into the binary via `go:embed` and served locally (see below).
- `README.md` — usage docs, containing two full example apps (a counter and
  a to-do list) inside fenced ` ```go ` blocks. These are the canonical
  reference for how consumers use the library, so keep them compiling.
- `TODO.md` — informal roadmap.

## Public API (lib.go)

- `Response`, `Request`, `ActionMap` — thin aliases over `net/http` types.
- `Run_app(title, port string, action ActionMap) error` — starts the local
  HTTP server on `127.0.0.1:<port>` and opens a WebView window pointed at
  it. Returns an error immediately if the port can't be bound (checked via
  a synchronous `net.Listen` before the WebView opens); the WebView then
  runs on the calling goroutine until the window is closed.
- `Display(w Response, n gomponents.Node)` — renders a gomponents `Node` to
  the response writer.
- `ButtonHTMX(action, target, id, text string) gomponents.Node` — an
  HTMX-wired `<button>`.
- `HeadHTMX() gomponents.Node` — a `<head>` with the vendored w3.css and
  htmx assets wired in.

## History of significant changes

### Security/correctness pass
The original code (`julvo/htmlgo`-based) had several issues, all fixed:
- **Stored XSS** in `ButtonHTMX` and in the README's to-do example — both
  built HTML via `fmt.Sprintf` on unescaped user input. Fixed by escaping
  (later: by switching to gomponents, which escapes by default).
- **Server bound to all interfaces** (`http.ListenAndServe(":8090", nil)`)
  instead of loopback-only, despite the WebView only ever talking to
  `127.0.0.1`. Fixed: server now binds `127.0.0.1:<port>`.
- **`Run_app`'s `port` argument was silently ignored** — the server always
  listened on a hardcoded `:8090` while the WebView navigated to whatever
  port was passed in, so any non-8090 port would fail to load. Fixed by
  threading `port` through to the listener.
- **No CSRF protection** on state-changing endpoints. Fixed by
  `csrfGuard` in `lib.go`, which rejects non-GET/HEAD requests whose
  Origin/Referer doesn't match the app's own loopback origin.
- **No server timeouts** — fixed with an explicit `http.Server` carrying
  read/write/idle timeouts.
- **Unpinned CDN assets** (w3.css floating version, htmx with no SRI) —
  resolved by vendoring both files locally (see below) rather than pinning,
  since w3.css has no versioned/hashed release to pin against anyway.
- **Data race** on the to-do example's shared `todoList` slice across
  concurrent HTTP handlers — fixed with a `sync.Mutex`.

### Templating migration: htmlgo → gomponents
Swapped `github.com/julvo/htmlgo` for `maragu.dev/gomponents`
(https://github.com/maragudk/gomponents). Reasons/effects:
- gomponents escapes `Text()`/`Attr()` values by default, which is what
  motivated dropping the manual `html.EscapeString` calls added during the
  security pass.
- `ButtonHTMX`/`HeadHTMX` now return `gomponents.Node` instead of
  `htmlgo.HTML`; `Display` takes a `Node` and calls `.Render`.
- Consumers dot-import `maragu.dev/gomponents/html` (elements/attributes)
  and alias `maragu.dev/gomponents` as `g` (core: `Text`, `Attr`, `Group`,
  `Node`) — see the README examples for the exact pattern.
- Fixed a latent bug surfaced during the rewrite: the to-do example had
  two nested elements both with `id="todo-list"` (invalid HTML, ambiguous
  `hx-target`). Now there's exactly one.

### Local asset vendoring
`w3.css` (w3schools, v4.15) and `htmx.min.js` (v1.9.8, integrity-verified
against unpkg's published sha256 before vendoring) now live in `assets/`
and are embedded into the binary with `go:embed`, served at
`/ginko-assets/w3.css` and `/ginko-assets/htmx.min.js` by `registerAssets`
in `lib.go`. This removes the runtime CDN dependency entirely — matches
the README's "does not require JS or CSS [from outside]" pitch and the
TODO item to vendor these assets, now checked off.

### `Run_app` returns an error instead of `log.Fatal`
Previously a failed `ListenAndServe` (e.g. port already in use) called
`log.Fatal`, killing the host process unconditionally — harsh for a
library function. Now `Run_app(title, port string, action ActionMap)
error`: it binds the listener synchronously with `net.Listen` before
opening the WebView, so a bind failure is returned to the caller before
any window appears; once bound, the server is served asynchronously via
`srv.Serve(ln)` and the WebView blocks on the calling goroutine as
before. Both README examples updated to check the returned error.

## Working notes

- Both README code examples are meant to compile as-is against the current
  `lib.go`. When changing the public API, update both examples and
  smoke-test them (extract the fenced block, build against the local
  module via a `replace` directive in a scratch `go.mod`) rather than just
  eyeballing them.
- `go.mod` requires a C/C++ toolchain (cgo) because of `webview_go`; expect
  harmless `-Wdeprecated-literal-operator` warnings from that dependency on
  build — unrelated to this package's own code.
- `lib_test.go` covers `csrfGuard` (allow/deny cases across Origin/Referer
  combinations) and `registerAssets` (content-type and body of the two
  vendored assets) via `httptest`. Run with `go test ./...`.
