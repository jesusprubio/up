<h1 align="center">up</h1>
<div align="center">
  <img alt="Logo" src="https://media.giphy.com/media/pYyFAHLW0zJL2/giphy.gif" width="40%">
</div>

Troubleshoot problems with your Internet connection based om different
[protocols](pkg/protocol.go) and [public servers](pkg/servers.go).

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] ![License](https://img.shields.io/github/license/jesusprubio/up)

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

### Library

Check [the examples](examples) to see how to use this project in your own code.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.

[doc-img]: https://pkg.go.dev/badge/github.com/jesusprubio/up
[doc]: https://pkg.go.dev/github.com/jesusprubio/up
[ci-img]: https://github.com/jesusprubio/up/workflows/CI/badge.svg
[ci]: https://github.com/jesusprubio/up/workflows/go.yml
