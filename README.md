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

TCP and [async-std](https://github.com/async-rs/async-std) based function which tries to connect to Chrome and Firefox (fallback) captive portal detection servers.

## Install

With [cargo-edit](https://github.com/killercup/cargo-edit) installed run:

```sh
cargo add online
```

## Use

📝 Please visit [the full documentation](https://docs.rs/online) if you want to learn the details.

<!-- cargo-sync-readme start -->

```rust
use online::*;

#[async_std::main]
async fn main() {
    assert_eq!(online(None).await.unwrap(), true);

    // with timeout
    assert_eq!(online(Some(6)).await.unwrap(), true);
}
```

<!-- cargo-sync-readme end -->

## Contributing

😎 If you want to help please take a look to [this file](.github/CONTRIBUTING.md).
