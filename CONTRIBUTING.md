# Developer Guide

## Dependencies

- [Task](https://taskfile.dev/installation)
- Linters:

```sh
task dep
```

## Develop

```sh
task
```

## Test

```sh
task vet # linters
task test
task fmt # formatters
```

## Build

```sh
task build
```

## Release

We use [GoReleaser](https://goreleaser.com/) and GitHub workflows to automate
the binary publishing process. Setup files:

- [GoReleaser](./.goreleaser.yml)
- [GitHub Workflows](./.github/workflows/release.yml)
