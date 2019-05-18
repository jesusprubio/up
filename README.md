# online

[![Build Status](https://travis-ci.org/jesusprubio/online.svg?branch=master)](https://travis-ci.org/jesusprubio/online)
[![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)](https://github.com/emersion/stability-badges#stable)
[![Latest version](https://img.shields.io/crates/v/log.svg)](https://crates.io/crates/online)
[![Documentation](https://docs.rs/online/badge.svg)](https://docs.rs/online)

📶 Library to check your Internet connectivity.

<!-- markdownlint-disable MD033 -->
<div align="center">
	<p><img src="https://media.giphy.com/media/pYyFAHLW0zJL2/giphy.gif" alt="gif icon"></p>
	<p><sub>🤙 Ping me on <a href="https://twitter.com/jesusprubio"><code>Twitter</code></a> if you like this project</sub></p>
</div>
<!-- markdownlint-enable MD033 -->

## Use

📝 Please visit [the full documentation](https://docs.rs/online) if you want to learn the details.

<!-- cargo-sync-readme start -->

```sh
extern crate online;
use online::*;

assert_eq!(online(None), Ok(true));
assert_eq!(online(Some(6)), Ok(true));
```

<!-- cargo-sync-readme end -->

## Contributing

😎 If you want to help please take a look to [this file](.github/CONTRIBUTING.md).
