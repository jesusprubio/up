use async_std::{
    io::{timeout as tout, Error},
    net::TcpStream,
};
use std::time::Duration;

// Captive portals:
// - https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/captivePortal
// - http://clients3.google.com/generate_204
// - http://detectportal.firefox.com/success.txt.
const ADDRS: &str = "clients3.google.com:80";
const ADDRS_BACK: &str = "detectportal.firefox.com:80";

async fn probe(addrs: &str, timeout: Option<Duration>) -> Result<TcpStream, Error> {
    match timeout {
        Some(d) => tout(d, TcpStream::connect(addrs)).await,
        _ => TcpStream::connect(addrs).await,
    }
}

/// Check if the internet connection is up. When the check succeeds, the returned future is resolved to `true`.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub async fn online(timeout: Option<u64>) -> Result<bool, Error> {
    //! ```rust
    //! use online::*;
    //!
    //! #[async_std::main]
    //! async fn main() {
    //!     assert_eq!(online(None).await.unwrap(), true);
    //!
    //!     // with timeout
    //!     assert_eq!(online(Some(6)).await.unwrap(), true);
    //! }
    //! ```
    let secs = match timeout {
        Some(n) => Some(Duration::from_secs(n)),
        _ => None,
    };

    match probe(ADDRS, secs).await {
        Ok(_) => Ok(true),
        Err(e) => match probe(ADDRS_BACK, secs).await {
            Ok(_) => Ok(true),
            // If both failed we return the first one.
            Err(_) => Err(e),
        },
    }
}
