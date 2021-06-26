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


*Features*

- Both asynchronous and blocking implementations.
- IPv4 and IPv6 support.

*How it works*

- Tries to connect to Chrome captive portal (using its domain name).
- If fails, tries the Firefox one.
- If both fail, the second error is returned to help with diagnostics.

## Install

The library is available on [crates.io](https://crates.io/crates/online). Simply add the next line to your project's `Cargo.toml`.

```toml
online = "3.0.1"
```

### Synchronous

The [`async-std`](https://crates.io/crates/async-std) runtime is supported by default. But you can explicitly choose the blocking alternative.

```toml
online = { version = "3.0.1",  default-features = false, features = ["sync"] }
```

## Use

📝 Please visit the [examples](examples) and the [full documentation](https://docs.rs/online) if you want to learn the details.

<!-- cargo-sync-readme start -->

```rust
use online::check;

#[async_std::main]
async fn main() {
    println!("Online? {}", check(None).await.is_ok());
    println!("Online (timeout)? {}", check(Some(5)).await.is_ok());
    println!("Online (`Result`)? {:?}", check(None).await.unwrap());
}
```

<!-- cargo-sync-readme end -->

- [Synchronous example](examples/sync.rs)
