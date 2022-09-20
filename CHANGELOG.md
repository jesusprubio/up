# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

> **Types of changes**:
>
> - 🎉 **Added**: for new features.
> - ✏️ **Changed**: for changes in existing functionality.
> - ⚠️ **Deprecated**: for soon-to-be removed features.
> - ❌ **Removed**: for now removed features.
> - 🐛 **Fixed**: for any bug fixes.
> - 👾 **Security**: in case of vulnerabilities.

## [Unreleased]

## [4.0.0] - 2022-09-21

### 🎉 Added

- [Tokio runtime](https://tokio.rs/) support (thanks to @abdulrahman1s).

### ✏️ Changed

- Synchronous version is the default now.

### ❌ Removed

- `async-std` runtime support. Details [here](https://github.com/async-rs/async-std/issues/992#issuecomment-1035223559)
- Examples folder.

## [3.0.2] - 2022-09-04

### ✏️ Changed

- Dependencies update.
- Rust 2021 migration.

## [3.0.1] - 2021-06-26

### ✏️ Changed

- Minor refactoring.

## [3.0.0] - 2021-06-26

### 🎉 Added

- Synchronous/blocking API.
- Examples.

### ✏️ Changed

- Major refactor moving to more rusty code.
- Dependencies update.

### ❌ Removed

- Minor documentation cleanup.

## [2.0.0] - 2021-01-12

### ✏️ Changed

- Simpler API.

## [1.0.0] - 2021-01-10

### 🎉 Added

- Concurrency (`async-std`).
- IPv6.

### 🐛 Fixed

- Problem with `sync-readme`.

### ❌ Removed

- Use of `simple-error`.

## [0.2.2] - 2019-10-24

### ✏️ Changed

- Dependencies update.

## [0.2.2] - 2019-10-24

### ✏️ Changed

- Dependencies update.

## [0.2.1] - 2019-06-05

### ✏️ Changed

- Using Chrome and Firefox captive portal detection servers instead.
- Minor improvements in the README.
- Dependencies update.

## [0.2.0] - 2019-05-21

### ✏️ Changed

- Parameter `timeout` is now `Duration` instead of `u64`.
- Minor improvements in the manifest.

### 🐛 Fixed

- Minor documentation changes.

## [0.1.2] - 2019-05-19

### 🐛 Fixed

- Emojis in crates.io site.
- Link to the documentation in crates.io site.
- Clippy setup update, failing in Travis.

## [0.1.1] - 2019-05-19

### 🐛 Fixed

- Now the README is shown in crates.io site.
- Other minor Cargo setup fields, ie: keywords, categories, etc.
- Travis setup.
- Developer guide dependencies install.

## [0.1.0] - 2019-05-18

### 🎉 Added

- First release.

[unreleased]: https://github.com/jesusprubio/online/compare/v4.0.0...HEAD
[4.0.0]: https://github.com/jesusprubio/online/compare/v3.0.2...v4.0.0
[3.0.2]: https://github.com/jesusprubio/online/compare/v3.0.1...v3.0.2
[3.0.1]: https://github.com/jesusprubio/online/compare/v3.0.0...v3.0.1
[3.0.0]: https://github.com/jesusprubio/online/compare/v2.0.0...v3.0.0
[2.0.0]: https://github.com/jesusprubio/online/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/jesusprubio/online/compare/v0.2.2...v1.0.0
[0.2.2]: https://github.com/jesusprubio/online/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/jesusprubio/online/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/jesusprubio/online/compare/v0.1.2...v0.2.0
[0.1.2]: https://github.com/jesusprubio/online/compare/v0.1.0...v0.1.2
[0.1.1]: https://github.com/jesusprubio/online/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/jesusprubio/online/compare/v.0.0.1...v0.1.0
[0.0.1]: https://github.com/jesusprubio/online/compare/f855db0341fd9e60f30c507ea5ac92d139b5b7b3...v0.0.1
