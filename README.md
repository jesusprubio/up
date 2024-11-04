# up

Troubleshoot problems with your Internet connection based on different
[protocols](pkg/protocol.go) and well-known [public servers](pkg/servers.go).

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] ![License](https://img.shields.io/github/license/jesusprubio/up)

<div align="center">
  <img alt="Logo" src="https://github.com/jesusprubio/up/assets/2753855/a9c6bdb5-ab53-4969-8b36-97896c09a090" width="70%">
</div>

## Install

```sh
go install -v github.com/jesusprubio/up@latest
```

### Dependencies

- [Go](https://go.dev/doc/install) stable version.

## Use

The default behavior is to verify all the [supported protocols](pkg/protocol.go)
against a randomly selected [public server](pkg/servers.go) for each one.

```sh
up
```

This will display help for the tool.

```sh
up -h
```

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full text.

[doc-img]: https://pkg.go.dev/badge/github.com/jesusprubio/up
[doc]: https://pkg.go.dev/github.com/jesusprubio/up
[ci-img]: https://github.com/jesusprubio/up/workflows/CI/badge.svg
[ci]: https://github.com/jesusprubio/up/workflows/go.yml
