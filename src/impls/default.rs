#[cfg(all(not(feature = "tokio-runtime"), feature = "async-std-runtime"))]
use async_std::{io::timeout as tout, net::TcpStream};
use std::{io::Error, time::Duration};
#[cfg(all(not(feature = "async-std-runtime"), feature = "tokio-runtime"))]
use tokio::{net::TcpStream, time::timeout as tout};

#[path = "../utils.rs"]
mod utils;

// Custom version of `TcpStream::connect_timeout` not present in `async-std`.
// https://docs.rs/async-std/1.9.0/async_std/net/struct.TcpStream.html#method.connect
async fn connect_timeout(addrs: &str, dur: Duration) -> Result<(), Error> {
    if let Err(e) = tout(dur, TcpStream::connect(addrs)).await {
        #[cfg(all(not(feature = "async-std-runtime"), feature = "tokio-runtime"))]
        return Err(Error::new(std::io::ErrorKind::TimedOut, "future timed out"));
        #[cfg(all(not(feature = "tokio-runtime"), feature = "async-std-runtime"))]
        Err(e)
    } else {
        Ok(())
    }
}

/// Asynchronous implementation.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub async fn check(timeout: Option<u64>) -> Result<(), Error> {
    // Avoiding `io:timeout` in this case to allow the OS decide for better diagnostics.
    if let Some(t) = timeout {
        let dur = utils::parse_timeout(t)?;

        // First try, ignoring error (if any).
        return if connect_timeout(utils::ADDRS[0], dur).await.is_ok() {
            Ok(())
        } else {
            // Fallback.
            connect_timeout(utils::ADDRS[1], dur).await
        };
    }

    // No timeout.
    if TcpStream::connect(utils::ADDRS[0]).await.is_ok() {
        Ok(())
    } else if let Err(e) = TcpStream::connect(utils::ADDRS[1]).await {
        Err(e)
    } else {
        Ok(())
    }
}
