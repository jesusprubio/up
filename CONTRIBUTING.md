# Developer Guide

## Dependencies

- [GNU Make](https://www.gnu.org/software/make/)
- Linters:

```sh
make dep
```

## Develop

```sh
make
```

## Test

```sh
make vet # linters
make test
make fmt # formatters
```

## Build

```sh
make build
```

## Release

We use [GoReleaser](https://goreleaser.com/) and GitHub workflows to automate
the binary publishing process. Setup files:

- [GoReleaser](./.goreleaser.yml)
- [GitHub Workflows](./.github/workflows/release.yml)
