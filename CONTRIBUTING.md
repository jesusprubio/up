# Developer guide

Thanks!

## Environment

Some development tools are needed to be ready.

```sh
git clone https://github.com/jesusprubio/online
cd online
cargo install cargo-make
cargo make dep
```

## Tests

We use different linters and formatters. Please run to be sure your code fits with them and the tests keep passing:

```sh
cargo make ci
```

## Publish

We use [cargo-release](https://github.com/sunng87/cargo-release) to make the process funnier.

```sh
cargo install cargo-release
cargo release
# cargo release minor
# cargo release major
```
