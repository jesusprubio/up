---
name: update-docs
description: Keep the project documentation in sync with the code. Use after changing flags, protocols, output formats, exit codes, public servers, Make targets or the Go version, or when asked to review or update the docs.
---

# Update documentation

The code is the source of truth. Update the docs to match it, never the other
way around.

## Where each fact lives

Each fact is documented in one place only. Link to it instead of repeating it.

| Fact | Source of truth (code) | Documented in |
| --- | --- | --- |
| Flags and their defaults | `internal/options.go` | `README.md` (Use), `appDesc` in `main.go` |
| Protocols | `internal/protocol.go` | `README.md` (linked), `-p` usage |
| Output formats and fields | `internal/report.go`, `internal/summary.go` | `appDesc` in `main.go` (OUTPUT), `README.md` (Output) |
| Exit codes | `main.go` | `appDesc` in `main.go` (EXIT STATUS), `README.md` (Output) |
| Tips for scripts and agents | `main.go`, `internal/options.go` | `README.md` (Scripts and AI agents) |
| Public servers | `internal/servers.go` | `README.md` (linked) |
| Make targets, tooling | `Makefile` | `CONTRIBUTING.md` |
| Go version | `go.mod` | `.github/workflows/*.yml` (`go-version`) |
| Code layout, conventions, agent tips | the code | `CLAUDE.md` |

`CLAUDE.md` imports `README.md` and `CONTRIBUTING.md`: do not copy their
content into it.

## Steps

1. Find what changed: `git diff` (and `git diff --staged`) or the commits of
   the branch against `main`.
2. For each changed fact, check every place listed in the table above.
3. Verify the examples really work, running them with a count so they end:

   ```sh
   go run . -h
   go run . -nstd -c 1
   go run . -nstd -p http -c 1 -j
   go run . -nstd -c 1 -sum -j; echo "exit=$?"
   ```

   Network errors are fine in sandboxed environments; check the flags are
   accepted and the output has the documented shape.
4. Check that relative links in the Markdown files point to existing files.
5. Keep the style: short sentences, lines up to 80 chars, `sh` code blocks.
6. Summarize what was updated and anything that could not be verified.
