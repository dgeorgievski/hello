# Repository Guidelines

## Project Overview

`hello` is a small Go HTTP service (`github.com/dgeorgievski/hello`) — a DevSecOps demo of a production-ish Go app: a gin-based server with health/version endpoints, a cobra CLI, Docker image, Helm chart, raw k8s manifests, GitHub Actions CI, and Backstage catalog entries.

- **Purpose**: runnable "hello world" that exercises the full build → container → Kubernetes pipeline (local `kind` + Helm/manifests).
- **Not** a library. The Go module is the application itself.

## Architecture & Data Flow

Entry point → CLI → HTTP server → handlers:

```
main.go (package main)
  └─ cmd.Execute()                      (cmd/root.go)
       └─ cobra: rootCmd
            ├─ "server"  (cmd/server.go) ─→ newRouter() → gin.Default()
            │    └─ routes: /hello, /version, /ping, /hey, /healthz, /readyz, /greet
            │       └─ middleware: requestLogger (stdlib log.Printf)
            │    └─ r.Run(":8080")   (hardcoded, error ignored)
            └─ "version" (cmd/version.go) ─→ version.GetPrettyVersion()
```

- `pkg/version/version.go` holds `VERSION` / `COMMIT` / `DATE` vars (default `1.1.8` / `none` / `none`), stamped at build time via `-ldflags -X` (see below).
- `cmd/server.go` builds the router in `newRouter()` — **lowercase, unexported**, but in-package tests call it directly (`cmd/server_test.go`). It is *not* exported despite its doc comment claiming otherwise.
- Handlers return 200 JSON unconditionally (`c.JSON(http.StatusOK, ...)`); no validation, no error responses, no graceful shutdown.

## Key Directories

| Path | Purpose |
|---|---|
| `main.go` | Entry point; just calls `cmd.Execute()` |
| `cmd/` | Cobra commands + gin router + handlers (`root.go`, `server.go`, `version.go`, `server_test.go`) |
| `pkg/version/` | Build-stamped version vars + `GetVersion()` / `GetPrettyVersion()` |
| `docker/` | `Dockerfile` (multi-stage: `golang:1.24` builder → scratch) |
| `deploy/helm/hello/` | Helm chart (deployment/hpa/ingress/service/serviceaccount/test pod) |
| `deploy/manifests/` | Raw (non-Helm) k8s manifests |
| `scripts/agit/` | Nushell git-workflow module (git-flow, commits, completions) |
| `scripts/githooks/` | `prepare-commit-msg` hook (Conventional Commits check) |
| `doc/` | `git-commands.md`, `git-hooks.md` |
| `.github/workflows/` | `go-build.yaml` (CI), `go-deploy.yml` (release) |
| `bin/` | Built binary (`bin/dgeorgievski-hello`, gitignored) |

## Development Commands

Requires Go 1.24+ (see `go.mod`: `go 1.24.0`; CI default `1.24`). `make` is optional — all targets are plain shell.

```sh
make run            # go run main.go (EXTRA_RUN_ARGS="..." to pass args)
make test           # go test ./... -coverprofile cover.out
make build          # stamped binary → bin/dgeorgievski-hello (TAG ?= dev)
make tidy           # go mod tidy
make fmt            # gofmt -l -s -w ./
make version-get    # echo current version from pkg/version/version.go
make get-env        # echo TAG/NAME/GIT_COMMIT/VERSION/etc.
```

Docker / release (see also README.md):

```sh
make docker-build           # docker build -f docker/Dockerfile -t dgeorgievski/dgeorgievski-hello:<VERSION> .
make kind-docker-upload     # docker-build + kind load
make tag                    # git tag $(VERSION) && push
```

Direct equivalents: `go run main.go server` (or `./bin/dgeorgievski-hello server`), `go run main.go version`, `go test ./...`.

**Version stamping** (`make build`): `-X <module>/pkg/version.VERSION=$(TAG)` (default `dev`), `.COMMIT=<git rev-list --abbrev-commit -1 HEAD>`, `.DATE=<epoch ms>`.

## Code Conventions & Common Patterns

- **Go 1.24**, gin v1.10 + cobra v1.9. Two commands: `server` (HTTP, `:8080`) and `version` (prints version). No config file, no env vars, no graceful shutdown, no structured logging (stdlib `log`).
- **Handlers are package-level funcs** in `cmd/server.go` returning `gin.HandlerFunc`-style JSON responses; the router is built in `newRouter()` and wired in `init()` via `rootCmd.AddCommand(...)`. Note: `greet` is defined at the *bottom* of `server.go`, after the routes that reference it.
- **Table-driven tests** (`cmd/server_test.go`) using `httptest.NewRecorder()` + in-process engine — no ports, no env, no filesystem. All handlers exercised through `newRouter()`.
- **Error handling is absent by design**: no 4xx/5xx paths, no validation, no `RunE` (all commands use `Run`), `r.Run(":8080")` error ignored, `os.Exit(1)` in `Execute()` without logging the error. Don't add error-path code paths without following the existing pattern — or intentionally change it.
- **`pkg/version/version.go`** is the single source of truth for version; `GetPrettyVersion()` prints only version+commit (drops DATE). Do not duplicate version strings elsewhere.

### Known gotchas (do not "fix" silently)

- **ldflags org mismatch**: Makefile `build`, Dockerfile and `.goreleaser.yaml` all stamp `-X github.com/altimetrik-digital-enablement-demo-hub/dgeorgievski-hello/pkg/version.VERSION=...`, but the module is `github.com/dgeorgievski/hello` — the `-X` likely silently no-ops. Don't trust stamped version output until the path matches `go.mod`.
- **Version drift**: `pkg/version` says 1.1.8, Helm `appVersion` says 1.16.0, raw manifests use image tag 0.0.1 — three different "versions".
- **Makefile `build` stamps VERSION from `TAG`** (default `dev`), while docker/tag targets use `pkg/version` (1.1.8) — the two can disagree.
- **Stub targets**: `version-set`, `docker-test`, `docker-push` echo "Coming soon"; `version-set`'s sed logic is commented out.
- **Stale Cobra boilerplate**: `Short`/`Long` strings, `toggle` flag (unused), and `Copyright © 2025 NAME HERE <EMAIL ADDRESS>` headers in all `cmd/*.go` are unedited scaffolding.
- **`scripts/agit/git-flow.nu` references undefined `_git_status`, `remote_branches`, `git_main_branch`**; `stat.nu` is 0 bytes. `git-merge` also has a typo flag `--fast-farward`. Treat the Nushell git-flow tooling as unreliable.
- **Catalog-info.yaml System description** says "A CSharp MVC API" (copy-paste leftover; it's a Go app).

## Important Files

| File | Role |
|---|---|
| `main.go` | Entry point |
| `cmd/root.go` | Cobra root command + `Execute()` |
| `cmd/server.go` | Router + all handlers + `requestLogger` middleware |
| `cmd/version.go` | `version` command |
| `pkg/version/version.go` | `VERSION`/`COMMIT`/`DATE` + `GetVersion()`/`GetPrettyVersion()` |
| `Makefile` | Build/test/release/Docker targets (see Development Commands) |
| `go.mod` / `go.sum` | Module `github.com/dgeorgievski/hello`, Go 1.24.0 |
| `docker/Dockerfile` | Multi-stage: builder (pinned `golang:1.24` digest) → `scratch`, `CMD ["/hello","server"]`, port 8080 |
| `.goreleaser.yaml` | Releases (`hello` binary, tar.gz/zip, `{{.Version}}`/`{{.ShortCommit}}`/`{{.Date}}` ldflags) |
| `.github/workflows/go-build.yaml` | CI: SHA-pinned actions + reusable `go-build.yml` (main/releases, PR) |
| `.github/workflows/go-deploy.yml` | Dispatch-only deploy (local-k8s / aws-eks) |
| `deploy/manifests/*.yaml` | Raw k8s: deployment (tag 0.0.1), serviceaccount, service |
| `deploy/helm/hello/` | Chart: deployment/hpa/ingress/service/serviceaccount/test pod; appVersion 1.16.0 |
| `catalog-info.yaml` | Backstage component + system entries |
| `scripts/agit/*.nu` | Nushell git-flow/commit helpers |
| `scripts/githooks/prepare-commit-msg` | Conventional Commits check (exit 1 aborts commit) |

## Runtime/Tooling Preferences

- **Go 1.24+** (module `go 1.24.0`) — not Node/Bun. Builds with `go build` (CGO_ENABLED=0 for container).
- **Package manager**: Go modules (`go.mod`/`go.sum`, `go mod tidy`); `make` is convenience, not required.
- **Runtime**: plain `go run` / built binary; container is `scratch` (no shell, static binary).
- **Tooling**: `make` (optional), `docker`, `kind` (local k8s), `helm`, `kubectl`, `nu` (Nushell, for `scripts/agit` + git hooks — hook shebang is `#!/opt/homebrew/bin/nu`), `git`, GitHub Actions for CI.
- **Constraints**: container runtime has no shell/CA certs/USER; use `curl`/HTTP inside k8s probes only. Do not commit `cover.out`/`coverage.out` (gitignored). `bin/` and `.worktrees/` are gitignored.

## Testing & QA

- **Framework**: stdlib `testing` only — no testify/ginkgo. HTTP tests use `httptest.NewRecorder()` with the in-process gin engine (`cmd/server_test.go`).
- **Run**: `make test` (`go test ./... -coverprofile cover.out`) or `go test ./...`. No `-race`, no `-count=1`, no coverage gate.
- **Coverage**: `cmd` package ≈ 74% (fresh run; `ping`/`hey`/`Execute`/`main`/`pkg/version` all 0.0%). `cover.out` / `coverage.out` are stale, untracked, and disagree with each other — regenerate with `make test` before trusting them.
- **Conventions**: table-driven subtests (`t.Run`), status-code + JSON-shape assertions, `newRouter()` for all handler coverage; no test opens a port or touches the filesystem; tests are order-independent. `TestRequestLoggerMiddleware` only asserts status through middleware — it does not verify log output.
- **Gaps to respect**: no tests for `/ping`, `/hey`, `Execute`, `main`, `version` command, or `pkg/version`; `/version` tests assert status only. If you add handlers, add table-driven tests in the same style (`httptest` + `newRouter()`).
