<h1 align="center">online</h1>

<h4 align="center">
  📶 Library to check your Internet connectivity
</h4>

<div align="center">
  <img alt="Logo" src="https://media.giphy.com/media/pYyFAHLW0zJL2/giphy.gif" width="40%">
</div>

<p align="center">
  <a href="https://github.com/jesusprubio/online/actions">
    <img alt="Workflow status" src="https://github.com/jesusprubio/online/workflows/CI/badge.svg">
  </a>
  <a href="https://crates.io/crates/online">
    <img alt="Latest version" src="https://img.shields.io/crates/v/online.svg">
  </a>
</p>

_Features_

- Both asynchronous and blocking implementations.
- IPv4 and IPv6 support.

_How it works_

- Tries to connect to Chrome captive portal (using its domain name).
- If fails, tries the Firefox one.
- If both fail, the second error is returned to help with diagnostics.

## Install

The library is available on [crates.io](https://crates.io/crates/online). In example,
through [cargo-edit](https://github.com/killercup/cargo-edit):

```sh
cargo add online
```

### Async

```toml
online = { version = "4.0.0",  default-features = false, features = ["tokio"] }
```

## Use

📝 Please visit the [examples](examples) and [documentation](https://docs.rs/online)
to check the details.

<!-- cargo-sync-readme start -->

```rust
use online::check;

println!("Online? {}", check(None).is_ok());
println!("Online (timeout)? {}", check(Some(5)).is_ok());
```

<!-- cargo-sync-readme end -->

```sh
cargo run --example sync
cargo run --features="tokio-runtime" --example tokio
```
