# up

Troubleshoot problems with your Internet connection based on different
[protocols](internal/protocol.go) and well-known [public servers](internal/servers.go).

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] ![License](https://img.shields.io/github/license/jesusprubio/up)

<div align="center">
  <img alt="Logo" src="https://github.com/jesusprubio/up/assets/2753855/a9c6bdb5-ab53-4969-8b36-97896c09a090" width="70%">
</div>

## Install

### Binary Release

You can manually download a binary release for Linux, OSX, Windows or FreeBSD
from the [releases](https://github.com/jesusprubio/up/releases) page.

### Go

Please notice `latest` will install the dev version.

```sh
go install -ldflags="-s -w" -v github.com/jesusprubio/up@latest
```

## Use

The default behavior is to verify all the [supported protocols](internal/protocol.go)
against a randomly selected [public server](internal/servers.go) for each one.

```sh
up
up -p http
up -p http,dns -c 3
up -p http -tg https://example.com
up -p dns example.com example.org
cat testdata/stdin-urls.txt | up -p http
```

Every flag has a long version too (`-c` and `--count`), run `up -h` to see all
of them. Targets (`-tg`, arguments after the flags or standard input) require a
single protocol.

### Output

```sh
up -c 3 -j       # JSON lines, one object per request
up -c 3 -g       # grepable, tab separated
up -c 3 -sum     # only a summary per protocol at the end
up -c 3 -sum -j  # summary as a single JSON object
```

Exit status: `0` at least one response was heard, `2` no responses were
received, `1` any other error.

### Scripts and AI agents

- Always set a count (`-c`), by default `up` runs forever.
- Use `-nstd` when not piping targets, to avoid waiting for the standard input.
- `up -nstd -c 3 -sum -j` gives a short JSON verdict, or check only the exit
  status.

[doc-img]: https://pkg.go.dev/badge/github.com/jesusprubio/up
[doc]: https://pkg.go.dev/github.com/jesusprubio/up
[ci-img]: https://github.com/jesusprubio/up/workflows/CI/badge.svg
[ci]: https://github.com/jesusprubio/up/workflows/go.yml
