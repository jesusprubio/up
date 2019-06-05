<h1 align="center">online</h1>

<h4 align="center">
  📶 Library to check your Internet connectivity
</h4>

<div align="center">
  <img alt="Logo" src="https://media.giphy.com/media/pYyFAHLW0zJL2/giphy.gif" width="40%">
</div>

<p align="center">
  <a href="https://travis-ci.org/jesusprubio/online">
    <img alt="Build Status" src="https://travis-ci.org/jesusprubio/online.svg?branch=master">
  </a>
  <a href="https://crates.io/crates/online">
    <img alt="Latest version" src="https://img.shields.io/crates/v/online.svg">
  </a>
  <img alt="Stability stable" src="https://img.shields.io/badge/stability-stable-green.svg">
</p>
<p align="center">
  <sub>🤙 Ping me on <a href="https://twitter.com/jesusprubio"><code>Twitter</code></a> if you like this project</sub>
</p>

It uses TCP to try to connect to Chrome and Firefox (fallback) captive portal detection servers.

## Use

📝 Please visit [the full documentation](https://docs.rs/online) if you want to learn the details.

<!-- cargo-sync-readme start -->

```rust
use std::time::Duration;

use online::*;

assert_eq!(online(None), Ok(true));

// with timeout
let timeout = Duration::new(6, 0);
assert_eq!(online(Some(timeout)), Ok(true));
```

<!-- cargo-sync-readme end -->

## Contributing

😎 If you want to help please take a look to [this file](.github/CONTRIBUTING.md).
