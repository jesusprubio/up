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
up -p http -c 3
up -p http -tg example.com
cat testdata/stdin-urls.txt | go run . -p http
```

[doc-img]: https://pkg.go.dev/badge/github.com/jesusprubio/up
[doc]: https://pkg.go.dev/github.com/jesusprubio/up
[ci-img]: https://github.com/jesusprubio/up/workflows/CI/badge.svg
[ci]: https://github.com/jesusprubio/up/workflows/go.yml
