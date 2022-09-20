//! ```rust
//! use online::check;
//!
//! println!("Online? {}", check(None).is_ok());
//! println!("Online (timeout)? {}", check(Some(5)).is_ok());
//! ```

#[cfg(feature = "sync-runtime")]
mod sync;
#[cfg(feature = "sync-runtime")]
pub use sync::check;

#[cfg(feature = "tokio-runtime")]
pub mod tokio;
