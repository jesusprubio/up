@README.md
@CONTRIBUTING.md

## Layout

- `main.go`: CLI wiring (flags → probes → report printing, signals, stdin).
- `internal/options.go`: command line flags.
- `internal/protocol.go`: supported protocols (`HTTP`, `TCP`, `DNS`), each
  implements the `Protocol` interface.
- `internal/probe.go`: runs a protocol `Count` times and sends a `Report` per
  attempt to a channel.
- `internal/report.go`: output formats (human, JSON lines, grepable).
- `internal/summary.go`: per protocol aggregation (`-sum`) and the "up"
  verdict used for the exit status.
- `internal/servers.go`: public servers used when no target is given.
- `internal/stdin.go`: reads targets from standard input.

## Conventions

- Go version is pinned in `go.mod` and `.github/workflows/ci.yml`; keep both in
  sync.
- Before finishing a change run `make vet` and `make test` (CI runs both).
- Lines up to 80 chars (`make fmt` uses golines).
- Wrap errors with context: `fmt.Errorf("doing x: %w", err)`.
- Tests live next to the code (`*_test.go`) and use table-driven subtests.
- Some tests hit real network services (DNS, HTTP); they can fail in sandboxed
  environments without direct Internet access. CI on GitHub is the reference.

## Running the CLI

See "Scripts and AI agents" in the README. When verifying a change locally,
targets on `127.0.0.1` (e.g. `python3 -m http.server`) avoid depending on the
network.

## Documentation

When a change affects user facing behavior (flags, protocols, output formats,
exit codes, servers) update the docs in the same change. Use the
`update-docs` skill.
