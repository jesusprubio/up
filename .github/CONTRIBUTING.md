# Contribution guidelines

Thanks! Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.

## Developer environment

Some development tools are needed to be ready.

```sh
git clone https://github.com/jesusprubio/online
cd online
cargo install cargo-cmd
cargo cmd deps
```

## Tests

We use [Clippy](https://github.com/rust-lang/rust-clippy) and [rustfmt](https://github.com/rust-lang/rustfmt) linters. Please run to be sure your code fits with them and the tests keep passing:

```sh
cargo cmd test
```

## Publish

We use [cargo-release](https://github.com/sunng87/cargo-release) to make the process funnier.

```sh
cargo install cargo-release
cargo release
# cargo release minor
# cargo release major
```
