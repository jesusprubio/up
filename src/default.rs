use async_std::{
    io::{timeout as tout, Error},
    net::TcpStream,
};
use std::time::Duration;

#[path = "./utils.rs"]
mod utils;

// Custom version of `TcpStream::connect_timeout` not present in `async-std`.
// https://docs.rs/async-std/1.9.0/async_std/net/struct.TcpStream.html#method.connect
async fn connect_timeout(addrs: &str, dur: Duration) -> Result<(), Error> {
    if let Err(e) = tout(dur, TcpStream::connect(addrs)).await {
        return Err(e);
    }

    Ok(())
}

/// Asynchronous implementation.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub async fn check(timeout: Option<u64>) -> Result<(), Error> {
    // Avoiding `io:timeout` in this case to allow the OS decide for better diagnostics.
    if let Some(t) = timeout {
        let dur = utils::parse_timeout(t)?;

        // First try, ignoring error (if any).
        if connect_timeout(utils::ADDRS, dur).await.is_ok() {
            return Ok(());
        }

        // Fallback.
        return connect_timeout(utils::ADDRS_BACK, dur).await;
    }

    // Second try.
    if TcpStream::connect(utils::ADDRS).await.is_ok() {
        return Ok(());
    }

    if let Err(e) = TcpStream::connect(utils::ADDRS_BACK).await {
        return Err(e);
    }

    Ok(())
}
